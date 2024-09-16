// Package jurisstat defines the Santa Clara County OA Jurisdiction Status Form
// message type.
package jurisstat

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for an OA jurisdiction status form.
var Type = message.Type{
	Tag:         "JurisStat",
	Name:        "OA jurisdiction status form",
	Article:     "an",
	PDFRenderV2: true,
}

// OldType is the previous type definition for an OA jurisdiction status form.
var OldType = message.Type{
	Tag:     "MuniStat",
	Name:    "OA municipal status form",
	Article: "an",
}

func init() {
	Type.Create = New
	Type.Decode = decode
	OldType.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*message.FormVersion{
	{HTML: "form-oa-muni-status.html", Version: "2.2", Tag: "JurisStat", FieldOrder: fieldOrder},
	{HTML: "form-oa-muni-status.html", Version: "2.1", Tag: "MuniStat", FieldOrder: fieldOrder},
	{HTML: "form-oa-muni-status.html", Version: "2.0", Tag: "MuniStat", FieldOrder: fieldOrder},
}
var fieldOrder = []string{
	"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "19.", "21.", "22.", "23.", "24.",
	"25.", "26.", "27.", "28.", "29.", "30.", "31.", "32.", "33.", "34.", "35.", "36.", "37.", "38.", "39.", "40.", "99.",
	"41.0.", "41.1.", "42.0.", "42.1.", "43.0.", "43.1.", "44.0.", "44.1.", "45.0.", "45.1.", "46.0.", "46.1.", "47.0.",
	"47.1.", "48.0.", "48.1.", "49.0.", "49.1.", "50.0.", "50.1.", "51.0.", "51.1.", "52.0.", "52.1.", "53.0.", "53.1.",
	"54.0.", "54.1.", "55.0.", "55.1.", "56.0.", "56.1.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
}
var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgIDR: message.PDFMultiRenderer{
		&message.PDFTextRenderer{Page: 1, X: 223, Y: 60, R: 349, B: 77, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 468, Y: 31, R: 573, B: 47},
		&message.PDFTextRenderer{Page: 3, X: 468, Y: 31, R: 573, B: 47},
	},
	DestinationMsgIDR: &message.PDFTextRenderer{Page: 1, X: 446, Y: 60, R: 573, B: 77, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDateR:      &message.PDFTextRenderer{Page: 1, X: 70, Y: 119, R: 138, B: 140, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageTimeR:      &message.PDFTextRenderer{Page: 1, X: 201, Y: 119, R: 236, B: 140, Style: message.PDFTextStyle{VAlign: "baseline"}},
	HandlingR: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
		"IMMEDIATE": {306, 129},
		"PRIORITY":  {408, 129},
		"ROUTINE":   {494, 129},
	}},
	ToICSPositionR:   &message.PDFTextRenderer{Page: 1, X: 127, Y: 141, R: 303, B: 162, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocationR:      &message.PDFTextRenderer{Page: 1, X: 127, Y: 164, R: 303, B: 185, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToNameR:          &message.PDFTextRenderer{Page: 1, X: 127, Y: 186, R: 303, B: 207, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContactR:       &message.PDFTextRenderer{Page: 1, X: 127, Y: 209, R: 303, B: 230, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPositionR: &message.PDFTextRenderer{Page: 1, X: 398, Y: 141, R: 573, B: 162, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocationR:    &message.PDFTextRenderer{Page: 1, X: 398, Y: 164, R: 573, B: 185, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromNameR:        &message.PDFTextRenderer{Page: 1, X: 398, Y: 186, R: 573, B: 207, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContactR:     &message.PDFTextRenderer{Page: 1, X: 398, Y: 209, R: 573, B: 230, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvdR:     &message.PDFTextRenderer{Page: 3, X: 110, Y: 509, R: 321, B: 526, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySentR:     &message.PDFTextRenderer{Page: 3, X: 358, Y: 509, R: 573, B: 526, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpNameR:          &message.PDFTextRenderer{Page: 3, X: 76, Y: 528, R: 250, B: 545, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCallR:          &message.PDFTextRenderer{Page: 3, X: 302, Y: 528, R: 367, B: 545, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDateR:          &message.PDFTextRenderer{Page: 3, X: 413, Y: 528, R: 473, B: 545, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTimeR:          &message.PDFTextRenderer{Page: 3, X: 539, Y: 528, R: 573, B: 545, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// JurisStat holds an OA jurisdiction status form.
type JurisStat struct {
	message.BaseMessage
	baseform.BaseForm
	ReportType                    string
	JurisdictionCode              string // added in 2.2
	Jurisdiction                  string
	EOCPhone                      string
	EOCFax                        string
	PriEMContactName              string
	PriEMContactPhone             string
	SecEMContactName              string
	SecEMContactPhone             string
	OfficeStatus                  string
	GovExpectedOpenDate           string
	GovExpectedOpenTime           string
	GovExpectedCloseDate          string
	GovExpectedCloseTime          string
	EOCOpen                       string
	EOCActivationLevel            string
	EOCExpectedOpenDate           string
	EOCExpectedOpenTime           string
	EOCExpectedCloseDate          string
	EOCExpectedCloseTime          string
	StateOfEmergency              string
	HowSOESent                    string
	Communications                string
	CommunicationsComments        string
	Debris                        string
	DebrisComments                string
	Flooding                      string
	FloodingComments              string
	Hazmat                        string
	HazmatComments                string
	EmergencyServices             string
	EmergencyServicesComments     string
	Casualties                    string
	CasualtiesComments            string
	UtilitiesGas                  string
	UtilitiesGasComments          string
	UtilitiesElectric             string
	UtilitiesElectricComments     string
	InfrastructurePower           string
	InfrastructurePowerComments   string
	InfrastructureWater           string
	InfrastructureWaterComments   string
	InfrastructureSewer           string
	InfrastructureSewerComments   string
	SearchAndRescue               string
	SearchAndRescueComments       string
	TransportationRoads           string
	TransportationRoadsComments   string
	TransportationBridges         string
	TransportationBridgesComments string
	CivilUnrest                   string
	CivilUnrestComments           string
	AnimalIssues                  string
	AnimalIssuesComments          string
}

func New() (f *JurisStat) {
	f = create(versions[0]).(*JurisStat)
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Handling = "IMMEDIATE"
	f.ToLocation = "County EOC"
	return f
}

func create(version *message.FormVersion) message.Message {
	const fieldCount = 80
	var f = JurisStat{BaseMessage: message.BaseMessage{
		Type: &Type,
		Form: version,
	}}
	if version.Version < "2.2" {
		f.Type = &OldType
	}
	f.BaseMessage.FSubject = &f.Jurisdiction
	f.BaseMessage.FBody = &f.CommunicationsComments
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  message.Choices{"Update", "Complete"},
			Presence: message.Required,
			PIFOTag:  "19.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Update":   {120, 238.5},
				"Complete": {176.5, 238.5},
			}},
			EditHelp: `This indicates whether the form should "Update" the previous status report for the jurisdiction, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
		}),
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:       "Jurisdiction Name",
				Value:       &f.Jurisdiction,
				Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:    message.Required,
				PIFOTag:     "21.",
				PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 334, Y: 232, R: 571, B: 258},
				EditWidth:   42,
				EditHelp:    `This is the name of the jurisdiction being described by the form.  It is required.`,
			}),
		)
	} else {
		f.Fields = append(f.Fields,
			message.NewCalculatedField(&message.Field{
				Label:   "Jurisdiction Code",
				Value:   &f.JurisdictionCode,
				PIFOTag: "21.",
			}),
			message.NewTextField(&message.Field{
				Label:       "Jurisdiction Name",
				Value:       &f.Jurisdiction,
				Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:    message.Required,
				PIFOTag:     "22.",
				PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 334, Y: 232, R: 571, B: 258, Style: message.PDFTextStyle{VAlign: "baseline"}},
				EditWidth:   42,
				EditHelp:    `This is the name of the jurisdiction being described by the form.  It is required.`,
				EditApply: func(field *message.Field, v string) {
					f.Jurisdiction = v
					if v == "" || field.Choices.IsPIFO(v) {
						f.JurisdictionCode = v
					} else {
						f.JurisdictionCode = "Unincorporated"
					}
				},
			}),
		)
	}
	f.Fields = append(f.Fields,
		message.NewPhoneNumberField(&message.Field{
			Label:       "EOC Phone",
			Value:       &f.EOCPhone,
			Presence:    f.requiredForComplete,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 150, Y: 287, R: 302, B: 303, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   34,
			EditHelp:    `This is the phone number of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "EOC Fax",
			Value:       &f.EOCFax,
			PIFOTag:     "24.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 420, Y: 287, R: 571, B: 303, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   37,
			EditHelp:    `This is the fax number of the jurisdiction's Emergency Operations Center (EOC).`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Primary EM Contact Name",
			Value:       &f.PriEMContactName,
			Presence:    f.requiredForComplete,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 150, Y: 305, R: 302, B: 321, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   27,
			EditHelp:    `This is the name of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Primary EM Contact Phone",
			Value:       &f.PriEMContactPhone,
			Presence:    f.requiredForComplete,
			PIFOTag:     "26.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 420, Y: 305, R: 571, B: 321, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   26,
			EditHelp:    `This is the phone number of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Secondary EM Contact Name",
			Value:       &f.SecEMContactName,
			PIFOTag:     "27.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 150, Y: 323, R: 302, B: 339, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   26,
			EditHelp:    `This is the name of the secondary emergency manager of the jurisdiction.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Secondary EM Contact Phone",
			Value:       &f.SecEMContactPhone,
			PIFOTag:     "28.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 420, Y: 323, R: 571, B: 339, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   26,
			EditHelp:    `This is the phone number of the secondary emergency manager of the jurisdiction.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Govt. Office Status",
			Value:    &f.OfficeStatus,
			Choices:  message.Choices{"Unknown", "Open", "Closed"},
			Presence: f.requiredForComplete,
			PIFOTag:  "29.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Unknown": {209, 376},
				"Open":    {321, 376},
				"Closed":  {443, 376},
			}},
			EditHelp: `This indicates whether the jurisdiction's regular business offices are open.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Govt. Office Expected Open Date",
			Value:       &f.GovExpectedOpenDate,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 154, Y: 387, R: 302, B: 402, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the jurisdiction's regular business offices are expected to open, in MM/DD/YYYY format.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Govt. Office Expected Open Time",
			Value:       &f.GovExpectedOpenTime,
			PIFOTag:     "31.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 424, Y: 387, R: 571, B: 402, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the jurisdiction's regular business offices are expected to open, in HH:MM format (24-hour clock).`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Govt. Office Expected to Open",
			EditHelp: `This is the date and time when the jurisdiction's regular business offices are expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
		}, &f.GovExpectedOpenDate, &f.GovExpectedOpenTime),
		message.NewDateField(true, &message.Field{
			Label:       "Govt. Office Expected Close Date",
			Value:       &f.GovExpectedCloseDate,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 154, Y: 404, R: 302, B: 420, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the jurisdiction's regular business offices are expected to close, in MM/DD/YYYY format.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Govt. Office Expected Close Time",
			Value:       &f.GovExpectedCloseTime,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 424, Y: 404, R: 571, B: 420, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the jurisdiction's regular business offices are expected to close, in HH:MM format (24-hour clock).`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Govt. Office Expected to Close",
			EditHelp: `This is the date and time when the jurisdiction's regular business offices are expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).`,
		}, &f.GovExpectedCloseDate, &f.GovExpectedCloseTime),
		message.NewRestrictedField(&message.Field{
			Label:    "EOC Open",
			Value:    &f.EOCOpen,
			Choices:  message.Choices{"Unknown", "Yes", "No"},
			Presence: f.requiredForComplete,
			PIFOTag:  "34.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Unknown": {209, 458},
				"Yes":     {321, 458},
				"No":      {443, 458},
			}},
			EditHelp: `This indicates whether the jurisdiction's Emergency Operations Center (EOC) is open.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "EOC Activation Level",
			Value:    &f.EOCActivationLevel,
			Choices:  message.Choices{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
			Presence: f.requiredForComplete,
			PIFOTag:  "35.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Normal":       {209, 474},
				"Duty Officer": {321, 474},
				"Monitor":      {443, 474},
				"Partial":      {209, 490},
				"Full":         {321, 490},
			}},
			EditHelp: `This indicates the activation level of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "EOC Expected to Open Date",
			Value:       &f.EOCExpectedOpenDate,
			PIFOTag:     "36.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 154, Y: 499, R: 302, B: 515, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the jurisdiction's Emergency Operations Center (EOC) is expected to open, in MM/DD/YYYY format.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "EOC Expected to Open Time",
			Value:       &f.EOCExpectedOpenTime,
			PIFOTag:     "37.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 424, Y: 499, R: 571, B: 515, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the jurisdiction's Emergency Operations Center (EOC) is expected to open, in HH:MM format (24-hour clock).`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "EOC Expected to Open",
			EditHelp: `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
		}, &f.EOCExpectedOpenDate, &f.EOCExpectedOpenTime),
		message.NewDateField(true, &message.Field{
			Label:       "EOC Expected to Close Date",
			Value:       &f.EOCExpectedCloseDate,
			PIFOTag:     "38.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 154, Y: 517, R: 302, B: 532, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the jurisdiction's Emergency Operations Center (EOC) is expected to close, in MM/DD/YYYY format.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "EOC Expected to Close Time",
			Value:       &f.EOCExpectedCloseTime,
			PIFOTag:     "39.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 424, Y: 517, R: 571, B: 532, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time when the jurisdiction's Emergency Operations Center (EOC) is expected to close, in HH:MM format (24-hour clock).`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "EOC Expected to Close",
			EditHelp: `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).`,
		}, &f.EOCExpectedCloseDate, &f.EOCExpectedCloseTime),
		message.NewRestrictedField(&message.Field{
			Label:    "State of Emergency",
			Value:    &f.StateOfEmergency,
			Choices:  message.Choices{"Unknown", "Yes", "No"},
			Presence: f.requiredForComplete,
			PIFOTag:  "40.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 1, Points: map[string][]float64{
				"Unknown": {209, 569},
				"Yes":     {321, 569},
				"No":      {443, 569},
			}},
			EditHelp: `This indicates whether the jurisdiction has a declared state of emergency.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewTextField(&message.Field{
			Label: "How SOE Sent",
			Value: &f.HowSOESent,
			Presence: func() (message.Presence, string) {
				if f.StateOfEmergency == "Yes" {
					return message.PresenceRequired, `when "State of Emergency" is "Yes"`
				} else {
					return message.PresenceNotAllowed, `when "State of Emergency" is not "Yes"`
				}
			},
			PIFOTag:     "99.",
			PDFRenderer: &message.PDFTextRenderer{Page: 1, X: 229, Y: 579, R: 571, B: 598},
			EditWidth:   58,
			EditHelp:    `This describes where and how the jurisdiction's "state of emergency" declaration was delivered.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Communications",
			Value:   &f.Communications,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "41.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 80},
				"Unknown":   {312.5, 80},
				"Problem":   {407, 80},
				"Failure":   {506, 80},
				"Delayed":   {218, 95},
				"Closed":    {312.5, 95},
				"Early Out": {407, 95},
			}},
			EditHelp: `This describes the current situation status with respect to communications.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Communications: Comments",
			Value:       &f.CommunicationsComments,
			PIFOTag:     "41.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 105, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to communications.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Debris",
			Value:   &f.Debris,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "42.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 141.5},
				"Unknown":   {312.5, 141.5},
				"Problem":   {407, 141.5},
				"Failure":   {506, 141.5},
				"Delayed":   {218, 157.5},
				"Closed":    {312.5, 157.5},
				"Early Out": {407, 157.5},
			}},
			EditHelp: `This describes the current situation status with respect to debris.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Debris: Comments",
			Value:       &f.DebrisComments,
			PIFOTag:     "42.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 167.5, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to debris.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Flooding",
			Value:   &f.Flooding,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "43.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 204},
				"Unknown":   {312.5, 204},
				"Problem":   {407, 204},
				"Failure":   {506, 204},
				"Delayed":   {218, 219.5},
				"Closed":    {312.5, 219.5},
				"Early Out": {407, 219.5},
			}},
			EditHelp: `This describes the current situation status with respect to flooding.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Flooding: Comments",
			Value:       &f.FloodingComments,
			PIFOTag:     "43.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 229.5, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to flooding.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Hazmat",
			Value:   &f.Hazmat,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "44.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 266},
				"Unknown":   {312.5, 266},
				"Problem":   {407, 266},
				"Failure":   {506, 266},
				"Delayed":   {218, 282},
				"Closed":    {312.5, 282},
				"Early Out": {407, 282},
			}},
			EditHelp: `This describes the current situation status with respect to hazmat.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Hazmat: Comments",
			Value:       &f.HazmatComments,
			PIFOTag:     "44.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 292, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to hazmat.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Emergency Services",
			Value:   &f.EmergencyServices,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "45.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 328.5},
				"Unknown":   {312.5, 328.5},
				"Problem":   {407, 328.5},
				"Failure":   {506, 328.5},
				"Delayed":   {218, 344},
				"Closed":    {312.5, 344},
				"Early Out": {407, 344},
			}},
			EditHelp: `This describes the current situation status with respect to emergency services.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Emergency Services: Comments",
			Value:       &f.EmergencyServicesComments,
			PIFOTag:     "45.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 354, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to emergency services.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Casualties",
			Value:   &f.Casualties,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "46.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 390.5},
				"Unknown":   {312.5, 390.5},
				"Problem":   {407, 390.5},
				"Failure":   {506, 390.5},
				"Delayed":   {218, 406},
				"Closed":    {312.5, 406},
				"Early Out": {407, 406},
			}},
			EditHelp: `This describes the current situation status with respect to casualties.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Casualties: Comments",
			Value:       &f.CasualtiesComments,
			PIFOTag:     "46.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 416, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to casualties.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Utilities Gas",
			Value:   &f.UtilitiesGas,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "47.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 453},
				"Unknown":   {312.5, 453},
				"Problem":   {407, 453},
				"Failure":   {506, 453},
				"Delayed":   {218, 468.5},
				"Closed":    {312.5, 468.5},
				"Early Out": {407, 468.5},
			}},
			EditHelp: `This describes the current situation status with respect to utilities (gas).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Utilities Gas: Comments",
			Value:       &f.UtilitiesGasComments,
			PIFOTag:     "47.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 478.5, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to utilities (gas).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Utilities Electric",
			Value:   &f.UtilitiesElectric,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "48.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 515},
				"Unknown":   {312.5, 515},
				"Problem":   {407, 515},
				"Failure":   {506, 515},
				"Delayed":   {218, 530},
				"Closed":    {312.5, 530},
				"Early Out": {407, 530},
			}},
			EditHelp: `This describes the current situation status with respect to utilities (electric).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Utilities Electric: Comments",
			Value:       &f.UtilitiesElectricComments,
			PIFOTag:     "48.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 540, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to utilities (electric).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure Power",
			Value:   &f.InfrastructurePower,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "49.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 577},
				"Unknown":   {312.5, 577},
				"Problem":   {407, 577},
				"Failure":   {506, 577},
				"Delayed":   {218, 592.5},
				"Closed":    {312.5, 592.5},
				"Early Out": {407, 592.5},
			}},
			EditHelp: `This describes the current situation status with respect to infrastructure (power).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Infrastructure Power: Comments",
			Value:       &f.InfrastructurePowerComments,
			PIFOTag:     "49.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 602.5, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to infrastructure (power).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure Water",
			Value:   &f.InfrastructureWater,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "50.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Points: map[string][]float64{
				"Normal":    {218, 639},
				"Unknown":   {312.5, 639},
				"Problem":   {407, 639},
				"Failure":   {506, 639},
				"Delayed":   {218, 655},
				"Closed":    {312.5, 655},
				"Early Out": {407, 655},
			}},
			EditHelp: `This describes the current situation status with respect to infrastructure (water).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Infrastructure Water: Comments",
			Value:       &f.InfrastructureWaterComments,
			PIFOTag:     "50.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 216, Y: 665, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to infrastructure (water).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure Sewer",
			Value:   &f.InfrastructureSewer,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "51.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Points: map[string][]float64{
				"Normal":    {218, 81},
				"Unknown":   {312.5, 81},
				"Problem":   {406.5, 81},
				"Failure":   {506, 81},
				"Delayed":   {218, 96.5},
				"Closed":    {312.5, 96.5},
				"Early Out": {406.5, 96.5},
			}},
			EditHelp: `This describes the current situation status with respect to infrastructure (sewer).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Infrastructure Sewer: Comments",
			Value:       &f.InfrastructureSewerComments,
			PIFOTag:     "51.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 216, Y: 106.5, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to infrastructure (sewer).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Search And Rescue",
			Value:   &f.SearchAndRescue,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "52.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Points: map[string][]float64{
				"Normal":    {218, 143},
				"Unknown":   {312.5, 143},
				"Problem":   {406.5, 143},
				"Failure":   {506, 143},
				"Delayed":   {218, 159},
				"Closed":    {312.5, 159},
				"Early Out": {406.5, 159},
			}},
			EditHelp: `This describes the current situation status with respect to search and rescue.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Search And Rescue: Comments",
			Value:       &f.SearchAndRescueComments,
			PIFOTag:     "52.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 216, Y: 169, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to search and rescue.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Transportation Roads",
			Value:   &f.TransportationRoads,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "53.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Points: map[string][]float64{
				"Normal":    {218, 205},
				"Unknown":   {312.5, 205},
				"Problem":   {406.5, 205},
				"Failure":   {506, 205},
				"Delayed":   {218, 221},
				"Closed":    {312.5, 221},
				"Early Out": {406.5, 221},
			}},
			EditHelp: `This describes the current situation status with respect to transportation (roads).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Transportation Roads: Comments",
			Value:       &f.TransportationRoadsComments,
			PIFOTag:     "53.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 216, Y: 231, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to transportation (roads).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Transportation Bridges",
			Value:   &f.TransportationBridges,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "54.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Points: map[string][]float64{
				"Normal":    {218, 268},
				"Unknown":   {312.5, 268},
				"Problem":   {406.5, 268},
				"Failure":   {506, 268},
				"Delayed":   {218, 283},
				"Closed":    {312.5, 283},
				"Early Out": {406.5, 283},
			}},
			EditHelp: `This describes the current situation status with respect to transportation (bridges).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Transportation Bridges: Comments",
			Value:       &f.TransportationBridgesComments,
			PIFOTag:     "54.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 216, Y: 293, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to transportation (bridges).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Civil Unrest",
			Value:   &f.CivilUnrest,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "55.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Points: map[string][]float64{
				"Normal":    {218, 330},
				"Unknown":   {312.5, 330},
				"Problem":   {406.5, 330},
				"Failure":   {506, 330},
				"Delayed":   {218, 345},
				"Closed":    {312.5, 345},
				"Early Out": {406.5, 345},
			}},
			EditHelp: `This describes the current situation status with respect to civil unrest.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Civil Unrest: Comments",
			Value:       &f.CivilUnrestComments,
			PIFOTag:     "55.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 216, Y: 355, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to civil unrest.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Animal Issues",
			Value:   &f.AnimalIssues,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "56.0.",
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Points: map[string][]float64{
				"Normal":    {218, 392},
				"Unknown":   {312.5, 392},
				"Problem":   {406.5, 392},
				"Failure":   {506, 392},
				"Delayed":   {218, 407.5},
				"Closed":    {312.5, 407.5},
				"Early Out": {406.5, 407.5},
			}},
			EditHelp: `This describes the current situation status with respect to animal issues.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Animal Issues: Comments",
			Value:       &f.AnimalIssuesComments,
			PIFOTag:     "56.1.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 216, Y: 417.5, R: 571, H: 27, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   60,
			EditHelp:    `These are comments on the current situation status with respect to animal issues.`,
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update JurisStat fieldCount")
	}
	return &f
}

func (f *JurisStat) requiredForComplete() (message.Presence, string) {
	if f.ReportType == "Complete" {
		return message.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return message.PresenceOptional, ""
}

func decode(subject, body string) (f *JurisStat) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-oa-muni-status.html") {
		return nil
	}
	if df, ok := message.DecodeForm(body, versions, create).(*JurisStat); ok {
		return df
	} else {
		return nil
	}
}
