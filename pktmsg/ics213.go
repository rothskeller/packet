package pktmsg

// This file defines TxICS213Form and RxICS213Form.

import (
	"fmt"
	"time"
)

// A TxICS213Form is an outgoing PackItForms-encoded message containing an
// SCCo ICS-213 form.
type TxICS213Form struct {
	TxForm
	DateTime         time.Time
	Severity         MessageSeverity
	TakeAction       string
	Reply            string
	ReplyBy          string
	FYI              bool
	ToICSPosition    string
	ToLocation       string
	ToName           string
	ToTelephone      string
	FromICSPosition  string
	FromLocation     string
	FromName         string
	FromTelephone    string
	Reference        string
	MessageBody      string
	RelayReceived    string
	RelaySent        string
	ReceiverSender   string
	OperatorCallSign string
	OperatorName     string
	OperatorMethod   string
	OperatorDateTime time.Time
}

// Encode returns the encoded subject line and body of the message.
func (i *TxICS213Form) Encode() (subject, body string, err error) {
	if i.DateTime.IsZero() ||
		i.Severity == 0 ||
		i.HandlingOrder == 0 ||
		i.ToICSPosition == "" ||
		i.ToLocation == "" ||
		i.FromICSPosition == "" ||
		i.FromLocation == "" ||
		i.Subject == "" ||
		i.MessageBody == "" ||
		i.ReceiverSender == "" ||
		i.OperatorCallSign == "" ||
		i.OperatorName == "" ||
		i.OperatorMethod == "" ||
		i.OperatorDateTime.IsZero() {
		return "", "", ErrIncomplete
	}
	if (i.TakeAction != "" && i.TakeAction != "Yes" && i.TakeAction != "No") ||
		(i.Reply != "" && i.Reply != "Yes" && i.Reply != "No") ||
		(i.ReceiverSender != "receiver" && i.ReceiverSender != "sender") {
		return "", "", ErrInvalid
	}
	i.FormName = "ICS213"
	i.FormHTML = "form-ics213.html"
	i.FormVersion = "2.1"
	i.SetField("MsgNo", i.MessageNumber)
	i.SetField("1a.", i.DateTime.Format("01/02/2006"))
	i.SetField("4.", i.Severity.String())
	i.SetField("5.", i.HandlingOrder.String())
	i.SetField("6a.", i.TakeAction)
	i.SetField("6b.", i.Reply)
	i.SetField("6d.", i.ReplyBy)
	i.SetField("6c.", boolToChecked(i.FYI))
	i.SetField("1b.", i.DateTime.Format("15:04"))
	i.SetField("7.", i.ToICSPosition)
	i.SetField("8.", i.FromICSPosition)
	i.SetField("9a.", i.ToLocation)
	i.SetField("9b.", i.FromLocation)
	i.SetField("ToName", i.ToName)
	i.SetField("FmName", i.FromName)
	i.SetField("ToTel", i.ToTelephone)
	i.SetField("FmTel", i.FromTelephone)
	i.SetField("10.", i.Subject)
	i.SetField("11.", i.Reference)
	i.SetField("12.", i.MessageBody)
	i.SetField("OpRelayRcvd", i.RelayReceived)
	i.SetField("OpRelaySent", i.RelaySent)
	i.SetField("Rec-Sent", i.ReceiverSender)
	i.SetField("OpCall", i.OperatorCallSign)
	i.SetField("OpName", i.OperatorName)
	switch i.OperatorMethod {
	case "Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio":
		i.SetField("Method", i.OperatorMethod)
		i.SetField("Other", "")
	default:
		i.SetField("Method", "Other")
		i.SetField("Other", i.OperatorMethod)
	}
	i.SetField("OpDate", i.OperatorDateTime.Format("01/02/2006"))
	i.SetField("OpTime", i.OperatorDateTime.Format("15:04"))
	return i.TxForm.Encode()
}

//------------------------------------------------------------------------------

// An RxICS213Form is a received PackItForms-encoded message containing an SCCo
// ICS-213 form.
type RxICS213Form struct {
	RxForm
	MessageNumber    string
	Date             string
	Time             string
	DateTime         time.Time
	Severity         MessageSeverity
	HandlingOrder    HandlingOrder
	TakeAction       string
	Reply            string
	ReplyBy          string
	FYI              bool
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
	MessageBody      string
	RelayReceived    string
	RelaySent        string
	ReceiverSender   string
	OperatorCallSign string
	OperatorName     string
	OperatorMethod   string
	OperatorDate     string
	OperatorTime     string
	OperatorDateTime time.Time
}

// parseRxICS213Form examines an RxForm to see if it contains an ICS-213 form,
// and if so, wraps it in an RxICS213Form and returns it.  If it is not, it
// returns nil.
func parseRxICS213Form(f *RxForm) *RxICS213Form {
	var i RxICS213Form

	if f.FormHTML != "form-ics213.html" {
		return nil
	}
	i.RxForm = *f
	i.MessageNumber = i.Fields["MsgNo"]
	i.Date = i.Fields["1a."]
	i.Time = i.Fields["1b."]
	i.DateTime = dateTimeParse(i.Date, i.Time)
	i.Severity, _ = ParseSeverity(i.Fields["4."])
	i.HandlingOrder, _ = ParseHandlingOrder(i.Fields["5."])
	i.TakeAction = i.Fields["6a."]
	i.Reply = i.Fields["6b."]
	i.ReplyBy = i.Fields["6d."]
	i.FYI = i.Fields["6c."] != ""
	i.ToICSPosition = i.Fields["7."]
	i.ToLocation = i.Fields["9a."]
	i.ToName = i.Fields["ToName"]
	i.ToTelephone = i.Fields["ToTel"]
	i.FromICSPosition = i.Fields["8."]
	i.FromLocation = i.Fields["9b."]
	i.FromName = i.Fields["FmName"]
	i.FromTelephone = i.Fields["FmTel"]
	i.Subject = i.Fields["10."]
	i.Reference = i.Fields["11."]
	i.MessageBody = i.Fields["12."]
	i.RelayReceived = i.Fields["OpRelayRcvd"]
	i.RelaySent = i.Fields["OpRelaySent"]
	i.ReceiverSender = i.Fields["Rec-Sent"]
	i.OperatorCallSign = i.Fields["OpCall"]
	i.OperatorName = i.Fields["OpName"]
	if i.OperatorMethod = i.Fields["Method"]; i.OperatorMethod == "Other" {
		i.OperatorMethod = i.Fields["Other"]
	}
	i.OperatorDate = i.Fields["OpDate"]
	i.OperatorTime = i.Fields["OpTime"]
	i.OperatorDateTime = dateTimeParse(i.OperatorDate, i.OperatorTime)
	return &i
}

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (i *RxICS213Form) Valid() bool {
	return (i.MessageNumber != "" &&
		i.Date != "" &&
		i.Time != "" &&
		i.Severity != 0 &&
		i.HandlingOrder != 0 &&
		(i.TakeAction == "" || i.TakeAction == "Yes" || i.TakeAction == "No") &&
		(i.Reply == "" || i.Reply == "Yes" || i.Reply == "No") &&
		i.ToICSPosition != "" &&
		i.ToLocation != "" &&
		i.FromICSPosition != "" &&
		i.FromLocation != "" &&
		i.Subject != "" &&
		i.MessageBody != "" &&
		(i.ReceiverSender == "receiver" || i.ReceiverSender == "sender") &&
		i.OperatorCallSign != "" &&
		i.OperatorName != "" &&
		i.OperatorMethod != "" &&
		i.OperatorDate != "" &&
		i.OperatorTime != "")
}

// EncodeSubjectLine returns what the subject line should be based on the
// received form contents.
func (i *RxICS213Form) EncodeSubjectLine() string {
	return fmt.Sprintf("%s_%s_ICS213_%s", i.MessageNumber, i.HandlingOrder.Code(), i.Subject)
}

// TypeCode returns the machine-readable code for the message type.
func (*RxICS213Form) TypeCode() string { return "ICS213" }

// TypeName returns the human-reading name of the message type.
func (*RxICS213Form) TypeName() string { return "ICS-213 form" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxICS213Form) TypeArticle() string { return "an" }
