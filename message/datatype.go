package message

// This file defines functions that add behavior to Field definitions based on
// the type of data expected to be stored in that field.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/envelope"
)

// Regular expressions for data type validation.  The ones with PIFO* names are
// taken from the PackItForms code, unmodified except for JavaScript-to-Go
// conversion.
var (
	PIFOCardinalNumberRE  = regexp.MustCompile(`^[0-9]+$`) // changed * to +
	dateLooseRE           = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)
	PIFODateRE            = regexp.MustCompile(`^(?:0[1-9]|1[012])/(?:0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]$`)
	fccCallSignRE         = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)
	PIFOFrequencyRE       = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)
	PIFOFrequencyOffsetRE = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+])$`)
	messageNumberLooseRE  = regexp.MustCompile(`^([A-Z0-9]{3})-(\d+)([A-Z]?)$`)
	messageNumberRE       = regexp.MustCompile(`^(?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-(?:[1-9][0-9]{3,}|[0-9]{3})[A-Z]$`)
	PIFOPhoneNumberRE     = regexp.MustCompile(`^[a-zA-Z ]*(?:[+][0-9]+ )?[0-9][0-9 -]*(?:[xX][0-9]+)?$`)
	PIFORealNumberRE      = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+)$`)
	tacticalCallSignRE    = regexp.MustCompile(`^[A-Z][A-Z0-9]{4,5}$`)
	timeLooseRE           = regexp.MustCompile(`^([1-9]:|[01][0-9]:?|2[0-4]:?)([0-5][0-9])$`)
	PIFOTimeRE            = regexp.MustCompile(`^(?:([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00)$`)
)

// NewAddressListField adds defaults to a Field that are appropriate for a field
// that contains a list of packet or email addresses.  It modifies its argument
// and returns it for chaining.
func NewAddressListField(f *Field) *Field {
	if f.EditApply == nil {
		f.EditApply = func(f *Field, s string) {
			if addrs, err := envelope.ParseAddressList(s); err == nil {
				// Normalize the syntax of each address by
				// calling String() on it.  This quotes or
				// unquotes things, adds or removes angle
				// brackets, etc.
				strs := make([]string, len(addrs))
				for i, a := range addrs {
					strs[i] = a.String()
				}
				// Join the list, which normalizes the
				// separation between addresses.
				s = strings.Join(strs, ", ")
			}
			*f.Value = s
		}
	}
	if f.EditValid == nil {
		f.EditValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if _, err := envelope.ParseAddressList(*f.Value); err != nil {
				return fmt.Sprintf("The %q field does not contain a valid address list.", f.Label)
			}
			return ""
		}
	}
	return AddFieldDefaults(f)
}

// NewAggregatorField adds defaults to a Field that are appropriate for a
// pseudo-field that displays and/or edits the contents of multiple other fields
// as a unit.  It modifies its argument and returns it for chaining.
func NewAggregatorField(f *Field) *Field { return AddFieldDefaults(f) }

// NewCalculatedField adds defaults to a Field that are appropriate for a field
// whose value is calculated from other fields, and is not displayed or edited
// directly.  It modifies its argument and returns it for chaining.
func NewCalculatedField(f *Field) *Field {
	if f.Compare == nil {
		f.Compare = CompareNone
	}
	if f.TableValue == nil {
		f.TableValue = TableOmit
	}
	return AddFieldDefaults(f)
}

// ApplyCardinalNumber applies an edited value to a cardinal number field.
func ApplyCardinalNumber(f *Field, v string) {
	v = strings.TrimSpace(v)
	if n, err := strconv.Atoi(v); err == nil {
		v = strconv.Itoa(n)
	}
	*f.Value = v
}

// ValidCardinalNumber verifies that the provided string is a valid cardinal
// number according to PackItForms.
func ValidCardinalNumber(f *Field) string {
	if *f.Value != "" && !PIFOCardinalNumberRE.MatchString(*f.Value) {
		return fmt.Sprintf("The %q field does not contain a valid number.", f.Label)
	}
	return ""
}

// NewCardinalNumberField adds defaults to a Field that are appropriate for a
// field that contains a non-negative integer.  It modifies its argument and
// returns it for chaining.
func NewCardinalNumberField(f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			return ValidCardinalNumber(f)
		}
	}
	if f.Compare == nil {
		f.Compare = CompareCardinal
	}
	if f.EditApply == nil {
		f.EditApply = ApplyCardinalNumber
	}
	return AddFieldDefaults(f)
}

// NewDateTimeField adds defaults to a Field that are appropriate for a
// pseudo-field that displays and edits a pair of DateWithTime and TimeWithDate
// fields together.  These fields must always come as a triplet: a DateWithTime
// field, a TimeWithDate field, and a DateTime field that refers to the first
// two.  NewDateTimeField modifies its argument and returns it for chaining.
func NewDateTimeField(f *Field, date, tval *string) *Field {
	if f.TableValue == nil {
		f.TableValue = func(*Field) string {
			return SmartJoin(*date, *tval, " ")
		}
	}
	if f.EditWidth == 0 {
		f.EditWidth = 16
	}
	if f.EditHint == "" {
		f.EditHint = "MM/DD/YYYY HH:MM"
	}
	if f.EditValue == nil {
		f.EditValue = func(*Field) string {
			return SmartJoin(*date, *tval, " ")
		}
	}
	if f.EditApply == nil {
		f.EditApply = func(_ *Field, v string) {
			words := strings.Fields(v)
			f := NewDateField(false, &Field{Value: date})
			if len(words) > 0 {
				f.EditApply(f, words[0])
			} else {
				f.EditApply(f, "")
			}
			f = NewTimeField(false, &Field{Value: tval})
			if len(words) > 1 {
				f.EditApply(f, strings.Join(words[1:], " "))
			} else {
				f.EditApply(f, "")
			}
		}
	}
	if f.EditValid == nil {
		f.EditValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			dtval := SmartJoin(*date, *tval, " ")
			if dtval == "" {
				return ""
			}
			if t, err := time.ParseInLocation("01/02/2006 15:04", dtval, time.Local); err != nil || dtval != t.Format("01/02/2006 15:04") {
				return fmt.Sprintf("The %q field does not contain a valid date and time in MM/DD/YYYY HH:MM format.", f.Label)
			}
			return ""
		}
	}
	return AddFieldDefaults(f)
}

// NewDateField adds defaults to a Field that are appropriate for a field that
// contains an MM/DD/YYYY date.  If paired is true, it means this field is part
// of a date/time field pair with a DateTimeField aggregator. It modifies its
// argument and returns it for chaining.
func NewDateField(paired bool, f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !PIFODateRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid date (MM/DD/YYYY).", f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareDate
	}
	if f.TableValue == nil && paired {
		f.TableValue = TableOmit
	}
	if f.EditWidth == 0 {
		f.EditWidth = 10
	}
	if f.EditHint == "" {
		f.EditHint = "MM/DD/YYYY"
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			v = strings.TrimSpace(v)
			if match := dateLooseRE.FindStringSubmatch(v); match != nil {
				// Add leading zeroes and set delimiter to slash.
				v = fmt.Sprintf("%02s/%02s/20%s", match[1], match[2], match[3])
				// Correct values that are out of range, e.g. 06/31 => 07/01.
				if t, err := time.ParseInLocation("01/02/2006", v, time.Local); err == nil {
					v = t.Format("01/02/2006")
				}
			}
			*f.Value = v
		}
	}
	if f.EditSkip == nil && paired {
		f.EditSkip = EditSkipAlways
	}
	return NewAggregatorField(f)
}

// NewFCCCallSignField adds defaults to a Field that are appropriate for a field
// that contains an FCC call sign.  It modifies its argument and returns it for
// chaining.
func NewFCCCallSignField(f *Field) *Field {
	if f.Compare == nil {
		f.Compare = CompareExact
	}
	if f.EditWidth == 0 {
		f.EditWidth = 6
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			*f.Value = strings.ToUpper(strings.TrimSpace(v))
		}
	}
	if f.EditValid == nil {
		f.EditValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !fccCallSignRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid FCC call sign.", f.Label)
			}
			return ""
		}
	}
	return AddFieldDefaults(f)
}

// NewFrequencyField adds defaults to a Field that are appropriate for a field
// that contains a frequency in MHz.  It modifies its argument and returns it
// for chaining.
func NewFrequencyField(f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !PIFOFrequencyRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid frequency.", f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareReal
	}
	if f.EditHint == "" {
		f.EditHint = "MHz"
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			v = strings.TrimSpace(v)
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				v = strconv.FormatFloat(n, 'f', -1, 64)
			}
			*f.Value = v
		}
	}
	return AddFieldDefaults(f)
}

// NewFrequencyOffsetField adds defaults to a Field that are appropriate for a
// field that contains a repeater frequency offset.  It modifies its argument
// and returns it for chaining.
func NewFrequencyOffsetField(f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !PIFOFrequencyOffsetRE.MatchString(*f.Value) {
				return fmt.Sprintf(`The %q field does not contain a valid frequency offset (a real number, a "+", or a "-").`, f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareReal
	}
	if f.EditHint == "" {
		f.EditHint = "MHz or +/-"
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			v = strings.TrimSpace(v)
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				v = strconv.FormatFloat(n, 'f', -1, 64)
			}
			*f.Value = v
		}
	}
	return AddFieldDefaults(f)
}

// NewMessageNumberField adds defaults to a Field that are appropriate for a
// field that contains a packet message number.  It modifies its argument and
// returns it for chaining.
func NewMessageNumberField(f *Field) *Field {
	if f.Compare == nil {
		f.Compare = CompareNone
	}
	if f.EditWidth == 0 {
		f.EditWidth = 9
	}
	if f.EditHint == "" {
		f.EditHint = "XXX-###P"
	}
	if f.EditApply == nil {
		f.EditApply = MessageNumberEditApply
	}
	if f.EditValid == nil {
		f.EditValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !messageNumberRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid packet message number.", f.Label)
			}
			return ""
		}
	}
	return AddFieldDefaults(f)
}

func MessageNumberEditApply(f *Field, v string) {
	v = strings.ToUpper(strings.TrimSpace(v))
	if match := messageNumberLooseRE.FindStringSubmatch(v); match != nil {
		num, _ := strconv.Atoi(match[2])
		v = fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	*f.Value = v
}

// NewMultilineField adds defaults to a Field that are appropriate for a field
// containing free-form text with possible newlines.  It modifies its argument
// and returns it for chaining.
func NewMultilineField(f *Field) *Field {
	f.Multiline = true
	return NewTextField(f)
}

// NewPhoneNumberField adds defaults to a Field that are appropriate for a field
// that contains a phone number.  It modifies its argument and returns it for
// chaining.
func NewPhoneNumberField(f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !PIFOPhoneNumberRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid phone number.", f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareExact
	}
	return AddFieldDefaults(f)
}

// NewRealNumberField adds defaults to a Field that are appropriate for a field
// that contains a real number.  It modifies its argument and returns it for
// chaining.
func NewRealNumberField(f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !PIFORealNumberRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid number.", f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareReal
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			v = strings.TrimSpace(v)
			if n, err := strconv.ParseFloat(v, 64); err == nil {
				v = strconv.FormatFloat(n, 'f', -1, 64)
			}
			*f.Value = v
		}
	}
	return AddFieldDefaults(f)
}

// NewRestrictedField adds defaults to a Field that are appropriate for a field
// that can contain only a restricted set of values.  It modifies its argument
// and returns it for chaining.
func NewRestrictedField(f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !f.Choices.IsPIFO(*f.Value) {
				return fmt.Sprintf("The %q field does not contain one of its allowed values.", f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareExact
	}
	return AddFieldDefaults(f)
}

// NewStaticPDFContentField adds defaults to a Field that are appropriate for a
// pseudo-field that renders static content in a generated PDF.  It modifies its
// argument and returns it for chaining.
func NewStaticPDFContentField(f *Field) *Field {
	return AddFieldDefaults(f)
}

// NewTacticalCallSignField adds defaults to a Field that are appropriate for a
// field that contains a tactical call sign.  It modifies its argument and
// returns it for chaining.
func NewTacticalCallSignField(f *Field) *Field {
	if f.Compare == nil {
		f.Compare = CompareExact
	}
	if f.EditWidth == 0 {
		f.EditWidth = 6
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			*f.Value = strings.ToUpper(strings.TrimSpace(v))
		}
	}
	if f.EditValid == nil {
		f.EditValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !tacticalCallSignRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid tactical call sign.", f.Label)
			}
			return ""
		}
	}
	return AddFieldDefaults(f)
}

// NewTextField adds defaults to a Field that are appropriate for a field
// containing free-form text, not expected to contain newlines.  It modifies its
// argument and returns it for chaining.
func NewTextField(f *Field) *Field {
	if f.Compare == nil {
		f.Compare = CompareText
	}
	return AddFieldDefaults(f)
}

// NewTimeField adds defaults to a Field that are appropriate for a field that
// contains an HH:MM time.  If paired is true, the field is part of a date/time
// field pair with a DateTimeField aggregator.  It modifies its argument and
// returns it for chaining.
func NewTimeField(paired bool, f *Field) *Field {
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string {
			if p := f.PresenceValid(); p != "" {
				return p
			}
			if *f.Value != "" && !PIFOTimeRE.MatchString(*f.Value) {
				return fmt.Sprintf("The %q field does not contain a valid time (HH:MM, 24-hour clock).", f.Label)
			}
			return ""
		}
	}
	if f.Compare == nil {
		f.Compare = CompareTime
	}
	if f.TableValue == nil && paired {
		f.TableValue = TableOmit
	}
	if f.EditWidth == 0 {
		f.EditWidth = 5
	}
	if f.EditHint == "" {
		f.EditHint = "HH:MM"
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, v string) {
			v = strings.TrimSpace(v)
			if match := timeLooseRE.FindStringSubmatch(v); match != nil {
				// Add colon if needed.
				if !strings.HasSuffix(match[1], ":") {
					match[1] += ":"
				}
				// Add leading zero to hour if needed.
				v = fmt.Sprintf("%03s%s", match[1], match[2])
			}
			*f.Value = v
		}
	}
	if f.EditSkip == nil && paired {
		f.EditSkip = EditSkipAlways
	}
	return AddFieldDefaults(f)
}
