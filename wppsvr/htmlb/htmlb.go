package htmlb

import (
	"fmt"
	"html"
	"io"
	"strings"
	"sync"
)

var (
	bsHTML   = []byte(`<!DOCTYPE html><html`)
	bsLT     = []byte{'<'}
	bsLTSL   = []byte{'<', '/'}
	bsGT     = []byte{'>'}
	bsAPOS   = []byte{'\''}
	bsQUOT   = []byte{'"'}
	bsSP     = []byte{' '}
	bsEQ     = []byte{'='}
	bsEQQ    = []byte(`="`)
	bsSPSLGT = []byte(` />`)
	bsCLSEQ  = []byte(`class=`)
	bsCLSEQQ = []byte(`class="`)
	escQUOT  = []byte("&#34;")
	escAPOS  = []byte("&#39;")
	escAMP   = []byte("&amp;")
	escLT    = []byte("&lt;")
	escGT    = []byte("&gt;")
)

var siblingsNoClose = map[string]map[string]bool{
	"li":       {"li": true},
	"dt":       {"dd": true, "dt": true},
	"optgroup": {"optgroup": true},
	"option":   {"optgroup": true, "option": true},
	"p": {
		"address":    true,
		"article":    true,
		"aside":      true,
		"blockquote": true,
		"details":    true,
		"div":        true,
		"dl":         true,
		"fieldset":   true,
		"figcaption": true,
		"figure":     true,
		"footer":     true,
		"form":       true,
		"h1":         true,
		"h2":         true,
		"h3":         true,
		"h4":         true,
		"h5":         true,
		"h6":         true,
		"header":     true,
		"hgroup":     true,
		"hr":         true,
		"main":       true,
		"menu":       true,
		"nav":        true,
		"ol":         true,
		"p":          true,
		"pre":        true,
		"section":    true,
		"table":      true,
		"ul":         true,
	},
	"rp":    {"rp": true, "rt": true},
	"rt":    {"rp": true, "rt": true},
	"tbody": {"tbody": true, "tfoot": true},
	"td":    {"td": true, "th": true},
	"th":    {"td": true, "th": true},
	"thead": {"tbody": true, "tfoot": true},
	"tr":    {"tr": true},
}

// Element represents an element in an HTML document being built.
type Element struct {
	// p points to this element's parent Element.
	p *Element
	// c points to the child Element of this Element that is currently being
	// written, if any.
	c *Element
	// w is the io.Writer to which we're writing the generated HTML.
	w io.Writer
	// err is the accumulated error from writing to w.
	err error
	// tag is the element tag name.
	tag string
	// open is a flag indicating that the starting tag is still open.
	open bool
	// closing is a flag indicating that the element is in the process of
	// being closed.
	closing bool
	// classes is a list of values for the class attribute.
	classes []string
}

// elementPool is a pool of unused Element objects, for efficiency.
var elementPool = &sync.Pool{
	New: func() any { return new(Element) },
}

// HTML starts a new HTML document, emitting "<!DOCTYPE html><html>".  It
// returns a pointer to the <html> element.  The HTML document must be finished
// by calling the Close method on the element returned by HTML.
func HTML(w io.Writer) (e *Element) {
	e = elementPool.Get().(*Element)
	e.w = w
	e.tag = "html"
	e.open = true
	e.writeBytes(bsHTML)
	return e
}

// Element creates a new HTML element as a child of the receiver Element.  Any
// preceding descendant elements of the receiver are closed. before creating the
// new element.  Element returns the Element representing the new child element.
//
// The first argument must be a string of the form
//
//	<tagname attribute-assignments>content
//
// The tagname is required, and must be in lower case.  If any content for the
// element is provided, the content must be preceded by '>'.  Everything else is
// optional.
//
// Attribute assignments have the form `attrname`, `attrname=value`,
// `attrname='value'`, or `attrname="value"`.  Multiple attribute assignments
// are separated by space characters.  See the description of Attr for how the
// "class" attribute is handled specially.
//
// Attribute values and content can make use of printf-style variable
// interpolation.  Values for the interpolated variables are taken from the
// subsequent arguments to Element.  Interpolated values are HTML-encoded.
// Attribute values and content that do not contain printf-style format strings
// are emitted verbatim (not HTML-encoded).
//
// Any arguments left over, after consuming the first argument and its
// interpolated variables, are passed to the Attr method (which see).
func (e *Element) Element(args ...any) *Element {
	var tag, attrspec string

	first := args[0].(string)
	if first != "" && first[0] == '<' {
		first = first[1:]
	}
	if idx := strings.IndexAny(first, " >"); idx >= 0 {
		tag, attrspec = first[:idx], first[idx:]
	} else {
		tag, args = first, args[1:]
	}
	e.prepForContent(false, tag)
	e.c = elementPool.Get().(*Element)
	e.c.p, e.c.w, e.c.err = e, e.w, e.err
	e = e.c
	e.writeBytes(bsLT)
	e.writeString(tag)
	e.tag, e.open = tag, true
	if attrspec != "" {
		args[0] = attrspec
		args, attrspec = e.attrspec(args, false)
	}
	if attrspec != "" {
		args = e.content(attrspec, args)
	}
	if len(args) != 0 {
		e.Attr(args...)
	}
	return e
}

// Attr adds attributes to the receiver Element.  It panics if any text content
// or child elements have been emitted for the receiver Element.  Attr returns
// its receiver for convenient chaining.
//
// The arguments must follow this sequence:
//  1. An optional Boolean argument, controlling whether the attribute
//     assignments in the next argument should be emitted.
//  2. A string argument, containing a set of attribute assignments.
//  3. Zero or more variable arguments to be interpolated into the preceding
//     attribute assignment string; there must be exactly as many of these as
//     called for by printf-style format codes in the attribute assignment
//     string.
//
// This entire sequence can repeat any number of times.
//
// Attribute assignments have the form `attrname`, `attrname=value`,
// `attrname='value'`, or `attrname="value"`.  Multiple attribute assignments
// are separated by space characters.
//
// Attribute values can make use of printf-style variable interpolation.  Values
// for the interpolated variables are taken from the subsequent arguments to
// Element.  Interpolated values are HTML-encoded.  Attribute values that do not
// contain printf-style format codes are emitted verbatim (not HTML-encoded).
//
// If the Boolean argument is present in the argument sequence, and is false,
// all attribute assignments in the subsequent string argument are ignored,
// along with the corresponding interpolated variable arguments.  If the Boolean
// argument is absent or is true, the attribute assignments are emitted.
//
// The "class" attribute is handled specially, in that it can be set multiple
// times, in calls to Element or Attr or both.  All values provided for it are
// concatenated, separated by spaces, into a single attribute setting in the
// output HTML document.  As such, the "class" attribute setting is always
// emitted last, regardless of where it may appear in the argument list.
func (e *Element) Attr(args ...any) *Element {
	var remainder string

	for len(args) != 0 {
		switch arg := args[0].(type) {
		case bool:
			args, remainder = e.attrspec(args[1:], !arg)
		default:
			args, remainder = e.attrspec(args, false)
		}
		if remainder != "" {
			panic(">content is not allowed when calling Attr")
		}
	}
	return e
}

// attrspec handles a single attribute specification string and its subsequent
// variable arguments.  It returns the remaining arguments.  If the attribute
// specification string ended with ">content", it returns the content.
func (e *Element) attrspec(args []any, skip bool) ([]any, string) {
	var name, value, attrspec string

	attrspec, args = args[0].(string), args[1:]
	for name, value, attrspec = nextassn(attrspec); name != ""; name, value, attrspec = nextassn(attrspec) {
		quote, encode := false, false

		if value == "%s" {
			if s, ok := args[0].(string); ok {
				value, args, encode = s, args[1:], true
			}
		} else if count := countParams(value); count != 0 {
			value, args, encode = fmt.Sprintf(value, args[:count]...), args[count:], true
		}
		if skip {
			continue
		}
		if name == "class" {
			if encode {
				value = html.EscapeString(value)
			}
			if value != "" {
				e.classes = append(e.classes, value)
			}
			continue
		}
		e.writeBytes(bsSP)
		e.writeString(name)
		if value == "" {
			continue
		}
		if quote = needsQuoting(value); quote {
			e.writeBytes(bsEQQ)
		} else {
			e.writeBytes(bsEQ)
		}
		if encode {
			e.writeStringEnc(value)
		} else {
			e.writeString(value)
		}
		if quote {
			e.writeBytes(bsQUOT)
		}
	}
	return args, attrspec
}
func nextassn(attrspec string) (name, value, _ string) {
	var idx int
	var quote byte

	for attrspec != "" && attrspec[0] == ' ' {
		attrspec = attrspec[1:]
	}
	if attrspec == "" {
		return "", "", ""
	}
	if attrspec[0] == '>' {
		return "", "", attrspec[1:]
	}
	for idx < len(attrspec) && attrspec[idx] != ' ' && attrspec[idx] != '=' && attrspec[idx] != '>' {
		idx++
	}
	name, attrspec = attrspec[:idx], attrspec[idx:]
	if attrspec == "" || attrspec[0] != '=' {
		return name, "", attrspec
	}
	attrspec = attrspec[1:]
	if attrspec == "" || (attrspec[0] != '\'' && attrspec[0] != '"') {
		if idx = strings.IndexAny(attrspec, " >"); idx >= 0 {
			return name, attrspec[:idx], attrspec[idx:]
		}
		return name, attrspec, ""
	}
	quote, attrspec = attrspec[0], attrspec[1:]
	if idx = strings.IndexByte(attrspec, quote); idx >= 0 {
		return name, attrspec[:idx], attrspec[idx+1:]
	}
	panic("unterminated quote in attribute specification")
}
func needsQuoting(s string) bool {
	return strings.IndexAny(s, " \t\f\r\n=`") >= 0
}

// content emits the content found at the end of the initial argument of
// Element.  It returns the remaining arguments after any variable
// interpolations are consumed.
func (e *Element) content(content string, args []any) []any {
	if count := countParams(content); count != 0 {
		e.Text(fmt.Sprintf(content, args[:count]...))
		return args[count:]
	}
	e.Raw(content)
	return args
}

// Text adds text to the receiver Element.  The text will be HTML-encoded
// before being emitted.  Text returns its receiver for convenient chaining.
func (e *Element) Text(text string) *Element {
	e.prepForContent(false, "")
	e.writeStringEnc(text)
	return e
}

// Textf adds formatted text to the receiver Element.  The text will be HTML
// encoded before being emitted.  Textf returns its receiver for convenient
// chaining.
func (e *Element) Textf(f string, v ...any) *Element {
	return e.T(fmt.Sprintf(f, v...))
}

// Raw adds raw HTML markup to the content of the receiver Element.  The markup
// will be emitted verbatim (no HTML encoding).  The markup should not close the
// receiver Element, and must close any other elements that it opens.  Raw
// returns its receiver for convenient chaining.
func (e *Element) Raw(html string) *Element {
	e.prepForContent(false, "")
	e.writeString(html)
	return e
}

// Parent returns the receiver Element's parent, or nil if the receiver Element
// is the initial <html> element.
func (e *Element) Parent() *Element { return e.p }

// Close closes the receiver Element and returns its parent.  If the receiver is
// the initial <html> element, Close returns nil.  Close also returns any
// accumulated error from writing to the output writer.
func (e *Element) Close() (p *Element, err error) { return e.close("") }
func (e *Element) close(nextsibling string) (p *Element, err error) {
	e.closing = true
	emitClosingTag := true

	switch e.tag {
	case "area", "base", "body", "br", "caption", "col", "colgroup", "embed", "head", "hr", "html", "img", "input", "link", "meta", "source", "track", "wbr":
		emitClosingTag = false
	case "svg", "g", "defs", "symbol", "use", "switch", "desc" /* "title", */, "metadata", "path", "rect", "circle", "ellipse", "line", "polyline", "polygon", "text", "tspan", "textPath", "image", "foreignObject", "marker", "view":
		emitClosingTag = !e.prepForContent(true, "")
	case "dd", "li", "optgroup", "option", "rp", "rt", "tbody", "td", "tfoot", "th", "tr":
		emitClosingTag = e.p == nil || !e.p.closing
	}
	if emitClosingTag {
		if _, ok := siblingsNoClose[e.tag]; ok {
			emitClosingTag = !siblingsNoClose[e.tag][nextsibling]
		}
	}
	e.prepForContent(false, "")
	if emitClosingTag {
		e.writeBytes(bsLTSL)
		e.writeString(e.tag)
		e.writeBytes(bsGT)
	}
	p, err = e.p, e.err
	if p != nil {
		p.err = err
		p.c = nil
	}
	*e = Element{}
	elementPool.Put(e)
	return p, err
}

// E is an alias for Element, for brevity in calling code.
func (e *Element) E(args ...any) *Element { return e.Element(args...) }

// A is an alias for Attr, for brevity in calling code.
func (e *Element) A(attrs ...any) *Element { return e.Attr(attrs...) }

// T is an alias for Text, for brevity in calling code.
func (e *Element) T(text string) *Element { return e.Text(text) }

// TF is an alias for Textf, for brevity in calling code.
func (e *Element) TF(f string, v ...any) *Element { return e.Textf(f, v...) }

// R is an alias for Raw, for brevity in calling code.
func (e *Element) R(html string) *Element { return e.Raw(html) }

// P is an alias for Parent, for brevity in calling code.
func (e *Element) P() *Element { return e.p }

// prepForContent prepares for emitting content for the receiver Element.  Its
// start tag is closed, if that hasn't already happened, and any open descendant
// elements are closed.
func (e *Element) prepForContent(selfclose bool, nextsibling string) bool {
	if e.open {
		if len(e.classes) != 0 {
			e.writeBytes(bsSP)
			e.writeBytes(bsCLSEQQ)
			for i, c := range e.classes {
				if i != 0 {
					e.writeBytes(bsSP)
				}
				e.writeString(c)
			}
			e.writeBytes(bsQUOT)
		}
		e.open = false
		if selfclose {
			e.writeBytes(bsSPSLGT)
			return true
		}
		e.writeBytes(bsGT)
	}
	if e.c != nil {
		e.c.close(nextsibling)
	}
	return false
}

// countParams returns the number of parameters used in a fmt.Sprintf format
// string.
func countParams(s string) (count int) {
	var incode, innum bool
	var index, c, num int
	for _, r := range s {
		switch {
		case !incode && r != '%':
			// nothing
		case !incode:
			incode, index = true, 0
		case index == 1 && r == '%': // escaped % sign
			incode = false
		case r == '*':
			c++
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
			incode, c = false, c+1
		case r == '[':
			innum, num = true, 0
		case innum && r >= '0' && r <= '9':
			num = num*10 + int(r-'0')
		case innum && r == ']':
			innum, c = false, num-1
		}
		index++
		if c > count {
			count = c
		}
	}
	return count
}
