package plaintext

import "github.com/rothskeller/packet/message/common"

func decode(subject, body string) *PlainText {
	var pt PlainText

	pt.OriginMsgID, _, pt.Handling, _, pt.Subject = common.DecodeSubject(subject)
	pt.Body = body
	return &pt
}
