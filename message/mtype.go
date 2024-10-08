package message

import "github.com/rothskeller/packet/envelope"

// This file contains the definition of message.Type, the registry of message
// types, and the Create and Decode calls that walk that registry.

// Type gives the details of a registered message type.
type Type struct {
	// Tag is the tag string that identifies the message type.
	Tag string
	// HTML is the HTML filename that identifies the type of form.  It is empty
	// for messages that are not forms.
	HTML string
	// Version is the version number of the form.  It is empty for messages that
	// are not forms.
	Version string
	// FieldOrder is an ordered list of form field tags.  When provided,
	// generated forms list the fields in that order.  (Should any fields be
	// omitted from the list, they are put at the start of the form, before
	// those in the list.)  This is nil for messages that are not forms.
	FieldOrder []string
	// Name is the English name of the message type, in prose case.
	Name string
	// Article is the indefinite article to use before the Name when needed;
	// it is always either "a" or "an".
	Article string
	// PDFBase is the PDF template (i.e., blank form) onto which we will
	// render the field values to create a PDF rendering of the message.
	PDFBase []byte
	// create is the function to create a new outgoing message of this type, with
	// appropriate default values for fields.  It is nil if new messages of this
	// type are not supported.
	create func() Message
}

// RegisteredTypes is a map from message type tag to Type structures for
// all registered message types.  It should be treated as read-only; to register
// a new type, call Register.
var RegisteredTypes = make(map[string][]*Type)

// decodeFunctions is the ordered list of functions to call to try decoding a
// message.
var decodeFunctions []func(env *envelope.Envelope, body string, form *PIFOForm, pass int) Message

// Register registers a message type.  The order of registration is significant
// for decoding messages:  catch-all decoders (e.g. UnknownForm and PlainText)
// must be registered last.  The decode function must examine the envelope,
// body, and form (which may be nil); if they are correct for the message type,
// it must return the decoded message; otherwise, it must return nil.
func Register(mtype *Type, decode func(env *envelope.Envelope, body string, form *PIFOForm, pass int) Message, create func() Message) {
	mtype.create = create
	RegisteredTypes[mtype.Tag] = append(RegisteredTypes[mtype.Tag], mtype)
	decodeFunctions = append(decodeFunctions, decode)
}

// Create creates a new, outgoing message with the specified type tag and
// version.  If version is empty, it creates the message with the first-
// registered version of the type (that supports message creation).  Create
// returns nil if the type tag is not registered or does not support message
// creation.  All message fields have default values for an outgoing message.
func Create(tag, version string) Message {
	for _, mtype := range RegisteredTypes[tag] {
		if mtype.create != nil && (version == "" || version == mtype.Version) {
			return mtype.create()
		}
	}
	return nil
}

// Decode decodes the supplied message and returns the typed, decoded message.
// It returns nil if no registered type can decode the message.  (Between them,
// the UnknownForm, Bulletin, and PlainText message types can decode any
// message, so if they are registered, a nil return is not possible.)
func Decode(env *envelope.Envelope, body string) (msg Message) {
	// Decode the PIFO form in the message if any.
	var form = DecodePIFO(body)
	for pass := 1; pass <= 2; pass++ {
		for _, fn := range decodeFunctions {
			if msg = fn(env, body, form, pass); msg != nil {
				return msg
			}
		}
	}
	return nil
}
