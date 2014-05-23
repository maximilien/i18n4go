package merge_strings

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/maximilien/i18n4cf/cmds"
	"github.com/maximilien/i18n4cf/common"
)

type MergeStrings struct {
	options cmds.Options

	RecurseFlag    bool
	SourceLanguage string
	Directory      string
}

func NewMergeStrings(options cmds.Options) MergeStrings {
	return MergeStrings{
		options:        options,
		RecurseFlag:    options.RecurseFlag,
		SourceLanguage: options.SourceLanguageFlag,
		Directory:      options.DirnameFlag,
	}
}

func (ms *MergeStrings) Options() cmds.Options {
	return ms.options
}

func (ms *MergeStrings) Println(a ...interface{}) (int, error) {
	if ms.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ms *MergeStrings) Printf(msg string, a ...interface{}) (int, error) {
	if ms.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (ms *MergeStrings) Run() error {
	return ms.combineStringInfosPerDirectory(ms.Directory)
}

func (ms *MergeStrings) combineStringInfosPerDirectory(directory string) error {
	files, directories := getFilesAndDir(directory)
	fileList := matchFileToSourceLanguage(files, ms.SourceLanguage)

	combinedMap := map[string]common.I18nStringInfo{}
	for _, file := range fileList {
		StringInfos, err := common.LoadI18nStringInfos(file)
		if err != nil {
			return err
		}

		combineStringInfo(StringInfos, combinedMap)
	}

	common.SaveI18nStringInfos(ms, mapValues2Array(combinedMap), filepath.Join(directory, ms.SourceLanguage+".all.json"))

	if ms.RecurseFlag {
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

func matchFileToSourceLanguage(files []string, lang string) (list []string) {
	languageMatcher := "go." + lang + ".json"
	for _, file := range files {
		if strings.Contains(file, languageMatcher) {
			list = append(list, file)
			fmt.Println(file)
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

func mapValues2Array(combinedMap map[string]common.I18nStringInfo) (stringInfoList []common.I18nStringInfo) {
	for _, stringInfo := range combinedMap {
		stringInfoList = append(stringInfoList, stringInfo)
	}
	return
}
