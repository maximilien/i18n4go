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

	combinedMap    map[string]common.I18nStringInfo
	OutputDirname  string
	SourceLanguage string
	Directory      string
}

func NewMergeStrings(options cmds.Options) MergeStrings {
	return MergeStrings{options: options,
		combinedMap:    make(map[string]common.I18nStringInfo),
		OutputDirname:  options.OutputDirFlag,
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
	files, _ := getFilesAndDir(ms.Directory)
	fileList := matchFileToSourceLanguage(files, ms.SourceLanguage)

	for _, file := range fileList {
		StringInfos, _ := common.LoadI18nStringInfos(file)
		ms.combineStringInfo(StringInfos)
	}

	common.SaveI18nStringInfos(ms, ms.mapValues2Array(), filepath.Join(ms.Directory, ms.SourceLanguage+".all.json"))

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

func (ms *MergeStrings) combineStringInfo(stringInfoList []common.I18nStringInfo) {
	for _, stringInfo := range stringInfoList {
		if _, ok := ms.combinedMap[stringInfo.ID]; !ok {
			ms.combinedMap[stringInfo.ID] = stringInfo
		}
	}
}

func (ms *MergeStrings) mapValues2Array() (stringInfoList []common.I18nStringInfo) {
	for _, stringInfo := range ms.combinedMap {
		stringInfoList = append(stringInfoList, stringInfo)
	}
	return
}
