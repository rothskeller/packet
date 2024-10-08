// Package ics213 defines the ICS-213 General Message Form message type.
package ics213

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type20 is the type definition for an ICS-213 general message form.
var Type20 = message.Type{
	Tag:        "ICS213",
	HTML:       "form-ics213.html",
	Version:    "2.0",
	Name:       "ICS-213 general message form",
	Article:    "an",
	FieldOrder: Type21.FieldOrder,
}

func init() {
	message.Register(&Type20, decode20, nil)
}

func decode20(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type20.HTML || form.FormVersion != Type20.Version {
		return nil
	}
	var df = make21()
	df.Type = &Type20
	message.DecodeForm(form, df)
	// If we got an OriginMsgID, or we are the receiver, move myMsgID to
	// DestinationMsgID.
	if df.OriginMsgID != "" || df.ReceivedSent == "receiver" {
		df.DestinationMsgID = df.myMsgID
		df.Fields[2].PIFOTag = ""
	} else {
		df.OriginMsgID = df.myMsgID
		df.Fields[0].PIFOTag = ""
	}
	return df
}
