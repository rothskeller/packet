// Package baseform provides shared code for all of the standard form message
// types.
package baseform

import (
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
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
	OriginMsgID      basemsg.PDFMapper
	DestinationMsgID basemsg.PDFMapper
	MessageDate      basemsg.PDFMapper
	MessageTime      basemsg.PDFMapper
	Handling         basemsg.PDFMapper
	ToICSPosition    basemsg.PDFMapper
	ToLocation       basemsg.PDFMapper
	ToName           basemsg.PDFMapper
	ToContact        basemsg.PDFMapper
	FromICSPosition  basemsg.PDFMapper
	FromLocation     basemsg.PDFMapper
	FromName         basemsg.PDFMapper
	FromContact      basemsg.PDFMapper
	OpRelayRcvd      basemsg.PDFMapper
	OpRelaySent      basemsg.PDFMapper
	OpName           basemsg.PDFMapper
	OpCall           basemsg.PDFMapper
	OpDate           basemsg.PDFMapper
	OpTime           basemsg.PDFMapper
}

var DefaultPDFMaps = BaseFormPDFMaps{
	OriginMsgID:      basemsg.PDFName("Origin Msg Nbr"),
	DestinationMsgID: basemsg.PDFName("Destination Msg Nbr"),
	MessageDate:      basemsg.PDFName("Date Created"),
	MessageTime:      basemsg.PDFName("Time Created"),
	Handling: basemsg.PDFNameMap{"Handling",
		"", "Off",
		"IMMEDIATE", "Immediate",
		"PRIORITY", "Priority",
		"ROUTINE", "Routine",
	},
	ToICSPosition:   basemsg.PDFName("To ICS Position"),
	ToLocation:      basemsg.PDFName("To Location"),
	ToName:          basemsg.PDFName("To Name"),
	ToContact:       basemsg.PDFName("To Contact Info"),
	FromICSPosition: basemsg.PDFName("From ICS Position"),
	FromLocation:    basemsg.PDFName("From Location"),
	FromName:        basemsg.PDFName("From Name"),
	FromContact:     basemsg.PDFName("From Contact Info"),
	OpRelayRcvd:     basemsg.PDFName("Relay Rcvd"),
	OpRelaySent:     basemsg.PDFName("Relay Sent"),
	OpName:          basemsg.PDFName("Op Name"),
	OpCall:          basemsg.PDFName("Op Call Sign"),
	OpDate:          basemsg.PDFName("Op Date"),
	OpTime:          basemsg.PDFName("Op Time"),
}

func (bf *BaseForm) AddHeaderFields(bm *basemsg.BaseMessage, pdf *BaseFormPDFMaps) {
	bm.Fields = append(bm.Fields,
		basemsg.NewMessageNumberField(&basemsg.Field{
			Label:    "Origin Message Number",
			Value:    &bf.OriginMsgID,
			Presence: basemsg.Required,
			PIFOTag:  "MsgNo",
			PDFMap:   pdf.OriginMsgID,
			EditHelp: "This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.",
		}),
		basemsg.NewMessageNumberField(&basemsg.Field{
			Label:   "Destination Message Number",
			Value:   &bf.DestinationMsgID,
			PIFOTag: "DestMsgNo",
			PDFMap:  pdf.DestinationMsgID,
		}),
		basemsg.NewDateWithTimeField(&basemsg.Field{
			Label:    "Message Date",
			Value:    &bf.MessageDate,
			Presence: basemsg.Required,
			PIFOTag:  "1a.",
			PDFMap:   pdf.MessageDate,
		}),
		basemsg.NewTimeWithDateField(&basemsg.Field{
			Label:    "Message Time",
			Value:    &bf.MessageTime,
			Presence: basemsg.Required,
			PIFOTag:  "1b.",
			PDFMap:   pdf.MessageTime,
		}),
		basemsg.NewDateTimeField(&basemsg.Field{
			Label:    "Message Date/Time",
			Presence: basemsg.Required,
			EditHelp: "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
		}, &bf.MessageDate, &bf.MessageTime),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Handling",
			Value:    &bf.Handling,
			Choices:  basemsg.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence: basemsg.Required,
			PIFOTag:  "5.",
			PDFMap:   pdf.Handling,
			EditHelp: `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "To ICS Position",
			Value:     &bf.ToICSPosition,
			Presence:  basemsg.Required,
			PIFOTag:   "7a.",
			PDFMap:    pdf.ToICSPosition,
			EditWidth: 30,
			EditHelp:  "This is the ICS position to which the message is addressed.  It is required.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "To Location",
			Value:     &bf.ToLocation,
			Presence:  basemsg.Required,
			PIFOTag:   "7b.",
			PDFMap:    pdf.ToLocation,
			EditWidth: 32,
			EditHelp:  "This is the location of the recipient ICS position.  It is required.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "To Name",
			Value:     &bf.ToName,
			PIFOTag:   "7c.",
			PDFMap:    pdf.ToName,
			EditWidth: 34,
			EditHelp:  "This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "To Contact Info",
			Value:     &bf.ToContact,
			PIFOTag:   "7d.",
			PDFMap:    pdf.ToContact,
			EditWidth: 29,
			EditHelp:  "This is contact information (phone number, email, etc.) for the receipient.  It is optional and rarely provided.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "From ICS Position",
			Value:     &bf.FromICSPosition,
			Presence:  basemsg.Required,
			PIFOTag:   "8a.",
			PDFMap:    pdf.FromICSPosition,
			EditWidth: 30,
			EditHelp:  "This is the ICS position of the message author.  It is required.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "From Location",
			Value:     &bf.FromLocation,
			Presence:  basemsg.Required,
			PIFOTag:   "8b.",
			PDFMap:    pdf.FromLocation,
			EditWidth: 32,
			EditHelp:  "This is the location of the message author.  It is required.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "From Name",
			Value:     &bf.FromName,
			PIFOTag:   "8c.",
			PDFMap:    pdf.FromName,
			EditWidth: 34,
			EditHelp:  "This is the name of the message author.  It is optional and rarely provided.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "From Contact Info",
			Value:     &bf.FromContact,
			PIFOTag:   "8d.",
			PDFMap:    pdf.FromContact,
			EditWidth: 29,
			EditHelp:  "This is contact information (phone number, email, etc.) for the message author.  It is optional and rarely provided.",
		}),
	)
	bm.FOriginMsgID = &bf.OriginMsgID
	bm.FHandling = &bf.Handling
	bm.FToICSPosition = &bf.ToICSPosition
	bm.FToLocation = &bf.ToLocation
}
func (bf *BaseForm) AddFooterFields(bm *basemsg.BaseMessage, pdf *BaseFormPDFMaps) {
	bm.Fields = append(bm.Fields,
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Operator: Relay Received",
			Value:     &bf.OpRelayRcvd,
			PIFOTag:   "OpRelayRcvd",
			PDFMap:    pdf.OpRelayRcvd,
			Compare:   basemsg.CompareNone,
			EditWidth: 36,
			EditHelp:  "This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Operator: Relay Sent",
			Value:     &bf.OpRelaySent,
			PIFOTag:   "OpRelaySent",
			PDFMap:    pdf.OpRelaySent,
			Compare:   basemsg.CompareNone,
			EditWidth: 36,
			EditHelp:  "This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.",
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Operator: Name",
			Value:      &bf.OpName,
			Presence:   basemsg.Required,
			PIFOTag:    "OpName",
			PDFMap:     pdf.OpName,
			TableValue: basemsg.TableOmit,
			Compare:    basemsg.CompareNone,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Operator: Call Sign",
			Value:      &bf.OpCall,
			Presence:   basemsg.Required,
			PIFOTag:    "OpCall",
			PDFMap:     pdf.OpCall,
			TableValue: basemsg.TableOmit,
			Compare:    basemsg.CompareNone,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "Operator",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(bf.OpCall, bf.OpName, " ")
			},
		}),
		basemsg.NewDateWithTimeField(&basemsg.Field{
			Label:    "Operator: Date",
			Value:    &bf.OpDate,
			Presence: basemsg.Required,
			PIFOTag:  "OpDate",
			PDFMap:   pdf.OpDate,
			Compare:  basemsg.CompareNone,
		}),
		basemsg.NewTimeWithDateField(&basemsg.Field{
			Label:    "Operator: Time",
			Value:    &bf.OpTime,
			Presence: basemsg.Required,
			PIFOTag:  "OpTime",
			PDFMap:   pdf.OpTime,
			Compare:  basemsg.CompareNone,
		}),
		basemsg.NewDateTimeField(&basemsg.Field{
			Label: "Operator: Date/Time",
		}, &bf.OpDate, &bf.OpTime),
	)
	bm.FOpCall = &bf.OpCall
	bm.FOpName = &bf.OpName
}
