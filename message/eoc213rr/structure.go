package eoc213rr

import (
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// Type is the type definition for an EOC-213RR resource request form.
var Type = message.Type{
	Tag:     "EOC213RR",
	Name:    "EOC-213RR resource request form",
	Article: "an",
	Create:  New,
	Decode:  decode,
}

// EOC213RR holds an EOC-213RR resource request form.
type EOC213RR struct {
	common.StdFields
	IncidentName        string
	DateInitiated       string
	TimeInitiated       string
	TrackingNumber      string
	RequestedBy         string
	PreparedBy          string
	ApprovedBy          string
	WithSignature       string // added in v2.4
	QtyUnit             string
	ResourceDescription string
	ResourceArrival     string
	Priority            string
	EstdCost            string
	DeliverTo           string
	DeliverToLocation   string
	Substitutes         string
	EquipmentOperator   string
	Lodging             string
	Fuel                string
	FuelType            string
	Power               string
	Meals               string
	Maintenance         string
	Water               string
	Other               string
	Instructions        string
	edit                *eoc213RREdit
}

// Type returns the message type definition.
func (*EOC213RR) Type() *message.Type { return &Type }
