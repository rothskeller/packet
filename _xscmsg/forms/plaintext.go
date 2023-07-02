package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// PlainText holds the details of a plain text message.
type PlainText struct {
	OriginMsgID string
	Handling    string
	Subject     string
	Body        string
}

// DecodePlainText decodes the supplied message contents as a plain text
// message.
func DecodePlainText(subject, body string) *PlainText {
	var pt PlainText

	pt.OriginMsgID, _, pt.Handling, _, pt.Subject = xscsubj.Decode(subject)
	pt.Body = body
	return &pt
}

// Encode encodes the message contents.
func (m *PlainText) Encode() (subject, body string) {
	if !strings.HasSuffix(m.Body, "\n") {
		m.Body += "\n"
	}
	return xscsubj.Encode(m.OriginMsgID, m.Handling, "", m.Subject), m.Body
}
