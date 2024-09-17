// Package bulletin handles bulletin messages.
package bulletin

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/plaintext"
)

// Type is the type definition for a bulletin message.
var Type = message.Type{
	Tag:     "bulletin",
	Name:    "bulletin",
	Article: "a",
}

func init() {
	Type.Create = New
}

// Bulletin holds the details of a bulletin message.
type Bulletin struct {
	message.BaseMessage
	Subject string
	Body    string
}

// New creates a new bulletin message.
func New() (m *Bulletin) {
	m = &Bulletin{BaseMessage: message.BaseMessage{Type: &Type}}
	m.BaseMessage.FSubject = &m.Subject
	m.BaseMessage.FBody = &m.Body
	m.Fields = []*message.Field{
		message.NewTextField(&message.Field{
			Label:     "Subject",
			Value:     &m.Subject,
			Presence:  message.Required,
			EditWidth: 80,
			EditHelp:  `This is the subject of the bulletin.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Message",
			Value:     &m.Body,
			Presence:  message.Required,
			EditWidth: 80,
			EditHelp:  `This is the body of the bulletin.  It is required.`,
		}),
	}
	return m
}

// FromPlainText converts a plain text message into a Bulletin.  It can be
// called when a plain text message is received by retrieval from a bulletin
// area.
func FromPlainText(pt *plaintext.PlainText) (m *Bulletin) {
	m = New()
	m.Subject, m.Body = pt.Subject, pt.Body
	return m
}

func (m *Bulletin) EncodeSubject() string { return m.Subject }
func (m *Bulletin) EncodeBody() string    { return m.Body }

func (m *Bulletin) RenderPDF(env *envelope.Envelope, filename string) (err error) {
	if plaintext.RenderPlainPDF == nil {
		return message.ErrNotSupported
	}
	return plaintext.RenderPlainPDF(env, "BULLETIN", filename[:len(filename)-4], m.Body, filename)
}
