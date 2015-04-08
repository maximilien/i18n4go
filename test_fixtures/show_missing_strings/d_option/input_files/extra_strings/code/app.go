package app

func main() {
	println("I am a string")
	println(T("This should be verified."))
	println(T("This should be also {{.Verified}}", map[string]interface{}{
		"Verified": T("verified"),
	}))
}
