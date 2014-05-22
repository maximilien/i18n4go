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

func (vs *VerifyStrings) Run() error {
	//Check input file
	//Get list of target files
	//verify() each target file

	return nil
}

func (vs *VerifyStrings) verify(inputFilename string, targetFilename string) error {
	//Check target file
	//Get list of StringInfo from input file as map
	//Get list of StringInfo from target file as array
	//for each stringInfo in target file StringInfo list array
	//  Check that stringInfo.ID is in input file StringInfo list
	//  Remove stringInfo.ID from input fileStringInfo list
	//Check if the input file StringInfo map is empty or not

	return nil
}
