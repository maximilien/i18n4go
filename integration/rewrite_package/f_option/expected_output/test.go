package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

type t struct {
	myString string
}

func Something() string {
	someString := T("hello")
	var anotherString string = T("world")
	println(someString, anotherString)

	yetAnotherString := []string{T("tricky tests")}
	var moreStrings []string
	moreStrings = []string{T("are"), T("tricky")}
	println(yetAnotherString, moreStrings)

	mappyMap := map[string]string{T("hello"): T("world")}
	println(mappyMap)
	println(mappyMap[T("hello")])

	myT := t{myString: T("my string")}
	println(myT)

	trickyT := t{T("this is a tricky case")}
	println(trickyT)

	concatenatedStrings := T("foo") + T(" ") + T("bar")
	println(concatenatedStrings)

	fmt.Printf(T("HAI"))
	if os.Getenv(T("SOMETHING")) {
		fmt.Printf(filepath.Clean(os.Getenv(T("something"))))
	}

	return T("enqueuedequeueenqueuebananapants")
}
