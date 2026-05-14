package form

import (
	"fmt"
	"iter"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/rothskeller/packet/message/msgifc"
	"github.com/rothskeller/packet/message/payload"
)

// ff2mf is an adapter that implements the message.Field interface for a
// formdef.FieldDef structure.
type ff2mf struct {
	fd *formdef.FieldDef
}

var _ field.Field = ff2mf{}

// Tag returns the PackItForms tag for the field, if any.
func (f ff2mf) Tag() string { return f.fd.Tag }

// Common is the tag string that identifies the field as one of the
// well-known common fields.  It allows software to address those
// fields without dependency on how they're stored in a particular
// message type.  It should be one of the field.CommonTag constants.
func (f ff2mf) Common() string { return f.fd.Common }

// Label returns the label for the field, if any.
func (f ff2mf) Label() string { return f.fd.Label }

// Parent returns the parent field that contains this field, if any.
// Generally only used for checkbox groups.
func (f ff2mf) Parent() field.Field {
	if f.fd.Parent != nil {
		return ff2mf{f.fd.Parent}
	}
	return nil
}

// Children returns the list of child fields contained by this field,
// if any.  Generally only used for checkbox groups.  A Field may not
// have both a Parent and Children.
func (f ff2mf) Children() (c []field.Field) {
	if len(f.fd.Children) == 0 {
		return nil
	}
	c = make([]field.Field, len(f.fd.Children))
	for i := range f.fd.Children {
		c[i] = ff2mf{f.fd.Children[i]}
	}
	return c
}

// Value returns the value of the field, in internal form.
func (f ff2mf) Value(m msgifc.Message) (val string) {
	switch f.fd.Type {
	case "static":
		val = f.fd.Value
		for _, vc := range f.fd.ValueCond {
			if cond(m, vc.OtherField, vc.OtherValue) {
				val = vc.ThisValue
				break
			}
		}
		return val
	case "checkboxGroup":
		var set []string
		for _, c := range f.fd.Children {
			if c.Tag != "" && fv(m, c) != "" {
				set = append(set, c.ChildLabel)
			}
		}
		return strings.Join(set, ", ")
	case "dateTime", "join":
		var set []string
		for _, c := range f.fd.Children {
			if c.Tag != "" { // should always be true, just being safe
				set = append(set, fv(m, c))
			}
		}
		val = f.fd.Value
		for _, vc := range f.fd.ValueCond {
			if cond(m, vc.OtherField, vc.OtherValue) {
				val = vc.ThisValue
				break
			}
		}
		if val != "" {
			return joinWithPattern(set, val)
		} else {
			set = slices.DeleteFunc(set, func(s string) bool { return s == "" })
			return strings.Join(set, " ")
		}
	default:
		if f.fd.Tag != "" {
			return fv(m, f.fd)
		}
		return ""
	}
}

var joinRefRE = regexp.MustCompile(`\{\d+\}`)

// joinWithPattern joins the values of the children of a "join" field according
// to the pattern in the field's definition, which has '{1}' references
// interspersed with separator text.
func joinWithPattern(ss []string, patt string) string {
	var refs, seps []string

	// Build a list of references and separator strings.  len(seps) will
	// always be one greater than len(refs).  seps[i] is the string that
	// precedes refs[i], and the last sep is the trailing string after the
	// last ref.
	match := joinRefRE.FindAllStringIndex(patt, -1)
	for i, pair := range match {
		if i != 0 {
			seps = append(seps, patt[match[i-1][1]:pair[0]])
		} else {
			seps = append(seps, patt[0:pair[0]])
		}
		idx, _ := strconv.Atoi(patt[pair[0]+1 : pair[1]-1])
		if idx > 0 && idx <= len(ss) {
			refs = append(refs, ss[idx-1])
		} else {
			// This is an error in the pattern, but we don't have
			// any reasonable path to report it, so we'll just
			// ignore it.
			refs = append(refs, "")
		}
	}
	if len(match) != 0 {
		seps = append(seps, patt[match[len(match)-1][1]:])
	} else {
		// A pattern without any refs in it is pretty useless...
		seps = append(seps, patt)
	}
	// If the string starts with a reference and that child's value is
	// empty, remove the reference and the separator string following it.
	// Repeat until false.
	for len(refs) != 0 && refs[0] == "" && seps[0] == "" {
		refs, seps = refs[1:], seps[1:]
		seps[0] = ""
	}
	// If the string ends with a reference and that child's value is empty,
	// remove the reference and the separator string preceding it.  Repeat
	// until false.
	for len(refs) != 0 && refs[len(refs)-1] == "" && seps[len(seps)-1] == "" {
		refs, seps = refs[:len(refs)-1], seps[:len(seps)-1]
		seps[len(seps)-1] = ""
	}
	// If any interior references are empty, remove their preceding
	// separator strings.
	for i := range refs {
		if refs[i] == "" {
			seps[i] = ""
		}
	}
	// Concatenate and return whatever's left.
	var sb strings.Builder
	for i := range refs {
		sb.WriteString(seps[i])
		sb.WriteString(refs[i])
	}
	sb.WriteString(seps[len(seps)-1])
	return sb.String()
}

// Default returns the default value of the field, in internal form.
func (f ff2mf) Default() string {
	if f.fd.Type == "join" {
		return ""
	}
	return f.fd.Value
}

// ToHuman converts the internal form of a value for the field into the
// human form appropriate for display and editing (often a no-op).
func (f ff2mf) ToHuman(msg message.Message, raw string) string {
	for _, c := range f.fd.Choices {
		if !cond(msg, c.CondField, c.CondValue) {
			continue
		}
		if c.Raw == raw {
			return c.Human
		}
	}
	return raw
}

// FromHuman converts the supplied value from human form to internal
// form, if possible; otherwise it makes no changes.  Implementations
// must not change the value if it is already in internal form.
func (f ff2mf) FromHuman(msg message.Message, human string) string {
	if human = strings.TrimSpace(human); human == "" {
		return human
	}
	switch f.fd.Type {
	case "cardinalNumber":
		if n, err := strconv.Atoi(human); err == nil {
			return strconv.Itoa(n)
		}
	case "date":
		return canonicalDate(human)
	case "dateTime":
		d, t, _ := strings.Cut(human, " ")
		d = canonicalDate(d)
		t = canonicalTime(t)
		if d != "" && t != "" {
			return d + " " + t
		} else if d != "" {
			return d
		} else {
			return t
		}
	case "fccCallSign", "tacticalCallSign":
		return strings.ToUpper(human)
	case "frequency", "frequencyOffset", "realNumber":
		if n, err := strconv.ParseFloat(human, 64); err == nil {
			return strconv.FormatFloat(n, 'f', -1, 64)
		}
	case "messageID":
		human, _ = messageid.Cleanup(human, false)
		return human
	case "phoneNumber":
		ext := phoneExtensionRE.FindString(human)
		human = strings.TrimSuffix(human, ext)
		if strings.IndexFunc(human, func(r rune) bool {
			return (r < '0' || r > '9') && r != ' ' && r != '-' && r != '(' && r != ')'
		}) < 0 {
			trim := strings.Map(func(r rune) rune {
				if r >= '0' && r <= '9' {
					return r
				}
				return -1
			}, human)
			if len(trim) == 10 {
				human = trim[0:3] + "-" + trim[3:6] + "-" + trim[6:10]
			}
		}
		return human + ext
	case "restricted":
		var match string
		var ambiguous bool
		for _, c := range f.fd.Choices {
			if !cond(msg, c.CondField, c.CondValue) {
				continue
			}
			if strings.EqualFold(human, c.Human) {
				return c.Raw
			} else if human != "" && len(human) < len(c.Human) && strings.EqualFold(human, c.Human[:len(human)]) {
				ambiguous = match != ""
				match = c.Raw
			}
		}
		if match != "" && !ambiguous {
			return match
		}
	case "time":
		return canonicalTime(human)
	default:
		for _, c := range f.fd.Choices {
			if !cond(msg, c.CondField, c.CondValue) {
				continue
			}
			if strings.EqualFold(human, c.Human) {
				return c.Raw
			}
		}
	}
	return human
}

func canonicalDate(s string) string {
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

func canonicalTime(s string) string {
	if match := timeLooseRE.FindStringSubmatch(s); match != nil {
		if !strings.HasSuffix(match[1], ":") {
			match[1] += ":"
		}
		s = fmt.Sprintf("%03s%s", match[1], match[2])
	}
	return s
}

// Settable returns whether the field is settable.
func (f ff2mf) Settable() bool { return f.fd.Tag != "" || f.fd.Type == "dateTime" }

// SetValue sets the value of the field.  The supplied value must be in
// internal form.
func (f ff2mf) SetValue(msg msgifc.Message, val string) {
	if _, ok := msg.(*message.DraftMessage); ok && f.fd.Tag != "" {
		// Look for other fields that have a value but that will be
		// disallowed after this change.  If any are found, remove
		// their values.
		for f2 := range msg.Type().(EditableFormType).AllFields() {
			if f2 == f.fd || f2.Tag == "" || fv(msg, f2) == "" {
				continue
			}
			presence := f2.Presence
			found := false
			for _, pc := range f2.PresenceCond {
				if pc.OtherField == f.fd.Tag {
					found = true
					if (pc.OtherValue == "" && val != "") || (pc.OtherValue != "" && pc.OtherValue == val) {
						presence = pc.Presence
						break
					}
				}
			}
			if found && presence == formdef.Blocked {
				reason := "form.FormBody.Field." + f2.Tag
				if f2.Common != "" {
					reason = "form.FormBody.Common." + f2.Common
				}
				msg.Body().(*FormBody).SetFieldR(f2.Tag, "", reason)
			}
		}
	}
	if f.fd.Tag != "" {
		reason := "form.FormBody.Field." + f.fd.Tag
		if f.fd.Common != "" {
			reason = "form.FormBody.Common." + f.fd.Common
		}
		msg.Body().(*FormBody).SetFieldR(f.fd.Tag, val, reason)
	} else if f.fd.Type == "dateTime" {
		d, t, _ := strings.Cut(val, " ")
		for _, c := range f.fd.Children {
			switch c.Type {
			case "date":
				ff2mf{c}.SetValue(msg, d)
			case "time":
				ff2mf{c}.SetValue(msg, t)
			}
		}
	}
	switch f.fd.Common {
	case field.COriginMessageID:
		msg.Subject().SetSubjectMessageID(f.Value(msg))
	case field.CHandling:
		msg.Subject().SetSubjectHandling(f.Value(msg))
		msg.Payload().(*payload.OutpostPayload).SetUrgent(f.Value(msg) == "IMMEDIATE")
	case field.CMessageSummary:
		msg.Subject().SetSubjectSummary(f.Value(msg))
	}
}

// Visible returns whether the field should be included when the
// message is displayed.  Note that fields with an empty value are
// never displayed no matter what this method returns.
func (f ff2mf) Visible(msgifc.Message) bool {
	return f.fd.Parent == nil
}

// Editable returns whether the field should be included when the
// message is edited.  The explicit flag is true if the user explicitly
// asked to edit this field by name.
func (f ff2mf) Editable(m msgifc.Message, explicit bool) bool {
	if f.fd.EditHelp == "" {
		return false
	}
	if f.fd.Parent != nil && f.fd.Parent.Type == "dateTime" && !explicit {
		return false
	}
	if p, _ := evalPresence(m, f.fd); p == formdef.Blocked {
		return false
	}
	return true
}

// EditHelp returns the help string for editing of the field.
func (f ff2mf) EditHelp() string { return f.fd.EditHelp }

// EditHint returns the hint string, if any, for editing of the field.
// This is displayed in the editing control when the control is
// otherwise empty (the equivalent of HTML <input placeholder="...">).
func (f ff2mf) EditHint() string {
	switch f.fd.Type {
	case "date":
		return "mm/dd/yyyy"
	case "dateTime":
		return "mm/dd/yyyy hh:mm"
	case "phoneNumber":
		return "000-000-0000 x00"
	case "time":
		return "hh:mm"
	default:
		return ""
	}
}

// Multiline returns whether the field is expected to contain a
// multiline value, i.e., a value containing newlines.  All values
// *can* include newlines, but this is whether it's *expected* to.
func (f ff2mf) Multiline() bool { return f.fd.Type == "multiline" }

// EditSize returns the width and height, in characters, of the text
// entry control for the field.  Zero return values mean unlimited.
// (Usually this should be set based on the size of the corresponding
// printable area of the PDF at the minimum font size.)
func (f ff2mf) EditSize() (int, int) {
	if f.fd.Type == "multiline" {
		return f.fd.EditWidth, 0
	}
	return f.fd.EditWidth, 1
}

// Obscured returns whether the field's value should be obscured for
// display or editing.  This is used for password fields.
func (f ff2mf) Obscured() bool { return f.fd.Type == "password" }

// Choices returns a list of allowed or recommended values for the
// field.  Each element is a pair with internal and human
// representations of the value.  The list may vary depending on the
// values of other fields of the message.
func (f ff2mf) Choices(m msgifc.Message) (cs []field.ChoicePair) {
	if f.fd.Type == "checkbox" {
		return []field.ChoicePair{{PIFO: "checked", Human: "checked"}}
	}
	for _, c := range f.fd.Choices {
		if c.CondField != "" {
			val := m.Body().(*FormBody).Field(c.CondField)
			if (c.CondValue == "" && val == "") || (c.CondValue != "" && c.CondValue != val) {
				continue
			}
		}
		cs = append(cs, field.ChoicePair{PIFO: c.Raw, Human: c.Human})
	}
	return cs
}

// Restricted returns whether the value of the field is restricted to
// one of the listed Choices (true), or whether they are just
// recommendations and any other value is accepted (false).
func (f ff2mf) Restricted() bool { return f.fd.Type == "restricted" || f.fd.Type == "checkbox" }

var (
	PIFOCardinalNumberRE  = regexp.MustCompile(`^[0-9]+$`) // changed * to +
	dateLooseRE           = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)
	PIFODateRE            = regexp.MustCompile(`^(?:0[1-9]|1[012])/(?:0[1-9]|1[0-9]|2[0-9]|3[01])/[1-2][0-9][0-9][0-9]$`)
	fccCallSignRE         = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)
	PIFOFrequencyRE       = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?$`)
	PIFOFrequencyOffsetRE = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+|[-+])$`)
	PIFOPhoneNumberRE     = regexp.MustCompile(`^[a-zA-Z ]*(?:[+][0-9]+ )?[0-9][0-9 -]*(?:[xX][0-9]+)?$`)
	phoneExtensionRE      = regexp.MustCompile(`[xX][0-9]+$`)
	PIFORealNumberRE      = regexp.MustCompile(`^(?:[-+]?[0-9]*\.[0-9]+|[-+]?[0-9]+)$`)
	tacticalCallSignRE    = regexp.MustCompile(`^[A-Z][A-Z0-9]{4,5}$`)
	timeLooseRE           = regexp.MustCompile(`^([1-9]:|[01][0-9]:?|2[0-4]:?)([0-5][0-9])$`)
	PIFOTimeRE            = regexp.MustCompile(`^(?:([01][0-9]|2[0-3]):?[0-5][0-9]|2400|24:00)$`)
	zipCodeRE             = regexp.MustCompile(`^\d{5}(?:-\d{4})?$`)
)

// Validate validates the value of the field and returns any problems
// with it.  If pifo is true, it restricts itself to those checks
// performed by PackItForms.
func (f ff2mf) Validate(m msgifc.Message, mf msgifc.Field, flags msgifc.ValidateFlags) (err error) {
	presence, why := evalPresence(m, f.fd)
	val := f.Value(m)
	switch {
	case presence == formdef.Blocked && val != "" && why != "":
		return errors.NewF("Field %q cannot be set %s.", f.fd.Label, why)
	case presence == formdef.Blocked && val != "":
		return errors.NewF("Field %q cannot be set.", f.fd.Label)
	case presence == formdef.Required && val == "" && why != "":
		return errors.NewF("Field %q is required %s.", f.fd.Label, why)
	case presence == formdef.Required && val == "":
		return errors.NewF("Field %q is required.", f.fd.Label)
	}
	// Type-specific validation.
	switch f.fd.Type {
	case "cardinalNumber":
		if val != "" && !PIFOCardinalNumberRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid number.`, f.fd.Label)
		}
	case "checkbox":
		if val != "" && val != "checked" {
			return errors.NewF(`Field %q must contain either "" or "checked".`, f.fd.Label)
		}
	case "date":
		if val != "" && !PIFODateRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid date (MM/DD/YYYY).`, f.fd.Label)
		}
	case "dateTime":
		for _, c := range f.fd.Children {
			err = errors.Join(err, ff2mf{c}.Validate(m, c.Field, flags))
		}
		return err
	case "fccCallSign":
		if val != "" && !fccCallSignRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a FCC call sign.`, f.fd.Label)
		}
	case "frequency":
		if val != "" && !PIFOFrequencyRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid frequency in MHz.`, f.fd.Label)
		}
	case "frequencyOffset":
		if val != "" && !PIFOFrequencyOffsetRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid frequency offset (a real number in MHz, a "+", or a "-").`, f.fd.Label)
		}
	case "messageID":
		if val != "" {
			if _, _, _, err := messageid.Decode(val, false, flags&msgifc.VPacket != 0); err != nil {
				if flags&msgifc.VPacket != 0 {
					return errors.NewF(`Field %q does not contain a valid packet message ID (XXX-###X).`, f.fd.Label)
				} else {
					return errors.NewF(`Field %q does not contain a valid message ID (XXX-### or XXX-###X).`, f.fd.Label)
				}
			}
		}
	case "phoneNumber":
		if val != "" && !PIFOPhoneNumberRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid phone number (###-###-####).`, f.fd.Label)
		}
	case "realNumber":
		if val != "" && !PIFORealNumberRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid number.`, f.fd.Label)
		}
	case "restricted":
		if val == "" {
			break
		}
		var choices []string
		var conderr error
		var found bool
		for _, c := range f.fd.Choices {
			if val == c.Raw {
				if cond(m, c.CondField, c.CondValue) {
					found = true
					break
				}
				conderr = errors.NewF(`Field %q can only be set to %q when %s.`, f.fd.Label, c.Human, condstr(m, c.CondField, c.CondValue))
			} else if cond(m, c.CondField, c.CondValue) {
				choices = append(choices, c.Human)
			}
		}
		if !found {
			if conderr != nil {
				return conderr
			} else if len(choices) <= 4 {
				return errors.NewF(`%q is not one of the allowed values for field %q ("%s").`, val, f.fd.Label, strings.Join(choices, `", "`))
			} else {
				return errors.NewF(`%q is not one of the allowed values for field %q.`, val, f.fd.Label)
			}
		}
	case "tacticalCallSign":
		if val != "" && !tacticalCallSignRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid tactical call sign.`, f.fd.Label)
		}
	case "time":
		if val != "" && !PIFOTimeRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid time (HH:MM, 24-hour clock).`, f.fd.Label)
		}
	case "zipCode":
		if val != "" && !zipCodeRE.MatchString(val) {
			return errors.NewF(`Field %q does not contain a valid ZIP code (##### or #####-####).`, f.fd.Label)
		}
	}
	return nil
}

// Compare compares the value of the field in two messages.
func (f ff2mf) Compare(label, expected, actual string) (cmp *msgifc.ComparedField) {
	var cmpType string

	switch f.fd.Type {
	case "addressList", "checkboxGroup", "dateTime", "join", "static":
		cmpType = "none"
	case "checkbox":
		cmpType = "checkbox"
	case "password", "restricted", "zipCode":
		cmpType = "exact"
	case "fccCallSign", "messageID", "tacticalCallSign":
		cmpType = "exact-ci"
	case "frequency", "frequencyOffset":
		cmpType = "realNumber"
	case "multiline":
		cmpType = "text"
	case "cardinalNumber", "date", "phoneNumber", "realNumber", "text", "time":
		cmpType = f.fd.Type
	}
	if f.fd.CompareMethod != "" {
		cmpType = f.fd.CompareMethod
	}
	switch cmpType {
	case "none":
		return nil
	case "cardinalNumber":
		return field.CompareCardinal(label, expected, actual)
	case "checkbox":
		return field.CompareCheckbox(label, expected, actual)
	case "date":
		return field.CompareDate(label, expected, actual)
	case "exact-ci":
		return field.CompareExactCI(label, expected, actual)
	case "phoneNumber":
		return field.ComparePhoneNumber(label, expected, actual)
	case "realNumber":
		return field.CompareReal(label, expected, actual)
	case "text":
		return field.CompareText(label, expected, actual)
	case "time":
		return field.CompareTime(label, expected, actual)
	default:
		return field.CompareExact(label, expected, actual)
	}
}

// evalPresence evaluates the presence state of a field based on the values of
// other message fields.  It also returns the reason string for the conditional
// if any.
func evalPresence(m msgifc.Message, fd *formdef.FieldDef) (formdef.Presence, string) {
	var condstrs []string
	for _, pc := range fd.PresenceCond {
		if cond(m, pc.OtherField, pc.OtherValue) {
			return pc.Presence, "when " + condstr(m, pc.OtherField, pc.OtherValue)
		}
		condstrs = append(condstrs, condstr(m, pc.OtherField, pc.OtherValue))
	}
	if len(condstrs) != 0 {
		return fd.Presence, "unless " + strings.Join(condstrs, " or ")
	}
	return fd.Presence, ""
}

// cond evaluates a conditional based on field values.  Conditionals appear in
// ValueCond, PresenceCond, and in Choices.
func cond(m msgifc.Message, tag, val string) bool {
	if tag == "" {
		return true
	}
	v := m.Body().(*FormBody).Field(tag)
	return (val == "" && v != "") || (val != "" && val == v)
}

// condstr returns a string describing a conditional, for use in messages.
func condstr(m msgifc.Message, tag, val string) string {
	var label, ftype string
	if tag == "" {
		return ""
	}
	for fd := range m.Type().(interface {
		AllFields() iter.Seq[*formdef.FieldDef]
	}).AllFields() {
		if fd.Tag == tag {
			label, ftype = fd.Label, fd.Type
			break
		}
	}
	if val != "" {
		switch ftype {
		case "cardinalNumber", "checkbox", "date", "dateTime", "fccCallSign", "frequency", "messageID", "realNumber", "tacticalCallSign", "time", "zipCode":
			return fmt.Sprintf("%q is %s", label, val)
		default:
			return fmt.Sprintf("%q is %q", label, val)
		}
	} else if ftype == "checkbox" {
		return fmt.Sprintf("%q is checked", label)
	} else {
		return fmt.Sprintf("%q is set", label)
	}
}

// fv is a shortcut for a common operation: retrieving the value of a field
// from a form message body.
func fv(m msgifc.Message, fd *formdef.FieldDef) string {
	return m.Body().(*FormBody).Field(fd.Tag)
}
