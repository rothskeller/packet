package pktmsg

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// now returns the current time.  It can be overridden by tests.
var now = time.Now

// ErrNoPlainTextBody is raised when a message being parsed has no plain text
// body content.
var ErrNoPlainTextBody = errors.New("message has no plain text body content")

// ParseMessage parses the supplied string as a stored message (which could be
// either incoming or outgoing).
func ParseMessage(raw string) (_ *Message, err error) {
	var (
		m    Message
		mm   *mail.Message
		body []byte
	)
	// Parse the message headers.  If we can't parse them, we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(raw)); err != nil {
		return nil, err
	}
	// Stored messages should not have a Content-Type or Content-Transfer-Encoding.
	if mm.Header["Content-Type"] != nil || mm.Header["Content-Transfer-Encoding"] != nil {
		return nil, errors.New("stored messages must be plain text only")
	}
	// Extract the plain text portion of the body and create the message.
	body, _ = io.ReadAll(mm.Body)
	if err = m.parseHeadersStored(mm.Header); err != nil {
		return nil, err
	}
	m.Body = string(body)
	if err = m.parseOutpost(); err != nil {
		return nil, err
	}
	// Parse the subject line.
	m.OriginMsgID, m.Severity, m.Handling, m.FormTag, m.Subject = splitSubject(m.SubjectHeader)
	// Parse the form if any.
	m.parsePIFOForm()
	return &m, nil
}

// ReceiveMessage parses the supplied string as a message that was just
// retrieved from the specified JNOS BBS.  If it is a bulletin, area should be
// set to the bulletin area from which it was retrieved; otherwise, area should
// be empty.
func ReceiveMessage(raw, bbs, area string) (_ *Message, err error) {
	var (
		m        Message
		envelope string
		mm       *mail.Message
		body     []byte
		err2     error
	)
	m.RxBBS, m.RxArea = bbs, area
	m.RxDate.SetTime(now())
	// If there is an envelope From line, remove it from the raw message.
	if strings.HasPrefix(raw, "From ") {
		if idx := strings.IndexByte(raw, '\n'); idx > 0 {
			envelope, raw = raw[:idx], raw[idx+1:]
		}
	}
	// Parse the message headers.  If we can't parse them, we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(raw)); err != nil {
		return nil, err
	}
	// Extract the plain text portion of the body.
	body, _ = io.ReadAll(mm.Body)
	body, m.NotPlainText, err = extractPlainText(textproto.MIMEHeader(mm.Header), body)
	if err == nil && m.NotPlainText && len(body) == 0 {
		err = ErrNoPlainTextBody
	}
	// Create the message.
	if err2 = m.parseHeadersReceived(mm.Header, envelope); err == nil {
		err = err2
	}
	m.Body = string(body)
	if err2 = m.parseOutpost(); err == nil {
		err = err2
	}
	// Parse the subject line.
	m.OriginMsgID, m.Severity, m.Handling, m.FormTag, m.Subject = splitSubject(m.SubjectHeader)
	// Parse the form if any.
	m.parsePIFOForm()
	return &m, err
}

func splitSubject(s string) (oid, sev, han, ftag, subj string) {
	var (
		codes string
		found bool
		parts []string
	)
	if codes, subj, found = strings.Cut(s, " "); found {
		subj = " " + subj
	}
	parts = strings.SplitN(codes, "_", 4)
	switch len(parts) {
	case 0, 1, 2:
		return "", "", "", "", s
	case 3:
		subj = parts[2] + subj
	case 4:
		ftag = parts[2]
		subj = parts[3] + subj
	}
	oid = parts[0]
	if idx := strings.IndexByte(parts[1], '/'); idx >= 0 {
		sev = decodeSeverity(parts[1][:idx])
		han = decodeHandling(parts[1][idx+1:])
	} else {
		han = decodeHandling(parts[1])
	}
	return
}

func decodeHandling(s string) string {
	switch s {
	case "I":
		return "IMMEDIATE"
	case "P":
		return "PRIORITY"
	case "R":
		return "ROUTINE"
	}
	return s
}

func decodeSeverity(s string) string {
	switch s {
	case "E":
		return "EMERGENCY"
	case "U":
		return "URGENT"
	case "O":
		return "OTHER"
	}
	return s
}

// extractPlainText extracts the plain text portion of a message from its body.
// It returns a nil body if there is none.  The returned boolean indicates
// whether the entire body was plain text.  This is a recursive function, to
// handled nested multipart bodies.
func extractPlainText(header textproto.MIMEHeader, body []byte) (nbody []byte, notplain bool, err error) {
	var (
		mediatype string
		params    map[string]string
	)
	// Decode any content transfer encoding.  If we come across an encoding
	// we can't handle, or we have an error decoding, return an empty body
	// with a notplain indicator.
	switch strings.ToLower(header.Get("Content-Transfer-Encoding")) {
	case "", "7bit", "8bit", "binary":
		break // no decoding needed
	case "quoted-printable":
		if body, err = io.ReadAll(quotedprintable.NewReader(bytes.NewReader(body))); err != nil {
			return nil, true, err
		}
		notplain = true
	case "base64":
		if body, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(body))); err != nil {
			return nil, true, err
		}
		notplain = true
	default:
		return nil, true, nil
	}
	// Decode the content type.
	if ct := header.Get("Content-Type"); ct != "" {
		if mediatype, params, err = mime.ParseMediaType(header.Get("Content-Type")); err != nil {
			return nil, false, err // can't decode Content-Type
		}
	} else {
		mediatype, params = "text/plain", map[string]string{}
	}
	// If the content type is multipart, look for the last plain text part
	// in it.  This is a recursive call.
	if strings.HasPrefix(mediatype, "multipart/") {
		var (
			mr       *multipart.Reader
			part     *multipart.Part
			partbody []byte
			found    []byte
		)
		mr = multipart.NewReader(bytes.NewReader(body), params["boundary"])
		for {
			part, err = mr.NextRawPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, false, err // Can't decode multipart body
			}
			partbody, _ = io.ReadAll(part)
			plain, _, err := extractPlainText(part.Header, partbody)
			if err != nil {
				return nil, false, err
			}
			if plain != nil {
				found = plain
			}
		}
		return found, true, nil
	}
	// If the content type is anything other than text/plain, we're out of
	// luck.
	if mediatype != "text/plain" {
		return nil, true, nil
	}
	// In theory we also ought to check the charset, but we'll elide that
	// until experience proves a need.
	return body, notplain, nil
}

// receivedRE is the regular expression for the "Received: " line that this
// package generates when saving a received message.
var receivedRE = regexp.MustCompile(`^FROM (\S+)\.ampr\.org BY pktmsg.local(?: FOR (\S+))?; (\w\w\w, \d\d \w\w\w \d\d\d\d \d\d:\d\d:\d\d [-+]\d\d\d\d)$`)

// parseHeadersStored parses the headers of a stored message.
func (m *Message) parseHeadersStored(h mail.Header) error {
	m.parseHeadersCommon(h)
	if recv := h.Get("Received"); recv != "" {
		if match := receivedRE.FindStringSubmatch(recv); match != nil {
			m.RxBBS = match[1]
			m.RxArea = match[2]
			m.RxDate.SetString(match[3])
		} else {
			// This shouldn't happen:  stored messages with a Received: header
			// should always have our Received: header format
			return errors.New("incorrect Received: header format for stored received message")
		}
	}
	return nil
}

// parseHeadersStored parses the headers of a stored message.
func (m *Message) parseHeadersReceived(h mail.Header, envelope string) error {
	m.parseHeadersCommon(h)
	// Parse the envelope "From " line if any.
	var hadEnvelope bool
	if envelope != "" {
		hadEnvelope = true
		envelope = envelope[5:] // skip "From "
		if idx := strings.IndexByte(envelope, ' '); idx >= 0 {
			m.ReturnAddr, envelope = envelope[:idx], envelope[idx+1:]
		} else {
			m.ReturnAddr, envelope = envelope, ""
		}
		if m.ReturnAddr == "" {
			m.Autoresponse = true
		}
	}
	if envelope != "" {
		// Looks like there's a timestamp on the envelope line.
		// RFC-4155 says it should be a ctime-style timestamp in UTC.
		// We'll try parsing it as that way, but we'll treat it as local
		// time because that's what JNOS BBSes do, and those are our
		// primary source of messages to parse.
		if t, err := time.ParseInLocation(time.ANSIC, envelope, time.Local); err == nil {
			m.BBSRxDate.SetTime(t)
		}
	}
	// Compute the return address if there wasn't one on the From line.
	if !hadEnvelope {
		var line string
		if line = h.Get("Return-Path"); line == "" {
			if line = h.Get("Reply-To"); line == "" {
				if line = h.Get("Sender"); line == "" {
					line = h.Get("From")
				}
			}
		}
		// Most of those sources can have a name comment in the address,
		// which we don't want.  Also, From can have more than one
		// address in it, and we only want the first.
		if addrs, err := mail.ParseAddressList(line); err == nil && len(addrs) > 0 {
			m.ReturnAddr = addrs[0].Address
		}
	}
	// If we didn't get a BBS Rx date from the envelope, get it from the
	// Received header.
	if m.BBSRxDate.String() == "" {
		_, date, _ := strings.Cut(h.Get("Received"), ";")
		if t, err := mail.ParseDate(strings.TrimSpace(date)); err == nil {
			m.BBSRxDate.SetTime(t)
		}
	}
	return nil
}

// parseHeadersCommon is the code that is common between parseHeadersStored and
// parseHeadersRetrieved.
func (m *Message) parseHeadersCommon(h mail.Header) error {
	var to string

	m.From.SetString(h.Get("From"))
	to = strings.Join(h["To"], ", ")
	if cc := h["Cc"]; len(cc) != 0 {
		if to != "" {
			to += ", "
		}
		to += strings.Join(cc, ", ")
	}
	if bcc := h["Bcc"]; len(bcc) != 0 {
		if to != "" {
			to += ", "
		}
		to += strings.Join(bcc, ", ")
	}
	m.To.SetString(to)
	m.SubjectHeader = h.Get("Subject")
	if t, err := mail.ParseDate(h.Get("Date")); err == nil {
		m.SentDate.SetTime(t)
	}
	return nil
}

func (m *Message) parseOutpost() error {
	var (
		found bool
		body  = m.Body
	)
	for {
		switch {
		case strings.HasPrefix(body, "\n"):
			// Remove newlines that might precede the first Outpost
			// code.  We'll keep this removal only if we actually
			// find an Outpost code.
			body = body[1:]
		case strings.HasPrefix(body, "!B64!"):
			// Message content has Base64 encoding.
			if dec, err := base64.StdEncoding.DecodeString(body[5:]); err == nil {
				body = string(dec)
				found = true
			} else {
				return err
			}
		case strings.HasPrefix(body, "!RRR!"):
			m.RequestReadReceipt, found = true, true
			body = body[5:]
		case strings.HasPrefix(body, "!RDR!"):
			m.RequestDeliveryReceipt, found = true, true
			body = body[5:]
		case strings.HasPrefix(body, "!URG!"):
			m.OutpostUrgent, found = true, true
			body = body[5:]
		default:
			if found {
				m.Body = body
			} else {
				// If we didn't find any Outpost codes, leave
				// the body unchanged so as not to remove
				// initial newlines that might be significant.
			}
			return nil
		}
	}
}

var (
	headerLine    = "!SCCoPIFO!\n"
	footerLine    = "!/ADDON!\n"
	typeLineRE    = regexp.MustCompile(`^#T: ([a-z][-a-z0-9]+\.html)\n`)
	versionLineRE = regexp.MustCompile(`^#V: (\d+(?:\.\d+)*)-(\d+(?:\.\d+)*)\n`)
	fieldLineRE   = regexp.MustCompile(`(?i)^([A-Z0-9][-A-Z0-9.]*): \[`)
)

// parsePIFOForm parses the PackItForms form encoded in the message, if any.  If
// it finds one and is able to decode it, it returns the form.  Otherwise it
// returns the input message.
func (m *Message) parsePIFOForm() {
	var (
		pifoVersion string
		formHTML    string
		formVersion string
		fields      []TaggedField
		body        = m.Body
		more        = true
	)
	if pifoVersion, formHTML, formVersion, body = m.parseFormHeader(body); pifoVersion == "" {
		return
	}
	for more {
		more = false
		if remainder, tag, value, ok := m.parseFormField(body); ok && tag != "" {
			fields = append(fields, TaggedField{tag, value})
			body, more = remainder, true
		} else if !ok {
			return
		}
	}
	if !m.parseFormFooter(body) {
		return
	}
	m.PIFOVersion, m.FormHTML, m.FormVersion, m.TaggedFields = pifoVersion, formHTML, formVersion, fields
	m.Body = ""
}

// parseFormHeader parses the three header lines of a form.  Any text can come
// before the header.  (Some BBSes add routing pseudo-headers, for example.)
func (m *Message) parseFormHeader(body string) (pifoVersion, formHTML, formVersion, remainder string) {
	idx := strings.Index(body, headerLine)
	if idx < 0 || (idx > 0 && body[idx-1] != '\n') {
		return "", "", "", body
	}
	body = body[idx+len(headerLine):]
	if match := typeLineRE.FindStringSubmatch(body); match != nil {
		formHTML = match[1]
		body = body[len(match[0]):]
	} else {
		return "", "", "", body
	}
	if match := versionLineRE.FindStringSubmatch(body); match != nil {
		pifoVersion = match[1]
		formVersion = match[2]
		body = body[len(match[0]):]
	} else {
		return "", "", "", body
	}
	return pifoVersion, formHTML, formVersion, body
}

// parseFormField parses a single field definition from the form.  It is usually
// a single line, but can be multiple lines if the value is long.
func (m *Message) parseFormField(body string) (remainder, tag, value string, ok bool) {
	if match := fieldLineRE.FindStringSubmatch(body); match != nil {
		tag = match[1]
		body = body[len(match[0]):]
	} else {
		return body, "", "", true
	}
	if value, body, ok = parseBracketedValue(body); !ok {
		return body, "", "", false // illegal bracketed value
	}
	return body, tag, value, true
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
func (m *Message) parseFormFooter(body string) bool {
	return strings.HasPrefix(body, footerLine)
}
