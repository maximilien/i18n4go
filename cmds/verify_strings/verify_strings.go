package verify_strings

import (
	"fmt"

	"github.com/maximilien/i18n4cf/cmds"
	"github.com/maximilien/i18n4cf/common"
)

type VerifyStrings struct {
	options cmds.Options

	LanguageFilenames []string
	Languages         []string
}

func NewVerifyStrings(options cmds.Options) VerifyStrings {
	languageFilenames := common.ParseStringList(options.LanguageFilesFlag, ",")
	languages := common.ParseStringList(options.LanguagesFlag, ",")

	return VerifyStrings{options: options,
		LanguageFilenames: languageFilenames,
		Languages:         languages,
	}
}

func (vs *VerifyStrings) Options() cmds.Options {
	return vs.options
}

func (vs *VerifyStrings) Println(a ...interface{}) (int, error) {
	if vs.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (vs *VerifyStrings) Printf(msg string, a ...interface{}) (int, error) {
	if vs.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}
