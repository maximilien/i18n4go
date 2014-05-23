package input_files

import (
	"fmt"
	"os"
	"path/filepath"
)

func Something() {
	fmt.Printf("HAI")
	if os.Getenv("SOMETHING") {
		fmt.Printf(filepath.Clean(os.Getenv("something")))
	}
}
