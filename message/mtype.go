package message

import (
	"iter"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/field"
)

// MType is the interface satisfied by all messages types.
type MType interface {
	// Recognize detects whether the argument message is of the type
	// described by this MType, and if so, calls its SetType method to
	// assign this MType to it.
	Recognize(Message)
	// CreateTag returns the tag(s) used to identify this message type on a
	// "new" command line.  Multiple tags may be returned, separated by
	// spaces.  Version numbers may be included in a tag; see the "new"
	// command documentation for details.  Message types that cannot be
	// created with "new" return an empty string.
	CreateTag() string
	// Name returns the name of the message type, as a phrase in lower case
	// (other than acronyms) starting with "a " or "an ".
	Name() string
	// Validate validates the contents of the message and returns any
	// problems.  If pifo is true, it returns only problems that
	// PackItForms would raise; otherwise the checks may be more extensive.
	Validate(m Message, pifo bool) error
	// Fields returns an iterator on the set of message fields.
	Fields() iter.Seq[field.Field]
	// RenderPDF creates a PDF representation of the message in the
	// specified file.  If copyname is not empty, it is placed in the
	// footer of each page.  The returned error may be a Warning, showing a
	// non-fatal rendering issue.
	RenderPDF(m Message, filename, copyname string) error
}

// Warning wraps a non-fatal "error".
type Warning struct{ Err error }

func (w Warning) Error() string { return w.Err.Error() }
func (w Warning) Unwrap() error { return w.Err }

//-----------------------------------------------------------------------------

// mtypes is the list of registered message types.
type registeredMType struct {
	mt       MType
	fallback bool
}

var mtypes []registeredMType

// RegisterType registers a message type.
func RegisterType(mt MType) {
	mtypes = append(mtypes, registeredMType{mt: mt})
}

// RegisterFallbackType registers a message type whose recognizer is guaranteed
// to be called after all types registered with RegisterType, even if they are
// registered later.
func RegisterFallbackType(mt MType) {
	mtypes = append(mtypes, registeredMType{mt: mt, fallback: true})
}

// SetType determines the type of a Message and sets its Type property to the
// first registered type whose Recognize method recognizes it, or otherwise to
// PlainMessage.
func SetType(m Message) {
	if m.Type() != nil {
		return
	}
	for mt := range AllTypes() {
		mt.Recognize(m)
		if m.Type() != nil {
			return
		}
	}
	m.SetType(PlainMessage)
}

// AllTypes returns an iterator on all of the registered message types.
func AllTypes() iter.Seq[MType] {
	return func(yield func(MType) bool) {
		for _, mt := range mtypes {
			if !mt.fallback && !yield(mt.mt) {
				return
			}
		}
		for _, mt := range mtypes {
			if mt.fallback && !yield(mt.mt) {
				return
			}
		}
		yield(PlainMessage)
	}
}

// FindType returns a registered message type that satisfies the supplied
// predicate, or nil if none does.
func FindType(pred func(MType) bool) MType {
	for mt := range AllTypes() {
		if pred(mt) {
			return mt
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

// BaseMType is the common core implementation for all message types.
type BaseMType struct {
	createTag string
	name      string
	fields    []field.Field
}

// NewBaseMType returns a new BaseMType with the specified details to be
// returned by Name and CreateTag methods.
func NewBaseMType(name, createTag string) *BaseMType {
	return &BaseMType{name: name, createTag: createTag}
}

// CreateTag returns the tag(s) used to identify this message type on a "new"
// command line.  Multiple tags may be returned, separated by spaces.  Version
// numbers may be included in a tag; see the "new" command documentation for
// details.  Message types that cannot be created with "new" return an empty
// string.
func (t *BaseMType) CreateTag() string {
	return t.createTag
}

// Name returns the name of the message type, as a phrase in lower case (other
// than acronyms) starting with "a " or "an ".
func (t *BaseMType) Name() string { return t.name }

// Validate validates the contents of the message and returns any problems.
// If pifo is true, it returns only problems that PackItForms would raise;
// otherwise the checks may be more extensive.
func (t *BaseMType) Validate(m Message, pifo bool) (err error) {
	for f := range t.Fields() {
		err = errors.Join(err, f.Validate(m, f, pifo))
	}
	return err
}

// AddField adds fields to a form definition.  It is usually called by the
// form definition's initFunc.  The arguments can be either field.Field
// implementations or *field.FieldFactory instances.
func (t *BaseMType) AddField(fs ...any) {
	for _, f := range fs {
		switch f := f.(type) {
		case field.Field:
			t.fields = append(t.fields, f)
		case *field.FieldFactory:
			t.fields = append(t.fields, f.MakeField())
		default:
			panic("AddField arguments must be field.Field or *field.FieldFactory")
		}
	}
}

// Fields returns an iterator on the set of message fields.
func (t *BaseMType) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		for _, f := range t.fields {
			if !yield(f) {
				return
			}
			for _, cf := range f.Children() {
				if !yield(cf) {
					return
				}
			}
		}
	}
}

// RenderPDF creates a PDF representation of the message in the specified file.
func (t *BaseMType) RenderPDF(m Message, filename, copyname string) (err error) {
	return RenderPlainPDF(m, filename, copyname)
}
