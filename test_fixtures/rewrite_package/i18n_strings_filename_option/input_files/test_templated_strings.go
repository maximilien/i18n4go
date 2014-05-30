package input_files

import (
	"fmt"
	"strings"
)

func Templated() string {
	name := "cruel"
	myName := "evil"
	fmt.Println("Hello {{.Name}} world!", name)
	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", name, myName)

	fmt.Println("Hello {{Not complex}} world! I am", name)
	fmt.Println("Hello {{}}", myName)

	fmt.Println("Hello {{.Name}} world!", strings.ToUpper(name))
	fmt.Println("Hello {{.Name}} world!", strings.ToUpper("Hi"))
	fmt.Println("Hello {{.Name}} world! {{.Number}} times", name, 10)

	fmt.Println("Hello {{.Name}} world!", strings.ToUpper("Hello {{.Name}} world!", strings.ToUpper(name)))
	fmt.Println("Hello {{.Name}} world!", strings.ToUpper("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), myName))

	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), strings.ToUpper("Hello {{.Name}} world!", strings.ToUpper(name)))
	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), strings.ToUpper("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), myName))
}
