package msgifc

// Field is the interface satisfied by all message fields.  Note that this is
// more than just PackItForms fields; there are fields relating to metadata as
// well, and non-forms messages have fields also.
type Field interface {
	// Tag returns the PackItForms tag for the field, if any.
	Tag() string
	// Common is the tag string that identifies the field as one of the
	// well-known common fields.  It allows software to address those
	// fields without dependency on how they're stored in a particular
	// message type.
	Common() string
	// Label returns the label for the field, if any.
	Label() string
	// Parent returns the parent field that contains this field, if any.
	Parent() Field
	// Children returns the list of child fields contained by this field,
	// if any.  A Field may not have both a Parent and Children.
	Children() []Field
	// Value returns the value of the field, in internal form.
	Value(Message) string
	// Default returns the default value of the field, in internal form.
	Default() string
	// ToHuman converts the internal form of a value for the field into the
	// human form appropriate for display and editing (often a no-op).
	ToHuman(Message, string) string
	// FromHuman converts the supplied value from human form to internal
	// form, if possible; otherwise it makes no changes.  Implementations
	// must not change the value if it is already in internal form.
	FromHuman(Message, string) string
	// SetValue sets the value of the field.  The supplied value must be in
	// internal form.
	SetValue(Message, string)
	// Visible returns whether the field should be included when the
	// message is displayed.  Note that fields with an empty value are
	// never displayed no matter what this method returns.
	Visible(Message) bool
	// Editable returns whether the field should be included when the
	// message is edited.  The explicit flag is true if the user explicitly
	// asked to edit this field by name.
	Editable(m Message, explicit bool) bool
	// EditHelp returns the help string for editing of the field.
	EditHelp() string
	// EditHint returns the hint string, if any, for editing of the field.
	// This is displayed in the editing control when the control is
	// otherwise empty (the equivalent of HTML <input placeholder="...">).
	EditHint() string
	// Multiline returns whether the field is expected to contain a
	// multiline value, i.e., a value containing newlines.  All values
	// *can* include newlines, but this is whether it's *expected* to.
	Multiline() bool
	// EditSize returns the width and height, in characters, of the text
	// entry control for the field.  Zero return values mean unlimited.
	// (Usually this should be set based on the size of the corresponding
	// printable area of the PDF at the minimum font size.)
	EditSize() (int, int)
	// Obscured returns whether the field's value should be obscured for
	// display or editing.  This is used for password fields.
	Obscured() bool
	// Choices returns a list of allowed or recommended values for the
	// field.  Each element is a pair with internal and human
	// representations of the value.  The list may vary depending on the
	// values of other fields of the message.
	Choices(Message) []ChoicePair
	// Restricted returns whether the value of the field is restricted to
	// one of the listed Choices (true), or whether they are just
	// recommendations and any other value is accepted (false).
	Restricted() bool
	// Validate validates the value of the field and returns any problems
	// with it.  flags customizes the validation.
	Validate(m Message, f Field, flags ValidateFlags) error
	// Compare compares the value of the field in the actual message to
	// the corresponding value in the expected message, and returns the
	// results of the comparison.
	// Compare(expected, actual Message) *ComparedField
	// RenderPDF renders the field onto the specified page of the specified
	// PDF file, if it belongs there.  It may return any errors in the
	// process (unsupported value, doesn't fit in the space, etc.).
	// RenderPDF(m Message, pdf *gofpdf.Pdf, page int) error
}

// ChoicePair is a pair of strings representing a choice for a field value.
type ChoicePair struct {
	PIFO  string
	Human string
}
