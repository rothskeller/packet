package delivrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
)

// Type is the MessageType for a delivery receipt.
var Type = typedmsg.MessageType{
	Tag:       "DELIVERED",
	Name:      "delivery receipt",
	Article:   "a",
	Create:    create,
	Recognize: recognize,
}

// DeliveryReceipt is a delivery receipt message.
type DeliveryReceipt struct {
	*pktmsg.Message
	DeliveredTo      string
	DeliveredSubject string
	LocalMessageID   string
	DeliveredTime    string
}

// NewDeliveryReceipt creates a new delivery receipt message.
func NewDeliveryReceipt() *DeliveryReceipt {
	return &DeliveryReceipt{Message: new(pktmsg.Message)}
}

func create() typedmsg.Message { return NewDeliveryReceipt() }

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)`)

func recognize(pm *pktmsg.Message) typedmsg.Message {
	var match []string

	if !strings.HasPrefix(pm.SubjectHeader, "DELIVERED: ") {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(pm.Body); match == nil {
		return nil
	}
	return &DeliveryReceipt{
		Message:          pm,
		DeliveredTo:      match[3],
		DeliveredSubject: pm.SubjectHeader[11:],
		LocalMessageID:   match[1],
		DeliveredTime:    match[2],
	}
}

// Type returns the type of the message.
func (m *DeliveryReceipt) Type() *typedmsg.MessageType { return &Type }

// Save returns the message encoded for saving to local storage.
func (m *DeliveryReceipt) Save() string {
	m.encode()
	return m.Message.Save()
}

// Transmit returns the message encoded for transmission.
func (m *DeliveryReceipt) Transmit() ([]string, string, string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *DeliveryReceipt) encode() {
	m.Subject = "DELIVERED: " + m.DeliveredSubject
	m.Body = fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n",
		m.LocalMessageID, m.DeliveredTime, m.DeliveredTo, m.DeliveredSubject,
	)
}
