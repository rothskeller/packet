package field

import (
	"github.com/rothskeller/packet/v4/message/msgifc"
)

// A FieldFactory is used to create a Field implementation customized for a
// specific message field.  To use it, call NewField, which returns a
// FieldFactory.  Then call methods on the FieldFactory to customze the field.
// Pass the resulting FieldFactory to a message's AddField method to finalize
// the Field implementation and add it to the message.
type FieldFactory struct {
	tf Field
	f  field
}

// NewField creates a FieldFactory that will be used to create a Field
// implementation for a specific message field.  The parameters are the
// PackItForms tag for the field, if any, and the display label for the field.
// Call methods on the returned FieldFactory to customize the field, then pass
// the FieldFactory to a message's AddField method to finalize the field
// implementation and add it to the message.
func NewField(tag, label string) (ff *FieldFactory) {
	ff = &FieldFactory{f: field{tag: tag, label: label}}
	ff.tf = &ff.f
	return ff
}

// AddField adds the supplied field as a child of the field being created by
// the receiver factory.  It returns the finalized child Field.
func (ff *FieldFactory) AddField(cf *FieldFactory) Field {
	if ff.f.parent != nil {
		panic("can't nest fields more than one layer")
	}
	if cf.f.parent != nil {
		panic("field can't be child of multiple parents")
	}
	if len(cf.f.children) != 0 {
		panic("can't nest fields more than one layer")
	}
	cf.f.parent = ff.tf
	ff.f.children = append(ff.f.children, cf.MakeField())
	return cf.tf
}

// Common sets the common field ID for the field, if any.
func (ff *FieldFactory) Common(c string) *FieldFactory {
	ff.f.common = c
	return ff
}

// ValueFunc provides a function that returns the value of the field, in
// internal form.  The default implementation returns an empty string unless
// the message is a form and the tag passed to NewField is non-empty; in that
// case, it returns the value of the corresponding PackItForms field.
func (ff *FieldFactory) ValueFunc(fn func(msgifc.Message) string) *FieldFactory {
	ff.f.valueFunc = fn
	return ff
}

// ToHumanFunc provides a function that converts the internal form of a value
// for the field into the human form appropriate for display and editing.  The
// default is an identity conversion.
func (ff *FieldFactory) ToHumanFunc(fn func(msgifc.Message, string) string) *FieldFactory {
	ff.f.toHumanFunc = fn
	return ff
}

// FromHumanFunc provides a function that converts the supplied value from
// human form to internal form, if possible; otherwise it makes no changes.
// The default is an identity conversion.
func (ff *FieldFactory) FromHumanFunc(fn func(msgifc.Message, string) string) *FieldFactory {
	ff.f.fromHumanFunc = fn
	return ff
}

// SetValueFunc provides a function that sets the value of the field.  The
// supplied value will be in internal form.  The default implementation panics
// unless the message is a form and the tag passed to NewField is non-empty; in
// that case, it sets the value of the corresponding PackItForms field.
func (ff *FieldFactory) SetValueFunc(fn func(msgifc.Message, string)) *FieldFactory {
	ff.f.setValueFunc = fn
	return ff
}

// VisibleWhen provides a predicate function that determines when the field
// should be included in displayed messages.  The default predicate always
// returns true.  Note that fields with empty values are never displayed.
func (ff *FieldFactory) VisibleWhen(pred func(msgifc.Message) bool) *FieldFactory {
	ff.f.visibleFunc = pred
	return ff
}

// Invisible is a predicate for VisibleWhen that always returns false.
func Invisible(msgifc.Message) bool { return false }

// EditableWhen provides a predicate function that determines when the field's
// value can be edited by a user.  The explicit flag passed to the predicate
// tells it whether the user is explicitly asking to edit this specific field
// by name (true) or whether editing it is being considered as part of the
// sequence of fields (false).  The default predicate returns true if EditHelp
// has been called for the field and false otherwise.
func (ff *FieldFactory) EditableWhen(pred func(msgifc.Message, bool) bool) *FieldFactory {
	ff.f.editableFunc = pred
	return ff
}

// NotEditable is a predicate for EditableWhen that always returns false.
func NotEditable(_ msgifc.Message, _ bool) bool { return false }

// OnlyExplicitlyEditable is a predicate for EditableWhen that only returns true
// when the field is explicitly requested.
func OnlyExplicitlyEditable(_ msgifc.Message, explicit bool) bool { return explicit }

// EditHelp sets the help string for editing the field.
func (ff *FieldFactory) EditHelp(help string) *FieldFactory {
	ff.f.editHelp = help
	return ff
}

// EditHint sets the hint string, if any, for editing of the field.
func (ff *FieldFactory) EditHint(hint string) *FieldFactory {
	ff.f.editHint = hint
	return ff
}

// Multiline marks the field as expected to contain a multiline value.
func (ff *FieldFactory) Multiline() *FieldFactory {
	ff.f.multiline = true
	return ff
}

// EditWidth sets the width, in characters, of the text entry control for the
// field.  It sets the height to 1.  The default is zero, meaning to use the
// full available width and height.
func (ff *FieldFactory) EditWidth(width int) *FieldFactory {
	ff.f.editWidth = width
	ff.f.editHeight = 1
	return ff
}

// EditSize sets the width and height, in characters, of the text entry control
// for the field.  The default is zero, meaning to use the full available width
// and height.
func (ff *FieldFactory) EditSize(width, height int) *FieldFactory {
	ff.f.editWidth = width
	ff.f.editHeight = height
	return ff
}

// Obscured marks the field as containing values that should be obscured for
// display or editing (e.g., password fields).
func (ff *FieldFactory) Obscured() *FieldFactory {
	ff.f.obscured = true
	return ff
}

// AllowedValues sets the list of allowed values for the field; other values
// cause validation failures.  Each supplied value must be either a ChoicePair
// or a string (used when the internal and human forms of the value are the
// same).
func (ff *FieldFactory) AllowedValues(values ...any) *FieldFactory {
	ff.setChoicesFunc(values)
	ff.f.restricted = true
	return ff
}

// AllowedValuesFunc provides a function that returns the (runtime-variable)
// list of allowed values for the field; other values cause validation
// failures.
func (ff *FieldFactory) AllowedValuesFunc(fn func(msgifc.Message) []msgifc.ChoicePair) *FieldFactory {
	ff.f.choicesFunc = fn
	ff.f.restricted = true
	return ff
}

// SuggestedValues sets the list of suggested values for the field; other
// values are also allowed.  Each supplied value must be either a ChoicePair
// or a string (used when the internal and human forms of the value are the
// same.)
func (ff *FieldFactory) SuggestedValues(values ...any) *FieldFactory {
	ff.setChoicesFunc(values)
	return ff
}

// SuggestedValuesFunc provides a function that returns the (runtime-variable)
// list of suggested values for the field; other values are also allowed.
func (ff *FieldFactory) SuggestedValuesFunc(fn func(msgifc.Message) []msgifc.ChoicePair) *FieldFactory {
	ff.f.choicesFunc = fn
	return ff
}

// Required marks a field as required: it must have a non-empty value.
func (ff *FieldFactory) Required() *FieldFactory {
	ff.f.requiredFunc = func(msgifc.Message) bool { return true }
	return ff
}

// RequiredWhen provides a predicate function to determine when the field is
// required.  It also provides a string to describe in English what the
// predicate checks.  The string should interpolate into the sentence "The XXX
// field is required when %s."
func (ff *FieldFactory) RequiredWhen(when string, pred func(msgifc.Message) bool) *FieldFactory {
	ff.f.requiredFunc = pred
	ff.f.requiredDesc = when
	return ff
}

// DisallowedUnless provides a predicate function to determine when the field
// is allowed to have a value.  It also provides a string to describe in
// English what the predicate checks.  The string should interpolate into the
// sentence "The XXX field cannot have a value unless %s."  If the field has
// been marked Required, the DisallowUnless predicate takes precedence, and the
// string should also interpolate into the sentence "The XXX field is required
// when %s."
func (ff *FieldFactory) DisallowedUnless(unless string, pred func(msgifc.Message) bool) *FieldFactory {
	ff.f.disallowedFunc = pred
	ff.f.disallowedDesc = unless
	return ff
}

// ValidateFunc provides a function to perform additional validation on the
// value of the field (beyond that performed by Required, RequiredWhen,
// DisallowedUnless, AllowedValues[Func], or one of the data type methods).  If
// the flag passed to the function is true, the function should restrict itself
// to those checks that PackItForms would enforce.  This method can be called
// multiple times to register multiple validation functions for the field.
func (ff *FieldFactory) ValidateFunc(fn func(msgifc.Message, msgifc.Field, msgifc.ValidateFlags) error) *FieldFactory {
	ff.f.validateFunc = fn
	return ff
}

// CompareFunc provides a function to compare the actual value of the field
// from one message to the expected value in another message.
func (ff *FieldFactory) CompareFunc(fn func(label, exp, act string) *ComparedField) *FieldFactory {
	ff.f.compareFunc = fn
	return ff
}

// MakeField resolves the factory and returns the constructed Field.
func (ff *FieldFactory) MakeField() Field {
	if ff.f.editHeight == 0 && !ff.f.multiline {
		ff.f.editHeight = 1
	}
	return ff.tf
}

func (ff *FieldFactory) setChoicesFunc(values []any) {
	var pairs []msgifc.ChoicePair
	var maxlen int

	for _, value := range values {
		switch value := value.(type) {
		case msgifc.ChoicePair:
			pairs = append(pairs, value)
			maxlen = max(maxlen, len(value.PIFO))
		case string:
			pairs = append(pairs, msgifc.ChoicePair{PIFO: value, Human: value})
			maxlen = max(maxlen, len(value))
		default:
			panic("Allowed/Suggested value is not string or ChoicePair")
		}
	}
	ff.f.editWidth = maxlen
	ff.f.editHeight = 1
	ff.f.choicesFunc = func(msgifc.Message) []msgifc.ChoicePair { return pairs }
}
