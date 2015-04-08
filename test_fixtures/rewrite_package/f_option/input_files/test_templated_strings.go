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
	fmt.Println("Hello {{.Name}} world!", "Evil")

	//These should not have a map[string]interface{}
	fmt.Println("Hello {{Not complex}} world! I am", name)
	fmt.Println("Hello {{}}", myName)

	fmt.Println("Hello {{.Name}} world!", strings.ToUpper(name))
	fmt.Println("Hello {{.Name}} world!", strings.ToUpper("Hi"))
	fmt.Println("Hello {{.Name}} world! {{.Number}} times", name, 10)
	fmt.Println("Hello {{.Name}} world! {{.Float}} times", name, 10.0)

	fmt.Println("Hello {{.Name}} world!", strings.ToUpper("Hello {{.Name}} world!", strings.ToUpper(name)))

	type something struct {
	}

	foo := something{}
	strz := []string{"one", "two", "buckle my shoe"}
	fmt.Println("Welp, that's a great {{.MyStruct}} how about a {{.Whatever}}", &foo, strz[2])

	println("Hello {{.Name}} world!", name)
	println("Hello {{.Name}} world! {{.Name}}", name, name)
}
