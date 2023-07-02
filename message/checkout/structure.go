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
}

// Type returns the message type definition.
func (*CheckOut) Type() *message.Type { return &Type }
