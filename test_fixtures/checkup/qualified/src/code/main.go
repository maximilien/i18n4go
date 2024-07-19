package code

import (
	"fmt"

	i18n "github.com/maximilien/i18n4go/i18n4go/cmds"
)

func main() {
	fmt.Println(i18n.T("Translated hello world!"))
}
