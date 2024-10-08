// Package jurisstat defines the Santa Clara County OA Jurisdiction Status Form
// message type.
package jurisstat

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type20 is the type definition for an OA jurisdiction status form.
var Type20 = message.Type{
	Tag:        "MuniStat",
	HTML:       "form-oa-muni-status.html",
	Version:    "2.0",
	Name:       "OA municipal status form",
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
	return df
}
