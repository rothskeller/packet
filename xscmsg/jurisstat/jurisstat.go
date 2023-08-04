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
	PDFBase:     pdfBase,
	PDFFontSize: 10,
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

var pdfBase []byte

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
	var basePDFMaps = baseform.DefaultPDFMaps
	basePDFMaps.OriginMsgID = message.PDFMapFunc(func(*message.Field) []message.PDFField {
		return []message.PDFField{
			{Name: "Origin Msg Nbr", Value: f.OriginMsgID},
			{Name: "Origin Msg Nbr Copy", Value: f.OriginMsgID},
		}
	})
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &basePDFMaps)
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  message.Choices{"Update", "Complete"},
			Presence: message.Required,
			PIFOTag:  "19.",
			PDFMap:   message.PDFNameMap{"Report Type", "", "Off"},
			EditHelp: `This indicates whether the form should "Update" the previous status report for the jurisdiction, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
		}),
	)
	if f.Form.Version < "2.2" {
		f.Fields = append(f.Fields,
			message.NewRestrictedField(&message.Field{
				Label:     "Jurisdiction Name",
				Value:     &f.Jurisdiction,
				Choices:   message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:  message.Required,
				PIFOTag:   "21.",
				PDFMap:    message.PDFName("Jurisdiction Name"),
				EditWidth: 42,
				EditHelp:  `This is the name of the jurisdiction being described by the form.  It is required.`,
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
				Label:     "Jurisdiction Name",
				Value:     &f.Jurisdiction,
				Choices:   message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Presence:  message.Required,
				PIFOTag:   "22.",
				PDFMap:    message.PDFName("Jurisdiction Name"),
				EditWidth: 42,
				EditHelp:  `This is the name of the jurisdiction being described by the form.  It is required.`,
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
			Label:     "EOC Phone",
			Value:     &f.EOCPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "23.",
			PDFMap:    message.PDFName("EOC Phone"),
			EditWidth: 34,
			EditHelp:  `This is the phone number of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:     "EOC Fax",
			Value:     &f.EOCFax,
			PIFOTag:   "24.",
			PDFMap:    message.PDFName("EOC Fax"),
			EditWidth: 37,
			EditHelp:  `This is the fax number of the jurisdiction's Emergency Operations Center (EOC).`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Primary EM Contact Name",
			Value:     &f.PriEMContactName,
			Presence:  f.requiredForComplete,
			PIFOTag:   "25.",
			PDFMap:    message.PDFName("Pri EM Contact Name"),
			EditWidth: 27,
			EditHelp:  `This is the name of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:     "Primary EM Contact Phone",
			Value:     &f.PriEMContactPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "26.",
			PDFMap:    message.PDFName("Pri EM Contact Phone"),
			EditWidth: 26,
			EditHelp:  `This is the phone number of the primary emergency manager of the jurisdiction.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewTextField(&message.Field{
			Label:     "Secondary EM Contact Name",
			Value:     &f.SecEMContactName,
			PIFOTag:   "27.",
			PDFMap:    message.PDFName("Sec EM Contact Name"),
			EditWidth: 26,
			EditHelp:  `This is the name of the secondary emergency manager of the jurisdiction.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:     "Secondary EM Contact Phone",
			Value:     &f.SecEMContactPhone,
			PIFOTag:   "28.",
			PDFMap:    message.PDFName("Sec EM Contact Phone"),
			EditWidth: 26,
			EditHelp:  `This is the phone number of the secondary emergency manager of the jurisdiction.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Govt. Office Status",
			Value:    &f.OfficeStatus,
			Choices:  message.Choices{"Unknown", "Open", "Closed"},
			Presence: f.requiredForComplete,
			PIFOTag:  "29.",
			PDFMap:   message.PDFNameMap{"Office Status", "", "Off"},
			EditHelp: `This indicates whether the jurisdiction's regular business offices are open.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewDateWithTimeField(&message.Field{
			Label:   "Govt. Office Expected Open Date",
			Value:   &f.GovExpectedOpenDate,
			PIFOTag: "30.",
			PDFMap:  message.PDFName("Office Open Date"),
		}),
		message.NewTimeWithDateField(&message.Field{
			Label:   "Govt. Office Expected Open Time",
			Value:   &f.GovExpectedOpenTime,
			PIFOTag: "31.",
			PDFMap:  message.PDFName("Office Open Time"),
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Govt. Office Expected to Open",
			EditHelp: `This is the date and time when the jurisdiction's regular business offices are expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
		}, &f.GovExpectedOpenDate, &f.GovExpectedOpenTime),
		message.NewDateWithTimeField(&message.Field{
			Label:   "Govt. Office Expected Close Date",
			Value:   &f.GovExpectedCloseDate,
			PIFOTag: "32.",
			PDFMap:  message.PDFName("Office Close Date"),
		}),
		message.NewTimeWithDateField(&message.Field{
			Label:   "Govt. Office Expected Close Time",
			Value:   &f.GovExpectedCloseTime,
			PIFOTag: "33.",
			PDFMap:  message.PDFName("Office Close Time"),
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
			PDFMap:   message.PDFNameMap{"EOC Open", "", "Off"},
			EditHelp: `This indicates whether the jurisdiction's Emergency Operations Center (EOC) is open.  It is required when "Report Type" is "Complete".`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "EOC Activation Level",
			Value:    &f.EOCActivationLevel,
			Choices:  message.Choices{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
			Presence: f.requiredForComplete,
			PIFOTag:  "35.",
			PDFMap:   message.PDFNameMap{"Activation", "", "Off"},
			EditHelp: `This indicates the activation level of the jurisdiction's Emergency Operations Center (EOC).  It is required when "Report Type" is "Complete".`,
		}),
		message.NewDateWithTimeField(&message.Field{
			Label:   "EOC Expected to Open Date",
			Value:   &f.EOCExpectedOpenDate,
			PIFOTag: "36.",
			PDFMap:  message.PDFName("EOC Open Date"),
		}),
		message.NewTimeWithDateField(&message.Field{
			Label:   "EOC Expected to Open Time",
			Value:   &f.EOCExpectedOpenTime,
			PIFOTag: "37.",
			PDFMap:  message.PDFName("EOC Open Time"),
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "EOC Expected to Open",
			EditHelp: `This is the date and time when the jurisdiction's Emergency Operations Center (EOC) is expected to open, in MM/DD/YYYY HH:MM format (24-hour clock).`,
		}, &f.EOCExpectedOpenDate, &f.EOCExpectedOpenTime),
		message.NewDateWithTimeField(&message.Field{
			Label:   "EOC Expected to Close Date",
			Value:   &f.EOCExpectedCloseDate,
			PIFOTag: "38.",
			PDFMap:  message.PDFName("EOC Close Date"),
		}),
		message.NewTimeWithDateField(&message.Field{
			Label:   "EOC Expected to Close Time",
			Value:   &f.EOCExpectedCloseTime,
			PIFOTag: "39.",
			PDFMap:  message.PDFName("EOC Close Time"),
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
			PDFMap:   message.PDFNameMap{"State of Emergency", "", "Off"},
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
			PIFOTag:   "99.",
			PDFMap:    message.PDFName("Attachment"),
			EditWidth: 58,
			EditHelp:  `This describes where and how the jurisdiction's "state of emergency" declaration was delivered.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Communications",
			Value:    &f.Communications,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "41.0.",
			PDFMap:   message.PDFNameMap{"Communications", "", "Off"},
			EditHelp: `This describes the current situation status with respect to communications.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Communications: Comments",
			Value:     &f.CommunicationsComments,
			PIFOTag:   "41.1.",
			PDFMap:    message.PDFName("Comm Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to communications.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Debris",
			Value:    &f.Debris,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "42.0.",
			PDFMap:   message.PDFNameMap{"Debris", "", "Off"},
			EditHelp: `This describes the current situation status with respect to debris.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Debris: Comments",
			Value:     &f.DebrisComments,
			PIFOTag:   "42.1.",
			PDFMap:    message.PDFName("Debris Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to debris.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Flooding",
			Value:    &f.Flooding,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "43.0.",
			PDFMap:   message.PDFNameMap{"Flooding", "", "Off"},
			EditHelp: `This describes the current situation status with respect to flooding.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Flooding: Comments",
			Value:     &f.FloodingComments,
			PIFOTag:   "43.1.",
			PDFMap:    message.PDFName("Flood Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to flooding.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Hazmat",
			Value:    &f.Hazmat,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "44.0.",
			PDFMap:   message.PDFNameMap{"Hazmat", "", "Off"},
			EditHelp: `This describes the current situation status with respect to hazmat.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Hazmat: Comments",
			Value:     &f.HazmatComments,
			PIFOTag:   "44.1.",
			PDFMap:    message.PDFName("Hazmat Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to hazmat.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Emergency Services",
			Value:    &f.EmergencyServices,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "45.0.",
			PDFMap:   message.PDFNameMap{"Em Svcs", "", "Off"},
			EditHelp: `This describes the current situation status with respect to emergency services.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Emergency Services: Comments",
			Value:     &f.EmergencyServicesComments,
			PIFOTag:   "45.1.",
			PDFMap:    message.PDFName("Em Svcs Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to emergency services.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Casualties",
			Value:    &f.Casualties,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "46.0.",
			PDFMap:   message.PDFNameMap{"Casualties", "", "Off"},
			EditHelp: `This describes the current situation status with respect to casualties.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Casualties: Comments",
			Value:     &f.CasualtiesComments,
			PIFOTag:   "46.1.",
			PDFMap:    message.PDFName("Casualties Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to casualties.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Utilities Gas",
			Value:    &f.UtilitiesGas,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "47.0.",
			PDFMap:   message.PDFNameMap{"Util Gas", "", "Off"},
			EditHelp: `This describes the current situation status with respect to utilities (gas).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Utilities Gas: Comments",
			Value:     &f.UtilitiesGasComments,
			PIFOTag:   "47.1.",
			PDFMap:    message.PDFName("Util Gas Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to utilities (gas).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Utilities Electric",
			Value:    &f.UtilitiesElectric,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "48.0.",
			PDFMap:   message.PDFNameMap{"Util Elec", "", "Off"},
			EditHelp: `This describes the current situation status with respect to utilities (electric).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Utilities Electric: Comments",
			Value:     &f.UtilitiesElectricComments,
			PIFOTag:   "48.1.",
			PDFMap:    message.PDFName("Util Elec Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to utilities (electric).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Infrastructure Power",
			Value:    &f.InfrastructurePower,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "49.0.",
			PDFMap:   message.PDFNameMap{"Infra Pwr", "", "Off"},
			EditHelp: `This describes the current situation status with respect to infrastructure (power).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Infrastructure Power: Comments",
			Value:     &f.InfrastructurePowerComments,
			PIFOTag:   "49.1.",
			PDFMap:    message.PDFName("Infra Pwr Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to infrastructure (power).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Infrastructure Water",
			Value:    &f.InfrastructureWater,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "50.0.",
			PDFMap:   message.PDFNameMap{"Infra Water", "", "Off"},
			EditHelp: `This describes the current situation status with respect to infrastructure (water).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Infrastructure Water: Comments",
			Value:     &f.InfrastructureWaterComments,
			PIFOTag:   "50.1.",
			PDFMap:    message.PDFName("Infra Water Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to infrastructure (water).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Infrastructure Sewer",
			Value:    &f.InfrastructureSewer,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "51.0.",
			PDFMap:   message.PDFNameMap{"Infra Sewer", "", "Off"},
			EditHelp: `This describes the current situation status with respect to infrastructure (sewer).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Infrastructure Sewer: Comments",
			Value:     &f.InfrastructureSewerComments,
			PIFOTag:   "51.1.",
			PDFMap:    message.PDFName("Infra Sewer Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to infrastructure (sewer).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Search And Rescue",
			Value:    &f.SearchAndRescue,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "52.0.",
			PDFMap:   message.PDFNameMap{"SAR", "", "Off"},
			EditHelp: `This describes the current situation status with respect to search and rescue.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Search And Rescue: Comments",
			Value:     &f.SearchAndRescueComments,
			PIFOTag:   "52.1.",
			PDFMap:    message.PDFName("SAR Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to search and rescue.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Transportation Roads",
			Value:    &f.TransportationRoads,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "53.0.",
			PDFMap:   message.PDFNameMap{"Trans Roads", "", "Off"},
			EditHelp: `This describes the current situation status with respect to transportation (roads).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Transportation Roads: Comments",
			Value:     &f.TransportationRoadsComments,
			PIFOTag:   "53.1.",
			PDFMap:    message.PDFName("Trans Roads Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to transportation (roads).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Transportation Bridges",
			Value:    &f.TransportationBridges,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "54.0.",
			PDFMap:   message.PDFNameMap{"Trans Bridges", "", "Off"},
			EditHelp: `This describes the current situation status with respect to transportation (bridges).`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Transportation Bridges: Comments",
			Value:     &f.TransportationBridgesComments,
			PIFOTag:   "54.1.",
			PDFMap:    message.PDFName("Trans Bridges Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to transportation (bridges).`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Civil Unrest",
			Value:    &f.CivilUnrest,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "55.0.",
			PDFMap:   message.PDFNameMap{"Civil", "", "Off"},
			EditHelp: `This describes the current situation status with respect to civil unrest.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Civil Unrest: Comments",
			Value:     &f.CivilUnrestComments,
			PIFOTag:   "55.1.",
			PDFMap:    message.PDFName("Civil Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to civil unrest.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Animal Issues",
			Value:    &f.AnimalIssues,
			Choices:  message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag:  "56.0.",
			PDFMap:   message.PDFNameMap{"Animal", "", "Off"},
			EditHelp: `This describes the current situation status with respect to animal issues.`,
		}),
		message.NewMultilineField(&message.Field{
			Label:     "Animal Issues: Comments",
			Value:     &f.AnimalIssuesComments,
			PIFOTag:   "56.1.",
			PDFMap:    message.PDFName("Animal Comment"),
			EditWidth: 60,
			EditHelp:  `These are comments on the current situation status with respect to animal issues.`,
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &basePDFMaps)
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
	return message.DecodeForm(body, versions, create).(*JurisStat)
}
