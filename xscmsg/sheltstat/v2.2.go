// Package sheltstat defines the Santa Clara County OA Shelter Status Form
// message type.
package sheltstat

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type22 is the type definition for an OA shelter status form.
var Type22 = message.Type{
	Tag:        "SheltStat",
	HTML:       "form-oa-shelter-status.html",
	Version:    "2.2",
	Name:       "OA shelter status form",
	Article:    "an",
	FieldOrder: Type23.FieldOrder,
}

func init() {
	message.Register(&Type22, decode22, nil)
}

// SheltStat22 holds an OA shelter status form.
type SheltStat22 struct {
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

func make22() (f *SheltStat22) {
	const fieldCount = 63
	f = &SheltStat22{BaseMessage: message.BaseMessage{Type: &Type22}}
	f.BaseMessage.FSubject = &f.ShelterName
	f.BaseMessage.FBody = &f.Comments
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, nil)
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  message.Choices{"Update", "Complete"},
			Presence: message.Required,
			PIFOTag:  "19.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Shelter Name",
			Value:    &f.ShelterName,
			Presence: message.Required,
			PIFOTag:  "32.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Shelter Type",
			Value:    &f.ShelterType,
			Choices:  message.Choices{"Type 1", "Type 2", "Type 3", "Type 4"},
			Presence: f.requiredForComplete,
			PIFOTag:  "30.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Shelter Status",
			Value:    &f.ShelterStatus,
			Choices:  message.Choices{"Open", "Closed", "Full"},
			Presence: f.requiredForComplete,
			PIFOTag:  "31.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Shelter Address",
			Value:    &f.ShelterAddress,
			Presence: f.requiredForComplete,
			PIFOTag:  "33a.",
		}),
		message.NewCalculatedField(&message.Field{
			Label:    "Shelter City Code",
			Value:    &f.ShelterCityCode,
			Presence: f.requiredForComplete,
			PIFOTag:  "33b.",
		}),
		message.NewTextField(&message.Field{
			Label:      "Shelter City",
			Value:      &f.ShelterCity,
			Choices:    message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
			Presence:   f.requiredForComplete,
			PIFOTag:    "34b.",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label:      "Shelter State",
			Value:      &f.ShelterState,
			Choices:    message.Choices{"CA"},
			Presence:   f.requiredForComplete,
			PIFOTag:    "33c.",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label:      "Shelter Zip",
			Value:      &f.ShelterZip,
			Presence:   f.requiredForComplete,
			PIFOTag:    "33d.",
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Shelter City",
			TableValue: func(*message.Field) string {
				return message.SmartJoin(f.ShelterCity, message.SmartJoin(f.ShelterState, f.ShelterZip, "  "), ", ")
			},
		}),
		message.NewRealNumberField(&message.Field{
			Label:   "Latitude",
			Value:   &f.Latitude,
			PIFOTag: "37a.",
			TableValue: func(*message.Field) string {
				if f.Longitude == "" {
					return f.Latitude
				}
				return ""
			},
		}),
		message.NewRealNumberField(&message.Field{
			Label:   "Longitude",
			Value:   &f.Longitude,
			PIFOTag: "37b.",
			TableValue: func(*message.Field) string {
				if f.Longitude == "" {
					return f.Latitude
				}
				return ""
			},
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
			Label:    "Capacity",
			Value:    &f.Capacity,
			Presence: f.requiredForComplete,
			PIFOTag:  "40a.",
			TableValue: func(*message.Field) string {
				if f.Occupancy == "" {
					return f.Capacity
				}
				return ""
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:    "Occupancy",
			Value:    &f.Occupancy,
			Presence: f.requiredForComplete,
			PIFOTag:  "40b.",
			TableValue: func(*message.Field) string {
				if f.Occupancy != "" && f.Capacity != "" {
					return f.Occupancy + " out of " + f.Capacity
				}
				return f.Occupancy
			},
		}),
		message.NewTextField(&message.Field{
			Label:   "Meals Served",
			Value:   &f.MealsServed,
			PIFOTag: "41.",
		}),
		message.NewTextField(&message.Field{
			Label:   "NSS Number",
			Value:   &f.NSSNumber,
			PIFOTag: "42.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Pet Friendly",
			Value:   &f.PetFriendly,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "43a.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Basic Safety Inspection",
			Value:   &f.BasicSafetyInspection,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "43b.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "ATC-20 Inspection",
			Value:   &f.ATC20Inspection,
			Choices: message.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag: "43c.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Available Services",
			Value:   &f.AvailableServices,
			PIFOTag: "44.",
		}),
		message.NewTextField(&message.Field{
			Label:   "MOU",
			Value:   &f.MOU,
			PIFOTag: "45.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Floor Plan",
			Value:   &f.FloorPlan,
			PIFOTag: "46.",
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
		}),
		message.NewTextField(&message.Field{
			Label:   "Managed By Detail",
			Value:   &f.ManagedByDetail,
			PIFOTag: "50b.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Primary Contact",
			Value:    &f.PrimaryContact,
			Presence: f.requiredForComplete,
			PIFOTag:  "51a.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:    "Primary Phone",
			Value:    &f.PrimaryPhone,
			Presence: f.requiredForComplete,
			PIFOTag:  "51b.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Secondary Contact",
			Value:   &f.SecondaryContact,
			PIFOTag: "52a.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "Secondary Phone",
			Value:   &f.SecondaryPhone,
			PIFOTag: "52b.",
		}),
		message.NewTacticalCallSignField(&message.Field{
			Label:   "Tactical Call Sign",
			Value:   &f.TacticalCallSign,
			PIFOTag: "60.",
		}),
		message.NewFCCCallSignField(&message.Field{
			Label:   "Repeater Call Sign",
			Value:   &f.RepeaterCallSign,
			PIFOTag: "61.",
		}),
		message.NewFrequencyField(&message.Field{
			Label:      "Repeater Input",
			Value:      &f.RepeaterInput,
			PIFOTag:    "62a.",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label:      "Repeater Input Tone",
			Value:      &f.RepeaterInputTone,
			PIFOTag:    "62b.",
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Repeater Input",
			TableValue: func(*message.Field) string {
				return formatFreq(f.RepeaterInput, f.RepeaterInputTone)
			},
		}),
		message.NewFrequencyField(&message.Field{
			Label:      "Repeater Output",
			Value:      &f.RepeaterOutput,
			PIFOTag:    "63a.",
			TableValue: message.TableOmit,
		}),
		message.NewTextField(&message.Field{
			Label:      "Repeater Output Tone",
			Value:      &f.RepeaterOutputTone,
			PIFOTag:    "63b.",
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Repeater Output",
			TableValue: func(*message.Field) string {
				return formatFreq(f.RepeaterOutput, f.RepeaterOutputTone)
			},
		}),
		message.NewFrequencyOffsetField(&message.Field{
			Label:   "Repeater Offset",
			Value:   &f.RepeaterOffset,
			PIFOTag: "62c.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Comments",
			Value:   &f.Comments,
			PIFOTag: "70.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Remove from List",
			Value:   &f.RemoveFromList,
			Choices: message.Choices{"checked"},
			PIFOTag: "71.",
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, nil)
	if len(f.Fields) > fieldCount {
		panic("update SheltStat22 fieldCount")
	}
	return f
}

func (f *SheltStat22) requiredForComplete() (message.Presence, string) {
	if f.ReportType == "Complete" {
		return message.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return message.PresenceOptional, ""
}

func decode22(_, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type22.HTML || form.FormVersion != Type22.Version {
		return nil
	}
	var df = make22()
	message.DecodeForm(form, df)
	return df
}

func (f *SheltStat22) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo23().Compare(actual)
}

func (f *SheltStat22) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo23().RenderPDF(env, filename)
}

func (f *SheltStat22) convertTo23() (c *SheltStat23) {
	c = make23()
	c.CopyHeaderFields(&f.BaseForm)
	c.ReportType = f.ReportType
	c.ShelterName = f.ShelterName
	c.ShelterType = f.ShelterType
	c.ShelterStatus = f.ShelterStatus
	c.ShelterAddress = f.ShelterAddress
	c.ShelterCityCode = f.ShelterCityCode
	c.ShelterCity = f.ShelterCity
	c.ShelterState = f.ShelterState
	c.ShelterZip = f.ShelterZip
	c.Latitude = f.Latitude
	c.Longitude = f.Longitude
	c.Capacity = f.Capacity
	c.Occupancy = f.Occupancy
	c.MealsServed = f.MealsServed
	c.NSSNumber = f.NSSNumber
	c.PetFriendly = f.PetFriendly
	c.BasicSafetyInspection = f.BasicSafetyInspection
	c.ATC20Inspection = f.ATC20Inspection
	c.AvailableServices = f.AvailableServices
	c.MOU = f.MOU
	c.FloorPlan = f.FloorPlan
	c.ManagedByCode = f.ManagedBy
	c.ManagedBy = f.ManagedByCode
	c.ManagedByDetail = f.ManagedByDetail
	c.PrimaryContact = f.PrimaryContact
	c.PrimaryPhone = f.PrimaryPhone
	c.SecondaryContact = f.SecondaryContact
	c.SecondaryPhone = f.SecondaryPhone
	c.TacticalCallSign = f.TacticalCallSign
	c.RepeaterCallSign = f.RepeaterCallSign
	c.RepeaterInput = f.RepeaterInput
	c.RepeaterInputTone = f.RepeaterInputTone
	c.RepeaterOutput = f.RepeaterOutput
	c.RepeaterOutputTone = f.RepeaterOutputTone
	c.RepeaterOffset = f.RepeaterOffset
	c.Comments = f.Comments
	c.RemoveFromList = f.RemoveFromList
	c.CopyFooterFields(&f.BaseForm)
	return c
}
