package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
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

	type something struct {
	}

	foo := something{}
	strz := []string{"one", "two", "buckle my shoe"}
	fmt.Println("Welp, that's a great {{.MyStruct}} how about a {{.Whatever}}", &foo, strz[2])
}
