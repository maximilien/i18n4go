package extract_strings

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"go/ast"
	"go/build"
	"go/parser"
	"go/token"

	"encoding/json"
	"io/ioutil"

	common "github.com/maximilien/i18n4cf/common"
)

type ExtractStrings struct {
	Options common.Options

	Filename      string
	I18nFilename  string
	PoFilename    string
	OutputDirname string

	ExtractedStrings map[string]common.StringInfo
	FilteredStrings  map[string]string
	FilteredRegexps  []*regexp.Regexp

	TotalStringsDir int
	TotalStrings    int
	TotalFiles      int

	IgnoreRegexp *regexp.Regexp
}

func NewExtractStrings(options common.Options) ExtractStrings {
	var compiledRegexp *regexp.Regexp
	if options.IgnoreRegexp != "" {
		compiledReg, err := regexp.Compile(options.IgnoreRegexp)
		if err != nil {
			fmt.Println("WARNING compiling ignore-regexp:", err)
		}
		compiledRegexp = compiledReg
	}

	return ExtractStrings{Options: options,
		Filename:         "extracted_strings.json",
		OutputDirname:    options.OutputDirFlag,
		ExtractedStrings: nil,
		FilteredStrings:  nil,
		FilteredRegexps:  nil,
		TotalStringsDir:  0,
		TotalStrings:     0,
		TotalFiles:       0,
		IgnoreRegexp:     compiledRegexp}
}

func (es *ExtractStrings) Println(a ...interface{}) (int, error) {
	if es.Options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (es *ExtractStrings) Printf(msg string, a ...interface{}) (int, error) {
	if es.Options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (es *ExtractStrings) InspectFile(filename string) error {
	es.Println("gi18n: extracting strings from file:", filename)

	es.ExtractedStrings = make(map[string]common.StringInfo)
	es.FilteredStrings = make(map[string]string)
	es.FilteredRegexps = []*regexp.Regexp{}

	es.setFilename(filename)
	es.setI18nFilename(filename)
	es.setPoFilename(filename)

	fset := token.NewFileSet()

	fileInfo, err := getFileName(filename)
	if err != nil {
		es.Println(err)
	}

	if strings.HasPrefix(fileInfo.Name(), ".") {
		es.Println("WARNING ignoring file:", filename)
		return nil
	}

	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		es.Println(err)
		return err
	}

	err = es.loadExcludedStrings()
	if err != nil {
		es.Println(err)
		return err
	}
	es.Println(fmt.Sprintf("Loaded %d excluded strings", len(es.FilteredStrings)))

	err = es.loadExcludedRegexps()
	if err != nil {
		es.Println(err)
		return err
	}
	es.Println(fmt.Sprintf("Loaded %d excluded regexps", len(es.FilteredRegexps)))

	es.excludeImports(astFile)

	es.extractString(astFile, fset)
	es.TotalStringsDir += len(es.ExtractedStrings)
	es.TotalStrings += len(es.ExtractedStrings)
	es.TotalFiles += 1

	es.Printf("Extracted %d strings from file: %s\n", len(es.ExtractedStrings), filename)

	var outputDirname = es.OutputDirname
	if es.Options.OutputDirFlag != "" {
		if es.Options.OutputMatchImportFlag {
			outputDirname, err = es.findImportPath(filename)
			if err != nil {
				es.Println(err)
				return err
			}
		} else if es.Options.OutputMatchPackageFlag {
			outputDirname, err = es.findPackagePath(filename)
			if err != nil {
				es.Println(err)
				return err
			}
		}
	} else {
		outputDirname, err = es.findFilePath(filename)
		if err != nil {
			es.Println(err)
			return err
		}
	}

	err = es.saveExtractedStrings(outputDirname)
	if err != nil {
		es.Println(err)
		return err
	}

	err = es.saveI18nStrings(outputDirname)
	if err != nil {
		es.Println(err)
		return err
	}

	if es.Options.PoFlag {
		err = es.saveI18nStringsInPo(outputDirname)
		if err != nil {
			es.Println(err)
			return err
		}
	}

	return nil
}

func (es *ExtractStrings) InspectDir(dirName string, recursive bool) error {
	es.Printf("gi18n: inspecting dir %s, recursive: %t\n", dirName, recursive)
	es.Println()

	fset := token.NewFileSet()
	es.TotalStringsDir = 0

	packages, err := parser.ParseDir(fset, dirName, nil, parser.ParseComments)
	if err != nil {
		es.Println(err)
		return err
	}

	for k, pkg := range packages {
		es.Println("Extracting strings in package:", k)
		for fileName, _ := range pkg.Files {
			if es.IgnoreRegexp != nil && es.IgnoreRegexp.MatchString(fileName) {
				es.Println("Using ignore-regexp:", es.Options.IgnoreRegexp)
				continue
			} else {
				es.Println("No match for ignore-regexp:", es.Options.IgnoreRegexp)
			}

			if strings.HasSuffix(fileName, ".go") {
				err = es.InspectFile(fileName)
				if err != nil {
					es.Println(err)
				}
			}
		}
	}
	es.Printf("Extracted total of %d strings\n\n", es.TotalStringsDir)

	if recursive {
		fileInfos, _ := ioutil.ReadDir(dirName)
		for _, fileInfo := range fileInfos {
			if fileInfo.IsDir() && !strings.HasPrefix(fileInfo.Name(), ".") {
				err = es.InspectDir(dirName+"/"+fileInfo.Name(), recursive)
				if err != nil {
					es.Println(err)
				}
			}
		}
	}

	return nil
}

func getFileName(filePath string) (os.FileInfo, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	defer file.Close()
	if err != nil {
		fmt.Println("ERROR opening file", err)
		return nil, err
	}

	return file.Stat()
}

func (es *ExtractStrings) createOutputDirsIfNeeded(outputDirname string) error {
	_, err := os.Stat(outputDirname)
	if os.IsNotExist(err) {
		es.Println("Creating output directory:", outputDirname)
		err = os.MkdirAll(outputDirname, 0777)
		if err != nil {
			fmt.Println("ERROR opening output directory", err)
			return err
		}
	}
	return nil
}

func (es *ExtractStrings) findFilePath(filename string) (string, error) {
	path := es.OutputDirname
	fileInfo, err := os.Stat(filename)
	if err != nil {
		fmt.Println("ERROR opening file", err)
		return "", err
	}
	path = filename[0 : len(filename)-len(fileInfo.Name())]
	return path, nil
}

func (es *ExtractStrings) findImportPath(filename string) (string, error) {
	path := es.OutputDirname
	filePath, err := es.findFilePath(filename)
	if err != nil {
		fmt.Println("ERROR opening file", err)
		return "", err
	}

	pkg, err := build.ImportDir(filePath, 0)
	if strings.HasPrefix(pkg.Dir, "src/") {
		path = path + "/" + pkg.Dir[len("src/"):len(pkg.Dir)]
	}

	return path, nil
}

func (es *ExtractStrings) findPackagePath(filename string) (string, error) {
	path := es.OutputDirname

	filePath, err := es.findFilePath(filename)
	if err != nil {
		fmt.Println("ERROR opening file", err)
		return "", err
	}

	pkg, err := build.ImportDir(filePath, 0)
	if err != nil {
		fmt.Println("ERROR opening file", err)
		return "", err
	}

	return path + "/" + pkg.Name, nil
}

func (es *ExtractStrings) saveExtractedStrings(outputDirname string) error {
	es.Println("Saving extracted strings to file:", es.Filename)
	es.createOutputDirsIfNeeded(outputDirname)

	stringInfos := make([]common.StringInfo, 0)
	for _, stringInfo := range es.ExtractedStrings {
		stringInfos = append(stringInfos, stringInfo)
	}

	jsonData, err := json.Marshal(stringInfos)
	if err != nil {
		es.Println(err)
		return err
	}

	file, err := os.Create(outputDirname + "/" + es.Filename[strings.LastIndex(es.Filename, "/")+1:len(es.Filename)])
	defer file.Close()
	if err != nil {
		es.Println(err)
		return err
	}

	file.Write(jsonData)

	return nil
}

func (es *ExtractStrings) saveI18nStrings(outputDirname string) error {
	es.Println("Saving extracted i18n strings to file:", es.I18nFilename)
	es.createOutputDirsIfNeeded(outputDirname)

	i18nStringInfos := make([]common.I18nStringInfo, len(es.ExtractedStrings))
	i := 0
	for _, stringInfo := range es.ExtractedStrings {
		i18nStringInfos[i] = common.I18nStringInfo{ID: stringInfo.Value, Translation: stringInfo.Value}
		i++
	}

	jsonData, err := json.Marshal(i18nStringInfos)
	if err != nil {
		es.Println(err)
		return err
	}

	file, err := os.Create(outputDirname + "/" + es.I18nFilename[strings.LastIndex(es.I18nFilename, "/")+1:len(es.I18nFilename)])
	if err != nil {
		es.Println(err)
		return err
	}

	file.Write(jsonData)
	defer file.Close()

	return nil
}

func (es *ExtractStrings) saveI18nStringsInPo(outputDirname string) error {
	es.Println("Creating and saving i18n strings to .po file:", es.PoFilename)
	es.createOutputDirsIfNeeded(outputDirname)

	file, err := os.Create(outputDirname + "/" + es.PoFilename[strings.LastIndex(es.PoFilename, "/")+1:len(es.PoFilename)])
	if err != nil {
		es.Println(err)
		return err
	}

	for _, stringInfo := range es.ExtractedStrings {
		file.Write([]byte("# filename: " + stringInfo.Filename +
			", offset: " + strconv.Itoa(stringInfo.Offset) +
			", line: " + strconv.Itoa(stringInfo.Line) +
			", column: " + strconv.Itoa(stringInfo.Column) + "\n"))
		file.Write([]byte("msgid " + strconv.Quote(stringInfo.Value) + "\n"))
		file.Write([]byte("msgstr " + strconv.Quote(stringInfo.Value) + "\n"))
		file.Write([]byte("\n"))
	}

	defer file.Close()

	return nil
}

func (es *ExtractStrings) setFilename(filename string) {
	es.Filename = filename + ".extracted.json"
}

func (es *ExtractStrings) setI18nFilename(filename string) {
	es.I18nFilename = filename + ".en.json"
}

func (es *ExtractStrings) setPoFilename(filename string) {
	es.PoFilename = filename + ".en.po"
}

func (es *ExtractStrings) loadExcludedStrings() error {
	es.Println("Excluding strings in file:", es.Options.ExcludedFilenameFlag)

	content, err := ioutil.ReadFile(es.Options.ExcludedFilenameFlag)
	if err != nil {
		fmt.Print(err)
		return err
	}

	var excludedStrings common.ExcludedStrings
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

func (es *ExtractStrings) loadExcludedRegexps() error {
	es.Println("Excluding regexps in file:", es.Options.ExcludedFilenameFlag)

	content, err := ioutil.ReadFile(es.Options.ExcludedFilenameFlag)
	if err != nil {
		fmt.Print(err)
		return err
	}

	var excludedRegexps common.ExcludedStrings
	err = json.Unmarshal(content, &excludedRegexps)
	if err != nil {
		fmt.Print(err)
		return err
	}

	for _, regexpString := range excludedRegexps.ExcludedRegexps {
		compiledRegexp, err := regexp.Compile(regexpString)
		if err != nil {
			fmt.Println("WARNING error compiling regexp:", regexpString)
		}

		es.FilteredRegexps = append(es.FilteredRegexps, compiledRegexp)
	}

	return nil
}

func (es *ExtractStrings) extractString(f *ast.File, fset *token.FileSet) error {
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s, _ = strconv.Unquote(x.Value)
			if len(s) > 0 && x.Kind == token.STRING && s != "\t" && s != "\n" && s != " " && !es.filter(s) { //TODO: fix to remove these: s != "\\t" && s != "\\n" && s != " "
				position := fset.Position(n.Pos())
				stringInfo := common.StringInfo{Value: s,
					Filename: position.Filename,
					Offset:   position.Offset,
					Line:     position.Line,
					Column:   position.Column}
				es.ExtractedStrings[s] = stringInfo
			}
		}
		return true
	})

	return nil
}

func (es *ExtractStrings) excludeImports(astFile *ast.File) {
	for i := range astFile.Imports {
		importString, _ := strconv.Unquote(astFile.Imports[i].Path.Value)
		es.FilteredStrings[importString] = importString
	}

}

func (es *ExtractStrings) filter(aString string) bool {
	for i := range common.BLANKS {
		if aString == common.BLANKS[i] {
			return true
		}
	}

	if es.FilteredStrings[aString] != "" {
		return true
	}

	for _, compiledRegexp := range es.FilteredRegexps {
		if compiledRegexp.MatchString(aString) {
			return true
		}
	}

	return false
}
