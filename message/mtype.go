package message

import (
	"cmp"
	"fmt"
	"io/fs"
	"iter"
	"net/http"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/field"
)

// MType is the interface satisfied by all messages types.
type MType interface {
	// Recognize detects whether the argument message is of the type
	// described by this MType, and if so, calls its SetType method to
	// assign this MType to it.
	Recognize(Message)
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

// EditableMType is the interface satisfied by a message type that allows
// creation and editing of messages.
type EditableMType interface {
	MType
	// CreateTag returns the tag used to identify this message type on a
	// "new" command line.  The tag is case insensitive.  It is an error
	// for two message types to have the same tag.
	CreateTag() string
	// CreateKey returns the key (usually one or two letters) used to
	// identify this message type in a GUI dialog box for creating a new
	// message (or also on a "new" command line, as an alternative to
	// CreateTag).  The key is case insensitive.  It is an error for two
	// message types to have the same key, or for one to have a key that is
	// a prefix of another's, or for any key to be the same as any
	// CreateTag.  This method may return an empty string, in which case
	//there is no shortcut key in the dialog and the message type must be
	// selected with mouse or arrow keys.)
	CreateKey() string
	// NewDraft returns a new *message.DraftMessage of this type.  It has
	// default values filled in but is otherwise empty.
	NewDraft() *DraftMessage
	// EditHTML returns the HTML for the edit page to edit the supplied
	// DraftMessage of this type.  vars customize the HTML based on the
	// the editing context.
	EditHTML(msg *DraftMessage, vars EditHTMLVars) ([]byte, error)
	// EditAssets returns the file system containing assets used by the
	// HTML returned by EditHTML.
	EditAssets() fs.FS
	// FromPOST interprets the form POSTed by the HTML returned by EditHTML,
	// and translates it into a DraftMessage of this message type.  If the
	// POSTed form is invalid, FromPOST may return an error instead.
	FromPOST(r *http.Request) (*DraftMessage, error)
}

type EditHTMLVars struct {
	// SubmitURL is the URL that the edit form should POST to.  The URL
	// should respond with either an error status with a text/plain body,
	// or an http.StatusSeeOther with a redirect.
	SubmitURL string
	// SubmitLabel is the label of the submit button.  It is required.
	SubmitLabel string
	// AltURL is the URL that the form's alternate submit button should
	// POST to.  It is optional; if not present, the alternate submit
	// button is not shown.  The URL should respond with either an error
	// status with a text/plain body containing an error message, or an
	// http.StatusSeeOther with a redirect.
	AltURL string
	// AltLabel is the label of the alternate submit button.  It is
	// required if AltURL is specified and ignored otherwise.
	AltLabel string
	// AssetBase is the URL base for any assets needed by the edit form.
	// It should correspond to the file system returned by EditAssets.
	AssetBase string
	// ShowAddressFields is a boolean indicating whether the To Address and
	// From Address fields should be shown and submitted.  (It's true for
	// manual/GUI editing and false for Outpost editing.)
	ShowAddressFields bool
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
		if tag != "" {
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
	fields []field.Field
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
