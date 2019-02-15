// Reference: https://hackernoon.com/how-a-go-program-compiles-down-to-machine-code-e4532dc8b8ca
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
)

func main() {
	src := []byte(`
	package main

	import "fmt"

	func main() {
		fmt.Println("Hello, world!")
	}
	`)

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		printer.Fprint(os.Stdout, fset, call.Fun)
		fmt.Println()

		return false
	})
}
