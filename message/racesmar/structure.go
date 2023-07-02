package racesmar

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for a RACES mutual aid request form.
var Type = message.Type{
	Tag:     "RACES-MAR",
	Name:    "RACES mutual aid request form",
	Article: "a",
	Create:  New,
	Decode:  decode,
}

// RACESMAR holds a RACES mutual aid request form.
type RACESMAR struct {
	common.StdFields
	AgencyName            string
	EventName             string
	EventNumber           string
	Assignment            string
	Resources             [5]Resource
	RequestedArrivalDates string
	RequestedArrivalTimes string
	NeededUntilDates      string
	NeededUntilTimes      string
	ReportingLocation     string
	ContactOnArrival      string
	TravelInfo            string
	RequestedByName       string
	RequestedByTitle      string
	RequestedByContact    string
	ApprovedByName        string
	ApprovedByTitle       string
	ApprovedByContact     string
	ApprovedByDate        string
	ApprovedByTime        string
}

// A Resource is the description of a single resource in a RACES mutual aid
// request form.
type Resource struct {
	Qty           string
	Role          string
	Position      string // Added in v2.3
	RolePos       string // Added in v2.3
	PreferredType string
	MinimumType   string
}

// Type returns the message type definition.
func (*RACESMAR) Type() *message.Type { return &Type }
