package input_files

import (
	"fmt"
	"strings"
)

func Interpolated() string {
	name := "cruel"
	myName := "evil"
	fmt.Printf("Hello %s world!", name)
	fmt.Printf("Hello %s world!, bye from %s", name, myName)

	fmt.Printf("Hello %d(%s) world!, bye from %s", 10, name, "Evil")

	fmt.Printf("Hello %s world!", name)
	fmt.Printf("Hello %s world! %s", name, name)
}
