// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type23 is the type definition for an allied health facility status form.
var Type23 = message.Type{
	Tag:        "AHFacStat",
	HTML:       "form-allied-health-facility-status.html",
	Version:    "2.3",
	Name:       "allied health facility status form",
	Article:    "an",
	FieldOrder: Type26.FieldOrder,
}

func init() {
	message.Register(&Type23, decode23, nil)
}

func decode23(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type23.HTML || form.FormVersion != Type23.Version {
		return nil
	}
	var df = make24()
	df.Type = &Type23
	message.DecodeForm(form, df)
	return df
}
