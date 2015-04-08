package app

import (
	"fmt"
)

func Issue4() {
	description := "Unsure how hard this is to fix, but I noticed that > would be translated to its unicode/ascii value."
	moreInfo := "I am assuming this means <> would not be correctly extracted as well."

	fmt.Println("GitHub: issue #4")
	fmt.Println(description)
	fmt.Println(moreInfo)
}
