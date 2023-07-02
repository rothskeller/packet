package ahfacstat

import "github.com/rothskeller/packet/message/common"

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (f *AHFacStat) Validate() (problems []string) {
	problems = append(problems, f.StdFields.Validate()...)
	switch f.ReportType {
	case "":
		problems = append(problems, `The "Report Type" field is required.`)
	case "Update", "Complete":
		break
	default:
		problems = append(problems, `The "Report Type" field does not contain a valid report type value.`)
	}
	if f.FacilityName == "" {
		problems = append(problems, `The "Facility Name" field is required.`)
	}
	if f.FacilityType == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Facility Type" field is required when the "Report Type" is "Complete".`)
	}
	if f.Date == "" {
		problems = append(problems, `The "Date" field is required.`)
	} else if !common.PIFODateRE.MatchString(f.Date) {
		problems = append(problems, `The "Date" field does not contain a valid date.`)
	}
	if f.Time == "" {
		problems = append(problems, `The "Time" field is required.`)
	} else if !common.PIFOTimeRE.MatchString(f.Time) {
		problems = append(problems, `The "Time" field does not contain a valid time.`)
	}
	if f.ContactName == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Contact Name" field is required when the "Report Type" is "Complete".`)
	}
	if f.ContactPhone == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Contact Phone #" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFOPhoneNumberRE.MatchString(f.ContactPhone) {
		problems = append(problems, `The "Contact Phone #" field does not contain a valid phone number.`)
	}
	if f.ContactFax != "" && !common.PIFOPhoneNumberRE.MatchString(f.ContactFax) {
		problems = append(problems, `The "Contact Fax #" field does not contain a valid phone number.`)
	}
	if f.IncidentName == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Incident Name" field is required when the "Report Type" is "Complete".`)
	}
	if f.IncidentDate == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Incident Date" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFODateRE.MatchString(f.IncidentDate) {
		problems = append(problems, `The "Incident Date" field does not contain a valid date.`)
	}
	switch f.FacilityStatus {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Facility Status" field is required when the "Report Type" is "Complete".`)
		}
	case "Fully Functional", "Limited Services", "Impaired/Closed":
		break
	default:
		problems = append(problems, `The "Facility Status" field does not contain a valid facility status.`)
	}
	if f.EOCPhone == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Facility EOC Main Contact Number" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFOPhoneNumberRE.MatchString(f.EOCPhone) {
		problems = append(problems, `The "Facility EOC Main Contact Number" field does not contain a valid phone number.`)
	}
	if f.EOCFax != "" && !common.PIFOPhoneNumberRE.MatchString(f.EOCFax) {
		problems = append(problems, `The "Facility EOC Main Contact Fax" field does not contain a valid phone number.`)
	}
	if f.LiaisonName == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Facility Liaison Officer Name" field is required when the "Report Type" is "Complete".`)
	}
	if f.LiaisonPhone != "" && !common.PIFOPhoneNumberRE.MatchString(f.LiaisonPhone) {
		problems = append(problems, `The "Facility Liaison Contact Number" field does not contain a valid phone number.`)
	}
	if f.InfoOfficerPhone != "" && !common.PIFOPhoneNumberRE.MatchString(f.InfoOfficerPhone) {
		problems = append(problems, `The "Facility Information Officer Contact Number" field does not contain a valid phone number.`)
	}
	if f.ClosedContactName == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "If Facility EOC is Not Activated, Who Should be Contacted for Questions/Requests" field is required when the "Report Type" is "Complete".`)
	}
	if f.ClosedContactPhone == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Facility Contact Number" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFOPhoneNumberRE.MatchString(f.ClosedContactPhone) {
		problems = append(problems, `The "Facility Contact Number" field does not contain a valid phone number.`)
	}
	if f.ClosedContactEmail == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Facility Contact Email" field is required when the "Report Type" is "Complete".`)
	}
	if f.PatientsToEvacuate != "" && !common.PIFOCardinalNumberRE.MatchString(f.PatientsToEvacuate) {
		problems = append(problems, `The "Facility Patients to Evacuate" field does not contain a valid number.`)
	}
	if f.PatientsInjuredMinor != "" && !common.PIFOCardinalNumberRE.MatchString(f.PatientsInjuredMinor) {
		problems = append(problems, `The "Facility Patients Injured - Minor" field does not contain a valid number.`)
	}
	if f.PatientsTransferred != "" && !common.PIFOCardinalNumberRE.MatchString(f.PatientsTransferred) {
		problems = append(problems, `The "Facility Patients Transfered Out of County" field does not contain a valid number.`)
	}
	switch f.AttachOrgChart {
	case "", "Yes", "No":
		break
	default:
		problems = append(problems, `The "NHICS/ICS Organization Chart" field does not contain a valid value.`)
	}
	switch f.AttachRR {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "DEOC-9A Resource Request Forms" field is required when "Report Type" is "Complete".`)
		}
	case "Yes", "No":
		break
	default:
		problems = append(problems, `The "DEOC-9A Resource Request Forms" field does not contain a valid value.`)
	}
	switch f.AttachStatus {
	case "", "Yes", "No":
		break
	default:
		problems = append(problems, `The "NHICS/ICS Status Report Form - Standard" field does not contain a valid value.`)
	}
	switch f.AttachActionPlan {
	case "", "Yes", "No":
		break
	default:
		problems = append(problems, `The "NHICS/ICS Incident Action Plan" field does not contain a valid value.`)
	}
	switch f.AttachDirectory {
	case "", "Yes", "No":
		break
	default:
		problems = append(problems, `The "Phone/Communications Directory" field does not contain a valid value.`)
	}
	if f.SkilledNursingBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.SkilledNursingBedsStaffedM) {
		problems = append(problems, `The "Skilled Nursing Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.SkilledNursingBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.SkilledNursingBedsStaffedF) {
		problems = append(problems, `The "Skilled Nursing Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.SkilledNursingBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.SkilledNursingBedsVacantM) {
		problems = append(problems, `The "Skilled Nursing Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.SkilledNursingBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.SkilledNursingBedsVacantF) {
		problems = append(problems, `The "Skilled Nursing Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.SkilledNursingBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.SkilledNursingBedsSurge) {
		problems = append(problems, `The "Skilled Nursing Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.SkilledNursingBedsStaffedM, f.SkilledNursingBedsStaffedF, f.SkilledNursingBedsVacantM, f.SkilledNursingBedsVacantF, f.SkilledNursingBedsSurge) {
		problems = append(problems, `Either all of the "Skilled Nursing" fields must have values, or none of them.`)
	}
	if f.AssistedLivingBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.AssistedLivingBedsStaffedM) {
		problems = append(problems, `The "Assisted Living Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.AssistedLivingBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.AssistedLivingBedsStaffedF) {
		problems = append(problems, `The "Assisted Living Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.AssistedLivingBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.AssistedLivingBedsVacantM) {
		problems = append(problems, `The "Assisted Living Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.AssistedLivingBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.AssistedLivingBedsVacantF) {
		problems = append(problems, `The "Assisted Living Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.AssistedLivingBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.AssistedLivingBedsSurge) {
		problems = append(problems, `The "Assisted Living Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.AssistedLivingBedsStaffedM, f.AssistedLivingBedsStaffedF, f.AssistedLivingBedsVacantM, f.AssistedLivingBedsVacantF, f.AssistedLivingBedsSurge) {
		problems = append(problems, `Either all of the "Assisted Living" fields must have values, or none of them.`)
	}
	if f.SubAcuteBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.SubAcuteBedsStaffedM) {
		problems = append(problems, `The "Sub-Acute Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.SubAcuteBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.SubAcuteBedsStaffedF) {
		problems = append(problems, `The "Sub-Acute Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.SubAcuteBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.SubAcuteBedsVacantM) {
		problems = append(problems, `The "Sub-Acute Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.SubAcuteBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.SubAcuteBedsVacantF) {
		problems = append(problems, `The "Sub-Acute Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.SubAcuteBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.SubAcuteBedsSurge) {
		problems = append(problems, `The "Sub-Acute Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.SubAcuteBedsStaffedM, f.SubAcuteBedsStaffedF, f.SubAcuteBedsVacantM, f.SubAcuteBedsVacantF, f.SubAcuteBedsSurge) {
		problems = append(problems, `Either all of the "Sub-Acute" fields must have values, or none of them.`)
	}
	if f.AlzheimersBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.AlzheimersBedsStaffedM) {
		problems = append(problems, `The "Alzeimers/Dimentia Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.AlzheimersBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.AlzheimersBedsStaffedF) {
		problems = append(problems, `The "Alzeimers/Dimentia Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.AlzheimersBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.AlzheimersBedsVacantM) {
		problems = append(problems, `The "Alzeimers/Dimentia Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.AlzheimersBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.AlzheimersBedsVacantF) {
		problems = append(problems, `The "Alzeimers/Dimentia Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.AlzheimersBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.AlzheimersBedsSurge) {
		problems = append(problems, `The "Alzeimers/Dimentia Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.AlzheimersBedsStaffedM, f.AlzheimersBedsStaffedF, f.AlzheimersBedsVacantM, f.AlzheimersBedsVacantF, f.AlzheimersBedsSurge) {
		problems = append(problems, `Either all of the "Alzeimers/Dimentia" fields must have values, or none of them.`)
	}
	if f.PedSubAcuteBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.PedSubAcuteBedsStaffedM) {
		problems = append(problems, `The "Pediatric-Sub Acute Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.PedSubAcuteBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.PedSubAcuteBedsStaffedF) {
		problems = append(problems, `The "Pediatric-Sub Acute Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.PedSubAcuteBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.PedSubAcuteBedsVacantM) {
		problems = append(problems, `The "Pediatric-Sub Acute Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.PedSubAcuteBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.PedSubAcuteBedsVacantF) {
		problems = append(problems, `The "Pediatric-Sub Acute Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.PedSubAcuteBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.PedSubAcuteBedsSurge) {
		problems = append(problems, `The "Pediatric-Sub Acute Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.PedSubAcuteBedsStaffedM, f.PedSubAcuteBedsStaffedF, f.PedSubAcuteBedsVacantM, f.PedSubAcuteBedsVacantF, f.PedSubAcuteBedsSurge) {
		problems = append(problems, `Either all of the "Pediatric-Sub Acute" fields must have values, or none of them.`)
	}
	if f.PsychiatricBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.PsychiatricBedsStaffedM) {
		problems = append(problems, `The "Psychiatric Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.PsychiatricBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.PsychiatricBedsStaffedF) {
		problems = append(problems, `The "Psychiatric Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.PsychiatricBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.PsychiatricBedsVacantM) {
		problems = append(problems, `The "Psychiatric Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.PsychiatricBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.PsychiatricBedsVacantF) {
		problems = append(problems, `The "Psychiatric Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.PsychiatricBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.PsychiatricBedsSurge) {
		problems = append(problems, `The "Psychiatric Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.PsychiatricBedsStaffedM, f.PsychiatricBedsStaffedF, f.PsychiatricBedsVacantM, f.PsychiatricBedsVacantF, f.PsychiatricBedsSurge) {
		problems = append(problems, `Either all of the "Psychiatric" fields must have values, or none of them.`)
	}
	if f.OtherCareBedsStaffedM != "" && !common.PIFOCardinalNumberRE.MatchString(f.OtherCareBedsStaffedM) {
		problems = append(problems, `The "(Other Care) Staffed Bed-M" field does not contain a valid number.`)
	}
	if f.OtherCareBedsStaffedF != "" && !common.PIFOCardinalNumberRE.MatchString(f.OtherCareBedsStaffedF) {
		problems = append(problems, `The "(Other Care) Staffed Bed-F" field does not contain a valid number.`)
	}
	if f.OtherCareBedsVacantM != "" && !common.PIFOCardinalNumberRE.MatchString(f.OtherCareBedsVacantM) {
		problems = append(problems, `The "(Other Care) Vacant Beds-M" field does not contain a valid number.`)
	}
	if f.OtherCareBedsVacantF != "" && !common.PIFOCardinalNumberRE.MatchString(f.OtherCareBedsVacantF) {
		problems = append(problems, `The "(Other Care) Vacant Beds-F" field does not contain a valid number.`)
	}
	if f.OtherCareBedsSurge != "" && !common.PIFOCardinalNumberRE.MatchString(f.OtherCareBedsSurge) {
		problems = append(problems, `The "(Other Care) Surge #" field does not contain a valid number.`)
	}
	if !allOrNone(f.OtherCareBedsStaffedM, f.OtherCareBedsStaffedF, f.OtherCareBedsVacantM, f.OtherCareBedsVacantF, f.OtherCareBedsSurge) {
		problems = append(problems, `Either all of the "(Other Care)" fields must have values, or none of them.`)
	}
	if f.DialysisChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.DialysisChairs) {
		problems = append(problems, `The "Dialysis Chairs/Rooms" field does not contain a valid number.`)
	}
	if f.DialysisVacantChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.DialysisVacantChairs) {
		problems = append(problems, `The "Dialysis Vancant Chairs/Room" field does not contain a valid number.`)
	}
	if f.DialysisFrontStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.DialysisFrontStaff) {
		problems = append(problems, `The "Dialysis Front Desk Staff" field does not contain a valid number.`)
	}
	if f.DialysisSupportStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.DialysisSupportStaff) {
		problems = append(problems, `The "Dialysis Medical Support Staff" field does not contain a valid number.`)
	}
	if f.DialysisProviders != "" && !common.PIFOCardinalNumberRE.MatchString(f.DialysisProviders) {
		problems = append(problems, `The "Dialysis Provider Staff" field does not contain a valid number.`)
	}
	if !allOrNone(f.DialysisChairs, f.DialysisVacantChairs, f.DialysisFrontStaff, f.DialysisSupportStaff, f.DialysisProviders) {
		problems = append(problems, `Either all of the "Dialysis" fields must have values, or none of them.`)
	}
	if f.SurgicalChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.SurgicalChairs) {
		problems = append(problems, `The "Surgical Chairs/Rooms" field does not contain a valid number.`)
	}
	if f.SurgicalVacantChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.SurgicalVacantChairs) {
		problems = append(problems, `The "Surgical Vancant Chairs/Room" field does not contain a valid number.`)
	}
	if f.SurgicalFrontStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.SurgicalFrontStaff) {
		problems = append(problems, `The "Surgical Front Desk Staff" field does not contain a valid number.`)
	}
	if f.SurgicalSupportStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.SurgicalSupportStaff) {
		problems = append(problems, `The "Surgical Medical Support Staff" field does not contain a valid number.`)
	}
	if f.SurgicalProviders != "" && !common.PIFOCardinalNumberRE.MatchString(f.SurgicalProviders) {
		problems = append(problems, `The "Surgical Provider Staff" field does not contain a valid number.`)
	}
	if !allOrNone(f.SurgicalChairs, f.SurgicalVacantChairs, f.SurgicalFrontStaff, f.SurgicalSupportStaff, f.SurgicalProviders) {
		problems = append(problems, `Either all of the "Surgical" fields must have values, or none of them.`)
	}
	if f.ClinicChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.ClinicChairs) {
		problems = append(problems, `The "Clinic Chairs/Rooms" field does not contain a valid number.`)
	}
	if f.ClinicVacantChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.ClinicVacantChairs) {
		problems = append(problems, `The "Clinic Vancant Chairs/Room" field does not contain a valid number.`)
	}
	if f.ClinicFrontStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.ClinicFrontStaff) {
		problems = append(problems, `The "Clinic Front Desk Staff" field does not contain a valid number.`)
	}
	if f.ClinicSupportStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.ClinicSupportStaff) {
		problems = append(problems, `The "Clinic Medical Support Staff" field does not contain a valid number.`)
	}
	if f.ClinicProviders != "" && !common.PIFOCardinalNumberRE.MatchString(f.ClinicProviders) {
		problems = append(problems, `The "Clinic Provider Staff" field does not contain a valid number.`)
	}
	if !allOrNone(f.ClinicChairs, f.ClinicVacantChairs, f.ClinicFrontStaff, f.ClinicSupportStaff, f.ClinicProviders) {
		problems = append(problems, `Either all of the "Clinic" fields must have values, or none of them.`)
	}
	if f.HomeHealthChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.HomeHealthChairs) {
		problems = append(problems, `The "Homehealth Chairs/Rooms" field does not contain a valid number.`)
	}
	if f.HomeHealthVacantChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.HomeHealthVacantChairs) {
		problems = append(problems, `The "Homehealth Vancant Chairs/Room" field does not contain a valid number.`)
	}
	if f.HomeHealthFrontStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.HomeHealthFrontStaff) {
		problems = append(problems, `The "Homehealth Front Desk Staff" field does not contain a valid number.`)
	}
	if f.HomeHealthSupportStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.HomeHealthSupportStaff) {
		problems = append(problems, `The "Homehealth Medical Support Staff" field does not contain a valid number.`)
	}
	if f.HomeHealthProviders != "" && !common.PIFOCardinalNumberRE.MatchString(f.HomeHealthProviders) {
		problems = append(problems, `The "Homehealth Provider Staff" field does not contain a valid number.`)
	}
	if !allOrNone(f.HomeHealthChairs, f.HomeHealthVacantChairs, f.HomeHealthFrontStaff, f.HomeHealthSupportStaff, f.HomeHealthProviders) {
		problems = append(problems, `Either all of the "Homehealth" fields must have values, or none of them.`)
	}
	if f.AdultDayCtrChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.AdultDayCtrChairs) {
		problems = append(problems, `The "Adult Day Center Chairs/Rooms" field does not contain a valid number.`)
	}
	if f.AdultDayCtrVacantChairs != "" && !common.PIFOCardinalNumberRE.MatchString(f.AdultDayCtrVacantChairs) {
		problems = append(problems, `The "Adult Day Center Vancant Chairs/Room" field does not contain a valid number.`)
	}
	if f.AdultDayCtrFrontStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.AdultDayCtrFrontStaff) {
		problems = append(problems, `The "Adult Day Center Front Desk Staff" field does not contain a valid number.`)
	}
	if f.AdultDayCtrSupportStaff != "" && !common.PIFOCardinalNumberRE.MatchString(f.AdultDayCtrSupportStaff) {
		problems = append(problems, `The "Adult Day Center Medical Support Staff" field does not contain a valid number.`)
	}
	if f.AdultDayCtrProviders != "" && !common.PIFOCardinalNumberRE.MatchString(f.AdultDayCtrProviders) {
		problems = append(problems, `The "Adult Day Center Provider Staff" field does not contain a valid number.`)
	}
	if !allOrNone(f.AdultDayCtrChairs, f.AdultDayCtrVacantChairs, f.AdultDayCtrFrontStaff, f.AdultDayCtrSupportStaff, f.AdultDayCtrProviders) {
		problems = append(problems, `Either all of the "Adult Day Center" fields must have values, or none of them.`)
	}
	return problems
}

func allOrNone(ss ...string) bool {
	var seenSet, seenUnset bool
	for _, s := range ss {
		if s == "" {
			if seenSet {
				return false
			}
			seenUnset = true
		} else {
			if seenUnset {
				return false
			}
			seenSet = true
		}
	}
	return true
}
