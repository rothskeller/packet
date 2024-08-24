// Package ics213 defines the ICS-213 General Message Form message type.
package ics213

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
)

// Type is the type definition for an ICS-213 general message form.
var Type = message.Type{
	Tag:         "ICS213",
	Name:        "ICS-213 general message form",
	Article:     "an",
	PDFRenderV2: true,
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*message.FormVersion{
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
	message.BaseMessage
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

func (f *ICS213) SetOperator(opcall, opname string, received bool) {
	if received {
		f.ReceivedSent = "receiver"
	} else {
		f.ReceivedSent = "sender"
	}
	f.TxMethod, f.OtherMethod = "Other", "Packet"
	f.BaseMessage.SetOperator(opcall, opname, received)
}

func create(version *message.FormVersion) message.Message {
	const fieldCount = 34
	var f = ICS213{BaseMessage: message.BaseMessage{
		Type: &Type,
		Form: version,
	}}
	f.BaseMessage.FOriginMsgID = &f.OriginMsgID
	f.BaseMessage.FDestinationMsgID = &f.DestinationMsgID
	f.BaseMessage.FMessageDate = &f.Date
	f.BaseMessage.FMessageTime = &f.Time
	f.BaseMessage.FHandling = &f.Handling
	f.BaseMessage.FToICSPosition = &f.ToICSPosition
	f.BaseMessage.FToLocation = &f.ToLocation
	f.BaseMessage.FFromICSPosition = &f.FromICSPosition
	f.BaseMessage.FFromLocation = &f.FromLocation
	f.BaseMessage.FSubject = &f.Subject
	f.BaseMessage.FReference = &f.Reference
	f.BaseMessage.FBody = &f.Message
	f.BaseMessage.FOpCall = &f.OpCall
	f.BaseMessage.FOpName = &f.OpName
	f.BaseMessage.FOpDate = &f.OpDate
	f.BaseMessage.FOpTime = &f.OpTime
	f.Fields = make([]*message.Field, 0, fieldCount)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			message.NewMessageNumberField(&message.Field{
				Label:       "Origin Message Number",
				Value:       &f.OriginMsgID,
				Presence:    message.Required,
				PIFOTag:     "2.", // may be changed to MsgNo by decode()
				PDFMap:      message.PDFName("Origin Msg #"),
				PDFRenderer: &message.PDFTextRenderer{X: 323, Y: 36, W: 64, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
				Compare:     message.CompareNone,
				EditHelp:    `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
			}),
			message.NewMessageNumberField(&message.Field{
				Label:      "My Message Number",
				Value:      &f.myMsgID,
				PIFOTag:    "MsgNo", // will be removed after decode()
				TableValue: message.TableOmit,
				Compare:    message.CompareNone,
			}),
		)
	} else {
		f.Fields = append(f.Fields,
			message.NewMessageNumberField(&message.Field{
				Label:       "Origin Message Number",
				Value:       &f.OriginMsgID,
				Presence:    message.Required,
				PIFOTag:     "MsgNo",
				PDFMap:      message.PDFName("Origin Msg #"),
				PDFRenderer: &message.PDFTextRenderer{X: 323, Y: 36, W: 64, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
				Compare:     message.CompareNone,
				EditHelp:    `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
			}),
		)
	}
	f.Fields = append(f.Fields,
		message.NewMessageNumberField(&message.Field{
			Label:       "Destination Message Number",
			Value:       &f.DestinationMsgID,
			PIFOTag:     "3.", // may be changed to MsgNo after decode()
			PDFMap:      message.PDFName("Destination Msg#"),
			PDFRenderer: &message.PDFTextRenderer{X: 520, Y: 36, W: 60, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date",
			Value:       &f.Date,
			Presence:    message.Required,
			PIFOTag:     "1a.",
			PDFMap:      message.PDFName("FormDate"),
			PDFRenderer: &message.PDFTextRenderer{X: 49, Y: 115, W: 55, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    "This is the date the message was written, in MM/DD/YYYY format.  It is required.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.Time,
			Presence:    message.Required,
			PIFOTag:     "1b.",
			PDFMap:      message.PDFName("FormTime"),
			PDFRenderer: &message.PDFTextRenderer{X: 143, Y: 115, W: 43, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    "This is the time the message was written, in HH:MM format (24-hour clock).  It is required.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			Presence: message.Required,
			EditHelp: "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
		}, &f.Date, &f.Time),
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:    "Situation Severity",
				Value:    &f.Severity,
				Choices:  message.Choices{"OTHER", "URGENT", "EMERGENCY"},
				Presence: message.Required,
				PIFOTag:  "4.",
				EditHelp: `This is the severity of the situation that caused the message to be sent.  It is required.`,
			}),
		)
	}
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:    "Handling",
			Value:    &f.Handling,
			Choices:  message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: message.Required,
			PIFOTag:  "5.",
			PDFMap:   message.PDFNameMap{"Immediate", "", "Off", "PRIORITY", "1", "ROUTINE", "2", "IMMEDIATE", "3"},
			PDFRenderer: &message.PDFRadioRenderer{
				Radius: 3.5,
				Points: map[string][]float64{
					"ROUTINE":    {509.5, 83.5},
					"PRIORITY":   {412, 83.5},
					"IMMEDIATE:": {304.5, 84},
				},
			},
			EditHelp: `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Take Action",
			Value:   &f.TakeAction,
			Choices: message.Choices{"Yes", "No"},
			PIFOTag: "6a.",
			PDFMap:  message.PDFNameMap{"TakeAction", "", "Off", "Yes", "1", "No", "2"},
			PDFRenderer: &message.PDFRadioRenderer{Points: map[string][]float64{
				"Yes": {421, 124},
				"No":  {522.5, 124},
			}},
			EditHelp: `This indicates whether the sender expects the recipient to take action based on this message.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Reply",
			Value:   &f.Reply,
			Choices: message.Choices{"Yes", "No"},
			PIFOTag: "6b.",
			TableValue: func(*message.Field) string {
				if f.Reply == "Yes" && f.ReplyBy != "" {
					return "Yes, by " + f.ReplyBy
				}
				return f.Reply
			},
			PDFMap: message.PDFNameMap{"Reply", "", "Off", "Yes", "Reply-Yes", "No", "Reply-No"},
			PDFRenderer: &message.PDFRadioRenderer{Points: map[string][]float64{
				"Yes": {421, 144.5},
				"No":  {522.5, 144.5},
			}},
			EditHelp: `This indicates whether the sender expects the recipient to reply to this message.`,
		}),
		message.NewTextField(&message.Field{
			Label: "Reply By",
			Value: &f.ReplyBy,
			Presence: func() (message.Presence, string) {
				if f.Reply != "Yes" {
					return message.PresenceNotAllowed, `unless "Reply" is "Yes"`
				}
				return message.PresenceOptional, ""
			},
			PIFOTag:     "6d.",
			Compare:     message.CompareExact,
			TableValue:  message.TableOmit,
			PDFMap:      message.PDFName("Reply_2"),
			PDFRenderer: &message.PDFTextRenderer{X: 463, Y: 138, W: 41, H: 13},
			EditWidth:   5,
			EditHelp:    `This is the time by which the sender expects a reply from the recipient.  It can be set only if "Reply" is "Yes".`,
		}),
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:    "For Your Information",
				Value:    &f.FYI,
				Choices:  message.Choices{"checked"},
				PIFOTag:  "6c.",
				EditHelp: `This indicates whether the message is merely informational.`,
			}),
		)
	}
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "To ICS Position",
			Value:       &f.ToICSPosition,
			Presence:    message.Required,
			PIFOTag:     "7.",
			PDFMap:      message.PDFName("TO ICS Position"),
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 167, W: 222, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the ICS position to which the message is addressed.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To Location",
			Value:       &f.ToLocation,
			Presence:    message.Required,
			PIFOTag:     "9a.",
			PDFMap:      message.PDFName("TO ICS Locatoin"),
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 202, W: 222, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the location of the recipient ICS position.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To Name",
			Value:       &f.ToName,
			PIFOTag:     "ToName",
			PDFMap:      message.PDFName("TO ICS Name"),
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 239, W: 222, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To Telephone #",
			Value:       &f.ToTelephone,
			PIFOTag:     "ToTel",
			Compare:     message.CompareExact,
			PDFMap:      message.PDFName("TO ICS Telephone"),
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 274, W: 222, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the telephone number for the receipient.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From ICS Position",
			Value:       &f.FromICSPosition,
			Presence:    message.Required,
			PIFOTag:     "8.",
			PDFMap:      message.PDFName("From ICS Position"),
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 167, W: 226, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the ICS position of the message author.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From Location",
			Value:       &f.FromLocation,
			Presence:    message.Required,
			PIFOTag:     "9b.",
			PDFMap:      message.PDFName("From ICS Location"),
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 202, W: 226, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the location of the message author.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From Name",
			Value:       &f.FromName,
			PIFOTag:     "FmName",
			PDFMap:      message.PDFName("From ICS Name"),
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 239, W: 226, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the name of the message author.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From Telephone #",
			Value:       &f.FromTelephone,
			PIFOTag:     "FmTel",
			Compare:     message.CompareExact,
			PDFMap:      message.PDFName("From ICS Telephone"),
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 274, W: 226, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the telephone number for the message author.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Subject",
			Value:       &f.Subject,
			Presence:    message.Required,
			PIFOTag:     "10.",
			PDFMap:      message.PDFName("Subject"),
			PDFRenderer: &message.PDFTextRenderer{X: 115, Y: 299, W: 452, H: 20},
			EditWidth:   80,
			EditHelp:    `This is the subject of the message.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Reference",
			Value:       &f.Reference,
			PIFOTag:     "11.",
			Compare:     message.CompareExact,
			PDFMap:      message.PDFName("Reference"),
			PDFRenderer: &message.PDFTextRenderer{X: 260, Y: 323, W: 305, H: 20},
			EditWidth:   55,
			EditHelp:    `This is the origin message number of the previous message to which this message is responding.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Message",
			Value:       &f.Message,
			Presence:    message.Required,
			PIFOTag:     "12.",
			PDFMap:      message.PDFName("Message"),
			PDFRenderer: &message.PDFTextRenderer{X: 40, Y: 365, W: 539, H: 121, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   80,
			Multiline:   true,
			EditHelp:    `This is the text of the message.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Relay Received",
			Value:       &f.OpRelayRcvd,
			PIFOTag:     "OpRelayRcvd",
			PDFMap:      message.PDFName("Relay Received"),
			PDFRenderer: &message.PDFTextRenderer{X: 123, Y: 590, W: 184, H: 18, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
			EditWidth:   36,
			EditHelp:    `This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Relay Sent",
			Value:       &f.OpRelaySent,
			PIFOTag:     "OpRelaySent",
			PDFMap:      message.PDFName("Relay Sent"),
			PDFRenderer: &message.PDFTextRenderer{X: 377, Y: 590, W: 183, H: 18, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
			EditWidth:   36,
			EditHelp:    `This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Operator: Receiver/Sender",
			Value:   &f.ReceivedSent,
			Choices: message.Choices{"sender", "receiver"},
			PIFOTag: "Rec-Sent",
			PDFMap:  message.PDFNameMap{"How: Received", "", "Off", "sender", "0", "receiver", "1"},
			PDFRenderer: &message.PDFRadioRenderer{Points: map[string][]float64{
				"sender":   {177.5, 619.5},
				"receiver": {87.5, 620.5},
			}},
			Compare: message.CompareNone,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Call Sign",
			Value:       &f.OpCall,
			Presence:    message.Required,
			PIFOTag:     "OpCall",
			TableValue:  message.TableOmit,
			PDFMap:      message.PDFName("Operation Call Sign"),
			PDFRenderer: &message.PDFTextRenderer{X: 422, Y: 612, W: 156, H: 19},
			Compare:     message.CompareNone,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Name",
			Value:       &f.OpName,
			Presence:    message.Required,
			PIFOTag:     "OpName",
			TableValue:  message.TableOmit,
			PDFMap:      message.PDFName("Relay Received_2"),
			PDFRenderer: &message.PDFTextRenderer{X: 408, Y: 637, W: 166, H: 15},
			Compare:     message.CompareNone,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Operator",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(f.OpCall, f.OpName, " ")
			},
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Operator: Tx Method",
			Value:    &f.TxMethod,
			Choices:  message.Choices{"Telephone", "Dispatch Center", "EOC Radio", "FAX", "Courier", "Amateur Radio", "Other"},
			Presence: message.Required,
			PIFOTag:  "Method",
			TableValue: func(*message.Field) string {
				if f.TxMethod == "Other" && f.OtherMethod != "" {
					return f.OtherMethod
				}
				return f.TxMethod
			},
			PDFMap: message.PDFNameMap{"Telephone", "", "Off", "Telephone", "1", "Dispatch Center", "2", "EOC Radio", "3", "FAX", "4", "Courier", "5", "Amateur Radio", "6", "Other", "7"},
			PDFRenderer: &message.PDFRadioRenderer{Points: map[string][]float64{
				"Telephone":       {42.5, 643},
				"Dispatch Center": {141.5, 643},
				"EOC Radio":       {42.5, 665.5},
				"FAX":             {141.5, 665.5},
				"Courier":         {218, 665.5},
				"Amateur Radio":   {42.5, 682.5},
				"Other":           {141.5, 682.5},
			}},
			Compare: message.CompareNone,
		}),
		message.NewTextField(&message.Field{
			Label: "Operator: Other Method",
			Value: &f.OtherMethod,
			Presence: func() (message.Presence, string) {
				if f.TxMethod != "Other" {
					return message.PresenceNotAllowed, `unless "Operator: Tx Method" is "Other"`
				}
				return message.PresenceOptional, ""
			},
			PIFOTag:     "Other",
			TableValue:  message.TableOmit,
			PDFMap:      message.PDFName("OtherText"),
			PDFRenderer: &message.PDFTextRenderer{X: 188, Y: 673, W: 100, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Operator: Date",
			Value:       &f.OpDate,
			Presence:    message.Required,
			PIFOTag:     "OpDate",
			TableValue:  message.TableOmit,
			PDFMap:      message.PDFName("OperatorDate"),
			PDFRenderer: &message.PDFTextRenderer{X: 354, Y: 673, W: 56, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Operator: Time",
			Value:       &f.OpTime,
			Presence:    message.Required,
			PIFOTag:     "OpTime",
			TableValue:  message.TableOmit,
			PDFMap:      message.PDFName("OperatorTime"),
			PDFRenderer: &message.PDFTextRenderer{X: 498, Y: 673, W: 40, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Operator: Date/Time",
			Presence: message.Required,
		}, &f.OpDate, &f.OpTime),
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
	if df, ok := message.DecodeForm(body, versions, create).(*ICS213); ok {
		f = df
	} else {
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
		f.Fields[2].PIFOTag = ""
	} else {
		f.OriginMsgID = f.myMsgID
		f.Fields[0].PIFOTag = ""
	}
	return f
}
