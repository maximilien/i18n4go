package app

import (
	"errors"
	"fmt"
)

func main() {
	translatedString := "translated"

	err := test()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = errors.New(T("This is another string which has been {{.Translated}}.", map[string]interface{}{"Translated": translatedString}))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func test49() {
	if true {
		return errors.New(T("This is a string which has been {{.Translated}}.", map[string]interface{}{"Translated": translatedString}))
	}
}
