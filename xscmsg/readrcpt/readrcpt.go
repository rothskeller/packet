// Package readrcpt handles read receipt messages.
package readrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type is the type definition for a read receipt.
var Type = message.Type{
	Tag:     "READ",
	Name:    "read receipt",
	Article: "a",
}

func init() {
	message.Register(&Type, decode, nil)
}

// ReadReceipt holds the details of an XSC-standard read receipt message.
type ReadReceipt struct {
	message.BaseMessage
	MessageTo      string
	MessageSubject string
	ReadTime       string
	ExtraText      string
}

// New creates a new read receipt message.
func New() (m *ReadReceipt) {
	m = &ReadReceipt{BaseMessage: message.BaseMessage{Type: &Type}}
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
			Label:    "Read At",
			Value:    &m.ReadTime,
			Presence: message.Required,
		}),
		message.NewMultilineField(&message.Field{
			Label: "Extra Text",
			Value: &m.ExtraText,
		}),
	}
	return m
}

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^\n*!RR!(.+)\n.*\n\nTo: (.+)`)

func decode(env *envelope.Envelope, body string, form *message.PIFOForm, _ int) message.Message {
	if !strings.HasPrefix(env.SubjectLine, "READ: ") || form != nil {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(body); match != nil {
		rr := New()
		rr.MessageSubject = env.SubjectLine[6:]
		rr.MessageTo = match[2]
		rr.ReadTime = match[2]
		rr.ExtraText = strings.TrimSpace(body[len(match[0]):])
		return rr
	}
	return nil
}

// EncodeSubject encodes the message subject.
func (m *ReadReceipt) EncodeSubject() string {
	return "READ: " + m.MessageSubject
}

// EncodeBody encodes the message body.
func (m *ReadReceipt) EncodeBody() string {
	return fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n",
		m.ReadTime, m.MessageTo, m.MessageSubject)
}
