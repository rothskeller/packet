package sheltstat

import "github.com/rothskeller/packet/message/common"

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (f *SheltStat) Validate() (problems []string) {
	problems = append(problems, f.StdFields.Validate()...)
	switch f.ReportType {
	case "":
		problems = append(problems, `The "Report Type" field is required.`)
	case "Update", "Complete":
		break
	default:
		problems = append(problems, `The "Report Type" field does not contain a valid report type value.`)
	}
	if f.ShelterName == "" {
		problems = append(problems, `The "Shelter Name" field is required.`)
	}
	switch f.ShelterType {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Shelter Type" field is required when the "Report Type" is "Complete".`)
		}
	case "Type 1", "Type 2", "Type 3", "Type 4":
		break
	default:
		problems = append(problems, `The "Shelter Type" field does not contain a valid shelter type value.`)
	}
	switch f.ShelterStatus {
	case "":
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Shelter Status" field is required when the "Report Type" is "Complete".`)
		}
	case "Open", "Closed", "Full":
		break
	default:
		problems = append(problems, `The "Shelter Status" field does not contain a valid shelter status value.`)
	}
	if f.ShelterAddress == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Shelter Address" field is required when "Report Type" is "Complete".`)
	}
	if f.FormVersion >= "2.2" {
		switch f.ShelterCityCode {
		case "":
			if f.ShelterCity != "" {
				problems = append(problems, `The hidden, required "Shelter City Code" field is not set.`)
			} // otherwise we let the problem with ShelterCity be reported.
		case "Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated":
			break
		default:
			problems = append(problems, `The hidden "Shelter City Code" field does not contain a valid shelter city code.`)
		}
		if f.ShelterCity == "" {
			problems = append(problems, `The "Shelter City" field is required.`)
		}
	} else {
		switch f.ShelterCity {
		case "":
			problems = append(problems, `The "Shelter City" field is required.`)
		case "Campbell", "Cupertino", "Gilroy", "Los Altos", "Los Altos Hills", "Los Gatos", "Milpitas", "Monte Sereno", "Morgan Hill", "Mountain View", "Palo Alto", "San Jose", "Santa Clara", "Saratoga", "Sunnyvale", "Unincorporated":
			break
		default:
			problems = append(problems, `The "Shelter City" field does not contain a valid shelter city name.`)
		}
	}
	if f.ShelterState == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Shelter State" field is required when "Report Type" is "Complete".`)
	}
	if f.Latitude != "" && !common.PIFORealNumberRE.MatchString(f.Latitude) {
		problems = append(problems, `The "Shelter Latitude" field does not contain a valid number.`)
	}
	if f.Longitude != "" && !common.PIFORealNumberRE.MatchString(f.Longitude) {
		problems = append(problems, `The "Shelter Longitude" field does not contain a valid number.`)
	}
	if f.Capacity == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Capacity" field is required when "Report Type" is "Complete".`)
	} else if f.Capacity != "" && !common.PIFOCardinalNumberRE.MatchString(f.Capacity) {
		problems = append(problems, `The "Capacity" field does not contain a valid number.`)
	}
	if f.Occupancy == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Occupancy" field is required when "Report Type" is "Complete".`)
	} else if f.Occupancy != "" && !common.PIFOCardinalNumberRE.MatchString(f.Occupancy) {
		problems = append(problems, `The "Occupancy" field does not contain a valid number.`)
	}
	switch f.PetFriendly {
	case "", "checked", "false":
		break
	default:
		problems = append(problems, `The "Pet Friendly" field does not contain a valid value.`)
	}
	switch f.BasicSafetyInspection {
	case "", "checked", "false":
		break
	default:
		problems = append(problems, `The "Basic Safety Inspection" field does not contain a valid value.`)
	}
	switch f.ATC20Inspection {
	case "", "checked", "false":
		break
	default:
		problems = append(problems, `The "ATC-20 Inspection" field does not contain a valid value.`)
	}
	if f.FormVersion >= "2.2" {
		switch f.ManagedByCode {
		case "":
			if f.ManagedBy != "" {
				problems = append(problems, `The hidden, required "Managed By Code" field is not set.`)
			} // otherwise we let the problem with ManagedBy be reported.
		case "American Red Cross", "Private", "Community", "Government", "Other":
			break
		default:
			problems = append(problems, `The hidden "Managed By Code" field does not contain a valid management code.`)
		}
		if f.ManagedBy == "" {
			problems = append(problems, `The "Managed By" field is required.`)
		}
	} else {
		switch f.ManagedBy {
		case "":
			problems = append(problems, `The "Managed By" field is required.`)
		case "American Red Cross", "Private", "Community", "Government", "Other":
			break
		default:
			problems = append(problems, `The "Managed By" field does not contain a valid management code.`)
		}
	}
	if f.PrimaryContact == "" && f.ReportType == "Complete" {
		problems = append(problems, `The "Primary Contact" field is required when "Report Type" is "Complete".`)
	}
	if f.PrimaryPhone == "" {
		if f.ReportType == "Complete" {
			problems = append(problems, `The "Primary Contact Phone" field is required when the "Report Type" is "Complete".`)
		}
	} else if !common.PIFOPhoneNumberRE.MatchString(f.PrimaryPhone) {
		problems = append(problems, `The "Primary Contact Phone" field does not contain a valid phone number.`)
	}
	if f.SecondaryPhone != "" && !common.PIFOPhoneNumberRE.MatchString(f.SecondaryPhone) {
		problems = append(problems, `The "Secondary Contact Phone" field does not contain a valid phone number.`)
	}
	if f.RepeaterInput != "" && !common.PIFOFrequencyRE.MatchString(f.RepeaterInput) {
		problems = append(problems, `The "Repeater Input" field does not contain a valid frequency value.`)
	}
	if f.RepeaterOutput != "" && !common.PIFOFrequencyRE.MatchString(f.RepeaterOutput) {
		problems = append(problems, `The "Repeater Output" field does not contain a valid frequency value.`)
	}
	if f.RepeaterOffset != "" && !common.PIFOFrequencyOffsetRE.MatchString(f.RepeaterOffset) {
		problems = append(problems, `The "Repeater Offset" field does not contain a valid frequency offset value.`)
	}
	switch f.RemoveFromList {
	case "", "checked":
		break
	case "false":
		if f.FormVersion >= "2.3" {
			break
		}
		fallthrough
	default:
		problems = append(problems, `The "Remove from List" field does not contain a valid checkbox value.`)
	}
	return problems
}
