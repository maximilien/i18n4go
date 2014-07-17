package cmds

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/maximilien/i18n4go/common"
)

type Checkup struct {
	options common.Options

	I18nStringInfos []common.I18nStringInfo
}

func NewCheckup(options common.Options) Checkup {
	return Checkup{
		options:         options,
		I18nStringInfos: []common.I18nStringInfo{},
	}
}

func (cu *Checkup) Options() common.Options {
	return cu.options
}

func (cu *Checkup) Println(a ...interface{}) (int, error) {
	if cu.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (cu *Checkup) Printf(msg string, a ...interface{}) (int, error) {
	if cu.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (cu *Checkup) Run() error {
	//FIND PROBLEMS HERE AND RETURN AN ERROR
	sourceStrings, err := cu.findSourceStrings()

	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't find any source strings: %s", err.Error()))
		return err
	}

	locales := findTranslationFiles(".")

	englishFiles := locales["en_US"]
	if englishFiles == nil {
		fmt.Println("Could not find an i18n file for locale: en_US")
		return errors.New("Could not find an i18n file for locale: en_US")
	}

	englishStrings, err := cu.findI18nStrings(englishFiles)

	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't find the english strings: %s", err.Error()))
		return err
	}

	err = cu.diffStrings("the code", "en_US", sourceStrings, englishStrings)

	for locale, i18nFiles := range locales {
		if locale == "en_US" {
			continue
		}

		translatedStrings, err := cu.findI18nStrings(i18nFiles)

		if err != nil {
			fmt.Println(fmt.Sprintf("Couldn't get the strings from %s: %s", locale, err.Error()))
			return err
		}

		err = cu.diffStrings("en_US", locale, englishStrings, translatedStrings)
	}

	if err == nil {
		fmt.Printf("OK")
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

func (cu *Checkup) inspectFile(file string) (translatedStrings []string, err error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		cu.Println(err)
		return
	}

	ast.Inspect(astFile, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.CallExpr:
			switch x.Fun.(type) {
			case *ast.Ident:
				funName := x.Fun.(*ast.Ident).Name

				if funName == "T" || funName == "t" {
					if stringArg, ok := x.Args[0].(*ast.BasicLit); ok {
						translatedString, err := strconv.Unquote(stringArg.Value)
						if err != nil {
							panic(err.Error())
						}
						translatedStrings = append(translatedStrings, translatedString)
					}
				}
			default:
				//Skip!
			}
		}
		return true
	})

	return
}

func (cu *Checkup) findSourceStrings() (sourceStrings map[string]string, err error) {
	sourceStrings = make(map[string]string)
	files := getGoFiles(".")

	for _, file := range files {
		fileStrings, err := cu.inspectFile(file)
		if err != nil {
			fmt.Println("Error when inspecting go file: ", file)
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

			if strings.HasSuffix(name, locale+".all.json") {
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

func findTranslationFiles(dir string) (locales map[string][]string) {
	locales = make(map[string][]string)
	contents, _ := ioutil.ReadDir(dir)

	for _, fileInfo := range contents {
		if !fileInfo.IsDir() {
			name := fileInfo.Name()

			if strings.HasSuffix(name, ".all.json") {
				locale := strings.Split(name, ".")[0]

				if locales[locale] == nil {
					locales[locale] = []string{}
				}

				locales[locale] = append(locales[locale], filepath.Join(dir, fileInfo.Name()))
			}
		} else {
			for locale, files := range findTranslationFiles(filepath.Join(dir, fileInfo.Name())) {
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
			fmt.Printf("\"%s\" exists in %s, but not in %s\n", key, sourceNameOne, sourceNameTwo)
			err = errors.New("Strings don't match")
		}
	}

	for key, _ := range stringsTwo {
		if stringsOne[key] == "" {
			fmt.Printf("\"%s\" exists in %s, but not in %s\n", key, sourceNameTwo, sourceNameOne)
			err = errors.New("Strings don't match")
		}
	}

	return
}
