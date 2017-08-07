// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package decl_test

import (
	"testing"

	"github.com/shurcooL/godecl/decl"
)

func TestGoToEnglish(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			"",
			"",
		},
		{
			"var x int",
			"declare variable x as int",
		},
		{
			"var x, y int",
			"declare variables x and y as int",
		},
		{
			"var x int = 1",
			"declare variable x as int with initial value 1",
		},
		{
			"var x, y int = 1, 2",
			"declare variables x and y as int with initial values 1 and 2",
		},
		{
			"var x = 1",
			"declare variable x with initial value 1",
		},
		{
			"var x, y = 1, 2",
			"declare variables x and y with initial values 1 and 2",
		},
		{
			"var (x int; y string)",
			"declare variable x as int and variable y as string",
		},
		{
			"var (x, y int; a, b string)",
			"declare variables x and y as int and variables a and b as string",
		},
		{
			"x := 1",
			"short declare variable x with initial value 1",
		},
		{
			"x, y := 1, 2",
			"short declare variables x and y with initial values 1 and 2",
		},

		{
			"*[]map[int]string",
			"pointer to slice of map of int to string",
		},
		{
			"var x *[]map[int][2]string",
			"declare variable x as pointer to slice of map of int to 2-element array of string",
		},
		{
			"**[][]*map[int32][][3]string",
			"pointer to pointer to slice of slice of pointer to map of int32 to slice of 3-element array of string",
		},
		{
			"func(string, bool) (int, error)",
			"function taking string and bool and returning int and error",
		},
		{
			"var x, y int = (2+5) / 3, 4",
			"declare variables x and y as int with initial values (2 plus 5) divided by 3 and 4",
		},
		{
			"var x func() *[5]*func() rune",
			"declare variable x as function returning pointer to 5-element array of pointer to function returning rune",
		},

		{
			"i = 5",
			"assign to i the value 5",
		},
		{
			"i, j = 5, 6",
			"assign to i and j the values 5 and 6",
		},

		{
			`import "fmt"`,
			`import package "fmt"`,
		},
		{
			`import myfmt "fmt"`,
			`import package "fmt" as myfmt`,
		},
		{
			`import ("fmt"; "net/http"; _ "image/png")`,
			`import packages "fmt", "net/http" and "image/png" for side-effects`,
		},

		{
			"func Foo()",
			"function Foo",
		},
		{
			"func Foo() {}",
			"function Foo",
		},
		{
			"func Foo(x int) string",
			"function Foo taking x int and returning string",
		},
	}
	for _, tc := range tests {
		got, err := decl.GoToEnglish(tc.in)
		if err != nil {
			t.Error("got error:", err)
			continue
		}
		if got != tc.want {
			t.Errorf("\ngot:  %q\nwant: %q\n", got, tc.want)
		}
	}
}
