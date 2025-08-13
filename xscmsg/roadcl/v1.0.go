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
	Version: "1.0",
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
	OriginMsgID:      &message.PDFTextRenderer{X: 224.52, Y: 62.04, W: 136.92, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 464.52, Y: 62.04, W: 96.24, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 72.72, Y: 116.76, W: 47.04, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 189.84, Y: 116.76, W: 37.92, H: 20, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {299.30, 126.57},
		"PRIORITY":  {401.42, 126.57},
		"ROUTINE":   {487.46, 126.57},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132.84, Y: 143.16, W: 144.48, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132.84, Y: 169.56, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132.84, Y: 195.96, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132.84, Y: 222.36, W: 144.48, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 375.84, Y: 143.16, W: 183.60, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 375.84, Y: 169.56, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 375.84, Y: 195.96, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 375.84, Y: 222.36, W: 183.60, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 111.00, Y: 581.40, W: 204.48, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 356.04, Y: 581.40, W: 204.72, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 77.04, Y: 600.24, W: 166.68, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 302.04, Y: 600.24, W: 58.20, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 402.00, Y: 600.24, W: 58.80, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 526.80, Y: 600.24, W: 33.96, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 134.16, Y: 267.12, W: 423.28, H: 11.64},
			EditHelp:    `This is the name of the jurisdiction originating the road closure.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Road/Intersection",
			Value:       &f.Road,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 134.16, Y: 285.60, W: 425.28, H: 15.48},
			EditWidth:   89,
			EditHelp:    `This is the road or intersection with the closure.  It should be unique among all road closure reports from this jurisdiction.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Location",
			Value:       &f.Location,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 319.92, W: 514.80, H: 62.40, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This is the detailed location of the closure.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Details",
			Value:       &f.Details,
			Presence:    message.Required,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 401.04, W: 514.80, H: 62.28, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This field has details of the closure.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Status",
			Value:    &f.Status,
			Presence: message.Required,
			PIFOTag:  "24.",
			Choices:  message.Choices{"Planned Closure", "Closed", "Reopened"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Planned Closure": {120.86, 475.29},
				"Closed":          {264.86, 475.29},
				"Reopened":        {372.86, 475.29},
			}},
			EditHelp: `This identifies the status of the closure.  It is required.`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Closure Start Date",
			Value:       &f.ClosureStartDate,
			PIFOTag:     "30d.",
			PDFRenderer: &message.PDFTextRenderer{X: 133.92, Y: 515.16, W: 125.28, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the closure started or will start.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Closure Start Time",
			Value:       &f.ClosureStartTime,
			PIFOTag:     "30t.",
			PDFRenderer: &message.PDFTextRenderer{X: 357.00, Y: 515.16, W: 202.44, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of day when the closure started or will start.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Closure Start",
			EditHelp: `This is the date and time when the closure started or will start, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ClosureStartDate, &f.ClosureStartTime),
		message.NewDateField(true, &message.Field{
			Label:       "Closure End Date",
			Value:       &f.ClosureEndDate,
			PIFOTag:     "31d.",
			PDFRenderer: &message.PDFTextRenderer{X: 133.92, Y: 532.80, W: 125.28, H: 14.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the closure ended or will end.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Closure End Time",
			Value:       &f.ClosureEndTime,
			PIFOTag:     "31t.",
			PDFRenderer: &message.PDFTextRenderer{X: 357.00, Y: 532.80, W: 202.44, H: 14.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of day when the closure ended or will end.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Closure End",
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
