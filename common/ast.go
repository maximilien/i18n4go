package common

import (
	"errors"
	"fmt"

	"go/ast"
)

func ImportsForASTFile(astFile *ast.File) (*ast.GenDecl, error) {
	for _, declaration := range astFile.Decls {
		decl, ok := declaration.(*ast.GenDecl)
		if !ok || len(decl.Specs) == 0 {
			continue
		}

		if _, ok = decl.Specs[0].(*ast.ImportSpec); ok {
			return decl, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Could not find imports for root node:\n\t%#v\n", astFile))
}
