package eoc213rr

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// RenderTable renders the message as a set of field label / field value
// pairs, intended for read-only display to a human.
func (f *EOC213RR) RenderTable() (lvs []message.LabelValue) {
	var approved = f.ApprovedBy
	if f.WithSignature != "" {
		approved = common.SmartJoin(approved, "[with signature]", "\n")
	}
	var supplemental []string
	if f.EquipmentOperator != "" {
		supplemental = append(supplemental, "EquipmentOperator")
	}
	if f.Lodging != "" {
		supplemental = append(supplemental, "Lodging")
	}
	if f.Fuel != "" {
		if f.FuelType != "" {
			supplemental = append(supplemental, fmt.Sprintf("Fuel (%s)", f.FuelType))
		} else {
			supplemental = append(supplemental, "Fuel")
		}
	}
	if f.Power != "" {
		supplemental = append(supplemental, "Power")
	}
	if f.Meals != "" {
		supplemental = append(supplemental, "Meals")
	}
	if f.Maintenance != "" {
		supplemental = append(supplemental, "Maintenance")
	}
	if f.Water != "" {
		supplemental = append(supplemental, "Water")
	}
	if f.Other != "" {
		supplemental = append(supplemental, "Other")
	}
	lvs = append(f.StdFields.RenderTable1(), []message.LabelValue{
		{Label: "Incident Name", Value: f.IncidentName},
		{Label: "Date/Time Initiated", Value: common.SmartJoin(f.DateInitiated, f.TimeInitiated, " ")},
		{Label: "Tracking Number", Value: f.TrackingNumber},
		{Label: "Requested By", Value: f.RequestedBy},
		{Label: "Prepared By", Value: f.PreparedBy},
		{Label: "Approved By", Value: approved},
		{Label: "Qty/Unit", Value: f.QtyUnit},
		{Label: "Resource Description", Value: f.ResourceDescription},
		{Label: "Resource Arrival", Value: f.ResourceArrival},
		{Label: "Priority", Value: f.Priority},
		{Label: "Estimated Cost", Value: f.EstdCost},
		{Label: "Deliver To", Value: f.DeliverTo},
		{Label: "Deliver To Location", Value: f.DeliverToLocation},
		{Label: "Substitutes/Sources", Value: f.Substitutes},
		{Label: "Supplemental Requirements", Value: strings.Join(supplemental, ", ")},
		{Label: "Special Instructions", Value: f.Instructions},
	}...)
	return append(lvs, f.StdFields.RenderTable2()...)
}
