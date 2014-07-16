package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"runtime/debug"

	"github.com/maximilien/i18n4go/cmds"
	"github.com/maximilien/i18n4go/common"
)

var options common.Options

func main() {
	defer handlePanic()

	if options.HelpFlag || options.LongHelpFlag {
		usage()
		return
	}

	switch options.CommandFlag {
	case "extract-strings":
		extractStringsCmd()
	case "create-translations":
		createTranslationsCmd()
	case "rewrite-package":
		rewritePackageCmd()
	case "verify-strings":
		verifyStringsCmd()
	case "merge-strings":
		mergeStringsCmd()
	case "show-missing-strings":
		showMissingStringsCmd()
	default:
		usage()
	}
}

func extractStringsCmd() {
	if options.HelpFlag || (options.FilenameFlag == "" && options.DirnameFlag == "") {
		usage()
		return
	}

	cmd := cmds.NewExtractStrings(options)

	startTime := time.Now()

	err := cmd.Run()
	if err != nil {
		cmd.Println("gi18n: Could not extract strings, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time:", duration)
}

func createTranslationsCmd() {
	if options.HelpFlag || (options.FilenameFlag == "") {
		usage()
		return
	}

	cmd := cmds.NewCreateTranslations(options)

	startTime := time.Now()

	err := cmd.Run()
	if err != nil {
		cmd.Println("gi18n: Could not create translation files, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time:", duration)
}

func verifyStringsCmd() {
	if options.HelpFlag || (options.FilenameFlag == "") {
		usage()
		return
	}

	cmd := cmds.NewVerifyStrings(options)

	startTime := time.Now()

	err := cmd.Run()
	if err != nil {
		cmd.Println("gi18n: Could not verify strings for input filename, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time:", duration)
}

func rewritePackageCmd() {
	if options.HelpFlag || (options.FilenameFlag == "" && options.DirnameFlag == "" &&
		(options.I18nStringsFilenameFlag == "" || options.I18nStringsDirnameFlag == "")) {
		usage()
		return
	}

	cmd := cmds.NewRewritePackage(options)

	startTime := time.Now()

	err := cmd.Run()
	if err != nil {
		cmd.Println("gi18n: Could not successfully rewrite package, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	cmd.Println("Total time:", duration)
}

func mergeStringsCmd() {
	if options.HelpFlag || (options.DirnameFlag == "") {
		usage()
		return
	}

	mergeStrings := cmds.NewMergeStrings(options)

	startTime := time.Now()

	err := mergeStrings.Run()
	if err != nil {
		mergeStrings.Println("gi18n: Could not merge strings, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	mergeStrings.Println("Total time:", duration)
}

func showMissingStringsCmd() {
	if options.HelpFlag || (options.DirnameFlag == "") {
		usage()
		return
	}

	showMissingStrings := cmds.NewShowMissingStrings(options)

	startTime := time.Now()

	err := showMissingStrings.Run()
	if err != nil {
		showMissingStrings.Println("gi18n: Could not show missing strings, err:", err)
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	showMissingStrings.Println("Total time:", duration)
}

func init() {
	flag.StringVar(&options.CommandFlag, "c", "", "the command, one of: extract-strings, create-translations, rewrite-package, verify-strings, merge-strings")

	flag.BoolVar(&options.HelpFlag, "h", false, "prints the usage")
	flag.BoolVar(&options.LongHelpFlag, "-help", false, "prints the usage")

	flag.StringVar(&options.SourceLanguageFlag, "source-language", "en", "the source language of the file, typically also part of the file name, e.g., \"en_US\"")
	flag.StringVar(&options.LanguagesFlag, "languages", "", "a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"")
	flag.StringVar(&options.GoogleTranslateApiKeyFlag, "google-translate-api-key", "", "[optional] your public Google Translate API key which is used to generate translations (charge is applicable)")

	flag.BoolVar(&options.VerboseFlag, "v", false, "verbose mode where lots of output is generated during execution")

	flag.BoolVar(&options.PoFlag, "po", false, "generate standard .po file for translation")

	flag.BoolVar(&options.MetaFlag, "meta", false, "[optional] create a *.extracted.json file with metadata such as: filename, directory, and positions of the strings in source file")
	flag.BoolVar(&options.DryRunFlag, "dry-run", false, "prevents any output files from being created")

	flag.StringVar(&options.ExcludedFilenameFlag, "e", "excluded.json", "[optional] the excluded JSON file name, all strings there will be excluded")

	flag.StringVar(&options.OutputDirFlag, "o", "", "output directory where the translation files will be placed")

	flag.BoolVar(&options.OutputFlatFlag, "output-flat", true, "generated files are created in the specified output directory")
	flag.BoolVar(&options.OutputMatchPackageFlag, "output-match-package", false, "generated files are created in directory to match the package name")

	flag.StringVar(&options.FilenameFlag, "f", "", "the file name for which strings are extracted")

	flag.StringVar(&options.DirnameFlag, "d", "", "the dir name for which all .go files will have their strings extracted")

	flag.BoolVar(&options.RecurseFlag, "r", false, "recursively extract strings from all files in the same directory as filename or dirName")

	flag.StringVar(&options.IgnoreRegexpFlag, "ignore-regexp", ".*test.*", "a perl-style regular expression for files to ignore, e.g., \".*test.*\"")

	flag.StringVar(&options.LanguageFilesFlag, "language-files", "", `[optional] a comma separated list of target files for different languages to compare,  e.g., \"en, en_US, fr_FR, es\"	                                                                  if not specified then the languages flag is used to find target files in same directory as source`)

	flag.StringVar(&options.I18nStringsFilenameFlag, "i18n-strings-filename", "", "a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command")
	flag.StringVar(&options.I18nStringsDirnameFlag, "i18n-strings-dirname", "", "a directory with the extracted JSON files, using -output-match-package with -extract-strings this directory should match the input files package name")
	flag.StringVar(&options.RootPathFlag, "root-path", "", "the root path to the Go source files whose packages are being rewritten, defaults to working directory, if not specified")

	flag.StringVar(&options.InitCodeSnippetFilenameFlag, "init-code-snippet-filename", "", "[optional] the path to a file containing the template snippet for the code that is used for go-i18n initialization")

	flag.Parse()
}

func usage() {
	usageString := `
usage: gi18n -c extract-strings [-vpe] [--dry-run] [--output-flat|--output-match-package|-o <outputDir>] -f <fileName>
   or: gi18n -c extract-strings [-vpe] [--dry-run] [--output-flat|--output-match-package|-o <outputDir>] -d <dirName> [-r] [--ignore-regexp <fileNameRegexp>]

usage: gi18n -c rewrite-package [-v] [-r] -d <dirName> [--i18n-strings-filename <fileName> | --i18n-strings-dirname <dirName>] [--init-code-snippet-filename <fileName>]
   or: gi18n -c rewrite-package [-v] [-r] -f <fileName> --i18n-strings-filename <fileName> [--init-code-snippet-filename <fileName>]

usage: gi18n -c create-translations [-v] [--google-translate-api-key <api key>] [--source-language <language>] -f <fileName> --languages <lang1,lang2,...> -o <outputDir>

usage: gi18n -c merge-strings [-v] [-r] [--source-language <language>] -d <dirName>

usage: gi18n -c verify-strings [-v] [--source-language <language>] -f <sourceFileName> --language-files <language files>
   or: gi18n -c verify-strings [-v] [--source-language <language>] -f <sourceFileName> --languages <lang1,lang2,...>

usage: gi18n -c show-missing-strings [-v] -d <dirName> --i18n-strings-filename <language file>

  -h | --help                prints the usage
  -v                         verbose

  EXTRACT-STRINGS:

  -c extract-strings         the extract strings command

  --po                       to generate standard .po files for translation
  -e                         [optional] the JSON file with strings to be excluded, defaults to excluded.json if present
  --meta                     [optional] create a *.extracted.json file with metadata such as: filename, directory, and positions of the strings in source file
  --dry-run                  [optional] prevents any output files from being created


  --output-flat              generated files are created in the specified output directory (default)
  --output-match-package     generated files are created in directory to match the package name
  -o                         the output directory where the translation files will be placed

  -f                         the go file name to extract strings

  -d                         the directory containing the go files to extract strings

  -r                         [optional] recursesively extract strings from all subdirectories
  --ignore-regexp            [optional] a perl-style regular expression for files to ignore, e.g., ".*test.*"

  REWRITE-PACKAGE:

  -c rewrite-package         the rewrite package command
  -f                         the source go file to be rewritten
  -d                         the directory containing the go files to rewrite

  --i18n-strings-filename    a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command
  --i18n-strings-dirname     a directory with the extracted JSON files, using -output-match-package with -extract-strings this directory should match the input files package name
  --root-path                the root path to the Go source files whose packages are being rewritten, defaults to working directory, if not specified

  --init-code-snippet-filename [optional] the path to a file containing the template snippet for the code that is used for go-i18n initialization"
  -o                           [optional] output diretory for rewritten file. If not specified, the original file will be overwritten
  
  MERGE STRINGS:

  -c merge-strings           the merge strings command which merges multiple <filename>.go.<language>.json files into a <language>.all.json

  -r                         [optional] recursesively combine files from all subdirectories
  --source-language          [optional] the source language of the file, typically also part of the file name, e.g., "en_US" (default to 'en')

  -d                         the directory containing the json files to combine

  CREATE-TRANSLATIONS:

  -c create-translations     the create translations command

  --google-translate-api-key [optional] your public Google Translate API key which is used to generate translations (charge is applicable)
  --source-language          [optional] the source language of the file, typically also part of the file name, e.g., \"en_US\"

  -f                         the source translation file
  --languages                a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"
  -o                         the output directory where the newly created translation files will be placed

  VERIFY-STRINGS:

  -c verify-strings          the verify strings command

  --source-language          [optional] the source language of the source translation file (default to 'en')

  -f                         the source translation file

  --language-files           a comma separated list of target files for different languages to compare, e.g., "en, en_US, fr_FR, es"
                             if not specified then the languages flag is used to find target files in same directory as source
  --languages                a comma separated list of valid languages with optional territory, e.g., "en, en_US, fr_FR, es"

  SHOW-MISSING-STRINGS:

  -c show-missing-strings    the missing strings command

  -d                         the directory containing the go files to validate
  --i18n-strings-filename    a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command
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
