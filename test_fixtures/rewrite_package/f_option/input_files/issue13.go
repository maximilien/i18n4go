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
	fmt.Println(someString, "Hello world %s", fmt.Printf("my_world"))
	fmt.Println(someString, "Hello world %s", "my_world")

	return fmt.Sprint(10, "world")
}
