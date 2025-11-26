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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/maximilien/i18n4go/i18n4go/common"
	"github.com/maximilien/i18n4go/i18n4go/i18n"
	"github.com/spf13/cobra"
)

type fixup struct {
	options common.Options

	I18nStringInfos []common.I18nStringInfo
	English         []common.I18nStringInfo
	Source          map[string]int
	Locales         map[string]map[string]string
	IgnoreRegexp    *regexp.Regexp
}

func NewFixup(options *common.Options) *fixup {
	return &fixup{
		options:         *options,
		I18nStringInfos: []common.I18nStringInfo{},
		IgnoreRegexp:    common.GetIgnoreRegexp(options.IgnoreRegexpFlag),
	}
}

// NewFixupCommand implements 'i18n fixup' command
func NewFixupCommand(options *common.Options) *cobra.Command {
	fixupCmd := &cobra.Command{
		Use:   "fixup",
		Long:  i18n.T("Add, update, or remove translation keys from source files and resources files"),
		Short: i18n.T("Fixup the translation files"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewFixup(options).Run()
		},
	}

	fixupCmd.Flags().StringVar(&options.IgnoreRegexpFlag, "ignore-regexp", ".*test.*", i18n.T("recursively extract strings from all files in the same directory as filename or dirName"))
	fixupCmd.Flags().StringVarP(&options.QualifierFlag, "qualifier", "q", "i18n", i18n.T("[optional] the qualifier string that is used when importing the package for i18n4go to use the i18n.T(...) function"))
	fixupCmd.Flags().StringVar(&options.SourceDirFlag, "source", ".", i18n.T("[optional] the directory where the source go files are located, defaults to current directory"))
	fixupCmd.Flags().StringVar(&options.ResourceDirFlag, "resource", ".", i18n.T("[optional] the directory where the translation files are located, defaults to current directory"))

	return fixupCmd
}

func (fix *fixup) Options() common.Options {
	return fix.options
}

func (fix *fixup) Println(a ...interface{}) (int, error) {
	if fix.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (fix *fixup) Printf(msg string, a ...interface{}) (int, error) {
	if fix.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (fix *fixup) Run() error {
	//FIND PROBLEMS HERE AND RETURN AN ERROR
	var (
		translationsAdded   bool
		translationsRemoved bool
		translationsUpdated bool
	)
	source, err := fix.findSourceStrings(fix.options.SourceDirFlag)
	fix.Source = source

	if err != nil {
		fmt.Println(i18n.T("Couldn't find any source strings: {{.Arg0}}", map[string]interface{}{
			"Arg0": err.Error(),
		}))
		return err
	}

	locales := findTranslationFiles(fix.options.ResourceDirFlag, fix.IgnoreRegexp, fix.options.VerboseFlag)
	englishFiles, ok := locales["en_US"]
	if !ok {
		fmt.Println(i18n.T("Unable to find english translation files"))
		return errors.New(i18n.T("Unable to find english translation files"))
	}

	englishFile := englishFiles[0]
	if englishFile == "" {
		fmt.Println(i18n.T("Could not find an i18n file for locale: en_US"))
		return errors.New(i18n.T("Could not find an i18n file for locale: en_US"))
	}

	englishStringInfos, err := fix.findI18nStrings(englishFile)

	if err != nil {
		fmt.Println(i18n.T("Couldn't find the english strings: {{.Arg0}}", map[string]interface{}{
			"Arg0": err.Error(),
		}))
		return err
	}

	//Check english to all other files before source
	for locale, i18nFile := range locales {
		if locale != "en_US" {
			foreignStringInfos, _ := fix.findI18nStrings(i18nFile[0])
			foreignAdditionalTranslations := getAdditionalForeignTranslations(englishStringInfos, foreignStringInfos)

			foreignMissingTranslations := getMissingForeignTranslations(englishStringInfos, foreignStringInfos)

			if len(foreignMissingTranslations) > 0 {
				addTranslations(foreignStringInfos, i18nFile[0], foreignMissingTranslations)
				translationsAdded = true
			} else {
				translationsAdded = false
			}

			if len(foreignAdditionalTranslations) > 0 {
				removeTranslations(foreignStringInfos, i18nFile[0], foreignAdditionalTranslations)
				translationsRemoved = true
			} else {
				translationsRemoved = false
			}

			// update only if there were changes
			if translationsAdded || translationsRemoved {
				writeStringInfoMapToJSON(foreignStringInfos, i18nFile[0])
			}

			// reset flags
			translationsRemoved = false
			translationsAdded = false
		}
	}

	//rewrite everything now
	potentialAdditionalTranslations := getAdditionalTranslations(source, englishStringInfos)
	removedTranslations := getRemovedTranslations(source, englishStringInfos)

	additionalTranslations := []string{}
	updatedTranslations := make(map[string]string)

	if len(potentialAdditionalTranslations) > 0 && len(removedTranslations) > 0 {
		for _, newUpdatedTranslation := range potentialAdditionalTranslations {
			if len(removedTranslations) > 0 {
				var input string

				escape := false
				updated := false

				for !escape {
					fmt.Printf(i18n.T("Is the string \"%s\" a new or updated string? [new/upd]\n"), newUpdatedTranslation)

					_, err := fmt.Scanf("%s\n", &input)
					if err != nil {
						panic(err)
					}

					input = strings.ToLower(input)

					switch input {
					case "new":
						additionalTranslations = append(additionalTranslations, newUpdatedTranslation)
						escape = true
					case "upd":
						fmt.Println(i18n.T("Select the number for the previous translation:"))
						for index, value := range removedTranslations {
							fmt.Printf("\t%d. %s\n", (index + 1), value)
						}

						var updSelection int
						for !updated {
							_, err := fmt.Scanf("%d\n", &updSelection)

							if err == nil && updSelection > 0 && updSelection <= len(removedTranslations) {
								updSelection = updSelection - 1

								updatedTranslations[removedTranslations[updSelection]] = newUpdatedTranslation

								removedTranslations = removeFromSlice(removedTranslations, updSelection)

								updated = true
							} else {
								fmt.Println(i18n.T("Invalid response."))
							}
						}
						escape = true
					case "exit":
						fmt.Println(i18n.T("Canceling fixup"))
						os.Exit(0)
					default:
						fmt.Println(i18n.T("Invalid response."))
					}
				}
			} else {
				additionalTranslations = append(additionalTranslations, newUpdatedTranslation)
			}
		}
	} else {
		additionalTranslations = potentialAdditionalTranslations
	}

	for locale, i18nFiles := range locales {
		translatedStrings, err := fix.findI18nStrings(i18nFiles[0])
		if err != nil {
			fmt.Println(fmt.Sprintf(i18n.T("Couldn't get the strings from {{.Arg0}}: {{.Arg1}}"), locale, err.Error()))
			return err
		}

		if len(updatedTranslations) > 0 {
			updateTranslations(translatedStrings, i18nFiles[0], locale, updatedTranslations)
			translationsUpdated = true
		} else {
			translationsUpdated = false
		}

		if len(additionalTranslations) > 0 {
			addTranslations(translatedStrings, i18nFiles[0], additionalTranslations)
			translationsAdded = true
		} else {
			translationsAdded = false
		}

		if len(removedTranslations) > 0 {
			removeTranslations(translatedStrings, i18nFiles[0], removedTranslations)
			translationsRemoved = true
		} else {
			translationsRemoved = false
		}

		// update only if there is a change
		if translationsAdded || translationsRemoved || translationsUpdated {
			err = writeStringInfoMapToJSON(translatedStrings, i18nFiles[0])
		}

		// reset flags
		translationsAdded = false
		translationsUpdated = false
		translationsRemoved = false

	}

	if err == nil {
		fmt.Print(i18n.T("OK"))
	}

	return err
}

func (fix *fixup) findSourceStrings(dir string) (sourceStrings map[string]int, err error) {
	sourceStrings = make(map[string]int)
	files := getGoFiles(dir)

	for _, file := range files {
		fileStrings, err := common.InspectFile(file, fix.options)
		if err != nil {
			fmt.Println(i18n.T("Error when inspecting go file: "), file)
			return sourceStrings, err
		}

		for _, string := range fileStrings {
			sourceStrings[string]++
		}
	}

	return
}

func (fix *fixup) findI18nStrings(i18nFile string) (i18nStrings map[string]common.I18nStringInfo, err error) {
	i18nStrings = make(map[string]common.I18nStringInfo)

	stringInfos, err := common.LoadI18nStringInfos(i18nFile)

	if err != nil {
		return nil, err
	}

	return common.CreateI18nStringInfoMap(stringInfos)
}

func getAdditionalTranslations(sourceTranslations map[string]int, englishTranslations map[string]common.I18nStringInfo) []string {
	additionalTranslations := []string{}

	for id, _ := range sourceTranslations {
		if _, ok := englishTranslations[id]; !ok {
			additionalTranslations = append(additionalTranslations, id)
		}
	}
	return additionalTranslations
}

func getAdditionalForeignTranslations(englishTranslations, foreignTranslations map[string]common.I18nStringInfo) []string {
	additionalForeignTranslations := []string{}
	for key, _ := range foreignTranslations {
		if (englishTranslations[key] == common.I18nStringInfo{}) {
			additionalForeignTranslations = append(additionalForeignTranslations, key)
		}
	}
	return additionalForeignTranslations
}

func getRemovedTranslations(sourceTranslations map[string]int, englishTranslations map[string]common.I18nStringInfo) []string {
	removedTranslations := []string{}

	for id := range englishTranslations {
		if _, ok := sourceTranslations[id]; !ok {
			removedTranslations = append(removedTranslations, id)
		}
	}

	return removedTranslations
}

func getMissingForeignTranslations(englishTranslations, foreignTranslations map[string]common.I18nStringInfo) []string {
	missingForeignTranslations := []string{}
	for key, _ := range englishTranslations {
		if (foreignTranslations[key] == common.I18nStringInfo{}) {
			missingForeignTranslations = append(missingForeignTranslations, key)
		}
	}
	return missingForeignTranslations
}

func writeStringInfoMapToJSON(localeMap map[string]common.I18nStringInfo, localeFile string) error {
	localeArray := common.I18nStringInfoMapValues2Array(localeMap)

	sort.Sort(array(localeArray))

	encodedLocale, err := json.MarshalIndent(localeArray, "", "   ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(localeFile, encodedLocale, 0644)
	if err != nil {
		return err
	}
	return nil
}

func addTranslations(localeMap map[string]common.I18nStringInfo, localeFile string, addTranslations []string) {
	fmt.Printf(i18n.T("Adding these strings to the %s translation file:\n"), localeFile)

	for _, id := range addTranslations {
		localeMap[id] = common.I18nStringInfo{ID: id, Translation: id}
		fmt.Println("\t", id)
	}
}

func removeTranslations(localeMap map[string]common.I18nStringInfo, localeFile string, remTranslations []string) error {
	var err error
	fmt.Printf(i18n.T("Removing these strings from the %s translation file:\n"), localeFile)

	for _, id := range remTranslations {
		delete(localeMap, id)
		fmt.Println("\t", id)
	}

	return err
}

func updateTranslations(localMap map[string]common.I18nStringInfo, localeFile string, locale string, updTranslations map[string]string) {
	fmt.Printf(i18n.T("Updating the following strings from the %s translation file:\n"), localeFile)

	for key, value := range updTranslations {
		fmt.Println("\t", key)

		if locale == "en_US" {
			localMap[value] = common.I18nStringInfo{ID: value, Translation: value}
		} else {
			localMap[value] = common.I18nStringInfo{ID: value, Translation: localMap[key].Translation}
		}
		delete(localMap, key)
	}
}

func removeFromSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

//Interface for sort

type array []common.I18nStringInfo

func (stringInfos array) Len() int {
	return len(stringInfos)
}

func (stringInfos array) Less(i, j int) bool {
	return stringInfos[i].ID < stringInfos[j].ID
}

func (stringInfos array) Swap(i, j int) {
	tmpI18nStringInfo := stringInfos[i]
	stringInfos[i] = stringInfos[j]
	stringInfos[j] = tmpI18nStringInfo
}
