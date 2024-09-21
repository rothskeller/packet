// Package ahfacstat defines the Allied Health Facility Status Form message
// type.
package ahfacstat

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/xscmsg/baseform"
)

// Type24 is the type definition for an allied health facility status form.
var Type24 = message.Type{
	Tag:        "AHFacStat",
	HTML:       "form-allied-health-facility-status.html",
	Version:    "2.4",
	Name:       "allied health facility status form",
	Article:    "an",
	FieldOrder: Type26.FieldOrder,
}

func init() {
	message.Register(&Type24, decode24, nil)
}

// AHFacStat24 holds an allied health facility status form.
type AHFacStat24 struct {
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
	SkilledNursingBeds   BedCounts24
	AssistedLivingBeds   BedCounts24
	SubAcuteBeds         BedCounts24
	AlzheimersBeds       BedCounts24
	PedSubAcuteBeds      BedCounts24
	PsychiatricBeds      BedCounts24
	OtherCareBedsType    string
	OtherCareBeds        BedCounts24
	DialysisResources    ResourceCounts24
	SurgicalResources    ResourceCounts24
	ClinicResources      ResourceCounts24
	HomeHealthResources  ResourceCounts24
	AdultDayCtrResources ResourceCounts24
}
type BedCounts24 struct {
	StaffedM string
	StaffedF string
	VacantM  string
	VacantF  string
	Surge    string
}
type ResourceCounts24 struct {
	Chairs       string
	VacantChairs string
	FrontStaff   string
	SupportStaff string
	Providers    string
}

func make24() (f *AHFacStat24) {
	const fieldCount = 129
	f = &AHFacStat24{BaseMessage: message.BaseMessage{Type: &Type24}}
	f.BaseMessage.FSubject = &f.FacilityName
	f.BaseMessage.FBody = &f.Summary
	f.Fields = make([]*message.Field, 0, fieldCount)
	f.BaseForm.AddHeaderFields(&f.BaseMessage, nil)
	f.Fields = append(f.Fields,
		message.NewRestrictedField(&message.Field{
			Label:    "Report Type",
			Value:    &f.ReportType,
			Choices:  message.Choices{"Update", "Complete"},
			Presence: message.Required,
			PIFOTag:  "19.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Facility Name",
			Value:    &f.FacilityName,
			Presence: message.Required,
			PIFOTag:  "20.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Facility Type",
			Value:    &f.FacilityType,
			Presence: f.requiredForComplete,
			PIFOTag:  "21.",
		}),
		message.NewDateField(true, &message.Field{
			Label:    "Date",
			Value:    &f.Date,
			Presence: message.Required,
			PIFOTag:  "22d.",
		}),
		message.NewTimeField(true, &message.Field{
			Label:    "Time",
			Value:    &f.Time,
			Presence: message.Required,
			PIFOTag:  "22t.",
		}),
		message.NewDateTimeField(&message.Field{
			Label:    "Date/Time",
			Presence: message.Required,
		}, &f.Date, &f.Time),
		message.NewTextField(&message.Field{
			Label:    "Contact Name",
			Value:    &f.ContactName,
			Presence: f.requiredForComplete,
			PIFOTag:  "23.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:    "Contact Phone",
			Value:    &f.ContactPhone,
			Presence: f.requiredForComplete,
			PIFOTag:  "23p.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "Contact Fax",
			Value:   &f.ContactFax,
			PIFOTag: "23f.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Other Contact",
			Value:   &f.OtherContact,
			PIFOTag: "24.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Incident Name",
			Value:    &f.IncidentName,
			Presence: f.requiredForComplete,
			PIFOTag:  "25.",
		}),
		message.NewDateField(false, &message.Field{
			Label:    "Incident Date",
			Value:    &f.IncidentDate,
			Presence: f.requiredForComplete,
			PIFOTag:  "25d.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:    "Facility Status",
			Value:    &f.FacilityStatus,
			Choices:  message.Choices{"Fully Functional", "Limited Services", "Impaired/Closed"},
			Presence: f.requiredForComplete,
			PIFOTag:  "35.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:    "EOC Main Contact Number",
			Value:    &f.EOCPhone,
			Presence: f.requiredForComplete,
			PIFOTag:  "27p.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "EOC Main Contact Fax",
			Value:   &f.EOCFax,
			PIFOTag: "27f.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Liaison Officer Name",
			Value:    &f.LiaisonName,
			Presence: f.requiredForComplete,
			PIFOTag:  "28.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "Liaison Contact Number",
			Value:   &f.LiaisonPhone,
			PIFOTag: "28p.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Info Officer Name",
			Value:   &f.InfoOfficerName,
			PIFOTag: "29.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:   "Info Officer Contact Number",
			Value:   &f.InfoOfficerPhone,
			PIFOTag: "29p.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Info Officer Contact Email",
			Value:   &f.InfoOfficerEmail,
			PIFOTag: "29e.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Not Active Contact Name",
			Value:    &f.ClosedContactName,
			Presence: f.requiredForComplete,
			PIFOTag:  "30.",
		}),
		message.NewPhoneNumberField(&message.Field{
			Label:    "Facility Contact Phone",
			Value:    &f.ClosedContactPhone,
			Presence: f.requiredForComplete,
			PIFOTag:  "30p.",
		}),
		message.NewTextField(&message.Field{
			Label:    "Facility Contact Email",
			Value:    &f.ClosedContactEmail,
			Presence: f.requiredForComplete,
			PIFOTag:  "30e.",
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Patients To Evacuate",
			Value:   &f.PatientsToEvacuate,
			PIFOTag: "31a.",
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Patients Injured - Minor",
			Value:   &f.PatientsInjuredMinor,
			PIFOTag: "31b.",
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Patients Transferred",
			Value:   &f.PatientsTransferred,
			PIFOTag: "31c.",
		}),
		message.NewTextField(&message.Field{
			Label:   "Other Patient Care Info",
			Value:   &f.OtherPatientCare,
			PIFOTag: "33.",
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Attached Org Chart",
			Value:      &f.AttachOrgChart,
			Choices:    message.Choices{"Yes", "No"},
			PIFOTag:    "26a.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Attached Resource Requests",
			Value:      &f.AttachRR,
			Choices:    message.Choices{"Yes", "No"},
			Presence:   f.requiredForComplete,
			PIFOTag:    "26b.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Attached Status Report",
			Value:      &f.AttachStatus,
			Choices:    message.Choices{"Yes", "No"},
			PIFOTag:    "26c.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Attached Incident Action Plan",
			Value:      &f.AttachActionPlan,
			Choices:    message.Choices{"Yes", "No"},
			PIFOTag:    "26d.",
			TableValue: message.TableOmit,
		}),
		message.NewRestrictedField(&message.Field{
			Label:      "Attached Phone Directory",
			Value:      &f.AttachDirectory,
			Choices:    message.Choices{"Yes", "No"},
			PIFOTag:    "26e.",
			TableValue: message.TableOmit,
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
			Label:   "General Summary",
			Value:   &f.Summary,
			PIFOTag: "34.",
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Skilled Nursing Beds: Staffed M",
			Value:      &f.SkilledNursingBeds.StaffedM,
			PIFOTag:    "40a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstSkilledNursing = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Staffed F",
			Value:   &f.SkilledNursingBeds.StaffedF,
			PIFOTag: "40b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSkilledNursing, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Vacant M",
			Value:   &f.SkilledNursingBeds.VacantM,
			PIFOTag: "40c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSkilledNursing, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Vacant F",
			Value:   &f.SkilledNursingBeds.VacantF,
			PIFOTag: "40d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSkilledNursing, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Skilled Nursing Beds: Surge",
			Value:   &f.SkilledNursingBeds.Surge,
			PIFOTag: "40e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSkilledNursing, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Skilled Nursing Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.SkilledNursingBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Assisted Living Beds: Staffed M",
			Value:      &f.AssistedLivingBeds.StaffedM,
			PIFOTag:    "41a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstAssistedLiving = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Staffed F",
			Value:   &f.AssistedLivingBeds.StaffedF,
			PIFOTag: "41b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAssistedLiving, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Vacant M",
			Value:   &f.AssistedLivingBeds.VacantM,
			PIFOTag: "41c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAssistedLiving, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Vacant F",
			Value:   &f.AssistedLivingBeds.VacantF,
			PIFOTag: "41d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAssistedLiving, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Assisted Living Beds: Surge",
			Value:   &f.AssistedLivingBeds.Surge,
			PIFOTag: "41e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAssistedLiving, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Assisted Living Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.AssistedLivingBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Sub-Acute Beds: Staffed M",
			Value:      &f.SubAcuteBeds.StaffedM,
			PIFOTag:    "42a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstSubAcute = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Staffed F",
			Value:   &f.SubAcuteBeds.StaffedF,
			PIFOTag: "42b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Vacant M",
			Value:   &f.SubAcuteBeds.VacantM,
			PIFOTag: "42c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Vacant F",
			Value:   &f.SubAcuteBeds.VacantF,
			PIFOTag: "42d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Sub-Acute Beds: Surge",
			Value:   &f.SubAcuteBeds.Surge,
			PIFOTag: "42e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Sub-Acute Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.SubAcuteBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Alzheimers Beds: Staffed M",
			Value:      &f.AlzheimersBeds.StaffedM,
			PIFOTag:    "43a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstAlzheimers = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Staffed F",
			Value:   &f.AlzheimersBeds.StaffedF,
			PIFOTag: "43b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAlzheimers, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Vacant M",
			Value:   &f.AlzheimersBeds.VacantM,
			PIFOTag: "43c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAlzheimers, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Vacant F",
			Value:   &f.AlzheimersBeds.VacantF,
			PIFOTag: "43d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAlzheimers, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Alzheimers Beds: Surge",
			Value:   &f.AlzheimersBeds.Surge,
			PIFOTag: "43e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAlzheimers, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Alzheimers/Dementia Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.AlzheimersBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Ped Sub-Acute Beds: Staffed M",
			Value:      &f.PedSubAcuteBeds.StaffedM,
			PIFOTag:    "44a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstPedSubAcute = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Staffed F",
			Value:   &f.PedSubAcuteBeds.StaffedF,
			PIFOTag: "44b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPedSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Vacant M",
			Value:   &f.PedSubAcuteBeds.VacantM,
			PIFOTag: "44c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPedSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Vacant F",
			Value:   &f.PedSubAcuteBeds.VacantF,
			PIFOTag: "44d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPedSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Ped Sub-Acute Beds: Surge",
			Value:   &f.PedSubAcuteBeds.Surge,
			PIFOTag: "44e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPedSubAcute, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Pediatric Sub-Acute Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.PedSubAcuteBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Psychiatric Beds: Staffed M",
			Value:      &f.PsychiatricBeds.StaffedM,
			PIFOTag:    "45a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstPsychiatric = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Staffed F",
			Value:   &f.PsychiatricBeds.StaffedF,
			PIFOTag: "45b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPsychiatric, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Vacant M",
			Value:   &f.PsychiatricBeds.VacantM,
			PIFOTag: "45c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPsychiatric, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Vacant F",
			Value:   &f.PsychiatricBeds.VacantF,
			PIFOTag: "45d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPsychiatric, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Psychiatric Beds: Surge",
			Value:   &f.PsychiatricBeds.Surge,
			PIFOTag: "45e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstPsychiatric, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Psychiatric Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.PsychiatricBeds)
			},
		}),
		message.NewTextField(&message.Field{
			Label:   "Other Care Beds Type",
			Value:   &f.OtherCareBedsType,
			PIFOTag: "46.",
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Other Care Beds: Staffed M",
			Value:      &f.OtherCareBeds.StaffedM,
			PIFOTag:    "46a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstOtherCare = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Staffed F",
			Value:   &f.OtherCareBeds.StaffedF,
			PIFOTag: "46b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstOtherCare, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Vacant M",
			Value:   &f.OtherCareBeds.VacantM,
			PIFOTag: "46c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstOtherCare, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Vacant F",
			Value:   &f.OtherCareBeds.VacantF,
			PIFOTag: "46d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstOtherCare, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Other Care Beds: Surge",
			Value:   &f.OtherCareBeds.Surge,
			PIFOTag: "46e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstOtherCare, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Other Care Beds",
			TableValue: func(*message.Field) string {
				return bedsTableValue24(&f.OtherCareBeds)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Dialysis: Chairs",
			Value:      &f.DialysisResources.Chairs,
			PIFOTag:    "50a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstDialysis = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Vacant Chairs",
			Value:   &f.DialysisResources.VacantChairs,
			PIFOTag: "50b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstDialysis, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Front Staff",
			Value:   &f.DialysisResources.FrontStaff,
			PIFOTag: "50c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstDialysis, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Support Staff",
			Value:   &f.DialysisResources.SupportStaff,
			PIFOTag: "50d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstDialysis, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Dialysis: Providers",
			Value:   &f.DialysisResources.Providers,
			PIFOTag: "50e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstDialysis, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Dialysis Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue24(&f.DialysisResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Surgical: Chairs",
			Value:      &f.SurgicalResources.Chairs,
			PIFOTag:    "51a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstSurgical = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Vacant Chairs",
			Value:   &f.SurgicalResources.VacantChairs,
			PIFOTag: "51b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSurgical, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Front Staff",
			Value:   &f.SurgicalResources.FrontStaff,
			PIFOTag: "51c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSurgical, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Support Staff",
			Value:   &f.SurgicalResources.SupportStaff,
			PIFOTag: "51d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSurgical, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Surgical: Providers",
			Value:   &f.SurgicalResources.Providers,
			PIFOTag: "51e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstSurgical, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Surgical Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue24(&f.SurgicalResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Clinic: Chairs",
			Value:      &f.ClinicResources.Chairs,
			PIFOTag:    "52a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstClinic = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Vacant Chairs",
			Value:   &f.ClinicResources.VacantChairs,
			PIFOTag: "52b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstClinic, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Front Staff",
			Value:   &f.ClinicResources.FrontStaff,
			PIFOTag: "52c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstClinic, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Support Staff",
			Value:   &f.ClinicResources.SupportStaff,
			PIFOTag: "52d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstClinic, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Clinic: Providers",
			Value:   &f.ClinicResources.Providers,
			PIFOTag: "52e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstClinic, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Clinic Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue24(&f.ClinicResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Home Health: Chairs",
			Value:      &f.HomeHealthResources.Chairs,
			PIFOTag:    "53a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstHomeHealth = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Vacant Chairs",
			Value:   &f.HomeHealthResources.VacantChairs,
			PIFOTag: "53b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstHomeHealth, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Front Staff",
			Value:   &f.HomeHealthResources.FrontStaff,
			PIFOTag: "53c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstHomeHealth, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Support Staff",
			Value:   &f.HomeHealthResources.SupportStaff,
			PIFOTag: "53d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstHomeHealth, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Home Health: Providers",
			Value:   &f.HomeHealthResources.Providers,
			PIFOTag: "53e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstHomeHealth, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Home Health Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue24(&f.HomeHealthResources)
			},
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:      "Adult Day Ctr: Chairs",
			Value:      &f.AdultDayCtrResources.Chairs,
			PIFOTag:    "54a.",
			TableValue: message.TableOmit,
		}),
	)
	var firstAdultDayCtr = f.Fields[len(f.Fields)-1]
	f.Fields = append(f.Fields,
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Vacant Chairs",
			Value:   &f.AdultDayCtrResources.VacantChairs,
			PIFOTag: "54b.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAdultDayCtr, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Front Staff",
			Value:   &f.AdultDayCtrResources.FrontStaff,
			PIFOTag: "54c.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAdultDayCtr, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Support Staff",
			Value:   &f.AdultDayCtrResources.SupportStaff,
			PIFOTag: "54d.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAdultDayCtr, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewCardinalNumberField(&message.Field{
			Label:   "Adult Day Ctr: Providers",
			Value:   &f.AdultDayCtrResources.Providers,
			PIFOTag: "54e.",
			PIFOValid: func(field *message.Field) string {
				return allOrNone24(firstAdultDayCtr, field)
			},
			TableValue: message.TableOmit,
		}),
		message.NewAggregatorField(&message.Field{
			Label: "Adult Day Center Resources",
			TableValue: func(*message.Field) string {
				return resourcesTableValue24(&f.AdultDayCtrResources)
			},
		}),
	)
	f.BaseForm.AddFooterFields(&f.BaseMessage, nil)
	if len(f.Fields) > fieldCount {
		panic("update AHFacStat24 fieldCount")
	}
	return f
}

func (f *AHFacStat24) requiredForComplete() (message.Presence, string) {
	if f.ReportType == "Complete" {
		return message.PresenceRequired, `the "Report Type" is "Complete"`
	}
	return message.PresenceOptional, ""
}

func decode24(_, _ string, form *message.PIFOForm, _ int) message.Message {
	if form == nil || form.HTMLIdent != Type24.HTML || form.FormVersion != Type24.Version {
		return nil
	}
	var df = make24()
	message.DecodeForm(form, df)
	return df
}

func allOrNone24(first, current *message.Field) string {
	if *first.Value == "" && *current.Value != "" {
		return fmt.Sprintf("The %q field must not have a value unless %q has a value.  (Either all fields on the row must have a value, or none.)", current.Label, first.Label)
	}
	if *first.Value != "" && *current.Value == "" {
		return fmt.Sprintf("The %q field is required when %q has a value.  (Either all fields on the row must have a value, or none.)", current.Label, first.Label)
	}
	return message.ValidCardinalNumber(current)
}

func bedsTableValue24(beds *BedCounts24) string {
	if beds.StaffedM == "" && beds.StaffedF == "" && beds.VacantM == "" && beds.VacantF == "" && beds.Surge == "" {
		return ""
	}
	return fmt.Sprintf("%3s %3s %3s %3s %3s", beds.StaffedM, beds.StaffedF, beds.VacantM, beds.VacantF, beds.Surge)
}

func resourcesTableValue24(resources *ResourceCounts24) string {
	if resources.Chairs == "" && resources.VacantChairs == "" && resources.FrontStaff == "" && resources.SupportStaff == "" && resources.Providers == "" {
		return ""
	}
	return fmt.Sprintf("%3s %3s %3s %3s %3s", resources.Chairs, resources.VacantChairs, resources.FrontStaff, resources.SupportStaff, resources.Providers)
}

func (f *AHFacStat24) Compare(actual message.Message) (int, int, []*message.CompareField) {
	return f.convertTo26().Compare(actual)
}

func (f *AHFacStat24) RenderPDF(env *envelope.Envelope, filename string) error {
	return f.convertTo26().RenderPDF(env, filename)
}

func (f *AHFacStat24) convertTo26() (c *AHFacStat26) {
	c = make26()
	c.CopyHeaderFields(&f.BaseForm)
	c.ReportType = f.ReportType
	c.FacilityName = f.FacilityName
	c.FacilityType = f.FacilityType
	c.Date = f.Date
	c.Time = f.Time
	c.ContactName = f.ContactName
	c.ContactPhone = f.ContactPhone
	c.ContactFax = f.ContactFax
	c.OtherContact = f.OtherContact
	c.IncidentName = f.IncidentName
	c.IncidentDate = f.IncidentDate
	c.FacilityStatus = f.FacilityStatus
	c.EOCPhone = f.EOCPhone
	c.EOCFax = f.EOCFax
	c.LiaisonName = f.LiaisonName
	c.LiaisonPhone = f.LiaisonPhone
	c.InfoOfficerName = f.InfoOfficerName
	c.InfoOfficerPhone = f.InfoOfficerPhone
	c.InfoOfficerEmail = f.InfoOfficerEmail
	c.ClosedContactName = f.ClosedContactName
	c.ClosedContactPhone = f.ClosedContactPhone
	c.ClosedContactEmail = f.ClosedContactEmail
	c.PatientsToEvacuate = f.PatientsToEvacuate
	c.PatientsInjuredMinor = f.PatientsInjuredMinor
	c.PatientsTransferred = f.PatientsTransferred
	c.OtherPatientCare = f.OtherPatientCare
	c.AttachOrgChart = f.AttachOrgChart
	c.AttachRR = f.AttachRR
	c.AttachStatus = f.AttachStatus
	c.AttachActionPlan = f.AttachActionPlan
	c.AttachDirectory = f.AttachDirectory
	c.Summary = f.Summary
	c.SkilledNursingBeds = f.SkilledNursingBeds.convertTo26()
	c.AssistedLivingBeds = f.AssistedLivingBeds.convertTo26()
	c.SubAcuteBeds = f.SubAcuteBeds.convertTo26()
	c.AlzheimersBeds = f.AlzheimersBeds.convertTo26()
	c.PedSubAcuteBeds = f.PedSubAcuteBeds.convertTo26()
	c.PsychiatricBeds = f.PsychiatricBeds.convertTo26()
	c.OtherCareBedsType = f.OtherCareBedsType
	c.OtherCareBeds = f.OtherCareBeds.convertTo26()
	c.DialysisResources = f.DialysisResources.convertTo26()
	c.SurgicalResources = f.SurgicalResources.convertTo26()
	c.ClinicResources = f.ClinicResources.convertTo26()
	c.HomeHealthResources = f.HomeHealthResources.convertTo26()
	c.AdultDayCtrResources = f.AdultDayCtrResources.convertTo26()
	c.CopyFooterFields(&f.BaseForm)
	return c
}

func (bc *BedCounts24) convertTo26() (c BedCounts26) {
	c.StaffedM = bc.StaffedM
	c.StaffedF = bc.StaffedF
	c.VacantM = bc.VacantM
	c.VacantF = bc.VacantF
	c.Surge = bc.Surge
	return c
}

func (rc *ResourceCounts24) convertTo26() (c ResourceCounts26) {
	c.Chairs = rc.Chairs
	c.VacantChairs = rc.VacantChairs
	c.FrontStaff = rc.FrontStaff
	c.SupportStaff = rc.SupportStaff
	c.Providers = rc.Providers
	return c
}
