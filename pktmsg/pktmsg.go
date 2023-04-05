// Package pktmsg handles encoding and decoding packet messages.  It understands
// RFC-4155 and RFC-5322 email encoding, PackItForms form encoding, and
// Outpost-specific feature encodings.
package pktmsg

// A Field represents a field of a message.
type Field interface {
	// Value returns the value of the Field as a string.
	Value() string
}

// A SettableField represents a field whose value can be changed.
type SettableField interface {
	Field
	// SetValue sets the string value of the Field.
	SetValue(string)
}

// A Message is an iterable, indexed collection of Fields.
type Message interface {
	// Body returns the Body field of the message.  It holds the body of the
	// message.
	Body() BodyField
	// BBSRxDate returns the BBSRxDate field of the message, if it has one.
	// It holds the time the message was received by the BBS.  It is present
	// only on instantly-received messages; it is not persisted in local
	// storage.
	BBSRxDate() BBSRxDateField
	// FormHTML returns the FormHTML field of the message, if it has one.
	// It holds the PackItForms HTML file for the form.  It is present only
	// on form messages.
	FormHTML() FormHTMLField
	// FormVersion returns the FormVersion field of the message, if it has
	// one.  It holds the form version number.  It is present only on form
	// messages.
	FormVersion() FormVersionField
	// FromAddr returns the FromAddr field of the message.  It holds the
	// origin address (From: header).  It may contain a name as well as an
	// address.
	FromAddr() FromAddrField
	// NotPlainText returns the NotPlainText field of the message, if it has
	// one.  Its presence indicates that the instantly-received message was
	// not in plain text.  It is not persisted in local storage.
	NotPlainText() NotPlainTextField
	// OutpostFlags returns the OutpostFlags field of the message, if it has
	// one.  It holds the Outpost message flags.
	OutpostFlags() OutpostFlagsField
	// PIFOVersion returns the PIFOVersion field of the message, if it has
	// one.  It holds the PackItForms encoding version number.  It is
	// present only on form messages.
	PIFOVersion() PIFOVersionField
	// ReturnAddr returns the ReturnAddr field of the message, if it has
	// one.  It holds the return address of the message.  It is present only
	// on instantly-received messages; it is not persisted in local storage.
	ReturnAddr() ReturnAddrField
	// RxArea returns the RxArea field of the message, if it has one.  It
	// holds the BBS bulletin area from which the message was retrieved.  It
	// is present only on received bulletin messages.
	RxArea() RxAreaField
	// RxBBS returns the RxBBS field of the message, if it has one.  It
	// holds the name of the BBS from which the message was retrieved.  It
	// is present only on received messages.
	RxBBS() RxBBSField
	// RxDate returns the RxDate field of the message, if it has one.  It
	// holds the time the message was received locally.  It is present only
	// on received messages.
	RxDate() RxDateField
	// SentDate returns the SentDate field of the message.  It holds the
	// time the message was sent (Date: header).  It is empty for outgoing
	// messages that have not yet been sent.
	SentDate() SentDateField
	// Subject returns the Subject field of the message.  It holds the
	// message subject (Subject: header).
	Subject() SubjectField
	// ToAddrs returns the ToAddrs field of the message.  It holds the list
	// of destination addresses (To: header).
	ToAddrs() ToAddrsField

	// TaggedField returns the Field with the specified tag string, or nil
	// if there is none.
	TaggedField(string) Field
	// TaggedFields calls the supplied function for each tagged field of the
	// message, in order.
	TaggedFields(func(string, Field))

	// Save returns the message, formatted for saving to local storage.
	// Note that this can be a lossy operation; some Fields are not
	// preserved in local storage.
	Save() string
	// Transmit returns the destination addresses, subject header, and body
	// of the message, suitable for transmission through JNOS.
	Transmit() (to []string, subject, body string)
}

// NewMessage creates a new, empty outgoing message.
func NewMessage() Message {
	return newOutpostMessage()
}
