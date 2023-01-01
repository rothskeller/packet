package xscmsg

import (
	"fmt"

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
	Tag:        string(FSubject),
	Key:        FSubject,
	Label:      "Subject",
	Validators: []Validator{ValidateRequired},
}

var plainTextBodyField = FieldDef{
	Tag:        string(FBody),
	Key:        FBody,
	Label:      "Body",
	Validators: []Validator{ValidateRequired},
}

// ValidateRequired is a Validator that verifies that the field has a value.
func ValidateRequired(f *Field, m *Message, strict bool) string {
	if f.Value == "" {
		return fmt.Sprintf("The %q field must have a value.", f.Def.Tag)
	}
	return ""
}
