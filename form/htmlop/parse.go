package htmlop

import (
	"io"
	"slices"
	"strings"

	"github.com/rothskeller/packet/v4/errors"
	"golang.org/x/net/html"
)

// Parse parses the HTML supplied by r into a DOM tree.  Unlike html.Parse, it
// does not do any tree manipulation; the node tree exactly matches the input.
// As a consequence, it does not accept input with omitted end tags.  It also
// keeps repeated attributes, which templates use to override a default value
// (e.g., value="Submit Form" value={{submit-label}}).
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
		var raw []byte
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			raw = slices.Clone(tz.Raw()) // Token modifies the buffer
		}
		tk := tz.Token()
		if raw != nil {
			// The tokenizer drops repeated attributes, as the HTML
			// spec requires.  Recover them from the raw tag text.
			if attrs := tagAttributes(raw); len(attrs) > len(tk.Attr) {
				tk.Attr = attrs
			}
		}
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

// tagAttributes returns all of the attributes, including repeated ones, of the
// start tag whose complete source text is raw.
func tagAttributes(raw []byte) (attrs []html.Attribute) {
	s := string(raw)
	i := 1 // skip "<"
	for i < len(s) && !isSpace(s[i]) && s[i] != '/' && s[i] != '>' {
		i++ // skip tag name
	}
	for {
		for i < len(s) && (isSpace(s[i]) || s[i] == '/') {
			i++
		}
		if i >= len(s) || s[i] == '>' {
			return attrs
		}
		start := i
		i++ // an attribute name can start with "="
		for i < len(s) && !isSpace(s[i]) && s[i] != '/' && s[i] != '>' && s[i] != '=' {
			i++
		}
		attr := html.Attribute{Key: strings.ToLower(s[start:i])}
		for i < len(s) && isSpace(s[i]) {
			i++
		}
		if i < len(s) && s[i] == '=' {
			i++
			for i < len(s) && isSpace(s[i]) {
				i++
			}
			if i < len(s) && (s[i] == '"' || s[i] == '\'') {
				quote := s[i]
				i++
				start = i
				for i < len(s) && s[i] != quote {
					i++
				}
				attr.Val = s[start:i]
				i++
			} else {
				start = i
				for i < len(s) && !isSpace(s[i]) && s[i] != '>' {
					i++
				}
				attr.Val = s[start:i]
			}
			attr.Val = html.UnescapeString(strings.ReplaceAll(strings.ReplaceAll(attr.Val, "\r\n", "\n"), "\r", "\n"))
		}
		attrs = append(attrs, attr)
	}
}

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\f' || c == '\r' }
