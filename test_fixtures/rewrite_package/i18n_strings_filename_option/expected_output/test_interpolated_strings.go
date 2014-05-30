package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := "cruel"
	myName := "evil"
	fmt.Println(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))
	fmt.Println("Hello %s world!, bye from %s", name, myName)

	fmt.Println(T("Hello {{.Arg0}}({{.Arg1}}) world!, bye from {{.Arg2}}", map[string]interface{}{"Arg0": 10, "Arg1": name, "Arg2": T("Evil")}))
}
