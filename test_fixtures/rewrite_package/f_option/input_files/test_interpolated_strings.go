package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := "cruel"
	myName := "evil"
	fmt.Println("Hello %s world!", name)
	fmt.Println("Hello %s world!, bye from %s", name, myName)

	fmt.Println("Hello %d(%s) world!, bye from %s", 10, name, "Evil")
}
