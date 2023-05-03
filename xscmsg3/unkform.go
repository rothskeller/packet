package xscmsg

import (
	"github.com/rothskeller/packet/pktmsg"
)

// Unknown forms may have a tag on the subject line, but we don't want to return
// it as part of the message type.  It could be the tag of a well-known type
// that for some reason isn't registered, and code keying off of the tag might
// assume it's recognized when it isn't.
var unknownFormType = MessageType{"UNKNOWN", "unrecognized form", "an"}

type unknownForm struct {
	Message
}

func recognizeUnknownForm(pm pktmsg.Message) Message {
	if pm.KeyedField(pktmsg.FPIFOVersion) == nil {
		return nil // not a PackItForms form
	}
	return unknownForm{NewBaseMessage(unknownFormType, pm)}
}

func (m unknownForm) Iterate(fn func(pktmsg.Field)) {
	m.Message.Iterate(func(f pktmsg.Field) {
		switch f := f.(type) {
		case pktmsg.KeyedField:
			if f.Key() != pktmsg.FBody {
				fn(f)
			}
		case pktmsg.TaggedField:
			fn(unknownTaggedField{f}) // make it Viewable
		}
	})
}

func (m unknownForm) KeyedField(key pktmsg.FieldKey) pktmsg.KeyedField {
	if key == pktmsg.FBody {
		return nil
	}
	return m.Message.KeyedField(key)
}

func (m unknownForm) TaggedField(tag string) pktmsg.TaggedField {
	if f := m.Message.TaggedField(tag); f != nil {
		return unknownTaggedField{f} // make it Viewable
	}
	return nil
}

type unknownTaggedField struct{ pktmsg.TaggedField }

func (f unknownTaggedField) Label() string { return f.Tag() }
