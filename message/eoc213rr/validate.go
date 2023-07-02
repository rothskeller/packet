package eoc213rr

import "github.com/rothskeller/packet/message/common"

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (f *EOC213RR) Validate() (problems []string) {
	problems = append(problems, f.StdFields.Validate()...)
	if f.IncidentName == "" {
		problems = append(problems, `The "Incident Name" field is required.`)
	}
	if f.DateInitiated == "" {
		problems = append(problems, `The "Date Initiated" field is required.`)
	} else if !common.PIFODateRE.MatchString(f.DateInitiated) {
		problems = append(problems, `The "Date Initiated" field does not contain a valid date.`)
	}
	if f.TimeInitiated == "" {
		problems = append(problems, `The "Time Initiated" field is required.`)
	} else if !common.PIFOTimeRE.MatchString(f.TimeInitiated) {
		problems = append(problems, `The "Time Initiated" field does not contain a valid time.`)
	}
	if f.RequestedBy == "" {
		problems = append(problems, `The "Requested by" field is required.`)
	}
	if f.QtyUnit == "" {
		problems = append(problems, `The "Qty/Unit" field is required.`)
	}
	if f.ResourceDescription == "" {
		problems = append(problems, `The "Resource Description" field is required.`)
	}
	if f.ResourceArrival == "" {
		problems = append(problems, `The "Arrival" field is required.`)
	}
	switch f.Priority {
	case "":
		problems = append(problems, `The "Priority" field is required.`)
	case "Now", "High", "Medium", "Low":
		break
	default:
		problems = append(problems, `The "Priority" field does not contain a valid priority.`)
	}
	if f.DeliverTo == "" {
		problems = append(problems, `The "Deliver to" field is required.`)
	}
	if f.DeliverToLocation == "" {
		problems = append(problems, `The "Location" field is required.`)
	}
	if f.EquipmentOperator != "" && f.EquipmentOperator != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Equipment Operator" field does not contain a valid checkbox value.`)
	}
	if f.Lodging != "" && f.Lodging != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Lodging" field does not contain a valid checkbox value.`)
	}
	if f.Fuel != "" && f.Fuel != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Fuel" field does not contain a valid checkbox value.`)
	} else {
		if f.Fuel == "" && f.FuelType != "" {
			problems = append(problems, `The "Supplemental Requirements: Fuel Type" field must not have a value unless the "Fuel" box is checked.`)
		}
		if f.Fuel != "" && f.FuelType == "" {
			problems = append(problems, `The "Supplemental Requirements: Fuel Type" field is required when the "Fuel" box is checked.`)
		}
	}
	if f.Power != "" && f.Power != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Power" field does not contain a valid checkbox value.`)
	}
	if f.Meals != "" && f.Meals != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Meals" field does not contain a valid checkbox value.`)
	}
	if f.Maintenance != "" && f.Maintenance != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Maintenance" field does not contain a valid checkbox value.`)
	}
	if f.Water != "" && f.Water != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Water" field does not contain a valid checkbox value.`)
	}
	if f.Other != "" && f.Other != "checked" {
		problems = append(problems, `The "Supplemental Requirements: Other" field does not contain a valid checkbox value.`)
	}
	return problems
}
