package pktmsg

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// Encode encodes the message for transmission or storage.
func (msg *Message) Encode(human bool) string {
	var sb strings.Builder

	// If the message has envelope data, emit an RFC-4155 "From " envelope
	// line.
	if msg.EnvelopeAddress != "" || !msg.EnvelopeDate.IsZero() {
		sb.WriteString("From ")
		sb.WriteString(msg.EnvelopeAddress)
		if !msg.EnvelopeDate.IsZero() {
			sb.WriteByte(' ')
			sb.WriteString(msg.EnvelopeDate.Format(time.ANSIC))
		}
		sb.WriteByte('\n')
	}
	// Emit the message headers.
	for key, vals := range msg.Header {
		fmt.Fprintf(&sb, "%s: %s\n", key, strings.Join(vals, ", "))
	}
	sb.WriteByte('\n')
	// Emit the body.
	sb.WriteString(msg.EncodeBody(human))
	return sb.String()
}

// EncodeBody encodes the message body for transmission.
func (msg *Message) EncodeBody(human bool) (body string) {
	// Add the Outpost flags to the body if needed.
	if msg.Flags != 0 {
		var sb strings.Builder
		if msg.Flags&OutpostUrgent != 0 {
			sb.WriteString("!URG!")
		}
		if msg.Flags&RequestDeliveryReceipt != 0 {
			sb.WriteString("!RDR!")
		}
		if msg.Flags&RequestReadReceipt != 0 {
			sb.WriteString("!RRR!")
		}
		sb.WriteString(msg.Body)
		body = sb.String()
	} else {
		body = msg.Body
	}
	// Encode in Base64 if needed.
	if !human && strings.IndexFunc(body, nonASCII) >= 0 {
		body = "!B64!" + base64.StdEncoding.EncodeToString([]byte(body)) + "\n"
	}
	return body
}
func nonASCII(r rune) bool {
	return r > 126 || (r < 32 && r != '\t' && r != '\n')
}
