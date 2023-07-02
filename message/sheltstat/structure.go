package sheltstat

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an OA shelter status form.
var Type = message.Type{
	Tag:     "SheltStat",
	Name:    "OA shelter status form",
	Article: "an",
	Create:  New,
	Decode:  decode,
}

// SheltStat holds an OA shelter status form.
type SheltStat struct {
	common.StdFields
	ReportType            string
	ShelterName           string
	ShelterType           string
	ShelterStatus         string
	ShelterAddress        string
	ShelterCityCode       string // added in v2.2
	ShelterCity           string
	ShelterState          string
	ShelterZip            string
	Latitude              string
	Longitude             string
	Capacity              string
	Occupancy             string
	MealsServed           string
	NSSNumber             string
	PetFriendly           string
	BasicSafetyInspection string
	ATC20Inspection       string
	AvailableServices     string
	MOU                   string
	FloorPlan             string
	ManagedByCode         string // added in v2.2
	ManagedBy             string
	ManagedByDetail       string
	PrimaryContact        string
	PrimaryPhone          string
	SecondaryContact      string
	SecondaryPhone        string
	TacticalCallSign      string
	RepeaterCallSign      string
	RepeaterInput         string
	RepeaterInputTone     string
	RepeaterOutput        string
	RepeaterOutputTone    string
	RepeaterOffset        string
	Comments              string
	RemoveFromList        string
}

// Type returns the message type definition.
func (*SheltStat) Type() *message.Type { return &Type }
