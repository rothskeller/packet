package pktmsg

// This file defines TxMuniStatForm and RxMuniStatForm.

import (
	"fmt"
	"time"
)

// A TxMuniStatForm is an outgoing PackItForms-encoded message containing an
// SCCo OA Municipal Status form.
type TxMuniStatForm struct {
	TxSCCoForm
	ReportType                      string
	Jurisdiction                    string
	EOCPhone                        string
	EOCFax                          string
	PrimaryEMContactName            string
	PrimaryEMContactPhone           string
	SecondaryEMContactName          string
	SecondaryEMContactPhone         string
	GovtOfficeStatus                string
	GovtOfficeExpectedOpenDateTime  time.Time
	GovtOfficeExpectedCloseDateTime time.Time
	EOCIsOpen                       string
	EOCActivationLevel              string
	EOCExpectedOpenDateTime         time.Time
	EOCExpectedCloseDateTime        time.Time
	StateOfEmergency                string
	StateOfEmergencySent            string
	CommunicationsStatus            string
	CommunicationsComments          string
	DebrisStatus                    string
	DebrisComments                  string
	FloodingStatus                  string
	FloodingComments                string
	HazmatStatus                    string
	HazmatComments                  string
	EmergencyServicesStatus         string
	EmergencyServicesComments       string
	CasualtiesStatus                string
	CasualtiesComments              string
	GasUtilitiesStatus              string
	GasUtilitiesComments            string
	ElectricUtilitiesStatus         string
	ElectricUtilitiesComments       string
	PowerInfraStatus                string
	PowerInfraComments              string
	WaterInfraStatus                string
	WaterInfraComments              string
	SewerInfraStatus                string
	SewerInfraComments              string
	SearchRescueStatus              string
	SearchRescueComments            string
	TransportRoadStatus             string
	TransportRoadComments           string
	TransportBridgeStatus           string
	TransportBridgeComments         string
	CivilUnrestStatus               string
	CivilUnrestComments             string
	AnimalIssueStatus               string
	AnimalIssueComments             string
}

var (
	validReportType         = map[string]bool{"Update": true, "Complete": true}
	validGovtOfficeStatus   = map[string]bool{"": true, "Unknown": true, "Open": true, "Closed": true}
	validEOCIsOpen          = map[string]bool{"": true, "Unknown": true, "Yes": true, "No": true}
	validEOCActivationLevel = map[string]bool{"": true, "Normal": true, "Duty Officer": true, "Monitor": true, "Partial": true, "Full": true}
	validStateOfEmergency   = map[string]bool{"": true, "Unknown": true, "Yes": true, "No": true}
	validSituationStatus    = map[string]bool{"": true, "Unknown": true, "Normal": true, "Problem": true, "Failure": true, "Delayed": true, "Closed": true, "Early Out": true}
)

// Encode returns the encoded subject line and body of the message.
func (ms *TxMuniStatForm) Encode() (subject, body string, err error) {
	if err = ms.checkHeaderFooterFields(); err != nil {
		return "", "", err
	}
	if ms.Subject != "" {
		return "", "", ErrDontSet
	}
	if ms.ReportType == "" || ms.Jurisdiction == "" {
		return "", "", ErrIncomplete
	}
	if !validReportType[ms.ReportType] ||
		!validGovtOfficeStatus[ms.GovtOfficeStatus] ||
		!validEOCIsOpen[ms.EOCIsOpen] ||
		!validEOCActivationLevel[ms.EOCActivationLevel] ||
		!validStateOfEmergency[ms.StateOfEmergency] ||
		!validSituationStatus[ms.CommunicationsStatus] ||
		!validSituationStatus[ms.DebrisStatus] ||
		!validSituationStatus[ms.FloodingStatus] ||
		!validSituationStatus[ms.HazmatStatus] ||
		!validSituationStatus[ms.EmergencyServicesStatus] ||
		!validSituationStatus[ms.CasualtiesStatus] ||
		!validSituationStatus[ms.GasUtilitiesStatus] ||
		!validSituationStatus[ms.ElectricUtilitiesStatus] ||
		!validSituationStatus[ms.PowerInfraStatus] ||
		!validSituationStatus[ms.WaterInfraStatus] ||
		!validSituationStatus[ms.SewerInfraStatus] ||
		!validSituationStatus[ms.SearchRescueStatus] ||
		!validSituationStatus[ms.TransportRoadStatus] ||
		!validSituationStatus[ms.TransportBridgeStatus] ||
		!validSituationStatus[ms.CivilUnrestStatus] ||
		!validSituationStatus[ms.AnimalIssueStatus] {
		return "", "", ErrInvalid
	}
	ms.FormName = "MuniStat"
	ms.FormHTML = "form-oa-muni-status.html"
	ms.FormVersion = "2.1"
	ms.Subject = ms.Jurisdiction
	ms.encodeHeaderFields()
	ms.SetField("19.", ms.ReportType)
	ms.SetField("21.", ms.Jurisdiction)
	ms.SetField("23.", ms.EOCPhone)
	ms.SetField("24.", ms.EOCFax)
	ms.SetField("25.", ms.PrimaryEMContactName)
	ms.SetField("26.", ms.PrimaryEMContactPhone)
	ms.SetField("27.", ms.SecondaryEMContactName)
	ms.SetField("28.", ms.SecondaryEMContactPhone)
	ms.SetField("29.", ms.GovtOfficeStatus)
	ms.SetField("30.", ms.GovtOfficeExpectedOpenDateTime.Format("01/02/2006"))
	ms.SetField("31.", ms.GovtOfficeExpectedOpenDateTime.Format("15:04"))
	ms.SetField("32.", ms.GovtOfficeExpectedCloseDateTime.Format("01/02/2006"))
	ms.SetField("33.", ms.GovtOfficeExpectedCloseDateTime.Format("15:04"))
	ms.SetField("34.", ms.EOCIsOpen)
	ms.SetField("35.", ms.EOCActivationLevel)
	ms.SetField("36.", ms.EOCExpectedOpenDateTime.Format("01/02/2006"))
	ms.SetField("37.", ms.EOCExpectedOpenDateTime.Format("15:04"))
	ms.SetField("38.", ms.EOCExpectedCloseDateTime.Format("01/02/2006"))
	ms.SetField("39.", ms.EOCExpectedCloseDateTime.Format("15:04"))
	ms.SetField("40.", ms.StateOfEmergency)
	ms.SetField("99.", ms.StateOfEmergencySent)
	ms.SetField("41.0.", ms.CommunicationsStatus)
	ms.SetField("41.1.", ms.CommunicationsComments)
	ms.SetField("42.0.", ms.DebrisStatus)
	ms.SetField("42.1.", ms.DebrisComments)
	ms.SetField("43.0.", ms.FloodingStatus)
	ms.SetField("43.1.", ms.FloodingComments)
	ms.SetField("44.0.", ms.HazmatStatus)
	ms.SetField("44.1.", ms.HazmatComments)
	ms.SetField("45.0.", ms.EmergencyServicesStatus)
	ms.SetField("45.1.", ms.EmergencyServicesComments)
	ms.SetField("46.0.", ms.CasualtiesStatus)
	ms.SetField("46.1.", ms.CasualtiesComments)
	ms.SetField("47.0.", ms.GasUtilitiesStatus)
	ms.SetField("47.1.", ms.GasUtilitiesComments)
	ms.SetField("48.0.", ms.ElectricUtilitiesStatus)
	ms.SetField("48.1.", ms.ElectricUtilitiesComments)
	ms.SetField("49.0.", ms.PowerInfraStatus)
	ms.SetField("49.1.", ms.PowerInfraComments)
	ms.SetField("50.0.", ms.WaterInfraStatus)
	ms.SetField("50.1.", ms.WaterInfraComments)
	ms.SetField("51.0.", ms.SewerInfraStatus)
	ms.SetField("51.1.", ms.SewerInfraComments)
	ms.SetField("52.0.", ms.SearchRescueStatus)
	ms.SetField("52.1.", ms.SearchRescueComments)
	ms.SetField("53.0.", ms.TransportRoadStatus)
	ms.SetField("53.1.", ms.TransportRoadComments)
	ms.SetField("54.0.", ms.TransportBridgeStatus)
	ms.SetField("54.1.", ms.TransportBridgeComments)
	ms.SetField("55.0.", ms.CivilUnrestStatus)
	ms.SetField("55.1.", ms.CivilUnrestComments)
	ms.SetField("56.0.", ms.AnimalIssueStatus)
	ms.SetField("56.1.", ms.AnimalIssueComments)
	ms.encodeFooterFields()
	return ms.TxSCCoForm.Encode()
}

//------------------------------------------------------------------------------

// An RxMuniStatForm is a received PackItForms-encoded message containing an
// SCCo OA Municipal Status form.
type RxMuniStatForm struct {
	RxSCCoForm
	ReportType                      string
	Jurisdiction                    string
	EOCPhone                        string
	EOCFax                          string
	PrimaryEMContactName            string
	PrimaryEMContactPhone           string
	SecondaryEMContactName          string
	SecondaryEMContactPhone         string
	GovtOfficeStatus                string
	GovtOfficeExpectedOpenDate      string
	GovtOfficeExpectedOpenTime      string
	GovtOfficeExpectedOpenDateTime  time.Time
	GovtOfficeExpectedCloseDate     string
	GovtOfficeExpectedCloseTime     string
	GovtOfficeExpectedCloseDateTime time.Time
	EOCIsOpen                       string
	EOCActivationLevel              string
	EOCExpectedOpenDate             string
	EOCExpectedOpenTime             string
	EOCExpectedOpenDateTime         time.Time
	EOCExpectedCloseDate            string
	EOCExpectedCloseTime            string
	EOCExpectedCloseDateTime        time.Time
	StateOfEmergency                string
	StateOfEmergencySent            string
	CommunicationsStatus            string
	CommunicationsComments          string
	DebrisStatus                    string
	DebrisComments                  string
	FloodingStatus                  string
	FloodingComments                string
	HazmatStatus                    string
	HazmatComments                  string
	EmergencyServicesStatus         string
	EmergencyServicesComments       string
	CasualtiesStatus                string
	CasualtiesComments              string
	GasUtilitiesStatus              string
	GasUtilitiesComments            string
	ElectricUtilitiesStatus         string
	ElectricUtilitiesComments       string
	PowerInfraStatus                string
	PowerInfraComments              string
	WaterInfraStatus                string
	WaterInfraComments              string
	SewerInfraStatus                string
	SewerInfraComments              string
	SearchRescueStatus              string
	SearchRescueComments            string
	TransportRoadStatus             string
	TransportRoadComments           string
	TransportBridgeStatus           string
	TransportBridgeComments         string
	CivilUnrestStatus               string
	CivilUnrestComments             string
	AnimalIssueStatus               string
	AnimalIssueComments             string
}

// parseRxMuniStatForm examines an RxForm to see if it contains an EOC-213RR
// form, and if so, wraps it in an RxMuniStatForm and returns it.  If it is not,
// it returns nil.
func parseRxMuniStatForm(f *RxForm) *RxMuniStatForm {
	var ms RxMuniStatForm

	if f.FormHTML != "form-oa-muni-status.html" {
		return nil
	}
	ms.RxSCCoForm.RxForm = *f
	ms.extractHeaderFields()
	ms.ReportType = ms.Fields["19."]
	ms.Jurisdiction = ms.Fields["21."]
	ms.EOCPhone = ms.Fields["23."]
	ms.EOCFax = ms.Fields["24."]
	ms.PrimaryEMContactName = ms.Fields["25."]
	ms.PrimaryEMContactPhone = ms.Fields["26."]
	ms.SecondaryEMContactName = ms.Fields["27."]
	ms.SecondaryEMContactPhone = ms.Fields["28."]
	ms.GovtOfficeStatus = ms.Fields["29."]
	ms.GovtOfficeExpectedOpenDate = ms.Fields["30."]
	ms.GovtOfficeExpectedOpenTime = ms.Fields["31."]
	ms.GovtOfficeExpectedOpenDateTime = dateTimeParse(ms.GovtOfficeExpectedOpenDate, ms.GovtOfficeExpectedOpenTime)
	ms.GovtOfficeExpectedCloseDate = ms.Fields["32."]
	ms.GovtOfficeExpectedCloseTime = ms.Fields["33."]
	ms.GovtOfficeExpectedCloseDateTime = dateTimeParse(ms.GovtOfficeExpectedCloseDate, ms.GovtOfficeExpectedCloseTime)
	ms.EOCIsOpen = ms.Fields["34."]
	ms.EOCActivationLevel = ms.Fields["35."]
	ms.EOCExpectedOpenDate = ms.Fields["36."]
	ms.EOCExpectedOpenTime = ms.Fields["37."]
	ms.EOCExpectedOpenDateTime = dateTimeParse(ms.EOCExpectedOpenDate, ms.EOCExpectedOpenTime)
	ms.EOCExpectedCloseDate = ms.Fields["38."]
	ms.EOCExpectedCloseTime = ms.Fields["39."]
	ms.EOCExpectedCloseDateTime = dateTimeParse(ms.EOCExpectedCloseDate, ms.EOCExpectedCloseTime)
	ms.StateOfEmergency = ms.Fields["40."]
	ms.StateOfEmergencySent = ms.Fields["99."]
	ms.CommunicationsStatus = ms.Fields["41.0."]
	ms.CommunicationsComments = ms.Fields["41.1."]
	ms.DebrisStatus = ms.Fields["42.0."]
	ms.DebrisComments = ms.Fields["42.1."]
	ms.FloodingStatus = ms.Fields["43.0."]
	ms.FloodingComments = ms.Fields["43.1."]
	ms.HazmatStatus = ms.Fields["44.0."]
	ms.HazmatComments = ms.Fields["44.1."]
	ms.EmergencyServicesStatus = ms.Fields["45.0."]
	ms.EmergencyServicesComments = ms.Fields["45.1."]
	ms.CasualtiesStatus = ms.Fields["46.0."]
	ms.CasualtiesComments = ms.Fields["46.1."]
	ms.GasUtilitiesStatus = ms.Fields["47.0."]
	ms.GasUtilitiesComments = ms.Fields["47.1."]
	ms.ElectricUtilitiesStatus = ms.Fields["48.0."]
	ms.ElectricUtilitiesComments = ms.Fields["48.1."]
	ms.PowerInfraStatus = ms.Fields["49.0."]
	ms.PowerInfraComments = ms.Fields["49.1."]
	ms.WaterInfraStatus = ms.Fields["50.0."]
	ms.WaterInfraComments = ms.Fields["50.1."]
	ms.SewerInfraStatus = ms.Fields["51.0."]
	ms.SewerInfraComments = ms.Fields["51.1."]
	ms.SearchRescueStatus = ms.Fields["52.0."]
	ms.SearchRescueComments = ms.Fields["52.1."]
	ms.TransportRoadStatus = ms.Fields["53.0."]
	ms.TransportRoadComments = ms.Fields["53.1."]
	ms.TransportBridgeStatus = ms.Fields["54.0."]
	ms.TransportBridgeComments = ms.Fields["54.1."]
	ms.CivilUnrestStatus = ms.Fields["55.0."]
	ms.CivilUnrestComments = ms.Fields["55.1."]
	ms.AnimalIssueStatus = ms.Fields["56.0."]
	ms.AnimalIssueComments = ms.Fields["56.1."]
	ms.extractFooterFields()
	return &ms
}

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (ms *RxMuniStatForm) Valid() bool {
	return ms.RxSCCoForm.Valid() &&
		validReportType[ms.ReportType] &&
		ms.Jurisdiction != "" &&
		validGovtOfficeStatus[ms.GovtOfficeStatus] &&
		validEOCIsOpen[ms.EOCIsOpen] &&
		validEOCActivationLevel[ms.EOCActivationLevel] &&
		validStateOfEmergency[ms.StateOfEmergency] &&
		validSituationStatus[ms.CommunicationsStatus] &&
		validSituationStatus[ms.DebrisStatus] &&
		validSituationStatus[ms.FloodingStatus] &&
		validSituationStatus[ms.HazmatStatus] &&
		validSituationStatus[ms.EmergencyServicesStatus] &&
		validSituationStatus[ms.CasualtiesStatus] &&
		validSituationStatus[ms.GasUtilitiesStatus] &&
		validSituationStatus[ms.ElectricUtilitiesStatus] &&
		validSituationStatus[ms.PowerInfraStatus] &&
		validSituationStatus[ms.WaterInfraStatus] &&
		validSituationStatus[ms.SewerInfraStatus] &&
		validSituationStatus[ms.SearchRescueStatus] &&
		validSituationStatus[ms.TransportRoadStatus] &&
		validSituationStatus[ms.TransportBridgeStatus] &&
		validSituationStatus[ms.CivilUnrestStatus] &&
		validSituationStatus[ms.AnimalIssueStatus]
}

// EncodeSubjectLine returns what the subject line should be based on the
// received form contents.
func (ms *RxMuniStatForm) EncodeSubjectLine() string {
	return fmt.Sprintf("%s_%s_MuniStat_%s", ms.OriginMessageNumber, ms.HandlingOrder.Code(), ms.Jurisdiction)
}

// TypeCode returns the machine-readable code for the message type.
func (*RxMuniStatForm) TypeCode() string { return "MuniStat" }

// TypeName returns the human-reading name of the message type.
func (*RxMuniStatForm) TypeName() string { return "OA Municipal Status form" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxMuniStatForm) TypeArticle() string { return "an" }
