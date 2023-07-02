// Package message contains the interfaces and registry for packet message
// types.  The definitions in this package can be used to register package
// message types and to itemize the registered types.
//
// Subpackages of this package provide definitions of all of the public,
// standard message types used in Santa Clara County.  These definitions must be
// registered at runtime with this package before they can be used.
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
	// Create is a function for creating a new message of the type.  If
	// Create is nil, end users are not allowed to create new messages of
	// the type.  Otherwise, create must be a function with one of the
	// following signatures:
	//   - func() «Message»
	//   - func(opcall, opname string) «Message»
	// where «Message» is any type that implements the Message interface.
	// The new message will have default values in all fields, and will have
	// the operator information filled in if the type supports it.
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
}

// IEncode is the interface for the Encode method supported by some (but not
// all) message types.
type IEncode interface {
	// EncodeSubject encodes the message subject line.
	EncodeSubject() string
	// EncodeBody encodes the message body, suitable for transmission or
	// storage.
	EncodeBody() string
}

// IKeyFields is the interface for the KeyFields method supported by some (but
// not all) message types.
type IKeyFields interface {
	// KeyFields returns a structure containing the values of certain key
	// fields of the message that are needed for message analysis.
	KeyFields() *KeyFields
}

// KeyFields is the structure returned by the KeyFields method of some message
// types.
type KeyFields struct {
	PIFOVersion   string
	FormVersion   string
	OriginMsgID   string
	Handling      string
	ToICSPosition string
	ToLocation    string
	Subject       string
	SubjectLabel  string
	OpCall        string
}

// IValidate is the interface for the Validate method supported by some (but not
// all) message types.
type IValidate interface {
	// Validate checks the contents of the message for compliance with rules
	// enforced by standard Santa Clara County packet software (Outpost and
	// PackItForms).  It returns a list of strings describing problems that
	// those programs would flag or block.
	Validate() (problems []string)
}

// ICompare is the interface for the Compare method supported by some (but not
// all) message types.
type ICompare interface {
	// Compare compares two messages, and returns a list of strings
	// describing differences between them.  Only "significant" differences
	// are reported.
	Compare(other any) (differences []string)
}

// IRenderPDF is the interface for the RenderPDF method supported by some (but
// not all) message types.
type IRenderPDF interface {
	// RenderPDF renders the message as a PDF file with the specified
	// filename, overwriting any existing file with that name.  Note that
	// the program needs to be built with the packet-pdf build tag in order
	// to include these methods.
	RenderPDF(filename string) error
}

// IRenderTable is the interface for the RenderTable method supported by some
// (but not all) message types.
type IRenderTable interface {
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
	// Edit TBD
	Edit()
}

// IUpdate is the interface for the UpdateReceived and UpdateSent methods
// supported by some (but not all) message types.
type IUpdate interface {
	// UpdateReceived updates the message contents to reflect the fact that
	// it has just been received, including saving the provided operator
	// call sign and name in the requisite fields if any.
	UpdateReceived(opcall, opname string)
	// UpdateSent updates the message contents to reflect the fact that it
	// is about to be sent, including saving the provided operator call sign
	// any name in the requisite fields if any.
	UpdateSent(opcall, opname string)
}
