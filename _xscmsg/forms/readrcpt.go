package xscmsg

import (
	"fmt"
	"regexp"
	"strings"
)

// ReadReceipt holds the details of an XSC-standard read receipt message.
type ReadReceipt struct {
	MessageTo      string
	MessageSubject string
	ReadTime       string
}

// readReceiptRE matches the first lines of a read receipt message.  Its
// substrings are the read time and the To address.
var readReceiptRE = regexp.MustCompile(`^!RR!(.+)\n.*\n\nTo: (.+)`)

// DecodeReadReceipt decodes the supplied message contents if they are a
// read receipt.  It returns nil if the message contents are not a
// well-formed read receipt.
func DecodeReadReceipt(subject, body string) *ReadReceipt {
	if !strings.HasPrefix(subject, "READ: ") {
		return nil
	}
	if match := readReceiptRE.FindStringSubmatch(body); match != nil {
		return &ReadReceipt{
			MessageSubject: subject[6:],
			MessageTo:      match[2],
			ReadTime:       match[1],
		}
	}
	return nil
}

// Encode encodes the message contents.
func (m *ReadReceipt) Encode() (subject, body string) {
	return "READ: " + m.MessageSubject,
		fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n",
			m.ReadTime, m.MessageTo, m.MessageSubject)
}
