package pktmsg

import "net/mail"

// RegisteredTypes is a map from message type tag to message type name for all
// registered message types.  It should be treated as read-only; to register a
// new type, call RegisterType.
var RegisteredTypes = map[string]string{PlainTextTag: plainTextType.Name}

// createFunctions is a map from message type tag to the function that creates
// new messages of that type.
var createFunctions = map[string]func() Message{PlainTextTag: createPlainText}

// recognizeFunctions is a list of functions that can recognize a message.
var recognizeFunctions = []func(mail.Header, string) Message{
	// Note that Recognize walks this list backwards, LIFO, so that plain
	// text is always called last and unknown form second-to-last.
	recognizePlainText,
	recognizeUnknownForm,
}

// Register registers a message type.
func Register(tag, name string, createFn func() Message, recognizeFn func(mail.Header, string) Message) {
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
// whether the message is PackItForms-encoded.  If returns an error if the
// message cannot be parsed or contains no plain text body.
//
// If the raw message was just received from JNOS, bbs, callsign, and area
// should be set to indicate the BBS name, callsign used to connect to the BBS,
// and bulletin area from which the message was retrieved if any.  If the raw
// message came from some other source (e.g., read from disk), those three
// fields can be left empty.
func Recognize(raw, bbs, callsign, area string) (Message, error) {
	var header, body, err = parseRawMessage(raw, bbs, callsign, area)
	if err != nil {
		return nil, err
	}
	// Walk the list in inverse order so that UnknownForm and PlainText are
	// the last two.
	for i := len(recognizeFunctions) - 1; i >= 0; i-- {
		if msg := recognizeFunctions[i](header, body); msg != nil {
			return msg, nil
		}
	}
	panic("not reachable") // recognizePlainText accepts anything
}
