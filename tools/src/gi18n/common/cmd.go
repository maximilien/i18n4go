package common

type Options struct {
	ExtractStringsCmdFlag bool
	VerboseFlag           bool
	PoFlag                bool

	ExcludedFilenameFlag string
	FilenameFlag         string
	DirnameFlag          string

	RecurseFlag bool

	IgnoreRegexp string
}

type I18nStringInfo struct {
	ID          string `json:"id"`
	Translation string `json:"translation"`
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

var BLANKS = []string{", ", "\t", "\n", "\n\t", "\t\n"}
