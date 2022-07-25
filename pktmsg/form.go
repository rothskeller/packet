package pktmsg

import (
	"fmt"
	"math"
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
func ParseForm(body string, strict bool) (f *Form) {
	f = new(Form)
	if f, body = parseFormHeader(body); f == nil {
		return nil
	}
	for {
		var done bool
		if f, body, done = parseFormField(f, body, strict); done {
			break
		}
	}
	if f != nil {
		f = parseFormFooter(f, body)
	}
	return f
}

// parseFormHeader parses the three header lines of a form.
func parseFormHeader(body string) (f *Form, _ string) {
	body = strings.TrimLeft(body, "\n")
	if !strings.HasPrefix(body, headerLine) {
		return nil, body
	}
	body = body[len(headerLine):]
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
func parseFormField(f *Form, body string, strict bool) (_ *Form, _ string, done bool) {
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
	if dot := strings.LastIndexByte(tag, '.'); dot >= 0 && dot < len(tag)-1 {
		if strict {
			return nil, body, true // field annotation not allowed in strict mode
		}
		tag = tag[:dot+1] // strip field annotation
	}
	if f.Has(tag) {
		return nil, body, true // multiple values for field
	}
	if strict && !strings.HasPrefix(body, " [") {
		return nil, body, true // strict mode requires bracket after exactly one space
	}
	body = strings.TrimLeft(body, " ")
	if strings.HasPrefix(body, "[") {
		value, body, ok = parseBracketedValue(body[1:], strict)
	} else {
		value, body, ok = parseUnbracketedValue(body)
	}
	if !ok {
		return nil, body, true
	}
	f.Set(tag, value)
	return f, body, false
}

// parseBracketedValue parses a field value in brackets.  Within the brackets,
// \n represents a newline, \\ represents a backslash, literal newlines are
// ignored, `] represents a close bracket, and `]]] represents a backtick at the
// end of the string.  The final close bracket must be followed by a newline.
// Spaces are allowed between the close bracket and the newline if we're not in
// strict mode.
func parseBracketedValue(body string, strict bool) (value, nbody string, ok bool) {
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
	if !strict {
		body = strings.TrimLeft(body, " ")
	}
	if body == "" || body[0] != '\n' { // no newline, or extra text after value
		return "", body, false
	}
	return value, body[1:], true
}

// parseUnbracketedValue parses an unbracketed value (which is only allowed in
// non-strict mode).  Leading and trailing blanks are ignored.  A pound sign (#)
// ends the value, and causes the rest of the line to be ignored; otherwise, the
// value ends at the end of the line.  \n represents a newline and \\ represents
// a backslash.
func parseUnbracketedValue(body string) (value, nbody string, ok bool) {
	for body != "" {
		if body[0] == '\n' { // newline ends the value
			return strings.TrimSpace(value), body[1:], true
		}
		if body[0] == '#' { // # ends the value and starts a comment
			value = strings.TrimSpace(value)
			if idx := strings.IndexByte(body, '\n'); idx >= 0 {
				return strings.TrimSpace(value), body[idx+1:], true
			}
			return "", "", false
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
		// anything else is copied literally
		r, sz := utf8.DecodeRuneInString(body)
		value += string(r)
		body = body[sz:]
	}
	return "", "", false // end of body before end of value
}

// parseFormFooter verifies that the only thing left in the body is the footer
// line, possibly followed by blank lines.
func parseFormFooter(f *Form, body string) *Form {
	if !strings.HasPrefix(body, footerLine) {
		return nil
	}
	body = body[len(footerLine):]
	if strings.Trim(body, "\n") != "" {
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
func (f *Form) EncodeToMessage(msg *Message, annotations, comments map[string]string, looseQuoting bool) {
	msg.Body = f.Encode(annotations, comments, looseQuoting)
}

// Encode returns the encoded form.  If annotations and/or comments are non-nil,
// they provide annotations to add to the field tags and comments to display in
// place of empty field values.  If looseQuoting is true, values are encoded
// with loose quoting; otherwise they are encoded with strict SCCoPIFO quoting.
// annotations and comments must be nil and looseQuoting must be false when
// encoding a message for transmission.
func (f *Form) Encode(annotations, comments map[string]string, looseQuoting bool) string {
	var sb strings.Builder

	if f.PIFOVersion == "" {
		f.PIFOVersion = "3.2"
	}
	sb.WriteString("!SCCoPIFO!\n#T: ")
	sb.WriteString(f.FormType)
	sb.WriteString("\n#V: ")
	sb.WriteString(f.PIFOVersion)
	sb.WriteByte('-')
	sb.WriteString(f.FormVersion)
	sb.WriteByte('\n')
	f.encodeFields(&sb, annotations, comments, looseQuoting)
	sb.WriteString("!/ADDON!\n")
	return sb.String()
}
func (f *Form) encodeFields(sb *strings.Builder, annotations, comments map[string]string, looseQuoting bool) {
	var (
		tags      []string
		values    []string
		fcomments []string
		taglens   []int
	)
	for _, field := range f.Fields {
		tag := field.Tag
		if annotations != nil {
			tag += annotations[tag]
		}
		if field.Value != "" || looseQuoting {
			tags = append(tags, tag)
			values = append(values, field.Value)
			taglens = append(taglens, len(tag)+1)
			if comments != nil {
				fcomments = append(fcomments, comments[field.Tag])
			}
		}
	}
	if looseQuoting {
		alignLengths(taglens)
	}
	for i, tag := range tags {
		if values[i] == "" && fcomments != nil && fcomments[i] != "" {
			fmt.Fprintf(sb, "%s:%*s# %s\n", tag, taglens[i]-len(tag), "", fcomments[i])
		} else if looseQuoting && !strings.HasPrefix(values[i], "[") && strings.TrimSpace(values[i]) == values[i] {
			fmt.Fprintf(sb, "%s:%*s%s\n", tag, taglens[i]-len(tag), "", quoteLoose.Replace(values[i]))
		} else {
			var s = fmt.Sprintf("%s: %s", tag, bracketQuote(values[i]))
			for len(s) > 128 {
				fmt.Fprintln(sb, s[:128])
				s = s[128:]
			}
			fmt.Fprintln(sb, s)
		}
	}
}

func alignLengths(lengths []int) {
	// This algorithm is stolen from go/printer.exprList().  It aligns items
	// (in this case, annotated field tags) that are of similar size, but
	// starts a new alignment block when the size changes significantly.
	// Specifically it starts a new alignment block when it comes across a
	// line whose size differs from the geometric mean of the previous line
	// sizes by greater than a threshold ratio.  A new alignment block is
	// never started when both lines are below a minimum size.
	const (
		minimum   = 12
		threshold = 2.5
	)
	var (
		start int
		max   int
		lnsum float64
	)
	for i := range lengths {
		var newblock bool
		if i > start && (lengths[i] > minimum || lengths[i-1] > minimum) {
			var mean = math.Exp(lnsum / float64(i-start))
			var ratio = float64(lengths[i]) / mean
			newblock = threshold*ratio <= 1 || threshold <= ratio
		}
		if newblock {
			for j := start; j < i; j++ {
				lengths[j] = max
			}
			start, max, lnsum = i, 0, 0
		}
		if lengths[i] > max {
			max = lengths[i]
		}
		lnsum += math.Log(float64(lengths[i]))
	}
	for i := start; i < len(lengths); i++ {
		lengths[i] = max
	}
}

func bracketQuote(value string) string {
	value = quoteSCCoPIFO.Replace(value)
	if strings.HasSuffix(value, "`") {
		value += "]]"
	}
	return "[" + value + "]"
}
