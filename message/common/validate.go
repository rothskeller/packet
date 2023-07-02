package common

import "regexp"

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (s *StdFields) Validate() (problems []string) {
	if s.OriginMsgID == "" {
		problems = append(problems, `The "Origin Message Number" field is required.`)
	}
	if s.MessageDate == "" {
		problems = append(problems, `The "Date" field is required.`)
	} else if !PIFODateRE.MatchString(s.MessageDate) {
		problems = append(problems, `The "Date" field does not contain a valid date.`)
	}
	if s.MessageTime == "" {
		problems = append(problems, `The "Time" field is required.`)
	} else if !PIFOTimeRE.MatchString(s.MessageTime) {
		problems = append(problems, `The "Time" field does not contain a valid time.`)
	}
	switch s.Handling {
	case "":
		problems = append(problems, `The "Handling" field is required.`)
	case "ROUTINE", "PRIORITY", "IMMEDIATE":
		break
	default:
		problems = append(problems, `The "Handling" field does not contain a valid message handling order.`)
	}
	if s.ToICSPosition == "" {
		problems = append(problems, `The "To ICS Position" field is required.`)
	}
	if s.FromICSPosition == "" {
		problems = append(problems, `The "From ICS Position" field is required.`)
	}
	if s.ToLocation == "" {
		problems = append(problems, `The "To Location" field is required.`)
	}
	if s.FromLocation == "" {
		problems = append(problems, `The "From Location" field is required.`)
	}
	if s.OpName == "" {
		problems = append(problems, `The "Radio Operator Only: Name" field is required.`)
	}
	if s.OpCall == "" {
		problems = append(problems, `The "Radio Operator Only: Call Sign" field is required.`)
	}
	if s.OpDate == "" {
		problems = append(problems, `The "Radio Operator Only: Date" field is required.`)
	} else if !PIFODateRE.MatchString(s.OpDate) {
		problems = append(problems, `The "Radio Operator Only: Date" field does not contain a valid date.`)
	}
	if s.OpTime == "" {
		problems = append(problems, `The "Radio Operator Only: Time" field is required.`)
	} else if !PIFOTimeRE.MatchString(s.OpTime) {
		problems = append(problems, `The "Radio Operator Only: Time" field does not contain a valid time.`)
	}
	return problems
}

// Regular expressions for data type validation, taken from the PackItForms
// code.  (Unmodified except for JavaScript-to-Go conversion.)
var (
	PIFODateRE            = regexp.MustCompile(`^(?:0[1-9]|1[012])/(?:0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]$`)
	PIFOTimeRE            = regexp.MustCompile(`^(?:([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00)$`)
	PIFOPhoneNumberRE     = regexp.MustCompile(`^[a-zA-Z ]*(?:[+][0-9]+ )?[0-9][0-9 -]*(?:[xX][0-9]+)?$`)
	PIFOCardinalNumberRE  = regexp.MustCompile(`^[0-9]*$`)
	PIFORealNumberRE      = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+)$`)
	PIFOFrequencyRE       = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)
	PIFOFrequencyOffsetRE = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+])$`)
)
