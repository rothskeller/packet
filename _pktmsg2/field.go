package pktmsg

// Value is a value of a form field in PackItForms encoding.  Its underlying
// type is string, but it is a separately defined type to distinguish it from
// human-readable representations of field values (which are plain strings).
// Use FromHuman to convert a human-readable representation to a Value and
// ToHuman to convert a Value to its human-readable representation.
type Value string

// FieldContainer is a container of fieldList.  It is a subset of Message, used to
// allow fields of a form to get the values of other fields duing validation,
// etc., but not to have access to the other methods of a Message.
type FieldContainer interface {
	// FieldValue returns the value of the field with the specified tag.  If
	// the provider has no such field, it returns an empty string.
	FieldValue(string) Value
	// KeyValue returns the value of the field with the specified key.  If
	// the provider has no such field, it returns an empty string.
	KeyValue(FieldKey) Value
}

// Field is the interface satisfied by all fields in all message types.
type Field interface {
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
	// therefore may change over time.  If pifo is true, only those problems
	// that would be flagged by PackItForms are returned.
	Validate(pifo bool) string
}

// fieldList is a list of fields, with methods on the list.
type fieldList []Field

func (fs fieldList) FieldCount() int { return len(fs) }
func (fs fieldList) FieldByIndex(idx int) Field {
	if idx < 0 || idx >= len(fs) {
		return nil
	}
	return fs[idx]
}
func (fs fieldList) FieldByTag(tag string) Field {
	for _, f := range fs {
		if f.Tag() == tag {
			return f
		}
	}
	return nil
}
func (fs fieldList) FieldValue(tag string) Value {
	if f := fs.FieldByTag(tag); f != nil {
		return f.Value()
	}
	return ""
}
func (fs fieldList) FieldByKey(key FieldKey) Field {
	for _, f := range fs {
		if f.Key() == key {
			return f
		}
	}
	return nil
}
func (fs fieldList) KeyValue(key FieldKey) Value {
	if f := fs.FieldByKey(key); f != nil {
		return f.Value()
	}
	return ""
}
func (fs fieldList) Calculate() {
	for _, f := range fs {
		f.Calculate()
	}
}
func (fs fieldList) Validate(pifo bool) (problems []string) {
	for _, f := range fs {
		if problem := f.Validate(pifo); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}

// BaseField returns a new base field with the specified characteristics.  None
// of the parameters have defaults, and all of them except key are required.
func BaseField(container FieldContainer, key FieldKey, tag, label string, width, height int, help string) Field {
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
func (bf *baseField) Validate(_ bool) string        { return "" }
