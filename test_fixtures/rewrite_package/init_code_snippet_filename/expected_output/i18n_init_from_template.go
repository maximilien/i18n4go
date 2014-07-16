package input_files

import (
	"path/filepath"

	"github.com/cloudfoundry/cf/cli/i18n"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
)

var T goi18n.TranslateFunc

func init() {
	T = i18n.Init(filepath.Join("test_fixtures", "rewrite_package", "f_option", "input_files"), i18n.GetResourcesPath())
}
