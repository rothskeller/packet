package form

import (
	"github.com/rothskeller/packet/form/formdef"
	"github.com/rothskeller/packet/message/field"
)

// ff2mf is an adapter that implements the message.Field interface for a
// formdef.FieldDef structure.
type ff2mf struct{ fd *formdef.FieldDef }

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
func (f ff2mf) Value(m field.Message) string {
	if f.fd.Tag != "" {
		return m.Body().(*FormBody).Field(f.fd.Tag)
	}
	return f.fd.Value
}

// Default returns the default value of the field, in internal form.
func (f ff2mf) Default() string { return f.fd.Value }

// ToHuman converts the internal form of a value for the field into the
// human form appropriate for display and editing (often a no-op).
func (f ff2mf) ToHuman(raw string) string {
	for _, c := range f.fd.Choices {
		if c.Raw == raw {
			return c.Human
		}
	}
	return raw
}

// FromHuman converts the supplied value from human form to internal
// form, if possible; otherwise it makes no changes.  Implementations
// must not change the value if it is already in internal form.
func (f ff2mf) FromHuman(human string) string {
	for _, c := range f.fd.Choices {
		if c.Human == human {
			return c.Raw
		}
	}
	return human
}

// SetValue sets the value of the field.  The supplied value must be in
// internal form.
func (f ff2mf) SetValue(msg field.Message, val string) {
	if f.fd.Tag != "" {
		reason := "body.FormBody.Field." + f.fd.Tag
		if f.fd.Common != "" {
			reason = "body.FormBody.Common." + f.fd.Common
		}
		msg.Body().(*FormBody).SetFieldR(f.fd.Tag, val, reason)
	}
}

// Visible returns whether the field should be included when the
// message is displayed.  Note that fields with an empty value are
// never displayed no matter what this method returns.
func (f ff2mf) Visible(field.Message) bool {
	return f.fd.Parent == nil
}

// Editable returns whether the field should be included when the
// message is edited.  The explicit flag is true if the user explicitly
// asked to edit this field by name.  Note: as a side effect, this
// method is allowed to clear the value of the field before returning
// false when that's more appropriate than a validation failure.
func (f ff2mf) Editable(m field.Message, explicit bool) bool {
	if f.fd.EditHelp == "" {
		return false
	}
	if f.fd.Type == "dateTime" && !explicit {
		return false
	}
	return true
	// TODO: apply presence
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
func (f ff2mf) EditSize() (int, int) { return f.fd.EditWidth, 1 }

// Obscured returns whether the field's value should be obscured for
// display or editing.  This is used for password fields.
func (f ff2mf) Obscured() bool { return f.fd.Type == "password" }

// Choices returns a list of allowed or recommended values for the
// field.  Each element is a pair with internal and human
// representations of the value.  The list may vary depending on the
// values of other fields of the message.
func (f ff2mf) Choices(field.Message) (cs []field.ChoicePair) {
	if f.fd.Type == "checkbox" {
		return []field.ChoicePair{{PIFO: "", Human: ""}, {PIFO: "checked", Human: "checked"}}
	}
	for _, c := range f.fd.Choices {
		// TODO: test condition
		cs = append(cs, field.ChoicePair{PIFO: c.Raw, Human: c.Human})
	}
	return cs
}

// Restricted returns whether the value of the field is restricted to
// one of the listed Choices (true), or whether they are just
// recommendations and any other value is accepted (false).
func (f ff2mf) Restricted() bool { return f.fd.Type == "restricted" || f.fd.Type == "checkbox" }

// Validate validates the value of the field and returns any problems
// with it.  If pifo is true, it restricts itself to those checks
// performed by PackItForms.
func (f ff2mf) Validate(m field.Message, mf field.Field, pifo bool) error {
	panic("not implemented")
}
