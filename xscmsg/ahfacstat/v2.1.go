// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type21 is the type definition for an allied health facility status form.
var Type21 = message.Type{
	Tag:        "AHFacStat",
	HTML:       "form-allied-health-facility-status.html",
	Version:    "2.1",
	Name:       "allied health facility status form",
	Article:    "an",
	FieldOrder: Type26.FieldOrder,
}

func init() {
	message.Register(&Type21, decode21, nil)
}

func decode21(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type21.HTML || form.FormVersion != Type21.Version {
		return nil
	}
	var df = make24()
	df.Type = &Type21
	message.DecodeForm(form, df)
	return df
}
