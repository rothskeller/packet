package message

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

// A JustReceivedMessage is for a message that the local system received from a
// BBS during the current session.  It has all of the methods of a
// ReceivedMessage, plus a few more.
type JustReceivedMessage struct {
	ReceivedMessage
	autoresponse bool
	returnAddr   string
	bbsRxDate    time.Time
	isMultipart  bool
}

var _ Message = (*JustReceivedMessage)(nil)

func (m *JustReceivedMessage) receivedMessage() *ReceivedMessage { return &m.ReceivedMessage }

// Autoresponse returns whether the message was sent as an auto-response (m.g.,
// a bounce message, a vacation responder, etc.).  Note that false negatives
// are possible depending on the message retrieval mechanics.
func (m *JustReceivedMessage) Autoresponse() bool { return m.autoresponse }

// ReturnAddr returns the return address of the message.  This is not always
// the same as the From address(es).
func (m *JustReceivedMessage) ReturnAddr() string { return m.returnAddr }

// BBSRxDate returns the timestamp of receipt of the message by the BBS that we
// retrieved it from.
func (m *JustReceivedMessage) BBSRxDate() time.Time { return m.bbsRxDate }

// IsMultipart returns whether the message, as received from the BBS, was a
// multipart message.  (The message Body will contain only the plain text
// portion of the received body.)
func (m *JustReceivedMessage) IsMultipart() bool { return m.isMultipart }

// SetLocalID sets the local ID of the message.
func (m *JustReceivedMessage) SetLocalID(localID string) { m.localID = localID }

// NewJustReceivedMessage creates a new JustReceivedMessage by decoding the
// provided retrieved message (i.m., the output of a JNOS "R" or "V" command).
// rxBBS specifies the BBS from which the message was retrieved, if known.
// rxArea specifies the bulletin area from which the message was retrieved, and
// will be an empty string for non-bulletin messages.  The function returns an
// error if the message cannot be parsed.
func NewJustReceivedMessage(retrieved, rxBBS, rxArea string) (m *JustReceivedMessage, err error) {
	var (
		efrom      string
		hadMessage bool
		by         []byte
		msg        *mail.Message
		to         []string
		headers    textproto.MIMEHeader
		issues     error
	)
	m = &JustReceivedMessage{ReceivedMessage: ReceivedMessage{
		common: new(common),
		rxBBS:  rxBBS,
		rxArea: rxArea,
		rxDate: time.Now(),
	}}
	// If there is an envelope From line, remove it from the raw message.
	if strings.HasPrefix(retrieved, "From ") {
		hadMessage = true
		if idx := strings.IndexByte(retrieved, '\n'); idx > 0 {
			efrom, retrieved = retrieved[5:idx], retrieved[idx+1:]
		}
		if idx := strings.IndexByte(efrom, ' '); idx >= 0 {
			m.returnAddr, efrom = efrom[:idx], efrom[idx+1:]
			// Looks like there's a timestamp on the envelope line.
			// RFC-4155 says it should be a ctime-style timestamp
			// in UTC.  We'll try parsing it that way, but we'll
			// treat it as local time because that's what JNOS
			// BBSes do, and those are our primary source of
			// messages to parse.
			if t, err := time.ParseInLocation(time.ANSIC, efrom, time.Local); err == nil {
				m.bbsRxDate = t
			}
		} else {
			m.returnAddr, efrom = efrom, ""
		}
		if m.returnAddr == "" {
			m.autoresponse = true
		}
	}
	// Parse the message headers.  If we can't parse them, we go no further.
	if msg, err = mail.ReadMessage(strings.NewReader(retrieved)); err != nil {
		return nil, err
	}
	// Handle the From header.
	m.from = strings.Join(msg.Header["From"], ", ")
	// Handle the To, Cc, and Bcc headers.
	to = msg.Header["To"]
	to = append(to, msg.Header["Cc"]...)
	to = append(to, msg.Header["Bcc"]...)
	m.to = strings.Join(to, ", ")
	// Handle the Date header.
	if t, err := mail.ParseDate(msg.Header.Get("Date")); err == nil {
		m.date = t
	}
	// Compute the return address if there wasn't an envelope From line.
	if !hadMessage {
		var line string
		if line = msg.Header.Get("Return-Path"); line == "" {
			if line = msg.Header.Get("Reply-To"); line == "" {
				if line = msg.Header.Get("Sender"); line == "" {
					line = msg.Header.Get("From")
				}
			}
		}
		// Most of those sources can have a name comment in the
		// address, which we don't want.  Also, From can have more than
		// one address in it, and we only want the first.
		if addrs, err := address.ParseList(line); err == nil && len(addrs) > 0 {
			m.returnAddr = addrs[0].Address
		}
	}
	// If we didn't get a BBS Rx date from the envelope, get it from the
	// Received header.
	if m.bbsRxDate.IsZero() {
		_, date, _ := strings.Cut(msg.Header.Get("Received"), ";")
		if t, err := mail.ParseDate(strings.TrimSpace(date)); err == nil {
			m.bbsRxDate = t
		}
	}
	// Extract the plain text portion of the body.
	if by, err = io.ReadAll(msg.Body); err != nil {
		return nil, fmt.Errorf("body: %w", err)
	}
	if by, headers, m.isMultipart, err = extractPlainText(textproto.MIMEHeader(msg.Header), by); err != nil {
		return m, err
	}
	// Unwrap and parse the payload.
	if m.payload, issues = payload.Decode(mail.Header(headers), string(by)); m.payload == nil {
		m.payload, _ = payload.Decode(nil, "")
		return m, issues
	}
	// Handle the Subject header.
	m.subject = subject.DecodePlainSubject(msg.Header.Get("Subject"))
	// Set the message type.
	SetType(m)
	// Nothing should change from this point on.
	m.OnDirty(func(reason string) {
		switch reason {
		case "body.FormBody.Field.DestMsgNo",
			"body.FormBody.Field.OpName",
			"body.FormBody.Field.OpCall",
			"body.FormBody.Field.OpDate",
			"body.FormBody.Field.OpTime",
			"body.FormBody.Field.RecSent",
			"body.FormBody.Field.Method",
			"body.FormBody.Field.Other",
			"body.FormBody.Common.destinationMessageID",
			"body.FormBody.Common.operatorName",
			"body.FormBody.Common.operatorCall",
			"body.FormBody.Common.operatorDate",
			"body.FormBody.Common.operatorTime",
			"body.FormBody.Common.receiverSender",
			"body.FormBody.Common.operatorMethod",
			"body.FormBody.Common.operatorMethodOther":
			// OK
		default:
			panic("JustReceivedMessage should not change: " + reason)
		}
	})
	return m, issues
}

// extractPlainText extracts the plain text portion of a message from its body.
// If the body is a multipart body, the last plain text portion is extracted.
// The function returns nil, and possibly an error, if there is no plain text
// portion.  The returned headers are the ones for the last plain text portion
// of the multipart body, or the argument header if the body was not multipart.
// The returned isMultipart flag indicates whether the plain text portion was
// found inside a multipart container.  This is a recursive function, to handle
// nested multipart bodies.
func extractPlainText(header textproto.MIMEHeader, body []byte) (nbody []byte, actual textproto.MIMEHeader, isMultipart bool, err error) {
	var (
		mediatype string
		params    map[string]string
		mr        *multipart.Reader
		part      *multipart.Part
		partbody  []byte
		found     []byte
		foundhdr  textproto.MIMEHeader
		parterrs  error
	)
	// Decode the content type.
	if ct := header.Get("Content-Type"); ct != "" {
		if mediatype, params, err = mime.ParseMediaType(ct); err != nil {
			return nil, nil, false, fmt.Errorf("invalid Content-Type %q", ct) // can't decode Content-Type
		}
	} else {
		mediatype, params = "text/plain", map[string]string{}
	}
	// If the message part is plain text, we're golden.
	if mediatype == "text/plain" {
		return body, header, false, nil
	}
	// The only other content type we can handle is multipart.
	if !strings.HasPrefix(mediatype, "multipart/") {
		return nil, nil, false, fmt.Errorf("unknown Content-Type %q", mediatype)
	}
	mr = multipart.NewReader(bytes.NewReader(body), params["boundary"])
	for {
		if part, err = mr.NextRawPart(); err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, true, fmt.Errorf("can't decode multipart body: %w", err)
		}
		partbody, _ = io.ReadAll(part)
		if plain, plainhdr, _, err := extractPlainText(part.Header, partbody); err != nil {
			parterrs = errors.Join(parterrs, err)
		} else if plain != nil {
			found, foundhdr = plain, plainhdr
		}
	}
	if found != nil {
		return found, foundhdr, true, nil
	} else {
		return nil, nil, true, parterrs
	}
}
