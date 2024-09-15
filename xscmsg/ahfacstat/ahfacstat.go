// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"fmt"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
	"golang.org/x/exp/slices"
)

// Type is the type definition for an allied health facility status form.
var Type = message.Type{
	Tag:         "AHFacStat",
	Name:        "allied health facility status form",
	Article:     "an",
	PDFRenderV2: true,
}

func init() {
	Type.Create = New
	Type.Decode = decode
}

// versions is the list of supported versions.  The first one is used when
// creating new forms.
var versions = []*message.FormVersion{
	{HTML: "form-allied-health-facility-status.html", Version: "2.6", Tag: "AHFacStat", FieldOrder: fieldOrder},
	{HTML: "form-allied-health-facility-status.html", Version: "2.4", Tag: "AHFacStat", FieldOrder: fieldOrder},
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
	message.BaseMessage
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

func create(version *message.FormVersion) message.Message {
	const fieldCount = 130
	var f = AHFacStat{BaseMessage: message.BaseMessage{
		Type: &Type,
		Form: version,
	}}
	var pdf = baseform.RoutingSlipPDFRenderers
	pdf.OriginMsgIDR = message.PDFMultiRenderer{
		pdf.OriginMsgIDR,
		&message.PDFTextRenderer{Page: 2, X: 492, Y: 36, W: 100, H: 17, Style: message.PDFTextStyle{HAlign: "right"}},
	}
	f.BaseMessage.FSubject = &f.FacilityName
	f.BaseMessage.FBody = &f.Summary
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, &pdf)
	f.Fields = append(f.Fields,
		message.NewStaticPDFContentField(&message.Field{
			PDFRenderer: &message.PDFStaticTextRenderer{
				Page: 1, X: 119, Y: 225, H: 17,
				Text: "Allied Health Facility Status",
			},
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  message.Choices{"Update", "Complete"},
			Presence: message.Required,
			PIFOTag:  "19.",
			EditHelp: `This indicates whether the form should "Update" the previous status report for the facility, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:    "Facility Name",
			Value:    &f.FacilityName,
			Presence: message.Required,
			PIFOTag:  "20.",
			PDFRenderer: message.PDFMultiRenderer{
				&message.PDFTextRenderer{Page: 1, X: 351, Y: 225, W: 222, H: 17, Style: message.PDFTextStyle{VAlign: "baseline"}},
				&message.PDFTextRenderer{Page: 2, X: 20, Y: 118, W: 320, H: 24, Style: message.PDFTextStyle{VAlign: "baseline"}},
			},
			EditWidth: 52,
			EditHelp:  `This is the name of the facility whose status is being reported.  It is required.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Facility Type",
			Value:       &f.FacilityType,
			Presence:    f.requiredForComplete,
			PIFOTag:     "21.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 344, Y: 118, W: 135, H: 24, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   21,
			EditHelp:    `This is the type of the facility, such as Skilled Nursing, Home Health, Dialysis, Community Health Center, Surgical Center, or Hospice.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewDateField(true, &message.Field{
			Label:       "Date",
			Value:       &f.Date,
			Presence:    message.Required,
			PIFOTag:     "22d.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 482, Y: 118, W: 67, H: 24, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date of the status report, in MM/DD/YYYY format.  It is required.`,
		}),
		message.NewTimeField(true, &message.Field{
			Label:       "Time",
			Value:       &f.Time,
			Presence:    message.Required,
			PIFOTag:     "22t.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 552, Y: 118, W: 37, H: 24, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the time of the status report, in HH:MM format (24-hour clock).  It is required.`,
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			Presence: message.Required,
			EditHelp: `This is the date and time of the status report, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.`,
		}, &f.Date, &f.Time),
		message.NewTextField(&message.Field{
			Label:       "Contact Name",
			Value:       &f.ContactName,
			Presence:    f.requiredForComplete,
			PIFOTag:     "23.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 20, Y: 152, W: 314, H: 23, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   52,
			EditHelp:    `This is the name of the person to be contacted with questions about this report.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Contact Phone",
			Value:       &f.ContactPhone,
			Presence:    f.requiredForComplete,
			PIFOTag:     "23p.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 339, Y: 152, W: 112, H: 23, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   19,
			EditHelp:    `This is the phone number of the person to be contacted with questions about this report.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Contact Fax",
			Value:       &f.ContactFax,
			PIFOTag:     "23f.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 152, W: 134, H: 23, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   22,
			EditHelp:    `This is the fax number of the person to be contacted with questions about this report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Other Contact",
			Value:       &f.OtherContact,
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 20, Y: 186, W: 314, H: 21, Style: message.PDFTextStyle{VAlign: "baseline"}},
			PIFOTag:     "24.",
			EditWidth:   53,
			EditHelp:    `This is additional contact information for the person to be contacted with questions about this report.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Incident Name",
			Value:       &f.IncidentName,
			Presence:    f.requiredForComplete,
			PIFOTag:     "25.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 339, Y: 186, W: 179, H: 21, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   42,
			EditHelp:    `This is the assigned incident name of the incident for which this report is filed.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewDateField(false, &message.Field{
			Label:       "Incident Date",
			Value:       &f.IncidentDate,
			Presence:    f.requiredForComplete,
			PIFOTag:     "25d.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 521, Y: 186, W: 68, H: 21, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditHelp:    `This is the date of the incident for which this report is filed, in MM/DD/YYYY format.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Facility Status",
			Value:    &f.FacilityStatus,
			Choices:  message.Choices{"Fully Functional", "Limited Services", "Impaired/Closed"},
			Presence: f.requiredForComplete,
			PIFOTag:  "35.",
			PDFRenderer: &message.PDFCheckRenderer{
				Page: 2, W: 12, H: 12,
				Points: map[string][]float64{
					"Fully Functional": {307, 227},
					"Limited Services": {307, 243},
					"Impaired/Closed":  {307, 260},
				},
			},
			EditHelp: `This indicates the status of the facility.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "EOC Main Contact Number",
			Value:       &f.EOCPhone,
			Presence:    f.requiredForComplete,
			PIFOTag:     "27p.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 294, W: 101, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   19,
			EditHelp:    `This is the main phone number for the facility's Emergency Operations Center (EOC).  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "EOC Main Contact Fax",
			Value:       &f.EOCFax,
			PIFOTag:     "27f.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 309, W: 101, H: 13},
			EditWidth:   20,
			EditHelp:    `This is the max fax number for the facility's Emergency Operations Center (EOC).`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Liaison Officer Name",
			Value:       &f.LiaisonName,
			Presence:    f.requiredForComplete,
			PIFOTag:     "28.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 324, W: 101, H: 20},
			EditWidth:   17,
			EditHelp:    `This is the name of the facility's liaison to the Public Health or Medical Health Branch.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Liaison Contact Number",
			Value:       &f.LiaisonPhone,
			PIFOTag:     "28p.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 347, W: 101, H: 13},
			EditWidth:   17,
			EditHelp:    `This is the phone number of the facility's liaison to the Public Health or Medical Health Branch.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Info Officer Name",
			Value:       &f.InfoOfficerName,
			PIFOTag:     "29.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 362, W: 101, H: 14},
			EditWidth:   17,
			EditHelp:    `This is the name of the facility's information officer.`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Info Officer Contact Number",
			Value:       &f.InfoOfficerPhone,
			PIFOTag:     "29p.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 378, W: 101, H: 13},
			EditWidth:   17,
			EditHelp:    `This is the phone number of the facility's information officer.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Info Officer Contact Email",
			Value:       &f.InfoOfficerEmail,
			PIFOTag:     "29e.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 393, W: 101, H: 14},
			EditWidth:   17,
			EditHelp:    `This is the email address of the facility's information officer.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Not Active Contact Name",
			Value:       &f.ClosedContactName,
			Presence:    f.requiredForComplete,
			PIFOTag:     "30.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 409, W: 101, H: 26},
			EditWidth:   17,
			EditHelp:    `This is the name of the person to be contacted with questions or requests when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:       "Facility Contact Phone",
			Value:       &f.ClosedContactPhone,
			Presence:    f.requiredForComplete,
			PIFOTag:     "30p.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 437, W: 101, H: 14, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   17,
			EditHelp:    `This is the phone number of the person to be contacted when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Facility Contact Email",
			Value:       &f.ClosedContactEmail,
			Presence:    f.requiredForComplete,
			PIFOTag:     "30e.",
			Compare:     message.CompareExact,
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 453, W: 101, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   17,
			EditHelp:    `This is the email address of the person to be contacted when the facility's EOC is not activated.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Patients To Evacuate",
			Value:       &f.PatientsToEvacuate,
			PIFOTag:     "31a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 484, W: 101, H: 14, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   17,
			EditHelp:    `This is the number of patients who need evacuation.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Patients Injured - Minor",
			Value:       &f.PatientsInjuredMinor,
			PIFOTag:     "31b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 501, W: 101, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   17,
			EditHelp:    `This is the number of patients with minor injuries.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Patients Transferred",
			Value:       &f.PatientsTransferred,
			PIFOTag:     "31c.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 516, W: 101, H: 14, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   17,
			EditHelp:    `This is the number of patients who have been transferred out of the county.`,
		}),
		message.NewTextField(&message.Field{
			Label:       "Other Patient Care Info",
			Value:       &f.OtherPatientCare,
			PIFOTag:     "33.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 233, Y: 532, W: 101, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   27,
			EditHelp:    `This field contains other patient care information as needed.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Attached Org Chart",
			Value:       &f.AttachOrgChart,
			Choices:     message.Choices{"Yes", "No"},
			PIFOTag:     "26a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 521, Y: 227, W: 68, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether an NHICS/ICS organization chart is attached to the status report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Attached Resource Requests",
			Value:       &f.AttachRR,
			Choices:     message.Choices{"Yes", "No"},
			Presence:    f.requiredForComplete,
			PIFOTag:     "26b.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 521, Y: 243, W: 68, H: 12, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether DEOC-9A resource request forms are attached to the status report.  It is required when the "Report Type" is "Complete".`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Attached Status Report",
			Value:       &f.AttachStatus,
			Choices:     message.Choices{"Yes", "No"},
			PIFOTag:     "26c.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 521, Y: 257, W: 68, H: 18, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether an NHICS/ICS standard status report form is attached to this status report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Attached Incident Action Plan",
			Value:       &f.AttachActionPlan,
			Choices:     message.Choices{"Yes", "No"},
			PIFOTag:     "26d.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 521, Y: 277, W: 68, H: 15, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether an NHICS/ICS incident action plan is attached to the status report.`,
		}),
		message.NewRestrictedField(&message.Field{
			Label:       "Attached Phone Directory",
			Value:       &f.AttachDirectory,
			Choices:     message.Choices{"Yes", "No"},
			PIFOTag:     "26e.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 521, Y: 294, W: 68, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This indicates whether a phone/communications directory is attached to the status report.`,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Attachments",
			TableValue: func(*message.Field) string {
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
		message.NewMultilineField(&message.Field{
			Label:       "General Summary",
			Value:       &f.Summary,
			PIFOTag:     "34.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 339, Y: 326, W: 250, H: 79, Style: message.PDFTextStyle{VAlign: "top"}},
			EditWidth:   41,
			EditHelp:    `This is a general summary of the situation and conditions at the facility.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Skilled Nursing Beds: Staffed M",
			Value:       &f.SkilledNursingBeds.StaffedM,
			PIFOTag:     "40a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 437, W: 25, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male skilled nursing beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstSkilledNursing = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Staffed F",
			Value:   &f.SkilledNursingBeds.StaffedF,
			PIFOTag: "40b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSkilledNursing, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 437, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female skilled nursing beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Vacant M",
			Value:   &f.SkilledNursingBeds.VacantM,
			PIFOTag: "40c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSkilledNursing, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 437, W: 22, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male skilled nursing beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Vacant F",
			Value:   &f.SkilledNursingBeds.VacantF,
			PIFOTag: "40d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSkilledNursing, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 437, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female skilled nursing beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Surge",
			Value:   &f.SkilledNursingBeds.Surge,
			PIFOTag: "40e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSkilledNursing, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 437, W: 29, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity skilled nursing beds at the facility (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Skilled Nursing Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.SkilledNursingBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of skilled nursing beds at the facility.  Enter five numbers separated by spaces: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.SkilledNursingBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.SkilledNursingBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.SkilledNursingBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Assisted Living Beds: Staffed M",
			Value:       &f.AssistedLivingBeds.StaffedM,
			PIFOTag:     "41a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 453, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male assisted living beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstAssistedLiving = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Staffed F",
			Value:   &f.AssistedLivingBeds.StaffedF,
			PIFOTag: "41b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAssistedLiving, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 453, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female assisted living beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Vacant M",
			Value:   &f.AssistedLivingBeds.VacantM,
			PIFOTag: "41c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAssistedLiving, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 453, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male assisted living beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Vacant F",
			Value:   &f.AssistedLivingBeds.VacantF,
			PIFOTag: "41d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAssistedLiving, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 453, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female assisted living beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Surge",
			Value:   &f.AssistedLivingBeds.Surge,
			PIFOTag: "41e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAssistedLiving, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 453, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity assisted living beds at the facility (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Assisted Living Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.AssistedLivingBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of assisted living beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.AssistedLivingBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.AssistedLivingBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.AssistedLivingBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Sub-Acute Beds: Staffed M",
			Value:       &f.SubAcuteBeds.StaffedM,
			PIFOTag:     "42a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 468, W: 25, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstSubAcute = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Staffed F",
			Value:   &f.SubAcuteBeds.StaffedF,
			PIFOTag: "42b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 468, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Vacant M",
			Value:   &f.SubAcuteBeds.VacantM,
			PIFOTag: "42c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 468, W: 22, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Vacant F",
			Value:   &f.SubAcuteBeds.VacantF,
			PIFOTag: "42d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 468, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Surge",
			Value:   &f.SubAcuteBeds.Surge,
			PIFOTag: "42e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 468, W: 29, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity sub-acute beds at the facility (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Sub-Acute Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.SubAcuteBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of sub-acute beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.SubAcuteBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.SubAcuteBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.SubAcuteBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Alzheimers Beds: Staffed M",
			Value:       &f.AlzheimersBeds.StaffedM,
			PIFOTag:     "43a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 484, W: 25, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male Alzheimers/dementia beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstAlzheimers = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Staffed F",
			Value:   &f.AlzheimersBeds.StaffedF,
			PIFOTag: "43b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAlzheimers, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 484, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female Alzheimers/dementia beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Vacant M",
			Value:   &f.AlzheimersBeds.VacantM,
			PIFOTag: "43c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAlzheimers, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 484, W: 22, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male Alzheimers/dementia beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Vacant F",
			Value:   &f.AlzheimersBeds.VacantF,
			PIFOTag: "43d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAlzheimers, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 484, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female Alzheimers/dementia beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Surge",
			Value:   &f.AlzheimersBeds.Surge,
			PIFOTag: "43e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAlzheimers, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 484, W: 29, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity Alzheimers/dementia beds at the facility (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Alzheimers/Dementia Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.AlzheimersBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of Alzheimers/dementia beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.AlzheimersBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.AlzheimersBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.AlzheimersBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Ped Sub-Acute Beds: Staffed M",
			Value:       &f.PedSubAcuteBeds.StaffedM,
			PIFOTag:     "44a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 501, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male pediatric sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstPedSubAcute = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Staffed F",
			Value:   &f.PedSubAcuteBeds.StaffedF,
			PIFOTag: "44b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPedSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 501, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female pediatric sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Vacant M",
			Value:   &f.PedSubAcuteBeds.VacantM,
			PIFOTag: "44c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPedSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 501, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male pediatric sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Vacant F",
			Value:   &f.PedSubAcuteBeds.VacantF,
			PIFOTag: "44d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPedSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 501, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female pediatric sub-acute beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Surge",
			Value:   &f.PedSubAcuteBeds.Surge,
			PIFOTag: "44e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPedSubAcute, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 501, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity pediatric sub-acute beds at the facility (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Pediatric Sub-Acute Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.PedSubAcuteBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of pediatric sub-acute beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.PedSubAcuteBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.PedSubAcuteBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.PedSubAcuteBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Psychiatric Beds: Staffed M",
			Value:       &f.PsychiatricBeds.StaffedM,
			PIFOTag:     "45a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 516, W: 25, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male psychiatric beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstPsychiatric = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Staffed F",
			Value:   &f.PsychiatricBeds.StaffedF,
			PIFOTag: "45b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPsychiatric, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 516, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female psychiatric beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Vacant M",
			Value:   &f.PsychiatricBeds.VacantM,
			PIFOTag: "45c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPsychiatric, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 516, W: 22, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male psychiatric beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Vacant F",
			Value:   &f.PsychiatricBeds.VacantF,
			PIFOTag: "45d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPsychiatric, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 516, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female psychiatric beds at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Surge",
			Value:   &f.PsychiatricBeds.Surge,
			PIFOTag: "45e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstPsychiatric, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 516, W: 29, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity psychiatric beds at the facility (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Psychiatric Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.PsychiatricBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of psychiatric beds at the facility.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.PsychiatricBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.PsychiatricBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.PsychiatricBeds)
			},
		}),
		message.NewTextField(&message.Field{
			Label:   "Other Care Beds Type",
			Value:   &f.OtherCareBedsType,
			PIFOTag: "46.",
			Compare: message.CompareText,
			// The PDF doesn't have a fillable field for this, so
			// its value is added to the Other Patient Care Info
			// field, above.
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 339, Y: 532, W: 112, H: 13, Style: message.PDFTextStyle{VAlign: "baseline"}},
			EditWidth:   17,
			EditHelp:    `This is the other type of beds available at the facility, if any.`,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Other Care Beds: Staffed M",
			Value:       &f.OtherCareBeds.StaffedM,
			PIFOTag:     "46a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 532, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed male beds at the facility of the named other type.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstOtherCare = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Staffed F",
			Value:   &f.OtherCareBeds.StaffedF,
			PIFOTag: "46b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstOtherCare, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 532, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of staffed female beds at the facility of the named other type.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Vacant M",
			Value:   &f.OtherCareBeds.VacantM,
			PIFOTag: "46c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstOtherCare, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 532, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant male beds at the facility of the named other type.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Vacant F",
			Value:   &f.OtherCareBeds.VacantF,
			PIFOTag: "46d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstOtherCare, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 532, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant female beds at the facility of the named other type.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Surge",
			Value:   &f.OtherCareBeds.Surge,
			PIFOTag: "46e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstOtherCare, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 532, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surge capacity beds at the facility of the named other type (over and above the vacant ones).`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Other Care Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue(&f.OtherCareBeds)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of beds at the facility of the named other type.  Enter five numbers separated by spaces or commas: the numbers of staffed beds for male patients, staffed beds for female patients, vacant beds for male patients, vacant beds for female patients, and surge beds (over and above the vacant ones).`,
			EditHint:  "M, F, V.M, V.F, Surge",
			EditValue: func(*message.Field) string {
				return bedsValue(&f.OtherCareBeds)
			},
			EditApply: func(_ *message.Field, value string) {
				bedsApply(&f.OtherCareBeds, value)
			},
			EditValid: func(field *message.Field) string {
				return bedsValid(field, &f.OtherCareBeds)
			},
			EditSkip: func(*message.Field) bool {
				return f.OtherCareBedsType == ""
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Dialysis: Chairs",
			Value:       &f.DialysisResources.Chairs,
			PIFOTag:     "50a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 584, W: 25, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of dialysis chairs at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstDialysis = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Vacant Chairs",
			Value:   &f.DialysisResources.VacantChairs,
			PIFOTag: "50b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstDialysis, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 584, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant dialysis chairs at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Front Staff",
			Value:   &f.DialysisResources.FrontStaff,
			PIFOTag: "50c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstDialysis, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 584, W: 22, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of dialysis front desk staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Support Staff",
			Value:   &f.DialysisResources.SupportStaff,
			PIFOTag: "50d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstDialysis, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 584, W: 21, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of dialysis support staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Providers",
			Value:   &f.DialysisResources.Providers,
			PIFOTag: "50e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstDialysis, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 584, W: 29, H: 14, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of dialysis provider staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Dialysis Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue(&f.DialysisResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of dialysis resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*message.Field) string {
				return resourcesValue(&f.DialysisResources)
			},
			EditApply: func(_ *message.Field, value string) {
				resourcesApply(&f.DialysisResources, value)
			},
			EditValid: func(field *message.Field) string {
				return resourcesValid(field, &f.DialysisResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Surgical: Chairs",
			Value:       &f.SurgicalResources.Chairs,
			PIFOTag:     "51a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 600, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surgical rooms at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstSurgical = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Vacant Chairs",
			Value:   &f.SurgicalResources.VacantChairs,
			PIFOTag: "51b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSurgical, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 600, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant surgical rooms at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Front Staff",
			Value:   &f.SurgicalResources.FrontStaff,
			PIFOTag: "51c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSurgical, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 600, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surgical front desk staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Support Staff",
			Value:   &f.SurgicalResources.SupportStaff,
			PIFOTag: "51d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSurgical, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 600, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surgical support staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Providers",
			Value:   &f.SurgicalResources.Providers,
			PIFOTag: "51e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstSurgical, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 600, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of surgical provider staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Surgical Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue(&f.SurgicalResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of surgical resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*message.Field) string {
				return resourcesValue(&f.SurgicalResources)
			},
			EditApply: func(_ *message.Field, value string) {
				resourcesApply(&f.SurgicalResources, value)
			},
			EditValid: func(field *message.Field) string {
				return resourcesValid(field, &f.SurgicalResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Clinic: Chairs",
			Value:       &f.ClinicResources.Chairs,
			PIFOTag:     "52a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 616, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of clinic rooms at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstClinic = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Vacant Chairs",
			Value:   &f.ClinicResources.VacantChairs,
			PIFOTag: "52b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstClinic, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 616, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant clinic rooms at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Front Staff",
			Value:   &f.ClinicResources.FrontStaff,
			PIFOTag: "52c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstClinic, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 616, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of clinic front desk staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Support Staff",
			Value:   &f.ClinicResources.SupportStaff,
			PIFOTag: "52d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstClinic, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 616, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of clinic support staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Providers",
			Value:   &f.ClinicResources.Providers,
			PIFOTag: "52e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstClinic, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 616, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of clinic provider staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Clinic Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue(&f.ClinicResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of clinic resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*message.Field) string {
				return resourcesValue(&f.ClinicResources)
			},
			EditApply: func(_ *message.Field, value string) {
				resourcesApply(&f.ClinicResources, value)
			},
			EditValid: func(field *message.Field) string {
				return resourcesValid(field, &f.ClinicResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Home Health: Chairs",
			Value:       &f.HomeHealthResources.Chairs,
			PIFOTag:     "53a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 632, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of home health rooms at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstHomeHealth = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Vacant Chairs",
			Value:   &f.HomeHealthResources.VacantChairs,
			PIFOTag: "53b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstHomeHealth, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 632, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant home health rooms at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Front Staff",
			Value:   &f.HomeHealthResources.FrontStaff,
			PIFOTag: "53c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstHomeHealth, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 632, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of home health front desk staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Support Staff",
			Value:   &f.HomeHealthResources.SupportStaff,
			PIFOTag: "53d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstHomeHealth, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 632, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of home health support staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Providers",
			Value:   &f.HomeHealthResources.Providers,
			PIFOTag: "53e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstHomeHealth, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 632, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of home health provider staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Home Health Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue(&f.HomeHealthResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of home health resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*message.Field) string {
				return resourcesValue(&f.HomeHealthResources)
			},
			EditApply: func(_ *message.Field, value string) {
				resourcesApply(&f.HomeHealthResources, value)
			},
			EditValid: func(field *message.Field) string {
				return resourcesValid(field, &f.HomeHealthResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:       "Adult Day Ctr: Chairs",
			Value:       &f.AdultDayCtrResources.Chairs,
			PIFOTag:     "54a.",
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 455, Y: 647, W: 25, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of adult day center chairs at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
	)
	var firstAdultDayCtr = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Vacant Chairs",
			Value:   &f.AdultDayCtrResources.VacantChairs,
			PIFOTag: "54b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAdultDayCtr, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 484, Y: 647, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of vacant adult day center chairs at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Front Staff",
			Value:   &f.AdultDayCtrResources.FrontStaff,
			PIFOTag: "54c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAdultDayCtr, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 509, Y: 647, W: 22, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of adult day center front desk staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Support Staff",
			Value:   &f.AdultDayCtrResources.SupportStaff,
			PIFOTag: "54d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAdultDayCtr, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 535, Y: 647, W: 21, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of adult day center support staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Providers",
			Value:   &f.AdultDayCtrResources.Providers,
			PIFOTag: "54e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone(firstAdultDayCtr, field)
			},
			PDFRenderer: &message.PDFTextRenderer{Page: 2, X: 560, Y: 647, W: 29, H: 13, Style: message.PDFTextStyle{HAlign: "right", VAlign: "baseline"}},
			TableValue:  message.TableOmit,
			EditHelp:    `This is the number of adult day center provider staff at the facility.`,
			EditSkip:    message.EditSkipAlways,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Adult Day Center Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue(&f.AdultDayCtrResources)
			},
			EditWidth: 20,
			EditHelp:  `This is the number of adult day center resources at the facility.  Enter five numbers separated by spaces or commas: the numbers of chairs or rooms, vacant chairs or rooms, front desk staff, medical support staff, and providers.`,
			EditHint:  "Ch, V.Ch, FDS, MSS, Prov.",
			EditValue: func(*message.Field) string {
				return resourcesValue(&f.AdultDayCtrResources)
			},
			EditApply: func(_ *message.Field, value string) {
				resourcesApply(&f.AdultDayCtrResources, value)
			},
			EditValid: func(field *message.Field) string {
				return resourcesValid(field, &f.AdultDayCtrResources)
			},
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, &pdf)
	if len(f.Fields) > fieldCount {
		panic("update AHFacStat fieldCount")
	}
	return &f
}

func (f *AHFacStat) requiredForComplete() (message.Presence, string) {
	if f.ReportType == "Complete" {
		return message.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return message.PresenceOptional, ""
}

func decode(subject, body string) (f *AHFacStat) {
	// Quick check to avoid overhead of creating the form object if it's not
	// our type of form.
	if !strings.Contains(body, "form-allied-health-facility-status.html") {
		return nil
	}
	if df, ok := message.DecodeForm(body, versions, create).(*AHFacStat); ok {
		return df
	} else {
		return nil
	}
}

func allOrNone(first, current *message.Field) string {
	if *first.Value == "" && *current.Value != "" {
		return fmt.Sprintf("The %q field must not have a value unless %q has a value.  (Either all fields on the row must have a value, or none.)", current.Label, first.Label)
	}
	if *first.Value != "" && *current.Value == "" {
		return fmt.Sprintf("The %q field is required when %q has a value.  (Either all fields on the row must have a value, or none.)", current.Label, first.Label)
	}
	return message.ValidCardinalNumber(current)
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
	return strings.Join(slices.DeleteFunc(
		[]string{beds.StaffedM, beds.StaffedF, beds.VacantM, beds.VacantF, beds.Surge},
		func(s string) bool { return s == "" },
	), " ")
}
func bedsApply(beds *BedCounts, value string) {
	var f message.Field
	values := strings.Fields(value)
	if len(values) > 0 {
		f.Value = &beds.StaffedM
		message.ApplyCardinalNumber(&f, values[0])
	} else {
		beds.StaffedM = ""
	}
	if len(values) > 1 {
		f.Value = &beds.StaffedF
		message.ApplyCardinalNumber(&f, values[1])
	} else {
		beds.StaffedF = ""
	}
	if len(values) > 2 {
		f.Value = &beds.VacantM
		message.ApplyCardinalNumber(&f, values[2])
	} else {
		beds.VacantM = ""
	}
	if len(values) > 3 {
		f.Value = &beds.VacantF
		message.ApplyCardinalNumber(&f, values[3])
	} else {
		beds.VacantF = ""
	}
	if len(values) > 4 {
		f.Value = &beds.Surge
		message.ApplyCardinalNumber(&f, strings.Join(values[4:], " "))
	} else {
		beds.Surge = ""
	}
}
func bedsValid(field *message.Field, beds *BedCounts) string {
	if beds.StaffedM == "" && beds.StaffedF == "" && beds.VacantM == "" && beds.VacantF == "" && beds.Surge == "" {
		return ""
	}
	if !message.PIFOCardinalNumberRE.MatchString(beds.StaffedM) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(beds.StaffedF) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(beds.VacantM) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(beds.VacantF) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(beds.Surge) {
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
	return strings.Join(slices.DeleteFunc(
		[]string{resources.Chairs, resources.VacantChairs, resources.FrontStaff, resources.SupportStaff, resources.Providers},
		func(s string) bool { return s == "" },
	), " ")
}
func resourcesApply(resources *ResourceCounts, value string) {
	var f message.Field
	values := strings.Fields(value)
	if len(values) > 0 {
		f.Value = &resources.Chairs
		message.ApplyCardinalNumber(&f, values[0])
	} else {
		resources.Chairs = ""
	}
	if len(values) > 1 {
		f.Value = &resources.VacantChairs
		message.ApplyCardinalNumber(&f, values[1])
	} else {
		resources.VacantChairs = ""
	}
	if len(values) > 2 {
		f.Value = &resources.FrontStaff
		message.ApplyCardinalNumber(&f, values[2])
	} else {
		resources.FrontStaff = ""
	}
	if len(values) > 3 {
		f.Value = &resources.SupportStaff
		message.ApplyCardinalNumber(&f, values[3])
	} else {
		resources.SupportStaff = ""
	}
	if len(values) > 4 {
		f.Value = &resources.Providers
		message.ApplyCardinalNumber(&f, strings.Join(values[4:], " "))
	} else {
		resources.Providers = ""
	}
}
func resourcesValid(field *message.Field, resources *ResourceCounts) string {
	if resources.Chairs == "" && resources.VacantChairs == "" && resources.FrontStaff == "" && resources.SupportStaff == "" && resources.Providers == "" {
		return ""
	}
	if !message.PIFOCardinalNumberRE.MatchString(resources.Chairs) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(resources.VacantChairs) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(resources.FrontStaff) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(resources.SupportStaff) {
		goto INVALID
	}
	if !message.PIFOCardinalNumberRE.MatchString(resources.Providers) {
		goto INVALID
	}
	return ""
INVALID:
	return fmt.Sprintf("The %q field does not contain a valid value.  It should contain five numbers separated by spaces.", field.Label)
}
