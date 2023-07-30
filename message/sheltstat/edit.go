package sheltstat

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
	"golang.org/x/exp/slices"
)

type sheltStatEdit struct {
	common.StdFieldsEdit
	ReportType            message.EditField
	ShelterName           message.EditField
	ShelterType           message.EditField
	ShelterStatus         message.EditField
	ShelterAddress        message.EditField
	ShelterCity           message.EditField
	ShelterState          message.EditField
	ShelterZip            message.EditField
	Latitude              message.EditField
	Longitude             message.EditField
	Capacity              message.EditField
	Occupancy             message.EditField
	MealsServed           message.EditField
	NSSNumber             message.EditField
	PetFriendly           message.EditField
	BasicSafetyInspection message.EditField
	ATC20Inspection       message.EditField
	AvailableServices     message.EditField
	MOU                   message.EditField
	FloorPlan             message.EditField
	ManagedBy             message.EditField
	ManagedByDetail       message.EditField
	PrimaryContact        message.EditField
	PrimaryPhone          message.EditField
	SecondaryContact      message.EditField
	SecondaryPhone        message.EditField
	TacticalCallSign      message.EditField
	RepeaterCallSign      message.EditField
	RepeaterInput         message.EditField
	RepeaterInputTone     message.EditField
	RepeaterOutput        message.EditField
	RepeaterOutputTone    message.EditField
	RepeaterOffset        message.EditField
	Comments              message.EditField
	RemoveFromList        message.EditField
	fields                []*message.EditField
}

// EditFields returns the set of editable fields of the message.
func (f *SheltStat) EditFields() []*message.EditField {
	if f.edit == nil {
		f.edit = &sheltStatEdit{
			StdFieldsEdit: common.StdFieldsEditTemplate,
			ReportType: message.EditField{
				Label:   "Report Type",
				Width:   8,
				Choices: []string{"Update", "Complete"},
				Help:    `This indicates whether the form should "Update" the previous status report for the shelter, or whether it is a "Complete" replacement of the previous report.  This field is required.`,
			},
			ShelterName: message.EditField{
				Label: "Shelter Name",
				Width: 44,
				Help:  `This is the name of the shelter whose status is being reported.  It is required.`,
			},
			ShelterType: message.EditField{
				Label:   "Shelter Type",
				Width:   6,
				Choices: []string{"Type 1", "Type 2", "Type 3", "Type 4"},
				Help:    `This is the shelter type.  It is required when "Report Type" is "Complete".`,
			},
			ShelterStatus: message.EditField{
				Label:   "Shelter Status",
				Width:   6,
				Choices: []string{"Open", "Closed", "Full"},
				Help:    `This indicates the status of the shelter.  It is required when "Report Type" is "Complete".`,
			},
			ShelterAddress: message.EditField{
				Label: "Shelter Address",
				Width: 75,
				Help:  `This is the street address of the shelter.  It is required when "Report Type" is "Complete".`,
			},
			ShelterCity: message.EditField{
				Label:   "Shelter City",
				Width:   30,
				Choices: []string{"Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated"},
				Help:    `This is the name of the city in which the shelter is located.  It is required when "Report Type" is "Complete".`,
			},
			ShelterState: message.EditField{
				Label: "Shelter State",
				Width: 12,
				Help:  `This is the name (or two-letter abbreviation) of the state in which the shelter is located.  It is required when "Report Type" is "Complete".`,
			},
			ShelterZip: message.EditField{
				Label: "Shelter Zip",
				Width: 12,
				Help:  `This is the shelter's ZIP code.  It is required when "Report Type" is "Complete".`,
			},
			Latitude: message.EditField{
				Label: "Latitude",
				Width: 30,
				Help:  `This is the latitude of the shelter location, expressed in fractional degrees.`,
			},
			Longitude: message.EditField{
				Label: "Longitude",
				Width: 29,
				Help:  `This is the longitude of the shelter location, expressed in fractional degrees.`,
			},
			Capacity: message.EditField{
				Label: "Capacity",
				Width: 6,
				Help:  `This is the number of people the shelter can accommodate.  It is required when "Report Type" is "Complete".`,
			},
			Occupancy: message.EditField{
				Label: "Occupancy",
				Width: 6,
				Help:  `This is the number of people currently using the shelter.  It is required when "Report Type" is "Complete".`,
			},
			MealsServed: message.EditField{
				Label: "Meals Served",
				Width: 65,
				Help:  `This is the number and/or description of meals served at the shelter in the last 24 hours.`,
			},
			NSSNumber: message.EditField{
				Label: "NSS Number",
				Width: 65,
				Help:  `This is the NSS number of the shelter.`,
			},
			PetFriendly: message.EditField{
				Label:   "Pet Friendly",
				Width:   3,
				Choices: []string{"Yes", "No"},
				Help:    `This indicates whether the shelter can accept pets.`,
			},
			BasicSafetyInspection: message.EditField{
				Label:   "Basic Safety Inspection",
				Width:   3,
				Choices: []string{"Yes", "No"},
				Help:    `This indicates whether the shelter has had a basic safety inspection.`,
			},
			ATC20Inspection: message.EditField{
				Label:   "ATC-20 Inspection",
				Width:   3,
				Choices: []string{"Yes", "No"},
				Help:    `This indicates whether the shelter has had an ATC-20 inspection.`,
			},
			AvailableServices: message.EditField{
				Label: "Available Services",
				Width: 85, Multiline: true,
				Help: `This is a list of services available at the shelter.`,
			},
			MOU: message.EditField{
				Label: "MOU",
				Width: 64,
				Help:  `This indicates where and how the shelter's Memorandum of Understanding (MOU) was reported.`,
			},
			FloorPlan: message.EditField{
				Label: "Floor Plan",
				Width: 64,
				Help:  `This indicates where and how the shelter's floor plan was reported.`,
			},
			ManagedBy: message.EditField{
				Label:   "Managed By",
				Width:   18,
				Choices: []string{"American Red Cross", "Private", "Community", "Government", "Other"},
				Help:    `This indicates what type of entity is managing the shelter.  It is required when "Report Type" is "Complete".`,
			},
			ManagedByDetail: message.EditField{
				Label: "Managed By Detail",
				Width: 65,
				Help:  `This is additional detail about who is managing the shelter (particularly if "Managed By" is "Other").`,
			},
			PrimaryContact: message.EditField{
				Label: "Primary Contact",
				Width: 65,
				Help:  `This is the name of the primary contact person for the shelter.  It is required when "Report Type" is "Complete".`,
			},
			PrimaryPhone: message.EditField{
				Label: "Primary Phone",
				Width: 65,
				Help:  `This is the phone number of the primary contact person for the shelter.  It is required when "Report Type" is "Complete".`,
			},
			SecondaryContact: message.EditField{
				Label: "Secondary Contact",
				Width: 65,
				Help:  `This is the name of the secondary contact person for the shelter.`,
			},
			SecondaryPhone: message.EditField{
				Label: "Secondary Phone",
				Width: 65,
				Help:  `This is the phone number of the secondary contact person for the shelter.`,
			},
			TacticalCallSign: message.EditField{
				Label: "Tactical Call Sign",
				Width: 29,
				Help:  `This is the tactical call sign assigned to the shelter for amateur radio communications.`,
			},
			RepeaterCallSign: message.EditField{
				Label: "Repeater Call Sign",
				Width: 29,
				Help:  `This is the call sign of the amateur radio repeater that the shelter is monitoring for communications.`,
			},
			RepeaterInput: message.EditField{
				Label: "Repeater Input",
				Width: 20,
				Help:  `This is the input frequency (in MHz) of the amateur radio repeater that the shelter is monitoring for communications.`,
			},
			RepeaterInputTone: message.EditField{
				Label: "Repeater Input Tone",
				Width: 30,
				Help:  `This is the frequency (in MHz) of the subaudible code or tone required by the amateur radio repeater that the shelter is monitoring for communications.`,
			},
			RepeaterOutput: message.EditField{
				Label: "Repeater Output",
				Width: 20,
				Help:  `This is the output frequency (in MHz) of the amateur radio repeater that the shelter is monitoring for communications.`,
			},
			RepeaterOutputTone: message.EditField{
				Label: "Repeater Output Tone",
				Width: 30,
				Help:  `This is the frequency (in MHz) of the subaudible code or tone transmitted by the amateur radio repeater that the shelter is monitoring for communications.`,
			},
			RepeaterOffset: message.EditField{
				Label: "Repeater Offset",
				Width: 15,
				Help:  `This is the offset for the amateur radio repeater that the shelter is monitoring for communications.  It can be either a number in MHz, a "+", or a "-".`,
			},
			Comments: message.EditField{
				Label: "Comments",
				Width: 85, Multiline: true,
				Help: `These are comments regarding the status of the shelter.`,
			},
			RemoveFromList: message.EditField{
				Label:   "Remove From List",
				Width:   3,
				Choices: []string{"Yes", "No"},
				Help:    `This indicates whether the shelter should be removed from the receiver's list of shelters.`,
			},
		}
		// Set the field list slice.
		f.edit.fields = append(f.edit.StdFieldsEdit.EditFields1(),
			&f.edit.ReportType,
			&f.edit.ShelterName,
			&f.edit.ShelterType,
			&f.edit.ShelterStatus,
			&f.edit.ShelterAddress,
			&f.edit.ShelterCity,
			&f.edit.ShelterState,
			&f.edit.ShelterZip,
			&f.edit.Latitude,
			&f.edit.Longitude,
			&f.edit.Capacity,
			&f.edit.Occupancy,
			&f.edit.MealsServed,
			&f.edit.NSSNumber,
			&f.edit.PetFriendly,
			&f.edit.BasicSafetyInspection,
			&f.edit.ATC20Inspection,
			&f.edit.AvailableServices,
			&f.edit.MOU,
			&f.edit.FloorPlan,
			&f.edit.ManagedBy,
			&f.edit.ManagedByDetail,
			&f.edit.PrimaryContact,
			&f.edit.PrimaryPhone,
			&f.edit.SecondaryContact,
			&f.edit.SecondaryPhone,
			&f.edit.TacticalCallSign,
			&f.edit.RepeaterCallSign,
			&f.edit.RepeaterInput,
			&f.edit.RepeaterInputTone,
			&f.edit.RepeaterOutput,
			&f.edit.RepeaterOutputTone,
			&f.edit.RepeaterOffset,
			&f.edit.Comments,
			&f.edit.RemoveFromList,
		)
		f.edit.fields = append(f.edit.fields, f.edit.StdFieldsEdit.EditFields2()...)
		f.toEdit()
		f.validate()
	}
	return f.edit.fields
}

// ApplyEdits applies the revised Values in the EditFields to the
// message.
func (f *SheltStat) ApplyEdits() {
	f.fromEdit()
	f.toEdit()
	f.validate()
}

func (f *SheltStat) fromEdit() {
	f.StdFields.FromEdit(&f.edit.StdFieldsEdit)
	f.ReportType = common.ExpandRestricted(&f.edit.ReportType)
	f.ShelterName = strings.TrimSpace(f.edit.ShelterName.Value)
	f.ShelterType = common.ExpandRestricted(&f.edit.ShelterType)
	f.ShelterStatus = common.ExpandRestricted(&f.edit.ShelterStatus)
	f.ShelterAddress = strings.TrimSpace(f.edit.ShelterAddress.Value)
	f.ShelterCity = common.ExpandRestricted(&f.edit.ShelterCity)
	if f.edit.ShelterCity.Value == "" || slices.Contains(f.edit.ShelterCity.Choices, f.edit.ShelterCity.Value) {
		f.ShelterCityCode = f.edit.ShelterCity.Value
	} else {
		f.ShelterCityCode = "Unincorporated"
	}
	f.ShelterState = strings.TrimSpace(f.edit.ShelterState.Value)
	f.ShelterZip = strings.TrimSpace(f.edit.ShelterZip.Value)
	f.Latitude = common.CleanReal(f.edit.Latitude.Value)
	f.Longitude = common.CleanReal(f.edit.Longitude.Value)
	f.Capacity = common.CleanCardinal(f.edit.Capacity.Value)
	f.Occupancy = common.CleanCardinal(f.edit.Occupancy.Value)
	f.MealsServed = strings.TrimSpace(f.edit.MealsServed.Value)
	f.NSSNumber = strings.TrimSpace(f.edit.NSSNumber.Value)
	f.PetFriendly = common.ExpandRestricted(&f.edit.PetFriendly)
	f.BasicSafetyInspection = common.ExpandRestricted(&f.edit.BasicSafetyInspection)
	f.ATC20Inspection = common.ExpandRestricted(&f.edit.ATC20Inspection)
	f.AvailableServices = strings.TrimSpace(f.edit.AvailableServices.Value)
	f.MOU = strings.TrimSpace(f.edit.MOU.Value)
	f.FloorPlan = strings.TrimSpace(f.edit.FloorPlan.Value)
	f.ManagedBy = common.ExpandRestricted(&f.edit.ManagedBy)
	if f.edit.ManagedBy.Value == "" || slices.Contains(f.edit.ManagedBy.Choices, f.edit.ManagedBy.Value) {
		f.ManagedByCode = f.edit.ManagedBy.Value
	} else {
		f.ManagedByCode = "Other"
	}
	f.ManagedByDetail = strings.TrimSpace(f.edit.ManagedByDetail.Value)
	f.PrimaryContact = strings.TrimSpace(f.edit.PrimaryContact.Value)
	f.PrimaryPhone = common.CleanPhoneNumber(f.edit.PrimaryPhone.Value)
	f.SecondaryContact = strings.TrimSpace(f.edit.SecondaryContact.Value)
	f.SecondaryPhone = common.CleanPhoneNumber(f.edit.SecondaryPhone.Value)
	f.TacticalCallSign = strings.ToUpper(strings.TrimSpace(f.edit.TacticalCallSign.Value))
	f.RepeaterCallSign = strings.ToUpper(strings.TrimSpace(f.edit.RepeaterCallSign.Value))
	f.RepeaterInput = common.CleanReal(strings.TrimSpace(f.edit.RepeaterInput.Value))
	f.RepeaterInputTone = common.CleanReal(strings.TrimSpace(f.edit.RepeaterInputTone.Value))
	f.RepeaterOutput = common.CleanReal(strings.TrimSpace(f.edit.RepeaterOutput.Value))
	f.RepeaterOutputTone = common.CleanReal(strings.TrimSpace(f.edit.RepeaterOutputTone.Value))
	f.RepeaterOffset = common.CleanReal(f.edit.RepeaterOffset.Value)
	f.Comments = strings.TrimSpace(f.edit.Comments.Value)
	f.RemoveFromList = common.ExpandRestricted(&f.edit.RemoveFromList)
}

func (f *SheltStat) toEdit() {
	f.StdFields.ToEdit(&f.edit.StdFieldsEdit)
	f.edit.ReportType.Value = f.ReportType
	f.edit.ShelterName.Value = f.ShelterName
	f.edit.ShelterType.Value = f.ShelterType
	f.edit.ShelterStatus.Value = f.ShelterStatus
	f.edit.ShelterAddress.Value = f.ShelterAddress
	f.edit.ShelterCity.Value = f.ShelterCity
	f.edit.ShelterState.Value = f.ShelterState
	f.edit.ShelterZip.Value = f.ShelterZip
	f.edit.Latitude.Value = f.Latitude
	f.edit.Longitude.Value = f.Longitude
	f.edit.Capacity.Value = f.Capacity
	f.edit.Occupancy.Value = f.Occupancy
	f.edit.MealsServed.Value = f.MealsServed
	f.edit.NSSNumber.Value = f.NSSNumber
	f.edit.PetFriendly.Value = f.PetFriendly
	f.edit.BasicSafetyInspection.Value = f.BasicSafetyInspection
	f.edit.ATC20Inspection.Value = f.ATC20Inspection
	f.edit.AvailableServices.Value = f.AvailableServices
	f.edit.MOU.Value = f.MOU
	f.edit.FloorPlan.Value = f.FloorPlan
	f.edit.ManagedBy.Value = f.ManagedBy
	f.edit.ManagedByDetail.Value = f.ManagedByDetail
	f.edit.PrimaryContact.Value = f.PrimaryContact
	f.edit.PrimaryPhone.Value = f.PrimaryPhone
	f.edit.SecondaryContact.Value = f.SecondaryContact
	f.edit.SecondaryPhone.Value = f.SecondaryPhone
	f.edit.TacticalCallSign.Value = f.TacticalCallSign
	f.edit.RepeaterCallSign.Value = f.RepeaterCallSign
	f.edit.RepeaterInput.Value = f.RepeaterInput
	f.edit.RepeaterInputTone.Value = f.RepeaterInputTone
	f.edit.RepeaterOutput.Value = f.RepeaterOutput
	f.edit.RepeaterOutputTone.Value = f.RepeaterOutputTone
	f.edit.RepeaterOffset.Value = f.RepeaterOffset
	f.edit.Comments.Value = f.Comments
	f.edit.RemoveFromList.Value = f.RemoveFromList
}

func (f *SheltStat) validate() {
	f.edit.StdFieldsEdit.Validate()
	if common.ValidateRequired(&f.edit.ReportType) {
		common.ValidateRestricted(&f.edit.ReportType)
	}
	common.ValidateRequired(&f.edit.ShelterName)
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ShelterType) && f.edit.ShelterType.Value != "" {
		common.ValidateRestricted(&f.edit.ShelterType)
	}
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ShelterStatus) && f.edit.ShelterStatus.Value != "" {
		common.ValidateRestricted(&f.edit.ShelterStatus)
	}
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ShelterAddress)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ShelterCity)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ShelterState)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ShelterZip)
	if f.edit.Latitude.Value != "" && !common.PIFORealNumberRE.MatchString(f.edit.Latitude.Value) {
		f.edit.Latitude.Problem = `The "Latitude" field does not contain a valid number.`
	} else {
		f.edit.Latitude.Problem = ""
	}
	if f.edit.Longitude.Value != "" && f.edit.Latitude.Value == "" {
		f.edit.Longitude.Problem = `The "Longitude" field may not have a value unless "Latitude" also has a value.`
	} else if f.edit.Longitude.Value != "" && !common.PIFORealNumberRE.MatchString(f.edit.Longitude.Value) {
		f.edit.Longitude.Problem = `The "Longitude" field does not contain a valid number.`
	} else if f.edit.Longitude.Value == "" && f.edit.Latitude.Value != "" {
		f.edit.Longitude.Problem = `The "Longitude" field must have a value when "Latitude" has a value.`
	} else {
		f.edit.Longitude.Problem = ""
	}
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.Capacity) {
		if f.edit.Capacity.Value != "" && !common.PIFOCardinalNumberRE.MatchString(f.edit.Capacity.Value) {
			f.edit.Capacity.Problem = `The "Capacity" field does not contain a valid number.`
		}
	}
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.Occupancy) {
		if f.edit.Occupancy.Value != "" && !common.PIFOCardinalNumberRE.MatchString(f.edit.Occupancy.Value) {
			f.edit.Occupancy.Problem = `The "Occupancy" field does not contain a valid number.`
		}
	}
	if f.edit.PetFriendly.Value != "" {
		common.ValidateRestricted(&f.edit.PetFriendly)
	}
	if f.edit.BasicSafetyInspection.Value != "" {
		common.ValidateRestricted(&f.edit.BasicSafetyInspection)
	}
	if f.edit.ATC20Inspection.Value != "" {
		common.ValidateRestricted(&f.edit.ATC20Inspection)
	}
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.ManagedBy)
	validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.PrimaryContact)
	if validateRequiredIfComplete(f.edit.ReportType.Value, &f.edit.PrimaryPhone) {
		if f.edit.PrimaryPhone.Value != "" && !common.PIFOPhoneNumberRE.MatchString(f.edit.PrimaryPhone.Value) {
			f.edit.PrimaryPhone.Problem = ""
		}
	}
	if f.edit.SecondaryPhone.Value != "" && !common.PIFOPhoneNumberRE.MatchString(f.edit.SecondaryPhone.Value) {
		f.edit.SecondaryPhone.Problem = `The "Secondary Contact Phone" field does not contain a valid phone number.`
	} else {
		f.edit.SecondaryPhone.Problem = ""
	}
	if f.edit.RepeaterInput.Value != "" && !common.PIFORealNumberRE.MatchString(f.edit.RepeaterInput.Value) {
		f.edit.RepeaterInput.Problem = `The "Repeater Input" field does not contain a valid number.`
	} else {
		f.edit.RepeaterInput.Problem = ""
	}
	if f.edit.RepeaterInputTone.Value != "" && !common.PIFORealNumberRE.MatchString(f.edit.RepeaterInputTone.Value) {
		f.edit.RepeaterInputTone.Problem = `The "Repeater Input Tone" field does not contain a valid number.`
	} else {
		f.edit.RepeaterInputTone.Problem = ""
	}
	if f.edit.RepeaterOutput.Value != "" && !common.PIFORealNumberRE.MatchString(f.edit.RepeaterOutput.Value) {
		f.edit.RepeaterOutput.Problem = `The "Repeater Output" field does not contain a valid number.`
	} else {
		f.edit.RepeaterOutput.Problem = ""
	}
	if f.edit.RepeaterOutputTone.Value != "" && !common.PIFORealNumberRE.MatchString(f.edit.RepeaterOutputTone.Value) {
		f.edit.RepeaterOutputTone.Problem = `The "Repeater Output Tone" field does not contain a valid number.`
	} else {
		f.edit.RepeaterOutputTone.Problem = ""
	}
	if f.edit.RepeaterOffset.Value != "" && f.edit.RepeaterOffset.Value != "+" && f.edit.RepeaterOffset.Value != "-" &&
		!common.PIFORealNumberRE.MatchString(f.edit.RepeaterOffset.Value) {
		f.edit.RepeaterOffset.Problem = `The "Repeater Offset" field does not contain a valid number or symbol.`
	} else {
		f.edit.RepeaterOffset.Problem = ""
	}
}

func validateRequiredIfComplete(reportType string, ef *message.EditField) bool {
	if reportType != "Complete" || ef.Value != "" {
		ef.Problem = ""
		return true
	}
	ef.Problem = fmt.Sprintf(`The %q field is required when "Report Type" is "Complete".`, ef.Label)
	return false

}
