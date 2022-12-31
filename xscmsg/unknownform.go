package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

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
	Tag:     "UNKNOWN",
	Name:    "unknown form",
	Article: "an",
	SubjectFunc: func(_ *Message) string {
		panic("unknown form messages cannot be rendered")
	},
	BodyFunc: func(_ *Message, _ bool) string {
		panic("unknown form messages cannot be rendered")
	},
}
