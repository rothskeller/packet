// Package msgui defines a Message that is viewable and editable by end users.
package msgui

// Message is a packet message that is viewable and editable by end users.
type Message struct {
	// Ident is the identification of the message (e.g., its local message
	// ID), displayed in the title bar of a message viewer or editor.
	Ident string
	// Type is a name for the type of message, displayed next to the Ident
	// in the title bar of a message viewer or editor.
	Type string
	// Fields is a set of fields in the message.
	Fields []*Field
	// Changed is called whenever any editable field is changed, with the
	// changed field passed as a parameter.  If Changed is nil, that
	// indicates that the message is viewable but not editable.
	Changed func(*Field)
}

// A Field is a single viewable and perhaps editable field of a message.
type Field struct {
	// Label is the label for the field.  It should be brief, 40 characters
	// maximum.
	Label string
	// Value is the value of the field.
	Value string
	// Width is the expected width of an input control for this field, or
	// zero for indeterminate.
	Width uint8
	// Height is the expected height of an input control for this field, or
	// zero for indeterminate.
	Height uint8
	// Editable is a flag indicating that the value of the field can be
	// edited by a user.
	Editable bool
	// Valid indicates whether the current value of the field is a valid
	// value.  It is significant only when Editable is true.  (If it is
	// false, Help should include a description of the problem.)
	Valid bool
	// Choices is an optional set of values for the field, from which the
	// user can choose.  It is significant only when Editable is true.
	Choices []string
	// Help is help text displayed for the field when requested.  It is
	// significant only when Editable is true.
	Help string
	// Hint is a hint displayed next to an empty field, generally to give
	// format suggestions, e.g., "MM/DD/YYYY" for a date field.  It is
	// significant only when Editable is true and Value is empty.
	Hint string
}
