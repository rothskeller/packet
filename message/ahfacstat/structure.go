package ahfacstat

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an allied health facility status form.
var Type = message.Type{
	Tag:     "AHFacStat",
	Name:    "allied health facility status form",
	Article: "an",
	Create:  New,
	Decode:  decode,
}

// AHFacStat holds an allied health facility status form.
type AHFacStat struct {
	common.StdFields
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
}

// Type returns the message type definition.
func (*AHFacStat) Type() *message.Type { return &Type }
