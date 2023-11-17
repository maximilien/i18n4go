package nested_dir

import (
	"path/filepath"

	i18n "github.com/maximilien/i18n4go/i18n4go/i18n"
)

var T i18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "f_option", "input_files", "nested_dir"), i18n.GetResourcesPath())
}
