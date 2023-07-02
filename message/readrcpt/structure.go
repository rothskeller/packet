package readrcpt

import "github.com/rothskeller/packet/message"

// Type is the type definition for a read receipt.
var Type = message.Type{
	Tag:     "READ",
	Name:    "read receipt",
	Article: "a",
	Create:  nil,
	Decode:  decode,
}

// ReadReceipt holds the details of an XSC-standard read receipt message.
type ReadReceipt struct {
	MessageTo      string
	MessageSubject string
	ReadTime       string
}

// Type returns the message type definition.
func (*ReadReceipt) Type() *message.Type { return &Type }
