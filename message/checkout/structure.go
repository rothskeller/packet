package checkout

import "github.com/rothskeller/packet/message"

// Type is the type definition for a check-out message.
var Type = message.Type{
	Tag:     "Check-Out",
	Name:    "check-out message",
	Article: "a",
	Create:  New,
	Decode:  decode,
}

// CheckOut holds the details of an XSC-standard check-out message.
type CheckOut struct {
	OriginMsgID         string
	Handling            string
	TacticalCallSign    string
	TacticalStationName string
	OperatorCallSign    string
	OperatorName        string
	edit                *checkOutEdit
}
type checkOutEdit struct {
	OriginMsgID         message.EditField
	Handling            message.EditField
	TacticalCallSign    message.EditField
	TacticalStationName message.EditField
	OperatorCallSign    message.EditField
	OperatorName        message.EditField
	fields              []*message.EditField
}

// Type returns the message type definition.
func (*CheckOut) Type() *message.Type { return &Type }
