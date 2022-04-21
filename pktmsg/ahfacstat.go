package pktmsg

// This file defines TxAHFacStatForm and RxAHFacStatForm.

import (
	"fmt"
	"strconv"
)

type BedType uint8

const (
	SkilledNursingBed BedType = iota
	AssistedLivingBed
	SubAcuteBed
	AlzheimersDementiaBed
	PediatricSubAcuteBed
	PsychiatricBed
	OtherBed
	maxBedType
)

type BedResource uint8

const (
	StaffedBedM BedResource = iota
	StaffedBedF
	VacantBedM
	VacantBedF
	SurgeBed
	maxBedResource
)

type FacilityType uint8

const (
	DialysisFacility FacilityType = iota
	SurgicalFacility
	ClinicFacility
	HomeHealthFacility
	AdultDayCenterFacility
	maxFacilityType
)

type FacilityResource uint8

const (
	ChairsPerRoom FacilityResource = iota
	VacantChairsPerRoom
	FrontDeskStaff
	MedicalSupportStaff
	ProviderStaff
	maxFacilityResource
)

// A TxAHFacStatForm is an outgoing PackItForms-encoded message containing an
// SCCo Allied Health Facility Status form.
type TxAHFacStatForm struct {
	TxSCCoForm
	ReportType        string
	Facility          string
	FacilityType      string
	Date              string
	Time              string
	Contact           string
	ContactPhone      string
	ContactFax        string
	OtherContact      string
	Incident          string
	IncidentDate      string
	Status            string
	AttachOrgChart    string
	AttachRR          string
	AttachStatus      string
	AttachActionPlan  string
	EOCPhone          string
	AttachDirectory   string
	EOCFax            string
	Liaison           string
	LiaisonPhone      string
	InfoOfficer       string
	InfoOfficerPhone  string
	InfoOfficerEmail  string
	EOC               string
	FacilityPhone     string
	Beds              [maxBedType][maxBedResource]int
	BedResource       string
	FacilityEmail     string
	Evac              int
	Injured           int
	Transferred       int
	OtherCare         string
	FacilityResources [maxFacilityType][maxFacilityResource]int
}

var (
	// validReportType is defined in munistat.go
	validStatus = map[string]bool{"": true, "Fully Functional": true, "Limited Services": true, "Impaired/Closed": true}
	validAttach = map[string]bool{"": true, "Yes": true, "No": true}
)

// Encode returns the encoded subject line and body of the message.
func (ahfs *TxAHFacStatForm) Encode() (subject, body string, err error) {
	if err = ahfs.checkHeaderFooterFields(); err != nil {
		return "", "", err
	}
	if ahfs.Subject != "" {
		return "", "", ErrDontSet
	}
	if ahfs.ReportType == "" || ahfs.Facility == "" || ahfs.Date == "" || ahfs.Time == "" {
		return "", "", ErrIncomplete
	}
	if ahfs.ReportType == "Complete" && (ahfs.FacilityType == "" ||
		ahfs.Contact == "" || ahfs.ContactPhone == "" || ahfs.Incident == "" || ahfs.IncidentDate == "" ||
		ahfs.Status == "" || ahfs.AttachRR == "" || ahfs.EOCPhone == "" || ahfs.Liaison == "" ||
		ahfs.EOC == "" || ahfs.FacilityPhone == "" || ahfs.FacilityEmail == "") {
		return "", "", ErrIncomplete
	}
	if !validReportType[ahfs.ReportType] ||
		!validStatus[ahfs.Status] ||
		!validAttach[ahfs.AttachOrgChart] ||
		!validAttach[ahfs.AttachRR] ||
		!validAttach[ahfs.AttachStatus] ||
		!validAttach[ahfs.AttachActionPlan] ||
		!validAttach[ahfs.AttachDirectory] {
		return "", "", ErrInvalid
	}
	ahfs.FormName = "AHFacStat"
	ahfs.FormHTML = "form-allied-health-facility-status.html"
	ahfs.FormVersion = "2.2"
	ahfs.Subject = ahfs.Facility
	ahfs.encodeHeaderFields()
	ahfs.SetField("19.", ahfs.ReportType)
	ahfs.SetField("20.", ahfs.Facility)
	ahfs.SetField("21.", ahfs.FacilityType)
	ahfs.SetField("22d.", ahfs.Date)
	ahfs.SetField("22t.", ahfs.Time)
	ahfs.SetField("23.", ahfs.Contact)
	ahfs.SetField("23p.", ahfs.ContactPhone)
	ahfs.SetField("23f.", ahfs.ContactFax)
	ahfs.SetField("24.", ahfs.OtherContact)
	ahfs.SetField("25.", ahfs.Incident)
	ahfs.SetField("25d.", ahfs.IncidentDate)
	ahfs.SetField("35.", ahfs.Status)
	ahfs.SetField("26a.", ahfs.AttachOrgChart)
	ahfs.SetField("26b.", ahfs.AttachRR)
	ahfs.SetField("26c.", ahfs.AttachStatus)
	ahfs.SetField("26d.", ahfs.AttachActionPlan)
	ahfs.SetField("27p.", ahfs.EOCPhone)
	ahfs.SetField("26e.", ahfs.AttachDirectory)
	ahfs.SetField("27f.", ahfs.EOCFax)
	ahfs.SetField("28.", ahfs.Liaison)
	ahfs.SetField("28p.", ahfs.LiaisonPhone)
	ahfs.SetField("29.", ahfs.InfoOfficer)
	ahfs.SetField("29p.", ahfs.InfoOfficerPhone)
	ahfs.SetField("29e.", ahfs.InfoOfficerEmail)
	ahfs.SetField("30.", ahfs.EOC)
	ahfs.SetField("30p.", ahfs.FacilityPhone)
	ahfs.SetField("40a.", intToText(ahfs.Beds[SkilledNursingBed][StaffedBedM]))
	ahfs.SetField("40b.", intToText(ahfs.Beds[SkilledNursingBed][StaffedBedF]))
	ahfs.SetField("40c.", intToText(ahfs.Beds[SkilledNursingBed][VacantBedM]))
	ahfs.SetField("40d.", intToText(ahfs.Beds[SkilledNursingBed][VacantBedF]))
	ahfs.SetField("40e.", intToText(ahfs.Beds[SkilledNursingBed][SurgeBed]))
	ahfs.SetField("30e.", ahfs.FacilityEmail)
	ahfs.SetField("41a.", intToText(ahfs.Beds[AssistedLivingBed][StaffedBedM]))
	ahfs.SetField("41b.", intToText(ahfs.Beds[AssistedLivingBed][StaffedBedF]))
	ahfs.SetField("41c.", intToText(ahfs.Beds[AssistedLivingBed][VacantBedM]))
	ahfs.SetField("41d.", intToText(ahfs.Beds[AssistedLivingBed][VacantBedF]))
	ahfs.SetField("41e.", intToText(ahfs.Beds[AssistedLivingBed][SurgeBed]))
	ahfs.SetField("42a.", intToText(ahfs.Beds[SubAcuteBed][StaffedBedM]))
	ahfs.SetField("42b.", intToText(ahfs.Beds[SubAcuteBed][StaffedBedF]))
	ahfs.SetField("42c.", intToText(ahfs.Beds[SubAcuteBed][VacantBedM]))
	ahfs.SetField("42d.", intToText(ahfs.Beds[SubAcuteBed][VacantBedF]))
	ahfs.SetField("42e.", intToText(ahfs.Beds[SubAcuteBed][SurgeBed]))
	ahfs.SetField("31a.", intToText(ahfs.Evac))
	ahfs.SetField("43a.", intToText(ahfs.Beds[AlzheimersDementiaBed][StaffedBedM]))
	ahfs.SetField("43b.", intToText(ahfs.Beds[AlzheimersDementiaBed][StaffedBedF]))
	ahfs.SetField("43c.", intToText(ahfs.Beds[AlzheimersDementiaBed][VacantBedM]))
	ahfs.SetField("43d.", intToText(ahfs.Beds[AlzheimersDementiaBed][VacantBedF]))
	ahfs.SetField("43e.", intToText(ahfs.Beds[AlzheimersDementiaBed][SurgeBed]))
	ahfs.SetField("31b.", intToText(ahfs.Injured))
	ahfs.SetField("44a.", intToText(ahfs.Beds[PediatricSubAcuteBed][StaffedBedM]))
	ahfs.SetField("44b.", intToText(ahfs.Beds[PediatricSubAcuteBed][StaffedBedF]))
	ahfs.SetField("44c.", intToText(ahfs.Beds[PediatricSubAcuteBed][VacantBedM]))
	ahfs.SetField("44d.", intToText(ahfs.Beds[PediatricSubAcuteBed][VacantBedF]))
	ahfs.SetField("44e.", intToText(ahfs.Beds[PediatricSubAcuteBed][SurgeBed]))
	ahfs.SetField("31c.", intToText(ahfs.Transferred))
	ahfs.SetField("45a.", intToText(ahfs.Beds[PsychiatricBed][StaffedBedM]))
	ahfs.SetField("45b.", intToText(ahfs.Beds[PsychiatricBed][StaffedBedF]))
	ahfs.SetField("45c.", intToText(ahfs.Beds[PsychiatricBed][VacantBedM]))
	ahfs.SetField("45d.", intToText(ahfs.Beds[PsychiatricBed][VacantBedF]))
	ahfs.SetField("45e.", intToText(ahfs.Beds[PsychiatricBed][SurgeBed]))
	ahfs.SetField("33.", ahfs.OtherCare)
	ahfs.SetField("46.", ahfs.BedResource)
	ahfs.SetField("46a.", intToText(ahfs.Beds[OtherBed][StaffedBedM]))
	ahfs.SetField("46b.", intToText(ahfs.Beds[OtherBed][StaffedBedF]))
	ahfs.SetField("46c.", intToText(ahfs.Beds[OtherBed][VacantBedM]))
	ahfs.SetField("46d.", intToText(ahfs.Beds[OtherBed][VacantBedF]))
	ahfs.SetField("46e.", intToText(ahfs.Beds[OtherBed][SurgeBed]))
	ahfs.SetField("50a.", intToText(ahfs.FacilityResources[DialysisFacility][ChairsPerRoom]))
	ahfs.SetField("50b.", intToText(ahfs.FacilityResources[DialysisFacility][VacantChairsPerRoom]))
	ahfs.SetField("50c.", intToText(ahfs.FacilityResources[DialysisFacility][FrontDeskStaff]))
	ahfs.SetField("50d.", intToText(ahfs.FacilityResources[DialysisFacility][MedicalSupportStaff]))
	ahfs.SetField("50e.", intToText(ahfs.FacilityResources[DialysisFacility][ProviderStaff]))
	ahfs.SetField("51a.", intToText(ahfs.FacilityResources[SurgicalFacility][ChairsPerRoom]))
	ahfs.SetField("51b.", intToText(ahfs.FacilityResources[SurgicalFacility][VacantChairsPerRoom]))
	ahfs.SetField("51c.", intToText(ahfs.FacilityResources[SurgicalFacility][FrontDeskStaff]))
	ahfs.SetField("51d.", intToText(ahfs.FacilityResources[SurgicalFacility][MedicalSupportStaff]))
	ahfs.SetField("51e.", intToText(ahfs.FacilityResources[SurgicalFacility][ProviderStaff]))
	ahfs.SetField("52a.", intToText(ahfs.FacilityResources[ClinicFacility][ChairsPerRoom]))
	ahfs.SetField("52b.", intToText(ahfs.FacilityResources[ClinicFacility][VacantChairsPerRoom]))
	ahfs.SetField("52c.", intToText(ahfs.FacilityResources[ClinicFacility][FrontDeskStaff]))
	ahfs.SetField("52d.", intToText(ahfs.FacilityResources[ClinicFacility][MedicalSupportStaff]))
	ahfs.SetField("52e.", intToText(ahfs.FacilityResources[ClinicFacility][ProviderStaff]))
	ahfs.SetField("53a.", intToText(ahfs.FacilityResources[HomeHealthFacility][ChairsPerRoom]))
	ahfs.SetField("53b.", intToText(ahfs.FacilityResources[HomeHealthFacility][VacantChairsPerRoom]))
	ahfs.SetField("53c.", intToText(ahfs.FacilityResources[HomeHealthFacility][FrontDeskStaff]))
	ahfs.SetField("53d.", intToText(ahfs.FacilityResources[HomeHealthFacility][MedicalSupportStaff]))
	ahfs.SetField("53e.", intToText(ahfs.FacilityResources[HomeHealthFacility][ProviderStaff]))
	ahfs.SetField("54a.", intToText(ahfs.FacilityResources[AdultDayCenterFacility][ChairsPerRoom]))
	ahfs.SetField("54b.", intToText(ahfs.FacilityResources[AdultDayCenterFacility][VacantChairsPerRoom]))
	ahfs.SetField("54c.", intToText(ahfs.FacilityResources[AdultDayCenterFacility][FrontDeskStaff]))
	ahfs.SetField("54d.", intToText(ahfs.FacilityResources[AdultDayCenterFacility][MedicalSupportStaff]))
	ahfs.SetField("54e.", intToText(ahfs.FacilityResources[AdultDayCenterFacility][ProviderStaff]))
	ahfs.encodeFooterFields()
	return ahfs.TxSCCoForm.Encode()
}

//------------------------------------------------------------------------------

// An RxAHFacStatForm is a received PackItForms-encoded message containing an
// SCCo Allied Health Facility Status form.
type RxAHFacStatForm struct {
	RxSCCoForm
	ReportType        string
	Facility          string
	FacilityType      string
	Date              string
	Time              string
	Contact           string
	ContactPhone      string
	ContactFax        string
	OtherContact      string
	Incident          string
	IncidentDate      string
	Status            string
	AttachOrgChart    string
	AttachRR          string
	AttachStatus      string
	AttachActionPlan  string
	EOCPhone          string
	AttachDirectory   string
	EOCFax            string
	Liaison           string
	LiaisonPhone      string
	InfoOfficer       string
	InfoOfficerPhone  string
	InfoOfficerEmail  string
	EOC               string
	FacilityPhone     string
	Beds              [maxBedType][maxBedResource]int
	BedResource       string
	FacilityEmail     string
	Evac              int
	Injured           int
	Transferred       int
	OtherCare         string
	FacilityResources [maxFacilityType][maxFacilityResource]int
}

// parseRxAHFacStatForm examines an RxForm to see if it contains an EOC-213RR
// form, and if so, wraps it in an RxAHFacStatForm and returns it.  If it is not,
// it returns nil.
func parseRxAHFacStatForm(f *RxForm) *RxAHFacStatForm {
	var ahfs RxAHFacStatForm

	if f.FormHTML != "form-allied-health-facility-status.html" {
		return nil
	}
	ahfs.RxSCCoForm.RxForm = *f
	ahfs.extractHeaderFields()
	ahfs.ReportType = ahfs.Fields["19."]
	ahfs.Facility = ahfs.Fields["20."]
	ahfs.FacilityType = ahfs.Fields["21."]
	ahfs.Date = ahfs.Fields["22d."]
	ahfs.Time = ahfs.Fields["22t."]
	ahfs.Contact = ahfs.Fields["23."]
	ahfs.ContactPhone = ahfs.Fields["23p."]
	ahfs.ContactFax = ahfs.Fields["23f."]
	ahfs.OtherContact = ahfs.Fields["24."]
	ahfs.Incident = ahfs.Fields["25."]
	ahfs.IncidentDate = ahfs.Fields["25d."]
	ahfs.Status = ahfs.Fields["35."]
	ahfs.AttachOrgChart = ahfs.Fields["26a."]
	ahfs.AttachRR = ahfs.Fields["26b."]
	ahfs.AttachStatus = ahfs.Fields["26c."]
	ahfs.AttachActionPlan = ahfs.Fields["26d."]
	ahfs.EOCPhone = ahfs.Fields["27p."]
	ahfs.AttachDirectory = ahfs.Fields["26e."]
	ahfs.EOCFax = ahfs.Fields["27f."]
	ahfs.Liaison = ahfs.Fields["28."]
	ahfs.LiaisonPhone = ahfs.Fields["28p."]
	ahfs.InfoOfficer = ahfs.Fields["29."]
	ahfs.InfoOfficerPhone = ahfs.Fields["29p."]
	ahfs.InfoOfficerEmail = ahfs.Fields["29e."]
	ahfs.EOC = ahfs.Fields["30."]
	ahfs.FacilityPhone = ahfs.Fields["30p."]
	ahfs.Beds[SkilledNursingBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["40a."])
	ahfs.Beds[SkilledNursingBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["40b."])
	ahfs.Beds[SkilledNursingBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["40c."])
	ahfs.Beds[SkilledNursingBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["40d."])
	ahfs.Beds[SkilledNursingBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["40e."])
	ahfs.FacilityEmail = ahfs.Fields["30e."]
	ahfs.Beds[AssistedLivingBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["41a."])
	ahfs.Beds[AssistedLivingBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["41b."])
	ahfs.Beds[AssistedLivingBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["41c."])
	ahfs.Beds[AssistedLivingBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["41d."])
	ahfs.Beds[AssistedLivingBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["41e."])
	ahfs.Beds[SubAcuteBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["42a."])
	ahfs.Beds[SubAcuteBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["42b."])
	ahfs.Beds[SubAcuteBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["42c."])
	ahfs.Beds[SubAcuteBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["42d."])
	ahfs.Beds[SubAcuteBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["42e."])
	ahfs.Evac, _ = strconv.Atoi(ahfs.Fields["31a."])
	ahfs.Beds[AlzheimersDementiaBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["43a."])
	ahfs.Beds[AlzheimersDementiaBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["43b."])
	ahfs.Beds[AlzheimersDementiaBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["43c."])
	ahfs.Beds[AlzheimersDementiaBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["43d."])
	ahfs.Beds[AlzheimersDementiaBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["43e."])
	ahfs.Injured, _ = strconv.Atoi(ahfs.Fields["31b."])
	ahfs.Beds[PediatricSubAcuteBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["44a."])
	ahfs.Beds[PediatricSubAcuteBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["44b."])
	ahfs.Beds[PediatricSubAcuteBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["44c."])
	ahfs.Beds[PediatricSubAcuteBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["44d."])
	ahfs.Beds[PediatricSubAcuteBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["44e."])
	ahfs.Transferred, _ = strconv.Atoi(ahfs.Fields["31c."])
	ahfs.Beds[PsychiatricBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["45a."])
	ahfs.Beds[PsychiatricBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["45b."])
	ahfs.Beds[PsychiatricBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["45c."])
	ahfs.Beds[PsychiatricBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["45d."])
	ahfs.Beds[PsychiatricBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["45e."])
	ahfs.OtherCare = ahfs.Fields["33."]
	ahfs.BedResource = ahfs.Fields["46."]
	ahfs.Beds[OtherBed][StaffedBedM], _ = strconv.Atoi(ahfs.Fields["46a."])
	ahfs.Beds[OtherBed][StaffedBedF], _ = strconv.Atoi(ahfs.Fields["46b."])
	ahfs.Beds[OtherBed][VacantBedM], _ = strconv.Atoi(ahfs.Fields["46c."])
	ahfs.Beds[OtherBed][VacantBedF], _ = strconv.Atoi(ahfs.Fields["46d."])
	ahfs.Beds[OtherBed][SurgeBed], _ = strconv.Atoi(ahfs.Fields["46e."])
	ahfs.FacilityResources[DialysisFacility][ChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["50a."])
	ahfs.FacilityResources[DialysisFacility][VacantChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["50b."])
	ahfs.FacilityResources[DialysisFacility][FrontDeskStaff], _ = strconv.Atoi(ahfs.Fields["50c."])
	ahfs.FacilityResources[DialysisFacility][MedicalSupportStaff], _ = strconv.Atoi(ahfs.Fields["50d."])
	ahfs.FacilityResources[DialysisFacility][ProviderStaff], _ = strconv.Atoi(ahfs.Fields["50e."])
	ahfs.FacilityResources[SurgicalFacility][ChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["51a."])
	ahfs.FacilityResources[SurgicalFacility][VacantChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["51b."])
	ahfs.FacilityResources[SurgicalFacility][FrontDeskStaff], _ = strconv.Atoi(ahfs.Fields["51c."])
	ahfs.FacilityResources[SurgicalFacility][MedicalSupportStaff], _ = strconv.Atoi(ahfs.Fields["51d."])
	ahfs.FacilityResources[SurgicalFacility][ProviderStaff], _ = strconv.Atoi(ahfs.Fields["51e."])
	ahfs.FacilityResources[ClinicFacility][ChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["52a."])
	ahfs.FacilityResources[ClinicFacility][VacantChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["52b."])
	ahfs.FacilityResources[ClinicFacility][FrontDeskStaff], _ = strconv.Atoi(ahfs.Fields["52c."])
	ahfs.FacilityResources[ClinicFacility][MedicalSupportStaff], _ = strconv.Atoi(ahfs.Fields["52d."])
	ahfs.FacilityResources[ClinicFacility][ProviderStaff], _ = strconv.Atoi(ahfs.Fields["52e."])
	ahfs.FacilityResources[HomeHealthFacility][ChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["53a."])
	ahfs.FacilityResources[HomeHealthFacility][VacantChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["53b."])
	ahfs.FacilityResources[HomeHealthFacility][FrontDeskStaff], _ = strconv.Atoi(ahfs.Fields["53c."])
	ahfs.FacilityResources[HomeHealthFacility][MedicalSupportStaff], _ = strconv.Atoi(ahfs.Fields["53d."])
	ahfs.FacilityResources[HomeHealthFacility][ProviderStaff], _ = strconv.Atoi(ahfs.Fields["53e."])
	ahfs.FacilityResources[AdultDayCenterFacility][ChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["54a."])
	ahfs.FacilityResources[AdultDayCenterFacility][VacantChairsPerRoom], _ = strconv.Atoi(ahfs.Fields["54b."])
	ahfs.FacilityResources[AdultDayCenterFacility][FrontDeskStaff], _ = strconv.Atoi(ahfs.Fields["54c."])
	ahfs.FacilityResources[AdultDayCenterFacility][MedicalSupportStaff], _ = strconv.Atoi(ahfs.Fields["54d."])
	ahfs.FacilityResources[AdultDayCenterFacility][ProviderStaff], _ = strconv.Atoi(ahfs.Fields["54e."])
	ahfs.extractFooterFields()
	return &ahfs
}

// Valid returns whether all of the fields of the form have valid values, and
// all required fields are filled in.
func (ahfs *RxAHFacStatForm) Valid() bool {
	return ahfs.RxSCCoForm.Valid() &&
		validReportType[ahfs.ReportType] &&
		ahfs.Facility != "" &&
		ahfs.Date != "" &&
		ahfs.Time != "" &&
		validStatus[ahfs.Status] &&
		validAttach[ahfs.AttachOrgChart] &&
		validAttach[ahfs.AttachRR] &&
		validAttach[ahfs.AttachStatus] &&
		validAttach[ahfs.AttachActionPlan] &&
		validAttach[ahfs.AttachDirectory] &&
		(ahfs.ReportType == "Update" ||
			(ahfs.FacilityType != "" &&
				ahfs.Contact != "" &&
				ahfs.ContactPhone != "" &&
				ahfs.Incident != "" &&
				ahfs.IncidentDate != "" &&
				ahfs.Status != "" &&
				ahfs.AttachRR != "" &&
				ahfs.EOCPhone != "" &&
				ahfs.Liaison != "" &&
				ahfs.EOC != "" &&
				ahfs.FacilityPhone != "" &&
				ahfs.FacilityEmail != ""))
}

// EncodeSubjectLine returns what the subject line should be based on the
// received form contents.
func (ahfs *RxAHFacStatForm) EncodeSubjectLine() string {
	return fmt.Sprintf("%s_%s_AHFacStat_%s", ahfs.OriginMessageNumber, ahfs.HandlingOrder.Code(), ahfs.Facility)
}

// TypeCode returns the machine-readable code for the message type.
func (*RxAHFacStatForm) TypeCode() string { return "AHFacStat" }

// TypeName returns the human-reading name of the message type.
func (*RxAHFacStatForm) TypeName() string { return "Allied Health Facility Status form" }

// TypeArticle returns "a" or "an", whichever is appropriate for the TypeName.
func (*RxAHFacStatForm) TypeArticle() string { return "an" }
