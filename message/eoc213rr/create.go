package eoc213rr

import (
	"time"

	"github.com/rothskeller/packet/message/common"
)

// New creates a new EOC-213RR resource request form with default values.
func New(opcall, opname string) *EOC213RR {
	return &EOC213RR{
		StdFields: common.StdFields{
			FormVersion:   "2.3",
			MessageDate:   time.Now().Format("01/02/2006"),
			ToICSPosition: "Planning Section",
			ToLocation:    "County EOC",
			OpCall:        opcall,
			OpName:        opname,
		},
		DateInitiated: time.Now().Format("01/02/2006"),
	}
}
