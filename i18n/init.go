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

package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pivotal-cf-experimental/jibber_jabber"

	go_i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	DEFAULT_LOCALE   = "en_US"
	DEFAULT_LANGUAGE = "en"
)

type TranslateFunc func(translateID string, args ...interface{}) string

var SUPPORTED_LOCALES = map[string]string{
	"de": "de_DE",
	"en": "en_US",
	"es": "es_ES",
	"fr": "fr_FR",
	"it": "it_IT",
	"ja": "ja_JA",
	"ko": "ko_KO",
	"pt": "pt_BR",
	"ru": "ru_RU",
	"zh": "zh_CN",
}
var (
	RESOUCES_PATH = filepath.Join("cf", "i18n", "resources")
	bundle        *go_i18n.Bundle
)

func GetResourcesPath() string {
	return RESOUCES_PATH
}

func init() {
	if bundle == nil {
		bundle = go_i18n.NewBundle(language.AmericanEnglish)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	}
}

func Init(packageName string, i18nDirname string) TranslateFunc {
	userLocale, err := initWithUserLocale(packageName, i18nDirname)
	if err != nil {
		userLocale = mustLoadDefaultLocale(packageName, i18nDirname)
	}

	return Tfunc(userLocale, DEFAULT_LOCALE)
}

func initWithUserLocale(packageName, i18nDirname string) (string, error) {
	userLocale, err := jibber_jabber.DetectIETF()
	if err != nil {
		userLocale = DEFAULT_LOCALE
	}

	language, err := jibber_jabber.DetectLanguage()
	if err != nil {
		language = DEFAULT_LANGUAGE
	}

	userLocale = strings.Replace(userLocale, "-", "_", 1)
	err = loadFromAsset(packageName, i18nDirname, userLocale, language)
	if err != nil {
		locale := SUPPORTED_LOCALES[language]
		if locale == "" {
			userLocale = DEFAULT_LOCALE
		} else {
			userLocale = locale
		}
		err = loadFromAsset(packageName, i18nDirname, userLocale, language)
	}

	return userLocale, err
}

func mustLoadDefaultLocale(packageName, i18nDirname string) string {
	userLocale := DEFAULT_LOCALE

	err := loadFromAsset(packageName, i18nDirname, DEFAULT_LOCALE, DEFAULT_LANGUAGE)
	if err != nil {
		panic("Could not load en_US language files. God save the queen. " + err.Error())
	}

	return userLocale
}

func loadFromAsset(packageName, assetPath, locale, language string) error {
	assetName := locale + ".all.json"
	assetKey := filepath.Join(assetPath, language, packageName, assetName)

	byteArray, err := resources.Asset(assetKey)
	if err != nil {
		return err
	}

	if len(byteArray) == 0 {
		return errors.New(fmt.Sprintf("Could not load i18n asset: %v", assetKey))
	}

	tmpDir, err := ioutil.TempDir("", "i18n4go_res")
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	fileName, err := saveLanguageFileToDisk(tmpDir, assetName, byteArray)
	if err != nil {
		return err
	}

	bundle.MustLoadMessageFile(fileName)

	os.RemoveAll(fileName)

	return nil
}

func saveLanguageFileToDisk(tmpDir, assetName string, byteArray []byte) (fileName string, err error) {
	fileName = filepath.Join(tmpDir, assetName)
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(byteArray)
	if err != nil {
		return
	}

	return
}

func isNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	}
	return false
}

func SupportedLocaleLanguageTags() []language.Tag {
	tags := []language.Tag{language.English}
	for _, locale := range SUPPORTED_LOCALES {
		tag, _ := language.Parse(locale)
		tags = append(tags, tag)
	}

	return tags
}

// translate is a wrapper function that is based on the translate method for v1.3.0
// To allow compatibility v2.0+ and older a wrapper method was created
// @see https://github.com/nicksnyder/go-i18n/blob/v1.3.0/i18n/bundle/bundle.go#L227-L257
func translate(loc *go_i18n.Localizer) TranslateFunc {
	return func(messageId string, args ...interface{}) string {
		var (
			count interface{}
			data  interface{}
		)

		if argc := len(args); argc > 0 {
			if isNumber(args[0]) {
				count = args[0]
				if argc > 1 {
					data = args[1]
				}
			} else {
				data = args[0]
			}
		}

		return loc.MustLocalize(&go_i18n.LocalizeConfig{
			MessageID:    messageId,
			TemplateData: data,
			PluralCount:  count,
		})
	}
}

// Tfunc will return a method of TranslateFunc type to be used to tranlation messages
func Tfunc(sources ...string) TranslateFunc {
	localizer := go_i18n.NewLocalizer(bundle, DEFAULT_LOCALE)
	tfunc := translate(localizer)
	for _, s := range sources {
		if s == "" {
			continue
		}

		if s == DEFAULT_LOCALE {
			return tfunc
		}
		localizer := go_i18n.NewLocalizer(bundle, s)
		t := translate(localizer)
		return func(translationID string, args ...interface{}) string {
			return t(translationID, args...)
		}
	}

	return tfunc
}
