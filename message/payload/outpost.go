package payload

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"net/mail"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
)

// OutpostPayload is a message payload interchangeable with Outpost.
type OutpostPayload struct {
	cachetrack.Tracker
	body      body.Body
	urgent    bool
	requestDR bool
	requestRR bool
	allowLong bool
	bbsRoutes string
}

var _ Payload = (*OutpostPayload)(nil)

// NewOutpostPayload creates a new Payload containing the supplied body.
func NewOutpostPayload(body body.Body, allowLong bool) *OutpostPayload {
	p := &OutpostPayload{body: body, allowLong: allowLong}
	if body.Dirty() {
		p.MarkDirty("payload.OutpostPayload.Body")
	}
	body.OnDirty(p.MarkDirty)
	return p
}

// Body returns the payload body.
func (p *OutpostPayload) Body() body.Body { return p.body }

// Urgent is a marker that a message is urgent.  Outpost shows such messages in
// red text in the message listing.  This generally corresponds to messages
// with an IMMEDIATE handling order.
func (p *OutpostPayload) Urgent() bool { return p.urgent }

// SetUrgent sets the Urgent flag.
func (p *OutpostPayload) SetUrgent(urgent bool) {
	if p.urgent != urgent {
		p.urgent = urgent
		p.MarkDirty("payload.OutpostPayload.Urgent")
	}
}

// RequestDR is a marker that the receiver should send a delivery receipt for
// the message when it is delivered to the recipient's system (i.e., fetched
// from JNOS by the receiver).
func (p *OutpostPayload) RequestDR() bool { return p.requestDR }

// SetRequestDR sets the RequestDR flag.
func (p *OutpostPayload) SetRequestDR(requestDR bool) {
	if p.requestDR != requestDR {
		p.requestDR = requestDR
		p.MarkDirty("payload.OutpostPayload.RequestDR")
	}
}

// RequestRR is a marker that the receiver should sent a read receipt for the
// message when it is read by the recipient.
func (p *OutpostPayload) RequestRR() bool { return p.requestRR }

// SetRequestRR sets the RequestRR flag.
func (p *OutpostPayload) SetRequestRR(requestRR bool) {
	if p.requestRR != requestRR {
		p.requestRR = requestRR
		p.MarkDirty("payload.OutpostPayload.RequestRR")
	}
}

// BBSRoutes returns the BBS routing lines at the top of the body (if any).
func (p *OutpostPayload) BBSRoutes() string { return p.bbsRoutes }

// SetAllowLong sets the allow-long-lines flag on the payload.  This prevents
// base64 encoding just because line lengths exceed the JNOS limit.  This flag
// should be set only when the payload is resilient to having extra line breaks
// added by JNOS.
func (p *OutpostPayload) SetAllowLong() { p.allowLong = true }

// Clone returns a copy of the Payload.
func (p *OutpostPayload) Clone() Payload {
	np := NewOutpostPayload(p.Body().Clone(), p.allowLong)
	np.SetRequestDR(p.RequestDR())
	np.SetRequestRR(p.RequestRR())
	np.SetUrgent(p.Urgent())
	np.bbsRoutes = p.BBSRoutes()
	return np
}

// encodedOutpostFlags returns the string encoding of the Outpost body envelope
// flags.
func (p *OutpostPayload) encodedOutpostFlags() (encoded string) {
	if p.urgent {
		encoded += "!URG!"
	}
	if p.requestDR {
		encoded += "!RDR!"
	}
	if p.requestRR {
		encoded += "!RRR!"
	}
	return encoded
}

// Decode decodes the supplied encoded payload and the body within it.  If
// successful, it returns the decoded payload, plus a possibly non-nil error
// that represents non-fatal errors decoding the payload or its body.  If
// unsuccessful, it returns a nil Payload and a non-nil error.
//
// The headers passed in may come from either the message headers or the headers
// of a part of a MIME multipart message.
func decodeOutpostPayload(headers mail.Header, payload string) (_ Payload, err error) {
	switch cte := headers.Get("Content-Transfer-Encoding"); cte {
	case "":
		// No Content-Transfer-Encoding.  If the message came from
		// WinLink and looks like it has quoted-printable encoding,
		// we'll decode that; otherwise, assume no encoding at all.
		// If the headers came from a part of a multipart message, there
		// will be no From header, but that's OK because we wouldn't
		// want that behavior for a multipart message anyway.
		if !strings.Contains(strings.ToLower(headers.Get("From")), "winlink") || !hasQuotedPrintable(payload) {
			break
		}
		fallthrough
	case "quoted-printable":
		qp := quotedprintable.NewReader(strings.NewReader(payload))
		if by, err := io.ReadAll(qp); err != nil {
			return nil, fmt.Errorf("decoding quoted-printable body: %w", err)
		} else {
			payload = string(by)
		}
	case "base64":
		if dec, err := base64.StdEncoding.DecodeString(payload); err != nil {
			return nil, fmt.Errorf("decoding base64-encoded body: %w", err)
		} else {
			payload = string(dec)
		}
	default:
		return nil, errors.NewF("This message has an unsupported Content-Transfer-Encoding %q.", cte)
	}
	p := new(OutpostPayload)
	payload = p.extractBBSRoutes(payload)
	if payload, err = p.decodeOutpostFlags(payload); err != nil {
		return nil, err
	}
	if p.body, err = body.Decode(payload); p.body == nil {
		return nil, err
	}
	p.body.OnDirty(p.MarkDirty)
	return p, err
}

// hasQuotedPrintable implements a heuristic to determine whether the wrapper
// is encoded with quoted-printable encoding.
func hasQuotedPrintable(wrapper string) bool {
	// We can't use Outpost's heuristic, which relies on CRLF line endings,
	// so we use our own instead.
	// First:  if there are no equals signs, it's not quoted-printable.
	if !strings.ContainsRune(wrapper, '=') {
		return false
	}
	// Second:  it's not QP if there are any nonASCII characters in it.
	if strings.ContainsFunc(wrapper, shouldBeQuotedInQP) {
		return false
	}
	// Third:  every '=' must be followed by newline or two hex digits.
	equals := strings.IndexByte(wrapper, '=')
	for equals >= 0 {
		switch {
		case len(wrapper) > equals+1 && wrapper[equals+1] == '\n':
			// nothing
		case len(wrapper) > equals+2 && isQPHex(wrapper[equals+1]) && isQPHex(wrapper[equals+2]):
			// nothing
		default:
			return false
		}
		if idx := strings.IndexByte(wrapper[equals+1:], '='); idx >= 0 {
			equals += idx + 1
		} else {
			break
		}
	}
	// Fourth:  lines must be no more than 76 characters long, and must not
	// end with space or tab.
	if exceedsLineLength(wrapper, 76, true) {
		return false
	}
	// Passed all checks.  We'll assume it's quoted-printable.
	_ = 0
	return true
}

// shouldBeQuotedInQP returns whether the supplied character should be quoted
// (i.e., is illegal) in a quoted-printable encoding.
func shouldBeQuotedInQP(r rune) bool {
	// Any control character except TAB, LF, CR, and any 8-bit character.
	return r < 9 || r == 11 || r == 12 || (r >= 14 && r <= 31) || r > 126
}

// isQPHex returnns whether the supplied character is a valid hex digit in a
// quoted-printable encoding.
func isQPHex(b byte) bool {
	// By spec, only uppercase hex digits are accepted.
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F')
}

var bbsRoutesRE = regexp.MustCompile(`^(?:R:\d{6}/\d{4}[zZ]? .*\n)+\n`)

// extractBBSRoutes removes any BBS routing headers from the top of the body.
func (p *OutpostPayload) extractBBSRoutes(payload string) string {
	if match := bbsRoutesRE.FindString(payload); match != "" {
		p.bbsRoutes = match[:len(match)-1]
		return payload[len(match):]
	}
	return payload
}

// decodeOutpostFlags decodes and removes Outpost flags from the payload.
func (p *OutpostPayload) decodeOutpostFlags(payload string) (decoded string, err error) {
	var (
		found    bool
		original = payload
	)
	for {
		switch {
		case strings.HasPrefix(payload, "\n"):
			// Remove newlines that might precede the first Outpost
			// code.  We'll keep this removal only if we actually
			// find an Outpost code.
			payload = payload[1:]
		case strings.HasPrefix(payload, "!B64!"):
			// Message content has Base64 encoding.
			if dec, err := base64.StdEncoding.DecodeString(payload[5:]); err != nil {
				return "", fmt.Errorf("unwrapping !B64! body: %w", err)
			} else {
				payload = string(dec)
				found = true
			}
		case strings.HasPrefix(payload, "!RRR!"):
			p.requestRR = true
			found = true
			payload = payload[5:]
		case strings.HasPrefix(payload, "!RDR!"):
			p.requestDR = true
			found = true
			payload = payload[5:]
		case strings.HasPrefix(payload, "!URG!"):
			p.urgent = true
			found = true
			payload = payload[5:]
		default:
			if !found {
				return original, nil
			}
			return payload, nil
		}
	}
}

// Encode encodes the payload (and its included body) according to the supplied
// bitmask of encodings.  It returns the encoded payload string.
func (p *OutpostPayload) Encode() (payload string) {
	// First, we need the encoded body.
	payload = p.body.EncodedBody()
	// Next, we'll prefix it with any Outpost body envelope flags.
	if flags := p.encodedOutpostFlags(); flags != "" {
		payload = flags + payload
	}
	p.MarkClean()
	// Next, encode the body.
	if !isJNOSSafe(payload, p.allowLong) {
		payload = "!B64!" + lineBreakEvery76(base64.StdEncoding.EncodeToString([]byte(payload)))
	}
	// If we have BBS routes, add them on top.
	if p.bbsRoutes != "" {
		payload = p.bbsRoutes + "\n" + payload
	}
	return payload
}

// isJNOSSafe returns whether the body can be sent to JNOS safely.  JNOS does
// not handle control characters, isn't guaranteed to handle high-bit
// characters, and silently breaks lines greater than 126 bytes.
func isJNOSSafe(body string, allowLong bool) bool {
	if strings.ContainsFunc(body, isNonASCII) {
		return false
	}
	if !allowLong && exceedsLineLength(body, 126, false) {
		return false
	}
	return true
}

func isNonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

// exceedsLineLength returns whether any line of the string has a length
// exceeding max.  If noTrailingWhitespace is set, it also returns true if any
// line ends with a tab or space character.
func exceedsLineLength(s string, max int, noTrailingWhitespace bool) bool {
	linestart := 0
	lineend := strings.IndexByte(s, '\n')
	if lineend < 0 {
		lineend = len(s)
	}
	for {
		if lineend-linestart > max {
			return true
		}
		if lineend > 0 && noTrailingWhitespace {
			if c := s[lineend-1]; c == '\t' || c == ' ' {
				return true
			}
		}
		if lineend == len(s) {
			return false
		}
		linestart = lineend + 1
		lineend = strings.IndexByte(s[linestart:], '\n')
		if lineend < 0 {
			lineend = len(s)
		} else {
			lineend += linestart
		}
	}
}

func lineBreakEvery76(s string) string {
	if s == "" {
		return s
	}
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	lines2 := make([]string, 0, len(lines))
	for _, line := range lines {
		for len(line) > 76 {
			lines2 = append(lines2, line[:76])
			line = line[76:]
		}
		if len(line) != 0 {
			lines2 = append(lines2, line)
		}
	}
	return strings.Join(lines2, "\r\n") + "\r\n"
}
