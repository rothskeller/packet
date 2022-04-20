package pktmsg

// This file defines TxBase and RxBase.

import (
	"encoding/base64"
	"errors"
	"io"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"
	"time"
)

// TxBase is the foundation for all outgoing messages.
type TxBase struct {
	// SubjectLine is the Subject line for the message.
	SubjectLine string
	// RequestDeliveryReceipt is true if the message should contain an
	// Outpost-encoded request for a delivery receipt.
	RequestDeliveryReceipt bool
	// RequestReadReceipt is true if the message should contain an
	// Outpost-encoded request for a read receipt.
	RequestReadReceipt bool
	// OutpostUrgent is true if the message should contain an
	// Outpost-encoded "urgent" flag.
	OutpostUrgent bool
	// Body is the plain-text body of the message.
	Body string
}

// OutboundMessage is an interface satisfied by any outbound (Tx*) message type.
type OutboundMessage interface {
	Encode() (subject, body string, err error)
}

// ErrIncomplete is the error returned by Encode when required fields are not
// set.
var ErrIncomplete = errors.New("required message field(s) not set")

// ErrInvalid is the error returned by Encode when fields are set to invalid
// values.
var ErrInvalid = errors.New("message field(s) have invalid values")

// ErrDontSet is the error returned by Encode when fields are set that shouldn't
// be set.
var ErrDontSet = errors.New("computed message field(s) must not have explicitly set values")

// Encode returns the encoded subject line and body of the message.
func (b *TxBase) Encode() (subject, body string, err error) {
	if b.SubjectLine == "" {
		return "", "", ErrIncomplete
	}
	if b.RequestReadReceipt || b.RequestDeliveryReceipt || b.OutpostUrgent {
		var sb strings.Builder
		if b.OutpostUrgent {
			sb.WriteString("!URG!")
		}
		if b.RequestDeliveryReceipt {
			sb.WriteString("!RDR!")
		}
		if b.RequestReadReceipt {
			sb.WriteString("!RRR!")
		}
		sb.WriteString(b.Body)
		body = sb.String()
	} else {
		body = b.Body
	}
	if strings.IndexFunc(body, nonASCII) >= 0 {
		body = "!B64!" + base64.StdEncoding.EncodeToString([]byte(body)) + "\n"
	}
	return b.SubjectLine, body, nil
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}

//------------------------------------------------------------------------------

// RxBase is the foundation for all received messages.
type RxBase struct {
	// RawMessage is the original message that was parsed.
	RawMessage string
	// ParseError is an error message describing why the message could not
	// be parsed.  It is empty for successfully-parsed messages.
	ParseError string
	// ReturnAddress is the return address of the message.
	ReturnAddress string
	// DeliveryTime indicates when the message was delivered to the BBS.
	DeliveryTime time.Time
	// SubjectLine is the Subject: header of the message.
	SubjectLine string
	// DateLine is the Date: header of the message.
	DateLine string
	// NotPlainText is true if the message is not entirely plain text.
	NotPlainText bool
	// RequestDeliveryReceipt is true if the message contains an
	// Outpost-encoded request for a delivery receipt.
	RequestDeliveryReceipt bool
	// RequestReadReceipt is true if the message contains an
	// Outpost-encoded request for a read receipt.
	RequestReadReceipt bool
	// OutpostUrgent is true if the message contains an Outpost-encoded
	// "urgent" flag.
	OutpostUrgent bool
	// Body is the plain text portoin of the message body.
	Body string
}

// Base returns a pointer to the RxBase portion of a message object.  It can be
// used to reach fields of the RxBase object that are occluded by types that
// embed RxBase.
func (b *RxBase) Base() *RxBase { return b }

// Message returns a pointer to the RxMessage portion of a message object.  It
// can be used to reach fields of the RxMessage object that are occluded by
// types that embed RxMessage.  It returns nil for messages that don't embed an
// RxMessage.
func (b *RxBase) Message() *RxMessage { return nil }

// Form returns a pointer to the RxForm portion of a message object.  It can be
// used to reach fields of the RxForm object that are occluded by types that
// embed RxForm.  It returns nil for messages that don't embed an RxForm.
func (b *RxBase) Form() *RxForm { return nil }

// SCCoForm returns a pointer to the RxSCCoForm portion of a message object.  It
// can be used to reach fields of the RxSCCoForm object that are occluded by
// types that embed RxSCCoForm.  It returns nil for messages that don't embed an
// RxSCCoForm.
func (b *RxBase) SCCoForm() *RxSCCoForm { return nil }

// TypeCode returns the machine-readable code for the message type.
func (b *RxBase) TypeCode() string {
	if b.ParseError == "" {
		return "BOUNCE"
	}
	return "CORRUPT"
}

// TypeName returns the human-reading name of the message type.
func (b *RxBase) TypeName() string {
	if b.ParseError == "" {
		return "auto-response message"
	}
	return "unparseable message"
}

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxBase) TypeArticle() string { return "an" }

// ParsedMessage is an interface satisfied by all message types returned from
// ParseMessage.
type ParsedMessage interface {
	Base() *RxBase
	Message() *RxMessage
	Form() *RxForm
	SCCoForm() *RxSCCoForm
	TypeCode() string
	TypeName() string
	TypeArticle() string
}

// fromLineRE is the regular expression that the "From " line at the beginning
// of an RFC-4155 "mbox" format email message is expected to match.  It has the
// word "From" followed by a space (not a colon, as in the RFC-5322 header);
// then a possibly-empty from address, and then optionally a space and a
// ctime-style timestamp.  Although RFC-4155 says the timestamp should be UTC,
// we treat it as local time because that's what JNOS BBSes do, and those are
// our primary source of messages to parse.
var fromLineRE = regexp.MustCompile(`^From (\S*)(?: (... ... .. ..:..:.. ....))?\n`)

// parseRxBase parses an incoming message and returns a filled-in RxBase.  If
// the ParseError field is non-empty, the RxBase is unusable.  If the Body field
// is empty and the NotPlainText flag is set, the message did not contain any
// plain text portion.
func parseRxBase(rawmsg string) *RxBase {
	var (
		b    RxBase
		loc  []int
		mm   *mail.Message
		body []byte
		err  error
	)
	b.RawMessage = rawmsg
	if loc = fromLineRE.FindStringSubmatchIndex(rawmsg); loc != nil {
		// The raw message starts with an RFC-4155 "From " envelope
		// line.  Extract the data from it and move past it.
		b.ReturnAddress = rawmsg[loc[2]:loc[3]]
		if loc[4] >= 0 {
			b.DeliveryTime, _ = time.ParseInLocation(time.ANSIC, rawmsg[loc[4]:loc[5]], time.Local)
		}
		rawmsg = rawmsg[loc[1]:]
	}
	if mm, err = mail.ReadMessage(strings.NewReader(rawmsg)); err != nil {
		b.ParseError = err.Error()
		return &b
	}
	if b.RawMessage == rawmsg {
		// We didn't find an RFC-4155 "From " envelope line above, so
		// extract the return address information from the message headers.
		b.ReturnAddress = returnAddressFromHeaders(mm.Header)
	}
	if b.DeliveryTime.IsZero() {
		// We didn't get a delivery time from the "From " envelope line,
		// either because there was no envelope line or because it
		// didn't contain a delivery time (that we could parse).  So,
		// take the time out of the first "Received:" header.
		received := mm.Header.Get("Received")
		if semi := strings.LastIndexByte(received, ';'); semi >= 0 {
			b.DeliveryTime, _ = mail.ParseDate(received[semi+1:])
		}
	}
	// Save the subject and date lines.
	b.SubjectLine = mm.Header.Get("Subject")
	b.DateLine = mm.Header.Get("Date")
	// Extract the plain text portion of the body, and decode it.
	body, _ = io.ReadAll(mm.Body)
	body, b.NotPlainText = extractPlainText(textproto.MIMEHeader(mm.Header), body)
	b.Body = string(body)
	// Handle any Outpost codes at the start of the message.
	b.parseOutpostCodes()
	return &b
}

// parseOutpostCodes looks for Outpost codes at the start of the message body.
// If it finds them, it removes them from the body and makes the appropriate
// notation in the RxBase structure.
func (b *RxBase) parseOutpostCodes() {
	for {
		switch {
		case strings.HasPrefix(b.Body, "\n"):
			b.Body = b.Body[1:]
		case strings.HasPrefix(b.Body, "!B64!"):
			if dec, err := base64.StdEncoding.DecodeString(b.Body[5:]); err == nil {
				b.Body = string(dec)
			} else {
				b.Body = ""
			}
			b.NotPlainText = true
		case strings.HasPrefix(b.Body, "!RRR!"):
			b.RequestReadReceipt = true
			b.Body = b.Body[5:]
		case strings.HasPrefix(b.Body, "!RDR!"):
			b.RequestDeliveryReceipt = true
			b.Body = b.Body[5:]
		case strings.HasPrefix(b.Body, "!URG!"):
			b.OutpostUrgent = true
			b.Body = b.Body[5:]
		default:
			return
		}
	}
}

// returnAddressFromHeaders extracts the return address of the message from its
// headers.  It looks for a Return-Path header; failing that, a Reply-To header;
// failing that, a Sender header; failing that, a From header.  It returns an
// empty string if no suitable address could be found.
func returnAddressFromHeaders(headers mail.Header) string {
	var header []string
	var ok bool

	if header, ok = headers["Return-Path"]; !ok {
		if header, ok = headers["Reply-To"]; !ok {
			if header, ok = headers["Sender"]; !ok {
				header = headers["From"]
			}
		}
	}
	if len(header) != 0 {
		if al, err := mail.ParseAddressList(header[0]); err == nil && len(al) != 0 {
			return al[0].Address
		}
	}
	return ""
}
