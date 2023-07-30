// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/baseform"
	"github.com/rothskeller/packet/message/basemsg"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an allied health facility status form.
var Type = message.Type{
	Tag:     "AHFacStat",
	Name:    "allied health facility status form",
	Article: "an",
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*basemsg.FormVersion{
	{HTML: "form-allied-health-facility-status.html", Version: "2.3", Tag: "AHFacStat", FieldOrder: fieldOrder},
	{HTML: "form-allied-health-facility-status.html", Version: "2.2", Tag: "AHFacStat", FieldOrder: fieldOrder},
	{HTML: "form-allied-health-facility-status.html", Version: "2.1", Tag: "AHFacStat", FieldOrder: fieldOrder},
	{HTML: "form-allied-health-facility-status.html", Version: "2.0", Tag: "AHFacStat", FieldOrder: fieldOrder},
}
var fieldOrder = []string{
	"MsgNo", "1a.", "1b.", "5.", "7a.", "8a.", "7b.", "8b.", "7c.", "8c.", "7d.", "8d.", "19.", "20.", "21.", "22d.", "22t.",
	"23.", "23p.", "23f.", "24.", "25.", "25d.", "35.", "26a.", "26b.", "26c.", "26d.", "27p.", "26e.", "27f.", "28.", "34.",
	"28p.", "29.", "29p.", "29e.", "30.", "30p.", "40a.", "40b.", "40c.", "40d.", "40e.", "30e.", "41a.", "41b.", "41c.",
	"41d.", "41e.", "42a.", "42b.", "42c.", "42d.", "42e.", "31a.", "43a.", "43b.", "43c.", "43d.", "43e.", "31b.", "44a.",
	"44b.", "44c.", "44d.", "44e.", "31c.", "45a.", "45b.", "45c.", "45d.", "45e.", "33.", "46.", "46a.", "46b.", "46c.",
	"46d.", "46e.", "50a.", "50b.", "50c.", "50d.", "50e.", "51a.", "51b.", "51c.", "51d.", "51e.", "52a.", "52b.", "52c.",
	"52d.", "52e.", "53a.", "53b.", "53c.", "53d.", "53e.", "54a.", "54b.", "54c.", "54d.", "54e.", "OpRelayRcvd",
	"OpRelaySent", "OpName", "OpCall", "OpDate", "OpTime",
}

// AHFacStat holds an allied health facility status form.
type AHFacStat struct {
	basemsg.BaseMessage
	baseform.BaseForm
	ReportType           string
	FacilityName         string
	FacilityType         string
	Date                 string
	Time                 string
	ContactName          string
	ContactPhone         string
	ContactFax           string
	OtherContact         string
	IncidentName         string
	IncidentDate         string
	FacilityStatus       string
	EOCPhone             string
	EOCFax               string
	LiaisonName          string
	LiaisonPhone         string
	InfoOfficerName      string
	InfoOfficerPhone     string
	InfoOfficerEmail     string
	ClosedContactName    string
	ClosedContactPhone   string
	ClosedContactEmail   string
	PatientsToEvacuate   string
	PatientsInjuredMinor string
	PatientsTransferred  string
	OtherPatientCare     string
	AttachOrgChart       string
	AttachRR             string
	AttachStatus         string
	AttachActionPlan     string
	AttachDirectory      string
	Summary              string
	SkilledNursingBeds   BedCounts
	AssistedLivingBeds   BedCounts
	SubAcuteBeds         BedCounts
	AlzheimersBeds       BedCounts
	PedSubAcuteBeds      BedCounts
	PsychiatricBeds      BedCounts
	OtherCareBedsType    string
	OtherCareBeds        BedCounts
	DialysisResources    ResourceCounts
	SurgicalResources    ResourceCounts
	ClinicResources      ResourceCounts
	HomeHealthResources  ResourceCounts
	AdultDayCtrResources ResourceCounts
}
type BedCounts struct {
	StaffedM string
	StaffedF string
	VacantM  string
	VacantF  string
	Surge    string
}
type ResourceCounts struct {
	Chairs       string
	VacantChairs string
	FrontStaff   string
	SupportStaff string
	Providers    string
}

func New() (f *AHFacStat) {
	f = create(versions[0]).(*AHFacStat)
	f.MessageDate = time.Now().Format("01/02/2006")
	f.Date = f.MessageDate
	f.Handling = "ROUTINE"
	return f
}

var pdfBase []byte

func create(version *basemsg.FormVersion) message.Message {
	var f = AHFacStat{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		PDFBase:     pdfBase,
		Form:        version,
	}}
	f.BaseMessage.FSubject = &f.FacilityName
	f.BaseMessage.FReportType = &f.ReportType
	f.BaseMessage.FBody = &f.Summary
	f.Fields = make([]*basemsg.Field, 0, 130)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &baseform.DefaultPDFMaps)
	f.Fields = append(f.Fields,
		&basemsg.Field{
			PDFMap: basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
				return []basemsg.PDFField{{Name: "Form Type", Value: "Allied Health Facility Status"}}
			}),
		},
		&basemsg.Field{
			Label:     "Report Type",
			PIFOTag:   "19.",
			Value:     &f.ReportType,
			Choices:   basemsg.Choices{"Update", "Complete"},
			Presence:  basemsg.Required,
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			EditWidth: 7,
			EditHelp:  `This indicates whether the form should "Update" the previous status report for the facility, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
		},
		&basemsg.Field{
			Label:    "Facility Name",
			PIFOTag:  "20.",
			Value:    &f.FacilityName,
			Presence: basemsg.Required,
			Compare:  common.CompareText,
			PDFMap: basemsg.PDFMapFunc(func(f *basemsg.Field) []basemsg.PDFField {
				return []basemsg.PDFField{
					{Name: "Form Topic", Value: *f.Value},
					{Name: "FACILTY TYPE TIME DATE FACILITY NAME", Value: *f.Value},
				}
			}),
			EditWidth: 52,
			EditHelp:  `This is the name of the facility whose status is being reported.  It is required.`,
		},
		&basemsg.Field{
			Label:     "Facility Type",
			PIFOTag:   "21.",
			Value:     &f.FacilityType,
			Presence:  f.requiredForComplete,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("facility type"),
			EditWidth: 21,
			EditHelp:  `This is the type of the facility, such as Skilled Nursing, Home Health, Dialysis, Community Health Center, Surgical Center, or Hospice.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:      "Date",
			PIFOTag:    "22d.",
			Value:      &f.Date,
			Presence:   basemsg.Required,
			PIFOValid:  basemsg.ValidDate,
			Compare:    common.CompareDate,
			PDFMap:     basemsg.PDFName("date"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:      "Time",
			PIFOTag:    "22t.",
			Value:      &f.Time,
			Presence:   basemsg.Required,
			PIFOValid:  basemsg.ValidTime,
			Compare:    common.CompareTime,
			PDFMap:     basemsg.PDFName("time"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Date/Time",
			TableValue: func(*basemsg.Field) string {
				return common.SmartJoin(f.Date, f.Time, " ")
			},
			EditWidth: 16,
			EditHelp:  `This is the date and time of the status report, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
			EditHint:  "MM/DD/YYYY HH:MM",
			EditValue: func(_ *basemsg.Field) string {
				return basemsg.ValueDateTime(f.Date, f.Time)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				basemsg.ApplyDateTime(&f.Date, &f.Time, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return basemsg.ValidDateTime(field, f.Date, f.Time)
			},
			Presence: basemsg.Required,
		},
		&basemsg.Field{
			Label:     "Contact Name",
			PIFOTag:   "23.",
			Value:     &f.ContactName,
			Presence:  f.requiredForComplete,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Contact Name"),
			EditWidth: 52,
			EditHelp:  `This is the name of the person to be contacted with questions about this report.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Contact Phone",
			PIFOTag:   "23p.",
			Value:     &f.ContactPhone,
			Presence:  f.requiredForComplete,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("Phone"),
			EditWidth: 19,
			EditHelp:  `This is the phone number of the person to be contacted with questions about this report.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Contact Fax",
			PIFOTag:   "23f.",
			Value:     &f.ContactFax,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("Fax"),
			EditWidth: 22,
			EditHelp:  `This is the fax number of the person to be contacted with questions about this report.`,
		},
		&basemsg.Field{
			Label:     "Other Contact",
			PIFOTag:   "24.",
			Value:     &f.OtherContact,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("Other Phone Fax Cell Phone Radio"),
			EditWidth: 53,
			EditHelp:  `This is additional contact information for the person to be contacted with questions about this report.`,
		},
		&basemsg.Field{
			Label:    "Incident Name",
			PIFOTag:  "25.",
			Value:    &f.IncidentName,
			Presence: f.requiredForComplete,
			Compare:  common.CompareText,
			PDFMap: basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
				return []basemsg.PDFField{{
					Name:  "Incident Name and Date",
					Value: common.SmartJoin(f.IncidentName, f.IncidentDate, " "),
				}}
			}),
			EditWidth: 42,
			EditHelp:  `This is the assigned incident name of the incident for which this report is filed.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Incident Date",
			PIFOTag:   "25d.",
			Value:     &f.IncidentDate,
			Presence:  f.requiredForComplete,
			PIFOValid: basemsg.ValidDate,
			Compare:   common.CompareDate,
			EditWidth: 10,
			EditHelp:  `This is the date of the incident for which this report is filed, in MM/DD/YYYY format.  It is required when the "Report Type" is "Complete".`,
			EditHint:  "MM/DD/YYYY",
			EditApply: basemsg.ApplyDate,
		},
		&basemsg.Field{
			Label:     "Facility Status",
			PIFOTag:   "35.",
			Value:     &f.FacilityStatus,
			Choices:   basemsg.Choices{"Fully Functional", "Limited Services", "Impaired/Closed"},
			Presence:  f.requiredForComplete,
			PIFOValid: basemsg.ValidRestricted,
			Compare:   common.CompareExact,
			PDFMap: basemsg.PDFMapFunc(func(f *basemsg.Field) []basemsg.PDFField {
				var name string
				switch *f.Value {
				case "Fully Functional":
					name = "CHECK ONEGREEN FULLY FUNCTIONAL"
				case "Limited Services":
					name = "CHECK ONERED LIMITED SERVICES"
				case "Impaired/Closed":
					name = "CHECK ONEBLACK IMPAIREDCLOSED"
				}
				if name != "" {
					return []basemsg.PDFField{{Name: name, Value: "X"}}
				}
				return nil
			}),
			EditWidth: 16,
			EditHelp:  `This indicates the status of the facility.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "EOC Main Contact Number",
			PIFOTag:   "27p.",
			Value:     &f.EOCPhone,
			Presence:  f.requiredForComplete,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("EOC MAIN CONTACT NUMBER"),
			EditWidth: 19,
			EditHelp:  `This is the main phone number for the facility's Emergency Operations Center (EOC).  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "EOC Main Contact Fax",
			PIFOTag:   "27f.",
			Value:     &f.EOCFax,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("EOC MAIN CONTACT FAX"),
			EditWidth: 20,
			EditHelp:  `This is the max fax number for the facility's Emergency Operations Center (EOC).`,
		},
		&basemsg.Field{
			Label:     "Liaison Officer Name",
			PIFOTag:   "28.",
			Value:     &f.LiaisonName,
			Presence:  f.requiredForComplete,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("NAME LIAISON TO PUBLIC HEALTHMEDICAL HEALTH BRANCH"),
			EditWidth: 17,
			EditHelp:  `This is the name of the facility's liaison to the Public Health or Medical Health Branch.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Liaison Contact Number",
			PIFOTag:   "28p.",
			Value:     &f.LiaisonPhone,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT NUMBER"),
			EditWidth: 17,
			EditHelp:  `This is the phone number of the facility's liaison to the Public Health or Medical Health Branch.`,
		},
		&basemsg.Field{
			Label:     "Info Officer Name",
			PIFOTag:   "29.",
			Value:     &f.InfoOfficerName,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("INFORMATION OFFICER NAME"),
			EditWidth: 17,
			EditHelp:  `This is the name of the facility's information officer.`,
		},
		&basemsg.Field{
			Label:     "Info Officer Contact Number",
			PIFOTag:   "29p.",
			Value:     &f.InfoOfficerPhone,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT NUMBER_2"),
			EditWidth: 17,
			EditHelp:  `This is the phone number of the facility's information officer.`,
		},
		&basemsg.Field{
			Label:     "Info Officer Contact Email",
			PIFOTag:   "29e.",
			Value:     &f.InfoOfficerEmail,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT EMAIL"),
			EditWidth: 17,
			EditHelp:  `This is the email address of the facility's information officer.`,
		},
		&basemsg.Field{
			Label:     "Not Active Contact Name",
			PIFOTag:   "30.",
			Value:     &f.ClosedContactName,
			Presence:  f.requiredForComplete,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("IF EOC IS NOT ACTIVATED WHO SHOULD BE CONTACTED FOR QUESTIONSREQUESTS"),
			EditWidth: 17,
			EditHelp:  `This is the name of the person to be contacted with questions or requests when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Facility Contact Phone",
			PIFOTag:   "30p.",
			Value:     &f.ClosedContactPhone,
			Presence:  f.requiredForComplete,
			PIFOValid: basemsg.ValidPhoneNumber,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT NUMBER_3"),
			EditWidth: 17,
			EditHelp:  `This is the phone number of the person to be contacted when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Facility Contact Email",
			PIFOTag:   "30e.",
			Value:     &f.ClosedContactEmail,
			Presence:  f.requiredForComplete,
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT EMAIL_2"),
			EditWidth: 17,
			EditHelp:  `This is the email address of the person to be contacted when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:     "Patients To Evacuate",
			PIFOTag:   "31a.",
			Value:     &f.PatientsToEvacuate,
			PIFOValid: basemsg.ValidCardinal,
			Compare:   common.CompareCardinal,
			PDFMap:    basemsg.PDFName("TOTALPATIENTS TO EVACUATE"),
			EditWidth: 17,
			EditHelp:  `This is the number of patients who need evacuation.`,
			EditApply: basemsg.ApplyCardinal,
		},
		&basemsg.Field{
			Label:     "Patients Injured - Minor",
			PIFOTag:   "31b.",
			Value:     &f.PatientsInjuredMinor,
			PIFOValid: basemsg.ValidCardinal,
			Compare:   common.CompareCardinal,
			PDFMap:    basemsg.PDFName("TOTALPATIENTS  INJURED  MINOR"),
			EditWidth: 17,
			EditHelp:  `This is the number of patients with minor injuries.`,
			EditApply: basemsg.ApplyCardinal,
		},
		&basemsg.Field{
			Label:     "Patients Transferred",
			PIFOTag:   "31c.",
			Value:     &f.PatientsTransferred,
			PIFOValid: basemsg.ValidCardinal,
			Compare:   common.CompareCardinal,
			PDFMap:    basemsg.PDFName("TOTALPATIENTS TRANSFERED OUT OF COUNTY"),
			EditWidth: 17,
			EditHelp:  `This is the number of patients who have been transferred out of the county.`,
			EditApply: basemsg.ApplyCardinal,
		},
		&basemsg.Field{
			Label:   "Other Patient Care Info",
			PIFOTag: "33.",
			Value:   &f.OtherPatientCare,
			Compare: common.CompareText,
			PDFMap: basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
				// We put OtherCareBedsType in this field as
				// well because the PDF doesn't have a fillable
				// field for it, but PackItForms does.
				return []basemsg.PDFField{{
					Name:  "OTHER PATIENT CARE INFORMATION",
					Value: common.SmartJoin(f.OtherPatientCare, f.OtherCareBedsType, " "),
				}}
			}),
			EditWidth: 27,
			EditHelp:  `This field contains other patient care information as needed.`,
		},
		&basemsg.Field{
			Label:      "Attached Org Chart",
			PIFOTag:    "26a.",
			Value:      &f.AttachOrgChart,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOValid:  basemsg.ValidRestricted,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("YesNoNHICSICS ORGANIZATION CHART"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  3,
			EditHelp:   `This indicates whether an NHICS/ICS organization chart is attached to the status report.`,
		},
		&basemsg.Field{
			Label:      "Attached Resource Requests",
			PIFOTag:    "26b.",
			Value:      &f.AttachRR,
			Choices:    basemsg.Choices{"Yes", "No"},
			Presence:   f.requiredForComplete,
			PIFOValid:  basemsg.ValidRestricted,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("YesNoDEOC9A RESOURCE REQUEST FORMS"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  3,
			EditHelp:   `This indicates whether DEOC-9A resource request forms are attached to the status report.  It is required when the "Report Type" is "Complete".`,
		},
		&basemsg.Field{
			Label:      "Attached Status Report",
			PIFOTag:    "26c.",
			Value:      &f.AttachStatus,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOValid:  basemsg.ValidRestricted,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("YesNoNHICSICS STATUS REPORT FORM  STANDARD"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  3,
			EditHelp:   `This indicates whether an NHICS/ICS standard status report form is attached to this status report.`,
		},
		&basemsg.Field{
			Label:      "Attached Incident Action Plan",
			PIFOTag:    "26d.",
			Value:      &f.AttachActionPlan,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOValid:  basemsg.ValidRestricted,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("YesNoNHICSICS INCIDENT ACTION PLAN"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  3,
			EditHelp:   `This indicates whether an NHICS/ICS incident action plan is attached to the status report.`,
		},
		&basemsg.Field{
			Label:      "Attached Phone Directory",
			PIFOTag:    "26e.",
			Value:      &f.AttachDirectory,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOValid:  basemsg.ValidRestricted,
			Compare:    common.CompareExact,
			PDFMap:     basemsg.PDFName("YesNoPHONECOMMUNICATIONS DIRECTORY"),
			TableValue: basemsg.OmitFromTable,
			EditWidth:  3,
			EditHelp:   `This indicates whether a phone/communications directory is attached to the status report.`,
		},
		&basemsg.Field{
			Label: "Attachments",
			TableValue: func(*basemsg.Field) string {
				var set []string
				if f.AttachOrgChart == "Yes" {
					set = append(set, "Organization Chart")
				}
				if f.AttachRR == "Yes" {
					set = append(set, "Resource Request Forms")
				}
				if f.AttachStatus == "Yes" {
					set = append(set, "Status Report Form")
				}
				if f.AttachActionPlan == "Yes" {
					set = append(set, "Incident Action Plan")
				}
				if f.AttachDirectory == "Yes" {
					set = append(set, "Phone/Communications Directory")
				}
				return strings.Join(set, ", ")
			},
		},
		&basemsg.Field{
			Label:     "General Summary",
			PIFOTag:   "34.",
			Value:     &f.Summary,
			Compare:   common.CompareText,
			PDFMap:    basemsg.PDFName("GENERAL SUMMARY OF SITUATIONCONDITIONSRow1"),
			EditWidth: 41,
			Multiline: true,
			EditHelp:  `This is a general summary of the situation and conditions at the facility.`,
		},
		&basemsg.Field{
			Label:      "Skilled Nursing Beds: Staffed M",
			PIFOTag:    "40a.",
			Value:      &f.SkilledNursingBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MSKILLED NURSING"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	var first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Skilled Nursing Beds: Staffed F",
			PIFOTag: "40b.",
			Value:   &f.SkilledNursingBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFSKILLED NURSING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Skilled Nursing Beds: Vacant M",
			PIFOTag: "40c.",
			Value:   &f.SkilledNursingBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMSKILLED NURSING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Skilled Nursing Beds: Vacant F",
			PIFOTag: "40d.",
			Value:   &f.SkilledNursingBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFSKILLED NURSING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Skilled Nursing Beds: Surge",
			PIFOTag: "40e.",
			Value:   &f.SkilledNursingBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge SKILLED NURSING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Skilled Nursing Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.SkilledNursingBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of skilled nursing beds at the facility.  Enter five numbers separated by spaces: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.SkilledNursingBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.SkilledNursingBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.SkilledNursingBeds)
			},
		},
		&basemsg.Field{
			Label:      "Assisted Living Beds: Staffed M",
			PIFOTag:    "41a.",
			Value:      &f.AssistedLivingBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MASSISTED LIVING"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Assisted Living Beds: Staffed F",
			PIFOTag: "41b.",
			Value:   &f.AssistedLivingBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFASSISTED LIVING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Assisted Living Beds: Vacant M",
			PIFOTag: "41c.",
			Value:   &f.AssistedLivingBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMASSISTED LIVING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Assisted Living Beds: Vacant F",
			PIFOTag: "41d.",
			Value:   &f.AssistedLivingBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFASSISTED LIVING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Assisted Living Beds: Surge",
			PIFOTag: "41e.",
			Value:   &f.AssistedLivingBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge ASSISTED LIVING"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Assisted Living Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.AssistedLivingBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of assisted living beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.AssistedLivingBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.AssistedLivingBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.AssistedLivingBeds)
			},
		},
		&basemsg.Field{
			Label:      "Sub-Acute Beds: Staffed M",
			PIFOTag:    "42a.",
			Value:      &f.SubAcuteBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Sub-Acute Beds: Staffed F",
			PIFOTag: "42b.",
			Value:   &f.SubAcuteBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Sub-Acute Beds: Vacant M",
			PIFOTag: "42c.",
			Value:   &f.SubAcuteBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Sub-Acute Beds: Vacant F",
			PIFOTag: "42d.",
			Value:   &f.SubAcuteBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Sub-Acute Beds: Surge",
			PIFOTag: "42e.",
			Value:   &f.SubAcuteBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge SURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Sub-Acute Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.SubAcuteBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of sub-acute beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.SubAcuteBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.SubAcuteBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.SubAcuteBeds)
			},
		},
		&basemsg.Field{
			Label:      "Alzheimers Beds: Staffed M",
			PIFOTag:    "43a.",
			Value:      &f.AlzheimersBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Alzheimers Beds: Staffed F",
			PIFOTag: "43b.",
			Value:   &f.AlzheimersBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Alzheimers Beds: Vacant M",
			PIFOTag: "43c.",
			Value:   &f.AlzheimersBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Alzheimers Beds: Vacant F",
			PIFOTag: "43d.",
			Value:   &f.AlzheimersBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Alzheimers Beds: Surge",
			PIFOTag: "43e.",
			Value:   &f.AlzheimersBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge SUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Alzheimers/Dementia Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.AlzheimersBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of Alzheimers/dementia beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.AlzheimersBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.AlzheimersBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.AlzheimersBeds)
			},
		},
		&basemsg.Field{
			Label:      "Ped Sub-Acute Beds: Staffed M",
			PIFOTag:    "44a.",
			Value:      &f.PedSubAcuteBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Staffed F",
			PIFOTag: "44b.",
			Value:   &f.PedSubAcuteBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Vacant M",
			PIFOTag: "44c.",
			Value:   &f.PedSubAcuteBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Vacant F",
			PIFOTag: "44d.",
			Value:   &f.PedSubAcuteBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Surge",
			PIFOTag: "44e.",
			Value:   &f.PedSubAcuteBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge ALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Pediatric Sub-Acute Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.PedSubAcuteBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of pediatric sub-acute beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.PedSubAcuteBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.PedSubAcuteBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.PedSubAcuteBeds)
			},
		},
		&basemsg.Field{
			Label:      "Psychiatric Beds: Staffed M",
			PIFOTag:    "45a.",
			Value:      &f.PsychiatricBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Psychiatric Beds: Staffed F",
			PIFOTag: "45b.",
			Value:   &f.PsychiatricBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Psychiatric Beds: Vacant M",
			PIFOTag: "45c.",
			Value:   &f.PsychiatricBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Psychiatric Beds: Vacant F",
			PIFOTag: "45d.",
			Value:   &f.PsychiatricBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Psychiatric Beds: Surge",
			PIFOTag: "45e.",
			Value:   &f.PsychiatricBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge PEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Psychiatric Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.PsychiatricBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of psychiatric beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.PsychiatricBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.PsychiatricBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.PsychiatricBeds)
			},
		},
		&basemsg.Field{
			Label:   "Other Care Beds Type",
			PIFOTag: "46.",
			Value:   &f.OtherCareBedsType,
			Compare: common.CompareText,
			// The PDF doesn't have a fillable field for this, so
			// its value is added to the Other Patient Care Info
			// field, above.
			EditWidth: 17,
			EditHelp:  `This is the other type of beds available at the facility, if any.`,
		},
		&basemsg.Field{
			Label:      "Other Care Beds: Staffed M",
			PIFOTag:    "46a.",
			Value:      &f.OtherCareBeds.StaffedM,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed Bed MPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Other Care Beds: Staffed F",
			PIFOTag: "46b.",
			Value:   &f.OtherCareBeds.StaffedF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Staffed BedFPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Other Care Beds: Vacant M",
			PIFOTag: "46c.",
			Value:   &f.OtherCareBeds.VacantM,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedsMPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Other Care Beds: Vacant F",
			PIFOTag: "46d.",
			Value:   &f.OtherCareBeds.VacantF,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Vacant BedFPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Other Care Beds: Surge",
			PIFOTag: "46e.",
			Value:   &f.OtherCareBeds.Surge,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("Surge PSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Other Care Beds",
			TableValue: func(*basemsg.Field) string {
				return bedsTableValue(&f.OtherCareBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of beds at the facility of the named other type.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*basemsg.Field) string {
				return bedsValue(&f.OtherCareBeds)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				bedsApply(&f.OtherCareBeds, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return bedsValid(field, &f.OtherCareBeds)
			},
			EditSkip: func() bool {
				return f.OtherCareBedsType == ""
			},
		},
		&basemsg.Field{
			Label:      "Dialysis: Chairs",
			PIFOTag:    "50a.",
			Value:      &f.DialysisResources.Chairs,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSDIALYSIS"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Dialysis: Vacant Chairs",
			PIFOTag: "50b.",
			Value:   &f.DialysisResources.VacantChairs,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMDIALYSIS"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Dialysis: Front Staff",
			PIFOTag: "50c.",
			Value:   &f.DialysisResources.FrontStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFDIALYSIS"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Dialysis: Support Staff",
			PIFOTag: "50d.",
			Value:   &f.DialysisResources.SupportStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFDIALYSIS"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Dialysis: Providers",
			PIFOTag: "50e.",
			Value:   &f.DialysisResources.Providers,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("PROVIDER STAFFDIALYSIS"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Dialysis Resources",
			TableValue: func(*basemsg.Field) string {
				return resourcesTableValue(&f.DialysisResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of dialysis resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*basemsg.Field) string {
				return resourcesValue(&f.DialysisResources)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				resourcesApply(&f.DialysisResources, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return resourcesValid(field, &f.DialysisResources)
			},
		},
		&basemsg.Field{
			Label:      "Surgical: Chairs",
			PIFOTag:    "51a.",
			Value:      &f.SurgicalResources.Chairs,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSSURGICAL"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Surgical: Vacant Chairs",
			PIFOTag: "51b.",
			Value:   &f.SurgicalResources.VacantChairs,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMSURGICAL"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Surgical: Front Staff",
			PIFOTag: "51c.",
			Value:   &f.SurgicalResources.FrontStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFSURGICAL"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Surgical: Support Staff",
			PIFOTag: "51d.",
			Value:   &f.SurgicalResources.SupportStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFSURGICAL"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Surgical: Providers",
			PIFOTag: "51e.",
			Value:   &f.SurgicalResources.Providers,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("PROVIDER STAFFSURGICAL"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Surgical Resources",
			TableValue: func(*basemsg.Field) string {
				return resourcesTableValue(&f.SurgicalResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of surgical resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*basemsg.Field) string {
				return resourcesValue(&f.SurgicalResources)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				resourcesApply(&f.SurgicalResources, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return resourcesValid(field, &f.SurgicalResources)
			},
		},
		&basemsg.Field{
			Label:      "Clinic: Chairs",
			PIFOTag:    "52a.",
			Value:      &f.ClinicResources.Chairs,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSCLINIC"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Clinic: Vacant Chairs",
			PIFOTag: "52b.",
			Value:   &f.ClinicResources.VacantChairs,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMCLINIC"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Clinic: Front Staff",
			PIFOTag: "52c.",
			Value:   &f.ClinicResources.FrontStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFCLINIC"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Clinic: Support Staff",
			PIFOTag: "52d.",
			Value:   &f.ClinicResources.SupportStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFCLINIC"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Clinic: Providers",
			PIFOTag: "52e.",
			Value:   &f.ClinicResources.Providers,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("PROVIDER STAFFCLINIC"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Clinic Resources",
			TableValue: func(*basemsg.Field) string {
				return resourcesTableValue(&f.ClinicResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of clinic resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*basemsg.Field) string {
				return resourcesValue(&f.ClinicResources)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				resourcesApply(&f.ClinicResources, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return resourcesValid(field, &f.ClinicResources)
			},
		},
		&basemsg.Field{
			Label:      "Home Health: Chairs",
			PIFOTag:    "53a.",
			Value:      &f.HomeHealthResources.Chairs,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSHOMEHEALTH"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Home Health: Vacant Chairs",
			PIFOTag: "53b.",
			Value:   &f.HomeHealthResources.VacantChairs,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMHOMEHEALTH"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Home Health: Front Staff",
			PIFOTag: "53c.",
			Value:   &f.HomeHealthResources.FrontStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFHOMEHEALTH"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Home Health: Support Staff",
			PIFOTag: "53d.",
			Value:   &f.HomeHealthResources.SupportStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFHOMEHEALTH"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Home Health: Providers",
			PIFOTag: "53e.",
			Value:   &f.HomeHealthResources.Providers,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("PROVIDER STAFFHOMEHEALTH"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Home Health Resources",
			TableValue: func(*basemsg.Field) string {
				return resourcesTableValue(&f.HomeHealthResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of home health resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*basemsg.Field) string {
				return resourcesValue(&f.HomeHealthResources)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				resourcesApply(&f.HomeHealthResources, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return resourcesValid(field, &f.HomeHealthResources)
			},
		},
		&basemsg.Field{
			Label:      "Adult Day Ctr: Chairs",
			PIFOTag:    "54a.",
			Value:      &f.AdultDayCtrResources.Chairs,
			PIFOValid:  basemsg.ValidCardinal,
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSADULT DAY CENTER"),
			TableValue: basemsg.OmitFromTable,
		},
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		&basemsg.Field{
			Label:   "Adult Day Ctr: Vacant Chairs",
			PIFOTag: "54b.",
			Value:   &f.AdultDayCtrResources.VacantChairs,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMADULT DAY CENTER"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Adult Day Ctr: Front Staff",
			PIFOTag: "54c.",
			Value:   &f.AdultDayCtrResources.FrontStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFADULT DAY CENTER"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Adult Day Ctr: Support Staff",
			PIFOTag: "54d.",
			Value:   &f.AdultDayCtrResources.SupportStaff,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFADULT DAY CENTER"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label:   "Adult Day Ctr: Providers",
			PIFOTag: "54e.",
			Value:   &f.AdultDayCtrResources.Providers,
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			Compare:    common.CompareCardinal,
			PDFMap:     basemsg.PDFName("PROVIDER STAFFADULT DAY CENTER"),
			TableValue: basemsg.OmitFromTable,
		},
		&basemsg.Field{
			Label: "Adult Day Center Resources",
			TableValue: func(*basemsg.Field) string {
				return resourcesTableValue(&f.AdultDayCtrResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of adult day center resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*basemsg.Field) string {
				return resourcesValue(&f.AdultDayCtrResources)
			},
			EditApply: func(_ *basemsg.Field, value string) {
				resourcesApply(&f.AdultDayCtrResources, value)
			},
			EditValid: func(field *basemsg.Field) string {
				return resourcesValid(field, &f.AdultDayCtrResources)
			},
		},
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &baseform.DefaultPDFMaps)
	return &f
}

func (f *AHFacStat) requiredForComplete() (basemsg.Presence, string) {
	if f.ReportType == "Complete" {
		return basemsg.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return basemsg.PresenceOptional, ""
}

func decode(subject, body string) (f *AHFacStat) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-allied-health-facility-status.html") {
		return nil
	}
	return basemsg.Decode(body, versions, create).(*AHFacStat)
}

func allOrNone(first, current *basemsg.Field) string {
	if *first.Value == "" && *current.Value != "" {
		return fmt.Sprintf("The %q field must not have a value unless %q has a value.  (Either all fields on the row must have a value, or none.)", current.Label, first.Label)
	}
	if *first.Value != "" && *current.Value == "" {
		return fmt.Sprintf("The %q field is required when %q has a value.  (Either all fields on the row must have a value, or none.)", current.Label, first.Label)
	}
	return basemsg.ValidCardinal(current)
}

func bedsTableValue(beds *BedCounts) string {
	if beds.StaffedM == "" && beds.StaffedF == "" && beds.VacantM == "" && beds.VacantF == "" && beds.Surge == "" {
		return ""
	}
	return fmt.Sprintf("%3s %3s %3s %3s %3s", beds.StaffedM, beds.StaffedF, beds.VacantM, beds.VacantF, beds.Surge)
}
func bedsValue(beds *BedCounts) string {
	if beds.StaffedM == "" && beds.StaffedF == "" && beds.VacantM == "" && beds.VacantF == "" && beds.Surge == "" {
		return ""
	}
	return fmt.Sprintf("%s %s %s %s %s", beds.StaffedM, beds.StaffedF, beds.VacantM, beds.VacantF, beds.Surge)
}
func bedsApply(beds *BedCounts, value string) {
	var f basemsg.Field
	values := strings.Fields(value)
	if len(values) > 0 {
		f.Value = &beds.StaffedM
		basemsg.ApplyCardinal(&f, values[0])
	} else {
		beds.StaffedM = ""
	}
	if len(values) > 1 {
		f.Value = &beds.StaffedF
		basemsg.ApplyCardinal(&f, values[1])
	} else {
		beds.StaffedF = ""
	}
	if len(values) > 2 {
		f.Value = &beds.VacantM
		basemsg.ApplyCardinal(&f, values[2])
	} else {
		beds.VacantM = ""
	}
	if len(values) > 3 {
		f.Value = &beds.VacantF
		basemsg.ApplyCardinal(&f, values[3])
	} else {
		beds.VacantF = ""
	}
	if len(values) > 4 {
		f.Value = &beds.Surge
		basemsg.ApplyCardinal(&f, strings.Join(values[4:], " "))
	} else {
		beds.Surge = ""
	}
}
func bedsValid(field *basemsg.Field, beds *BedCounts) string {
	if beds.StaffedM == "" && beds.StaffedF == "" && beds.VacantM == "" && beds.VacantF == "" && beds.Surge == "" {
		return ""
	}
	if !common.PIFOCardinalNumberRE.MatchString(beds.StaffedM) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(beds.StaffedF) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(beds.VacantM) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(beds.VacantF) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(beds.Surge) {
		goto INVALID
	}
	return ""
INVALID:
	return fmt.Sprintf("The %q field does not contain a valid value.  It should contain five numbers separated by spaces.", field.Label)
}

func resourcesTableValue(resources *ResourceCounts) string {
	if resources.Chairs == "" && resources.VacantChairs == "" && resources.FrontStaff == "" && resources.SupportStaff == "" && resources.Providers == "" {
		return ""
	}
	return fmt.Sprintf("%3s %3s %3s %3s %3s", resources.Chairs, resources.VacantChairs, resources.FrontStaff, resources.SupportStaff, resources.Providers)
}
func resourcesValue(resources *ResourceCounts) string {
	if resources.Chairs == "" && resources.VacantChairs == "" && resources.FrontStaff == "" && resources.SupportStaff == "" && resources.Providers == "" {
		return ""
	}
	return fmt.Sprintf("%s %s %s %s %s", resources.Chairs, resources.VacantChairs, resources.FrontStaff, resources.SupportStaff, resources.Providers)
}
func resourcesApply(resources *ResourceCounts, value string) {
	var f basemsg.Field
	values := strings.Fields(value)
	if len(values) > 0 {
		f.Value = &resources.Chairs
		basemsg.ApplyCardinal(&f, values[0])
	} else {
		resources.Chairs = ""
	}
	if len(values) > 1 {
		f.Value = &resources.VacantChairs
		basemsg.ApplyCardinal(&f, values[1])
	} else {
		resources.VacantChairs = ""
	}
	if len(values) > 2 {
		f.Value = &resources.FrontStaff
		basemsg.ApplyCardinal(&f, values[2])
	} else {
		resources.FrontStaff = ""
	}
	if len(values) > 3 {
		f.Value = &resources.SupportStaff
		basemsg.ApplyCardinal(&f, values[3])
	} else {
		resources.SupportStaff = ""
	}
	if len(values) > 4 {
		f.Value = &resources.Providers
		basemsg.ApplyCardinal(&f, strings.Join(values[4:], " "))
	} else {
		resources.Providers = ""
	}
}
func resourcesValid(field *basemsg.Field, resources *ResourceCounts) string {
	if resources.Chairs == "" && resources.VacantChairs == "" && resources.FrontStaff == "" && resources.SupportStaff == "" && resources.Providers == "" {
		return ""
	}
	if !common.PIFOCardinalNumberRE.MatchString(resources.Chairs) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(resources.VacantChairs) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(resources.FrontStaff) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(resources.SupportStaff) {
		goto INVALID
	}
	if !common.PIFOCardinalNumberRE.MatchString(resources.Providers) {
		goto INVALID
	}
	return ""
INVALID:
	return fmt.Sprintf("The %q field does not contain a valid value.  It should contain five numbers separated by spaces.", field.Label)
}
