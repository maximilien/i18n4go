package doption

import (
	"fmt"
)

func DOption2() string {
	name := T("cruel")
	fmt.Printf(T("Bye from {{.Arg0}}", map[string]interface{}{"Arg0": T("Evil")}))

	for i := range 10 {
		fmt.Printf(T("Hello {{.Arg0}} world!", map[string]interface{}{"Arg0": name}))
	}
}
