// Package xscmsg handles recognition, parsing, validation, and encoding of all
// of the Santa Clara County (XSC) standard message types.  Each message type is
// represented by its own concrete type; all of them implement the Message
// interface.
package xscmsg

// Implementation note:
//
// It would be more conceptually correct to define Message and Field as
// interfaces.  xscform could have its own implementation of them, and fields of
// those types that are specific to forms (such as Message.Subject and
// Field.Annotation) could be private to xscform.  However, that results in
// considerably greater code complexity.  Instead, the code has one common
// implementation of Message and Field as structures, with function pointers in
// the MessageType and FieldDef structures that provide hooks for type-specific
// functionality.  It's vastly easier to work with field definitions that are
// literal structure values rather than composing chains of nested objects.

import (
	"strconv"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

// Message is the interface implemented by all XSC messages.
type Message struct {
	// Type is the message type definition.
	Type *MessageType
	// RawMessage is the underlying raw pktmsg.Message, if any.
	RawMessage *pktmsg.Message
	// RawForm is the underlying raw pktmsg.Form, if any.
	RawForm *pktmsg.Form
	// Fields returns the list of fields in the message.
	Fields []*Field
}

// Field returns the field with the specified name, or nil if there is
// no such field.  It matches against both FieldDef.Tag and FieldDef.Canonical.
func (m *Message) Field(name string) *Field {
	for _, f := range m.Fields {
		if f.Def.Tag == name || f.Def.Canonical == name {
			return f
		}
	}
	return nil
}

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (m *Message) Validate(strict bool) (problems []string) {
	for _, f := range m.Fields {
		if prob := f.Validate(m, strict); prob != "" {
			problems = append(problems, prob)
		}
	}
	return problems
}

// Subject returns the encoded subject of the message.
func (m *Message) Subject() string {
	if m.Type.SubjectFunc != nil {
		return m.Type.SubjectFunc(m)
	}
	if f := m.Field(FSubject); f != nil {
		return f.Value
	}
	return ""
}

// Body returns the encoded body of the message.  If human is
// true, it is encoded for human reading or editing; if false, it is
// encoded for transmission.
func (m *Message) Body(human bool) string {
	if m.Type.BodyFunc != nil {
		return m.Type.BodyFunc(m, human)
	}
	if f := m.Field(FBody); f != nil {
		return f.Value
	}
	return ""
}

// MessageType defines a message type.
type MessageType struct {
	// Tag is the string used to identify the message type.  For form
	// messages, this is the form tag that appears encoded in the subject
	// line of the message.
	Tag string
	// Name is the English name of the message type.  It is a noun phrase in
	// prose case, such as "foo message" or "bar form".
	Name string
	// Article is the indefinite article ("a" or "an") to be used preceding
	// the Name, in a sentence that needs one.
	Article string
	// HTML is the HTML file listed in the PackItForms header for the
	// message type, if any.
	HTML string
	// Version is the version number of the message type, if any.
	Version string
	// SubjectFunc is a function that returns the encoded subject of the
	// message.  If it is nil, the subject is taken from the "SUBJECT"
	// message field.
	SubjectFunc func(msg *Message) string
	// BodyFunc is a function that returns the encoded body of the message.
	// If human is true, it is encoded for human reading or editing; if
	// false, it is encoded for transmission.  If it is nil, the body is
	// taken from the "BODY" message field.
	BodyFunc func(msg *Message, human bool) string
}

// Field is a single field within an XSC message.
type Field struct {
	// Def is the definition of the field.
	Def *FieldDef
	// Value is the value of the field.
	Value string
}

// Validate verifies the correctness of the field value, returning an error
// message if it is invalid or an empty string if it is valid.  If strict is
// true, correctable problems are reported; otherwise they are corrected.  For
// fields whose values are computed based on the values of other fields,
// Validate performs the computation.
func (f *Field) Validate(msg *Message, strict bool) string {
	for _, vfn := range f.Def.Validators {
		if prob := vfn(f, msg, strict); prob != "" {
			return prob
		}
	}
	return ""
}

// Default returns the default value of the field.
func (f *Field) Default() string {
	if f.Def.DefaultFunc != nil {
		return f.Def.DefaultFunc()
	}
	return f.Def.DefaultValue
}

// FieldDef contains the definition of a Field within a Message.
type FieldDef struct {
	// Tag is the unique identifier of the field within its message.  For
	// most form fields, this is a field number followed by a period.
	Tag string
	// Canonical is a name for a well-known field that is common across all
	// forms containing the field, even if different message types have
	// different tags for it.  Canonical is empty for less common fields.
	Canonical string
	// Label is the English label of the field.  For form fields, it is the
	// label of the field as it appears on the PDF form.
	Label string
	// Annotation is a short textual annotation added to the Tag when
	// rendering a form field in human-readable mode.  It gives a string
	// name to the field when the tag is a field number.
	Annotation string
	// Comment is a comment displayed in the human-readable rendering of a
	// form when the field has no assigned value.  It is generally a textual
	// reminder of the validation requirements for the field value.
	Comment string
	// ReadOnly is a flag indicating that the field is read-only.  It may
	// have a value when a message is received and decoded, or it may have
	// one computed by its Validate method based on other message fields,
	// but it should not be presented to the user for editing.
	ReadOnly bool
	// DefaultFunc is a function to compute the default value of the field.
	DefaultFunc func() string
	// DefaultValue is the default value of the field, if DefaultFunc is not
	// set.
	DefaultValue string
	// Validators is the list of functions called to validate the field
	// value.  See the comment on Field.Validate() for details.
	Validators []Validator
	// Choices is the set of allowed values for a field with a restricted
	// set.
	Choices []string
}

// Validator is a function called to validate a field value.  See the
// comment on Field.Validate() for details.
type Validator func(f *Field, msg *Message, strict bool) (problem string)

// Values for FieldID.Canonical.  The values are in all caps so that they do not
// overlap with PackItForms field tags.
const (
	FBody             = "BODY"
	FDestinationMsgNo = "DESTINATION_MESSAGE_NUMBER"
	FHandling         = "HANDLING"
	FMessageDate      = "MESSAGE_DATE"
	FMessageTime      = "MESSAGE_TIME"
	FOpCall           = "OPERATOR_CALL_SIGN"
	FOpDate           = "OPERATOR_DATE"
	FOpName           = "OPERATOR_NAME"
	FOpTime           = "OPERATOR_TIME"
	FOriginMsgNo      = "ORIGIN_MESSAGE_NUMBER"
	FSubject          = "SUBJECT"
	FToICSPosition    = "TO_ICS_POSITION"
	FToLocation       = "TO_LOCATION"
)

var createFuncs = map[string]func() *Message{}
var recognizeFuncs []func(*pktmsg.Message, *pktmsg.Form) *Message

// RegisterCreate registers a function to create a message with the specified
// type tag.  (It is normally called by an init function in the package that
// implements the message type.  To make a message type usable, simply import
// its implementing package.)
func RegisterCreate(tag string, fn func() *Message) {
	createFuncs[tag] = fn
}

// RegisterType registers a message type to be recognizable by the Recognize
// function.  (It is normally called by an init function in the package that
// implements the message type.  To make a message type recognizable, simply
// import its implementing package.)
func RegisterType(recognize func(*pktmsg.Message, *pktmsg.Form) *Message) {
	recognizeFuncs = append(recognizeFuncs, recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized, Create returns nil.
//
// For this to work, the message type you want to create must have been
// registered.  The standard message types can be registered by importing the
// appropriate message-type-specific subpackages of xscmsg.  Alternatively, all
// standard message types can be registered by importing xscmsg/all.
func Create(tag string) *Message {
	if fn := createFuncs[tag]; fn != nil {
		return fn()
	}
	return nil
}

// Recognize determines which registered XSC message type to use for the
// supplied pktmsg.Message, and returns the corresponding xscmsg.Message.  The
// strict flag indicates whether any form embedded in the message should be
// parsed strictly; see pktmsg.ParseForm for details.  Recognize always returns
// an xscmsg.Message; if the supplied message isn't recognized as any more
// specific type, it is returned as an "unknown form" message or a "plain text"
// message.
//
// For a message to be recognized as a specific type, that type must have been
// registered.  The standard message types can be registered by importing the
// appropriate message-type-specific subpackages of xscmsg.  Alternatively, all
// standard message types can be registered by importing xscmsg/all.
func Recognize(msg *pktmsg.Message, strict bool) *Message {
	var form *pktmsg.Form
	if pktmsg.IsForm(msg.Body) {
		form = pktmsg.ParseForm(msg.Body, strict)
	}
	for _, recognizeFunc := range recognizeFuncs {
		if xsc := recognizeFunc(msg, form); xsc != nil {
			return xsc
		}
	}
	if form != nil {
		return makeUnknownFormMessage(msg, form)
	}
	return makePlainTextMessage(msg)
}

// OlderVersion compares two version numbers, and returns true if the first one
// is older than the second one.  Each version number is a dot-separated
// sequence of parts, each of which is compared independently with the
// corresponding part in the other version number.  The parts are compared
// numerically if they parse as integers, and as strings otherwise.  However, a
// part that starts with a digit is always "newer" than a part that doesn't.
// (This results in "undefined" being treated as infinitely old, which is what
// we want.)
func OlderVersion(a, b string) bool {
	aparts := strings.Split(a, ".")
	bparts := strings.Split(b, ".")
	for len(aparts) != 0 && len(bparts) != 0 {
		anum, aerr := strconv.Atoi(aparts[0])
		bnum, berr := strconv.Atoi(bparts[0])
		if aerr == nil && berr == nil {
			if anum < bnum {
				return true
			}
			if anum > bnum {
				return false
			}
		} else if startsWithDigit(aparts[0]) && !startsWithDigit(bparts[0]) {
			return false
		} else if !startsWithDigit(aparts[0]) && startsWithDigit(bparts[0]) {
			return true
		} else {
			if aparts[0] < bparts[0] {
				return true
			}
			if aparts[0] > bparts[0] {
				return false
			}
		}
		aparts = aparts[1:]
		bparts = bparts[1:]
	}
	if len(bparts) != 0 {
		return true
	}
	return false
}
func startsWithDigit(s string) bool {
	return s != "" && s[0] >= '0' && s[0] <= '9'
}
