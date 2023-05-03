// Package xscmsg handles recognition, parsing, validation, and encoding of
// Santa Clara County (XSC) standard messages.  It supports plain text messages
// and "unknown form" PackItForms messages directly, and allows registration of
// other message types.
package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscsubj"
)

// A Message is a pktmsg.Message with a defined message type.
type Message interface {

	// Non-overridden methods of pktmsg.Message:

	// BBSRxDate returns the BBSRxDate field of the message, if it has one.
	// It holds the time the message was received by the BBS.  It is present
	// only on instantly-received messages; it is not persisted in local
	// storage.
	BBSRxDate() pktmsg.BBSRxDateField
	// FormHTML returns the FormHTML field of the message, if it has one.
	// It holds the PackItForms HTML file for the form.  It is present only
	// on form messages.
	FormHTML() pktmsg.FormHTMLField
	// FormVersion returns the FormVersion field of the message, if it has
	// one.  It holds the form version number.  It is present only on form
	// messages.
	FormVersion() pktmsg.FormVersionField
	// FromAddr returns the FromAddr field of the message.  It holds the
	// origin address (From: header).  It may contain a name as well as an
	// address.
	FromAddr() pktmsg.FromAddrField
	// NotPlainText returns the NotPlainText field of the message, if it has
	// one.  Its presence indicates that the instantly-received message was
	// not in plain text.  It is not persisted in local storage.
	NotPlainText() pktmsg.NotPlainTextField
	// OutpostFlags returns the OutpostFlags field of the message, if it has
	// one.  It holds the Outpost message flags.
	OutpostFlags() pktmsg.OutpostFlagsField
	// PIFOVersion returns the PIFOVersion field of the message, if it has
	// one.  It holds the PackItForms encoding version number.  It is
	// present only on form messages.
	PIFOVersion() pktmsg.PIFOVersionField
	// ReturnAddr returns the ReturnAddr field of the message, if it has
	// one.  It holds the return address of the message.  It is present only
	// on instantly-received messages; it is not persisted in local storage.
	ReturnAddr() pktmsg.ReturnAddrField
	// RxArea returns the RxArea field of the message, if it has one.  It
	// holds the BBS bulletin area from which the message was retrieved.  It
	// is present only on received bulletin messages.
	RxArea() pktmsg.RxAreaField
	// RxBBS returns the RxBBS field of the message, if it has one.  It
	// holds the name of the BBS from which the message was retrieved.  It
	// is present only on received messages.
	RxBBS() pktmsg.RxBBSField
	// RxDate returns the RxDate field of the message, if it has one.  It
	// holds the time the message was received locally.  It is present only
	// on received messages.
	RxDate() pktmsg.RxDateField
	// SentDate returns the SentDate field of the message.  It holds the
	// time the message was sent (Date: header).  It is empty for outgoing
	// messages that have not yet been sent.
	SentDate() pktmsg.SentDateField
	// ToAddrs returns the ToAddrs field of the message.  It holds the list
	// of destination addresses (To: header).
	ToAddrs() pktmsg.ToAddrsField

	// TaggedField returns the Field with the specified tag string, or nil
	// if there is none.
	TaggedField(string) pktmsg.Field
	// TaggedFields calls the supplied function for each tagged field of the
	// message, in order.
	TaggedFields(func(string, pktmsg.Field))

	// Save returns the message, formatted for saving to local storage.
	// Note that this can be a lossy operation; some Fields are not
	// preserved in local storage.
	Save() string
	// Transmit returns the destination addresses, subject header, and body
	// of the message, suitable for transmission through JNOS.
	Transmit() (to []string, subject, body string)

	// Non-overridden methods of xscsubj.Message:

	// FormTag returns the FormTag field of the message.  For PackItForms
	// messages, this contains a tag identifying the form type.
	FormTag() xscsubj.FormTagField
	// SubjectHeader returns the SubjectHeader field of the message.  This
	// contains the unparsed Subject: header line of the message.
	SubjectHeader() xscsubj.SubjectHeaderField

	// New methods for xscmsg.Message:

	// Type returns the details of the message type.
	Type() MessageType
	// Iterate calls the supplied function for each field of the message, in
	// user-visible order (non-visible fields last).
	Iterate(func(pktmsg.Field))

	// Extended interfaces for fields of pktmsg.Message or xscsubj.Message:

	// Body returns the Body field of the message.  It holds the body of the
	// message.
	Body() BodyField
	// Handling returns the Handling field of the message.  This contains
	// the message handling order.
	Handling() HandlingField
	// OriginMsgID returns the OriginMessageID field of the message.
	// It contains the message ID assigned by the sender of the message.
	OriginMsgID() OriginMsgIDField
	// Severity returns the Severity field of the message.  It contains the
	// severity of the situation to which the message pertains.  This is an
	// obsolete field.
	Severity() SeverityField
	// Subject returns the Subject field of the message.  It holds the
	// message subject (Subject: header).
	Subject() SubjectField

	// New fields of xscmsg.Message:

	DestinationMsgID() pktmsg.SettableField
	MessageDate() pktmsg.SettableField
	MessageTime() pktmsg.SettableField
	OpCall() pktmsg.SettableField
	OpDate() pktmsg.SettableField
	OpName() pktmsg.SettableField
	OpTime() pktmsg.SettableField
	Reference() pktmsg.SettableField
	Retrieved() pktmsg.Field
	TacCall() pktmsg.SettableField
	TacName() pktmsg.SettableField
	ToICSPosition() pktmsg.SettableField
	ToLocation() pktmsg.SettableField
}

// MessageType describes a defined message type.
type MessageType struct {
	// Tag is a unique identifier of the message type, generally a single
	// word in upper case.  For PackItForms messages, the tag appears on the
	// subject line after the handling order.
	Tag string
	// Name is the English name of the message type, in prose case.
	Name string
	// Article is the indefinite article that should precede Name, when
	// needed.  It must be either "a" or "an".
	Article string
}

// Field capability interfaces:
type (
	// A CalculatedField is a pktmsg.Field whose value is calculated based
	// on the values of other fields.  A CalculatedField must not also be an
	// EditableField.
	CalculatedField interface {
		// Calculate calculates the value of the field.
		Calculate(m pktmsg.Message)
	}
	// A ChoicesField is a field for which there is a discrete set of
	// preferred or allowed values.  (Whether values other than these are
	// allowed depends on the behavior of the Validate method.)
	ChoicesField interface {
		// Choices returns the set of preferred or allowed values, in
		// order by preference if there is a preferred order, or
		// alphabetically otherwise.
		Choices(m pktmsg.Message) []string
	}
	// A DefaultedField is a field for which there is a default value, which
	// gets applied to the field when creating new messages.
	DefaultedField interface {
		// Default returns the default value for the field.
		Default() string
	}
	// An EditableField is a ViewableField that can also be edited by end
	// users in a message editor or similar input.
	EditableField interface {
		// Help returns help text describing the field and its allowed
		// values.  The returned text will be word-wrapped for
		// presentation.
		Help(m pktmsg.Message) string
	}
	// A HintedField is a field that can display a hint about allowed
	// values.
	HintedField interface {
		// Hint returns a hint about the of data that can appear in the
		// field (e.g., "MM/DD/YYYY" for a date field).
		Hint() string
	}
	// A SizedField is a field with a known size for display and editing.
	SizedField interface {
		// Size returns the size of the field for display and editing.
		// This corresponds to the amount of characters and lines that
		// fit in the field in a paper form.
		Size() (width, height int)
	}
	// A ValidatableField is one whose contents can be checked for validity.
	ValidatableField interface {
		// Validate checks the validity of the field value.  If it is
		// valid, it returns an empty string.  Otherwise, it returns a
		// description of the problem with the field.  The problem
		// description should identify which field has the problem; that
		// context is not added externally.  If the pifo flag is true,
		// Validate only reports values that PackItForms considers
		// invalid.
		Validate(m pktmsg.Message, pifo bool) string
	}
	// A ViewableField is a pktmsg.Field that is displayed to the end user
	// in a message viewer or similar output.
	ViewableField interface {
		// Label returns the English name of the field, in title case.
		// It should be relatively short (40 chars maximum, less is
		// better).
		Label() string
	}
)

// RegisteredTypes is a map from message type tag to message type name for all
// registered message types.  It should be treated as read-only; to register a
// new type, call RegisterType.
var RegisteredTypes = map[string]string{PlainTextTag: plainTextName}

// createFunctions is a map from message type tag to the function that creates
// new messages of that type.
var createFunctions = map[string]func() Message{PlainTextTag: createPlainText}

// recognizeFunctions is a list of functions that can recognize a message.
var recognizeFunctions = []func(pktmsg.Message) Message{
	// Note that Recognize walks this list backwards, LIFO, so that plain
	// text is always called last and unknown form second-to-last.
	recognizePlainText,
	recognizeUnknownForm,
}

// Register registers a message type.
func Register(tag, name string, createFn func() Message, recognizeFn func(pktmsg.Message) Message) {
	RegisteredTypes[tag] = name
	if createFn != nil {
		createFunctions[tag] = createFn
	}
	if recognizeFn != nil {
		recognizeFunctions = append(recognizeFunctions, recognizeFn)
	}
}

// NewMessage creates a new, outgoing message with the specified type.  It
// returns nil if the type tag is not registered.
func NewMessage(tag string) Message {
	if createFn := createFunctions[tag]; createFn != nil {
		return createFn()
	}
	return nil
}

// ParseMessage parses the supplied string as a stored message (which could be
// either incoming or outgoing).
func ParseMessage(raw string) (m Message, err error) {
	var pm pktmsg.Message

	if pm, err = pktmsg.ParseMessage(raw); err != nil {
		return nil, err
	}
	pm = xscsubj.ParseSubject(pm)
	for i := len(recognizeFunctions) - 1; i >= 0; i-- {
		if msg := recognizeFunctions[i](pm); msg != nil {
			return msg, nil
		}
	}
	panic("not reachable") // recognizePlainText accepts anything
}

// ReceiveMessage parses the supplied string as a message that was just
// retrieved from the specified JNOS BBS.  If it is a bulletin, area should be
// set to the bulletin area from which it was retrieved; otherwise, area should
// be empty.
func ReceiveMessage(raw, bbs, area string) (m Message, err error) {
	var pm pktmsg.Message

	if pm, err = pktmsg.ReceiveMessage(raw, bbs, area); err != nil {
		return nil, err
	}
	pm = xscsubj.ParseSubject(pm)
	for i := len(recognizeFunctions) - 1; i >= 0; i-- {
		if msg := recognizeFunctions[i](pm); msg != nil {
			return msg, nil
		}
	}
	panic("not reachable") // recognizePlainText accepts anything
}
