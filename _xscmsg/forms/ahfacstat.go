package xscmsg

import (
	"strings"

	"github.com/rothskeller/packet/xscmsg/forms/pifo"
	"github.com/rothskeller/packet/xscmsg/forms/xscsubj"
)

// AHFacStat form metadata:
const (
	AHFacStatTag     = "AHFacStat"
	AHFacStatHTML    = "form-allied-health-facility-status.html"
	AHFacStatVersion = "2.3"
)

// AHFacStat holds an allied health facility status form.
type AHFacStat struct {
	StdHeader
	ReportType                 string
	FacilityName               string
	FacilityType               string
	Date                       string
	Time                       string
	ContactName                string
	ContactPhone               string
	ContactFax                 string
	OtherContact               string
	IncidentName               string
	IncidentDate               string
	FacilityStatus             string
	EOCPhone                   string
	EOCFax                     string
	LiaisonName                string
	LiaisonPhone               string
	InfoOfficerName            string
	InfoOfficerPhone           string
	InfoOfficerEmail           string
	ClosedContactName          string
	ClosedContactPhone         string
	ClosedContactEmail         string
	PatientsToEvacuate         string
	PatientsInjuredMinor       string
	PatientsTransferred        string
	OtherPatientCare           string
	AttachOrgChart             string
	AttachRR                   string
	AttachStatus               string
	AttachActionPlan           string
	AttachDirectory            string
	Summary                    string
	SkilledNursingBedsStaffedM string
	SkilledNursingBedsStaffedF string
	SkilledNursingBedsVacantM  string
	SkilledNursingBedsVacantF  string
	SkilledNursingBedsSurge    string
	AssistedLivingBedsStaffedM string
	AssistedLivingBedsStaffedF string
	AssistedLivingBedsVacantM  string
	AssistedLivingBedsVacantF  string
	AssistedLivingBedsSurge    string
	SubAcuteBedsStaffedM       string
	SubAcuteBedsStaffedF       string
	SubAcuteBedsVacantM        string
	SubAcuteBedsVacantF        string
	SubAcuteBedsSurge          string
	AlzheimersBedsStaffedM     string
	AlzheimersBedsStaffedF     string
	AlzheimersBedsVacantM      string
	AlzheimersBedsVacantF      string
	AlzheimersBedsSurge        string
	PedSubAcuteBedsStaffedM    string
	PedSubAcuteBedsStaffedF    string
	PedSubAcuteBedsVacantM     string
	PedSubAcuteBedsVacantF     string
	PedSubAcuteBedsSurge       string
	PsychiatricBedsStaffedM    string
	PsychiatricBedsStaffedF    string
	PsychiatricBedsVacantM     string
	PsychiatricBedsVacantF     string
	PsychiatricBedsSurge       string
	OtherCareBedsType          string
	OtherCareBedsStaffedM      string
	OtherCareBedsStaffedF      string
	OtherCareBedsVacantM       string
	OtherCareBedsVacantF       string
	OtherCareBedsSurge         string
	DialysisChairs             string
	DialysisVacantChairs       string
	DialysisFrontStaff         string
	DialysisSupportStaff       string
	DialysisProviders          string
	SurgicalChairs             string
	SurgicalVacantChairs       string
	SurgicalFrontStaff         string
	SurgicalSupportStaff       string
	SurgicalProviders          string
	ClinicChairs               string
	ClinicVacantChairs         string
	ClinicFrontStaff           string
	ClinicSupportStaff         string
	ClinicProviders            string
	HomeHealthChairs           string
	HomeHealthVacantChairs     string
	HomeHealthFrontStaff       string
	HomeHealthSupportStaff     string
	HomeHealthProviders        string
	AdultDayCtrChairs          string
	AdultDayCtrVacantChairs    string
	AdultDayCtrFrontStaff      string
	AdultDayCtrSupportStaff    string
	AdultDayCtrProviders       string
	StdFooter
}

// DecodeAHFacStat decodes the supplied form if it is an AHFacStat form.  It
// returns the decoded form and strings describing any non-fatal decoding
// problems.  It returns nil, nil if the form is not an AHFacStat form or has an
// unknown version.
func DecodeAHFacStat(form *pifo.Form) (f *AHFacStat, problems []string) {
	if form.HTMLIdent != AHFacStatHTML {
		return nil, nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.2", "2.3":
		break
	default:
		return nil, nil
	}
	f = new(AHFacStat)
	f.FormVersion = form.FormVersion
	f.StdHeader.PullTags(form.TaggedValues)
	f.ReportType = PullTag(form.TaggedValues, "19.")
	f.FacilityName = PullTag(form.TaggedValues, "20.")
	f.FacilityType = PullTag(form.TaggedValues, "21.")
	f.Date = PullTag(form.TaggedValues, "22d.")
	f.Time = PullTag(form.TaggedValues, "22t.")
	f.ContactName = PullTag(form.TaggedValues, "23.")
	f.ContactPhone = PullTag(form.TaggedValues, "23p.")
	f.ContactFax = PullTag(form.TaggedValues, "23f.")
	f.OtherContact = PullTag(form.TaggedValues, "24.")
	f.IncidentName = PullTag(form.TaggedValues, "25.")
	f.IncidentDate = PullTag(form.TaggedValues, "25d.")
	f.FacilityStatus = PullTag(form.TaggedValues, "35.")
	f.EOCPhone = PullTag(form.TaggedValues, "27p.")
	f.EOCFax = PullTag(form.TaggedValues, "27f.")
	f.LiaisonName = PullTag(form.TaggedValues, "28.")
	f.LiaisonPhone = PullTag(form.TaggedValues, "28p.")
	f.InfoOfficerName = PullTag(form.TaggedValues, "29.")
	f.InfoOfficerPhone = PullTag(form.TaggedValues, "29p.")
	f.InfoOfficerEmail = PullTag(form.TaggedValues, "29e.")
	f.ClosedContactName = PullTag(form.TaggedValues, "30.")
	f.ClosedContactPhone = PullTag(form.TaggedValues, "30p.")
	f.ClosedContactEmail = PullTag(form.TaggedValues, "30e.")
	f.PatientsToEvacuate = PullTag(form.TaggedValues, "31a.")
	f.PatientsInjuredMinor = PullTag(form.TaggedValues, "31b.")
	f.PatientsTransferred = PullTag(form.TaggedValues, "31c.")
	f.OtherPatientCare = PullTag(form.TaggedValues, "33.")
	f.AttachOrgChart = PullTag(form.TaggedValues, "26a.")
	f.AttachRR = PullTag(form.TaggedValues, "26b.")
	f.AttachStatus = PullTag(form.TaggedValues, "26c.")
	f.AttachActionPlan = PullTag(form.TaggedValues, "26d.")
	f.AttachDirectory = PullTag(form.TaggedValues, "26e.")
	f.Summary = PullTag(form.TaggedValues, "34.")
	f.SkilledNursingBedsStaffedM = PullTag(form.TaggedValues, "40a.")
	f.SkilledNursingBedsStaffedF = PullTag(form.TaggedValues, "40b.")
	f.SkilledNursingBedsVacantM = PullTag(form.TaggedValues, "40c.")
	f.SkilledNursingBedsVacantF = PullTag(form.TaggedValues, "40d.")
	f.SkilledNursingBedsSurge = PullTag(form.TaggedValues, "40e.")
	f.AssistedLivingBedsStaffedM = PullTag(form.TaggedValues, "41a.")
	f.AssistedLivingBedsStaffedF = PullTag(form.TaggedValues, "41b.")
	f.AssistedLivingBedsVacantM = PullTag(form.TaggedValues, "41c.")
	f.AssistedLivingBedsVacantF = PullTag(form.TaggedValues, "41d.")
	f.AssistedLivingBedsSurge = PullTag(form.TaggedValues, "41e.")
	f.SubAcuteBedsStaffedM = PullTag(form.TaggedValues, "42a.")
	f.SubAcuteBedsStaffedF = PullTag(form.TaggedValues, "42b.")
	f.SubAcuteBedsVacantM = PullTag(form.TaggedValues, "42c.")
	f.SubAcuteBedsVacantF = PullTag(form.TaggedValues, "42d.")
	f.SubAcuteBedsSurge = PullTag(form.TaggedValues, "42e.")
	f.AlzheimersBedsStaffedM = PullTag(form.TaggedValues, "43a.")
	f.AlzheimersBedsStaffedF = PullTag(form.TaggedValues, "43b.")
	f.AlzheimersBedsVacantM = PullTag(form.TaggedValues, "43c.")
	f.AlzheimersBedsVacantF = PullTag(form.TaggedValues, "43d.")
	f.AlzheimersBedsSurge = PullTag(form.TaggedValues, "43e.")
	f.PedSubAcuteBedsStaffedM = PullTag(form.TaggedValues, "44a.")
	f.PedSubAcuteBedsStaffedF = PullTag(form.TaggedValues, "44b.")
	f.PedSubAcuteBedsVacantM = PullTag(form.TaggedValues, "44c.")
	f.PedSubAcuteBedsVacantF = PullTag(form.TaggedValues, "44d.")
	f.PedSubAcuteBedsSurge = PullTag(form.TaggedValues, "44e.")
	f.PsychiatricBedsStaffedM = PullTag(form.TaggedValues, "45a.")
	f.PsychiatricBedsStaffedF = PullTag(form.TaggedValues, "45b.")
	f.PsychiatricBedsVacantM = PullTag(form.TaggedValues, "45c.")
	f.PsychiatricBedsVacantF = PullTag(form.TaggedValues, "45d.")
	f.PsychiatricBedsSurge = PullTag(form.TaggedValues, "45e.")
	f.OtherCareBedsType = PullTag(form.TaggedValues, "46.")
	f.OtherCareBedsStaffedM = PullTag(form.TaggedValues, "46a.")
	f.OtherCareBedsStaffedF = PullTag(form.TaggedValues, "46b.")
	f.OtherCareBedsVacantM = PullTag(form.TaggedValues, "46c.")
	f.OtherCareBedsVacantF = PullTag(form.TaggedValues, "46d.")
	f.OtherCareBedsSurge = PullTag(form.TaggedValues, "46e.")
	f.DialysisChairs = PullTag(form.TaggedValues, "50a.")
	f.DialysisVacantChairs = PullTag(form.TaggedValues, "50b.")
	f.DialysisFrontStaff = PullTag(form.TaggedValues, "50c.")
	f.DialysisSupportStaff = PullTag(form.TaggedValues, "50d.")
	f.DialysisProviders = PullTag(form.TaggedValues, "50e.")
	f.SurgicalChairs = PullTag(form.TaggedValues, "51a.")
	f.SurgicalVacantChairs = PullTag(form.TaggedValues, "51b.")
	f.SurgicalFrontStaff = PullTag(form.TaggedValues, "51c.")
	f.SurgicalSupportStaff = PullTag(form.TaggedValues, "51d.")
	f.SurgicalProviders = PullTag(form.TaggedValues, "51e.")
	f.ClinicChairs = PullTag(form.TaggedValues, "52a.")
	f.ClinicVacantChairs = PullTag(form.TaggedValues, "52b.")
	f.ClinicFrontStaff = PullTag(form.TaggedValues, "52c.")
	f.ClinicSupportStaff = PullTag(form.TaggedValues, "52d.")
	f.ClinicProviders = PullTag(form.TaggedValues, "52e.")
	f.HomeHealthChairs = PullTag(form.TaggedValues, "53a.")
	f.HomeHealthVacantChairs = PullTag(form.TaggedValues, "53b.")
	f.HomeHealthFrontStaff = PullTag(form.TaggedValues, "53c.")
	f.HomeHealthSupportStaff = PullTag(form.TaggedValues, "53d.")
	f.HomeHealthProviders = PullTag(form.TaggedValues, "53e.")
	f.AdultDayCtrChairs = PullTag(form.TaggedValues, "54a.")
	f.AdultDayCtrVacantChairs = PullTag(form.TaggedValues, "54b.")
	f.AdultDayCtrFrontStaff = PullTag(form.TaggedValues, "54c.")
	f.AdultDayCtrSupportStaff = PullTag(form.TaggedValues, "54d.")
	f.AdultDayCtrProviders = PullTag(form.TaggedValues, "54e.")
	f.StdFooter.PullTags(form.TaggedValues)
	return f, LeftoverTagProblems(AHFacStatTag, form.FormVersion, form.TaggedValues)
}

// Encode encodes the message content.
func (f *AHFacStat) Encode() (subject, body string) {
	var (
		sb  strings.Builder
		enc *pifo.Encoder
	)
	subject = xscsubj.Encode(f.OriginMsgID, f.Handling, AHFacStatTag, f.FacilityName)
	if f.FormVersion == "" {
		f.FormVersion = "2.3"
	}
	enc = pifo.NewEncoder(&sb, AHFacStatHTML, f.FormVersion)
	f.StdHeader.EncodeBody(enc)
	enc.Write("19.", f.ReportType)
	enc.Write("20.", f.FacilityName)
	enc.Write("21.", f.FacilityType)
	enc.Write("22d.", f.Date)
	enc.Write("22t.", f.Time)
	enc.Write("23.", f.ContactName)
	enc.Write("23p.", f.ContactPhone)
	enc.Write("23f.", f.ContactFax)
	enc.Write("24.", f.OtherContact)
	enc.Write("25.", f.IncidentName)
	enc.Write("25d.", f.IncidentDate)
	enc.Write("35.", f.FacilityStatus)
	enc.Write("26a.", f.AttachOrgChart)
	enc.Write("26b.", f.AttachRR)
	enc.Write("26c.", f.AttachStatus)
	enc.Write("26d.", f.AttachActionPlan)
	enc.Write("27p.", f.EOCPhone)
	enc.Write("26e.", f.AttachDirectory)
	enc.Write("27f.", f.EOCFax)
	enc.Write("28.", f.LiaisonName)
	enc.Write("34.", f.Summary)
	enc.Write("28p.", f.LiaisonPhone)
	enc.Write("29.", f.InfoOfficerName)
	enc.Write("29p.", f.InfoOfficerPhone)
	enc.Write("29e.", f.InfoOfficerEmail)
	enc.Write("30.", f.ClosedContactName)
	enc.Write("30p.", f.ClosedContactPhone)
	enc.Write("40a.", f.SkilledNursingBedsStaffedM)
	enc.Write("40b.", f.SkilledNursingBedsStaffedF)
	enc.Write("40c.", f.SkilledNursingBedsVacantM)
	enc.Write("40d.", f.SkilledNursingBedsVacantF)
	enc.Write("40e.", f.SkilledNursingBedsSurge)
	enc.Write("30e.", f.ClosedContactEmail)
	enc.Write("41a.", f.AssistedLivingBedsStaffedM)
	enc.Write("41b.", f.AssistedLivingBedsStaffedF)
	enc.Write("41c.", f.AssistedLivingBedsVacantM)
	enc.Write("41d.", f.AssistedLivingBedsVacantF)
	enc.Write("41e.", f.AssistedLivingBedsSurge)
	enc.Write("42a.", f.SubAcuteBedsStaffedM)
	enc.Write("42b.", f.SubAcuteBedsStaffedF)
	enc.Write("42c.", f.SubAcuteBedsVacantM)
	enc.Write("42d.", f.SubAcuteBedsVacantF)
	enc.Write("42e.", f.SubAcuteBedsSurge)
	enc.Write("31a.", f.PatientsToEvacuate)
	enc.Write("43a.", f.AlzheimersBedsStaffedM)
	enc.Write("43b.", f.AlzheimersBedsStaffedF)
	enc.Write("43c.", f.AlzheimersBedsVacantM)
	enc.Write("43d.", f.AlzheimersBedsVacantF)
	enc.Write("43e.", f.AlzheimersBedsSurge)
	enc.Write("31b.", f.PatientsInjuredMinor)
	enc.Write("44a.", f.PedSubAcuteBedsStaffedM)
	enc.Write("44b.", f.PedSubAcuteBedsStaffedF)
	enc.Write("44c.", f.PedSubAcuteBedsVacantM)
	enc.Write("44d.", f.PedSubAcuteBedsVacantF)
	enc.Write("44e.", f.PedSubAcuteBedsSurge)
	enc.Write("31c.", f.PatientsTransferred)
	enc.Write("45a.", f.PsychiatricBedsStaffedM)
	enc.Write("45b.", f.PsychiatricBedsStaffedF)
	enc.Write("45c.", f.PsychiatricBedsVacantM)
	enc.Write("45d.", f.PsychiatricBedsVacantF)
	enc.Write("45e.", f.PsychiatricBedsSurge)
	enc.Write("33.", f.OtherPatientCare)
	enc.Write("46.", f.OtherCareBedsType)
	enc.Write("46a.", f.OtherCareBedsStaffedM)
	enc.Write("46b.", f.OtherCareBedsStaffedF)
	enc.Write("46c.", f.OtherCareBedsVacantM)
	enc.Write("46d.", f.OtherCareBedsVacantF)
	enc.Write("46e.", f.OtherCareBedsSurge)
	enc.Write("50a.", f.DialysisChairs)
	enc.Write("50b.", f.DialysisVacantChairs)
	enc.Write("50c.", f.DialysisFrontStaff)
	enc.Write("50d.", f.DialysisSupportStaff)
	enc.Write("50e.", f.DialysisProviders)
	enc.Write("51a.", f.SurgicalChairs)
	enc.Write("51b.", f.SurgicalVacantChairs)
	enc.Write("51c.", f.SurgicalFrontStaff)
	enc.Write("51d.", f.SurgicalSupportStaff)
	enc.Write("51e.", f.SurgicalProviders)
	enc.Write("52a.", f.ClinicChairs)
	enc.Write("52b.", f.ClinicVacantChairs)
	enc.Write("52c.", f.ClinicFrontStaff)
	enc.Write("52d.", f.ClinicSupportStaff)
	enc.Write("52e.", f.ClinicProviders)
	enc.Write("53a.", f.HomeHealthChairs)
	enc.Write("53b.", f.HomeHealthVacantChairs)
	enc.Write("53c.", f.HomeHealthFrontStaff)
	enc.Write("53d.", f.HomeHealthSupportStaff)
	enc.Write("53e.", f.HomeHealthProviders)
	enc.Write("54a.", f.AdultDayCtrChairs)
	enc.Write("54b.", f.AdultDayCtrVacantChairs)
	enc.Write("54c.", f.AdultDayCtrFrontStaff)
	enc.Write("54d.", f.AdultDayCtrSupportStaff)
	enc.Write("54e.", f.AdultDayCtrProviders)
	f.StdFooter.EncodeBody(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return subject, sb.String()
}
