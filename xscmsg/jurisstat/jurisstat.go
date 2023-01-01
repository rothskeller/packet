package jurisstat

import (
	"fmt"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies jurisdiction status forms.
const (
	Tag21 = "MuniStat"
	Tag22 = "JurisStat"
	Tag   = Tag22
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef.DefaultValue = "IMMEDIATE"
	toICSPositionDef.DefaultValue = "Situation Analysis Unit"
	toICSPositionDef.Comment = "required: Situation Analysis Unit, Planning Section, ..."
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Comment = "required: County EOC, ..."
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype22, fieldDefsV22)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil || form.FormType != formtype22.HTML {
		return nil
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.2") {
		return xscform.AdoptForm(formtype22, fieldDefsV22, msg, form)
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.1") {
		return xscform.AdoptForm(formtype21, fieldDefsV21, msg, form)
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return xscform.AdoptForm(formtype21, fieldDefsV20, msg, form)
	}
	return nil
}

var formtype22 = &xscmsg.MessageType{
	Tag:         Tag22,
	Name:        "OA jurisdiction status form",
	Article:     "an",
	HTML:        "form-oa-muni-status.html",
	Version:     "2.2",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

var formtype21 = &xscmsg.MessageType{
	Tag:         Tag21,
	Name:        "OA municipal status form",
	Article:     "an",
	HTML:        "form-oa-muni-status.html",
	Version:     "2.1",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

// fieldDefsV22 adds the jurisdictionCodeDef and chainges jurisdictionDefV21 to
// jurisdictionDefV22.
var fieldDefsV22 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// Jurisdiction Status fields
	reportTypeDef, jurisdictionCodeDef, jurisdictionDefV22, eocPhoneDef, eocFaxDef, priEmContactNameDef, priEmContactPhoneDef,
	secEmContactNameDef, secEmContactPhoneDef, officeStatusDef, govExpectedOpenDateDef, govExpectedOpenTimeDef,
	govExpectedCloseDateDef, govExpectedCloseTimeDef, eocOpenDef, eocActivationLevelDef, eocExpectedOpenDateDef,
	eocExpectedOpenTimeDef, eocExpectedCloseDateDef, eocExpectedCloseTimeDef, stateOfEmergencyDef, /* howToSendDef */
	howSentDef, communicationsDef, communicationsCommentsDef, debrisDef, debrisCommentsDef, floodingDef, floodingCommentsDef,
	hazmatDef, hazmatCommentsDef, emergencyServicesDef, emergencyServicesCommentsDef, casualtiesDef, casualtiesCommentsDef,
	utilitiesGasDef, utilitiesGasCommentsDef, utilitiesElectricDef, utilitiesElectricCommentsDef, infrastructurePowerDef,
	infrastructurePowerCommentsDef, infrastructureWaterSystemsDef, infrastructureWaterSystemsCommentsDef,
	infrastructureSewerSystemsDef, infrastructureSewerSystemsCommentsDef, searchAndRescueDef, searchAndRescueCommentsDef,
	transportationRoadsDef, transportationRoadsCommentsDef, transportationBridgesDef, transportationBridgesCommentsDef,
	civilUnrestDef, civilUnrestCommentsDef, animalIssuesDef, animalIssuesCommentsDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

// fieldDefsV21 removes howToSendDef.
var fieldDefsV21 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// Jurisdiction Status fields
	reportTypeDef, jurisdictionDefV21, eocPhoneDef, eocFaxDef, priEmContactNameDef, priEmContactPhoneDef,
	secEmContactNameDef, secEmContactPhoneDef, officeStatusDef, govExpectedOpenDateDef, govExpectedOpenTimeDef,
	govExpectedCloseDateDef, govExpectedCloseTimeDef, eocOpenDef, eocActivationLevelDef, eocExpectedOpenDateDef,
	eocExpectedOpenTimeDef, eocExpectedCloseDateDef, eocExpectedCloseTimeDef, stateOfEmergencyDef, /* howToSendDef */
	howSentDef, communicationsDef, communicationsCommentsDef, debrisDef, debrisCommentsDef, floodingDef, floodingCommentsDef,
	hazmatDef, hazmatCommentsDef, emergencyServicesDef, emergencyServicesCommentsDef, casualtiesDef, casualtiesCommentsDef,
	utilitiesGasDef, utilitiesGasCommentsDef, utilitiesElectricDef, utilitiesElectricCommentsDef, infrastructurePowerDef,
	infrastructurePowerCommentsDef, infrastructureWaterSystemsDef, infrastructureWaterSystemsCommentsDef,
	infrastructureSewerSystemsDef, infrastructureSewerSystemsCommentsDef, searchAndRescueDef, searchAndRescueCommentsDef,
	transportationRoadsDef, transportationRoadsCommentsDef, transportationBridgesDef, transportationBridgesCommentsDef,
	civilUnrestDef, civilUnrestCommentsDef, animalIssuesDef, animalIssuesCommentsDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var fieldDefsV20 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// Jurisdiction Status fields
	reportTypeDef, jurisdictionDefV21, eocPhoneDef, eocFaxDef, priEmContactNameDef, priEmContactPhoneDef,
	secEmContactNameDef, secEmContactPhoneDef, officeStatusDef, govExpectedOpenDateDef, govExpectedOpenTimeDef,
	govExpectedCloseDateDef, govExpectedCloseTimeDef, eocOpenDef, eocActivationLevelDef, eocExpectedOpenDateDef,
	eocExpectedOpenTimeDef, eocExpectedCloseDateDef, eocExpectedCloseTimeDef, stateOfEmergencyDef, howToSendDef,
	howSentDef, communicationsDef, communicationsCommentsDef, debrisDef, debrisCommentsDef, floodingDef, floodingCommentsDef,
	hazmatDef, hazmatCommentsDef, emergencyServicesDef, emergencyServicesCommentsDef, casualtiesDef, casualtiesCommentsDef,
	utilitiesGasDef, utilitiesGasCommentsDef, utilitiesElectricDef, utilitiesElectricCommentsDef, infrastructurePowerDef,
	infrastructurePowerCommentsDef, infrastructureWaterSystemsDef, infrastructureWaterSystemsCommentsDef,
	infrastructureSewerSystemsDef, infrastructureSewerSystemsCommentsDef, searchAndRescueDef, searchAndRescueCommentsDef,
	transportationRoadsDef, transportationRoadsCommentsDef, transportationBridgesDef, transportationBridgesCommentsDef,
	civilUnrestDef, civilUnrestCommentsDef, animalIssuesDef, animalIssuesCommentsDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	handlingDef      = *xscform.HandlingDef      // modified in func init
	toICSPositionDef = *xscform.ToICSPositionDef // modified in func init
	toLocationDef    = *xscform.ToLocationDef    // modified in func init
	reportTypeDef    = &xscmsg.FieldDef{
		Tag:        "19.",
		Annotation: "report-type",
		Label:      "Report Type",
		Comment:    "required: Update, Complete",
		Choices:    []string{"Update", "Complete"},
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
	}
	jurisdictionCodeDef = &xscmsg.FieldDef{
		Tag:        "21.",
		Annotation: "jurisdiction-code",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setJurisdictionCode, xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
	}
	jurisdictionDefV21 = &xscmsg.FieldDef{
		Tag:        "21.",
		Annotation: "jurisdiction",
		Label:      "Jurisdiction Name",
		Comment:    "required: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, Unincorporated",
		Canonical:  xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
	}
	jurisdictionDefV22 = &xscmsg.FieldDef{
		Tag:        "22.",
		Annotation: "jurisdiction",
		Label:      "Jurisdiction Name",
		Comment:    "required: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, ...",
		Canonical:  xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	eocPhoneDef = &xscmsg.FieldDef{
		Tag:        "23.",
		Annotation: "eoc-phone",
		Label:      "EOC Phone",
		Comment:    "phone-number required-for-complete",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidatePhoneNumber},
	}
	eocFaxDef = &xscmsg.FieldDef{
		Tag:        "24.",
		Annotation: "eoc-fax",
		Label:      "EOC Fax",
		Comment:    "phone-number",
	}
	priEmContactNameDef = &xscmsg.FieldDef{
		Tag:        "25.",
		Annotation: "pri-em-contact-name",
		Label:      "Primary EM Contact Name",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{validateRequiredForComplete},
	}
	priEmContactPhoneDef = &xscmsg.FieldDef{
		Tag:        "26.",
		Annotation: "pri-em-contact-phone",
		Label:      "Primary EM Contact Phone",
		Comment:    "phone-number required-for-complete",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidatePhoneNumber},
	}
	secEmContactNameDef = &xscmsg.FieldDef{
		Tag:        "27.",
		Annotation: "sec-em-contact-name",
		Label:      "Secondary EM Contact Name",
	}
	secEmContactPhoneDef = &xscmsg.FieldDef{
		Tag:        "28.",
		Annotation: "sec-em-contact-phone",
		Label:      "Secondary EM Contact Phone",
		Comment:    "phone-number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	officeStatusDef = &xscmsg.FieldDef{
		Tag:        "29.",
		Annotation: "office-status",
		Label:      "Office Status",
		Comment:    "required-for-complete: Unknown, Open, Closed",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Open", "Closed"},
	}
	govExpectedOpenDateDef = &xscmsg.FieldDef{
		Tag:        "30.",
		Annotation: "gov-expected-open-date",
		Label:      "Office Expected to Open Date",
		Comment:    "date",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	govExpectedOpenTimeDef = &xscmsg.FieldDef{
		Tag:        "31.",
		Annotation: "gov-expected-open-time",
		Label:      "Office Expected to Open Time",
		Comment:    "time",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	govExpectedCloseDateDef = &xscmsg.FieldDef{
		Tag:        "32.",
		Annotation: "gov-expected-close-date",
		Label:      "Office Expected to Close Date",
		Comment:    "date",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	govExpectedCloseTimeDef = &xscmsg.FieldDef{
		Tag:        "33.",
		Annotation: "gov-expected-close-time",
		Label:      "Office Expected to Close Time",
		Comment:    "time",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	eocOpenDef = &xscmsg.FieldDef{
		Tag:        "34.",
		Annotation: "eoc-open",
		Label:      "EOC Open",
		Comment:    "required-for-complete: Unknown, Yes, No",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Yes", "No"},
	}
	eocActivationLevelDef = &xscmsg.FieldDef{
		Tag:        "35.",
		Annotation: "eoc-activation-level",
		Label:      "Activation Level",
		Comment:    "required-for-complete: Normal, Duty Officer, Monitor, Partial, Full",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
	}
	eocExpectedOpenDateDef = &xscmsg.FieldDef{
		Tag:        "36.",
		Annotation: "eoc-expected-open-date",
		Label:      "EOC Expected to Open Date",
		Comment:    "date",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	eocExpectedOpenTimeDef = &xscmsg.FieldDef{
		Tag:        "37.",
		Annotation: "eoc-expected-open-time",
		Label:      "EOC Expected to Open Time",
		Comment:    "time",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	eocExpectedCloseDateDef = &xscmsg.FieldDef{
		Tag:        "38.",
		Annotation: "eoc-expected-close-date",
		Label:      "EOC Expected to Close Date",
		Comment:    "date",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	eocExpectedCloseTimeDef = &xscmsg.FieldDef{
		Tag:        "39.",
		Annotation: "eoc-expected-close-time",
		Label:      "EOC Expected to Close Time",
		Comment:    "time",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	stateOfEmergencyDef = &xscmsg.FieldDef{
		Tag:        "40.",
		Annotation: "state-of-emergency",
		Label:      "State of Emergency",
		Comment:    "required-for-complete: Unknown, Yes, No",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Yes", "No"},
	}
	howToSendDef = &xscmsg.FieldDef{
		Tag:        "98.",
		Annotation: "how-to-send",
		Label:      "SOE Proclamation: How to Send",
		ReadOnly:   true,
	}
	howSentDef = &xscmsg.FieldDef{
		Tag:        "99.",
		Annotation: "how-sent",
		Label:      "SOE Proclamation: Indicate how sent (method/to)",
	}
	communicationsDef = &xscmsg.FieldDef{
		Tag:        "41.0.",
		Annotation: "communications",
		Label:      "Communications",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	communicationsCommentsDef = &xscmsg.FieldDef{
		Tag:        "41.1.",
		Annotation: "communications-comments",
		Label:      "Communications Comments",
	}
	debrisDef = &xscmsg.FieldDef{
		Tag:        "42.0.",
		Annotation: "debris",
		Label:      "Debris",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	debrisCommentsDef = &xscmsg.FieldDef{
		Tag:        "42.1.",
		Annotation: "debris-comments",
		Label:      "Debris Comments",
	}
	floodingDef = &xscmsg.FieldDef{
		Tag:        "43.0.",
		Annotation: "flooding",
		Label:      "Flooding",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	floodingCommentsDef = &xscmsg.FieldDef{
		Tag:        "43.1.",
		Annotation: "flooding-comments",
		Label:      "Flooding Comments",
	}
	hazmatDef = &xscmsg.FieldDef{
		Tag:        "44.0.",
		Annotation: "hazmat",
		Label:      "Hazmat",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	hazmatCommentsDef = &xscmsg.FieldDef{
		Tag:        "44.1.",
		Annotation: "hazmat-comments",
		Label:      "Hazmat Comments",
	}
	emergencyServicesDef = &xscmsg.FieldDef{
		Tag:        "45.0.",
		Annotation: "emergency-services",
		Label:      "Emergency Services",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	emergencyServicesCommentsDef = &xscmsg.FieldDef{
		Tag:        "45.1.",
		Annotation: "emergency-services-comments",
		Label:      "Emergency Services Comments",
	}
	casualtiesDef = &xscmsg.FieldDef{
		Tag:        "46.0.",
		Annotation: "casualties",
		Label:      "Casualties",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	casualtiesCommentsDef = &xscmsg.FieldDef{
		Tag:        "46.1.",
		Annotation: "casualties-comments",
		Label:      "Casualties Comments",
	}
	utilitiesGasDef = &xscmsg.FieldDef{
		Tag:        "47.0.",
		Annotation: "utilities-gas",
		Label:      "Utilities (Gas)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	utilitiesGasCommentsDef = &xscmsg.FieldDef{
		Tag:        "47.1.",
		Annotation: "utilities-gas-comments",
		Label:      "Utilities (Gas) Comments",
	}
	utilitiesElectricDef = &xscmsg.FieldDef{
		Tag:        "48.0.",
		Annotation: "utilities-electric",
		Label:      "Utilities (Electric)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	utilitiesElectricCommentsDef = &xscmsg.FieldDef{
		Tag:        "48.1.",
		Annotation: "utilities-electric-comments",
		Label:      "Utilities (Electric) Comments",
	}
	infrastructurePowerDef = &xscmsg.FieldDef{
		Tag:        "49.0.",
		Annotation: "infrastructure-power",
		Label:      "Infrastructure (Power)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	infrastructurePowerCommentsDef = &xscmsg.FieldDef{
		Tag:        "49.1.",
		Annotation: "infrastructure-power-comments",
		Label:      "Infrastructure (Power) Comments",
	}
	infrastructureWaterSystemsDef = &xscmsg.FieldDef{
		Tag:        "50.0.",
		Annotation: "infrastructure-water-systems",
		Label:      "Infrastructure (Water Systems)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	infrastructureWaterSystemsCommentsDef = &xscmsg.FieldDef{
		Tag:        "50.1.",
		Annotation: "infrastructure-water-systems-comments",
		Label:      "Infrastructure (Water Systems) Comments",
	}
	infrastructureSewerSystemsDef = &xscmsg.FieldDef{
		Tag:        "51.0.",
		Annotation: "infrastructure-sewer-systems",
		Label:      "Infrastructure (Sewer Systems)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	infrastructureSewerSystemsCommentsDef = &xscmsg.FieldDef{
		Tag:        "51.1.",
		Annotation: "infrastructure-sewer-systems-comments",
		Label:      "Infrastructure (Sewer Systems) Comments",
	}
	searchAndRescueDef = &xscmsg.FieldDef{
		Tag:        "52.0.",
		Annotation: "search-and-rescue",
		Label:      "Search and Rescue",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	searchAndRescueCommentsDef = &xscmsg.FieldDef{
		Tag:        "52.1.",
		Annotation: "search-and-rescue-comments",
		Label:      "Search and Rescue Comments",
	}
	transportationRoadsDef = &xscmsg.FieldDef{
		Tag:        "53.0.",
		Annotation: "transportation-roads",
		Label:      "Transportation (Roads)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	transportationRoadsCommentsDef = &xscmsg.FieldDef{
		Tag:        "53.1.",
		Annotation: "transportation-roads-comments",
		Label:      "Transportation (Roads) Comments",
	}
	transportationBridgesDef = &xscmsg.FieldDef{
		Tag:        "54.0.",
		Annotation: "transportation-bridges",
		Label:      "Transportation (Bridges)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	transportationBridgesCommentsDef = &xscmsg.FieldDef{
		Tag:        "54.1.",
		Annotation: "transportation-bridges-comments",
		Label:      "Transportation (Bridges) Comments",
	}
	civilUnrestDef = &xscmsg.FieldDef{
		Tag:        "55.0.",
		Annotation: "civil-unrest",
		Label:      "Civil Unrest",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	civilUnrestCommentsDef = &xscmsg.FieldDef{
		Tag:        "55.1.",
		Annotation: "civil-unrest-comments",
		Label:      "Civil Unrest Comments",
	}
	animalIssuesDef = &xscmsg.FieldDef{
		Tag:        "56.0.",
		Annotation: "animal-issues",
		Label:      "Animal Issues",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
		Validators: []xscmsg.Validator{validateRequiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
	}
	animalIssuesCommentsDef = &xscmsg.FieldDef{
		Tag:        "56.1.",
		Annotation: "animal-issues-comments",
		Label:      "Animal Issues Comments",
	}
)

func validateRequiredForComplete(f *xscmsg.Field, m *xscmsg.Message, _ bool) string {
	if m.Field("19.").Value != "Complete" {
		return ""
	}
	if f.Value == "" {
		return fmt.Sprintf("The %q field must have a value when the \"19.\" field is set to \"Complete\".", f.Def.Tag)
	}
	return ""
}

func setJurisdictionCode(f *xscmsg.Field, m *xscmsg.Message, strict bool) string {
	if !strict {
		xscform.SetChoiceFieldFromValue(f, m.Field("22.").Value)
	}
	return ""
}
