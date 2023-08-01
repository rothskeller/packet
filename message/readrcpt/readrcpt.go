// Package readrcpt handles read receipt messages.
package readrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/basemsg"
)

// Type is the type definition for a read receipt.
var Type = message.Type{
	Tag:     "READ",
	Name:    "read receipt",
	Article: "a",
}

func init() {
	Type.Decode = decode
}

// ReadReceipt holds the details of an XSC-standard read receipt message.
type ReadReceipt struct {
	basemsg.BaseMessage
	MessageTo      string
	MessageSubject string
	ReadTime       string
	ExtraText      string
}

// New creates a new read receipt message.
func New() (m *ReadReceipt) {
	return &ReadReceipt{BaseMessage: basemsg.BaseMessage{MessageType: &Type}}
}

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

func decode(subject, body string) *ReadReceipt {
	if !strings.HasPrefix(subject, "READ: ") {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(body); match != nil {
		rr := New()
		rr.MessageSubject = subject[6:]
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
