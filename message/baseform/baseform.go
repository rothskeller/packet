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
		&basemsg.Field{
			Label:     "Origin Message Number",
			PIFOTag:   "MsgNo",
			Value:     &bf.OriginMsgID,
			Presence:  basemsg.Required,
			PDFMap:    pdf.OriginMsgID,
			EditWidth: 9,
			EditHelp:  "This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.",
			EditHint:  "XXX-###P",
			EditApply: basemsg.ApplyMessageNumber,
			EditValid: basemsg.ValidMessageNumber,
		},
		&basemsg.Field{
			Label:   "Destination Message Number",
			PIFOTag: "DestMsgNo",
			Value:   &bf.DestinationMsgID,
			PDFMap:  pdf.DestinationMsgID,
		},
		&basemsg.Field{
			Label:      "Message Date",
			PIFOTag:    "1a.",
			Value:      &bf.MessageDate,
			Presence:   basemsg.Required,
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareDate,
			PDFMap:     pdf.MessageDate,
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Message Time",
			PIFOTag:    "1b.",
			Value:      &bf.MessageTime,
			Presence:   basemsg.Required,
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareTime,
			PDFMap:     pdf.MessageTime,
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:    "Message Date/Time",
			Presence: basemsg.Required,
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(bf.MessageDate, bf.MessageTime, " ")
			},
			EditWidth: 16,
			EditHelp:  "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(bf.MessageDate, bf.MessageTime)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&bf.MessageDate, &bf.MessageTime, value)
			},
			EditValid: func(f *basemsg.Field) string {
				return basemsg.ValidDateTime(f, bf.MessageDate, bf.MessageTime)
			},
		},
		&basemsg.Field{
			Label:     "Handling",
			PIFOTag:   "5.",
			Value:     &bf.Handling,
			Choices:   basemsg.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence:  basemsg.Required,
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    pdf.Handling,
			EditWidth: 9,
			EditHelp:  `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		},
		&basemsg.Field{
			Label:     "To ICS Position",
			PIFOTag:   "7a.",
			Value:     &bf.ToICSPosition,
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    pdf.ToICSPosition,
			EditWidth: 30,
			EditHelp:  "This is the ICS position to which the message is addressed.  It is required.",
		},
		&basemsg.Field{
			Label:     "To Location",
			PIFOTag:   "7b.",
			Value:     &bf.ToLocation,
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    pdf.ToLocation,
			EditWidth: 32,
			EditHelp:  "This is the location of the recipient ICS position.  It is required.",
		},
		&basemsg.Field{
			Label:     "To Name",
			PIFOTag:   "7c.",
			Value:     &bf.ToName,
			Compare:   common.CompareText,
			PDFMap:    pdf.ToName,
			EditWidth: 34,
			EditHelp:  "This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.",
		},
		&basemsg.Field{
			Label:     "To Contact Info",
			PIFOTag:   "7d.",
			Value:     &bf.ToContact,
			Compare:   common.CompareText,
			PDFMap:    pdf.ToContact,
			EditWidth: 29,
			EditHelp:  "This is contact information (phone number, email, etc.) for the receipient.  It is optional and rarely provided.",
		},
		&basemsg.Field{
			Label:     "From ICS Position",
			PIFOTag:   "8a.",
			Value:     &bf.FromICSPosition,
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    pdf.FromICSPosition,
			EditWidth: 30,
			EditHelp:  "This is the ICS position of the message author.  It is required.",
		},
		&basemsg.Field{
			Label:     "From Location",
			PIFOTag:   "8b.",
			Value:     &bf.FromLocation,
			Presence:  basemsg.Required,
			Compare:   common.CompareText,
			PDFMap:    pdf.FromLocation,
			EditWidth: 32,
			EditHelp:  "This is the location of the message author.  It is required.",
		},
		&basemsg.Field{
			Label:     "From Name",
			PIFOTag:   "8c.",
			Value:     &bf.FromName,
			Compare:   common.CompareText,
			PDFMap:    pdf.FromName,
			EditWidth: 34,
			EditHelp:  "This is the name of the message author.  It is optional and rarely provided.",
		},
		&basemsg.Field{
			Label:     "From Contact Info",
			PIFOTag:   "8d.",
			Value:     &bf.FromContact,
			Compare:   common.CompareText,
			PDFMap:    pdf.FromContact,
			EditWidth: 29,
			EditHelp:  "This is contact information (phone number, email, etc.) for the message author.  It is optional and rarely provided.",
		},
	)
	bm.FOriginMsgID = &bf.OriginMsgID
	bm.FHandling = &bf.Handling
	bm.FToICSPosition = &bf.ToICSPosition
	bm.FToLocation = &bf.ToLocation
}
func (bf *BaseForm) AddFooterFields(bm *basemsg.BaseMessage, pdf *BaseFormPDFMaps) {
	bm.Fields = append(bm.Fields,
		&basemsg.Field{
			Label:     "Operator: Relay Received",
			PIFOTag:   "OpRelayRcvd",
			Value:     &bf.OpRelayRcvd,
			Compare:   common.CompareText,
			PDFMap:    pdf.OpRelayRcvd,
			EditWidth: 36,
			EditHelp:  "This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.",
		},
		&basemsg.Field{
			Label:     "Operator: Relay Sent",
			PIFOTag:   "OpRelaySent",
			Value:     &bf.OpRelaySent,
			Compare:   common.CompareText,
			PDFMap:    pdf.OpRelaySent,
			EditWidth: 36,
			EditHelp:  "This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.",
		},
		&basemsg.Field{
			Label:      "Operator: Name",
			PIFOTag:    "OpName",
			Value:      &bf.OpName,
			Presence:   basemsg.Required,
			PDFMap:     pdf.OpName,
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Operator: Call Sign",
			PIFOTag:    "OpCall",
			Value:      &bf.OpCall,
			Presence:   basemsg.Required,
			PDFMap:     pdf.OpCall,
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Operator",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(bf.OpCall, bf.OpName, " ")
			},
		},
		&basemsg.Field{
			Label:      "Operator: Date",
			PIFOTag:    "OpDate",
			Value:      &bf.OpDate,
			Presence:   basemsg.Required,
			PIFOValid:  basemsg.ValidDate,
			PDFMap:     pdf.OpDate,
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Operator: Time",
			PIFOTag:    "OpTime",
			Value:      &bf.OpTime,
			Presence:   basemsg.Required,
			PIFOValid:  basemsg.ValidTime,
			PDFMap:     pdf.OpTime,
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Operator: Date/Time",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(bf.OpDate, bf.OpTime, " ")
			},
		},
	)
	bm.FOpCall = &bf.OpCall
	bm.FOpName = &bf.OpName
}
