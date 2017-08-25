// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// examples is a list of example input values. On page load,
// one of these is randomly chosen to display to the user.
// These should be gofmted, to keep the examples nicer.
var examples = []string{
	"var x *[]map[int][2]string",
	"var x func() *[5]*func() rune",
	"var x, y int = 1, 2",
	"var x = (2+5)/3.0 + 4",
	"type T1 = T2",

	// TODO: Add more fun and interesting example inputs.
	//       See decl tests for inspiration.
}
