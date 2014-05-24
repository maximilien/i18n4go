package input_files

import (
	"fmt"
	"github.com/cloudfoundry/cli/cf/i18n"
	goi18n "github.com/nicksnyder/go-i18n/i18n"
	"os"
	"path/filepath"
)

var T goi18n.TranslateFunc

func init() {
	var err error
	T, err = i18n.Init("input_files", i18n.GetResourcesPath())
	if err != nil {
		panic(err)
	}
}
func Something() {
	fmt.Printf("HAI")
	if os.Getenv("SOMETHING") {
		fmt.Printf(filepath.Clean(os.Getenv("something")))
	}
}
