package field

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/rothskeller/packet/message/msgifc"
	"k8s.io/apimachinery/pkg/util/sets"
)

//------------------------------------------------------------------------------

// NewAddressList creates a factory for a new address list field with the
// specified label.
func NewAddressList(label string) (ff *FieldFactory) {
	ff = NewField("", label)
	ff.tf = addressList{&ff.f}
	ff.f.editHeight = 1
	return ff
}

type addressList struct{ *field }

func (f addressList) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if !pifo {
		if _, err := address.ParseList(f.Value(m)); err != nil {
			return errors.NewF("The %q field does not contain a valid address list.", f.label)
		}
	}
	return f.validateCustom(m, fi, pifo)
}

func (f addressList) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	if addrs, err := address.ParseList(s); err == nil {
		strs := make([]string, len(addrs))
		for i, a := range addrs {
			strs[i] = a.String()
		}
		s = strings.Join(strs, ", ")
	}
	return s
}

//------------------------------------------------------------------------------

var PIFOCardinalNumberRE = regexp.MustCompile(`^[0-9]+$`) // changed * to +

// NewCardinalNumber creates a factory for a new cardinal number field with the
// specified tag and label.
func NewCardinalNumber(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.tf = cardinalNumber{&ff.f}
	ff.f.editHeight = 1
	return ff
}

type cardinalNumber struct{ *field }

func (f cardinalNumber) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFOCardinalNumberRE.MatchString(value) {
		return errors.NewF("The %q field does not contain a valid number.", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f cardinalNumber) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if n, err := strconv.Atoi(s); err == nil {
		s = strconv.Itoa(n)
	}
	return s
}

//------------------------------------------------------------------------------

// NewCheckbox creates a factory for a new checkbox field with the specified tag
// and label.
func NewCheckbox(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.AllowedValues("checked")
	return ff
}

// NewCheckboxGroup creates a factory for a group of related checkbox fields.
// The label is the label for the group.  The remaining strings are pairs: tag
// label for each checkbox in the group.  Note: marking the group Required
// means that at least one checkbox in the group must be checked.
func NewCheckboxGroup(label string, tagLabelPairs ...string) (ff *FieldFactory) {
	ff = NewField("", label)
	ff.tf = checkboxGroup{&ff.f}
	ff.f.editableFunc = NotEditable
	for i := 0; i < len(tagLabelPairs)-1; i += 2 {
		cb := NewCheckbox(tagLabelPairs[i], tagLabelPairs[i+1])
		cb.f.visibleFunc = Invisible
		cb.f.parent = ff.tf
		ff.f.children = append(ff.f.children, cb.MakeField())
	}
	return ff
}

type checkboxGroup struct{ *field }

func (f checkboxGroup) Value(m msgifc.Message) string {
	var checked []string

	// The internal value of a checkbox group is a comma-separated list of
	// the tags of the checked boxes.  This behaves correctly in Required
	// checks.
	for _, cb := range f.children {
		if cb.Value(m) != "" {
			checked = append(checked, cb.Tag())
		}
	}
	return strings.Join(checked, ",")
}

func (f checkboxGroup) ToHuman(s string) string {
	if f.toHumanFunc != nil {
		return f.toHumanFunc(s)
	}
	// The human form of a checkbox group value is a comma-separated list of
	// the labels of the checked boxes.
	if s == "" {
		return s
	}
	var labels []string
	var tags = sets.New(strings.Split(s, ",")...)
	for _, cb := range f.children {
		if tags.Has(cb.Tag()) {
			labels = append(labels, cb.Label())
			tags.Delete(cb.Tag())
		}
	}
	labels = append(labels, tags.UnsortedList()...)
	return strings.Join(labels, ", ")
}

//------------------------------------------------------------------------------

var (
	dateLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)
	PIFODateRE  = regexp.MustCompile(`^(?:0[1-9]|1[012])/(?:0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]$`)
)

// NewDate creates a factory for a new date field with the specified tag and
// label.
func NewDate(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 10
	ff.f.editHeight = 1
	ff.f.editHint = "MM/DD/YYYY"
	ff.tf = date{&ff.f}
	return ff
}

type date struct{ *field }

func (f date) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFODateRE.MatchString(value) {
		return errors.NewF("The %q field does not contain a valid date (MM/DD/YYYY).", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f date) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if match := dateLooseRE.FindStringSubmatch(s); match != nil {
		// Add leading zeroes and set delimiter to slash.
		s = fmt.Sprintf("%02s/%02s/20%s", match[1], match[2], match[3])
		// Correct values that are out of range, e.g. 06/31 => 07/01.
		if t, err := time.ParseInLocation("01/02/2006", s, time.Local); err == nil {
			s = t.Format("01/02/2006")
		}
	}
	return s
}

//------------------------------------------------------------------------------

// NewDateTime creates a factory for a new date/time field with the specified
// label.  It requires factories for a date and a time field, which will become
// children of the date/time field.
func NewDateTime(label string, dateField, timeField *FieldFactory) (ff *FieldFactory) {
	ff = NewField("", label)
	ff.tf = datetime{&ff.f}
	// Connect the fields.
	df := dateField.MakeField().(date)
	tf := timeField.MakeField().(timef)
	df.parent = ff.tf
	tf.parent = ff.tf
	ff.f.children = []Field{df, tf}
	// Mark the date and time not visible, and not editable except by
	// explicit field name.
	df.visibleFunc = Invisible
	tf.visibleFunc = Invisible
	df.editableFunc = OnlyExplicitlyEditable
	tf.editableFunc = OnlyExplicitlyEditable
	// Set up the date/time field.
	ff.f.editWidth = 16
	ff.f.editHeight = 1
	ff.f.editHint = "MM/DD/YYYY HH:MM"
	ff.f.requiredFunc = ff.tf.(datetime).required
	if df.requiredDesc == tf.requiredDesc {
		ff.f.requiredDesc = df.requiredDesc
	} else {
		ff.f.requiredDesc = smartJoin(df.requiredDesc, tf.requiredDesc)
	}
	ff.f.disallowedFunc = ff.tf.(datetime).disallowed
	if df.disallowedDesc == tf.disallowedDesc {
		ff.f.disallowedDesc = df.disallowedDesc
	} else {
		ff.f.disallowedDesc = smartJoin(df.disallowedDesc, tf.disallowedDesc)
	}
	return ff
}

type datetime struct{ *field }

func (f datetime) Value(m msgifc.Message) string {
	d := f.children[0].Value(m)
	t := f.children[1].Value(m)
	return smartJoin(d, t)
}

func (f datetime) FromHuman(s string) string {
	d, t, _ := strings.Cut(s, " ")
	d = f.children[0].FromHuman(d)
	t = f.children[1].FromHuman(t)
	return smartJoin(d, t)
}

func (f datetime) SetValue(m msgifc.Message, s string) {
	d, t, _ := strings.Cut(s, " ")
	f.children[0].SetValue(m, d)
	f.children[1].SetValue(m, t)
}

func (f datetime) required(m msgifc.Message) bool {
	if fn := f.children[0].(date).requiredFunc; fn != nil && fn(m) {
		return true
	}
	if fn := f.children[1].(timef).requiredFunc; fn != nil && fn(m) {
		return true
	}
	return false
}

func (f datetime) disallowed(m msgifc.Message) bool {
	if fn := f.children[0].(date).disallowedFunc; fn != nil && !fn(m) {
		return false
	}
	if fn := f.children[1].(timef).disallowedFunc; fn != nil && !fn(m) {
		return false
	}
	return true
}

func smartJoin(a, b string) string {
	if a != "" && b != "" {
		return a + " " + b
	}
	return a + b
}

//------------------------------------------------------------------------------

var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)

// NewFCCCallSign creates a factory for a new FCC call sign field with the
// specified tag and label.
func NewFCCCallSign(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 6
	ff.f.editHeight = 1
	ff.tf = fccCallSign{&ff.f}
	return ff
}

type fccCallSign struct{ *field }

func (f fccCallSign) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !pifo && !fccCallSignRE.MatchString(f.Value(m)) {
		return errors.NewF("The %q field does not contain a valid FCC call sign.", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f fccCallSign) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	return strings.ToUpper(strings.TrimSpace(s))
}

//------------------------------------------------------------------------------

var PIFOFrequencyRE = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)

// NewFrequency creates a factory for a new frequency field with the specified
// tag and label.
func NewFrequency(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editHint = "MHz"
	ff.f.editHeight = 1
	ff.tf = frequency{&ff.f}
	return ff
}

type frequency struct{ *field }

func (f frequency) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFOFrequencyRE.MatchString(f.Value(m)) {
		return errors.NewF("The %q field does not contain a valid frequency in MHz.", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f frequency) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		s = strconv.FormatFloat(n, 'f', -1, 64)
	}
	return s
}

//------------------------------------------------------------------------------

var PIFOFrequencyOffsetRE = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+])$`)

// NewFrequencyOffset creates a factory for a new frequency offset field with
// the specified tag and label.
func NewFrequencyOffset(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editHint = "MHz or +/-"
	ff.f.editHeight = 1
	ff.tf = frequencyOffset{&ff.f}
	return ff
}

type frequencyOffset struct{ *field }

func (f frequencyOffset) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFOFrequencyOffsetRE.MatchString(f.Value(m)) {
		return errors.NewF(`The %q field does not contain a valid frequency offset (a real number in MHz, a "+", or a "-").`, f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f frequencyOffset) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		s = strconv.FormatFloat(n, 'f', -1, 64)
	}
	return s
}

//------------------------------------------------------------------------------

// NewMessageID creates a factory for a new message ID field with the specified
// tag and label.  If packet is true, the message ID is validated to have a
// letter suffix.
func NewMessageID(tag, label string, packet bool) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 9
	ff.f.editHeight = 1
	ff.f.editHint = "XXX-###P"
	ff.tf = messageID{&ff.f, packet}
	return ff
}

type messageID struct {
	*field
	packet bool
}

func (f messageID) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !pifo {
		if _, _, _, err := messageid.Decode(value, false, f.packet); err != nil {
			if f.packet {
				return errors.NewF("The %q field does not contain a valid packet message ID (XXX-###X).", f.label)
			} else {
				return errors.NewF("The %q field does not contain a valid message ID (XXX-### or XXX-###X).", f.label)
			}
		}
	}
	return f.validateCustom(m, fi, pifo)
}

func (f messageID) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if prefix, seq, suffix, err := messageid.Decode(s, true, f.packet); err == nil {
		s, _ = messageid.Encode(prefix, seq, suffix)
	}
	return s
}

//------------------------------------------------------------------------------

var (
	PIFOPhoneNumberRE = regexp.MustCompile(`^[a-zA-Z ]*(?:[+][0-9]+ )?[0-9][0-9 -]*(?:[xX][0-9]+)?$`)
	phoneExtensionRE  = regexp.MustCompile(`[xX][0-9]+$`)
)

// NewPhoneNumber creates a factory for a new phone number field with the
// specified tag and label.
func NewPhoneNumber(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editHint = "###-###-####"
	ff.f.editHeight = 1
	ff.tf = phoneNumber{&ff.f}
	return ff
}

type phoneNumber struct{ *field }

func (f phoneNumber) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFOPhoneNumberRE.MatchString(f.Value(m)) {
		return errors.NewF("The %q field does not contain a valid phone number (###-###-####).", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f phoneNumber) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	ext := phoneExtensionRE.FindString(s)
	s = s[:len(s)-len(ext)]
	if strings.IndexFunc(s, func(r rune) bool {
		return r < '0' && r > '9' && r != ' ' && r != '-' && r != '(' && r != ')'
	}) < 0 {
		trim := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, s)
		if len(trim) == 10 {
			s = trim[0:3] + "-" + trim[3:6] + "-" + trim[6:10]
		}
	}
	return s + ext
}

//------------------------------------------------------------------------------

var PIFORealNumberRE = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+)$`)

// NewRealNumber creates a factory for a new real number field with the
// specified tag and label.
func NewRealNumber(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.tf = realNumber{&ff.f}
	ff.f.editHeight = 1
	return ff
}

type realNumber struct{ *field }

func (f realNumber) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFORealNumberRE.MatchString(f.Value(m)) {
		return errors.NewF("The %q field does not contain a valid number.", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f realNumber) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		s = strconv.FormatFloat(n, 'f', -1, 64)
	}
	return s
}

//------------------------------------------------------------------------------

var tacticalCallSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{4,5}$`)

// NewTacticalCallSign creates a factory for a new tactical call sign field with
// the specified tag and label.
func NewTacticalCallSign(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 6
	ff.f.editHeight = 1
	ff.tf = tacticalCallSign{&ff.f}
	return ff
}

type tacticalCallSign struct{ *field }

func (f tacticalCallSign) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !pifo && !tacticalCallSignRE.MatchString(f.Value(m)) {
		return errors.NewF("The %q field does not contain a valid tactical call sign.", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f tacticalCallSign) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	return strings.ToUpper(strings.TrimSpace(s))
}

//------------------------------------------------------------------------------

var (
	timeLooseRE = regexp.MustCompile(`^([1-9]:|[01][0-9]:?|2[0-4]:?)([0-5][0-9])$`)
	PIFOTimeRE  = regexp.MustCompile(`^(?:([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00)$`)
)

// NewTime creates a factory for a new time field with the specified tag and
// label.
func NewTime(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 5
	ff.f.editHeight = 1
	ff.f.editHint = "HH:MM"
	ff.tf = timef{&ff.f}
	return ff
}

type timef struct{ *field }

func (f timef) Validate(m msgifc.Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if value := f.Value(m); value != "" && !PIFOTimeRE.MatchString(f.Value(m)) {
		return errors.NewF("The %q field does not contain a valid time (HH:MM, 24-hour clock).", f.label)
	}
	return f.validateCustom(m, fi, pifo)
}

func (f timef) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if match := timeLooseRE.FindStringSubmatch(s); match != nil {
		if !strings.HasSuffix(match[1], ":") {
			match[1] += ":"
		}
		s = fmt.Sprintf("%03s%s", match[1], match[2])
	}
	return s
}
