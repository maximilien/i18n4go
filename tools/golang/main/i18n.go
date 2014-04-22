package main

import (
	"fmt"
	"flag"
	"os"
	"strings"
	"runtime/debug"

	"../extract_strings"
)

var fileNameFlag string
var dirNameFlag string
var recurseFlag	bool

func main() {
	defer handlePanic()

	if fileNameFlag == "" && dirNameFlag == "" {
		usage()
		return
	}

	if fileNameFlag != "" {
		extract_strings.InspectFile(fileNameFlag)
	} else {
		extract_strings.InspectDir(dirNameFlag, recurseFlag)
	}
}

func init() {
	flag.StringVar(&fileNameFlag, "f", "", "the file name for which strings are extracted")
	flag.StringVar(&dirNameFlag, "d", "", "the dir name for which all .go files will have their strings extracted")
	flag.BoolVar(&recurseFlag, "r", false, "recursesively extract strings from all files in the same directory as filename or dirName")

	flag.Parse()
}

func usage() {
	usageString := `
gi18n -f <fileName> | [-d <dirName> | -r -d <dirName>]
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
