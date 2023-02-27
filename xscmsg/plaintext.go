package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

// PlainTextTag identifies a plain text message.
const PlainTextTag = "plain"

func makePlainTextMessage(msg *pktmsg.Message) *Message {
	return &Message{
		Type:       &plainTextMessageType,
		RawMessage: msg,
		Fields: []*Field{
			{
				Def:   &plainTextSubjectField,
				Value: msg.Header.Get("Subject"),
			},
			{
				Def:   &plainTextBodyField,
				Value: msg.Body,
			},
		},
	}
}

// CreatePlainTextMessage creates a new, empty plain text message.
func CreatePlainTextMessage() *Message {
	return &Message{
		Type: &plainTextMessageType,
		Fields: []*Field{
			{Def: &plainTextSubjectField},
			{Def: &plainTextBodyField},
		},
	}
}

var plainTextMessageType = MessageType{
	Tag:     PlainTextTag,
	Name:    "plain text message",
	Article: "a",
}

var plainTextSubjectField = FieldDef{
	Tag:   string(FSubject),
	Key:   FSubject,
	Label: "Subject",
	Flags: Required,
}

var plainTextBodyField = FieldDef{
	Tag:   string(FBody),
	Key:   FBody,
	Label: "Body",
	Flags: Required | Multiline,
}
