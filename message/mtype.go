package message

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/msgifc"
)

// MType is the interface satisfied by all messages types.
type MType = msgifc.MType

// EditableMType is the interface satisfied by a message type that allows
// creation and editing of messages.
type EditableMType = msgifc.EditableMType

type EditHTMLVars = msgifc.EditHTMLVars

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
func RegisterType(mt MType) error {
	return registerType(mt, false)
}

// RegisterFallbackType registers a message type whose recognizer is guaranteed
// to be called after all types registered with RegisterType, even if they are
// registered later.
func RegisterFallbackType(mt MType) error {
	return registerType(mt, true)
}

func registerType(mt MType, fallback bool) error {
	var tag, key string

	if mt, ok := mt.(EditableMType); ok {
		tag, key = strings.ToLower(mt.CreateTag()), strings.ToLower(mt.CreateKey())
		if tag == "" {
			return errors.New("editable message types must have a CreateTag")
		}
		for _, exist := range mtypes {
			if exist, ok := exist.mt.(EditableMType); ok {
				etag, ekey := strings.ToLower(exist.CreateTag()), strings.ToLower(exist.CreateKey())
				if tag == etag {
					return fmt.Errorf("the CreateTag %q is already in use", mt.CreateTag())
				}
				if key == etag {
					return fmt.Errorf("the CreateKey %q is already in use", mt.CreateKey())
				}
				if key != "" && ekey != "" {
					if strings.HasPrefix(key, ekey) || strings.HasPrefix(ekey, key) {
						return fmt.Errorf("the CreateKey %q conflicts with another CreateKey %q", mt.CreateKey(), exist.CreateKey())
					}
				}
			}
		}
	}
	mtypes = append(mtypes, registeredMType{mt: mt, fallback: fallback})
	return nil
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

// FindTypeTag returns the registered message type with the specified
// CreateTag, or nil if none exists.
func FindTypeTag(tag string) EditableMType {
	for mt := range AllTypes() {
		if emt, ok := mt.(EditableMType); ok && emt.CreateTag() == tag {
			return emt
		}
	}
	return nil
}

// CompareTypes returns -1/0/+1 for sorting a list of types.  Types are sorted
// by name, except that plain text message always comes first.
func CompareTypes(a, b MType) int {
	if a == PlainMessage && b != PlainMessage {
		return -1
	}
	if b == PlainMessage && a != PlainMessage {
		return +1
	}
	an, bn := strings.ToLower(a.Name()), strings.ToLower(b.Name())
	_, an, _ = strings.Cut(an, " ") // remove "a" or "an"
	_, bn, _ = strings.Cut(bn, " ") // remove "a" or "an"
	return cmp.Compare(an, bn)
}

//-----------------------------------------------------------------------------

// BaseMType is the common core implementation for all message types.
type BaseMType struct {
	name   string
	fields []msgifc.Field
}

// NewBaseMType returns a new BaseMType with the specified details to be
// returned by Name and CreateTag methods.
func NewBaseMType(name string) *BaseMType {
	return &BaseMType{name: name}
}

// Name returns the name of the message type, as a phrase in lower case (other
// than acronyms) starting with "a " or "an ".
func (t *BaseMType) Name() string { return t.name }

// Validate validates the contents of the message and returns any problems.
// If pifo is true, it returns only problems that PackItForms would raise;
// otherwise the checks may be more extensive.
func (t *BaseMType) Validate(m Message, flags msgifc.ValidateFlags) (err error) {
	for f := range t.Fields(m) {
		err = errors.Join(err, f.Validate(m, f, flags))
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
func (t *BaseMType) Fields(m Message) iter.Seq[field.Field] {
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

//-----------------------------------------------------------------------------

// BaseEditableMType is the common core implementation for all editable message
// types.  Note that it is not a complete implementation; message types must
// implement their own NewDraft and NewDraftCopy methods.
type BaseEditableMType struct {
	BaseMType
	createTag string
	createKey string
}

// NewBaseEditableMType returns a new BaseEditableMType with the specified
// details to be returned by Name and CreateTag methods.
func NewBaseEditableMType(name, createTag, createKey string) *BaseEditableMType {
	return &BaseEditableMType{BaseMType: BaseMType{name: name}, createTag: createTag, createKey: createKey}
}

// CreateTag returns the tag used to identify this message type on a
// "new" command line.  The tag is case insensitive.  It is an error
// for two message types to have the same tag.
func (t *BaseEditableMType) CreateTag() string { return t.createTag }

// CreateKey returns the key (usually one or two letters) used to identify this
// message type in a GUI dialog box for creating a new message (or also on a
// "new" command line, as an alternative to CreateTag).  The key is case
// insensitive.  It is an error for two message types to have the same key, or
// for one to have a key that is a prefix of another's, or for any key to be
// the same as any CreateTag.  This method may return an empty string, in which
// case there is no shortcut key in the dialog and the message type must be
// selected with mouse or arrow keys.)
func (t *BaseEditableMType) CreateKey() string { return t.createKey }
