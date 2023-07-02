package ics213

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) (f *ICS213) {
	if idx := strings.Index(body, "form-ics213.html"); idx < 0 {
		return nil
	}
	form := common.DecodePIFO(body)
	if form == nil || form.HTMLIdent != "form-ics213.html" {
		return nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2":
		break
	default:
		return nil
	}
	f = new(ICS213)
	f.PIFOVersion = form.PIFOVersion
	f.FormVersion = form.FormVersion
	f.Date = form.TaggedValues["1a."]
	f.Time = form.TaggedValues["1b."]
	if f.FormVersion < "2.2" {
		f.Severity = form.TaggedValues["4."]
	}
	f.Handling = form.TaggedValues["5."]
	f.TakeAction = form.TaggedValues["6a."]
	f.Reply = form.TaggedValues["6b."]
	f.ReplyBy = form.TaggedValues["6d."]
	if f.FormVersion < "2.2" {
		f.FYI = form.TaggedValues["6c."]
	}
	f.ToICSPosition = form.TaggedValues["7."]
	f.ToLocation = form.TaggedValues["9a."]
	f.ToName = form.TaggedValues["ToName"]
	f.ToTelephone = form.TaggedValues["ToTel"]
	f.FromICSPosition = form.TaggedValues["8."]
	f.FromLocation = form.TaggedValues["9b."]
	f.FromName = form.TaggedValues["FmName"]
	f.FromTelephone = form.TaggedValues["FmTel"]
	f.Subject = form.TaggedValues["10."]
	f.Reference = form.TaggedValues["11."]
	f.Message = form.TaggedValues["12."]
	f.OpRelayRcvd = form.TaggedValues["OpRelayRcvd"]
	f.OpRelaySent = form.TaggedValues["OpRelaySent"]
	f.ReceivedSent = form.TaggedValues["Rec-Sent"]
	f.OpCall = form.TaggedValues["OpCall"]
	f.OpName = form.TaggedValues["OpName"]
	f.TxMethod = form.TaggedValues["Method"]
	f.OtherMethod = form.TaggedValues["Other"]
	f.OpDate = form.TaggedValues["OpDate"]
	f.OpTime = form.TaggedValues["OpTime"]
	if f.FormVersion >= "2.2" {
		f.OriginMsgID = form.TaggedValues["MsgNo"]
		f.DestinationMsgID = form.TaggedValues["3."]
	} else if f.OriginMsgID = form.TaggedValues["2."]; f.OriginMsgID != "" || f.ReceivedSent == "receiver" {
		f.DestinationMsgID = form.TaggedValues["MsgNo"]
	} else {
		f.OriginMsgID = form.TaggedValues["MsgNo"]
		f.DestinationMsgID = form.TaggedValues["3."]
	}
	return f
}
