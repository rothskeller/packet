// Package xscmsg handles recognition, parsing, validation, and encoding of all
// of the Santa Clara County (XSC) standard message types.  Each message type is
// represented by its own concrete type; all of them implement the Message
// interface.
package xscmsg

import "github.com/rothskeller/packet/pktmsg"

// Message is the interface implemented by all XSC messages.
type Message interface {
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
	EncodedBody() string
}

// RegisteredTypes is a map from message type tag to message type name for all
// registered message types.  It should be treated as read-only; to register a
// new type, call RegisterType.
var RegisteredTypes = map[string]string{PlainTextTag: plainTextName}

// createFunctions is a map from message type tag to the function that creates
// new messages of that type.
var createFunctions = map[string]func() Message{PlainTextTag: createPlainText}

// recognizeFunctions is a list of functions that can recognize a message.
var recognizeFunctions = []func(*pktmsg.Message, *pktmsg.Form) Message{
	// Note that Recognize walks this list backwards, LIFO, so that plain
	// text is always called last and unknown form second-to-last.
	recognizePlainText,
	recognizeUnknownForm,
}

// Register registers a message type.
func Register(tag, name string, createFn func() Message, recognizeFn func(*pktmsg.Message, *pktmsg.Form) Message) {
	RegisteredTypes[tag] = name
	if createFn != nil {
		createFunctions[tag] = createFn
	}
	if recognizeFn != nil {
		recognizeFunctions = append(recognizeFunctions, recognizeFn)
	}
}

// Create creates a new message of the specified type.  It returns nil if the
// specified tag is not registered.  Create fills in all fields with their
// default values.
func Create(tag string) Message {
	if createFn := createFunctions[tag]; createFn != nil {
		return createFn()
	}
	return nil
}

// Recognize recognizes the provided raw message as one of the registered types,
// and returns the corresponding Message.  If the message type is not
// registered, it will return an UnknownForm or a PlainText, depending on
// whether the message is PackItForms-encoded.
func Recognize(raw *pktmsg.Message) Message {
	var form = pktmsg.ParseForm(raw.Body)
	// Walk the list in inverse order so that UnknownForm and PlainText are
	// the last two.
	for i := len(recognizeFunctions) - 1; i >= 0; i-- {
		if msg := recognizeFunctions[i](raw, form); msg != nil {
			return msg
		}
	}
	panic("not reachable") // recognizePlainText accepts anything
}

/******************************************************************************

// OMessage is the interface implemented by all XSC messages.
type OMessage struct {
	// Type is the message type definition.
	Type *MessageType
	// RawMessage is the underlying raw pktmsg.Message, if any.
	RawMessage *pktmsg.Message
	// RawForm is the underlying raw pktmsg.Form, if any.
	RawForm *pktmsg.Form
	// Fields returns the list of fields in the message.
	Fields []*OField
}

// Field returns the field with the specified tag, or nil if there is
// no such field.
func (m *OMessage) Field(tag string) *OField {
	for _, f := range m.Fields {
		if f.Def.Tag == tag {
			return f
		}
	}
	return nil
}

// KeyField returns the field with the specified Key, or nil if there is no
// such field.
func (m *OMessage) KeyField(key FieldKey) *OField {
	for _, f := range m.Fields {
		if f.Def.KeyFunc != nil {
			if f.Def.KeyFunc(f, m) == key {
				return f
			}
		} else if f.Def.Key == key {
			return f
		}
	}
	return nil
}

// Validate ensures that the contents of the message are correct.  It returns a
// list of problems, which is empty if the message is fine.  If strict is true,
// the message must be exactly correct; otherwise, some trivial issues are
// corrected and not reported.
func (m *OMessage) Validate(strict bool) (problems []string) {
	for _, f := range m.Fields {
		if prob := f.Validate(m, strict); prob != "" {
			problems = append(problems, prob)
		}
	}
	return problems
}

// Subject returns the encoded subject of the message.
func (m *OMessage) Subject() string {
	if m.Type.SubjectFunc != nil {
		return m.Type.SubjectFunc(m)
	}
	if f := m.KeyField(FSubject); f != nil {
		return f.Value
	}
	return ""
}

// Body returns the encoded body of the message.  If human is
// true, it is encoded for human reading or editing; if false, it is
// encoded for transmission.
func (m *OMessage) Body() string {
	if m.Type.BodyFunc != nil {
		return m.Type.BodyFunc(m)
	}
	if f := m.KeyField(FBody); f != nil {
		return f.Value
	}
	return ""
}

// SetDestinationMessageNumber sets the destination message number for a message
// (i.e., the local message number for a received message).  It handles the
// special cases for old versions of ICS-213 forms.  It is a no-op for message
// types that don't have a destination message number field.
func (m *OMessage) SetDestinationMessageNumber(msgno string) {
	if f := m.KeyField(FOldICS213TxMsgNo); f != nil {
		f.Value = m.KeyField(FOriginMsgNo).Value
	}
	if f := m.KeyField(FDestinationMsgNo); f != nil {
		f.Value = msgno
	}
}

// MessageType defines a message type.
type MessageType struct {
	// Tag is the string used to identify the message type.  For form
	// messages, this is the form tag that appears encoded in the subject
	// line of the message.
	Tag string
	// Name is the English name of the message type.  It is a noun phrase in
	// prose case, such as "foo message" or "bar form".
	Name string
	// Article is the indefinite article ("a" or "an") to be used preceding
	// the Name, in a sentence that needs one.
	Article string
	// HTML is the HTML file listed in the PackItForms header for the
	// message type, if any.
	HTML string
	// Version is the version number of the message type, if any.
	Version string
	// SubjectFunc is a function that returns the encoded subject of the
	// message.  If it is nil, the subject is taken from the "SUBJECT"
	// message field.
	SubjectFunc func(msg *OMessage) string
	// BodyFunc is a function that returns the encoded body of the message.
	// If it is nil, the body is taken from the "BODY" message field.
	BodyFunc func(msg *OMessage) string
}

// OField is a single field within an XSC message.
type OField struct {
	// Def is the definition of the field.
	Def *FieldDef
	// Value is the value of the field.
	Value string
}

// Validate verifies the correctness of the field value, returning an error
// message if it is invalid or an empty string if it is valid.  If strict is
// true, correctable problems are reported; otherwise they are corrected.  For
// fields whose values are computed based on the values of other fields,
// Validate performs the computation.
func (f *OField) Validate(msg *OMessage, strict bool) (prob string) {
	if f.Def.Flags&Required != 0 && f.Value == "" {
		prob = fmt.Sprintf("The %q field must have a value.", f.Def.Tag)
		// don't return it yet
	}
	if prob == "" && f.Def.Flags&RequiredForComplete != 0 && f.Value == "" {
		if c := msg.KeyField(FComplete); c != nil && c.Value == "Complete" {
			prob = fmt.Sprintf("The %q field must have a value when the %q field is Complete.", f.Def.Tag, c.Def.Tag)
			// don't return it yet
		}
	}
	for _, vfn := range f.Def.Validators {
		if prob2 := vfn(f, msg, strict); prob2 != "" {
			if prob != "" {
				return prob // prefer the requirement problem
			}
			return prob2
		}
	}
	// Return the deferred requirement problem only if none of the
	// validator functions added a value.
	if prob != "" && f.Value == "" {
		return prob
	}
	return ""
}

// Default returns the default value of the field.
func (f *OField) Default() string {
	if f.Def.DefaultFunc != nil {
		return f.Def.DefaultFunc()
	}
	return f.Def.DefaultValue
}

// FieldDef contains the definition of a OField within a OMessage.
type FieldDef struct {
	// Tag is the unique identifier of the field within its message.  For
	// most form fields, this is a field number followed by a period.
	Tag string
	// KeyFunc is a function that returns the well-known-field key for this
	// field, if any.
	KeyFunc func(f *OField, m *OMessage) FieldKey
	// Key is the well-known-field key for this field, if KeyFunc is not
	// set.
	Key FieldKey
	// Label is the English label of the field.  For form fields, it is the
	// label of the field as it appears on the PDF form.
	Label string
	// Comment is a comment displayed in the human-readable rendering of a
	// form when the field has no assigned value.  It is generally a textual
	// reminder of the validation requirements for the field value.
	Comment string
	// Flags is a set of flags describing the field.
	Flags FieldFlag
	// DefaultFunc is a function to compute the default value of the field.
	DefaultFunc func() string
	// DefaultValue is the default value of the field, if DefaultFunc is not
	// set.
	DefaultValue string
	// Validators is the list of functions called to validate the field
	// value.  See the comment on OField.Validate() for details.
	Validators []Validator
	// Choices is the set of allowed values for a field with a restricted
	// set.
	Choices []string
}

// FieldFlag is a flag, or bitmask of flags, describing a field.
type FieldFlag uint

// Values for FieldFlag:
const (
	// Readonly indicates that the field is read-only.  It may have a value
	// when a message is received and decoded, or it may have one computed
	// by its Validate method based on other message fields, but it should
	// not be presented to the user for editing.
	Readonly FieldFlag = 1 << iota
	// Multiline indicates that the field is allowed to have a multi-line
	// value.
	Multiline
	// Required indicates that a value for the field is required.
	Required
	// RequiredForComplete indicates that a value for the field is required
	// if the FComplete field contains "Complete".
	RequiredForComplete
)

// Validator is a function called to validate a field value.  See the
// comment on OField.Validate() for details.
type Validator func(f *OField, msg *OMessage, strict bool) (problem string)

// RegisteredTypes is a map from type tag to type definition for all registered
// types.  Callers should treat it as read-only; to register a type, use
// RegisterCreate and RegisterType.
var RegisteredTypes = map[string]*MessageType{
	PlainTextTag: &plainTextMessageType, // plain text is always registered
}

var createFuncs = map[string]func() *OMessage{
	PlainTextTag: CreatePlainTextMessage, // always registered for create
}
var recognizeFuncs []func(*pktmsg.Message, *pktmsg.Form) *OMessage

// RegisterCreate registers a function to create a message with the specified
// type.  (It is normally called by an init function in the package that
// implements the message type.  To make a message type usable, simply import
// its implementing package.)
func RegisterCreate(mtype *MessageType, fn func() *OMessage) {
	RegisteredTypes[mtype.Tag] = mtype
	createFuncs[mtype.Tag] = fn
}

// RegisterType registers a message type to be recognizable by the Recognize
// function.  (It is normally called by an init function in the package that
// implements the message type.  To make a message type recognizable, simply
// import its implementing package.)
func RegisterType(recognize func(*pktmsg.Message, *pktmsg.Form) *OMessage) {
	recognizeFuncs = append(recognizeFuncs, recognize)
}

// Create creates a new message of the type identified by the supplied tag.  If
// the tag is not recognized, Create returns nil.
//
// For this to work, the message type you want to create must have been
// registered.  The standard message types can be registered by importing the
// appropriate message-type-specific subpackages of xscmsg.  Alternatively, all
// standard message types can be registered by importing xscmsg/all.
func Create(tag string) *OMessage {
	if fn := createFuncs[tag]; fn != nil {
		return fn()
	}
	return nil
}

// Recognize determines which registered XSC message type to use for the
// supplied pktmsg.Message, and returns the corresponding xscmsg.OMessage.
// Recognize always returns an xscmsg.OMessage; if the supplied message isn't
// recognized as any more specific type, it is returned as an "unknown form"
// message or a "plain text" message.
//
// For a message to be recognized as a specific type, that type must have been
// registered.  The standard message types can be registered by importing the
// appropriate message-type-specific subpackages of xscmsg.  Alternatively, all
// standard message types can be registered by importing xscmsg/all.
func Recognize(msg *pktmsg.Message) *OMessage {
	var form *pktmsg.Form
	if pktmsg.IsForm(msg.Body) {
		form = pktmsg.ParseForm(msg.Body)
	}
	for _, recognizeFunc := range recognizeFuncs {
		if xsc := recognizeFunc(msg, form); xsc != nil {
			return xsc
		}
	}
	if form != nil {
		return makeUnknownFormMessage(msg, form)
	}
	return makePlainTextMessage(msg)
}

// OlderVersion compares two version numbers, and returns true if the first one
// is older than the second one.  Each version number is a dot-separated
// sequence of parts, each of which is compared independently with the
// corresponding part in the other version number.  The parts are compared
// numerically if they parse as integers, and as strings otherwise.  However, a
// part that starts with a digit is always "newer" than a part that doesn't.
// (This results in "undefined" being treated as infinitely old, which is what
// we want.)
func OlderVersion(a, b string) bool {
	aparts := strings.Split(a, ".")
	bparts := strings.Split(b, ".")
	for len(aparts) != 0 && len(bparts) != 0 {
		anum, aerr := strconv.Atoi(aparts[0])
		bnum, berr := strconv.Atoi(bparts[0])
		if aerr == nil && berr == nil {
			if anum < bnum {
				return true
			}
			if anum > bnum {
				return false
			}
		} else if startsWithDigit(aparts[0]) && !startsWithDigit(bparts[0]) {
			return false
		} else if !startsWithDigit(aparts[0]) && startsWithDigit(bparts[0]) {
			return true
		} else {
			if aparts[0] < bparts[0] {
				return true
			}
			if aparts[0] > bparts[0] {
				return false
			}
		}
		aparts = aparts[1:]
		bparts = bparts[1:]
	}
	if len(bparts) != 0 {
		return true
	}
	return false
}
func startsWithDigit(s string) bool {
	return s != "" && s[0] >= '0' && s[0] <= '9'
}
*/
