// Package delivrcpt handles delivery receipt messages.
package delivrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
)

// Type is the type definition for a delivery receipt.
var Type = message.Type{
	Tag:     "DELIVERED",
	Name:    "delivery receipt",
	Article: "a",
}

func init() {
	Type.Decode = decode
}

// DeliveryReceipt holds the details of an XSC-standard delivery receipt
// message.
type DeliveryReceipt struct {
	message.BaseMessage
	MessageTo      string
	MessageSubject string
	LocalMessageID string
	DeliveredTime  string
	ExtraText      string
}

// New creates a new delivery receipt message.
func New() (m *DeliveryReceipt) {
	m = &DeliveryReceipt{BaseMessage: message.BaseMessage{Type: &Type}}
	m.Fields = []*message.Field{
		message.NewTextField(&message.Field{
			Label:    "Message To",
			Value:    &m.MessageTo,
			Presence: message.Required,
		}),
		message.NewTextField(&message.Field{
			Label: "Message Subject",
			Value: &m.MessageSubject,
		}),
		message.NewTextField(&message.Field{
			Label:    "Delivered At",
			Value:    &m.DeliveredTime,
			Presence: message.Required,
		}),
		message.NewMessageNumberField(&message.Field{
			Label:    "Message Number",
			Value:    &m.LocalMessageID,
			Presence: message.Required,
		}),
		message.NewMultilineField(&message.Field{
			Label: "Extra Text",
			Value: &m.ExtraText,
		}),
	}
	return m
}

// deliveryReceiptRE matches the first lines of a delivery receipt message.  Its
// substrings are the local message ID, the delivery time, and the To address.
var deliveryReceiptRE = regexp.MustCompile(`^!LMI!([^!]+)!DR!(.+)\n.*\nTo: (.+)\nSubject:.*\nwas delivered on.*\nRecipient's Local.*\n`)

func decode(subject, body string) *DeliveryReceipt {
	if !strings.HasPrefix(subject, "DELIVERED: ") {
		return nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(body); match != nil {
		dr := New()
		dr.MessageSubject = subject[11:]
		dr.MessageTo = match[3]
		dr.LocalMessageID = match[1]
		dr.DeliveredTime = match[2]
		dr.ExtraText = strings.TrimSpace(body[len(match[0]):])
		return dr
	}
	return nil
}

// EncodeSubject encodes the message subject.
func (m *DeliveryReceipt) EncodeSubject() string {
	return "DELIVERED: " + m.MessageSubject
}

// EncodeBody encodes the message body.
func (m *DeliveryReceipt) EncodeBody() string {
	et := m.ExtraText
	if et != "" {
		et = "\n" + et
		if !strings.HasSuffix(et, "\n") {
			et += "\n"
		}
	}
	return fmt.Sprintf(
		"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n%[5]s",
		m.LocalMessageID, m.DeliveredTime, m.MessageTo, m.MessageSubject, et)
}
