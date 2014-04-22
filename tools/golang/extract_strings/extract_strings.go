package extract_strings

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func InspectFile(fileName string) error {
	fmt.Println("Extracting strings from file:", fileName)

	fset := token.NewFileSet() // positions are relative to fset

	astFile, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		fmt.Println(err)
		return err
	}

	extractString(astFile, fset)

	return nil
}

func InspectDir(dirName string, recursive bool) error {
	fmt.Println("Inspecting dir:", dirName)
	fmt.Println("recursive:", recursive)

	fset := token.NewFileSet() // positions are relative to fset

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
			if s != "" && x.Kind == token.STRING {
				fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
			}
		}
		return true
	})
}