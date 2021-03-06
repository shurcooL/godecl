// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"go/format"
	"testing"
)

// Ensure all examples are nicely gofmted.
func TestExamplesGofmt(t *testing.T) {
	for _, example := range examples {
		gofmted, err := format.Source([]byte(example))
		if err != nil {
			t.Errorf("failed to gofmt example %q:\n%v", example, err)
			continue
		}
		if example != string(gofmted) {
			t.Errorf("\nexample %q is not gofmted\n  want: %q", example, gofmted)
		}
	}
}
