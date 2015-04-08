package input_files

import (
	"fmt"
	"strings"
)

func Templated() string {
	name := T("cruel")
	myName := T("evil")
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": name}))
	fmt.Println(T("Hello {{.Name}} world!, bye from {{.MyName}}", map[string]interface{}{"Name": name, "MyName": myName}))
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": T("Evil")}))

	//These should not have a map[string]interface{}
	fmt.Println(T("Hello {{Not complex}} world! I am"), name)
	fmt.Println(T("Hello {{}}"), myName)

	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)}))
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(T("Hi"))}))
	fmt.Println(T("Hello {{.Name}} world! {{.Number}} times", map[string]interface{}{"Name": name, "Number": 10}))
	fmt.Println(T("Hello {{.Name}} world! {{.Float}} times", map[string]interface{}{"Name": name, "Float": 10.0}))

	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(T("Hello {{.Name}} world!", map[string]interface{}{"Name": strings.ToUpper(name)}))}))

	type something struct {
	}

	foo := something{}
	strz := []string{T("one"), T("two"), T("buckle my shoe")}
	fmt.Println(T("Welp, that's a great {{.MyStruct}} how about a {{.Whatever}}", map[string]interface{}{"MyStruct": &foo, "Whatever": strz[2]}))

	println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": name}))
	println(T("Hello {{.Name}} world! {{.Name}}", map[string]interface{}{"Name": name}))
}
