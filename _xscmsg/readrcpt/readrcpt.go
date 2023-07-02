package readrcpt

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
)

// Type is the MessageType for a read receipt.
var Type = typedmsg.MessageType{
	Tag:       "READ",
	Name:      "read receipt",
	Article:   "a",
	Create:    NewReadReceipt,
	Recognize: recognize,
}

// ReadReceipt is a read receipt message.
type ReadReceipt struct {
	*pktmsg.Message
	ReadTo      string
	ReadSubject string
	ReadTime    string
}

// NewReadReceipt creates a new read receipt message.
func NewReadReceipt() typedmsg.Message {
	return &ReadReceipt{Message: new(pktmsg.Message)}
}

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

func recognize(pm *pktmsg.Message) typedmsg.Message {
	var match []string

	if !strings.HasPrefix(pm.SubjectHeader, "READ: ") {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(pm.Body); match == nil {
		return nil
	}
	return &ReadReceipt{
		Message:     pm,
		ReadTo:      match[2],
		ReadSubject: pm.SubjectHeader[6:],
		ReadTime:    match[1],
	}
}

// Type returns the type of the message.
func (m *ReadReceipt) Type() *typedmsg.MessageType { return &Type }

// Save returns the message encoded for saving to local storage.
func (m *ReadReceipt) Save() string {
	m.encode()
	return m.Message.Save()
}

// Transmit returns the message encoded for transmission.
func (m *ReadReceipt) Transmit() ([]string, string, string) {
	m.encode()
	return m.Message.Transmit()
}

func (m *ReadReceipt) encode() {
	m.Subject = "READ: " + m.ReadSubject
	m.Body = fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n",
		m.ReadTime, m.ReadTo, m.ReadSubject,
	)
}
