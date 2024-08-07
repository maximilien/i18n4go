package nested_dir

import (
	"path/filepath"

	"github.com/maximilien/i18n4go/i18n4go/i18n"
)

var T i18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "f_option", "input_files", "nested_dir"), i18n.GetResourcesPath(), func(asset string) ([]byte, error) {
		return Asset(asset)
	})
}
