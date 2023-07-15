package checkout

import (
	"regexp"
	"strings"

	"github.com/rothskeller/packet/message/common"
)

var checkOutBodyRE = regexp.MustCompile(`(?im)^Check-Out\s+([A-Z][A-Z0-9]{2,5})\s*,(.*)(?:\n(A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})\s*,(.*))?`)

func decode(subject, body string) *CheckOut {
	var msgid, _, handling, formtag, realsubj = common.DecodeSubject(subject)
	if formtag != "" || !strings.HasPrefix(strings.ToLower(realsubj), "check-out ") {
		return nil
	}
	if h := common.DecodeHandlingMap[handling]; h != "" {
		handling = h
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
