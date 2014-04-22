package extract_strings

import (
	"fmt"

	"go/ast"
	"go/parser"
	"go/token"
)

type StringInfo struct {
	Value string
	Filename string
	Pos int
}

type ExtactString struct {
	ExtractedStrings map[string]StringInfo
	FilteredStrings map[string]StringInfo
}

var BLANKS = []string{"\"\"", "\" \"", "\"\\t\"", "\"\\n\"", "\"\\n\\t\"", "\"\\t\\n\""}

var filteredStrings map[string]string

func InspectFile(fileName string) error {
	fmt.Println("Extracting strings from file:", fileName)

	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	filteredStrings = make(map[string]string)
	addImports(astFile)

	extractString(astFile, fset)

	return nil
}

func InspectDir(dirName string, recursive bool) error {
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
			InspectFile(file)
		}		
	}

	return nil
}

func extractString(f *ast.File, fset *token.FileSet) {
		ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:			
			s = x.Value
			if len(s) > 0 && x.Kind == token.STRING && s != "\t" && s != "\n" && s != " " && !filter(s) {
				fmt.Println("length:", len(s))
				fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
			}
		}
		return true
	})
}

func addImports(astFile *ast.File) {
		for i := range astFile.Imports {
		filteredStrings[astFile.Imports[i].Path.Value] = astFile.Imports[i].Path.Value
	}

}

func filter(aString string) bool {
	for i := range BLANKS {
		if aString == BLANKS[i] {
			return true
		}
	}

	if filteredStrings[aString] != "" {
		return true
	}
	return false
}