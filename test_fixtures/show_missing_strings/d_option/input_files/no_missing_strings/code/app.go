package app

import "strings"

func main() {
	println("I am a string")
	println(T("This should be verified."))
	println(T("This should be also {{.Verified}}", map[string]interface{}{
		"Verified": T("verified"),
	}))
	println(T(strings.Join("a", "b"))) //ast.Expr is *ast.CallExpr, not *ast.BasicLit
}
