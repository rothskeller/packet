package htmlop

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"k8s.io/apimachinery/pkg/util/sets"
)

type expandState struct {
	definitions map[string]*html.Node
	variables   map[string]string
	expanding   *html.Node
}

// Expand implements HTML template expansion according to the following
// syntax.
//
// # Variable Expansion
//
// Text nodes and attribute values in the HTML can contain `{{variable}}`
// constructs. These are replaced with the value of the named variable.
// Variable names can be any sequence of characters not including whitespace or
// a close brace, and not ending in "++".
//
// In an element tag, if an attribute setting contains a variable reference,
// and the attribute's value is empty after variable expansion, the attribute
// setting is removed.  Otherwise, any preceding settings of the attribute in
// the same element tag are removed.  This means an element tag like
// `<div attr="foo" attr="{{bar}}">` will give `attr` the value of `bar` if it
// is set, and the static value "foo" otherwise.
//
// The `{{variable++}}` construct, used when the variable's value is an
// integer, increments the variable's value after expansion.
//
// Element tags can also contain `{{variable}}` constructs (possibly with
// increment mark). They are shorthand for `variable="{{variable}}"`.
//
// # Template Expansion
//
// HTML can make use of template expansion. Templates can be defined using
// `<def:template-name>` and then later expanded using `<template-name>`.
//
// Attributes specified in the `<template-name>` element become available as
// variables that can be interpolated in the `<def:template-name>` definition.
// For example, if the definition is expanded using
// `<template-name my-var=foo>`, then instances of `{{my-var}}` in the
// definition will be replaced with `foo` in the expansion.  Attributes
// specified without a value are set to their attribute name, so if the
// definition is expanded using `<template-name my-var>`, instances of
// `{{my-var}}` in the definition will be replaced with `my-var` in the
// expansion.
//
// Attributes specified in the `<def:template-name>` element are ignored.  By
// convention, they act as documentation of the variables that the template
// uses, with their values being short descriptions of their usage.
//
// Within a template definition, several additional constructs are available.
//   - A `<slot>` element is replaced with the contents of the template
//     reference element.
//   - The contents of an `<if:variable>` element are rendered only if the
//     named variable has a non-empty definition.
//   - The contents of an `<ifnot:variable>` element are rendered only if the
//     named variable is empty (or unset).
//   - A `<for:variable>` element assumes that the named variable is a
//     semicolon-separated list of values, and renders its contents once for
//     each such value.  The value is available within that rendering as
//     `{{.}}`.  Each value has whitespace trimmed before use.  If the variable
//     contains nothing except whitespace, or is not set, the contents of the
//     `<for:variable>` are not rendered.
//
// Expand expands any templates used in the provided document, and removes the
// template definitions from it.  The provided variables, if any, are available
// for expansion in the document.
func Expand(doc *html.Node, variables map[string]string) {
	var state expandState

	state.definitions = make(map[string]*html.Node)
	if variables == nil {
		state.variables = make(map[string]string)
	} else {
		state.variables = variables
		state.expanding = new(html.Node)
	}
	state.expand(doc)
}

func (state *expandState) expand(node *html.Node) (next *html.Node) {
	// Is this a template definition?
	if node.Type == html.ElementNode && strings.HasPrefix(node.Data, "def:") {
		// Yes.  Add it to the set of definitions and take it out of the
		// DOM tree.
		state.definitions[node.Data[4:]] = node
		next = node.NextSibling
		removeNode(node)
		return next
	}
	// Is this the expansion of a template?
	if node.Type == html.ElementNode {
		if def, ok := state.definitions[node.Data]; ok {
			return state.expandTemplate(def, node)
		}
	}
	// Are we expanding a template at the moment?
	if state.expanding != nil {
		// Yes.  Is this a text node?
		if node.Type == html.TextNode {
			// Yes.  Replace any variable references in the text.
			node.Data = state.expandVariables(node.Data)
		}
		// Is this an element node?
		if node.Type == html.ElementNode {
			// Yes.  Handle variables in attributes.
			state.expandAttributes(node, true)
		}
		// Is this a <for:XXX> node?
		if node.Type == html.ElementNode && strings.HasPrefix(node.Data, "for:") {
			return state.expandFor(node)
		}
		// Is this an <if:XXX> or <ifnot:XXX> node?
		if node.Type == html.ElementNode && (strings.HasPrefix(node.Data, "if:") || strings.HasPrefix(node.Data, "ifnot:")) {
			return state.expandIf(node)
		}
		// Is this a <slot> node?
		if node.Type == html.ElementNode && node.DataAtom == atom.Slot {
			return state.expandSlot(node)
		}
	}
	// If the node has children, expand those.
	for child := node.FirstChild; child != nil; child = state.expand(child) {
	}
	return node.NextSibling
}

func (state *expandState) expandAttributes(node *html.Node, emptyBools bool) {
	for i := range node.Attr {
		// Translate an "attribute" of {{x}} into x={{x}}.
		if match := variableRE.FindStringSubmatch(node.Attr[i].Key); match != nil && len(match[0]) == len(node.Attr[i].Key) && node.Attr[i].Val == "" {
			node.Attr[i].Val = node.Attr[i].Key
			node.Attr[i].Key = match[1]
		}
		// Replace variable references in the value.
		// handle conditionals.
		if variableRE.MatchString(node.Attr[i].Val) {
			if node.Attr[i].Val = state.expandVariables(node.Attr[i].Val); node.Attr[i].Val == "" {
				node.Attr[i].Key = "" // removes attribute later
			} else {
				for prev := 0; prev < i; prev++ {
					if node.Attr[prev].Key == node.Attr[i].Key {
						node.Attr[prev].Key = ""
					}
				}
				if emptyBools && booleanAttributes.Has(strings.ToLower(node.Attr[i].Key)) {
					node.Attr[i].Val = "" // Boolean attributes should have empty values
				}
			}
		}
	}
	// The above loop marks attributes for deletion by setting their keys to
	// an empty string.  Remove any such.
	node.Attr = slices.DeleteFunc(node.Attr, func(attr html.Attribute) bool { return attr.Key == "" })
}

func (state *expandState) expandTemplate(def, node *html.Node) (next *html.Node) {
	var (
		saveVars      = make(map[string]string)
		saveExpanding = state.expanding
	)
	state.expandAttributes(node, false)
	for _, attr := range node.Attr {
		if prevval, ok := state.variables[attr.Key]; ok {
			saveVars[attr.Key] = prevval
		}
		if attr.Val != "" {
			state.variables[attr.Key] = attr.Val
		} else {
			state.variables[attr.Key] = attr.Key
		}
	}
	state.expanding = node
	for child := range def.ChildNodes() {
		clone := cloneTree(child)
		insertBefore(node, clone)
		state.expand(clone)
	}
	for _, attr := range node.Attr {
		if prevval, ok := saveVars[attr.Key]; ok {
			state.variables[attr.Key] = prevval
		} else {
			delete(state.variables, attr.Key)
		}
	}
	state.expanding = saveExpanding
	next = node.NextSibling
	removeNode(node)
	return next
}

func (state *expandState) expandSlot(node *html.Node) (next *html.Node) {
	for child := state.expanding.FirstChild; child != nil; child = next {
		next = child.NextSibling
		insertBefore(node, child)
		state.expand(child)
	}
	next = node.NextSibling
	removeNode(node)
	return next
}

func (state *expandState) expandFor(node *html.Node) (next *html.Node) {
	var prevDot = state.variables["."]

	if list := strings.TrimSpace(state.variables[node.Data[4:]]); list != "" {
		for value := range strings.SplitSeq(list, ";") {
			state.variables["."] = strings.TrimSpace(value)
			for child := range node.ChildNodes() {
				clone := cloneTree(child)
				state.expand(clone)
				insertBefore(node, clone)
			}
		}
		if prevDot != "" {
			state.variables["."] = prevDot
		} else {
			delete(state.variables, ".")
		}
	}
	next = node.NextSibling
	removeNode(node)
	return next
}

func (state *expandState) expandIf(node *html.Node) (next *html.Node) {
	var (
		value     string
		satisfied bool
	)
	if strings.HasPrefix(node.Data, "ifnot:") {
		value = state.variables[node.Data[6:]]
	} else {
		value = state.variables[node.Data[3:]]
	}
	if strings.HasPrefix(node.Data, "ifnot:") {
		satisfied = value == ""
	} else {
		satisfied = value != ""
	}
	if satisfied {
		next = node.FirstChild
		for node.FirstChild != nil {
			insertBefore(node, node.FirstChild)
		}
	}
	if next == nil {
		next = node.NextSibling
	}
	removeNode(node)
	return next
}

var variableRE = regexp.MustCompile(`\{\{([^}\s]+)\}\}`)

func (state *expandState) expandVariables(s string) (result string) {
	return variableRE.ReplaceAllStringFunc(s, state.expandVariable)
}

func (state *expandState) expandVariable(s string) (val string) {
	name := s[2 : len(s)-2]
	increment := strings.HasSuffix(name, "++")
	if increment {
		name = strings.TrimSuffix(name, "++")
	}
	val = state.variables[name]
	if increment {
		if intv, err := strconv.Atoi(val); err == nil {
			state.variables[name] = strconv.Itoa(intv + 1)
		}
	}
	return val
}

// cloneTree returns a clone of the DOM tree rooted at the supplied node.
func cloneTree(node *html.Node) (clone *html.Node) {
	clone = &html.Node{
		Type:      node.Type,
		DataAtom:  node.DataAtom,
		Data:      node.Data,
		Namespace: node.Namespace,
		Attr:      slices.Clone(node.Attr),
	}
	for child := range node.ChildNodes() {
		appendNode(clone, cloneTree(child))
	}
	return clone
}

// appendNode adds the child node as the new LastChild of the parent node.
func appendNode(parent, child *html.Node) {
	removeNode(child)
	child.Parent = parent
	if parent.LastChild != nil {
		child.PrevSibling = parent.LastChild
		parent.LastChild.NextSibling = child
	} else {
		parent.FirstChild = child
	}
	parent.LastChild = child
}

// insertBefore inserts the operant node before the anchor node in the DOM tree.
func insertBefore(anchor, operant *html.Node) {
	removeNode(operant)
	operant.Parent = anchor.Parent
	operant.PrevSibling = anchor.PrevSibling
	operant.NextSibling = anchor
	if operant.PrevSibling != nil {
		operant.PrevSibling.NextSibling = operant
	}
	anchor.PrevSibling = operant
	if operant.Parent != nil && operant.Parent.FirstChild == anchor {
		operant.Parent.FirstChild = operant
	}
}

// removeNode removes a node from the containing DOM tree.
func removeNode(node *html.Node) {
	if node.Parent != nil {
		if node.Parent.FirstChild == node {
			node.Parent.FirstChild = node.NextSibling
		}
		if node.Parent.LastChild == node {
			node.Parent.LastChild = node.PrevSibling
		}
		node.Parent = nil
	}
	if node.PrevSibling != nil {
		node.PrevSibling.NextSibling = node.NextSibling
	}
	if node.NextSibling != nil {
		node.NextSibling.PrevSibling = node.PrevSibling
	}
	node.NextSibling = nil
	node.PrevSibling = nil
}

var booleanAttributes = sets.New(
	"allowfullscreen",
	"alpha",
	"async",
	"autofocus",
	"autoplay",
	"checked",
	"controls",
	"default",
	"defer",
	"disabled",
	"formnovalidate",
	"inert",
	"ismap",
	"itemscope",
	"loop",
	"multiple",
	"muted",
	"nomodule",
	"novalidate",
	"open",
	"playsinline",
	"readonly",
	"required",
	"reversed",
	"selected",
	"shadowrootclonable",
	"shadowrootcustomelementregistry",
	"shadowrootdelegatesfocus",
	"shadowrootserializable",
)
