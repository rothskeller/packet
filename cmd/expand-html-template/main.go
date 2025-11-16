package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/form/htmlop"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	var (
		filename string
		fields   map[string]string
		formFile []byte
		formHTML *html.Node
		formBuf  bytes.Buffer
		err      error
	)
	filename = os.Args[1]
	fields = make(map[string]string)
	for _, arg := range os.Args[2:] {
		name, val, _ := strings.Cut(arg, "=")
		fields[name] = val
	}
	// Read and parse the HTML for the form.
	if formFile, err = os.ReadFile(filename); err != nil {
		panic("x")
	}
	if formHTML, err = html.Parse(bytes.NewReader(formFile)); err != nil {
		panic("y")
	}
	// If there is a definitions.html in the same directory, read and parse
	// it too, and prepend it to the form HTML.
	if formFile, err = os.ReadFile(filepath.Join(filepath.Dir(filename), "definitions.html")); err == nil {
		var defHTML *html.Node
		if defHTML, err = html.Parse(bytes.NewReader(formFile)); err != nil {
			panic("z")
		}
		formBody := findBody(formHTML)
		defBody := findBody(defHTML)
		for c := defBody.LastChild; c != nil; c = defBody.LastChild {
			defBody.RemoveChild(c)
			formBody.InsertBefore(c, formBody.FirstChild)
		}
	}
	// Expand the templates in the form HTML, using the supplied fields.
	htmlop.Expand(formHTML, fields)
	// Render and minimize the result.
	htmlop.Minify(&formBuf, formHTML)
	os.Stdout.Write(formBuf.Bytes())
}

func findBody(doc *html.Node) *html.Node {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.Body {
			return n
		}
	}
	return nil
}
