package verify_strings

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/maximilien/i18n4cf/cmds"
	"github.com/maximilien/i18n4cf/common"
)

type verifyStrings struct {
	options cmds.Options

	InputFilename string
	OutputDirname string

	SourceLanguage    string
	LanguageFilenames []string
	Languages         []string
}

func NewVerifyStrings(options cmds.Options) verifyStrings {
	languageFilenames := common.ParseStringList(options.LanguageFilesFlag, ",")
	languages := common.ParseStringList(options.LanguagesFlag, ",")

	return verifyStrings{options: options,
		InputFilename:     options.FilenameFlag,
		OutputDirname:     options.OutputDirFlag,
		LanguageFilenames: languageFilenames,
		Languages:         languages,
		SourceLanguage:    options.SourceLanguageFlag,
	}
}

func (vs *verifyStrings) Options() cmds.Options {
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
		vs.Println("gi18n: Error checking input filename: ", vs.InputFilename)
		return err
	}

	targetFilenames := vs.determineTargetFilenames(fileName, filePath)
	vs.Println("targetFilenames:", targetFilenames)
	for _, targetFilename := range targetFilenames {
		err = vs.verify(vs.InputFilename, targetFilename)
		if err != nil {
			vs.Println("gi18n: Error verifying target filename: ", targetFilename)
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
		vs.Println("gi18n: Error loading the i18n strings from input filename:", inputFilename)
		return err
	}

	if len(inputI18nStringInfos) == 0 {
		return fmt.Errorf("gi18n: Error input file: %s is empty", inputFilename)
	}

	inputMap := common.CreateI18nStringInfoMap(inputI18nStringInfos)

	targetI18nStringInfos, err := common.LoadI18nStringInfos(targetFilename)
	if err != nil {
		vs.Println("gi18n: Error loading the i18n strings from target filename:", targetFilename)
		return err
	}

	var targetExtraStringInfos []common.I18nStringInfo
	for _, stringInfo := range targetI18nStringInfos {
		if _, ok := inputMap[stringInfo.ID]; ok {
			delete(inputMap, stringInfo.ID)
		} else {
			vs.Println("gi18n: WARNING target file has extra key with ID: ", stringInfo.ID)
			targetExtraStringInfos = append(targetExtraStringInfos, stringInfo)
		}
	}

	var verficationError error
	if len(targetExtraStringInfos) > 0 {
		vs.Println("gi18n: WARNING target file contains total of extra keys:", len(targetExtraStringInfos))

		diffFilename, err := vs.generateExtraKeysDiffFile(targetExtraStringInfos, targetFilename)
		if err != nil {
			vs.Println("gi18n: ERROR could not create the diff file:", err)
			return err
		}
		vs.Println("gi18n: generated diff file:", diffFilename)
		verficationError = fmt.Errorf("gi18n: target file has extra i18n strings with IDs: %s", strings.Join(keysForI18nStringInfos(targetExtraStringInfos), ","))
	}

	if len(inputMap) > 0 {
		vs.Println("gi18n: ERROR input file does not match target file:", targetFilename)

		diffFilename, err := vs.generateMissingKeysDiffFile(valuesForI18nStringInfoMap(inputMap), targetFilename)
		if err != nil {
			vs.Println("gi18n: ERROR could not create the diff file:", err)
			return err
		}
		vs.Println("gi18n: generated diff file:", diffFilename)
		verficationError = fmt.Errorf("gi18n: target file is missing i18n strings with IDs: %s", strings.Join(keysForI18nStringInfoMap(inputMap), ","))
	}

	return verficationError
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

	return diffFilename, common.SaveI18nStringInfos(vs, missingStringInfos, diffFilename)
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

	return diffFilename, common.SaveI18nStringInfos(vs, extraStringInfos, diffFilename)
}
