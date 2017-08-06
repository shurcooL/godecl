// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A godecl experiment. Like cdecl, but for Go.
//
// Inspired by @bradfitz at https://twitter.com/bradfitz/status/833048466456600576.
package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gopherjs/gopherjs/js"
	"github.com/shurcooL/godecl/decl"
	"honnef.co/go/js/dom"
)

func main() {
	if js.Global == nil || js.Global.Get("document") == js.Undefined {
		fmt.Fprintln(os.Stderr, `Where's the DOM at? It looks like you're running godecl in an environment
where the DOM is not available. You'll need to run it inside a browser.`)
		os.Exit(1)
	}

	document := dom.GetWindow().Document()
	c := context{
		input:  document.GetElementByID("input").(*dom.HTMLInputElement),
		output: document.GetElementByID("output").(*dom.HTMLDivElement),
		issue:  document.GetElementByID("issue").(*dom.HTMLAnchorElement),
	}
	c.input.AddEventListener("input", false, func(dom.Event) { c.OnInput() })
	c.OnInput()
}

type context struct {
	input  *dom.HTMLInputElement
	output *dom.HTMLDivElement
	issue  *dom.HTMLAnchorElement
}

func (c context) OnInput() {
	out, err := decl.GoToEnglish(c.input.Value)
	if err != nil {
		c.output.SetTextContent("error: " + err.Error())
		c.updateIssueHref()
		return
	}
	c.output.SetTextContent(out)
	c.updateIssueHref()
}

func (c context) updateIssueHref() {
	v := url.Values{}
	v.Set("title", fmt.Sprintf("decl: Unexpected handling of %q.", c.input.Value))
	v.Set("body", fmt.Sprintf(`### What did you do?

I typed the following input at godecl.org:

`+"```Go"+`
%v
`+"```"+`

### What did you expect to see?

I expected to see ...

### What did you see instead?

`+"```"+`
%v
`+"```"+`
`, c.input.Value, c.output.TextContent()))
	url := url.URL{
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/shurcooL/godecl/issues/new",
		RawQuery: v.Encode(),
	}
	c.issue.Href = url.String()
}
