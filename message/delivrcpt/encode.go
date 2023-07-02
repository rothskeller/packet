package delivrcpt

import (
	"fmt"
)

// EncodeSubject encodes the message subject.
func (m *DeliveryReceipt) EncodeSubject() string {
	return "DELIVERED: " + m.MessageSubject
}

// EncodeBody encodes the message body.
func (m *DeliveryReceipt) EncodeBody() string {
	return fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n",
		m.LocalMessageID, m.DeliveredTime, m.MessageTo, m.MessageSubject)
}
