// Package ics213 defines the ICS-213 General Message Form message type.
package ics213

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an ICS-213 general message form.
var Type = message.Type{
	Tag:     "ICS213",
	Name:    "ICS-213 general message form",
	Article: "an",
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*basemsg.FormVersion{
	{HTML: "form-ics213.html", Version: "2.2", Tag: "ICS213", FieldOrder: fieldOrder},
	{HTML: "form-ics213.html", Version: "2.1", Tag: "ICS213", FieldOrder: fieldOrder},
	{HTML: "form-ics213.html", Version: "2.0", Tag: "ICS213", FieldOrder: fieldOrder},
}
var fieldOrder = []string{
	"2.", "MsgNo", "3.", "1a.", "1b.", "4.", "5.", "6a.", "6b.", "6c.", "6d.", "7.", "8.", "9a.", "9b.", "ToName", "FmName",
	"ToTel", "FmTel", "10.", "11.", "12.", "OpRelayRcvd", "OpRelaySent", "Rec-Sent", "OpCall", "OpName", "Method", "Other",
	"OpDate", "OpTime",
}

// ICS213 holds an ICS-213 general message form.
type ICS213 struct {
	basemsg.BaseMessage
	OriginMsgID      string
	myMsgID          string // removed in v2.2
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
}

func New() (f *ICS213) {
	f = create(versions[0]).(*ICS213)
	f.Date = time.Now().Format("01/02/2006")
	f.ReceivedSent = "sender"
	f.TxMethod = "Other"
	f.OtherMethod = "Packet"
	return f
}

var pdfBase []byte

func create(version *basemsg.FormVersion) message.Message {
	const fieldCount = 34
	var f = ICS213{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		PDFBase:     pdfBase,
		Form:        version,
	}}
	f.BaseMessage.FOriginMsgID = &f.OriginMsgID
	f.BaseMessage.FHandling = &f.Handling
	f.BaseMessage.FToICSPosition = &f.ToICSPosition
	f.BaseMessage.FToLocation = &f.ToLocation
	f.BaseMessage.FSubject = &f.Subject
	f.BaseMessage.FBody = &f.Message
	f.BaseMessage.FOpCall = &f.OpCall
	f.BaseMessage.FOpName = &f.OpName
	f.Fields = make([]*basemsg.Field, 0, fieldCount)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "Origin Message Number",
				Value:     &f.OriginMsgID,
				Presence:  basemsg.Required,
				PIFOTag:   "2.", // may be changed to MsgNo by decode()
				PDFMap:    basemsg.PDFName("Origin Msg #"),
				EditWidth: 9,
				EditHelp:  `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
			},
			&basemsg.Field{
				Label:      "My Message Number",
				Value:      &f.myMsgID,
				PIFOTag:    "MsgNo", // will be removed after decode()
				TableValue: basemsg.OmitFromTable,
			},
		)
	} else {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "Origin Message Number",
				Value:     &f.OriginMsgID,
				Presence:  basemsg.Required,
				PIFOTag:   "MsgNo",
				PDFMap:    basemsg.PDFName("Origin Msg #"),
				EditWidth: 9,
				EditHelp:  `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
			},
		)
	}
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Destination Message Number",
			Value:   &f.DestinationMsgID,
			PIFOTag: "3.", // may be changed to MsgNo after decode()
			PDFMap:  basemsg.PDFName("Destination Msg#"),
		},
		&basemsg.Field{
			Label:      "Date",
			Value:      &f.Date,
			Presence:   basemsg.Required,
			PIFOTag:    "1a.",
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareDate,
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("FormDate"),
		},
		&basemsg.Field{
			Label:      "Time",
			Value:      &f.Time,
			Presence:   basemsg.Required,
			PIFOTag:    "1b.",
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareTime,
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("FormTime"),
		},
		&basemsg.Field{
			Label:    "Date/Time",
			Presence: basemsg.Required,
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.Date, f.Time, " ")
			},
			EditWidth: 16,
			EditHelp:  "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.Date, f.Time)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.Date, &f.Time, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.Date, f.Time)
			},
		},
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "Situation Severity",
				Value:     &f.Severity,
				Choices:   basemsg.Choices{"OTHER", "URGENT", "EMERGENCY"},
				Presence:  basemsg.Required,
				PIFOTag:   "4.",
				PIFOValid: basemsg.ValidRestricted,
				Compare:   common.CompareExact,
				EditWidth: 9,
				EditHelp:  `This is the severity of the situation that caused the message to be sent.  It is required.`,
			},
		)
	}
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:     "Handling",
			Value:     &f.Handling,
			Choices:   basemsg.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence:  basemsg.Required,
			PIFOTag:   "5.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Immediate", "", "Off", "PRIORITY", "1", "ROUTINE", "2", "IMMEDIATE", "3"},
			EditWidth: 9,
			EditHelp:  `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		},
		&basemsg.Field{
			Label:     "Take Action",
			Value:     &f.TakeAction,
			Choices:   basemsg.Choices{"Yes", "No"},
			PIFOTag:   "6a.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"TakeAction", "", "Off", "Yes", "1", "No", "2"},
			EditWidth: 3,
			EditHelp:  `This indicates whether the sender expects the recipient to take action based on this message.`,
		},
		&basemsg.Field{
			Label:     "Reply",
			Value:     &f.Reply,
			Choices:   basemsg.Choices{"Yes", "No"},
			PIFOTag:   "6b.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			TableValue: func(*basemsg.Field) string {
				if f.Reply == "Yes" && f.ReplyBy != "" {
					return "Yes, by " + f.ReplyBy
				}
				return f.Reply
			},
			PDFMap:    basemsg.PDFNameMap{"Reply", "", "Off", "Yes", "Reply-Yes", "No", "Reply-No"},
			EditWidth: 3,
			EditHelp:  `This indicates whether the sender expects the recipient to reply to this message.`,
		},
		&basemsg.Field{
			Label: "Reply By",
			Value: &f.ReplyBy,
			Presence: func() (basemsg.Presence, string) {
				if f.Reply != "Yes" {
					return basemsg.PresenceNotAllowed, `unless "Reply" is "Yes"`
				}
				return basemsg.PresenceOptional, ""
			},
			PIFOTag:    "6d.",
			Compare:    common.CompareExact,
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("Reply_2"),
			EditWidth:  5,
			EditHelp:   `This is the time by which the sender expects a reply from the recipient.  It can be set only if "Reply" is "Yes".`,
		},
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "For Your Information",
				Value:     &f.FYI,
				Choices:   basemsg.Choices{"checked"},
				PIFOTag:   "6c.",
				PIFOValid: basemsg.ValidRestricted,
				Compare:   common.CompareExact,
				EditWidth: 7,
				EditHelp:  `This indicates whether the message is merely informational.`,
			},
		)
	}
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:     "To ICS Position",
			Value:     &f.ToICSPosition,
			Presence:  basemsg.Required,
			PIFOTag:   "7.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("TO ICS Position"),
			EditWidth: 40,
			EditHelp:  `This is the ICS position to which the message is addressed.  It is required.`,
		},
		&basemsg.Field{
			Label:     "To Location",
			Value:     &f.ToLocation,
			Presence:  basemsg.Required,
			PIFOTag:   "9a.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("TO ICS Locatoin"),
			EditWidth: 40,
			EditHelp:  `This is the location of the recipient ICS position.  It is required.`,
		},
		&basemsg.Field{
			Label:     "To Name",
			Value:     &f.ToName,
			PIFOTag:   "ToName",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("TO ICS Name"),
			EditWidth: 40,
			EditHelp:  `This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.`,
		},
		&basemsg.Field{
			Label:     "To Telephone #",
			Value:     &f.ToTelephone,
			PIFOTag:   "ToTel",
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("TO ICS Telephone"),
			EditWidth: 40,
			EditHelp:  `This is the telephone number for the receipient.  It is optional and rarely provided.`,
		},
		&basemsg.Field{
			Label:     "From ICS Position",
			Value:     &f.FromICSPosition,
			Presence:  basemsg.Required,
			PIFOTag:   "8.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("From ICS Position"),
			EditWidth: 40,
			EditHelp:  `This is the ICS position of the message author.  It is required.`,
		},
		&basemsg.Field{
			Label:     "From Location",
			Value:     &f.FromLocation,
			Presence:  basemsg.Required,
			PIFOTag:   "9b.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("From ICS Location"),
			EditWidth: 40,
			EditHelp:  `This is the location of the message author.  It is required.`,
		},
		&basemsg.Field{
			Label:     "From Name",
			Value:     &f.FromName,
			PIFOTag:   "FmName",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("From ICS Name"),
			EditWidth: 40,
			EditHelp:  `This is the name of the message author.  It is optional and rarely provided.`,
		},
		&basemsg.Field{
			Label:     "From Telephone #",
			Value:     &f.FromTelephone,
			PIFOTag:   "FmTel",
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("From ICS Telephone"),
			EditWidth: 40,
			EditHelp:  `This is the telephone number for the message author.  It is optional and rarely provided.`,
		},
		&basemsg.Field{
			Label:     "Subject",
			Value:     &f.Subject,
			Presence:  basemsg.Required,
			PIFOTag:   "10.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Subject"),
			EditWidth: 80,
			EditHelp:  `This is the subject of the message.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Reference",
			Value:     &f.Reference,
			PIFOTag:   "11.",
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("Reference"),
			EditWidth: 55,
			EditHelp:  `This is the origin message number of the previous message to which this message is responding.`,
		},
		&basemsg.Field{
			Label:     "Message",
			Value:     &f.Message,
			Presence:  basemsg.Required,
			PIFOTag:   "12.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Message"),
			EditWidth: 80,
			Multiline: true,
			EditHelp:  `This is the text of the message.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Operator: Relay Received",
			Value:     &f.OpRelayRcvd,
			PIFOTag:   "OpRelayRcvd",
			PDFMap:    basemsg.PDFName("Relay Received"),
			EditWidth: 36,
			EditHelp:  `This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.`,
		},
		&basemsg.Field{
			Label:     "Operator: Relay Sent",
			Value:     &f.OpRelaySent,
			PIFOTag:   "OpRelaySent",
			PDFMap:    basemsg.PDFName("Relay Sent"),
			EditWidth: 36,
			EditHelp:  `This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.`,
		},
		&basemsg.Field{
			Label:     "Operator: Receiver/Sender",
			Value:     &f.ReceivedSent,
			Choices:   basemsg.Choices{"sender", "receiver"},
			Presence:  basemsg.Required,
			PIFOTag:   "Rec-Sent",
			PIFOValid: basemsg.ValidRestricted,
			PDFMap:    basemsg.PDFNameMap{"How: Received", "", "Off", "sender", "0", "receiver", "1"},
		},
		&basemsg.Field{
			Label:      "Operator: Call Sign",
			Value:      &f.OpCall,
			Presence:   basemsg.Required,
			PIFOTag:    "OpCall",
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("Operation Call Sign"),
		},
		&basemsg.Field{
			Label:      "Operator: Name",
			Value:      &f.OpName,
			Presence:   basemsg.Required,
			PIFOTag:    "OpName",
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("Relay Received_2"),
		},
		&basemsg.Field{
			Label: "Operator",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.OpCall, f.OpName, " ")
			},
		},
		&basemsg.Field{
			Label:    "Operator: Tx Method",
			Value:    &f.TxMethod,
			Choices:  basemsg.Choices{"Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio", "Other"},
			Presence: basemsg.Required,
			PIFOTag:  "Method",
			TableValue: func(*basemsg.Field) string {
				if f.TxMethod == "Other" && f.OtherMethod != "" {
					return f.OtherMethod
				}
				return f.TxMethod
			},
			PDFMap: basemsg.PDFNameMap{"Telephone", "", "Off", "Telephone", "1", "Dispatch Center", "2", "EOC Radio", "3", "FAX", "4", "Courier", "5", "Amateur Radio", "6", "Other", "7"},
		},
		&basemsg.Field{
			Label: "Operator: Other Method",
			Value: &f.OtherMethod,
			Presence: func() (basemsg.Presence, string) {
				if f.TxMethod != "Other" {
					return basemsg.PresenceNotAllowed, `unless "Operator: Tx Method" is "Other"`
				}
				return basemsg.PresenceOptional, ""
			},
			PIFOTag:    "Other",
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("OtherText"),
		},
		&basemsg.Field{
			Label:      "Operator: Date",
			Value:      &f.OpDate,
			Presence:   basemsg.Required,
			PIFOTag:    "OpDate",
			PIFOValid:  basemsg.ValidDate,
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("OperatorDate"),
		},
		&basemsg.Field{
			Label:      "Operator: Time",
			Value:      &f.OpTime,
			Presence:   basemsg.Required,
			PIFOTag:    "OpTime",
			PIFOValid:  basemsg.ValidTime,
			TableValue: basemsg.OmitFromTable,
			PDFMap:     basemsg.PDFName("OperatorTime"),
		},
		&basemsg.Field{
			Label:    "Operator: Date/Time",
			Presence: basemsg.Required,
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.OpDate, f.OpTime, " ")
			},
		},
	)
	if len(f.Fields) > fieldCount {
		panic("update ICS213 fieldCount")
	}
	return &f
}

func decode(subject, body string) (f *ICS213) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-ics213.html") {
		return nil
	}
	if f = basemsg.Decode(body, versions, create).(*ICS213); f == nil {
		return nil
	}
	// We need to fix up the origin/my/destination message numbers in the
	// old version of this form.
	if f.Form.Version >= "2.2" {
		return f
	}
	// If we got an OriginMsgID, or we are the receiver, move myMsgID to
	// DestinationMsgID.
	if f.OriginMsgID != "" || f.ReceivedSent == "receiver" {
		f.DestinationMsgID = f.myMsgID
		f.Fields[2].PIFOTag = "MsgNo"
	} else {
		f.OriginMsgID = f.myMsgID
		f.Fields[0].PIFOTag = "MsgNo"
	}
	f.Fields[1].PIFOTag = ""
	return f
}
