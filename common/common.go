package common

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"io/ioutil"
	"strconv"
	"strings"

	"path/filepath"

	"encoding/json"

	"github.com/maximilien/i18n4cf/cmds"
)

const (
	TEMPLATED_STRING_REGEXP    = `\{\{\.[[:alnum:][:punct:][:print:]]+?\}\}`
	INTERPOLATED_STRING_REGEXP = `%(?:[#v]|[%EGUTXbcdefgopqstvx])`
)

var templatedStringRegexp, interpolatedStringRegexp *regexp.Regexp

func ParseStringList(stringList string, delimiter string) []string {
	stringArray := strings.Split(stringList, delimiter)
	var parsedStrings []string
	for _, aString := range stringArray {
		if aString != "" {
			parsedStrings = append(parsedStrings, strings.Trim(strings.Trim(aString, " "), "\""))
		}
	}
	return parsedStrings
}

func CreateTmpFile(content string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	ioutil.WriteFile(tmpFile.Name(), []byte(content), 666)

	return tmpFile, nil
}

func CheckFile(fileName string) (string, string, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", "", err
	}

	if !fileInfo.Mode().IsRegular() {
		return "", "", fmt.Errorf("Non-regular source file %s (%s)", fileInfo.Name(), fileInfo.Mode().String())
	}

	return filepath.Base(fileName), filepath.Dir(fileName), nil
}

func CopyFileContents(src, dst string) error {
	err := CreateOutputDirsIfNeeded(filepath.Dir(dst))
	if err != nil {
		return err
	}

	byteArray, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dst, byteArray, 0644)
}

func GetAbsFileInfo(fileNamePath string) (os.FileInfo, error) {
	var absFilePath = fileNamePath
	if !filepath.IsAbs(absFilePath) {
		absFilePath = filepath.Join(os.Getenv("PWD"), absFilePath)
	}

	file, err := os.OpenFile(absFilePath, os.O_RDONLY, 0)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return file.Stat()
}

func FindFilePath(filename string) (string, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return "", err
	}
	path := filename[0 : len(filename)-len(fileInfo.Name())]
	return path, nil
}

func CreateOutputDirsIfNeeded(outputDirname string) error {
	_, err := os.Stat(outputDirname)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputDirname, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveStrings(cmd cmds.CommandInterface, stringInfos map[string]StringInfo, outputDirname string, fileName string) error {
	if len(stringInfos) != 0 {
		cmd.Println("Saving extracted i18n strings to file:", fileName)
	}

	CreateOutputDirsIfNeeded(outputDirname)

	i18nStringInfos := make([]I18nStringInfo, len(stringInfos))
	i := 0
	for _, stringInfo := range stringInfos {
		i18nStringInfos[i] = I18nStringInfo{ID: stringInfo.Value, Translation: stringInfo.Value}
		i++
	}

	jsonData, err := json.MarshalIndent(i18nStringInfos, "", "   ")
	if err != nil {
		cmd.Println(err)
		return err
	}

	if !cmd.Options().DryRunFlag && len(i18nStringInfos) != 0 {
		file, err := os.Create(filepath.Join(outputDirname, fileName[strings.LastIndex(fileName, string(os.PathSeparator))+1:len(fileName)]))
		defer file.Close()
		if err != nil {
			cmd.Println(err)
			return err
		}

		file.Write(jsonData)
	}

	return nil
}

func SaveStringsInPo(cmd cmds.CommandInterface, stringInfos map[string]StringInfo, outputDirname string, fileName string) error {
	if len(stringInfos) != 0 {
		cmd.Println("Creating and saving i18n strings to .po file:", fileName)
	}

	if !cmd.Options().DryRunFlag && len(stringInfos) != 0 {
		CreateOutputDirsIfNeeded(outputDirname)
		file, err := os.Create(filepath.Join(outputDirname, fileName[strings.LastIndex(fileName, string(os.PathSeparator))+1:len(fileName)]))
		defer file.Close()
		if err != nil {
			cmd.Println(err)
			return err
		}

		for _, stringInfo := range stringInfos {
			file.Write([]byte("# filename: " + strings.Split(fileName, ".en.po")[0] +
				", offset: " + strconv.Itoa(stringInfo.Offset) +
				", line: " + strconv.Itoa(stringInfo.Line) +
				", column: " + strconv.Itoa(stringInfo.Column) + "\n"))
			file.Write([]byte("msgid " + strconv.Quote(stringInfo.Value) + "\n"))
			file.Write([]byte("msgstr " + strconv.Quote(stringInfo.Value) + "\n"))
			file.Write([]byte("\n"))
		}
	}
	return nil
}

func SaveI18nStringsInPo(cmd cmds.CommandInterface, i18nStrings []I18nStringInfo, fileName string) error {
	cmd.Println("gi18n: creating and saving i18n strings to .po file:", fileName)

	if !cmd.Options().DryRunFlag && len(i18nStrings) != 0 {
		file, err := os.Create(fileName)
		defer file.Close()
		if err != nil {
			cmd.Println(err)
			return err
		}

		for _, stringInfo := range i18nStrings {
			file.Write([]byte("msgid " + strconv.Quote(stringInfo.ID) + "\n"))
			file.Write([]byte("msgstr " + strconv.Quote(stringInfo.Translation) + "\n"))
			file.Write([]byte("\n"))
		}
	}
	return nil
}

func SaveI18nStringInfos(cmd cmds.CommandInterface, i18nStringInfos []I18nStringInfo, fileName string) error {
	jsonData, err := json.MarshalIndent(i18nStringInfos, "", "   ")
	if err != nil {
		cmd.Println(err)
		return err
	}

	if !cmd.Options().DryRunFlag && len(i18nStringInfos) != 0 {
		err := ioutil.WriteFile(fileName, jsonData, 0700)
		if err != nil {
			cmd.Println(err)
			return err
		}
	}

	return nil
}

func LoadI18nStringInfos(fileName string) ([]I18nStringInfo, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var i18nStringInfos []I18nStringInfo
	err = json.Unmarshal(content, &i18nStringInfos)
	if err != nil {
		return nil, err
	}

	return i18nStringInfos, nil
}

func CreateI18nStringInfoMap(i18nStringInfos []I18nStringInfo) (map[string]I18nStringInfo, error) {
	inputMap := make(map[string]I18nStringInfo, len(i18nStringInfos))

	for _, i18nStringInfo := range i18nStringInfos {

		if _, ok := inputMap[i18nStringInfo.ID]; !ok {
			inputMap[i18nStringInfo.ID] = i18nStringInfo
		} else {
			return nil, errors.New("Duplicated key found: " + i18nStringInfo.ID)
		}

	}

	return inputMap, nil
}

func GetTemplatedStringArgs(aString string) []string {
	re, err := getTemplatedStringRegexp()
	if err != nil {
		fmt.Errorf("gi18n: Error compiling templated string Regexp: %s", err.Error())
		return []string{}
	}

	matches := re.FindAllStringSubmatch(aString, -1)
	var stringMatches []string
	for _, match := range matches {
		stringMatch := match[0]
		stringMatch = stringMatch[3 : len(stringMatch)-2]
		stringMatches = append(stringMatches, stringMatch)
	}

	return stringMatches
}

func IsTemplatedString(aString string) bool {
	re, err := getTemplatedStringRegexp()
	if err != nil {
		fmt.Errorf("gi18n: Error compiling templated string Regexp: %s", err.Error())
		return false
	}

	return re.Match([]byte(aString))
}

func IsInterpolatedString(aString string) bool {
	re, err := getInterpolatedStringRegexp()
	if err != nil {
		fmt.Errorf("gi18n: Error compiling interpolated string Regexp: %s", err.Error())
		return false
	}

	return re.Match([]byte(aString))
}

func ConvertToTemplatedString(aString string) string {
	if !IsInterpolatedString(aString) {
		return aString
	}

	re, err := getInterpolatedStringRegexp()
	if err != nil {
		fmt.Errorf("gi18n: Error compiling interpolated string Regexp: %s", err.Error())
		return ""
	}

	matches := re.FindAllStringSubmatch(aString, -1)
	templatedString := aString
	for i, match := range matches {
		argName := "{{.Arg" + strconv.Itoa(i) + "}}"
		templatedString = strings.Replace(templatedString, match[0], argName, 1)
	}

	return templatedString
}

// Private

func getTemplatedStringRegexp() (*regexp.Regexp, error) {
	var err error
	if templatedStringRegexp == nil {
		templatedStringRegexp, err = regexp.Compile(TEMPLATED_STRING_REGEXP)
	}

	return templatedStringRegexp, err
}

func getInterpolatedStringRegexp() (*regexp.Regexp, error) {
	var err error
	if interpolatedStringRegexp == nil {
		interpolatedStringRegexp, err = regexp.Compile(INTERPOLATED_STRING_REGEXP)
	}

	return interpolatedStringRegexp, err
}
