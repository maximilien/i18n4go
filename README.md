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
gi18n -extract-strings [-vpe] [-o <outputDir>] -f <fileName> | -d [-r] [-ignore-regexp <regex>] <dirName>
  -h                      prints the usage

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
```

Extracting Strings
==================

Extracting strings from go files

```
$ gi18n -extract-strings -v -p -o src/i18n -f ../tmp/cli/cf/app/app.go

gi18n: extracting strings from file: ../tmp/cli/cf/app/app.go
Excluding strings in file: excluded.json
Loaded 38 excluded strings
Excluding regexps in file: excluded.json
Loaded 4 excluded regexps
Extracted 7 strings from file: ../tmp/cli/cf/app/app.go
Saving extracted strings to file: ../tmp/cli/cf/app/app.go.extracted.json
Saving extracted i18n strings to file: ../tmp/cli/cf/app/app.go.en.json
Creating and saving i18n strings to .po file: ../tmp/cli/cf/app/app.go.en.po
Total time: 1.229666ms
```

The output for the command above are three files, of which two are important for translation:

1. `../tmp/cli/cf/app/app.go.en.json`

This file is the JSON formatted translation file for English. Some of its content is as follows:

```json
[
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

b. `../tmp/cli/cf/app/app.go.en.po`

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

Specifying `excluded.json` File
===============================

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
