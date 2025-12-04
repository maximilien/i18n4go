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
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/maximilien/i18n4go/i18n4go/common"
	"github.com/maximilien/i18n4go/i18n4go/i18n"
)

type Checkup struct {
	options common.Options

	I18nStringInfos []common.I18nStringInfo
	IgnoreRegexp    *regexp.Regexp
}

func NewCheckup(options *common.Options) *Checkup {
	return &Checkup{
		options:         *options,
		I18nStringInfos: []common.I18nStringInfo{},
		IgnoreRegexp:    common.GetIgnoreRegexp(options.IgnoreRegexpFlag),
	}
}

// NewCheckupCommand implements 'i18n4go checkup' command
func NewCheckupCommand(options *common.Options) *cobra.Command {
	checkupCmd := &cobra.Command{
		Use:   "checkup",
		Short: i18n.T("Checks the translated files"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewCheckup(options).Run()
		},
	}

	checkupCmd.Flags().StringVarP(&options.QualifierFlag, "qualifier", "q", "", i18n.T("[optional] the qualifier string that is used when using the i18n.T(...) function, default to nothing but could be set to `i18n` so that all calls would be: i18n.T(...)"))
	// TODO: Optional flags shouldn't have set defaults. We should look into removing the default
	checkupCmd.Flags().StringVar(&options.IgnoreRegexpFlag, "ignore-regexp", ".*test.*", i18n.T("recursively extract strings from all files in the same directory as filename or dirName"))
	return checkupCmd
}

func (cu *Checkup) Options() common.Options {
	return cu.options
}

func (cu *Checkup) Println(a ...any) (int, error) {
	if cu.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (cu *Checkup) Printf(msg string, a ...any) (int, error) {
	if cu.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (cu *Checkup) Run() error {
	//FIND PROBLEMS HERE AND RETURN AN ERROR
	sourceStrings, err := cu.findSourceStrings()

	if err != nil {
		cu.Println(i18n.T("Couldn't find any source strings: {{.Arg0}}", map[string]any{"Arg0": err.Error()}))
		return err
	}

	locales := findTranslationFiles(".", cu.IgnoreRegexp, false)

	englishFiles := locales["en_US"]
	if englishFiles == nil {
		cu.Println(i18n.T("Could not find an i18n file for locale: en_US"))
		return errors.New(i18n.T("Could not find an i18n file for locale: en_US"))
	}

	englishStrings, err := cu.findI18nStrings(englishFiles)

	if err != nil {
		cu.Println(i18n.T("Couldn't find the english strings: {{.Arg0}}", map[string]any{"Arg0": err.Error()}))
		return err
	}

	err = cu.diffStrings(i18n.T("the code"), "en_US", sourceStrings, englishStrings)

	for locale, i18nFiles := range locales {
		if locale == "en_US" {
			continue
		}

		translatedStrings, err := cu.findI18nStrings(i18nFiles)

		if err != nil {
			cu.Println(i18n.T("Couldn't get the strings from {{.Arg0}}: {{.Arg1}}", map[string]any{"Arg0": locale, "Arg1": err.Error()}))
			return err
		}

		err = cu.diffStrings("en_US", locale, englishStrings, translatedStrings)
	}

	if err == nil {
		cu.Printf(i18n.T("OK"))
	}

	return err
}

func getGoFiles(dir string) (files []string) {
	contents, _ := ioutil.ReadDir(dir)

	for _, fileInfo := range contents {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()

			if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
				files = append(files, filepath.Join(dir, fileInfo.Name()))
			}
		} else {
			moreFiles := getGoFiles(filepath.Join(dir, fileInfo.Name()))
			files = append(files, moreFiles...)
		}
	}
	return
}

func (cu *Checkup) findSourceStrings() (sourceStrings map[string]string, err error) {
	sourceStrings = make(map[string]string)
	files := getGoFiles(".")

	for _, file := range files {
		fileStrings, err := common.InspectFile(file, cu.options)
		if err != nil {
			cu.Println(i18n.T("Error when inspecting go file: "), file)
			return sourceStrings, err
		}

		for _, string := range fileStrings {
			sourceStrings[string] = string
		}
	}

	return
}

// Thought: Implement a function that searches directories recursively
// and finds all files that fit a certain pattern (e.g. regex).
//
// This would be fairly useful in finding translation files/source files.
func getI18nFile(locale, dir string) (filePath string) {
	contents, _ := ioutil.ReadDir(dir)

	for _, fileInfo := range contents {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()

			// assume the file path is a json file and the path contains the locale
			if strings.HasSuffix(name, ".json") && strings.Contains(name, fmt.Sprintf("{{.Arg0}}.", map[string]any{"Arg0": locale})) {
				filePath = filepath.Join(dir, fileInfo.Name())
				break
			}
		} else {
			filePath = getI18nFile(locale, filepath.Join(dir, fileInfo.Name()))

			if filePath != "" {
				break
			}
		}
	}

	return
}

func findTranslationFiles(dir string, ignoreRegexp *regexp.Regexp, verbose bool) (locales map[string][]string) {
	locales = make(map[string][]string)
	contents, _ := ioutil.ReadDir(dir)

	for _, fileInfo := range contents {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()

			if strings.HasSuffix(name, ".json") {
				parts := strings.Split(name, ".")
				var locale string

				for _, part := range parts {
					invalidLangRegexp, _ := regexp.Compile("excluded|json|all")
					if !invalidLangRegexp.MatchString(part) {
						locale = part
					}
				}

				// No locale found so skipping
				if locale == "" {
					continue
				}

				if locales[locale] == nil {
					locales[locale] = []string{}
				}

				locales[locale] = append(locales[locale], filepath.Join(dir, fileInfo.Name()))
			}
		} else {
			if ignoreRegexp != nil {
				if ignoreRegexp.MatchString(filepath.Join(dir, fileInfo.Name())) {
					continue
				}
			}
			for locale, files := range findTranslationFiles(filepath.Join(dir, fileInfo.Name()), ignoreRegexp, verbose) {
				if locales[locale] == nil {
					locales[locale] = []string{}
				}

				locales[locale] = append(locales[locale], files...)
			}
		}
	}

	return
}

func (cu *Checkup) findI18nStrings(i18nFiles []string) (i18nStrings map[string]string, err error) {
	i18nStrings = make(map[string]string)

	for _, i18nFile := range i18nFiles {
		stringInfos, err := common.LoadI18nStringInfos(i18nFile)

		if err != nil {
			return nil, err
		}

		for _, info := range stringInfos {
			i18nStrings[info.ID] = info.Translation
		}
	}

	return
}

func (cu *Checkup) diffStrings(sourceNameOne, sourceNameTwo string, stringsOne, stringsTwo map[string]string) (err error) {
	for key, _ := range stringsOne {
		if stringsTwo[key] == "" {
			cu.Printf(i18n.T("\"{{.Arg0}}\" exists in {{.Arg1}}, but not in {{.Arg2}}\n", map[string]any{"Arg0": key, "Arg1": sourceNameOne, "Arg2": sourceNameTwo}))
			err = errors.New(i18n.T("Strings don't match"))
		}
	}

	for key, _ := range stringsTwo {
		if stringsOne[key] == "" {
			cu.Printf(i18n.T("\"{{.Arg0}}\" exists in {{.Arg1}}, but not in {{.Arg2}}\n", map[string]any{"Arg0": key, "Arg1": sourceNameTwo, "Arg2": sourceNameOne}))
			err = errors.New(i18n.T("Strings don't match"))
		}
	}

	return
}
