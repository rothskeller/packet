// Package plaintext defines the plain text message type.
package plaintext

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Type is the type definition for a plain text message.
var Type = typedmsg.MessageType{
	Tag:       "plain",
	Name:      "plain text message",
	Article:   "a",
	Create:    create,
	Recognize: recognize,
}

// PlainText is a plain text message.
type PlainText struct {
	*xscmsg.Message
	fields []xscmsg.Field
}

// NewPlainText creates a new plain text message.
func NewPlainText() *PlainText {
	return &PlainText{Message: &xscmsg.Message{Message: new(pktmsg.Message)}}
}

func create() typedmsg.Message { return NewPlainText() }

func recognize(base *pktmsg.Message) typedmsg.Message {
	return &PlainText{Message: &xscmsg.Message{Message: base}}
}

// Type returns the type of the message.
func (m *PlainText) Type() *typedmsg.MessageType { return &Type }

// View returns the set of viewable fields of the message.
func (m *PlainText) View() []xscmsg.LabelValue {
	var lvs = m.ViewHeaders()
	lvs = append(lvs,
		xscmsg.LV("Subject", m.SubjectHeader),
		xscmsg.LV("Body", m.Body),
	)
	return lvs
}

// Edit returns the set of editable fields of the message.
func (m *PlainText) Edit() []xscmsg.Field {
	if m.fields == nil {
		m.fields = []xscmsg.Field{
			xscmsg.NewToField(&m.To),
			xscmsg.NewOriginMsgIDField(&m.OriginMsgID),
			xscmsg.NewHandlingField(&m.Handling),
			&subjectField{xscmsg.NewBaseField(&m.Subject)},
			&bodyField{xscmsg.NewBaseField(&m.Body)},
		}
	}
	return m.fields
}

type subjectField struct{ *xscmsg.BaseField }

func (f subjectField) Label() string    { return "Subject" }
func (f subjectField) Size() (w, h int) { return 80, 1 }
func (f subjectField) Problem() string {
	if f.Value() == "" {
		return "The message subject is required."
	}
	return ""
}
func (f subjectField) Help() string {
	return "This is the message subject.  It is required."
}

type bodyField struct{ *xscmsg.BaseField }

func (f bodyField) Label() string    { return "Body" }
func (f bodyField) Size() (w, h int) { return 80, 10 }
func (f bodyField) Problem() string {
	if f.Value() == "" {
		return "The message body is required."
	}
	return ""
}
func (f bodyField) Help() string {
	return "This is the message body.  It is required."
}
