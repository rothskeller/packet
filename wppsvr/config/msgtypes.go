package config

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/eoc213rr"
	"github.com/rothskeller/packet/message/ics213"
)

// ComputeRecommendedHandlingOrder computes the recommended handling order for a
// message.  Only message types with computed (non-static) recommended handling
// orders are handled by this function.
func ComputeRecommendedHandlingOrder(msg message.Message) string {
	switch msg := msg.(type) {
	case *ics213.ICS213:
		switch msg.Severity {
		case "EMERGENCY":
			return "IMMEDIATE"
		case "URGENT":
			return "PRIORITY"
		case "OTHER":
			return "ROUTINE"
		}
	case *eoc213rr.EOC213RR:
		switch msg.Priority {
		case "Now", "High":
			return "IMMEDIATE"
		case "Medium":
			return "PRIORITY"
		case "Low":
			return "ROUTINE"
		}
	}
	return ""
}
