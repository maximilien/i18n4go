package input_files

import (
	"path/filepath"

	i18n "github.com/maximilien/i18n4go/i18n4go/i18n"
)

var T i18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "init_code_snippet_filename", "input_files"), i18n.GetResourcesPath())
}
