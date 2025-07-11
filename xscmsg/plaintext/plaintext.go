// Package plaintext handles plain text messages.
package plaintext

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type is the type definition for a plain text message.
var Type = message.Type{
	Tag:     "plain",
	Name:    "plain text message",
	Article: "a",
}

func init() {
	message.Register(&Type, decode, New)
}

// PlainText holds the details of a plain text message.
type PlainText struct {
	message.BaseMessage
	OriginMsgID string
	Handling    string
	Subject     string
	Body        string
}

// New creates a new plain text message.
func New() message.Message {
	m := &PlainText{BaseMessage: message.BaseMessage{Type: &Type}}
	m.FOriginMsgID = &m.OriginMsgID
	m.FHandling = &m.Handling
	m.FSubject = &m.Subject
	m.FBody = &m.Body
	m.Fields = []*message.Field{
		message.NewMessageNumberField(&message.Field{
			Label:    "Origin Message Number",
			Value:    &m.OriginMsgID,
			Presence: message.Required,
			EditHelp: `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Handling",
			Value:    &m.Handling,
			Choices:  message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: message.Required,
			EditHelp: `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Subject",
			Value:     &m.Subject,
			Presence:  message.Required,
			EditWidth: 80,
			EditHelp:  `This is the subject of the message.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Message",
			Value:     &m.Body,
			Presence:  message.Required,
			EditWidth: 80,
			EditHelp:  `This is the body of the message.  It is required.`,
		}),
	}
	return m
}

// This function is called to find out whether an incoming message matches this
// type.  It should return the decoded message if it belongs to this type, or
// nil if it doesn't.
func decode(env *envelope.Envelope, body string, form *message.PIFOForm, pass int) message.Message {
	if pass != 2 || form != nil || env.Bulletin {
		return nil
	}
	f := New().(*PlainText)
	var fauxFormType string
	f.OriginMsgID, _, f.Handling, fauxFormType, f.Subject = message.DecodeSubject(env.SubjectLine)
	if fauxFormType != "" {
		f.Subject = fauxFormType + "_" + f.Subject
	}
	if h := message.DecodeHandlingMap[f.Handling]; h != "" {
		f.Handling = h
	}
	f.Body = body
	return f
}

func (m *PlainText) EncodeBody() string { return m.Body }

var RenderPlainPDF func(env *envelope.Envelope, label, lmi, body, filename string) error

func (m *PlainText) RenderPDF(env *envelope.Envelope, filename string) (err error) {
	if RenderPlainPDF == nil {
		return message.ErrNotSupported
	}
	return RenderPlainPDF(env, "PLAIN TEXT MESSAGE", filename[:len(filename)-4], m.Body, filename)
}
