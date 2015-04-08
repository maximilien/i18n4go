package input_files

import (
	"fmt"
	"strings"
)

func Templated() string {
	name := T("cruel")
	myName := "evil"
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": name}))
	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", name, myName)

	fmt.Println(T("Hello {{Not complex}} world! I am"), name)
	fmt.Println("Hello {{}}", myName)

	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)}))
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper("Hi")}))
	fmt.Println(T("Hello {{.Name}} world! {{.Number}} times", map[string]interface{}{"Name": name, "Number": 10}))

	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)}))}))
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), myName)}))

	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), strings.ToUpper(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)})))
	fmt.Println("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), strings.ToUpper("Hello {{.Name}} world!, bye from {{.MyName}}", strings.ToUpper(name), myName))
}
