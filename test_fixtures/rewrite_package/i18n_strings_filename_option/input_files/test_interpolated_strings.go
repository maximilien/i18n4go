package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := "cruel"
	myName := "evil"
	fmt.Println("Hello %s world!", name)
	fmt.Printf("Bye %s world!\n", name)
	fmt.Println("Hello %s world!, bye from %s", name, myName)
	fmt.Printf("Hello again:\t %s world!\n", name)

	fmt.Println("Hello %d(%s) world!, bye from %s", 10, name, "Evil")
}
