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

package common

import (
	"errors"
	"fmt"
	"strconv"

	"go/ast"
	"go/parser"
	"go/token"

	"github.com/maximilien/i18n4go/i18n4go/i18n"
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

	return nil, errors.New(fmt.Sprintf(i18n.T("Could not find imports for root node:\n\t{{.Arg0}}v\n", map[string]interface{}{"Arg0": astFile})))
}

func InspectFile(file string, options Options) (translatedStrings []string, err error) {
	defineAssignStmtMap := make(map[string][]ast.AssignStmt)
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, file, nil, parser.AllErrors)
	if err != nil {
		Println(options, err)
		return
	}

	ast.Inspect(astFile, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.AssignStmt:
			// inspect any potential translation string in defined / assigned statement nodes
			// add node to map if variable contains a translation string
			// eg: translation := "Hello {{.FirstName}}"
			//     T(translation)
			//     translation = "Hello {{.LastName}}"
			//     T(translation)
			inspectAssignStmt(defineAssignStmtMap, x)
		case *ast.CallExpr:
			// inspect any T()/t() or <MODULE>.T()/<MODULE>.t() (eg. i18n.T()) method calls using map
			/// then retrieve a list of translation strings that were passed into method
			translatedStrings = inspectCallExpr(translatedStrings, defineAssignStmtMap, x, options.QualifierFlag)
		}
		return true
	})

	return
}
func inspectCallExpr(translatedStrings []string, stmtMap map[string][]ast.AssignStmt, node *ast.CallExpr, qualifier string) []string {
	switch node.Fun.(type) {
	case *ast.Ident:
		funName := node.Fun.(*ast.Ident).Name
		// inspect any T() or t() method calls
		if funName == "T" || funName == "t" {
			translatedStrings = inspectTFunc(translatedStrings, stmtMap, *node)
		}

	case *ast.SelectorExpr:
		expr := node.Fun.(*ast.SelectorExpr)
		if ident, ok := expr.X.(*ast.Ident); ok {
			funName := expr.Sel.Name
			// inspect any <MODULE>.T() or <MODULE>.t() method calls (eg. i18n.T())
			if (ident.Name == qualifier || ident.Name == "i18n") && (funName == "T" || funName == "t") {
				translatedStrings = inspectTFunc(translatedStrings, stmtMap, *node)
			}
		}
	default:
		//Skip!
	}

	return translatedStrings
}

func inspectAssignStmt(stmtMap map[string][]ast.AssignStmt, node *ast.AssignStmt) {
	// use a hashmap for defined variables to a list of reassigned variables sharing the same var name
	if assignStmt, okIdent := node.Lhs[0].(*ast.Ident); okIdent {
		varName := assignStmt.Name
		switch node.Tok {
		case token.DEFINE:
			stmtMap[varName] = []ast.AssignStmt{}
		case token.ASSIGN:
			if _, exists := stmtMap[varName]; exists {
				stmtMap[varName] = append(stmtMap[varName], *node)
			}

		}
	}
}

func inspectStmt(translatedStrings []string, stmtMap map[string][]ast.AssignStmt, node ast.AssignStmt) []string {
	if strStmtArg, ok := node.Rhs[0].(*ast.BasicLit); ok {
		varName := node.Lhs[0].(*ast.Ident).Name
		translatedString, err := strconv.Unquote(strStmtArg.Value)
		if err != nil {
			panic(err.Error())
		}
		translatedStrings = append(translatedStrings, translatedString)
		// apply all translation ids from reassigned variables
		if _, exists := stmtMap[varName]; exists {
			for _, assignStmt := range stmtMap[varName] {
				strVarVal := assignStmt.Rhs[0].(*ast.BasicLit).Value
				translatedString, err := strconv.Unquote(strVarVal)
				if err != nil {
					panic(err.Error())
				}
				translatedStrings = append(translatedStrings, translatedString)

			}
		}
	}

	return translatedStrings
}

func inspectTFunc(translatedStrings []string, stmtMap map[string][]ast.AssignStmt, node ast.CallExpr) []string {
	if stringArg, ok := node.Args[0].(*ast.BasicLit); ok {
		translatedString, err := strconv.Unquote(stringArg.Value)
		if err != nil {
			panic(err.Error())
		}
		translatedStrings = append(translatedStrings, translatedString)
	}
	if idt, okIdt := node.Args[0].(*ast.Ident); okIdt {
		if obj := idt.Obj; obj != nil {
			if stmtArg, okStmt := obj.Decl.(*ast.AssignStmt); okStmt {
				translatedStrings = inspectStmt(translatedStrings, stmtMap, *stmtArg)
			}
		}
	}

	return translatedStrings
}
