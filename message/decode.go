package message

// The code in this file handles decoding messages from their saved form into a
// message.Message.

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"slices"
)

var (
	// DecodeSeverityMap maps severity codes, as returned by DecodeSubject,
	// into Situation Severity values.
	DecodeSeverityMap = map[string]string{"E": "EMERGENCY", "U": "URGENT", "O": "OTHER"}
	// DecodeHandlingMap maps handling codes, as returned by DecodeSubject,
	// into Handling values.
	DecodeHandlingMap = map[string]string{"I": "IMMEDIATE", "P": "PRIORITY", "R": "ROUTINE"}
)

// DecodeForm decodes the message.  It returns whether the message was
// recognized and decoded.
func DecodeForm(body string, versions []*FormVersion, create func(*FormVersion) Message) (msg Message) {
	var (
		form *PIFOForm
		bm   *BaseMessage
	)
	// Decode the form and check for an HTML/version combo we recognize.
	if form = DecodePIFO(body); form == nil {
		return nil // not a form or not encoded properly
	}
	if idx := slices.IndexFunc(versions, func(v *FormVersion) bool {
		return form.HTMLIdent == v.HTML && form.FormVersion == v.Version
	}); idx < 0 {
		return nil // not an HTML/version combo we recognize
	} else {
		// Record which form version we actually saw.
		msg = create(versions[idx])
	}
	bm = msg.(interface{ Base() *BaseMessage }).Base()
	bm.PIFOVersion = form.PIFOVersion
	for _, f := range bm.Fields {
		if f.PIFOTag == "" {
			continue // field is not part of PIFO encoding
		}
		*f.Value = form.TaggedValues[f.PIFOTag]
	}
	// TODO really should make sure there aren't any fields unaccounted for.
	return msg
}

// DecodeSubject decodes an XSC-standard message subject line into its component
// parts.  If the subject line does not follow the XSC standard, the function
// returns "", "", "", "", line.
func DecodeSubject(line string) (msgid, severity, handling, formtag, subject string) {
	var (
		codes string
		found bool
		parts []string
	)
	if codes, subject, found = strings.Cut(line, " "); found {
		subject = " " + subject
	}
	parts = strings.SplitN(codes, "_", 4)
	switch len(parts) {
	case 0, 1, 2:
		return "", "", "", "", line
	case 3:
		subject = parts[2] + subject
	case 4:
		formtag = parts[2]
		subject = parts[3] + subject
	}
	msgid = parts[0]
	if idx := strings.IndexByte(parts[1], '/'); idx >= 0 {
		severity = parts[1][:idx]
		handling = parts[1][idx+1:]
	} else {
		handling = parts[1]
	}
	return
}

var (
	headerRE    = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n#V: (\d+(?:\.\d+)*[A-Za-z]?)-(\d+(?:\.\d+)*[A-Za-z]*)\n`)
	fieldLineRE = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*): \[`)
)

// PIFOForm is a decoded PackItForms form.
type PIFOForm struct {
	TextBefore   string
	HTMLIdent    string
	PIFOVersion  string
	FormVersion  string
	TaggedValues map[string]string
	TextAfter    string
}

// DecodePIFO decodes a message body and returns the decoded form contents.  If the
// body does not contain a valid encoded form, Decode returns nil.
func DecodePIFO(body string) (f *PIFOForm) {
	if strings.HasPrefix(body, "!SCCoPIFO!\n") {
		f = new(PIFOForm)
		body = body[11:]
	} else if idx := strings.Index(body, "\n!SCCoPIFO!\n"); idx >= 0 {
		f = new(PIFOForm)
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
	f.TaggedValues = make(map[string]string)
	for {
		var (
			match []string
			tag   string
			value string
			ok    bool
		)
		if strings.HasPrefix(body, "\n") {
			// Blank lines are allowed.  They can result from JNOS
			// inserting line breaks in the wrong place.
			body = body[1:]
			continue
		}
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
	if !strings.HasPrefix(body, "!/ADDON!") || (len(body) > 8 && body[8] != '\n') {
		return nil
	}
	f.TextAfter = body[8:]
	if len(f.TextAfter) != 0 && f.TextAfter[0] == '\n' {
		f.TextAfter = f.TextAfter[1:]
	}
	return f
}

// parseBracketedValue parses a field value in brackets.  Within the brackets,
// \n represents a newline, \\ represents a backslash, `] represents a close
// bracket, and `]]] represents a backtick at the end of the string.  Literal
// newlines are ignored, even in the middle of the above sequences.  The final
// close bracket must be followed by a newline.
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
		if len(body) > 2 && body[0] == '\\' && body[1] == '\n' && body[2] == '\\' { // escaped backslashes are a single backslash
			value += "\\"
			body = body[3:]
			continue
		}
		if len(body) > 1 && body[0] == '\\' && body[1] == 'n' { // escaped 'n's are newlines
			value += "\n"
			body = body[2:]
			continue
		}
		if len(body) > 2 && body[0] == '\\' && body[1] == '\n' && body[2] == 'n' { // escaped 'n's are newlines
			value += "\n"
			body = body[3:]
			continue
		}
		if len(body) > 3 && body[0] == '`' && body[1] == ']' && body[2] == ']' && body[3] == ']' {
			// backtick followed by three close brackets is a backtick that ends the string
			ok = true
			value += "`"
			body = body[4:]
			break
		}
		if len(body) > 4 &&
			((body[0] == '`' && body[1] == '\n' && body[2] == ']' && body[3] == ']' && body[4] == ']') ||
				(body[0] == '`' && body[1] == ']' && body[2] == '\n' && body[3] == ']' && body[4] == ']') ||
				(body[0] == '`' && body[1] == ']' && body[2] == ']' && body[3] == '\n' && body[4] == ']')) {
			// backtick followed by three close brackets is a backtick that ends the string
			ok = true
			value += "`"
			body = body[5:]
			break
		}
		if len(body) > 1 && body[0] == '`' && body[1] == ']' { // backtick, close bracket is a literal close bracket
			value += "]"
			body = body[2:]
			continue
		}
		if len(body) > 2 && body[0] == '`' && body[1] == '\n' && body[2] == ']' { // backtick, close bracket is a literal close bracket
			value += "]"
			body = body[3:]
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
