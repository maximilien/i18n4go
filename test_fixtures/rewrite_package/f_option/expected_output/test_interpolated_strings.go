package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := T("cruel")
	myName := T("evil")
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": name}))
	fmt.Println(T("Hello {{.Name}} world!, bye from {{.MyName}}", map[string]interface{}{"Name": name, "MyName": myName}))

	fmt.Println(T("Hello {{Not complex}} world! I am"), name)
	fmt.Println(T("Hello {{}}"), myName)

	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)}))
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(T("Hi"))}))
	fmt.Println(T("Hello {{.Name}} world! {{.Number}} times", map[string]interface{}{"Name": name, "Number": 10}))

	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)}))}))
}
