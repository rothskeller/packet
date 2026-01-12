package msgifc

import (
	"iter"
)

// Message is the interface satisfied by all messages.
type Message interface {
	CacheTracker
	MType
	// Type returns the message type for the message.
	Type() MType
	// SetType sets the message type for the message.  (This should be
	// called only during message.Recognize.)
	SetType(MType)
	// RFC5322 returns the message encoded in RFC-5322 format for storage
	// or email transmission.
	RFC5322() string
	// To returns the list of recipients for the message.  Use
	// address.ParseList to decode it (but note that To: lines in received
	// messages might not be syntactically correct).
	To() string
	// Subject is the subject of the message.
	Subject() Subject
	// SetSubject sets the subject object of the message.  (This should be
	// called only during message.Recognize, to change the subject type.)
	SetSubject(Subject)
	// Bulletin returns whether the message is a BBS bulletin (as opposed
	// to a private message).
	Bulletin() bool
	// Payload is the payload of the message.
	Payload() Payload
	// Body is the body of the message.
	Body() Body
	// Fields returns an iterator the field definitions for the message.
	Fields() iter.Seq[Field]
}
