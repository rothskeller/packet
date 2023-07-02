package eoc213rr

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (f *EOC213RR) EncodeSubject() string {
	return common.EncodeSubject(f.OriginMsgID, f.Handling, Type.Tag, f.IncidentName)
}

// EncodeBody encodes the message body.
func (f *EOC213RR) EncodeBody() string {
	var (
		sb  strings.Builder
		enc *common.PIFOEncoder
	)
	if f.FormVersion == "" {
		f.FormVersion = "2.3"
	}
	enc = common.NewPIFOEncoder(&sb, "form-scco-eoc-213rr.html", f.FormVersion)
	f.StdFields.EncodeHeader(enc)
	enc.Write("21.", f.IncidentName)
	enc.Write("22.", f.DateInitiated)
	enc.Write("23.", f.TimeInitiated)
	enc.Write("24.", f.TrackingNumber)
	enc.Write("25.", f.RequestedBy)
	enc.Write("26.", f.PreparedBy)
	enc.Write("27.", f.ApprovedBy)
	enc.Write("28.", f.QtyUnit)
	enc.Write("29.", f.ResourceDescription)
	enc.Write("30.", f.ResourceArrival)
	enc.Write("31.", f.Priority)
	enc.Write("32.", f.EstdCost)
	enc.Write("33.", f.DeliverTo)
	enc.Write("34.", f.DeliverToLocation)
	enc.Write("35.", f.Substitutes)
	enc.Write("36a.", f.EquipmentOperator)
	enc.Write("36b.", f.Lodging)
	enc.Write("36c.", f.Fuel)
	enc.Write("36d.", f.FuelType)
	enc.Write("36e.", f.Power)
	enc.Write("36f.", f.Meals)
	enc.Write("36g.", f.Maintenance)
	enc.Write("36h.", f.Water)
	enc.Write("36i.", f.Other)
	enc.Write("37.", f.Instructions)
	f.StdFields.EncodeFooter(enc)
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}
