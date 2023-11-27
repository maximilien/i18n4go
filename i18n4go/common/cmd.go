// Copyright © 2015-2023 The Knative Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

type Options struct {
	CommandFlag string

	HelpFlag     bool
	LongHelpFlag bool

	VerboseFlag bool
	DryRunFlag  bool
	PoFlag      bool
	MetaFlag    bool

	SourceLanguageFlag        string
	LanguagesFlag             string
	GoogleTranslateApiKeyFlag string

	OutputDirFlag          string
	OutputMatchImportFlag  bool
	OutputMatchPackageFlag bool
	OutputFlatFlag         bool

	ExcludedFilenameFlag  string
	SubstringFilenameFlag string
	FilenameFlag          string
	DirnameFlag           string

	RecurseFlag bool

	IgnoreRegexpFlag string

	LanguageFilesFlag string

	I18nStringsFilenameFlag string
	I18nStringsDirnameFlag  string

	RootPathFlag string

	InitCodeSnippetFilenameFlag string

	QualifierFlag string
}

type I18nStringInfo struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
	Modified    bool   `json:"modified"`
}

type StringInfo struct {
	Filename string `json:"filename"`
	Value    string `json:"value"`
	Offset   int    `json:"offset"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
}

type ExcludedStrings struct {
	ExcludedStrings []string `json:"excludedStrings"`
	ExcludedRegexps []string `json:"excludedRegexps"`
}

type PrinterInterface interface {
	Println(a ...interface{}) (int, error)
	Printf(msg string, a ...interface{}) (int, error)
}

var BLANKS = []string{", ", "\t", "\n", "\n\t", "\t\n"}
