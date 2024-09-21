// Package eoc213rr defines the Santa Clara County EOC-213RR Resource Request
// Form message type.
package eoc213rr

import (
	"github.com/rothskeller/packet/message"
)

// Type20 is the type definition for an EOC-213RR resource request form.
var Type20 = message.Type{
	Tag:        "EOC213RR",
	HTML:       "form-scco-eoc-213rr.html",
	Version:    "2.0",
	Name:       "EOC-213RR resource request form",
	Article:    "an",
	FieldOrder: Type23.FieldOrder,
}

func init() {
	message.Register(&Type20, decode20, nil)
}

func decode20(_, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type20.HTML || form.FormVersion != Type20.Version {
		return nil
	}
	var df = make23()
	df.Type = &Type20
	message.DecodeForm(form, df)
	return df
}
