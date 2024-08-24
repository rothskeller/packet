package message

// This file contains the definition of message.Type, the registry of message
// types, and the Create and Decode calls that walk that registry.

import (
	"reflect"
)

// Type gives the details of a registered message type.
type Type struct {
	// Tag is the tag string that identifies the message type.
	Tag string
	// Name is the English name of the message type, in prose case.
	Name string
	// Article is the indefinite article to use before the Name when needed;
	// it is always either "a" or "an".
	Article string
	// PDFBase is the form-fillable PDF file whose fields will be filled in
	// to create a PDF rendering of the message.  It is nil if PDF rendering
	// is not supported.
	PDFBase []byte
	// PDFFontSize is default font size for the fillable fields in the PDF
	// file.  It can be zero if all of the fields already have assigned
	// sizes in the PDF file.
	PDFFontSize float64
	PDFRenderV2 bool
	// Create is a function for creating a new message of the type.  If
	// Create is nil, end users are not allowed to create new messages of
	// the type.  Otherwise, create must be a function with the signature
	// func() «Message», where «Message» is any type that implements the
	// Message interface.  The new message will have default values in all
	// fields.
	Create any
	// Decode is a function for decoding messages of the type.  If the input
	// message conforms to the type, Decode will return the decoded message;
	// otherwise, it will return nil.  The function must have the signature
	// func(subject string, body string) «Message», where «Message» is any
	// type that implements the Message interface.
	Decode any
}

var stringType = reflect.TypeOf("")
var messageType = reflect.TypeOf((*Message)(nil)).Elem()

// RegisteredTypes is a map from message type tag to Type structure for
// all registered message types.  It should be treated as read-only; to register
// a new type, call Register.
var RegisteredTypes = make(map[string]*Type)

// decodeFunctions is the ordered list of functions to call to try decoding a
// message.
var decodeFunctions []any

// Register registers a message type.  The order of registration is significant
// for decoding messages:  catch-all decoders (e.g. UnknownForm and PlainText)
// must be registered last.
func Register(mtype *Type) {
	var fntype reflect.Type

	// Check the validity of the decode function now, where it's easy to
	// report which function is bad if we have to.
	fntype = reflect.TypeOf(mtype.Decode)
	if fntype.Kind() != reflect.Func {
		goto BADDECODE
	}
	if fntype.NumIn() != 2 || fntype.In(0) != stringType || fntype.In(1) != stringType {
		goto BADDECODE
	}
	if fntype.NumOut() != 1 || !fntype.Out(0).Implements(messageType) {
		goto BADDECODE
	}
	// Also check the validity of the create function, if any.
	if mtype.Create != nil {
		fntype = reflect.TypeOf(mtype.Create)
		if fntype.Kind() != reflect.Func {
			goto BADCREATE
		}
		if fntype.NumIn() != 0 {
			goto BADCREATE
		}
		if fntype.NumOut() != 1 || !fntype.Out(0).Implements(messageType) {
			goto BADCREATE
		}
	}
	// Now register the type.
	RegisteredTypes[mtype.Tag] = mtype
	decodeFunctions = append(decodeFunctions, mtype.Decode)
	return
BADDECODE:
	panic("illegal decode function signature for message type " + mtype.Tag)
BADCREATE:
	panic("illegal create function signature for message type " + mtype.Tag)
}

// Create creates a new, outgoing message with the specified type tag.  It
// returns nil if the type tag is not registered.  All message fields have
// default values for an outgoing message.
func Create(tag string) Message {
	if mtype := RegisteredTypes[tag]; mtype != nil {
		if createFn := mtype.Create; createFn != nil {
			// We use reflection to build and make the call to the create
			// function.  It's not very efficient, but it allows for each
			// create function to have a clean, self-descriptive signature.
			// In particular it allows them to declare their concrete return
			// type rather than "any".
			rets := reflect.ValueOf(createFn).Call(nil)
			return rets[0].Interface().(Message)
		}
	}
	return nil
}

// Decode decodes the supplied message and returns the typed, decoded message.
// It returns nil if no registered type can decode the message.  (Between them,
// he UnknownForm and PlainText message types can decode any message, so if they
// are registered, a nil return is not possible.)
func Decode(subject string, body string) (msg Message) {
	// We use reflection to build and make the call to each decode function.
	// It's not very efficient, but it allows for each decode function to
	// have a clean, self-descriptive signature.  In particular it allows
	// them to declare their concrete return type rather than "any".
	var args = []reflect.Value{reflect.ValueOf(subject), reflect.ValueOf(body)}
	for _, fn := range decodeFunctions {
		if rets := reflect.ValueOf(fn).Call(args); !rets[0].IsNil() {
			return rets[0].Interface().(Message)
		}
	}
	return nil
}
