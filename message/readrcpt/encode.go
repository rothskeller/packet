package readrcpt

import (
	"fmt"
)

// EncodeSubject encodes the message subject.
func (m *ReadReceipt) EncodeSubject() string {
	return "READ: " + m.MessageSubject
}

// EncodeBody encodes the message body.
func (m *ReadReceipt) EncodeBody() string {
	return fmt.Sprintf("!RR!%s\nYour Message\n\nTo: %s\nSubject: %s\n\nwas read on %[1]s\n",
		m.ReadTime, m.MessageTo, m.MessageSubject)
}
