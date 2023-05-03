package xscmsg

import "github.com/rothskeller/packet/pktmsg"

// Message2 is the interface honored by all XSC message types.
type Message2 interface {
	// TypeTag returns the string that identifies the message type.  For
	// PackItForms messages, this is also the string that appears in the
	// subject line.
	TypeTag() string
	// TypeName returns the English name of the message type, in prose case.
	TypeName() string
	// TypeArticle returns "a" or "an", whichever is appropriate as the
	// indefinite article preceding TypeName.
	TypeArticle() string
	// FieldCount returns the number of fields in the message type.  It is
	// used as the bounds around interative calls to FieldByIndex.
	FieldCount() int
	// FieldByIndex returns the field at the specified index in the ordered
	// list of fields in the message.  This is the order of human
	// presentation, which may not be identical to the order of encoding.
	// The index is zero-based.  If the index is out of range, nil is
	// returned.
	FieldByIndex(int) FormField
	// FieldByTag returns the message field with the specified tag, or nil
	// if there is none.
	FieldByTag(string) FormField
	// FieldByKey returns the message field with the specified key, or nil
	// if there is none.
	FieldByKey(FieldKey) FormField
	// FieldValue returns the value of the field with the specified tag, or
	// an empty string if there is none.
	FieldValue(string) Value
	// KeyValue returns the value of the field with the specified key, or an
	// empty string if there is none.
	KeyValue(FieldKey) Value
	// Calculate recalculates the values of all fields in the message.
	Calculate()
	// Validate validates the values of all fields in the message, returning
	// a list of problem strings that is empty if the message is valid.
	Validate() []string
	// ValidatePIFO validates the values of all fields in the message,
	// restricting itself to only those requirements enforced by
	// PackItForms.  It returns a list of problem strings that is empty if
	// PackItForms would consider the message valid.
	ValidatePIFO() []string
	// EncodedSubject returns the encoded message subject line.
	EncodedSubject() string
	// EncodedBody returns the encoded message body.
	EncodeBody() string
}

// basePIFOForm is the base implementation for PIFO forms.  It shouldn't be used
// directly; use fixedPIFOForm or editablePIFOForm instead.
type basePIFOForm struct {
	typeTag     string
	typeName    string
	typeArticle string
	typeHTML    string
	pifoVersion string
	typeVersion string
	pifo        []FormField
	human       []FormField
}

func (rf basePIFOForm) TypeTag() string     { return rf.typeTag }
func (rf basePIFOForm) TypeName() string    { return rf.typeName }
func (rf basePIFOForm) TypeArticle() string { return rf.typeArticle }
func (rf basePIFOForm) FieldCount() int     { return len(rf.human) }

func (rf basePIFOForm) FieldByIndex(idx int) FormField {
	if idx < 0 || idx >= len(rf.human) {
		return nil
	}
	return rf.human[idx]
}

func (rf basePIFOForm) FieldByTag(tag string) FormField {
	for _, f := range rf.pifo {
		if f.Tag() == tag {
			return f
		}
	}
	return nil
}

func (rf basePIFOForm) FieldByKey(key FieldKey) FormField {
	for _, f := range rf.pifo {
		if f.Key() == key {
			return f
		}
	}
	return nil
}

func (rf basePIFOForm) FieldValue(tag string) Value {
	if f := rf.FieldByTag(tag); f != nil {
		return f.Value()
	}
	return ""
}

func (rf basePIFOForm) KeyValue(key FieldKey) Value {
	if f := rf.FieldByKey(key); f != nil {
		return f.Value()
	}
	return ""
}

func (rf basePIFOForm) Calculate() {
	for _, f := range rf.pifo {
		f.Calculate()
	}
}

func (rf basePIFOForm) Validate() (problems []string) {
	for _, f := range rf.human {
		if problem := f.Validate(); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}

// ValidatePIFO validates the values of all fields in the message,
// restricting itself to only those requirements enforced by
// PackItForms.  It returns a list of problem strings that is empty if
// PackItForms would consider the message valid.
func (rf basePIFOForm) ValidatePIFO() (problems []string) {
	for _, f := range rf.pifo {
		if problem := f.ValidatePIFO(); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}

func (rf basePIFOForm) EncodeBody() string {
	var form = pktmsg.Form{
		PIFOVersion: rf.pifoVersion,
		FormType:    rf.typeHTML,
		FormVersion: rf.typeVersion,
	}
	for _, f := range rf.pifo {
		if v := f.Value(); v != "" {
			form.Fields = append(form.Fields, pktmsg.FormField{Tag: f.Tag(), Value: string(v)})
		}
	}
	return form.Encode()
}

// FixedPIFOForm returns a Message2 for a PackItForms message that has been sent
// over the air and therefore should not be freely changed.
func FixedPIFOForm(m *pktmsg.Message) Message2 {

}

// A fixedPIFOForm is a message that has been sent over the air and therefore
// shouldn't be changed.  (This isn't strictly true; we will change it to add
// a destination message number when a receipt is received, and things like
// that, but it's not freely modifiable by the user.)
type fixedPIFOForm struct {
	basePIFOForm
	subject string
}

// For a fixed message, EncodedSubject returns the subject line as it was
// transmitted; it doesn't try to re-calculate it.
func (rf fixedPIFOForm) EncodedSubject() string { return rf.subject }
