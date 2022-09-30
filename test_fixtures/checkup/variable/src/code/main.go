package code

import "fmt"

func main() {
	locale := "Translated hello world!"
	fmt.Println(T(locale))
	locale = "I like bananas"
	fmt.Println(T(locale))
}
