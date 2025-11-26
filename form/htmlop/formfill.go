package htmlop

import (
	"net/url"
	"slices"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// FillForm scans the supplied HTML document looking for form controls.  It
// adds or removes the "checked" attribute on <input type="checkbox"> and
// <input type="radio"> elements and the "selected" attribute on <option>
// elements, and sets the value of the "value" attribute on all other form
// controls, to match the supplied form values.
func FillForm(doc *html.Node, values url.Values) {
	var name string

	if doc.Type != html.ElementNode {
		goto CHILDREN
	}
	if name = getAttr(doc, "name"); name == "" {
		goto CHILDREN
	}
	switch doc.DataAtom {
	case atom.Input:
		switch getAttr(doc, "type") {
		case "hidden":
			// nothing
		case "checkbox", "radio":
			var sbchecked bool
			if hasAttr(doc, "value") {
				sbchecked = slices.Contains(values[name], getAttr(doc, "value"))
			} else {
				_, sbchecked = values[name]
			}
			if hasAttr(doc, "checked") {
				if !sbchecked {
					removeAttr(doc, "checked")
				}
			} else if sbchecked {
				setAttr(doc, "checked", "")
			}
		default:
			if val := values.Get(name); val != "" {
				setAttr(doc, "value", val)
			} else {
				removeAttr(doc, "value")
			}
		}
	case atom.Option:
		if values.Get(name) == getAttr(doc, "value") {
			if !hasAttr(doc, "selected") {
				doc.Attr = append(doc.Attr, html.Attribute{Key: "selected"})
			}
		} else {
			if hasAttr(doc, "selected") {
				removeAttr(doc, "selected")
			}
		}
	case atom.Textarea:
		if val := values.Get(name); val != "" {
			doc.FirstChild = &html.Node{Type: html.TextNode, Data: val}
			doc.LastChild = doc.FirstChild
		} else {
			doc.FirstChild, doc.LastChild = nil, nil
		}
	}
CHILDREN:
	for child := range doc.ChildNodes() {
		FillForm(child, values)
	}
}

func hasAttr(n *html.Node, key string) bool {
	return slices.ContainsFunc(n.Attr, func(a html.Attribute) bool { return strings.EqualFold(a.Key, key) })
}
func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if strings.EqualFold(attr.Key, key) {
			return attr.Val
		}
	}
	return ""
}
func setAttr(n *html.Node, key, val string) {
	for i, attr := range n.Attr {
		if strings.EqualFold(attr.Key, key) {
			n.Attr[i].Val = val
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{Key: key, Val: val})
}
func removeAttr(n *html.Node, key string) {
	n.Attr = slices.DeleteFunc(n.Attr, func(a html.Attribute) bool { return strings.EqualFold(a.Key, key) })
}
