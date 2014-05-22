package common

import (
	"os"

	"io/ioutil"
	"strconv"
	"strings"

	"path/filepath"

	"encoding/json"

	"github.com/maximilien/i18n4cf/cmds"
)

func ParseStringList(stringList string, delimiter string) []string {
	stringArray := strings.Split(stringList, delimiter)
	parsedStrings := make([]string, len(stringArray))
	for i, aString := range stringArray {
		parsedStrings[i] = strings.Trim(aString, " ")
	}
	return parsedStrings
}

func CopyFileContents(src, dst string) error {
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
	cmd.Println("fileName:", fileName)
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
