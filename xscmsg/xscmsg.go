// Package xscmsg handles recognition, parsing, validation, and encoding of all
// of the Santa Clara County (XSC) standard message types.  Each message type is
// represented by its own concrete type; all of them implement the Message
// interface.
package xscmsg

import (
	"strconv"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

// Message is the interface implemented by all XSC messages.
type Message interface {
	// Type returns the message type definition.
	Type() *MessageType
	// Fields returns the list of fields in the message.
	Fields() []Field
	// Field returns the field with the specified name, or nil if there is
	// no such field.
	Field(name string) Field
	// Validate ensures that the contents of the message are correct.  It
	// returns a list of problems, which is empty if the message is fine.
	// If strict is true, the message must be exactly correct; otherwise,
	// some trivial issues are corrected and not reported.
	Validate(strict bool) (problems []string)
	// EncodeSubject returns the encoded subject of the message.
	EncodeSubject() string
	// EncodeBody returns the encoded body of the message.  If human is
	// true, it is encoded for human reading or editing; if false, it is
	// encoded for transmission.
	EncodeBody(human bool) string
}

// MessageType contains the static data that identifies a particular type of
// message.
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
}

// Field is the interface implemented by all fields in XSC messages.
type Field interface {
	// ID returns the identification of the field.
	ID() *FieldID
	// Validate verifies the correctness of the field value, returning an
	// error message if it is invalid or an empty string if it is valid.
	// If strict is true, correctable problems are reported as errors;
	// otherwise they are corrected.
	Validate(msg Message, strict bool) string
	// Get returns the value of the field.
	Get() string
	// Set sets the value of the field.
	Set(string)
	// Default returns the default value of the field.
	Default() string
}

// FieldID contains the static data that identifies a particular field within
// a Message.
type FieldID struct {
	// Tag is the string that identifies the field.  For most form fields,
	// this is a field number followed by a period.
	Tag string
	// Annotation is a short textual annotation added to the Tag when
	// rendering a form field in human-readable mode.  It gives a string
	// name to the field when the tag is a field number.
	Annotation string
	// Label is the label of the field as it appears on the PDF form.
	Label string
	// Comment is a comment displayed in the human-readable rendering of a
	// form when the field has no assigned value.
	Comment string
	// Canonical is a name for the field that is common across all forms
	// containing the field.
	Canonical string
	// ReadOnly is a flag indicating that the field is read-only.  It may
	// have a value when a message is received and decoded, or it may have
	// one computed by its Validate method based on other message fields,
	// but it should not be presented to the user for editing.
	ReadOnly bool
}

// Values for FieldID.Canonical:
const (
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
	// TODO more
)

var createFuncs = map[string]func() Message{}
var recognizeFuncs []func(*pktmsg.Message, *pktmsg.Form) Message

// RegisterCreate registers a function to create a message with the specified
// type tag.  (It is normally called by an init function in the package that
// implements the message type.  To make a message type usable, simply import
// its implementing package.)
func RegisterCreate(tag string, fn func() Message) {
	createFuncs[tag] = fn
}

// RegisterType registers a message type to be accessible through the Create and
// Recognize functions.  (It is normally called by an init function in the
// package that implements the message type.  To make a message type
// recognizable, simply import its implementing package.)
func RegisterType(recognize func(*pktmsg.Message, *pktmsg.Form) Message) {
	recognizeFuncs = append(recognizeFuncs, recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized, Create returns nil.
//
// For this to work, the message type you want to create must have been
// registered.  The standard message types can be registered by importing the
// appropriate message-type-specific subpackages of xscmsg.  Alternatively, all
// standard message types can be registered by importing xscmsg/all.
func Create(tag string) Message {
	if fn := createFuncs[tag]; fn != nil {
		return fn()
	}
	return nil
}

// Recognize examines the supplied pktmsg.Message to see if it is one of the
// registered XSC message types.  If so, it returns the appropriate Message
// implementation wrapping it.  If not, it returns nil.  The strict flag
// indicates whether any form embedded in the message should be parsed strictly;
// see pktmsg.ParseForm for details.
//
// For a message to be recognized, its message type must have been registered.
// The standard message types can be registered by importing the appropriate
// message-type-specific subpackages of xscmsg.  Alternatively, all standard
// message types can be registered by importing xscmsg/all.
func Recognize(msg *pktmsg.Message, strict bool) Message {
	var form *pktmsg.Form
	if pktmsg.IsForm(msg.Body) {
		form = pktmsg.ParseForm(msg.Body, strict)
	}
	for _, recognizeFunc := range recognizeFuncs {
		if xsc := recognizeFunc(msg, form); xsc != nil {
			return xsc
		}
	}
	return nil
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
