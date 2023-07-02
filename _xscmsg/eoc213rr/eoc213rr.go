package eoc213rr

import (
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Type is the type definition for a resource request form.
var Type = typedmsg.MessageType{
	Tag:       "EOC213RR",
	Name:      "EOC-213RR resource request form",
	Article:   "an",
	Create:    create,
	Recognize: recognize,
}

// EOC213RR is a resource request form.
type EOC213RR struct {
	*xscmsg.StdForm
	IncidentName        string
	DateInitiated       string
	TimeInitiated       string
	TrackingNumber      string
	RequestedBy         string
	PreparedBy          string
	ApprovedBy          string
	QtyUnit             string
	ResourceDescription string
	ResourceArrival     string
	Priority            string
	EstdCost            string
	DeliverTo           string
	DeliverToLocation   string
	Substitutes         string
	EquipmentOperator   bool
	Lodging             bool
	Fuel                bool
	FuelType            string
	Power               bool
	Meals               bool
	Maintenance         bool
	Water               bool
	Other               bool
	Instructions        string
	UnknownFields       []pktmsg.TaggedField
	editFields          []xscmsg.Field
}

// NewEOC213RR creates a new resource request form.
func NewEOC213RR() *EOC213RR {
	m := &EOC213RR{StdForm: xscmsg.NewStdForm(nil)}
	m.FormHTML = "form-scco-eoc-213rr.html"
	m.FormVersion = "2.3"
	return m
}

func create() typedmsg.Message { return NewEOC213RR() }

func recognize(base *pktmsg.Message) typedmsg.Message {
	if base.FormHTML != "form-scco-eoc-213rr.html" || xscmsg.OlderVersion(base.FormVersion, "2.0") {
		return nil
	}
	m := &EOC213RR{StdForm: xscmsg.NewStdForm(xscmsg.NewBaseMessage(base))}
	m.readTaggedFields()
	return m
}

func (m *EOC213RR) readTaggedFields() {
	// Make a copy of the tagged fields since this method is destructive.
	fields := make([]pktmsg.TaggedField, len(m.TaggedFields))
	copy(fields, m.TaggedFields)
	// Read the standard form fields.
	fields = m.StdForm.ReadTaggedFields(fields)
	// Read the form-specific tagged fields.
	j := 0
	for _, tf := range fields {
		switch tf.Tag {
		case "21.":
			m.IncidentName = tf.Value
		case "22.":
			m.DateInitiated = tf.Value
		case "23.":
			m.TimeInitiated = tf.Value
		case "24.":
			m.TrackingNumber = tf.Value
		case "25.":
			m.RequestedBy = tf.Value
		case "26.":
			m.PreparedBy = tf.Value
		case "27.":
			m.ApprovedBy = tf.Value
		case "28.":
			m.QtyUnit = tf.Value
		case "29.":
			m.ResourceDescription = tf.Value
		case "30.":
			m.ResourceArrival = tf.Value
		case "31.":
			m.Priority = tf.Value
		case "32.":
			m.EstdCost = tf.Value
		case "33.":
			m.DeliverTo = tf.Value
		case "34.":
			m.DeliverToLocation = tf.Value
		case "35.":
			m.Substitutes = tf.Value
		case "36a.":
			m.EquipmentOperator = tf.Value != ""
		case "36b.":
			m.Lodging = tf.Value != ""
		case "36c.":
			m.Fuel = tf.Value != ""
		case "36d.":
			m.FuelType = tf.Value
		case "36e.":
			m.Power = tf.Value != ""
		case "36f.":
			m.Meals = tf.Value != ""
		case "36g.":
			m.Maintenance = tf.Value != ""
		case "36h.":
			m.Water = tf.Value != ""
		case "36i.":
			m.Other = tf.Value != ""
		case "37.":
			m.Instructions = tf.Value
		default:
			fields[j], j = tf, j+1
		}
	}
	m.UnknownFields = fields[:j]
}

func (m *EOC213RR) makeTaggedFields() {
	m.TaggedFields = m.StdForm.MakeLeadingTaggedFields()
	m.TaggedFields = append(m.TaggedFields, []pktmsg.TaggedField{
		{Tag: "21.", Value: m.IncidentName},
		{Tag: "22.", Value: m.DateInitiated},
		{Tag: "23.", Value: m.TimeInitiated},
		{Tag: "24.", Value: m.TrackingNumber},
		{Tag: "25.", Value: m.RequestedBy},
		{Tag: "26.", Value: m.PreparedBy},
		{Tag: "27.", Value: m.ApprovedBy},
		{Tag: "28.", Value: m.QtyUnit},
		{Tag: "29.", Value: m.ResourceDescription},
		{Tag: "30.", Value: m.ResourceArrival},
		{Tag: "31.", Value: m.Priority},
		{Tag: "32.", Value: m.EstdCost},
		{Tag: "33.", Value: m.DeliverTo},
		{Tag: "34.", Value: m.DeliverToLocation},
		{Tag: "35.", Value: m.Substitutes},
		{Tag: "36a.", Value: checkedIfTrue(m.EquipmentOperator)},
		{Tag: "36b.", Value: checkedIfTrue(m.Lodging)},
		{Tag: "36c.", Value: checkedIfTrue(m.Fuel)},
		{Tag: "36d.", Value: m.FuelType},
		{Tag: "36e.", Value: checkedIfTrue(m.Power)},
		{Tag: "36f.", Value: checkedIfTrue(m.Meals)},
		{Tag: "36g.", Value: checkedIfTrue(m.Maintenance)},
		{Tag: "36h.", Value: checkedIfTrue(m.Water)},
		{Tag: "36i.", Value: checkedIfTrue(m.Other)},
		{Tag: "37.", Value: m.Instructions},
	}...)
	m.TaggedFields = append(m.TaggedFields, m.StdForm.MakeTrailingTaggedFields()...)
}
func checkedIfTrue(value bool) string {
	if value {
		return "checked"
	}
	return ""
}

// Type returns the type of the message.
func (m *EOC213RR) Type() *typedmsg.MessageType { return &Type }

// Validate checks the validity of the message and returns strings describing
// the issues.  This is used for received messages, and only checks for validity
// issues that Outpost and/or PackItForms check for.
func (m *EOC213RR) Validate() (problems []string) {
	problems = m.StdForm.Validate()
	if m.IncidentName == "" {
		problems = append(problems, "The incident name is required.")
	}
	if m.DateInitiated == "" {
		problems = append(problems, "The date initiated is required.")
	} else if !xscmsg.ValidDate(m.DateInitiated) {
		problems = append(problems, "The date initiated is not a valid date.")
	}
	if m.TimeInitiated == "" {
		problems = append(problems, "The time initiated is required.")
	} else if !xscmsg.ValidTime(m.TimeInitiated) {
		problems = append(problems, "The time initiated is not a valid time.")
	}
	if m.RequestedBy == "" {
		problems = append(problems, "The requested by field is required.")
	}
	if m.QtyUnit == "" {
		problems = append(problems, "The qty/unit field is required.")
	}
	if m.ResourceDescription == "" {
		problems = append(problems, "The resource description field is required.")
	}
	if m.ResourceArrival == "" {
		problems = append(problems, "The resource arrival field is required.")
	}
	if m.Priority == "" {
		problems = append(problems, "The priority field is required.")
	} else if !xscmsg.ValidRestrictedValue(m.Priority, priorityChoices) {
		problems = append(problems, "The priority is not one of the recognized priority choices (Now, High, Medium, or Low).")
	}
	if m.DeliverTo == "" {
		problems = append(problems, "The deliver to field is required.")
	}
	if m.DeliverToLocation == "" {
		problems = append(problems, "The delivery location field is required.")
	}
	return problems
}

// View returns the set of viewable fields of the message.
func (m *EOC213RR) View() []xscmsg.LabelValue {
	var supplemental []string
	if m.EquipmentOperator {
		supplemental = append(supplemental, "Equipment Operator")
	}
	if m.Lodging {
		supplemental = append(supplemental, "Lodging")
	}
	if m.Fuel {
		if m.FuelType != "" {
			supplemental = append(supplemental, "Fuel ("+m.FuelType+")")
		} else {
			supplemental = append(supplemental, "Fuel")
		}
	}
	if m.Power {
		supplemental = append(supplemental, "Power")
	}
	if m.Meals {
		supplemental = append(supplemental, "Meals")
	}
	if m.Maintenance {
		supplemental = append(supplemental, "Maintenance")
	}
	if m.Water {
		supplemental = append(supplemental, "Water")
	}
	if m.Other {
		supplemental = append(supplemental, "Other")
	}
	var lvs = m.StdForm.MakeLeadingViewFields()
	lvs = append(lvs, []xscmsg.LabelValue{
		{Label: "Incident Name", Value: m.IncidentName},
		{Label: "Date/Time Initiated", Value: m.DateInitiated + " " + m.TimeInitiated},
		{Label: "Tracking Number", Value: m.TrackingNumber},
		{Label: "Requested By", Value: m.RequestedBy},
		{Label: "Prepared By", Value: m.PreparedBy},
		{Label: "Approved By", Value: m.ApprovedBy},
		{Label: "Qty/Unit", Value: m.QtyUnit},
		{Label: "Resource Description", Value: m.ResourceDescription},
		{Label: "Resource Arrival", Value: m.ResourceArrival},
		{Label: "Priority", Value: m.Priority},
		{Label: "Estimated Cost", Value: m.EstdCost},
		{Label: "Deliver To", Value: m.DeliverTo},
		{Label: "Delivery Location", Value: m.DeliverToLocation},
		{Label: "Substitute/Suggested Sources", Value: m.Substitutes},
		{Label: "Supplemental", Value: strings.Join(supplemental, ", ")},
		{Label: "Special Instructions", Value: m.Instructions},
	}...)
	lvs = append(lvs, m.StdForm.MakeTrailingViewFields()...)
	return lvs
}

// Edit returns the set of editable fields of the message.
func (m *EOC213RR) Edit() []xscmsg.Field {
	if m.editFields == nil {
		fuel := &fuelField{xscmsg.CheckboxField(&m.Fuel)}
		m.editFields = m.StdForm.MakeLeadingEditFields()
		m.editFields = append(m.editFields, []xscmsg.Field{
			xscmsg.WrapRequiredField(&incidentNameField{xscmsg.BaseField(&m.IncidentName)}),
			xscmsg.WrapRequiredField(xscmsg.WrapDateField(&dateInitiatedField{xscmsg.BaseField(&m.DateInitiated)})),
			xscmsg.WrapRequiredField(xscmsg.WrapTimeField(&timeInitiatedField{xscmsg.BaseField(&m.TimeInitiated)})),
			xscmsg.WrapRequiredField(&requestedByField{xscmsg.BaseField(&m.RequestedBy)}),
			&preparedByField{xscmsg.BaseField(&m.PreparedBy)},
			&approvedByField{xscmsg.BaseField(&m.ApprovedBy)},
			xscmsg.WrapRequiredField(&qtyUnitField{xscmsg.BaseField(&m.QtyUnit)}),
			xscmsg.WrapRequiredField(&resourceDescriptionField{xscmsg.BaseField(&m.ResourceDescription)}),
			xscmsg.WrapRequiredField(&resourceArrivalField{xscmsg.BaseField(&m.ResourceArrival)}),
			xscmsg.WrapRequiredField(xscmsg.WrapRestrictedField(&priorityField{xscmsg.BaseField(&m.Priority)})),
			&estdCostField{xscmsg.BaseField(&m.EstdCost)},
			xscmsg.WrapRequiredField(&deliverToField{xscmsg.BaseField(&m.DeliverTo)}),
			xscmsg.WrapRequiredField(&deliverToLocationField{xscmsg.BaseField(&m.DeliverToLocation)}),
			&substitutesField{xscmsg.BaseField(&m.Substitutes)},
			&equipmentOperatorField{xscmsg.CheckboxField(&m.EquipmentOperator)},
			&lodgingField{xscmsg.CheckboxField(&m.Lodging)},
			fuel,
			&fuelTypeField{xscmsg.BaseField(&m.FuelType), fuel},
			&powerField{xscmsg.CheckboxField(&m.Power)},
			&mealsField{xscmsg.CheckboxField(&m.Meals)},
			&maintenanceField{xscmsg.CheckboxField(&m.Maintenance)},
			&waterField{xscmsg.CheckboxField(&m.Water)},
			&otherField{xscmsg.CheckboxField(&m.Other)},
			&instructionsField{xscmsg.BaseField(&m.Instructions)},
		}...)
		m.editFields = append(m.editFields, m.StdForm.MakeTrailingEditFields()...)
	}
	return m.editFields
}

// Save returns the message encoded in a format suitable for saving to local
// storage.
func (m *EOC213RR) Save() string {
	m.makeTaggedFields()
	return m.StdForm.Save()
}

// Transmit returns the destination addresses, subject header, and body
// of the message, suitable for transmission through JNOS.
func (m *EOC213RR) Transmit() ([]string, string, string) {
	m.makeTaggedFields()
	return m.StdForm.Transmit()
}

// TODO recommended To: County EOC Planning Section

type incidentNameField struct{ xscmsg.Field }

func (f incidentNameField) Label() string             { return "Incident Name" }
func (f incidentNameField) Size() (width, height int) { return 34, 1 }
func (f incidentNameField) Help() string {
	return "This is the name of the incident for which resources are being requested."
}

type dateInitiatedField struct{ xscmsg.Field }

func (f dateInitiatedField) Label() string { return "Date Initiated" }
func (f dateInitiatedField) Help() string {
	return "This is the date when the resource request is initiated."
}

type timeInitiatedField struct{ xscmsg.Field }

func (f timeInitiatedField) Label() string { return "Time Initiated" }
func (f timeInitiatedField) Help() string {
	return "This is the time when the resource request is initiated."
}

type requestedByField struct{ xscmsg.Field }

func (f requestedByField) Label() string             { return "Requested By" }
func (f requestedByField) Size() (width, height int) { return 34, 1 }
func (f requestedByField) Help() string {
	return "This is the person requesting the resources.  This field should include the person's name, agency, position, and contact information."
}

type preparedByField struct{ xscmsg.Field }

func (f preparedByField) Label() string             { return "Prepared By" }
func (f preparedByField) Size() (width, height int) { return 34, 1 }
func (f preparedByField) Help() string {
	return "This is the person who prepared the resource request form.  This field should include the person's name, agency, position, and contact information."
}

type approvedByField struct{ xscmsg.Field }

func (f approvedByField) Label() string             { return "Approved By" }
func (f approvedByField) Size() (width, height int) { return 34, 1 }
func (f approvedByField) Help() string {
	return "This is the person who approved the resource request (a section chief in the requesting EOC, or supervising official at the requesting agency).  This field should include the person's name, agency, position, and contact information."
}

type qtyUnitField struct{ xscmsg.Field }

func (f qtyUnitField) Label() string             { return "Qty/Unit" }
func (f qtyUnitField) Size() (width, height int) { return 8, 5 }
func (f qtyUnitField) Help() string {
	return "This is the count of each item requested, with units when appropriate.  When requesting multiple items, enter one count per line."
}

type resourceDescriptionField struct{ xscmsg.Field }

func (f resourceDescriptionField) Label() string             { return "Resource Description" }
func (f resourceDescriptionField) Size() (width, height int) { return 30, 5 }
func (f resourceDescriptionField) Help() string {
	return "This is the description of each item requested.  When requesting multiple items, enter one description per line, parallel with the counts in the Qty/Unit field."
}

type resourceArrivalField struct{ xscmsg.Field }

func (f resourceArrivalField) Label() string             { return "Arrival" }
func (f resourceArrivalField) Size() (width, height int) { return 16, 5 }
func (f resourceArrivalField) Help() string {
	return "This is the date and time when the requested resources need to arrive."
}

type priorityField struct{ xscmsg.Field }

var priorityChoices = []string{"Now", "High", "Medium", "Low"}

func (f priorityField) Label() string             { return "Priority" }
func (f priorityField) Size() (width, height int) { return 0, 0 }
func (f priorityField) Choices() []string         { return priorityChoices }
func (f priorityField) Help() string {
	return `This is the priority of the resource request.  Use "High" for resources needed within 4 hours, "Medium" for resources needed within 12 hours, and "Low" for resources not needed in the next 12 hours.`
}

type estdCostField struct{ xscmsg.Field }

func (f estdCostField) Label() string             { return "Estimated Cost" }
func (f estdCostField) Size() (width, height int) { return 10, 5 }
func (f estdCostField) Help() string {
	return "This is the estimated cost of the resources requested."
}

type deliverToField struct{ xscmsg.Field }

func (f deliverToField) Label() string             { return "Deliver To" }
func (f deliverToField) Size() (width, height int) { return 40, 1 }
func (f deliverToField) Help() string {
	return "This is the person to whom the resources should be delivered.  This field should include the person's name, agency, position, and contact information."
}

type deliverToLocationField struct{ xscmsg.Field }

func (f deliverToLocationField) Label() string             { return "Delivery Location" }
func (f deliverToLocationField) Size() (width, height int) { return 38, 1 }
func (f deliverToLocationField) Help() string {
	return "This is the location to which the resources should be delivered.  Include the address and/or GPS coordinates.  Also include the site type."
}

type substitutesField struct{ xscmsg.Field }

func (f substitutesField) Label() string             { return "Substitute/Suggested Sources" }
func (f substitutesField) Size() (width, height int) { return 80, 1 }
func (f substitutesField) Help() string {
	return "This field suggests sources from which the resources can be acquired.  Include name, phone number, and website."
}

type equipmentOperatorField struct{ xscmsg.Field }

func (f equipmentOperatorField) Label() string { return "Supplemental: Equipment Operator" }
func (f equipmentOperatorField) Help() string {
	return "This field indicates a supplemental requirement for an equipment operator."
}

type lodgingField struct{ xscmsg.Field }

func (f lodgingField) Label() string { return "Supplemental: Lodging" }
func (f lodgingField) Help() string {
	return "This field indicates a supplemental requirement for lodging."
}

type fuelField struct{ xscmsg.Field }

func (f fuelField) Label() string { return "Supplemental: Fuel" }
func (f fuelField) Help() string {
	return "This field indicates a supplemental requirement for fuel."
}

type fuelTypeField struct {
	xscmsg.Field
	checkbox xscmsg.Field
}

func (f fuelTypeField) Label() string             { return "Supplemental: Fuel Type" }
func (f fuelTypeField) Size() (width, height int) { return 10, 1 }
func (f fuelTypeField) Problem() string {
	if f.Value() != "" && f.checkbox.Value() == "" {
		return `The "Fuel Type" field must not have a value unless the Fuel box is checked.`
	}
	return ""
}
func (f fuelTypeField) Help() string {
	return "This field identifies the type of fuel needed."
}

type powerField struct{ xscmsg.Field }

func (f powerField) Label() string { return "Supplemental: Power" }
func (f powerField) Help() string {
	return "This field indicates a supplemental requirement for power."
}

type mealsField struct{ xscmsg.Field }

func (f mealsField) Label() string { return "Supplemental: Meals" }
func (f mealsField) Help() string {
	return "This field indicates a supplemental requirement for meals."
}

type maintenanceField struct{ xscmsg.Field }

func (f maintenanceField) Label() string { return "Supplemental: Maintenance" }
func (f maintenanceField) Help() string {
	return "This field indicates a supplemental requirement for maintenance of the requested resources."
}

type waterField struct{ xscmsg.Field }

func (f waterField) Label() string { return "Supplemental: Water" }
func (f waterField) Help() string {
	return "This field indicates a supplemental requirement for water."
}

type otherField struct{ xscmsg.Field }

func (f otherField) Label() string { return "Supplemental: Other" }
func (f otherField) Help() string {
	return "This field indicates that other supplemental requirements exist, and are described in the Special Instructions field."
}

type instructionsField struct{ xscmsg.Field }

func (f instructionsField) Label() string             { return "Special Instructions" }
func (f instructionsField) Size() (width, height int) { return 38, 8 }
func (f instructionsField) Help() string {
	return "This is any special instructions or other additional information about the resource request."
}

// ...

func (m *EOC213RR) GetBody() string         { return m.Instructions }
func (m *EOC213RR) SetBody(value string)    { m.Instructions = value }
func (m *EOC213RR) GetSubject() string      { return m.IncidentName }
func (m *EOC213RR) SetSubject(value string) { m.IncidentName = value }
