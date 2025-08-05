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
	OriginMsgID:      &message.PDFTextRenderer{X: 225.36, Y: 62.04, W: 112.80, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 444.00, Y: 62.04, W: 116.76, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 72.72, Y: 119.40, W: 47.04, H: 20, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	MessageTime:      &message.PDFTextRenderer{X: 189.84, Y: 119.40, W: 42.48, H: 20, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {298.44, 123.93},
		"PRIORITY":  {400.23, 123.93},
		"ROUTINE":   {485.70, 123.93},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132.84, Y: 147.72, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132.84, Y: 174.12, W: 144.48, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132.84, Y: 200.52, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132.84, Y: 226.92, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 375.84, Y: 147.72, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 375.84, Y: 174.12, W: 183.60, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 375.84, Y: 200.52, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 375.84, Y: 226.92, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 111.00, Y: 573.72, W: 204.48, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 356.04, Y: 573.72, W: 204.72, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 77.04, Y: 592.68, W: 166.68, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 302.04, Y: 592.68, W: 58.20, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 402.00, Y: 592.68, W: 59.88, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 528.00, Y: 592.68, W: 32.76, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// WindshieldSurvey holds a shelter form.
type WindshieldSurvey struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction        string
	Team                string
	Location            string
	ItemBuilding        string
	ItemRoad            string
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
			PDFRenderer: &message.PDFTextRenderer{X: 103.64, Y: 271.53, W: 455.80, H: 12.03},
			EditHelp:    `This is the name of the jurisdiction where the windshield survey was performed.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Team",
			Value:       &f.Team,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 103.64, Y: 290.04, W: 455.80, H: 11.04},
			EditWidth:   95,
			EditHelp:    `This is the name of the team (or individual) who performed the windshield survey.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Location",
			Value:       &f.Location,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 103.64, Y: 308.04, W: 455.80, H: 10.20},
			EditWidth:   95,
			EditHelp:    `This is the location that was surveyed.  It should be unique among all windshield survey reports in the jurisdiction.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Item: Building",
			Value:       &f.ItemBuilding,
			PIFOTag:     "23a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {151.79, 325.49}}},
			EditHelp:    `This indicates that the item surveyed was a building.  At least one item type is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Item: Road",
			Value:       &f.ItemRoad,
			PIFOTag:     "23b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {259.79, 325.49}}},
			EditHelp:    `This indicates that the item surveyed was a road.  At least one item type is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label: "Building Type",
			Value: &f.BuildingType,
			Presence: func() (message.Presence, string) {
				if f.ItemBuilding != "" {
					return message.PresenceRequired, `"Item: Building" is checked`
				}
				return message.PresenceNotAllowed, `"Item: Building" is not checked`
			},
			PIFOTag: "24.",
			Choices: message.ChoicePairs{
				"Single", "Single Family Home",
				"Townhouse", "Townhouse or Condo",
				"Apartment", "Apartment",
				"Mobile", "Mobile Home or Trailer",
				"Business", "Business",
				"Other", "Other",
			},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Single":    {156.86, 347.85},
				"Townhouse": {156.86, 363.33},
				"Apartment": {156.86, 378.93},
				"Mobile":    {336.86, 347.85},
				"Business":  {336.86, 363.33},
				"Other":     {336.86, 378.93},
			}},
			EditHelp: `This indicates the type of building surveyed, if applicable.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Damage Categorization",
			Value:    &f.DamageCategory,
			Presence: message.Required,
			PIFOTag:  "26.",
			Choices:  message.Choices{"Affected", "Minor", "Major", "Destroyed"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Affected":  {192.86, 396.49},
				"Minor":     {192.86, 411.97},
				"Major":     {372.86, 396.49},
				"Destroyed": {372.86, 411.97},
			}},
			EditHelp: `This indicates the level of damage to the property surveyed.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Other Damage Observed",
			Value:       &f.OtherDamageObserved,
			PIFOTag:     "27.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 434.88, W: 514.80, H: 90.60, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This describes any other damage observed during the survey.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update WindshieldSurvey fieldCount")
	}
	return &f
}

func (f *WindshieldSurvey) atLeastOneItem() (message.Presence, string) {
	if f.ItemRoad != "" {
		return message.PresenceOptional, ""
	}
	return message.PresenceRequired, "no other item type is checked"
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
