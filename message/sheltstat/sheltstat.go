// Package sheltstat defines the Santa Clara County OA Shelter Status Form
// message type.
package sheltstat

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/baseform"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an OA shelter status form.
var Type = message.Type{
	Tag:     "SheltStat",
	Name:    "OA shelter status form",
	Article: "an",
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*basemsg.FormVersion{
	{HTML: "form-oa-shelter-status.html", Version: "2.3", Tag: "SheltStat", FieldOrder: fieldOrder},
	{HTML: "form-oa-shelter-status.html", Version: "2.2", Tag: "SheltStat", FieldOrder: fieldOrder},
	{HTML: "form-oa-shelter-status.html", Version: "2.1", Tag: "SheltStat", FieldOrder: fieldOrder},
	{HTML: "form-oa-shelter-status.html", Version: "2.0", Tag: "SheltStat", FieldOrder: fieldOrder},
}
var fieldOrder = []string{
	"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "19.", "32.", "30.", "31.", "33a.",
	"33b.", "34b.", "33c.", "33d.", "37a.", "37b.", "40a.", "40b.", "41.", "42.", "43a.", "43b.", "43c.", "44.", "45.", "46.",
	"50a.", "49a.", "50b.", "51a.", "51b.", "52a.", "52b.", "60.", "61.", "62a.", "62b.", "63a.", "63b.", "62c.", "70.", "71.",
	"OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
}

// SheltStat holds an OA shelter status form.
type SheltStat struct {
	basemsg.BaseMessage
	baseform.BaseForm
	ReportType            string
	ShelterName           string
	ShelterType           string
	ShelterStatus         string
	ShelterAddress        string
	ShelterCityCode       string // added in v2.2
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
	ManagedByCode         string // added in v2.2
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

func New() (f *SheltStat) {
	f = create(versions[0]).(*SheltStat)
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "PRIORITY"
	return f
}

var pdfBase []byte

func create(version *basemsg.FormVersion) message.Message {
	const fieldCount = 63
	var f = SheltStat{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		PDFBase:     pdfBase,
		PDFFontSize: 10,
		Form:        version,
	}}
	f.BaseMessage.FSubject = &f.ShelterName
	f.BaseMessage.FReportType = &f.ReportType
	f.BaseMessage.FBody = &f.Comments
	var basePDFMaps = baseform.DefaultPDFMaps
	basePDFMaps.OriginMsgID = basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
		return []basemsg.PDFField{
			{Name: "Origin Msg Nbr", Value: f.OriginMsgID},
			{Name: "Origin Msg Nbr Copy", Value: f.OriginMsgID},
		}
	})
	f.Fields = make([]*basemsg.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFMaps)
	f.Fields = append(f.Fields,
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  basemsg.Choices{"Update", "Complete"},
			Presence: basemsg.Required,
			PIFOTag:  "19.",
			EditHelp: `This indicates whether the form should "Update" the previous status report for the shelter, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
			PDFMap:   basemsg.PDFNameMap{"Report Type", "", "Off"},
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Shelter Name",
			Value:     &f.ShelterName,
			Presence:  basemsg.Required,
			PIFOTag:   "32.",
			PDFMap:    basemsg.PDFName("Shelter Name"),
			EditWidth: 44,
			EditHelp:  `This is the name of the shelter whose status is being reported.  It is required.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Shelter Type",
			Value:    &f.ShelterType,
			Choices:  basemsg.Choices{"Type 1", "Type 2", "Type 3", "Type 4"},
			Presence: f.requiredForComplete,
			PIFOTag:  "30.",
			PDFMap:   basemsg.PDFNameMap{"Shelter Type", "", "Off"},
			EditHelp: `This is the shelter type.  It is required when "Report Type" is "Complete".`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Shelter Status",
			Value:    &f.ShelterStatus,
			Choices:  basemsg.Choices{"Open", "Closed", "Full"},
			Presence: f.requiredForComplete,
			PIFOTag:  "31.",
			PDFMap:   basemsg.PDFNameMap{"Shelter Status", "", "Off"},
			EditHelp: `This indicates the status of the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Shelter Address",
			Value:     &f.ShelterAddress,
			Presence:  f.requiredForComplete,
			PIFOTag:   "33a.",
			PDFMap:    basemsg.PDFName("Address"),
			EditWidth: 75,
			EditHelp:  `This is the street address of the shelter.  It is required when "Report Type" is "Complete".`,
		}),
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			basemsg.NewRestrictedField(&basemsg.Field{
				Label:      "Shelter City",
				Value:      &f.ShelterCity,
				Choices:    basemsg.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:   f.requiredForComplete,
				PIFOTag:    "33b.",
				PDFMap:     basemsg.PDFName("City"),
				TableValue: basemsg.TableOmit,
				EditHelp:   `This is the name of the city in which the shelter is located.  It is required when "Report Type" is "Complete".`,
				EditSkip:   func(*basemsg.Field) bool { return f.ShelterAddress == "" },
			}),
		)
	} else {
		f.Fields = append(f.Fields,
			basemsg.NewCalculatedField(&basemsg.Field{
				Label:    "Shelter City Code",
				Value:    &f.ShelterCityCode,
				Presence: f.requiredForComplete,
				PIFOTag:  "33b.",
			}),
			basemsg.NewTextField(&basemsg.Field{
				Label:      "Shelter City",
				Value:      &f.ShelterCity,
				Choices:    basemsg.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:   f.requiredForComplete,
				PIFOTag:    "34b.",
				PDFMap:     basemsg.PDFName("City"),
				TableValue: basemsg.TableOmit,
				EditWidth:  30,
				EditHelp:   `This is the name of the city in which the shelter is located.  It is required when "Report Type" is "Complete".`,
				EditApply: func(field *basemsg.Field, s string) {
					f.ShelterCity = field.Choices.ToPIFO(s)
					if f.ShelterCity == "" || field.Choices.IsPIFO(f.ShelterCity) {
						f.ShelterCityCode = f.ShelterCity
					} else {
						f.ShelterCityCode = "Unincorporated"
					}
				},
				EditSkip: func(*basemsg.Field) bool { return f.ShelterAddress == "" },
			}),
		)
	}
	f.Fields = append(f.Fields,
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Shelter State",
			Value:      &f.ShelterState,
			Choices:    basemsg.Choices{"CA"},
			Presence:   f.requiredForComplete,
			PIFOTag:    "33c.",
			PDFMap:     basemsg.PDFName("State"),
			TableValue: basemsg.TableOmit,
			EditWidth:  12,
			EditHelp:   `This is the name (or two-letter abbreviation) of the state in which the shelter is located.  It is required when "Report Type" is "Complete".`,
			EditSkip:   func(*basemsg.Field) bool { return f.ShelterCity == "" },
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Shelter Zip",
			Value:      &f.ShelterZip,
			Presence:   f.requiredForComplete,
			PIFOTag:    "33d.",
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("Zip"),
			TableValue: basemsg.TableOmit,
			EditWidth:  12,
			EditHelp:   `This is the shelter's ZIP code.  It is required when "Report Type" is "Complete".`,
			EditSkip:   func(*basemsg.Field) bool { return f.ShelterState == "" },
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "Shelter City",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.ShelterCity, common.SmartJoin(f.ShelterState, f.ShelterZip, "  "), ", ")
			},
		}),
		basemsg.NewRealNumberField(&basemsg.Field{
			Label:   "Latitude",
			Value:   &f.Latitude,
			PIFOTag: "37a.",
			PDFMap:  basemsg.PDFName("Latitude"),
			TableValue: func(*basemsg.Field) string {
				if f.Longitude == "" {
					return f.Latitude
				}
				return ""
			},
			EditWidth: 30,
			EditHelp:  `This is the latitude of the shelter location, expressed in fractional degrees.`,
		}),
		basemsg.NewRealNumberField(&basemsg.Field{
			Label:   "Longitude",
			Value:   &f.Longitude,
			PIFOTag: "37b.",
			PDFMap:  basemsg.PDFName("Longitude"),
			TableValue: func(*basemsg.Field) string {
				if f.Longitude == "" {
					return f.Latitude
				}
				return ""
			},
			EditWidth: 29,
			EditHelp:  `This is the longitude of the shelter location, expressed in fractional degrees.`,
			EditValid: func(field *basemsg.Field) string {
				if f.Latitude != "" && f.Longitude == "" {
					return `The "Longitude" field must have a value when "Latitude" has a value.`
				}
				return basemsg.ValidReal(field)
			},
			EditSkip: func(*basemsg.Field) bool { return f.Latitude == "" },
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "GPS Coordinates",
			TableValue: func(*basemsg.Field) string {
				if f.Latitude != "" && f.Longitude != "" {
					return f.Latitude + ", " + f.Longitude
				}
				return ""
			},
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:    "Capacity",
			Value:    &f.Capacity,
			Presence: f.requiredForComplete,
			PIFOTag:  "40a.",
			PDFMap:   basemsg.PDFName("Capacity"),
			TableValue: func(*basemsg.Field) string {
				if f.Occupancy == "" {
					return f.Capacity
				}
				return ""
			},
			EditWidth: 6,
			EditHelp:  `This is the number of people the shelter can accommodate.  It is required when "Report Type" is "Complete".`,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:    "Occupancy",
			Value:    &f.Occupancy,
			Presence: f.requiredForComplete,
			PIFOTag:  "40b.",
			PDFMap:   basemsg.PDFName("Occupancy"),
			TableValue: func(*basemsg.Field) string {
				if f.Occupancy != "" && f.Capacity != "" {
					return f.Occupancy + " out of " + f.Capacity
				}
				return f.Occupancy
			},
			EditWidth: 6,
			EditHelp:  `This is the number of people currently using the shelter.  It is required when "Report Type" is "Complete".`,
			EditSkip:  func(*basemsg.Field) bool { return f.ShelterAddress == "" },
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Meals Served",
			Value:     &f.MealsServed,
			PIFOTag:   "41.",
			PDFMap:    basemsg.PDFName("Meals"),
			EditWidth: 65,
			EditHelp:  `This is the number and/or description of meals served at the shelter in the last 24 hours.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "NSS Number",
			Value:     &f.NSSNumber,
			PIFOTag:   "42.",
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("NSS Number"),
			EditWidth: 65,
			EditHelp:  `This is the NSS number of the shelter.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Pet Friendly",
			Value:    &f.PetFriendly,
			Choices:  basemsg.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag:  "43a.",
			PDFMap:   basemsg.PDFNameMap{"Pet Friendly", "", "Off", "false", "No", "checked", "Yes"},
			EditHelp: `This indicates whether the shelter can accept pets.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Basic Safety Inspection",
			Value:    &f.BasicSafetyInspection,
			Choices:  basemsg.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag:  "43b.",
			PDFMap:   basemsg.PDFNameMap{"Basic Safety Insp", "", "Off", "false", "No", "checked", "Yes"},
			EditHelp: `This indicates whether the shelter has had a basic safety inspection.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "ATC-20 Inspection",
			Value:    &f.ATC20Inspection,
			Choices:  basemsg.ChoicePairs{"checked", "Yes", "false", "No"},
			PIFOTag:  "43c.",
			PDFMap:   basemsg.PDFNameMap{"ATC20 Insp", "", "Off", "false", "No", "checked", "Yes"},
			EditHelp: `This indicates whether the shelter has had an ATC-20 inspection.`,
		}),
		basemsg.NewMultilineField(&basemsg.Field{
			Label:     "Available Services",
			Value:     &f.AvailableServices,
			PIFOTag:   "44.",
			PDFMap:    basemsg.PDFName("Available Services"),
			EditWidth: 85,
			EditHelp:  `This is a list of services available at the shelter.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "MOU",
			Value:     &f.MOU,
			PIFOTag:   "45.",
			PDFMap:    basemsg.PDFName("MOU"),
			EditWidth: 64,
			EditHelp:  `This indicates where and how the shelter's Memorandum of Understanding (MOU) was reported.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Floor Plan",
			Value:     &f.FloorPlan,
			PIFOTag:   "46.",
			PDFMap:    basemsg.PDFName("Floorplan"),
			EditWidth: 64,
			EditHelp:  `This indicates where and how the shelter's floor plan was reported.`,
		}),
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			basemsg.NewRestrictedField(&basemsg.Field{
				Label:    "Managed By",
				Value:    &f.ManagedBy,
				Choices:  basemsg.Choices{"American Red Cross", "Private", "Community", "Government", "Other"},
				Presence: f.requiredForComplete,
				PIFOTag:  "50a.",
				PDFMap:   basemsg.PDFName("Managed By"),
				EditHelp: `This indicates what type of entity is managing the shelter.  It is required when "Report Type" is "Complete".`,
			}),
		)
	} else {
		f.Fields = append(f.Fields,
			basemsg.NewCalculatedField(&basemsg.Field{
				Label:    "Managed By Code",
				Value:    &f.ManagedByCode,
				Presence: f.requiredForComplete,
				PIFOTag:  "50a.",
			}),
			basemsg.NewTextField(&basemsg.Field{
				Label:     "Managed By",
				Value:     &f.ManagedBy,
				Choices:   basemsg.Choices{"American Red Cross", "Private", "Community", "Government", "Other"},
				Presence:  f.requiredForComplete,
				PIFOTag:   "49a.",
				PDFMap:    basemsg.PDFName("Managed By"),
				EditWidth: 18,
				EditHelp:  `This indicates what type of entity is managing the shelter.  It is required when "Report Type" is "Complete".`,
				EditApply: func(field *basemsg.Field, s string) {
					f.ManagedBy = field.Choices.ToPIFO(s)
					if f.ManagedBy == "" || field.Choices.IsPIFO(f.ManagedBy) {
						f.ManagedByCode = f.ManagedBy
					} else {
						f.ManagedByCode = "Other"
					}
				},
				EditValid: func(f *basemsg.Field) string {
					if *f.Value != "" && !f.Choices.IsPIFO(*f.Value) {
						return `The "Managed By" field should contain one of the recommended values ("American Red Cross", "Private", "Community", "Government", or "Other").  Other values can be successfully sent in a PackItForms message, but they cannot be rendered in the generated PDF.`
					}
					return ""
				},
			}),
		)
	}
	f.Fields = append(f.Fields,
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Managed By Detail",
			Value:     &f.ManagedByDetail,
			PIFOTag:   "50b.",
			PDFMap:    basemsg.PDFName("Managed By Detail"),
			EditWidth: 65,
			EditHelp:  `This is additional detail about who is managing the shelter (particularly if "Managed By" is "Other").`,
			EditValid: func(*basemsg.Field) string {
				if f.ManagedBy == "Other" && f.ManagedByDetail == "" {
					return `The "Managed By Detail" field is required when "Managed By" is "Other".`
				}
				return ""
			},
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Primary Contact",
			Value:     &f.PrimaryContact,
			Presence:  f.requiredForComplete,
			PIFOTag:   "51a.",
			PDFMap:    basemsg.PDFName("Pri Contact"),
			EditWidth: 65,
			EditHelp:  `This is the name of the primary contact person for the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Primary Phone",
			Value:     &f.PrimaryPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "51b.",
			PDFMap:    basemsg.PDFName("Pri Contact Phone"),
			EditWidth: 65,
			EditHelp:  `This is the phone number of the primary contact person for the shelter.  It is required when "Report Type" is "Complete".`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Secondary Contact",
			Value:     &f.SecondaryContact,
			PIFOTag:   "52a.",
			PDFMap:    basemsg.PDFName("Sec Contact"),
			EditWidth: 65,
			EditHelp:  `This is the name of the secondary contact person for the shelter.`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Secondary Phone",
			Value:     &f.SecondaryPhone,
			PIFOTag:   "52b.",
			PDFMap:    basemsg.PDFName("Sec Contact Phone"),
			EditWidth: 65,
			EditHelp:  `This is the phone number of the secondary contact person for the shelter.`,
		}),
		basemsg.NewTacticalCallSignField(&basemsg.Field{
			Label:     "Tactical Call Sign",
			Value:     &f.TacticalCallSign,
			PIFOTag:   "60.",
			PDFMap:    basemsg.PDFName("Tactical Call Sign"),
			EditWidth: 29,
			EditHelp:  `This is the tactical call sign assigned to the shelter for amateur radio communications.`,
		}),
		basemsg.NewFCCCallSignField(&basemsg.Field{
			Label:     "Repeater Call Sign",
			Value:     &f.RepeaterCallSign,
			PIFOTag:   "61.",
			PDFMap:    basemsg.PDFName("Repeater Call Sign"),
			EditWidth: 29,
			EditHelp:  `This is the call sign of the amateur radio repeater that the shelter is monitoring for communications.`,
		}),
		basemsg.NewFrequencyField(&basemsg.Field{
			Label:      "Repeater Input",
			Value:      &f.RepeaterInput,
			PIFOTag:    "62a.",
			PDFMap:     basemsg.PDFName("Input Freq"),
			TableValue: basemsg.TableOmit,
			EditWidth:  20,
			EditHelp:   `This is the input frequency (in MHz) of the amateur radio repeater that the shelter is using for communications.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Repeater Input Tone",
			Value:      &f.RepeaterInputTone,
			PIFOTag:    "62b.",
			PDFMap:     basemsg.PDFName("Input Tone"),
			TableValue: basemsg.TableOmit,
			EditWidth:  30,
			EditHelp:   `This is the analog CTCSS tone, P25 NAC, DMR TS/TG/CC, or other access details required by the amateur radio repeater that the shelter is using for communications.`,
			EditSkip:   func(*basemsg.Field) bool { return f.RepeaterInput == "" },
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "Repeater Input",
			TableValue: func(*basemsg.Field) string {
				return formatFreq(f.RepeaterInput, f.RepeaterInputTone)
			},
		}),
		basemsg.NewFrequencyField(&basemsg.Field{
			Label:      "Repeater Output",
			Value:      &f.RepeaterOutput,
			PIFOTag:    "63a.",
			PDFMap:     basemsg.PDFName("Output Freq"),
			TableValue: basemsg.TableOmit,
			EditWidth:  20,
			EditHelp:   `This is the output frequency (in MHz) of the amateur radio repeater that the shelter is using for communications.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:      "Repeater Output Tone",
			Value:      &f.RepeaterOutputTone,
			PIFOTag:    "63b.",
			PDFMap:     basemsg.PDFName("Output Tone"),
			TableValue: basemsg.TableOmit,
			EditWidth:  30,
			EditHelp:   `This is the analog CTCSS tone, P25 NAC, DMR TS/TG/CC, or other access details for the transmission from the amateur radio repeater that the shelter is using for communications.`,
			EditSkip:   func(*basemsg.Field) bool { return f.RepeaterOutput == "" },
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
			Label: "Repeater Output",
			TableValue: func(*basemsg.Field) string {
				return formatFreq(f.RepeaterOutput, f.RepeaterOutputTone)
			},
		}),
		basemsg.NewFrequencyOffsetField(&basemsg.Field{
			Label:     "Repeater Offset",
			Value:     &f.RepeaterOffset,
			PIFOTag:   "62c.",
			PDFMap:    basemsg.PDFName("Offset"),
			EditWidth: 15,
			EditHelp:  `This is the offset for the amateur radio repeater that the shelter is using for communications.  It can be either a number in MHz, a "+", or a "-".`,
			EditSkip:  func(*basemsg.Field) bool { return (f.RepeaterInput == "") == (f.RepeaterOutput == "") },
		}),
		basemsg.NewMultilineField(&basemsg.Field{
			Label:     "Comments",
			Value:     &f.Comments,
			PIFOTag:   "70.",
			PDFMap:    basemsg.PDFName("Comments"),
			EditWidth: 85,
			EditHelp:  `These are comments regarding the status of the shelter.`,
		}),
	)
	if f.Form.Version < "2.3" {
		f.Fields = append(f.Fields,
			basemsg.NewRestrictedField(&basemsg.Field{
				Label:    "Remove from List",
				Value:    &f.RemoveFromList,
				Choices:  basemsg.Choices{"checked"},
				PIFOTag:  "71.",
				PDFMap:   basemsg.PDFNameMap{"Remove", "", "No", "checked", "Yes"},
				EditHelp: `This indicates whether the shelter should be removed from the receiver's list of shelters.`,
			}),
		)
	} else {
		f.Fields = append(f.Fields,
			basemsg.NewRestrictedField(&basemsg.Field{
				Label:    "Remove from List",
				Value:    &f.RemoveFromList,
				Choices:  basemsg.ChoicePairs{"checked", "Yes", "false", "No"},
				PIFOTag:  "71.",
				PDFMap:   basemsg.PDFNameMap{"Remove", "", "Off", "false", "No", "checked", "Yes"},
				EditHelp: `This indicates whether the shelter should be removed from the receiver's list of shelters.`,
			}),
		)
	}
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFMaps)
	if len(f.Fields) > fieldCount {
		panic("update SheltStat fieldCount")
	}
	return &f
}

func (f *SheltStat) requiredForComplete() (basemsg.Presence, string) {
	if f.ReportType == "Complete" {
		return basemsg.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return basemsg.PresenceOptional, ""
}

func decode(subject, body string) (f *SheltStat) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-oa-shelter-status.html") {
		return nil
	}
	return basemsg.Decode(body, versions, create).(*SheltStat)
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
