// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package decl implements functionality to convert
// fragments of Go code to an English representation.
package decl

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/shurcooL/go/parserutil"
)

// GoToEnglish returns a (possibly simplified) English representation
// for the fragment of Go code.
func GoToEnglish(frag string) (english string, err error) {
	var errors []string
	if expr, err := parser.ParseExpr(frag); err == nil {
		return exprString(expr), nil
	} else {
		errors = append(errors, "as an expression: "+err.Error())
	}
	if decl, err := parserutil.ParseDecl(frag); err == nil {
		return declString(decl), nil
	} else {
		errors = append(errors, "as a declaration: "+err.Error())
	}
	if stmt, err := parserutil.ParseStmt(frag); err == nil {
		return stmtString(stmt), nil
	} else {
		errors = append(errors, "as a statement: "+err.Error())
	}
	return "", fmt.Errorf("failed to parse fragment of Go code:\n%v", strings.Join(errors, "\n"))
}

// stmtString returns the (possibly simplified) English representation for x.
func stmtString(x ast.Stmt) string {
	var buf bytes.Buffer
	writeStmt(&buf, x)
	return buf.String()
}

// declString returns the (possibly simplified) English representation for x.
func declString(x ast.Decl) string {
	var buf bytes.Buffer
	writeDecl(&buf, x)
	return buf.String()
}

// exprString returns the (possibly simplified) English representation for x.
func exprString(x ast.Expr) string {
	var buf bytes.Buffer
	writeExpr(&buf, x)
	return buf.String()
}

// writeStmt writes the (possibly simplified) English representation for x to buf.
func writeStmt(buf *bytes.Buffer, x ast.Stmt) {
	switch x := x.(type) {
	default:
		fmt.Fprintf(buf, "<TODO: %T>", x)

	case *ast.EmptyStmt:
		// Do nothing.

	case *ast.DeclStmt:
		writeDecl(buf, x.Decl)

	case *ast.ExprStmt:
		writeExpr(buf, x.X)

	case *ast.AssignStmt:
		switch x.Tok {
		case token.DEFINE:
			buf.WriteString("short declare variable")
			if len(x.Lhs) > 1 {
				buf.WriteByte('s') // Plural.
			}
			buf.WriteByte(' ')
			for i, e := range x.Lhs {
				if i > 0 {
					buf.WriteString(" and ")
				}
				writeExpr(buf, e)
			}
			switch len(x.Rhs) {
			case 0: // Do nothing.
			case 1:
				buf.WriteString(" with initial value ")
			default:
				buf.WriteString(" with initial values ")
			}
			for i, e := range x.Rhs {
				if i > 0 {
					buf.WriteString(" and ")
				}
				writeExpr(buf, e)
			}
		case token.ASSIGN:
			buf.WriteString("assign to ")
			for i, e := range x.Lhs {
				if i > 0 {
					buf.WriteString(" and ")
				}
				writeExpr(buf, e)
			}
			switch len(x.Rhs) {
			case 0: // Do nothing.
			case 1:
				buf.WriteString(" the value ")
			default:
				buf.WriteString(" the values ")
			}
			for i, e := range x.Rhs {
				if i > 0 {
					buf.WriteString(" and ")
				}
				writeExpr(buf, e)
			}
		default:
			fmt.Fprintf(buf, "<TODO: '%v'>", x.Tok)
		}
	}
}

// writeDecl writes the (possibly simplified) English representation for x to buf.
func writeDecl(buf *bytes.Buffer, x ast.Decl) {
	switch x := x.(type) {
	default:
		fmt.Fprintf(buf, "<TODO: %T>", x)

	case *ast.GenDecl:
		switch x.Tok {
		default:
			buf.WriteString("declare ")
		case token.IMPORT:
			buf.WriteString("import package")
			if len(x.Specs) > 1 {
				buf.WriteByte('s') // Plural.
			}
			buf.WriteByte(' ')
		}
		for i, s := range x.Specs {
			writeSep(buf, i, len(x.Specs))
			switch x.Tok {
			case token.VAR:
				buf.WriteString("variable")
			case token.CONST:
				buf.WriteString("constant")
			case token.TYPE:
				buf.WriteString("type")
			case token.IMPORT:
				// Do nothing.
			default:
				fmt.Fprintf(buf, "<unexpected *ast.GenDecl.Tok: '%v'>", x.Tok)
			}
			if isValueSpecPlural(s) {
				buf.WriteByte('s') // Plural.
			}
			if x.Tok != token.IMPORT {
				buf.WriteByte(' ')
			}
			writeSpec(buf, s)
		}
	case *ast.FuncDecl:
		buf.WriteString("function ")
		buf.WriteString(x.Name.Name)
		if x.Type.Params.NumFields() > 0 || x.Type.Results.NumFields() > 0 {
			buf.WriteByte(' ')
		}
		writeSigExpr(buf, x.Type)
	}
}

// writeSep writes a separator that comes before entry i out of total.
// It creates English phrases like "first, second, third, forth and fifth".
func writeSep(buf *bytes.Buffer, i int, total int) {
	switch {
	case i > 0 && i < total-1:
		buf.WriteString(", ")
	case i > 0 && i == total-1:
		buf.WriteString(" and ")
	}
}

// isValueSpecPlural reports whether spec x is a value spec containing multiple names.
func isValueSpecPlural(x ast.Spec) bool {
	v, ok := x.(*ast.ValueSpec)
	if !ok {
		return false
	}
	return len(v.Names) > 1
}

// writeSpec writes the (possibly simplified) English representation for x to buf.
func writeSpec(buf *bytes.Buffer, x ast.Spec) {
	switch x := x.(type) {
	case *ast.ValueSpec:
		for i, n := range x.Names {
			if i > 0 {
				buf.WriteString(" and ")
			}
			buf.WriteString(n.Name)
		}

		if x.Type != nil {
			buf.WriteString(" as ")
			writeExpr(buf, x.Type)
		}

		switch len(x.Values) {
		case 0: // Do nothing.
		case 1:
			buf.WriteString(" with initial value ")
		default:
			buf.WriteString(" with initial values ")
		}
		for i, v := range x.Values {
			if i > 0 {
				buf.WriteString(" and ")
			}
			writeExpr(buf, v)
		}
	case *ast.TypeSpec:
		buf.WriteString(x.Name.Name)
		buf.WriteString(" as ")
		writeExpr(buf, x.Type)

	case *ast.ImportSpec:
		buf.WriteString(x.Path.Value)
		if x.Name == nil {
			break
		}
		switch x.Name.Name {
		default:
			buf.WriteString(" as ")
			buf.WriteString(x.Name.Name)
		case "_":
			buf.WriteString(" for side-effects")
		}

	default:
		fmt.Fprintf(buf, "<unexpected ast.Spec: %T>", x)
	}
}

// writeExpr writes the (possibly simplified) English representation for x to buf.
func writeExpr(buf *bytes.Buffer, x ast.Expr) {
	// The AST preserves source-level parentheses so there is
	// no need to introduce them here to correct for different
	// operator precedences. (This assumes that the AST was
	// generated by a Go parser.)

	switch x := x.(type) {
	default:
		buf.WriteString("(bad expr)") // nil, ast.BadExpr, ast.KeyValueExpr

	case *ast.Ident:
		buf.WriteString(x.Name)

	case *ast.Ellipsis:
		buf.WriteString("...")
		if x.Elt != nil {
			writeExpr(buf, x.Elt)
		}

	case *ast.BasicLit:
		buf.WriteString(x.Value)

	case *ast.FuncLit:
		buf.WriteByte('(')
		writeExpr(buf, x.Type)
		buf.WriteString(" literal)") // simplified

	case *ast.CompositeLit:
		buf.WriteByte('(')
		writeExpr(buf, x.Type)
		buf.WriteString(" literal)") // simplified

	case *ast.ParenExpr:
		buf.WriteByte('(')
		writeExpr(buf, x.X)
		buf.WriteByte(')')

	case *ast.SelectorExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('.')
		buf.WriteString(x.Sel.Name)

	case *ast.IndexExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('[')
		writeExpr(buf, x.Index)
		buf.WriteByte(']')

	case *ast.SliceExpr:
		writeExpr(buf, x.X)
		buf.WriteByte('[')
		if x.Low != nil {
			writeExpr(buf, x.Low)
		}
		buf.WriteByte(':')
		if x.High != nil {
			writeExpr(buf, x.High)
		}
		if x.Slice3 {
			buf.WriteByte(':')
			if x.Max != nil {
				writeExpr(buf, x.Max)
			}
		}
		buf.WriteByte(']')

	case *ast.TypeAssertExpr:
		writeExpr(buf, x.X)
		buf.WriteString(".(")
		writeExpr(buf, x.Type)
		buf.WriteByte(')')

	case *ast.CallExpr:
		writeExpr(buf, x.Fun)
		buf.WriteByte('(')
		for i, arg := range x.Args {
			if i > 0 {
				buf.WriteString(", ")
			}
			writeExpr(buf, arg)
		}
		if x.Ellipsis.IsValid() {
			buf.WriteString("...")
		}
		buf.WriteByte(')')

	case *ast.StarExpr:
		buf.WriteString("pointer to ")
		writeExpr(buf, x.X)

	case *ast.UnaryExpr:
		buf.WriteString(x.Op.String())
		writeExpr(buf, x.X)

	case *ast.BinaryExpr:
		writeExpr(buf, x.X)
		buf.WriteByte(' ')
		switch x.Op {
		default:
			buf.WriteString(x.Op.String())
		case token.ADD:
			buf.WriteString("plus")
		case token.SUB:
			buf.WriteString("minus")
		case token.QUO:
			buf.WriteString("divided by")
		}
		buf.WriteByte(' ')
		writeExpr(buf, x.Y)

	case *ast.ArrayType:
		if x.Len == nil {
			buf.WriteString("slice of ")
		} else {
			writeExpr(buf, x.Len)
			buf.WriteString("-element array of ")
		}
		writeExpr(buf, x.Elt)

	case *ast.StructType:
		buf.WriteString("struct{")
		writeFieldList(buf, x.Fields, "; ", false)
		buf.WriteByte('}')

	case *ast.FuncType:
		buf.WriteString("function")
		if x.Params.NumFields() > 0 || x.Results.NumFields() > 0 {
			buf.WriteByte(' ')
		}
		writeSigExpr(buf, x)

	case *ast.InterfaceType:
		buf.WriteString("interface{")
		writeFieldList(buf, x.Methods, "; ", true)
		buf.WriteByte('}')

	case *ast.MapType:
		buf.WriteString("map of ")
		writeExpr(buf, x.Key)
		buf.WriteString(" to ")
		writeExpr(buf, x.Value)

	case *ast.ChanType:
		var s string
		switch x.Dir {
		case ast.SEND:
			s = "chan<- "
		case ast.RECV:
			s = "<-chan "
		default:
			s = "chan "
		}
		buf.WriteString(s)
		writeExpr(buf, x.Value)
	}
}

func writeSigExpr(buf *bytes.Buffer, sig *ast.FuncType) {
	if sig.Params.NumFields() > 0 {
		buf.WriteString("taking ")
		writeFieldList(buf, sig.Params, " and ", false)
	}
	if sig.Params.NumFields() > 0 && sig.Results.NumFields() > 0 {
		buf.WriteString(" and ")
	}
	if sig.Results.NumFields() > 0 {
		buf.WriteString("returning ")
		writeFieldList(buf, sig.Results, " and ", false)
	}
}

func writeFieldList(buf *bytes.Buffer, fields *ast.FieldList, sep string, iface bool) {
	for i, f := range fields.List {
		if i > 0 {
			buf.WriteString(sep)
		}

		// Field list names.
		for i, name := range f.Names {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(name.Name)
		}

		// Types of interface methods consist of signatures only.
		if sig, _ := f.Type.(*ast.FuncType); sig != nil && iface {
			writeSigExpr(buf, sig)
			continue
		}

		// Named fields are separated with a blank from the field type.
		if len(f.Names) > 0 {
			buf.WriteByte(' ')
		}

		writeExpr(buf, f.Type)

		// Ignore tag.
	}
}
