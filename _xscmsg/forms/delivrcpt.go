package xscmsg

import (
	"fmt"
	"regexp"
	"strings"
)

// DeliveryReceipt holds the details of an XSC-standard delivery receipt
// message.
type DeliveryReceipt struct {
	MessageTo      string
	MessageSubject string
	LocalMessageID string
	DeliveredTime  string
}

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)`)

// DecodeDeliveryReceipt decodes the supplied message contents if they are a
// delivery receipt.  It returns nil if the message contents are not a
// well-formed delivery receipt.
func DecodeDeliveryReceipt(subject, body string) *DeliveryReceipt {
	if !strings.HasPrefix(subject, "DELIVERED: ") {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(body); match != nil {
		return &DeliveryReceipt{
			MessageSubject: subject[11:],
			MessageTo:      match[3],
			LocalMessageID: match[1],
			DeliveredTime:  match[2],
		}
	}
	return nil
}

// Encode encodes the message contents.
func (m *DeliveryReceipt) Encode() (subject, body string) {
	return "DELIVERED: " + m.MessageSubject, fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n",
		m.LocalMessageID, m.DeliveredTime, m.MessageTo, m.MessageSubject)
}
