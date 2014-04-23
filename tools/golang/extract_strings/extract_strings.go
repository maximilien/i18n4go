package extract_strings

import (
	"fmt"

	"go/ast"
	"go/parser"
	"go/token"

	"io/ioutil"
	"encoding/json"
)

type StringInfo struct {
	Value string "json:value"
	LineNo int "json:lineNo"
	Pos int "json:pos"
}

type ExtractStrings struct {
	Filename string
	ExcludedFilename string
	ExtractedStrings map[string]StringInfo
	FilteredStrings map[string]string
}

type ExcludedStrings struct {
	ExcludedStrings []string "json:excludedStrings"
}

var BLANKS = []string{"\"\"", "\" \"", "\"\\t\"", "\"\\n\"", "\"\\n\\t\"", "\"\\t\\n\""}

func NewExtractStrings(excludedFilename string) ExtractStrings {
	return ExtractStrings{Filename : "extracted_strings.json", 
					      ExcludedFilename : excludedFilename, 
					      ExtractedStrings : make(map[string]StringInfo), 
					      FilteredStrings : make(map[string]string)}
}

func (es * ExtractStrings) InspectFile(filename string) error {
	fmt.Println("Extracting strings from file:", filename)

	es.setFilename(filename)

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

	//TODO: save
	fmt.Println("TODO: Saving to file:", es.Filename)
	//END

	return nil
}

func (es * ExtractStrings)  InspectDir(dirName string, recursive bool) error {
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

func (es *ExtractStrings) setFilename(filename string) {
	es.Filename = filename + ".extracted.json"
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

func (es * ExtractStrings) extractString(f *ast.File, fset *token.FileSet) {
		ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:			
			s = x.Value
			if len(s) > 0 && x.Kind == token.STRING && s != "\t" && s != "\n" && s != " " && !es.filter(s) {
				fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
			}
		}
		return true
	})
}

func (es * ExtractStrings) excludeImports(astFile *ast.File) {
		for i := range astFile.Imports {
		es.FilteredStrings[astFile.Imports[i].Path.Value] = astFile.Imports[i].Path.Value
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