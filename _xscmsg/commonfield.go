package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

// NewHandlingField creates a new handling order field.
func NewHandlingField(vp *string) Field {
	return WrapRequiredField(WrapRestrictedField(&handlingField{&baseField{vp}}))
}

type handlingField struct{ *baseField }

var handlingChoices = []string{"ROUTINE", "PRIORITY", "IMMEDIATE"}

func (f handlingField) Label() string     { return "Handling Order" }
func (f handlingField) Choices() []string { return handlingChoices }
func (f handlingField) Help() string {
	return "This is the message handling order, which determines how fast it needs to be delivered."
}

////////////////////////////////////////////////////////////////////////////////

// NewOriginMsgIDField creates a new origin message ID field.
func NewOriginMsgIDField(vp *string) Field {
	return WrapRequiredField(WrapPacketMessageIDField(&originMsgIDField{&baseField{vp}}))
}

type originMsgIDField struct{ *baseField }

func (f originMsgIDField) Label() string { return "Origin Message Number" }
func (f originMsgIDField) Help() string {
	return "This is the message number assigned to the message by the origin station."
}

////////////////////////////////////////////////////////////////////////////////

// NewToField creates a new To: address list field.
func NewToField(vp *pktmsg.Addresses) Field {
	return WrapRequiredField(WrapAddressListField(&toField{&baseField{(*string)(vp)}}))
}

type toField struct{ *baseField }

func (f toField) Label() string { return "To" }
func (f toField) Help() string {
	return "This is the list of addresses to which the message is sent."
}
