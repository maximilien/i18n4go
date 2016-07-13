package main

import "fmt"

func main() {
	s := "a string"
	fmt.Println(s)
	newStruct := mystruct{
		mystring: s,
	}
	fmt.Println(newStruct.mystring)
}

type mystruct struct {
	mystring string `command:"app" description:"show the app details"`
}
