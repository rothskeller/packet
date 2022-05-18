package config

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/eoc213rr"
	"github.com/rothskeller/packet/xscmsg/ics213"
)

// validMessageTypes contains an empty message of each type that can be used
// for packet practice.
var validMessageTypes []xscmsg.XSCMessage

// ValidMessageTypes returns a slice containing an empty message of each type
// that can be used for packet practice.
func ValidMessageTypes() []xscmsg.XSCMessage {
	if validMessageTypes == nil {
		// Initialize this on first use rather than at init time, so we
		// can be sure that all of the message types have been
		// registered.
		validMessageTypes = []xscmsg.XSCMessage{
			new(PlainTextMessage),
			xscmsg.Create("AHFacStat"),
			xscmsg.Create("EOC213RR"),
			xscmsg.Create("ICS213"),
			xscmsg.Create("JurisStat"),
			xscmsg.Create("MuniStat"),
			xscmsg.Create("RACES-MAR"),
			xscmsg.Create("SheltStat"),
		}
	}
	return validMessageTypes
}

// LookupMessageType finds the message type with the specified tag, if it
// exists in ValidMessageTypes.
func LookupMessageType(tag string) xscmsg.XSCMessage {
	for _, msg := range ValidMessageTypes() {
		if msg.TypeTag() == tag {
			return msg
		}
	}
	return nil
}

// ComputedRecommendedHandlingOrder is a map from message type tags to functions
// that compute the recommended handling order for messages of that type.  Only
// message types with computed (non-static) recommended handling orders have
// entries in this map.
var ComputedRecommendedHandlingOrder = map[string](func(xscmsg.XSCMessage) xscmsg.HandlingOrder){
	"ICS213": func(msg xscmsg.XSCMessage) xscmsg.HandlingOrder {
		severity := msg.(*ics213.Form).Form().Get("4.")
		switch sev, _ := xscmsg.ParseSeverity(severity); sev {
		case xscmsg.SeverityEmergency:
			return xscmsg.HandlingImmediate
		case xscmsg.SeverityUrgent:
			return xscmsg.HandlingPriority
		case xscmsg.SeverityOther:
			return xscmsg.HandlingRoutine
		default:
			return 0
		}
	},
	"EOC213RR": func(msg xscmsg.XSCMessage) xscmsg.HandlingOrder {
		priority := msg.(*eoc213rr.Form).Form().Get("31.")
		switch priority {
		case "Now", "High":
			return xscmsg.HandlingImmediate
		case "Medium":
			return xscmsg.HandlingPriority
		case "Low":
			return xscmsg.HandlingRoutine
		default:
			return 0
		}
	},
}

// A PlainTextMessage is an XSCMessage implementation for plain text messages.
// They aren't a true XSC message type, but for the purposes of packet checkins
// it's convenient to treat them as one.
type PlainTextMessage struct {
	M *pktmsg.Message
}

// TypeTag returns the tag string used to identify the message type.
func (m *PlainTextMessage) TypeTag() string { return "plain" }

// TypeName returns the English name of the message type.  It is a noun
// phrase in prose case, such as "foo message" or "bar form".
func (m *PlainTextMessage) TypeName() string { return "plain text message" }

// TypeArticle returns the indefinite article ("a" or "an") to be used
// preceding the TypeName, in a sentence that needs one.
func (m *PlainTextMessage) TypeArticle() string { return "a" }

// Validate ensures that the contents of the message are correct.  It
// returns a list of problems, which is empty if the message is fine.
// If strict is true, the message must be exactly correct; otherwise,
// some trivial issues are corrected and not reported.
func (m *PlainTextMessage) Validate(strict bool) (problems []string) { return nil }

// Message returns the encoded message.  If human is true, it is encoded
// for human reading or editing; if false, it is encoded for
// transmission.  If the XSCMessage was originally created by a call to
// Recognize, the Message structure passed to it is updated and reused;
// otherwise, a new Message structure is created and filled in.
func (m *PlainTextMessage) Message(human bool) *pktmsg.Message { return m.M }

// An UnknownForm is an XSCMessage implementation for messages that contain a
// form of a type we don't recognize.  This isn't a true XSC message type, but
// for the purposes of packet checkins it's convenient to treat it as one.
type UnknownForm struct {
	M *pktmsg.Message
	F *pktmsg.Form
}

// TypeTag returns the tag string used to identify the message type.
func (m *UnknownForm) TypeTag() string { return "UNKNOWN" }

// TypeName returns the English name of the message type.  It is a noun
// phrase in prose case, such as "foo message" or "bar form".
func (m *UnknownForm) TypeName() string { return "form of unknown type" }

// TypeArticle returns the indefinite article ("a" or "an") to be used
// preceding the TypeName, in a sentence that needs one.
func (m *UnknownForm) TypeArticle() string { return "a" }

// Validate ensures that the contents of the message are correct.  It
// returns a list of problems, which is empty if the message is fine.
// If strict is true, the message must be exactly correct; otherwise,
// some trivial issues are corrected and not reported.
func (m *UnknownForm) Validate(strict bool) (problems []string) { return nil }

// Message returns the encoded message.  If human is true, it is encoded
// for human reading or editing; if false, it is encoded for
// transmission.  If the XSCMessage was originally created by a call to
// Recognize, the Message structure passed to it is updated and reused;
// otherwise, a new Message structure is created and filled in.
func (m *UnknownForm) Message(human bool) *pktmsg.Message { return m.M }
