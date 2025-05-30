// Package unkform handles PackItForms messages with an unrecognized form type.
package unkform

import (
	"sort"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type is the type definition for an unrecognized form message.
var Type = message.Type{
	Tag:     "UNKNOWN",
	Name:    "unrecognized form message",
	Article: "an",
}

func init() {
	message.Register(&Type, decode, nil)
}

// UnknownForm holds the details of an unrecognized form message.
type UnknownForm struct {
	message.BaseMessage
	OriginMsgID  string
	Handling     string
	Subject      string
	TaggedValues map[string]string
}

// This function is called to find out whether an incoming message matches this
// type.  It should return the decoded message if it belongs to this type, or
// nil if it doesn't.
func decode(env *envelope.Envelope, body string, form *message.PIFOForm, pass int) message.Message {
	if pass != 2 || form == nil {
		return nil
	}
	var typeCopy = Type
	typeCopy.HTML = form.HTMLIdent
	typeCopy.Version = form.FormVersion
	var f = &UnknownForm{BaseMessage: message.BaseMessage{Type: &typeCopy}}
	f.BaseMessage.FOriginMsgID = &f.OriginMsgID
	f.BaseMessage.FHandling = &f.Handling
	f.BaseMessage.FSubject = &f.Subject
	f.OriginMsgID, _, f.Handling, f.Type.Tag, f.Subject = message.DecodeSubject(env.SubjectLine)
	if h := message.DecodeHandlingMap[f.Handling]; h != "" {
		f.Handling = h
	}
	f.Fields = []*message.Field{
		message.NewCalculatedField(&message.Field{
			Label: "Form Type",
			TableValue: func(*message.Field) string {
				return f.Type.HTML + " v" + f.Type.Version
			},
		}),
		message.NewMessageNumberField(&message.Field{
			Label:    "Origin Message Number",
			Value:    &f.OriginMsgID,
			Presence: message.Required,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Handling",
			Value:    &f.Handling,
			Choices:  message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: message.Required,
		}),
		message.NewTextField(&message.Field{
			Label:    "Subject",
			Value:    &f.Subject,
			Presence: message.Required,
		}),
	}
	var tags = make([]string, 0, len(form.TaggedValues))
	for tag := range form.TaggedValues {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	for _, tag := range tags {
		value := form.TaggedValues[tag]
		f.Fields = append(f.Fields, message.NewTextField(&message.Field{
			Label:   tag,
			Value:   &value,
			PIFOTag: tag,
		}))
	}
	return f
}
