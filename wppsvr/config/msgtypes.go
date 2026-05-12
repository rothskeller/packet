package config

import (
	"github.com/rothskeller/packet/message"
)

// HasComputedHandlingOrder returns whether the message type with the specified
// tag can support a "computed" handling order.
func HasComputedHandlingOrder(tag string) bool {
	switch tag {
	case "NotRep", "ResReq":
		return true
	default:
		return false
	}
}

// ComputeRecommendedHandlingOrder computes the recommended handling order for a
// message.  Only message types with computed (non-static) recommended handling
// orders are handled by this function.
func ComputeRecommendedHandlingOrder(msg message.Message) string {
	switch msg.Type().Tag() {
	case "NotRep":
		for f := range msg.Fields() {
			if f.Tag() == "22." {
				switch f.Value(msg) {
				case "High":
					return "IMMEDIATE"
				case "Medium":
					return "PRIORITY"
				case "Low":
					return "ROUTINE"
				}
				break
			}
		}
	case "ResReq":
		for f := range msg.Fields() {
			if f.Tag() == "30." {
				switch f.Value(msg) {
				case "Urgent":
					return "IMMEDIATE"
				case "High", "Medium", "Low":
					return "ROUTINE"
				}
				break
			}
		}
	}
	return ""
}
