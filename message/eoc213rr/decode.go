package eoc213rr

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

func decode(subject, body string) (f *EOC213RR) {
	if idx := strings.Index(body, "form-scco-eoc-213rr.html"); idx < 0 {
		return nil
	}
	form := common.DecodePIFO(body)
	if form == nil || form.HTMLIdent != "form-scco-eoc-213rr.html" {
		return nil
	}
	switch form.FormVersion {
	case "2.0", "2.1", "2.3":
		break
	default:
		return nil
	}
	f = new(EOC213RR)
	f.PIFOVersion = form.PIFOVersion
	f.FormVersion = form.FormVersion
	f.StdFields.Decode(form.TaggedValues)
	f.IncidentName = form.TaggedValues["21."]
	f.DateInitiated = form.TaggedValues["22."]
	f.TimeInitiated = form.TaggedValues["23."]
	f.TrackingNumber = form.TaggedValues["24."]
	f.RequestedBy = form.TaggedValues["25."]
	f.PreparedBy = form.TaggedValues["26."]
	f.ApprovedBy = form.TaggedValues["27."]
	f.QtyUnit = form.TaggedValues["28."]
	f.ResourceDescription = form.TaggedValues["29."]
	f.ResourceArrival = form.TaggedValues["30."]
	f.Priority = form.TaggedValues["31."]
	f.EstdCost = form.TaggedValues["32."]
	f.DeliverTo = form.TaggedValues["33."]
	f.DeliverToLocation = form.TaggedValues["34."]
	f.Substitutes = form.TaggedValues["35."]
	f.EquipmentOperator = form.TaggedValues["36a."]
	f.Lodging = form.TaggedValues["36b."]
	f.Fuel = form.TaggedValues["36c."]
	f.FuelType = form.TaggedValues["36d."]
	f.Power = form.TaggedValues["36e."]
	f.Meals = form.TaggedValues["36f."]
	f.Maintenance = form.TaggedValues["36g."]
	f.Water = form.TaggedValues["36h."]
	f.Other = form.TaggedValues["36i."]
	f.Instructions = form.TaggedValues["37."]
	return f
}
