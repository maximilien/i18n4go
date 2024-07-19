// Copyright © 2015-2023 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmds

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/maximilien/i18n4go/i18n4go/common"
	"github.com/maximilien/i18n4go/i18n4go/i18n"
	"github.com/spf13/cobra"
)

type verifyStrings struct {
	options common.Options

	InputFilename string
	OutputDirname string

	SourceLanguage    string
	LanguageFilenames []string
	Languages         []string
}

func NewVerifyStrings(options *common.Options) *verifyStrings {
	languageFilenames := common.ParseStringList(options.LanguageFilesFlag, ",")
	languages := common.ParseStringList(options.LanguagesFlag, ",")

	return &verifyStrings{options: *options,
		InputFilename:     options.FilenameFlag,
		OutputDirname:     options.OutputDirFlag,
		LanguageFilenames: languageFilenames,
		Languages:         languages,
		SourceLanguage:    options.SourceLanguageFlag,
	}
}

// NewVerifyStringsCommand implements 'i18n4go verify-strings' command
func NewVerifyStringsCommand(options *common.Options) *cobra.Command {
	verifyStringsCmd := &cobra.Command{
		Use:   "verify-strings",
		Short: i18n.T("Verify strings in translations"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewVerifyStrings(options).Run()
		},
	}

	verifyStringsCmd.Flags().StringVarP(&options.SourceLanguageFlag, "source-language", "s", "en", i18n.T("the source language of the file, typically also part of the file name, e.g., \"en_US\""))
	verifyStringsCmd.Flags().StringVar(&options.LanguagesFlag, "languages", "", i18n.T("a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\""))
	verifyStringsCmd.Flags().StringVar(&options.LanguageFilesFlag, "language-files", "", i18n.T(`a comma separated list of target files for different languages to compare,  e.g., \"en, en_US, fr_FR, es\"	                                                                  if not specified then the languages flag is used to find target files in same directory as source`))
	verifyStringsCmd.Flags().StringVarP(&options.OutputDirFlag, "output", "o", "", i18n.T("the output directory where the missing translation keys will be placed"))
	verifyStringsCmd.Flags().StringVarP(&options.FilenameFlag, "file", "f", "", i18n.T("the source translation file"))
	return verifyStringsCmd
}

func (vs *verifyStrings) Options() common.Options {
	return vs.options
}

func (vs *verifyStrings) Println(a ...interface{}) (int, error) {
	if vs.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (vs *verifyStrings) Printf(msg string, a ...interface{}) (int, error) {
	if vs.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (vs *verifyStrings) Run() error {
	fileName, filePath, err := common.CheckFile(vs.InputFilename)
	if err != nil {
		vs.Println(i18n.T("i18n4go: Error checking input filename: "), vs.InputFilename)
		return err
	}

	targetFilenames := vs.determineTargetFilenames(fileName, filePath)
	vs.Println(i18n.T("targetFilenames:"), targetFilenames)
	for _, targetFilename := range targetFilenames {
		err = vs.verify(vs.InputFilename, targetFilename)
		if err != nil {
			vs.Println(i18n.T("i18n4go: Error verifying target filename: "), targetFilename)
		}
	}

	return err
}

func (vs *verifyStrings) determineTargetFilenames(inputFilename string, inputFilePath string) []string {
	if len(vs.LanguageFilenames) != 0 {
		return vs.LanguageFilenames
	}

	var targetFilename string
	targetFilenames := make([]string, len(vs.Languages))
	for i, lang := range vs.Languages {
		targetFilename = strings.Replace(inputFilename, vs.SourceLanguage, lang, -1)
		targetFilenames[i] = filepath.Join(inputFilePath, targetFilename)
	}

	return targetFilenames
}

func (vs *verifyStrings) verify(inputFilename string, targetFilename string) error {
	common.CheckFile(targetFilename)

	inputI18nStringInfos, err := common.LoadI18nStringInfos(inputFilename)
	if err != nil {
		vs.Println(i18n.T("i18n4go: Error loading the i18n strings from input filename:"), inputFilename)
		return err
	}

	if len(inputI18nStringInfos) == 0 {
		return fmt.Errorf(i18n.T("i18n4go: Error input file: {{.Arg0}} is empty", map[string]interface{}{"Arg0": inputFilename}))
	}

	inputMap, err := common.CreateI18nStringInfoMap(inputI18nStringInfos)
	if err != nil {
		return fmt.Errorf(i18n.T("File has duplicated key: {{.Arg0}}\n{{.Arg1}}", map[string]interface{}{"Arg0": inputFilename, "Arg1": err}))
	}

	targetI18nStringInfos, err := common.LoadI18nStringInfos(targetFilename)
	if err != nil {
		vs.Println(i18n.T("i18n4go: Error loading the i18n strings from target filename:"), targetFilename)
		return err
	}

	var targetExtraStringInfos, targetInvalidStringInfos []common.I18nStringInfo
	for _, stringInfo := range targetI18nStringInfos {
		if _, ok := inputMap[stringInfo.ID]; ok {
			if common.IsTemplatedString(stringInfo.ID) && vs.isTemplatedStringTranslationInvalid(stringInfo) {
				vs.Println(i18n.T("i18n4go: WARNING target file has invalid templated translations with key ID: "), stringInfo.ID)
				targetInvalidStringInfos = append(targetInvalidStringInfos, stringInfo)
			}
			delete(inputMap, stringInfo.ID)
		} else {
			vs.Println(i18n.T("i18n4go: WARNING target file has extra key with ID: "), stringInfo.ID)
			targetExtraStringInfos = append(targetExtraStringInfos, stringInfo)
		}
	}

	var verficationError error
	if len(targetExtraStringInfos) > 0 {
		vs.Println(i18n.T("i18n4go: WARNING target file contains total of extra keys:"), len(targetExtraStringInfos))

		diffFilename, err := vs.generateExtraKeysDiffFile(targetExtraStringInfos, targetFilename)
		if err != nil {
			vs.Println(i18n.T("i18n4go: ERROR could not create the diff file:"), err)
			return err
		}
		vs.Println(i18n.T("i18n4go: generated diff file:"), diffFilename)
		verficationError = fmt.Errorf(i18n.T("i18n4go: target file has extra i18n strings with IDs: {{.Arg0}}", map[string]interface{}{"Arg0": strings.Join(keysForI18nStringInfos(targetExtraStringInfos), ",")}))
	}

	if len(targetInvalidStringInfos) > 0 {
		vs.Println(i18n.T("i18n4go: WARNING target file contains total of invalid translations:"), len(targetInvalidStringInfos))

		diffFilename, err := vs.generateInvalidTranslationDiffFile(targetInvalidStringInfos, targetFilename)
		if err != nil {
			vs.Println(i18n.T("i18n4go: ERROR could not create the diff file:"), err)
			return err
		}
		vs.Println(i18n.T("i18n4go: generated diff file:"), diffFilename)
		verficationError = fmt.Errorf(i18n.T("i18n4go: target file has invalid i18n strings with IDs: {{.Arg0}}", map[string]interface{}{"Arg0": strings.Join(keysForI18nStringInfos(targetInvalidStringInfos), ",")}))
	}

	if len(inputMap) > 0 {
		vs.Println(i18n.T("i18n4go: ERROR input file does not match target file:"), targetFilename)

		diffFilename, err := vs.generateMissingKeysDiffFile(valuesForI18nStringInfoMap(inputMap), targetFilename)
		if err != nil {
			vs.Println(i18n.T("i18n4go: ERROR could not create the diff file:"), err)
			return err
		}
		vs.Println(i18n.T("i18n4go: generated diff file:"), diffFilename)
		verficationError = fmt.Errorf(i18n.T("i18n4go: target file is missing i18n strings with IDs: {{.Arg0}}", map[string]interface{}{"Arg0": strings.Join(keysForI18nStringInfoMap(inputMap), ",")}))
	}

	return verficationError
}

func (vs *verifyStrings) isTemplatedStringTranslationInvalid(stringInfo common.I18nStringInfo) bool {
	if !common.IsTemplatedString(stringInfo.ID) || !common.IsTemplatedString(stringInfo.Translation) {
		return false
	}

	translationArgs := common.GetTemplatedStringArgs(stringInfo.Translation)
	argsMap := make(map[string]string)
	for _, translationArg := range translationArgs {
		argsMap[translationArg] = translationArg
	}

	var missingArgs []string
	idArgs := common.GetTemplatedStringArgs(stringInfo.ID)
	for _, idArg := range idArgs {
		if _, ok := argsMap[idArg]; !ok {
			missingArgs = append(missingArgs, idArg)
		}
	}

	if len(missingArgs) > 0 {
		vs.Println(i18n.T("i18n4go: templated string is invalid, missing args in translation:"), strings.Join(missingArgs, ","))
		return true
	}

	return false
}

func keysForI18nStringInfos(in18nStringInfos []common.I18nStringInfo) []string {
	var keys []string
	for _, stringInfo := range in18nStringInfos {
		keys = append(keys, stringInfo.ID)
	}
	return keys
}

func keysForI18nStringInfoMap(inputMap map[string]common.I18nStringInfo) []string {
	var keys []string
	for k, _ := range inputMap {
		keys = append(keys, k)
	}
	return keys
}

func valuesForI18nStringInfoMap(inputMap map[string]common.I18nStringInfo) []common.I18nStringInfo {
	var values []common.I18nStringInfo
	for _, v := range inputMap {
		values = append(values, v)
	}
	return values
}

func (vs *verifyStrings) generateMissingKeysDiffFile(missingStringInfos []common.I18nStringInfo, fileName string) (string, error) {
	name, pathName, err := common.CheckFile(fileName)
	if err != nil {
		return "", err
	}

	diffFilename := name + ".missing.diff.json"
	if vs.OutputDirname != "" {
		common.CreateOutputDirsIfNeeded(vs.OutputDirname)
		diffFilename = filepath.Join(vs.OutputDirname, diffFilename)
	} else {
		diffFilename = filepath.Join(pathName, diffFilename)
	}

	return diffFilename, common.SaveI18nStringInfos(vs, vs.Options(), missingStringInfos, diffFilename)
}

func (vs *verifyStrings) generateExtraKeysDiffFile(extraStringInfos []common.I18nStringInfo, fileName string) (string, error) {
	name, pathName, err := common.CheckFile(fileName)
	if err != nil {
		return "", err
	}

	diffFilename := name + ".extra.diff.json"
	if vs.OutputDirname != "" {
		common.CreateOutputDirsIfNeeded(vs.OutputDirname)
		diffFilename = filepath.Join(vs.OutputDirname, diffFilename)
	} else {
		diffFilename = filepath.Join(pathName, diffFilename)
	}

	return diffFilename, common.SaveI18nStringInfos(vs, vs.Options(), extraStringInfos, diffFilename)
}

func (vs *verifyStrings) generateInvalidTranslationDiffFile(invalidStringInfos []common.I18nStringInfo, fileName string) (string, error) {
	name, pathName, err := common.CheckFile(fileName)
	if err != nil {
		return "", err
	}

	diffFilename := name + ".invalid.diff.json"
	if vs.OutputDirname != "" {
		common.CreateOutputDirsIfNeeded(vs.OutputDirname)
		diffFilename = filepath.Join(vs.OutputDirname, diffFilename)
	} else {
		diffFilename = filepath.Join(pathName, diffFilename)
	}

	return diffFilename, common.SaveI18nStringInfos(vs, vs.Options(), invalidStringInfos, diffFilename)
}
