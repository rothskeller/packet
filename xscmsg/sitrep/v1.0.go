// Package sitrep defines the Situation Report Form message type.
package sitrep

import (
	"time"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type is the type definition for a situation report form.
var Type = message.Type{
	Tag:     "SitRep",
	HTML:    "form-situation-report.html",
	Version: "1.0",
	Name:    "situation report form",
	Article: "a",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.",
		"8c.", "7d.", "8d.", "20d.", "20t.", "21.", "22.", "23.", "24.",
		"25.", "26.", "30.", "31.", "32.", "50.", "51d.", "51t.",
		"52d.", "52t.", "60.", "61.", "62.", "63.", "70a.", "70b.",
		"71a.", "71b.", "72a.", "72b.", "73a.", "73b.", "74a.", "74b.",
		"75a.", "75b.", "76a.", "76b.", "77a.", "77b.", "78a.", "78b.",
		"79a.", "79b.", "80a.", "80b.", "81a.", "81b.", "82a.", "82b.",
		"83a.", "83b.", "84a.", "84b.", "85a.", "85b.", "90.", "91.",
		"92.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall",
		"OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type, decode, create)
}

var basePDFRenderers = baseform.BaseFormPDF{
	OriginMsgID: &message.PDFMultiRenderer{
		&message.PDFTextRenderer{X: 225.36, Y: 62.04, W: 112.80, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
		&message.PDFTextRenderer{Page: 2, X: 457.46, Y: 35.31, W: 109.14, H: 12.00},
		&message.PDFTextRenderer{Page: 3, X: 457.46, Y: 35.31, W: 109.14, H: 12.00},
	},
	DestinationMsgID: &message.PDFTextRenderer{X: 444.00, Y: 62.04, W: 116.76, H: 16.08, Style: message.PDFTextStyle{VAlign: "baseline"}},
	MessageDate:      &message.PDFTextRenderer{X: 72.72, Y: 116.76, W: 47.04, H: 20, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	MessageTime:      &message.PDFTextRenderer{X: 189.84, Y: 116.76, W: 42.48, H: 20, Style: message.PDFTextStyle{VAlign: "baseline", FontSize: 9}},
	Handling: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
		"IMMEDIATE": {299.32, 127.29},
		"PRIORITY":  {401.04, 127.29},
		"ROUTINE":   {486.10, 127.29},
	}},
	ToICSPosition:   &message.PDFTextRenderer{X: 132.84, Y: 145.08, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToLocation:      &message.PDFTextRenderer{X: 132.84, Y: 171.48, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToName:          &message.PDFTextRenderer{X: 132.84, Y: 197.88, W: 144.48, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	ToContact:       &message.PDFTextRenderer{X: 132.84, Y: 224.28, W: 144.48, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromICSPosition: &message.PDFTextRenderer{X: 375.84, Y: 145.08, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromLocation:    &message.PDFTextRenderer{X: 375.84, Y: 171.48, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromName:        &message.PDFTextRenderer{X: 375.84, Y: 197.88, W: 183.60, H: 19.80, Style: message.PDFTextStyle{VAlign: "baseline"}},
	FromContact:     &message.PDFTextRenderer{X: 375.84, Y: 224.28, W: 183.60, H: 19.92, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelayRcvd:     &message.PDFTextRenderer{Page: 3, X: 111.00, Y: 418.44, W: 204.48, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpRelaySent:     &message.PDFTextRenderer{Page: 3, X: 356.04, Y: 418.44, W: 204.72, H: 12.36, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpName:          &message.PDFTextRenderer{Page: 3, X: 77.04, Y: 437.28, W: 166.68, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpCall:          &message.PDFTextRenderer{Page: 3, X: 302.04, Y: 437.28, W: 58.20, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpDate:          &message.PDFTextRenderer{Page: 3, X: 399.00, Y: 437.28, W: 68.88, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
	OpTime:          &message.PDFTextRenderer{Page: 3, X: 532.92, Y: 437.28, W: 27.84, H: 12.60, Style: message.PDFTextStyle{VAlign: "baseline"}},
}

// SituationReport holds a shelter form.
type SituationReport struct {
	message.BaseMessage
	baseform.BaseForm
	PreparedDate           string
	PreparedTime           string
	Jurisdiction           string
	IncidentName           string
	EmergencyDeclaration   string
	EOCActivation          string
	IncidentDescription    string
	AdditionalInformation  string
	PreparedBy             string
	PreparedByPhone        string
	PreparedByEmail        string
	OfficeStatus           string
	ExpectedToOpenDate     string
	ExpectedToOpenTime     string
	ExpectedToCloseDate    string
	ExpectedToCloseTime    string
	ApprovedBy             string
	ApprovedByPhone        string
	ApprovedByPosition     string
	ApprovedByLocation     string
	CommunicationsStatus   string
	CommunicationsComments string
	DebrisStatus           string
	DebrisComments         string
	FloodingStatus         string
	FloodingComments       string
	HazMatStatus           string
	HazMatComments         string
	EmergencySvcsStatus    string
	EmergencySvcsComments  string
	CasualtiesStatus       string
	CasualtiesComments     string
	GasStatus              string
	GasComments            string
	ElectricStatus         string
	ElectricComments       string
	PowerStatus            string
	PowerComments          string
	WaterStatus            string
	WaterComments          string
	SewerStatus            string
	SewerComments          string
	SARStatus              string
	SARComments            string
	RoadsStatus            string
	RoadsComments          string
	BridgesStatus          string
	BridgesComments        string
	UnrestStatus           string
	UnrestComments         string
	AnimalStatus           string
	AnimalComments         string
	LifelineStatus         string
	LifelineUpdate         string
	UnmetNeeds             string
}

func create() message.Message {
	f := makeF()
	f.MessageDate = time.Now().Format("01/02/2006")
	f.PreparedDate = f.MessageDate
	f.Handling = "IMMEDIATE"
	f.ToLocation = "County EOC"
	return f
}

func makeF() *SituationReport {
	const fieldCount = 80
	f := SituationReport{BaseMessage: message.BaseMessage{Type: &Type}}
	f.FSubject = &f.Jurisdiction
	f.FBody = &f.AdditionalInformation
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.AddHeaderFields(&f.BaseMessage, &basePDFRenderers)
	f.Fields = append(f.Fields,
		message.NewDateField(true, &message.Field{
			Label:       "Prepared Date",
			Value:       &f.PreparedDate,
			Presence:    message.Required,
			PIFOTag:     "20d.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.40, Y: 268.92, W: 105.84, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the situation report was prepared.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Prepared Time",
			Value:       &f.PreparedTime,
			Presence:    message.Required,
			PIFOTag:     "20t.",
			PDFRenderer: &message.PDFTextRenderer{X: 307.92, Y: 268.92, W: 251.52, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of date when the situation report was prepared.  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Prepared Date/Time",
			EditHelp: `This is the date and time when the situation report was prepared, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.PreparedDate, &f.PreparedTime),
		message.NewTextField(&message.Field{
			Label:       "Jurisdiction",
			Value:       &f.Jurisdiction,
			Choices:     message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara (City)", "Saratoga", "Sunnyvale", "Santa Clara County", "County unincorporated"},
			Presence:    message.Required,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{X: 116.40, Y: 287.40, W: 172.80, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the name of the jurisdiction reporting the situation.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Incident Name",
			Value:       &f.IncidentName,
			Presence:    message.Required,
			PIFOTag:     "22.",
			PDFRenderer: &message.PDFTextRenderer{X: 378.36, Y: 287.40, W: 181.08, H: 11.04, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   38,
			EditHelp:    `This is the name of the incident for which the situation is being reported.  It is required.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Emergency Declaration",
			Value:   &f.EmergencyDeclaration,
			PIFOTag: "23.",
			Choices: message.Choices{"Unknown", "Yes", "No"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Unknown": {192.86, 310.77},
				"Yes":     {300.86, 310.77},
				"No":      {408.86, 310.77},
			}},
			EditHelp: `This indicates whether a declaration of emergency has been proclaimed.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "EOC Activation",
			Value:   &f.EOCActivation,
			PIFOTag: "24.",
			Choices: message.Choices{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Normal":       {156.86, 328.05},
				"Duty Officer": {264.86, 328.05},
				"Monitor":      {408.86, 328.05},
				"Partial":      {156.86, 343.65},
				"Full":         {264.86, 343.65},
			}},
			EditHelp: `This indicates the jurisdiction EOC's activation level.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Incident Description",
			Value:       &f.IncidentDescription,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 367.92, W: 514.80, H: 60.00, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This is a description of the incident.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Additional Incident Information",
			Value:       &f.AdditionalInformation,
			PIFOTag:     "26.",
			PDFRenderer: &message.PDFTextRenderer{X: 44.64, Y: 445.32, W: 514.80, H: 37.32, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This field provides additional information about the incident.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Prepared By",
			Value:       &f.PreparedBy,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{X: 106.92, Y: 516.48, W: 183.84, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   39,
			EditHelp:    `This is the name of the person who prepared the situation report.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Prepared By Phone",
			Value:       &f.PreparedByPhone,
			PIFOTag:     "31.",
			PDFRenderer: &message.PDFTextRenderer{X: 386.16, Y: 516.48, W: 173.28, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the phone number of the person who prepared the situation report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Prepared By Email",
			Value:       &f.PreparedByEmail,
			PIFOTag:     "32.",
			PDFRenderer: &message.PDFTextRenderer{X: 106.92, Y: 534.48, W: 452.52, H: 11.52, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   94,
			EditHelp:    `This is the email address of the person who prepared the situation report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Government Office Status",
			Value:   &f.OfficeStatus,
			PIFOTag: "50.",
			Choices: message.Choices{"Unknown", "Open", "Closed"},
			PDFRenderer: &message.PDFRadioRenderer{Radius: 3, Points: map[string][]float64{
				"Unknown": {156.86, 584.01},
				"Open":    {264.86, 584.01},
				"Closed":  {372.86, 584.01},
			}},
			EditHelp: `This is the status of the jurisdiction government office(s).`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Expected to Open Date",
			Value:       &f.ExpectedToOpenDate,
			PIFOTag:     "51d.",
			PDFRenderer: &message.PDFTextRenderer{X: 156.00, Y: 595.44, W: 49.80, H: 14.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the government office is expected to open.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.ExpectedToOpenTime,
			PIFOTag:     "51t.",
			PDFRenderer: &message.PDFTextRenderer{X: 242.52, Y: 595.44, W: 48.24, H: 14.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of day when the government office is expected to open.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Expected to Open",
			EditHelp: `This is the date and time when the government office is expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ExpectedToOpenDate, &f.ExpectedToOpenTime),
		message.NewDateField(true, &message.Field{
			Label:       "Expected to Close Date",
			Value:       &f.ExpectedToCloseDate,
			PIFOTag:     "52d.",
			PDFRenderer: &message.PDFTextRenderer{X: 411.55, Y: 595.44, W: 50, H: 14.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date when the government office is expected to close.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.ExpectedToCloseTime,
			PIFOTag:     "52t.",
			PDFRenderer: &message.PDFTextRenderer{X: 492.12, Y: 595.44, W: 67.32, H: 14.64, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of day when the government office is expected to close.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Expected to Close",
			EditHelp: `This is the date and time when the government office is expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.ExpectedToCloseDate, &f.ExpectedToCloseTime),
		message.NewTextField(&message.Field{
			Label:       "Approved By",
			Value:       &f.ApprovedBy,
			PIFOTag:     "60.",
			PDFRenderer: &message.PDFTextRenderer{X: 109.92, Y: 644.16, W: 180.84, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   38,
			EditHelp:    `This is the name of the person who approved the situation report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Phone",
			Value:       &f.ApprovedByPhone,
			PIFOTag:     "61.",
			PDFRenderer: &message.PDFTextRenderer{X: 377.76, Y: 644.16, W: 181.68, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   38,
			EditHelp:    `This is the phone number of the person who approved the situation report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Position",
			Value:       &f.ApprovedByPosition,
			PIFOTag:     "62.",
			PDFRenderer: &message.PDFTextRenderer{X: 109.92, Y: 662.16, W: 180.84, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   38,
			EditHelp:    `This is the ICS position of the person who approved the situation report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Approved By Location",
			Value:       &f.ApprovedByLocation,
			PIFOTag:     "63.",
			PDFRenderer: &message.PDFTextRenderer{X: 377.76, Y: 662.16, W: 181.68, H: 11.52, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   38,
			EditHelp:    `This is the location of the person who approved the situation report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Communications Status",
			Value:   &f.CommunicationsStatus,
			PIFOTag: "70a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 84.45},
				"Problem": {284.90, 84.45},
				"Failure": {356.78, 84.45},
				"Unknown": {420.98, 84.45},
			}},
			EditHelp: `This is the current situation with respect to communications.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Communications Comments",
			Value:       &f.CommunicationsComments,
			PIFOTag:     "70b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 96.00, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to communications.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Debris Status",
			Value:   &f.DebrisStatus,
			PIFOTag: "71a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 127.84},
				"Problem": {284.90, 127.84},
				"Failure": {356.78, 127.84},
				"Unknown": {420.98, 127.84},
			}},
			EditHelp: `This is the current situation with respect to debris.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Debris Comments",
			Value:       &f.DebrisComments,
			PIFOTag:     "71b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 139.39, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to debris.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Flooding Status",
			Value:   &f.FloodingStatus,
			PIFOTag: "72a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 171.24},
				"Problem": {284.90, 171.24},
				"Failure": {356.78, 171.24},
				"Unknown": {420.98, 171.24},
			}},
			EditHelp: `This is the current situation with respect to flooding.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Flooding Comments",
			Value:       &f.FloodingComments,
			PIFOTag:     "72b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 182.79, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to flooding.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "HazMat Status",
			Value:   &f.HazMatStatus,
			PIFOTag: "73a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 214.63},
				"Problem": {284.90, 214.63},
				"Failure": {356.78, 214.63},
				"Unknown": {420.98, 214.63},
			}},
			EditHelp: `This is the current situation with respect to hazardous materials.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "HazMat Comments",
			Value:       &f.HazMatComments,
			PIFOTag:     "73b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 226.18, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to hazardous materials.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Emergency Services Status",
			Value:   &f.EmergencySvcsStatus,
			PIFOTag: "74a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 258.03},
				"Problem": {284.90, 258.03},
				"Failure": {356.78, 258.03},
				"Unknown": {420.98, 258.03},
			}},
			EditHelp: `This is the current situation with respect to emergency services.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Emergency Services Comments",
			Value:       &f.EmergencySvcsComments,
			PIFOTag:     "74b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 269.58, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to emergency services.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Casualties Status",
			Value:   &f.CasualtiesStatus,
			PIFOTag: "75a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 301.42},
				"Problem": {284.90, 301.42},
				"Failure": {356.78, 301.42},
				"Unknown": {420.98, 301.42},
			}},
			EditHelp: `This is the current situation with respect to casualties.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Casualties Comments",
			Value:       &f.CasualtiesComments,
			PIFOTag:     "75b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 312.97, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to casualties.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Utilities (Gas) Status",
			Value:   &f.GasStatus,
			PIFOTag: "76a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 344.81},
				"Problem": {284.90, 344.81},
				"Failure": {356.78, 344.81},
				"Unknown": {420.98, 344.81},
			}},
			EditHelp: `This is the current situation with respect to gas utilities.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Utilities (Gas) Comments",
			Value:       &f.GasComments,
			PIFOTag:     "76b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 356.36, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to gas utilities.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Utilities (Electric) Status",
			Value:   &f.ElectricStatus,
			PIFOTag: "77a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 388.21},
				"Problem": {284.90, 388.21},
				"Failure": {356.78, 388.21},
				"Unknown": {420.98, 388.21},
			}},
			EditHelp: `This is the current situation with respect to electric utilities.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Utilities (Electric) Comments",
			Value:       &f.ElectricComments,
			PIFOTag:     "77b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 399.76, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to electric utilities.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure (Power) Status",
			Value:   &f.PowerStatus,
			PIFOTag: "78a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 431.60},
				"Problem": {284.90, 431.60},
				"Failure": {356.78, 431.60},
				"Unknown": {420.98, 431.60},
			}},
			EditHelp: `This is the current situation with respect to power infrastructure.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Infrastructure (Power) Comments",
			Value:       &f.PowerComments,
			PIFOTag:     "78b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 443.15, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to power infrastructure.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure (Water) Status",
			Value:   &f.WaterStatus,
			PIFOTag: "79a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 474.99},
				"Problem": {284.90, 474.99},
				"Failure": {356.78, 474.99},
				"Unknown": {420.98, 474.99},
			}},
			EditHelp: `This is the current situation with respect to water infrastructure.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Infrastructure (Water) Comments",
			Value:       &f.WaterComments,
			PIFOTag:     "79b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 486.54, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to water infrastructure.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure (Sewer) Status",
			Value:   &f.SewerStatus,
			PIFOTag: "80a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 518.39},
				"Problem": {284.90, 518.39},
				"Failure": {356.78, 518.39},
				"Unknown": {420.98, 518.39},
			}},
			EditHelp: `This is the current situation with respect to sewer infrastructure.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Infrastructure (Sewer) Comments",
			Value:       &f.SewerComments,
			PIFOTag:     "80b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 529.94, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to sewer infrastructure.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Search and Rescue Status",
			Value:   &f.SARStatus,
			PIFOTag: "81a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 561.78},
				"Problem": {284.90, 561.78},
				"Failure": {356.78, 561.78},
				"Unknown": {420.98, 561.78},
			}},
			EditHelp: `This is the current situation with respect to search and rescue.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Search and Rescue Comments",
			Value:       &f.SARComments,
			PIFOTag:     "81b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 573.33, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to search and rescue.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Transportation (Roads) Status",
			Value:   &f.RoadsStatus,
			PIFOTag: "82a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 605.18},
				"Problem": {284.90, 605.18},
				"Failure": {356.78, 605.18},
				"Unknown": {420.98, 605.18},
			}},
			EditHelp: `This is the current situation with respect to roads.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Transportation (Roads) Comments",
			Value:       &f.RoadsComments,
			PIFOTag:     "82b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 616.73, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to roads.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Transportation (Bridges) Status",
			Value:   &f.BridgesStatus,
			PIFOTag: "83a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 2, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 648.57},
				"Problem": {284.90, 648.57},
				"Failure": {356.78, 648.57},
				"Unknown": {420.98, 648.57},
			}},
			EditHelp: `This is the current situation with respect to bridges.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Transportation (Bridges) Comments",
			Value:       &f.BridgesComments,
			PIFOTag:     "83b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 226.25, Y: 660.12, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to bridges.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Civil Unrest Status",
			Value:   &f.UnrestStatus,
			PIFOTag: "84a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 57.45},
				"Problem": {284.90, 57.45},
				"Failure": {356.78, 57.45},
				"Unknown": {420.98, 57.45},
			}},
			EditHelp: `This is the current situation with respect to civil unrest.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Civil Unrest Comments",
			Value:       &f.UnrestComments,
			PIFOTag:     "84b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 226.25, Y: 69.00, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to civil unrest.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Animal Issues Status",
			Value:   &f.AnimalStatus,
			PIFOTag: "85a.",
			Choices: message.Choices{"Normal", "Problem", "Failure", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Radius: 3, Points: map[string][]float64{
				"Normal":  {215.18, 100.89},
				"Problem": {284.90, 100.89},
				"Failure": {356.78, 100.89},
				"Unknown": {420.98, 100.89},
			}},
			EditHelp: `This is the current situation with respect to animal issues.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Animal Issues Comments",
			Value:       &f.AnimalComments,
			PIFOTag:     "85b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 226.25, Y: 112.39, W: 333.07, H: 20.88, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   69,
			EditHelp:    `These are comments on the current situation with respect to animal issues.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Lifeline Status",
			Value:    &f.LifelineStatus,
			Presence: message.Required,
			PIFOTag:  "90.",
			Choices:  message.Choices{"Stable", "Stabilizing", "Unstable", "Unknown"},
			PDFRenderer: &message.PDFRadioRenderer{Page: 3, Radius: 3, Points: map[string][]float64{
				"Stable":      {156.86, 171.33},
				"Stabilizing": {220.94, 171.33},
				"Unstable":    {303.14, 171.33},
				"Unknown":     {382.10, 171.33},
			}},
			EditHelp: `This is the status of the jurisdiction's response to this incident.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Lifeline Update",
			Value:       &f.LifelineUpdate,
			Presence:    message.Required,
			PIFOTag:     "91.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 44.64, Y: 195.12, W: 514.68, H: 75.84, B: 999, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This provides details of the jurisdiction's response and recovery to the incident.  It is required.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:       "Unmet Needs",
			Value:       &f.UnmetNeeds,
			Presence:    message.Required,
			PIFOTag:     "92.",
			PDFRenderer: &message.PDFTextRenderer{Page: 3, X: 44.64, Y: 290.16, W: 514.68, H: 89.16, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   108,
			EditHelp:    `This lists shortcomings from the jurisdiction's response and recovery.  It is required.`,
		}),
	)
	f.AddFooterFields(&f.BaseMessage, &basePDFRenderers)
	if len(f.Fields) > fieldCount {
		panic("update SituationReport fieldCount")
	}
	return &f
}

func decode(_ *envelope.Envelope, _ string, form *message.PIFOForm, _ int) message.Message {
	var df *SituationReport

	if form == nil || form.HTMLIdent != Type.HTML || form.FormVersion != Type.Version {
		return nil
	}
	df = makeF()
	message.DecodeForm(form, df)
	return df
}
