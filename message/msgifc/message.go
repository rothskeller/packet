package msgifc

// Message is the interface satisfied by all messages.
type Message interface {
	CacheTracker
	MType
	// Type returns the message type for the message.
	Type() MType
	// SetType sets the message type for the message.
	SetType(MType)
	// RFC5322 returns the message encoded in RFC-5322 format for storage
	// or email transmission.
	RFC5322() string
	// To returns the list of recipients for the message.  Use
	// ParseAddressList to decode it (but note that To: lines in received
	// messages might not be syntactically correct).
	To() string
	// Subject is the subject of the message.
	Subject() Subject
	// Bulletin returns whether the message is a BBS bulletin (as opposed
	// to a private message).
	Bulletin() bool
	// Draft returns whether the message is a draft.
	Draft() bool
	// Received returns whether the message has been received by the local
	// system.  It returns false for a message that has been sent or is
	// being prepared to be sent.
	Received() bool
	// Payload is the payload of the message.
	Payload() Payload
	// Body is the body of the message.
	Body() Body
}
