package racesmar

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies shelter status forms.
const (
	Tag = "RACES-MAR"
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)
}

func create() xscmsg.Message {
	return xscform.CreateForm(formtype23, makeFields23())
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.Message {
	if form == nil {
		return nil
	}
	if form.FormType == formtype23.HTML && !xscmsg.OlderVersion(form.FormVersion, "2.3") {
		return xscform.AdoptForm(formtype23, makeFields23(), msg, form)
	}
	if form.FormType == formtype21.HTML && !xscmsg.OlderVersion(form.FormVersion, "2.1") {
		return xscform.AdoptForm(formtype21, makeFields21(), msg, form)
	}
	if form.FormType != formtype16.HTML && !xscmsg.OlderVersion(form.FormVersion, "1.6") {
		return xscform.AdoptForm(formtype16, makeFields16(), msg, form)
	}
	return nil
}

var formtype23 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "RACES mutual aid request form",
	Article: "a",
	HTML:    "form-oa-mutual-aid-request-v2.html",
	Version: "2.3",
}

func makeFields23() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		xscform.NewField(agencyID, true),
		xscform.NewField(eventNameID, true),
		xscform.NewField(eventNumberID, false),
		xscform.NewField(assignmentID, true),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty1ID, true)},
		&xscform.ChoicesField{Field: *xscform.NewField(role1ID, true), Choices: roleChoices},
		xscform.NewField(position1ID, false),
		&resourcesRoleField{*xscform.NewField(resourcesRole1ID23, false)},
		&typeField{*xscform.NewField(preferredType1ID23, false)},
		&typeField{*xscform.NewField(minimumType1ID23, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty2ID, false)},
		&xscform.ChoicesField{Field: *xscform.NewField(role2ID, false), Choices: roleChoices},
		xscform.NewField(position2ID, false),
		&resourcesRoleField{*xscform.NewField(resourcesRole2ID23, false)},
		&typeField{*xscform.NewField(preferredType2ID23, false)},
		&typeField{*xscform.NewField(minimumType2ID23, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty3ID, false)},
		&xscform.ChoicesField{Field: *xscform.NewField(role3ID, false), Choices: roleChoices},
		xscform.NewField(position3ID, false),
		&resourcesRoleField{*xscform.NewField(resourcesRole3ID23, false)},
		&typeField{*xscform.NewField(preferredType3ID23, false)},
		&typeField{*xscform.NewField(minimumType3ID23, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty4ID, false)},
		&xscform.ChoicesField{Field: *xscform.NewField(role4ID, false), Choices: roleChoices},
		xscform.NewField(position4ID, false),
		&resourcesRoleField{*xscform.NewField(resourcesRole4ID23, false)},
		&typeField{*xscform.NewField(preferredType4ID23, false)},
		&typeField{*xscform.NewField(minimumType4ID23, false)},
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty5ID, false)},
		&xscform.ChoicesField{Field: *xscform.NewField(role5ID, false), Choices: roleChoices},
		xscform.NewField(position5ID, false),
		&resourcesRoleField{*xscform.NewField(resourcesRole5ID23, false)},
		&typeField{*xscform.NewField(preferredType5ID23, false)},
		&typeField{*xscform.NewField(minimumType5ID23, false)},
		xscform.NewField(arrivalDatesID, true),
		xscform.NewField(arrivalTimesID, true),
		xscform.NewField(neededDatesID, true),
		xscform.NewField(neededTimesID, true),
		xscform.NewField(reportingLocationID, true),
		xscform.NewField(contactOnArrivalID, true),
		xscform.NewField(travelInfoID, true),
		xscform.NewField(requesterNameID, true),
		xscform.NewField(requesterTitleID, true),
		xscform.NewField(requesterContactID, true),
		xscform.NewField(agencyApproverNameID, true),
		xscform.NewField(agencyApproverTitleID, true),
		xscform.NewField(agencyApproverContactID, true),
		&xscform.DateField{Field: *xscform.NewField(agencyApprovedDateID, true)},
		&xscform.TimeField{Field: *xscform.NewField(agencyApprovedTimeID, true)},
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var formtype21 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "RACES mutual aid request form",
	Article: "a",
	HTML:    "form-oa-mutual-aid-request-v2.html",
	Version: "2.1",
}

func makeFields21() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		xscform.NewField(agencyID, true),
		xscform.NewField(eventNameID, true),
		xscform.NewField(eventNumberID, false),
		xscform.NewField(assignmentID, true),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty1ID, true)},
		xscform.NewField(resourcesRole1ID21, true),
		xscform.NewField(preferredType1ID21, true),
		xscform.NewField(minimumType1ID21, true),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty2ID, false)},
		xscform.NewField(resourcesRole2ID21, false),
		xscform.NewField(preferredType2ID21, false),
		xscform.NewField(minimumType2ID21, false),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty3ID, false)},
		xscform.NewField(resourcesRole3ID21, false),
		xscform.NewField(preferredType3ID21, false),
		xscform.NewField(minimumType3ID21, false),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty4ID, false)},
		xscform.NewField(resourcesRole4ID21, false),
		xscform.NewField(preferredType4ID21, false),
		xscform.NewField(minimumType4ID21, false),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQty5ID, false)},
		xscform.NewField(resourcesRole5ID21, false),
		xscform.NewField(preferredType5ID21, false),
		xscform.NewField(minimumType5ID21, false),
		xscform.NewField(arrivalDatesID, true),
		xscform.NewField(arrivalTimesID, true),
		xscform.NewField(neededDatesID, true),
		xscform.NewField(neededTimesID, true),
		xscform.NewField(reportingLocationID, true),
		xscform.NewField(contactOnArrivalID, true),
		xscform.NewField(travelInfoID, true),
		xscform.NewField(requesterNameID, true),
		xscform.NewField(requesterTitleID, true),
		xscform.NewField(requesterContactID, true),
		xscform.NewField(agencyApproverNameID, true),
		xscform.NewField(agencyApproverTitleID, true),
		xscform.NewField(agencyApproverContactID, true),
		&xscform.DateField{Field: *xscform.NewField(agencyApprovedDateID, true)},
		&xscform.TimeField{Field: *xscform.NewField(agencyApprovedTimeID, true)},
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var formtype16 = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "RACES mutual aid request form",
	Article: "a",
	HTML:    "form-oa-mutual-aid-request.html",
	Version: "1.6",
}

func makeFields16() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		&handlingField{ChoicesField: *xscform.FHandling().(*xscform.ChoicesField)},
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		xscform.NewField(agencyID, true),
		xscform.NewField(eventNameID, true),
		xscform.NewField(eventNumberID, false),
		xscform.NewField(assignmentID, true),
		&xscform.CardinalNumberField{Field: *xscform.NewField(resourcesQtyID, true)},
		xscform.NewField(resourcesRoleID, true),
		xscform.NewField(preferredTypeID, true),
		xscform.NewField(minimumTypeID, true),
		xscform.NewField(arrivalDatesID, true),
		xscform.NewField(arrivalTimesID, true),
		xscform.NewField(neededDatesID, true),
		xscform.NewField(neededTimesID, true),
		xscform.NewField(reportingLocationID, true),
		xscform.NewField(contactOnArrivalID, true),
		xscform.NewField(travelInfoID, true),
		xscform.NewField(requesterNameID, true),
		xscform.NewField(requesterTitleID, true),
		xscform.NewField(requesterContactID, true),
		xscform.NewField(agencyApproverNameID, true),
		xscform.NewField(agencyApproverTitleID, true),
		xscform.NewField(agencyApproverContactID, true),
		&xscform.DateField{Field: *xscform.NewField(agencyApprovedDateID, true)},
		&xscform.TimeField{Field: *xscform.NewField(agencyApprovedTimeID, true)},
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var (
	agencyID = &xscmsg.FieldID{
		Tag:        "15.",
		Annotation: "agency",
		Label:      "Agency Name",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
	}
	eventNameID = &xscmsg.FieldID{
		Tag:        "16a.",
		Annotation: "event-name",
		Label:      "Event / Incident Name",
		Comment:    "required",
	}
	eventNumberID = &xscmsg.FieldID{
		Tag:        "16b.",
		Annotation: "event-number",
		Label:      "Event / Incident Number",
	}
	assignmentID = &xscmsg.FieldID{
		Tag:        "17.",
		Annotation: "assignment",
		Label:      "Assignment",
		Comment:    "required",
	}
	resourcesQtyID = &xscmsg.FieldID{
		Tag:        "18a.",
		Annotation: "resources-qty",
		Label:      "Qty",
		Comment:    "required cardinal-number",
	}
	resourcesRoleID = &xscmsg.FieldID{
		Tag:        "18b.",
		Annotation: "resources-role",
		Label:      "Role/Position",
		Comment:    "required",
	}
	preferredTypeID = &xscmsg.FieldID{
		Tag:        "18c.",
		Annotation: "preferred-type",
		Label:      "Preferred Type",
		Comment:    "required",
	}
	minimumTypeID = &xscmsg.FieldID{
		Tag:        "18d.",
		Annotation: "minimum-type",
		Label:      "Minimum Type",
		Comment:    "required",
	}
	resourcesQty1ID = &xscmsg.FieldID{
		Tag:        "18.1a.",
		Annotation: "resources-qty",
		Label:      "Resource 1 Qty",
		Comment:    "required cardinal-number",
	}
	role1ID = &xscmsg.FieldID{
		Tag:        "18.1e.",
		Annotation: "role",
		Label:      "Resource 1 Role",
		Comment:    "required: Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
	}
	roleChoices = []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"}
	position1ID = &xscmsg.FieldID{
		Tag:        "18.1f.",
		Annotation: "position",
		Label:      "Resource 1 Position (for example, Checkpoint)",
	}
	resourcesRole1ID23 = &xscmsg.FieldID{
		Tag:        "18.1b.",
		Annotation: "resources-role",
		ReadOnly:   true,
	}
	resourcesRole1ID21 = &xscmsg.FieldID{
		Tag:        "18.1b.",
		Annotation: "resources-role",
		Label:      "Resource 1 Role/Position",
		Comment:    "required",
	}
	preferredType1ID23 = &xscmsg.FieldID{
		Tag:        "18.1c.",
		Annotation: "preferred-type",
		Label:      "Resource 1 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	preferredType1ID21 = &xscmsg.FieldID{
		Tag:        "18.1c.",
		Annotation: "preferred-type",
		Label:      "Resource 1 Preferred Type",
		Comment:    "required",
	}
	minimumType1ID23 = &xscmsg.FieldID{
		Tag:        "18.1d.",
		Annotation: "minimum-type",
		Label:      "Resource 1 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	minimumType1ID21 = &xscmsg.FieldID{
		Tag:        "18.1d.",
		Annotation: "minimum-type",
		Label:      "Resource 1 Minimum Type",
		Comment:    "required",
	}
	resourcesQty2ID = &xscmsg.FieldID{
		Tag:        "18.2a.",
		Annotation: "resources-qty",
		Label:      "Resource 2 Qty",
		Comment:    "cardinal-number",
	}
	role2ID = &xscmsg.FieldID{
		Tag:        "18.2e.",
		Annotation: "role",
		Label:      "Resource 2 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
	}
	position2ID = &xscmsg.FieldID{
		Tag:        "18.2f.",
		Annotation: "position",
		Label:      "Resource 2 Position (for example, Checkpoint)",
	}
	resourcesRole2ID23 = &xscmsg.FieldID{
		Tag:        "18.2b.",
		Annotation: "resources-role",
		ReadOnly:   true,
	}
	resourcesRole2ID21 = &xscmsg.FieldID{
		Tag:        "18.2b.",
		Annotation: "resources-role",
		Label:      "Resource 2 Role/Position",
	}
	preferredType2ID23 = &xscmsg.FieldID{
		Tag:        "18.2c.",
		Annotation: "preferred-type",
		Label:      "Resource 2 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	preferredType2ID21 = &xscmsg.FieldID{
		Tag:        "18.2c.",
		Annotation: "preferred-type",
		Label:      "Resource 2 Preferred Type",
	}
	minimumType2ID23 = &xscmsg.FieldID{
		Tag:        "18.2d.",
		Annotation: "minimum-type",
		Label:      "Resource 2 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	minimumType2ID21 = &xscmsg.FieldID{
		Tag:        "18.2d.",
		Annotation: "minimum-type",
		Label:      "Resource 2 Minimum Type",
	}
	resourcesQty3ID = &xscmsg.FieldID{
		Tag:        "18.3a.",
		Annotation: "resources-qty",
		Label:      "Resource 3 Qty",
		Comment:    "cardinal-number",
	}
	role3ID = &xscmsg.FieldID{
		Tag:        "18.3e.",
		Annotation: "role",
		Label:      "Resource 3 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
	}
	position3ID = &xscmsg.FieldID{
		Tag:        "18.3f.",
		Annotation: "position",
		Label:      "Resource 3 Position (for example, Checkpoint)",
	}
	resourcesRole3ID23 = &xscmsg.FieldID{
		Tag:        "18.3b.",
		Annotation: "resources-role",
		ReadOnly:   true,
	}
	resourcesRole3ID21 = &xscmsg.FieldID{
		Tag:        "18.3b.",
		Annotation: "resources-role",
		Label:      "Resource 3 Role/Position",
	}
	preferredType3ID23 = &xscmsg.FieldID{
		Tag:        "18.3c.",
		Annotation: "preferred-type",
		Label:      "Resource 3 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	preferredType3ID21 = &xscmsg.FieldID{
		Tag:        "18.3c.",
		Annotation: "preferred-type",
		Label:      "Resource 3 Preferred Type",
	}
	minimumType3ID23 = &xscmsg.FieldID{
		Tag:        "18.3d.",
		Annotation: "minimum-type",
		Label:      "Resource 3 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	minimumType3ID21 = &xscmsg.FieldID{
		Tag:        "18.3d.",
		Annotation: "minimum-type",
		Label:      "Resource 3 Minimum Type",
	}
	resourcesQty4ID = &xscmsg.FieldID{
		Tag:        "18.4a.",
		Annotation: "resources-qty",
		Label:      "Resource 4 Qty",
		Comment:    "cardinal-number",
	}
	role4ID = &xscmsg.FieldID{
		Tag:        "18.4e.",
		Annotation: "role",
		Label:      "Resource 4 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
	}
	position4ID = &xscmsg.FieldID{
		Tag:        "18.4f.",
		Annotation: "position",
		Label:      "Resource 4 Position (for example, Checkpoint)",
	}
	resourcesRole4ID23 = &xscmsg.FieldID{
		Tag:        "18.4b.",
		Annotation: "resources-role",
		ReadOnly:   true,
	}
	resourcesRole4ID21 = &xscmsg.FieldID{
		Tag:        "18.4b.",
		Annotation: "resources-role",
		Label:      "Resource 4 Role/Position",
	}
	preferredType4ID23 = &xscmsg.FieldID{
		Tag:        "18.4c.",
		Annotation: "preferred-type",
		Label:      "Resource 4 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	preferredType4ID21 = &xscmsg.FieldID{
		Tag:        "18.4c.",
		Annotation: "preferred-type",
		Label:      "Resource 4 Preferred Type",
	}
	minimumType4ID23 = &xscmsg.FieldID{
		Tag:        "18.4d.",
		Annotation: "minimum-type",
		Label:      "Resource 4 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	minimumType4ID21 = &xscmsg.FieldID{
		Tag:        "18.4d.",
		Annotation: "minimum-type",
		Label:      "Resource 4 Minimum Type",
	}
	resourcesQty5ID = &xscmsg.FieldID{
		Tag:        "18.5a.",
		Annotation: "resources-qty",
		Label:      "Resource 5 Qty",
		Comment:    "cardinal-number",
	}
	role5ID = &xscmsg.FieldID{
		Tag:        "18.5e.",
		Annotation: "role",
		Label:      "Resource 5 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
	}
	position5ID = &xscmsg.FieldID{
		Tag:        "18.5f.",
		Annotation: "position",
		Label:      "Resource 5 Position (for example, Checkpoint)",
	}
	resourcesRole5ID23 = &xscmsg.FieldID{
		Tag:        "18.5b.",
		Annotation: "resources-role",
		ReadOnly:   true,
	}
	resourcesRole5ID21 = &xscmsg.FieldID{
		Tag:        "18.5b.",
		Annotation: "resources-role",
		Label:      "Resource 5 Role/Position",
	}
	preferredType5ID23 = &xscmsg.FieldID{
		Tag:        "18.5c.",
		Annotation: "preferred-type",
		Label:      "Resource 5 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	preferredType5ID21 = &xscmsg.FieldID{
		Tag:        "18.5c.",
		Annotation: "preferred-type",
		Label:      "Resource 5 Preferred Type",
	}
	minimumType5ID23 = &xscmsg.FieldID{
		Tag:        "18.5d.",
		Annotation: "minimum-type",
		Label:      "Resource 5 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
	}
	minimumType5ID21 = &xscmsg.FieldID{
		Tag:        "18.5d.",
		Annotation: "minimum-type",
		Label:      "Resource 5 Minimum Type",
	}
	arrivalDatesID = &xscmsg.FieldID{
		Tag:        "19a.",
		Annotation: "arrival-dates",
		Label:      "Requested Arrival Date(s)",
		Comment:    "required",
	}
	arrivalTimesID = &xscmsg.FieldID{
		Tag:        "19b.",
		Annotation: "arrival-times",
		Label:      "Requested Arrival Time(s)",
		Comment:    "required",
	}
	neededDatesID = &xscmsg.FieldID{
		Tag:        "20a.",
		Annotation: "needed-dates",
		Label:      "Needed Until Date(s)",
		Comment:    "required",
	}
	neededTimesID = &xscmsg.FieldID{
		Tag:        "20b.",
		Annotation: "needed-times",
		Label:      "Needed Until Time(s)",
		Comment:    "required",
	}
	reportingLocationID = &xscmsg.FieldID{
		Tag:        "21.",
		Annotation: "reporting-location",
		Label:      "Reporting Location (Street Address, Parking, Entry Instructions)",
		Comment:    "required",
	}
	contactOnArrivalID = &xscmsg.FieldID{
		Tag:        "22.",
		Annotation: "contact-on-arrival",
		Label:      "Contact on Arrival (Name/Position and contact info)",
		Comment:    "required",
	}
	travelInfoID = &xscmsg.FieldID{
		Tag:        "23.",
		Annotation: "travel-info",
		Label:      "Travel Info (Routes, Hazards, Lodging)",
		Comment:    "required",
	}
	requesterNameID = &xscmsg.FieldID{
		Tag:        "24a.",
		Annotation: "requester-name",
		Label:      "Requested By Name",
		Comment:    "required",
	}
	requesterTitleID = &xscmsg.FieldID{
		Tag:        "24b.",
		Annotation: "requester-title",
		Label:      "Requested By Title",
		Comment:    "required",
	}
	requesterContactID = &xscmsg.FieldID{
		Tag:        "24c.",
		Annotation: "requester-contact",
		Label:      "Requested By Contact (E-mail, phone, frequency)",
		Comment:    "required",
	}
	agencyApproverNameID = &xscmsg.FieldID{
		Tag:        "25a.",
		Annotation: "agency-approver-name",
		Label:      "Approved By Name",
		Comment:    "required",
	}
	agencyApproverTitleID = &xscmsg.FieldID{
		Tag:        "25b.",
		Annotation: "agency-approver-title",
		Label:      "Approved By Title",
		Comment:    "required",
	}
	agencyApproverContactID = &xscmsg.FieldID{
		Tag:        "25c.",
		Annotation: "agency-approver-contact",
		Label:      "Approved By Contact (E-mail, phone, frequency)",
		Comment:    "required",
	}
	agencyApprovedDateID = &xscmsg.FieldID{
		Tag:        "26a.",
		Annotation: "agency-approved-date",
		Label:      "Approved By Date",
		Comment:    "required date",
	}
	agencyApprovedTimeID = &xscmsg.FieldID{
		Tag:        "26b.",
		Annotation: "agency-approved-time",
		Label:      "Approved By Time",
		Comment:    "required time",
	}
)

type handlingField struct{ xscform.ChoicesField }

func (f *handlingField) Default() string { return "ROUTINE" }

type toICSPositionField struct{ xscform.Field }

func (f *toICSPositionField) Default() string { return "RACES Chief Radio Officer" }

type toLocationField struct{ xscform.Field }

func (f *toLocationField) Default() string { return "County EOC" }

type resourcesRoleField struct{ xscform.Field }

func (f *resourcesRoleField) Validate(msg xscmsg.Message, strict bool) string {
	var btag = f.ID().Tag
	var etag = strings.Replace(btag, "b", "e", 1)
	var ftag = strings.Replace(btag, "b", "f", 1)
	var eval = msg.Field(etag).Get()
	var fval = msg.Field(ftag).Get()
	var bval = eval
	if fval != "" {
		bval += " / " + fval
	}
	f.Set(fval)
	return f.Field.Validate(msg, strict)
}

type typeField struct{ xscform.Field }

func (f *typeField) Validate(msg xscmsg.Message, strict bool) string {
	if err := f.Field.Validate(msg, strict); err != "" {
		return err
	}
	var tag = f.ID().Tag
	var etag = tag[:len(tag)-2] + "e."
	var eval = msg.Field(etag).Get()
	if eval == "" {
		if f.Get() != "" {
			return fmt.Sprintf("field %q must not have a value unless field %q has a value", tag, etag)
		}
		return ""
	}
	var val = f.Get()
	if val == "" || val == "Type IV" || val == "Type V" {
		return ""
	}
	if !strict {
		if strings.EqualFold(val, "Type IV") {
			f.Set("Type IV")
			return ""
		}
		if strings.EqualFold(val, "Type V") {
			f.Set("Type V")
			return ""
		}
	}
	if len(val) != 2 || (val[1] != '1' && val[1] != '2' && val[1] != '3') {
		return fmt.Sprintf("%q is not a valid value for field %q", val, tag)
	}
	tcode := val[:1]
	if !strict {
		tcode = strings.ToUpper(tcode)
	}
	if tcode != "F" && tcode != "N" && tcode != "P" && tcode != "S" {
		return fmt.Sprintf("%q is not a valid value for field %q", val, tag)
	}
	if tcode != eval[:1] {
		return fmt.Sprintf("the value %q for field %q is not compatible with the value %q for field %q", val, tag, eval, etag)
	}
	f.Set(tcode + val[1:2])
	return ""
}
