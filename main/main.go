package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"runtime/debug"

	"github.com/maximilien/i18n4cf/cmds/create_translations"
	"github.com/maximilien/i18n4cf/cmds/extract_strings"

	common "github.com/maximilien/i18n4cf/common"
)

var options common.Options

func main() {
	defer handlePanic()

	if options.ExtractStringsCmdFlag {
		extractStringsCmd()
	} else if options.CreateTranslationsCmdFlag {
		createTranslationsCmd()
	} else {
		usage()
		return
	}
}

func createTranslationsCmd() {
	if options.HelpFlag || (options.FilenameFlag == "") {
		usage()
		return
	}

	createTranslations := create_translations.NewCreateTranslations(options)

	startTime := time.Now()

	err := createTranslations.CreateTranslationFiles(options.FilenameFlag)
	if err != nil {
		createTranslations.Println("gi18n: Could not create translation files, err:", err)
	}

	duration := time.Now().Sub(startTime)
	createTranslations.Println("Total time:", duration)
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

	flag.BoolVar(&options.ExtractStringsCmdFlag, "extract-strings", false, "want to extract strings from file or directory")
	flag.BoolVar(&options.CreateTranslationsCmdFlag, "create-translations", false, "create translation files for different languages using a source file")

	flag.StringVar(&options.SourceLanguageFlag, "source-language", "en", "the source language of the file, typically also part of the file name, e.g., \"en_US\"")
	flag.StringVar(&options.LanguagesFlag, "languages", "", "a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"")
	flag.StringVar(&options.GoogleTranslateApiKeyFlag, "google-translate-api-key", "", "your public Google Translate API key which is used to generate translations (charge is applicable)")

	flag.BoolVar(&options.VerboseFlag, "v", false, "verbose mode where lots of output is generated during execution")
	flag.BoolVar(&options.PoFlag, "p", true, "generate standard .po file for translation")
	flag.BoolVar(&options.DryRunFlag, "dry-run", false, "prevents any output files from being created")

	flag.StringVar(&options.ExcludedFilenameFlag, "e", "excluded.json", "the excluded JSON file name, all strings there will be excluded")

	flag.StringVar(&options.OutputDirFlag, "o", "", "output directory where the translation files will be placed")
	flag.BoolVar(&options.OutputFlatFlag, "output-flat", true, "generated files are created in the specified output directory")
	flag.BoolVar(&options.OutputMatchPackageFlag, "output-match-package", false, "generated files are created in directory to match the package name")

	flag.StringVar(&options.FilenameFlag, "f", "", "the file name for which strings are extracted")
	flag.StringVar(&options.DirnameFlag, "d", "", "the dir name for which all .go files will have their strings extracted")
	flag.BoolVar(&options.RecurseFlag, "r", false, "recursively extract strings from all files in the same directory as filename or dirName")

	flag.StringVar(&options.IgnoreRegexp, "ignore-regexp", "", "a perl-style regular expression for files to ignore, e.g., \".*test.*\"")

	flag.Parse()
}

func usage() {
	usageString := `
gi18n [-command] [-vpe] [-o <outputDir>] -f <fileName> | -d [-r] [-ignore-regexp <regex>] <dirName>
  -h                        prints the usage

  -v                        verbose
  -dry-run                  prevents any output files from being created
  -p                        to generate standard .po files for translation

  EXTRACT-STRINGS:

  -extract-strings          the extract strings command flag

  -o                        the output directory where the translation files will be placed
  -output-flat              generated files are created in the specified output directory (default)
  -output-match-package     generated files are created in directory to match the package name

  -e                        the JSON file with strings to be excluded, defaults to excluded.json if present

  -f                        the go file name to extract strings

  -r                        recursesively extract strings from all subdirectories
  -d                        the directory containing the go files to extract strings

  -ignore-regexp            a perl-style regular expression for files to ignore, e.g., ".*test.*"

  CREATE-TRANSLATIONS:

  -create-translations      the create translations command flag

  -google-translate-api-key your public Google Translate API key which is used to generate translations (charge is applicable)

  -source-language		    the source language of the file, typically also part of the file name, e.g., \"en_US\"
  -languages 	            a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"
  -o                        the output directory where the newly created translation files will be placed
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
