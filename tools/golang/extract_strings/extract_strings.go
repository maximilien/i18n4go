package extract_strings

import (
	"os"
	"fmt"
	"strconv"

	"go/ast"
	"go/parser"
	"go/token"

	"io/ioutil"
	"encoding/json"
)

type I18nStringInfo struct {
	ID string `json:"id"`
	Translation string `json:"translation"`
}

type StringInfo struct {
	Filename string `json:"filename"`
	Value string `json:"value"`
	Offset int `json:"offset"`
	Line int `json:"line"`
	Column int `json:"column"`
}

type ExtractStrings struct {
	Filename string
	I18nFilename string
	ExcludedFilename string
	ExtractedStrings map[string]StringInfo
	FilteredStrings map[string]string
}

type ExcludedStrings struct {
	ExcludedStrings []string `json:"excludedStrings"`
}

var BLANKS = []string{", " , "\t", "\n", "\n\t", "\t\n"}

func NewExtractStrings(excludedFilename string) ExtractStrings {
	return ExtractStrings{Filename : "extracted_strings.json", 
					      ExcludedFilename : excludedFilename, 
					      ExtractedStrings : make(map[string]StringInfo), 
					      FilteredStrings : make(map[string]string)}
}

func (es * ExtractStrings) InspectFile(filename string) error {
	fmt.Println("Extracting strings from file:", filename)

	es.setFilename(filename)
	es.setI18nFilename(filename)

	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, filename, nil, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = es.loadExcludedStrings()
	if err != nil {
		fmt.Println(err)
		return err
	}

	es.excludeImports(astFile)

	es.extractString(astFile, fset)

	fmt.Printf("Extracted %d strings\n", len(es.ExtractedStrings))

	fmt.Println("Saving extracted strings to file:", es.Filename)
	err = es.saveExtractedStrings()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Saving extracted i18n strings to file:", es.I18nFilename)
	err = es.saveI18nStrings()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (es * ExtractStrings) InspectDir(dirName string, recursive bool) error {
	fmt.Println("Inspecting dir:", dirName)
	fmt.Println("recursive:", recursive)

	fset := token.NewFileSet()

	packages, err := parser.ParseDir(fset, dirName, nil, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for k, pkg := range packages {
		fmt.Println("Extracting string in package:", k)
		for file := range pkg.Files {
			es.InspectFile(file)
		}		
	}

	return nil
}

func (es *ExtractStrings) saveExtractedStrings() error {
	stringInfos := make([]StringInfo, 0)
	for _, stringInfo := range es.ExtractedStrings {
		stringInfos = append(stringInfos, stringInfo)
	}

	jsonData, err := json.Marshal(stringInfos)
    if err != nil {
        fmt.Println(err)
        return err
    }

    file, err := os.Create(es.Filename)
    if err != nil {
        fmt.Println(err)
        return err
    }

    file.Write(jsonData)
    defer file.Close()

    return nil
}

func (es *ExtractStrings) saveI18nStrings() error {
	i18nStringInfos := make([]I18nStringInfo, len(es.ExtractedStrings))
	i := 0
	for _, stringInfo := range es.ExtractedStrings {
		i18nStringInfos[i] = I18nStringInfo{ID: stringInfo.Value, Translation: stringInfo.Value}
		i++
	}

	jsonData, err := json.Marshal(i18nStringInfos)
    if err != nil {
        fmt.Println(err)
        return err
    }

    file, err := os.Create(es.I18nFilename)
    if err != nil {
        fmt.Println(err)
        return err
    }

    file.Write(jsonData)
    defer file.Close()

    return nil
}

func (es *ExtractStrings) setFilename(filename string) {
	es.Filename = filename + ".extracted.json"
}

func (es *ExtractStrings) setI18nFilename(filename string) {
	es.I18nFilename = filename + ".en.json"
}

func (es *ExtractStrings) loadExcludedStrings() error {
	fmt.Println("Excluding strings in file:", es.ExcludedFilename)

	content, err := ioutil.ReadFile(es.ExcludedFilename)
    if err != nil {
        fmt.Print(err)
        return err
    }

    var excludedStrings ExcludedStrings
    err = json.Unmarshal(content, &excludedStrings)
    if err != nil {
        fmt.Print(err)
        return err
    }

    for i := range excludedStrings.ExcludedStrings {
    	es.FilteredStrings[excludedStrings.ExcludedStrings[i]] = excludedStrings.ExcludedStrings[i]
    }

	return nil
}

func (es * ExtractStrings) extractString(f *ast.File, fset *token.FileSet) error {
		ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:			
			s, _ = strconv.Unquote(x.Value)
			if len(s) > 0 && x.Kind == token.STRING && s != "\t" && s != "\n" && s != " " && !es.filter(s) {
				position := fset.Position(n.Pos())
				stringInfo := StringInfo{Value: s, 
										 Filename: position.Filename, 
										 Offset: position.Offset, 
										 Line: position.Line, 
										 Column: position.Column}
				es.ExtractedStrings[s] = stringInfo
			}
		}
		return true
	})

	return nil
}

func (es * ExtractStrings) excludeImports(astFile *ast.File) {
	for i := range astFile.Imports {
		importString, _ := strconv.Unquote(astFile.Imports[i].Path.Value)
		es.FilteredStrings[importString] = importString
	}

}

func (es * ExtractStrings) filter(aString string) bool {
	for i := range BLANKS {
		if aString == BLANKS[i] {
			return true
		}
	}

	if es.FilteredStrings[aString] != "" {
		return true
	}
	return false
}