// Package ics213 defines the ICS-213 General Message Form message type.
package ics213

import (
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type22 is the type definition for an ICS-213 general message form.
var Type22 = message.Type{
	Tag:     "ICS213",
	HTML:    "form-ics213.html",
	Version: "2.2",
	Name:    "ICS-213 general message form",
	Article: "an",
	FieldOrder: []string{
		"MsgNo", "3.", "1a.", "1b.", "5.", "6a.", "6b.", "6d.", "7.",
		"8.", "9a.", "9b.", "ToName", "FmName", "ToTel", "FmTel", "10.",
		"11.", "12.", "OpRelayRcvd", "OpRelaySent", "Rec-Sent",
		"OpCall", "OpName", "Method", "Other", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type22, decode22, create22)
}

// ICS213v22 holds an ICS-213 general message form.
type ICS213v22 struct {
	message.BaseMessage
	OriginMsgID      string
	DestinationMsgID string
	Date             string
	Time             string
	Handling         string
	TakeAction       string
	Reply            string
	ReplyBy          string
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

func create22() message.Message {
	f := make22()
	f.Date = time.Now().Format("01/02/2006")
	f.ReceivedSent = "sender"
	f.TxMethod = "Other"
	f.OtherMethod = "Packet"
	return f
}

func (f *ICS213v22) SetOperator(opcall, opname string, received bool) {
	if received {
		f.ReceivedSent = "receiver"
	} else {
		f.ReceivedSent = "sender"
	}
	f.TxMethod, f.OtherMethod = "Other", "Packet"
	f.BaseMessage.SetOperator(opcall, opname, received)
}

func make22() (f *ICS213v22) {
	const fieldCount = 31
	f = &ICS213v22{BaseMessage: message.BaseMessage{Type: &Type22}}
	f.FOriginMsgID = &f.OriginMsgID
	f.FDestinationMsgID = &f.DestinationMsgID
	f.FMessageDate = &f.Date
	f.FMessageTime = &f.Time
	f.FHandling = &f.Handling
	f.FToICSPosition = &f.ToICSPosition
	f.FToLocation = &f.ToLocation
	f.FFromICSPosition = &f.FromICSPosition
	f.FFromLocation = &f.FromLocation
	f.FSubject = &f.Subject
	f.FReference = &f.Reference
	f.FBody = &f.Message
	f.FOpCall = &f.OpCall
	f.FOpName = &f.OpName
	f.FOpDate = &f.OpDate
	f.FOpTime = &f.OpTime
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.Fields = append(f.Fields,
		message.NewMessageNumberField(&message.Field{
			Label:       "Origin Message Number",
			Value:       &f.OriginMsgID,
			Presence:    message.Required,
			PIFOTag:     "MsgNo",
			PDFRenderer: &message.PDFTextRenderer{X: 323, Y: 36, W: 64, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
			EditHelp:    `This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.`,
		}),
		message.NewMessageNumberField(&message.Field{
			Label:       "Destination Message Number",
			Value:       &f.DestinationMsgID,
			PIFOTag:     "3.",
			PDFRenderer: &message.PDFTextRenderer{X: 520, Y: 36, W: 60, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date",
			Value:       &f.Date,
			Presence:    message.Required,
			PIFOTag:     "1a.",
			PDFRenderer: &message.PDFTextRenderer{X: 49, Y: 115, W: 55, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    "This is the date the message was written, in MM/DD/YYYY format.  It is required.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.Time,
			Presence:    message.Required,
			PIFOTag:     "1b.",
			PDFRenderer: &message.PDFTextRenderer{X: 143, Y: 115, W: 43, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    "This is the time the message was written, in HH:MM format (24-hour clock).  It is required.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			Presence: message.Required,
			EditHelp: "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
		}, &f.Date, &f.Time),
		message.NewRestrictedField(&message.Field{
			Label:    "Handling",
			Value:    &f.Handling,
			Choices:  message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: message.Required,
			PIFOTag:  "5.",
			PDFRenderer: &message.PDFRadioRenderer{
				Radius: 3.5,
				Points: map[string][]float64{
					"ROUTINE":   {509.5, 83.5},
					"PRIORITY":  {412, 83.5},
					"IMMEDIATE": {304.5, 84},
				},
			},
			EditHelp: `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Take Action",
			Value:   &f.TakeAction,
			Choices: message.Choices{"Yes", "No"},
			PIFOTag: "6a.",
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
			PDFRenderer: &message.PDFTextRenderer{X: 463, Y: 138, W: 41, H: 13},
			EditWidth:   5,
			EditHelp:    `This is the time by which the sender expects a reply from the recipient.  It can be set only if "Reply" is "Yes".`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To ICS Position",
			Value:       &f.ToICSPosition,
			Presence:    message.Required,
			PIFOTag:     "7.",
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 167, W: 222, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the ICS position to which the message is addressed.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To Location",
			Value:       &f.ToLocation,
			Presence:    message.Required,
			PIFOTag:     "9a.",
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 202, W: 222, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the location of the recipient ICS position.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To Name",
			Value:       &f.ToName,
			PIFOTag:     "ToName",
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 239, W: 222, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To Telephone #",
			Value:       &f.ToTelephone,
			PIFOTag:     "ToTel",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{X: 61, Y: 274, W: 222, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the telephone number for the receipient.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From ICS Position",
			Value:       &f.FromICSPosition,
			Presence:    message.Required,
			PIFOTag:     "8.",
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 167, W: 226, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the ICS position of the message author.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From Location",
			Value:       &f.FromLocation,
			Presence:    message.Required,
			PIFOTag:     "9b.",
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 202, W: 226, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the location of the message author.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From Name",
			Value:       &f.FromName,
			PIFOTag:     "FmName",
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 239, W: 226, H: 19, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the name of the message author.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "From Telephone #",
			Value:       &f.FromTelephone,
			PIFOTag:     "FmTel",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{X: 340, Y: 274, W: 226, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   40,
			EditHelp:    `This is the telephone number for the message author.  It is optional and rarely provided.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Subject",
			Value:       &f.Subject,
			Presence:    message.Required,
			PIFOTag:     "10.",
			PDFRenderer: &message.PDFTextRenderer{X: 115, Y: 299, W: 452, H: 20},
			EditWidth:   80,
			EditHelp:    `This is the subject of the message.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Reference",
			Value:       &f.Reference,
			PIFOTag:     "11.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{X: 260, Y: 323, W: 305, H: 20},
			EditWidth:   55,
			EditHelp:    `This is the origin message number of the previous message to which this message is responding.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Message",
			Value:       &f.Message,
			Presence:    message.Required,
			PIFOTag:     "12.",
			PDFRenderer: &message.PDFTextRenderer{X: 40, Y: 365, W: 539, H: 121, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   80,
			Multiline:   true,
			EditHelp:    `This is the text of the message.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Relay Received",
			Value:       &f.OpRelayRcvd,
			PIFOTag:     "OpRelayRcvd",
			PDFRenderer: &message.PDFTextRenderer{X: 123, Y: 590, W: 184, H: 18, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
			EditWidth:   36,
			EditHelp:    `This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Relay Sent",
			Value:       &f.OpRelaySent,
			PIFOTag:     "OpRelaySent",
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
			PDFRenderer: &message.PDFTextRenderer{X: 422, Y: 612, W: 156, H: 19},
			Compare:     message.CompareNone,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Name",
			Value:       &f.OpName,
			Presence:    message.Required,
			PIFOTag:     "OpName",
			TableValue:  message.TableOmit,
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
			PDFRenderer: &message.PDFTextRenderer{X: 188, Y: 673, W: 100, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Operator: Date",
			Value:       &f.OpDate,
			Presence:    message.Required,
			PIFOTag:     "OpDate",
			TableValue:  message.TableOmit,
			PDFRenderer: &message.PDFTextRenderer{X: 354, Y: 673, W: 56, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Operator: Time",
			Value:       &f.OpTime,
			Presence:    message.Required,
			PIFOTag:     "OpTime",
			TableValue:  message.TableOmit,
			PDFRenderer: &message.PDFTextRenderer{X: 498, Y: 673, W: 40, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
			Compare:     message.CompareNone,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Operator: Date/Time",
			Presence: message.Required,
		}, &f.OpDate, &f.OpTime),
	)
	if len(f.Fields) > fieldCount {
		panic("update ICS213v22 fieldCount")
	}
	return f
}

func decode22(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type22.HTML || form.FormVersion != Type22.Version {
		return nil
	}
	df := make22()
	message.DecodeForm(form, df)
	return df
}

func (f *ICS213v22) Compare(actual message.Message) (int, int, []*message.CompareField) {
	switch act := actual.(type) {
	case *ICS213v21:
		actual = act.convertTo22()
	}
	return f.BaseMessage.Compare(actual)
}
