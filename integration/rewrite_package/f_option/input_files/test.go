package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

func Something() {
	someString := "hello"
	var anotherString string = "world"
	println(someString, anotherString)

	yetAnotherString := []string{"tricky tests"}
	var moreStrings []string
	moreStrings = []string{"are", "tricky"}
	println(yetAnotherString, moreStrings)

	mappyMap := map[string]string{"hello": "world"}
	println(mappyMap)

	fmt.Printf("HAI")
	if os.Getenv("SOMETHING") {
		fmt.Printf(filepath.Clean(os.Getenv("something")))
	}
}
