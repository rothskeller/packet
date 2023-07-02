package plaintext

import "github.com/rothskeller/packet/message"

// Type is the type definition for a plain text message.
var Type = message.Type{
	Tag:     "plain",
	Name:    "plain text message",
	Article: "a",
	Create:  New,
	Decode:  decode,
}

// PlainText holds the details of a plain text message.
type PlainText struct {
	OriginMsgID string
	Handling    string
	Subject     string
	Body        string
}

// Type returns the message type definition.
func (*PlainText) Type() *message.Type { return &Type }
