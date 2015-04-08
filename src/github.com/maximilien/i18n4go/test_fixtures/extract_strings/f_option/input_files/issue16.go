package app

import "fmt"

func main() {
	translatedString := "translated"

	T("This is a string which has been {{.Translated}}.", map[string]interface{}{"Translated": translatedString})

	inputs := map[string]string{
		"test1": "foo",
	}

	fmt.Println(inputs["test1"])
}
