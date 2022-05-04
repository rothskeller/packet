package pktmsg

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

const lastLineMarker = "Ω"

type parseStateFunc func(*Form, string, bool) parseStateFunc

var (
	typeLineRE      = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)$`)
	versionLineRE   = regexp.MustCompile(`^#V: (\d+(?:\.\d+)*)-(\d+(?:\.\d+)*)$`)
	fieldLineRE     = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*):(.*)$`)
	unquoteSCCoPIFO = strings.NewReplacer(`\\`, `\`, `\n`, "\n", "`]", "]")
	unquoteLoose    = strings.NewReplacer(`\\`, `\`, `\n`, "\n")
	quoteSCCoPIFO   = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "]", "`]")
	quoteLoose      = strings.NewReplacer(`\`, `\\`, "\n", `\n`)
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
	state := expectHeader
	for _, line := range strings.Split(body, "\n") {
		if state = state(f, line, strict); state == nil {
			return nil
		}
	}
	if state(f, lastLineMarker, strict) == nil {
		return nil
	}
	return f
}
func expectHeader(f *Form, line string, _ bool) parseStateFunc {
	switch line {
	case "":
		return expectHeader
	case "!SCCoPIFO!":
		return expectType
	default:
		return nil
	}
}
func expectType(f *Form, line string, _ bool) parseStateFunc {
	if match := typeLineRE.FindStringSubmatch(line); match != nil {
		f.FormType = match[1]
		return expectVersion
	}
	return nil
}
func expectVersion(f *Form, line string, _ bool) parseStateFunc {
	if match := versionLineRE.FindStringSubmatch(line); match != nil {
		f.PIFOVersion, f.FormVersion = match[1], match[2]
		return expectField
	}
	return nil
}
func expectField(f *Form, line string, strict bool) parseStateFunc {
	if line == "!/ADDON!" {
		return expectEOF
	}
	if match := fieldLineRE.FindStringSubmatch(line); match != nil {
		var tag, value = match[1], match[2]
		if dot := strings.LastIndexByte(tag, '.'); dot >= 0 && dot < len(tag)-1 {
			if strict {
				return nil // field annotation not allowed in strict mode
			}
			tag = tag[:dot+1] // strip field annotation
		}
		if f.Has(tag) {
			return nil // multiple values for field
		}
		if value, ok := unquote(value, strict); ok {
			f.Set(tag, value)
			return expectField
		}
	}
	return nil
}
func expectEOF(_ *Form, line string, _ bool) parseStateFunc {
	if line == "" || line == lastLineMarker {
		return expectEOF
	}
	return nil
}
func unquote(value string, strict bool) (_ string, ok bool) {
	// Strict SCCoPIFO has a single space after the colon, and then a value
	// in square brackets.  Inside the brackets, newlines and backslashes
	// are escaped with a backslash, and close brackets are escaped with a
	// backtick.  A backtick at the end of the value is followed by three
	// close brackets (including the one that ends the value).
	if strings.HasPrefix(value, " [") && strings.HasSuffix(value, "]") {
		value = value[2 : len(value)-1]
		if strings.HasSuffix(value, "`]]") {
			value = value[:len(value)-2]
		}
		return unquoteSCCoPIFO.Replace(value), true
	}
	if strict {
		return "", false
	}
	// If we're not in strict mode, we don't care how much whitespace
	// appears around the square brackets.
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = value[1 : len(value)-1]
		if strings.HasSuffix(value, "`]]") {
			value = value[:len(value)-2]
		}
		return unquoteSCCoPIFO.Replace(value), true
	}
	// We also allow a value that doesn't have brackets (including an empty
	// string).  If the value is not bracketed, it's allowed to have a
	// trailing comment starting with a '#'.  Remove that.
	if hash := strings.IndexByte(value, '#'); hash >= 0 {
		value = strings.TrimSpace(value[:hash])
	}
	// For unbracketed values, newlines and backslashes have to be escaped
	// with backslashes, but nothing else is escaped.
	return unquoteLoose.Replace(value), true
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
		tags    = make([]string, len(f.Fields))
		taglens = make([]int, len(f.Fields))
	)
	for i, field := range f.Fields {
		tag := field.Tag
		if annotations != nil {
			tag += annotations[tag]
		}
		tags[i] = tag
		taglens[i] = len(tag) + 1
	}
	if looseQuoting {
		alignLengths(taglens)
	}
	for i, tag := range tags {
		if f.Fields[i].Value == "" && comments != nil {
			if comment := comments[f.Fields[i].Tag]; comment != "" {
				fmt.Fprintf(sb, "%s:%*s# %s\n", tag, taglens[i]-len(tag), "", comment)
				continue
			}
		}
		fmt.Fprintf(sb, "%s:%*s%s\n", tag, taglens[i]-len(tag), "", quote(f.Fields[i].Value, looseQuoting))
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
func quote(value string, loose bool) string {
	if !loose || strings.HasPrefix(value, "[") || strings.TrimSpace(value) != value {
		value = quoteSCCoPIFO.Replace(value)
		if strings.HasSuffix(value, "`") {
			value += "]]"
		}
		return "[" + value + "]"
	}
	return quoteLoose.Replace(value)
}
