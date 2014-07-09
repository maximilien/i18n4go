Cloud Foundry CLI i18n Tooling [![Build Status](https://travis-ci.org/cloudfoundry/cli.png?branch=master)](https://travis-ci.org/maximilien/i18n4go)
==============================

This is a general purpose i18n tooling for Go language prorams. It allows you to prepare Go language code for internationalization (i18n). This tool was extracted while we worked on enabling the Cloud Foundry CLI for i18n.

## Getting Started
==================
Download and run the installer for your platform from the section below. If you are on OS X, you can also install the gi18n tool
with homebrew--run `brew install gi18n` *.

Once installed, you can use it to issue some of the typical i18n tooling processes.

### Printing Usage
==================

Printing the usage help

```
usage: gi18n -c extract-strings [-vpe] [--dry-run] [--output-flat|--output-match-package|-o <outputDir>] -f <fileName>
   or: gi18n -c extract-strings [-vpe] [--dry-run] [--output-flat|--output-match-package|-o <outputDir>] -d <dirName> [-r] [--ignore-regexp <fileNameRegexp>]

usage: gi18n -c rewrite-package [-v] [-r] -d <dirName> [--i18n-strings-filename <fileName> | --i18n-strings-dirname <dirName>]
   or: gi18n -c rewrite-package [-v] [-r] -f <fileName> --i18n-strings-filename <fileName>

usage: gi18n -c merge-strings [-v] [-r] [--source-language <language>] -d <dirName>

usage: gi18n -c verify-strings [-v] [--source-language <language>] -f <sourceFileName> --language-files <language files>
   or: gi18n -c verify-strings [-v] [--source-language <language>] -f <sourceFileName> --languages <lang1,lang2,...>

usage: gi18n -c create-translations [-v] [--google-translate-api-key <api key>] [--source-language <language>] -f <fileName> --languages <lang1,lang2,...> -o <outputDir>

  -h | --help                prints the usage
  -v                         verbose
...
```

## extract-strings

The general usage for `-extract-strings` command is:

```
  ...
  EXTRACT-STRINGS:

  -c extract-strings         the extract strings command

  -e                         [optional] the JSON file with strings to be excluded, defaults to excluded.json if present
  
  --po                       to generate standard .po files for translation
  --meta                     [optional] create a *.extracted.json file with metadata such as: filename, directory, and positions of the strings in source file
  --dry-run                  [optional] prevents any output files from being created


  -o                         the output directory where the translation files will be placed
  -f                         the go file name to extract strings
  -d                         the directory containing the go files to extract strings
  -r                         [optional] recursesively extract strings from all subdirectories

  --output-flat              generated files are created in the specified output directory (default)
  --output-match-package     generated files are created in directory to match the package name

  --ignore-regexp            [optional] a perl-style regular expression for files to ignore, e.g., ".*test.*"

```

The command `-c extract-strings` pulls strings out of go files.  For the examples below we are running the tool on a copy of the the [CloudFoundry CLI](https://github.com/cloudfoundry/cli) cloned in the `./tmp`

```
$ gi18n -c extract-strings -v -p -f ./tmp/cli/cf/app/app.go -o ./tmp/cli/i18n -output-match-package

gi18n: extracting strings from file: ./tmp/cli/cf/app/app.go
Could not find: excluded.json
Loaded 0 excluded strings
Could not find: excluded.json
Loaded 0 excluded regexps
Extracted 10 strings from file: ./tmp/cli/cf/app/app.go
Saving extracted i18n strings to file: tmp/cli/i18n/app/app.go.en.json
Creating and saving i18n strings to .po file: ./tmp/cli/cf/app/app.go.en.po
Total time: 3.962ms
```

The output for the command above are three files, of which two are important for translation:

a. `./tmp/cli/i18n/app/app.go.en.json`

This file is the JSON formatted translation file for English. Some of its content is as follows:

```json
[
    ...
    {
        "id": "Show help",
        "translation": "Show help"
    },
    {
        "id": "%s help [COMMAND]",
        "translation": "%s help [COMMAND]"
    },
  ...
]
```

b. Optionaly, using -p flag, it will generate `./tmp/cli/i18n/app/app.go.en.po`

This file is the [PO](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html) formatted translation file for English. Some of its content is as follows:

```
# filename: ../tmp/cli/cf/app/app.go, offset: 1617, line: 48, column: 16
msgid "Show help"
msgstr "Show help"

# filename: ../tmp/cli/cf/app/app.go, offset: 1657, line: 49, column: 28
msgid "%s help [COMMAND]"
msgstr "%s help [COMMAND]"
...
```

To extract multiples files that are in one directory, use the following:

```
$ gi18n -c extract-strings -v -p -d ./tmp/cli/cf/app/ -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*"

...
```

The generated output JSON files are in: `./tmp/cli/i18n/app`

## merge-strings

The general usage for `-merge-strings` command is:

```
  ...
  MERGE STRINGS:

  -c merge-strings            merges multiple <filename>.go.<language>.json files into a <language>.all.json

  -d                         the directory containing the json files to combine
  -r                         [optional] recursesively combine files from all subdirectories

  --source-language          [optional] the source language of the file, typically also part of the file name, e.g., "en_US" (default to 'en')

```

The command `-c merge-strings` combines strings in multiple `*.go.[lang].json` files generated by `Extract Strings` into one file. Using the same example source as above.

```
$ gi18n -c merge-strings -v -d ./tmp/cli/i18n/app -source-language en

gi18n: scanning file: tmp/cli/i18n/app/app.go.en.json
gi18n: scanning file: tmp/cli/i18n/app/flag_helper.go.en.json
gi18n: scanning file: tmp/cli/i18n/app/help.go.en.json
gi18n: saving combined language file: tmp/cli/i18n/app/en.all.json
Total time: 1.283116ms
```

The output for the command above is one file placed in the same directory as the JSON files being merged: `en.all.json`.
This file containes one formatted translation for each translation generated by extract-strings for English. The `-source-language` flag
must match the language portion of the files in the directory, e.g., app.go.en.json, where the language is "en".

## rewrite-package

The general usage for `-rewrite-package` command is:

```
  ...
  REWRITE-PACKAGE:
  
  -c rewrite-package         the rewrite package command
  
  -f                         the source go file to be rewritten
  -d                         the directory containing the go files to rewrite
  -o                         [optional] output diretory for rewritten file. If not specified, the original file will be overwritten
  
  --i18n-strings-filename    a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command
  --i18n-strings-dirname     a directory with the extracted JSON files, using -output-match-package with -extract-strings this directory should match the input files package name
  --root-path                the root path to the Go source files whose packages are being rewritten, defaults to working directory, if not specified

```

The command `-c rewrite-package` will modify the go source files such that every string identified in the JSON translation files are wrapped with the `T()` function. There are two cases:

a. running it on one source file

```
$ gi18n -c rewrite-package -v -f tmp/cli/cf/app/help.go -i18n-strings-dirname tmp/cli/i18n/app/ -o tmp/cli/cf/app/

gi18n: rewriting strings for source file: tmp/cli/cf/app/help.go
gi18n: adding init func to package: app  to output dir: tmp/cli/cf/app
gi18n: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/help.go

Total files parsed: 1
Total extracted strings: 17
Total time: 9.986963ms
```

b. running it on a directory

```
$ gi18n -c rewrite-package -v -d tmp/cli/cf/app/ -i18n-strings-dirname tmp/cli/i18n/app/ -o tmp/cli/cf/app/

gi18n: rewriting strings in dir tmp/cli/cf/app/, recursive: false

gi18n: loading JSON strings from file: tmp/cli/i18n/app/app.go.en.json
gi18n: rewriting strings for source file: tmp/cli/cf/app/app.go
gi18n: adding init func to package: app  to output dir: tmp/cli/cf/app
gi18n: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/app.go
gi18n: loading JSON strings from file: tmp/cli/i18n/app/flag_helper.go.en.json
gi18n: rewriting strings for source file: tmp/cli/cf/app/flag_helper.go
gi18n: adding init func to package: app  to output dir: tmp/cli/cf/app
gi18n: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/flag_helper.go
gi18n: loading JSON strings from file: tmp/cli/i18n/app/help.go.en.json
gi18n: rewriting strings for source file: tmp/cli/cf/app/help.go
gi18n: adding init func to package: app  to output dir: tmp/cli/cf/app
gi18n: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/help.go

Total files parsed: 3
Total extracted strings: 21
Total time: 16.648105ms
```

In both cases above the `-i18n-strings-dirname` specifies the directory containing the `<source.go>.en.json` file with the strings to process.
However, this can be replaced with `-i18n-strings-filename` and specify one JSON file (e.g., `en.all.json`) which contains all the strings.

---------

The result in each case is that the source files are rewritten with the wrapped `T()` function but also dealing with converting interpolated strings into Go-style templated strings. For instance:

The following interpolated string: `"%s help [COMMAND]"` is templated to: `"{{.Arg0}} help [COMMAND]"` and rewritten automaticall as:

```
T("{{.Arg0}} help [COMMAND]", map[string]interface{}{"Arg0": cf.Name()})
```

So in essence the strings in the JSON files that where interpolated become templated, that is new IDs for the default language.

## create-translations

The general usage for `-c create-translations` command is:

```
  ...
  CREATE-TRANSLATIONS:

  -c create-translations     the create translations command

  -f                         the source translation file
  -o                         the output directory where the newly created translation files will be placed

  --languages                a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"
  --source-language          [optional] the source language of the file, typically also part of the file name, e.g., \"en_US\"
  --google-translate-api-key [optional] your public Google Translate API key which is used to generate translations (charge is applicable)

```

The command `-create-translations` generates copies of the `-source-language` file, one per language specified in the `-languages` flag (seperated by comma).

```
$ gi18n -create-translations -v -f tmp/cli/i18n/app/en.all.json -source-language en -languages "en_US,fr_FR,es_ES,de_DE" -o tmp/cli/i18n/app/

gi18n: creating translation files for: tmp/cli/i18n/app/en.all.json

gi18n: creating translation file copy for language: en_US
gi18n: creating translation file: tmp/cli/i18n/app/en_US.all.json
gi18n: created default translation file: tmp/cli/i18n/app/en_US.all.json
gi18n: creating translation file copy for language: fr_FR
gi18n: creating translation file: tmp/cli/i18n/app/fr_FR.all.json
gi18n: created default translation file: tmp/cli/i18n/app/fr_FR.all.json
gi18n: creating translation file copy for language: es_ES
gi18n: creating translation file: tmp/cli/i18n/app/es_ES.all.json
gi18n: created default translation file: tmp/cli/i18n/app/es_ES.all.json
gi18n: creating translation file copy for language: de_DE
gi18n: creating translation file: tmp/cli/i18n/app/de_DE.all.json
gi18n: created default translation file: tmp/cli/i18n/app/de_DE.all.json

Total time: 2.143251ms
```

Optionally, we can create automated translations for the generated copies using Google Translate[link] passing the `google-translate-api-key` flag.

## verify-strings

The general usage for `-c verify-strings` command is:

```
  ...
  VERIFY-STRINGS:

  -c verify-strings          the verify strings command


  -f                         the source translation file

  --source-language          [optional] the source language of the source translation file (default to 'en')
  --languages                a comma separated list of valid languages with optional territory, e.g., "en, en_US, fr_FR, es"
  --language-files           a comma separated list of target files for different languages to compare, e.g., "en, en_US, fr_FR, es"
                             if not specified then the languages flag is used to find target files in same directory as source

```

The command `-verify-strings` assures that combined language files have exactly the same keys.

For instance, in the example in `merge-strings` we created a combined language file called `./tmp/cli/i18n/app/en.all.json` and if we also
had a `./tmp/cli/i18n/app/fr.all.json` for French and that file had missing strings then running the `verify-strings` would generate a
`tmp/cli/i18n/app/fr.all.json.missing.diff.json`, as in the following:

```
$ gi18n -c verify-strings -v -f tmp/cli/i18n/app/en.all.json -languages "fr"

targetFilenames: [tmp/cli/i18n/app/fr.all.json]
gi18n: ERROR input file does not match target file: tmp/cli/i18n/app/fr.all.json
gi18n: generated diff file: tmp/cli/i18n/app/fr.all.json.missing.diff.json
gi18n: Error verifying target filename:  tmp/cli/i18n/app/fr.all.json
gi18n: Could not verify strings for input filename, err: gi18n: target file is missing i18n strings with IDs: --,'%v',-
```

Similarly, `verify-strings` will make sure that no additonal strings are added. So if we had an additional German `de.all.json` file that included additional strings
running `verify-strings` would include a `tmp/cli/i18n/app/de.all.json.extra.diff.json`.

```
$ gi18n -v -verify-strings -f tmp/cli/i18n/app/en.all.json -languages "fr,de"

targetFilenames: [tmp/cli/i18n/app/fr.all.json tmp/cli/i18n/app/de.all.json]
gi18n: ERROR input file does not match target file: tmp/cli/i18n/app/fr.all.json
gi18n: generated diff file: tmp/cli/i18n/app/fr.all.json.missing.diff.json
gi18n: Error verifying target filename:  tmp/cli/i18n/app/fr.all.json
gi18n: WARNING target file has extra key with ID:  advanced
gi18n: WARNING target file has extra key with ID:  apps
gi18n: WARNING target file contains total of extra keys: 2
gi18n: generated diff file: tmp/cli/i18n/app/de.all.json.extra.diff.json
gi18n: Error verifying target filename:  tmp/cli/i18n/app/de.all.json
gi18n: Could not verify strings for input filename, err: gi18n: target file has extra i18n strings with IDs: advanced,apps
```

Finally, if a combined language file contains both extra and missing keys then `verify-strings` will generate two diff files: `missing` and `extra`.

## Specifying `excluded.json` File

The exclude.json file can be used to manage which strings should not be extract with the `extracting-strings` command. In the `excluded.json` file,
you can specifie string literals to ignore as well as classes of strings using a Perl-style regular expression. We have provided an example file
[exclude](https://github.com/maximilien/i18n4go/blob/master/example/excluded.json) to demonstrate the string and regexp cases.

### string literals

String literals are defined within the `excludedStrings` array. Any strings in your source files that exactly matches one of these will not be extracted
to the generated files from extracted strings.

```
}
  "excludedStrings" : [
     "",
    " ",
    "\n",
    "help",
    ...
  ] ...
}
```

As an example run, generate an extracted string files useing the command:

```
$ gi18n -c extract-strings -p -d ./tmp/cli/cf/app/ -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*" -e ./example/excluded.json
```

If we inspect the `./tmp/cli/i18n/app/app.go.en.json` file there should not be an entry for `"id": "help"`, but you should still see an entry for `"id": "show help"`

### regular expressions 

Regular expressions can be defined using the same JSON annotation as string literals with the tag `"excludedRegexps"`.

```
{
...
"excludedRegexps" : [
   "^\\w$",
   "^json:"
 ]
}
```

As an example for regular expressions, let us consider the `^json:`. This expression will remove any string containg `json:` which would be useful when parsing the
`./tmp/cli/cf/api/resources/events.go` file such as: `ExitDescription string `json:"exit_description"`. After running the command:

```
$ gi18n -c extract-strings -v -d ./tmp/cli/cf/api/resources -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*" -e ./example/excluded.json
```

We can inspect the `./tmp/cli/i18n/resources/events.go.en.json` file and see that there are no strings with the expression `json:`.

## Stable Release
=================

### Installers *
----------------
- [Debian 32 bit]()
- [Debian 64 bit]()
- [Redhat 32 bit]()
- [Redhat 64 bit]()
- [Mac OS X 64 bit]()
- [Windows 32 bit]()
- [Windows 64 bit]()

### Edge Releases (master)
==========================

Get latest code here on Github and build it: `./bin/build` *

The binary will be in the `./out` * directory.

You can follow our development progress on [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1071880).

## Troubleshooting / FAQs
=========================

### Linux
---------

TBD

### Filing Bugs
===============

##### For simple bugs (eg: text formatting, help messages, etc), please provide

- the command options you ran
- what occurred
- what you expected to occur

##### For panics and other crashes, please provide

- the command you ran
- the stack trace generated (if any)
- any other relevant information

## Cloning the repository
=========================

1. Install [Go](https://golang.org)
1. Clone (Forking beforehand for development).
1. [Ensure your $GOPATH is set correctly](http://golang.org/cmd/go/#hdr-GOPATH_environment_variable)

## Building *
=============

1. Run `./bin/build`
1. The binary will be built into the `./out` directory.

Optionally, you can use `bin/run` to compile and run the executable in one step.

## Developing *
===============

1. Run `go get code.google.com/p/go.tools/cmd/vet`
2. Run `go get github.com/cloudfoundry/cli ...` to install test dependencies
1. Write a Ginkgo test.
1. Run `bin/test` and watch the test fail.
1. Make the test pass.
1. Submit a pull request.

## Contributing
===============

### Architecture overview
------------------------

TODO

### Managing dependencies
-------------------------

TODO

### Example command
-------------------

TODO

### Current conventions
========================

TODO

(*) these items are in the works, we will remove the * once they are available
