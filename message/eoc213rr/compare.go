package eoc213rr

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Compare compares two messages.  It returns a score indicating how
// closely they match, and the detailed comparisons of each field in the
// message.  The comparison is not symmetric:  the receiver of the call
// is the "expected" message and the argument is the "actual" message.
func (exp *EOC213RR) Compare(actual message.Message) (score int, outOf int, fields []*message.CompareField) {
	var (
		act *EOC213RR
		ok  bool
	)
	if act, ok = actual.(*EOC213RR); !ok {
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
		common.CompareText("Incident Name", exp.IncidentName, act.IncidentName),
		common.CompareDate("Date Initiated", exp.DateInitiated, act.DateInitiated),
		common.CompareTime("Time Initiated", exp.TimeInitiated, act.TimeInitiated),
		common.CompareText("Requested By", exp.RequestedBy, act.RequestedBy),
		common.CompareText("Prepared By", exp.PreparedBy, act.PreparedBy),
		common.CompareText("Approved By", exp.ApprovedBy, act.ApprovedBy),
		common.CompareCheckbox("Approved with Signature", exp.WithSignature, act.WithSignature),
		common.CompareText("Qty/Unit", exp.QtyUnit, act.QtyUnit),
		common.CompareText("Resource Description", exp.ResourceDescription, act.ResourceDescription),
		common.CompareText("Resource Arrival", exp.ResourceArrival, act.ResourceArrival),
		common.CompareExact("Priority", exp.Priority, act.Priority),
		common.CompareText("Estimated Cost", exp.EstdCost, act.EstdCost),
		common.CompareText("Deliver To", exp.DeliverTo, act.DeliverTo),
		common.CompareText("Deliver To Location", exp.DeliverToLocation, act.DeliverToLocation),
		common.CompareText("Substitutes/Sources", exp.Substitutes, act.Substitutes),
		common.CompareCheckbox("Supplemental: Equipment Operator", exp.EquipmentOperator, act.EquipmentOperator),
		common.CompareCheckbox("Supplemental: Lodging", exp.Lodging, act.Lodging),
		common.CompareCheckbox("Supplemental: Fuel", exp.Fuel, act.Fuel),
		common.CompareText("Supplemental: Fuel Type", exp.FuelType, act.FuelType),
		common.CompareCheckbox("Supplemental: Power", exp.Power, act.Power),
		common.CompareCheckbox("Supplemental: Meals", exp.Meals, act.Meals),
		common.CompareCheckbox("Supplemental: Maintenance", exp.Maintenance, act.Maintenance),
		common.CompareCheckbox("Supplemental: Water", exp.Water, act.Water),
		common.CompareCheckbox("Supplemental: Other", exp.Other, act.Other),
		common.CompareText("Special Instructions", exp.Instructions, act.Instructions),
	)
	return common.ConsolidateCompareFields(fields)
}
