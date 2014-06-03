Cloud Foundry CLI i18n Tooling [![Build Status](https://travis-ci.org/cloudfoundry/cli.png?branch=master)](https://travis-ci.org/maximilien/i18n4cf)
==============================

This is the official i18n tooling for the Cloud Foundry command line client. It allows you to prepare Go language code for internationalization (i18n).

Getting Started
===============
Download and run the installer for your platform from the section below. If you are on OS X, you can also install the gi18n tool
with homebrew--run `brew install gi18n` *.

Once installed, you can use it to issue some of the typical i18n tooling processes.

Printing Usage
===========

Printing the usage help

```
usage: gi18n -extract-strings [-vpe] [-dry-run] [-output-flat|-output-match-package|-o <outputDir>] -f <fileName>
   or: gi18n -extract-strings [-vpe] [-dry-run] [-output-flat|-output-match-package|-o <outputDir>] -d <dirName> [-r] [-ignore-regexp <fileNameRegexp>]

usage: gi18n -merge-strings [-v] [-r] [-source-language <language>] -d <dirName>

usage: gi18n -verify-strings [-v] [-source-language <language>] -f <sourceFileName> -language-files <language files>
   or: gi18n -verify-strings [-v] [-source-language <language>] -f <sourceFileName> -languages <lang1,lang2,...>

usage: gi18n -create-translations [-v] [-google-translate-api-key <api key>] [-source-language <language>] -f <fileName> -languages <lang1,lang2,...> -o <outputDir>

  -h                        prints the usage
  -v                        verbose

  EXTRACT-STRINGS:

  -extract-strings          the extract strings command

  -p                        to generate standard .po files for translation
  -e                        [optional] the JSON file with strings to be excluded, defaults to excluded.json if present
  -meta                     [optional] create a *.extracted.json file with metadata such as: filename, directory, and positions of the strings in source file
  -dry-run                  [optional] prevents any output files from being created


  -output-flat              generated files are created in the specified output directory (default)
  -output-match-package     generated files are created in directory to match the package name
  -o                        the output directory where the translation files will be placed

  -f                        the go file name to extract strings

  -d                        the directory containing the go files to extract strings

  -r                        [optional] recursesively extract strings from all subdirectories
  -ignore-regexp            [optional] a perl-style regular expression for files to ignore, e.g., ".*test.*"

  MERGE STRINGS:

  -merge-strings            merges multiple <filename>.go.<language>.json files into a <language>.all.json

  -r                        [optional] recursesively combine files from all subdirectories
  -source-language          [optional] the source language of the file, typically also part of the file name, e.g., \"en_US\" (default to 'en')

  -d                        the directory containing the json files to combine

  VERIFY-STRINGS:

  -verify-strings           the verify strings command

  -source-language          [optional] the source language of the source translation file (default to 'en')

  -f                        the source translation file

  -language-files           a comma separated list of target files for different languages to compare, e.g., \"en, en_US, fr_FR, es\"
                            if not specified then the languages flag is used to find target files in same directory as source
  -languages                a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"

  REWRITE-PACKAGE:

  -f                        the source go file to be rewritten
  -d                        the directory containing the go files to rewrite
  -i18n-strings-filename    a JSON file with the strings that should be i18n enabled, typically the output of -extract-strings command
  -o                        [optional] output diretory for rewritten file. If not specified, the original file will be overwritten

  CREATE-TRANSLATIONS:

  -create-translations      the create translations command

  -google-translate-api-key [optional] your public Google Translate API key which is used to generate translations (charge is applicable)
  -source-language          [optional] the source language of the file, typically also part of the file name, e.g., \"en_US\"

  -f                        the source translation file
  -languages                a comma separated list of valid languages with optional territory, e.g., \"en, en_US, fr_FR, es\"
  -o                        the output directory where the newly created translation files will be placed
```

## extract-strings
Extracting strings from go files.  For the examples below we are running the tool on a copy of the the CloudFoundry CLI cloned in the `./tmp`

```
$ gi18n -extract-strings -v -p -f ./tmp/cli/cf/app/app.go -o ./tmp/cli/i18n -output-match-package

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

1. `./tmp/cli/i18n/app/app.go.en.json`

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
$ gi18n -extract-strings -v -p -d ./tmp/cli/cf/app/ -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*"

...
```

The generated output JSON files are in: `./tmp/cli/i18n/app`

## merge-strings

Combining strings in multiple `*.go.[lang].json` files generated by `Extract Strings` into one file. Using the same example source as above.

```
$ gi18n -merge-strings -v -d ./tmp/cli/i18n/app -source-language en

gi18n: scanning file: tmp/cli/i18n/app/app.go.en.json
gi18n: scanning file: tmp/cli/i18n/app/flag_helper.go.en.json
gi18n: scanning file: tmp/cli/i18n/app/help.go.en.json
gi18n: saving combined language file: tmp/cli/i18n/app/en.all.json
Total time: 1.283116ms
```

The output for the command above is one file placed in the same directory as the JSON files being merged: `en.all.json`.
This file containes one formatted translation for each translation generated by extract-strings for English. The `-source-language` flag
must match the language portion of the files in the directory, e.g., app.go.en.json, where the language is "en".

## verify-strings

TODO

## rewrite-package

TODO

## create-translations

TODO

# Specifying `excluded.json` File

TODO

Stable Release
==============

Installers *
-----------
- [Debian 32 bit]()
- [Debian 64 bit]()
- [Redhat 32 bit]()
- [Redhat 64 bit]()
- [Mac OS X 64 bit]()
- [Windows 32 bit]()
- [Windows 64 bit]()

Edge Releases (master)
======================

Get latest code here on Github and build it: `./bin/build` *

The binary will be in the `./out` * directory.

You can follow our development progress on [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1071880).

Troubleshooting / FAQs
======================

Linux
-----

TBD

Filing Bugs
===========

##### For simple bugs (eg: text formatting, help messages, etc), please provide

- the command options you ran
- what occurred
- what you expected to occur

##### For panics and other crashes, please provide

- the command you ran
- the stack trace generated (if any)
- any other relevant information

Cloning the repository
======================

1. Install [Go](https://golang.org)
1. Clone (Forking beforehand for development).
1. [Ensure your $GOPATH is set correctly](http://golang.org/cmd/go/#hdr-GOPATH_environment_variable)

Building *
==========

1. Run `./bin/build`
1. The binary will be built into the `./out` directory.

Optionally, you can use `bin/run` to compile and run the executable in one step.

Developing *
============

1. Run `go get code.google.com/p/go.tools/cmd/vet`
2. Run `go get github.com/cloudfoundry/cli ...` to install test dependencies
1. Write a Ginkgo test.
1. Run `bin/test` and watch the test fail.
1. Make the test pass.
1. Submit a pull request.

Contributing
============

Architecture overview
---------------------
TODO

Managing dependencies
---------------------

TODO

Example command
---------------

TODO

Current conventions
===================

TODO

(*) these items are in the works, we will remove the * once they are available
