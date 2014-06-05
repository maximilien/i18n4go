package nested_dir

import (
	"path/filepath"

	"github.com/cloudfoundry/cli/cf/i18n"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

var T goi18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "f_option", "input_files", "nested_dir"), i18n.GetResourcesPath())
}
