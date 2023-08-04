// Package unkform handles PackItForms messages with an unrecognized form type.
package unkform

import (
	"sort"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an unrecognized form message.
var Type = message.Type{
	Tag:     "UNKNOWN",
	Name:    "unrecognized form message",
	Article: "an",
}

func init() {
	Type.Decode = decode
}

// UnknownForm holds the details of an unrecognized form message.
type UnknownForm struct {
	basemsg.BaseMessage
	OriginMsgID  string
	Handling     string
	Subject      string
	TaggedValues map[string]string
}

// This function is called to find out whether an incoming message matches this
// type.  It should return the decoded message if it belongs to this type, or
// nil if it doesn't.
func decode(subject, body string) (f *UnknownForm) {
	form := common.DecodePIFO(body)
	if form == nil {
		return nil
	}
	f = &UnknownForm{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		Form: &basemsg.FormVersion{
			HTML:    form.HTMLIdent,
			Version: form.FormVersion,
		},
	}}
	f.OriginMsgID, _, f.Handling, f.Form.Tag, f.Subject = common.DecodeSubject(subject)
	if h := common.DecodeHandlingMap[f.Handling]; h != "" {
		f.Handling = h
	}
	f.Fields = []*basemsg.Field{
		basemsg.NewCalculatedField(&basemsg.Field{
			Label: "Form Type",
			TableValue: func(*basemsg.Field) string {
				return f.Form.HTML + " v" + f.Form.Version
			},
		}),
		basemsg.NewMessageNumberField(&basemsg.Field{
			Label:    "Origin Message Number",
			Value:    &f.OriginMsgID,
			Presence: basemsg.Required,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Handling",
			Value:    &f.Handling,
			Choices:  basemsg.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: basemsg.Required,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:    "Subject",
			Value:    &f.Subject,
			Presence: basemsg.Required,
		}),
	}
	var tags = make([]string, 0, len(form.TaggedValues))
	for tag := range form.TaggedValues {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	for _, tag := range tags {
		value := form.TaggedValues[tag]
		f.Fields = append(f.Fields, basemsg.NewTextField(&basemsg.Field{
			Label:   tag,
			Value:   &value,
			PIFOTag: tag,
		}))
	}
	return f
}
