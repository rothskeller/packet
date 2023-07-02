package unkform

import "github.com/rothskeller/packet/message"

// Type is the type definition for an unrecognized form message.
var Type = message.Type{
	Tag:     "UNKNOWN",
	Name:    "unrecognized form message",
	Article: "an",
	Create:  nil,
	Decode:  decode,
}

// UnknownForm holds the details of an unrecognized form message.
type UnknownForm struct {
	OriginMsgID  string
	Handling     string
	FormTag      string
	Subject      string
	FormHTML     string
	FormVersion  string
	TaggedValues map[string]string
}

// Type returns the message type definition.
func (*UnknownForm) Type() *message.Type { return &Type }
