package xscform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/xscmsg"
)

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

// NewField creates a new base field.
func NewField(id *xscmsg.FieldID, required bool) *Field {
	return &Field{id: id, required: required}
}

// newField creates a new base field, and returns it by value.
func newField(id *xscmsg.FieldID, required bool) Field {
	return Field{id: id, required: required}
}

// Field is the base for all form fields.
type Field struct {
	id       *xscmsg.FieldID
	required bool
	value    string
}

// ID returns the ID of the field.
func (f *Field) ID() *xscmsg.FieldID { return f.id }

// Get returns the value of the field.
func (f *Field) Get() string { return f.value }

// Set sets the value of the field.
func (f *Field) Set(v string) { f.value = v }

// Default returns an empty string for base fields.
func (*Field) Default() string { return "" }

// Validate ensures that required fields have a value.
func (f *Field) Validate(_ xscmsg.Message, _ bool) string {
	if f.value == "" {
		return fmt.Sprintf("field %q needs a value", f.id.Tag)
	}
	return ""
}

// BooleanField is a field with a Boolean value.
type BooleanField struct{ Field }

// Validate ensures the value is a proper Boolean.
func (f *BooleanField) Validate(msg xscmsg.Message, strict bool) string {
	if f.value == "" || f.value == "checked" {
		return ""
	}
	if !strict {
		switch strings.ToLower(f.value) {
		case "f", "false", "n", "no":
			f.value = ""
			return ""
		case "t", "true", "y", "yes":
			f.value = "checked"
			return ""
		}
	}
	return fmt.Sprintf("%q is not a valid boolean value for field %q", f.value, f.id.Tag)
}

// CallSignField is a field with an FCC call sign value.
type CallSignField struct{ Field }

// Validate ensures the value is a valid FCC call sign.  In non-strict mode, it
// converts the call sign to upper case.  (In strict mode, it does not required
// that the value is in upper case, because PackItForms doesn't.)
func (f *CallSignField) Validate(msg xscmsg.Message, strict bool) string {
	if err := f.Field.Validate(msg, strict); err != "" {
		return err
	}
	if f.value == "" {
		return ""
	}
	if !callSignRE.MatchString(f.value) {
		return fmt.Sprintf("%q is not a valid call sign for field %q", f.value, f.id.Tag)
	}
	if !strict {
		f.value = strings.ToUpper(f.value)
	}
	return ""
}

// CardinalNumberField is a field with a cardinal number value.
type CardinalNumberField struct{ Field }

// Validate ensures that the value is a properly-formatted cardinal number.
func (f *CardinalNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if f.value != "" && !cardinalNumberRE.MatchString(f.value) {
		return fmt.Sprintf("%q is not a valid integer value for field %q", f.value, f.id.Tag)
	}
	return ""
}

// ChoicesField is a field with a discrete set of allowed values.
type ChoicesField struct {
	Field
	Choices []string
}

// Validate ensures the field has one of the allowed values.  In non-strict
// mode, the values are case-insensitive, and any unambiguous prefix of an
// allowed value is accepted.
func (f *ChoicesField) Validate(msg xscmsg.Message, strict bool) string {
	var prefixOf string

	if err := f.Field.Validate(msg, strict); err != "" {
		return err
	}
	if f.value == "" {
		return f.Field.Validate(msg, strict)
	}
	for _, allowed := range f.Choices {
		if strict && f.value == allowed {
			return ""
		}
		if !strict && strings.EqualFold(allowed, f.value) {
			return ""
		}
		if !strict && len(f.value) < len(allowed) && strings.EqualFold(allowed[:len(f.value)], f.value) {
			if prefixOf == "" {
				prefixOf = allowed
			} else {
				prefixOf = "∅"
			}
		}
	}
	if prefixOf != "" && prefixOf != "∅" {
		f.value = prefixOf
		return ""
	}
	return fmt.Sprintf("%q is not a valid value for field %q", f.value, f.id.Tag)
}

// DateField is a field with a date value.
type DateField struct{ Field }

// Validate ensures the value is a valid date in MM/DD/YYYY format.
func (f *DateField) Validate(msg xscmsg.Message, strict bool) string {
	if err := f.Field.Validate(msg, strict); err != "" {
		return err
	}
	if !strict {
		if t, err := time.ParseInLocation("1/2/2006", f.value, time.Local); err == nil {
			// Add leading zeroes.  Also corrects 6/31 to 7/1, etc.
			f.value = t.Format("01/02/2006")
		}
	}
	if f.value != "" && !dateRE.MatchString(f.value) {
		return fmt.Sprintf("%q is not a valid date value for field %q", f.value, f.id.Tag)
	}
	return ""
}

// DateFieldDefaultNow is a field with a date value, which defaults to the
// current date.
type DateFieldDefaultNow struct{ DateField }

// Default returns the current date.
func (*DateFieldDefaultNow) Default() string {
	return time.Now().Format("01/02/2006")
}

// MessageNumberField is a field with a message number value.
type MessageNumberField struct{ Field }

// Validate doesn't actually validate the message number, because PackItForms
// doesn't validate it, and we don't want to be raising errors that PackItForms
// doesn't.  However, in non-strict mode, if the message number is well formed,
// we can at least canonicalize it.
func (f *MessageNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if err := f.Field.Validate(msg, strict); err != "" {
		return err
	}
	if !strict {
		if match := messageNumberRE.FindStringSubmatch(f.value); match != nil {
			num, _ := strconv.Atoi(match[2])
			f.value = fmt.Sprintf("%s%03d%s", strings.ToUpper(match[1]), num, strings.ToUpper(match[3]))
		}
	}
	return ""
}

// PhoneNumberField is a field with a phone number value.
type PhoneNumberField struct{ Field }

// Validate ensures that the value is a properly-formatted phone number.
func (f *PhoneNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if f.value != "" && !phoneNumberRE.MatchString(f.value) {
		return fmt.Sprintf("%q is not a valid phone number value for field %q", f.value, f.id.Tag)
	}
	return ""
}

// RealNumberField is a field with a real number value.
type RealNumberField struct{ Field }

// Validate ensures that the value is a properly-formatted real number.
func (f *RealNumberField) Validate(msg xscmsg.Message, strict bool) string {
	if f.value != "" && !realNumberRE.MatchString(f.value) {
		return fmt.Sprintf("%q is not a valid number value for field %q", f.value, f.id.Tag)
	}
	return ""
}

// TimeField is a field with a time value.
type TimeField struct{ Field }

// Validate ensures the time is a valid timestamp in 24-hour HH:MM format.
func (f *TimeField) Validate(msg xscmsg.Message, strict bool) string {
	if err := f.Field.Validate(msg, strict); err != "" {
		return err
	}
	if f.value != "" && !timeRE.MatchString(f.value) {
		return fmt.Sprintf("%q is not a valid time value for field %q", f.value, f.id.Tag)
	}
	return ""
}

// TimeFieldDefaultNow is a time field that defaults to the current time.
type TimeFieldDefaultNow struct{ TimeField }

// Default returns the current time.
func (*TimeFieldDefaultNow) Default() string {
	return time.Now().Format("15:04")
}
