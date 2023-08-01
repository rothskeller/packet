// Package plaintext handles plain text messages.
package plaintext

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for a plain text message.
var Type = message.Type{
	Tag:     "plain",
	Name:    "plain text message",
	Article: "a",
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// PlainText holds the details of a plain text message.
type PlainText struct {
	basemsg.BaseMessage
	OriginMsgID string
	Handling    string
	Subject     string
	Body        string
}

// New creates a new plain text message.
func New() (m *PlainText) {
	m = &PlainText{BaseMessage: basemsg.BaseMessage{MessageType: &Type}}
	m.BaseMessage.FOriginMsgID = &m.OriginMsgID
	m.BaseMessage.FHandling = &m.Handling
	m.BaseMessage.FSubject = &m.Subject
	m.BaseMessage.FBody = &m.Body
	m.Fields = []*basemsg.Field{
		{
			Label:     "Origin Message Number",
			Value:     &m.OriginMsgID,
			Presence:  basemsg.Required,
			Compare:   common.CompareExact,
			EditWidth: 9,
			EditHelp:  `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
			EditHint:  "XXX-###P",
			EditApply: basemsg.ApplyMessageNumber,
			EditValid: basemsg.ValidMessageNumber,
		},
		{
			Label:     "Handling",
			Value:     &m.Handling,
			Choices:   basemsg.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence:  basemsg.Required,
			Compare:   common.CompareExact,
			EditWidth: 9,
			EditHelp:  `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
			EditValid: basemsg.ValidRestricted,
		},
		{
			Label:     "Subject",
			Value:     &m.Subject,
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			EditWidth: 80,
			EditHelp:  `This is the subject of the message.  It is required.`,
		},
		{
			Label:     "Message",
			Value:     &m.Body,
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			EditWidth: 80,
			Multiline: true,
			EditHelp:  `This is the body of the message.  It is required.`,
		},
	}
	return m
}

// This function is called to find out whether an incoming message matches this
// type.  It should return the decoded message if it belongs to this type, or
// nil if it doesn't.
func decode(subject, body string) (f *PlainText) {
	f = New()
	f.OriginMsgID, _, f.Handling, _, f.Subject = common.DecodeSubject(subject)
	f.Body = body
	return f
}

func (m *PlainText) EncodeBody() string { return m.Body }
