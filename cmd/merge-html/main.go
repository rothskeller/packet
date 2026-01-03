// merge-html merges multiple HTML files into one, merging the <head>, <script>,
// <style>, and <body> sections separately.  The result is minified.
//
// usage: merge-html inputfile... > outputfile
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	minhtml "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
)

var heads, scripts, styles, bodies, escripts *html.Node

func main() {
	heads = &html.Node{Type: html.ElementNode, Data: "head", DataAtom: atom.Head}
	scripts = &html.Node{Type: html.ElementNode, Data: "script", DataAtom: atom.Script}
	scripts.AppendChild(&html.Node{Type: html.TextNode})
	styles = &html.Node{Type: html.ElementNode, Data: "style", DataAtom: atom.Style}
	styles.AppendChild(&html.Node{Type: html.TextNode})
	bodies = &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	escripts = &html.Node{Type: html.ElementNode, Data: "script", DataAtom: atom.Script}
	escripts.AppendChild(&html.Node{Type: html.TextNode})
	for _, fname := range os.Args[1:] {
		if fh, err := os.Open(fname); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		} else if doc, err := html.Parse(fh); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", fname, err)
			os.Exit(1)
		} else {
			extractFrom(doc)
		}
	}
	writeResult()
}

func extractFrom(doc *html.Node) {
	var next *html.Node

	switch doc.Type {
	case html.ElementNode:
		switch doc.DataAtom {
		case atom.Head:
			extractFromHead(doc)
			return
		case atom.Body:
			extractFromBody(doc)
			return
		}
	}
	for c := doc.FirstChild; c != nil; c = next {
		next = c.NextSibling
		extractFrom(c)
	}
}

func extractFromHead(head *html.Node) {
	var next *html.Node

	for c := head.FirstChild; c != nil; c = next {
		next = c.NextSibling
		switch c.Type {
		case html.ElementNode:
			switch c.DataAtom {
			case atom.Script:
				appendText(scripts, c)
			case atom.Style:
				appendText(styles, c)
			default:
				moveUnder(c, heads)
			}
		default:
			moveUnder(c, heads)
		}
	}
}

func extractFromBody(body *html.Node) {
	var next *html.Node

	for c := body.FirstChild; c != nil; c = next {
		next = c.NextSibling
		if c.Type == html.ElementNode && c.DataAtom == atom.Script {
			appendText(escripts, c)
		} else {
			moveUnder(c, bodies)
		}
	}
}

func moveUnder(from, to *html.Node) {
	if from.Parent != nil {
		from.Parent.RemoveChild(from)
	}
	to.AppendChild(from)
}

func appendText(to, from *html.Node) {
	if from.FirstChild == nil || from.FirstChild.Type != html.TextNode || from.LastChild != from.FirstChild {
		fmt.Fprintln(os.Stderr, "ERROR: script or style contains something other than text")
		os.Exit(1)
	}
	to.FirstChild.Data += from.FirstChild.Data
}

func writeResult() {
	doc := &html.Node{Type: html.DocumentNode}
	doc.AppendChild(&html.Node{Type: html.DoctypeNode, Data: "html"})
	top := &html.Node{Type: html.ElementNode, Data: "html", DataAtom: atom.Html}
	doc.AppendChild(top)
	top.AppendChild(heads)
	heads.AppendChild(scripts)
	heads.AppendChild(styles)
	top.AppendChild(bodies)
	bodies.AppendChild(escripts)
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", minhtml.Minify)
	m.AddFunc("application/javascript", js.Minify)
	mw := m.Writer("text/html", os.Stdout)
	html.Render(mw, doc)
	mw.Close()
}
