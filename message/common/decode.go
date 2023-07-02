package common

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Decode pulls the standard form header and footer fields from the tags map.
func (s *StdFields) Decode(tags map[string]string) {
	s.OriginMsgID = tags["MsgNo"]
	s.DestinationMsgID = tags["DestMsgNo"]
	s.MessageDate = tags["1a."]
	s.MessageTime = tags["1b."]
	s.Handling = tags["5."]
	s.ToICSPosition = tags["7a."]
	s.ToLocation = tags["7b."]
	s.ToName = tags["7c."]
	s.ToContact = tags["7d."]
	s.FromICSPosition = tags["8a."]
	s.FromLocation = tags["8b."]
	s.FromName = tags["8c."]
	s.FromContact = tags["8d."]
	s.OpRelayRcvd = tags["OpRelayRcvd"]
	s.OpRelaySent = tags["OpRelaySent"]
	s.OpName = tags["OpName"]
	s.OpCall = tags["OpCall"]
	s.OpDate = tags["OpDate"]
	s.OpTime = tags["OpTime"]
}

var (
	decodeSeverityMap = map[string]string{"E": "EMERGENCY", "U": "URGENT", "O": "OTHER"}
	decodeHandlingMap = map[string]string{"I": "IMMEDIATE", "P": "PRIORITY", "R": "ROUTINE"}
)

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
	if s, ok := decodeSeverityMap[severity]; ok {
		severity = s
	}
	if h, ok := decodeHandlingMap[handling]; ok {
		handling = h
	}
	return
}

var (
	headerRE    = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n#V: (\d+(?:\.\d+)*)-(\d+(?:\.\d+)*)\n`)
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
	} else if idx := strings.Index(body, "\n!SCCOPIFO!\n"); idx >= 0 {
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
