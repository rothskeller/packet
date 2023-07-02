package xscmsg

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/english"
)

// WrapAddressListField wraps a Field that contains a list of addresses.
func WrapAddressListField(f Field) Field {
	return &addressListField{f}
}

type addressListField struct{ Field }

var jnosMailboxRE = regexp.MustCompile(`(?i)^[A-Z][A-Z0-9]{0,5}$`)

func (f *addressListField) SetValue(value string) {
	if value = strings.TrimSpace(value); value != "" {
		addresses := strings.Split(value, ",")
		j := 0
		for _, address := range addresses {
			if trim := strings.TrimSpace(address); trim != "" {
				addresses[j], j = trim, j+1
			}
		}
		value = strings.Join(addresses[:j], ", ")
	}
	f.Field.SetValue(value)
}
func (f addressListField) Size() (width, height int) { return 80, 1 }
func (f addressListField) Problem() string {
	if value := f.Value(); value != "" {
		addresses := strings.Split(value, ", ")
		for _, address := range addresses {
			if jnosMailboxRE.MatchString(address) {
				continue
			}
			if _, err := mail.ParseAddress(address); err != nil {
				return fmt.Sprintf("The %q field contains %q, which is not a valid JNOS mailbox name, BBS network address, or email address.", f.Label(), address)
			}
		}
	}
	return ""
}
func (f addressListField) Help() string {
	return f.Field.Help() + "  Each address must be a JNOS mailbox name, a BBS network address, or an email address.  The addresses must be separated by commas."
}

////////////////////////////////////////////////////////////////////////////////

// CheckboxField creates a Field for a checkbox control.  It must be wrapped by
// another type that provides Label and Help implementations.
func CheckboxField(vp *bool) Field {
	return &checkboxField{vp}
}

type checkboxField struct{ vp *bool }

func (f checkboxField) Label() string {
	panic("CheckboxField.Label must be overridden")
}
func (f checkboxField) Value() string {
	if *f.vp {
		return "checked"
	}
	return ""
}
func (f *checkboxField) SetBool(value bool) { *f.vp = value }
func (f *checkboxField) SetValue(value string) {
	switch strings.ToLower(value) {
	case "", "f", "false", "n", "no":
		*f.vp = false
	default:
		*f.vp = true
	}
}
func (f *checkboxField) Size() (width int, height int) { return 3, 1 }
func (f *checkboxField) Problem() string               { return "" }
func (f *checkboxField) Choices() []string             { return nil }
func (f checkboxField) Help() string {
	panic("CheckboxField.Help must be overridden")
}
func (f *checkboxField) Hint() string { return "" }

////////////////////////////////////////////////////////////////////////////////

// WrapDateField wraps a Field that should contain a date.
func WrapDateField(f Field) Field {
	return &dateField{f}
}

type dateField struct{ Field }

var dateLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)

func (f *dateField) SetValue(value string) {
	if match := dateLooseRE.FindStringSubmatch(value); match != nil {
		// Add leading zeroes and set delimiter to slash.
		value = fmt.Sprintf("%02s/%02s/20%s", match[1], match[2], match[3])
		// Correct values that are out of range, e.g. 06/31 => 07/01.
		if t, err := time.ParseInLocation("01/02/2006", value, time.Local); err == nil {
			value = t.Format("01/02/2006")
		}
	}
	f.Field.SetValue(value)
}
func (f dateField) Size() (width, height int) { return 10, 1 }
func (f dateField) Problem() string {
	if value := f.Value(); value != "" && !ValidDate(value) {
		return fmt.Sprintf("The %q field does not contain a valid date.", f.Label())
	}
	return f.Field.Problem()
}
func (f dateField) Help() string {
	return f.Field.Help() + "  The date must be in MM/DD/YYYY format."
}
func (f dateField) Hint() string { return "MM/DD/YYYY" }

////////////////////////////////////////////////////////////////////////////////

// WrapFCCCallSignField wraps a Field that should contain an FCC call sign.
func WrapFCCCallSignField(f Field) Field {
	return &fccCallSignField{f}
}

type fccCallSignField struct{ Field }

func (f *fccCallSignField) SetValue(value string) {
	f.Field.SetValue(strings.ToUpper(value))
}
func (f *fccCallSignField) Size() (w, h int) { return 6, 1 }
func (f *fccCallSignField) Problem() string {
	if value := f.Value(); value != "" && !ValidFCCCallSign(value) {
		return fmt.Sprintf("The %q field does not contain a valid FCC call sign.", f.Label())
	}
	return f.Field.Problem()
}

////////////////////////////////////////////////////////////////////////////////

// WrapMessageIDField wraps a Field that must contain a message ID.  The message
// ID need not be a packet message ID; for that, see WrapPacketMessageIDField.
func WrapMessageIDField(f Field) Field { return &messageIDField{f} }

type messageIDField struct{ Field }

var messageNumberLooseRE = regexp.MustCompile(`^([A-Z0-9]{3})-(\d+)([A-Z]?)$`)

func (f *messageIDField) SetValue(value string) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if match := messageNumberLooseRE.FindStringSubmatch(value); match != nil {
		num, _ := strconv.Atoi(match[2])
		value = fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	f.Field.SetValue(value)
}
func (f messageIDField) Size() (width, height int) { return 10, 1 }
func (f messageIDField) Problem() string {
	if value := f.Value(); value == "" || ValidMessageID(value) {
		return f.Field.Problem()
	}
	return fmt.Sprintf("The %q field does not contain a valid message number.", f.Label())
}
func (f messageIDField) Help() string {
	return f.Field.Help() + "  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter."
}

////////////////////////////////////////////////////////////////////////////////

// WrapPacketMessageIDField wraps a Field that must contain a packet message ID.
func WrapPacketMessageIDField(f Field) Field { return &packetMessageIDField{WrapMessageIDField(f)} }

type packetMessageIDField struct{ Field }

func (f packetMessageIDField) Problem() string {
	if value := f.Value(); value == "" || ValidPacketMessageID(value) {
		return f.Field.Problem()
	}
	return fmt.Sprintf("The %q field does not contain a valid packet message number.", f.Label())
}
func (f packetMessageIDField) Help() string {
	return f.Field.Help() + "  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is a suffix letter (usually P, but sometimes M or R)."
}

////////////////////////////////////////////////////////////////////////////////

// WrapRequiredField wraps a Field and makes it required (i.e., it must have a
// non-empty value).
func WrapRequiredField(f Field) Field {
	return &requiredField{f}
}

type requiredField struct {
	Field
}

func (f *requiredField) Problem() string {
	if f.Value() == "" {
		return fmt.Sprintf("The %q field is required.", f.Label())
	}
	return f.Field.Problem()
}
func (f *requiredField) Help() string {
	return f.Field.Help() + "  This field is required."
}

////////////////////////////////////////////////////////////////////////////////

// WrapRestrictedField wraps a Field and enforces that, if it contains any value
// at all, the value it contains is one of the values returned by the Choices
// method.
func WrapRestrictedField(f Field) Field {
	return &restrictedField{f}
}

type restrictedField struct{ Field }

func (f *restrictedField) SetValue(value string) {
	value = strings.TrimSpace(value)
	var match string
	for _, c := range f.Choices() {
		if len(value) <= len(c) && strings.EqualFold(value, c[:len(value)]) {
			if match != "" {
				match = ""
				break
			}
			match = c
		}
	}
	if match != "" {
		value = match
	}
	f.Field.SetValue(value)
}
func (f *restrictedField) Size() (width, height int) {
	for _, c := range f.Choices() {
		if len(c) > width {
			width = len(c)
		}
	}
	return width, 1
}
func (f *restrictedField) Problem() string {
	if value := f.Value(); value == "" || ValidRestrictedValue(value, f.Choices()) {
		return f.Field.Problem()
	}
	return fmt.Sprintf("The %q field does not contain an allowed value.  Allowed values are %s.",
		f.Label(), english.Conjoin(f.Choices(), "and"))
}
func (f *restrictedField) Help() string {
	return fmt.Sprintf("%s  Allowed values are %s.", f.Field.Help(), english.Conjoin(f.Choices(), "and"))
}

////////////////////////////////////////////////////////////////////////////////

// WrapTacticalCallSignField wraps a Field that should contain a tactical call
// sign.
func WrapTacticalCallSignField(f Field) Field {
	return &tacticalCallSignField{f}
}

type tacticalCallSignField struct{ Field }

func (f *tacticalCallSignField) SetValue(value string) {
	f.Field.SetValue(strings.ToUpper(value))
}
func (f *tacticalCallSignField) Size() (w, h int) { return 6, 1 }
func (f *tacticalCallSignField) Problem() string {
	if value := f.Value(); value != "" && !ValidCallSign(value) {
		return fmt.Sprintf("The %q field does not contain a valid tactical call sign.", f.Label())
	}
	return f.Field.Problem()
}
func (f *tacticalCallSignField) Help() string {
	return f.Field.Help() + "  Valid tactical call signs contain three to six letters or digits, starting with a letter."
}

////////////////////////////////////////////////////////////////////////////////

// WrapTimeField wraps a Field that should contain a time.
func WrapTimeField(f Field) Field {
	return &timeField{f}
}

type timeField struct{ Field }

var timeLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-9]|2[0-3]):([0-5][0-9])$`)

func (f *timeField) SetValue(value string) {
	if match := timeLooseRE.FindStringSubmatch(value); match != nil {
		// Add leading zero to hour if needed.
		value = fmt.Sprintf("%02s:%s", match[1], match[2])
	}
	f.Field.SetValue(value)
}
func (f timeField) Size() (width, height int) { return 5, 1 }
func (f timeField) Problem() string {
	if value := f.Value(); value != "" && !ValidTime(value) {
		return fmt.Sprintf("The %q field does not contain a valid time.", f.Label())
	}
	return f.Field.Problem()
}
func (f timeField) Help() string {
	return f.Field.Help() + "  The time must be in HH:MM notation (24-hour clock)."
}
func (f timeField) Hint() string { return "HH:MM" }
