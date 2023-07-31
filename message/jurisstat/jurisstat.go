// Package jurisstat defines the Santa Clara County OA Jurisdiction Status Form
// message type.
package jurisstat

import (
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/baseform"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an OA jurisdiction status form.
var Type = message.Type{
	Tag:     "JurisStat",
	Name:    "OA jurisdiction status form",
	Article: "an",
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
var versions = []*basemsg.FormVersion{
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

// JurisStat holds an OA jurisdiction status form.
type JurisStat struct {
	basemsg.BaseMessage
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

var pdfBase []byte

func create(version *basemsg.FormVersion) message.Message {
	var f = JurisStat{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		PDFBase:     pdfBase,
		PDFFontSize: 10,
		Form:        version,
	}}
	if version.Version < "2.2" {
		f.MessageType = &OldType
	}
	f.BaseMessage.FSubject = &f.Jurisdiction
	f.BaseMessage.FReportType = &f.ReportType
	f.BaseMessage.FBody = &f.CommunicationsComments
	var basePDFMaps = baseform.DefaultPDFMaps
	basePDFMaps.OriginMsgID = basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
		return []basemsg.PDFField{
			{Name: "Origin Msg Nbr", Value: f.OriginMsgID},
			{Name: "Origin Msg Nbr Copy", Value: f.OriginMsgID},
		}
	})
	f.Fields = make([]*basemsg.Field, 0, 80)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFMaps)
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:     "Report Type",
			Value:     &f.ReportType,
			Choices:   basemsg.Choices{"Update", "Complete"},
			Presence:  basemsg.Required,
			PIFOTag:   "19.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Report Type", "", "Off"},
			EditWidth: 7,
			EditHelp:  `This indicates whether the form should "Update" the previous status report for the jurisdiction, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
		},
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:     "Jurisdiction Name",
				Value:     &f.Jurisdiction,
				Choices:   basemsg.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:  basemsg.Required,
				PIFOTag:   "21.",
				Compare:   common.CompareText,
				PDFMap:    basemsg.PDFName("Jurisdiction Name"),
				EditWidth: 42,
				EditHelp:  ``,
			},
		)
	} else {
		f.Fields = append(f.Fields,
			&basemsg.Field{
				Label:      "Jurisdiction Code",
				Value:      &f.JurisdictionCode,
				Choices:    basemsg.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:   basemsg.Required,
				PIFOTag:    "21.",
				PIFOValid:  basemsg.ValidRestricted,
				TableValue: basemsg.OmitFromTable,
			},
			&basemsg.Field{
				Label:     "Jurisdiction Name",
				Value:     &f.Jurisdiction,
				Choices:   basemsg.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:  basemsg.Required,
				PIFOTag:   "22.",
				Compare:   common.CompareText,
				PDFMap:    basemsg.PDFName("Jurisdiction Name"),
				EditWidth: 42,
				EditHelp:  `This is the name of the jurisdiction being described by the form.  It is required.`,
				EditApply: func(field *basemsg.Field, v string) {
					f.Jurisdiction = v
					if v == "" || field.Choices.IsPIFO(v) {
						f.JurisdictionCode = v
					} else {
						f.JurisdictionCode = "Unincorporated"
					}
				},
			},
		)
	}
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:     "EOC Phone",
			Value:     &f.EOCPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "23.",
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("EOC Phone"),
			EditWidth: 34,
			EditHelp:  `This is the phone number of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "EOC Fax",
			Value:     &f.EOCFax,
			PIFOTag:   "24.",
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("EOC Fax"),
			EditWidth: 37,
			EditHelp:  `This is the fax number of the jurisdiction's Emergency Operations Center (EOC).`,
		},
		&basemsg.Field{
			Label:     "Primary EM Contact Name",
			Value:     &f.PriEMContactName,
			Presence:  f.requiredForComplete,
			PIFOTag:   "25.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Pri EM Contact Name"),
			EditWidth: 27,
			EditHelp:  `This is the name of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Primary EM Contact Phone",
			Value:     &f.PriEMContactPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "26.",
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("Pri EM Contact Phone"),
			EditWidth: 26,
			EditHelp:  `This is the phone number of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Secondary EM Contact Name",
			Value:     &f.SecEMContactName,
			PIFOTag:   "27.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Sec EM Contact Name"),
			EditWidth: 26,
			EditHelp:  `This is the name of the secondary emergency manager of the jurisdiction.`,
		},
		&basemsg.Field{
			Label:     "Secondary EM Contact Phone",
			Value:     &f.SecEMContactPhone,
			PIFOTag:   "28.",
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("Sec EM Contact Phone"),
			EditWidth: 26,
			EditHelp:  `This is the phone number of the secondary emergency manager of the jurisdiction.`,
		},
		&basemsg.Field{
			Label:     "Govt. Office Status",
			Value:     &f.OfficeStatus,
			Choices:   basemsg.Choices{"Unknown", "Open", "Closed"},
			Presence:  f.requiredForComplete,
			PIFOTag:   "29.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Office Status", "", "Off"},
			EditWidth: 7,
			EditHelp:  `This indicates whether the jurisdiction's regular business offices are open.  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:      "Govt. Office Expected Open Date",
			Value:      &f.GovExpectedOpenDate,
			PIFOTag:    "30.",
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("Office Open Date"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Govt. Office Expected Open Time",
			Value:      &f.GovExpectedOpenTime,
			PIFOTag:    "31.",
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("Office Open Time"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Govt. Office Expected to Open",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.GovExpectedOpenDate, f.GovExpectedOpenTime, " ")
			},
			EditWidth: 16,
			EditHelp:  `This is the date and time when the jurisdiction's regular business offices are expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.GovExpectedOpenDate, f.GovExpectedOpenTime)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.GovExpectedOpenDate, &f.GovExpectedOpenTime, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.GovExpectedOpenDate, f.GovExpectedOpenTime)
			},
		},
		&basemsg.Field{
			Label:      "Govt. Office Expected Close Date",
			Value:      &f.GovExpectedCloseDate,
			PIFOTag:    "32.",
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("Office Close Date"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Govt. Office Expected Close Time",
			Value:      &f.GovExpectedCloseTime,
			PIFOTag:    "33.",
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("Office Close Time"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Govt. Office Expected to Close",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.GovExpectedCloseDate, f.GovExpectedCloseTime, " ")
			},
			EditWidth: 16,
			EditHelp:  `This is the date and time when the jurisdiction's regular business offices are expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.GovExpectedCloseDate, f.GovExpectedCloseTime)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.GovExpectedCloseDate, &f.GovExpectedCloseTime, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.GovExpectedCloseDate, f.GovExpectedCloseTime)
			},
		},
		&basemsg.Field{
			Label:     "EOC Open",
			Value:     &f.EOCOpen,
			Choices:   basemsg.Choices{"Unknown", "Yes", "No"},
			Presence:  f.requiredForComplete,
			PIFOTag:   "34.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"EOC Open", "", "Off"},
			EditWidth: 7,
			EditHelp:  `This indicates whether the jurisdiction's Emergency Operations Center (EOC) is open.  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "EOC Activation Level",
			Value:     &f.EOCActivationLevel,
			Choices:   basemsg.Choices{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
			Presence:  f.requiredForComplete,
			PIFOTag:   "35.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Activation", "", "Off"},
			EditWidth: 12,
			EditHelp:  `This indicates the activation level of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:      "EOC Expected to Open Date",
			Value:      &f.EOCExpectedOpenDate,
			PIFOTag:    "36.",
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("EOC Open Date"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "EOC Expected to Open Time",
			Value:      &f.EOCExpectedOpenTime,
			PIFOTag:    "37.",
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("EOC Open Time"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "EOC Expected to Open",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.EOCExpectedOpenDate, f.EOCExpectedOpenTime, " ")
			},
			EditWidth: 16,
			EditHelp:  `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.EOCExpectedOpenDate, f.EOCExpectedOpenTime)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.EOCExpectedOpenDate, &f.EOCExpectedOpenTime, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.EOCExpectedOpenDate, f.EOCExpectedOpenTime)
			},
		},
		&basemsg.Field{
			Label:      "EOC Expected to Close Date",
			Value:      &f.EOCExpectedCloseDate,
			PIFOTag:    "38.",
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("EOC Close Date"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "EOC Expected to Close Time",
			Value:      &f.EOCExpectedCloseTime,
			PIFOTag:    "39.",
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("EOC Close Time"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "EOC Expected to Close",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.EOCExpectedCloseDate, f.EOCExpectedCloseTime, " ")
			},
			EditWidth: 16,
			EditHelp:  `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to close, in MM/DD/YYYY HH:MM format (24-hour clock).`,
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.EOCExpectedCloseDate, f.EOCExpectedCloseTime)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.EOCExpectedCloseDate, &f.EOCExpectedCloseTime, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.EOCExpectedCloseDate, f.EOCExpectedCloseTime)
			},
		},
		&basemsg.Field{
			Label:     "State of Emergency",
			Value:     &f.StateOfEmergency,
			Choices:   basemsg.Choices{"Unknown", "Yes", "No"},
			Presence:  f.requiredForComplete,
			PIFOTag:   "40.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"State of Emergency", "", "Off"},
			EditWidth: 7,
			EditHelp:  `This indicates whether the jurisdiction has a declared state of emergency.  It is required when "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label: "How SOE Sent",
			Value: &f.HowSOESent,
			Presence: func() (basemsg.Presence, string) {
				if f.StateOfEmergency == "Yes" {
					return basemsg.PresenceRequired, `when "State of Emergency" is "Yes"`
				} else {
					return basemsg.PresenceNotAllowed, `when "State of Emergency" is not "Yes"`
				}
			},
			PIFOTag:   "99.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Attachment"),
			EditWidth: 58,
			EditHelp:  `This describes where and how the jurisdiction's "state of emergency" declaration was delivered.`,
		},
		&basemsg.Field{
			Label:     "Communications",
			Value:     &f.Communications,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "41.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Communications", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to communications.`,
		},
		&basemsg.Field{
			Label:     "Communications: Comments",
			Value:     &f.CommunicationsComments,
			PIFOTag:   "41.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Comm Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to communications.`,
		},
		&basemsg.Field{
			Label:     "Debris",
			Value:     &f.Debris,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "42.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Debris", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to debris.`,
		},
		&basemsg.Field{
			Label:     "Debris: Comments",
			Value:     &f.DebrisComments,
			PIFOTag:   "42.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Debris Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to debris.`,
		},
		&basemsg.Field{
			Label:     "Flooding",
			Value:     &f.Flooding,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "43.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Flooding", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to flooding.`,
		},
		&basemsg.Field{
			Label:     "Flooding: Comments",
			Value:     &f.FloodingComments,
			PIFOTag:   "43.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Flood Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to flooding.`,
		},
		&basemsg.Field{
			Label:     "Hazmat",
			Value:     &f.Hazmat,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "44.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Hazmat", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to hazmat.`,
		},
		&basemsg.Field{
			Label:     "Hazmat: Comments",
			Value:     &f.HazmatComments,
			PIFOTag:   "44.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Hazmat Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to hazmat.`,
		},
		&basemsg.Field{
			Label:     "Emergency Services",
			Value:     &f.EmergencyServices,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "45.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Em Svcs", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to emergency services.`,
		},
		&basemsg.Field{
			Label:     "Emergency Services: Comments",
			Value:     &f.EmergencyServicesComments,
			PIFOTag:   "45.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Em Svcs Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to emergency services.`,
		},
		&basemsg.Field{
			Label:     "Casualties",
			Value:     &f.Casualties,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "46.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Casualties", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to casualties.`,
		},
		&basemsg.Field{
			Label:     "Casualties: Comments",
			Value:     &f.CasualtiesComments,
			PIFOTag:   "46.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Casualties Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to casualties.`,
		},
		&basemsg.Field{
			Label:     "Utilities Gas",
			Value:     &f.UtilitiesGas,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "47.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Util Gas", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to utilities (gas).`,
		},
		&basemsg.Field{
			Label:     "Utilities Gas: Comments",
			Value:     &f.UtilitiesGasComments,
			PIFOTag:   "47.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Util Gas Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to utilities (gas).`,
		},
		&basemsg.Field{
			Label:     "Utilities Electric",
			Value:     &f.UtilitiesElectric,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "48.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Util Elec", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to utilities (electric).`,
		},
		&basemsg.Field{
			Label:     "Utilities Electric: Comments",
			Value:     &f.UtilitiesElectricComments,
			PIFOTag:   "48.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Util Elec Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to utilities (electric).`,
		},
		&basemsg.Field{
			Label:     "Infrastructure Power",
			Value:     &f.InfrastructurePower,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "49.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Infra Pwr", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to infrastructure (power).`,
		},
		&basemsg.Field{
			Label:     "Infrastructure Power: Comments",
			Value:     &f.InfrastructurePowerComments,
			PIFOTag:   "49.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Infra Pwr Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to infrastructure (power).`,
		},
		&basemsg.Field{
			Label:     "Infrastructure Water",
			Value:     &f.InfrastructureWater,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "50.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Infra Water", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to infrastructure (water).`,
		},
		&basemsg.Field{
			Label:     "Infrastructure Water: Comments",
			Value:     &f.InfrastructureWaterComments,
			PIFOTag:   "50.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Infra Water Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to infrastructure (water).`,
		},
		&basemsg.Field{
			Label:     "Infrastructure Sewer",
			Value:     &f.InfrastructureSewer,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "51.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Infra Sewer", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to infrastructure (sewer).`,
		},
		&basemsg.Field{
			Label:     "Infrastructure Sewer: Comments",
			Value:     &f.InfrastructureSewerComments,
			PIFOTag:   "51.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Infra Sewer Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to infrastructure (sewer).`,
		},
		&basemsg.Field{
			Label:     "Search And Rescue",
			Value:     &f.SearchAndRescue,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "52.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"SAR", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to search and rescue.`,
		},
		&basemsg.Field{
			Label:     "Search And Rescue: Comments",
			Value:     &f.SearchAndRescueComments,
			PIFOTag:   "52.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("SAR Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to search and rescue.`,
		},
		&basemsg.Field{
			Label:     "Transportation Roads",
			Value:     &f.TransportationRoads,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "53.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Trans Roads", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to transportation (roads).`,
		},
		&basemsg.Field{
			Label:     "Transportation Roads: Comments",
			Value:     &f.TransportationRoadsComments,
			PIFOTag:   "53.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Trans Roads Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to transportation (roads).`,
		},
		&basemsg.Field{
			Label:     "Transportation Bridges",
			Value:     &f.TransportationBridges,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "54.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Trans Bridges", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to transportation (bridges).`,
		},
		&basemsg.Field{
			Label:     "Transportation Bridges: Comments",
			Value:     &f.TransportationBridgesComments,
			PIFOTag:   "54.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Trans Bridges Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to transportation (bridges).`,
		},
		&basemsg.Field{
			Label:     "Civil Unrest",
			Value:     &f.CivilUnrest,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "55.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Civil", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to civil unrest.`,
		},
		&basemsg.Field{
			Label:     "Civil Unrest: Comments",
			Value:     &f.CivilUnrestComments,
			PIFOTag:   "55.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Civil Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to civil unrest.`,
		},
		&basemsg.Field{
			Label:     "Animal Issues",
			Value:     &f.AnimalIssues,
			Choices:   basemsg.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:   "56.0.",
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFNameMap{"Animal", "", "Off"},
			EditWidth: 9,
			EditHelp:  `This describes the current situation status with respect to animal issues.`,
		},
		&basemsg.Field{
			Label:     "Animal Issues: Comments",
			Value:     &f.AnimalIssuesComments,
			PIFOTag:   "56.1.",
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Animal Comment"),
			EditWidth: 60,
			Multiline: true,
			EditHelp:  `These are comments on the current situation status with respect to animal issues.`,
		},
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFMaps)
	return &f
}

func (f *JurisStat) requiredForComplete() (basemsg.Presence, string) {
	if f.ReportType == "Complete" {
		return basemsg.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return basemsg.PresenceOptional, ""
}

func decode(subject, body string) (f *JurisStat) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-oa-muni-status.html") {
		return nil
	}
	return basemsg.Decode(body, versions, create).(*JurisStat)
}
