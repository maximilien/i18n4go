i18n Tooling for the Go Language [![Build Status](https://travis-ci.org/maximilien/i18n4go.svg?branch=master)](https://travis-ci.org/maximilien/i18n4go#)
==============================

This is a general purpose internationalization (i18n) tooling for Go language programs. It allows you to prepare Go language code for internationalization and localization (l11n). You can also use it to help you maintain the resulting i18n-enabled Go code so that it remains internationalized. This tool was extracted while we worked on enabling the [Cloud Foundry CLI](https://github.com/cloudfoundry/cli) with i18n support. 

This tool is licensed under the [Apache 2.0 OSS license](https://github.com/maximilien/i18n4go/blob/master/LICENSE). We'd love to hear from you if you are using, attempting to use, or planning to use this tool. Feel free to send email to `i18n4go` at the Gmail domain.

## Getting Started
------------------

### Overview Presentations
--------------------------

* Talk at [GoSF Meetup](http://www.meetup.com/golangsf/events/220603955/) on April, 2015. Slides ([PDF](https://github.com/maximilien/presentations/blob/master/2015/i18n4go-gosf-meetup/releases/i18n4go-v0.4.1.pdf) and [PPTX](https://github.com/maximilien/presentations/blob/master/2015/i18n4go-gosf-meetup/releases/i18n4go-v0.4.1.pptx)), video, [demo](https://github.com/maximilien/i18n4go/tree/master/examples/demo1)

### Cloning and Building
------------------------

Clone this repo and build it. Using the following commands on a Linux or Mac OS X system:

```
$ mkdir -p i18n4go/src/github.com/maximilien
$ export GOPATH=$(pwd)/i18n4go:$GOPATH
$ cd i18n4go/src/github.com/maximilien
$ git clone https://github.com/maximilien/i18n4go.git 
$ cd i18n4go
$ ./bin/build
```

NOTE: if you get any dependency errors, then use `go get path/to/dependency` to get it, e.g., `go get github.com/onsi/ginkgo` and `go get github.com/onsi/gomega`

The executable output should now be located in: `out/i18n4go`. Place it wherever you want, e.g., `/usr/local/bin` on Linux or Mac OS X.

You can now use the `i18n4go` executable to issue some of the typical i18n tooling processes.

### Running Tests
-----------------

You should run the tests to make sure all is well, do this with: `$ ./bin/test` in your cloned repository.

The output should be similar to:

```
$ bin/test

 Cleaning build artifacts...

 Formatting packages...

 Integration Testing packages:
ok  	github.com/maximilien/i18n4go/integration/checkup	1.571s
ok  	github.com/maximilien/i18n4go/integration/create_translations	1.542s
ok  	github.com/maximilien/i18n4go/integration/extract_strings	1.694s
ok  	github.com/maximilien/i18n4go/integration/fixup	1.657s
ok  	github.com/maximilien/i18n4go/integration/merge_strings	1.645s
ok  	github.com/maximilien/i18n4go/integration/rewrite_package	1.853s
ok  	github.com/maximilien/i18n4go/integration/show_missing_strings	1.590s
?   	github.com/maximilien/i18n4go/integration/test_helpers	[no test files]
ok  	github.com/maximilien/i18n4go/integration/verify_strings	1.701s

 Vetting packages for potential issues...

SWEET SUITE SUCCESS
```

### Typical Workflow
--------------------

The recommended workflow is to use the commands (documented below) in the following order. For each command, use the command's help or this README for details and to experiment for your project. So the nine steps are:

1. **extract-strings** which will automatically extract every string from your Go source files and create a JSON and optionally a PO file.

2. **merge-strings** to create one file and removing what is not needed and strings you do not want to i18n. This could be important and time consuming but to help this process, we've found that it's good to keep a list of all the strings that you do not want to i18n as well as string patterns (as regex). Take a look at the CF CLI [excluded.json](https://github.com/cloudfoundry/cli/blob/master/cf/i18n/excluded.json) for a real world example file you might end up with. The regex in there might be useful to reuse.

3. might need to do 1 again, but using `excluded.json`. The outcome should be the file or files for `en_US` for all the strings that will be i18n for your app. So for instance, if you decide to combine all into one: `en_US.all.json`

4. **rewrite-package** using the file or files in 3. This will rewrite your code to use the `T(...)` function and also deal with parameters to your strings, using the pattern: `Arg0`, `Arg1`, etc. *NOTE* that this step will rewrite (yes, modify) your code. You can always use `go fmt` so the code will look fine. All files that contain strings that need to be i18n will be rewritten. You can do this step one package at a time.

5. **create-translations** to create initial translation file or files for each language that you want to support. 
For instance to create `fr_FR` file(s) for French and every other locale_Language you specify. This could be done manually. The reason to use tool is optional next step and also because the tool may help streamline your build process... The resulting files can be sent to human translators to be officially completed.

6. [optional] **create-translations** with [Google Translate API](https://cloud.google.com/translate/docs). You will need to have a Google Translate API key (*NOTE*: might require you to pay or at least enter your credit card if usage is above some threshold). Generally the strings generated by Google Translate are OK, but not great. They usually require additional work, however, we have found that they can be a good start when sending files to be officially translated by human translator team(s).

7. **verify-strings** this will help you ensure that your translation files, e.g., `en_US.all.json` and `fr_FR.all.json`, and others, all have the same keys. This is *important* since if you are missing a key then for that language you might crash your app. We recommend using this during your build and for CI and not build resulting app in 8 (next step) if this step fails.

8. package your app with your i18n resource files. The packaging is slightly tricky since one of the great value of Golang is to have one binary file distribution for your app. This means you need to convert your i18n resource files (the JSON files) into binary that can be loaded in code (as source code). We've been using [go-bindata](https://github.com/jteeuwen/go-bindata) for the CF CLI and that seems to work pretty well. See [this script](https://github.com/cloudfoundry/cli/blob/fa7bcb07cdb6c6960f0907022bcef83ec4363a47/bin/generate-language-resources) on how we used it in the CF CLI. Other alternatives exist but we have not tried them.

9. ship and profit :)

### Typical Workflow Diagram

![Typical i18n4go workflow diagram](https://github.com/maximilien/i18n4go/blob/master/docs/images/typical-workflow.png)

### Help
--------

Printing the usage help: `$ i18n4go -h` or `$ i18n4go --help`

```
usage: i18n4go -c extract-strings [-vpe] [--dry-run] [--output-flat|--output-match-package|-o <outputDir>] -f <fileName>
   or: i18n4go -c extract-strings [-vpe] [--dry-run] [--output-flat|--output-match-package|-o <outputDir>] -d <dirName> [-r] [--ignore-regexp <fileNameRegexp>]

usage: i18n4go -c rewrite-package [-v] [-r] -d <dirName> [--i18n-strings-filename <fileName> | --i18n-strings-dirname <dirName>]
   or: i18n4go -c rewrite-package [-v] [-r] -f <fileName> --i18n-strings-filename <fileName>

usage: i18n4go -c merge-strings [-v] [-r] [--source-language <language>] -d <dirName>

usage: i18n4go -c verify-strings [-v] [--source-language <language>] -f <sourceFileName> --language-files <language files>
   or: i18n4go -c verify-strings [-v] [--source-language <language>] -f <sourceFileName> --languages <lang1,lang2,...>

usage: i18n4go -c create-translations [-v] [--google-translate-api-key <api key>] [--source-language <language>] -f <fileName> --languages <lang1,lang2,...> -o <outputDir>

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
$ i18n4go -c extract-strings -v --po -f ./tmp/cli/cf/app/app.go -o ./tmp/cli/i18n -output-match-package

i18n4go: extracting strings from file: ./tmp/cli/cf/app/app.go
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
$ i18n4go -c extract-strings -v --po -d ./tmp/cli/cf/app/ -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*"

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
$ i18n4go -c merge-strings -v -d ./tmp/cli/i18n/app -source-language en

i18n4go: scanning file: tmp/cli/i18n/app/app.go.en.json
i18n4go: scanning file: tmp/cli/i18n/app/flag_helper.go.en.json
i18n4go: scanning file: tmp/cli/i18n/app/help.go.en.json
i18n4go: saving combined language file: tmp/cli/i18n/app/en.all.json
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
$ i18n4go -c rewrite-package -v -f tmp/cli/cf/app/help.go -i18n-strings-dirname tmp/cli/i18n/app/ -o tmp/cli/cf/app/

i18n4go: rewriting strings for source file: tmp/cli/cf/app/help.go
i18n4go: adding init func to package: app  to output dir: tmp/cli/cf/app
i18n4go: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/help.go

Total files parsed: 1
Total extracted strings: 17
Total time: 9.986963ms
```

b. running it on a directory

```
$ i18n4go -c rewrite-package -v -d tmp/cli/cf/app/ -i18n-strings-dirname tmp/cli/i18n/app/ -o tmp/cli/cf/app/

i18n4go: rewriting strings in dir tmp/cli/cf/app/, recursive: false

i18n4go: loading JSON strings from file: tmp/cli/i18n/app/app.go.en.json
i18n4go: rewriting strings for source file: tmp/cli/cf/app/app.go
i18n4go: adding init func to package: app  to output dir: tmp/cli/cf/app
i18n4go: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/app.go
i18n4go: loading JSON strings from file: tmp/cli/i18n/app/flag_helper.go.en.json
i18n4go: rewriting strings for source file: tmp/cli/cf/app/flag_helper.go
i18n4go: adding init func to package: app  to output dir: tmp/cli/cf/app
i18n4go: inserting T() calls for strings that need to be translated
saving file to path tmp/cli/cf/app/flag_helper.go
i18n4go: loading JSON strings from file: tmp/cli/i18n/app/help.go.en.json
i18n4go: rewriting strings for source file: tmp/cli/cf/app/help.go
i18n4go: adding init func to package: app  to output dir: tmp/cli/cf/app
i18n4go: inserting T() calls for strings that need to be translated
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
$ i18n4go -create-translations -v -f tmp/cli/i18n/app/en.all.json -source-language en -languages "en_US,fr_FR,es_ES,de_DE" -o tmp/cli/i18n/app/

i18n4go: creating translation files for: tmp/cli/i18n/app/en.all.json

i18n4go: creating translation file copy for language: en_US
i18n4go: creating translation file: tmp/cli/i18n/app/en_US.all.json
i18n4go: created default translation file: tmp/cli/i18n/app/en_US.all.json
i18n4go: creating translation file copy for language: fr_FR
i18n4go: creating translation file: tmp/cli/i18n/app/fr_FR.all.json
i18n4go: created default translation file: tmp/cli/i18n/app/fr_FR.all.json
i18n4go: creating translation file copy for language: es_ES
i18n4go: creating translation file: tmp/cli/i18n/app/es_ES.all.json
i18n4go: created default translation file: tmp/cli/i18n/app/es_ES.all.json
i18n4go: creating translation file copy for language: de_DE
i18n4go: creating translation file: tmp/cli/i18n/app/de_DE.all.json
i18n4go: created default translation file: tmp/cli/i18n/app/de_DE.all.json

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
$ i18n4go -c verify-strings -v -f tmp/cli/i18n/app/en.all.json -languages "fr"

targetFilenames: [tmp/cli/i18n/app/fr.all.json]
i18n4go: ERROR input file does not match target file: tmp/cli/i18n/app/fr.all.json
i18n4go: generated diff file: tmp/cli/i18n/app/fr.all.json.missing.diff.json
i18n4go: Error verifying target filename:  tmp/cli/i18n/app/fr.all.json
i18n4go: Could not verify strings for input filename, err: i18n4go: target file is missing i18n strings with IDs: --,'%v',-
```

Similarly, `verify-strings` will make sure that no additonal strings are added. So if we had an additional German `de.all.json` file that included additional strings
running `verify-strings` would include a `tmp/cli/i18n/app/de.all.json.extra.diff.json`.

```
$ i18n4go -v -verify-strings -f tmp/cli/i18n/app/en.all.json -languages "fr,de"

targetFilenames: [tmp/cli/i18n/app/fr.all.json tmp/cli/i18n/app/de.all.json]
i18n4go: ERROR input file does not match target file: tmp/cli/i18n/app/fr.all.json
i18n4go: generated diff file: tmp/cli/i18n/app/fr.all.json.missing.diff.json
i18n4go: Error verifying target filename:  tmp/cli/i18n/app/fr.all.json
i18n4go: WARNING target file has extra key with ID:  advanced
i18n4go: WARNING target file has extra key with ID:  apps
i18n4go: WARNING target file contains total of extra keys: 2
i18n4go: generated diff file: tmp/cli/i18n/app/de.all.json.extra.diff.json
i18n4go: Error verifying target filename:  tmp/cli/i18n/app/de.all.json
i18n4go: Could not verify strings for input filename, err: i18n4go: target file has extra i18n strings with IDs: advanced,apps
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
$ i18n4go -c extract-strings -p -d ./tmp/cli/cf/app/ -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*" -e ./example/excluded.json
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
$ i18n4go -c extract-strings -v -d ./tmp/cli/cf/api/resources -o ./tmp/cli/i18n -output-match-package -ignore-regexp ".*test.*" -e ./example/excluded.json
```

We can inspect the `./tmp/cli/i18n/resources/events.go.en.json` file and see that there are no strings with the expression `json:`.

---------

### Edge Releases (master)
--------------------------

Get latest code here on Github and build it: `./bin/build` *

The binary will be in the `./out` * directory.

You can follow our development progress on [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1071880).

## Troubleshooting / FAQs
-------------------------

None for now. Submit questions/comments as issues and we will update here

### Filing Bugs
---------------

##### For simple bugs (eg: text formatting, help messages, etc), please provide

- the command options you ran
- what occurred
- what you expected to occur

##### For panics and other crashes, please provide

- the command you ran
- the stack trace generated (if any)
- any other relevant information

## Cloning the repository
-------------------------

1. Install [Go](https://golang.org)
1. Clone (Forking beforehand for development).
1. [Ensure your $GOPATH is set correctly](http://golang.org/cmd/go/#hdr-GOPATH_environment_variable)

## Building *
-------------

1. Run `./bin/build`
1. The binary will be built into the `./out` directory.

Optionally, you can use `bin/run` to compile and run the executable in one step.

## Developing *
---------------

1. Run `go get code.google.com/p/go.tools/cmd/vet`
2. Run `go get github.com/cloudfoundry/cli ...` to install test dependencies
1. Write a [Ginkgo](https://github.com/onsi/ginkgo) test.
1. Run `bin/test` and watch the test fail.
1. Make the test pass.
1. Submit a pull request.

## Contributing
---------------

* We welcome any and all contributions as Pull Requests (PR)
* We also welcome issues and bug report and new feature request. We will address as time permits
* Follow the steps above in Developing to get your system setup correctly
* Please make sure your PR is passing Travis before submitting
* Feel free to email me or the current collaborators if you have additional questions about contributions

### Managing dependencies
-------------------------

* All dependencies managed via [Godep](https://github.com/tools/godep). See [Godeps/_workspace](https://github.com/maximilien/i18n4go/tree/master/Godeps/_workspace) directory on master

### Current conventions
-----------------------

* Basic Go conventions
* Strict TDD for any code added or changed
* Go fakes when needing to mock objects

(*) these items are in the works, we will remove the * once they are available
