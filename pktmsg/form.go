package pktmsg

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// pifoVersion is the PIFO version number used in new outgoing forms.
const pifoVersion = "3.9"

var (
	headerLine    = "!SCCoPIFO!\n"
	footerLine    = "!/ADDON!\n"
	typeLineRE    = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n`)
	versionLineRE = regexp.MustCompile(`^#V: (\d+(?:\.\d+)*)-(\d+(?:\.\d+)*)\n`)
	fieldLineRE   = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*): \[`)
	quoteSCCoPIFO = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")
)

type pifoMessage struct {
	Message
	pifoVersion field
	formHTML    field
	formVersion field
	fields      []*taggedField
}

// NewForm creates a new, outgoing form with the specified form HTML, version
// number, and fields.
func NewForm(html, version string, fields []*taggedField) Message {
	var m = pifoMessage{Message: newOutpostMessage()}
	m.pifoVersion = field(pifoVersion)
	m.formHTML = field(html)
	m.formVersion = field(version)
	m.fields = fields
	return &m
}

// parsePIFOForm parses the PackItForms form encoded in the message, if any.  If
// it finds one and is able to decode it, it returns the form.  Otherwise it
// returns the input message.
func parsePIFOForm(in Message) Message {
	var (
		ok   bool
		body = in.Body().Value()
		m    = pifoMessage{Message: in}
	)
	if body, ok = m.parseFormHeader(body); !ok {
		return in
	}
	for {
		var done bool
		if body, done, ok = m.parseFormField(body); done {
			break
		}
	}
	if ok {
		ok = m.parseFormFooter(body)
	}
	if !ok {
		return in
	}
	return &m
}

// parseFormHeader parses the three header lines of a form.  Any text can come
// before the header.  (Some BBSes add routing pseudo-headers, for example.)
func (m *pifoMessage) parseFormHeader(body string) (string, bool) {
	idx := strings.Index(body, headerLine)
	if idx < 0 || (idx > 0 && body[idx-1] != '\n') {
		return body, false
	}
	body = body[idx+len(headerLine):]
	if match := typeLineRE.FindStringSubmatch(body); match != nil {
		m.formHTML = field(match[1])
		body = body[len(match[0]):]
	} else {
		return body, false
	}
	if match := versionLineRE.FindStringSubmatch(body); match != nil {
		m.pifoVersion = field(match[1])
		m.formVersion = field(match[2])
		body = body[len(match[0]):]
	} else {
		return body, false
	}
	return body, true
}

// parseFormField parses a single field definition from the form.  It is usually
// a single line, but can be multiple lines if the value is long.
func (m *pifoMessage) parseFormField(body string) (_ string, done, ok bool) {
	var (
		tf    taggedField
		value string
	)
	if match := fieldLineRE.FindStringSubmatch(body); match != nil {
		tf.tag = match[1]
		body = body[len(match[0]):]
	} else {
		return body, true, true
	}
	for _, f := range m.fields {
		if f.Tag() == tf.tag {
			return body, true, false // multiple values for field
		}
	}
	if value, body, ok = parseBracketedValue(body); !ok {
		return body, true, false // illegal bracketed value
	}
	tf.SetValue(value)
	m.fields = append(m.fields, &tf)
	return body, false, true
}

// parseBracketedValue parses a field value in brackets.  Within the brackets,
// \n represents a newline, \\ represents a backslash, literal newlines are
// ignored, `] represents a close bracket, and `]]] represents a backtick at the
// end of the string.  The final close bracket must be followed by a newline.
func parseBracketedValue(body string) (value, nbody string, ok bool) {
	for body != "" {
		if body[0] == ']' { // close bracket ends the value
			ok = true
			body = body[1:]
			break
		}
		if body[0] == '\n' { // newlines are ignored
			body = body[1:]
			continue
		}
		if len(body) > 1 && body[0] == '\\' && body[1] == '\\' { // escaped backslashes are a single backslash
			value += "\\"
			body = body[2:]
			continue
		}
		if len(body) > 1 && body[0] == '\\' && body[1] == 'n' { // escaped 'n's are newlines
			value += "\n"
			body = body[2:]
			continue
		}
		if len(body) > 3 && body[0] == '`' && body[1] == ']' && body[2] == ']' && body[3] == ']' {
			// backtick followed by three close brackets is a backtick that ends the string
			ok = true
			value += "`"
			body = body[4:]
			break
		}
		if len(body) > 1 && body[0] == '`' && body[1] == ']' { // backtick, close bracket is a literal close bracket
			value += "]"
			body = body[2:]
			continue
		}
		// anything else is copied literally
		r, sz := utf8.DecodeRuneInString(body)
		value += string(r)
		body = body[sz:]
	}
	if !ok { // end of body before end of value
		return "", body, false
	}
	if body == "" || body[0] != '\n' { // no newline, or extra text after value
		return "", body, false
	}
	return value, body[1:], true
}

// parseFormFooter verifies that the next thing in the body is the footer line.
// Anything after that is ignored.  (This is because some email clients add
// their own footers that the sender can't control.)
func (m *pifoMessage) parseFormFooter(body string) bool {
	return strings.HasPrefix(body, footerLine)
}

// Accessors for pifoMessage fields.
func (pifoMessage) Body() BodyField                  { return nil }
func (m *pifoMessage) FormHTML() FormHTMLField       { return &m.formHTML }
func (m *pifoMessage) FormVersion() FormVersionField { return &m.formVersion }
func (m *pifoMessage) PIFOVersion() PIFOVersionField { return &m.pifoVersion }

func (m *pifoMessage) TaggedField(tag string) Field {
	for _, f := range m.fields {
		if f.Tag() == tag {
			return f
		}
	}
	return nil
}

func (m *pifoMessage) TaggedFields(fn func(string, Field)) {
	for _, f := range m.fields {
		fn(f.tag, f)
	}
}

func (m *pifoMessage) Save() string {
	m.encode()
	return m.Message.Save()
}

func (m *pifoMessage) Transmit() (to []string, subject string, body string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *pifoMessage) encode() {
	var sb strings.Builder

	sb.WriteString("!SCCoPIFO!\n#T: ")
	sb.WriteString(m.formHTML.Value())
	sb.WriteString("\n#V: ")
	sb.WriteString(m.pifoVersion.Value())
	sb.WriteByte('-')
	sb.WriteString(m.formVersion.Value())
	sb.WriteByte('\n')
	for _, f := range m.fields {
		sb.WriteString(pifoEncode(f.tag, f))
	}
	sb.WriteString("!/ADDON!\n")
	m.Message.Body().SetValue(sb.String())
}

func pifoEncode(tag string, f Field) string {
	if f.Value() == "" {
		return "" // omit empty fields
	}
	var s = fmt.Sprintf("%s: %s", tag, bracketQuote(f.Value()))
	var w string
	for len(s) > 128 {
		w += s[:128] + "\n"
		s = s[128:]
	}
	return w + s + "\n"
}

func bracketQuote(value string) string {
	value = quoteSCCoPIFO.Replace(value)
	if strings.HasSuffix(value, "`") {
		value += "]]"
	}
	return "[" + value + "]"
}

// The PIFOVersionField holds the PackItForms encoding version number.  It is
// present only on form messages.
type PIFOVersionField interface{ Field }

// The FormHTMLField holds the PackItForms HTML file for the form.  It is
// present only on form messages.
type FormHTMLField interface{ Field }

// The FormVersionField holds the form version number.  It is present only on
// form messages.
type FormVersionField interface{ Field }

// taggedField is a field with a tag.
type taggedField struct {
	settableField
	tag string
}

func (f taggedField) Tag() string { return f.tag }
