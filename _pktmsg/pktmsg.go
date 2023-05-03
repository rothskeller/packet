// Package pktmsg handles encoding and decoding packet messages.  It understands
// RFC-4155 and RFC-5322 email encoding and Outpost-specific feature encodings.
// There is nothing specific to Santa Clara County in package pktmsg.
package pktmsg

import (
	"net/textproto"
	"time"
)

// A Message represents a packet message.
type Message struct {
	// EnvelopeAddress is the return address on the RFC-4155 "From "
	// envelope line, if any.
	EnvelopeAddress string
	// EnvelopeDate is the message date from the RFC-4155 "From " envelope
	// line, if any.
	EnvelopeDate time.Time
	// Header contains the RFC-5322 header of the message.
	Header textproto.MIMEHeader
	// Body contains the plain text body of the message.
	Body string
	// Flags contains flags associated with the message.
	Flags MessageFlag
}

// MessageFlag contains a flag, or a bitmask of flags, describing a Message.
type MessageFlag uint8

const (
	// NotPlainText indicates that the message was not entirely plain text.
	NotPlainText MessageFlag = 1 << iota
	// AutoResponse indicates that the message appears to be a bounce
	// message or other auto-responder message.
	AutoResponse
	// RequestDeliveryReceipt indicates that the recipient should send a
	// delivery receipt to the sender when the message is delivered.
	RequestDeliveryReceipt
	// RequestReadReceipt indicates that the recipient should send a read
	// receipt to the sender when the message has been read by a human.
	RequestReadReceipt
	// OutpostUrgent indicates that the Outpost messaging software should
	// display this message as an urgent message.
	OutpostUrgent
)

// New creates a new, empty message.
func New() *Message {
	return &Message{Header: make(textproto.MIMEHeader)}
}
