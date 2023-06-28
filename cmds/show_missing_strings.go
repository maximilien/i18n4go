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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/spf13/cobra"

	"github.com/maximilien/i18n4go/common"
)

type showMissingStrings struct {
	options common.Options

	I18nStringInfos     []common.I18nStringInfo
	TranslatedStrings   []string
	I18nStringsFilename string
	Directory           string
}

func NewShowMissingStrings(options *common.Options) *showMissingStrings {
	return &showMissingStrings{
		options:             *options,
		Directory:           options.DirnameFlag,
		I18nStringsFilename: options.I18nStringsFilenameFlag,
		TranslatedStrings:   []string{},
	}
}

// NewShowMissingStringsCommand implements 'i18n show-missing-strings' command
func NewShowMissingStringsCommand(options *common.Options) *cobra.Command {
	showMissingStringsCmd := &cobra.Command{
		Use:   "show-missing-strings",
		Short: "Shows missing strings in translations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewShowMissingStrings(options).Run()
		},
	}

	showMissingStringsCmd.Flags().StringVarP(&options.DirnameFlag, "directory", "d", "", "the directory containing the go files to validate")
	showMissingStringsCmd.Flags().StringVar(&options.I18nStringsFilenameFlag, "i18n-strings-filename", "", "a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command")

	// TODO: setup options and params for Cobra command here using common.Options

	return showMissingStringsCmd
}

func (sms *showMissingStrings) Options() common.Options {
	return sms.options
}

func (sms *showMissingStrings) Println(a ...interface{}) (int, error) {
	if sms.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (sms *showMissingStrings) Printf(msg string, a ...interface{}) (int, error) {
	if sms.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (sms *showMissingStrings) Run() error {
	return sms.showMissingStrings()
}

func (sms *showMissingStrings) showMissingStrings() error {

	//Load en_US.all.json

	stringInfos, err := common.LoadI18nStringInfos(sms.I18nStringsFilename)
	if err != nil {
		return err
	}
	sms.I18nStringInfos = stringInfos

	//Run AST to get list of strings
	err = sms.parseFiles()
	if err != nil {
		return err
	}

	//Compare list of strings with <lang>.all.json
	err = sms.showMissingTranslatedStrings()
	if err != nil {
		return err
	}

	//Compare list of translated strings with strings in codebase
	return sms.showExtraStrings()
}

func (sms *showMissingStrings) parseFiles() error {
	sourceFiles, _ := getFilesAndDir(sms.Directory)
	for _, sourceFile := range sourceFiles {
		err := sms.inspectFile(sourceFile)
		if err != nil {
			return err
		}
	}

	return nil
}
func (sms *showMissingStrings) inspectFile(filename string) error {
	fset := token.NewFileSet()

	var absFilePath = filename
	if !filepath.IsAbs(absFilePath) {
		absFilePath = filepath.Join(os.Getenv("PWD"), absFilePath)
	}

	fileInfo, err := common.GetAbsFileInfo(absFilePath)
	if err != nil {
		sms.Println(err)
	}

	if strings.HasPrefix(fileInfo.Name(), ".") || !strings.HasSuffix(fileInfo.Name(), ".go") {
		sms.Println("WARNING ignoring file:", absFilePath)
		return nil
	}

	astFile, err := parser.ParseFile(fset, absFilePath, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		sms.Println(err)
		return err
	}

	return sms.extractString(astFile, fset, filename)
}

func (sms *showMissingStrings) extractString(f *ast.File, fset *token.FileSet, filename string) error {
	ast.Inspect(f, func(n ast.Node) bool {
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

						sms.Println("Adding to translated strings:", translatedString)
						sms.TranslatedStrings = append(sms.TranslatedStrings, filename+": "+translatedString)
					}
				}
			default:
				//Skip!
			}
		}
		return true
	})

	return nil
}

func (sms *showMissingStrings) showMissingTranslatedStrings() error {
	missingStrings := false
	for _, codeString := range sms.TranslatedStrings {
		if !sms.stringInStringInfos(codeString, sms.I18nStringInfos) {
			fmt.Println("Missing:", codeString)
			missingStrings = true
		}
	}

	if missingStrings {
		return errors.New("Missing Strings!")
	}

	return nil
}

func splitFilePathAndString(str string) (string, string) {
	splitFileStr := strings.SplitAfterN(str, ": ", 2)
	return splitFileStr[0], splitFileStr[1]
}

func (sms *showMissingStrings) stringInStringInfos(str string, list []common.I18nStringInfo) bool {
	_, translatedStr := splitFilePathAndString(str)
	for _, stringInfo := range list {
		if translatedStr == stringInfo.ID {
			sms.Println("Found", stringInfo.ID, "UNDER", str)
			return true
		}
	}

	return false
}

// Compares translated strings with strings in codebase.
func (sms *showMissingStrings) showExtraStrings() error {
	additionalStrings := false
	for _, stringInfo := range sms.I18nStringInfos {
		if !stringInTranslatedStrings(stringInfo.ID, sms.TranslatedStrings) {
			fmt.Println("Additional:", stringInfo.ID)
			additionalStrings = true
		}
	}
	if additionalStrings {
		return errors.New("Additional Strings!")
	}
	return nil
}

func stringInTranslatedStrings(stringInfoID string, list []string) bool {
	for _, fileAndStr := range list {
		_, translatedStr := splitFilePathAndString(fileAndStr)
		if translatedStr == stringInfoID {
			return true
		}
	}

	return false
}
