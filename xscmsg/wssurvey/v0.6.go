// Package wssurvey defines the Windshield Survey Form message type.
package wssurvey

import (
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a windshield survey form.
var Type = message.Type{
	Tag:     "WSSurvey",
	HTML:    "form-windshield-survey.html",
	Version: "0.6",
	Name:    "windshield survey form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21.", "22.", "23.", "24.", "25.",
		"26.", "27.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall",
		"OpDate", "OpTime",
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

// WindshieldSurvey holds a shelter form.
type WindshieldSurvey struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction        string
	Team                string
	Location            string
	Item                string
	BuildingType        string
	NumberOfStories     string
	DamageCategory      string
	OtherDamageObserved string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "ROUTINE"
	return f
}

func makeF() *WindshieldSurvey {
	const fieldCount = 30
	f := WindshieldSurvey{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.Location
	f.FBody = &f.OtherDamageObserved
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
			EditHelp:    `This is the name of the jurisdiction where the windshield survey was performed.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Team",
			Value:       &f.Team,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the team (or individual) who performed the windshield survey.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Location",
			Value:       &f.Location,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the location that was surveyed.  It should be unique among all windshield survey reports in the jurisdiction.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Item",
			Value:       &f.Item,
			Presence:    message.Required,
			PIFOTag:     "23.",
			Choices:     message.Choices{"Building", "Road", "Other"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates the type of property surveyed.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Building Type",
			Value:   &f.BuildingType,
			PIFOTag: "24.",
			Choices: message.ChoicePairs{
				"Single", "Single Family Home",
				"Townhouse", "Townhouse or Condo",
				"Apartment", "Apartment",
				"Mobile", "Mobile Home or Trailer",
				"Business", "Business",
				"Other", "Other",
			},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates the type of building surveyed, if applicable.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Number of Stories",
			Value:       &f.NumberOfStories,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   2,
			EditHelp:    `This is the number of stories in the building surveyed, if applicable.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Categorization",
			Value:       &f.DamageCategory,
			Presence:    message.Required,
			PIFOTag:     "26.",
			Choices:     message.Choices{"Affected", "Minor", "Major", "Destroyed"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates the level of damage to the property surveyed.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Other Damage Observed",
			Value:       &f.OtherDamageObserved,
			PIFOTag:     "27.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This describes any other damage observed during the survey.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update WindshieldSurvey fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *WindshieldSurvey

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}
