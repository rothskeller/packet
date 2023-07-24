package eoc213rr

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

type eoc213RREdit struct {
	common.StdFieldsEdit
	IncidentName        message.EditField
	DateInitiated       message.EditField
	TimeInitiated       message.EditField
	RequestedBy         message.EditField
	PreparedBy          message.EditField
	ApprovedBy          message.EditField
	WithSignature       message.EditField
	QtyUnit             message.EditField
	ResourceDescription message.EditField
	ResourceArrival     message.EditField
	Priority            message.EditField
	EstdCost            message.EditField
	DeliverTo           message.EditField
	DeliverToLocation   message.EditField
	Substitutes         message.EditField
	EquipmentOperator   message.EditField
	Lodging             message.EditField
	Fuel                message.EditField
	FuelType            message.EditField
	Power               message.EditField
	Meals               message.EditField
	Maintenance         message.EditField
	Water               message.EditField
	Other               message.EditField
	Instructions        message.EditField
	fields              []*message.EditField
}

// EditFields returns the set of editable fields of the message.
func (f *EOC213RR) EditFields() []*message.EditField {
	if f.edit == nil {
		f.edit = &eoc213RREdit{
			StdFieldsEdit: common.StdFieldsEditTemplate,
			IncidentName: message.EditField{
				Label: "Incident Name",
				Width: 39,
				Help:  "This is the name of the incident for which resources are being requested.  It is required.",
			},
			DateInitiated: message.EditField{
				Label: "Date Initiated",
				Width: 10,
				Help:  "This is the date on which the resource request was initiated, in MM/DD/YYYY format.  It is required.",
				Hint:  "MM/DD/YYYY",
			},
			TimeInitiated: message.EditField{
				Label: "Time Initiated",
				Width: 5,
				Help:  "This is the time at which the resource requested was initiated, in HH:MM format (24-hour clock).  It is required.",
				Hint:  "HH:MM",
			},
			RequestedBy: message.EditField{
				Label:     "Requested By",
				Width:     37,
				Multiline: true,
				Help:      "This is the name, agency, position, email, and phone number of the person requesting resources.  It is required.",
			},
			PreparedBy: message.EditField{
				Label:     "Prepared By",
				Width:     37,
				Multiline: true,
				Help:      "This is the name, position, email, and phone number of the person who prepared the resource request form.",
			},
			ApprovedBy: message.EditField{
				Label:     "Approved By",
				Width:     37,
				Multiline: true,
				Help:      "This is the name, position, email, and phone number of the person who approved the resource request.",
			},
			WithSignature: message.EditField{
				Label:   "With Signature",
				Width:   3,
				Help:    "If checked, this indicates that the original paper resource request form was signed.",
				Choices: []string{"checked"},
			},
			QtyUnit: message.EditField{
				Label:     "Qty/Unit",
				Width:     9,
				Multiline: true,
				Help:      "This is the quantity (with units where applicable) of the resource requested.  If multiple resources are being requested, enter the quantity of each on a separate line.  This field is required.",
			},
			ResourceDescription: message.EditField{
				Label:     "Resource Description",
				Width:     34,
				Multiline: true,
				Help:      "This is the description of the resource requested.  If multiple resources are being requested, enter the description of each on a separate line, parallel to the corresponding quantity.  This field is required.",
			},
			ResourceArrival: message.EditField{
				Label:     "Resource Arrival",
				Width:     18,
				Multiline: true,
				Help:      "This is the date and time by which the resource needs to have arrived.  It is required.",
			},
			Priority: message.EditField{
				Label:   "Priority",
				Width:   6,
				Help:    `This is the priority of the resource request.  It must have the value "Now", "High" (meaning within the next 4 hours), "Medium" (meaning between 5 and 12 hours), or "Low" (meaning more than 12 hours).  It is required.`,
				Choices: []string{"Now", "High", "Medium", "Low"},
			},
			EstdCost: message.EditField{
				Label:     "Estimated Cost",
				Width:     11,
				Multiline: true,
				Help:      "This is the estimated cost of the resources requested.",
			},
			DeliverTo: message.EditField{
				Label:     "Deliver To",
				Width:     44,
				Multiline: true,
				Help:      "This is the name, agency, position, email, and phone number of the person to whom the requested resources should be delivered.  It is required.",
			},
			DeliverToLocation: message.EditField{
				Label:     "Deliver To Location",
				Width:     42,
				Multiline: true,
				Help:      "This is the address and/or GPS coordinates of the location to which the requested resources should be delivered.  It is required.",
			},
			Substitutes: message.EditField{
				Label:     "Substitutes/Sources",
				Width:     86,
				Multiline: true,
				Help:      "This is the names, phone numbers, and/or websites of suggested or substitute sources for the requested resources.",
			},
			EquipmentOperator: message.EditField{
				Label:   "Supplemental: Equipment Operator",
				Width:   3,
				Help:    "If checked, this indicates a supplemental requirement for an equipment operator.",
				Choices: []string{"checked"},
			},
			Lodging: message.EditField{
				Label:   "Supplemental: Lodging",
				Width:   3,
				Help:    "If checked, this indicates a supplemental requirement for lodging.",
				Choices: []string{"checked"},
			},
			Fuel: message.EditField{
				Label:   "Supplemental: Fuel",
				Width:   3,
				Help:    `If checked, this indicates a supplemental requirement for fuel.  The fuel type must be specified in the "Fuel Type" field.`,
				Choices: []string{"checked"},
			},
			FuelType: message.EditField{
				Label: "Supplemental: Fuel Type",
				Width: 11,
				Help:  `This is the type of fuel required.  It must be and can be set only when "Supplemental: Fuel" is checked.`,
			},
			Power: message.EditField{
				Label:   "Supplemental: Power",
				Width:   3,
				Help:    "If checked, this indicates a supplemental requirement for power.",
				Choices: []string{"checked"},
			},
			Meals: message.EditField{
				Label:   "Supplemental: Meals",
				Width:   3,
				Help:    "If checked, this indicates a supplemental requirement for meals.",
				Choices: []string{"checked"},
			},
			Maintenance: message.EditField{
				Label:   "Supplemental: Maintenance",
				Width:   3,
				Help:    "If checked, this indicates a supplemental requirement for maintenance.",
				Choices: []string{"checked"},
			},
			Water: message.EditField{
				Label:   "Supplemental: Water",
				Width:   3,
				Help:    "If checked, this indicates a supplemental requirement for water.",
				Choices: []string{"checked"},
			},
			Other: message.EditField{
				Label:   "Supplemental: Other",
				Width:   3,
				Help:    `If checked, this indicates an additional supplemental requirement (described in the "Special Instructions" field).`,
				Choices: []string{"checked"},
			},
			Instructions: message.EditField{
				Label:     "Special Instructions",
				Width:     42,
				Help:      "This is any special requirements or instructions for the resource request.",
				Multiline: true,
			},
		}
		// Set the field list slice.
		f.edit.fields = append(f.edit.StdFieldsEdit.EditFields1(),
			&f.edit.IncidentName,
			&f.edit.DateInitiated,
			&f.edit.TimeInitiated,
			&f.edit.RequestedBy,
			&f.edit.PreparedBy,
			&f.edit.ApprovedBy,
			&f.edit.WithSignature,
			&f.edit.QtyUnit,
			&f.edit.ResourceDescription,
			&f.edit.ResourceArrival,
			&f.edit.Priority,
			&f.edit.EstdCost,
			&f.edit.DeliverTo,
			&f.edit.DeliverToLocation,
			&f.edit.Substitutes,
			&f.edit.EquipmentOperator,
			&f.edit.Lodging,
			&f.edit.Fuel,
			&f.edit.FuelType,
			&f.edit.Power,
			&f.edit.Meals,
			&f.edit.Maintenance,
			&f.edit.Water,
			&f.edit.Other,
			&f.edit.Instructions,
		)
		f.edit.fields = append(f.edit.fields, f.edit.StdFieldsEdit.EditFields2()...)
		f.toEdit()
		f.validate()
	}
	return f.edit.fields
}

// ApplyEdits applies the revised Values in the EditFields to the
// message.
func (f *EOC213RR) ApplyEdits() {
	f.fromEdit()
	f.validate()
	f.toEdit()
}

func (f *EOC213RR) fromEdit() {
	f.StdFields.FromEdit(&f.edit.StdFieldsEdit)
	f.IncidentName = strings.TrimSpace(f.edit.IncidentName.Value)
	f.DateInitiated = common.CleanDate(f.edit.DateInitiated.Value)
	f.TimeInitiated = common.CleanTime(f.edit.TimeInitiated.Value)
	f.RequestedBy = strings.TrimSpace(f.edit.RequestedBy.Value)
	f.PreparedBy = strings.TrimSpace(f.edit.PreparedBy.Value)
	f.ApprovedBy = strings.TrimSpace(f.edit.ApprovedBy.Value)
	f.WithSignature = common.CleanCheckbox(f.edit.WithSignature.Value)
	f.QtyUnit = strings.TrimSpace(f.edit.QtyUnit.Value)
	f.ResourceDescription = strings.TrimSpace(f.edit.ResourceDescription.Value)
	f.ResourceArrival = strings.TrimSpace(f.edit.ResourceArrival.Value)
	f.Priority = common.ExpandRestricted(&f.edit.Priority)
	f.EstdCost = strings.TrimSpace(f.edit.EstdCost.Value)
	f.DeliverTo = strings.TrimSpace(f.edit.DeliverTo.Value)
	f.DeliverToLocation = strings.TrimSpace(f.edit.DeliverToLocation.Value)
	f.Substitutes = strings.TrimSpace(f.edit.Substitutes.Value)
	f.EquipmentOperator = common.CleanCheckbox(f.edit.EquipmentOperator.Value)
	f.Lodging = common.CleanCheckbox(f.edit.Lodging.Value)
	f.Fuel = common.CleanCheckbox(f.edit.Fuel.Value)
	f.FuelType = strings.TrimSpace(f.edit.FuelType.Value)
	f.Power = common.CleanCheckbox(f.edit.Power.Value)
	f.Meals = common.CleanCheckbox(f.edit.Meals.Value)
	f.Maintenance = common.CleanCheckbox(f.edit.Maintenance.Value)
	f.Water = common.CleanCheckbox(f.edit.Water.Value)
	f.Other = common.CleanCheckbox(f.edit.Other.Value)
	f.Instructions = strings.TrimSpace(f.edit.Instructions.Value)
}

func (f *EOC213RR) toEdit() {
	f.StdFields.ToEdit(&f.edit.StdFieldsEdit)
	f.edit.IncidentName.Value = f.IncidentName
	f.edit.DateInitiated.Value = f.DateInitiated
	f.edit.TimeInitiated.Value = f.TimeInitiated
	f.edit.RequestedBy.Value = f.RequestedBy
	f.edit.PreparedBy.Value = f.PreparedBy
	f.edit.ApprovedBy.Value = f.ApprovedBy
	f.edit.WithSignature.Value = f.WithSignature
	f.edit.QtyUnit.Value = f.QtyUnit
	f.edit.ResourceDescription.Value = f.ResourceDescription
	f.edit.ResourceArrival.Value = f.ResourceArrival
	f.edit.Priority.Value = f.Priority
	f.edit.EstdCost.Value = f.EstdCost
	f.edit.DeliverTo.Value = f.DeliverTo
	f.edit.DeliverToLocation.Value = f.DeliverToLocation
	f.edit.Substitutes.Value = f.Substitutes
	f.edit.EquipmentOperator.Value = f.EquipmentOperator
	f.edit.Lodging.Value = f.Lodging
	f.edit.Fuel.Value = f.Fuel
	f.edit.FuelType.Value = f.FuelType
	f.edit.Power.Value = f.Power
	f.edit.Meals.Value = f.Meals
	f.edit.Maintenance.Value = f.Maintenance
	f.edit.Water.Value = f.Water
	f.edit.Other.Value = f.Other
	f.edit.Instructions.Value = f.Instructions
}

func (f *EOC213RR) validate() {
	f.edit.StdFieldsEdit.Validate()
	common.ValidateRequired(&f.edit.IncidentName)
	if common.ValidateRequired(&f.edit.DateInitiated) {
		common.ValidateDate(&f.edit.DateInitiated)
	}
	if common.ValidateRequired(&f.edit.TimeInitiated) {
		common.ValidateTime(&f.edit.TimeInitiated)
	}
	common.ValidateRequired(&f.edit.RequestedBy)
	common.ValidateRequired(&f.edit.QtyUnit)
	common.ValidateRequired(&f.edit.ResourceDescription)
	common.ValidateRequired(&f.edit.ResourceArrival)
	if common.ValidateRequired(&f.edit.Priority) {
		common.ValidateRestricted(&f.edit.Priority)
	}
	common.ValidateRequired(&f.edit.DeliverTo)
	common.ValidateRequired(&f.edit.DeliverToLocation)
	if f.edit.Fuel.Value == "" {
		if f.edit.FuelType.Value != "" {
			f.edit.FuelType.Problem = `The "Fuel Type" field cannot be set unless "Fuel" is checked.`
		} else {
			f.edit.FuelType.Problem = ""
		}
	} else {
		if f.edit.FuelType.Value == "" {
			f.edit.FuelType.Problem = `The "Fuel Type" field is required when "Fuel" is checked.`
		} else {
			f.edit.FuelType.Problem = ""
		}
	}
}
