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
	const fieldCount = 130
	var f = AHFacStat{BaseMessage: basemsg.BaseMessage{
		MessageType: &Type,
		PDFBase:     pdfBase,
		Form:        version,
	}}
	f.BaseMessage.FSubject = &f.FacilityName
	f.BaseMessage.FReportType = &f.ReportType
	f.BaseMessage.FBody = &f.Summary
	f.Fields = make([]*basemsg.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &baseform.DefaultPDFMaps)
	f.Fields = append(f.Fields,
		basemsg.NewStaticPDFContentField(&basemsg.Field{
			PDFMap: basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
				return []basemsg.PDFField{{Name: "Form Type", Value: "Allied Health Facility Status"}}
			}),
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  basemsg.Choices{"Update", "Complete"},
			Presence: basemsg.Required,
			PIFOTag:  "19.",
			EditHelp: `This indicates whether the form should "Update" the previous status report for the facility, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:    "Facility Name",
			Value:    &f.FacilityName,
			Presence: basemsg.Required,
			PIFOTag:  "20.",
			PDFMap: basemsg.PDFMapFunc(func(f *basemsg.Field) []basemsg.PDFField {
				return []basemsg.PDFField{
					{Name: "Form Topic", Value: *f.Value},
					{Name: "FACILTY TYPE TIME DATE FACILITY NAME", Value: *f.Value},
				}
			}),
			EditWidth: 52,
			EditHelp:  `This is the name of the facility whose status is being reported.  It is required.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Facility Type",
			Value:     &f.FacilityType,
			Presence:  f.requiredForComplete,
			PIFOTag:   "21.",
			PDFMap:    basemsg.PDFName("facility type"),
			EditWidth: 21,
			EditHelp:  `This is the type of the facility, such as Skilled Nursing, Home Health, Dialysis, Community Health Center, Surgical Center, or Hospice.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewDateWithTimeField(&basemsg.Field{
			Label:    "Date",
			Value:    &f.Date,
			Presence: basemsg.Required,
			PIFOTag:  "22d.",
			PDFMap:   basemsg.PDFName("date"),
		}),
		basemsg.NewTimeWithDateField(&basemsg.Field{
			Label:    "Time",
			Value:    &f.Time,
			Presence: basemsg.Required,
			PIFOTag:  "22t.",
			PDFMap:   basemsg.PDFName("time"),
		}),
		basemsg.NewDateTimeField(&basemsg.Field{
			Label:    "Date/Time",
			Presence: basemsg.Required,
			EditHelp: `This is the date and time of the status report, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.Date, &f.Time),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Contact Name",
			Value:     &f.ContactName,
			Presence:  f.requiredForComplete,
			PIFOTag:   "23.",
			PDFMap:    basemsg.PDFName("Contact Name"),
			EditWidth: 52,
			EditHelp:  `This is the name of the person to be contacted with questions about this report.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Contact Phone",
			Value:     &f.ContactPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "23p.",
			PDFMap:    basemsg.PDFName("Phone"),
			EditWidth: 19,
			EditHelp:  `This is the phone number of the person to be contacted with questions about this report.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Contact Fax",
			Value:     &f.ContactFax,
			PIFOTag:   "23f.",
			PDFMap:    basemsg.PDFName("Fax"),
			EditWidth: 22,
			EditHelp:  `This is the fax number of the person to be contacted with questions about this report.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Other Contact",
			Value:     &f.OtherContact,
			PDFMap:    basemsg.PDFName("Other Phone Fax Cell Phone Radio"),
			PIFOTag:   "24.",
			EditWidth: 53,
			EditHelp:  `This is additional contact information for the person to be contacted with questions about this report.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:    "Incident Name",
			Value:    &f.IncidentName,
			Presence: f.requiredForComplete,
			PIFOTag:  "25.",
			PDFMap: basemsg.PDFMapFunc(func(*basemsg.Field) []basemsg.PDFField {
				return []basemsg.PDFField{{
					Name:  "Incident Name and Date",
					Value: common.SmartJoin(f.IncidentName, f.IncidentDate, " "),
				}}
			}),
			EditWidth: 42,
			EditHelp:  `This is the assigned incident name of the incident for which this report is filed.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewDateWithoutTimeField(&basemsg.Field{
			Label:    "Incident Date",
			Value:    &f.IncidentDate,
			Presence: f.requiredForComplete,
			PIFOTag:  "25d.",
			EditHelp: `This is the date of the incident for which this report is filed, in MM/DD/YYYY format.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:    "Facility Status",
			Value:    &f.FacilityStatus,
			Choices:  basemsg.Choices{"Fully Functional", "Limited Services", "Impaired/Closed"},
			Presence: f.requiredForComplete,
			PIFOTag:  "35.",
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
			EditHelp: `This indicates the status of the facility.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "EOC Main Contact Number",
			Value:     &f.EOCPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "27p.",
			PDFMap:    basemsg.PDFName("EOC MAIN CONTACT NUMBER"),
			EditWidth: 19,
			EditHelp:  `This is the main phone number for the facility's Emergency Operations Center (EOC).  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "EOC Main Contact Fax",
			Value:     &f.EOCFax,
			PIFOTag:   "27f.",
			PDFMap:    basemsg.PDFName("EOC MAIN CONTACT FAX"),
			EditWidth: 20,
			EditHelp:  `This is the max fax number for the facility's Emergency Operations Center (EOC).`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Liaison Officer Name",
			Value:     &f.LiaisonName,
			Presence:  f.requiredForComplete,
			PIFOTag:   "28.",
			PDFMap:    basemsg.PDFName("NAME LIAISON TO PUBLIC HEALTHMEDICAL HEALTH BRANCH"),
			EditWidth: 17,
			EditHelp:  `This is the name of the facility's liaison to the Public Health or Medical Health Branch.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Liaison Contact Number",
			Value:     &f.LiaisonPhone,
			PIFOTag:   "28p.",
			PDFMap:    basemsg.PDFName("CONTACT NUMBER"),
			EditWidth: 17,
			EditHelp:  `This is the phone number of the facility's liaison to the Public Health or Medical Health Branch.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Info Officer Name",
			Value:     &f.InfoOfficerName,
			PIFOTag:   "29.",
			PDFMap:    basemsg.PDFName("INFORMATION OFFICER NAME"),
			EditWidth: 17,
			EditHelp:  `This is the name of the facility's information officer.`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Info Officer Contact Number",
			Value:     &f.InfoOfficerPhone,
			PIFOTag:   "29p.",
			PDFMap:    basemsg.PDFName("CONTACT NUMBER_2"),
			EditWidth: 17,
			EditHelp:  `This is the phone number of the facility's information officer.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Info Officer Contact Email",
			Value:     &f.InfoOfficerEmail,
			PIFOTag:   "29e.",
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT EMAIL"),
			EditWidth: 17,
			EditHelp:  `This is the email address of the facility's information officer.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Not Active Contact Name",
			Value:     &f.ClosedContactName,
			Presence:  f.requiredForComplete,
			PIFOTag:   "30.",
			PDFMap:    basemsg.PDFName("IF EOC IS NOT ACTIVATED WHO SHOULD BE CONTACTED FOR QUESTIONSREQUESTS"),
			EditWidth: 17,
			EditHelp:  `This is the name of the person to be contacted with questions or requests when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewPhoneNumberField(&basemsg.Field{
			Label:     "Facility Contact Phone",
			Value:     &f.ClosedContactPhone,
			Presence:  f.requiredForComplete,
			PIFOTag:   "30p.",
			PDFMap:    basemsg.PDFName("CONTACT NUMBER_3"),
			EditWidth: 17,
			EditHelp:  `This is the phone number of the person to be contacted when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:     "Facility Contact Email",
			Value:     &f.ClosedContactEmail,
			Presence:  f.requiredForComplete,
			PIFOTag:   "30e.",
			Compare:   common.CompareExact,
			PDFMap:    basemsg.PDFName("CONTACT EMAIL_2"),
			EditWidth: 17,
			EditHelp:  `This is the email address of the person to be contacted when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:     "Patients To Evacuate",
			Value:     &f.PatientsToEvacuate,
			PIFOTag:   "31a.",
			PDFMap:    basemsg.PDFName("TOTALPATIENTS TO EVACUATE"),
			EditWidth: 17,
			EditHelp:  `This is the number of patients who need evacuation.`,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:     "Patients Injured - Minor",
			Value:     &f.PatientsInjuredMinor,
			PIFOTag:   "31b.",
			PDFMap:    basemsg.PDFName("TOTALPATIENTS  INJURED  MINOR"),
			EditWidth: 17,
			EditHelp:  `This is the number of patients with minor injuries.`,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:     "Patients Transferred",
			Value:     &f.PatientsTransferred,
			PIFOTag:   "31c.",
			PDFMap:    basemsg.PDFName("TOTALPATIENTS TRANSFERED OUT OF COUNTY"),
			EditWidth: 17,
			EditHelp:  `This is the number of patients who have been transferred out of the county.`,
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:   "Other Patient Care Info",
			Value:   &f.OtherPatientCare,
			PIFOTag: "33.",
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
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:      "Attached Org Chart",
			Value:      &f.AttachOrgChart,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOTag:    "26a.",
			PDFMap:     basemsg.PDFName("YesNoNHICSICS ORGANIZATION CHART"),
			TableValue: basemsg.TableOmit,
			EditHelp:   `This indicates whether an NHICS/ICS organization chart is attached to the status report.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:      "Attached Resource Requests",
			Value:      &f.AttachRR,
			Choices:    basemsg.Choices{"Yes", "No"},
			Presence:   f.requiredForComplete,
			PIFOTag:    "26b.",
			PDFMap:     basemsg.PDFName("YesNoDEOC9A RESOURCE REQUEST FORMS"),
			TableValue: basemsg.TableOmit,
			EditHelp:   `This indicates whether DEOC-9A resource request forms are attached to the status report.  It is required when the "Report Type" is "Complete".`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:      "Attached Status Report",
			Value:      &f.AttachStatus,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOTag:    "26c.",
			PDFMap:     basemsg.PDFName("YesNoNHICSICS STATUS REPORT FORM  STANDARD"),
			TableValue: basemsg.TableOmit,
			EditHelp:   `This indicates whether an NHICS/ICS standard status report form is attached to this status report.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:      "Attached Incident Action Plan",
			Value:      &f.AttachActionPlan,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOTag:    "26d.",
			PDFMap:     basemsg.PDFName("YesNoNHICSICS INCIDENT ACTION PLAN"),
			TableValue: basemsg.TableOmit,
			EditHelp:   `This indicates whether an NHICS/ICS incident action plan is attached to the status report.`,
		}),
		basemsg.NewRestrictedField(&basemsg.Field{
			Label:      "Attached Phone Directory",
			Value:      &f.AttachDirectory,
			Choices:    basemsg.Choices{"Yes", "No"},
			PIFOTag:    "26e.",
			PDFMap:     basemsg.PDFName("YesNoPHONECOMMUNICATIONS DIRECTORY"),
			TableValue: basemsg.TableOmit,
			EditHelp:   `This indicates whether a phone/communications directory is attached to the status report.`,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewMultilineField(&basemsg.Field{
			Label:     "General Summary",
			Value:     &f.Summary,
			PIFOTag:   "34.",
			PDFMap:    basemsg.PDFName("GENERAL SUMMARY OF SITUATIONCONDITIONSRow1"),
			EditWidth: 41,
			EditHelp:  `This is a general summary of the situation and conditions at the facility.`,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Skilled Nursing Beds: Staffed M",
			Value:      &f.SkilledNursingBeds.StaffedM,
			PIFOTag:    "40a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MSKILLED NURSING"),
			TableValue: basemsg.TableOmit,
		}),
	)
	var first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Skilled Nursing Beds: Staffed F",
			Value:   &f.SkilledNursingBeds.StaffedF,
			PIFOTag: "40b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFSKILLED NURSING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Skilled Nursing Beds: Vacant M",
			Value:   &f.SkilledNursingBeds.VacantM,
			PIFOTag: "40c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMSKILLED NURSING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Skilled Nursing Beds: Vacant F",
			Value:   &f.SkilledNursingBeds.VacantF,
			PIFOTag: "40d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFSKILLED NURSING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Skilled Nursing Beds: Surge",
			Value:   &f.SkilledNursingBeds.Surge,
			PIFOTag: "40e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge SKILLED NURSING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Assisted Living Beds: Staffed M",
			Value:      &f.AssistedLivingBeds.StaffedM,
			PIFOTag:    "41a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MASSISTED LIVING"),
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Assisted Living Beds: Staffed F",
			Value:   &f.AssistedLivingBeds.StaffedF,
			PIFOTag: "41b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFASSISTED LIVING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Assisted Living Beds: Vacant M",
			Value:   &f.AssistedLivingBeds.VacantM,
			PIFOTag: "41c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMASSISTED LIVING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Assisted Living Beds: Vacant F",
			Value:   &f.AssistedLivingBeds.VacantF,
			PIFOTag: "41d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFASSISTED LIVING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Assisted Living Beds: Surge",
			Value:   &f.AssistedLivingBeds.Surge,
			PIFOTag: "41e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge ASSISTED LIVING"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Sub-Acute Beds: Staffed M",
			Value:      &f.SubAcuteBeds.StaffedM,
			PIFOTag:    "42a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Sub-Acute Beds: Staffed F",
			Value:   &f.SubAcuteBeds.StaffedF,
			PIFOTag: "42b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Sub-Acute Beds: Vacant M",
			Value:   &f.SubAcuteBeds.VacantM,
			PIFOTag: "42c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Sub-Acute Beds: Vacant F",
			Value:   &f.SubAcuteBeds.VacantF,
			PIFOTag: "42d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFSURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Sub-Acute Beds: Surge",
			Value:   &f.SubAcuteBeds.Surge,
			PIFOTag: "42e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge SURGICAL BEDS"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Alzheimers Beds: Staffed M",
			Value:      &f.AlzheimersBeds.StaffedM,
			PIFOTag:    "43a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Alzheimers Beds: Staffed F",
			Value:   &f.AlzheimersBeds.StaffedF,
			PIFOTag: "43b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Alzheimers Beds: Vacant M",
			Value:   &f.AlzheimersBeds.VacantM,
			PIFOTag: "43c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Alzheimers Beds: Vacant F",
			Value:   &f.AlzheimersBeds.VacantF,
			PIFOTag: "43d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFSUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Alzheimers Beds: Surge",
			Value:   &f.AlzheimersBeds.Surge,
			PIFOTag: "43e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge SUBACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Ped Sub-Acute Beds: Staffed M",
			Value:      &f.PedSubAcuteBeds.StaffedM,
			PIFOTag:    "44a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Staffed F",
			Value:   &f.PedSubAcuteBeds.StaffedF,
			PIFOTag: "44b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Vacant M",
			Value:   &f.PedSubAcuteBeds.VacantM,
			PIFOTag: "44c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Vacant F",
			Value:   &f.PedSubAcuteBeds.VacantF,
			PIFOTag: "44d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Ped Sub-Acute Beds: Surge",
			Value:   &f.PedSubAcuteBeds.Surge,
			PIFOTag: "44e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge ALZEIMERSDIMENTIA"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Psychiatric Beds: Staffed M",
			Value:      &f.PsychiatricBeds.StaffedM,
			PIFOTag:    "45a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Psychiatric Beds: Staffed F",
			Value:   &f.PsychiatricBeds.StaffedF,
			PIFOTag: "45b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Psychiatric Beds: Vacant M",
			Value:   &f.PsychiatricBeds.VacantM,
			PIFOTag: "45c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Psychiatric Beds: Vacant F",
			Value:   &f.PsychiatricBeds.VacantF,
			PIFOTag: "45d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFPEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Psychiatric Beds: Surge",
			Value:   &f.PsychiatricBeds.Surge,
			PIFOTag: "45e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge PEDIATRICSUB ACUTE"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewTextField(&basemsg.Field{
			Label:   "Other Care Beds Type",
			Value:   &f.OtherCareBedsType,
			PIFOTag: "46.",
			Compare: common.CompareText,
			// The PDF doesn't have a fillable field for this, so
			// its value is added to the Other Patient Care Info
			// field, above.
			EditWidth: 17,
			EditHelp:  `This is the other type of beds available at the facility, if any.`,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Other Care Beds: Staffed M",
			Value:      &f.OtherCareBeds.StaffedM,
			PIFOTag:    "46a.",
			PDFMap:     basemsg.PDFName("Staffed Bed MPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Other Care Beds: Staffed F",
			Value:   &f.OtherCareBeds.StaffedF,
			PIFOTag: "46b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Staffed BedFPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Other Care Beds: Vacant M",
			Value:   &f.OtherCareBeds.VacantM,
			PIFOTag: "46c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedsMPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Other Care Beds: Vacant F",
			Value:   &f.OtherCareBeds.VacantF,
			PIFOTag: "46d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Vacant BedFPSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Other Care Beds: Surge",
			Value:   &f.OtherCareBeds.Surge,
			PIFOTag: "46e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("Surge PSYCHIATRIC"), // name is wrong in PDF
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
			EditSkip: func(*basemsg.Field) bool {
				return f.OtherCareBedsType == ""
			},
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Dialysis: Chairs",
			Value:      &f.DialysisResources.Chairs,
			PIFOTag:    "50a.",
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSDIALYSIS"),
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Dialysis: Vacant Chairs",
			Value:   &f.DialysisResources.VacantChairs,
			PIFOTag: "50b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMDIALYSIS"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Dialysis: Front Staff",
			Value:   &f.DialysisResources.FrontStaff,
			PIFOTag: "50c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFDIALYSIS"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Dialysis: Support Staff",
			Value:   &f.DialysisResources.SupportStaff,
			PIFOTag: "50d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFDIALYSIS"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Dialysis: Providers",
			Value:   &f.DialysisResources.Providers,
			PIFOTag: "50e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("PROVIDER STAFFDIALYSIS"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Surgical: Chairs",
			Value:      &f.SurgicalResources.Chairs,
			PIFOTag:    "51a.",
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSSURGICAL"),
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Surgical: Vacant Chairs",
			Value:   &f.SurgicalResources.VacantChairs,
			PIFOTag: "51b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMSURGICAL"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Surgical: Front Staff",
			Value:   &f.SurgicalResources.FrontStaff,
			PIFOTag: "51c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFSURGICAL"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Surgical: Support Staff",
			Value:   &f.SurgicalResources.SupportStaff,
			PIFOTag: "51d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFSURGICAL"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Surgical: Providers",
			Value:   &f.SurgicalResources.Providers,
			PIFOTag: "51e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("PROVIDER STAFFSURGICAL"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Clinic: Chairs",
			Value:      &f.ClinicResources.Chairs,
			PIFOTag:    "52a.",
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSCLINIC"),
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Clinic: Vacant Chairs",
			Value:   &f.ClinicResources.VacantChairs,
			PIFOTag: "52b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMCLINIC"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Clinic: Front Staff",
			Value:   &f.ClinicResources.FrontStaff,
			PIFOTag: "52c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFCLINIC"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Clinic: Support Staff",
			Value:   &f.ClinicResources.SupportStaff,
			PIFOTag: "52d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFCLINIC"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Clinic: Providers",
			Value:   &f.ClinicResources.Providers,
			PIFOTag: "52e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("PROVIDER STAFFCLINIC"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Home Health: Chairs",
			Value:      &f.HomeHealthResources.Chairs,
			PIFOTag:    "53a.",
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSHOMEHEALTH"),
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Home Health: Vacant Chairs",
			Value:   &f.HomeHealthResources.VacantChairs,
			PIFOTag: "53b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMHOMEHEALTH"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Home Health: Front Staff",
			Value:   &f.HomeHealthResources.FrontStaff,
			PIFOTag: "53c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFHOMEHEALTH"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Home Health: Support Staff",
			Value:   &f.HomeHealthResources.SupportStaff,
			PIFOTag: "53d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFHOMEHEALTH"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Home Health: Providers",
			Value:   &f.HomeHealthResources.Providers,
			PIFOTag: "53e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("PROVIDER STAFFHOMEHEALTH"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:      "Adult Day Ctr: Chairs",
			Value:      &f.AdultDayCtrResources.Chairs,
			PIFOTag:    "54a.",
			PDFMap:     basemsg.PDFName("CHAIRS ROOMSADULT DAY CENTER"),
			TableValue: basemsg.TableOmit,
		}),
	)
	first = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Adult Day Ctr: Vacant Chairs",
			Value:   &f.AdultDayCtrResources.VacantChairs,
			PIFOTag: "54b.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("VANCANT CHAIRS ROOMADULT DAY CENTER"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Adult Day Ctr: Front Staff",
			Value:   &f.AdultDayCtrResources.FrontStaff,
			PIFOTag: "54c.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("FRONT DESK STAFFADULT DAY CENTER"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Adult Day Ctr: Support Staff",
			Value:   &f.AdultDayCtrResources.SupportStaff,
			PIFOTag: "54d.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("MEDICAL SUPPORT STAFFADULT DAY CENTER"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewCardinalNumberField(&basemsg.Field{
			Label:   "Adult Day Ctr: Providers",
			Value:   &f.AdultDayCtrResources.Providers,
			PIFOTag: "54e.",
			PIFOValid: func(field *basemsg.Field) string {
				return allOrNone(first, field)
			},
			PDFMap:     basemsg.PDFName("PROVIDER STAFFADULT DAY CENTER"),
			TableValue: basemsg.TableOmit,
		}),
		basemsg.NewAggregatorField(&basemsg.Field{
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
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &baseform.DefaultPDFMaps)
	if len(f.Fields) > fieldCount {
		panic("update AHFacStat fieldCount")
	}
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
	return basemsg.ValidCardinalNumber(current)
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
		basemsg.ApplyCardinalNumber(&f, values[0])
	} else {
		beds.StaffedM = ""
	}
	if len(values) > 1 {
		f.Value = &beds.StaffedF
		basemsg.ApplyCardinalNumber(&f, values[1])
	} else {
		beds.StaffedF = ""
	}
	if len(values) > 2 {
		f.Value = &beds.VacantM
		basemsg.ApplyCardinalNumber(&f, values[2])
	} else {
		beds.VacantM = ""
	}
	if len(values) > 3 {
		f.Value = &beds.VacantF
		basemsg.ApplyCardinalNumber(&f, values[3])
	} else {
		beds.VacantF = ""
	}
	if len(values) > 4 {
		f.Value = &beds.Surge
		basemsg.ApplyCardinalNumber(&f, strings.Join(values[4:], " "))
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
		basemsg.ApplyCardinalNumber(&f, values[0])
	} else {
		resources.Chairs = ""
	}
	if len(values) > 1 {
		f.Value = &resources.VacantChairs
		basemsg.ApplyCardinalNumber(&f, values[1])
	} else {
		resources.VacantChairs = ""
	}
	if len(values) > 2 {
		f.Value = &resources.FrontStaff
		basemsg.ApplyCardinalNumber(&f, values[2])
	} else {
		resources.FrontStaff = ""
	}
	if len(values) > 3 {
		f.Value = &resources.SupportStaff
		basemsg.ApplyCardinalNumber(&f, values[3])
	} else {
		resources.SupportStaff = ""
	}
	if len(values) > 4 {
		f.Value = &resources.Providers
		basemsg.ApplyCardinalNumber(&f, strings.Join(values[4:], " "))
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
