// Package roadcl defines the Road Closure Form message type.
package roadcl

import (
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a road closure form.
var Type = message.Type{
	Tag:     "RoadCl",
	HTML:    "form-road-closure.html",
	Version: "0.6",
	Name:    "road closure form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21.", "22.", "23.", "24.", "30d.",
		"30t.", "31d.", "31t.", "OpRelayRcvd", "OpRelaySent", "OpName",
		"OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID: &message.PDFMultiRenderer{
		&message.PDFTextRenderer{X: 223, Y: 50, R: 348, B: 67, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 420, Y: 36, R: 574, B: 48},
	},
	DestinationMsgID: &message.PDFTextRenderer{X: 452, Y: 50, R: 574, B: 67, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 69, Y: 87, R: 128, B: 104, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	MessageTime:      &message.PDFTextRenderer{X: 162, Y: 87, R: 203, B: 104, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	Handling: &message.PDFRadioRenderer{Radius: 5, Points: map[string][]float64{
		"IMMEDIATE": {277, 96},
		"PRIORITY":  {388, 96},
		"ROUTINE":   {493, 96},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132, Y: 106, R: 292, B: 123, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132, Y: 125, R: 292, B: 142, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132, Y: 144, R: 292, B: 161, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132, Y: 163, R: 292, B: 181, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 382, Y: 106, R: 572, B: 123, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 382, Y: 125, R: 572, B: 142, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 382, Y: 144, R: 572, B: 161, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 382, Y: 163, R: 572, B: 181, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 109, Y: 711, R: 320, B: 728, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 356, Y: 711, R: 574, B: 728, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 76, Y: 730, R: 249, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 301, Y: 730, R: 366, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 401, Y: 730, R: 479, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 542, Y: 730, R: 574, B: 747, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// RoadClosure holds a road closure form.
type RoadClosure struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction     string
	Road             string
	Location         string
	Details          string
	Status           string
	ClosureStartDate string
	ClosureStartTime string
	ClosureEndDate   string
	ClosureEndTime   string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "ROUTINE"
	f.ToLocation = "County EOC"
	return f
}

func makeF() *RoadClosure {
	const fieldCount = 33
	f := RoadClosure{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.Road
	f.FBody = &f.Details
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "20.",
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   21,
			EditHelp:    `This is the name of the jurisdiction originating the road closure.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Road/Intersection",
			Value:       &f.Road,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the road or intersection with the closure.  It should be unique among all road closure reports from this jurisdiction.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Location",
			Value:       &f.Location,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the detailed location of the closure.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Details",
			Value:       &f.Details,
			Presence:    message.Required,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This field has details of the closure.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Status",
			Value:       &f.Status,
			Presence:    message.Required,
			PIFOTag:     "24.",
			Choices:     message.Choices{"Planned Closure", "Closed", "Reopened"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This identifies the status of the closure.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Closure Start Date",
			Value:       &f.ClosureStartDate,
			PIFOTag:     "30d.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the date when the closure started or will start.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Closure Start Time",
			Value:       &f.ClosureStartTime,
			PIFOTag:     "30t.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the time of day when the closure started or will start.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			EditHelp: `This is the date and time when the closure started or will start, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ClosureStartDate, &f.ClosureStartTime),
		message.NewDateField(true, &message.Field{
			Label:       "Closure End Date",
			Value:       &f.ClosureEndDate,
			PIFOTag:     "31d.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the date when the closure ended or will end.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Closure End Time",
			Value:       &f.ClosureEndTime,
			PIFOTag:     "31t.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the time of day when the closure ended or will end.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			EditHelp: `This is the date and time when the closure ended or will end, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ClosureEndDate, &f.ClosureEndTime),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update RoadClosure fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *RoadClosure

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}
