package ahfacstat

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) (f *AHFacStat) {
	if idx := strings.Index(body, "form-allied-health-facility-status.html"); idx < 0 {
		return nil
	}
	form := common.DecodePIFO(body)
	if form == nil || form.HTMLIdent != "form-allied-health-facility-status.html" {
		return nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2", "2.3":
		break
	default:
		return nil
	}
	f = new(AHFacStat)
	f.PIFOVersion = form.PIFOVersion
	f.FormVersion = form.FormVersion
	f.StdFields.Decode(form.TaggedValues)
	f.ReportType = form.TaggedValues["19."]
	f.FacilityName = form.TaggedValues["20."]
	f.FacilityType = form.TaggedValues["21."]
	f.Date = form.TaggedValues["22d."]
	f.Time = form.TaggedValues["22t."]
	f.ContactName = form.TaggedValues["23."]
	f.ContactPhone = form.TaggedValues["23p."]
	f.ContactFax = form.TaggedValues["23f."]
	f.OtherContact = form.TaggedValues["24."]
	f.IncidentName = form.TaggedValues["25."]
	f.IncidentDate = form.TaggedValues["25d."]
	f.FacilityStatus = form.TaggedValues["35."]
	f.EOCPhone = form.TaggedValues["27p."]
	f.EOCFax = form.TaggedValues["27f."]
	f.LiaisonName = form.TaggedValues["28."]
	f.LiaisonPhone = form.TaggedValues["28p."]
	f.InfoOfficerName = form.TaggedValues["29."]
	f.InfoOfficerPhone = form.TaggedValues["29p."]
	f.InfoOfficerEmail = form.TaggedValues["29e."]
	f.ClosedContactName = form.TaggedValues["30."]
	f.ClosedContactPhone = form.TaggedValues["30p."]
	f.ClosedContactEmail = form.TaggedValues["30e."]
	f.PatientsToEvacuate = form.TaggedValues["31a."]
	f.PatientsInjuredMinor = form.TaggedValues["31b."]
	f.PatientsTransferred = form.TaggedValues["31c."]
	f.OtherPatientCare = form.TaggedValues["33."]
	f.AttachOrgChart = form.TaggedValues["26a."]
	f.AttachRR = form.TaggedValues["26b."]
	f.AttachStatus = form.TaggedValues["26c."]
	f.AttachActionPlan = form.TaggedValues["26d."]
	f.AttachDirectory = form.TaggedValues["26e."]
	f.Summary = form.TaggedValues["34."]
	f.SkilledNursingBedsStaffedM = form.TaggedValues["40a."]
	f.SkilledNursingBedsStaffedF = form.TaggedValues["40b."]
	f.SkilledNursingBedsVacantM = form.TaggedValues["40c."]
	f.SkilledNursingBedsVacantF = form.TaggedValues["40d."]
	f.SkilledNursingBedsSurge = form.TaggedValues["40e."]
	f.AssistedLivingBedsStaffedM = form.TaggedValues["41a."]
	f.AssistedLivingBedsStaffedF = form.TaggedValues["41b."]
	f.AssistedLivingBedsVacantM = form.TaggedValues["41c."]
	f.AssistedLivingBedsVacantF = form.TaggedValues["41d."]
	f.AssistedLivingBedsSurge = form.TaggedValues["41e."]
	f.SubAcuteBedsStaffedM = form.TaggedValues["42a."]
	f.SubAcuteBedsStaffedF = form.TaggedValues["42b."]
	f.SubAcuteBedsVacantM = form.TaggedValues["42c."]
	f.SubAcuteBedsVacantF = form.TaggedValues["42d."]
	f.SubAcuteBedsSurge = form.TaggedValues["42e."]
	f.AlzheimersBedsStaffedM = form.TaggedValues["43a."]
	f.AlzheimersBedsStaffedF = form.TaggedValues["43b."]
	f.AlzheimersBedsVacantM = form.TaggedValues["43c."]
	f.AlzheimersBedsVacantF = form.TaggedValues["43d."]
	f.AlzheimersBedsSurge = form.TaggedValues["43e."]
	f.PedSubAcuteBedsStaffedM = form.TaggedValues["44a."]
	f.PedSubAcuteBedsStaffedF = form.TaggedValues["44b."]
	f.PedSubAcuteBedsVacantM = form.TaggedValues["44c."]
	f.PedSubAcuteBedsVacantF = form.TaggedValues["44d."]
	f.PedSubAcuteBedsSurge = form.TaggedValues["44e."]
	f.PsychiatricBedsStaffedM = form.TaggedValues["45a."]
	f.PsychiatricBedsStaffedF = form.TaggedValues["45b."]
	f.PsychiatricBedsVacantM = form.TaggedValues["45c."]
	f.PsychiatricBedsVacantF = form.TaggedValues["45d."]
	f.PsychiatricBedsSurge = form.TaggedValues["45e."]
	f.OtherCareBedsType = form.TaggedValues["46."]
	f.OtherCareBedsStaffedM = form.TaggedValues["46a."]
	f.OtherCareBedsStaffedF = form.TaggedValues["46b."]
	f.OtherCareBedsVacantM = form.TaggedValues["46c."]
	f.OtherCareBedsVacantF = form.TaggedValues["46d."]
	f.OtherCareBedsSurge = form.TaggedValues["46e."]
	f.DialysisChairs = form.TaggedValues["50a."]
	f.DialysisVacantChairs = form.TaggedValues["50b."]
	f.DialysisFrontStaff = form.TaggedValues["50c."]
	f.DialysisSupportStaff = form.TaggedValues["50d."]
	f.DialysisProviders = form.TaggedValues["50e."]
	f.SurgicalChairs = form.TaggedValues["51a."]
	f.SurgicalVacantChairs = form.TaggedValues["51b."]
	f.SurgicalFrontStaff = form.TaggedValues["51c."]
	f.SurgicalSupportStaff = form.TaggedValues["51d."]
	f.SurgicalProviders = form.TaggedValues["51e."]
	f.ClinicChairs = form.TaggedValues["52a."]
	f.ClinicVacantChairs = form.TaggedValues["52b."]
	f.ClinicFrontStaff = form.TaggedValues["52c."]
	f.ClinicSupportStaff = form.TaggedValues["52d."]
	f.ClinicProviders = form.TaggedValues["52e."]
	f.HomeHealthChairs = form.TaggedValues["53a."]
	f.HomeHealthVacantChairs = form.TaggedValues["53b."]
	f.HomeHealthFrontStaff = form.TaggedValues["53c."]
	f.HomeHealthSupportStaff = form.TaggedValues["53d."]
	f.HomeHealthProviders = form.TaggedValues["53e."]
	f.AdultDayCtrChairs = form.TaggedValues["54a."]
	f.AdultDayCtrVacantChairs = form.TaggedValues["54b."]
	f.AdultDayCtrFrontStaff = form.TaggedValues["54c."]
	f.AdultDayCtrSupportStaff = form.TaggedValues["54d."]
	f.AdultDayCtrProviders = form.TaggedValues["54e."]
	return f
}
