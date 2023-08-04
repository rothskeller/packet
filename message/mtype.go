// Package message contains the interfaces and registry for packet message
// types.  The definitions in this package can be used to register package
// message types and to itemize the registered types.
package message

// Type gives the details of a registered message type.
type Type struct {
	// Tag is the tag string that identifies the message type.
	Tag string
	// Name is the English name of the message type, in prose case.
	Name string
	// Article is the indefinite article to use before the Name when needed;
	// it is always either "a" or "an".
	Article string
	// an empty string for message types that don't have one.
	// Create is a function for creating a new message of the type.  If
	// Create is nil, end users are not allowed to create new messages of
	// the type.  Otherwise, create must be a function with the signature
	// func() «Message», where «Message» is any type that implements the
	// Message interface.  The new message will have default values in all
	// fields.
	Create any
	// Decode is a function for decoding messages of the type.  If the input
	// message conforms to the type, Decode will return the decoded message;
	// otherwise, it will return nil.  The function must have the signature
	// func(subject string, body string) «Message», where «Message» is any
	// type that implements the Message interface.
	Decode any
}

// Message is the interface that all message types implement.
type Message interface {
	// Type returns the message type definition.
	Type() *Type
	// EncodeSubject encodes the message subject line.
	EncodeSubject() string
	// EncodeBody encodes the message body, suitable for transmission or
	// storage.
	EncodeBody() string
}

// HumanMessage is the interface satisfied by all human-originated messages, or
// in other words, all messages that are assigned an origin message ID.
type HumanMessage interface {
	Message
	// GetHandling returns the handling order of the message.
	GetHandling() string
	// SetHandling sets the handling order of the message.
	SetHandling(string)
	// GetOriginID returns the origin message ID of the message.
	GetOriginID() string
	// SetOriginID sets the origin message ID of the message.
	SetOriginID(string)
	// GetSubject returns the subject of the message (not including
	// message ID, handling, etc.).
	GetSubject() string
}

// KnownForm is the interface satisfied by all messages that are recognized
// form types and have the standard fields used in all known forms.
type KnownForm interface {
	HumanMessage
	// KeyFields returns a structure containing the values of certain key
	// fields of the message that are needed for message analysis.
	KeyFields() *KeyFields
}

// KeyFields is the structure returned by the KeyFields method of known form
// types.
type KeyFields struct {
	PIFOVersion   string
	FormVersion   string
	ToICSPosition string
	ToLocation    string
	OpCall        string
}

// IValidate is the interface for the Validate method supported by some (but not
// all) message types.
type IValidate interface {
	HumanMessage
	// Validate checks the contents of the message for compliance with rules
	// enforced by standard Santa Clara County packet software (Outpost and
	// PackItForms).  It returns a list of strings describing problems that
	// those programs would flag or block.
	Validate() (problems []string)
}

// ICompare is the interface for the Compare method supported by some (but not
// all) message types.
type ICompare interface {
	HumanMessage
	// Compare compares two messages.  It returns a score indicating how
	// closely they match, and the detailed comparisons of each field in the
	// message.  The comparison is not symmetric:  the receiver of the call
	// is the "expected" message and the argument is the "actual" message.
	Compare(actual Message) (score, outOf int, fields []*CompareField)
}

// A CompareField structure represents a single field in the comparison of two
// messages as performed by the Compare method.
type CompareField struct {
	// Label is the field label.
	Label string
	// Score is the comparison score for this field.  0 <= Score <= OutOf.
	Score int
	// OutOf is the maximum possible score for this field, i.e., the score
	// for this field if its contents match exactly.
	OutOf int
	// Expected is the value of this field in the expected message (i.e.,
	// the receiver of the Compare method), formatted for human viewing.
	Expected string
	// ExpectedMask is a string describing which characters of Expected are
	// different from those in Actual.  Space characters in the mask
	// correspond to characters in Expected that are properly matched by
	// Actual.  "~" characters in the mask correspond to characters in
	// Expected that have minor differences in Actual.  All other characters
	// in the mask correspond to significant differences.  If ExpectedMask
	// is shorter than Expected, the last character of ExpectedMask is
	// implicitly repeated.
	ExpectedMask string
	// Actual is the value of this field in the actual message (i.e., the
	// argument of the Compare method), formatted for human viewing.
	Actual string
	// ActualMask is a string describing which characters of Actual are
	// different from those in Expected.  Space characters in the mask
	// correspond to characters in Actual that properly match Expected.  "~"
	// characters in the mask correspond to characters in Actual that have
	// minor differences with Expected.  All other characters in the mask
	// correspond to significant differences.  If ActualMask is shorter than
	// Actual, the last character of ActualMask is implicitly repeated.
	ActualMask string
}

// IRenderPDF is the interface for the RenderPDF method supported by some (but
// not all) message types.
type IRenderPDF interface {
	HumanMessage
	// RenderPDF renders the message as a PDF file with the specified
	// filename, overwriting any existing file with that name.  Note that
	// the program needs to be built with the packet-pdf build tag in order
	// to include these methods.
	RenderPDF(filename string) error
}

// IRenderTable is the interface for the RenderTable method supported by some
// (but not all) message types.
type IRenderTable interface {
	Message
	// RenderTable renders the message as a set of field label / field value
	// pairs, intended for read-only display to a human.
	RenderTable() []LabelValue
}

// A LabelValue describes a single field label / field value pair returned from
// a message's RenderTable method.
type LabelValue struct{ Label, Value string }

// IEdit is the interface for the Edit method supported by some (but not
// all) message types.
type IEdit interface {
	HumanMessage
	// EditFields returns the set of editable fields of the message.
	// Callers may change the Value in each field, but must otherwise treat
	// the set as read-only.  Changing the Value in a field does not affect
	// the underlying message until ApplyEdits is called.
	EditFields() []*EditField
	// ApplyEdits applies the revised Values in the EditFields to the
	// message.
	ApplyEdits()
}

// An EditField is an editable field of a message.
type EditField struct {
	// Label is the label for the field.  It should be brief, 40 characters
	// maximum.
	Label string
	// Value is the value of the field, formatted for human presentation.
	Value string
	// Width is the width of the input control for the field.  It should
	// match the amount of space available for the field in the fillable
	// PDF, or the maximum length of any allowed value for the field.
	Width int
	// Multiline is a flag indicating that the value can contain newlines.
	Multiline bool
	// LocalMessageID is a flag indicating that this field contains the
	// local message ID of the message.  Exactly one field must have this
	// flag set.
	LocalMessageID bool
	// Problem is a string describing the issue with the current value of
	// the field, if it is not valid.  If it is valid, Problem is an empty
	// string.
	Problem string
	// Choices is the set of recommended or restricted values for the field,
	// if any.  If there are none, it is nil.
	Choices []string
	// Help is the help text to display for the field when requested.  Note
	// that this must include any text necessary to explain issues with the
	// current value of the field, if any.  The returned text will be
	// word-wrapped for display.
	Help string
	// Hint is a hint to be displayed next to the input control for the
	// field, generally to give format suggestions, e.g., "MM/DD/YYYY" for a
	// date field.
	Hint string
}

// IUpdate is the interface for the UpdateReceived and UpdateSent methods
// supported by some (but not all) message types.
type IUpdate interface {
	HumanMessage
	// UpdateReceived updates the message contents to reflect the fact that
	// it has just been received, including saving the provided destination
	// (local) message ID, operator call sign and name in the requisite
	// fields if any.
	UpdateReceived(dmi, opcall, opname string)
	// UpdateSent updates the message contents to reflect the fact that it
	// is about to be sent, including saving the provided operator call sign
	// any name in the requisite fields if any.
	UpdateSent(opcall, opname string)
	// UpdateDelivered updates the message contents to reflect the fact that
	// it has been delivered, by setting the provided destination (remote)
	// message ID in the requisite field if any.
	UpdateDelivered(dmi string)
}

// ISetSubject is the interface for the SetSubject method.  It is implemented by
// message types where the subject is a free-form text field.  For those message
// types, GetSubject should return the value set by SetSubject.
type ISetSubject interface {
	HumanMessage
	// SetSubject sets the value of the subject field of the message.
	SetSubject(string)
}

// ISetBody is the interface for the SetBody method.  It is implemented by any
// message type that has at least one free-form text field.  (Note this is not
// parallel to ISetBody; that's why they have separate intefaces.)
type ISetBody interface {
	HumanMessage
	// SetBody sets the value of the most prominent free-form text field of
	// the message.
	SetBody(string)
}

// IGetBody is the interface for the GetBody method.  It is implemented by those
// message types that have a single body text field (basically plain text and
// ICS-213).  For those message types, it returns the value of the same field
// that SetBody sets.
type IGetBody interface {
	ISetBody
	// GetBody retrieves the value of the body field of the message.
	GetBody() string
}
