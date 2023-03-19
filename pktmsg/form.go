package pktmsg

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	headerLine    = "!SCCoPIFO!\n"
	footerLine    = "!/ADDON!\n"
	typeLineRE    = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n`)
	versionLineRE = regexp.MustCompile(`^#V: (\d+(?:\.\d+)*)-(\d+(?:\.\d+)*)\n`)
	fieldLineRE   = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*):`)
	quoteSCCoPIFO = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")
	quoteLoose    = strings.NewReplacer(`\`, `\\`, "\n", `\n`)
)

// IsForm returns whether the supplied message body looks like it has an
// embedded form.
func IsForm(body string) bool {
	return strings.Contains(body, "!SCCoPIFO!") || strings.Contains(body, "!PACF!") || strings.Contains(body, "!/ADDON!")
}

// ParseForm parses the supplied message body as a form.  It returns the decoded
// form, or nil if it could not be parsed.  The strict flag indicates whether
// strict SCCoPIFO syntax must be used.  If it is false, field annotations,
// comments, and loose quoting are allowed.
func ParseForm(body string) (f *Form) {
	f = new(Form)
	if f, body = parseFormHeader(body); f == nil {
		return nil
	}
	for {
		var done bool
		if f, body, done = parseFormField(f, body); done {
			break
		}
	}
	if f != nil {
		f = parseFormFooter(f, body)
	}
	return f
}

// parseFormHeader parses the three header lines of a form.  Any text can come
// before the header.  (Some BBSes add routing pseudo-headers, for example.)
func parseFormHeader(body string) (f *Form, _ string) {
	idx := strings.Index(body, headerLine)
	if idx < 0 || (idx > 0 && body[idx-1] != '\n') {
		return nil, body
	}
	body = body[idx+len(headerLine):]
	if match := typeLineRE.FindStringSubmatch(body); match != nil {
		f = &Form{FormType: match[1]}
		body = body[len(match[0]):]
	} else {
		return nil, body
	}
	if match := versionLineRE.FindStringSubmatch(body); match != nil {
		f.PIFOVersion, f.FormVersion = match[1], match[2]
		body = body[len(match[0]):]
	} else {
		return nil, body
	}
	return f, body
}

// parseFormField parses a single field definition from the form.  It is usually
// a single line, but can be multiple lines if the value is long.
func parseFormField(f *Form, body string) (_ *Form, _ string, done bool) {
	var (
		tag   string
		value string
		ok    bool
	)
	if match := fieldLineRE.FindStringSubmatch(body); match != nil {
		tag = match[1]
		body = body[len(match[0]):]
	} else {
		return f, body, true
	}
	if f.Has(tag) {
		return nil, body, true // multiple values for field
	}
	if !strings.HasPrefix(body, " [") {
		return nil, body, true // must have bracket after exactly one space
	}
	if value, body, ok = parseBracketedValue(body[2:]); !ok {
		return nil, body, true
	}
	f.Set(tag, value)
	return f, body, false
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
func parseFormFooter(f *Form, body string) *Form {
	if !strings.HasPrefix(body, footerLine) {
		return nil
	}
	return f
}

// A Form represents a form encoded in a packet message.
type Form struct {
	// PIFOVersion identifies the version of the PackItForms encoding of the
	// form.
	PIFOVersion string
	// FormType is a PackItForms HTML file name that identifies the form
	// type.
	FormType string
	// FormVersion identifies the form version.
	FormVersion string
	// Fields is a list of form fields.
	Fields []FormField
}

// A FormField is a single field of a Form.
type FormField struct {
	// Tag is the field tag as it appears in the PackItForms encoding of the
	// form.  In most cases this is a number followed by a period.
	Tag string
	// Value is the value of the field, in string form.
	Value string
}

// Has returns whether the form has a setting (even if empty) for the field with
// the specified tag.
func (f *Form) Has(tag string) bool {
	for _, ff := range f.Fields {
		if ff.Tag == tag {
			return true
		}
	}
	return false
}

// Get retrieves the value of the field with the specified tag.  If the field
// is not in the form, it returns an empty string.
func (f *Form) Get(tag string) string {
	for _, ff := range f.Fields {
		if ff.Tag == tag {
			return ff.Value
		}
	}
	return ""
}

// Set sets a field value.
func (f *Form) Set(tag, value string) {
	for i, ff := range f.Fields {
		if ff.Tag == tag {
			f.Fields[i].Value = value
			return
		}
	}
	f.Fields = append(f.Fields, FormField{Tag: tag, Value: value})
}

// EncodeToMessage encodes the form into the supplied message.  If annotations
// and/or comments are non-nil, they provide annotations to add to the field
// tags and comments to display in place of empty field values.  If looseQuoting
// is true, values are encoded with loose quoting; otherwise they are encoded
// with strict SCCoPIFO quoting.  annotations and comments must be nil and
// looseQuoting must be false when encoding a message for transmission.
func (f *Form) EncodeToMessage(msg *Message) {
	msg.Body = f.Encode()
}

// Encode returns the encoded form.
func (f *Form) Encode() string {
	var sb strings.Builder

	if f.PIFOVersion == "" {
		f.PIFOVersion = "3.9"
	}
	sb.WriteString("!SCCoPIFO!\n#T: ")
	sb.WriteString(f.FormType)
	sb.WriteString("\n#V: ")
	sb.WriteString(f.PIFOVersion)
	sb.WriteByte('-')
	sb.WriteString(f.FormVersion)
	sb.WriteByte('\n')
	f.encodeFields(&sb)
	sb.WriteString("!/ADDON!\n")
	return sb.String()
}
func (f *Form) encodeFields(sb *strings.Builder) {
	var (
		tags   []string
		values []string
	)
	for _, field := range f.Fields {
		tag := field.Tag
		if field.Value != "" {
			tags = append(tags, tag)
			values = append(values, field.Value)
		}
	}
	for i, tag := range tags {
		var s = fmt.Sprintf("%s: %s", tag, bracketQuote(values[i]))
		for len(s) > 128 {
			fmt.Fprintln(sb, s[:128])
			s = s[128:]
		}
		fmt.Fprintln(sb, s)
	}
}

func bracketQuote(value string) string {
	value = quoteSCCoPIFO.Replace(value)
	if strings.HasSuffix(value, "`") {
		value += "]]"
	}
	return "[" + value + "]"
}
