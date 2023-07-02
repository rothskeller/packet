package ics213

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (f *ICS213) EncodeSubject() string {
	return common.EncodeSubject(f.OriginMsgID, f.Handling, Type.Tag, f.Subject)
}

// EncodeBody encodes the message body.
func (f *ICS213) EncodeBody() string {
	var (
		sb  strings.Builder
		enc *common.PIFOEncoder
	)
	if f.FormVersion == "" {
		f.FormVersion = "2.2"
	}
	enc = common.NewPIFOEncoder(&sb, "form-ics213.html", f.FormVersion)
	if f.FormVersion < "2.2" && f.ReceivedSent == "receiver" {
		enc.Write("2.", f.OriginMsgID)
		enc.Write("MsgNo", f.DestinationMsgID)
	} else {
		enc.Write("MsgNo", f.OriginMsgID)
		enc.Write("3.", f.DestinationMsgID)
	}
	enc.Write("1a.", f.Date)
	enc.Write("1b.", f.Time)
	if f.FormVersion < "2.2" {
		enc.Write("4.", f.Severity)
	}
	enc.Write("5.", f.Handling)
	enc.Write("6a.", f.TakeAction)
	enc.Write("6b.", f.Reply)
	if f.FormVersion < "2.2" {
		enc.Write("6c.", f.FYI)
	}
	enc.Write("6d.", f.ReplyBy)
	enc.Write("7.", f.ToICSPosition)
	enc.Write("8.", f.FromICSPosition)
	enc.Write("9a.", f.ToLocation)
	enc.Write("9b.", f.FromLocation)
	enc.Write("ToName", f.ToName)
	enc.Write("FmName", f.FromName)
	enc.Write("ToTel", f.ToTelephone)
	enc.Write("FmTel", f.FromTelephone)
	enc.Write("10.", f.Subject)
	enc.Write("11.", f.Reference)
	enc.Write("12.", f.Message)
	enc.Write("OpRelayRcvd", f.OpRelayRcvd)
	enc.Write("OpRelaySent", f.OpRelaySent)
	enc.Write("Rec-Sent", f.ReceivedSent)
	enc.Write("OpCall", f.OpCall)
	enc.Write("OpName", f.OpName)
	enc.Write("Method", f.TxMethod)
	enc.Write("Other", f.OtherMethod)
	enc.Write("OpDate", f.OpDate)
	enc.Write("OpTime", f.OpTime)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}
