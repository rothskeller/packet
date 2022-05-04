package pktmsg

import (
	"bytes"
	"encoding/base64"
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

// fromLineRE is the regular expression that the "From " line at the beginning
// of an RFC-4155 "mbox" format email message is expected to match.  It has the
// word "From" followed by a space (not a colon, as in the RFC-5322 header);
// then a possibly-empty from address, and then optionally a space and a
// timestamp.
var fromLineRE = regexp.MustCompile(`^From (\S*)(?: (.*))?\n`)

// ParseMessage parses an encoded message.  It returns an error if the message
// could not be decoded or parsed.  Even if it returns a non-nil error, the
// returned Message will be non-nil and may contain incomplete message data.
func ParseMessage(rawmsg string) (msg *Message, err error) {
	var (
		loc      []int
		mm       *mail.Message
		body     []byte
		notplain bool
	)
	msg = new(Message)
	// Extract information from the RFC-4155 "From " envelope line if any.
	// If we have such a line but it has no return address on it, that's an
	// auto-responder message.
	if loc = fromLineRE.FindStringSubmatchIndex(rawmsg); loc != nil {
		if msg.EnvelopeAddress = rawmsg[loc[2]:loc[3]]; msg.EnvelopeAddress == "" {
			msg.Flags |= AutoResponse
		}
		if loc[4] >= 0 {
			// Looks like there's a timestamp on the envelope line.
			// RFC-4155 says it should be a ctime-style timestamp in
			// UTC.  We'll try parsing it as that way, but we'll
			// treat it as local time because that's what JNOS BBSes
			// do, and those are our primary source of messages to
			// parse.
			msg.EnvelopeDate, _ = time.ParseInLocation(time.ANSIC, rawmsg[loc[4]:loc[5]], time.Local)
		}
		rawmsg = rawmsg[loc[1]:]
	}
	// Parse the message headers.  If we can't parse them, it's an
	// unparseable message and we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(rawmsg)); err != nil {
		return msg, err
	}
	msg.Header = textproto.MIMEHeader(mm.Header)
	// Extract the plain text portion of the body, and decode it.
	body, _ = io.ReadAll(mm.Body)
	body, notplain, err = extractPlainText(textproto.MIMEHeader(mm.Header), body)
	if err != nil {
		return msg, err
	}
	if notplain {
		msg.Flags |= NotPlainText
	}
	msg.Body = string(body)
	// Handle Outpost flags.
	err = msg.parseOutpostFlags()
	return msg, err
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
		body, _ = io.ReadAll(quotedprintable.NewReader(bytes.NewReader(body)))
		notplain = true
	case "base64":
		body, _ = io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(body)))
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
func (msg *Message) parseOutpostFlags() error {
	var (
		found bool
		body  = msg.Body
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
				msg.Flags |= NotPlainText
			} else {
				return err
			}
		case strings.HasPrefix(body, "!RRR!"):
			msg.Flags |= RequestReadReceipt
			body = body[5:]
			found = true
		case strings.HasPrefix(body, "!RDR!"):
			msg.Flags |= RequestDeliveryReceipt
			body = body[5:]
			found = true
		case strings.HasPrefix(body, "!URG!"):
			msg.Flags |= OutpostUrgent
			body = body[5:]
			found = true
		default:
			if found {
				msg.Body = body
				// If we didn't find any Outpost codes, leave
				// the body unchanged so as not to remove
				// initial newlines that might be significant.
			}
			return nil
		}
	}
}
