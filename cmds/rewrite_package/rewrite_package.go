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
	"strings"
)

const (
	INIT_CODE_SNIPPET = `package __PACKAGE__NAME__

import (
	"github.com/cloudfoundry/cli/cf/i18n"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

var T goi18n.TranslateFunc

func init() {
	var err error
	T, err = i18n.Init("__PACKAGE__NAME__", i18n.GetResourcesPath())
	if err != nil {
		panic(err)
	}
}`
)

type rewritePackage struct {
	options cmds.Options

	Filename            string
	OutputDirname       string
	I18nStringsFilename string

	Dirname string
	Recurse bool

	ExtractedStrings map[string]common.I18nStringInfo

	TotalStrings int
	TotalFiles   int
}

func NewRewritePackage(options cmds.Options) rewritePackage {
	return rewritePackage{options: options,
		Filename:            options.FilenameFlag,
		OutputDirname:       options.OutputDirFlag,
		I18nStringsFilename: options.I18nStringsFilenameFlag,

		Dirname: options.DirnameFlag,
		Recurse: options.RecurseFlag,
	}
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
	var err error

	if err = rp.loadStringsToBeTranslated(); err != nil {
		return err
	}

	if rp.options.FilenameFlag != "" {
		err = rp.processFilename(rp.options.FilenameFlag)
	} else {
		err = rp.processDir(rp.options.DirnameFlag, rp.options.RecurseFlag)
	}

	rp.Println()
	rp.Println("Total files parsed:", rp.TotalFiles)
	rp.Println("Total extracted strings:", rp.TotalStrings)
	return err
}

func (rp *rewritePackage) loadStringsToBeTranslated() error {
	if rp.I18nStringsFilename != "" {
		stringList, err := common.LoadI18nStringInfos(rp.I18nStringsFilename)
		if err != nil {
			return err
		}

		rp.ExtractedStrings, err = common.CreateI18nStringInfoMap(stringList)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rp *rewritePackage) processDir(dirName string, recursive bool) error {
	rp.Printf("gi18n: rewriting strings in dir %s, recursive: %t\n", dirName, recursive)
	rp.Println()

	fileInfos, _ := ioutil.ReadDir(dirName)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			if recursive {
				rp.processDir(filepath.Join(dirName, fileInfo.Name()), recursive)
			} else {
				continue
			}
		} else if filepath.Base(fileInfo.Name()) != "i18n_init.go" {
			err := rp.processFilename(filepath.Join(dirName, fileInfo.Name()))
			if err != nil {
				rp.Println(err)
				return err
			}
		}
	}

	return nil
}

func (rp *rewritePackage) processFilename(fileName string) error {
	rp.TotalFiles += 1
	rp.Println("gi18n: rewriting strings for source file:", fileName)

	fileSet := token.NewFileSet()

	var absFilePath = fileName
	if !filepath.IsAbs(absFilePath) {
		absFilePath = filepath.Join(os.Getenv("PWD"), absFilePath)
	}

	astFile, err := parser.ParseFile(fileSet, absFilePath, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		rp.Println(err)
		return err
	}

	outputDir := filepath.Join(rp.OutputDirname, filepath.Dir(rp.relativePathForFile(fileName)))
	err = rp.addInitFuncToPackage(astFile.Name.Name, outputDir)
	if err != nil {
		rp.Println("gi18n: error adding init() func to package:", err.Error())
		return err
	}

	err = rp.insertTFuncCall(astFile)
	if err != nil {
		rp.Println("gi18n: error appending T() to AST file:", err.Error())
		return err
	}

	relativeFilePath := rp.relativePathForFile(fileName)
	err = rp.saveASTFile(relativeFilePath, fileName, astFile, fileSet)
	if err != nil {
		rp.Println("gi18n: error saving AST file:", err.Error())
		return err
	}

	return err
}

func (rp *rewritePackage) insertTFuncCall(astFile *ast.File) error {
	rp.Println("gi18n: inserting T() calls for strings that need to be translated")
	var declarations []ast.Decl
	if len(astFile.Imports) > 0 {
		declarations = astFile.Decls[1:]
	} else {
		declarations = astFile.Decls[0:]
	}

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
			case *ast.KeyValueExpr:
				rp.keyValueExprTFunc(node.(*ast.KeyValueExpr))
			case *ast.ReturnStmt:
				rp.returnStmtTFunc(node.(*ast.ReturnStmt))
			case *ast.BinaryExpr:
				rp.binaryExprTFunc(node.(*ast.BinaryExpr))
			case *ast.IndexExpr:
				rp.indexExprTFunc(node.(*ast.IndexExpr))
			}

			return true
		})
	}

	return nil
}

func (rp *rewritePackage) indexExprTFunc(indexExpr *ast.IndexExpr) {
	if index, ok := indexExpr.Index.(*ast.BasicLit); ok {
		indexExpr.Index = rp.wrapBasicLitWithT(index)
	}
}

func (rp *rewritePackage) binaryExprTFunc(binaryExpr *ast.BinaryExpr) {
	if x, ok := binaryExpr.X.(*ast.BasicLit); ok {
		binaryExpr.X = rp.wrapBasicLitWithT(x)
	}

	if y, ok := binaryExpr.Y.(*ast.BasicLit); ok {
		binaryExpr.Y = rp.wrapBasicLitWithT(y)
	}
}

func (rp *rewritePackage) returnStmtTFunc(returnStmt *ast.ReturnStmt) {
	for index, arg := range returnStmt.Results {
		if asLit, ok := arg.(*ast.BasicLit); ok {
			returnStmt.Results[index] = rp.wrapBasicLitWithT(asLit)
		}
	}
}

func (rp *rewritePackage) keyValueExprTFunc(keyValueExpr *ast.KeyValueExpr) {
	if key, ok := keyValueExpr.Key.(*ast.BasicLit); ok {
		keyValueExpr.Key = rp.wrapBasicLitWithT(key)
	}

	if value, ok := keyValueExpr.Value.(*ast.BasicLit); ok {
		keyValueExpr.Value = rp.wrapBasicLitWithT(value)
	}
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

func (rp *rewritePackage) wrapBasicLitWithT(basicLit *ast.BasicLit) ast.Expr {
	if basicLit.Kind != token.STRING {
		return basicLit
	}

	valueWithoutQuotes := basicLit.Value[1 : len(basicLit.Value)-1]
	_, ok := rp.ExtractedStrings[valueWithoutQuotes]
	if !ok && rp.ExtractedStrings != nil {
		return basicLit
	}

	rp.TotalStrings++
	tIdent := &ast.Ident{Name: "T"}
	return &ast.CallExpr{Fun: tIdent, Args: []ast.Expr{basicLit}}
}

func (rp *rewritePackage) addInitFuncToPackage(packageName, outputDir string) error {
	rp.Println("gi18n: adding init func to package:", packageName, " to output dir:", outputDir)

	common.CreateOutputDirsIfNeeded(outputDir)
	content := strings.Replace(INIT_CODE_SNIPPET, "__PACKAGE__NAME__", packageName, -1)
	return ioutil.WriteFile(filepath.Join(outputDir, "i18n_init.go"), []byte(content), 0666)
}

func (rp *rewritePackage) appendInitFunc(astFile *ast.File) error {
	fileSet := token.NewFileSet()

	file, err := common.CreateTmpFile(INIT_CODE_SNIPPET)
	if err != nil {
		return err
	}

	pathToFile := file.Name()
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

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

func (rp *rewritePackage) saveASTFile(relativeFilePath, fileName string, astFile *ast.File, fileSet *token.FileSet) error {
	var buffer bytes.Buffer
	if err := format.Node(&buffer, fileSet, astFile); err != nil {
		return err
	}

	pathToFile := filepath.Join(rp.OutputDirname, relativeFilePath)
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	common.CreateOutputDirsIfNeeded(filepath.Dir(pathToFile))

	rp.Println("saving file to path", pathToFile)
	ioutil.WriteFile(pathToFile, buffer.Bytes(), fileInfo.Mode())

	return nil
}

func (rp *rewritePackage) relativePathForFile(fileName string) string {
	if rp.Dirname != "" {
		return strings.Replace(fileName, rp.Dirname, "", -1)
	} else {
		return filepath.Base(fileName)
	}
}
