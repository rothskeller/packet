package plaintext

import "github.com/rothskeller/packet/message/common"

func decode(subject, body string) *PlainText {
	var pt PlainText

	pt.OriginMsgID, _, pt.Handling, _, pt.Subject = common.DecodeSubject(subject)
	if h := common.DecodeHandlingMap[pt.Handling]; h != "" {
		pt.Handling = h
	}
	pt.Body = body
	return &pt
}
