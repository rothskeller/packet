package ahfacstat

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/xscform"
)

// Tag identifies allied health facility status forms.
const Tag = "AHFacStat"

func init() {
	xscmsg.RegisterCreate(formtype, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef.DefaultValue = "ROUTINE"
	toICSPositionDef.DefaultValue = "EMS Unit"
	toICSPositionDef.Choices = []string{"EMS Unit", "Medical Health Branch", "Operations Section", "Public Health Unit"}
	toLocationDef.DefaultValue = "MHJOC"
	toLocationDef.Choices = []string{"County EOC", "MHJOC"}
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
	Name:        "Allied Health Facility status report",
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
		Key:        xscmsg.FComplete,
		Label:      "Report Type",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Update", "Complete"},
		Flags:      xscmsg.Required,
	}
	facilityDef = &xscmsg.FieldDef{
		Tag:   "20.",
		Label: "Facility Name",
		Key:   xscmsg.FSubject,
		Flags: xscmsg.Required,
	}
	facilityTypeDef = &xscmsg.FieldDef{
		Tag:   "21.",
		Label: "Facility Type",
		Flags: xscmsg.RequiredForComplete,
	}
	dateDef = &xscmsg.FieldDef{
		Tag:         "22d.",
		Label:       "Date",
		Comment:     "MM/DD/YYYY",
		Validators:  []xscmsg.Validator{xscform.ValidateDate},
		DefaultFunc: xscform.DefaultDate,
		Flags:       xscmsg.Required,
	}
	timeDef = &xscmsg.FieldDef{
		Tag:        "22t.",
		Label:      "Time",
		Comment:    "HH:MM",
		Validators: []xscmsg.Validator{xscform.ValidateTime},
		Flags:      xscmsg.Required,
	}
	contactDef = &xscmsg.FieldDef{
		Tag:   "23.",
		Label: "Contact Name",
		Flags: xscmsg.RequiredForComplete,
	}
	contactPhoneDef = &xscmsg.FieldDef{
		Tag:        "23p.",
		Label:      "Contact Phone #",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	contactFaxDef = &xscmsg.FieldDef{
		Tag:        "23f.",
		Label:      "Contact Fax #",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	otherContactDef = &xscmsg.FieldDef{
		Tag:   "24.",
		Label: "Other Phone, Fax, Cell Phone, Radio",
	}
	incidentDef = &xscmsg.FieldDef{
		Tag:   "25.",
		Label: "Incident Name",
		Flags: xscmsg.RequiredForComplete,
	}
	incidentDateDef = &xscmsg.FieldDef{
		Tag:        "25d.",
		Label:      "Incident Date",
		Comment:    "MM/DD/YYYY",
		Validators: []xscmsg.Validator{xscform.ValidateDate},
		Flags:      xscmsg.RequiredForComplete,
	}
	statusDef = &xscmsg.FieldDef{
		Tag:        "35.",
		Label:      "Facility Status",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Fully Functional", "Limited Services", "Impaired/Closed"},
		Flags:      xscmsg.RequiredForComplete,
	}
	attachOrgChartDef = &xscmsg.FieldDef{
		Tag:        "26a.",
		Label:      "NHICS/ICS Organization Chart Attached",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	attachRrDef = &xscmsg.FieldDef{
		Tag:        "26b.",
		Label:      "DEOC-9A Resource Request Forms Attached",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
		Flags:      xscmsg.RequiredForComplete,
	}
	attachStatusDef = &xscmsg.FieldDef{
		Tag:        "26c.",
		Label:      "NHICS/ICS Status Report Form Attached",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	attachActionPlanDef = &xscmsg.FieldDef{
		Tag:        "26d.",
		Label:      "NHICS/ICS Incident Action Plan Attached",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	eocPhoneDef = &xscmsg.FieldDef{
		Tag:        "27p.",
		Label:      "Facility EOC Main Contact Number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	attachDirectoryDef = &xscmsg.FieldDef{
		Tag:        "26e.",
		Label:      "Phone/Communications Directory Attached",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Yes", "No"},
	}
	eocFaxDef = &xscmsg.FieldDef{
		Tag:        "27f.",
		Label:      "Facility EOC Main Contact Fax",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	liaisonDef = &xscmsg.FieldDef{
		Tag:   "28.",
		Label: "Facility Liaison Officer Name",
		Flags: xscmsg.RequiredForComplete,
	}
	summaryDef = &xscmsg.FieldDef{
		Tag:   "34.",
		Label: "General Summary of Situation/Conditions",
		Key:   xscmsg.FBody,
		Flags: xscmsg.Multiline,
	}
	liaisonPhoneDef = &xscmsg.FieldDef{
		Tag:        "28p.",
		Label:      "Facility Liaison Contact Number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	infoOfficerDef = &xscmsg.FieldDef{
		Tag:   "29.",
		Label: "Facility Information Officer Name",
	}
	infoOfficerPhoneDef = &xscmsg.FieldDef{
		Tag:        "29p.",
		Label:      "Facility Information Officer Number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
	}
	infoOfficerEmailDef = &xscmsg.FieldDef{
		Tag:   "29e.",
		Label: "Facility Information Officer Email",
	}
	eocClosedContactDef = &xscmsg.FieldDef{
		Tag:   "30.",
		Label: "If EOC is Not Activated, Who to Contact",
		Flags: xscmsg.RequiredForComplete,
	}
	eocPhone2Def = &xscmsg.FieldDef{
		Tag:        "30p.",
		Label:      "Facility Contact Number",
		Validators: []xscmsg.Validator{xscform.ValidatePhoneNumber},
		Flags:      xscmsg.RequiredForComplete,
	}
	skilledNursingBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "40a.",
		Label:      "Skilled Nursing Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "40b.",
		Label:      "Skilled Nursing Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "40c.",
		Label:      "Skilled Nursing Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "40d.",
		Label:      "Skilled Nursing Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	skilledNursingBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "40e.",
		Label:      "Skilled Nursing Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	eocEmailDef = &xscmsg.FieldDef{
		Tag:   "30e.",
		Label: "Facility Contact Email",
		Flags: xscmsg.RequiredForComplete,
	}
	assistedLivingBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "41a.",
		Label:      "Assisted Living Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "41b.",
		Label:      "Assisted Living Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "41c.",
		Label:      "Assisted Living Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "41d.",
		Label:      "Assisted Living Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	assistedLivingBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "41e.",
		Label:      "Assisted Living Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "42a.",
		Label:      "Sub-Acute Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "42b.",
		Label:      "Sub-Acute Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "42c.",
		Label:      "Sub-Acute Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "42d.",
		Label:      "Sub-Acute Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	subAcuteBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "42e.",
		Label:      "Sub-Acute Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	toEvacuateDef = &xscmsg.FieldDef{
		Tag:        "31a.",
		Label:      "Facility Patients to Evacuate",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "43a.",
		Label:      "Alzheimers/Dementia Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "43b.",
		Label:      "Alzheimers/Dementia Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "43c.",
		Label:      "Alzheimers/Dementia Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "43d.",
		Label:      "Alzheimers/Dementia Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	alzheimersBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "43e.",
		Label:      "Alzheimers/Dementia Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	injuredDef = &xscmsg.FieldDef{
		Tag:        "31b.",
		Label:      "Facility Patients Injured - Minor",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "44a.",
		Label:      "Pediatric-Sub Acute Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "44b.",
		Label:      "Pediatric-Sub Acute Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "44c.",
		Label:      "Pediatric-Sub Acute Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "44d.",
		Label:      "Pediatric-Sub Acute Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	pedSubAcuteBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "44e.",
		Label:      "Pediatric-Sub Acute Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	transferredDef = &xscmsg.FieldDef{
		Tag:        "31c.",
		Label:      "Patients Transferred Out of County",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "45a.",
		Label:      "Psychiatric Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "45b.",
		Label:      "Psychiatric Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "45c.",
		Label:      "Psychiatric Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "45d.",
		Label:      "Psychiatric Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	psychiatricBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "45e.",
		Label:      "Psychiatric Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareDef = &xscmsg.FieldDef{
		Tag:   "33.",
		Label: "Other Facility Patient Care Information",
	}
	bedResourceDef = &xscmsg.FieldDef{
		Tag:   "46.",
		Label: "Other Care",
	}
	otherCareBedsStaffedMDef = &xscmsg.FieldDef{
		Tag:        "46a.",
		Label:      "Other Care Beds Staffed M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsStaffedFDef = &xscmsg.FieldDef{
		Tag:        "46b.",
		Label:      "Other Care Beds Staffed F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsVacantMDef = &xscmsg.FieldDef{
		Tag:        "46c.",
		Label:      "Other Care Beds Vacant M",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsVacantFDef = &xscmsg.FieldDef{
		Tag:        "46d.",
		Label:      "Other Care Beds Vacant F",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	otherCareBedsSurgeDef = &xscmsg.FieldDef{
		Tag:        "46e.",
		Label:      "Other Care Beds Surge",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisChairsDef = &xscmsg.FieldDef{
		Tag:        "50a.",
		Label:      "Dialysis Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "50b.",
		Label:      "Dialysis Vacant Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "50c.",
		Label:      "Dialysis Front Desk Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "50d.",
		Label:      "Dialysis Medical Support Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	dialysisProvidersDef = &xscmsg.FieldDef{
		Tag:        "50e.",
		Label:      "Dialysis Provider Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalChairsDef = &xscmsg.FieldDef{
		Tag:        "51a.",
		Label:      "Surgical Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "51b.",
		Label:      "Surgical Vacant Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "51c.",
		Label:      "Surgical Front Desk Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "51d.",
		Label:      "Surgical Medical Support Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	surgicalProvidersDef = &xscmsg.FieldDef{
		Tag:        "51e.",
		Label:      "Surgical Provider Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicChairsDef = &xscmsg.FieldDef{
		Tag:        "52a.",
		Label:      "Clinic Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "52b.",
		Label:      "Clinic Vacant Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "52c.",
		Label:      "Clinic Front Desk Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "52d.",
		Label:      "Clinic Medical Support Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	clinicProvidersDef = &xscmsg.FieldDef{
		Tag:        "52e.",
		Label:      "Clinic Provider Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthChairsDef = &xscmsg.FieldDef{
		Tag:        "53a.",
		Label:      "Home Health Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "53b.",
		Label:      "Home Health Vacant Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "53c.",
		Label:      "Home Health Front Desk Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "53d.",
		Label:      "Home Health Medical Support Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	homeHealthProvidersDef = &xscmsg.FieldDef{
		Tag:        "53e.",
		Label:      "Home Health Provider Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrChairsDef = &xscmsg.FieldDef{
		Tag:        "54a.",
		Label:      "Adult Day Ctr Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrVacantChairsDef = &xscmsg.FieldDef{
		Tag:        "54b.",
		Label:      "Adult Day Ctr Vacant Chairs/Room",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrFrontStaffDef = &xscmsg.FieldDef{
		Tag:        "54c.",
		Label:      "Adult Day Ctr Front Desk Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrSupportStaffDef = &xscmsg.FieldDef{
		Tag:        "54d.",
		Label:      "Adult Day Ctr Medical Support Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	adultDayCtrProvidersDef = &xscmsg.FieldDef{
		Tag:        "54e.",
		Label:      "Adult Day Ctr Provider Staff",
		Comment:    "count",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
)
