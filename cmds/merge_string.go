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
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/maximilien/i18n4go/common"
)

type mergeStrings struct {
	options common.Options

	I18nStringInfos []common.I18nStringInfo

	Recurse        bool
	SourceLanguage string
	Directory      string
}

func NewMergeStrings(options common.Options) *mergeStrings {
	return &mergeStrings{
		options:         options,
		I18nStringInfos: []common.I18nStringInfo{},
		Recurse:         options.RecurseFlag,
		SourceLanguage:  options.SourceLanguageFlag,
		Directory:       options.DirnameFlag,
	}
}

// NewMergeStringsCommand implements 'i18n merge-strings' command
func NewMergeStringsCommand(p *I18NParams, options common.Options) *cobra.Command {
	mergeStringsCmd := &cobra.Command{
		Use:   "merge-strings",
		Short: "Merge translation strings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewMergeStrings(options).Run()
		},
	}

	// TODO: setup options and params for Cobra command here using common.Options

	return mergeStringsCmd
}

func (ms *mergeStrings) Options() common.Options {
	return ms.options
}

func (ms *mergeStrings) Println(a ...interface{}) (int, error) {
	if ms.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ms *mergeStrings) Printf(msg string, a ...interface{}) (int, error) {
	if ms.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (ms *mergeStrings) Run() error {
	return ms.combineStringInfosPerDirectory(ms.Directory)
}

func (ms *mergeStrings) combineStringInfosPerDirectory(directory string) error {
	files, directories := getFilesAndDir(directory)
	fileList := ms.matchFileToSourceLanguage(files, ms.SourceLanguage)

	combinedMap := map[string]common.I18nStringInfo{}
	for _, file := range fileList {
		StringInfos, err := common.LoadI18nStringInfos(file)
		if err != nil {
			return err
		}

		combineStringInfo(StringInfos, combinedMap)
	}

	filePath := filepath.Join(directory, ms.SourceLanguage+".all.json")
	ms.I18nStringInfos = common.I18nStringInfoMapValues2Array(combinedMap)
	sort.Sort(ms)
	common.SaveI18nStringInfos(ms, ms.Options(), ms.I18nStringInfos, filePath)
	ms.Println("i18n4go: saving combined language file: " + filePath)

	if ms.Recurse {
		for _, directory = range directories {
			err := ms.combineStringInfosPerDirectory(directory)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getFilesAndDir(dir string) (files []string, dirs []string) {
	contents, _ := ioutil.ReadDir(dir)

	for _, fileInfo := range contents {
		if !fileInfo.IsDir() {
			files = append(files, filepath.Join(dir, fileInfo.Name()))
		} else {
			dirs = append(dirs, filepath.Join(dir, fileInfo.Name()))
		}
	}
	return
}

func (ms mergeStrings) matchFileToSourceLanguage(files []string, lang string) (list []string) {
	languageMatcher := "go." + lang + ".json"
	for _, file := range files {
		if strings.Contains(file, languageMatcher) {
			list = append(list, file)
			ms.Println("i18n4go: scanning file: " + file)
		}
	}
	return
}

func combineStringInfo(stringInfoList []common.I18nStringInfo, combinedMap map[string]common.I18nStringInfo) {
	for _, stringInfo := range stringInfoList {
		if _, ok := combinedMap[stringInfo.ID]; !ok {
			combinedMap[stringInfo.ID] = stringInfo
		}
	}
}

// sort.Interface methods

func (ms *mergeStrings) Len() int {
	return len(ms.I18nStringInfos)
}

func (ms *mergeStrings) Less(i, j int) bool {
	return ms.I18nStringInfos[i].ID < ms.I18nStringInfos[j].ID
}

func (ms *mergeStrings) Swap(i, j int) {
	tmpI18nStringInfo := ms.I18nStringInfos[i]
	ms.I18nStringInfos[i] = ms.I18nStringInfos[j]
	ms.I18nStringInfos[j] = tmpI18nStringInfo
}
