package rewrite_package

import (
	"bytes"
	"fmt"
	"os"

	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"

	"github.com/maximilien/i18n4cf/cmds"
	"github.com/maximilien/i18n4cf/common"

	"path/filepath"
)

var (
	IMPORT_MAP = map[string]string{
		"":       `"github.com/cloudfoundry/cli/cf/i18n"`,
		"goi18n": `"github.com/nicksnyder/go-i18n/i18n"`,
	}
)

type rewritePackage struct {
	options cmds.Options

	Filename            string
	OutputDirname       string
	I18nStringsFilename string

	ExtractedStrings map[string]common.StringInfo

	TotalStrings int
	TotalFiles   int
}

func NewRewritePackage(options cmds.Options) rewritePackage {
	return rewritePackage{options: options,
		Filename:            options.FilenameFlag,
		OutputDirname:       options.OutputDirFlag,
		I18nStringsFilename: options.I18nStringsFilenameFlag,
		TotalStrings:        0,
		TotalFiles:          0}
}

func (rp *rewritePackage) Options() cmds.Options {
	return rp.options
}

func (rp *rewritePackage) Println(a ...interface{}) (int, error) {
	if rp.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (rp *rewritePackage) Printf(msg string, a ...interface{}) (int, error) {
	if rp.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (rp *rewritePackage) Run() error {
	rp.Println("gi18n: rewriting strings for source file:", rp.Filename)
	if rp.options.DryRunFlag {
		rp.Println("WARNING running in -dry-run mode")
	}

	fileSet := token.NewFileSet()

	var absFilePath = rp.Filename
	if !filepath.IsAbs(absFilePath) {
		absFilePath = filepath.Join(os.Getenv("PWD"), absFilePath)
	}

	astFile, err := parser.ParseFile(fileSet, absFilePath, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		rp.Println(err)
		return err
	}

	err = rp.rewriteImports(astFile)
	if err != nil {
		rp.Println("gi18n: error rewriting the imports on AST file:", err.Error())
		return err
	}

	err = rp.appendInitFunc(astFile)
	if err != nil {
		rp.Println("gi18n: error appending init() to AST file:", err.Error())
		return err
	}

	err = rp.insertTFuncCall(astFile)
	if err != nil {
		rp.Println("gi18n: error appending T() to AST file:", err.Error())
		return err
	}

	err = rp.saveASTFile(astFile, fileSet)
	if err != nil {
		rp.Println("gi18n: error saving AST file:", err.Error())
		return err
	}

	return err
}

func (rp *rewritePackage) insertTFuncCall(astFile *ast.File) error {
	declarations := astFile.Decls[1:]

	for _, decl := range declarations {
		ast.Inspect(decl, func(node ast.Node) bool {
			switch node.(type) {
			case *ast.CallExpr:
				if !rp.callExprTFunc(node.(*ast.CallExpr)) {
					return false // don't recurse infinitely
				}
			case *ast.AssignStmt:
				rp.assignStmtTFunc(node.(*ast.AssignStmt))
			case *ast.ValueSpec:
				rp.valueSpecTFunc(node.(*ast.ValueSpec))
			case *ast.CompositeLit:
				rp.compositeLitTFunc(node.(*ast.CompositeLit))
			}

			return true
		})
	}

	return nil
}

func (rp *rewritePackage) compositeLitTFunc(compositeLit *ast.CompositeLit) bool {
	for index, arg := range compositeLit.Elts {
		if asLit, ok := arg.(*ast.BasicLit); ok {
			compositeLit.Elts[index] = rp.wrapBasicLitWithT(asLit)
		}
	}

	return true
}

func (rp *rewritePackage) assignStmtTFunc(assignStmt *ast.AssignStmt) bool {
	for index, arg := range assignStmt.Rhs {
		if asLit, ok := arg.(*ast.BasicLit); ok {
			assignStmt.Rhs[index] = rp.wrapBasicLitWithT(asLit)
		}
	}

	return true
}

func (rp *rewritePackage) valueSpecTFunc(valueSpec *ast.ValueSpec) bool {
	for index, arg := range valueSpec.Values {
		if asLit, ok := arg.(*ast.BasicLit); ok {
			valueSpec.Values[index] = rp.wrapBasicLitWithT(asLit)
		}
	}

	return true
}

func (rp *rewritePackage) callExprTFunc(callExpr *ast.CallExpr) bool {
	callFuncIdent, ok := callExpr.Fun.(*ast.Ident)
	if ok && callFuncIdent.Name == "T" { // yeah, not the best
		return false
	}

	for index, arg := range callExpr.Args {
		if asLit, ok := arg.(*ast.BasicLit); ok {
			callExpr.Args[index] = rp.wrapBasicLitWithT(asLit)
		}
	}

	return true
}

func (rp *rewritePackage) wrapBasicLitWithT(basicLit *ast.BasicLit) *ast.CallExpr {
	tIdent := &ast.Ident{Name: "T"}
	return &ast.CallExpr{Fun: tIdent, Args: []ast.Expr{basicLit}}
}

func (rp *rewritePackage) appendInitFunc(astFile *ast.File) error {
	fileSet := token.NewFileSet()
	cwd, err := os.Getwd()
	pathToFile := filepath.Join(cwd, "..", "..", "cmds", "rewrite_package", "code_snippets", "init_func.go.example")
	astSnippet, err := parser.ParseFile(fileSet, pathToFile, nil, 0)
	if err != nil {
		return err
	}

	for _, x := range astSnippet.Decls[0:len(astSnippet.Decls)] {
		ast.Inspect(x, func(node ast.Node) bool {
			if node == nil {
				return false
			}

			switch node.(type) {
			case *ast.Ident:
				node.(*ast.Ident).NamePos = 1
			case *ast.BasicLit:
				if node.(*ast.BasicLit).Value == `"__PACKAGE__NAME__"` {
					packageName := "\"" + astFile.Name.Name + "\""
					node.(*ast.BasicLit).Value = packageName
				}
			}
			return true
		})
	}

	astFile.Decls = append(astFile.Decls[:1], append(astSnippet.Decls[0:len(astSnippet.Decls)], astFile.Decls[1:]...)...)

	return nil
}

func (rp *rewritePackage) rewriteImports(astFile *ast.File) error {
	importDecl, err := common.ImportsForASTFile(astFile)
	if err != nil {
		return err
	}

	for importName, importPath := range IMPORT_MAP {
		importSpec := &ast.ImportSpec{
			Name: &ast.Ident{Name: importName},
			Path: &ast.BasicLit{Value: importPath, Kind: token.STRING},
		}

		importDecl.Specs = append(importDecl.Specs, importSpec)
	}

	return nil
}

func (rp *rewritePackage) saveASTFile(astFile *ast.File, fileSet *token.FileSet) error {
	var buffer bytes.Buffer
	if err := format.Node(&buffer, fileSet, astFile); err != nil {
		return err
	}

	fileName := filepath.Base(rp.Filename)
	pathToFile := filepath.Join(rp.OutputDirname, fileName)

	fileInfo, err := os.Stat(rp.Filename)
	if err != nil {
		return err
	}

	rp.Println("saving file to path", pathToFile)
	ioutil.WriteFile(pathToFile, buffer.Bytes(), fileInfo.Mode())

	return nil
}
