// Package shelter defines the Shelter Form message type.
package shelter

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a shelter form.
var Type = message.Type{
	Tag:     "Shelter",
	HTML:    "form-shelter.html",
	Version: "0.5",
	Name:    "shelter form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20.", "21.", "22.", "23a.", "23c.",
		"23e.", "23b.", "23d.", "24.", "25a.", "25c.", "25b.", "25d.",
		"26.", "27.", "28.", "29.", "30.", "31.", "32.", "40.", "41.",
		"42.", "43.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall",
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

// Shelter holds a shelter form.
type Shelter struct {
	message.BaseMessage
	baseform.BaseForm
	Jurisdiction              string
	ShelterName               string
	Location                  string
	TypeCongregate            string
	TypeNonCongregate         string
	TypeTEP                   string
	TypePetsAccepted          string
	TypeLargeAnimals          string
	ShelterStatus             string
	ResourceWarming           string
	ResourceCooling           string
	ResourceFoodDrink         string
	ResourceCharging          string
	Notes                     string
	ADACompliant              string
	PetsAllowed               string
	MaximumCapacity           string
	TotalRegistered           string
	CotNumbers                string
	AFNConsiderations         string
	ManagedBy                 string
	PrimaryContactName        string
	PrimaryContactPhoneNumber string
	PrimaryContactEmail       string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "PRIORITY"
	return f
}

func makeF() *Shelter {
	const fieldCount = 48
	f := Shelter{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.ShelterName
	f.FBody = &f.Notes
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
			EditHelp:    `This is the name of the jurisdiction originating the shelter.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter Name",
			Value:       &f.ShelterName,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the shelter being described.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Location",
			Value:       &f.Location,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the location of the shelter.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Congregate",
			Value:       &f.TypeCongregate,
			Presence:    f.atLeastOneType,
			PIFOTag:     "23a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a congregate shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Non-Congregate",
			Value:       &f.TypeNonCongregate,
			PIFOTag:     "23b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a non-congregate shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Temporary Evacuation Point",
			Value:       &f.TypeTEP,
			PIFOTag:     "23c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a temporary evacuation point.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Pets Accepted",
			Value:       &f.TypePetsAccepted,
			PIFOTag:     "23d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter accepts pets.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Large Animals",
			Value:       &f.TypeLargeAnimals,
			PIFOTag:     "23e.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter accepts large animals.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Type",
			TableValue: func(*message.Field) string {
				var stype []string
				if f.TypeCongregate != "" {
					stype = append(stype, "Congregate")
				}
				if f.TypeNonCongregate != "" {
					stype = append(stype, "Non-Congregate")
				}
				if f.TypeTEP != "" {
					stype = append(stype, "Temporary Evacuation Point")
				}
				if f.TypePetsAccepted != "" {
					stype = append(stype, "Pets Accepted")
				}
				if f.TypeLargeAnimals != "" {
					stype = append(stype, "Large Animals")
				}
				return strings.Join(stype, ", ")
			},
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Shelter Status",
			Value:       &f.ShelterStatus,
			Presence:    message.Required,
			PIFOTag:     "24.",
			Choices:     message.Choices{"Setup", "Open", "Full", "Closed"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates the status of the shelter.  Setup means setup is in progress, not yet accepting guests.  This field is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Warming Center",
			Value:       &f.ResourceWarming,
			PIFOTag:     "25a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicate that the shelter is a warming center.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Cooling Center",
			Value:       &f.ResourceCooling,
			PIFOTag:     "25b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a cooling center.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Food & Drink",
			Value:       &f.ResourceFoodDrink,
			PIFOTag:     "25c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter offers food and drink.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Charging Stations",
			Value:       &f.ResourceCharging,
			PIFOTag:     "25d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter has charging stations.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Resources",
			TableValue: func(*message.Field) string {
				var resources []string
				if f.ResourceWarming != "" {
					resources = append(resources, "Warming Center")
				}
				if f.ResourceCooling != "" {
					resources = append(resources, "Cooling Center")
				}
				if f.ResourceFoodDrink != "" {
					resources = append(resources, "Food & Drink")
				}
				if f.ResourceCharging != "" {
					resources = append(resources, "Charging Stations")
				}
				return strings.Join(resources, ", ")
			},
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Notes",
			Value:       &f.Notes,
			PIFOTag:     "26.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This field contains notes about the shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "ADA Compliant?",
			Value:       &f.ADACompliant,
			PIFOTag:     "27.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the shelter is ADA compliant.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Pets Allowed?",
			Value:       &f.PetsAllowed,
			PIFOTag:     "28.",
			Choices:     message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates whether the shelter allows pets.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Maximum Capacity",
			Value:       &f.MaximumCapacity,
			Presence:    message.Required,
			PIFOTag:     "29.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   3,
			EditHelp:    `This is the maximum capacity of the shelter.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Total Registered",
			Value:       &f.TotalRegistered,
			Presence:    message.Required,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   3,
			EditHelp:    `This is the number of people currently registered in the shelter.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Cot Numbers",
			Value:       &f.CotNumbers,
			Presence:    message.Required,
			PIFOTag:     "31.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   3,
			EditHelp:    `This is the number of cots in the shelter.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "AFN Considerations",
			Value:       &f.AFNConsiderations,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This field describes shelter considerations relating to access and functional needs.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Managed By",
			Value:    &f.ManagedBy,
			Presence: message.Required,
			PIFOTag:  "40.",
			Choices: message.ChoicePairs{
				"ARC", "American Red Cross",
				"Government", "Government",
				"Community", "Community",
				"Private", "Private",
				"Other", "Other",
			},
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This indicates what type of entity the shelter is managed by.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Primary Contact Name",
			Value:       &f.PrimaryContactName,
			Presence:    message.Required,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the name of the primary contact for the shelter.  It is required.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Primary Contact Phone Number",
			Value:       &f.PrimaryContactPhoneNumber,
			Presence:    message.Required,
			PIFOTag:     "42.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditHelp:    `This is the phone number of the primary contact for the shelter.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Primary Contact Email",
			Value:       &f.PrimaryContactEmail,
			PIFOTag:     "43.",
			PDFRenderer: &message.PDFTextRenderer{X: 999, Y: 999, R: 999, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   999,
			EditHelp:    `This is the email address of the primary contact for the shelter.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update Shelter fieldCount")
	}
	return &f
}

func (f *Shelter) atLeastOneType() (message.Presence, string) {
	if f.TypeNonCongregate != "" || f.TypeTEP != "" || f.TypePetsAccepted != "" || f.TypeLargeAnimals != "" {
		return message.PresenceOptional, ""
	}
	return message.PresenceRequired, "no other shelter type is checked"
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *Shelter

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}
