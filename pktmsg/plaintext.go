package pktmsg

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
)

// extractPlainText extracts the plain text portion of a message from its body.
// It returns a nil body if there is none.  The returned boolean indicates
// whether the entire body was plain text.  This is a recursive function, to
// handled nested multipart bodies.
func extractPlainText(header textproto.MIMEHeader, body []byte) (nbody []byte, notplain bool) {
	var (
		mediatype string
		params    map[string]string
		err       error
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
		return nil, true
	}
	// Decode the content type.
	if ct := header.Get("Content-Type"); ct != "" {
		if mediatype, params, err = mime.ParseMediaType(header.Get("Content-Type")); err != nil {
			return nil, true // can't decode Content-Type
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
				return nil, true // Can't decode multipart body
			}
			partbody, _ = io.ReadAll(part)
			if plain, _ := extractPlainText(part.Header, partbody); plain != nil {
				found = plain
			}
		}
		return found, true
	}
	// If the content type is anything other than text/plain, we're out of
	// luck.
	if mediatype != "text/plain" {
		return nil, true
	}
	// In theory we also ought to check the charset, but we'll elide that
	// until experience proves a need.
	return body, notplain
}
