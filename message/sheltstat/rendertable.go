package sheltstat

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// RenderTable renders the message as a set of field label / field value
// pairs, intended for read-only display to a human.
func (f *SheltStat) RenderTable() (lvs []message.LabelValue) {
	lvs = append(f.StdFields.RenderTable1(), []message.LabelValue{
		{Label: "Report Type", Value: f.ReportType},
		{Label: "Shelter Name", Value: f.ShelterName},
		{Label: "Shelter Type", Value: f.ShelterType},
		{Label: "Shelter Status", Value: f.ShelterStatus},
		{Label: "Shelter Address", Value: f.ShelterAddress},
		{Label: "Shelter City", Value: common.SmartJoin(f.ShelterCity, common.SmartJoin(f.ShelterState, f.ShelterZip, "  "), ", ")},
		{Label: "GPS Coordinates", Value: common.SmartJoin(f.Latitude, f.Longitude, ", ")},
		{Label: "Occupancy", Value: common.SmartJoin(f.Occupancy, f.Capacity, " out of ")},
		{Label: "Meals Served", Value: f.MealsServed},
		{Label: "NSS Number", Value: f.NSSNumber},
		{Label: "Pet Friendly", Value: f.PetFriendly},
		{Label: "Basic Safety Inspection", Value: f.BasicSafetyInspection},
		{Label: "ATC-20 Inspection", Value: f.ATC20Inspection},
		{Label: "Available Services", Value: f.AvailableServices},
		{Label: "MOU", Value: f.MOU},
		{Label: "Floor Plan", Value: f.FloorPlan},
		{Label: "Managed By", Value: common.SmartJoin(f.ManagedBy, f.ManagedByDetail, ": ")},
		{Label: "Primary Contact", Value: f.PrimaryContact},
		{Label: "Primary Phone", Value: f.PrimaryPhone},
		{Label: "Secondary Contact", Value: f.SecondaryContact},
		{Label: "Secondary Phone", Value: f.SecondaryPhone},
		{Label: "Tactical Call Sign", Value: f.TacticalCallSign},
		{Label: "Repeater Call Sign", Value: f.RepeaterCallSign},
		{Label: "Repeater Input", Value: f.RepeaterInput},
		{Label: "Repeater Input Tone", Value: f.RepeaterInputTone},
		{Label: "Repeater Output", Value: f.RepeaterOutput},
		{Label: "Repeater Output Tone", Value: f.RepeaterOutputTone},
		{Label: "Repeater Offset", Value: f.RepeaterOffset},
		{Label: "Comments", Value: f.Comments},
		{Label: "Remove from List", Value: f.RemoveFromList},
	}...)
	return append(lvs, f.StdFields.RenderTable2()...)
}
