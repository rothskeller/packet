// Package dmgasmt defines the Damage Assessment Form message type.
package dmgasmt

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a damage assessment form.
var Type = message.Type{
	Tag:     "DmgAsmt",
	HTML:    "form-damage-assessment.html",
	Version: "0.6",
	Name:    "damage assessment form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21.", "22.", "23.", "24.", "25.",
		"26.", "27a.", "27c.", "27b.", "27d.", "28.", "29.", "30.",
		"31.", "32.", "33.", "34.", "35.", "OpRelayRcvd", "OpRelaySent",
		"OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID:      &message.PDFTextRenderer{X: 224.52, Y: 62.04, W: 75.24, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 405.60, Y: 62.04, W: 155.16, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 72.72, Y: 114.24, W: 53.04, H: 15.24, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 190.44, Y: 114.24, W: 37.32, H: 15.24, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {299.30, 121.65},
		"PRIORITY":  {401.42, 121.65},
		"ROUTINE":   {487.46, 121.65},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132.84, Y: 136.32, W: 144.48, H: 15.12, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132.84, Y: 158.40, W: 144.48, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132.84, Y: 181.32, W: 144.48, H: 15.12, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132.84, Y: 203.40, W: 144.48, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 375.84, Y: 136.32, W: 183.60, H: 15.12, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 375.84, Y: 158.40, W: 183.60, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 375.84, Y: 181.32, W: 183.60, H: 15.12, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 375.84, Y: 203.40, W: 183.60, H: 15.96, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 111.00, Y: 615.96, W: 205.56, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 357.12, Y: 615.96, W: 208.20, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 77.04, Y: 634.80, W: 155.28, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 290.52, Y: 634.80, W: 71.04, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 403.32, Y: 634.80, W: 40.44, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 513.84, Y: 634.80, W: 51.48, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// DamageAssessment holds a damage assessment form.
type DamageAssessment struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction         string
	IncidentName         string
	Address              string
	UnitSuite            string
	TypeOfStructure      string
	NumberOfStories      string
	OwnRent              string
	DamageFlooding       string
	DamageStructural     string
	DamageExterior       string
	DamageOther          string
	BasementType         string
	DamageClassification string
	Tag                  string
	Insurance            string
	EstimatedDamage      string
	Comments             string
	ContactName          string
	ContactPhone         string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "ROUTINE"
	return f
}

func makeF() *DamageAssessment {
	const fieldCount = 42
	f := DamageAssessment{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.Address
	f.FBody = &f.Comments
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "20.",
			PDFRenderer: &message.PDFTextRenderer{X: 104.70, Y: 245.04, W: 190.74, H: 11.28, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the name of the jurisdiction originating the damage report.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Incident Name",
			Value:       &f.IncidentName,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 383.64, Y: 244.80, W: 175.80, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   37,
			EditHelp:    `This is the name of the incident that caused the damage being reported.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Address",
			Value:       &f.Address,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 104.70, Y: 263.28, W: 190.74, H: 11.40, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   22,
			EditHelp:    `This is the street address of the property whose damage is being reported.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Unit/Suite",
			Value:       &f.UnitSuite,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{X: 383.64, Y: 263.28, W: 175.80, H: 11.40, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   37,
			EditHelp:    `This is the unit number or suite number of the property whose damage is being reported, if applicable.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Type of Structure",
			Value:    &f.TypeOfStructure,
			Presence: message.Required,
			PIFOTag:  "24.",
			Choices:  message.Choices{"Single Family", "Multi-Family", "Mobile Home", "Business", "Non-Profit Orgs", "Outbuilding"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Single Family":   {156.86, 286.17},
				"Multi-Family":    {156.86, 301.77},
				"Mobile Home":     {264.86, 286.17},
				"Business":        {264.86, 301.77},
				"Non-Profit Orgs": {372.86, 286.17},
				"Outbuilding":     {372.86, 301.77},
			}},
			EditHelp: `This is the type of structure that sustained damage.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Number of Stories",
			Value:       &f.NumberOfStories,
			Presence:    message.Required,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{X: 151.79, Y: 313.68, W: 20.00, H: 11.16},
			EditWidth:   4,
			EditHelp:    `This is the number of stories of the building that sustained damage.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Own/Rent",
			Value:   &f.OwnRent,
			PIFOTag: "26.",
			Choices: message.Choices{"Own", "Rent"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Own":  {156.86, 337.05},
				"Rent": {300.86, 337.05},
			}},
			EditHelp: `This indicates whether the property is owned or rented by its residents.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Flooding",
			Value:       &f.DamageFlooding,
			PIFOTag:     "27a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {151.79, 349.25}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained flooding damage.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Roof/Windows/Exterior",
			Value:       &f.DamageExterior,
			PIFOTag:     "27b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {151.79, 364.85}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained damage to the roof, windows, or other exterior features.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Structural",
			Value:       &f.DamageStructural,
			PIFOTag:     "27c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {331.79, 349.25}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained structural damage.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Other",
			Value:       &f.DamageOther,
			Presence:    f.atLeastOneDamage,
			PIFOTag:     "27d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {331.79, 364.85}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained other uncategorized damage.  (Describe in the Comments field.)`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Types of Damage",
			TableValue: func(*message.Field) string {
				var damage []string
				if f.DamageFlooding != "" {
					damage = append(damage, "Flooding")
				}
				if f.DamageExterior != "" {
					damage = append(damage, "Roof / Windows / Other Exterior")
				}
				if f.DamageStructural != "" {
					damage = append(damage, "Structural")
				}
				if f.DamageOther != "" {
					damage = append(damage, "Other")
				}
				return strings.Join(damage, ", ")
			},
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Basement Type",
			Value:    &f.BasementType,
			Presence: message.Required,
			PIFOTag:  "28.",
			Choices:  message.Choices{"Basement", "Basement and Crawlspace", "Crawlspace", "N/A"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Basement":                {156.86, 386.49},
				"Basement and Crawlspace": {300.86, 386.49},
				"Crawlspace":              {156.86, 401.97},
				"N/A":                     {300.86, 401.97},
			}},
			EditHelp: `This indicates whether the property has a basement and/or crawlspace.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Damage Classification",
			Value:   &f.DamageClassification,
			PIFOTag: "29.",
			Choices: message.Choices{"Destroyed", "Major", "Minor", "Affected", "No Visible Damage"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Destroyed":         {120.86, 431.97},
				"Major":             {120.86, 447.45},
				"Minor":             {228.86, 431.97},
				"Affected":          {228.86, 447.45},
				"No Visible Damage": {372.86, 431.97},
			}},
			EditHelp: `This indicates the level of damage to the property.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Tag",
			Value:   &f.Tag,
			PIFOTag: "30.",
			Choices: message.Choices{"Green", "Yellow", "Red"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Green":  {120.86, 464.73},
				"Yellow": {228.86, 464.73},
				"Red":    {336.86, 464.73},
			}},
			EditHelp: `This indicates what building department occupancy tag has been placed on the property, if any.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Insurance",
			Value:   &f.Insurance,
			PIFOTag: "31.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Yes": {120.86, 482.85},
				"No":  {192.86, 482.85},
			}},
			EditHelp: `This indicates whether the property damage is insured, if known.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Estimated Damage",
			Value:       &f.EstimatedDamage,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{X: 416.88, Y: 477.48, W: 40, H: 11.16},
			EditWidth:   8,
			EditHelp:    `This indicates the estimated damage cost to the property, in dollars, if known.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Comments",
			Value:       &f.Comments,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 506.16, W: 514.80, H: 50.64, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This gives additional comments about the damage to the property.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Contact Name",
			Value:       &f.ContactName,
			Presence:    message.Required,
			PIFOTag:     "34.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.64, Y: 563.64, W: 183.24, H: 12, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   999,
			EditHelp:    `This gives the name of the contact person for this damage report.  It is required.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Contact Phone",
			Value:       &f.ContactPhone,
			Presence:    message.Required,
			PIFOTag:     "35.",
			PDFRenderer: &message.PDFTextRenderer{X: 387.96, Y: 563.64, W: 171.48, H: 12, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This gives the phone number of the contact person for this damage report.  It is required.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update DamageAssessment fieldCount")
	}
	return &f
}

func (f *DamageAssessment) atLeastOneDamage() (message.Presence, string) {
	if f.DamageFlooding != "" || f.DamageExterior != "" || f.DamageStructural != "" {
		return message.PresenceOptional, ""
	}
	return message.PresenceRequired, "no other damage type is checked"
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *DamageAssessment

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}
