// Package baseform provides shared code for all of the standard form message
// types.
package baseform

import (
	"github.com/rothskeller/packet/message"
)

type BaseForm struct {
	// Metadata Fields

	PIFOVersion string
	FormVersion string

	// Header Fields

	OriginMsgID      string
	DestinationMsgID string
	MessageDate      string
	MessageTime      string
	Handling         string
	ToICSPosition    string
	ToLocation       string
	ToName           string
	ToContact        string
	FromICSPosition  string
	FromLocation     string
	FromName         string
	FromContact      string

	// Footer Fields

	OpRelayRcvd string
	OpRelaySent string
	OpName      string
	OpCall      string
	OpDate      string
	OpTime      string
}

type BaseFormPDFMaps struct {
	OriginMsgID      message.PDFMapper
	DestinationMsgID message.PDFMapper
	MessageDate      message.PDFMapper
	MessageTime      message.PDFMapper
	Handling         message.PDFMapper
	ToICSPosition    message.PDFMapper
	ToLocation       message.PDFMapper
	ToName           message.PDFMapper
	ToContact        message.PDFMapper
	FromICSPosition  message.PDFMapper
	FromLocation     message.PDFMapper
	FromName         message.PDFMapper
	FromContact      message.PDFMapper
	OpRelayRcvd      message.PDFMapper
	OpRelaySent      message.PDFMapper
	OpName           message.PDFMapper
	OpCall           message.PDFMapper
	OpDate           message.PDFMapper
	OpTime           message.PDFMapper
}

var DefaultPDFMaps = BaseFormPDFMaps{
	OriginMsgID:      message.PDFName("Origin Msg Nbr"),
	DestinationMsgID: message.PDFName("Destination Msg Nbr"),
	MessageDate:      message.PDFName("Date Created"),
	MessageTime:      message.PDFName("Time Created"),
	Handling: message.PDFNameMap{"Handling",
		"", "Off",
		"IMMEDIATE", "Immediate",
		"PRIORITY", "Priority",
		"ROUTINE", "Routine",
	},
	ToICSPosition:   message.PDFName("To ICS Position"),
	ToLocation:      message.PDFName("To Location"),
	ToName:          message.PDFName("To Name"),
	ToContact:       message.PDFName("To Contact Info"),
	FromICSPosition: message.PDFName("From ICS Position"),
	FromLocation:    message.PDFName("From Location"),
	FromName:        message.PDFName("From Name"),
	FromContact:     message.PDFName("From Contact Info"),
	OpRelayRcvd:     message.PDFName("Relay Rcvd"),
	OpRelaySent:     message.PDFName("Relay Sent"),
	OpName:          message.PDFName("Op Name"),
	OpCall:          message.PDFName("Op Call Sign"),
	OpDate:          message.PDFName("Op Date"),
	OpTime:          message.PDFName("Op Time"),
}

func (bf *BaseForm) AddHeaderFields(bm *message.BaseMessage, pdf *BaseFormPDFMaps) {
	bm.Fields = append(bm.Fields,
		message.NewMessageNumberField(&message.Field{
			Label:    "Origin Message Number",
			Value:    &bf.OriginMsgID,
			Presence: message.Required,
			PIFOTag:  "MsgNo",
			PDFMap:   pdf.OriginMsgID,
			EditHelp: "This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.",
		}),
		message.NewMessageNumberField(&message.Field{
			Label:   "Destination Message Number",
			Value:   &bf.DestinationMsgID,
			PIFOTag: "DestMsgNo",
			PDFMap:  pdf.DestinationMsgID,
		}),
		message.NewDateField(true, &message.Field{
			Label:    "Message Date",
			Value:    &bf.MessageDate,
			Presence: message.Required,
			PIFOTag:  "1a.",
			PDFMap:   pdf.MessageDate,
			EditHelp: "This is the date the message was written, in MM/DD/YYYY format.  It is required.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:    "Message Time",
			Value:    &bf.MessageTime,
			Presence: message.Required,
			PIFOTag:  "1b.",
			PDFMap:   pdf.MessageTime,
			EditHelp: "This is the time the message was written, in HH:MM format (24-hour clock).  It is required.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Message Date/Time",
			Presence: message.Required,
			EditHelp: "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
		}, &bf.MessageDate, &bf.MessageTime),
		message.NewRestrictedField(&message.Field{
			Label:    "Handling",
			Value:    &bf.Handling,
			Choices:  message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: message.Required,
			PIFOTag:  "5.",
			PDFMap:   pdf.Handling,
			EditHelp: `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:     "To ICS Position",
			Value:     &bf.ToICSPosition,
			Presence:  message.Required,
			PIFOTag:   "7a.",
			PDFMap:    pdf.ToICSPosition,
			EditWidth: 30,
			EditHelp:  "This is the ICS position to which the message is addressed.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:     "To Location",
			Value:     &bf.ToLocation,
			Presence:  message.Required,
			PIFOTag:   "7b.",
			PDFMap:    pdf.ToLocation,
			EditWidth: 32,
			EditHelp:  "This is the location of the recipient ICS position.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:     "To Name",
			Value:     &bf.ToName,
			PIFOTag:   "7c.",
			PDFMap:    pdf.ToName,
			EditWidth: 34,
			EditHelp:  "This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.",
		}),
		message.NewTextField(&message.Field{
			Label:     "To Contact Info",
			Value:     &bf.ToContact,
			PIFOTag:   "7d.",
			PDFMap:    pdf.ToContact,
			EditWidth: 29,
			EditHelp:  "This is contact information (phone number, email, etc.) for the receipient.  It is optional and rarely provided.",
		}),
		message.NewTextField(&message.Field{
			Label:     "From ICS Position",
			Value:     &bf.FromICSPosition,
			Presence:  message.Required,
			PIFOTag:   "8a.",
			PDFMap:    pdf.FromICSPosition,
			EditWidth: 30,
			EditHelp:  "This is the ICS position of the message author.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:     "From Location",
			Value:     &bf.FromLocation,
			Presence:  message.Required,
			PIFOTag:   "8b.",
			PDFMap:    pdf.FromLocation,
			EditWidth: 32,
			EditHelp:  "This is the location of the message author.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:     "From Name",
			Value:     &bf.FromName,
			PIFOTag:   "8c.",
			PDFMap:    pdf.FromName,
			EditWidth: 34,
			EditHelp:  "This is the name of the message author.  It is optional and rarely provided.",
		}),
		message.NewTextField(&message.Field{
			Label:     "From Contact Info",
			Value:     &bf.FromContact,
			PIFOTag:   "8d.",
			PDFMap:    pdf.FromContact,
			EditWidth: 29,
			EditHelp:  "This is contact information (phone number, email, etc.) for the message author.  It is optional and rarely provided.",
		}),
	)
	bm.FOriginMsgID = &bf.OriginMsgID
	bm.FDestinationMsgID = &bf.DestinationMsgID
	bm.FMessageDate = &bf.MessageDate
	bm.FMessageTime = &bf.MessageTime
	bm.FHandling = &bf.Handling
	bm.FToICSPosition = &bf.ToICSPosition
	bm.FToLocation = &bf.ToLocation
	bm.FFromICSPosition = &bf.FromICSPosition
	bm.FFromLocation = &bf.FromLocation
}
func (bf *BaseForm) AddFooterFields(bm *message.BaseMessage, pdf *BaseFormPDFMaps) {
	bm.Fields = append(bm.Fields,
		message.NewTextField(&message.Field{
			Label:     "Operator: Relay Received",
			Value:     &bf.OpRelayRcvd,
			PIFOTag:   "OpRelayRcvd",
			PDFMap:    pdf.OpRelayRcvd,
			Compare:   message.CompareNone,
			EditWidth: 36,
			EditHelp:  "This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.",
		}),
		message.NewTextField(&message.Field{
			Label:     "Operator: Relay Sent",
			Value:     &bf.OpRelaySent,
			PIFOTag:   "OpRelaySent",
			PDFMap:    pdf.OpRelaySent,
			Compare:   message.CompareNone,
			EditWidth: 36,
			EditHelp:  "This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.",
		}),
		message.NewTextField(&message.Field{
			Label:      "Operator: Name",
			Value:      &bf.OpName,
			Presence:   message.Required,
			PIFOTag:    "OpName",
			PDFMap:     pdf.OpName,
			TableValue: message.TableOmit,
			Compare:    message.CompareNone,
		}),
		message.NewTextField(&message.Field{
			Label:      "Operator: Call Sign",
			Value:      &bf.OpCall,
			Presence:   message.Required,
			PIFOTag:    "OpCall",
			PDFMap:     pdf.OpCall,
			TableValue: message.TableOmit,
			Compare:    message.CompareNone,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Operator",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(bf.OpCall, bf.OpName, " ")
			},
		}),
		message.NewDateField(true, &message.Field{
			Label:    "Operator: Date",
			Value:    &bf.OpDate,
			Presence: message.Required,
			PIFOTag:  "OpDate",
			PDFMap:   pdf.OpDate,
			Compare:  message.CompareNone,
		}),
		message.NewTimeField(true, &message.Field{
			Label:    "Operator: Time",
			Value:    &bf.OpTime,
			Presence: message.Required,
			PIFOTag:  "OpTime",
			PDFMap:   pdf.OpTime,
			Compare:  message.CompareNone,
		}),
		message.NewDateTimeField(&message.Field{
			Label: "Operator: Date/Time",
		}, &bf.OpDate, &bf.OpTime),
	)
	bm.FOpCall = &bf.OpCall
	bm.FOpName = &bf.OpName
	bm.FOpDate = &bf.OpDate
	bm.FOpTime = &bf.OpTime
}
