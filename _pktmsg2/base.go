package pktmsg

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

// FieldUrgent is the tag of the BaseMessage field containing the message
// urgency flag; any non-empty value indicates that this message should be
// flagged as urgent (i.e., displayed in red) in Outpost.  (In Santa Clara
// County, that means messages with an IMMEDATE handling order.  But this
// package carefully doesn't know anything about Santa Clara County protocols.)
const FieldUrgent = "Urgent"

// BaseMessage represents a message with only the envelope metadata processed.
type BaseMessage struct {
	received *ReceivedInfo
	fieldList
}

// Type returns the message type.
func (m BaseMessage) Type() MessageType { return plainTextType }

// Received returns information about the reception of the message, or
// nil if it is an outgoing message.
func (m BaseMessage) Received() *ReceivedInfo { return m.received }

// Encode returns the RFC-5322-encoded form of the message, suitable for
// saving on disk.
func (m BaseMessage) Encode() string {
	var sb strings.Builder

	// Write the Received: header if appropriate.
	if m.received != nil {
		fmt.Fprintf(&sb, "Received: FROM %s BY pktmsg.local", m.received.BBS)
		linelen := len(m.received.BBS) + 64
		if m.received.Area != "" {
			fmt.Fprintf(&sb, " FOR %s", m.received.Area)
			linelen += len(m.received.Area) + 5
		}
		if linelen > 78 {
			sb.WriteString(";\n\t")
		} else {
			sb.WriteString("; ")
		}
		fmt.Fprintln(&sb, m.received.Retrieved.Format(time.RFC1123Z))
	}
	// Write the From header if appropriate.
	if f := m.KeyValue(FFromAddr); f != "" {
		fmt.Fprintf(&sb, "From: %s\n", f)
	}
	// Write the To header if appropriate.
	if f := m.KeyValue(FToAddrs); f != "" {
		fmt.Fprintf(&sb, "To: %s\n", f)
	}
	// Write the Subject header if appropriate.
	if f := m.KeyValue(FSubject); f != "" {
		fmt.Fprintf(&sb, "Subject: %s\n", f)
	}
	// Write the Date header if appropriate.
	if f := m.KeyValue(FSent); f != "" {
		fmt.Fprintf(&sb, "Date: %s\n", f)
	}
	// Newline to end the headers.
	fmt.Fprintln(&sb)
	// Compute the message body.
	var body = string(m.KeyValue(FBody))
	if m.FieldValue(FieldUrgent) != "" {
		body = "!URG!" + body
	}
	if strings.IndexFunc(body, nonASCII) >= 0 {
		body = "!B64!" + base64.StdEncoding.EncodeToString([]byte(body)) + "\n"
	}
	// Write the message body.
	fmt.Fprint(&sb, body)
	return sb.String()
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

// EncodeTx returns the message, encoded as needed for transmission via
// JNOS.  Specifically, it returns a list of destination addresses, an
// encoded subject line, and an encoded body.
func (m BaseMessage) EncodeTx() (to []string, subject string, body string) {
	to = splitTrim(string(m.KeyValue(FToAddrs)))
	subject = string(m.KeyValue(FSubject))
	body = string(m.KeyValue(FBody))
	return
}
func splitTrim(list string) (slice []string) {
	for _, item := range strings.Split(list, ",") {
		if item = strings.TrimSpace(item); item != "" {
			slice = append(slice, item)
		}
	}
	return slice
}

// fromLineRE is the regular expression that the "From " line at the beginning
// of an RFC-4155 "mbox" format email message is expected to match.  It has the
// word "From" followed by a space (not a colon, as in the RFC-5322 header);
// then a possibly-empty from address, and then optionally a space and a
// timestamp.
// var fromLineRE = regexp.MustCompile(`^From (\S*)(?: (.*))?\n`)
var fromLineRE = regexp.MustCompile(`^From (\S*)(?: (.*))?\n`)

// parseMessage parses a raw message and returns the corresponding BaseMessage.
// If the message cannot be parsed, parseError returns a partially filled out
// BaseMessage and an error.  bbs, callsign, and area should be set when parsing
// messages retrieved from JNOS; otherwise they should be empty.
func parseMessage(raw, bbs, callsign, area string) (Message, error) {
	var (
		bm       BaseMessage
		mm       *mail.Message
		hdr      textproto.MIMEHeader
		bodyb    []byte
		notplain bool
		err      error
	)
	// Initialize the fields.
	bm.fieldList = fieldList{
		newBMFromField(bm),
		newBMToField(bm),
		newBMSubjectField(bm),
		newBMSentField(bm),
		newBMBodyField(bm),
		newBMUrgentField(bm),
	}
	// Save the caller-supplied received information if any.
	if bbs != "" {
		bm.received = &ReceivedInfo{BBS: bbs, CallSign: callsign, Area: area}
	}
	// Extract information from the RFC-4155 "From " envelope line if any.
	// If we have such a line but it has no return address on it, that's an
	// auto-responder message.
	if match := fromLineRE.FindStringSubmatch(raw); match != nil {
		if bm.received == nil {
			bm.received = new(ReceivedInfo)
		}
		bm.received.ReturnAddress = match[1]
		bm.received.Arrived, _ = time.ParseInLocation(time.ANSIC, match[2], time.Local)
		raw = raw[len(match[0]):]
	}
	// Parse the message headers.  If we can't parse them, it's an
	// unparseable message and we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(raw)); err != nil {
		return bm, err
	}
	// Save values from the headers.
	if f := mm.Header.Get("From"); f != "" {
		if addr, err := mail.ParseAddress(f); err == nil {
			bm.FieldByKey(FFromAddr).SetValue(Value(addr.Address))
		} else {
			bm.FieldByKey(FFromAddr).SetValue(Value(f))
		}
	}
	var to = strings.Join(mm.Header["To"], ", ")
	to += ", " + strings.Join(mm.Header["Cc"], ", ")
	var to2 []string
	for _, t := range splitTrim(to) {
		if addr, err := mail.ParseAddress(t); err == nil {
			to2 = append(to2, addr.Address)
		} else {
			to2 = append(to2, t)
		}
	}
	if len(to2) != 0 {
		bm.FieldByKey(FToAddrs).SetValue(Value(strings.Join(to2, ", ")))
	}
	if f := mm.Header.Get("Subject"); f != "" {
		bm.FieldByKey(FSubject).SetValue(Value(f))
	}
	if f := mm.Header.Get("Date"); f != "" {
		bm.FieldByKey(FSent).SetValue(Value(f))
	}
	// Extract the plain text portion of the body, and decode it.
	bodyb, _ = io.ReadAll(mm.Body)
	if bodyb, notplain, err = extractPlainText(hdr, bodyb); err != nil {
		return bm, err
	}
	bm.FieldByKey(FBody).SetValue(Value(bodyb))
	if err = parseOutpostFlags(bm); err != nil {
		return bm, err
	}
	if notplain {
		return bm, errors.New("not plain text")
	}
	return bm, nil

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

// parseOutpostFlags looks for Outpost flags at the beginning of the message
// body, and handles them if they are found.
func parseOutpostFlags(bm BaseMessage) error {
	var (
		found bool
		body  = string(bm.KeyValue(FBody))
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
				bm.received.Flags |= NotPlainText
			} else {
				return err
			}
		case strings.HasPrefix(body, "!RRR!"):
			bm.received.Flags |= RequestReadReceipt
			body = body[5:]
			found = true
		case strings.HasPrefix(body, "!RDR!"):
			bm.received.Flags |= RequestDeliveryReceipt
			body = body[5:]
			found = true
		case strings.HasPrefix(body, "!URG!"):
			bm.received.Flags |= OutpostUrgent
			body = body[5:]
			found = true
		default:
			if found {
				bm.FieldByKey(FBody).SetValue(Value(body))
				// If we didn't find any Outpost codes, leave
				// the body unchanged so as not to remove
				// initial newlines that might be significant.
			}
			return nil
		}
	}
}
