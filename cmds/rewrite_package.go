package cmds

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"

	"github.com/maximilien/i18n4go/common"

	"path/filepath"
	"strconv"
	"strings"
)

const (
	INIT_CODE_SNIPPET = `package __PACKAGE__NAME__

import (
	"path/filepath"

	i18n "github.com/maximilien/i18n4go/i18n"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

var T goi18n.TranslateFunc

func init() {
	T = i18n.Init(__FULL_IMPORT_PATH__, i18n.GetResourcesPath())
}`
)

type rewritePackage struct {
	options common.Options

	Filename                string
	OutputDirname           string
	I18nStringsFilename     string
	I18nStringsDirname      string
	RootPath                string
	InitCodeSnippetFilename string

	Dirname string
	Recurse bool

	ExtractedStrings        map[string]common.I18nStringInfo
	UpdatedExtractedStrings map[string]common.I18nStringInfo
	SaveExtractedStrings    bool

	TotalStrings int
	TotalFiles   int

	IgnoreRegexp *regexp.Regexp
}

func NewRewritePackage(options common.Options) rewritePackage {
	var compiledRegexp *regexp.Regexp
	if options.IgnoreRegexpFlag != "" {
		compiledReg, err := regexp.Compile(options.IgnoreRegexpFlag)
		if err != nil {
			fmt.Println("WARNING compiling ignore-regexp:", err)
		}
		compiledRegexp = compiledReg
	}

	return rewritePackage{options: options,
		Filename:                options.FilenameFlag,
		OutputDirname:           options.OutputDirFlag,
		I18nStringsFilename:     options.I18nStringsFilenameFlag,
		I18nStringsDirname:      options.I18nStringsDirnameFlag,
		RootPath:                options.RootPathFlag,
		InitCodeSnippetFilename: options.InitCodeSnippetFilenameFlag,

		ExtractedStrings:        nil,
		UpdatedExtractedStrings: nil,
		SaveExtractedStrings:    false,

		Dirname:      options.DirnameFlag,
		Recurse:      options.RecurseFlag,
		IgnoreRegexp: compiledRegexp,
	}
}

func (rp *rewritePackage) Options() common.Options {
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

	if rp.options.FilenameFlag != "" {
		if err = rp.loadStringsToBeTranslated(rp.I18nStringsFilename); err != nil {
			return err
		}
		err = rp.processFilename(rp.options.FilenameFlag)
	} else {
		err = rp.processDir(rp.options.DirnameFlag, rp.options.RecurseFlag)
	}

	rp.Println()
	rp.Println("Total files parsed:", rp.TotalFiles)
	rp.Println("Total rewritten strings:", rp.TotalStrings)
	return err
}

func (rp *rewritePackage) loadStringsToBeTranslated(fileName string) error {
	if fileName != "" {
		stringList, err := common.LoadI18nStringInfos(fileName)
		if err != nil {
			return err
		}

		rp.ExtractedStrings, err = common.CreateI18nStringInfoMap(stringList)
		if err != nil {
			return err
		}

		rp.UpdatedExtractedStrings = common.CopyI18nStringInfoMap(rp.ExtractedStrings)
	}

	return nil
}

func (rp *rewritePackage) processDir(dirName string, recursive bool) error {
	rp.Printf("i18n4go: rewriting strings in dir %s, recursive: %t\n", dirName, recursive)
	rp.Println()

	fileInfos, _ := ioutil.ReadDir(dirName)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			if recursive {
				rp.processDir(filepath.Join(dirName, fileInfo.Name()), recursive)
			} else {
				continue
			}
		} else if rp.ignoreFile(filepath.Base(fileInfo.Name())) {
			i18nFilename := rp.I18nStringsFilename
			if rp.I18nStringsDirname != "" {
				i18nFilename = filepath.Base(fileInfo.Name()) + "." + rp.options.SourceLanguageFlag + ".json"
			}

			rp.I18nStringsFilename = filepath.Join(rp.I18nStringsDirname, i18nFilename)
			rp.Printf("i18n4go: loading JSON strings from file: %s\n", rp.I18nStringsFilename)
			if err := rp.loadStringsToBeTranslated(rp.I18nStringsFilename); err != nil {
				rp.Println("i18n4go: WARNING could not find JSON file:", rp.I18nStringsFilename, err.Error())
				rp.resetProcessing()
				continue
			}
			err := rp.processFilename(filepath.Join(dirName, fileInfo.Name()))
			if err != nil {
				rp.Println(err)
				return err
			}

			if rp.I18nStringsDirname != "" {
				rp.resetProcessing()
			}
		}
	}

	return nil
}

func (rp *rewritePackage) resetProcessing() {
	rp.ExtractedStrings = nil
	rp.UpdatedExtractedStrings = nil
	rp.I18nStringsFilename = ""
	rp.SaveExtractedStrings = false
}

func (rp *rewritePackage) ignoreFile(fileName string) bool {
	return fileName != "i18n_init.go" &&
		!strings.HasPrefix(fileName, ".") &&
		strings.HasSuffix(fileName, ".go") &&
		rp.IgnoreRegexp != nil && !rp.IgnoreRegexp.MatchString(fileName)
}

func (rp *rewritePackage) processFilename(fileName string) error {
	rp.TotalFiles += 1
	rp.Println("i18n4go: rewriting strings for source file:", fileName)

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

	if strings.HasSuffix(fileName, "_test.go") {
		rp.Println("cowardly refusing to translate the strings in test file:", fileName)
		return nil
	}

	importPath, err := rp.determineImportPath(absFilePath)
	if err != nil {
		rp.Println("i18n4go: error determining the import path:", err.Error())
		return err
	}

	if rp.OutputDirname == "" {
		rp.OutputDirname = filepath.Dir(fileName)
	}

	outputDir := filepath.Join(rp.OutputDirname, filepath.Dir(rp.relativePathForFile(fileName)))
	err = rp.addInitFuncToPackage(astFile.Name.Name, outputDir, importPath)
	if err != nil {
		rp.Println("i18n4go: error adding init() func to package:", err.Error())
		return err
	}

	err = rp.insertTFuncCall(astFile)
	if err != nil {
		rp.Println("i18n4go: error appending T() to AST file:", err.Error())
		return err
	}

	relativeFilePath := rp.relativePathForFile(fileName)
	err = rp.saveASTFile(relativeFilePath, fileName, astFile, fileSet)
	if err != nil {
		rp.Println("i18n4go: error saving AST file:", err.Error())
		return err
	}

	if rp.SaveExtractedStrings {
		i18nStringInfos := common.I18nStringInfoMapValues2Array(rp.UpdatedExtractedStrings)
		err := common.SaveI18nStringInfos(rp, rp.Options(), i18nStringInfos, rp.I18nStringsFilename)
		if err != nil {
			rp.Println("i18n4go: error saving updated i18n strings file:", err.Error())
			return err
		}
	}

	return err
}

func (rp *rewritePackage) determineImportPath(filePath string) (string, error) {
	dirName := filepath.Dir(filePath)
	if rp.options.RootPathFlag == "" {
		rp.Println("i18n4go: using the PWD as the rootPath:", os.Getenv("PWD"))
		rp.RootPath = os.Getenv("PWD")
	}
	rp.Println("i18n4go: determining import path using root path:", rp.RootPath)
	pkg, err := build.Default.ImportDir(rp.RootPath, build.ImportMode(1))
	if err != nil {
		rp.Println("i18n4go: error getting root path import:", err.Error)
		return "", err
	}
	rp.Println("i18n4go: got a root pkg with import path:", pkg.ImportPath)

	otherPkg, err := build.Default.ImportDir(dirName, build.ImportMode(0))
	if err != nil {
		rp.Println("i18n4go: error getting root path import:", err.Error)
		return "", err
	}
	rp.Println("i18n4go: got a pkg with import:", otherPkg.ImportPath)

	importPath := otherPkg.ImportPath
	importPath = strings.Replace(importPath, pkg.ImportPath, "", 1)
	if strings.HasPrefix(importPath, "/") {
		importPath = strings.TrimLeft(importPath, "/")
	}
	rp.Println("i18n4go: using import path as:", importPath)

	return importPath, nil
}

func (rp *rewritePackage) insertTFuncCall(astFile *ast.File) error {
	rp.Println("i18n4go: inserting T() calls for strings that need to be translated")
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
	if ok && (callFuncIdent.Name == "T") { // yeah, not the best
		return false
	}

	switch len(callExpr.Args) {
	case 0:
		return false
	case 1:
		if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
			callExpr.Args[0] = rp.wrapBasicLitWithT(basicLit)
		}
	default:
		rp.wrapMultiArgsCallExpr(callExpr)
	}

	return true
}

func (rp *rewritePackage) wrapMultiArgsCallExpr(callExpr *ast.CallExpr) {
	for i, arg := range callExpr.Args {
		if basicLit, ok := arg.(*ast.BasicLit); ok {
			if basicLit.Kind == token.STRING {
				valueWithoutQuotes, _ := strconv.Unquote(basicLit.Value) //basicLit.Value[1 : len(basicLit.Value)-1]

				if common.IsTemplatedString(valueWithoutQuotes) {
					rp.wrapCallExprWithTemplatedT(basicLit, callExpr, i)
				} else if common.IsInterpolatedString(valueWithoutQuotes) {
					rp.wrapCallExprWithInterpolatedT(basicLit, callExpr, i)
				} else {
					rp.wrapExprArgs(callExpr.Args)
				}

			} else {
				rp.wrapExprArgs(callExpr.Args)
			}
		}
	}
}

func (rp *rewritePackage) wrapCallExprWithInterpolatedT(basicLit *ast.BasicLit, callExpr *ast.CallExpr, argIndex int) {
	valueWithoutQuotes, _ := strconv.Unquote(basicLit.Value)

	i18nStringInfo, ok := rp.ExtractedStrings[valueWithoutQuotes]
	if !ok && rp.ExtractedStrings != nil {
		rp.wrapExprArgs(callExpr.Args)
		return
	}

	templatedString := common.ConvertToTemplatedString(valueWithoutQuotes)
	basicLit.Value = strconv.Quote(templatedString)

	if rp.ExtractedStrings != nil {
		rp.updateExtractedStrings(i18nStringInfo, templatedString)
	}

	rp.wrapCallExprWithTemplatedT(basicLit, callExpr, argIndex)
}

func (rp *rewritePackage) wrapCallExprWithTemplatedT(basicLit *ast.BasicLit, callExpr *ast.CallExpr, argIndex int) {
	templatedCallExpr := rp.wrapBasicLitWithTemplatedT(basicLit, callExpr.Args, callExpr, argIndex)
	if templatedCallExpr != callExpr {
		newArgs := []ast.Expr{}

		if argIndex != 0 {
			for i, arg := range callExpr.Args {
				if i < argIndex {
					newArgs = append(newArgs, arg)
				}
			}
		}

		callExpr.Args = append(newArgs, templatedCallExpr)
	}
}

func (rp *rewritePackage) wrapExprArgs(exprArgs []ast.Expr) {
	for i, _ := range exprArgs {
		if basicLit, ok := exprArgs[i].(*ast.BasicLit); ok {
			exprArgs[i] = rp.wrapBasicLitWithT(basicLit)
		} else if callExpr, ok := exprArgs[i].(*ast.CallExpr); ok {
			rp.callExprTFunc(callExpr)
		}
	}
}

func (rp *rewritePackage) wrapBasicLitWithTemplatedT(basicLit *ast.BasicLit, args []ast.Expr, callExpr *ast.CallExpr, argIndex int) ast.Expr {
	valueWithoutQuotes, _ := strconv.Unquote(basicLit.Value) //basicLit.Value[1 : len(basicLit.Value)-1]

	_, ok := rp.ExtractedStrings[valueWithoutQuotes]
	if !ok && rp.ExtractedStrings != nil {
		return callExpr
	}

	rp.TotalStrings++
	tIdent := &ast.Ident{Name: "T"}
	argNames := common.GetTemplatedStringArgs(valueWithoutQuotes)

	compositeExpr := []ast.Expr{}
	processedArgsMap := make(map[string]bool)

	for i, argName := range argNames {
		if callExpr, ok := args[argIndex+i+1].(*ast.CallExpr); ok {
			rp.callExprTFunc(callExpr)
		} else if basicLit, ok := args[argIndex+i+1].(*ast.BasicLit); ok {
			args[argIndex+i] = rp.wrapBasicLitWithT(basicLit)
		}

		if processedArgsMap[argName] != true {
			quotedArgName := "\"" + argName + "\""
			basicLit.ValuePos = 0

			valueExpr := args[argIndex+i+1]
			if basicLit, ok := args[argIndex+i+1].(*ast.BasicLit); ok {
				valueExpr = rp.wrapBasicLitWithT(basicLit)
			}

			keyValueExpr := &ast.KeyValueExpr{Key: &ast.BasicLit{Kind: 9, Value: quotedArgName}, Value: valueExpr}
			processedArgsMap[argName] = true
			compositeExpr = append(compositeExpr, keyValueExpr)
		}
	}

	mapInterfaceType := &ast.InterfaceType{Interface: 142, Methods: &ast.FieldList{List: nil, Opening: 1, Closing: 2}, Incomplete: false}
	mapType := &ast.MapType{Map: 131, Key: &ast.Ident{Name: "string"}, Value: mapInterfaceType}
	compositeLit := &ast.CompositeLit{Type: mapType, Elts: compositeExpr}

	return &ast.CallExpr{Fun: tIdent, Args: []ast.Expr{basicLit, compositeLit}}
}

func (rp *rewritePackage) wrapBasicLitWithT(basicLit *ast.BasicLit) ast.Expr {
	if basicLit.Kind != token.STRING {
		return basicLit
	}

	valueWithoutQuotes, _ := strconv.Unquote(basicLit.Value) //basicLit.Value[1 : len(basicLit.Value)-1]
	_, ok := rp.ExtractedStrings[valueWithoutQuotes]
	if !ok && rp.ExtractedStrings != nil {
		return basicLit
	}

	rp.TotalStrings++
	tIdent := &ast.Ident{Name: "T"}
	return &ast.CallExpr{Fun: tIdent, Args: []ast.Expr{basicLit}}
}

func (rp *rewritePackage) addInitFuncToPackage(packageName, outputDir, importPath string) error {
	rp.Println("i18n4go: adding init func to package:", packageName, " to output dir:", outputDir)

	common.CreateOutputDirsIfNeeded(outputDir)

	pieces := strings.Split(importPath, "/")
	for index, str := range pieces {
		pieces[index] = `"` + str + `"`
	}

	joinedImportPath := "filepath.Join(" + strings.Join(pieces, ", ") + ")"
	content := rp.getInitFuncCodeSnippetContent(packageName, joinedImportPath)

	return ioutil.WriteFile(filepath.Join(outputDir, "i18n_init.go"), []byte(content), 0666)
}

func (rp *rewritePackage) getInitFuncCodeSnippetContent(packageName, importPath string) string {
	snippetContent := INIT_CODE_SNIPPET
	if rp.InitCodeSnippetFilename != "" {
		bytes, err := ioutil.ReadFile(rp.InitCodeSnippetFilename)
		if err != nil {
			rp.Printf("i18n4go: error reading content of init code snippet file: %s\n, using default", rp.InitCodeSnippetFilename)
		} else {
			snippetContent = string(bytes)
		}
	}

	content := strings.Replace(snippetContent, "__PACKAGE__NAME__", packageName, -1)
	content = strings.Replace(content, "__FULL_IMPORT_PATH__", importPath, -1)
	return content
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

func (rp *rewritePackage) updateExtractedStrings(i18nStringInfo common.I18nStringInfo, templatedString string) {
	oldID := i18nStringInfo.ID

	i18nStringInfo.ID = templatedString
	i18nStringInfo.Translation = templatedString

	rp.ExtractedStrings[templatedString] = i18nStringInfo
	rp.UpdatedExtractedStrings[templatedString] = i18nStringInfo
	delete(rp.UpdatedExtractedStrings, oldID)

	rp.SaveExtractedStrings = true
}
