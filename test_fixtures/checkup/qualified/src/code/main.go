package code

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
)

func main() {
	fmt.Println(i18n.T("Translated hello world!"))
}
