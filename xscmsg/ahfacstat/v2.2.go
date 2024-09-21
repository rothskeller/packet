// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"github.com/rothskeller/packet/message"
)

// Type22 is the type definition for an allied health facility status form.
var Type22 = message.Type{
	Tag:        "AHFacStat",
	HTML:       "form-allied-health-facility-status.html",
	Version:    "2.2",
	Name:       "allied health facility status form",
	Article:    "an",
	FieldOrder: Type26.FieldOrder,
}

func init() {
	message.Register(&Type22, decode22, nil)
}

func decode22(_, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type22.HTML || form.FormVersion != Type22.Version {
		return nil
	}
	var df = make24()
	df.Type = &Type22
	message.DecodeForm(form, df)
	return df
}
