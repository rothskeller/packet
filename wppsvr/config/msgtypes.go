package config

import (
	"github.com/rothskeller/packet/xscmsg"
)

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
