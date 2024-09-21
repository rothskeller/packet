// Package jurisstat defines the Santa Clara County OA Jurisdiction Status Form
// message type.
package jurisstat

import (
	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type21 is the type definition for an OA jurisdiction status form.
var Type21 = message.Type{
	Tag:     "MuniStat",
	HTML:    "form-oa-muni-status.html",
	Version: "2.1",
	Name:    "OA municipal status form",
	Article: "an",
	FieldOrder: []string{
		"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "19.", "21.", "23.", "24.",
		"25.", "26.", "27.", "28.", "29.", "30.", "31.", "32.", "33.", "34.", "35.", "36.", "37.", "38.", "39.", "40.", "99.",
		"41.0.", "41.1.", "42.0.", "42.1.", "43.0.", "43.1.", "44.0.", "44.1.", "45.0.", "45.1.", "46.0.", "46.1.", "47.0.",
		"47.1.", "48.0.", "48.1.", "49.0.", "49.1.", "50.0.", "50.1.", "51.0.", "51.1.", "52.0.", "52.1.", "53.0.", "53.1.",
		"54.0.", "54.1.", "55.0.", "55.1.", "56.0.", "56.1.", "OpRelayRcvd", "OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
	},
}

func init() {
	message.Register(&Type21, decode21, nil)
}

// JurisStat21 holds an OA jurisdiction status form.
type JurisStat21 struct {
	message.BaseMessage
	baseform.BaseForm
	ReportType                    string
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

func make21() (f *JurisStat21) {
	const fieldCount = 79
	f = &JurisStat21{BaseMessage: message.BaseMessage{Type: &Type21}}
	f.BaseMessage.FSubject = &f.Jurisdiction
	f.BaseMessage.FBody = &f.CommunicationsComments
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
		message.NewRestrictedField(&message.Field{
			Label:    "Jurisdiction Name",
			Value:    &f.Jurisdiction,
			Choices:  message.Choices{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
			Presence: message.Required,
			PIFOTag:  "21.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:    "EOC Phone",
			Value:    &f.EOCPhone,
			Presence: f.requiredForComplete,
			PIFOTag:  "23.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "EOC Fax",
			Value:   &f.EOCFax,
			PIFOTag: "24.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Primary EM Contact Name",
			Value:    &f.PriEMContactName,
			Presence: f.requiredForComplete,
			PIFOTag:  "25.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:    "Primary EM Contact Phone",
			Value:    &f.PriEMContactPhone,
			Presence: f.requiredForComplete,
			PIFOTag:  "26.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Secondary EM Contact Name",
			Value:   &f.SecEMContactName,
			PIFOTag: "27.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "Secondary EM Contact Phone",
			Value:   &f.SecEMContactPhone,
			PIFOTag: "28.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Govt. Office Status",
			Value:    &f.OfficeStatus,
			Choices:  message.Choices{"Unknown", "Open", "Closed"},
			Presence: f.requiredForComplete,
			PIFOTag:  "29.",
		}),
		message.NewDateField(true, &message.Field{
			Label:   "Govt. Office Expected Open Date",
			Value:   &f.GovExpectedOpenDate,
			PIFOTag: "30.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:   "Govt. Office Expected Open Time",
			Value:   &f.GovExpectedOpenTime,
			PIFOTag: "31.",
		}),
		message.NewDateTimeField(&message.Field{
			Label: "Govt. Office Expected to Open",
		}, &f.GovExpectedOpenDate, &f.GovExpectedOpenTime),
		message.NewDateField(true, &message.Field{
			Label:   "Govt. Office Expected Close Date",
			Value:   &f.GovExpectedCloseDate,
			PIFOTag: "32.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:   "Govt. Office Expected Close Time",
			Value:   &f.GovExpectedCloseTime,
			PIFOTag: "33.",
		}),
		message.NewDateTimeField(&message.Field{
			Label: "Govt. Office Expected to Close",
		}, &f.GovExpectedCloseDate, &f.GovExpectedCloseTime),
		message.NewRestrictedField(&message.Field{
			Label:    "EOC Open",
			Value:    &f.EOCOpen,
			Choices:  message.Choices{"Unknown", "Yes", "No"},
			Presence: f.requiredForComplete,
			PIFOTag:  "34.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "EOC Activation Level",
			Value:    &f.EOCActivationLevel,
			Choices:  message.Choices{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
			Presence: f.requiredForComplete,
			PIFOTag:  "35.",
		}),
		message.NewDateField(true, &message.Field{
			Label:   "EOC Expected to Open Date",
			Value:   &f.EOCExpectedOpenDate,
			PIFOTag: "36.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:   "EOC Expected to Open Time",
			Value:   &f.EOCExpectedOpenTime,
			PIFOTag: "37.",
		}),
		message.NewDateTimeField(&message.Field{
			Label: "EOC Expected to Open",
		}, &f.EOCExpectedOpenDate, &f.EOCExpectedOpenTime),
		message.NewDateField(true, &message.Field{
			Label:   "EOC Expected to Close Date",
			Value:   &f.EOCExpectedCloseDate,
			PIFOTag: "38.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:   "EOC Expected to Close Time",
			Value:   &f.EOCExpectedCloseTime,
			PIFOTag: "39.",
		}),
		message.NewDateTimeField(&message.Field{
			Label: "EOC Expected to Close",
		}, &f.EOCExpectedCloseDate, &f.EOCExpectedCloseTime),
		message.NewRestrictedField(&message.Field{
			Label:    "State of Emergency",
			Value:    &f.StateOfEmergency,
			Choices:  message.Choices{"Unknown", "Yes", "No"},
			Presence: f.requiredForComplete,
			PIFOTag:  "40.",
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
			PIFOTag: "99.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Communications",
			Value:   &f.Communications,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "41.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Communications: Comments",
			Value:   &f.CommunicationsComments,
			PIFOTag: "41.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Debris",
			Value:   &f.Debris,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "42.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Debris: Comments",
			Value:   &f.DebrisComments,
			PIFOTag: "42.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Flooding",
			Value:   &f.Flooding,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "43.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Flooding: Comments",
			Value:   &f.FloodingComments,
			PIFOTag: "43.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Hazmat",
			Value:   &f.Hazmat,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "44.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Hazmat: Comments",
			Value:   &f.HazmatComments,
			PIFOTag: "44.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Emergency Services",
			Value:   &f.EmergencyServices,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "45.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Emergency Services: Comments",
			Value:   &f.EmergencyServicesComments,
			PIFOTag: "45.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Casualties",
			Value:   &f.Casualties,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "46.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Casualties: Comments",
			Value:   &f.CasualtiesComments,
			PIFOTag: "46.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Utilities Gas",
			Value:   &f.UtilitiesGas,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "47.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Utilities Gas: Comments",
			Value:   &f.UtilitiesGasComments,
			PIFOTag: "47.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Utilities Electric",
			Value:   &f.UtilitiesElectric,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "48.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Utilities Electric: Comments",
			Value:   &f.UtilitiesElectricComments,
			PIFOTag: "48.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure Power",
			Value:   &f.InfrastructurePower,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "49.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Infrastructure Power: Comments",
			Value:   &f.InfrastructurePowerComments,
			PIFOTag: "49.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure Water",
			Value:   &f.InfrastructureWater,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "50.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Infrastructure Water: Comments",
			Value:   &f.InfrastructureWaterComments,
			PIFOTag: "50.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Infrastructure Sewer",
			Value:   &f.InfrastructureSewer,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "51.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Infrastructure Sewer: Comments",
			Value:   &f.InfrastructureSewerComments,
			PIFOTag: "51.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Search And Rescue",
			Value:   &f.SearchAndRescue,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "52.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Search And Rescue: Comments",
			Value:   &f.SearchAndRescueComments,
			PIFOTag: "52.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Transportation Roads",
			Value:   &f.TransportationRoads,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "53.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Transportation Roads: Comments",
			Value:   &f.TransportationRoadsComments,
			PIFOTag: "53.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Transportation Bridges",
			Value:   &f.TransportationBridges,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "54.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Transportation Bridges: Comments",
			Value:   &f.TransportationBridgesComments,
			PIFOTag: "54.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Civil Unrest",
			Value:   &f.CivilUnrest,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "55.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Civil Unrest: Comments",
			Value:   &f.CivilUnrestComments,
			PIFOTag: "55.1.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:   "Animal Issues",
			Value:   &f.AnimalIssues,
			Choices: message.Choices{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
			PIFOTag: "56.0.",
		}),
		message.NewMultilineField(&message.Field{
			Label:   "Animal Issues: Comments",
			Value:   &f.AnimalIssuesComments,
			PIFOTag: "56.1.",
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, nil)
	if len(f.Fields) > fieldCount {
		panic("update JurisStat21 fieldCount")
	}
	return f
}

func (f *JurisStat21) requiredForComplete() (message.Presence, string) {
	if f.ReportType == "Complete" {
		return message.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return message.PresenceOptional, ""
}

func decode21(_, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type21.HTML || form.FormVersion != Type21.Version {
		return nil
	}
	var df = make21()
	message.DecodeForm(form, df)
	return df
}

func (f *JurisStat21) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo22().Compare(actual)
}

func (f *JurisStat21) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo22().RenderPDF(env, filename)
}

func (f *JurisStat21) convertTo22() (c *JurisStat22) {
	c = make22()
	c.CopyHeaderFields(&f.BaseForm)
	c.ReportType = f.ReportType
	c.JurisdictionCode = f.Jurisdiction
	c.Jurisdiction = f.Jurisdiction
	c.EOCPhone = f.EOCPhone
	c.EOCFax = f.EOCFax
	c.PriEMContactName = f.PriEMContactName
	c.PriEMContactPhone = f.PriEMContactPhone
	c.SecEMContactName = f.SecEMContactName
	c.SecEMContactPhone = f.SecEMContactPhone
	c.OfficeStatus = f.OfficeStatus
	c.GovExpectedOpenDate = f.GovExpectedOpenDate
	c.GovExpectedOpenTime = f.GovExpectedOpenTime
	c.GovExpectedCloseDate = f.GovExpectedCloseDate
	c.GovExpectedCloseTime = f.GovExpectedCloseTime
	c.EOCOpen = f.EOCOpen
	c.EOCActivationLevel = f.EOCActivationLevel
	c.EOCExpectedOpenDate = f.EOCExpectedOpenDate
	c.EOCExpectedOpenTime = f.EOCExpectedOpenTime
	c.EOCExpectedCloseDate = f.EOCExpectedCloseDate
	c.EOCExpectedCloseTime = f.EOCExpectedCloseTime
	c.StateOfEmergency = f.StateOfEmergency
	c.HowSOESent = f.HowSOESent
	c.Communications = f.Communications
	c.CommunicationsComments = f.CommunicationsComments
	c.Debris = f.Debris
	c.DebrisComments = f.DebrisComments
	c.Flooding = f.Flooding
	c.FloodingComments = f.FloodingComments
	c.Hazmat = f.Hazmat
	c.HazmatComments = f.HazmatComments
	c.EmergencyServices = f.EmergencyServices
	c.EmergencyServicesComments = f.EmergencyServicesComments
	c.Casualties = f.Casualties
	c.CasualtiesComments = f.CasualtiesComments
	c.UtilitiesGas = f.UtilitiesGas
	c.UtilitiesGasComments = f.UtilitiesGasComments
	c.UtilitiesElectric = f.UtilitiesElectric
	c.UtilitiesElectricComments = f.UtilitiesElectricComments
	c.InfrastructurePower = f.InfrastructurePower
	c.InfrastructurePowerComments = f.InfrastructurePowerComments
	c.InfrastructureWater = f.InfrastructureWater
	c.InfrastructureWaterComments = f.InfrastructureWaterComments
	c.InfrastructureSewer = f.InfrastructureSewer
	c.InfrastructureSewerComments = f.InfrastructureSewerComments
	c.SearchAndRescue = f.SearchAndRescue
	c.SearchAndRescueComments = f.SearchAndRescueComments
	c.TransportationRoads = f.TransportationRoads
	c.TransportationRoadsComments = f.TransportationRoadsComments
	c.TransportationBridges = f.TransportationBridges
	c.TransportationBridgesComments = f.TransportationBridgesComments
	c.CivilUnrest = f.CivilUnrest
	c.CivilUnrestComments = f.CivilUnrestComments
	c.AnimalIssues = f.AnimalIssues
	c.AnimalIssuesComments = f.AnimalIssuesComments
	c.CopyFooterFields(&f.BaseForm)
	return c
}
