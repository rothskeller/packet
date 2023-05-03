package xscmsg

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// BaseField is the base implementation of a Field.  Note that it is not a
// complete implementation.  It omits Label, Size, and Help.
type BaseField struct{ vp *string }

// NewBaseField creates a new BaseField operating on the supplied string.
func NewBaseField(vp *string) *BaseField { return &BaseField{vp} }

// Value returns the value of the field.
func (f BaseField) Value() string { return *f.vp }

// SetValue sets the value of the field.
func (f *BaseField) SetValue(value string) { *f.vp = strings.TrimSpace(value) }

// Problem returns the validation problem with the field, if any.
func (f BaseField) Problem() string { return "" }

// Choices returns the set of recommended or required values for the field.
func (f BaseField) Choices() []string { return nil }

// Hint returns a hint about what to type in the field.
func (f BaseField) Hint() string { return "" }

////////////////////////////////////////////////////////////////////////////////

// AddressesClean cleans an input address list by changing all of the separators
// to ", " and removing any empty items from the list.
func AddressesClean(value string) string {
	if value = strings.TrimSpace(value); value == "" {
		return value
	}
	addresses := strings.Split(value, ",")
	j := 0
	for _, address := range addresses {
		if trim := strings.TrimSpace(address); trim != "" {
			addresses[j], j = trim, j+1
		}
	}
	return strings.Join(addresses[:j], ", ")
}

var jnosMailboxRE = regexp.MustCompile(`(?i)^[A-Z][A-Z0-9]{0,5}$`)

// AddressesValid returns whether an address list is valid.  (An empty list is
// valid.)
func AddressesValid(value string) bool {
	if value == "" {
		return true
	}
	addresses := strings.Split(value, ", ")
	for _, address := range addresses {
		if jnosMailboxRE.MatchString(address) {
			continue
		}
		if _, err := mail.ParseAddress(address); err != nil {
			return false
		}
	}
	return true
}

// AddressesHelp is a string describing the format of an address list.
const AddressesHelp = "Each address must be a JNOS mailbox name, a BBS network address, or an email address.  The addresses must be separated by commas."

////////////////////////////////////////////////////////////////////////////////

// CallSignWidth is the width of a call sign field.  (The height is 1.)
const CallSignWidth = 6

var callSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)

// CallSignValid returns whether the input value is a valid FCC call sign.
func CallSignValid(value string) bool {
	return callSignRE.MatchString(value)
}

////////////////////////////////////////////////////////////////////////////////

// ChoicesExpand expands an input value to make it match one of the supplied
// choices.  It trims the value, and if it is a case-insensitive prefix match of
// exactly one of the supplied choices, it is replaced with that choice.
func ChoicesExpand(value string, choices []string) string {
	value = strings.TrimSpace(value)
	var match string
	for _, c := range choices {
		if len(value) <= len(c) && strings.EqualFold(value, c[:len(value)]) {
			if match != "" {
				match = ""
				break
			}
			match = c
		}
	}
	if match != "" {
		return match
	}
	return value
}

// ChoicesValid returns whether the input value is one of the supplied choices.
func ChoicesValid(value string, choices []string) bool {
	for _, c := range choices {
		if c == value {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////

// DateValid returns whether the input value is a valid date.
func DateValid(value string) bool {
	if t, err := time.ParseInLocation("01/02/2006", value, time.Local); err == nil {
		return value == t.Format("01/02/2006")
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////

var msgNumberLooseRE = regexp.MustCompile(`^([A-Z0-9]{3})-(\d+)([A-Z]?)$`)

// MessageIDClean cleans a message ID input value by trimming it, upcasing it,
// and making the sequence number the proper number of digits.
func MessageIDClean(value string) string {
	value = strings.ToUpper(strings.TrimSpace(value))
	if match := msgNumberLooseRE.FindStringSubmatch(value); match != nil {
		num, _ := strconv.Atoi(match[2])
		value = fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	return value
}

// MessageIDWidth is the width of a message ID field.  (The height is 1.)
const MessageIDWidth = 10

var msgNumberRE = regexp.MustCompile(`^(?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-(?:[1-9][0-9]{3,}|[0-9]{3})[PMR]$`)

// MessageIDValid returns whether the value is a valid message ID.
func MessageIDValid(value string) bool {
	return msgNumberRE.MatchString(value)
}

// MessageIDHelp is a string describing the syntax of a valid message number.
const MessageIDHelp = "Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to your station, ### is a sequence number (any number of digits), and P is a suffix character (usually 'P', but sometimes 'M' or 'R')."

////////////////////////////////////////////////////////////////////////////////

var timeRE = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):[0-5][0-9]$`)

// TimeValid returns whether the input value is a valid time.
func TimeValid(value string) bool {
	return timeRE.MatchString(value)
}
