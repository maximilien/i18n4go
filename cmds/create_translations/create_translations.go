package create_translations

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"encoding/json"
	"io/ioutil"

	"net/http"
	"net/url"

	"path/filepath"

	common "github.com/maximilien/i18n4cf/common"
)

type CreateTranslations struct {
	Options common.Options

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

func NewCreateTranslations(options common.Options) CreateTranslations {
	languages := common.ParseLanguages(options.LanguagesFlag)

	return CreateTranslations{Options: options,
		Filename:       options.FilenameFlag,
		OutputDirname:  options.OutputDirFlag,
		SourceLanguage: options.SourceLanguageFlag,
		Languages:      languages,
		TotalStrings:   0,
		TotalFiles:     0}
}

func (ct *CreateTranslations) Println(a ...interface{}) (int, error) {
	if ct.Options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ct *CreateTranslations) Printf(msg string, a ...interface{}) (int, error) {
	if ct.Options.VerboseFlag {
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

		if ct.Options.GoogleTranslateApiKeyFlag != "" {
			fileName, _, err := ct.checkFile(sourceFilename)
			if err != nil {
				return err
			}

			destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.Options.SourceLanguageFlag, language, -1))
			i18nStringInfos, err := ct.loadI18nStringInfos(sourceFilename)
			if err != nil {
				ct.Println(err)
				return fmt.Errorf("gi18n: could not load i18n strings from file: ", destFilename)
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

			err = ct.saveI18nStringInfos(destFilename, modifiedI18nStringInfos)
			if err != nil {
				ct.Println(err)
				return fmt.Errorf("gi18n: could not save Google Translate i18n strings to file: ", destFilename)
			}

			if ct.Options.PoFlag {
				poFilename := destFilename[:len(destFilename)-len(".json")] + ".po"
				err = ct.saveI18nStringsInPo(poFilename, modifiedI18nStringInfos)
				if err != nil {
					ct.Println(err)
					return fmt.Errorf("gi18n: could not save PO file: ", poFilename)
				}
			}

			ct.Println()
		} else {
			_, err := ct.createTranslationFile(sourceFilename, language)
			if err != nil {
				return fmt.Errorf("gi18n: could not create default translation file for language: ", language)
			}
		}
	}

	ct.Println()

	return nil
}

func (ct *CreateTranslations) saveI18nStringsInPo(fileName string, i18nStrings []common.I18nStringInfo) error {
	ct.Println("gi18n: creating and saving i18n strings to .po file:", fileName)

	if !ct.Options.DryRunFlag && len(i18nStrings) != 0 {
		file, err := os.Create(fileName)
		defer file.Close()
		if err != nil {
			ct.Println(err)
			return err
		}

		for _, stringInfo := range i18nStrings {
			file.Write([]byte("msgid " + strconv.Quote(stringInfo.ID) + "\n"))
			file.Write([]byte("msgstr " + strconv.Quote(stringInfo.Translation) + "\n"))
			file.Write([]byte("\n"))
		}
	}
	return nil
}

func (ct *CreateTranslations) saveI18nStringInfos(fileName string, i18nStringInfos []common.I18nStringInfo) error {
	ct.Println("fileName:", fileName)
	jsonData, err := json.MarshalIndent(i18nStringInfos, "", "   ")
	if err != nil {
		ct.Println(err)
		return err
	}

	if !ct.Options.DryRunFlag && len(i18nStringInfos) != 0 {
		err := ioutil.WriteFile(fileName, jsonData, 0700)
		if err != nil {
			ct.Println(err)
			return err
		}
	}

	return nil
}

func (ct *CreateTranslations) googleTranslate(translateString string, language string) (string, string, error) {
	escapedTranslateString := url.QueryEscape(translateString)
	googleTranslateUrl := "https://www.googleapis.com/language/translate/v2?key=" + ct.Options.GoogleTranslateApiKeyFlag + "&target=" + language + "&q=" + escapedTranslateString

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

func (ct *CreateTranslations) createTranslationFile(sourceFilename string, language string) (string, error) {
	fileName, _, err := ct.checkFile(sourceFilename)
	if err != nil {
		return "", err
	}

	destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.Options.SourceLanguageFlag, language, -1))
	ct.Println("gi18n: creating translation file:", destFilename)

	return destFilename, common.CopyFileContents(sourceFilename, destFilename)
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
