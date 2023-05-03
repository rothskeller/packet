package xscmsg

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/pktmsg"
)

type bodyField struct{ pktmsg.SettableField }

func (f bodyField) Label() string                 { return "Body" }
func (f bodyField) Size() (width int, height int) { return 80, 10 }
func (f bodyField) Help(m pktmsg.Message) string {
	return "This is the body of the message.  It is required."
}
func (f bodyField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The message body is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// A DateField is a base implementation for a pktmsg.SettableField that contains
// a date.
type DateField struct{ StringField }

var pifoDateRE = regexp.MustCompile(`^(0[1-9]|1[012])/(0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]$`)

// Size returns the size of a date field.
func (DateField) Size() (width, height int) { return 10, 1 }

// Hint returns a hint about the format of a date.
func (DateField) Hint() string { return "MM/DD/YYYY" }

// PHelp returns date format help that the containing field can add to its Help
// string.
func (DateField) PHelp() string {
	return "The date must be in MM/DD/YYYY format."
}

// SetValue sets the string value of the date field.
func (f *DateField) SetValue(value string) {
	// Zero-pad the month and day, and extend the year to four digits.
	value = strings.TrimSpace(value)
	t, err := time.ParseInLocation(value, "1/2/2006", time.Local)
	if err != nil {
		t, err = time.ParseInLocation(value, "1/2/06", time.Local)
	}
	if err == nil {
		value = t.Format("01/02/2006")
	}
	f.StringField.SetValue(value)
}

// PValidate checks the validity of the date value.  The caller should prefix
// the returned problem (if any) with text identifying the field.
func (f DateField) PValidate(pifo bool) string {
	if f.Value() == "" {
		return ""
	}
	if pifo {
		if !pifoDateRE.MatchString(f.Value()) {
			return "The date is not in MM/DD/YYYY format."
		}
	} else if t, err := time.ParseInLocation("01/02/2006", f.Value(), time.Local); err != nil {
		return "The date is not in MM/DD/YYYY format."
	} else if t.Format("01/02/2006") != f.Value() {
		return "The date is not a valid date."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// DestinationMsgIDField implements the XSC-standard destination message ID
// field.
type DestinationMsgIDField struct{ MessageIDField }

// Label returns the field label.
func (DestinationMsgIDField) Label() string { return "Destination Message Number" }

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f DestinationMsgIDField) Validate(_ pktmsg.Message, pifo bool) string {
	if problem := f.MessageIDField.PValidate(pifo); problem != "" {
		return "The destination message number is not valid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// An FCCCallSignField is a base implementation for a field that contains an FCC
// call sign.
type FCCCallSignField struct{ StringField }

var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)

// Size returns the size of a call sign field.
func (FCCCallSignField) Size() (width, height int) { return 6, 1 }

// SetValue sets the string value of the field.
func (f *FCCCallSignField) SetValue(value string) {
	f.StringField.SetValue(strings.ToUpper(value))
}

// PValidate checks the validity of the call sign.  The containing field should
// prefix any returned problem with identification of the field.
func (f FCCCallSignField) PValidate(pifo bool) string {
	if !pifo && f.Value() != "" && !fccCallSignRE.MatchString(f.Value()) {
		return "It is not a valid FCC call sign."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

type fromAddrField struct{ pktmsg.FromAddrField }

func (f fromAddrField) Label() string { return "From" }

////////////////////////////////////////////////////////////////////////////////

// FromContactField implements the XSC-standard "From Contact" field.
type FromContactField struct{ StringField }

// FromContactTag is the tag for the XSC-standard "From Contact" field.
const FromContactTag = "8d."

// Label returns the field label.
func (FromContactField) Label() string { return "From Contact" }

// Help returns the help text for the field.
func (FromContactField) Help(_ pktmsg.Message) string {
	return "This is the contact information (e.g., telephone number) of the person who wrote the message."
}

////////////////////////////////////////////////////////////////////////////////

// FromICSPositionField implements the XSC-standard "From ICS Position" field.
type FromICSPositionField struct{ StringField }

// FromICSPositionTag is the tag for the XSC-standard "From ICS Position" field.
const FromICSPositionTag = "8a."

// Label returns the field label.
func (FromICSPositionField) Label() string { return "From ICS Position" }

// Help returns the help text for the field.
func (FromICSPositionField) Help(_ pktmsg.Message) string {
	return "This is the ICS position (section, branch, unit, etc.) from which the message was sent.  It is required."
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f FromICSPositionField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The \"From ICS Position\" field is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// FromLocationField implements the XSC-standard "From Location" field.
type FromLocationField struct{ StringField }

// FromLocationTag is the tag for the XSC-standard "From Location" field.
const FromLocationTag = "8b."

// Label returns the field label.
func (FromLocationField) Label() string { return "From Location" }

// Help returns the help text for the field.
func (FromLocationField) Help(_ pktmsg.Message) string {
	return "This is the location (EOC, DOC, JOC, etc.) from which the message was sent.  It is required."
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f FromLocationField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The \"From Location\" field is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// FromNameField implements the XSC-standard "From Name" field.
type FromNameField struct{ StringField }

// FromNameTag is the tag for the XSC-standard "From Name" field.
const FromNameTag = "8c."

// Label returns the field label.
func (FromNameField) Label() string { return "From Name" }

// Help returns the help text for the field.
func (FromNameField) Help(_ pktmsg.Message) string {
	return "This is the name of the person who wrote the message."
}

////////////////////////////////////////////////////////////////////////////////

// HandlingField implements the XSC-standard message handling order field.
type HandlingField struct{ SelectionField }

// HandlingTag is the tag for the XSC-standard message handling order field.
const HandlingTag = "5."

var handlingValues = []string{"ROUTINE", "PRIORITY", "IMMEDIATE"}

// Label returns the field label.
func (HandlingField) Label() string { return "Handling" }

// Size returns the display size of the field.
func (f HandlingField) Size() (width, height int) { return f.SelectionField.Size(handlingValues) }

// SetValue sets the value of the field.
func (f *HandlingField) SetValue(value string) { f.SelectionField.SetValue(value, handlingValues) }

// Choices returns the set of allowed values for the field.
func (f HandlingField) Choices() []string { return handlingValues }

// Help returns the help text for the field.
func (f HandlingField) Help(_ pktmsg.Message) string {
	return "This is the handling order for the message.  It is required.  " + f.SelectionField.PHelp(handlingValues)
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f HandlingField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The message handling order is required."
	}
	if problem := f.SelectionField.PValidate(handlingValues); problem != "" {
		return "The message handling order is not valid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// MessageDateField implements the XSC-standard message date field.
type MessageDateField struct{ DateField }

// MessageDateTag is the tag for the XSC-standard message date field.
const MessageDateTag = "1a."

// Label returns the field label.
func (MessageDateField) Label() string { return "Date" }

// Default returns the default value for the field.
func (MessageDateField) Default() string { return time.Now().Format("01/02/2006") }

// Help returns the help text for the field.
func (f MessageDateField) Help(_ pktmsg.Message) string {
	return "This is the date on which the message was written (as opposed to when it was sent).  It is required." + f.DateField.PHelp()
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f MessageDateField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The message date is required."
	}
	if problem := f.DateField.PValidate(pifo); problem != "" {
		return "The message date is not valid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// A MessageIDField is a base implementation for settable fields that contain a
// Santa Clara County standard packet message ID.
type MessageIDField struct{ StringField }

var messageIDRE = regexp.MustCompile(`^(?:[0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-\d{3,}[PMR]$`)
var messageIDLooseRE = regexp.MustCompile(`^([0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-(\d+)([PMR])?$`)

// SetValue sets the string value of the message ID.
func (f *MessageIDField) SetValue(value string) {
	// Uppercase the string, and format the sequence number with the minimal
	// number of digits >= 3.
	value = strings.ToUpper(strings.TrimSpace(value))
	if match := messageIDLooseRE.FindStringSubmatch(value); match != nil {
		num, _ := strconv.Atoi(match[2])
		value = fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	f.StringField.SetValue(value)
}

// Size returns the size of a message ID field.
func (MessageIDField) Size() (width, height int) { return 10, 1 }

// PHelp returns message ID help text that the containing Field can add to its
// Help string.
func (MessageIDField) PHelp() string {
	return "Message numbers have the form XXX-###P, where XXX is a three-character prefix assigned to the origin station, ### is a sequence number (any number of digits), and P is a suffix character (usually 'P', but sometimes 'M' or 'R')."
}

// PValidate checks the validity of a message ID.  The containing field should
// prefix any returned problem with identification of the field.
func (f MessageIDField) PValidate(pifo bool) string {
	if !pifo && f.Value() != "" && !messageIDRE.MatchString(f.Value()) {
		return "Message numbers have the form XXX-###P, where XXX is a three-character prefix, ### is a sequence number (three or more digits), and P is a suffix character ('P', 'M', or 'R')."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// MessageTimeField implements the XSC-standard message time field.
type MessageTimeField struct{ TimeField }

// MessageTimeTag is the tag for the XSC-standard message time field.
const MessageTimeTag = "1b."

// Label returns the field label.
func (MessageTimeField) Label() string { return "Time" }

// Help returns the help text for the field.
func (f MessageTimeField) Help(_ pktmsg.Message) string {
	return "This is the time at which the message was written (as opposed to when it was sent).  It is required." + f.TimeField.PHelp()
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f MessageTimeField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The message time is required."
	}
	if problem := f.TimeField.PValidate(pifo); problem != "" {
		return "The message time is not valid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// OpCallField implements the XSC-standard operator call sign field.
type OpCallField struct{ FCCCallSignField }

// OpCallTag is the tag for the XSC-standard operator call sign field.
const OpCallTag = "OpCall"

// Label returns the field label.
func (OpCallField) Label() string { return "Operator Call Sign" }

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f OpCallField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The operator call sign is required."
	}
	return f.FCCCallSignField.PValidate(pifo)
}

////////////////////////////////////////////////////////////////////////////////

// OpDateField implements the XSC-standard operator date field.
type OpDateField struct{ DateField }

// OpDateTag is the tag for the XSC-standard operator date field.
const OpDateTag = "OpDate"

// Label returns the field label.
func (OpDateField) Label() string { return "Operator Date" }

// Default returns the default value for the field.
func (OpDateField) Default() string { return time.Now().Format("01/02/2006") }

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f OpDateField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The operator date field is required."
	}
	if problem := f.DateField.PValidate(pifo); problem != "" {
		return "The operator date field is invalid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// OpNameField implements the XSC-standard operator name field.
type OpNameField struct{ StringField }

// OpNameTag is the tag for the XSC-standard operator name field.
const OpNameTag = "OpName"

// Label returns the field label.
func (OpNameField) Label() string { return "Operator Name" }

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f OpNameField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The operator name field is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// OpRelayRcvdField implements the XSC-standard operator relay received field.
type OpRelayRcvdField struct{ StringField }

// OpRelayRcvdTag is the tag for the XSC-standard operator relay received field.
const OpRelayRcvdTag = "OpRelayRcvd"

// Label returns the field label.
func (OpRelayRcvdField) Label() string { return "Relay Rcvd" }

// Help returns the help text for the field.
func (OpRelayRcvdField) Help(_ pktmsg.Message) string {
	return "This is the name of the relay station from which this message was received, if it was not received directly from the origin station."
}

////////////////////////////////////////////////////////////////////////////////

// OpRelaySentField implements the XSC-standard operator relay sent field.
type OpRelaySentField struct{ StringField }

// OpRelaySentTag is the tag for the XSC-standard operator relay sent field.
const OpRelaySentTag = "OpRelaySent"

// Label returns the field label.
func (OpRelaySentField) Label() string { return "Relay Sent" }

// Help returns the help text for the field.
func (OpRelaySentField) Help(_ pktmsg.Message) string {
	return "This is the name of the relay station to which this message was sent, if it was not sent directly to the destination station."
}

////////////////////////////////////////////////////////////////////////////////

// OpTimeField implements the XSC-standard operator time field.
type OpTimeField struct{ TimeField }

// OpTimeTag is the tag for the XSC-standard operator time field.
const OpTimeTag = "OpTime"

// Label returns the field label.
func (OpTimeField) Label() string { return "Operator Time" }

// Default returns the default value for the field.
func (OpTimeField) Default() string { return time.Now().Format("15:04") }

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f OpTimeField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The operator time field is required."
	}
	if problem := f.TimeField.PValidate(pifo); problem != "" {
		return "The operator time field is invalid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// OriginMsgIDField implements the XSC-standard origin message ID field.
type OriginMsgIDField struct{ MessageIDField }

// OriginMsgIDTag is the tag for the XSC-standard origin message ID field.
const OriginMsgIDTag = "MsgNo"

// Label returns the field label.
func (OriginMsgIDField) Label() string { return "Origin Message Number" }

// Help returns the help text for the field.
func (f OriginMsgIDField) Help(_ pktmsg.Message) string {
	return "This is the message number assigned by the origin station.  It is required.  " + f.MessageIDField.PHelp()
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f OriginMsgIDField) Validate(_ pktmsg.Message, pifo bool) string {
	if f.Value() == "" {
		return "The origin message number is required."
	}
	if problem := f.MessageIDField.PValidate(pifo); problem != "" {
		return "The origin message number is not valid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

type retrievedField struct {
	rxBBS  pktmsg.RxBBSField
	rxArea pktmsg.RxAreaField
	rxDate pktmsg.RxDateField
}

func (f retrievedField) Value() string {
	var date = f.rxDate.Time().Format("01/02/2006 15:04")
	if f.rxArea != nil {
		return fmt.Sprintf("%s from %s on %s", date, f.rxArea.Value(), f.rxBBS.Value())
	}
	return fmt.Sprintf("%s from %s", date, f.rxBBS.Value())
}
func (f retrievedField) Label() string                 { return "Retrieved" }
func (f retrievedField) Size() (width int, height int) { return 44, 1 }

////////////////////////////////////////////////////////////////////////////////

// A SelectionField is a base implementation for a settable field that can only
// contain one of a specific set of allowed values.
type SelectionField struct{ StringField }

// SetValue sets the string value of the selection field.
func (f *SelectionField) SetValue(value string, allowed []string) {
	value = strings.TrimSpace(value)
	var match string
	for _, c := range allowed {
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
	f.StringField.SetValue(value)
}

// Size returns the size of longest allowed value in the field.
func (f SelectionField) Size(allowed []string) (width, height int) {
	for _, c := range allowed {
		if len(c) > width {
			width = len(c)
		}
	}
	return width, 1
}

// PHelp returns help text that the containing field can add to its Help string.
// Essentially, it formats the list of allowed values.
func (f SelectionField) PHelp(allowed []string) string {
	switch len(allowed) {
	case 0, 1:
		panic("SelectionField with < 2 Choices")
	case 2:
		return fmt.Sprintf("The allowed values are %s and %s.", allowed[0], allowed[1])
	default:
		return fmt.Sprintf("The allowed values are %s, and %s.", strings.Join(allowed[:len(allowed)-1], ", "), allowed[len(allowed)-1])
	}
}

// PValidate checks the validity of the selection field value.  The containing
// field should prefix any returned problem with identification of the field.
func (f SelectionField) PValidate(allowed []string) string {
	var v = f.Value()
	if v == "" {
		return ""
	}
	for _, c := range allowed {
		if c == v {
			return ""
		}
	}
	return fmt.Sprintf("%q is not one of the allowed values.  %s", v, f.PHelp(allowed))
}

////////////////////////////////////////////////////////////////////////////////

type sentDateField struct{ pktmsg.SentDateField }

func (f sentDateField) Value() string {
	if f.SentDateField.Value() == "" {
		return ""
	}
	return f.SentDateField.Time().Format("01/02/2006 15:04")
}
func (f sentDateField) Label() string                 { return "Sent" }
func (f sentDateField) Size() (width int, height int) { return 16, 1 }

////////////////////////////////////////////////////////////////////////////////

// SeverityField implements a read-only field for access to the obsolete
// situation severity code embedded in the message subject, if any.
type SeverityField struct{ SelectionField }

var severityValues = []string{"OTHER", "URGENT", "EMERGENCY"}

// Label returns the field label.
func (SeverityField) Label() string { return "Situation Severity" }

// Size returns the display size of the field.
func (f SeverityField) Size() (width, height int) { return f.SelectionField.Size(severityValues) }

// Choices returns the set of allowed values for the field.
func (SeverityField) Choices() []string { return severityValues }

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f SeverityField) Validate(_ pktmsg.Message, _ bool) string {
	if problem := f.SelectionField.PValidate(severityValues); problem != "" {
		return "The situation severity is not valid.  " + problem
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// A StringField is a pktmsg.SettableField that contains a string.  It is the
// base implementation for nearly all Fields.
type StringField string

// Value returns the value of the Field as a string.
func (f StringField) Value() string { return string(f) }

// SetValue sets the string value of the Field.
func (f *StringField) SetValue(value string) { *f = StringField(strings.TrimSpace(value)) }

////////////////////////////////////////////////////////////////////////////////

type SubjectField struct{ StringField }

func (f SubjectField) Label() string { return "Subject" }
func (f SubjectField) Help(m pktmsg.Message) string {
	return "This is the subject of the message.  It is required."
}
func (f SubjectField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The message subject is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// A TacCallSignField is a base implementation of a field that contains a
// tactical call sign.
type TacCallSignField struct{ StringField }

var tacCallSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{0,5}$`)

// Size returns the size of a call sign field.
func (TacCallSignField) Size() (width, height int) { return 6, 1 }

// SetValue sets the string value of the field.
func (f *TacCallSignField) SetValue(value string) {
	f.StringField.SetValue(strings.ToUpper(value))
}

// PValidate checks the validity of the call sign.  The containing field should
// prefix any returned problem with identification of the field.
func (f TacCallSignField) PValidate(pifo bool) string {
	if !pifo && f.Value() != "" && !tacCallSignRE.MatchString(f.Value()) {
		return "It is not a valid tactical call sign."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// A TimeField is a base implementation for a settable field that contains a
// time of day.
type TimeField struct{ StringField }

var pifoTimeRE = regexp.MustCompile(`^([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00$`)
var timeLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-9]|2[0-4]):?([0-5][0-9])$`)

// Size returns the size of a time field.
func (TimeField) Size() (width, height int) { return 5, 1 }

// Hint returns a hint about the format of a time value.
func (TimeField) Hint() string { return "HH:MM" }

// PHelp returns time formatting help text that the containing field can add to
// its Help string.
func (TimeField) PHelp() string {
	return "The time must be in HH:MM format (24-hour clock)."
}

// SetValue sets the string value of the time field.
func (f TimeField) SetValue(value string) {
	// Zero-pad the hour, and add a colon if missing.
	value = strings.TrimSpace(value)
	if match := timeLooseRE.FindStringSubmatch(value); match != nil {
		h := match[1]
		if len(h) == 1 {
			h = "0" + h
		}
		value = h + ":" + match[2]
	}
	f.StringField.SetValue(value)
}

// PValidate checks the validity of the time value.  The containing field should
// prefix any returned problem with identification of the field.
func (f TimeField) PValidate(_ bool) string {
	if v := f.Value(); v != "" && !pifoTimeRE.MatchString(v) {
		return "The time is not in HH:MM format."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

type toAddrsField struct{ pktmsg.ToAddrsField }

// A mailbox name in an address list must start with a letter and contain 4 to 6
// letters or digits.
var mailboxRE = regexp.MustCompile(`^(?i)[A-Z][A-Z0-9]{3,5}$`)

func (f *toAddrsField) SetValue(value string) {
	// Remove empty fields and rationalize separators to ", ".
	addrs := strings.Split(value, ",")
	j := 0
	for _, a := range addrs {
		if a := strings.TrimSpace(a); a != "" {
			addrs[j] = a
			j++
		}
	}
	addrs = addrs[:j]
	f.ToAddrsField.SetValue(strings.Join(addrs, ", "))
}
func (f toAddrsField) Label() string { return "To" }
func (f toAddrsField) Help(m pktmsg.Message) string {
	return "This is the comma-separated list of addresses to which the message should be sent.  At least one address is required.  Each address must be a JNOS mailbox name, a valid BBS network address, or a valid email address."
}
func (f toAddrsField) Validate(_ pktmsg.Message, pifo bool) string {
	if pifo {
		return ""
	}
	if f.Value() == "" {
		return "The message must have at least one To: address."
	}
	addrs := strings.Split(f.Value(), ", ")
	for _, a := range addrs {
		if !mailboxRE.MatchString(a) {
			if _, err := mail.ParseAddress(a); err != nil {
				return fmt.Sprintf("The To: address %q is not a valid address.", a)
			}
		}
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// ToContactField implements the XSC-standard "To Contact" field.
type ToContactField struct{ StringField }

// ToContactTag is the tag for the XSC-standard "To Contact" field.
const ToContactTag = "7d."

// Label returns the field label.
func (ToContactField) Label() string { return "To Contact" }

// Help returns the help text for the field.
func (ToContactField) Help(_ pktmsg.Message) string {
	return "This is the contact information (e.g., telephone number) of the person to whom the missed is addressed.  It is optional and rarely specified."
}

////////////////////////////////////////////////////////////////////////////////

// ToICSPositionField implements the XSC-standard "To ICS Position" field.
type ToICSPositionField struct{ StringField }

// ToICSPositionTag is the tag for the XSC-standard "To ICS Position" field.
const ToICSPositionTag = "7a."

// Label returns the field label.
func (ToICSPositionField) Label() string { return "To ICS Position" }

// Help returns the help text for the field.
func (ToICSPositionField) Help(_ pktmsg.Message) string {
	return "This is the ICS position (section, branch, unit, etc.) to which the message is addressed.  It is required."
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f ToICSPositionField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The \"To ICS Position\" field is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// ToLocationField implements the XSC-standard "To Location" field.
type ToLocationField struct{ StringField }

// ToLocationTag is the tag for the XSC-standard "To Location" field.
const ToLocationTag = "7b."

// Label returns the field label.
func (ToLocationField) Label() string { return "To Location" }

// Help returns the help text for the field.
func (ToLocationField) Help(_ pktmsg.Message) string {
	return "This is the location (EOC, DOC, JOC, etc.) to which the message is addressed.  It is required."
}

// Validate validates the value of the field.  It returns a string describing
// the problem with the value, or an empty string if the value is valid.
func (f ToLocationField) Validate(_ pktmsg.Message, _ bool) string {
	if f.Value() == "" {
		return "The \"To Location\" field is required."
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////

// ToNameField implements the XSC-standard "To Name" field.
type ToNameField struct{ StringField }

// ToNameTag is the tag for the XSC-standard "To Name" field.
const ToNameTag = "7c."

// Label returns the field label.
func (ToNameField) Label() string { return "To Name" }

// Help returns the help text for the field.
func (ToNameField) Help(_ pktmsg.Message) string {
	return "This is the name of the person to whom the message is addressed.  It is optional and rarely specified; addressing to an ICS position alone is typical."
}
