package basemsg

import (
	"fmt"
	"time"

	"github.com/rothskeller/packet/message/common"
)

// TODO defining this functions on BaseMessage means that only messages for
// which it's appropriate can leverage BaseMessage.

// Validate checks the contents of the message for compliance with rules
// enforced by standard Santa Clara County packet software (Outpost and
// PackItForms).  It returns a list of strings describing problems that those
// programs would flag or block.
func (bm *BaseMessage) Validate() (problems []string) {
	for _, f := range bm.Fields {
		if p := validatePresence(f); p != "" {
			problems = append(problems, p)
			continue
		}
		if problem := f.PIFOValid(f); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}
func validatePresence(f *Field) string {
	if f.Presence == nil {
		return ""
	}
	presence, when := f.Presence()
	if when != "" {
		when = " when " + when
	}
	var value string
	if f.Value != nil {
		value = *f.Value
	} else if f.EditValue != nil {
		value = f.EditValue(f)
	}
	switch presence {
	case PresenceNotAllowed:
		if value != "" {
			return fmt.Sprintf("The %q field cannot have a value%s.", f.Label, when)
		}
	case PresenceRequired:
		if value == "" {
			return fmt.Sprintf(`The %q field is required%s.`, f.Label, when)
		}
	}
	return ""
}

// ValidCardinal is a validation function for a field that can contain a
// cardinal number.
func ValidCardinal(f *Field) string {
	if *f.Value != "" && !common.PIFOCardinalNumberRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid number.", f.Label)
	}
	return ""
}

// ValidDate is a validation function for a field that can contain a date.
func ValidDate(f *Field) string {
	if *f.Value != "" && !common.PIFODateRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid date (MM/DD/YYYY).", f.Label)
	}
	return ""
}

// ValidDateTime is a validation function for a field that can contain a
// date/time value.
func ValidDateTime(f *Field, date, tval string) string {
	var dtval = date + " " + tval
	if t, err := time.ParseInLocation("01/02/2006 15:04", dtval, time.Local); err != nil || dtval != t.Format("01/02/2006 15:04") {
		return fmt.Sprintf("The %q field does not contain a valid date and time in MM/DD/YYYY HH:MM format.", f.Label)
	}
	return ""
}

// ValidFCCCallSign is a validation function for a field that can contain an FCC
// call sign.
func ValidFCCCallSign(f *Field) string {
	if *f.Value != "" && !fccCallSignRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid FCC call sign.", f.Label)
	}
	return ""
}

// ValidFrequency is a validation function for a field that can contain a
// frequency in MHz.
func ValidFrequency(f *Field) string {
	if *f.Value != "" && !common.PIFOFrequencyRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid frequency.", f.Label)
	}
	return ""
}

// ValidFrequencyOffset is a validation function for a field that can contain a
// frequency offset in MHz.
func ValidFrequencyOffset(f *Field) string {
	if *f.Value != "" && !common.PIFOFrequencyOffsetRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid frequency offset.", f.Label)
	}
	return ""
}

// ValidMessageNumber is a validation function for a field that can contain a
// message number.
func ValidMessageNumber(f *Field) string {
	if *f.Value != "" && !messageNumberRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid message number.", f.Label)
	}
	return ""

}

// ValidPhoneNumber is a validation function for a field that can contain a
// phone number.
func ValidPhoneNumber(f *Field) string {
	if *f.Value != "" && !common.PIFOPhoneNumberRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a phone number.", f.Label)
	}
	return ""
}

// ValidReal is a validation function for a field that can contain a real
// number.
func ValidReal(f *Field) string {
	if *f.Value != "" && !common.PIFORealNumberRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid number.", f.Label)
	}
	return ""
}

// ValidRestricted is a validation function for a field that can contain a value
// from a restricted set.
func ValidRestricted(f *Field) string {
	if *f.Value != "" && !f.Choices.IsPIFO(*f.Value) {
		return fmt.Sprintf("The %q field does not contain one of its allowed values.", f.Label)
	}
	return ""
}

// ValidTacticalCallSign is a validation function for a field that can contain a
// tactical call sign.
func ValidTacticalCallSign(f *Field) string {
	if *f.Value != "" && !tacticalCallSignRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid tactical call sign.", f.Label)
	}
	return ""
}

// ValidTime is a validation function for a field that can contain a time.
func ValidTime(f *Field) string {
	if *f.Value != "" && !common.PIFOTimeRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid time (HH:MM, 24-hour clock).", f.Label)
	}
	return ""
}
