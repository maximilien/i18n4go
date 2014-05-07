package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"runtime/debug"

	"gi18n/cmds/extract_strings"
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
	if options.HelpFlag || (options.FilenameFlag == "" && options.DirnameFlag == "") {
		usage()
		return
	}

	extractStrings := extract_strings.NewExtractStrings(options)

	startTime := time.Now()
	if options.FilenameFlag != "" {
		extractStrings.InspectFile(options.FilenameFlag)
	} else {
		extractStrings.InspectDir(options.DirnameFlag, options.RecurseFlag)
		extractStrings.Println()
		extractStrings.Println("Total files parsed:", extractStrings.TotalFiles)
		extractStrings.Println("Total extracted strings:", extractStrings.TotalStrings)
	}
	duration := time.Now().Sub(startTime)
	extractStrings.Println("Total time:", duration)
}

func init() {
	flag.BoolVar(&options.HelpFlag, "h", false, "prints the usage")

	flag.BoolVar(&options.ExtractStringsCmdFlag, "extract-strings", true, "want to extract strings from file or directory")

	flag.BoolVar(&options.VerboseFlag, "v", false, "verbose mode where lots of output is generated during execution")
	flag.BoolVar(&options.PoFlag, "p", true, "generate standard .po file for translation")

	flag.StringVar(&options.ExcludedFilenameFlag, "e", "excluded.json", "the excluded JSON file name, all strings there will be excluded")

	flag.StringVar(&options.OutputDirFlag, "o", "", "output directory where the translation files will be placed")
	flag.BoolVar(&options.OutputMatchImportFlag, "output-match-import", false, "generated files are created in directory to match the import structure")
	flag.BoolVar(&options.OutputMatchPackageFlag, "output-match-package", false, "generated files are created in directory to match the package name")
	flag.BoolVar(&options.OutputFlatFlag, "output-flat", true, "generated files are created in the specified output directory")

	flag.StringVar(&options.FilenameFlag, "f", "", "the file name for which strings are extracted")
	flag.StringVar(&options.DirnameFlag, "d", "", "the dir name for which all .go files will have their strings extracted")
	flag.BoolVar(&options.RecurseFlag, "r", false, "recursesively extract strings from all files in the same directory as filename or dirName")

	flag.StringVar(&options.IgnoreRegexp, "ignore-regexp", "", "a perl-style regular expression for files to ignore, e.g., \".*test.*\"")

	flag.Parse()
}

func usage() {
	usageString := `
gi18n -extract-strings [-vpe] [-o <outputDir>] -f <fileName> | -d [-r] [-ignore-regexp <regex>] <dirName>
        -h   	                prints the usage

        -v                      verbose
	-p                      to generate standard .po files for translation

	-o                      the output directory where the translation files will be placed
	-output-match-import    generated files are created in directory to match the import structure
	-output-match-package   generated files are created in directory to match the package name
	-output-flat            generated files are created in the specified output directory

	-e                      the JSON file with strings to be excluded, defaults to excluded.json if present

	-f                      the go file name to extract strings

	-r                      recursesively extract strings from all subdirectories
	-d                      the directory containing the go files to extract strings

	-ignore-regexp          a perl-style regular expression for files to ignore, e.g., ".*test.*"
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
