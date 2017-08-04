// A godecl experiment. Like cdecl, but for Go.
//
// Inspired by @bradfitz at https://twitter.com/bradfitz/status/833048466456600576.
package main

import (
	"fmt"
	"os"

	"github.com/gopherjs/gopherjs/js"
	"github.com/shurcooL/play/228/godecl/decl"
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
	}
	c.input.AddEventListener("input", false, func(dom.Event) { c.OnInput() })
	c.OnInput()
}

type context struct {
	input  *dom.HTMLInputElement
	output *dom.HTMLDivElement
}

func (c context) OnInput() {
	out, err := decl.GoToEnglish(c.input.Value)
	if err != nil {
		c.output.SetTextContent("error: " + err.Error())
		return
	}
	c.output.SetTextContent(out)
}
