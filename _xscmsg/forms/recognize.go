package xscmsg

import (
	"github.com/rothskeller/packet/xscmsg/forms/pifo"
)

// Decode decodes the supplied message contents.  It returns the decoded message
// contents and strings describing any non-fatal decoding problems.
func Decode(subject, body string) (m any, problems []string) {
	if form := pifo.Decode(body); form != nil {
		if m, problems = DecodeAHFacStat(form); m != nil {
			return m, problems
		}
		if m, problems = DecodeEOC213RR(form); m != nil {
			return m, problems
		}
		if m, problems = DecodeICS213(form); m != nil {
			return m, problems
		}
		if m, problems = DecodeJurisStat(form); m != nil {
			return m, problems
		}
		if m, problems = DecodeRACESMAR(form); m != nil {
			return m, problems
		}
		if m, problems = DecodeSheltStat(form); m != nil {
			return m, problems
		}
		return DecodeUnknownForm(subject, form), nil
	}
	if m = DecodeCheckIn(subject, body); m != nil {
		return m, nil
	}
	if m = DecodeCheckOut(subject, body); m != nil {
		return m, nil
	}
	if m = DecodeDeliveryReceipt(subject, body); m != nil {
		return m, nil
	}
	if m = DecodeReadReceipt(subject, body); m != nil {
		return m, nil
	}
	return DecodePlainText(subject, body), nil
}
