package cmds

type Options struct {
	HelpFlag bool

	ExtractStringsCmdFlag     bool
	CreateTranslationsCmdFlag bool
	VerifyStringsCmdFlag      bool
	RewritePackageCmdFlag     bool
	MergeStringsCmdFlag       bool

	VerboseFlag bool
	DryRunFlag  bool
	PoFlag      bool

	SourceLanguageFlag        string
	LanguagesFlag             string
	GoogleTranslateApiKeyFlag string

	OutputDirFlag          string
	OutputMatchImportFlag  bool
	OutputMatchPackageFlag bool
	OutputFlatFlag         bool

	ExcludedFilenameFlag string
	FilenameFlag         string
	DirnameFlag          string

	RecurseFlag bool

	IgnoreRegexpFlag string

	LanguageFilesFlag string

	I18nStringsFilenameFlag string
}

type CommandInterface interface {
	Println(a ...interface{}) (int, error)
	Printf(msg string, a ...interface{}) (int, error)
	Options() Options
	Run() error
}
