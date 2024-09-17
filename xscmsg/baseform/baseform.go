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

type BaseFormPDF struct {
	OriginMsgID      message.PDFRenderer
	DestinationMsgID message.PDFRenderer
	MessageDate      message.PDFRenderer
	MessageTime      message.PDFRenderer
	Handling         message.PDFRenderer
	ToICSPosition    message.PDFRenderer
	ToLocation       message.PDFRenderer
	ToName           message.PDFRenderer
	ToContact        message.PDFRenderer
	FromICSPosition  message.PDFRenderer
	FromLocation     message.PDFRenderer
	FromName         message.PDFRenderer
	FromContact      message.PDFRenderer
	OpRelayRcvd      message.PDFRenderer
	OpRelaySent      message.PDFRenderer
	OpName           message.PDFRenderer
	OpCall           message.PDFRenderer
	OpDate           message.PDFRenderer
	OpTime           message.PDFRenderer
}
type BaseFormPDFMaps = BaseFormPDF // temporary

var RoutingSlipPDFRenderers = BaseFormPDF{
	OriginMsgID:      &message.PDFTextRenderer{X: 223, Y: 65, W: 136, H: 18, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 454, Y: 65, W: 120, H: 18, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 74, Y: 126, W: 67, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 211, Y: 126, W: 34, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Points: map[string][]float64{
		"IMMEDIATE": {313.5, 134.5},
		"PRIORITY":  {413, 134.5},
		"ROUTINE":   {497, 134.5},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132, Y: 146, W: 170, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132, Y: 166, W: 170, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132, Y: 186, W: 170, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132, Y: 206, W: 170, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 404, Y: 146, W: 169, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 404, Y: 166, W: 169, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 404, Y: 186, W: 169, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 404, Y: 206, W: 169, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 110, Y: 369, W: 211, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 356, Y: 369, W: 217, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 76, Y: 388, W: 174, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 302, Y: 388, W: 65, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 403, Y: 388, W: 71, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 540, Y: 388, W: 33, H: 16, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

func (bf *BaseForm) AddHeaderFields(bm *message.BaseMessage, pdf *BaseFormPDF) {
	bm.Fields = append(bm.Fields,
		message.NewMessageNumberField(&message.Field{
			Label:       "Origin Message Number",
			Value:       &bf.OriginMsgID,
			Presence:    message.Required,
			PIFOTag:     "MsgNo",
			PDFRenderer: pdf.OriginMsgID,
			EditHelp:    "This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.",
		}),
		message.NewMessageNumberField(&message.Field{
			Label:       "Destination Message Number",
			Value:       &bf.DestinationMsgID,
			PIFOTag:     "DestMsgNo",
			PDFRenderer: pdf.DestinationMsgID,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Message Date",
			Value:       &bf.MessageDate,
			Presence:    message.Required,
			PIFOTag:     "1a.",
			PDFRenderer: pdf.MessageDate,
			EditHelp:    "This is the date the message was written, in MM/DD/YYYY format.  It is required.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Message Time",
			Value:       &bf.MessageTime,
			Presence:    message.Required,
			PIFOTag:     "1b.",
			PDFRenderer: pdf.MessageTime,
			EditHelp:    "This is the time the message was written, in HH:MM format (24-hour clock).  It is required.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Message Date/Time",
			Presence: message.Required,
			EditHelp: "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
		}, &bf.MessageDate, &bf.MessageTime),
		message.NewRestrictedField(&message.Field{
			Label:       "Handling",
			Value:       &bf.Handling,
			Choices:     message.Choices{"ROUTINE", "PRIORITY", "IMMEDIATE"},
			Presence:    message.Required,
			PIFOTag:     "5.",
			PDFRenderer: pdf.Handling,
			EditHelp:    `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "To ICS Position",
			Value:       &bf.ToICSPosition,
			Presence:    message.Required,
			PIFOTag:     "7a.",
			PDFRenderer: pdf.ToICSPosition,
			EditWidth:   30,
			EditHelp:    "This is the ICS position to which the message is addressed.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:       "To Location",
			Value:       &bf.ToLocation,
			Presence:    message.Required,
			PIFOTag:     "7b.",
			PDFRenderer: pdf.ToLocation,
			EditWidth:   32,
			EditHelp:    "This is the location of the recipient ICS position.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:       "To Name",
			Value:       &bf.ToName,
			PIFOTag:     "7c.",
			PDFRenderer: pdf.ToName,
			EditWidth:   34,
			EditHelp:    "This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.",
		}),
		message.NewTextField(&message.Field{
			Label:       "To Contact Info",
			Value:       &bf.ToContact,
			PIFOTag:     "7d.",
			PDFRenderer: pdf.ToContact,
			EditWidth:   29,
			EditHelp:    "This is contact information (phone number, email, etc.) for the receipient.  It is optional and rarely provided.",
		}),
		message.NewTextField(&message.Field{
			Label:       "From ICS Position",
			Value:       &bf.FromICSPosition,
			Presence:    message.Required,
			PIFOTag:     "8a.",
			PDFRenderer: pdf.FromICSPosition,
			EditWidth:   30,
			EditHelp:    "This is the ICS position of the message author.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:       "From Location",
			Value:       &bf.FromLocation,
			Presence:    message.Required,
			PIFOTag:     "8b.",
			PDFRenderer: pdf.FromLocation,
			EditWidth:   32,
			EditHelp:    "This is the location of the message author.  It is required.",
		}),
		message.NewTextField(&message.Field{
			Label:       "From Name",
			Value:       &bf.FromName,
			PIFOTag:     "8c.",
			PDFRenderer: pdf.FromName,
			EditWidth:   34,
			EditHelp:    "This is the name of the message author.  It is optional and rarely provided.",
		}),
		message.NewTextField(&message.Field{
			Label:       "From Contact Info",
			Value:       &bf.FromContact,
			PIFOTag:     "8d.",
			PDFRenderer: pdf.FromContact,
			EditWidth:   29,
			EditHelp:    "This is contact information (phone number, email, etc.) for the message author.  It is optional and rarely provided.",
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
func (bf *BaseForm) AddFooterFields(bm *message.BaseMessage, pdf *BaseFormPDF) {
	bm.Fields = append(bm.Fields,
		message.NewTextField(&message.Field{
			Label:       "Operator: Relay Received",
			Value:       &bf.OpRelayRcvd,
			PIFOTag:     "OpRelayRcvd",
			PDFRenderer: pdf.OpRelayRcvd,
			Compare:     message.CompareNone,
			EditWidth:   36,
			EditHelp:    "This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.",
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Relay Sent",
			Value:       &bf.OpRelaySent,
			PIFOTag:     "OpRelaySent",
			PDFRenderer: pdf.OpRelaySent,
			Compare:     message.CompareNone,
			EditWidth:   36,
			EditHelp:    "This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.",
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Name",
			Value:       &bf.OpName,
			Presence:    message.Required,
			PIFOTag:     "OpName",
			PDFRenderer: pdf.OpName,
			TableValue:  message.TableOmit,
			Compare:     message.CompareNone,
		}),
		message.NewTextField(&message.Field{
			Label:       "Operator: Call Sign",
			Value:       &bf.OpCall,
			Presence:    message.Required,
			PIFOTag:     "OpCall",
			PDFRenderer: pdf.OpCall,
			TableValue:  message.TableOmit,
			Compare:     message.CompareNone,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Operator",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(bf.OpCall, bf.OpName, " ")
			},
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Operator: Date",
			Value:       &bf.OpDate,
			Presence:    message.Required,
			PIFOTag:     "OpDate",
			PDFRenderer: pdf.OpDate,
			Compare:     message.CompareNone,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Operator: Time",
			Value:       &bf.OpTime,
			Presence:    message.Required,
			PIFOTag:     "OpTime",
			PDFRenderer: pdf.OpTime,
			Compare:     message.CompareNone,
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
