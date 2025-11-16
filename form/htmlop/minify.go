package htmlop

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Minify writes the supplied HTML node tree to the supplied writer, in a
// minimized form.
func Minify(w io.Writer, n *html.Node) (err error) {
	if n == nil {
		return nil
	}
	switch n.Type {
	case html.DocumentNode:
		for c := range n.ChildNodes() {
			if err = Minify(w, c); err != nil {
				break
			}
		}
	case html.DoctypeNode:
		_, err = fmt.Fprintf(w, "<!DOCTYPE %s>", n.Data)
	case html.TextNode:
		_, err = io.WriteString(w, minifyStringReplacer.Replace(n.Data))
	case html.CommentNode:
		err = nil
	default:
		err = minifyElement(w, n)
	}
	return err
}

var minifyStringReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
)

// minifyElement minifies an element and its descendants.
func minifyElement(w io.Writer, n *html.Node) (err error) {
	if _, err = io.WriteString(w, startTag(n)); err != nil {
		return err
	}
	for c := range n.ChildNodes() {
		if err = Minify(w, c); err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, endTag(n))
	return err
}

// startTag returns the start tag for an element, including its attributes.
func startTag(n *html.Node) string {
	var sb strings.Builder

	if alwaysOmit.Has(n.DataAtom) {
		return ""
	}
	sb.WriteByte('<')
	sb.WriteString(n.Data)
	for _, a := range n.Attr {
		minifyAttribute(&sb, a.Key, a.Val)
	}
	sb.WriteByte('>')
	return sb.String()
}

func minifyAttribute(sb *strings.Builder, name, val string) {
	sb.WriteByte(' ')
	sb.WriteString(name)
	if val != "" {
		sb.WriteByte('=')
		val = minifyAttrReplacer.Replace(val)
		if strings.ContainsAny(val, "\t\n\f\r '=<>`") || strings.HasSuffix(val, "/") {
			sb.WriteByte('"')
			sb.WriteString(val)
			sb.WriteByte('"')
		} else {
			sb.WriteString(val)
		}
	}
}

var minifyAttrReplacer = strings.NewReplacer(
	"&", "&amp;",
	"\"", "&quot;",
)

func endTag(n *html.Node) string {
	if alwaysOmit.Has(n.DataAtom) {
		return ""
	}
	if n.FirstChild == nil && voidElements.Has(n.DataAtom) {
		return ""
	}
	if ns := nextSiblingSkipWS(n); ns != nil && startImpliesClose.Has(elementPair{n.DataAtom, ns.DataAtom}) {
		return ""
	} else if ns == nil && n.Parent != nil && closeImpliesClose.Has(elementPair{n.DataAtom, n.Parent.DataAtom}) {
		return ""
	}
	return "</" + n.Data + ">"
}

// nextSiblingSkipWS returns the next sibling node of the supplied node,
// skipping over any text nodes that contain only whitespace.
func nextSiblingSkipWS(n *html.Node) *html.Node {
	for n.NextSibling != nil {
		n = n.NextSibling
		if n.Type != html.TextNode || strings.TrimSpace(n.Data) != "" {
			return n
		}
	}
	return nil
}

// voidElements are the elements that have no end tag (unless they have content,
// which they're not supposed to).
var voidElements = sets.New(
	atom.Area, atom.Base, atom.Br, atom.Col, atom.Embed, atom.Hr, atom.Img,
	atom.Input, atom.Keygen, atom.Link, atom.Meta, atom.Param, atom.Source,
	atom.Track, atom.Wbr,
)

// alwaysOmit are the elements that we never emit either start or end tags for.
var alwaysOmit = sets.New(
	atom.Html, atom.Head, atom.Body,
)

type elementPair struct{ e1, e2 atom.Atom }

// startImpliesClose is a set of element pairs.  If the current element matches
// the first of a pair, and the element being opened matches the second of the
// pair, the first is closed.
var startImpliesClose = sets.New(
	elementPair{atom.Caption, atom.Col},
	elementPair{atom.Caption, atom.Colgroup},
	elementPair{atom.Caption, atom.Thead},
	elementPair{atom.Caption, atom.Tbody},
	elementPair{atom.Caption, atom.Tfoot},
	elementPair{atom.Caption, atom.Tr},
	elementPair{atom.Colgroup, atom.Colgroup},
	elementPair{atom.Colgroup, atom.Thead},
	elementPair{atom.Colgroup, atom.Tbody},
	elementPair{atom.Colgroup, atom.Tfoot},
	elementPair{atom.Colgroup, atom.Tr},
	elementPair{atom.Dd, atom.Dd},
	elementPair{atom.Dd, atom.Dt},
	elementPair{atom.Dt, atom.Dd},
	elementPair{atom.Dt, atom.Dt},
	elementPair{atom.Li, atom.Li},
	elementPair{atom.Optgroup, atom.Optgroup},
	elementPair{atom.Option, atom.Optgroup},
	elementPair{atom.Option, atom.Option},
	elementPair{atom.P, atom.Address},
	elementPair{atom.P, atom.Article},
	elementPair{atom.P, atom.Aside},
	elementPair{atom.P, atom.Blockquote},
	elementPair{atom.P, atom.Details},
	elementPair{atom.P, atom.Div},
	elementPair{atom.P, atom.Dd}, // non-spec
	elementPair{atom.P, atom.Dl},
	elementPair{atom.P, atom.Dt}, // non-spec
	elementPair{atom.P, atom.Fieldset},
	elementPair{atom.P, atom.Figcaption},
	elementPair{atom.P, atom.Figure},
	elementPair{atom.P, atom.Footer},
	elementPair{atom.P, atom.Form},
	elementPair{atom.P, atom.H1},
	elementPair{atom.P, atom.H2},
	elementPair{atom.P, atom.H3},
	elementPair{atom.P, atom.H4},
	elementPair{atom.P, atom.H5},
	elementPair{atom.P, atom.H6},
	elementPair{atom.P, atom.Header},
	elementPair{atom.P, atom.Hgroup},
	elementPair{atom.P, atom.Hr},
	elementPair{atom.P, atom.Li}, // non-spec
	elementPair{atom.P, atom.Main},
	elementPair{atom.P, atom.Menu},
	elementPair{atom.P, atom.Nav},
	elementPair{atom.P, atom.Ol},
	elementPair{atom.P, atom.P},
	elementPair{atom.P, atom.Pre},
	elementPair{atom.P, atom.Search},
	elementPair{atom.P, atom.Section},
	elementPair{atom.P, atom.Table},
	elementPair{atom.P, atom.Ul},
	elementPair{atom.Rp, atom.Rp},
	elementPair{atom.Rp, atom.Rt},
	elementPair{atom.Rt, atom.Rp},
	elementPair{atom.Rt, atom.Rt},
	elementPair{atom.Tbody, atom.Tfoot},
	elementPair{atom.Td, atom.Tbody},
	elementPair{atom.Td, atom.Td},
	elementPair{atom.Td, atom.Tfoot},
	elementPair{atom.Td, atom.Th},
	elementPair{atom.Td, atom.Tr},
	elementPair{atom.Th, atom.Tbody},
	elementPair{atom.Th, atom.Td},
	elementPair{atom.Th, atom.Tfoot},
	elementPair{atom.Th, atom.Th},
	elementPair{atom.Th, atom.Tr},
	elementPair{atom.Thead, atom.Tbody},
	elementPair{atom.Thead, atom.Tfoot},
	elementPair{atom.Tr, atom.Tbody},
	elementPair{atom.Tr, atom.Tfoot},
	elementPair{atom.Tr, atom.Tr},
)

// closeImpliesClose is a set of element pairs.  If the current element matches
// the first of a pair, and the element being closed matches the second of the
// pair, the first is closed.
var closeImpliesClose = sets.New(
	elementPair{atom.Caption, atom.Table},
	elementPair{atom.Colgroup, atom.Table},
	elementPair{atom.Dd, atom.Dl},
	elementPair{atom.Dt, atom.Dl},
	elementPair{atom.Li, atom.Menu},
	elementPair{atom.Li, atom.Ol},
	elementPair{atom.Li, atom.Ul},
	elementPair{atom.Optgroup, atom.Select},
	elementPair{atom.Option, atom.Datalist},
	elementPair{atom.Option, atom.Optgroup},
	elementPair{atom.Option, atom.Select},
	elementPair{atom.P, atom.Address},
	elementPair{atom.P, atom.Article},
	elementPair{atom.P, atom.Aside},
	elementPair{atom.P, atom.Blockquote},
	elementPair{atom.P, atom.Body},
	elementPair{atom.P, atom.Canvas},
	elementPair{atom.P, atom.Caption},
	elementPair{atom.P, atom.Dd},
	elementPair{atom.P, atom.Details},
	elementPair{atom.P, atom.Dialog},
	elementPair{atom.P, atom.Div},
	elementPair{atom.P, atom.Dl},
	elementPair{atom.P, atom.Dt},
	elementPair{atom.P, atom.Fieldset},
	elementPair{atom.P, atom.Figcaption},
	elementPair{atom.P, atom.Figure},
	elementPair{atom.P, atom.Footer},
	elementPair{atom.P, atom.Form},
	elementPair{atom.P, atom.Header},
	elementPair{atom.P, atom.Hgroup},
	elementPair{atom.P, atom.Html},
	elementPair{atom.P, atom.Li},
	elementPair{atom.P, atom.Main},
	elementPair{atom.P, atom.Nav},
	elementPair{atom.P, atom.Object},
	elementPair{atom.P, atom.Ol},
	elementPair{atom.P, atom.Search},
	elementPair{atom.P, atom.Section},
	elementPair{atom.P, atom.Slot},
	elementPair{atom.P, atom.Td},
	elementPair{atom.P, atom.Template},
	elementPair{atom.P, atom.Th},
	elementPair{atom.P, atom.Tr},
	elementPair{atom.P, atom.Ul},
	elementPair{atom.Rp, atom.Ruby},
	elementPair{atom.Rt, atom.Ruby},
	elementPair{atom.Tbody, atom.Table},
	elementPair{atom.Td, atom.Table},
	elementPair{atom.Td, atom.Tbody},
	elementPair{atom.Td, atom.Tfoot},
	elementPair{atom.Td, atom.Thead},
	elementPair{atom.Td, atom.Tr},
	elementPair{atom.Tfoot, atom.Table},
	elementPair{atom.Th, atom.Table},
	elementPair{atom.Th, atom.Tbody},
	elementPair{atom.Th, atom.Tfoot},
	elementPair{atom.Th, atom.Thead},
	elementPair{atom.Th, atom.Tr},
	elementPair{atom.Thead, atom.Table},
	elementPair{atom.Tr, atom.Table},
	elementPair{atom.Tr, atom.Tbody},
	elementPair{atom.Tr, atom.Tfoot},
	elementPair{atom.Tr, atom.Thead},
)
