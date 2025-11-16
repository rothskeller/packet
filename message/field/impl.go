package field

import (
	"strings"

	"github.com/phpdave11/gofpdf"
	"github.com/rothskeller/packet/errors"
)

// field is the implementation of Field created by a FieldFactory.
type field struct {
	tag            string
	label          string
	parent         Field
	children       []Field
	editHelp       string
	editHint       string
	multiline      bool
	editWidth      int
	editHeight     int
	obscured       bool
	restricted     bool
	valueFunc      func(Message) string
	toHumanFunc    func(string) string
	fromHumanFunc  func(string) string
	setValueFunc   func(Message, string)
	visibleFunc    func(Message) bool
	editableFunc   func(Message, bool) bool
	choicesFunc    func(Message) []ChoicePair
	requiredFunc   func(Message) bool
	requiredDesc   string
	disallowedFunc func(Message) bool
	disallowedDesc string
	validateFuncs  []func(Message, bool) error

	// internal transient data
	disallowed bool
}

var _ Field = (*field)(nil)

func (f *field) Tag() string          { return f.tag }
func (f *field) Label() string        { return f.label }
func (f *field) Parent() Field        { return f.parent }
func (f *field) Children() []Field    { return f.children }
func (f *field) EditHelp() string     { return f.editHelp }
func (f *field) EditHint() string     { return f.editHint }
func (f *field) Multiline() bool      { return f.multiline }
func (f *field) EditSize() (int, int) { return f.editWidth, f.editHeight }
func (f *field) Obscured() bool       { return f.obscured }
func (f *field) Restricted() bool     { return f.restricted }

func (f *field) Value(m Message) string {
	if f.valueFunc != nil {
		return f.valueFunc(m)
	}
	return ""
}

func (f *field) ToHuman(s string) string {
	if f.toHumanFunc != nil {
		return f.toHumanFunc(s)
	}
	return s
}

func (f *field) FromHuman(s string) string {
	if f.fromHumanFunc != nil {
		return f.fromHumanFunc(s)
	}
	return strings.TrimSpace(s)
}

func (f *field) SetValue(m Message, val string) {
	if f.setValueFunc != nil {
		f.setValueFunc(m, val)
	} else {
		panic("SetValue called on non-settable field")
	}
	// If any other fields were newly made disallowed, remove their values.
	for of := range m.Fields() {
		if of, ok := of.(*field); ok && !of.disallowed && of.disallowedFunc != nil {
			if of.disallowedFunc(m) {
				of.SetValue(m, "")
			}
		}
	}
}

func (f *field) Visible(m Message) bool {
	if f.visibleFunc != nil {
		return f.visibleFunc(m)
	}
	return true
}

func (f *field) Editable(m Message, explicit bool) bool {
	if f.editableFunc != nil {
		return f.editableFunc(m, explicit)
	}
	return f.editHelp != ""
}

func (f *field) Choices(m Message) []ChoicePair {
	if f.choicesFunc != nil {
		return f.choicesFunc(m)
	}
	return nil
}

// Validate validates the value of the field and returns any problems with it.
// If pifo is true, it restricts itself to those checks performed by
// PackItForms.
func (f *field) Validate(m Message, fi Field, pifo bool) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	return f.validateCustom(m, fi, pifo)
}

// validatePresence validates the value of the field for compliance with
// required or disallowed rules.
func (f *field) validatePresence(m Message, fi Field) (err error) {
	val := fi.Value(m)
	// If we have a value, check whether any value is allowed.
	f.disallowed = false
	if f.disallowedFunc != nil && !f.disallowedFunc(m) {
		if val == "" {
			return nil
		}
		f.disallowed = true
		if f.disallowedDesc != "" {
			return errors.NewF("The %q field cannot have a value unless %s.", f.label, f.disallowedDesc)
		}
		return errors.NewF("The %q field cannot have a value.", f.label)
	}
	// If we don't have a value, check whether one is required.
	if val == "" && f.requiredFunc != nil && f.requiredFunc(m) {
		if f.requiredDesc != "" {
			return errors.NewF("The %q field is required when %s.", f.label, f.requiredDesc)
		} else if f.disallowedDesc != "" {
			return errors.NewF("The %q field is required when %s.", f.label, f.disallowedDesc)
		} else {
			return errors.NewF("The %q field is required.", f.label)
		}
	}
	return nil
}

// validateCustom validates the value of the field for compliance with
// restricted choices and with any custom validation handlers.
func (f *field) validateCustom(m Message, fi Field, pifo bool) (err error) {
	val := fi.Value(m)
	// If we have a value and restricted choices, check to be sure the value
	// is one of the allowed ones.
	if val != "" && f.restricted {
		var found bool

		for _, choice := range f.Choices(m) {
			if val == choice.PIFO {
				found = true
				break
			}
		}
		if !found {
			return errors.NewF("The %q field does not contain one of its allowed values.", f.label)
		}
	}
	// Run any custom validation checks.
	for _, fn := range f.validateFuncs {
		err = errors.Join(err, fn(m, pifo))
	}
	return err
}

// Compare compares the value of the field in the actual message to
// the corresponding value in the expected message, and returns the
// results of the comparison.
func (f *field) Compare(expected Message, actual Message) *ComparedField {
	panic("not implemented") // TODO: Implement
}

// RenderPDF renders the field onto the specified page of the specified
// PDF file, if it belongs there.  It may return any errors in the
// process (unsupported value, doesn't fit in the space, etc.).
func (f *field) RenderPDF(m Message, pdf *gofpdf.Pdf, page int) error {
	panic("not implemented") // TODO: Implement
}
