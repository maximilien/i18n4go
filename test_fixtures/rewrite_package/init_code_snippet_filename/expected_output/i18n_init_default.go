package input_files

import (
	"path/filepath"

	"github.com/maximilien/i18n4go/i18n4go/i18n"
)

var T i18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "init_code_snippet_filename", "input_files"), i18n.GetResourcesPath(), func(asset string) ([]byte, error) {
		return Asset(asset)
	})
}
