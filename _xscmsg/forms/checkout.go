package xscmsg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// CheckOut holds the details of an XSC-standard check-out message.
type CheckOut struct {
	OriginMsgID         string
	Handling            string
	TacticalCallSign    string
	TacticalStationName string
	OperatorCallSign    string
	OperatorName        string
}

var checkOutBodyRE = regexp.MustCompile(`(?im)^Check-Out\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n(A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})\s*,(.*))?`)

// DecodeCheckOut decodes the supplied message contents if they are an
// XSC-standard check-out message.  It returns nil if the message contents are
// not a well-formed check-out message.
func DecodeCheckOut(subject, body string) *CheckOut {
	var msgid, _, handling, formtag, realsubj = xscsubj.Decode(subject)
	if formtag != "" || !strings.HasPrefix(strings.ToLower(realsubj), "check-out ") {
		return nil
	}
	if match := checkOutBodyRE.FindStringSubmatch(body); match != nil {
		var ci = CheckOut{OriginMsgID: msgid, Handling: handling}
		if match[3] != "" {
			ci.TacticalCallSign = match[1]
			ci.TacticalStationName = strings.TrimSpace(match[2])
			ci.OperatorCallSign = match[3]
			ci.OperatorName = strings.TrimSpace(match[4])
		} else {
			ci.OperatorCallSign = match[1]
			ci.OperatorName = strings.TrimSpace(match[2])
		}
		return &ci
	}
	return nil
}

// Encode encodes the message contents.
func (m *CheckOut) Encode() (subject, body string) {
	if m.TacticalCallSign != "" {
		subject = xscsubj.Encode(m.OriginMsgID, m.Handling, "", fmt.Sprintf("Check-Out %s, %s\n",
			m.TacticalCallSign, m.TacticalStationName))
		body = fmt.Sprintf("Check-Out %s, %s\n%s, %s\n",
			m.TacticalCallSign, m.TacticalStationName, m.OperatorCallSign, m.OperatorName)
	} else {
		subject = xscsubj.Encode(m.OriginMsgID, m.Handling, "", fmt.Sprintf("Check-Out %s, %s\n",
			m.OperatorCallSign, m.OperatorName))
		body = fmt.Sprintf("Check-Out %s, %s\n", m.OperatorCallSign, m.OperatorName)
	}
	return
}
