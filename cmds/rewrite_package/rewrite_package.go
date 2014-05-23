package rewrite_package

import (
	"fmt"

	"github.com/maximilien/i18n4cf/cmds"
	"github.com/maximilien/i18n4cf/common"
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

func (ct *rewritePackage) Options() cmds.Options {
	return ct.options
}

func (ct *rewritePackage) Println(a ...interface{}) (int, error) {
	if ct.options.VerboseFlag {
		return fmt.Println(a...)
	}

	return 0, nil
}

func (ct *rewritePackage) Printf(msg string, a ...interface{}) (int, error) {
	if ct.options.VerboseFlag {
		return fmt.Printf(msg, a...)
	}

	return 0, nil
}

func (ct *rewritePackage) Run() error {
	return nil
}
