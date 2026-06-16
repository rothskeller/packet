package message

import (
	"cmp"
	"fmt"
	"iter"
	"strings"

	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/message/msgifc"
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

// FindCreateTag returns the registered message type with the specified
// CreateTag, or nil if none exists.
func FindCreateTag(tag string) EditableMType {
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

// BaseMType is a common core implementation for some message types.
type BaseMType struct {
	tag  string
	name string
}

// NewBaseMType returns a new BaseMType with the specified name.
func NewBaseMType(tag, name string) *BaseMType {
	return &BaseMType{tag: tag, name: name}
}

func (t *BaseMType) Tag() string  { return t.tag }
func (t *BaseMType) Name() string { return t.name }

// RenderPDF creates a PDF representation of the message in the specified file.
func (t *BaseMType) RenderPDF(m Message, filename, copyname string) (err error) {
	return RenderPlainPDF(m, filename, copyname)
}

//-----------------------------------------------------------------------------

// BaseEditableMType is a common core implementation for editable message
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
	return &BaseEditableMType{BaseMType: BaseMType{tag: createTag, name: name}, createTag: createTag, createKey: createKey}
}

func (t *BaseEditableMType) CreateTag() string { return t.createTag }

func (t *BaseEditableMType) CreateKey() string { return t.createKey }
