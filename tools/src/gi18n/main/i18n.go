package main

import (
	"fmt"
	"flag"
	"os"
	"strings"
	"runtime/debug"

	"gi18n/cmds"
	common "gi18n/common"
)

var options common.Options

func main() {
	defer handlePanic()

	if options.ExtractStringsCmdFlag {
		extractStringsCmd()
	} else {
		usage()
		return
	}
}

func extractStringsCmd() {
	if options.FilenameFlag == "" && options.DirnameFlag == "" {
		usage()
		return
	}

	extractStrings := extract_strings.NewExtractStrings(options)

	if options.FilenameFlag != "" {
		extractStrings.InspectFile(options.FilenameFlag)
	} else {
		extractStrings.InspectDir(options.DirnameFlag, options.RecurseFlag)
		extractStrings.Println()
		extractStrings.Println("Total files parsed:", extractStrings.TotalFiles)
		extractStrings.Println("Total extracted strings:", extractStrings.TotalStrings)
	}	
}

func init() {
	flag.BoolVar(&options.ExtractStringsCmdFlag, "extract_strings", true, "want to extract strings from file or directory")
	flag.BoolVar(&options.VerboseFlag, "v", false, "verbose mode where lots of output is generated during execution")
	flag.BoolVar(&options.PoFlag, "po", true, "generate standard .po file for translation")
	flag.StringVar(&options.ExcludedFilenameFlag, "e", "excluded.json", "the excluded JSON file name, all strings there will be excluded")
	flag.StringVar(&options.FilenameFlag, "f", "", "the file name for which strings are extracted")
	flag.StringVar(&options.DirnameFlag, "d", "", "the dir name for which all .go files will have their strings extracted")
	flag.BoolVar(&options.RecurseFlag, "r", false, "recursesively extract strings from all files in the same directory as filename or dirName")

	flag.Parse()
}

func usage() {
	usageString := `
gi18n -extract_strings [-po] -f <fileName> | [-d <dirName> | -r -d <dirName>]
	-po to generate standard .po files for translation
	-e the JSON file with strings to be excluded, defaults to excluded.json if present
	-f the go file name to extract strings
	-d the directory containing the go files to extract strings
	-r recursesively extract strings from all subdirectories
	`
	fmt.Println(usageString)
}

func handlePanic() {
	err := recover()
	if err != nil {
		switch err := err.(type) {
		case error:
			displayCrashDialog(err.Error())
		case string:
			displayCrashDialog(err)
		default:
			displayCrashDialog("An unexpected type of error")
		}
	}

	if err != nil {
		os.Exit(1)
	}
}

func displayCrashDialog(errorMessage string) {
	formattedString := `
Something completely unexpected happened. This is a bug in %s.
Please file this bug : https://github.com/maximilien/gi18n/issues
Tell us that you ran this command:

	%s

this error occurred:

	%s

and this stack trace:

%s
	`

	stackTrace := "\t" + strings.Replace(string(debug.Stack()), "\n", "\n\t", -1)
	println(fmt.Sprintf(formattedString, "gi18n", strings.Join(os.Args, " "), errorMessage, stackTrace))
}
