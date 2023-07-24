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
	edit        *plainTextEdit
	editFields  []*message.EditField
}
type plainTextEdit struct {
	OriginMsgID message.EditField
	Handling    message.EditField
	Subject     message.EditField
	Body        message.EditField
}

// Type returns the message type definition.
func (*PlainText) Type() *message.Type { return &Type }
