// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"

	"github.com/shurcooL/go/parserutil"
	"github.com/shurcooL/go/printerutil"
)

func main() {
	/*out, err := run2(input.Value)
	if err != nil {
		output.SetTextContent("error: " + err.Error())
		return
	}
	output.SetTextContent(out)

	out, err = run3(input.Value)
	if err != nil {
		output.SetTextContent("error: " + err.Error())
		return
	}
	output.SetTextContent(output.TextContent() + "\n" + out)

	out, err := run4(input.Value)
	if err != nil {
		output.SetTextContent("error: " + err.Error())
		return
	}
	output.SetTextContent(output.TextContent() + "\n" + out)*/
}

func run2(s string) (string, error) {
	decl, err := parserutil.ParseDecl(s)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = ast.Fprint(&buf, nil, decl, ast.NotNilFilter)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func run3(x string) (string, error) {
	decl, err := parserutil.ParseDecl(x)
	if err != nil {
		return "", err
	}
	var s string
	ast.Inspect(decl, func(n ast.Node) bool {
		if n != nil {
			s += fmt.Sprintf("%T: %v\n", n, printerutil.SprintAstBare(n))
		}
		return true
	})
	return s, nil
}

func goToEnglish(x string) (string, error) {
	expr, err := parser.ParseExpr(x)
	if err != nil {
		return "", err
	}
	return ExprString(expr), nil
}

// ExprString returns the (possibly simplified) string representation for x.
func ExprString(x ast.Expr) string {
	var buf bytes.Buffer
	WriteExpr(&buf, x)
	return buf.String()
}

/*var typeString func(t map[string]interface{}) string
typeString = func(t map[string]interface{}) string {
	switch t["kind"] {
	case "NON_NULL":
		s := typeString(t["ofType"].(map[string]interface{}))
		if !strings.HasPrefix(s, "*") {
			panic(fmt.Errorf("nullable type %q doesn't begin with '*'", s))
		}
		return s[1:] // Strip star from nullable type to make it non-null.
	case "LIST":
		return "*[]" + typeString(t["ofType"].(map[string]interface{}))
	default:
		return "*" + t["name"].(string)
	}
}*/
