package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := T("cruel")
	myName := T("evil")
	fmt.Println(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))
	fmt.Println(T("Hello {{.Arg0}} world!, bye from {{.Arg1}}", map[string]interface{}{"Arg0": name, "Arg1": myName}))

	fmt.Println(T("Hello {{.Arg0}}({{.Arg1}}) world!, bye from {{.Arg2}}", map[string]interface{}{"Arg0": 10, "Arg1": name, "Arg2": T("Evil")}))

	println(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))
}
