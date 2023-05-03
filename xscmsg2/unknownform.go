package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

func recognizeUnknownForm(*pktmsg.Message, *pktmsg.Form) Message {}
func makeUnknownFormMessage(msg *pktmsg.Message, form *pktmsg.Form) *Message {
	m := Message{
		Type:       &unknownFormMessageType,
		RawMessage: msg,
		RawForm:    form,
	}
	for _, rf := range form.Fields {
		m.Fields = append(m.Fields, &Field{
			Def:   &FieldDef{Tag: rf.Tag, Label: rf.Tag},
			Value: rf.Value,
		})
	}
	return &m
}

var unknownFormMessageType = MessageType{
	Tag:         "UNKNOWN",
	Name:        "form of unknown type",
	Article:     "a",
	SubjectFunc: func(m *Message) string { return m.RawMessage.Header.Get("Subject") },
	BodyFunc: func(m *Message) string {
		// This is a copy of xscform.EncodeBody.  It was easier to copy
		// it than to work around the import loop to reuse it.
		var form = pktmsg.Form{
			PIFOVersion: m.RawForm.PIFOVersion,
			FormType:    m.RawForm.FormType,
			FormVersion: m.RawForm.FormVersion,
		}
		for _, f := range m.Fields {
			if f.Value != "" {
				form.Fields = append(form.Fields, pktmsg.FormField{Tag: f.Def.Tag, Value: f.Value})
			}
		}
		return form.Encode()
	},
}
