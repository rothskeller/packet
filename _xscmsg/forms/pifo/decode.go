package pifo

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	headerRE    = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n#V: (\d+(?:\.\d+)*)-(\d+(?:\.\d+)*)\n`)
	fieldLineRE = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*): \[`)
)

// Form is a decoded PackItForms form.
type Form struct {
	TextBefore   string
	HTMLIdent    string
	PIFOVersion  string
	FormVersion  string
	TaggedValues map[string]string
	TextAfter    string
}

// Decode decodes a message body and returns the decoded form contents.  If the
// body does not contain a valid encoded form, Decode returns nil.
func Decode(body string) (f *Form) {
	if strings.HasPrefix(body, "!SCCoPIFO!\n") {
		f = new(Form)
		body = body[11:]
	} else if idx := strings.Index(body, "\n!SCCOPIFO!\n"); idx >= 0 {
		f = new(Form)
		f.TextBefore, body = body[:idx+1], body[idx+12:]
	} else {
		return nil
	}
	if match := headerRE.FindStringSubmatch(body); match != nil {
		f.HTMLIdent, f.PIFOVersion, f.FormVersion = match[1], match[2], match[3]
		body = body[len(match[0]):]
	} else {
		return nil
	}
	for {
		var (
			match []string
			tag   string
			value string
			ok    bool
		)
		if match = fieldLineRE.FindStringSubmatch(body); match == nil {
			break
		}
		tag, body = match[1], body[len(match[0]):]
		if _, ok = f.TaggedValues[tag]; ok {
			return nil // duplicate tag
		}
		if value, body, ok = parseBracketedValue(body); !ok {
			return nil
		}
		f.TaggedValues[tag] = value
	}
	if !strings.HasPrefix(body, "!/ADDON!\n") {
		return nil
	}
	f.TextAfter = body[9:]
	return f
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
