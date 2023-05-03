// Package xscmsg defines interfaces honored by the XSC-standard message types,
// as well as base classes and common functions used by their implementations.
// The message types themselves are defined in other packages (mostly
// subpackages of this one).
package xscmsg

import (
	"github.com/rothskeller/packet/typedmsg"
)

// IMessage is the interface supported by all human-originated messages.  It is
// not supported by machine-originated messages like delivery receipts.
type IMessage interface {
	typedmsg.Message
	// GetOpCall returns the value of the operator call sign field if the
	// message type has such a field, or an empty string otherwise.
	GetOpCall() string
	// GetToICSPosition returns the value of the "To ICS Position" field if
	// the message type has such a field, or an empty string otherwise.
	GetToICSPosition() string
	// GetToLocation returns the value of the "To Location" field if the
	// message type has such a field, or an empty string otherwise.
	GetToLocation() string
	// SetDestinationMsgID sets the destination message ID for the message,
	// if the message type has such a field.  It is a no-op otherwise.
	SetDestinationMsgID(string)
	// SetOperator sets the OpCall, OpName, OpDate, and OpTime fields of the
	// message if it has them.  (OpDate and OpTime are set to the current
	// date and time.)  For messages without these fields, this is a no-op.
	SetOperator(received bool, call, name string)
	// SetReference sets the reference field of the message, if it has one.
	// It's a no-op otherwise.
	SetReference(string)
	// SetTactical sets the TacCall and TacName fields of the message if it
	// has them; otherwise, this is a no-op.
	SetTactical(call, name string)
	// Validate checks the validity of the message and returns strings
	// describing the issues.  This is used for received messages, and only
	// checks for validity issues that Outpost and/or PackItForms check for.
	Validate() []string
	// View returns the set of viewable fields in the message.
	View() []LabelValue
}

// A LabelValue is a single (label, value) pair in the list returned by
// Message.View.
type LabelValue struct{ Label, Value string }

// LV is a shortcut for creating a LabelValue.
func LV(label, value string) LabelValue { return LabelValue{label, value} }

// Editable is the interface honored by message types that can be edited by
// humans.
type Editable interface {
	IMessage
	// Edit returns the ordered list of editable fields of the message.
	Edit() []Field
}

// Field is the interface honored by a single editable field of a message.
type Field interface {
	// Label returns the label for the field.  It should be brief, 40
	// characters maximum.
	Label() string
	// Value returns the value of the field.
	Value() string
	// SetValue sets the value of the field.  It also performs any
	// necessary validations and/or recalculations of any affected fields.
	SetValue(string)
	// Size returns the width and height of an input control for the field.
	// It should match the amount of space available for the field in the
	// fillable PDF, or the maximum length of any allowed value for the
	// field.
	Size() (width, height int)
	// Problem returns a string describing the issue with the current value
	// of the field, if it is not valid.  If it is valid, Problem returns an
	// empty string.
	Problem() string
	// Choices returns the set of recommended or restricted values for the
	// field, if any.  If there are none, it returns nil.
	Choices() []string
	// Help returns the help text displayed for the field when requested.
	// Note that the Problem() string is also displayed in the help box if
	// it is not empty, after the text returned by Help().
	Help() string
	// Hint returns a hint to be displayed next to an empty field, generally
	// to give format suggestions, e.g., "MM/DD/YYYY" for a date field.
	Hint() string
}
