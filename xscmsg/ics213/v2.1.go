// Package ics213 defines the ICS-213 General Message Form message type.
package ics213

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// Type21 is the type definition for an ICS-213 general message form.
var Type21 = message.Type{
	Tag:     "ICS213",
	HTML:    "form-ics213.html",
	Version: "2.1",
	Name:    "ICS-213 general message form",
	Article: "an",
	FieldOrder: []string{
		"2.", "MsgNo", "3.", "1a.", "1b.", "4.", "5.", "6a.", "6b.", "6d.", "6c.", "7.",
		"8.", "9a.", "9b.", "ToName", "FmName", "ToTel", "FmTel", "10.",
		"11.", "12.", "OpRelayRcvd", "OpRelaySent", "Rec-Sent",
		"OpCall", "OpName", "Method", "Other", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type21, decode21, nil)
}

// ICS213v21 holds an ICS-213 general message form.
type ICS213v21 struct {
	message.BaseMessage
	OriginMsgID      string
	myMsgID          string
	DestinationMsgID string
	Date             string
	Time             string
	Severity         string
	Handling         string
	TakeAction       string
	Reply            string
	ReplyBy          string
	FYI              string
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

func (f *ICS213v21) SetOperator(opcall, opname string, received bool) {
	if received {
		f.ReceivedSent = "receiver"
	} else {
		f.ReceivedSent = "sender"
	}
	f.TxMethod, f.OtherMethod = "Other", "Packet"
	f.BaseMessage.SetOperator(opcall, opname, received)
}

func make21() (f *ICS213v21) {
	const fieldCount = 34
	f = &ICS213v21{BaseMessage: message.BaseMessage{Type: &Type21}}
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
			Label:    "Origin Message Number",
			Value:    &f.OriginMsgID,
			Presence: message.Required,
			PIFOTag:  "2.", // may be changed to MsgNo by decode21()
		}),
		message.NewMessageNumberField(&message.Field{
			Label:      "My Message Number",
			Value:      &f.myMsgID,
			PIFOTag:    "MsgNo", // will be removed after decode21()
			TableValue: message.TableOmit,
		}),
		message.NewMessageNumberField(&message.Field{
			Label:   "Destination Message Number",
			Value:   &f.DestinationMsgID,
			PIFOTag: "3.",
		}),
		message.NewDateField(true, &message.Field{
			Label:    "Date",
			Value:    &f.Date,
			Presence: message.Required,
			PIFOTag:  "1a.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:    "Time",
			Value:    &f.Time,
			Presence: message.Required,
			PIFOTag:  "1b.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			Presence: message.Required,
		}, &f.Date, &f.Time),
		message.NewRestrictedField(&message.Field{
			Label:    "Situation Severity",
			Value:    &f.Severity,
			Choices:  message.Choices{"OTHER", "URGENT", "EMERGENCY"},
			Presence: message.Required,
			PIFOTag:  "4.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Handling",
			Value:    &f.Handling,
			Choices:  message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: message.Required,
			PIFOTag:  "5.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Take Action",
			Value:   &f.TakeAction,
			Choices: message.Choices{"Yes", "No"},
			PIFOTag: "6a.",
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
			PIFOTag:    "6d.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "For Your Information",
			Value:   &f.FYI,
			Choices: message.Choices{"checked"},
			PIFOTag: "6c.",
		}),
		message.NewTextField(&message.Field{
			Label:    "To ICS Position",
			Value:    &f.ToICSPosition,
			Presence: message.Required,
			PIFOTag:  "7.",
		}),
		message.NewTextField(&message.Field{
			Label:    "To Location",
			Value:    &f.ToLocation,
			Presence: message.Required,
			PIFOTag:  "9a.",
		}),
		message.NewTextField(&message.Field{
			Label:   "To Name",
			Value:   &f.ToName,
			PIFOTag: "ToName",
		}),
		message.NewTextField(&message.Field{
			Label:   "To Telephone #",
			Value:   &f.ToTelephone,
			PIFOTag: "ToTel",
		}),
		message.NewTextField(&message.Field{
			Label:    "From ICS Position",
			Value:    &f.FromICSPosition,
			Presence: message.Required,
			PIFOTag:  "8.",
		}),
		message.NewTextField(&message.Field{
			Label:    "From Location",
			Value:    &f.FromLocation,
			Presence: message.Required,
			PIFOTag:  "9b.",
		}),
		message.NewTextField(&message.Field{
			Label:   "From Name",
			Value:   &f.FromName,
			PIFOTag: "FmName",
		}),
		message.NewTextField(&message.Field{
			Label:   "From Telephone #",
			Value:   &f.FromTelephone,
			PIFOTag: "FmTel",
		}),
		message.NewTextField(&message.Field{
			Label:    "Subject",
			Value:    &f.Subject,
			Presence: message.Required,
			PIFOTag:  "10.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Reference",
			Value:   &f.Reference,
			PIFOTag: "11.",
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Message",
			Value:     &f.Message,
			Presence:  message.Required,
			PIFOTag:   "12.",
			Multiline: true,
		}),
		message.NewTextField(&message.Field{
			Label:   "Operator: Relay Received",
			Value:   &f.OpRelayRcvd,
			PIFOTag: "OpRelayRcvd",
		}),
		message.NewTextField(&message.Field{
			Label:   "Operator: Relay Sent",
			Value:   &f.OpRelaySent,
			PIFOTag: "OpRelaySent",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Operator: Receiver/Sender",
			Value:   &f.ReceivedSent,
			Choices: message.Choices{"sender", "receiver"},
			PIFOTag: "Rec-Sent",
		}),
		message.NewTextField(&message.Field{
			Label:      "Operator: Call Sign",
			Value:      &f.OpCall,
			Presence:   message.Required,
			PIFOTag:    "OpCall",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label:      "Operator: Name",
			Value:      &f.OpName,
			Presence:   message.Required,
			PIFOTag:    "OpName",
			TableValue: message.TableOmit,
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
			PIFOTag:    "Other",
			TableValue: message.TableOmit,
		}),
		message.NewDateField(true, &message.Field{
			Label:      "Operator: Date",
			Value:      &f.OpDate,
			Presence:   message.Required,
			PIFOTag:    "OpDate",
			TableValue: message.TableOmit,
		}),
		message.NewTimeField(true, &message.Field{
			Label:      "Operator: Time",
			Value:      &f.OpTime,
			Presence:   message.Required,
			PIFOTag:    "OpTime",
			TableValue: message.TableOmit,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Operator: Date/Time",
			Presence: message.Required,
		}, &f.OpDate, &f.OpTime),
	)
	if len(f.Fields) > fieldCount {
		panic("update ICS213v21 fieldCount")
	}
	return f
}

func decode21(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type21.HTML || form.FormVersion != Type21.Version {
		return nil
	}
	df := make21()
	message.DecodeForm(form, df)
	// If we got an OriginMsgID, or we are the receiver, move myMsgID to
	// DestinationMsgID.
	if df.OriginMsgID != "" || df.ReceivedSent == "receiver" {
		df.DestinationMsgID = df.myMsgID
		df.Fields[2].PIFOTag = ""
	} else {
		df.OriginMsgID = df.myMsgID
		df.Fields[0].PIFOTag = ""
	}
	return df
}

func (f *ICS213v21) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo22().Compare(actual)
}

func (f *ICS213v21) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo22().RenderPDF(env, filename)
}

func (f *ICS213v21) convertTo22() (c *ICS213v22) {
	c = make22()
	c.OriginMsgID = f.OriginMsgID
	c.DestinationMsgID = f.DestinationMsgID
	c.Date = f.Date
	c.Time = f.Time
	c.Handling = f.Handling
	c.TakeAction = f.TakeAction
	c.Reply = f.Reply
	c.ReplyBy = f.ReplyBy
	c.ToICSPosition = f.ToICSPosition
	c.ToLocation = f.ToLocation
	c.ToName = f.ToName
	c.ToTelephone = f.ToTelephone
	c.FromICSPosition = f.FromICSPosition
	c.FromLocation = f.FromLocation
	c.FromName = f.FromName
	c.FromTelephone = f.FromTelephone
	c.Subject = f.Subject
	c.Reference = f.Reference
	c.Message = f.Message
	c.OpRelayRcvd = f.OpRelayRcvd
	c.OpRelaySent = f.OpRelaySent
	c.ReceivedSent = f.ReceivedSent
	c.OpCall = f.OpCall
	c.OpName = f.OpName
	c.TxMethod = f.TxMethod
	c.OtherMethod = f.OtherMethod
	c.OpDate = f.OpDate
	c.OpTime = f.OpTime
	return c
}
