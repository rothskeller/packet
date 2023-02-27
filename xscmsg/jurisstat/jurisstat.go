package jurisstat

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies jurisdiction status forms.
const (
	Tag21 = "MuniStat"
	Tag22 = "JurisStat"
	Tag   = Tag22
)

func init() {
	xscmsg.RegisterCreate(formtype22, create22)
	xscmsg.RegisterCreate(formtype21, create21)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef.DefaultValue = "IMMEDIATE"
	toICSPositionDef.DefaultValue = "Situation Analysis Unit"
	toICSPositionDef.Choices = []string{"Planning Section", "Situation Analysis Unit"}
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Choices = []string{"County EOC"}
}

func create22() *xscmsg.Message {
	return xscform.CreateForm(formtype22, fieldDefsV22)
}

func create21() *xscmsg.Message {
	return xscform.CreateForm(formtype21, fieldDefsV21)
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
		Key:        xscmsg.FComplete,
		Label:      "Report Type",
		Choices:    []string{"Update", "Complete"},
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Flags:      xscmsg.Required,
	}
	jurisdictionCodeDef = &xscmsg.FieldDef{
		Tag:        "21.",
		Flags:      xscmsg.Readonly,
		Validators: []xscmsg.Validator{setJurisdictionCode, xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
	}
	jurisdictionDefV21 = &xscmsg.FieldDef{
		Tag:        "21.",
		Label:      "Jurisdiction Name",
		Key:        xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
		Flags:      xscmsg.Required,
	}
	jurisdictionDefV22 = &xscmsg.FieldDef{
		Tag:   "22.",
		Label: "Jurisdiction Name",
		Key:   xscmsg.FSubject,
		Flags: xscmsg.Required,
	}
	eocPhoneDef = &xscmsg.FieldDef{
		Tag:        "23.",
		Label:      "EOC Phone",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	eocFaxDef = &xscmsg.FieldDef{
		Tag:   "24.",
		Label: "EOC Fax",
	}
	priEmContactNameDef = &xscmsg.FieldDef{
		Tag:   "25.",
		Label: "Primary EM Contact Name",
		Flags: xscmsg.RequiredForComplete,
	}
	priEmContactPhoneDef = &xscmsg.FieldDef{
		Tag:        "26.",
		Label:      "Primary EM Contact Phone",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	secEmContactNameDef = &xscmsg.FieldDef{
		Tag:   "27.",
		Label: "Secondary EM Contact Name",
	}
	secEmContactPhoneDef = &xscmsg.FieldDef{
		Tag:        "28.",
		Label:      "Secondary EM Contact Phone",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	officeStatusDef = &xscmsg.FieldDef{
		Tag:        "29.",
		Label:      "Office Status",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Open", "Closed"},
		Flags:      xscmsg.RequiredForComplete,
	}
	govExpectedOpenDateDef = &xscmsg.FieldDef{
		Tag:        "30.",
		Label:      "Office Expected to Open Date",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	govExpectedOpenTimeDef = &xscmsg.FieldDef{
		Tag:        "31.",
		Label:      "Office Expected to Open Time",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	govExpectedCloseDateDef = &xscmsg.FieldDef{
		Tag:        "32.",
		Label:      "Office Expected to Close Date",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	govExpectedCloseTimeDef = &xscmsg.FieldDef{
		Tag:        "33.",
		Label:      "Office Expected to Close Time",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	eocOpenDef = &xscmsg.FieldDef{
		Tag:        "34.",
		Label:      "EOC Open",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Yes", "No"},
		Flags:      xscmsg.RequiredForComplete,
	}
	eocActivationLevelDef = &xscmsg.FieldDef{
		Tag:        "35.",
		Label:      "Activation Level",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Normal", "Duty Officer", "Monitor", "Partial", "Full"},
		Flags:      xscmsg.RequiredForComplete,
	}
	eocExpectedOpenDateDef = &xscmsg.FieldDef{
		Tag:        "36.",
		Label:      "EOC Expected to Open Date",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	eocExpectedOpenTimeDef = &xscmsg.FieldDef{
		Tag:        "37.",
		Label:      "EOC Expected to Open Time",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	eocExpectedCloseDateDef = &xscmsg.FieldDef{
		Tag:        "38.",
		Label:      "EOC Expected to Close Date",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
	}
	eocExpectedCloseTimeDef = &xscmsg.FieldDef{
		Tag:        "39.",
		Label:      "EOC Expected to Close Time",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
	}
	stateOfEmergencyDef = &xscmsg.FieldDef{
		Tag:        "40.",
		Label:      "State of Emergency",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Yes", "No"},
		Flags:      xscmsg.RequiredForComplete,
	}
	howToSendDef = &xscmsg.FieldDef{
		Tag:   "98.",
		Label: "SOE Proclamation: How to Send",
		Flags: xscmsg.Readonly,
	}
	howSentDef = &xscmsg.FieldDef{
		Tag:   "99.",
		Label: "SOE Proclamation: Indicate how sent (method/to)",
	}
	communicationsDef = &xscmsg.FieldDef{
		Tag:        "41.0.",
		Label:      "Communications",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	communicationsCommentsDef = &xscmsg.FieldDef{
		Tag:   "41.1.",
		Label: "Communications Comments",
		Key:   xscmsg.FBody,
		Flags: xscmsg.Multiline,
	}
	debrisDef = &xscmsg.FieldDef{
		Tag:        "42.0.",
		Label:      "Debris",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	debrisCommentsDef = &xscmsg.FieldDef{
		Tag:   "42.1.",
		Label: "Debris Comments",
		Flags: xscmsg.Multiline,
	}
	floodingDef = &xscmsg.FieldDef{
		Tag:        "43.0.",
		Label:      "Flooding",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	floodingCommentsDef = &xscmsg.FieldDef{
		Tag:   "43.1.",
		Label: "Flooding Comments",
		Flags: xscmsg.Multiline,
	}
	hazmatDef = &xscmsg.FieldDef{
		Tag:        "44.0.",
		Label:      "Hazmat",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	hazmatCommentsDef = &xscmsg.FieldDef{
		Tag:   "44.1.",
		Label: "Hazmat Comments",
		Flags: xscmsg.Multiline,
	}
	emergencyServicesDef = &xscmsg.FieldDef{
		Tag:        "45.0.",
		Label:      "Emergency Services",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	emergencyServicesCommentsDef = &xscmsg.FieldDef{
		Tag:   "45.1.",
		Label: "Emergency Services Comments",
		Flags: xscmsg.Multiline,
	}
	casualtiesDef = &xscmsg.FieldDef{
		Tag:        "46.0.",
		Label:      "Casualties",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	casualtiesCommentsDef = &xscmsg.FieldDef{
		Tag:   "46.1.",
		Label: "Casualties Comments",
		Flags: xscmsg.Multiline,
	}
	utilitiesGasDef = &xscmsg.FieldDef{
		Tag:        "47.0.",
		Label:      "Utilities (Gas)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	utilitiesGasCommentsDef = &xscmsg.FieldDef{
		Tag:   "47.1.",
		Label: "Utilities (Gas) Comments",
		Flags: xscmsg.Multiline,
	}
	utilitiesElectricDef = &xscmsg.FieldDef{
		Tag:        "48.0.",
		Label:      "Utilities (Electric)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	utilitiesElectricCommentsDef = &xscmsg.FieldDef{
		Tag:   "48.1.",
		Label: "Utilities (Electric) Comments",
		Flags: xscmsg.Multiline,
	}
	infrastructurePowerDef = &xscmsg.FieldDef{
		Tag:        "49.0.",
		Label:      "Infrastructure (Power)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	infrastructurePowerCommentsDef = &xscmsg.FieldDef{
		Tag:   "49.1.",
		Label: "Infrastructure (Power) Comments",
		Flags: xscmsg.Multiline,
	}
	infrastructureWaterSystemsDef = &xscmsg.FieldDef{
		Tag:        "50.0.",
		Label:      "Infrastructure (Water Systems)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	infrastructureWaterSystemsCommentsDef = &xscmsg.FieldDef{
		Tag:   "50.1.",
		Label: "Infrastructure (Water Systems) Comments",
		Flags: xscmsg.Multiline,
	}
	infrastructureSewerSystemsDef = &xscmsg.FieldDef{
		Tag:        "51.0.",
		Label:      "Infrastructure (Sewer Systems)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	infrastructureSewerSystemsCommentsDef = &xscmsg.FieldDef{
		Tag:   "51.1.",
		Label: "Infrastructure (Sewer Systems) Comments",
		Flags: xscmsg.Multiline,
	}
	searchAndRescueDef = &xscmsg.FieldDef{
		Tag:        "52.0.",
		Label:      "Search and Rescue",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	searchAndRescueCommentsDef = &xscmsg.FieldDef{
		Tag:   "52.1.",
		Label: "Search and Rescue Comments",
		Flags: xscmsg.Multiline,
	}
	transportationRoadsDef = &xscmsg.FieldDef{
		Tag:        "53.0.",
		Label:      "Transportation (Roads)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	transportationRoadsCommentsDef = &xscmsg.FieldDef{
		Tag:   "53.1.",
		Label: "Transportation (Roads) Comments",
		Flags: xscmsg.Multiline,
	}
	transportationBridgesDef = &xscmsg.FieldDef{
		Tag:        "54.0.",
		Label:      "Transportation (Bridges)",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	transportationBridgesCommentsDef = &xscmsg.FieldDef{
		Tag:   "54.1.",
		Label: "Transportation (Bridges) Comments",
		Flags: xscmsg.Multiline,
	}
	civilUnrestDef = &xscmsg.FieldDef{
		Tag:        "55.0.",
		Label:      "Civil Unrest",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	civilUnrestCommentsDef = &xscmsg.FieldDef{
		Tag:   "55.1.",
		Label: "Civil Unrest Comments",
		Flags: xscmsg.Multiline,
	}
	animalIssuesDef = &xscmsg.FieldDef{
		Tag:        "56.0.",
		Label:      "Animal Issues",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"},
		Flags:      xscmsg.RequiredForComplete,
	}
	animalIssuesCommentsDef = &xscmsg.FieldDef{
		Tag:   "56.1.",
		Label: "Animal Issues Comments",
		Flags: xscmsg.Multiline,
	}
)

func setJurisdictionCode(f *xscmsg.Field, m *xscmsg.Message, strict bool) string {
	if !strict {
		xscform.SetChoiceFieldFromValue(f, m.Field("22.").Value)
	}
	return ""
}
