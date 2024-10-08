// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type20 is the type definition for an allied health facility status form.
var Type20 = message.Type{
	Tag:        "AHFacStat",
	HTML:       "form-allied-health-facility-status.html",
	Version:    "2.0",
	Name:       "allied health facility status form",
	Article:    "an",
	FieldOrder: Type26.FieldOrder,
}

func init() {
	message.Register(&Type20, decode20, nil)
}

func decode20(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type20.HTML || form.FormVersion != Type20.Version {
		return nil
	}
	var df = make24()
	df.Type = &Type20
	message.DecodeForm(form, df)
	return df
}
