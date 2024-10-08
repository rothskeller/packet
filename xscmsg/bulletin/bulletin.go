// Package bulletin handles plain text messages.
package bulletin

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/plaintext"
)

// Type is the type definition for a bulletin message.
var Type = message.Type{
	Tag:     "bulletin",
	Name:    "bulletin message",
	Article: "a",
}

func init() {
	message.Register(&Type, decode, New)
}

// Bulletin holds the details of a bulletin message.
type Bulletin struct {
	message.BaseMessage
	Subject string
	Body    string
}

// New creates a new plain text message.
func New() message.Message {
	var m = &Bulletin{BaseMessage: message.BaseMessage{Type: &Type}}
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

// This function is called to find out whether an incoming message matches this
// type.  It should return the decoded message if it belongs to this type, or
// nil if it doesn't.
func decode(env *envelope.Envelope, body string, form *message.PIFOForm, pass int) message.Message {
	if pass != 2 || form != nil || !env.Bulletin {
		return nil
	}
	var f = New().(*Bulletin)
	f.Subject = env.SubjectLine
	f.Body = body
	return f
}

func (m *Bulletin) EncodeBody() string { return m.Body }

func (m *Bulletin) RenderPDF(env *envelope.Envelope, filename string) (err error) {
	if plaintext.RenderPlainPDF == nil {
		return message.ErrNotSupported
	}
	return plaintext.RenderPlainPDF(env, "BULLETIN", filename[:len(filename)-4], m.Body, filename)
}
