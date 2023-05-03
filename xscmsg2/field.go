package xscmsg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Value is a value of a form field in PackItForms encoding.  Its underlying
// type is string, but it is a separately defined type to distinguish it from
// human-readable representations of field values (which are plain strings).
// Use FromHuman to convert a string to a Value and ToHuman to convert a Value
// to a string.
type Value string

// FieldContainer is a container of FormFields (i.e., usually a Form).
type FieldContainer interface {
	// FieldValue returns the value of the field with the specified tag.  If
	// the provider has no such field, it returns an empty string.
	FieldValue(string) Value
	// KeyValue returns the value of the field with the specified key.  If
	// the provider has no such field, it returns an empty string.
	KeyValue(FieldKey) Value
}

// FormField is the interface satisfied by all fields in all message types.
type FormField interface {
	// Container returns the container of this field.
	Container() FieldContainer
	// Tag returns the unique identifier of the field in the message type.
	// For fields included in PIFO-encoded forms, this is the identifier
	// used in the PIFO encoding.
	Tag() string
	// Key returns the well-known-field key for this field, if it is a well
	// known field.  It returns an empty string otherwise.
	Key() FieldKey
	// Label returns the English label of the field.  For form fields, it
	// should be the same as (or at least an easily recognizable
	// abbreviation of) the name of the field on the form.  But, it should
	// be no more than a few dozen characters.
	Label() string
	// Hint returns a short (one or two word) hint about the type of data
	// that should appear in the field (e.g., "MM/DD/YYYY").  It returns
	// the empty string if there is no hint.
	Hint() string
	// Help returns a full description of the meaning of the field and what
	// values are allowed in it.  It can be long and should be word-wrapped
	// for presentation.
	Help() string
	// Default returns the default value for the field.  It should be
	// assigned to the field when creating new outgoing messages.
	Default() Value
	// Editable returns whether the user is allowed to edit the value of the
	// field.  This is false for computed fields and operator-only fields.
	Editable() bool
	// Size returns a hint of the appropriate size of an input field for the
	// form, i.e., the number of characters and rows that the field on the
	// printed form has room for.  For fields with restricted/suggested
	// values, Size should always be large enough to contain any of them.
	Size() (width, height int)
	// Choices returns a list of suggested values for the field, or nil if
	// there is no such list.  Note that for some fields, the list can
	// change over time because it is derived from the current values of
	// other fields.
	Choices() []Value
	// Value returns the current value of the field.
	Value() Value
	// SetValue sets the current value of the field.  It does not validate
	// the value.
	SetValue(Value)
	// ToHuman converts a Value of the field into the corresponding value in
	// human-readable form.  If the Value is not recognized or valid,
	// ToHuman returns its input.
	ToHuman(Value) string
	// FromHuman converts a value of the field in human-readable form into
	// the corresponding V.  If it cannot do so, it returns its input.
	FromHuman(string) Value
	// Calculate updates the value of a calculated field.  It is a no-op for
	// non-calculated fields.
	Calculate()
	// Validate returns an empty string if the current value is valid as
	// human input, or a problem description string if it is not.  Note that
	// the result may depend on the current values of other fields, and
	// therefore may change over time.
	Validate() string
	// Validate returns an empty string if the current value is valid in a
	// PackItForms-encoded message, or a problem description string if it is
	// not.  ValidatePIFO may be less restrictive than Validate for some
	// fields.  Note that the result may depend on the current values of
	// other fields, and therefore may change over time.
	ValidatePIFO() string
}

// BaseField returns a new base field with the specified characteristics.  None
// of the parameters have defaults, and all of them except key are required.
func BaseField(container FieldContainer, key FieldKey, tag, label string, width, height int, help string) FormField {
	return &baseField{
		container: container,
		key:       key,
		tag:       tag,
		label:     label,
		help:      help,
		width:     width,
		height:    height,
	}
}

type baseField struct {
	container FieldContainer
	tag       string
	key       FieldKey
	label     string
	help      string
	width     int
	height    int
	value     Value
}

func (bf *baseField) Container() FieldContainer     { return bf.container }
func (bf *baseField) Tag() string                   { return bf.tag }
func (bf *baseField) Key() FieldKey                 { return bf.key }
func (bf *baseField) Label() string                 { return bf.label }
func (bf *baseField) Hint() string                  { return "" }
func (bf *baseField) Help() string                  { return bf.help }
func (bf *baseField) Default() Value                { return "" }
func (bf *baseField) Editable() bool                { return true }
func (bf *baseField) Size() (width int, height int) { return bf.width, bf.height }
func (bf *baseField) Choices() []Value              { return nil }
func (bf *baseField) Value() Value                  { return bf.value }
func (bf *baseField) SetValue(v Value)              { bf.value = v }
func (bf *baseField) ToHuman(v Value) string        { return string(v) }
func (bf *baseField) FromHuman(v string) Value      { return Value(v) }
func (bf *baseField) Calculate()                    {}
func (bf *baseField) Validate() string              { return "" }
func (bf *baseField) ValidatePIFO() string          { return "" }

// RequiredField wraps a form field and makes it required.
func RequiredField(f FormField) FormField { return requiredField{f} }

type requiredField struct {
	FormField
}

func (rf requiredField) Help() string {
	return rf.FormField.Help() + "  This is a required field."
}

func (rf requiredField) Validate() string {
	if rf.Value() == "" {
		return fmt.Sprintf("A value for the %q field is required.", rf.Label())
	}
	return rf.FormField.Validate()
}

func (rf requiredField) ValidatePIFO() string {
	if rf.Value() == "" {
		return fmt.Sprintf("A value for the %q field is required.", rf.Label())
	}
	return rf.FormField.ValidatePIFO()
}

// MessageNumberField wraps a form field and enforces that it contains a valid
// XSC message number.
func MessageNumberField(packet bool, f FormField) FormField {
	return messageNumberField{f, packet}
}

type messageNumberField struct {
	FormField
	packet bool
}

func (mf messageNumberField) Hint() string { return "XXX-###S" }

func (mf messageNumberField) Help() string {
	return mf.FormField.Help() + `  Message numbers have the format XXX-###S, where XXX is a message number prefix (usually the last three characters of your call sign), ### is a sequence number (three or more digits), and S is a suffix letter "P" or "M".`
}

func (mf messageNumberField) Size() (width int, height int) { return 10, 1 }

var msgnoFromHumanRE = regexp.MustCompile(`^(?i)(...-)([0-9]+)[A-Z]?$`)

func (mf messageNumberField) FromHuman(h string) Value {
	// If it looks like a message number, convert it to upper case and make
	// sure the sequence number has at least three digits, and no more
	// digits than that unless needed.
	if match := msgnoFromHumanRE.FindStringSubmatch(h); match != nil {
		num, _ := strconv.Atoi(match[2])
		return Value(strings.ToUpper(fmt.Sprintf("%s%03d%s", match[1], num, match[3])))
	}
	return Value(h)
}

var msgnoValidateRE = regexp.MustCompile(`^([0-9][A-Z]{2}|[A-Z][A-Z0-9]{2})-([1-9][0-9]{3,}|[0-9]{3})([MPR]?)$`)

func (mf messageNumberField) Validate() string {
	if match := msgnoValidateRE.FindStringSubmatch(string(mf.Value())); match == nil {
		if match[3] != "" || !mf.packet {
			return mf.FormField.Validate()
		}
	}
	return fmt.Sprintf("The %q field does not contain a valid message number.", mf.Label())
}

func (mf messageNumberField) ValidatePIFO() string {
	return mf.FormField.ValidatePIFO() // PIFO never validates message number format
}

// ReadOnlyField wraps a form field and marks it as being read-only.  It can be
// updated by software but not by a user directly.
func ReadOnlyField(f FormField) FormField { return readOnlyField{f} }

type readOnlyField struct{ FormField }

func (rf readOnlyField) Editable() bool { return false }

// CallSignField wraps a form field and requires it to contain a call sign.  If
// tactical is true, it can be either a tactical or an FCC call sign; otherwise,
// it must be an FCC call sign.
func CallSignField(tactical bool, f FormField) FormField {
	return callSignField{f, tactical}
}

type callSignField struct {
	FormField
	tactical bool
}

func (cf callSignField) Help() string {
	if cf.tactical {
		return cf.FormField.Help() + " The field must contain a valid FCC or tactical call sign."
	}
	return cf.FormField.Help() + " The field must contain a valid FCC call sign."
}

func (cf callSignField) Size() (width int, height int) { return 6, 1 }

func (cf callSignField) FromHuman(h string) Value { return Value(strings.ToUpper(h)) }

var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][A-Z][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3})$`)
var tacCallSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{5}$`)

func (cf callSignField) Validate() string {
	v := string(cf.Value())
	if fccCallSignRE.MatchString(v) || (cf.tactical && tacCallSignRE.MatchString(v)) {
		return cf.FormField.Validate()
	}
	if cf.tactical {
		return fmt.Sprintf(`The %q field does not contain a valid FCC or tactical call sign.`, cf.Label())
	}
	return fmt.Sprintf(`The %q field does not contain a valid FCC call sign.`, cf.Label())
}

func (cf callSignField) ValidatePIFO() string {
	return cf.FormField.ValidatePIFO() // PIFO never validates call sign format
}

// DateField wraps a form field and requires it to contain a date.
func DateField(f FormField) FormField { return dateField{f} }

type dateField struct{ FormField }

func (df dateField) Hint() string { return "MM/DD/YYYY" }

func (df dateField) Help() string {
	return df.FormField.Help() + " The date must be in MM/DD/YYYY format."
}

func (df dateField) Size() (width int, height int) { return 10, 1 }

func (df dateField) Choices() []Value {
	return []Value{Value(time.Now().Format("01/02/2006"))}
}

var dateFromHumanRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[.-/](0?[1-9]|[12][0-9]|3[01])[.-/]((?:20)?[0-9][0-9])$`)

func (df dateField) FromHuman(h string) Value {
	// Add leading zeroes where needed and change separators to slashes.
	if match := dateFromHumanRE.FindStringSubmatch(h); match != nil {
		m, _ := strconv.Atoi(match[1])
		d, _ := strconv.Atoi(match[2])
		y, _ := strconv.Atoi(match[3])
		if y < 100 {
			y += 2000
		}
		return Value(fmt.Sprintf("%02d/%02d/%04d", m, d, y))
	}
	return Value(h)
}

var dateValidateRE = regexp.MustCompile(`^(?:0[1-9]|1[012])/(?:0[1-9]|[12][0-9]|3[01])/[1-2][0-9][0-9][0-9]$`) // from PIFO

func (df dateField) Validate() string {
	if v := string(df.Value()); v == "" || dateValidateRE.MatchString(v) {
		return df.FormField.Validate()
	}
	return fmt.Sprintf(`The %q field does not contain a valid MM/DD/YYYY date.`, df.Label())
}

func (df dateField) ValidatePIFO() string {
	if v := string(df.Value()); v == "" || dateValidateRE.MatchString(v) {
		return df.FormField.ValidatePIFO()
	}
	return fmt.Sprintf(`The %q field does not contain a valid MM/DD/YYYY date.`, df.Label())
}

// TimeField wraps a form field and requires it to contain a time.
func TimeField(f FormField) FormField { return timeField{f} }

type timeField struct{ FormField }

func (tf timeField) Hint() string { return "HH:MM" }

func (tf timeField) Help() string {
	return tf.FormField.Help() + " The time must be in HH:MM format (24-hour clock)."
}

func (tf timeField) Size() (width int, height int) { return 5, 1 }

var timeFromHumanRE = regexp.MustCompile(`^([0-9]{1,2}):?([0-9]{2})$`)

func (tf timeField) FromHuman(h string) Value {
	// Add leading zero and colon if needed.
	if match := timeFromHumanRE.FindStringSubmatch(h); match != nil {
		h, _ := strconv.Atoi(match[1])
		m, _ := strconv.Atoi(match[2])
		return Value(fmt.Sprintf("%02d:%02d", h, m))
	}
	return Value(h)
}

var timeValidateRE = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00$`) // from PIFO

func (tf timeField) Validate() string {
	if v := string(tf.Value()); v == "" || timeValidateRE.MatchString(v) {
		return tf.FormField.Validate()
	}
	return fmt.Sprintf(`The %q field does not contain a valid HH:MM time.`, tf.Label())
}

func (tf timeField) ValidatePIFO() string {
	if v := string(tf.Value()); v == "" || timeValidateRE.MatchString(v) {
		return tf.FormField.ValidatePIFO()
	}
	return fmt.Sprintf(`The %q field does not contain a valid HH:MM time.`, tf.Label())
}

// ChoicesField wraps a form field and requires its value to be one of a
// discrete set of possibilities.  The supplied choices is a list of Values and
// strings, starting with a Value, which are interpreted as follows:
//   - Any Value in the list is allowed as a Value for the field.
//   - When called for a particular Value in the list, ToHuman will return the
//     string immediately after it in the choices list.  If there isn't one
//     (i.e., the Value is followed by another Value, or is at the end of the
//     list), ToHuman will return the Value itself converted to a string.
//   - FromHuman will look for a case-insensitive, unambiguous prefix match to
//     any Value or string in the choices list.  If it matches a Value, it will
//     return that Value.  If it matches a string, it will return the
//     immediately preceding Value.
//
// The choices list can alternatively consist of a single element of type
// func() []any, which returns the list of choices.  It will be called any time
// the list is needed.
func ChoicesField(f FormField, choices ...any) FormField {
	if fn, ok := choices[0].(func() []any); ok {
		return choicesField{f, nil, fn}
	}
	return choicesField{f, choices, nil}
}

type choicesField struct {
	FormField
	choices   []any
	choicesFn func() []any
}

func (cf choicesField) Help() string {
	var list []string
	if cf.choicesFn != nil {
		cf.choices = cf.choicesFn()
	}
	for i, choice := range cf.choices {
		if v, ok := choice.(Value); ok {
			if i < len(cf.choices)-1 {
				if h, ok := cf.choices[i+1].(string); ok {
					list = append(list, h)
				} else {
					list = append(list, string(v))
				}
			} else {
				list = append(list, string(v))
			}
		}
	}
	switch len(list) {
	case 0, 1:
		panic("ChoicesField must have at least two choices")
	case 2:
		return fmt.Sprintf("%s Allowed values are %s and %s.", cf.FormField.Help(), list[0], list[1])
	default:
		return fmt.Sprintf("%s Allowed values are %s, and %s.", cf.FormField.Help(), strings.Join(list[:len(list)-1], ", "), list[len(list)-1])
	}
}

func (cf choicesField) Size() (width int, height int) {
	if cf.choicesFn != nil {
		cf.choices = cf.choicesFn()
	}
	for _, choice := range cf.choices {
		var l int
		switch choice := choice.(type) {
		case Value:
			l = len(choice)
		case string:
			l = len(choice)
		default:
			panic("ChoicesField choices list must have only Values and strings")
		}
		if l > width {
			width = l
		}
	}
	return width, 1
}

func (cf choicesField) Choices() (list []Value) {
	if cf.choicesFn != nil {
		cf.choices = cf.choicesFn()
	}
	for _, choice := range cf.choices {
		if choice, ok := choice.(Value); ok {
			list = append(list, choice)
		}
	}
	return list
}

func (cf choicesField) ToHuman(v Value) string {
	if cf.choicesFn != nil {
		cf.choices = cf.choicesFn()
	}
	for i, choice := range cf.choices {
		if choice == v {
			if i < len(cf.choices)-1 {
				if h, ok := cf.choices[i+1].(string); ok {
					return h
				}
			}
			return string(v)
		}
	}
	return string(v)
}

func (cf choicesField) FromHuman(h string) (v Value) {
	var (
		lastValue Value
		found     bool
	)
	if cf.choicesFn != nil {
		cf.choices = cf.choicesFn()
	}
	for _, choice := range cf.choices {
		var c string
		switch choice := choice.(type) {
		case Value:
			lastValue = choice
			c = string(choice)
		case string:
			c = choice
		default:
			panic("ChoicesField choices list must have only Values and strings")
		}
		if len(c) >= len(h) && strings.EqualFold(c[:len(h)], h) {
			if found {
				return Value(h)
			}
			v, found = lastValue, true
		}
	}
	if found {
		return v
	}
	return Value(h)
}

func (cf choicesField) Validate() string {
	if prob := cf.validate(); prob != "" {
		return prob
	}
	return cf.FormField.Validate()
}
func (cf choicesField) ValidatePIFO() string {
	if prob := cf.validate(); prob != "" {
		return prob
	}
	return cf.FormField.ValidatePIFO()
}
func (cf choicesField) validate() string {
	v := cf.Value()
	if cf.choicesFn != nil {
		cf.choices = cf.choicesFn()
	}
	for _, choice := range cf.choices {
		if choice == v {
			return ""
		}
	}
	return fmt.Sprintf("The value %q is not valid for the %q field.", v, cf.Label())
}

// NoValidatePIFO applies a wrapper to a form field, but only for methods other
// than ValidatePIFO.
func NoValidatePIFO(wrapper func(FormField) FormField, field FormField) FormField {
	return noValidatePIFO{wrapper(field), field}
}

type noValidatePIFO struct {
	FormField
	unwrapped FormField
}

func (nf noValidatePIFO) ValidatePIFO() string {
	return nf.unwrapped.ValidatePIFO()
}

// A FieldKey is an identifier of a well-known field that is constant across all
// messages containing that field, even if it has different tags in different
// message types.
type FieldKey string

// Values for FieldKey.  Generally, these are all of the fields that
// non-message-type-specific code needs to interact with.
const (
	// FOriginMsgNo is the origin message number field.  It is set by code
	// generating a new outgoing message, and read by code generating
	// subject lines.
	FOriginMsgNo FieldKey = "ORIGIN_MESSAGE_NUMBER"
	// FDestinationMsgNo is the destination (receiver) message number field.
	// Code that is receiving messages will set this to a local message
	// number.
	FDestinationMsgNo FieldKey = "DESTINATION_MESSAGE_NUMBER"
	// FHandling is the handling order for the message.  It gets used in
	// generating subject lines, and is read by wppsvr to verify correct
	// handling.
	FHandling FieldKey = "HANDLING"
	// FToICSPosition is the To ICS Position field.  It gets read by wppsvr
	// to verify correct routing.
	FToICSPosition FieldKey = "TO_ICS_POSITION"
	// FToLocation is the To Location field.  It gets read by wppsvr to
	// verify correct handling.
	FToLocation FieldKey = "TO_LOCATION"
	// FSubject is the Subject field.  It is the field whose contents are
	// returned by OMessage.Subject() if the message type does not have a
	// SubjectFunc.  If the message type does have a SubjectFunc, that
	// function often uses the contents of this field as part of the subject
	// line.
	FSubject FieldKey = "SUBJECT"
	// FReference is the Reference field.  It is the field that contains the
	// origin message ID of the message to which the instant message is a
	// reply.
	FReference FieldKey = "REFERENCE"
	// FBody is the field whose contents are returned by OMessage.Body() if
	// the message type does not have a BodyFunc.  It is also the field into
	// which default message body text is placed.
	FBody FieldKey = "BODY"
	// FComplete is the field that, when set to "Complete", triggers
	// conditional requirement of other fields marked with the
	// RequiredForComplete flag.
	FComplete FieldKey = "COMPLETE"
	// FOpCall is the operator call sign field.  It gets set by code that
	// creates a new outgoing message, or by code receiving a message.
	FOpCall FieldKey = "OPERATOR_CALL_SIGN"
	// FOpName is the operator name field.  It gets set by code that
	// creates a new outgoing message, or by code receiving a message.
	FOpName FieldKey = "OPERATOR_NAME"
	// FTacCall is the tactical call sign field.  It gets set by code that
	// creates a new outgoing message.
	FTacCall FieldKey = "TACTICAL_CALL_SIGN"
	// FTacName is the tactical name field.  It gets set by code that
	// creates a new outgoing message.
	FTacName FieldKey = "TACTICAL_NAME"
	// FOpDate is the transmission date (for outgoing messages) or the
	// reception date (for incoming messages).  It gets set when a message
	// is sent or received.
	FOpDate FieldKey = "OPERATOR_DATE"
	// FOpTime is the transmission time (for outgoing messages) or the
	// reception date (for incoming messages).  It gets set when a message
	// is sent or received.
	FOpTime FieldKey = "OPERATOR_TIME"
)
