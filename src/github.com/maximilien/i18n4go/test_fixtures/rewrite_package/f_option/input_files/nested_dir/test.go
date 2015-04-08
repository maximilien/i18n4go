package nested_dir

import (
	"fmt"
	"os"
	"path/filepath"
)

type t struct {
	myString string
}

func Something() string {
	someString := "hello"
	var anotherString string = "world"
	println(someString, anotherString)

	yetAnotherString := []string{"tricky tests"}
	var moreStrings []string
	moreStrings = []string{"are", "tricky"}
	println(yetAnotherString, moreStrings)

	mappyMap := map[string]string{"hello": "world"}
	println(mappyMap)

	myT := t{myString: "my string"}
	println(myT)

	trickyT := t{"this is a tricky case"}
	println(trickyT)

	fmt.Printf("HAI")
	if os.Getenv("SOMETHING") {
		fmt.Printf(filepath.Clean(os.Getenv("something")))
	}

	return "enqueuedequeueenqueuebananapants"
}
