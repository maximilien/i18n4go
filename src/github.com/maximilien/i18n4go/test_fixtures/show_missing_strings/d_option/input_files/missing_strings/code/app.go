package app

func main() {
	println(T("I am a missing string"))
	println(T("This should be verified."))
	println(T("This should be also {{.Verified}}", map[string]interface{}{
		"Verified": T("verified"),
	}))
}
