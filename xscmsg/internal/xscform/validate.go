package xscform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/xscmsg"
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
	messageNumberRE = regexp.MustCompile(`(?i)((?:[A-Z]{3}|[0-9][A-Z]{2}|[A-Z][0-9][A-Z])-)(\d+)([A-Z]?)$`)
)

// ValidateBoolean ensures the value is a proper Boolean.
func ValidateBoolean(f *xscmsg.Field, _ *xscmsg.Message, strict bool) string {
	if f.Value == "" || f.Value == "checked" {
		return ""
	}
	if !strict {
		switch strings.ToLower(f.Value) {
		case "f", "false", "n", "no":
			f.Value = ""
			return ""
		case "t", "true", "y", "yes":
			f.Value = "checked"
			return ""
		}
	}
	return fmt.Sprintf("%q is not a valid boolean value for field %q.", f.Value, f.Def.Tag)
}

// ValidateCallSign doesn't actually validate call signs, since PackItForms
// doesn't, and we don't want to raise any errors that it doesn't.  However, in
// non-strict mode, we can at least upcase the callsign.
func ValidateCallSign(f *xscmsg.Field, _ *xscmsg.Message, strict bool) string {
	if !strict {
		f.Value = strings.ToUpper(f.Value)
	}
	return ""
}

// ValidateCardinalNumber ensures that the value is a properly-formatted
// cardinal number.
func ValidateCardinalNumber(f *xscmsg.Field, _ *xscmsg.Message, _ bool) string {
	if f.Value != "" && !cardinalNumberRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid integer value for field %q.", f.Value, f.Def.Tag)
	}
	return ""
}

// ValidateChoices ensures the field has one of the allowed values.  In
// non-strict mode, the values are case-insensitive, and any unambiguous prefix
// of an allowed value is accepted.
func ValidateChoices(f *xscmsg.Field, _ *xscmsg.Message, strict bool) string {
	var prefixOf string

	if f.Value == "" {
		return ""
	}
	for _, allowed := range f.Def.Choices {
		if strict && f.Value == allowed {
			return ""
		}
		if !strict && strings.EqualFold(allowed, f.Value) {
			return ""
		}
		if !strict && len(f.Value) < len(allowed) && strings.EqualFold(allowed[:len(f.Value)], f.Value) {
			if prefixOf == "" {
				prefixOf = allowed
			} else {
				prefixOf = "∅"
			}
		}
	}
	if prefixOf != "" && prefixOf != "∅" {
		f.Value = prefixOf
		return ""
	}
	return fmt.Sprintf("%q is not a valid value for field %q.", f.Value, f.Def.Tag)
}

// ValidateDate ensures the value is a valid date in MM/DD/YYYY format.
func ValidateDate(f *xscmsg.Field, _ *xscmsg.Message, strict bool) string {
	if !strict {
		if t, err := time.ParseInLocation("1/2/2006", f.Value, time.Local); err == nil {
			// Add leading zeroes.  Also corrects 6/31 to 7/1, etc.
			f.Value = t.Format("01/02/2006")
		}
	}
	if f.Value != "" && !dateRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid date value for field %q.", f.Value, f.Def.Tag)
	}
	return ""
}

// DefaultDate returns the current date.
func DefaultDate() string {
	return time.Now().Format("01/02/2006")
}

// ValidateMessageNumber doesn't actually validate the message number, because
// PackItForms doesn't validate it, and we don't want to be raising errors that
// PackItForms doesn't.  However, in non-strict mode, if the message number is
// well formed, we can at least canonicalize it.
func ValidateMessageNumber(f *xscmsg.Field, _ *xscmsg.Message, strict bool) string {
	if !strict {
		if match := messageNumberRE.FindStringSubmatch(f.Value); match != nil {
			num, _ := strconv.Atoi(match[2])
			f.Value = fmt.Sprintf("%s%03d%s", strings.ToUpper(match[1]), num, strings.ToUpper(match[3]))
		}
	}
	return ""
}

// ValidatePhoneNumber ensures that the value is a properly-formatted phone
// number.
func ValidatePhoneNumber(f *xscmsg.Field, _ *xscmsg.Message, _ bool) string {
	if f.Value != "" && !phoneNumberRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid phone number value for field %q.", f.Value, f.Def.Tag)
	}
	return ""
}

// ValidateRealNumber ensures that the value is a properly-formatted real
// number.
func ValidateRealNumber(f *xscmsg.Field, _ *xscmsg.Message, _ bool) string {
	if f.Value != "" && !realNumberRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid number value for field %q.", f.Value, f.Def.Tag)
	}
	return ""
}

// ValidateRequired verifies that a field has a value.
var ValidateRequired = xscmsg.ValidateRequired

// ValidateTime ensures the time is a valid timestamp in 24-hour HH:MM format.
func ValidateTime(f *xscmsg.Field, _ *xscmsg.Message, _ bool) string {
	if f.Value != "" && !timeRE.MatchString(f.Value) {
		return fmt.Sprintf("%q is not a valid time value for field %q.", f.Value, f.Def.Tag)
	}
	return ""
}

// DefaultTime returns the current time.
func DefaultTime() string {
	return time.Now().Format("15:04")
}

// ValidateUnknownField reports a problem with a field that isn't defined for
// the message type.
func ValidateUnknownField(f *xscmsg.Field, m *xscmsg.Message, _ bool) string {
	return fmt.Sprintf("This form has a field %q, which is not defined for %s forms.", f.Def.Tag, m.Type.Tag)
}

// SetChoiceFieldFromValue sets the value of the specified field (whose
// definition must have a Choices list) to the specified value if it is one of
// the choices, or to the last choice on the list if the specified value is
// something else.
func SetChoiceFieldFromValue(f *xscmsg.Field, value string) {
	for _, c := range f.Def.Choices {
		if c == value {
			f.Value = c
			return
		}
	}
	f.Value = f.Def.Choices[len(f.Def.Choices)-1]
}

/*
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
*/
