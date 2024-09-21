// Package message contains the interfaces and registry for packet message
// types.  The definitions in this package can be used to register package
// message types and to itemize the registered types.
package message

import (
	"time"

	"github.com/rothskeller/packet/envelope"
)

// Message is the interface that all message types implement.  In addition to
// implementing this interface, all message types must embed BaseMessage, which
// provides shared functionality.
type Message interface {
	// Base returns the BaseMessage embedded in the message.
	Base() *BaseMessage
	// EncodeSubject encodes the message subject line.
	EncodeSubject() string
	// EncodeBulletinSubject encodes the message subject line, appropriately
	// for posting as a bulletin.
	EncodeBulletinSubject() string
	// EncodeBody encodes the message body, suitable for transmission or
	// storage.
	EncodeBody() string
	// Validate checks the contents of the message for compliance with rules
	// enforced by standard Santa Clara County packet software (Outpost and
	// PackItForms).  It returns a list of strings describing problems that
	// those programs would flag or block.
	PIFOValid() (problems []string)
	// Compare compares two messages.  It returns a score indicating how
	// closely they match, and the detailed comparisons of each field in the
	// message.  The comparison is not symmetric:  the receiver of the call
	// is the "expected" message and the argument is the "actual" message.
	Compare(actual Message) (score, outOf int, fields []*CompareField)
	// RenderPDF renders the message as a PDF file with the specified
	// filename, overwriting any existing file with that name.  This method
	// will return ErrNotSupported for message types that do not support PDF
	// rendering.  Note that the program needs to be built with "-tags
	// packetpdf" in order for any message types to support PDF rendering.
	RenderPDF(env *envelope.Envelope, filename string) error
	// SetOperator sets the operator only fields of the message, if it has
	// them.
	SetOperator(opcall, opname string, received bool)
	// Editable returns whether the message type supports editing.
	Editable() bool
}

// BaseMessage is the type underlying all packet messages, providing their
// shared functionality.  Every message, regardless of type, embeds a
// BaseMessage and provides access to it through the Base() method on the
// message type.
type BaseMessage struct {
	// Type is the type definition for the message type.
	Type *Type
	// PIFOVersion is the PIFO version found when decoding the message.  It
	// is set only for messages with PIFO encoding.
	PIFOVersion string
	// Fields is an ordered list of fields in the message.  This is the
	// core of the shared message functionality:  most operations are
	// implemented by iterating through these fields.
	Fields []*Field

	// Pointers to key fields.

	// FOriginMsgID points to the value of the Origin Message ID field.  It
	// is nil for message types that do not have that field.
	FOriginMsgID *string
	// FDestinationMsgID points to the value of the Destination Message ID
	// field.  It is nil for message types that do not have that field.
	FDestinationMsgID *string
	// FMessageDate points to the value of the message date field.  It is
	// nil for message types that do not have that field.
	FMessageDate *string
	// FMessageTime points to the value of the message time field.  It is
	// nil for message types that do not have that field.
	FMessageTime *string
	// FHandling points to the value of the Handling field.  It
	// is nil for message types that do not have that field.
	FHandling *string
	// FSubject points to the value of the field of the message that will
	// get propagated to the message's subject line.  It is nil for message
	// types that do not have any such field.
	FSubject *string
	// RestrictedSubject is a flag indicating that the FSubject field allows
	// only certain restricted values.  (It will not be populated with the
	// subject of a message being replied to, unless that message is of the
	// same type.)
	RestrictedSubject bool
	// FToICSPosition points to the value of the To ICS Position field.  It
	// is nil for message types that do not have that field.
	FToICSPosition *string
	// FToLocation points to the value of the To Location field.  It
	// is nil for message types that do not have that field.
	FToLocation *string
	// FFromICSPosition points to the value of the From ICS Position field.
	// It is nil for message types that do not have that field.
	FFromICSPosition *string
	// FFromLocation points to the value of the From Location field.  It
	// is nil for message types that do not have that field.
	FFromLocation *string
	// FReference points to the value of the Reference field.  It is nil for
	// message types that do not have that field.
	FReference *string
	// FTacCall points to the value of the Tactical Call Sign field.  It is
	// nil for message types that do not have that field.
	FTacCall *string
	// FTacName points to the value of the Tactical Station Name field.  It
	// is nil for message types that do not have that field.
	FTacName *string
	// FOpCall points to the value of the Operator Call Sign field.  It is
	// nil for message types that do not have that field.
	FOpCall *string
	// FOpName points to the value of the Operator Name field.  It is nil
	// for message types that do not have that field.
	FOpName *string
	// FOpDate points to the value of the Operator Date field.  It is nil
	// for message types that do not have that field.
	FOpDate *string
	// FOpTime points to the value of the Operator Time field.  It is nil
	// for message types that do not have that field.
	FOpTime *string
	// FBody points to the value of the most prominent, or first, multi-line
	// text field of the message.  It is nil for message types that do not
	// have any such field.
	FBody *string
}

// Base returns the BaseMessage structure for the message.
func (bm *BaseMessage) Base() *BaseMessage { return bm }

// SetOperator sets the operator only fields of the message, if it has them.
func (bm *BaseMessage) SetOperator(opcall, opname string, received bool) {
	if bm.FOpCall != nil {
		*bm.FOpCall = opcall
	}
	if bm.FOpName != nil {
		*bm.FOpName = opname
	}
	if bm.FOpDate != nil {
		*bm.FOpDate = time.Now().Format("01/02/2006")
	}
	if bm.FOpTime != nil {
		*bm.FOpTime = time.Now().Format("15:04")
	}
}

// Editable returns whether the message type supports editing.
func (bm *BaseMessage) Editable() bool {
	for _, f := range bm.Fields {
		if f.EditHelp != "" {
			return true
		}
	}
	return false
}
