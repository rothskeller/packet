package message

import (
	"fmt"
	"strings"
)

// A Field describes a single field within a message.  Generally, a message type
// has one Field for each field in the PackItForms encoding of the message, plus
// occasionally other Fields for special purposes (e.g. aggregating the
// underlying fields for display or editing).
type Field struct {
	// Label is the name of the field, as it is displayed to the user.  It
	// should be short, definitely no more than 40 characters.
	Label string
	// Value is a pointer to where the value of the field is stored.  Not
	// all fields have a stored value, so this pointer may be nil.
	Value *string
	// Choices is a set of recommended or allowed values for the field.
	// (Whether other values are allowed is up to the validation functions.)
	Choices ChoiceMapper
	// Presence is a function that returns whether the field is allowed or
	// required.  The function may optionally return a reason, which is
	// interpolated into validation problem strings when needed.
	Presence func() (Presence, string)
	// PIFOTag is the tag for this field in a PackItForms encoding.  If this
	// is empty, the field will not be rendered in PackItForms encoding nor
	// populated from PackItForms decoding.
	PIFOTag string
	// PIFOValid checks the value of the field against the restrictions
	// enforced by the PackItForms software.  It returns a problem
	// description if the value is one that PackItForms would reject, and an
	// empty string otherwise.
	PIFOValid func(*Field) string
	// Compare compares an expected value of this field against an actual
	// value of this field, and returns a description of the comparison.  To
	// disable comparison for a field, set this to CompareNone.
	Compare func(label, exp, act string) *CompareField
	// PDFMap is the mapper that tells how to render this field into a
	// form-fillable PDF file.
	PDFMap PDFMapper
	// TableValue returns the value of this field when rendered in flat text
	// table form.  To omit a field from the table rendering, set this to
	// TableOmit.
	TableValue func(*Field) string
	// EditWidth is the width in characters of the input control for this
	// field.  It should correspond to the number of characters that will
	// fit in the PDF rendering of the field, if applicable.
	EditWidth int
	// Multiline indicates that this field can contain multiple lines, i.e.,
	// can contain newline characters.
	Multiline bool
	// EditHelp is the help text for the form field, describing its contents
	// and its validity rules.  If this is empty, the field is not editable.
	EditHelp string
	// EditHint is a short string giving a model for the field value (e.g.,
	// "MM/DD/YYYY" for a date field).  It is optional, and will only be
	// displayed if there is room for it.
	EditHint string
	// EditValue returns the editable representation of the value of the
	// field.
	EditValue func(*Field) string
	// EditApply stores the supplied edited value into the field, revising
	// it if need be to convert from human to internal (PIFO) form.
	EditApply func(*Field, string)
	// EditValid checks the value of the field and returns a problem
	// description, or an empty string if there are no problems.
	EditValid func(*Field) string
	// EditSkip returns whether the field should be skipped while editing
	// the message (e.g., no entry in this field is valid because of the
	// value of some earlier field).
	EditSkip func(*Field) bool
}

// Presence is a enumeration indicating whether a field is allowed or required.
type Presence uint8

// Values for Presence:
const (
	// PresenceNotAllowed means a value for this field is not allowed.
	// (This is generally because some parent field is not set, or some
	// conflicting field is set.)
	PresenceNotAllowed Presence = iota
	// PresenceOptional means a value for this field is allowed but not
	// required.
	PresenceOptional
	// PresenceRequired means a value for this field is required.
	PresenceRequired
)

func NotAllowed() (Presence, string) { return PresenceNotAllowed, "" }
func Optional() (Presence, string)   { return PresenceOptional, "" }
func Required() (Presence, string)   { return PresenceRequired, "" }

// PresenceValid returns a problem string if the field's value is incompatible
// with its presence requirement (i.e., empty when a value is required, or
// non-empty when a value is not allowed).  It returns an empty string
// otherwise.
func (f *Field) PresenceValid() string {
	presence, when := f.Presence()
	if when != "" {
		when = " when " + when
	}
	var value = f.EditValue(f)
	switch presence {
	case PresenceNotAllowed:
		if value != "" {
			return fmt.Sprintf("The %q field cannot have a value%s.", f.Label, when)
		}
	case PresenceRequired:
		if value == "" {
			return fmt.Sprintf(`The %q field is required%s.`, f.Label, when)
		}
	}
	return ""
}

// AddFieldDefaults adds defaults to a Field.  It replaces nil functions with
// default implementations.
func AddFieldDefaults(f *Field) *Field {
	if f.Choices == nil {
		f.Choices = NoChoices{}
	}
	if f.Presence == nil {
		f.Presence = func() (Presence, string) { return PresenceOptional, "" }
	}
	if f.PIFOValid == nil {
		f.PIFOValid = func(f *Field) string { return "" }
	}
	if f.Compare == nil {
		f.Compare = func(label, exp, act string) *CompareField { return nil }
	}
	if f.PDFMap == nil {
		f.PDFMap = NoPDFField{}
	}
	if f.TableValue == nil {
		f.TableValue = func(f *Field) string {
			if f.Value == nil {
				return ""
			} else {
				return f.Choices.ToHuman(*f.Value)
			}
		}
	}
	if f.EditValue == nil {
		f.EditValue = func(f *Field) string {
			if f.Value == nil {
				return ""
			} else {
				return f.Choices.ToHuman(*f.Value)
			}
		}
	}
	if f.EditApply == nil {
		f.EditApply = func(f *Field, s string) {
			if f.Value != nil {
				*f.Value = f.Choices.ToPIFO(strings.TrimSpace(s))
			}
		}
	}
	if f.EditValid == nil {
		f.EditValid = f.PIFOValid
	}
	if f.EditSkip == nil {
		f.EditSkip = func(f *Field) bool {
			p, _ := f.Presence()
			return p == PresenceNotAllowed
		}
	}
	return f
}
