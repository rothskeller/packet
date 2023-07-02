package message

import (
	"reflect"
)

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
		switch fntype.NumIn() {
		case 0:
			break
		case 2:
			if fntype.In(0) != stringType || fntype.In(1) != stringType {
				goto BADCREATE
			}
		default:
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
// returns nil if the type tag is not registered.  The operator call sign and
// name are filled into the message if it supports them.  All other message
// fields have default values for an outgoing message.
func Create(tag, opcall, opname string) Message {
	if createFn := RegisteredTypes[tag].Create; createFn != nil {
		// We use reflection to build and make the call to the create
		// function.  It's not very efficient, but it allows for each
		// create function to have a clean, self-descriptive signature.
		// In particular it allows them to declare their concrete return
		// type rather than "any".
		var args []reflect.Value
		if reflect.TypeOf(createFn).NumIn() == 2 {
			args = []reflect.Value{reflect.ValueOf(opcall), reflect.ValueOf(opname)}
		}
		rets := reflect.ValueOf(createFn).Call(args)
		return rets[0].Interface().(Message)
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
