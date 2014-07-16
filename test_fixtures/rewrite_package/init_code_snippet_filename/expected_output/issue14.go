package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

func Issue14() string {
	someString := T("hello issue 14")
	fmt.Println(someString)
}
