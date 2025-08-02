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
	Version: "0.5",
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
			PDFRenderer: &message.PDFTextRenderer{X: 190, Y: 184, R: 373, B: 206, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   21,
			EditHelp:    `This is the name of the jurisdiction originating the damage report.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Incident Name",
			Value:       &f.IncidentName,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the incident that caused the damage being reported.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Address",
			Value:       &f.Address,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the street address of the property whose damage is being reported.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Unit/Suite",
			Value:       &f.UnitSuite,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the unit number or suite number of the property whose damage is being reported, if applicable.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type of Structure",
			Value:       &f.TypeOfStructure,
			Presence:    message.Required,
			PIFOTag:     "24.",
			Choices:     message.Choices{"Single Family", "Multi-Family", "Mobile Home", "Business", "Non-Profit Orgs", "Outbuilding"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the type of structure that sustained damage.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Number of Stories",
			Value:       &f.NumberOfStories,
			Presence:    message.Required,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   2,
			EditHelp:    `This is the number of stories of the building that sustained damage.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Own/Rent",
			Value:       &f.OwnRent,
			PIFOTag:     "26.",
			Choices:     message.Choices{"Own", "Rent"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the property is owned or rented by its residents.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Flooding",
			Value:       &f.DamageFlooding,
			PIFOTag:     "27a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained flooding damage.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Roof / Windows / Other Exterior",
			Value:       &f.DamageExterior,
			PIFOTag:     "27b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained damage to the roof, windows, or other exterior features.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Structural",
			Value:       &f.DamageStructural,
			PIFOTag:     "27c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the property sustained structural damage.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Type: Other",
			Value:       &f.DamageOther,
			Presence:    f.atLeastOneDamage,
			PIFOTag:     "27d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
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
			Label:       "Basement Type",
			Value:       &f.BasementType,
			Presence:    message.Required,
			PIFOTag:     "28.",
			Choices:     message.Choices{"Basement", "Basement and Crawlspace", "Crawlspace", "N/A"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the property has a basement and/or crawlspace.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Damage Classification",
			Value:       &f.DamageClassification,
			PIFOTag:     "29.",
			Choices:     message.Choices{"Destroyed", "Major", "Minor", "Affected", "No Visible Damage"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates the level of damage to the property.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Tag",
			Value:       &f.Tag,
			PIFOTag:     "30.",
			Choices:     message.Choices{"Green", "Yellow", "Red"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates what building department occupancy tag has been placed on the property, if any.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Insurance",
			Value:       &f.Insurance,
			PIFOTag:     "31.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the property damage is insured, if known.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Estimated Damage",
			Value:       &f.EstimatedDamage,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   8,
			EditHelp:    `This indicates the estimated damage cost to the property, in dollars, if known.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Comments",
			Value:       &f.Comments,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This gives additional comments about the damage to the property.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Contact Name",
			Value:       &f.ContactName,
			Presence:    message.Required,
			PIFOTag:     "34.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This gives the name of the contact person for this damage report.  It is required.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Contact Phone",
			Value:       &f.ContactPhone,
			Presence:    message.Required,
			PIFOTag:     "35.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
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
