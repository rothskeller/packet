package racesmar

import (
	"fmt"

	"github.com/rothskeller/packet/message/common"
)

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (f *RACESMAR) Validate() (problems []string) {
	problems = append(problems, f.StdFields.Validate()...)
	if f.AgencyName == "" {
		problems = append(problems, `The "Agency Name" field is required.`)
	}
	if f.EventName == "" {
		problems = append(problems, `The "Event Name" field is required.`)
	}
	if f.Assignment == "" {
		problems = append(problems, `The "Assignment" field is required.`)
	}
	switch f.FormVersion {
	case "1.6":
		if f.Resources[0].Qty == "" {
			problems = append(problems, `The "Resources Requested: Qty" field is required.`)
		}
		if f.Resources[0].RolePos == "" {
			problems = append(problems, `The "Resources Requested: Role/Position" field is required.`)
		}
		if f.Resources[0].PreferredType == "" {
			problems = append(problems, `The "Resources Requested: Preferred Type" field is required.`)
		}
		if f.Resources[0].MinimumType == "" {
			problems = append(problems, `The "Resources Requested: Minimum Type" field is required.`)
		}
	case "2.1":
		for i, r := range f.Resources {
			if i == 0 && r.Qty == "" {
				problems = append(problems, `The "Qty" field in row 1 of the "Resources Requested" section is required.`)
			} else if r.Qty != "" && !common.PIFOCardinalNumberRE.MatchString(r.Qty) {
				problems = append(problems, fmt.Sprintf(`The "Qty" field in row %d of the "Resources Requested" section does not contain a valid number.`, i+1))
			}
		}
		if f.Resources[0].RolePos == "" {
			problems = append(problems, `The "Role/Position" field in row 1 of the "Resources Requested" section is required.`)
		}
		if f.Resources[0].PreferredType == "" {
			problems = append(problems, `The "Preferred Type" field in row 1 of the "Resources Requested" section is required.`)
		}
		if f.Resources[0].MinimumType == "" {
			problems = append(problems, `The "Minimum Type" field in row 1 of the "Resources Requested" section is required.`)
		}
	case "2.3":
		for i, r := range f.Resources {
			if i == 0 && r.Qty == "" {
				problems = append(problems, `The "Qty" field in row 1 of the "Resources Requested" section is required.`)
			} else if r.Qty == "" && r.Role != "" {
				problems = append(problems, fmt.Sprintf(`The "Qty" field in row %d of the "Resources Requested" section is required when the "Role" field in that row is set.`, i+1))
			} else if r.Qty != "" && !common.PIFOCardinalNumberRE.MatchString(r.Qty) {
				problems = append(problems, fmt.Sprintf(`The "Qty" field in row %d of the "Resources Requested" section does not contain a valid number.`, i+1))
			}
			if (r.Position == "" && r.RolePos != r.Role) || (r.Position != "" && r.RolePos != r.Role+" / "+r.Position) {
				problems = append(problems, `The hidden "Role/Position" field in row %d of the "Resources Requested" section is not computed correctly.`)
			}
			switch r.Role {
			case "":
				if i == 0 {
					problems = append(problems, `The "Role" field in row 1 of the "Resources Requested" section is required.`)
				} else {
					if r.PreferredType != "" {
						problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section is not allowed unless the "Role" in that row is set.`, i+1))
					}
					if r.MinimumType != "" {
						problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section is not allowed unless the "Role" in that row is set.`, i+1))
					}
				}
			case "Field Communicator":
				switch r.PreferredType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "F1", "F2", "F3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
				switch r.MinimumType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "F1", "F2", "F3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
			case "Net Control Operator":
				switch r.PreferredType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "N1", "N2", "N3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
				switch r.MinimumType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "N1", "N2", "N3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
			case "Packet Operator":
				switch r.PreferredType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "P1", "P2", "P3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
				switch r.MinimumType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "P1", "P2", "P3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
			case "Shadow Communicator":
				switch r.PreferredType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "S1", "S2", "S3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Preferred Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
				switch r.MinimumType {
				case "":
					if i != 0 {
						problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section is required when the "Role" in that row is set.`, i+1))
					}
				case "S1", "S2", "S3", "Type IV", "Type V":
					break
				default:
					problems = append(problems, fmt.Sprintf(`The "Minimum Type" field in row %d of the "Resources Requested" section does not contain a valid type for the "Role" in that row.`, i+1))
				}
			default:
				problems = append(problems, fmt.Sprintf(`The "Role" field in row %d of the "Resources Requested" section does not contain a valid role.`, i+1))
			}
		}
	}
	if f.RequestedArrivalDates == "" {
		problems = append(problems, `The "Requested Arrival Date(s)" field is required.`)
	}
	if f.RequestedArrivalTimes == "" {
		problems = append(problems, `The "Requested Arrival Time(s)" field is required.`)
	}
	if f.NeededUntilDates == "" {
		problems = append(problems, `The "Needed Until Date(s)" field is required.`)
	}
	if f.NeededUntilTimes == "" {
		problems = append(problems, `The "Needed Until Time(s)" field is required.`)
	}
	if f.ReportingLocation == "" {
		problems = append(problems, `The "Reporting Location" field is required.`)
	}
	if f.ContactOnArrival == "" {
		problems = append(problems, `The "Contact on Arrival" field is required.`)
	}
	if f.TravelInfo == "" {
		problems = append(problems, `The "Travel Info" field is required.`)
	}
	if f.RequestedByName == "" {
		problems = append(problems, `The "Requested By Name" field is required.`)
	}
	if f.RequestedByTitle == "" {
		problems = append(problems, `The "Requested By Title" field is required.`)
	}
	if f.RequestedByContact == "" {
		problems = append(problems, `The "Requested By Contact" field is required.`)
	}
	if f.ApprovedByName == "" {
		problems = append(problems, `The "Approved By Name" field is required.`)
	}
	if f.ApprovedByTitle == "" {
		problems = append(problems, `The "Approved By Title" field is required.`)
	}
	if f.ApprovedByContact == "" {
		problems = append(problems, `The "Approved By Contact" field is required.`)
	}
	if f.ApprovedByDate == "" {
		problems = append(problems, `The "Approved By Date" field is required.`)
	} else if !common.PIFODateRE.MatchString(f.ApprovedByDate) {
		problems = append(problems, `The "Approved By Date" field does not contain a valid date.`)
	}
	if f.ApprovedByTime == "" {
		problems = append(problems, `The "Approved By Time" field is required.`)
	} else if !common.PIFOTimeRE.MatchString(f.ApprovedByTime) {
		problems = append(problems, `The "Approved By Time" field does not contain a valid time.`)
	}
	return problems
}
