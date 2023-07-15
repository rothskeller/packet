package delivrcpt

import "github.com/rothskeller/packet/message"

// Type is the type definition for a delivery receipt.
var Type = message.Type{
	Tag:     "DELIVERED",
	Name:    "delivery receipt",
	Article: "a",
	Create:  nil,
	Decode:  decode,
}

// DeliveryReceipt holds the details of an XSC-standard delivery receipt
// message.
type DeliveryReceipt struct {
	MessageTo      string
	MessageSubject string
	LocalMessageID string
	DeliveredTime  string
	ExtraText      string
}

// Type returns the message type definition.
func (*DeliveryReceipt) Type() *message.Type { return &Type }
