package sheltstat

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Compare compares two messages.  It returns a score indicating how
// closely they match, and the detailed comparisons of each field in the
// message.  The comparison is not symmetric:  the receiver of the call
// is the "expected" message and the argument is the "actual" message.
func (exp *SheltStat) Compare(actual message.Message) (score int, outOf int, fields []*message.CompareField) {
	var (
		act *SheltStat
		ok  bool
	)
	if act, ok = actual.(*SheltStat); !ok {
		return 0, 1, []*message.CompareField{{
			Label: "Message Type",
			Score: 0, OutOf: 1,
			Expected:     Type.Name,
			ExpectedMask: strings.Repeat("*", len(Type.Name)),
			Actual:       actual.Type().Name,
			ActualMask:   strings.Repeat("*", len(actual.Type().Name)),
		}}
	}
	fields = append(exp.StdFields.Compare(&act.StdFields),
		common.CompareExact("Report Type", exp.ReportType, act.ReportType),
		common.CompareText("Shelter Name", exp.ShelterName, act.ShelterName),
		common.CompareExact("Shelter Type", exp.ShelterType, act.ShelterType),
		common.CompareExact("Shelter Status", exp.ShelterStatus, act.ShelterStatus),
		common.CompareText("Shelter Address", exp.ShelterAddress, act.ShelterAddress),
		common.CompareText("Shelter City", exp.ShelterCity, act.ShelterCity),
		common.CompareText("Shelter State", exp.ShelterState, act.ShelterState),
		common.CompareText("Shelter Zip", exp.ShelterZip, act.ShelterZip),
		common.CompareReal("Latitude", exp.Latitude, act.Latitude),
		common.CompareReal("Longitude", exp.Longitude, act.Longitude),
		common.CompareCardinal("Capacity", exp.Capacity, act.Capacity),
		common.CompareCardinal("Occupancy", exp.Occupancy, act.Occupancy),
		common.CompareCardinal("Meals Served", exp.MealsServed, act.MealsServed),
		common.CompareText("NSS Number", exp.NSSNumber, act.NSSNumber),
		common.CompareExactMap("Pet-Friendly", exp.PetFriendly, act.PetFriendly, yesNoMap),
		common.CompareExactMap("Basic Safety Inspection", exp.BasicSafetyInspection, act.BasicSafetyInspection, yesNoMap),
		common.CompareExactMap("ATC-20 Inspection", exp.ATC20Inspection, act.ATC20Inspection, yesNoMap),
		common.CompareText("Available Services", exp.AvailableServices, act.AvailableServices),
		common.CompareText("MOU", exp.MOU, act.MOU),
		common.CompareText("Floor Plan", exp.FloorPlan, act.FloorPlan),
		common.CompareExact("Managed By", exp.ManagedBy, act.ManagedBy),
		common.CompareText("Managed By Detail", exp.ManagedByDetail, act.ManagedByDetail),
		common.CompareText("Primary Contact", exp.PrimaryContact, act.PrimaryContact),
		common.CompareExact("Primary Phone", exp.PrimaryPhone, act.PrimaryPhone),
		common.CompareText("Secondary Contact", exp.SecondaryContact, act.SecondaryContact),
		common.CompareExact("Secondary Phone", exp.SecondaryPhone, act.SecondaryPhone),
		common.CompareExact("Tactical Call Sign", exp.TacticalCallSign, act.TacticalCallSign),
		common.CompareExact("Repeater Call Sign", exp.RepeaterCallSign, act.RepeaterCallSign),
		common.CompareExact("Repeater Input", exp.RepeaterInput, act.RepeaterInput),
		common.CompareExact("Repeater Input Tone", exp.RepeaterInputTone, act.RepeaterInputTone),
		common.CompareExact("Repeater Output", exp.RepeaterOutput, act.RepeaterOutput),
		common.CompareExact("Repeater Output Tone", exp.RepeaterOutputTone, act.RepeaterOutputTone),
		common.CompareExact("Repeater Offset", exp.RepeaterOffset, act.RepeaterOffset),
		common.CompareText("Comments", exp.Comments, act.Comments),
		common.CompareExactMap("Remove from List", exp.RemoveFromList, act.RemoveFromList, yesNoMap),
	)
	return common.ConsolidateCompareFields(fields)
}

var yesNoMap = map[string]string{
	"":        "(not set)",
	"false":   "No",
	"checked": "Yes",
}
