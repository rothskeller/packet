package typedmsg

import "github.com/rothskeller/packet/pktmsg"

// RegisteredTypes is a map from message type tag to MessageType structure for
// all registered message types.  It should be treated as read-only; to register
// a new type, call RegisterType.
var RegisteredTypes = make(map[string]*MessageType)

// recognizeFunctions is a list of functions that can recognize a message.
var recognizeFunctions []func(*pktmsg.Message) Message

// Register registers a message type.
func Register(mtype *MessageType) {
	RegisteredTypes[mtype.Tag] = mtype
	if mtype.Recognize != nil {
		recognizeFunctions = append(recognizeFunctions, mtype.Recognize)
	}
}

// Create creates a new, outgoing message with the specified type.  It returns
// nil if the type tag is not registered.
func Create(tag string) Message {
	if createFn := RegisteredTypes[tag].Create; createFn != nil {
		return createFn()
	}
	return nil
}

// Recognize returns the typed message corresponding to the supplied packet
// message, or nil, if no registered type recognizes it.
func Recognize(msg *pktmsg.Message) (tm Message) {
	for _, fn := range recognizeFunctions {
		if tm = fn(msg); tm != nil {
			break
		}
	}
	return tm
}
