package checkin

import "github.com/rothskeller/packet/message"

// Type is the type definition for a check-in message.
var Type = message.Type{
	Tag:     "Check-In",
	Name:    "check-in message",
	Article: "a",
	Create:  New,
	Decode:  decode,
}

// CheckIn holds the details of an XSC-standard check-in message.
type CheckIn struct {
	OriginMsgID         string
	Handling            string
	TacticalCallSign    string
	TacticalStationName string
	OperatorCallSign    string
	OperatorName        string
	edit                *checkInEdit
}
type checkInEdit struct {
	OriginMsgID         message.EditField
	Handling            message.EditField
	TacticalCallSign    message.EditField
	TacticalStationName message.EditField
	OperatorCallSign    message.EditField
	OperatorName        message.EditField
	fields              []*message.EditField
}

// Type returns the message type definition.
func (*CheckIn) Type() *message.Type { return &Type }
