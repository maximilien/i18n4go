// Code generated for package i18n by go-bindata DO NOT EDIT. (@generated)
// sources:
// i18n4go/i18n/resources/all.en_US.json
package i18n

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _i18n4goI18nResourcesAllEn_usJson = []byte(`[
  {
    "id": "i18n4go: Error compiling interpolated string Regexp: {{.Arg0}}\n",
    "translation": "i18n4go: Error compiling interpolated string Regexp: {{.Arg0}}\n"
  },
  {
    "id": "i18n4go: error reading content of init code snippet file: {{.Arg0}}\n, using default",
    "translation": "i18n4go: error reading content of init code snippet file: {{.Arg0}}\n, using default"
  },
  {
    "id": "verbose mode where lots of output is generated during execution",
    "translation": "verbose mode where lots of output is generated during execution"
  },
  {
    "id": "i18n4go: rewriting strings for source file:",
    "translation": "i18n4go: rewriting strings for source file:"
  },
  {
    "id": "generated files are created in the specified output directory",
    "translation": "generated files are created in the specified output directory"
  },
  {
    "id": "the root path to the Go source files whose packages are being rewritten, defaults to working directory, if not specified",
    "translation": "the root path to the Go source files whose packages are being rewritten, defaults to working directory, if not specified"
  },
  {
    "id": "i18n4go: Error checking input filename: ",
    "translation": "i18n4go: Error checking input filename: "
  },
  {
    "id": "i18n4go: attempting to use Google Translate to translate source strings in: ",
    "translation": "i18n4go: attempting to use Google Translate to translate source strings in: "
  },
  {
    "id": "i18n4go: WARNING target file contains total of invalid translations:",
    "translation": "i18n4go: WARNING target file contains total of invalid translations:"
  },
  {
    "id": "src",
    "translation": "src"
  },
  {
    "id": "i18n4go: error invoking Google Translate for string:",
    "translation": "i18n4go: error invoking Google Translate for string:"
  },
  {
    "id": "i18n4go: Could not show missing strings, err:",
    "translation": "i18n4go: Could not show missing strings, err:"
  },
  {
    "id": "i18n4go: creating and saving i18n strings to .po file:",
    "translation": "i18n4go: creating and saving i18n strings to .po file:"
  },
  {
    "id": "i18n4go: ERROR invoking Google Translate: ",
    "translation": "i18n4go: ERROR invoking Google Translate: "
  },
  {
    "id": "Verify strings in translations",
    "translation": "Verify strings in translations"
  },
  {
    "id": "i18n4go: Could not successfully rewrite package, err:",
    "translation": "i18n4go: Could not successfully rewrite package, err:"
  },
  {
    "id": "i18n4go: scanning file: ",
    "translation": "i18n4go: scanning file: "
  },
  {
    "id": "a directory with the extracted JSON files, using -output-match-package with -extract-strings this directory should match the input files package name",
    "translation": "a directory with the extracted JSON files, using -output-match-package with -extract-strings this directory should match the input files package name"
  },
  {
    "id": "i18n4go: inspecting dir {{.Arg0}}, recursive: {{.Arg1}}\n",
    "translation": "i18n4go: inspecting dir {{.Arg0}}, recursive: {{.Arg1}}\n"
  },
  {
    "id": "Removing these strings from the %s translation file:\n",
    "translation": "Removing these strings from the %s translation file:\n"
  },
  {
    "id": "Total time:",
    "translation": "Total time:"
  },
  {
    "id": "i18n4go: could not load i18n strings from file: {{.Arg0}}",
    "translation": "i18n4go: could not load i18n strings from file: {{.Arg0}}"
  },
  {
    "id": "the source go file to be rewritten",
    "translation": "the source go file to be rewritten"
  },
  {
    "id": "the output directory where the missing translation keys will be placed",
    "translation": "the output directory where the missing translation keys will be placed"
  },
  {
    "id": "ERROR opening file",
    "translation": "ERROR opening file"
  },
  {
    "id": "Updating the following strings from the %s translation file:\n",
    "translation": "Updating the following strings from the %s translation file:\n"
  },
  {
    "id": "Total files parsed:",
    "translation": "Total files parsed:"
  },
  {
    "id": "Adding to translated strings:",
    "translation": "Adding to translated strings:"
  },
  {
    "id": "prints the usage",
    "translation": "prints the usage"
  },
  {
    "id": "i18n4go: generated diff file:",
    "translation": "i18n4go: generated diff file:"
  },
  {
    "id": "output directory where the translation files will be placed",
    "translation": "output directory where the translation files will be placed"
  },
  {
    "id": "\"{{.Arg0}}\" exists in {{.Arg1}}, but not in {{.Arg2}}\n",
    "translation": "\"{{.Arg0}}\" exists in {{.Arg1}}, but not in {{.Arg2}}\n"
  },
  {
    "id": "a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command",
    "translation": "a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command"
  },
  {
    "id": "UNDER",
    "translation": "UNDER"
  },
  {
    "id": "i18n4go: saving combined language file: ",
    "translation": "i18n4go: saving combined language file: "
  },
  {
    "id": "Shows missing strings in translations",
    "translation": "Shows missing strings in translations"
  },
  {
    "id": "[optional] the substring capturing JSON file name, all strings there will only have their first capturing group saved as a translation",
    "translation": "[optional] the substring capturing JSON file name, all strings there will only have their first capturing group saved as a translation"
  },
  {
    "id": "saving file to path",
    "translation": "saving file to path"
  },
  {
    "id": "i18n4go: loading JSON strings from file: {{.Arg0}}\n",
    "translation": "i18n4go: loading JSON strings from file: {{.Arg0}}\n"
  },
  {
    "id": "Couldn't find any source strings: {{.Arg0}}",
    "translation": "Couldn't find any source strings: {{.Arg0}}"
  },
  {
    "id": "the command, one of: extract-strings, create-translations, rewrite-package, verify-strings, merge-strings, checkup, fixup",
    "translation": "the command, one of: extract-strings, create-translations, rewrite-package, verify-strings, merge-strings, checkup, fixup"
  },
  {
    "id": "cowardly refusing to translate the strings in test file:",
    "translation": "cowardly refusing to translate the strings in test file:"
  },
  {
    "id": "[optional] a comma separated list of target files for different languages to compare,  e.g., \\\"en, en_US, fr_FR, es\\\"\t                                                                  if not specified then the languages flag is used to find target files in same directory as source",
    "translation": "[optional] a comma separated list of target files for different languages to compare,  e.g., \\\"en, en_US, fr_FR, es\\\"\t                                                                  if not specified then the languages flag is used to find target files in same directory as source"
  },
  {
    "id": "General purpose tool for i18n",
    "translation": "General purpose tool for i18n"
  },
  {
    "id": "i18n4go: got a local import {{.Arg0}} so using {{.Arg1}} instead for pkg",
    "translation": "i18n4go: got a local import {{.Arg0}} so using {{.Arg1}} instead for pkg"
  },
  {
    "id": "i18n4go: could not save Google Translate i18n strings to file: {{.Arg0}}",
    "translation": "i18n4go: could not save Google Translate i18n strings to file: {{.Arg0}}"
  },
  {
    "id": "i18n4go: error determining the import path:",
    "translation": "i18n4go: error determining the import path:"
  },
  {
    "id": "i18n4go: could not create default translation file for language: {{.Arg0}}\nerr:{{.Arg1}}",
    "translation": "i18n4go: could not create default translation file for language: {{.Arg0}}\nerr:{{.Arg1}}"
  },
  {
    "id": "Could not load i18n asset: {{.Arg0}}",
    "translation": "Could not load i18n asset: {{.Arg0}}"
  },
  {
    "id": "i18n4go: rewriting strings in dir {{.Arg0}}, recursive: {{.Arg1}}\n",
    "translation": "i18n4go: rewriting strings in dir {{.Arg0}}, recursive: {{.Arg1}}\n"
  },
  {
    "id": "i18n4go: error saving updated i18n strings file:",
    "translation": "i18n4go: error saving updated i18n strings file:"
  },
  {
    "id": "targetFilenames:",
    "translation": "targetFilenames:"
  },
  {
    "id": "i18n4go: got a root pkg with import path:",
    "translation": "i18n4go: got a root pkg with import path:"
  },
  {
    "id": "i18n4go: Could not extract strings, err:",
    "translation": "i18n4go: Could not extract strings, err:"
  },
  {
    "id": "i18n4go: target file has invalid i18n strings with IDs: {{.Arg0}}",
    "translation": "i18n4go: target file has invalid i18n strings with IDs: {{.Arg0}}"
  },
  {
    "id": "i18n4go: created default translation file:",
    "translation": "i18n4go: created default translation file:"
  },
  {
    "id": "recursively extract strings from all files in the same directory as filename or dirName",
    "translation": "recursively extract strings from all files in the same directory as filename or dirName"
  },
  {
    "id": "Duplicated key found: ",
    "translation": "Duplicated key found: "
  },
  {
    "id": "WARNING compiling ignore-regexp:",
    "translation": "WARNING compiling ignore-regexp:"
  },
  {
    "id": "Couldn't get the strings from {{.Arg0}}: {{.Arg1}}",
    "translation": "Couldn't get the strings from {{.Arg0}}: {{.Arg1}}"
  },
  {
    "id": "Loaded {{.Arg0}} excluded regexps",
    "translation": "Loaded {{.Arg0}} excluded regexps"
  },
  {
    "id": "Extracting strings in package:",
    "translation": "Extracting strings in package:"
  },
  {
    "id": "Creates the transation files",
    "translation": "Creates the transation files"
  },
  {
    "id": "Loaded {{.Arg0}} excluded strings",
    "translation": "Loaded {{.Arg0}} excluded strings"
  },
  {
    "id": "i18n4go: could not save PO file: {{.Arg0}}",
    "translation": "i18n4go: could not save PO file: {{.Arg0}}"
  },
  {
    "id": "i18n4go: Error loading the i18n strings from input filename:",
    "translation": "i18n4go: Error loading the i18n strings from input filename:"
  },
  {
    "id": "WARNING running in -dry-run mode",
    "translation": "WARNING running in -dry-run mode"
  },
  {
    "id": "i18n4go: adding init func to package:",
    "translation": "i18n4go: adding init func to package:"
  },
  {
    "id": "Merge translation strings",
    "translation": "Merge translation strings"
  },
  {
    "id": " to output dir:",
    "translation": " to output dir:"
  },
  {
    "id": "i18n4go: Error input file: {{.Arg0}} is empty",
    "translation": "i18n4go: Error input file: {{.Arg0}} is empty"
  },
  {
    "id": "i18n4go: Could not verify strings for input filename, err:",
    "translation": "i18n4go: Could not verify strings for input filename, err:"
  },
  {
    "id": "Found",
    "translation": "Found"
  },
  {
    "id": "i18n4go: creating translation file copy for language:",
    "translation": "i18n4go: creating translation file copy for language:"
  },
  {
    "id": "i18n4go: Non-regular source file {{.Arg0}} ({{.Arg1}})\n",
    "translation": "i18n4go: Non-regular source file {{.Arg0}} ({{.Arg1}})\n"
  },
  {
    "id": "Show the version of the i18n client",
    "translation": "Show the version of the i18n client"
  },
  {
    "id": "[optional] create a *.extracted.json file with metadata such as: filename, directory, and positions of the strings in source file",
    "translation": "[optional] create a *.extracted.json file with metadata such as: filename, directory, and positions of the strings in source file"
  },
  {
    "id": "generate standard .po file for translation",
    "translation": "generate standard .po file for translation"
  },
  {
    "id": "i18n4go: ERROR could not create the diff file:",
    "translation": "i18n4go: ERROR could not create the diff file:"
  },
  {
    "id": "Missing Strings!",
    "translation": "Missing Strings!"
  },
  {
    "id": "i18n4go: WARNING target file contains total of extra keys:",
    "translation": "i18n4go: WARNING target file contains total of extra keys:"
  },
  {
    "id": "Extracted {{.Arg0}} strings from file: {{.Arg1}}\n",
    "translation": "Extracted {{.Arg0}} strings from file: {{.Arg1}}\n"
  },
  {
    "id": "a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"",
    "translation": "a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\""
  },
  {
    "id": "i18n4go: could not create translation file for language: {{.Arg0}} with Google Translate",
    "translation": "i18n4go: could not create translation file for language: {{.Arg0}} with Google Translate"
  },
  {
    "id": "i18n4go: created translation file with Google Translate:",
    "translation": "i18n4go: created translation file with Google Translate:"
  },
  {
    "id": "i18n4go: could not extract strings from directory:",
    "translation": "i18n4go: could not extract strings from directory:"
  },
  {
    "id": "OK",
    "translation": "OK"
  },
  {
    "id": "i18n4go: error adding init() func to package:",
    "translation": "i18n4go: error adding init() func to package:"
  },
  {
    "id": "i18n4go: creating translation file:",
    "translation": "i18n4go: creating translation file:"
  },
  {
    "id": "i18n4go: ERROR parsing Google Translate response body",
    "translation": "i18n4go: ERROR parsing Google Translate response body"
  },
  {
    "id": "Could not find imports for root node:\n\t{{.Arg0}}v\n",
    "translation": "Could not find imports for root node:\n\t{{.Arg0}}v\n"
  },
  {
    "id": "No match for ignore-regexp:",
    "translation": "No match for ignore-regexp:"
  },
  {
    "id": "Checks the transated files",
    "translation": "Checks the transated files"
  },
  {
    "id": "i18n4go: extracting strings from file:",
    "translation": "i18n4go: extracting strings from file:"
  },
  {
    "id": "the JSON file with strings to be excluded, defaults to excluded.json if present",
    "translation": "the JSON file with strings to be excluded, defaults to excluded.json if present"
  },
  {
    "id": "i18n4go: WARNING could not find JSON file:",
    "translation": "i18n4go: WARNING could not find JSON file:"
  },
  {
    "id": "i18n4go: input file: {{.Arg0}} is empty",
    "translation": "i18n4go: input file: {{.Arg0}} is empty"
  },
  {
    "id": "the source translation file",
    "translation": "the source translation file"
  },
  {
    "id": "the file name for which strings are extracted",
    "translation": "the file name for which strings are extracted"
  },
  {
    "id": "Add, update, or remove translation keys from source files and resources files",
    "translation": "Add, update, or remove translation keys from source files and resources files"
  },
  {
    "id": "Invalid response.",
    "translation": "Invalid response."
  },
  {
    "id": "i18n4go: got a pkg with import:",
    "translation": "i18n4go: got a pkg with import:"
  },
  {
    "id": "[optional] the excluded JSON file name, all strings there will be excluded",
    "translation": "[optional] the excluded JSON file name, all strings there will be excluded"
  },
  {
    "id": "a directory with the extracted JSON files, using -output-match-package with extract-strings command this directory should match the input files package name",
    "translation": "a directory with the extracted JSON files, using -output-match-package with extract-strings command this directory should match the input files package name"
  },
  {
    "id": "i18n4go: templated string is invalid, missing args in translation:",
    "translation": "i18n4go: templated string is invalid, missing args in translation:"
  },
  {
    "id": "An unexpected type of error",
    "translation": "An unexpected type of error"
  },
  {
    "id": "a comma separated list of target files for different languages to compare,  e.g., \\\"en, en_US, fr_FR, es\\\"\t                                                                  if not specified then the languages flag is used to find target files in same directory as source",
    "translation": "a comma separated list of target files for different languages to compare,  e.g., \\\"en, en_US, fr_FR, es\\\"\t                                                                  if not specified then the languages flag is used to find target files in same directory as source"
  },
  {
    "id": "Using ignore-regexp:",
    "translation": "Using ignore-regexp:"
  },
  {
    "id": "Creating and saving i18n strings to .po file:",
    "translation": "Creating and saving i18n strings to .po file:"
  },
  {
    "id": "Rewrite translated packages from go source files",
    "translation": "Rewrite translated packages from go source files"
  },
  {
    "id": "recursively rewrite packages from all files in the same directory as filename or dirName",
    "translation": "recursively rewrite packages from all files in the same directory as filename or dirName"
  },
  {
    "id": "the source language of the file, typically also part of the file name, e.g., \"en_US\"",
    "translation": "the source language of the file, typically also part of the file name, e.g., \"en_US\""
  },
  {
    "id": "Capturing substrings in file:",
    "translation": "Capturing substrings in file:"
  },
  {
    "id": "i18n4go: Could not create translation files, err:",
    "translation": "i18n4go: Could not create translation files, err:"
  },
  {
    "id": "Total rewritten strings:",
    "translation": "Total rewritten strings:"
  },
  {
    "id": "WARNING No capturing group found in {{.Arg0}}",
    "translation": "WARNING No capturing group found in {{.Arg0}}"
  },
  {
    "id": "Version:      {{.Arg0}}\n",
    "translation": "Version:      {{.Arg0}}\n"
  },
  {
    "id": "WARNING ignoring file:",
    "translation": "WARNING ignoring file:"
  },
  {
    "id": "Missing:",
    "translation": "Missing:"
  },
  {
    "id": "Total extracted strings:",
    "translation": "Total extracted strings:"
  },
  {
    "id": "WARNING: fail to compile ignore-regexp:",
    "translation": "WARNING: fail to compile ignore-regexp:"
  },
  {
    "id": "i18n4go: using import path as:",
    "translation": "i18n4go: using import path as:"
  },
  {
    "id": "i18n4go: creating translation files for:",
    "translation": "i18n4go: creating translation files for:"
  },
  {
    "id": "Excluding regexps in file:",
    "translation": "Excluding regexps in file:"
  },
  {
    "id": "Strings don't match",
    "translation": "Strings don't match"
  },
  {
    "id": "prevents any output files from being created",
    "translation": "prevents any output files from being created"
  },
  {
    "id": "generated files are created in directory to match the package name",
    "translation": "generated files are created in directory to match the package name"
  },
  {
    "id": "Git Revision: {{.Arg0}}\n",
    "translation": "Git Revision: {{.Arg0}}\n"
  },
  {
    "id": "Extract the translation strings from go source files",
    "translation": "Extract the translation strings from go source files"
  },
  {
    "id": "Saving extracted strings to file:",
    "translation": "Saving extracted strings to file:"
  },
  {
    "id": "Could not find:",
    "translation": "Could not find:"
  },
  {
    "id": "Excluding strings in file:",
    "translation": "Excluding strings in file:"
  },
  {
    "id": "i18n4go: target file is missing i18n strings with IDs: {{.Arg0}}",
    "translation": "i18n4go: target file is missing i18n strings with IDs: {{.Arg0}}"
  },
  {
    "id": "Saving extracted i18n strings to file:",
    "translation": "Saving extracted i18n strings to file:"
  },
  {
    "id": "a perl-style regular expression for files to ignore, e.g., \".*test.*\"",
    "translation": "a perl-style regular expression for files to ignore, e.g., \".*test.*\""
  },
  {
    "id": "WARNING error compiling regexp:",
    "translation": "WARNING error compiling regexp:"
  },
  {
    "id": "i18n4go: Could not merge strings, err:",
    "translation": "i18n4go: Could not merge strings, err:"
  },
  {
    "id": "Additional Strings!",
    "translation": "Additional Strings!"
  },
  {
    "id": "the substring capturing JSON file name, all strings there will only have their first capturing group saved as a translation",
    "translation": "the substring capturing JSON file name, all strings there will only have their first capturing group saved as a translation"
  },
  {
    "id": "i18n4go: error saving AST file:",
    "translation": "i18n4go: error saving AST file:"
  },
  {
    "id": "i18n4go: got a local import {{.Arg0}} so using {{.Arg1}} instead for root pkg",
    "translation": "i18n4go: got a local import {{.Arg0}} so using {{.Arg1}} instead for root pkg"
  },
  {
    "id": "Fixup the transation files",
    "translation": "Fixup the transation files"
  },
  {
    "id": "[optional] the path to a file containing the template snippet for the code that is used for go-i18n initialization",
    "translation": "[optional] the path to a file containing the template snippet for the code that is used for go-i18n initialization"
  },
  {
    "id": "i18n4go: ERROR input file does not match target file:",
    "translation": "i18n4go: ERROR input file does not match target file:"
  },
  {
    "id": "the dir name for which all .go files will have their strings extracted",
    "translation": "the dir name for which all .go files will have their strings extracted"
  },
  {
    "id": "Build Date:   {{.Arg0}}\n",
    "translation": "Build Date:   {{.Arg0}}\n"
  },
  {
    "id": "i18n4go: inserting i18n.T() calls for strings that need to be translated",
    "translation": "i18n4go: inserting i18n.T() calls for strings that need to be translated"
  },
  {
    "id": "Loaded {{.Arg0}} substring regexps",
    "translation": "Loaded {{.Arg0}} substring regexps"
  },
  {
    "id": "Extracted total of {{.Arg0}} strings\n\n",
    "translation": "Extracted total of {{.Arg0}} strings\n\n"
  },
  {
    "id": "Select the number for the previous translation:",
    "translation": "Select the number for the previous translation:"
  },
  {
    "id": "Additional:",
    "translation": "Additional:"
  },
  {
    "id": "[optional] your public Google Translate API key which is used to generate translations (charge is applicable)",
    "translation": "[optional] your public Google Translate API key which is used to generate translations (charge is applicable)"
  },
  {
    "id": "the code",
    "translation": "the code"
  },
  {
    "id": "i18n4go: Error verifying target filename: ",
    "translation": "i18n4go: Error verifying target filename: "
  },
  {
    "id": "i18n4go: WARNING target file has extra key with ID: ",
    "translation": "i18n4go: WARNING target file has extra key with ID: "
  },
  {
    "id": "the output directory where the newly created translation files will be placed",
    "translation": "the output directory where the newly created translation files will be placed"
  },
  {
    "id": "i18n4go: Error compiling templated string Regexp: {{.Arg0}}\n",
    "translation": "i18n4go: Error compiling templated string Regexp: {{.Arg0}}\n"
  },
  {
    "id": "i18n4go: Could not checkup, err:",
    "translation": "i18n4go: Could not checkup, err:"
  },
  {
    "id": "Couldn't find the english strings: {{.Arg0}}",
    "translation": "Couldn't find the english strings: {{.Arg0}}"
  },
  {
    "id": "i18n4go: error appending i18n.T() to AST file:",
    "translation": "i18n4go: error appending i18n.T() to AST file:"
  },
  {
    "id": "i18n4go: WARNING target file has invalid templated translations with key ID: ",
    "translation": "i18n4go: WARNING target file has invalid templated translations with key ID: "
  },
  {
    "id": "Unable to find english translation files",
    "translation": "Unable to find english translation files"
  },
  {
    "id": "Adding these strings to the %s translation file:\n",
    "translation": "Adding these strings to the %s translation file:\n"
  },
  {
    "id": "a JSON file with the strings that should be i18n enabled, typically the output of the extract-strings command",
    "translation": "a JSON file with the strings that should be i18n enabled, typically the output of the extract-strings command"
  },
  {
    "id": "i18n4go: error getting root path import:",
    "translation": "i18n4go: error getting root path import:"
  },
  {
    "id": "the directory containing the go files to validate",
    "translation": "the directory containing the go files to validate"
  },
  {
    "id": "Canceling fixup",
    "translation": "Canceling fixup"
  },
  {
    "id": "[optional] the qualifier string that is used when using the i18n.T(...) function, default to nothing but could be set to ` + "`" + `i18n` + "`" + ` so that all calls would be: i18n.T(...)",
    "translation": "[optional] the qualifier string that is used when using the i18n.T(...) function, default to nothing but could be set to ` + "`" + `i18n` + "`" + ` so that all calls would be: i18n.T(...)"
  },
  {
    "id": "i18n4go: could not create output directory: {{.Arg0}}",
    "translation": "i18n4go: could not create output directory: {{.Arg0}}"
  },
  {
    "id": "Could not find an i18n file for locale: en_US",
    "translation": "Could not find an i18n file for locale: en_US"
  },
  {
    "id": "Error when inspecting go file: ",
    "translation": "Error when inspecting go file: "
  },
  {
    "id": "i18n4go: Error loading the i18n strings from target filename:",
    "translation": "i18n4go: Error loading the i18n strings from target filename:"
  },
  {
    "id": "Is the string \"%s\" a new or updated string? [new/upd]\n",
    "translation": "Is the string \"%s\" a new or updated string? [new/upd]\n"
  },
  {
    "id": "i18n4go: determining import path using root path:",
    "translation": "i18n4go: determining import path using root path:"
  },
  {
    "id": "\nSomething completely unexpected happened. This is a bug in %s.\nPlease file this bug : https://github.com/maximilien/i18n4go/issues\nTell us that you ran this command:\n\n\t%s\n\nthis error occurred:\n\n\t%s\n\nand this stack trace:\n\n%s\n\t",
    "translation": "\nSomething completely unexpected happened. This is a bug in %s.\nPlease file this bug : https://github.com/maximilien/i18n4go/issues\nTell us that you ran this command:\n\n\t%s\n\nthis error occurred:\n\n\t%s\n\nand this stack trace:\n\n%s\n\t"
  },
  {
    "id": "{{.Arg0}}\nVersion {{.Arg1}}",
    "translation": "{{.Arg0}}\nVersion {{.Arg1}}"
  },
  {
    "id": "i18n4go: Could not fixup, err:",
    "translation": "i18n4go: Could not fixup, err:"
  },
  {
    "id": "File has duplicated key: {{.Arg0}}\n{{.Arg1}}",
    "translation": "File has duplicated key: {{.Arg0}}\n{{.Arg1}}"
  },
  {
    "id": "i18n4go: target file has extra i18n strings with IDs: {{.Arg0}}",
    "translation": "i18n4go: target file has extra i18n strings with IDs: {{.Arg0}}"
  },
  {
    "id": "i18n4go: using the PWD as the rootPath:",
    "translation": "i18n4go: using the PWD as the rootPath:"
  }
]
`)

func i18n4goI18nResourcesAllEn_usJsonBytes() ([]byte, error) {
	return _i18n4goI18nResourcesAllEn_usJson, nil
}

func i18n4goI18nResourcesAllEn_usJson() (*asset, error) {
	bytes, err := i18n4goI18nResourcesAllEn_usJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "i18n4go/i18n/resources/all.en_US.json", size: 27145, mode: os.FileMode(420), modTime: time.Unix(1719593700, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"i18n4go/i18n/resources/all.en_US.json": i18n4goI18nResourcesAllEn_usJson,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//
//	data/
//	  foo.txt
//	  img/
//	    a.png
//	    b.png
//
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"i18n4go": &bintree{nil, map[string]*bintree{
		"i18n": &bintree{nil, map[string]*bintree{
			"resources": &bintree{nil, map[string]*bintree{
				"all.en_US.json": &bintree{i18n4goI18nResourcesAllEn_usJson, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
