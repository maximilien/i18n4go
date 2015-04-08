package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

func Issue13() string {
	someString := T("hello")
	fmt.Println(someString, T("world"))
	fmt.Println(someString, T("hello"), T("world"))

	fmt.Println(someString, T("Hello world {{.Arg0}}", map[string]interface{}{"Arg0": someString}))
	fmt.Println(someString, T("Hello world {{.Arg0}}", map[string]interface{}{"Arg0": fmt.Printf(T("my_world"))}))
	fmt.Println(someString, T("Hello world {{.Arg0}}", map[string]interface{}{"Arg0": T("my_world")}))

	return fmt.Sprint(10, T("world"))
}
