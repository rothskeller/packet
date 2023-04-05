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
	"strings"
)

// ErrNoPlainTextBody is raised when a message being parsed has no plain text
// body content.
var ErrNoPlainTextBody = errors.New("message has no plain text body content")

// ParseMessage parses the supplied string as a stored message (which could be
// either incoming or outgoing).
func ParseMessage(raw string) (m Message, err error) {
	var (
		mm   *mail.Message
		body []byte
		tx   *baseTx
		om   *outpostMessage
	)
	// Parse the message headers.  If we can't parse them, we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(raw)); err != nil {
		return nil, err
	}
	// Stored messages should not have a Content-Type or Content-Transfer-Encoding.
	if mm.Header["Content-Type"] != nil || mm.Header["Content-Transfer-Encoding"] != nil {
		return nil, errors.New("stored messages must be plain text only")
	}
	// Extract the plain text portion of the body and create the baseTx.
	body, _ = io.ReadAll(mm.Body)
	tx = parseBaseTx(mm.Header, string(body))
	// Decode outpost flags.
	if om, err = parseOutpostFlags(tx); err != nil {
		return nil, err
	}
	// If the message was received, transform it into a baseRx.
	if recv := mm.Header.Get("Received"); recv != "" {
		if m, err = parseBaseRx(om, recv); err != nil {
			return nil, err
		}
	} else {
		m = om
	}
	// Parse the form if any.
	m = parsePIFOForm(m)
	return m, nil
}

// ReceiveMessage parses the supplied string as a message that was just
// retrieved from the specified JNOS BBS.  If it is a bulletin, area should be
// set to the bulletin area from which it was retrieved; otherwise, area should
// be empty.
func ReceiveMessage(raw, bbs, area string) (m Message, err error) {
	var (
		envelope string
		mm       *mail.Message
		body     []byte
		notplain bool
		tx       *baseTx
		om       *outpostMessage
		err2     error
	)
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
	body, notplain, err = extractPlainText(textproto.MIMEHeader(mm.Header), body)
	if err == nil && notplain && len(body) == 0 {
		err = ErrNoPlainTextBody
	}
	// Create the message.
	tx = parseBaseTx(mm.Header, string(body))
	if om, err2 = parseOutpostFlags(tx); err == nil {
		err = err2
	}
	m = parseBaseRetrieved(om, bbs, area, envelope, mm.Header, notplain)
	m = parsePIFOForm(m)
	return m, err
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
