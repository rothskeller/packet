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
	Version: "0.6",
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
	OriginMsgID:      &message.PDFTextRenderer{X: 224.52, Y: 62.04, W: 102.12, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	DestinationMsgID: &message.PDFTextRenderer{X: 432.48, Y: 62.04, W: 132.84, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 72.56, Y: 116.76, W: 50, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{X: 187.08, Y: 116.76, W: 38.88, H: 12.48, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {297.50, 122.73},
		"PRIORITY":  {399.62, 122.73},
		"ROUTINE":   {485.66, 122.73},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 137.76, Y: 135.72, W: 165.84, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 137.76, Y: 154.56, W: 165.84, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 137.76, Y: 173.52, W: 165.84, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 137.76, Y: 192.36, W: 165.84, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 402.48, Y: 135.72, W: 161.52, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 402.48, Y: 154.56, W: 161.52, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 402.48, Y: 173.52, W: 161.52, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 402.48, Y: 192.36, W: 161.52, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{X: 116.76, Y: 654.72, W: 199.80, H: 12.24, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{X: 357.12, Y: 654.72, W: 208.20, H: 12.24, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{X: 81.48, Y: 673.56, W: 163.08, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{X: 302.88, Y: 673.56, W: 58.68, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{X: 403.32, Y: 673.56, W: 66.24, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{X: 539.64, Y: 673.56, W: 25.68, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 116.64, Y: 231.12, W: 447.36, H: 11.64},
			EditHelp:    `This is the name of the jurisdiction originating the shelter.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter Name",
			Value:       &f.ShelterName,
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.64, Y: 249.60, W: 447.36, H: 11.52},
			EditWidth:   94,
			EditHelp:    `This is the name of the shelter being described.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Location",
			Value:       &f.Location,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.64, Y: 268.08, W: 447.36, H: 11.52},
			EditWidth:   94,
			EditHelp:    `This is the location of the shelter.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Congregate",
			Value:       &f.TypeCongregate,
			Presence:    f.atLeastOneType,
			PIFOTag:     "23a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {119.03, 286.13}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a congregate shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Non-Congregate",
			Value:       &f.TypeNonCongregate,
			PIFOTag:     "23b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {119.03, 301.61}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a non-congregate shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Temporary Evacuation Point",
			Value:       &f.TypeTEP,
			PIFOTag:     "23c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {226.94, 286.13}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a temporary evacuation point.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Pets Accepted",
			Value:       &f.TypePetsAccepted,
			PIFOTag:     "23d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {226.94, 301.61}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter accepts pets.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Type: Large Animals",
			Value:       &f.TypeLargeAnimals,
			PIFOTag:     "23e.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {407.03, 286.13}}},
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
			Label:    "Shelter Status",
			Value:    &f.ShelterStatus,
			Presence: message.Required,
			PIFOTag:  "24.",
			Choices:  message.Choices{"Setup", "Open", "Full", "Closed"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 2, Points: map[string][]float64{
				"Setup":  {122.76, 322.17},
				"Open":   {122.76, 335.61},
				"Full":   {338.76, 322.17},
				"Closed": {338.76, 335.61},
			}},
			EditHelp: `This indicates the status of the shelter.  Setup means setup is in progress, not yet accepting guests.  This field is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Warming Center",
			Value:       &f.ResourceWarming,
			PIFOTag:     "25a.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {119.03, 346.01}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicate that the shelter is a warming center.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Cooling Center",
			Value:       &f.ResourceCooling,
			PIFOTag:     "25b.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {119.03, 361.61}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter is a cooling center.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Food & Drink",
			Value:       &f.ResourceFoodDrink,
			PIFOTag:     "25c.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {299.03, 346.01}}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates that the shelter offers food and drink.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Resources: Charging Stations",
			Value:       &f.ResourceCharging,
			PIFOTag:     "25d.",
			Choices:     message.Choices{"checked"},
			PDFRenderer: &message.PDFCheckRenderer{W: 10.15, H: 10.15, Points: map[string][]float64{"checked": {299.03, 361.61}}},
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
			PDFRenderer: &message.PDFTextRenderer{X: 48.00, Y: 389.16, W: 516.00, H: 50.76, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This field contains notes about the shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "ADA Compliant?",
			Value:   &f.ADACompliant,
			PIFOTag: "27.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 2, Points: map[string][]float64{
				"Yes": {158.76, 452.25},
				"No":  {230.76, 452.25},
			}},
			EditHelp: `This indicates whether the shelter is ADA compliant.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Pets Allowed?",
			Value:   &f.PetsAllowed,
			PIFOTag: "28.",
			Choices: message.Choices{"Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 2, Points: map[string][]float64{
				"Yes": {424.32, 452.25},
				"No":  {496.32, 452.25},
			}},
			EditHelp: `This indicates whether the shelter allows pets.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Maximum Capacity",
			Value:       &f.MaximumCapacity,
			Presence:    message.Required,
			PIFOTag:     "29.",
			PDFRenderer: &message.PDFTextRenderer{X: 142.32, Y: 465.24, W: 68.76, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   3,
			EditHelp:    `This is the maximum capacity of the shelter.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Total Registered",
			Value:       &f.TotalRegistered,
			Presence:    message.Required,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 305.88, Y: 465.24, W: 82.20, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   3,
			EditHelp:    `This is the number of people currently registered in the shelter.  It is required.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Cot Numbers",
			Value:       &f.CotNumbers,
			Presence:    message.Required,
			PIFOTag:     "31.",
			PDFRenderer: &message.PDFTextRenderer{X: 468.24, Y: 465.24, W: 95.76, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   3,
			EditHelp:    `This is the number of cots in the shelter.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "AFN Considerations",
			Value:       &f.AFNConsiderations,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{X: 48.00, Y: 494.40, W: 516.00, H: 37.20, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
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
			PDFRenderer: &message.PDFRadioRenderer{Radius: 2, Points: map[string][]float64{
				"ARC":        {158.76, 569.01},
				"Government": {302.76, 569.01},
				"Community":  {410.76, 569.01},
				"Private":    {158.76, 582.33},
				"Other":      {302.76, 582.33},
			}},
			EditHelp: `This indicates what type of entity the shelter is managed by.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Primary Contact Name",
			Value:       &f.PrimaryContactName,
			Presence:    message.Required,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{X: 158.28, Y: 593.28, W: 141.48, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   30,
			EditHelp:    `This is the name of the primary contact for the shelter.  It is required.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Primary Contact Phone",
			Value:       &f.PrimaryContactPhoneNumber,
			Presence:    message.Required,
			PIFOTag:     "42.",
			PDFRenderer: &message.PDFTextRenderer{X: 465.72, Y: 593.28, W: 98.28, H: 11.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the phone number of the primary contact for the shelter.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Primary Contact Email",
			Value:       &f.PrimaryContactEmail,
			PIFOTag:     "43.",
			PDFRenderer: &message.PDFTextRenderer{X: 158.28, Y: 611.88, W: 405.72, H: 11.52},
			EditWidth:   84,
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
