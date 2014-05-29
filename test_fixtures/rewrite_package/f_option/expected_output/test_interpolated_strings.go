package input_files

import (
	"fmt"
)

func Interpolated() string {
	name := T("cruel")
	myName := T("evil")
	fmt.Println(T("Hello {{.Name}} world!", map[string]interface{}{"Name": name}))
	fmt.Println(T("Hello {{.Name}} world!, bye from {{.MyName}}", map[string]interface{}{"Name": name, "MyName": myName}))

	fmt.Println(T("Hello {{Not complex}} world! I am"), name)
	fmt.Println(T("Hello {{}}"), myName)
}
