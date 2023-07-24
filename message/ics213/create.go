package ics213

import "time"

// New creates a new ICS-213 general message form with default values.
func New() *ICS213 {
	return &ICS213{
		FormVersion: "2.2",
		Date:        time.Now().Format("01/02/2006"),
		TxMethod:    "Other",
		OtherMethod: "Packet",
	}
}
