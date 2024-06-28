// Copyright © 2015-2023 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmds

import (
	"fmt"
	"strings"

	"encoding/json"
	"io/ioutil"

	"net/http"
	"net/url"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/maximilien/i18n4go/i18n4go/common"
	"github.com/maximilien/i18n4go/i18n4go/i18n"
)

type createTranslations struct {
	options common.Options

	Filename       string
	OutputDirname  string
	SourceLanguage string

	Languages []string

	ExtractedStrings map[string]common.StringInfo

	TotalStrings int
	TotalFiles   int
}

type GoogleTranslateData struct {
	Data GoogleTranslateTranslations `json:"data"`
}

type GoogleTranslateTranslations struct {
	Translations []GoogleTranslateTranslation `json:"translations"`
}

type GoogleTranslateTranslation struct {
	TranslatedText         string `json:"translatedText"`
	DetectedSourceLanguage string `json:"detectedSourceLanguage"`
}

func NewCreateTranslations(options *common.Options) *createTranslations {
	languages := common.ParseStringList(options.LanguagesFlag, ",")

	return &createTranslations{options: *options,
		Filename:       options.FilenameFlag,
		OutputDirname:  options.OutputDirFlag,
		SourceLanguage: options.SourceLanguageFlag,
		Languages:      languages,
		TotalStrings:   0,
		TotalFiles:     0}
}

// NewCreateTranslationsCommand implements 'i18n4go create-translations' command
func NewCreateTranslationsCommand(options *common.Options) *cobra.Command {
	createTranslationsCmd := &cobra.Command{
		Use:   "create-translations",
		Short: i18n.T("Creates the transation files"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewCreateTranslations(options).Run()
		},
	}

	// TODO: --google-translate-api-key is too long of an optional flag
	// might want to shorten it or add an alias for usability
	createTranslationsCmd.Flags().StringVar(&options.GoogleTranslateApiKeyFlag, "google-translate-api-key", "", i18n.T("[optional] your public Google Translate API key which is used to generate translations (charge is applicable)"))
	createTranslationsCmd.Flags().StringVarP(&options.SourceLanguageFlag, "source-language", "s", "en", i18n.T("the source language of the file, typically also part of the file name, e.g., \"en_US\""))
	createTranslationsCmd.Flags().StringVarP(&options.FilenameFlag, "file", "f", "", i18n.T("the source translation file"))
	createTranslationsCmd.Flags().StringVarP(&options.LanguagesFlag, "languages", "l", "", i18n.T("a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\""))
	createTranslationsCmd.Flags().StringVarP(&options.OutputDirFlag, "output", "o", "", i18n.T("the output directory where the newly created translation files will be placed"))

	return createTranslationsCmd

}

func (ct *createTranslations) Options() common.Options {
	return ct.options
}

func (ct *createTranslations) Println(a ...interface{}) (int, error) {
	if ct.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ct *createTranslations) Printf(msg string, a ...interface{}) (int, error) {
	if ct.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (ct *createTranslations) Run() error {
	ct.Println(i18n.T("i18n4go: creating translation files for:"), ct.Filename)
	ct.Println()

	for _, language := range ct.Languages {
		ct.Println(i18n.T("i18n4go: creating translation file copy for language:"), language)

		if ct.options.GoogleTranslateApiKeyFlag != "" {
			destFilename, err := ct.createTranslationFileWithGoogleTranslate(language)
			if err != nil {
				return fmt.Errorf(i18n.T("i18n4go: could not create translation file for language: {{.Arg0}} with Google Translate", map[string]interface{}{"Arg0": language}))
			}
			ct.Println(i18n.T("i18n4go: created translation file with Google Translate:"), destFilename)
		} else {
			destFilename, err := ct.createTranslationFile(ct.Filename, language)
			if err != nil {
				return fmt.Errorf(i18n.T("i18n4go: could not create default translation file for language: {{.Arg0}}\nerr:{{.Arg1}}", map[string]interface{}{"Arg0": language, "Arg1": err.Error()}))
			}
			ct.Println(i18n.T("i18n4go: created default translation file:"), destFilename)
		}
	}

	ct.Println()

	return nil
}

func (ct *createTranslations) createTranslationFileWithGoogleTranslate(language string) (string, error) {
	fileName, _, err := common.CheckFile(ct.Filename)
	if err != nil {
		return "", err
	}

	err = common.CreateOutputDirsIfNeeded(ct.OutputDirname)
	if err != nil {
		ct.Println(err)
		return "", fmt.Errorf(i18n.T("i18n4go: could not create output directory: {{.Arg0}}", map[string]interface{}{"Arg0": ct.OutputDirname}))
	}

	destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.options.SourceLanguageFlag, language, -1))

	i18nStringInfos, err := common.LoadI18nStringInfos(ct.Filename)
	if err != nil {
		ct.Println(err)
		return "", fmt.Errorf(i18n.T("i18n4go: could not load i18n strings from file: {{.Arg0}}", map[string]interface{}{"Arg0": ct.Filename}))
	}

	if len(i18nStringInfos) == 0 {
		return "", fmt.Errorf(i18n.T("i18n4go: input file: {{.Arg0}} is empty", map[string]interface{}{"Arg0": ct.Filename}))
	}

	ct.Println(i18n.T("i18n4go: attempting to use Google Translate to translate source strings in: "), language)
	modifiedI18nStringInfos := make([]common.I18nStringInfo, len(i18nStringInfos))
	for i, i18nStringInfo := range i18nStringInfos {
		translation, _, err := ct.googleTranslate(i18nStringInfo.Translation, language)
		if err != nil {
			ct.Println(i18n.T("i18n4go: error invoking Google Translate for string:"), i18nStringInfo.Translation)
		} else {
			modifiedI18nStringInfos[i] = common.I18nStringInfo{ID: i18nStringInfo.ID, Translation: translation}
		}
	}

	err = common.SaveI18nStringInfos(ct, ct.Options(), modifiedI18nStringInfos, destFilename)
	if err != nil {
		ct.Println(err)
		return "", fmt.Errorf(i18n.T("i18n4go: could not save Google Translate i18n strings to file: {{.Arg0}}", map[string]interface{}{"Arg0": destFilename}))
	}

	if ct.options.PoFlag {
		poFilename := destFilename[:len(destFilename)-len(".json")] + ".po"
		err = common.SaveI18nStringsInPo(ct, ct.Options(), modifiedI18nStringInfos, poFilename)
		if err != nil {
			ct.Println(err)
			return "", fmt.Errorf(i18n.T("i18n4go: could not save PO file: {{.Arg0}}", map[string]interface{}{"Arg0": poFilename}))
		}
	}

	ct.Println()

	return destFilename, nil
}

func (ct *createTranslations) createTranslationFile(sourceFilename string, language string) (string, error) {
	fileName, _, err := common.CheckFile(sourceFilename)
	if err != nil {
		return "", err
	}

	i18nStringInfos, err := common.LoadI18nStringInfos(sourceFilename)
	if err != nil {
		ct.Println(err)
		return "", fmt.Errorf(i18n.T("i18n4go: could not load i18n strings from file: {{.Arg0}}", map[string]interface{}{"Arg0": sourceFilename}))
	}

	if len(i18nStringInfos) == 0 {
		return "", fmt.Errorf(i18n.T("i18n4go: input file: {{.Arg0}} is empty", map[string]interface{}{"Arg0": sourceFilename}))
	}

	destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.options.SourceLanguageFlag, language, -1))
	ct.Println(i18n.T("i18n4go: creating translation file:"), destFilename)

	return destFilename, common.CopyFileContents(sourceFilename, destFilename)
}

func (ct *createTranslations) googleTranslate(translateString string, language string) (string, string, error) {
	escapedTranslateString := url.QueryEscape(translateString)
	googleTranslateUrl := "https://www.googleapis.com/language/translate/v2?key=" + ct.options.GoogleTranslateApiKeyFlag + "&target=" + language + "&q=" + escapedTranslateString

	// REMOVEME: Do not commit
	fmt.Printf("\ngoogleTranslateUrl: %s\n", googleTranslateUrl)
	response, err := http.Get(googleTranslateUrl)
	if err != nil {
		ct.Println(i18n.T("i18n4go: ERROR invoking Google Translate: "), googleTranslateUrl)
		return "", "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ct.Println(i18n.T("i18n4go: ERROR parsing Google Translate response body"))
		return "", "", err
	}

	var googleTranslateData GoogleTranslateData
	err = json.Unmarshal(body, &googleTranslateData)
	if err != nil {
		ct.Println(i18n.T("i18n4go: ERROR parsing Google Translate response body"))
		return "", "", err
	}

	return googleTranslateData.Data.Translations[0].TranslatedText, googleTranslateData.Data.Translations[0].DetectedSourceLanguage, err
}
