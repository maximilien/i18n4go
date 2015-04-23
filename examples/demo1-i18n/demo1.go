package main

import (
	"fmt"
)

const VERSION = "v0.0.1"

func main() {
	fmt.Println(T("Hello from Goffer land and i18n4go"))
	fmt.Println("")
	fmt.Printf(T("Version {{.Arg0}}\n", map[string]interface{}{"Arg0": VERSION}))
}
