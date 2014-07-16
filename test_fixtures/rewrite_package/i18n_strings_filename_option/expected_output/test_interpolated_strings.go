package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := "cruel"
	myName := "evil"
	fmt.Printf(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))
	fmt.Printf(T("Bye {{.Arg0}} world!\n", map[string]interface{}{"Arg0": name}))
	fmt.Printf("Hello %s world!, bye from %s", name, myName)
	fmt.Printf(T("Hello again:\t {{.Arg0}} world!\n", map[string]interface{}{"Arg0": name}))

	fmt.Printf(T("Hello {{.Arg0}}({{.Arg1}}) world!, bye from {{.Arg2}}", map[string]interface{}{"Arg0": 10, "Arg1": name, "Arg2": T("Evil")}))
}
