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

	err = rp.saveASTFile(astFile, fileSet)
	if err != nil {
		rp.Println("gi18n: error saving AST file:", err.Error())
		return err
	}

	return err
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
