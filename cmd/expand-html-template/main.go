package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/rothskeller/packet/v4/form/htmlop"
	"golang.org/x/net/html"
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
	if formHTML, err = htmlop.Parse(bytes.NewReader(formFile)); err != nil {
		panic("y")
	}
	// If there is a definitions.html in the same directory, read and parse
	// it too, and prepend it to the form HTML.
	if formFile, err = os.ReadFile(filepath.Join(filepath.Dir(filename), "definitions.html")); err == nil {
		var defHTML *html.Node
		if defHTML, err = htmlop.Parse(bytes.NewReader(formFile)); err != nil {
			panic("z")
		}
		for c := defHTML.LastChild; c != nil; c = defHTML.LastChild {
			defHTML.RemoveChild(c)
			formHTML.InsertBefore(c, formHTML.FirstChild)
		}
	}
	// Expand the templates in the form HTML, using the supplied fields.
	htmlop.Expand(formHTML, fields)
	// Render and minimize the result.
	htmlop.Minify(&formBuf, formHTML)
	os.Stdout.Write(formBuf.Bytes())
}
