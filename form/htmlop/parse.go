package htmlop

import (
	"io"

	"github.com/rothskeller/packet/v4/errors"
	"golang.org/x/net/html"
)

// Parse parses the HTML supplied by r into a DOM tree.  Unlike html.Parse, it
// does not do any tree manipulation; the node tree exactly matches the input.
// As a consequence, it does not accept input with omitted end tags.
func Parse(r io.Reader) (doc *html.Node, err error) {
	var tz = html.NewTokenizer(r)
	doc = &html.Node{Type: html.DocumentNode}
	var ptr = doc
	for {
		tt := tz.Next()
		if tt == html.ErrorToken {
			if tz.Err() != io.EOF {
				return nil, err
			}
			if ptr != doc {
				return nil, errors.NewF("unexpected EOF: missing end tag </%s>", ptr.Data)
			}
			return doc, nil
		}
		tk := tz.Token()
		switch tt {
		case html.CommentToken:
			ptr.AppendChild(&html.Node{Type: html.CommentNode, Data: tk.Data})
		case html.DoctypeToken:
			ptr.AppendChild(&html.Node{Type: html.DoctypeNode, Data: tk.Data})
		case html.EndTagToken:
			if tk.Data == ptr.Data {
				ptr = ptr.Parent
			} else {
				return nil, errors.NewF("unexpected </%s>, expecting </%s> %d", tk.Data, ptr.Data, ptr.DataAtom)
			}
		case html.SelfClosingTagToken:
			ptr.AppendChild(&html.Node{Type: html.ElementNode, Data: tk.Data, DataAtom: tk.DataAtom, Attr: tk.Attr})
		case html.StartTagToken:
			n := &html.Node{Type: html.ElementNode, Data: tk.Data, DataAtom: tk.DataAtom, Attr: tk.Attr}
			ptr.AppendChild(n)
			if !voidElements.Has(n.DataAtom) {
				ptr = n
			}
		case html.TextToken:
			ptr.AppendChild(&html.Node{Type: html.TextNode, Data: tk.Data})
		default:
			panic("unexpected token type")
		}
	}
}
