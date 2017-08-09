// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A godecl experiment. Like cdecl, but for Go.
//
// Inspired by @bradfitz at https://twitter.com/bradfitz/status/833048466456600576.
package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/shurcooL/godecl/decl"
	"honnef.co/go/js/dom"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if js.Global == nil || js.Global.Get("document") == js.Undefined {
		fmt.Fprintln(os.Stderr, `Where's the DOM at? It looks like you're running godecl in an environment
where the DOM is not available. You'll need to run it inside a browser.`)
		os.Exit(1)
	}

	document := dom.GetWindow().Document()
	c := context{
		input:     document.GetElementByID("input").(*dom.HTMLInputElement),
		output:    document.GetElementByID("output").(*dom.HTMLDivElement),
		permalink: document.GetElementByID("permalink").(*dom.HTMLAnchorElement),
		issue:     document.GetElementByID("issue").(*dom.HTMLAnchorElement),
	}
	c.SetInitialInput()
	c.Update()
	c.input.AddEventListener("input", false, func(dom.Event) {
		deleteQuery()
		c.Update()
	})
	c.permalink.AddEventListener("click", false, func(e dom.Event) {
		me := e.(*dom.MouseEvent)
		if me.CtrlKey || me.AltKey || me.MetaKey || me.ShiftKey {
			// Only override normal clicks.
			return
		}
		setQuery(c.input.Value)
		e.PreventDefault()
	})
	document.GetElementByID("random").AddEventListener("click", false, func(e dom.Event) {
		me := e.(*dom.MouseEvent)
		if me.CtrlKey || me.AltKey || me.MetaKey || me.ShiftKey {
			// Only override normal clicks.
			return
		}
		deleteQuery()
		// Pick a random example (but avoid picking same as current input).
		for {
			new := examples[rand.Intn(len(examples))]
			if new == c.input.Value {

				continue
			}
			c.input.Value = new
			break
		}
		c.Update()
		e.PreventDefault()
	})
}

type context struct {
	input     *dom.HTMLInputElement
	output    *dom.HTMLDivElement
	permalink *dom.HTMLAnchorElement
	issue     *dom.HTMLAnchorElement
}

// SetInitialInput sets initial input value.
func (c context) SetInitialInput() {
	if c.input.Value != "" {
		// If the user has managed to already type some input, don't steamroll over it.
		return
	}
	query, _ := url.ParseQuery(strings.TrimPrefix(dom.GetWindow().Location().Search, "?"))
	q, ok := query["q"]
	if !ok {
		// Random initial example.
		c.input.Value = examples[rand.Intn(len(examples))]
		c.input.Focus()
		return
	}
	c.input.Value = q[0]
	c.input.Focus()
}

// Update updates the output, permalink anchor href,
// and issue anchor href, based on current input.
func (c context) Update() {
	c.updateOutput()
	c.updatePermalink()
	c.updateIssue()
}

// updateOutput updates the output text based on current input.
func (c context) updateOutput() {
	out, err := decl.GoToEnglish(c.input.Value)
	if err != nil {
		c.output.SetTextContent("error: " + err.Error())
		return
	}
	c.output.SetTextContent(out)
}

// updatePermalink updates the permalink anchor href based on current input.
func (c context) updatePermalink() {
	v := url.Values{}
	v.Set("q", c.input.Value)
	url := url.URL{RawQuery: v.Encode()}
	c.permalink.Href = url.String()
}

// updateIssue updates the "report an issue" anchor href based on current input.
func (c context) updateIssue() {
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
		Scheme: "https", Host: "github.com", Path: "/shurcooL/godecl/issues/new",
		RawQuery: v.Encode(),
	}
	c.issue.Href = url.String()
}

// setQuery sets q in the window URL query to value.
func setQuery(value string) {
	url, err := url.Parse(dom.GetWindow().Location().Href)
	if err != nil {
		// We don't expect this can ever happen, so treat it as an internal error if it does.
		panic(fmt.Errorf("internal error: parsing window.location.href as URL failed: %v", err))
	}
	query := url.Query()
	query.Set("q", value)
	url.RawQuery = query.Encode()
	// TODO: dom.GetWindow().History().ReplaceState(...), blocked on https://github.com/dominikh/go-js-dom/issues/41.
	js.Global.Get("window").Get("history").Call("replaceState", nil, nil, url.String())
}

// deleteQuery deletes q in the window URL query.
func deleteQuery() {
	url, err := url.Parse(dom.GetWindow().Location().Href)
	if err != nil {
		// We don't expect this can ever happen, so treat it as an internal error if it does.
		panic(fmt.Errorf("internal error: parsing window.location.href as URL failed: %v", err))
	}
	query := url.Query()
	query.Del("q")
	url.RawQuery = query.Encode()
	// TODO: dom.GetWindow().History().ReplaceState(...), blocked on https://github.com/dominikh/go-js-dom/issues/41.
	js.Global.Get("window").Get("history").Call("replaceState", nil, nil, url.String())
}
