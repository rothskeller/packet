// Package typedmsg builds on package pktmsg to allow for the definition of
// specific types of messages that may have different behaviors.  Message types
// are identified by a tag string.
package typedmsg

import "github.com/rothskeller/packet/pktmsg"

// MessageType gives the details of a registered message type.
type MessageType struct {
	// Tag is the tag string that identifies the message type.
	Tag string
	// Name is the English name of the message type, in prose case.
	Name string
	// Article is the indefinite article to use before the Name when needed;
	// it is always either "a" or "an".
	Article string
	// Create is a function for creating a new message of the type.  If
	// Create is nil, end users are not allowed to create new messages of
	// the type.  The new message will have default values in all fields.
	Create func() Message
	// Recognize is a function for recognizing existing messages of the
	// type.  If the supplied message belongs to this type, Recognize will
	// return a Message for it; otherwise, it will return nil.
	Recognize func(*pktmsg.Message) Message
}

// Message is the interface supported by all typed messages.
type Message interface {
	pktmsg.IMessage
	Type() *MessageType
}
