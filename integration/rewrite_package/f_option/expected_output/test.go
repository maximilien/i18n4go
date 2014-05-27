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
	T, err = i18n.Init(T("input_files"), i18n.GetResourcesPath())
	if err != nil {
		panic(err)
	}
}
func Something() {
	someString := T("hello")
	var anotherString string = T("world")
	println(someString, anotherString)

	yetAnotherString := []string{T("tricky tests")}
	var moreStrings []string
	moreStrings = []string{T("are"), T("tricky")}
	println(yetAnotherString, moreStrings)

	mappyMap := map[string]string{T("hello"): T("world")}
	println(mappyMap)

	fmt.Printf(T("HAI"))
	if os.Getenv(T("SOMETHING")) {
		fmt.Printf(filepath.Clean(os.Getenv(T("something"))))
	}
}
