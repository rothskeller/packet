package jurisstat

import (
	"fmt"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies jurisdiction status forms.
const (
	Tag20 = "MuniStat"
	Tag22 = "JurisStat"
	Tag   = Tag22
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)
}

func create() xscmsg.Message {
	return xscform.CreateForm(formtype22, makeFields22())
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.Message {
	if form == nil || form.FormType != formtype22.HTML {
		return nil
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.2") {
		return xscform.AdoptForm(formtype22, makeFields22(), msg, form)
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.1") {
		return xscform.AdoptForm(formtype21, makeFields21(), msg, form)
	}
	if !xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return xscform.AdoptForm(formtype20, makeFields20(), msg, form)
	}
	return nil
}

var formtype22 = &xscmsg.MessageType{
	Tag:     Tag22,
	Name:    "OA jurisdiction status form",
	Article: "an",
	HTML:    "form-oa-muni-status.html",
	Version: "2.2",
}

func makeFields22() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		&xscform.ChoicesField{Field: *xscform.NewField(reportTypeID, true), Choices: reportTypeChoices},
		&jurisdictionCodeField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(jurisdictionCodeID, true), Choices: jurisdictionChoices}},
		xscform.NewField(jurisdictionID22, true),
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(eocPhoneID, false)}},
		&xscform.PhoneNumberField{Field: *xscform.NewField(eocFaxID, false)},
		&requiredForCompleteField{Field: *xscform.NewField(priEmContactNameID, false)},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(priEmContactPhoneID, false)}},
		xscform.NewField(secEmContactNameID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(secEmContactPhoneID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(officeStatusID, false), Choices: officeStatusChoices}},
		&xscform.DateField{Field: *xscform.NewField(govExpectedOpenDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(govExpectedOpenTimeID, false)},
		&xscform.DateField{Field: *xscform.NewField(govExpectedCloseDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(govExpectedCloseTimeID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(eocOpenID, false), Choices: eocOpenChoices}},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(eocActivationLevelID, false), Choices: eocActivationLevelChoices}},
		&xscform.DateField{Field: *xscform.NewField(eocExpectedOpenDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(eocExpectedOpenTimeID, false)},
		&xscform.DateField{Field: *xscform.NewField(eocExpectedCloseDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(eocExpectedCloseTimeID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(stateOfEmergencyID, false), Choices: stateofEmergencyChoices}},
		xscform.NewField(howSentID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(communicationsID, false), Choices: situationChoices}},
		xscform.NewField(communicationsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(debrisID, false), Choices: situationChoices}},
		xscform.NewField(debrisCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(floodingID, false), Choices: situationChoices}},
		xscform.NewField(floodingCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(hazmatID, false), Choices: situationChoices}},
		xscform.NewField(hazmatCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(emergencyServicesID, false), Choices: situationChoices}},
		xscform.NewField(emergencyServicesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(casualtiesID, false), Choices: situationChoices}},
		xscform.NewField(casualtiesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(utilitiesGasID, false), Choices: situationChoices}},
		xscform.NewField(utilitiesGasCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(utilitiesElectricID, false), Choices: situationChoices}},
		xscform.NewField(utilitiesElectricCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructurePowerID, false), Choices: situationChoices}},
		xscform.NewField(infrastructurePowerCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructureWaterSystemsID, false), Choices: situationChoices}},
		xscform.NewField(infrastructureWaterSystemsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructureSewerSystemsID, false), Choices: situationChoices}},
		xscform.NewField(infrastructureSewerSystemsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(searchAndRescueID, false), Choices: situationChoices}},
		xscform.NewField(searchAndRescueCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(transportationRoadsID, false), Choices: situationChoices}},
		xscform.NewField(transportationRoadsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(transportationBridgesID, false), Choices: situationChoices}},
		xscform.NewField(transportationBridgesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(civilUnrestID, false), Choices: situationChoices}},
		xscform.NewField(civilUnrestCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(animalIssuesID, false), Choices: situationChoices}},
		xscform.NewField(animalIssuesCommentsID, false),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var formtype21 = &xscmsg.MessageType{
	Tag:     Tag20,
	Name:    "OA municipal status form",
	Article: "an",
	HTML:    "form-oa-muni-status.html",
	Version: "2.1",
}

func makeFields21() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		&xscform.ChoicesField{Field: *xscform.NewField(reportTypeID, true), Choices: reportTypeChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(jurisdictionID20, true), Choices: jurisdictionChoices},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(eocPhoneID, false)}},
		&xscform.PhoneNumberField{Field: *xscform.NewField(eocFaxID, false)},
		&requiredForCompleteField{Field: *xscform.NewField(priEmContactNameID, false)},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(priEmContactPhoneID, false)}},
		xscform.NewField(secEmContactNameID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(secEmContactPhoneID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(officeStatusID, false), Choices: officeStatusChoices}},
		&xscform.DateField{Field: *xscform.NewField(govExpectedOpenDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(govExpectedOpenTimeID, false)},
		&xscform.DateField{Field: *xscform.NewField(govExpectedCloseDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(govExpectedCloseTimeID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(eocOpenID, false), Choices: eocOpenChoices}},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(eocActivationLevelID, false), Choices: eocActivationLevelChoices}},
		&xscform.DateField{Field: *xscform.NewField(eocExpectedOpenDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(eocExpectedOpenTimeID, false)},
		&xscform.DateField{Field: *xscform.NewField(eocExpectedCloseDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(eocExpectedCloseTimeID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(stateOfEmergencyID, false), Choices: stateofEmergencyChoices}},
		xscform.NewField(howSentID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(communicationsID, false), Choices: situationChoices}},
		xscform.NewField(communicationsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(debrisID, false), Choices: situationChoices}},
		xscform.NewField(debrisCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(floodingID, false), Choices: situationChoices}},
		xscform.NewField(floodingCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(hazmatID, false), Choices: situationChoices}},
		xscform.NewField(hazmatCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(emergencyServicesID, false), Choices: situationChoices}},
		xscform.NewField(emergencyServicesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(casualtiesID, false), Choices: situationChoices}},
		xscform.NewField(casualtiesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(utilitiesGasID, false), Choices: situationChoices}},
		xscform.NewField(utilitiesGasCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(utilitiesElectricID, false), Choices: situationChoices}},
		xscform.NewField(utilitiesElectricCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructurePowerID, false), Choices: situationChoices}},
		xscform.NewField(infrastructurePowerCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructureWaterSystemsID, false), Choices: situationChoices}},
		xscform.NewField(infrastructureWaterSystemsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructureSewerSystemsID, false), Choices: situationChoices}},
		xscform.NewField(infrastructureSewerSystemsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(searchAndRescueID, false), Choices: situationChoices}},
		xscform.NewField(searchAndRescueCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(transportationRoadsID, false), Choices: situationChoices}},
		xscform.NewField(transportationRoadsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(transportationBridgesID, false), Choices: situationChoices}},
		xscform.NewField(transportationBridgesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(civilUnrestID, false), Choices: situationChoices}},
		xscform.NewField(civilUnrestCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(animalIssuesID, false), Choices: situationChoices}},
		xscform.NewField(animalIssuesCommentsID, false),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var formtype20 = &xscmsg.MessageType{
	Tag:     Tag20,
	Name:    "OA municipal status form",
	Article: "an",
	HTML:    "form-oa-muni-status.html",
	Version: "2.0",
}

func makeFields20() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		&xscform.ChoicesField{Field: *xscform.NewField(reportTypeID, true), Choices: reportTypeChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(jurisdictionID20, true), Choices: jurisdictionChoices},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(eocPhoneID, false)}},
		&xscform.PhoneNumberField{Field: *xscform.NewField(eocFaxID, false)},
		&requiredForCompleteField{Field: *xscform.NewField(priEmContactNameID, false)},
		&requiredForCompletePhoneNumberField{PhoneNumberField: xscform.PhoneNumberField{Field: *xscform.NewField(priEmContactPhoneID, false)}},
		xscform.NewField(secEmContactNameID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(secEmContactPhoneID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(officeStatusID, false), Choices: officeStatusChoices}},
		&xscform.DateField{Field: *xscform.NewField(govExpectedOpenDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(govExpectedOpenTimeID, false)},
		&xscform.DateField{Field: *xscform.NewField(govExpectedCloseDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(govExpectedCloseTimeID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(eocOpenID, false), Choices: eocOpenChoices}},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(eocActivationLevelID, false), Choices: eocActivationLevelChoices}},
		&xscform.DateField{Field: *xscform.NewField(eocExpectedOpenDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(eocExpectedOpenTimeID, false)},
		&xscform.DateField{Field: *xscform.NewField(eocExpectedCloseDateID, false)},
		&xscform.TimeField{Field: *xscform.NewField(eocExpectedCloseTimeID, false)},
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(stateOfEmergencyID, false), Choices: stateofEmergencyChoices}},
		xscform.NewField(howToSendID, false),
		xscform.NewField(howSentID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(communicationsID, false), Choices: situationChoices}},
		xscform.NewField(communicationsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(debrisID, false), Choices: situationChoices}},
		xscform.NewField(debrisCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(floodingID, false), Choices: situationChoices}},
		xscform.NewField(floodingCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(hazmatID, false), Choices: situationChoices}},
		xscform.NewField(hazmatCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(emergencyServicesID, false), Choices: situationChoices}},
		xscform.NewField(emergencyServicesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(casualtiesID, false), Choices: situationChoices}},
		xscform.NewField(casualtiesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(utilitiesGasID, false), Choices: situationChoices}},
		xscform.NewField(utilitiesGasCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(utilitiesElectricID, false), Choices: situationChoices}},
		xscform.NewField(utilitiesElectricCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructurePowerID, false), Choices: situationChoices}},
		xscform.NewField(infrastructurePowerCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructureWaterSystemsID, false), Choices: situationChoices}},
		xscform.NewField(infrastructureWaterSystemsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(infrastructureSewerSystemsID, false), Choices: situationChoices}},
		xscform.NewField(infrastructureSewerSystemsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(searchAndRescueID, false), Choices: situationChoices}},
		xscform.NewField(searchAndRescueCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(transportationRoadsID, false), Choices: situationChoices}},
		xscform.NewField(transportationRoadsCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(transportationBridgesID, false), Choices: situationChoices}},
		xscform.NewField(transportationBridgesCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(civilUnrestID, false), Choices: situationChoices}},
		xscform.NewField(civilUnrestCommentsID, false),
		&requiredForCompleteChoicesField{ChoicesField: xscform.ChoicesField{Field: *xscform.NewField(animalIssuesID, false), Choices: situationChoices}},
		xscform.NewField(animalIssuesCommentsID, false),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var (
	reportTypeID = &xscmsg.FieldID{
		Tag:        "19.",
		Annotation: "report-type",
		Label:      "Report Type",
		Comment:    "required: Update, Complete",
	}
	reportTypeChoices  = []string{"Update", "Complete"}
	jurisdictionCodeID = &xscmsg.FieldID{
		Tag:        "21.",
		Annotation: "jurisdiction-code",
		ReadOnly:   true,
	}
	jurisdictionID20 = &xscmsg.FieldID{
		Tag:        "21.",
		Annotation: "jurisdiction",
		Label:      "Jurisdiction Name",
		Comment:    "required: Campbell, Cupertino, Gilroy, Los Altos, Los Altos Hills, Los Gatos, Milpitas, Monte Sereno, Morgan Hill, Mountain View, Palo Alto, San Jose, Santa Clara, Saratoga, Sunnyvale, Unincorporated",
		Canonical:  xscmsg.FSubject,
	}
	jurisdictionChoices = []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"}
	jurisdictionID22    = &xscmsg.FieldID{
		Tag:        "22.",
		Annotation: "jurisdiction",
		Label:      "Jurisdiction Name",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
	}
	eocPhoneID = &xscmsg.FieldID{
		Tag:        "23.",
		Annotation: "eoc-phone",
		Label:      "EOC Phone",
		Comment:    "phone-number required-for-complete",
	}
	eocFaxID = &xscmsg.FieldID{
		Tag:        "24.",
		Annotation: "eoc-fax",
		Label:      "EOC Fax",
		Comment:    "phone-number",
	}
	priEmContactNameID = &xscmsg.FieldID{
		Tag:        "25.",
		Annotation: "pri-em-contact-name",
		Label:      "Primary EM Contact Name",
		Comment:    "required-for-complete",
	}
	priEmContactPhoneID = &xscmsg.FieldID{
		Tag:        "26.",
		Annotation: "pri-em-contact-phone",
		Label:      "Primary EM Contact Phone",
		Comment:    "phone-number required-for-complete",
	}
	secEmContactNameID = &xscmsg.FieldID{
		Tag:        "27.",
		Annotation: "sec-em-contact-name",
		Label:      "Secondary EM Contact Name",
	}
	secEmContactPhoneID = &xscmsg.FieldID{
		Tag:        "28.",
		Annotation: "sec-em-contact-phone",
		Label:      "Secondary EM Contact Phone",
		Comment:    "phone-number",
	}
	officeStatusID = &xscmsg.FieldID{
		Tag:        "29.",
		Annotation: "office-status",
		Label:      "Office Status",
		Comment:    "required-for-complete: Unknown, Open, Closed",
	}
	officeStatusChoices   = []string{"Unknown", "Open", "Closed"}
	govExpectedOpenDateID = &xscmsg.FieldID{
		Tag:        "30.",
		Annotation: "gov-expected-open-date",
		Label:      "Office Expected to Open Date",
		Comment:    "date",
	}
	govExpectedOpenTimeID = &xscmsg.FieldID{
		Tag:        "31.",
		Annotation: "gov-expected-open-time",
		Label:      "Office Expected to Open Time",
		Comment:    "time",
	}
	govExpectedCloseDateID = &xscmsg.FieldID{
		Tag:        "32.",
		Annotation: "gov-expected-close-date",
		Label:      "Office Expected to Close Date",
		Comment:    "date",
	}
	govExpectedCloseTimeID = &xscmsg.FieldID{
		Tag:        "33.",
		Annotation: "gov-expected-close-time",
		Label:      "Office Expected to Close Time",
		Comment:    "time",
	}
	eocOpenID = &xscmsg.FieldID{
		Tag:        "34.",
		Annotation: "eoc-open",
		Label:      "EOC Open",
		Comment:    "required-for-complete: Unknown, Yes, No",
	}
	eocOpenChoices       = []string{"Unknown", "Yes", "No"}
	eocActivationLevelID = &xscmsg.FieldID{
		Tag:        "35.",
		Annotation: "eoc-activation-level",
		Label:      "Activation Level",
		Comment:    "required-for-complete: Normal, Duty Officer, Monitor, Partial, Full",
	}
	eocActivationLevelChoices = []string{"Normal", "Duty Officer", "Monitor", "Partial", "Full"}
	eocExpectedOpenDateID     = &xscmsg.FieldID{
		Tag:        "36.",
		Annotation: "eoc-expected-open-date",
		Label:      "EOC Expected to Open Date",
		Comment:    "date",
	}
	eocExpectedOpenTimeID = &xscmsg.FieldID{
		Tag:        "37.",
		Annotation: "eoc-expected-open-time",
		Label:      "EOC Expected to Open Time",
		Comment:    "time",
	}
	eocExpectedCloseDateID = &xscmsg.FieldID{
		Tag:        "38.",
		Annotation: "eoc-expected-close-date",
		Label:      "EOC Expected to Close Date",
		Comment:    "date",
	}
	eocExpectedCloseTimeID = &xscmsg.FieldID{
		Tag:        "39.",
		Annotation: "eoc-expected-close-time",
		Label:      "EOC Expected to Close Time",
		Comment:    "time",
	}
	stateOfEmergencyID = &xscmsg.FieldID{
		Tag:        "40.",
		Annotation: "state-of-emergency",
		Label:      "State of Emergency",
		Comment:    "required-for-complete: Unknown, Yes, No",
	}
	stateofEmergencyChoices = []string{"Unknown", "Yes", "No"}
	howToSendID             = &xscmsg.FieldID{
		Tag:        "98.",
		Annotation: "how-to-send",
		Label:      "SOE Proclamation: How to Send",
		ReadOnly:   true,
	}
	howSentID = &xscmsg.FieldID{
		Tag:        "99.",
		Annotation: "how-sent",
		Label:      "SOE Proclamation: Indicate how sent (method/to)",
	}
	situationChoices = []string{"Unknown", "Normal", "Problem", "Failure", "Delayed", "Closed", "Early Out"}
	communicationsID = &xscmsg.FieldID{
		Tag:        "41.0.",
		Annotation: "communications",
		Label:      "Communications",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	communicationsCommentsID = &xscmsg.FieldID{
		Tag:        "41.1.",
		Annotation: "communications-comments",
		Label:      "Communications Comments",
	}
	debrisID = &xscmsg.FieldID{
		Tag:        "42.0.",
		Annotation: "debris",
		Label:      "Debris",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	debrisCommentsID = &xscmsg.FieldID{
		Tag:        "42.1.",
		Annotation: "debris-comments",
		Label:      "Debris Comments",
	}
	floodingID = &xscmsg.FieldID{
		Tag:        "43.0.",
		Annotation: "flooding",
		Label:      "Flooding",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	floodingCommentsID = &xscmsg.FieldID{
		Tag:        "43.1.",
		Annotation: "flooding-comments",
		Label:      "Flooding Comments",
	}
	hazmatID = &xscmsg.FieldID{
		Tag:        "44.0.",
		Annotation: "hazmat",
		Label:      "Hazmat",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	hazmatCommentsID = &xscmsg.FieldID{
		Tag:        "44.1.",
		Annotation: "hazmat-comments",
		Label:      "Hazmat Comments",
	}
	emergencyServicesID = &xscmsg.FieldID{
		Tag:        "45.0.",
		Annotation: "emergency-services",
		Label:      "Emergency Services",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	emergencyServicesCommentsID = &xscmsg.FieldID{
		Tag:        "45.1.",
		Annotation: "emergency-services-comments",
		Label:      "Emergency Services Comments",
	}
	casualtiesID = &xscmsg.FieldID{
		Tag:        "46.0.",
		Annotation: "casualties",
		Label:      "Casualties",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	casualtiesCommentsID = &xscmsg.FieldID{
		Tag:        "46.1.",
		Annotation: "casualties-comments",
		Label:      "Casualties Comments",
	}
	utilitiesGasID = &xscmsg.FieldID{
		Tag:        "47.0.",
		Annotation: "utilities-gas",
		Label:      "Utilities (Gas)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	utilitiesGasCommentsID = &xscmsg.FieldID{
		Tag:        "47.1.",
		Annotation: "utilities-gas-comments",
		Label:      "Utilities (Gas) Comments",
	}
	utilitiesElectricID = &xscmsg.FieldID{
		Tag:        "48.0.",
		Annotation: "utilities-electric",
		Label:      "Utilities (Electric)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	utilitiesElectricCommentsID = &xscmsg.FieldID{
		Tag:        "48.1.",
		Annotation: "utilities-electric-comments",
		Label:      "Utilities (Electric) Comments",
	}
	infrastructurePowerID = &xscmsg.FieldID{
		Tag:        "49.0.",
		Annotation: "infrastructure-power",
		Label:      "Infrastructure (Power)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	infrastructurePowerCommentsID = &xscmsg.FieldID{
		Tag:        "49.1.",
		Annotation: "infrastructure-power-comments",
		Label:      "Infrastructure (Power) Comments",
	}
	infrastructureWaterSystemsID = &xscmsg.FieldID{
		Tag:        "50.0.",
		Annotation: "infrastructure-water-systems",
		Label:      "Infrastructure (Water Systems)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	infrastructureWaterSystemsCommentsID = &xscmsg.FieldID{
		Tag:        "50.1.",
		Annotation: "infrastructure-water-systems-comments",
		Label:      "Infrastructure (Water Systems) Comments",
	}
	infrastructureSewerSystemsID = &xscmsg.FieldID{
		Tag:        "51.0.",
		Annotation: "infrastructure-sewer-systems",
		Label:      "Infrastructure (Sewer Systems)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	infrastructureSewerSystemsCommentsID = &xscmsg.FieldID{
		Tag:        "51.1.",
		Annotation: "infrastructure-sewer-systems-comments",
		Label:      "Infrastructure (Sewer Systems) Comments",
	}
	searchAndRescueID = &xscmsg.FieldID{
		Tag:        "52.0.",
		Annotation: "search-and-rescue",
		Label:      "Search and Rescue",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	searchAndRescueCommentsID = &xscmsg.FieldID{
		Tag:        "52.1.",
		Annotation: "search-and-rescue-comments",
		Label:      "Search and Rescue Comments",
	}
	transportationRoadsID = &xscmsg.FieldID{
		Tag:        "53.0.",
		Annotation: "transportation-roads",
		Label:      "Transportation (Roads)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	transportationRoadsCommentsID = &xscmsg.FieldID{
		Tag:        "53.1.",
		Annotation: "transportation-roads-comments",
		Label:      "Transportation (Roads) Comments",
	}
	transportationBridgesID = &xscmsg.FieldID{
		Tag:        "54.0.",
		Annotation: "transportation-bridges",
		Label:      "Transportation (Bridges)",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	transportationBridgesCommentsID = &xscmsg.FieldID{
		Tag:        "54.1.",
		Annotation: "transportation-bridges-comments",
		Label:      "Transportation (Bridges) Comments",
	}
	civilUnrestID = &xscmsg.FieldID{
		Tag:        "55.0.",
		Annotation: "civil-unrest",
		Label:      "Civil Unrest",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	civilUnrestCommentsID = &xscmsg.FieldID{
		Tag:        "55.1.",
		Annotation: "civil-unrest-comments",
		Label:      "Civil Unrest Comments",
	}
	animalIssuesID = &xscmsg.FieldID{
		Tag:        "56.0.",
		Annotation: "animal-issues",
		Label:      "Animal Issues",
		Comment:    "required-for-complete: Unknown, Normal, Problem, Failure, Delayed, Closed, Early Out",
	}
	animalIssuesCommentsID = &xscmsg.FieldID{
		Tag:        "56.1.",
		Annotation: "animal-issues-comments",
		Label:      "Animal Issues Comments",
	}
)

type handlingField struct{ xscform.ChoicesField }

func (f *handlingField) Default() string { return "IMMEDIATE" }

type toICSPositionField struct{ xscform.Field }

func (f *toICSPositionField) Default() string { return "Situation Analysis Unit" }

type toLocationField struct{ xscform.Field }

func (f *toLocationField) Default() string { return "County EOC" }

type jurisdictionCodeField struct{ xscform.ChoicesField }

func (f *jurisdictionCodeField) Validate(msg xscmsg.Message, strict bool) string {
	if f.Get() == "" {
		j := msg.Field("22.").Get()
		if j != "" {
			for _, choice := range jurisdictionChoices {
				if j == choice {
					f.Set(j)
					break
				}
			}
			if f.Get() == "" {
				f.Set(jurisdictionChoices[len(jurisdictionChoices)-1])
			}
		}
	}
	return f.ChoicesField.Validate(msg, strict)
}

type requiredForCompleteField struct{ xscform.Field }

func (f *requiredForCompleteField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.Field.Validate(msg, strict)
}

type requiredForCompleteChoicesField struct{ xscform.ChoicesField }

func (f *requiredForCompleteChoicesField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.ChoicesField.Validate(msg, strict)
}

type requiredForCompletePhoneNumberField struct{ xscform.PhoneNumberField }

func (f *requiredForCompletePhoneNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if err := validateRequiredIfComplete(f, msg); err != "" {
		return err
	}
	return f.PhoneNumberField.Validate(msg, strict)
}

func validateRequiredIfComplete(f xscmsg.Field, msg xscmsg.Message) string {
	if f.Get() == "" && msg.Field("19.").Get() == "Complete" {
		return fmt.Sprintf("field %q needs a value when field \"19.\" is \"Complete\"", f.ID().Tag)
	}
	return ""
}
