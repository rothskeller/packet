package ics213

import "github.com/rothskeller/packet/message"

// Type is the type definition for an ICS-213 general message form.
var Type = message.Type{
	Tag:     "ICS213",
	Name:    "ICS-213 general message form",
	Article: "an",
	Create:  New,
	Decode:  decode,
}

// ICS213 holds an ICS-213 general message form.
type ICS213 struct {
	PIFOVersion      string
	FormVersion      string
	OriginMsgID      string
	DestinationMsgID string
	Date             string
	Time             string
	Severity         string // removed in v2.2
	Handling         string
	TakeAction       string
	Reply            string
	ReplyBy          string
	FYI              string // removed in v2.2
	ToICSPosition    string
	ToLocation       string
	ToName           string
	ToTelephone      string
	FromICSPosition  string
	FromLocation     string
	FromName         string
	FromTelephone    string
	Subject          string
	Reference        string
	Message          string
	OpRelayRcvd      string
	OpRelaySent      string
	ReceivedSent     string
	OpCall           string
	OpName           string
	TxMethod         string
	OtherMethod      string
	OpDate           string
	OpTime           string
	edit             *ics213Edit
}

// Type returns the message type definition.
func (*ICS213) Type() *message.Type { return &Type }
