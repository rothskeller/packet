package envelope

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
)

// now returns the current time.  It can be overridden by tests.
var now = time.Now

// ErrNoPlainTextBody is raised when a message being parsed has no plain text
// body content.
var ErrNoPlainTextBody = errors.New("message has no plain text body content")

// ParseSaved parses the supplied string as a saved message (which could be
// either incoming or outgoing).
func ParseSaved(saved string) (_ *Envelope, body string, err error) {
	var (
		env  Envelope
		mm   *mail.Message
		bbuf []byte
	)
	// Parse the message headers.  If we can't parse them, we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(saved)); err != nil {
		return nil, "", err
	}
	// Saved messages should not have a Content-Type or Content-Transfer-Encoding.
	if mm.Header["Content-Type"] != nil || mm.Header["Content-Transfer-Encoding"] != nil {
		return nil, "", errors.New("saved messages must be plain text only")
	}
	// Create the message.
	bbuf, _ = io.ReadAll(mm.Body)
	if err = env.parseHeadersStored(mm.Header); err != nil {
		return nil, "", err
	}
	if body, err = env.parseOutpost(string(bbuf)); err != nil {
		return nil, "", err
	}
	return &env, body, nil
}

// ParseRetrieved parses the supplied string as a message that was just
// retrieved from the specified JNOS BBS.  If it is a bulletin, area should be
// set to the bulletin area from which it was retrieved; otherwise, area should
// be empty.
func ParseRetrieved(retrieved, bbs, area string) (_ *Envelope, body string, err error) {
	var (
		env   Envelope
		efrom string
		mm    *mail.Message
		bbuf  []byte
		err2  error
	)
	env.ReceivedBBS, env.ReceivedArea = bbs, area
	env.ReceivedDate = now()
	// If there is an envelope From line, remove it from the raw message.
	if strings.HasPrefix(retrieved, "From ") {
		if idx := strings.IndexByte(retrieved, '\n'); idx > 0 {
			efrom, retrieved = retrieved[:idx], retrieved[idx+1:]
		}
	}
	// Parse the message headers.  If we can't parse them, we go no further.
	if mm, err = mail.ReadMessage(strings.NewReader(retrieved)); err != nil {
		return &env, "", err
	}
	// Extract the plain text portion of the body.
	bbuf, _ = io.ReadAll(mm.Body)
	bbuf, env.NotPlainText, err = extractPlainText(textproto.MIMEHeader(mm.Header), bbuf)
	if err == nil && env.NotPlainText && len(bbuf) == 0 {
		err = ErrNoPlainTextBody
	}
	// Create the message.
	env.parseHeadersRetrieved(mm.Header, efrom)
	if body, err2 = env.parseOutpost(string(bbuf)); err == nil {
		err = err2
	}
	return &env, body, err
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

// parseHeadersSaved parses the headers of a saved message.
func (env *Envelope) parseHeadersStored(h mail.Header) error {
	env.parseHeadersCommon(h)
	if recv := h.Get("Received"); recv != "" {
		if match := receivedRE.FindStringSubmatch(recv); match != nil {
			env.ReceivedBBS = match[1]
			env.ReceivedArea = match[2]
			env.ReceivedDate, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", match[3])
		} else {
			// This shouldn't happen:  stored messages with a Received: header
			// should always have our Received: header format
			return errors.New("incorrect Received: header format for stored received message")
		}
	}
	env.ReadyToSend = h.Get("X-Packet-Queued") != ""
	return nil
}

// parseHeadersRetrieved parses the headers of a retrieved message.
func (env *Envelope) parseHeadersRetrieved(h mail.Header, efrom string) {
	env.parseHeadersCommon(h)
	// Parse the envelope "From " line if any.
	var hadEnvelope bool
	if efrom != "" {
		hadEnvelope = true
		efrom = efrom[5:] // skip "From "
		if idx := strings.IndexByte(efrom, ' '); idx >= 0 {
			env.ReturnAddr, efrom = efrom[:idx], efrom[idx+1:]
		} else {
			env.ReturnAddr, efrom = efrom, ""
		}
		if env.ReturnAddr == "" {
			env.Autoresponse = true
		}
	}
	if efrom != "" {
		// Looks like there's a timestamp on the envelope line.
		// RFC-4155 says it should be a ctime-style timestamp in UTC.
		// We'll try parsing it as that way, but we'll treat it as local
		// time because that's what JNOS BBSes do, and those are our
		// primary source of messages to parse.
		if t, err := time.ParseInLocation(time.ANSIC, efrom, time.Local); err == nil {
			env.BBSReceivedDate = t
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
		if addrs, err := ParseAddressList(line); err == nil && len(addrs) > 0 {
			env.ReturnAddr = addrs[0].Address
		}
	}
	// If we didn't get a BBS Rx date from the envelope, get it from the
	// Received header.
	if env.BBSReceivedDate.IsZero() {
		_, date, _ := strings.Cut(h.Get("Received"), ";")
		if t, err := mail.ParseDate(strings.TrimSpace(date)); err == nil {
			env.BBSReceivedDate = t
		}
	}
}

// parseHeadersCommon is the code that is common between parseHeadersStored and
// parseHeadersRetrieved.
func (env *Envelope) parseHeadersCommon(h mail.Header) {
	env.From = h.Get("From")
	if addrs, err := ParseAddressList(env.From); err == nil && len(addrs) >= 1 {
		env.From = addrs[0].String()
	}
	for _, list := range h["To"] {
		if addrs, err := ParseAddressList(list); err == nil {
			for _, addr := range addrs {
				env.To = append(env.To, addr.String())
			}
		}
	}
	for _, list := range h["Cc"] {
		if addrs, err := ParseAddressList(list); err == nil {
			for _, addr := range addrs {
				env.To = append(env.To, addr.String())
			}
		}
	}
	for _, list := range h["Bcc"] {
		if addrs, err := ParseAddressList(list); err == nil {
			for _, addr := range addrs {
				env.To = append(env.To, addr.String())
			}
		}
	}
	if t, err := mail.ParseDate(h.Get("Date")); err == nil {
		env.Date = t
	}
	env.SubjectLine = h.Get("Subject")
}

// parseOutpost handles the Outpost codes at the start of the body, if any.
func (env *Envelope) parseOutpost(original string) (string, error) {
	var (
		found bool
		body  = original
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
				return "", err
			}
		case strings.HasPrefix(body, "!RRR!"):
			env.RequestReadReceipt, found = true, true
			body = body[5:]
		case strings.HasPrefix(body, "!RDR!"):
			env.RequestDeliveryReceipt, found = true, true
			body = body[5:]
		case strings.HasPrefix(body, "!URG!"):
			env.OutpostUrgent, found = true, true
			body = body[5:]
		default:
			if !found {
				// If we didn't find any Outpost codes, leave
				// the body unchanged so as not to remove
				// initial newlines that might be significant.
				body = original
			}
			return body, nil
		}
	}
}
