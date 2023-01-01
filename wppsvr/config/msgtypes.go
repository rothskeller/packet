package config

import (
	"github.com/rothskeller/packet/xscmsg"
)

// validMessageTypes contains an empty message of each type that can be used
// for packet practice.
var validMessageTypes []*xscmsg.Message

// ValidMessageTypes returns a slice containing an empty message of each type
// that can be used for packet practice.
func ValidMessageTypes() []*xscmsg.Message {
	if validMessageTypes == nil {
		// Initialize this on first use rather than at init time, so we
		// can be sure that all of the message types have been
		// registered.
		validMessageTypes = []*xscmsg.Message{
			xscmsg.CreatePlainTextMessage(),
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
func LookupMessageType(tag string) *xscmsg.Message {
	for _, msg := range ValidMessageTypes() {
		if msg.Type.Tag == tag {
			return msg
		}
	}
	return nil
}

// ComputedRecommendedHandlingOrder is a map from message type tags to functions
// that compute the recommended handling order for messages of that type.  Only
// message types with computed (non-static) recommended handling orders have
// entries in this map.
var ComputedRecommendedHandlingOrder = map[string](func(*xscmsg.Message) xscmsg.HandlingOrder){
	"ICS213": func(msg *xscmsg.Message) xscmsg.HandlingOrder {
		var sev xscmsg.MessageSeverity
		if f := msg.Field("4."); f != nil {
			sev, _ = xscmsg.ParseSeverity(f.Value)
		}
		switch sev {
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
	"EOC213RR": func(msg *xscmsg.Message) xscmsg.HandlingOrder {
		switch msg.Field("31.").Value {
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
