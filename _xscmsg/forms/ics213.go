package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// ICS213 form metadata:
const (
	ICS213Tag     = "ICS213"
	ICS213HTML    = "form-ics213.html"
	ICS213Version = "2.2"
)

// ICS213 holds an ICS-213 general message form.
type ICS213 struct {
	FormVersion      string
	OriginMsgID      string
	DestinationMsgID string
	Date             string
	Time             string
	Severity         string // removed in v2.2
	Handling         string
	TakeAction       string
	Reply            string
	ReplyBy          string
	FYI              string // removed in v2.2
	ToICSPosition    string
	ToLocation       string
	ToName           string
	ToTelephone      string
	FromICSPosition  string
	FromLocation     string
	FromName         string
	FromTelephone    string
	Subject          string
	Reference        string
	Message          string
	OpRelayRcvd      string
	OpRelaySent      string
	ReceivedSent     string
	OpCall           string
	OpName           string
	TxMethod         string
	OtherMethod      string
	OpDate           string
	OpTime           string
	UnknownFields    map[string]string
}

// DecodeICS213 decodes the supplied form if it is an ICS213 form.  It returns
// the decoded form and strings describing any non-fatal decoding problems.  It
// returns nil, nil if the form is not an ICS213 form or has an unknown version.
func DecodeICS213(form *pifo.Form) (f *ICS213, problems []string) {
	if form.HTMLIdent != ICS213HTML {
		return nil, nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2":
		break
	default:
		return nil, nil
	}
	f = new(ICS213)
	f.FormVersion = form.FormVersion
	f.Date = PullTag(form.TaggedValues, "1a.")
	f.Time = PullTag(form.TaggedValues, "1b.")
	if f.FormVersion != "2.2" {
		f.Severity = PullTag(form.TaggedValues, "4.")
	}
	f.Handling = PullTag(form.TaggedValues, "5.")
	f.TakeAction = PullTag(form.TaggedValues, "6a.")
	f.Reply = PullTag(form.TaggedValues, "6b.")
	f.ReplyBy = PullTag(form.TaggedValues, "6d.")
	if f.FormVersion != "2.2" {
		f.FYI = PullTag(form.TaggedValues, "6c.")
	}
	f.ToICSPosition = PullTag(form.TaggedValues, "7.")
	f.ToLocation = PullTag(form.TaggedValues, "9a.")
	f.ToName = PullTag(form.TaggedValues, "ToName")
	f.ToTelephone = PullTag(form.TaggedValues, "ToTel")
	f.FromICSPosition = PullTag(form.TaggedValues, "8.")
	f.FromLocation = PullTag(form.TaggedValues, "9b.")
	f.FromName = PullTag(form.TaggedValues, "FmName")
	f.FromTelephone = PullTag(form.TaggedValues, "FmTel")
	f.Subject = PullTag(form.TaggedValues, "10.")
	f.Reference = PullTag(form.TaggedValues, "11.")
	f.Message = PullTag(form.TaggedValues, "12.")
	f.OpRelayRcvd = PullTag(form.TaggedValues, "OpRelayRcvd")
	f.OpRelaySent = PullTag(form.TaggedValues, "OpRelaySent")
	f.ReceivedSent = PullTag(form.TaggedValues, "Rec-Sent")
	f.OpCall = PullTag(form.TaggedValues, "OpCall")
	f.OpName = PullTag(form.TaggedValues, "OpName")
	f.TxMethod = PullTag(form.TaggedValues, "Method")
	f.OtherMethod = PullTag(form.TaggedValues, "Other")
	f.OpDate = PullTag(form.TaggedValues, "OpDate")
	f.OpTime = PullTag(form.TaggedValues, "OpTime")
	if f.FormVersion == "2.2" {
		f.OriginMsgID = PullTag(form.TaggedValues, "MsgNo")
		f.DestinationMsgID = PullTag(form.TaggedValues, "3.")
	} else if f.OriginMsgID = PullTag(form.TaggedValues, "2."); f.OriginMsgID != "" || f.ReceivedSent == "receiver" {
		f.DestinationMsgID = PullTag(form.TaggedValues, "MsgNo")
	} else {
		f.OriginMsgID = PullTag(form.TaggedValues, "MsgNo")
		f.DestinationMsgID = PullTag(form.TaggedValues, "3.")
	}
	return f, LeftoverTagProblems(ICS213Tag, form.FormVersion, form.TaggedValues)
}

// Encode encodes the message contents.
func (f *ICS213) Encode() (subject, body string) {
	var (
		sb  strings.Builder
		enc *pifo.Encoder
	)
	subject = xscsubj.Encode(f.OriginMsgID, f.Handling, ICS213Tag, f.Subject)
	if f.FormVersion == "" {
		f.FormVersion = "2.2"
	}
	enc = pifo.NewEncoder(&sb, ICS213HTML, f.FormVersion)
	if f.FormVersion != "2.2" && f.ReceivedSent == "receiver" {
		enc.Write("2.", f.OriginMsgID)
		enc.Write("MsgNo", f.DestinationMsgID)
	} else {
		enc.Write("MsgNo", f.OriginMsgID)
		enc.Write("3.", f.DestinationMsgID)
	}
	enc.Write("1a.", f.Date)
	enc.Write("1b.", f.Time)
	if f.FormVersion != "2.2" {
		enc.Write("4.", f.Severity)
	}
	enc.Write("5.", f.Handling)
	enc.Write("6a.", f.TakeAction)
	enc.Write("6b.", f.Reply)
	if f.FormVersion != "2.2" {
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
	for tag, value := range f.UnknownFields {
		enc.Write(tag, value)
	}
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return subject, sb.String()
}
