// Package sheltstat defines the Santa Clara County OA Shelter Status Form
// message type.
package sheltstat

import (
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type23 is the type definition for an OA shelter status form.
var Type23 = message.Type{
	Tag:     "SheltStat",
	HTML:    "form-oa-shelter-status.html",
	Version: "2.3",
	Name:    "OA shelter status form",
	Article: "an",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "19.", "32.", "30.", "31.", "33a.",
		"33b.", "34b.", "33c.", "33d.", "37a.", "37b.", "40a.", "40b.", "41.", "42.", "43a.", "43b.", "43c.", "44.", "45.", "46.",
		"50a.", "49a.", "50b.", "51a.", "51b.", "52a.", "52b.", "60.", "61.", "62a.", "62b.", "63a.", "63b.", "62c.", "70.", "71.",
		"OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type23, decode23, create23)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID: message.PDFMultiRenderer{
		&message.PDFTextRenderer{Page: 1, X: 223, Y: 60, R: 349, B: 77, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 468, Y: 31, R: 573, B: 47},
	},
	DestinationMsgID: &message.PDFTextRenderer{Page: 1, X: 446, Y: 60, R: 573, B: 77, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{Page: 1, X: 70, Y: 123, R: 138, B: 144, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTime:      &message.PDFTextRenderer{Page: 1, X: 201, Y: 123, R: 236, B: 144, Style: message.PDFTextStyle{VAlign: "baseline"}},
	Handling: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
		"IMMEDIATE": {306, 133},
		"PRIORITY":  {408, 133},
		"ROUTINE":   {494, 133},
	}},
	ToICSPosition:   &message.PDFTextRenderer{Page: 1, X: 127, Y: 145, R: 303, B: 162, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{Page: 1, X: 127, Y: 163, R: 303, B: 180, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{Page: 1, X: 127, Y: 181, R: 303, B: 198, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{Page: 1, X: 127, Y: 199, R: 303, B: 216, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{Page: 1, X: 398, Y: 145, R: 573, B: 162, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{Page: 1, X: 398, Y: 163, R: 573, B: 180, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{Page: 1, X: 398, Y: 181, R: 573, B: 198, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{Page: 1, X: 398, Y: 199, R: 573, B: 216, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{Page: 2, X: 110, Y: 489, R: 321, B: 507, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{Page: 2, X: 358, Y: 489, R: 573, B: 507, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{Page: 2, X: 76, Y: 508, R: 250, B: 525, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{Page: 2, X: 302, Y: 508, R: 367, B: 525, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{Page: 2, X: 413, Y: 508, R: 473, B: 525, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{Page: 2, X: 539, Y: 508, R: 573, B: 525, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// SheltStat23 holds an OA shelter status form.
type SheltStat23 struct {
	message.BaseMessage
	baseform.BaseForm
	ReportType            string
	ShelterName           string
	ShelterType           string
	ShelterStatus         string
	ShelterAddress        string
	ShelterCityCode       string
	ShelterCity           string
	ShelterState          string
	ShelterZip            string
	Latitude              string
	Longitude             string
	Capacity              string
	Occupancy             string
	MealsServed           string
	NSSNumber             string
	PetFriendly           string
	BasicSafetyInspection string
	ATC20Inspection       string
	AvailableServices     string
	MOU                   string
	FloorPlan             string
	ManagedByCode         string
	ManagedBy             string
	ManagedByDetail       string
	PrimaryContact        string
	PrimaryPhone          string
	SecondaryContact      string
	SecondaryPhone        string
	TacticalCallSign      string
	RepeaterCallSign      string
	RepeaterInput         string
	RepeaterInputTone     string
	RepeaterOutput        string
	RepeaterOutputTone    string
	RepeaterOffset        string
	Comments              string
	RemoveFromList        string
}

func create23() message.Message {
	var f = make23()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "PRIORITY"
	return f
}

func make23() (f *SheltStat23) {
	const fieldCount = 63
	f = &SheltStat23{BaseMessage: message.BaseMessage{Type: &Type23}}
	f.BaseMessage.FSubject = &f.ShelterName
	f.BaseMessage.FBody = &f.Comments
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  message.Choices{"Update", "Complete"},
			Presence: message.Required,
			PIFOTag:  "19.",
			EditHelp: `This indicates whether the form should "Update" the previous status report for the shelter, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Update":   {120, 225},
				"Complete": {176, 225},
			}},
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter Name",
			Value:       &f.ShelterName,
			Presence:    message.Required,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 310, Y: 218, R: 573, B: 244},
			EditWidth:   44,
			EditHelp:    `This is the name of the shelter whose status is being reported.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Shelter Type",
			Value:    &f.ShelterType,
			Choices:  message.Choices{"Type 1", "Type 2", "Type 3", "Type 4"},
			Presence: f.requiredForComplete,
			PIFOTag:  "30.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Type 1": {195, 277},
				"Type 2": {276, 277},
				"Type 3": {357, 277},
				"Type 4": {443, 277},
			}},
			EditHelp: `This is the shelter type.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Shelter Status",
			Value:    &f.ShelterStatus,
			Choices:  message.Choices{"Open", "Closed", "Full"},
			Presence: f.requiredForComplete,
			PIFOTag:  "31.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Open":   {195, 295},
				"Closed": {276, 295},
				"Full":   {357, 295},
			}},
			EditHelp: `This indicates the status of the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter Address",
			Value:       &f.ShelterAddress,
			Presence:    f.requiredForComplete,
			PIFOTag:     "33a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 87, Y: 305, R: 573, B: 321},
			EditWidth:   75,
			EditHelp:    `This is the street address of the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewCalculatedField(&message.Field{
			Label:    "Shelter City Code",
			Value:    &f.ShelterCityCode,
			Presence: f.requiredForComplete,
			PIFOTag:  "33b.",
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter City",
			Value:       &f.ShelterCity,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
			Presence:    f.requiredForComplete,
			PIFOTag:     "34b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 87, Y: 323, R: 573, B: 339},
			TableValue:  message.TableOmit,
			EditWidth:   30,
			EditHelp:    `This is the name of the city in which the shelter is located.  It is required when "Report Type" is "Complete".`,
			EditApply: func(field *message.Field, s string) {
				f.ShelterCity = field.Choices.ToPIFO(s)
				if f.ShelterCity == "" || field.Choices.IsPIFO(f.ShelterCity) {
					f.ShelterCityCode = f.ShelterCity
				} else {
					f.ShelterCityCode = "Unincorporated"
				}
			},
			EditSkip: func(*message.Field) bool { return f.ShelterAddress == "" },
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter State",
			Value:       &f.ShelterState,
			Choices:     message.Choices{"CA"},
			Presence:    f.requiredForComplete,
			PIFOTag:     "33c.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 87, Y: 341, R: 573, B: 357},
			TableValue:  message.TableOmit,
			EditWidth:   12,
			EditHelp:    `This is the name (or two-letter abbreviation) of the state in which the shelter is located.  It is required when "Report Type" is "Complete".`,
			EditSkip:    func(*message.Field) bool { return f.ShelterCity == "" },
		}),
		message.NewTextField(&message.Field{
			Label:       "Shelter Zip",
			Value:       &f.ShelterZip,
			Presence:    f.requiredForComplete,
			PIFOTag:     "33d.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 87, Y: 359, R: 573, B: 375},
			TableValue:  message.TableOmit,
			EditWidth:   12,
			EditHelp:    `This is the shelter's ZIP code.  It is required when "Report Type" is "Complete".`,
			EditSkip:    func(*message.Field) bool { return f.ShelterState == "" },
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Shelter City",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(f.ShelterCity, message.SmartJoin(f.ShelterState, f.ShelterZip, "  "), ", ")
			},
		}),
		message.NewRealNumberField(&message.Field{
			Label:       "Latitude",
			Value:       &f.Latitude,
			PIFOTag:     "37a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 123, Y: 377, R: 302, B: 392, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue: func(*message.Field) string {
				if f.Longitude == "" {
					return f.Latitude
				}
				return ""
			},
			EditWidth: 30,
			EditHelp:  `This is the latitude of the shelter location, expressed in fractional degrees.`,
		}),
		message.NewRealNumberField(&message.Field{
			Label:       "Longitude",
			Value:       &f.Longitude,
			PIFOTag:     "37b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 399, Y: 377, R: 573, B: 392, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue: func(*message.Field) string {
				if f.Longitude == "" {
					return f.Latitude
				}
				return ""
			},
			EditWidth: 29,
			EditHelp:  `This is the longitude of the shelter location, expressed in fractional degrees.`,
			EditValid: func(field *message.Field) string {
				if f.Latitude != "" && f.Longitude == "" {
					return `The "Longitude" field must have a value when "Latitude" has a value.`
				}
				if f.Longitude != "" && !message.PIFORealNumberRE.MatchString(f.Longitude) {
					return `The "Longitude" field does not contain a valid number.`
				}
				return ""
			},
			EditSkip: func(*message.Field) bool { return f.Latitude == "" },
		}),
		message.NewAggregatorField(&message.Field{
			Label: "GPS Coordinates",
			TableValue: func(*message.Field) string {
				if f.Latitude != "" && f.Longitude != "" {
					return f.Latitude + ", " + f.Longitude
				}
				return ""
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Capacity",
			Value:       &f.Capacity,
			Presence:    f.requiredForComplete,
			PIFOTag:     "40a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 190, Y: 418, R: 573, B: 433},
			TableValue: func(*message.Field) string {
				if f.Occupancy == "" {
					return f.Capacity
				}
				return ""
			},
			EditWidth: 6,
			EditHelp:  `This is the number of people the shelter can accommodate.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Occupancy",
			Value:       &f.Occupancy,
			Presence:    f.requiredForComplete,
			PIFOTag:     "40b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 190, Y: 436, R: 573, B: 451},
			TableValue: func(*message.Field) string {
				if f.Occupancy != "" && f.Capacity != "" {
					return f.Occupancy + " out of " + f.Capacity
				}
				return f.Occupancy
			},
			EditWidth: 6,
			EditHelp:  `This is the number of people currently using the shelter.  It is required when "Report Type" is "Complete".`,
			EditSkip:  func(*message.Field) bool { return f.ShelterAddress == "" },
		}),
		message.NewTextField(&message.Field{
			Label:       "Meals Served",
			Value:       &f.MealsServed,
			PIFOTag:     "41.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 190, Y: 454, R: 573, B: 469, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   65,
			EditHelp:    `This is the number and/or description of meals served at the shelter in the last 24 hours.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "NSS Number",
			Value:       &f.NSSNumber,
			PIFOTag:     "42.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 190, Y: 472, R: 573, B: 487, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   65,
			EditHelp:    `This is the NSS number of the shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Pet Friendly",
			Value:   &f.PetFriendly,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "43a.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"checked": {195, 496},
				"false":   {231, 496},
			}},
			EditHelp: `This indicates whether the shelter can accept pets.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Basic Safety Inspection",
			Value:   &f.BasicSafetyInspection,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "43b.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"checked": {195, 514},
				"false":   {231, 514},
			}},
			EditHelp: `This indicates whether the shelter has had a basic safety inspection.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "ATC-20 Inspection",
			Value:   &f.ATC20Inspection,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "43c.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"checked": {195, 532},
				"false":   {231, 532},
			}},
			EditHelp: `This indicates whether the shelter has had an ATC-20 inspection.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Available Services",
			Value:       &f.AvailableServices,
			PIFOTag:     "44.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 42, Y: 555, R: 573, B: 642, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   85,
			EditHelp:    `This is a list of services available at the shelter.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "MOU",
			Value:       &f.MOU,
			PIFOTag:     "45.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 178, Y: 646, R: 573, B: 661, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   64,
			EditHelp:    `This indicates where and how the shelter's Memorandum of Understanding (MOU) was reported.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Floor Plan",
			Value:       &f.FloorPlan,
			PIFOTag:     "46.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 178, Y: 664, R: 573, B: 679, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   64,
			EditHelp:    `This indicates where and how the shelter's floor plan was reported.`,
		}),
		message.NewCalculatedField(&message.Field{
			Label:    "Managed By Code",
			Value:    &f.ManagedByCode,
			Presence: f.requiredForComplete,
			PIFOTag:  "50a.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Managed By",
			Value:    &f.ManagedBy,
			Choices:  message.Choices{"American Red Cross", "Private", "Community", "Government", "Other"},
			Presence: f.requiredForComplete,
			PIFOTag:  "49a.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"American Red Cross": {195, 80},
				"Private":            {339, 80},
				"Community":          {438, 80},
				"Government":         {195, 96},
				"Other":              {339, 96},
			}},
			EditWidth: 18,
			EditHelp:  `This indicates what type of entity is managing the shelter.  It is required when "Report Type" is "Complete".`,
			EditApply: func(field *message.Field, s string) {
				f.ManagedBy = field.Choices.ToPIFO(s)
				if f.ManagedBy == "" || field.Choices.IsPIFO(f.ManagedBy) {
					f.ManagedByCode = f.ManagedBy
				} else {
					f.ManagedByCode = "Other"
				}
			},
			EditValid: func(f *message.Field) string {
				if p := f.PresenceValid(); p != "" {
					return p
				}
				if *f.Value != "" && !f.Choices.IsPIFO(*f.Value) {
					return `The "Managed By" field should contain one of the recommended values ("American Red Cross", "Private", "Community", "Government", or "Other").  Other values can be successfully sent in a PackItForms message, but they cannot be rendered in the generated PDF.`
				}
				return ""
			},
		}),
		message.NewTextField(&message.Field{
			Label:       "Managed By Detail",
			Value:       &f.ManagedByDetail,
			PIFOTag:     "50b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 167, Y: 105, R: 573, B: 121},
			EditWidth:   65,
			EditHelp:    `This is additional detail about who is managing the shelter (particularly if "Managed By" is "Other").`,
			EditValid: func(*message.Field) string {
				if f.ManagedBy == "Other" && f.ManagedByDetail == "" {
					return `The "Managed By Detail" field is required when "Managed By" is "Other".`
				}
				return ""
			},
		}),
		message.NewTextField(&message.Field{
			Label:       "Primary Contact",
			Value:       &f.PrimaryContact,
			Presence:    f.requiredForComplete,
			PIFOTag:     "51a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 167, Y: 123, R: 573, B: 139},
			EditWidth:   65,
			EditHelp:    `This is the name of the primary contact person for the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Primary Phone",
			Value:       &f.PrimaryPhone,
			Presence:    f.requiredForComplete,
			PIFOTag:     "51b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 167, Y: 141, R: 573, B: 157},
			EditWidth:   65,
			EditHelp:    `This is the phone number of the primary contact person for the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Secondary Contact",
			Value:       &f.SecondaryContact,
			PIFOTag:     "52a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 167, Y: 159, R: 573, B: 175},
			EditWidth:   65,
			EditHelp:    `This is the name of the secondary contact person for the shelter.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Secondary Phone",
			Value:       &f.SecondaryPhone,
			PIFOTag:     "52b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 167, Y: 177, R: 573, B: 192},
			EditWidth:   65,
			EditHelp:    `This is the phone number of the secondary contact person for the shelter.`,
		}),
		message.NewTacticalCallSignField(&message.Field{
			Label:       "Tactical Call Sign",
			Value:       &f.TacticalCallSign,
			PIFOTag:     "60.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 132, Y: 216, R: 573, B: 232},
			EditWidth:   29,
			EditHelp:    `This is the tactical call sign assigned to the shelter for amateur radio communications.`,
		}),
		message.NewFCCCallSignField(&message.Field{
			Label:       "Repeater Call Sign",
			Value:       &f.RepeaterCallSign,
			PIFOTag:     "61.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 132, Y: 234, R: 573, B: 249},
			EditWidth:   29,
			EditHelp:    `This is the call sign of the amateur radio repeater that the shelter is monitoring for communications.`,
		}),
		message.NewFrequencyField(&message.Field{
			Label:       "Repeater Input",
			Value:       &f.RepeaterInput,
			PIFOTag:     "62a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 187, Y: 252, R: 311, B: 267, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditWidth:   20,
			EditHelp:    `This is the input frequency (in MHz) of the amateur radio repeater that the shelter is using for communications.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Repeater Input Tone",
			Value:       &f.RepeaterInputTone,
			PIFOTag:     "62b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 389, Y: 252, R: 573, B: 267, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditWidth:   30,
			EditHelp:    `This is the analog CTCSS tone, P25 NAC, DMR TS/TG/CC, or other access details required by the amateur radio repeater that the shelter is using for communications.`,
			EditSkip:    func(*message.Field) bool { return f.RepeaterInput == "" },
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Repeater Input",
			TableValue: func(*message.Field) string {
				return formatFreq(f.RepeaterInput, f.RepeaterInputTone)
			},
		}),
		message.NewFrequencyField(&message.Field{
			Label:       "Repeater Output",
			Value:       &f.RepeaterOutput,
			PIFOTag:     "63a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 187, Y: 271, R: 311, B: 285, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditWidth:   20,
			EditHelp:    `This is the output frequency (in MHz) of the amateur radio repeater that the shelter is using for communications.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Repeater Output Tone",
			Value:       &f.RepeaterOutputTone,
			PIFOTag:     "63b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 389, Y: 271, R: 573, B: 285, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditWidth:   30,
			EditHelp:    `This is the analog CTCSS tone, P25 NAC, DMR TS/TG/CC, or other access details for the transmission from the amateur radio repeater that the shelter is using for communications.`,
			EditSkip:    func(*message.Field) bool { return f.RepeaterOutput == "" },
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Repeater Output",
			TableValue: func(*message.Field) string {
				return formatFreq(f.RepeaterOutput, f.RepeaterOutputTone)
			},
		}),
		message.NewFrequencyOffsetField(&message.Field{
			Label:       "Repeater Offset",
			Value:       &f.RepeaterOffset,
			PIFOTag:     "62c.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 219, Y: 287, R: 311, B: 303},
			EditWidth:   15,
			EditHelp:    `This is the offset for the amateur radio repeater that the shelter is using for communications.  It can be either a number in MHz, a "+", or a "-".`,
			EditSkip:    func(*message.Field) bool { return (f.RepeaterInput == "") == (f.RepeaterOutput == "") },
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Comments",
			Value:       &f.Comments,
			PIFOTag:     "70.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 42, Y: 340, R: 573, B: 430, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   85,
			EditHelp:    `These are comments regarding the status of the shelter.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Remove from List",
			Value:   &f.RemoveFromList,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "71.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"checked": {195, 441},
				"false":   {231, 441},
			}},
			EditHelp: `This indicates whether the shelter should be removed from the receiver's list of shelters.`,
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update SheltStat23 fieldCount")
	}
	return f
}

func (f *SheltStat23) requiredForComplete() (message.Presence, string) {
	if f.ReportType == "Complete" {
		return message.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return message.PresenceOptional, ""
}

func decode23(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type23.HTML || form.FormVersion != Type23.Version {
		return nil
	}
	var df = make23()
	message.DecodeForm(form, df)
	return df
}

func formatFreq(freq, tone string) string {
	switch {
	case freq == "" && tone == "":
		return ""
	case freq != "" && tone == "":
		return freq
	case freq == "" && tone != "":
		return "Tone: " + tone
	default: // case freq != "" && tone != "":
		return freq + ", Tone: " + tone
	}
}
