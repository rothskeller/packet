package ahfacstat

import (
	"fmt"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies allied health facility status forms.
const Tag = "AHFacStat"

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef.DefaultValue = "ROUTINE"
	toICSPositionDef.DefaultValue = "EMS Unit"
	toICSPositionDef.Comment = "required: EMS Unit, Public Health Unit, Medical Health Branch, Operations Section, ..."
	toLocationDef.DefaultValue = "MHJOC"
	toLocationDef.Comment = "required: MHJOC, County EOC, ..."
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype, fieldDefs)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil || form.FormType != formtype.HTML || xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return nil
	}
	return xscform.AdoptForm(formtype, fieldDefs, msg, form)
}

var formtype = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "allied health facility status report",
	Article:     "an",
	HTML:        "form-allied-health-facility-status.html",
	Version:     "2.3",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

var fieldDefs = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// Allied Health Facility Status fields
	reportTypeDef, facilityDef, facilityTypeDef, dateDef, timeDef, contactDef, contactPhoneDef, contactFaxDef, otherContactDef,
	incidentDef, incidentDateDef, statusDef, attachOrgChartDef, attachRrDef, attachStatusDef, attachActionPlanDef, eocPhoneDef,
	attachDirectoryDef, eocFaxDef, liaisonDef, summaryDef, liaisonPhoneDef, infoOfficerDef, infoOfficerPhoneDef,
	infoOfficerEmailDef, eocClosedContactDef, eocPhone2Def, skilledNursingBedsStaffedMDef, skilledNursingBedsStaffedFDef,
	skilledNursingBedsVacantMDef, skilledNursingBedsVacantFDef, skilledNursingBedsSurgeDef, eocEmailDef,
	assistedLivingBedsStaffedMDef, assistedLivingBedsStaffedFDef, assistedLivingBedsVacantMDef, assistedLivingBedsVacantFDef,
	assistedLivingBedsSurgeDef, subAcuteBedsStaffedMDef, subAcuteBedsStaffedFDef, subAcuteBedsVacantMDef,
	subAcuteBedsVacantFDef, subAcuteBedsSurgeDef, toEvacuateDef, alzheimersBedsStaffedMDef, alzheimersBedsStaffedFDef,
	alzheimersBedsVacantMDef, alzheimersBedsVacantFDef, alzheimersBedsSurgeDef, injuredDef, pedSubAcuteBedsStaffedMDef,
	pedSubAcuteBedsStaffedFDef, pedSubAcuteBedsVacantMDef, pedSubAcuteBedsVacantFDef, pedSubAcuteBedsSurgeDef, transferredDef,
	psychiatricBedsStaffedMDef, psychiatricBedsStaffedFDef, psychiatricBedsVacantMDef, psychiatricBedsVacantFDef,
	psychiatricBedsSurgeDef, otherCareDef, bedResourceDef, otherCareBedsStaffedMDef, otherCareBedsStaffedFDef,
	otherCareBedsVacantMDef, otherCareBedsVacantFDef, otherCareBedsSurgeDef, dialysisChairsDef, dialysisVacantChairsDef,
	dialysisFrontStaffDef, dialysisSupportStaffDef, dialysisProvidersDef, surgicalChairsDef, surgicalVacantChairsDef,
	surgicalFrontStaffDef, surgicalSupportStaffDef, surgicalProvidersDef, clinicChairsDef, clinicVacantChairsDef,
	clinicFrontStaffDef, clinicSupportStaffDef, clinicProvidersDef, homeHealthChairsDef, homeHealthVacantChairsDef,
	homeHealthFrontStaffDef, homeHealthSupportStaffDef, homeHealthProvidersDef, adultDayCtrChairsDef,
	adultDayCtrVacantChairsDef, adultDayCtrFrontStaffDef, adultDayCtrSupportStaffDef, adultDayCtrProvidersDef,
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
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"Update", "Complete"},
	}
	facilityDef = &xscmsg.FieldDef{
		Tag:        "20.",
		Annotation: "facility",
		Label:      "Facility Name",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	facilityTypeDef = &xscmsg.FieldDef{
		Tag:        "21.",
		Annotation: "facility-type",
		Label:      "Facility Type",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	dateDef = &xscmsg.FieldDef{
		Tag:         "22d.",
		Annotation:  "date",
		Label:       "Date",
		Comment:     "required date",
		Validators:  []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateDate},
		DefaultFunc: xscform.DefaultDate,
	}
	timeDef = &xscmsg.FieldDef{
		Tag:        "22t.",
		Annotation: "time",
		Label:      "Time",
		Comment:    "required time",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateTime},
	}
	contactDef = &xscmsg.FieldDef{
		Tag:        "23.",
		Annotation: "contact",
		Label:      "Contact Name",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	contactPhoneDef = &xscmsg.FieldDef{
		Tag:        "23p.",
		Annotation: "contact-phone",
		Label:      "Contact Phone #",
		Comment:    "phone-number required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidatePhoneNumber},
	}
	contactFaxDef = &xscmsg.FieldDef{
		Tag:        "23f.",
		Annotation: "contact-fax",
		Label:      "Contact Fax #",
		Comment:    "phone-number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	otherContactDef = &xscmsg.FieldDef{
		Tag:        "24.",
		Annotation: "other-contact",
		Label:      "Other Phone, Fax, Cell Phone, Radio",
	}
	incidentDef = &xscmsg.FieldDef{
		Tag:        "25.",
		Annotation: "incident",
		Label:      "Incident Name",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	incidentDateDef = &xscmsg.FieldDef{
		Tag:        "25d.",
		Annotation: "incident-date",
		Label:      "Incident Date",
		Comment:    "date required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateDate},
	}
	statusDef = &xscmsg.FieldDef{
		Tag:        "35.",
		Annotation: "status",
		Label:      "Facility Status",
		Comment:    "required-for-complete: Fully Functional, Limited Services, Impaired/Closed",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Fully Functional", "Limited Services", "Impaired/Closed"},
	}
	attachOrgChartDef = &xscmsg.FieldDef{
		Tag:        "26a.",
		Annotation: "attach-org-chart",
		Label:      "NHICS/ICS Organization Chart Attached",
		Comment:    "Yes, No",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	attachRrDef = &xscmsg.FieldDef{
		Tag:        "26b.",
		Annotation: "attach-RR",
		Label:      "DEOC-9A Resource Request Forms Attached",
		Comment:    "required-for-complete: Yes, No",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	attachStatusDef = &xscmsg.FieldDef{
		Tag:        "26c.",
		Annotation: "attach-status",
		Label:      "NHICS/ICS Status Report Form - Standard Attached",
		Comment:    "Yes, No",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	attachActionPlanDef = &xscmsg.FieldDef{
		Tag:        "26d.",
		Annotation: "attach-action-plan",
		Label:      "NHICS/ICS Incident Action Plan Attached",
		Comment:    "Yes, No",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	eocPhoneDef = &xscmsg.FieldDef{
		Tag:        "27p.",
		Annotation: "eoc-phone",
		Label:      "Facility EOC Main Contact Number",
		Comment:    "phone-number required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidatePhoneNumber},
	}
	attachDirectoryDef = &xscmsg.FieldDef{
		Tag:        "26e.",
		Annotation: "attach-directory",
		Label:      "Phone/Communications Directory Attached",
		Comment:    "Yes, No",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	eocFaxDef = &xscmsg.FieldDef{
		Tag:        "27f.",
		Annotation: "eoc-fax",
		Label:      "Facility EOC Main Contact Fax",
		Comment:    "phone-number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	liaisonDef = &xscmsg.FieldDef{
		Tag:        "28.",
		Annotation: "liaison",
		Label:      "Facility Liaison Officer Name: Liaison to Public Health/Medical Health Branch",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	summaryDef = &xscmsg.FieldDef{
		Tag:        "34.",
		Annotation: "summary",
		Label:      "General Summary of Situation/Conditions",
	}
	liaisonPhoneDef = &xscmsg.FieldDef{
		Tag:        "28p.",
		Annotation: "liaison-phone",
		Label:      "Facility Liaison Contact Number",
		Comment:    "phone-number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	infoOfficerDef = &xscmsg.FieldDef{
		Tag:        "29.",
		Annotation: "info-officer",
		Label:      "Facility Information Officer Name",
	}
	infoOfficerPhoneDef = &xscmsg.FieldDef{
		Tag:        "29p.",
		Annotation: "info-officer-phone",
		Label:      "Facility Information Officer Contact Number",
		Comment:    "phone-number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	infoOfficerEmailDef = &xscmsg.FieldDef{
		Tag:        "29e.",
		Annotation: "info-officer-email",
		Label:      "Facility Information Officer Contact Email",
	}
	eocClosedContactDef = &xscmsg.FieldDef{
		Tag:        "30.",
		Annotation: "eoc-closed-contact",
		Label:      "If Facility EOC is Not Activated, Who Should be Contacted for Questions/Requests",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	eocPhone2Def = &xscmsg.FieldDef{
		Tag:        "30p.",
		Annotation: "eoc-phone",
		Label:      "Facility Contact Number",
		Comment:    "phone-number required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete, xscform.ValidatePhoneNumber},
	}
	skilledNursingBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "40a.",
		Annotation: "skilled-nursing-beds-staffed-m",
		Label:      "Skilled Nursing Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "40b.",
		Annotation: "skilled-nursing-beds-staffed-f",
		Label:      "Skilled Nursing Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "40c.",
		Annotation: "skilled-nursing-beds-vacant-m",
		Label:      "Skilled Nursing Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "40d.",
		Annotation: "skilled-nursing-beds-vacant-f",
		Label:      "Skilled Nursing Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "40e.",
		Annotation: "skilled-nursing-beds-surge",
		Label:      "Skilled Nursing Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	eocEmailDef = &xscmsg.FieldDef{
		Tag:        "30e.",
		Annotation: "eoc-email",
		Label:      "Facility Contact Email",
		Comment:    "required-for-complete",
		Validators: []xscmsg.Validator{requiredForComplete},
	}
	assistedLivingBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "41a.",
		Annotation: "assisted-living-beds-staffed-m",
		Label:      "Assisted Living Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "41b.",
		Annotation: "assisted-living-beds-staffed-f",
		Label:      "Assisted Living Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "41c.",
		Annotation: "assisted-living-beds-vacant-m",
		Label:      "Assisted Living Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "41d.",
		Annotation: "assisted-living-beds-vacant-f",
		Label:      "Assisted Living Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "41e.",
		Annotation: "assisted-living-beds-surge",
		Label:      "Assisted Living Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "42a.",
		Annotation: "sub-acute-beds-staffed-m",
		Label:      "Sub-Acute Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "42b.",
		Annotation: "sub-acute-beds-staffed-f",
		Label:      "Sub-Acute Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "42c.",
		Annotation: "sub-acute-beds-vacant-m",
		Label:      "Sub-Acute Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "42d.",
		Annotation: "sub-acute-beds-vacant-f",
		Label:      "Sub-Acute Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "42e.",
		Annotation: "sub-acute-beds-surge",
		Label:      "Sub-Acute Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	toEvacuateDef = &xscmsg.FieldDef{
		Tag:        "31a.",
		Annotation: "to-evacuate",
		Label:      "Facility Patients to Evacuate",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "43a.",
		Annotation: "alzheimers-beds-staffed-m",
		Label:      "Alzheimers/Dementia Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "43b.",
		Annotation: "alzheimers-beds-staffed-f",
		Label:      "Alzheimers/Dementia Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "43c.",
		Annotation: "alzheimers-beds-vacant-m",
		Label:      "Alzheimers/Dementia Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "43d.",
		Annotation: "alzheimers-beds-vacant-f",
		Label:      "Alzheimers/Dementia Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "43e.",
		Annotation: "alzheimers-beds-surge",
		Label:      "Alzheimers/Dementia Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	injuredDef = &xscmsg.FieldDef{
		Tag:        "31b.",
		Annotation: "injured",
		Label:      "Facility Patients Injured - Minor",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "44a.",
		Annotation: "ped-sub-acute-beds-staffed-m",
		Label:      "Pediatric-Sub Acute Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "44b.",
		Annotation: "ped-sub-acute-beds-staffed-f",
		Label:      "Pediatric-Sub Acute Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "44c.",
		Annotation: "ped-sub-acute-beds-vacant-m",
		Label:      "Pediatric-Sub Acute Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "44d.",
		Annotation: "ped-sub-acute-beds-vacant-f",
		Label:      "Pediatric-Sub Acute Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "44e.",
		Annotation: "ped-sub-acute-beds-surge",
		Label:      "Pediatric-Sub Acute Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	transferredDef = &xscmsg.FieldDef{
		Tag:        "31c.",
		Annotation: "transferred",
		Label:      "Facility Patients Transferred Out of County",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "45a.",
		Annotation: "psychiatric-beds-staffed-m",
		Label:      "Psychiatric Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "45b.",
		Annotation: "psychiatric-beds-staffed-f",
		Label:      "Psychiatric Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "45c.",
		Annotation: "psychiatric-beds-vacant-m",
		Label:      "Psychiatric Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "45d.",
		Annotation: "psychiatric-beds-vacant-f",
		Label:      "Psychiatric Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "45e.",
		Annotation: "psychiatric-beds-surge",
		Label:      "Psychiatric Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareDef = &xscmsg.FieldDef{
		Tag:        "33.",
		Annotation: "other-care",
		Label:      "Other Facility Patient Care Information",
	}
	bedResourceDef = &xscmsg.FieldDef{
		Tag:        "46.",
		Annotation: "bed-resource",
		Label:      "Other Care",
	}
	otherCareBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "46a.",
		Annotation: "other-care-beds-staffed-m",
		Label:      "Other Care Beds Staffed M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "46b.",
		Annotation: "other-care-beds-staffed-f",
		Label:      "Other Care Beds Staffed F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "46c.",
		Annotation: "other-care-beds-vacant-m",
		Label:      "Other Care Beds Vacant M",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "46d.",
		Annotation: "other-care-beds-vacant-f",
		Label:      "Other Care Beds Vacant F",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "46e.",
		Annotation: "other-care-beds-surge",
		Label:      "Other Care Beds Surge",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisChairsDef = &xscmsg.FieldDef{
		Tag:        "50a.",
		Annotation: "dialysis-chairs",
		Label:      "Dialysis Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "50b.",
		Annotation: "dialysis-vacant-chairs",
		Label:      "Dialysis Vacant Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "50c.",
		Annotation: "dialysis-front-staff",
		Label:      "Dialysis Front Desk Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "50d.",
		Annotation: "dialysis-support-staff",
		Label:      "Dialysis Medical Support Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisProvidersDef = &xscmsg.FieldDef{
		Tag:        "50e.",
		Annotation: "dialysis-providers",
		Label:      "Dialysis Provider Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalChairsDef = &xscmsg.FieldDef{
		Tag:        "51a.",
		Annotation: "surgical-chairs",
		Label:      "Surgical Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "51b.",
		Annotation: "surgical-vacant-chairs",
		Label:      "Surgical Vacant Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "51c.",
		Annotation: "surgical-front-staff",
		Label:      "Surgical Front Desk Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "51d.",
		Annotation: "surgical-support-staff",
		Label:      "Surgical Medical Support Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalProvidersDef = &xscmsg.FieldDef{
		Tag:        "51e.",
		Annotation: "surgical-providers",
		Label:      "Surgical Provider Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicChairsDef = &xscmsg.FieldDef{
		Tag:        "52a.",
		Annotation: "clinic-chairs",
		Label:      "Clinic Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "52b.",
		Annotation: "clinic-vacant-chairs",
		Label:      "Clinic Vacant Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "52c.",
		Annotation: "clinic-front-staff",
		Label:      "Clinic Front Desk Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "52d.",
		Annotation: "clinic-support-staff",
		Label:      "Clinic Medical Support Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicProvidersDef = &xscmsg.FieldDef{
		Tag:        "52e.",
		Annotation: "clinic-providers",
		Label:      "Clinic Provider Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthChairsDef = &xscmsg.FieldDef{
		Tag:        "53a.",
		Annotation: "home-health-chairs",
		Label:      "Home Health Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "53b.",
		Annotation: "home-health-vacant-chairs",
		Label:      "Home Health Vacant Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "53c.",
		Annotation: "home-health-front-staff",
		Label:      "Home Health Front Desk Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "53d.",
		Annotation: "home-health-support-staff",
		Label:      "Home Health Medical Support Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthProvidersDef = &xscmsg.FieldDef{
		Tag:        "53e.",
		Annotation: "home-health-providers",
		Label:      "Home Health Provider Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrChairsDef = &xscmsg.FieldDef{
		Tag:        "54a.",
		Annotation: "adult-day-ctr-chairs",
		Label:      "Adult Day Ctr Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "54b.",
		Annotation: "adult-day-ctr-vacant-chairs",
		Label:      "Adult Day Ctr Vacant Chairs/Room",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "54c.",
		Annotation: "adult-day-ctr-front-staff",
		Label:      "Adult Day Ctr Front Desk Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "54d.",
		Annotation: "adult-day-ctr-support-staff",
		Label:      "Adult Day Ctr Medical Support Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrProvidersDef = &xscmsg.FieldDef{
		Tag:        "54e.",
		Annotation: "adult-day-ctr-providers",
		Label:      "Adult Day Ctr Provider Staff",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
)

func requiredForComplete(f *xscmsg.Field, m *xscmsg.Message, _ bool) string {
	if m.Field("19.").Value != "Complete" {
		return ""
	}
	if f.Value == "" {
		return fmt.Sprintf("The %q field must have a value when the \"19.\" field is set to \"Complete\".", f.Def.Tag)
	}
	return ""
}
