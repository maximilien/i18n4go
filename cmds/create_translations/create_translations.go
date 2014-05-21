package create_translations

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"io/ioutil"

	"net/http"
	"net/url"
	"path/filepath"

	"github.com/maximilien/i18n4cf/cmds"
	"github.com/maximilien/i18n4cf/common"
)

type CreateTranslations struct {
	options cmds.Options

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

func NewCreateTranslations(options cmds.Options) CreateTranslations {
	languages := common.ParseLanguages(options.LanguagesFlag)

	return CreateTranslations{options: options,
		Filename:       options.FilenameFlag,
		OutputDirname:  options.OutputDirFlag,
		SourceLanguage: options.SourceLanguageFlag,
		Languages:      languages,
		TotalStrings:   0,
		TotalFiles:     0}
}

func (ct *CreateTranslations) Options() cmds.Options {
	return ct.options
}

func (ct *CreateTranslations) Println(a ...interface{}) (int, error) {
	if ct.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ct *CreateTranslations) Printf(msg string, a ...interface{}) (int, error) {
	if ct.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (ct *CreateTranslations) CreateTranslationFiles(sourceFilename string) error {
	ct.Println("gi18n: creating translation files for:", sourceFilename)
	ct.Println()
	ct.Filename = sourceFilename

	for _, language := range ct.Languages {
		ct.Println("gi18n: creating translation file copy for language:", language)

		if ct.options.GoogleTranslateApiKeyFlag != "" {
			fileName, _, err := ct.checkFile(sourceFilename)
			if err != nil {
				return err
			}

			destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.options.SourceLanguageFlag, language, -1))
			i18nStringInfos, err := ct.loadI18nStringInfos(sourceFilename)
			if err != nil {
				ct.Println(err)
				return fmt.Errorf("gi18n: could not load i18n strings from file: %s", destFilename)
			}

			ct.Println("gi18n: attempting to use Google Translate to translate source strings in: ", language)
			modifiedI18nStringInfos := make([]common.I18nStringInfo, len(i18nStringInfos))
			for i, i18nStringInfo := range i18nStringInfos {
				translation, _, err := ct.googleTranslate(i18nStringInfo.Translation, language)
				if err != nil {
					ct.Println("gi18n: error invoking Google Translate for string:", i18nStringInfo.Translation)
				} else {
					modifiedI18nStringInfos[i] = common.I18nStringInfo{ID: i18nStringInfo.ID, Translation: translation}
				}
			}

			err = common.SaveI18nStringInfos(ct, modifiedI18nStringInfos, destFilename)
			if err != nil {
				ct.Println(err)
				return fmt.Errorf("gi18n: could not save Google Translate i18n strings to file: %s", destFilename)
			}

			if ct.options.PoFlag {
				poFilename := destFilename[:len(destFilename)-len(".json")] + ".po"
				err = common.SaveI18nStringsInPo(ct, modifiedI18nStringInfos, poFilename)
				if err != nil {
					ct.Println(err)
					return fmt.Errorf("gi18n: could not save PO file: %s", poFilename)
				}
			}

			ct.Println()
		} else {
			_, err := ct.createTranslationFile(sourceFilename, language)
			if err != nil {
				return fmt.Errorf("gi18n: could not create default translation file for language: %s", language)
			}
		}
	}

	ct.Println()

	return nil
}

func (ct *CreateTranslations) checkFile(fileName string) (string, string, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", "", err
	}

	if !fileInfo.Mode().IsRegular() {
		return "", "", fmt.Errorf("gi18n: non-regular source file %s (%s)", fileInfo.Name(), fileInfo.Mode().String())
	}

	return fileInfo.Name(), fileName[:len(fileName)-len(fileInfo.Name())-1], nil
}

func (ct *CreateTranslations) loadI18nStringInfos(fileName string) ([]common.I18nStringInfo, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		ct.Println("gi18n: could not find file:", fileName)
		return nil, err
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		ct.Println(err)
		return nil, err
	}

	var i18nStringInfos []common.I18nStringInfo
	err = json.Unmarshal(content, &i18nStringInfos)
	if err != nil {
		ct.Println(err)
		return nil, err
	}

	return i18nStringInfos, nil
}

func (ct *CreateTranslations) googleTranslate(translateString string, language string) (string, string, error) {
	escapedTranslateString := url.QueryEscape(translateString)
	googleTranslateUrl := "https://www.googleapis.com/language/translate/v2?key=" + ct.options.GoogleTranslateApiKeyFlag + "&target=" + language + "&q=" + escapedTranslateString

	response, err := http.Get(googleTranslateUrl)
	if err != nil {
		ct.Println("gi18n: ERROR invoking Google Translate: ", googleTranslateUrl)
		return "", "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		ct.Println("gi18n: ERROR parsing Google Translate response body")
		return "", "", err
	}

	var googleTranslateData GoogleTranslateData
	err = json.Unmarshal(body, &googleTranslateData)
	if err != nil {
		ct.Println("gi18n: ERROR parsing Google Translate response body")
		return "", "", err
	}

	return googleTranslateData.Data.Translations[0].TranslatedText, googleTranslateData.Data.Translations[0].DetectedSourceLanguage, err
}

func (ct *CreateTranslations) createTranslationFile(sourceFilename string, language string) (string, error) {
	fileName, _, err := ct.checkFile(sourceFilename)
	if err != nil {
		return "", err
	}

	destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.options.SourceLanguageFlag, language, -1))
	ct.Println("gi18n: creating translation file:", destFilename)

	return destFilename, common.CopyFileContents(sourceFilename, destFilename)
}
