package create_translations

import (
	"fmt"
	// "os"
	// "path/filepath"

	// "encoding/json"
	// "io/ioutil"

	common "github.com/maximilien/i18n4cf/common"
)

type CreateTranslations struct {
	Options common.Options

	Filename      string
	OutputDirname string

	ExtractedStrings map[string]common.StringInfo

	TotalStrings int
	TotalFiles   int
}

func NewCreateTranslations(options common.Options) CreateTranslations {
	fmt.Println("New CreateTranslations with options:", options)

	return CreateTranslations{Options: options,
		Filename:      options.FilenameFlag,
		OutputDirname: options.OutputDirFlag,
		TotalStrings:  0,
		TotalFiles:    0}
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
