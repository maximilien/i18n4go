package doption

import (
	"fmt"
)

func DOption() string {
	name := T("cruel")
	fmt.Printf(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))
	printf(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))

	fmt.Printf(T("Bye from {{.Arg0}}", map[string]interface{}{"Arg0": T("Evil")}))
}
