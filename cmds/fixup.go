package cmds

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/maximilien/i18n4go/common"
)

type Fixup struct {
	options common.Options

	I18nStringInfos []common.I18nStringInfo
	English         []common.I18nStringInfo
	Source          map[string]int
	Locales         map[string]map[string]string
}

func NewFixup(options common.Options) Fixup {
	return Fixup{
		options:         options,
		I18nStringInfos: []common.I18nStringInfo{},
	}
}

func (fu *Fixup) Options() common.Options {
	return fu.options
}

func (fu *Fixup) Println(a ...interface{}) (int, error) {
	if fu.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (fu *Fixup) Printf(msg string, a ...interface{}) (int, error) {
	if fu.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (fu *Fixup) Run() error {
	//FIND PROBLEMS HERE AND RETURN AN ERROR
	source, err := fu.findSourceStrings()
	fu.Source = source

	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't find any source strings: %s", err.Error()))
		return err
	}

	locales := findTranslationFiles(".")

	englishFile := locales["en_US"][0]
	if englishFile == "" {
		fmt.Println("Could not find an i18n file for locale: en_US")
		return errors.New("Could not find an i18n file for locale: en_US")
	}

	englishStringInfos, err := fu.findI18nStrings(englishFile)

	if err != nil {
		fmt.Println(fmt.Sprintf("Couldn't find the english strings: %s", err.Error()))
		return err
	}

	additionalTranslations := getAdditionalTranslations(source, englishStringInfos)
	//	deletedTranslations :=
	//	updatedTranslations :=

	for locale, i18nFiles := range locales {
		translatedStrings, err := fu.findI18nStrings(i18nFiles[0])

		if err != nil {
			fmt.Println(fmt.Sprintf("Couldn't get the strings from %s: %s", locale, err.Error()))
			return err
		}
		err = addTranslations(translatedStrings, i18nFiles[0], additionalTranslations)
	}

	if err == nil {
		fmt.Printf("OK")
	}

	return err
}

func (fu *Fixup) inspectFile(file string) (translatedStrings []string, err error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		fu.Println(err)
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

func (fu *Fixup) findSourceStrings() (sourceStrings map[string]int, err error) {
	sourceStrings = make(map[string]int)
	files := getGoFiles(".")

	for _, file := range files {
		fileStrings, err := fu.inspectFile(file)
		if err != nil {
			fmt.Println("Error when inspecting go file: ", file)
			return sourceStrings, err
		}

		for _, string := range fileStrings {
			sourceStrings[string]++
		}
	}

	return
}

func (fu *Fixup) findI18nStrings(i18nFile string) (i18nStrings map[string]common.I18nStringInfo, err error) {
	i18nStrings = make(map[string]common.I18nStringInfo)
	//make(map[string]string)

	stringInfos, err := common.LoadI18nStringInfos(i18nFile)

	if err != nil {
		return nil, err
	}

	return common.CreateI18nStringInfoMap(stringInfos)
}

func getAdditionalTranslations(sourceTranslations map[string]int, englishTranslations map[string]common.I18nStringInfo) []string {
	additionalTranslations := []string{}
	//diffSlaveToMaster := []string{}
	//masterUpdates := map[string]string{} //[old] -> new

	//get master -> slave diff
	for id, _ := range sourceTranslations {
		if _, ok := englishTranslations[id]; !ok {
			additionalTranslations = append(additionalTranslations, id)
		}
	}
	return additionalTranslations
}

func addTranslations(localeMap map[string]common.I18nStringInfo, localeFile string, addTranslations []string) error {
	var err error
	/*
		//get slave -> master diff
		for id, translation := range diff {
			if master[id] == 0 {
				diffSlaveToMaster = append(diffSlaveToMaster, id)
			}
		}
		//updates
		for _, id := range additionalTranslations {
			var input rune
		addOrModified:
			fmt.Printf("Is %s an updated key [y,n]? ", id)
			fmt.Scanf("%c\n", &input)
			switch input {
			case 'y':
			chooseKey:
				fmt.Printf("What was the previous key?\n")
				for i, key := range diffSlaveToMaster {
					fmt.Printf("%d: %s\n", i+1, key)
				}

				var number int
				fmt.Scanf("%d\n", &number)
				// Check input
				if number <= 0 || number > len(diffSlaveToMaster) {
					goto chooseKey
				}

				// Get old ID, add the new ID to english, delete old ID.
				updateString := diffSlaveToMaster[number-1]
				fu.English[id] = id
				delete(fu.English, updateString)
				masterUpdates[updateString] = id

				diffSlaveToMaster = removeFromSlice(diffSlaveToMaster, number-1)
				// Take previous key and change to new key in all translations
				// (change translation in english also), mark as dirty.
			case 'n':
			default:
				goto addOrModified
			}
		}
		//new
	*/
	fmt.Println("Adding these strings to the translation file:")
	for _, id := range addTranslations {
		localeMap[id] = common.I18nStringInfo{ID: id, Translation: id}
		fmt.Println(id)
	}
	localeArray := common.I18nStringInfoMapValues2Array(localeMap)
	encodedLocale, err := json.MarshalIndent(localeArray, "", "   ")
	if err != nil {
		return err
	}
	fmt.Println(localeFile)
	fmt.Println(string(encodedLocale))
	err = ioutil.WriteFile(localeFile, encodedLocale, 0644)
	if err != nil {
		return err
	}

	//delete
	return err
}

func removeFromSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}
