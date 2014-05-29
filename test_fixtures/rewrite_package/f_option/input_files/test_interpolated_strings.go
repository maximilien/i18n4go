package input_files

import (
	"fmt"
)

func Interpolated() string {
	name := "cruel"
	myName := "evil"
	fmt.Println("Hello {{.Name}} world!", name)
	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", name, myName)

	fmt.Println("Hello {{Not complex}} world! I am", name)
	fmt.Println("Hello {{}}", myName)
}
