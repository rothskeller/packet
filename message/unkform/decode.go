package unkform

import (
	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) *UnknownForm {
	form := common.DecodePIFO(body)
	if form == nil {
		return nil
	}
	f := UnknownForm{
		FormHTML:     form.HTMLIdent,
		FormVersion:  form.FormVersion,
		TaggedValues: form.TaggedValues,
	}
	f.OriginMsgID, _, f.Handling, f.FormTag, f.Subject = common.DecodeSubject(subject)
	if h := common.DecodeHandlingMap[f.Handling]; h != "" {
		f.Handling = h
	}
	return &f
}
