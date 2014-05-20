package create_translations

import (
	"fmt"
	"os"
	"strings"

	"path/filepath"

	common "github.com/maximilien/i18n4cf/common"
)

type CreateTranslations struct {
	Options common.Options

	Filename       string
	OutputDirname  string
	SourceLanguage string

	Languages []string

	ExtractedStrings map[string]common.StringInfo

	TotalStrings int
	TotalFiles   int
}

func NewCreateTranslations(options common.Options) CreateTranslations {
	languages := common.ParseLanguages(options.LanguagesFlag)

	return CreateTranslations{Options: options,
		Filename:       options.FilenameFlag,
		OutputDirname:  options.OutputDirFlag,
		SourceLanguage: options.SourceLanguageFlag,
		Languages:      languages,
		TotalStrings:   0,
		TotalFiles:     0}
}

func (ct *CreateTranslations) Println(a ...interface{}) (int, error) {
	if ct.Options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ct *CreateTranslations) Printf(msg string, a ...interface{}) (int, error) {
	if ct.Options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (ct *CreateTranslations) CreateTranslationFiles(sourceFilename string) error {
	ct.Println("gi18n: creating translation files for:", sourceFilename)
	ct.Filename = sourceFilename

	for _, language := range ct.Languages {
		ct.Println("gi18n: creating translation file copy for language:", language)
		err := ct.createTranslationFile(sourceFilename, language)
		if err != nil {
			return fmt.Errorf("gi18n: could not create translation file for language: ", language)
		}
		ct.Println()
	}

	return nil
}

func (ct *CreateTranslations) checkFile(fileName string) (string, string, error) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return "", "", err
	}

	if !fileInfo.Mode().IsRegular() {
		return "", "", fmt.Errorf("gi18n: non-regular source file %s (%s)", fileInfo.Name(), fileInfo.Mode().String())
	}

	return fileInfo.Name(), fileName[:len(fileName)-len(fileInfo.Name())-1], nil
}

func (ct *CreateTranslations) createTranslationFile(sourceFilename string, language string) error {
	fileName, _, err := ct.checkFile(sourceFilename)
	if err != nil {
		return err
	}

	destFilename := filepath.Join(ct.OutputDirname, strings.Replace(fileName, ct.Options.SourceLanguageFlag, language, -1))
	ct.Println("gi18n: creating translation file:", destFilename)

	return common.CopyFileContents(sourceFilename, destFilename)
}
