package xscmsg

import (
	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// UnknownForm holds the details of a form message containing an unknown form.
type UnknownForm struct {
	OriginMsgID  string
	Handling     string
	FormTag      string
	Subject      string
	FormHTML     string
	FormVersion  string
	TaggedValues map[string]string
}

// DecodeUnknownForm decodes the supplied form as an unknown form.  It
// returns the decoded form and strings describing any non-fatal decoding
// problems.
func DecodeUnknownForm(subject string, form *pifo.Form) *UnknownForm {
	f := UnknownForm{
		FormHTML:     form.HTMLIdent,
		FormVersion:  form.FormVersion,
		TaggedValues: form.TaggedValues,
	}
	f.OriginMsgID, _, f.Handling, f.FormTag, f.Subject = xscsubj.Decode(subject)
	return &f
}
