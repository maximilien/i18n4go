package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

func Issue13() string {
	someString := "hello"
	fmt.Println(someString, "world")
	fmt.Println(someString, "hello", "world")

	fmt.Println(someString, "Hello world %s", someString)

	return fmt.Sprint(10, "world")
}
