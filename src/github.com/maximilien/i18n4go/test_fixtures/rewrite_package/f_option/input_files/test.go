package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

type t struct {
	myString string
}

func myPrint() {
	println("isn't that grand")
}

func Something() string {
	myPrint() // a trivial case

	someString := "hello"
	var anotherString string = "world"
	println(someString, anotherString)

	yetAnotherString := []string{"tricky tests"}
	var moreStrings []string
	moreStrings = []string{"are", "tricky"}
	println(yetAnotherString, moreStrings)

	mappyMap := map[string]string{"hello": "world"}
	println(mappyMap)
	println(mappyMap["hello"])

	myT := t{myString: "my string"}
	println(myT.myString)

	trickyT := t{"this is a tricky case"}
	println(trickyT.myString)

	concatenatedStrings := "foo" + " " + "bar"
	println(concatenatedStrings)

	fmt.Printf("HAI")
	if os.Getenv("SOMETHING") != "" {
		fmt.Printf(filepath.Clean(os.Getenv("SOMETHING")))
	}

	fmt.Println("hello" + "world")

	return "enqueuedequeueenqueuebananapants"
}
