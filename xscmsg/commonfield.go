package xscmsg

import "github.com/rothskeller/packet/pktmsg"

// NewHandlingField creates a new handling order field.
func NewHandlingField(vp *string) *HandlingField {
	return &HandlingField{&BaseField{vp}}
}

// HandlingField is the message handling order field.
type HandlingField struct{ *BaseField }

var handlingChoices = []string{"ROUTINE", "PRIORITY", "IMMEDIATE"}

// Label returns the display label for the field.
func (f HandlingField) Label() string { return "Handling Order" }

// SetValue sets the value for the field.
func (f *HandlingField) SetValue(value string) {
	f.BaseField.SetValue(ChoicesExpand(value, handlingChoices))
}

// Size returns the display size of the field.
func (f HandlingField) Size() (width, height int) { return 9, 1 }

// Problem returns a string describing the validation problem with the field, if
// any.
func (f HandlingField) Problem() string {
	value := f.Value()
	if value == "" {
		return "The handling order is required."
	}
	if !ChoicesValid(value, handlingChoices) {
		return "The handling order must be one of ROUTINE, PRIORITY, or IMMEDIATE."
	}
	return ""
}

// Choices returns the set of recommended or required values for the field.
func (f HandlingField) Choices() []string { return handlingChoices }

// Help returns help text for the field.
func (f HandlingField) Help() string {
	return "This is the message handling order, which determines how fast it needs to be delivered.  It must be one of ROUTINE, PRIORITY, or IMMEDIATE."
}

////////////////////////////////////////////////////////////////////////////////

// NewOriginMsgIDField creates a new origin message ID field.
func NewOriginMsgIDField(vp *string) *OriginMsgIDField {
	return &OriginMsgIDField{&BaseField{vp}}
}

// OriginMsgIDField is the origin message ID field.
type OriginMsgIDField struct{ *BaseField }

// Label returns the display label for the field.
func (f OriginMsgIDField) Label() string { return "Origin Message Number" }

// SetValue sets the value for the field.
func (f *OriginMsgIDField) SetValue(value string) {
	f.BaseField.SetValue(MessageIDClean(value))
}

// Size returns the display size of the field.
func (f OriginMsgIDField) Size() (width, height int) { return MessageIDWidth, 1 }

// Problem returns a string describing the validation problem with the field, if
// any.
func (f OriginMsgIDField) Problem() string {
	value := f.Value()
	if value == "" {
		return "The origin message number is required."
	}
	if !MessageIDValid(value) {
		return "The origin message number is not valid.  " + MessageIDHelp
	}
	return ""
}

// Help returns help text for the field.
func (f OriginMsgIDField) Help() string {
	return "This is the message number assigned to the message by the origin station.  It is required.  " + MessageIDHelp
}

////////////////////////////////////////////////////////////////////////////////

// NewToField creates a new To: address list field.
func NewToField(vp *pktmsg.Addresses) *ToField {
	return &ToField{&BaseField{(*string)(vp)}}
}

// ToField is the To: address list field.
type ToField struct{ *BaseField }

// Label returns the display label for the field.
func (f ToField) Label() string { return "To" }

// SetValue sets the value for the field.
func (f *ToField) SetValue(value string) {
	f.BaseField.SetValue(AddressesClean(value))
}

// Size returns the display size of the field.
func (f ToField) Size() (width, height int) { return 80, 1 }

// Problem returns a string describing the validation problem with the field, if
// any.
func (f ToField) Problem() string {
	value := f.Value()
	if value == "" {
		return "The \"To\" address is required."
	}
	if !AddressesValid(value) {
		return "The \"To\" address list is not valid.  " + AddressesHelp
	}
	return ""
}

// Help returns help text for the field.
func (f ToField) Help() string {
	return "This is the list of addresses to which the message is sent.  At least one address is required.  " + AddressesHelp
}
