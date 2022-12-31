package ahfacstat

import (
	"fmt"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies allied health facility status forms.
const (
	Tag = "AHFacStat"
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)
}

func create() xscmsg.Message {
	return xscform.CreateForm(formtype, makeFields())
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.Message {
	if form == nil || form.FormType != formtype.HTML || xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return nil
	}
	return xscform.AdoptForm(formtype, makeFields(), msg, form)
}

var formtype = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "allied health facility status report",
	Article: "an",
	HTML:    "form-allied-health-facility-status.html",
	Version: "2.3",
}

func makeFields() []xscmsg.Field {
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
		xscform.NewField(facilityID, true),
		&requiredForComplete[*xscform.Field]{f: xscform.NewField(facilityTypeID, false)},
		&xscform.DateFieldDefaultNow{DateField: xscform.DateField{Field: *xscform.NewField(dateID, true)}},
		&xscform.TimeField{Field: *xscform.NewField(timeID, true)},
		&requiredForComplete[*xscform.Field]{f: xscform.NewField(contactID, false)},
		&requiredForComplete[*xscform.PhoneNumberField]{f: &xscform.PhoneNumberField{Field: *xscform.NewField(contactPhoneID, false)}},
		&xscform.PhoneNumberField{Field: *xscform.NewField(contactFaxID, false)},
		xscform.NewField(otherContactID, false),
		&requiredForComplete[*xscform.Field]{f: xscform.NewField(incidentID, false)},
		&requiredForComplete[*xscform.DateField]{f: &xscform.DateField{Field: *xscform.NewField(incidentDateID, false)}},
		&requiredForComplete[*xscform.ChoicesField]{f: &xscform.ChoicesField{Field: *xscform.NewField(statusID, false), Choices: statusChoices}},
		&xscform.ChoicesField{Field: *xscform.NewField(attachOrgChartID, false), Choices: yesNoChoices},
		&requiredForComplete[*xscform.ChoicesField]{f: &xscform.ChoicesField{Field: *xscform.NewField(attachRrID, false), Choices: yesNoChoices}},
		&xscform.ChoicesField{Field: *xscform.NewField(attachStatusID, false), Choices: yesNoChoices},
		&xscform.ChoicesField{Field: *xscform.NewField(attachActionPlanID, false), Choices: yesNoChoices},
		&requiredForComplete[*xscform.PhoneNumberField]{f: &xscform.PhoneNumberField{Field: *xscform.NewField(eocPhoneID, false)}},
		&xscform.ChoicesField{Field: *xscform.NewField(attachDirectoryID, false), Choices: yesNoChoices},
		&xscform.PhoneNumberField{Field: *xscform.NewField(eocFaxID, false)},
		&requiredForComplete[*xscform.Field]{f: xscform.NewField(liaisonID, false)},
		xscform.NewField(summaryID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(liaisonPhoneID, false)},
		xscform.NewField(infoOfficerID, false),
		&xscform.PhoneNumberField{Field: *xscform.NewField(infoOfficerPhoneID, false)},
		xscform.NewField(infoOfficerEmailID, false),
		&requiredForComplete[*xscform.Field]{f: xscform.NewField(eocClosedContactID, false)},
		&requiredForComplete[*xscform.PhoneNumberField]{f: &xscform.PhoneNumberField{Field: *xscform.NewField(eocPhone2ID, false)}},
		&xscform.CardinalNumberField{Field: *xscform.NewField(skilledNursingBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(skilledNursingBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(skilledNursingBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(skilledNursingBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(skilledNursingBedsSurgeID, false)},
		&requiredForComplete[*xscform.Field]{f: xscform.NewField(eocEmailID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(assistedLivingBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(assistedLivingBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(assistedLivingBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(assistedLivingBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(assistedLivingBedsSurgeID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(subAcuteBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(subAcuteBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(subAcuteBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(subAcuteBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(subAcuteBedsSurgeID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(toEvacuateID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(alzheimersBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(alzheimersBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(alzheimersBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(alzheimersBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(alzheimersBedsSurgeID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(injuredID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(pedSubAcuteBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(pedSubAcuteBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(pedSubAcuteBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(pedSubAcuteBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(pedSubAcuteBedsSurgeID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(transferredID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(psychiatricBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(psychiatricBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(psychiatricBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(psychiatricBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(psychiatricBedsSurgeID, false)},
		xscform.NewField(otherCareID, false),
		xscform.NewField(bedResourceID, false),
		&xscform.CardinalNumberField{Field: *xscform.NewField(otherCareBedsStaffedMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(otherCareBedsStaffedFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(otherCareBedsVacantMID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(otherCareBedsVacantFID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(otherCareBedsSurgeID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(dialysisChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(dialysisVacantChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(dialysisFrontStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(dialysisSupportStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(dialysisProvidersID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(surgicalChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(surgicalVacantChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(surgicalFrontStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(surgicalSupportStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(surgicalProvidersID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(clinicChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(clinicVacantChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(clinicFrontStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(clinicSupportStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(clinicProvidersID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(homeHealthChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(homeHealthVacantChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(homeHealthFrontStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(homeHealthSupportStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(homeHealthProvidersID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(adultDayCtrChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(adultDayCtrVacantChairsID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(adultDayCtrFrontStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(adultDayCtrSupportStaffID, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(adultDayCtrProvidersID, false)},
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
	reportTypeChoices = []string{"Update", "Complete"}
	facilityID        = &xscmsg.FieldID{
		Tag:        "20.",
		Annotation: "facility",
		Label:      "Facility Name",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
	}
	facilityTypeID = &xscmsg.FieldID{
		Tag:        "21.",
		Annotation: "facility-type",
		Label:      "Facility Type",
		Comment:    "required-for-complete",
	}
	dateID = &xscmsg.FieldID{
		Tag:        "22d.",
		Annotation: "date",
		Label:      "Date",
		Comment:    "required date",
	}
	timeID = &xscmsg.FieldID{
		Tag:        "22t.",
		Annotation: "time",
		Label:      "Time",
		Comment:    "required time",
	}
	contactID = &xscmsg.FieldID{
		Tag:        "23.",
		Annotation: "contact",
		Label:      "Contact Name",
		Comment:    "required-for-complete",
	}
	contactPhoneID = &xscmsg.FieldID{
		Tag:        "23p.",
		Annotation: "contact-phone",
		Label:      "Contact Phone #",
		Comment:    "phone-number required-for-complete",
	}
	contactFaxID = &xscmsg.FieldID{
		Tag:        "23f.",
		Annotation: "contact-fax",
		Label:      "Contact Fax #",
		Comment:    "phone-number",
	}
	otherContactID = &xscmsg.FieldID{
		Tag:        "24.",
		Annotation: "other-contact",
		Label:      "Other Phone, Fax, Cell Phone, Radio",
	}
	incidentID = &xscmsg.FieldID{
		Tag:        "25.",
		Annotation: "incident",
		Label:      "Incident Name",
		Comment:    "required-for-complete",
	}
	incidentDateID = &xscmsg.FieldID{
		Tag:        "25d.",
		Annotation: "incident-date",
		Label:      "Incident Date",
		Comment:    "date required-for-complete",
	}
	statusID = &xscmsg.FieldID{
		Tag:        "35.",
		Annotation: "status",
		Label:      "Facility Status",
		Comment:    "required-for-complete: Fully Functional, Limited Services, Impaired/Closed",
	}
	statusChoices    = []string{"Fully Functional", "Limited Services", "Impaired/Closed"}
	attachOrgChartID = &xscmsg.FieldID{
		Tag:        "26a.",
		Annotation: "attach-org-chart",
		Label:      "NHICS/ICS Organization Chart Attached",
		Comment:    "Yes, No",
	}
	attachRrID = &xscmsg.FieldID{
		Tag:        "26b.",
		Annotation: "attach-RR",
		Label:      "DEOC-9A Resource Request Forms Attached",
		Comment:    "required-for-complete: Yes, No",
	}
	attachStatusID = &xscmsg.FieldID{
		Tag:        "26c.",
		Annotation: "attach-status",
		Label:      "NHICS/ICS Status Report Form - Standard Attached",
		Comment:    "Yes, No",
	}
	attachActionPlanID = &xscmsg.FieldID{
		Tag:        "26d.",
		Annotation: "attach-action-plan",
		Label:      "NHICS/ICS Incident Action Plan Attached",
		Comment:    "Yes, No",
	}
	eocPhoneID = &xscmsg.FieldID{
		Tag:        "27p.",
		Annotation: "eoc-phone",
		Label:      "Facility EOC Main Contact Number",
		Comment:    "phone-number required-for-complete",
	}
	attachDirectoryID = &xscmsg.FieldID{
		Tag:        "26e.",
		Annotation: "attach-directory",
		Label:      "Phone/Communications Directory Attached",
		Comment:    "Yes, No",
	}
	yesNoChoices = []string{"Yes", "No"}
	eocFaxID     = &xscmsg.FieldID{
		Tag:        "27f.",
		Annotation: "eoc-fax",
		Label:      "Facility EOC Main Contact Fax",
		Comment:    "phone-number",
	}
	liaisonID = &xscmsg.FieldID{
		Tag:        "28.",
		Annotation: "liaison",
		Label:      "Facility Liaison Officer Name: Liaison to Public Health/Medical Health Branch",
		Comment:    "required-for-complete",
	}
	summaryID = &xscmsg.FieldID{
		Tag:        "34.",
		Annotation: "summary",
		Label:      "General Summary of Situation/Conditions",
	}
	liaisonPhoneID = &xscmsg.FieldID{
		Tag:        "28p.",
		Annotation: "liaison-phone",
		Label:      "Facility Liaison Contact Number",
		Comment:    "phone-number",
	}
	infoOfficerID = &xscmsg.FieldID{
		Tag:        "29.",
		Annotation: "info-officer",
		Label:      "Facility Information Officer Name",
	}
	infoOfficerPhoneID = &xscmsg.FieldID{
		Tag:        "29p.",
		Annotation: "info-officer-phone",
		Label:      "Facility Information Officer Contact Number",
		Comment:    "phone-number",
	}
	infoOfficerEmailID = &xscmsg.FieldID{
		Tag:        "29e.",
		Annotation: "info-officer-email",
		Label:      "Facility Information Officer Contact Email",
	}
	eocClosedContactID = &xscmsg.FieldID{
		Tag:        "30.",
		Annotation: "eoc-closed-contact",
		Label:      "If Facility EOC is Not Activated, Who Should be Contacted for Questions/Requests",
		Comment:    "required-for-complete",
	}
	eocPhone2ID = &xscmsg.FieldID{
		Tag:        "30p.",
		Annotation: "eoc-phone",
		Label:      "Facility Contact Number",
		Comment:    "phone-number required-for-complete",
	}
	skilledNursingBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "40a.",
		Annotation: "skilled-nursing-beds-staffed-m",
		Label:      "Skilled Nursing Beds Staffed M",
		Comment:    "cardinal-number",
	}
	skilledNursingBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "40b.",
		Annotation: "skilled-nursing-beds-staffed-f",
		Label:      "Skilled Nursing Beds Staffed F",
		Comment:    "cardinal-number",
	}
	skilledNursingBedsVacantMID = &xscmsg.FieldID{
		Tag:        "40c.",
		Annotation: "skilled-nursing-beds-vacant-m",
		Label:      "Skilled Nursing Beds Vacant M",
		Comment:    "cardinal-number",
	}
	skilledNursingBedsVacantFID = &xscmsg.FieldID{
		Tag:        "40d.",
		Annotation: "skilled-nursing-beds-vacant-f",
		Label:      "Skilled Nursing Beds Vacant F",
		Comment:    "cardinal-number",
	}
	skilledNursingBedsSurgeID = &xscmsg.FieldID{
		Tag:        "40e.",
		Annotation: "skilled-nursing-beds-surge",
		Label:      "Skilled Nursing Beds Surge",
		Comment:    "cardinal-number",
	}
	eocEmailID = &xscmsg.FieldID{
		Tag:        "30e.",
		Annotation: "eoc-email",
		Label:      "Facility Contact Email",
		Comment:    "required-for-complete",
	}
	assistedLivingBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "41a.",
		Annotation: "assisted-living-beds-staffed-m",
		Label:      "Assisted Living Beds Staffed M",
		Comment:    "cardinal-number",
	}
	assistedLivingBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "41b.",
		Annotation: "assisted-living-beds-staffed-f",
		Label:      "Assisted Living Beds Staffed F",
		Comment:    "cardinal-number",
	}
	assistedLivingBedsVacantMID = &xscmsg.FieldID{
		Tag:        "41c.",
		Annotation: "assisted-living-beds-vacant-m",
		Label:      "Assisted Living Beds Vacant M",
		Comment:    "cardinal-number",
	}
	assistedLivingBedsVacantFID = &xscmsg.FieldID{
		Tag:        "41d.",
		Annotation: "assisted-living-beds-vacant-f",
		Label:      "Assisted Living Beds Vacant F",
		Comment:    "cardinal-number",
	}
	assistedLivingBedsSurgeID = &xscmsg.FieldID{
		Tag:        "41e.",
		Annotation: "assisted-living-beds-surge",
		Label:      "Assisted Living Beds Surge",
		Comment:    "cardinal-number",
	}
	subAcuteBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "42a.",
		Annotation: "sub-acute-beds-staffed-m",
		Label:      "Sub-Acute Beds Staffed M",
		Comment:    "cardinal-number",
	}
	subAcuteBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "42b.",
		Annotation: "sub-acute-beds-staffed-f",
		Label:      "Sub-Acute Beds Staffed F",
		Comment:    "cardinal-number",
	}
	subAcuteBedsVacantMID = &xscmsg.FieldID{
		Tag:        "42c.",
		Annotation: "sub-acute-beds-vacant-m",
		Label:      "Sub-Acute Beds Vacant M",
		Comment:    "cardinal-number",
	}
	subAcuteBedsVacantFID = &xscmsg.FieldID{
		Tag:        "42d.",
		Annotation: "sub-acute-beds-vacant-f",
		Label:      "Sub-Acute Beds Vacant F",
		Comment:    "cardinal-number",
	}
	subAcuteBedsSurgeID = &xscmsg.FieldID{
		Tag:        "42e.",
		Annotation: "sub-acute-beds-surge",
		Label:      "Sub-Acute Beds Surge",
		Comment:    "cardinal-number",
	}
	toEvacuateID = &xscmsg.FieldID{
		Tag:        "31a.",
		Annotation: "to-evacuate",
		Label:      "Facility Patients to Evacuate",
		Comment:    "cardinal-number",
	}
	alzheimersBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "43a.",
		Annotation: "alzheimers-beds-staffed-m",
		Label:      "Alzheimers/Dementia Beds Staffed M",
		Comment:    "cardinal-number",
	}
	alzheimersBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "43b.",
		Annotation: "alzheimers-beds-staffed-f",
		Label:      "Alzheimers/Dementia Beds Staffed F",
		Comment:    "cardinal-number",
	}
	alzheimersBedsVacantMID = &xscmsg.FieldID{
		Tag:        "43c.",
		Annotation: "alzheimers-beds-vacant-m",
		Label:      "Alzheimers/Dementia Beds Vacant M",
		Comment:    "cardinal-number",
	}
	alzheimersBedsVacantFID = &xscmsg.FieldID{
		Tag:        "43d.",
		Annotation: "alzheimers-beds-vacant-f",
		Label:      "Alzheimers/Dementia Beds Vacant F",
		Comment:    "cardinal-number",
	}
	alzheimersBedsSurgeID = &xscmsg.FieldID{
		Tag:        "43e.",
		Annotation: "alzheimers-beds-surge",
		Label:      "Alzheimers/Dementia Beds Surge",
		Comment:    "cardinal-number",
	}
	injuredID = &xscmsg.FieldID{
		Tag:        "31b.",
		Annotation: "injured",
		Label:      "Facility Patients Injured - Minor",
		Comment:    "cardinal-number",
	}
	pedSubAcuteBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "44a.",
		Annotation: "ped-sub-acute-beds-staffed-m",
		Label:      "Pediatric-Sub Acute Beds Staffed M",
		Comment:    "cardinal-number",
	}
	pedSubAcuteBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "44b.",
		Annotation: "ped-sub-acute-beds-staffed-f",
		Label:      "Pediatric-Sub Acute Beds Staffed F",
		Comment:    "cardinal-number",
	}
	pedSubAcuteBedsVacantMID = &xscmsg.FieldID{
		Tag:        "44c.",
		Annotation: "ped-sub-acute-beds-vacant-m",
		Label:      "Pediatric-Sub Acute Beds Vacant M",
		Comment:    "cardinal-number",
	}
	pedSubAcuteBedsVacantFID = &xscmsg.FieldID{
		Tag:        "44d.",
		Annotation: "ped-sub-acute-beds-vacant-f",
		Label:      "Pediatric-Sub Acute Beds Vacant F",
		Comment:    "cardinal-number",
	}
	pedSubAcuteBedsSurgeID = &xscmsg.FieldID{
		Tag:        "44e.",
		Annotation: "ped-sub-acute-beds-surge",
		Label:      "Pediatric-Sub Acute Beds Surge",
		Comment:    "cardinal-number",
	}
	transferredID = &xscmsg.FieldID{
		Tag:        "31c.",
		Annotation: "transferred",
		Label:      "Facility Patients Transferred Out of County",
		Comment:    "cardinal-number",
	}
	psychiatricBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "45a.",
		Annotation: "psychiatric-beds-staffed-m",
		Label:      "Psychiatric Beds Staffed M",
		Comment:    "cardinal-number",
	}
	psychiatricBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "45b.",
		Annotation: "psychiatric-beds-staffed-f",
		Label:      "Psychiatric Beds Staffed F",
		Comment:    "cardinal-number",
	}
	psychiatricBedsVacantMID = &xscmsg.FieldID{
		Tag:        "45c.",
		Annotation: "psychiatric-beds-vacant-m",
		Label:      "Psychiatric Beds Vacant M",
		Comment:    "cardinal-number",
	}
	psychiatricBedsVacantFID = &xscmsg.FieldID{
		Tag:        "45d.",
		Annotation: "psychiatric-beds-vacant-f",
		Label:      "Psychiatric Beds Vacant F",
		Comment:    "cardinal-number",
	}
	psychiatricBedsSurgeID = &xscmsg.FieldID{
		Tag:        "45e.",
		Annotation: "psychiatric-beds-surge",
		Label:      "Psychiatric Beds Surge",
		Comment:    "cardinal-number",
	}
	otherCareID = &xscmsg.FieldID{
		Tag:        "33.",
		Annotation: "other-care",
		Label:      "Other Facility Patient Care Information",
	}
	bedResourceID = &xscmsg.FieldID{
		Tag:        "46.",
		Annotation: "bed-resource",
		Label:      "Other Care",
	}
	otherCareBedsStaffedMID = &xscmsg.FieldID{
		Tag:        "46a.",
		Annotation: "other-care-beds-staffed-m",
		Label:      "Other Care Beds Staffed M",
		Comment:    "cardinal-number",
	}
	otherCareBedsStaffedFID = &xscmsg.FieldID{
		Tag:        "46b.",
		Annotation: "other-care-beds-staffed-f",
		Label:      "Other Care Beds Staffed F",
		Comment:    "cardinal-number",
	}
	otherCareBedsVacantMID = &xscmsg.FieldID{
		Tag:        "46c.",
		Annotation: "other-care-beds-vacant-m",
		Label:      "Other Care Beds Vacant M",
		Comment:    "cardinal-number",
	}
	otherCareBedsVacantFID = &xscmsg.FieldID{
		Tag:        "46d.",
		Annotation: "other-care-beds-vacant-f",
		Label:      "Other Care Beds Vacant F",
		Comment:    "cardinal-number",
	}
	otherCareBedsSurgeID = &xscmsg.FieldID{
		Tag:        "46e.",
		Annotation: "other-care-beds-surge",
		Label:      "Other Care Beds Surge",
		Comment:    "cardinal-number",
	}
	dialysisChairsID = &xscmsg.FieldID{
		Tag:        "50a.",
		Annotation: "dialysis-chairs",
		Label:      "Dialysis Chairs/Room",
		Comment:    "cardinal-number",
	}
	dialysisVacantChairsID = &xscmsg.FieldID{
		Tag:        "50b.",
		Annotation: "dialysis-vacant-chairs",
		Label:      "Dialysis Vacant Chairs/Room",
		Comment:    "cardinal-number",
	}
	dialysisFrontStaffID = &xscmsg.FieldID{
		Tag:        "50c.",
		Annotation: "dialysis-front-staff",
		Label:      "Dialysis Front Desk Staff",
		Comment:    "cardinal-number",
	}
	dialysisSupportStaffID = &xscmsg.FieldID{
		Tag:        "50d.",
		Annotation: "dialysis-support-staff",
		Label:      "Dialysis Medical Support Staff",
		Comment:    "cardinal-number",
	}
	dialysisProvidersID = &xscmsg.FieldID{
		Tag:        "50e.",
		Annotation: "dialysis-providers",
		Label:      "Dialysis Provider Staff",
		Comment:    "cardinal-number",
	}
	surgicalChairsID = &xscmsg.FieldID{
		Tag:        "51a.",
		Annotation: "surgical-chairs",
		Label:      "Surgical Chairs/Room",
		Comment:    "cardinal-number",
	}
	surgicalVacantChairsID = &xscmsg.FieldID{
		Tag:        "51b.",
		Annotation: "surgical-vacant-chairs",
		Label:      "Surgical Vacant Chairs/Room",
		Comment:    "cardinal-number",
	}
	surgicalFrontStaffID = &xscmsg.FieldID{
		Tag:        "51c.",
		Annotation: "surgical-front-staff",
		Label:      "Surgical Front Desk Staff",
		Comment:    "cardinal-number",
	}
	surgicalSupportStaffID = &xscmsg.FieldID{
		Tag:        "51d.",
		Annotation: "surgical-support-staff",
		Label:      "Surgical Medical Support Staff",
		Comment:    "cardinal-number",
	}
	surgicalProvidersID = &xscmsg.FieldID{
		Tag:        "51e.",
		Annotation: "surgical-providers",
		Label:      "Surgical Provider Staff",
		Comment:    "cardinal-number",
	}
	clinicChairsID = &xscmsg.FieldID{
		Tag:        "52a.",
		Annotation: "clinic-chairs",
		Label:      "Clinic Chairs/Room",
		Comment:    "cardinal-number",
	}
	clinicVacantChairsID = &xscmsg.FieldID{
		Tag:        "52b.",
		Annotation: "clinic-vacant-chairs",
		Label:      "Clinic Vacant Chairs/Room",
		Comment:    "cardinal-number",
	}
	clinicFrontStaffID = &xscmsg.FieldID{
		Tag:        "52c.",
		Annotation: "clinic-front-staff",
		Label:      "Clinic Front Desk Staff",
		Comment:    "cardinal-number",
	}
	clinicSupportStaffID = &xscmsg.FieldID{
		Tag:        "52d.",
		Annotation: "clinic-support-staff",
		Label:      "Clinic Medical Support Staff",
		Comment:    "cardinal-number",
	}
	clinicProvidersID = &xscmsg.FieldID{
		Tag:        "52e.",
		Annotation: "clinic-providers",
		Label:      "Clinic Provider Staff",
		Comment:    "cardinal-number",
	}
	homeHealthChairsID = &xscmsg.FieldID{
		Tag:        "53a.",
		Annotation: "home-health-chairs",
		Label:      "Home Health Chairs/Room",
		Comment:    "cardinal-number",
	}
	homeHealthVacantChairsID = &xscmsg.FieldID{
		Tag:        "53b.",
		Annotation: "home-health-vacant-chairs",
		Label:      "Home Health Vacant Chairs/Room",
		Comment:    "cardinal-number",
	}
	homeHealthFrontStaffID = &xscmsg.FieldID{
		Tag:        "53c.",
		Annotation: "home-health-front-staff",
		Label:      "Home Health Front Desk Staff",
		Comment:    "cardinal-number",
	}
	homeHealthSupportStaffID = &xscmsg.FieldID{
		Tag:        "53d.",
		Annotation: "home-health-support-staff",
		Label:      "Home Health Medical Support Staff",
		Comment:    "cardinal-number",
	}
	homeHealthProvidersID = &xscmsg.FieldID{
		Tag:        "53e.",
		Annotation: "home-health-providers",
		Label:      "Home Health Provider Staff",
		Comment:    "cardinal-number",
	}
	adultDayCtrChairsID = &xscmsg.FieldID{
		Tag:        "54a.",
		Annotation: "adult-day-ctr-chairs",
		Label:      "Adult Day Ctr Chairs/Room",
		Comment:    "cardinal-number",
	}
	adultDayCtrVacantChairsID = &xscmsg.FieldID{
		Tag:        "54b.",
		Annotation: "adult-day-ctr-vacant-chairs",
		Label:      "Adult Day Ctr Vacant Chairs/Room",
		Comment:    "cardinal-number",
	}
	adultDayCtrFrontStaffID = &xscmsg.FieldID{
		Tag:        "54c.",
		Annotation: "adult-day-ctr-front-staff",
		Label:      "Adult Day Ctr Front Desk Staff",
		Comment:    "cardinal-number",
	}
	adultDayCtrSupportStaffID = &xscmsg.FieldID{
		Tag:        "54d.",
		Annotation: "adult-day-ctr-support-staff",
		Label:      "Adult Day Ctr Medical Support Staff",
		Comment:    "cardinal-number",
	}
	adultDayCtrProvidersID = &xscmsg.FieldID{
		Tag:        "54e.",
		Annotation: "adult-day-ctr-providers",
		Label:      "Adult Day Ctr Provider Staff",
		Comment:    "cardinal-number",
	}
)

type handlingField struct{ xscform.ChoicesField }

func (f *handlingField) Default() string { return "ROUTINE" }

type toICSPositionField struct{ xscform.Field }

func (f *toICSPositionField) Default() string { return "EMS Unit" }

type toLocationField struct{ xscform.Field }

func (f *toLocationField) Default() string { return "MHJOC" }

type requiredForComplete[T xscmsg.Field] struct{ f T }

func (f *requiredForComplete[T]) Validate(msg xscmsg.Message, strict bool) string {
	if f.f.Get() == "" && msg.Field("19.").Get() == "Complete" {
		return fmt.Sprintf("field %q needs a value when field \"19.\" is \"Complete\"", f.ID().Tag)
	}
	return f.f.Validate(msg, strict)
}
func (f *requiredForComplete[T]) ID() *xscmsg.FieldID { return f.f.ID() }
func (f *requiredForComplete[T]) Get() string         { return f.f.Get() }
func (f *requiredForComplete[T]) Set(v string)        { f.f.Set(v) }
func (f *requiredForComplete[T]) Default() string     { return f.f.Default() }
