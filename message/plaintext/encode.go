package plaintext

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (m *PlainText) EncodeSubject() string {
	return common.EncodeSubject(m.OriginMsgID, m.Handling, "", m.Subject)
}

// EncodeBody encodes the message body.
func (m *PlainText) EncodeBody() string {
	var body = m.Body
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return body
}
