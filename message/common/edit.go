package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
)

// StdFieldsEdit is the set of edit fields for the standard form fields.
type StdFieldsEdit struct {
	OriginMsgID     message.EditField
	DateTime        message.EditField
	Handling        message.EditField
	ToICSPosition   message.EditField
	ToLocation      message.EditField
	ToName          message.EditField
	ToContact       message.EditField
	FromICSPosition message.EditField
	FromLocation    message.EditField
	FromName        message.EditField
	FromContact     message.EditField
	OpRelayRcvd     message.EditField
	OpRelaySent     message.EditField
}

// StdFieldsEditTemplate is the template for edit fields for a standard form.
var StdFieldsEditTemplate = StdFieldsEdit{
	OriginMsgID:     OriginMsgIDEditField,
	DateTime:        MessageDateTimeEditField,
	Handling:        HandlingEditField,
	ToICSPosition:   ToICSPositionEditField,
	ToLocation:      ToLocationEditField,
	ToName:          ToNameEditField,
	ToContact:       ToContactEditField,
	FromICSPosition: FromICSPositionEditField,
	FromLocation:    FromLocationEditField,
	FromName:        FromNameEditField,
	FromContact:     FromContactEditField,
	OpRelayRcvd:     OpRelayRcvdEditField,
	OpRelaySent:     OpRelaySentEditField,
}

// EditFields1 returns pointers to the initial set of edit fields.
func (sfe *StdFieldsEdit) EditFields1() []*message.EditField {
	return []*message.EditField{
		&sfe.OriginMsgID,
		&sfe.DateTime,
		&sfe.Handling,
		&sfe.ToICSPosition,
		&sfe.ToLocation,
		&sfe.ToName,
		&sfe.ToContact,
		&sfe.FromICSPosition,
		&sfe.FromLocation,
		&sfe.FromName,
		&sfe.FromContact,
	}
}

// EditFields2 returns pointers to the final set of edit fields.
func (sfe *StdFieldsEdit) EditFields2() []*message.EditField {
	return []*message.EditField{
		&sfe.OpRelayRcvd,
		&sfe.OpRelaySent,
	}
}

// Common field definitions.
var (
	FromContactEditField = message.EditField{
		Label: "From Contact Info",
		Width: 29,
		Help:  "This is contact information (phone number, email, etc.) for the message author.  It is optional and rarely provided.",
	}
	FromICSPositionEditField = message.EditField{
		Label: "From ICS Position",
		Width: 30,
		Help:  "This is the ICS position of the message author.  It is required.",
	}
	FromLocationEditField = message.EditField{
		Label: "From Location",
		Width: 32,
		Help:  "This is the location of the message author.  It is required.",
	}
	FromNameEditField = message.EditField{
		Label: "From Name",
		Width: 34,
		Help:  "This is the name of the message author.  It is optional and rarely provided.",
	}
	HandlingEditField = message.EditField{
		Label:   "Handling Order",
		Width:   9,
		Choices: []string{"ROUTINE", "PRIORITY", "IMMEDIATE"},
		Help:    `This is the message handling order, which specifies how fast it needs to be delivered.  Allowed values are "ROUTINE" (within 2 hours), "PRIORITY" (within 1 hour), and "IMMEDIATE".  This field is required.`,
	}
	MessageDateTimeEditField = message.EditField{
		Label: "Message Date/Time",
		Width: 16,
		Help:  "This is the date and time the message was written, in MM/DD/YYYY HH:MM format (24-hour clock).  It is required.",
	}
	OpRelayRcvdEditField = message.EditField{
		Label: "Operator: Relay Received",
		Width: 36,
		Help:  "This is the name of the station from which this message was directly received.  It is filled in for messages that go through a relay station.",
	}
	OpRelaySentEditField = message.EditField{
		Label: "Operator: Relay Sent",
		Width: 36,
		Help:  "This is the name of the station to which this message was directly sent.  It is filled in for messages that go through a relay station.",
	}
	OriginMsgIDEditField = message.EditField{
		Label:          "Origin Message Number",
		Width:          9,
		LocalMessageID: true,
		Help:           "This is the message number assigned to the message by the origin station.  Valid message numbers have the form XXX-###P, where XXX is the three-character message number prefix assigned to the station, ### is a sequence number (any number of digits), and P is an optional suffix letter.  This field is required.",
	}
	ToContactEditField = message.EditField{
		Label: "To Contact Info",
		Width: 29,
		Help:  "This is contact information (phone number, email, etc.) for the receipient.  It is optional and rarely provided.",
	}
	ToICSPositionEditField = message.EditField{
		Label: "To ICS Position",
		Width: 30,
		Help:  "This is the ICS position to which the message is addressed.  It is required.",
	}
	ToLocationEditField = message.EditField{
		Label: "To Location",
		Width: 32,
		Help:  "This is the location of the recipient ICS position.  It is required.",
	}
	ToNameEditField = message.EditField{
		Label: "To Name",
		Width: 34,
		Help:  "This is the name of the person holding the recipient ICS position.  It is optional and rarely provided.",
	}
)

// FromEdit takes values from the edit fields, cleans them up, converts them to
// canonical form, and puts them in the message.
func (sf *StdFields) FromEdit(sfe *StdFieldsEdit) {
	sf.OriginMsgID = CleanMessageNumber(sfe.OriginMsgID.Value)
	sf.MessageDate, sf.MessageTime = CleanDateTime(sfe.DateTime.Value)
	sf.Handling = ExpandRestricted(&sfe.Handling)
	sf.ToICSPosition = strings.TrimSpace(sfe.ToICSPosition.Value)
	sf.ToLocation = strings.TrimSpace(sfe.ToLocation.Value)
	sf.ToName = strings.TrimSpace(sfe.ToName.Value)
	sf.ToContact = strings.TrimSpace(sfe.ToContact.Value)
	sf.FromICSPosition = strings.TrimSpace(sfe.FromICSPosition.Value)
	sf.FromLocation = strings.TrimSpace(sfe.FromLocation.Value)
	sf.FromName = strings.TrimSpace(sfe.FromName.Value)
	sf.FromContact = strings.TrimSpace(sfe.FromContact.Value)
	sf.OpRelayRcvd = strings.TrimSpace(sfe.OpRelayRcvd.Value)
	sf.OpRelaySent = strings.TrimSpace(sfe.OpRelaySent.Value)
}

// ToEdit takes values from the message, converts them to human form, and puts
// them in the edit fields.
func (sf *StdFields) ToEdit(sfe *StdFieldsEdit) {
	sfe.OriginMsgID.Value = sf.OriginMsgID
	sfe.DateTime.Value = SmartJoin(sf.MessageDate, sf.MessageTime, " ")
	sfe.Handling.Value = sf.Handling
	sfe.ToICSPosition.Value = sf.ToICSPosition
	sfe.ToLocation.Value = sf.ToLocation
	sfe.ToName.Value = sf.ToName
	sfe.ToContact.Value = sf.ToContact
	sfe.FromICSPosition.Value = sf.FromICSPosition
	sfe.FromLocation.Value = sf.FromLocation
	sfe.FromName.Value = sf.FromName
	sfe.FromContact.Value = sf.FromContact
	sfe.OpRelayRcvd.Value = sf.OpRelayRcvd
	sfe.OpRelaySent.Value = sf.OpRelaySent
}

// Validate validates the standard field values.
func (sfe *StdFieldsEdit) Validate() {
	if ValidateRequired(&sfe.OriginMsgID) {
		ValidateMessageNumber(&sfe.OriginMsgID)
	}
	if ValidateRequired(&sfe.DateTime) {
		ValidateDateTime(&sfe.DateTime)
	}
	if ValidateRequired(&sfe.Handling) {
		ValidateRestricted(&sfe.Handling)
	}
	ValidateRequired(&sfe.ToICSPosition)
	ValidateRequired(&sfe.ToLocation)
	ValidateRequired(&sfe.FromICSPosition)
	ValidateRequired(&sfe.FromLocation)
}

// CleanCardinal reduces any cardinal number to its canonical format.
func CleanCardinal(s string) string {
	if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
		return strconv.Itoa(n)
	}
	return s
}

// CleanCheckbox converts any non-empty string to "checked".
func CleanCheckbox(s string) string {
	if s != "" {
		return "checked"
	}
	return s
}

var dateLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)

// CleanDate accepts loose formatting of dates and converts them into strictly
// valid formatting.  Anything unusable remains unchanged.
func CleanDate(loose string) string {
	if match := dateLooseRE.FindStringSubmatch(strings.TrimSpace(loose)); match != nil {
		// Add leading zeroes and set delimiter to slash.
		value := fmt.Sprintf("%02s/%02s/20%s", match[1], match[2], match[3])
		// Correct values that are out of range, e.g. 06/31 => 07/01.
		if t, err := time.ParseInLocation("01/02/2006", value, time.Local); err == nil {
			value = t.Format("01/02/2006")
		}
		return value
	}
	return loose
}

// ValidateDate verifies that a string contains a valid date.
func ValidateDate(ef *message.EditField) {
	if t, err := time.ParseInLocation("01/02/2006", ef.Value, time.Local); err == nil && ef.Value == t.Format("01/02/2006") {
		ef.Problem = ""
	} else {
		ef.Problem = fmt.Sprintf("The %q field does not contain a valid MM/DD/YYYY date.", ef.Label)
	}
}

// CleanDateTime accepts loose formatting of date/time strings and converts them
// into strictly valid formatting.  Anything unusable remains unchanged.
func CleanDateTime(loose string) (date, time string) {
	fields := strings.Fields(loose)
	if len(fields) > 0 {
		date = CleanDate(fields[0])
	}
	if len(fields) > 1 {
		fields[1] = CleanTime(fields[1])
		time = strings.Join(fields[1:], " ")
	}
	return date, time
}

// ValidateDateTime verifies that a string contains a valid date/time string.
func ValidateDateTime(ef *message.EditField) {
	if t, err := time.ParseInLocation("01/02/2006 15:04", ef.Value, time.Local); err == nil && ef.Value == t.Format("01/02/2006 15:04") {
		ef.Problem = ""
	} else {
		ef.Problem = fmt.Sprintf("The %q field does not contain a valid date and time in MM/DD/YYYY HH:MM format.", ef.Label)
	}
}

var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)

// ValidateFCCCallSign verifies that a string contains a valid FCC call sign.
func ValidateFCCCallSign(ef *message.EditField) {
	if fccCallSignRE.MatchString(ef.Value) {
		ef.Problem = ""
	} else {
		ef.Problem = fmt.Sprintf(`The %q field does not contain a valid FCC call sign.`, ef.Label)
	}
}

var messageNumberLooseRE = regexp.MustCompile(`^([A-Z0-9]{3})-(\d+)([A-Z]?)$`)

// CleanMessageNumber accepts loose formatting of message numbers and converts
// them into strictly valid formatting.  Anything unusable remains unchanged.
func CleanMessageNumber(loose string) string {
	strict := strings.ToUpper(strings.TrimSpace(loose))
	if match := messageNumberLooseRE.FindStringSubmatch(strict); match != nil {
		num, _ := strconv.Atoi(match[2])
		return fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	return loose
}

var messageIDRE = regexp.MustCompile(`^(?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-(?:[1-9][0-9]{3,}|[0-9]{3})[A-Z]?$`)

// ValidateMessageNumber validates an edit field that is supposed to contain a
// message number.
func ValidateMessageNumber(ef *message.EditField) {
	if messageIDRE.MatchString(ef.Value) {
		ef.Problem = ""
	} else {
		ef.Problem = fmt.Sprintf("The %q field does not contain a valid packet message number.", ef.Label)
	}
}

var phoneNumberLooseReplacer = strings.NewReplacer("(", "", ") ", "-", ")", "-", ".", "-")

// CleanPhoneNumber accepts loose formatting of phone numbers and converts them
// into strictly valid formatting.  Anything unusable remains unchanged.
func CleanPhoneNumber(loose string) string {
	var value = phoneNumberLooseReplacer.Replace(strings.TrimSpace(loose))
	if PIFOPhoneNumberRE.MatchString(value) {
		return value
	}
	return loose
}

// CleanReal reduces any cardinal number to its canonical format.
func CleanReal(s string) string {
	if n, err := strconv.ParseFloat(strings.TrimSpace(s), 64); err == nil {
		return strconv.FormatFloat(n, 'f', -1, 64)
	}
	return s
}

// ValidateRequired validates a field that must contain a value.  It returns
// true if the field has a value.
func ValidateRequired(ef *message.EditField) bool {
	if ef.Value != "" {
		ef.Problem = ""
		return true
	}
	ef.Problem = fmt.Sprintf("The %q field is required.", ef.Label)
	return false
}

// ExpandRestricted expands a partially entered value for a restricted field.
// If the entered value is a case-insensitive prefix of exactly one allowed
// value of the field, that full allowed value is returned.  Otherwise the
// entered value is returned unchanged.
func ExpandRestricted(ef *message.EditField) string {
	var lc, match string

	lc = strings.ToLower(strings.TrimSpace(ef.Value))
	if lc == "" {
		return ""
	}
	for _, allowed := range ef.Choices {
		if len(allowed) >= len(lc) {
			if lc == strings.ToLower(allowed[:len(lc)]) {
				if match != "" {
					return ef.Value // multiple matches
				}
				match = allowed
			}
		}
	}
	if match != "" {
		return match
	}
	return ef.Value
}

// ValidateRestricted validates a field with a defined set of allowed values.
func ValidateRestricted(ef *message.EditField) {
	for _, c := range ef.Choices {
		if c == ef.Value {
			ef.Problem = ""
			return
		}
	}
	ef.Problem = fmt.Sprintf("The %q field does not contain an allowed value.  Allowed values are %s.",
		ef.Label, Conjoin(ef.Choices, "and"))
}

var tacticalCallSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{5}$`)

// ValidateTacticalCallSign verifies that a string contains a valid tactical call sign.
func ValidateTacticalCallSign(ef *message.EditField) {
	if tacticalCallSignRE.MatchString(ef.Value) {
		ef.Problem = ""
	} else {
		ef.Problem = fmt.Sprintf(`The %q field does not contain a valid tactical call sign.`, ef.Label)
	}
}

var timeLooseRE = regexp.MustCompile(`^([1-9]:|[01][0-9]:?|2[0-4]:?)([0-5][0-9])$`)

// CleanTime accepts loose formatting of times and converts them into strictly
// valid formatting.  Anything unusable remains unchanged.
func CleanTime(loose string) string {
	if match := timeLooseRE.FindStringSubmatch(strings.TrimSpace(loose)); match != nil {
		// Add colon if needed.
		if !strings.HasSuffix(match[1], ":") {
			match[1] += ":"
		}
		// Add leading zero to hour if needed.
		return fmt.Sprintf("%03s%s", match[1], match[2])
	}
	return loose
}

var vtimeRE = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):[0-5][0-9]$`)

// ValidateTime verifies that a string contains a valid time.
func ValidateTime(ef *message.EditField) {
	if ef.Value == "24:00" || vtimeRE.MatchString(ef.Value) {
		ef.Problem = ""
	} else {
		ef.Problem = fmt.Sprintf("The %q field does not contain a valid HH:MM time.", ef.Label)
	}
}

// Conjoin returns a set of strings joined together with the specified
// conjunction and (of course) the Oxford comma.
func Conjoin(ss []string, conj string) string {
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return ss[0]
	case 2:
		return fmt.Sprintf("%s %s %s", ss[0], conj, ss[1])
	default:
		return fmt.Sprintf("%s %s, %s", strings.Join(ss[:len(ss)-1], ", "), conj, ss[len(ss)-1])
	}
}
