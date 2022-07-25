package xscform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Most of the functions in this file correspond to the validation checks called
// out in the "required" and "class" attributes of the various input controls in
// PackItForms HTML files.

var (
	// these REs are taken from the PackItForms source code
	cardinalNumberRE  = regexp.MustCompile(`^[0-9]*$`)
	dateRE            = regexp.MustCompile(`^(0[1-9]|1[012])/(0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]$`)
	frequencyRE       = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	frequencyOffsetRE = regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+]$`)
	phoneNumberRE     = regexp.MustCompile(`^[a-zA-Z ]*([+][0-9]+ )?[0-9][0-9 -]*([xX][0-9]+)?$`)
	realNumberRE      = regexp.MustCompile(`^[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+$`)
	timeRE            = regexp.MustCompile(`^([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00$`)
	// these are defined locally
	callSignRE      = regexp.MustCompile(`(?i)^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)
	messageNumberRE = regexp.MustCompile(`(?i)((?:[A-Z]{3}|[0-9][A-Z]{2}|[A-Z][0-9][A-Z])-)(\d+)([A-Z]?)$`)
)

// ValidateBoolean verifies that the value is a Boolean (i.e., checkbox in the
// SCCoPIFO HTML).
func ValidateBoolean(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	switch strings.ToLower(value) {
	case "", "f", "false", "n", "no":
		return "", ""
	case "checked", "t", "true", "y", "yes":
		return "checked", ""
	default:
		return value, fmt.Sprintf("%q is not a valid boolean value for field %q", value, fd.Tag)
	}
}

// ValidateCallSign verifies that the value is an FCC call sign.
func ValidateCallSign(_ *XSCForm, fd *FieldDefinition, value string, strict bool) (newval, problem string) {
	if value == "" {
		return value, ""
	}
	if !callSignRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid call sign for field %q", value, fd.Tag)
	}
	if strict {
		return value, ""
	}
	return strings.ToUpper(value), ""
}

// ValidateCardinalNumber verifies that the value is a cardinal number (i.e.,
// non-negative integer).
func ValidateCardinalNumber(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !cardinalNumberRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid integer value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateComputedChoice handles a common pattern where the value of the target
// field is computed based on the value of another field.  Specifically, if the
// value of the other field is one of the allowed values for this field, it is
// kept, and otherwise, the last allowed value of this field is used.  This
// "validation" is used for fields that are implemented in the PackItForms HTML
// with a combo box and a "compatible_values" function call.
func ValidateComputedChoice(f *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, _ string) {
	other := f.Get(fd.ComputedFromField)
	other, _ = ValidateSelect(f, fd, other, false)
	for _, allowed := range fd.Values {
		if other == allowed {
			return other, ""
		}
	}
	return fd.Values[len(fd.Values)-1], ""
}

// ValidateDate verifies that the value is a date.
func ValidateDate(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if t, err := time.ParseInLocation("1/2/2006", value, time.Local); err == nil {
		value = t.Format("01/02/2006")
	}
	if value != "" && !dateRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid date value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateFrequency verifies that the value is a frequency.
func ValidateFrequency(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !frequencyRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid frequency value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateFrequencyOffset verifies that the value is a frequency offset.
func ValidateFrequencyOffset(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !frequencyOffsetRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid frequency offset value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateMessageNumber doesn't actually validate the message number, because
// PackItForms doesn't validate it, and we don't want to be raising errors that
// PackItForms doesn't.  Instead, we'll raise an error about the message number
// in the subject line.  However, if the message number is well formed, we can
// at least canonicalize it.
func ValidateMessageNumber(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if match := messageNumberRE.FindStringSubmatch(value); match != nil {
		num, _ := strconv.Atoi(match[2])
		return fmt.Sprintf("%s%03d%s", strings.ToUpper(match[1]), num, strings.ToUpper(match[3])), ""
	}
	return value, ""
}

// ValidatePhoneNumber verifies that the value is a phone number.
func ValidatePhoneNumber(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !phoneNumberRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid phone number value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateRealNumber verifies that the value is a real number.
func ValidateRealNumber(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !realNumberRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid number value for field %q", value, fd.Tag)
	}
	return value, ""
}

// ValidateRequired verifies that the value not empty.
func ValidateRequired(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value == "" {
		return value, fmt.Sprintf("field %q needs a value", fd.Tag)
	}
	return value, ""
}

// ValidateRequiredForComplete verifies that, if the form type is Complete, the
// value is not empty.
func ValidateRequiredForComplete(xf *XSCForm, fd *FieldDefinition, value string, strict bool) (newval, problem string) {
	if xf.form.Get("19.") == "Complete" {
		return ValidateRequired(xf, fd, value, strict)
	}
	return value, ""
}

// ValidateSelect verifies that the value, if not empty, is one of the allowed
// choices.
func ValidateSelect(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	var prefixOf string

	if value == "" {
		return value, ""
	}
	for _, allowed := range fd.Values {
		if strings.EqualFold(allowed, value) {
			return allowed, ""
		}
		if len(value) < len(allowed) && strings.EqualFold(allowed[:len(value)], value) {
			if prefixOf == "" {
				prefixOf = allowed
			} else {
				prefixOf = "∅"
			}
		}
	}
	if prefixOf != "" && prefixOf != "∅" {
		return prefixOf, ""
	}
	return value, fmt.Sprintf("%q is not a valid value for field %q", value, fd.Tag)
}

// ValidateTime verifies that the value is a time.
func ValidateTime(_ *XSCForm, fd *FieldDefinition, value string, _ bool) (newval, problem string) {
	if value != "" && !timeRE.MatchString(value) {
		return value, fmt.Sprintf("%q is not a valid time value for field %q", value, fd.Tag)
	}
	return value, ""
}
