package racesmar

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies RACES MAR forms.
const (
	Tag = "RACES-MAR"
)

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)

	// Our handling, toICSPosition, and toLocation fields are variants of
	// the standard ones, adding default values to them.
	handlingDef = *xscform.HandlingDef
	handlingDef.DefaultValue = "ROUTINE"
	toICSPositionDef = *xscform.ToICSPositionDef
	toICSPositionDef.DefaultValue = "RACES Chief Radio Officer"
	toICSPositionDef.Comment = "required: RACES Chief Radio Officer, RACES Unit, Operations Section, ..."
	toLocationDef = *xscform.ToLocationDef
	toLocationDef.DefaultValue = "County EOC"
	toLocationDef.Comment = "required: County EOC, ..."
}

func create() *xscmsg.Message {
	return xscform.CreateForm(formtype23, fieldDefsV23)
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) *xscmsg.Message {
	if form == nil {
		return nil
	}
	if form.FormType == formtype23.HTML && !xscmsg.OlderVersion(form.FormVersion, "2.3") {
		return xscform.AdoptForm(formtype23, fieldDefsV23, msg, form)
	}
	if form.FormType == formtype23.HTML && !xscmsg.OlderVersion(form.FormVersion, "2.1") {
		return xscform.AdoptForm(formtype23, fieldDefsV21, msg, form)
	}
	if form.FormType == formtype16.HTML && !xscmsg.OlderVersion(form.FormVersion, "1.6") {
		return xscform.AdoptForm(formtype16, fieldDefsV16, msg, form)
	}
	return nil
}

var formtype23 = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "RACES mutual aid request form",
	Article:     "a",
	HTML:        "form-oa-mutual-aid-request-v2.html",
	Version:     "2.3",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

var formtype16 = &xscmsg.MessageType{
	Tag:         Tag,
	Name:        "RACES mutual aid request form",
	Article:     "a",
	HTML:        "form-oa-mutual-aid-request.html",
	Version:     "1.6",
	SubjectFunc: xscform.EncodeSubject,
	BodyFunc:    xscform.EncodeBody,
}

// fieldDefsV23 replaces the single resourcesRole field for each resource with
// separate role and position fields, plus a computed resourcesRole field.  It
// also adds validation to the preferredType and mininumType fields.
var fieldDefsV23 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// RACES MAR fields
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQty1Def, role1Def, position1Def, resourcesRole1DefV23, preferredType1DefV23, minimumType1DefV23,
	resourcesQty2Def, role2Def, position2Def, resourcesRole2DefV23, preferredType2DefV23, minimumType2DefV23,
	resourcesQty3Def, role3Def, position3Def, resourcesRole3DefV23, preferredType3DefV23, minimumType3DefV23,
	resourcesQty4Def, role4Def, position4Def, resourcesRole4DefV23, preferredType4DefV23, minimumType4DefV23,
	resourcesQty5Def, role5Def, position5Def, resourcesRole5DefV23, preferredType5DefV23, minimumType5DefV23,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

// fieldDefsV21 replaces a single set of resource definition fields with five
// sets of fields.
var fieldDefsV21 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// RACES MAR fields
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQtyDef, resourcesRoleDef, preferredTypeDef, minimumTypeDef,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQty1Def, resourcesRole1DefV21, preferredType1DefV21, minimumType1DefV21,
	resourcesQty2Def, resourcesRole2DefV21, preferredType2DefV21, minimumType2DefV21,
	resourcesQty3Def, resourcesRole3DefV21, preferredType3DefV21, minimumType3DefV21,
	resourcesQty4Def, resourcesRole4DefV21, preferredType4DefV21, minimumType4DefV21,
	resourcesQty5Def, resourcesRole5DefV21, preferredType5DefV21, minimumType5DefV21,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var fieldDefsV16 = []*xscmsg.FieldDef{
	// Standard header
	xscform.OriginMessageNumberDef, xscform.DestinationMessageNumberDef, xscform.MessageDateDef, xscform.MessageTimeDef,
	&handlingDef, &toICSPositionDef, xscform.FromICSPositionDef, &toLocationDef, xscform.FromLocationDef, xscform.ToNameDef,
	xscform.FromNameDef, xscform.ToContactDef, xscform.FromContactDef,
	// RACES MAR fields
	agencyDef, eventNameDef, eventNumberDef, assignmentDef,
	resourcesQtyDef, resourcesRoleDef, preferredTypeDef, minimumTypeDef,
	arrivalDatesDef, arrivalTimesDef, neededDatesDef, neededTimesDef, reportingLocationDef, contactOnArrivalDef,
	travelInfoDef, requesterNameDef, requesterTitleDef, requesterContactDef, agencyApproverNameDef, agencyApproverTitleDef,
	agencyApproverContactDef, agencyApprovedDateDef, agencyApprovedTimeDef,
	// Standard footer
	xscform.OpRelayRcvdDef, xscform.OpRelaySentDef, xscform.OpNameDef, xscform.OpCallDef, xscform.OpDateDef, xscform.OpTimeDef,
}

var (
	handlingDef      = *xscform.HandlingDef      // modified in func init
	toICSPositionDef = *xscform.ToICSPositionDef // modified in func init
	toLocationDef    = *xscform.ToLocationDef    // modified in func init
	agencyDef        = &xscmsg.FieldDef{
		Tag:        "15.",
		Annotation: "agency",
		Label:      "Agency Name",
		Comment:    "required",
		Key:        xscmsg.FSubject,
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	eventNameDef = &xscmsg.FieldDef{
		Tag:        "16a.",
		Annotation: "event-name",
		Label:      "Event / Incident Name",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	eventNumberDef = &xscmsg.FieldDef{
		Tag:        "16b.",
		Annotation: "event-number",
		Label:      "Event / Incident Number",
		Validators: []xscmsg.Validator{},
	}
	assignmentDef = &xscmsg.FieldDef{
		Tag:        "17.",
		Annotation: "assignment",
		Label:      "Assignment",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	resourcesQtyDef = &xscmsg.FieldDef{
		Tag:        "18a.",
		Annotation: "resources-qty",
		Label:      "Qty",
		Comment:    "required cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateCardinalNumber},
	}
	resourcesRoleDef = &xscmsg.FieldDef{
		Tag:        "18b.",
		Annotation: "resources-role",
		Label:      "Role/Position",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	preferredTypeDef = &xscmsg.FieldDef{
		Tag:        "18c.",
		Annotation: "preferred-type",
		Label:      "Preferred Type",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	minimumTypeDef = &xscmsg.FieldDef{
		Tag:        "18d.",
		Annotation: "minimum-type",
		Label:      "Minimum Type",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	resourcesQty1Def = &xscmsg.FieldDef{
		Tag:        "18.1a.",
		Annotation: "resources-qty",
		Label:      "Resource 1 Qty",
		Comment:    "required cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateCardinalNumber},
	}
	role1Def = &xscmsg.FieldDef{
		Tag:        "18.1e.",
		Annotation: "role",
		Label:      "Resource 1 Role",
		Comment:    "required: Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position1Def = &xscmsg.FieldDef{
		Tag:        "18.1f.",
		Annotation: "position",
		Label:      "Resource 1 Position (for example, Checkpoint)",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole1DefV23 = &xscmsg.FieldDef{
		Tag:        "18.1b.",
		Annotation: "resources-role",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setResourcesRole, xscform.ValidateRequired},
	}
	resourcesRole1DefV21 = &xscmsg.FieldDef{
		Tag:        "18.1b.",
		Annotation: "resources-role",
		Label:      "Resource 1 Role/Position",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	preferredType1DefV23 = &xscmsg.FieldDef{
		Tag:        "18.1c.",
		Annotation: "preferred-type",
		Label:      "Resource 1 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType1DefV21 = &xscmsg.FieldDef{
		Tag:        "18.1c.",
		Annotation: "preferred-type",
		Label:      "Resource 1 Preferred Type",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	minimumType1DefV23 = &xscmsg.FieldDef{
		Tag:        "18.1d.",
		Annotation: "minimum-type",
		Label:      "Resource 1 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType1DefV21 = &xscmsg.FieldDef{
		Tag:        "18.1d.",
		Annotation: "minimum-type",
		Label:      "Resource 1 Minimum Type",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	resourcesQty2Def = &xscmsg.FieldDef{
		Tag:        "18.2a.",
		Annotation: "resources-qty",
		Label:      "Resource 2 Qty",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role2Def = &xscmsg.FieldDef{
		Tag:        "18.2e.",
		Annotation: "role",
		Label:      "Resource 2 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position2Def = &xscmsg.FieldDef{
		Tag:        "18.2f.",
		Annotation: "position",
		Label:      "Resource 2 Position (for example, Checkpoint)",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole2DefV23 = &xscmsg.FieldDef{
		Tag:        "18.2b.",
		Annotation: "resources-role",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole2DefV21 = &xscmsg.FieldDef{
		Tag:        "18.2b.",
		Annotation: "resources-role",
		Label:      "Resource 2 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType2DefV23 = &xscmsg.FieldDef{
		Tag:        "18.2c.",
		Annotation: "preferred-type",
		Label:      "Resource 2 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType2DefV21 = &xscmsg.FieldDef{
		Tag:        "18.2c.",
		Annotation: "preferred-type",
		Label:      "Resource 2 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType2DefV23 = &xscmsg.FieldDef{
		Tag:        "18.2d.",
		Annotation: "minimum-type",
		Label:      "Resource 2 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType2DefV21 = &xscmsg.FieldDef{
		Tag:        "18.2d.",
		Annotation: "minimum-type",
		Label:      "Resource 2 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	resourcesQty3Def = &xscmsg.FieldDef{
		Tag:        "18.3a.",
		Annotation: "resources-qty",
		Label:      "Resource 3 Qty",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role3Def = &xscmsg.FieldDef{
		Tag:        "18.3e.",
		Annotation: "role",
		Label:      "Resource 3 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position3Def = &xscmsg.FieldDef{
		Tag:        "18.3f.",
		Annotation: "position",
		Label:      "Resource 3 Position (for example, Checkpoint)",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole3DefV23 = &xscmsg.FieldDef{
		Tag:        "18.3b.",
		Annotation: "resources-role",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole3DefV21 = &xscmsg.FieldDef{
		Tag:        "18.3b.",
		Annotation: "resources-role",
		Label:      "Resource 3 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType3DefV23 = &xscmsg.FieldDef{
		Tag:        "18.3c.",
		Annotation: "preferred-type",
		Label:      "Resource 3 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType3DefV21 = &xscmsg.FieldDef{
		Tag:        "18.3c.",
		Annotation: "preferred-type",
		Label:      "Resource 3 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType3DefV23 = &xscmsg.FieldDef{
		Tag:        "18.3d.",
		Annotation: "minimum-type",
		Label:      "Resource 3 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType3DefV21 = &xscmsg.FieldDef{
		Tag:        "18.3d.",
		Annotation: "minimum-type",
		Label:      "Resource 3 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	resourcesQty4Def = &xscmsg.FieldDef{
		Tag:        "18.4a.",
		Annotation: "resources-qty",
		Label:      "Resource 4 Qty",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role4Def = &xscmsg.FieldDef{
		Tag:        "18.4e.",
		Annotation: "role",
		Label:      "Resource 4 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position4Def = &xscmsg.FieldDef{
		Tag:        "18.4f.",
		Annotation: "position",
		Label:      "Resource 4 Position (for example, Checkpoint)",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole4DefV23 = &xscmsg.FieldDef{
		Tag:        "18.4b.",
		Annotation: "resources-role",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole4DefV21 = &xscmsg.FieldDef{
		Tag:        "18.4b.",
		Annotation: "resources-role",
		Label:      "Resource 4 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType4DefV23 = &xscmsg.FieldDef{
		Tag:        "18.4c.",
		Annotation: "preferred-type",
		Label:      "Resource 4 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType4DefV21 = &xscmsg.FieldDef{
		Tag:        "18.4c.",
		Annotation: "preferred-type",
		Label:      "Resource 4 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType4DefV23 = &xscmsg.FieldDef{
		Tag:        "18.4d.",
		Annotation: "minimum-type",
		Label:      "Resource 4 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType4DefV21 = &xscmsg.FieldDef{
		Tag:        "18.4d.",
		Annotation: "minimum-type",
		Label:      "Resource 4 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	resourcesQty5Def = &xscmsg.FieldDef{
		Tag:        "18.5a.",
		Annotation: "resources-qty",
		Label:      "Resource 5 Qty",
		Comment:    "cardinal-number",
		Validators: []xscmsg.Validator{xscform.ValidateCardinalNumber},
	}
	role5Def = &xscmsg.FieldDef{
		Tag:        "18.5e.",
		Annotation: "role",
		Label:      "Resource 5 Role",
		Comment:    "Field Communicator, Net Control Operator, Packet Operator, Shadow Communicator",
		Validators: []xscmsg.Validator{xscform.ValidateChoices},
		Choices:    []string{"Field Communicator", "Net Control Operator", "Packet Operator", "Shadow Communicator"},
	}
	position5Def = &xscmsg.FieldDef{
		Tag:        "18.5f.",
		Annotation: "position",
		Label:      "Resource 5 Position (for example, Checkpoint)",
		Validators: []xscmsg.Validator{},
	}
	resourcesRole5DefV23 = &xscmsg.FieldDef{
		Tag:        "18.5b.",
		Annotation: "resources-role",
		ReadOnly:   true,
		Validators: []xscmsg.Validator{setResourcesRole},
	}
	resourcesRole5DefV21 = &xscmsg.FieldDef{
		Tag:        "18.5b.",
		Annotation: "resources-role",
		Label:      "Resource 5 Role/Position",
		Validators: []xscmsg.Validator{},
	}
	preferredType5DefV23 = &xscmsg.FieldDef{
		Tag:        "18.5c.",
		Annotation: "preferred-type",
		Label:      "Resource 5 Preferred Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	preferredType5DefV21 = &xscmsg.FieldDef{
		Tag:        "18.5c.",
		Annotation: "preferred-type",
		Label:      "Resource 5 Preferred Type",
		Validators: []xscmsg.Validator{},
	}
	minimumType5DefV23 = &xscmsg.FieldDef{
		Tag:        "18.5d.",
		Annotation: "minimum-type",
		Label:      "Resource 5 Minimum Type",
		Comment:    "[FNPS][123], Type IV, Type V",
		Validators: []xscmsg.Validator{validateType},
	}
	minimumType5DefV21 = &xscmsg.FieldDef{
		Tag:        "18.5d.",
		Annotation: "minimum-type",
		Label:      "Resource 5 Minimum Type",
		Validators: []xscmsg.Validator{},
	}
	arrivalDatesDef = &xscmsg.FieldDef{
		Tag:        "19a.",
		Annotation: "arrival-dates",
		Label:      "Requested Arrival Date(s)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	arrivalTimesDef = &xscmsg.FieldDef{
		Tag:        "19b.",
		Annotation: "arrival-times",
		Label:      "Requested Arrival Time(s)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	neededDatesDef = &xscmsg.FieldDef{
		Tag:        "20a.",
		Annotation: "needed-dates",
		Label:      "Needed Until Date(s)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	neededTimesDef = &xscmsg.FieldDef{
		Tag:        "20b.",
		Annotation: "needed-times",
		Label:      "Needed Until Time(s)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	reportingLocationDef = &xscmsg.FieldDef{
		Tag:        "21.",
		Annotation: "reporting-location",
		Label:      "Reporting Location (Street Address, Parking, Entry Instructions)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	contactOnArrivalDef = &xscmsg.FieldDef{
		Tag:        "22.",
		Annotation: "contact-on-arrival",
		Label:      "Contact on Arrival (Name/Position and contact info)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	travelInfoDef = &xscmsg.FieldDef{
		Tag:        "23.",
		Annotation: "travel-info",
		Label:      "Travel Info (Routes, Hazards, Lodging)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	requesterNameDef = &xscmsg.FieldDef{
		Tag:        "24a.",
		Annotation: "requester-name",
		Label:      "Requested By Name",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	requesterTitleDef = &xscmsg.FieldDef{
		Tag:        "24b.",
		Annotation: "requester-title",
		Label:      "Requested By Title",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	requesterContactDef = &xscmsg.FieldDef{
		Tag:        "24c.",
		Annotation: "requester-contact",
		Label:      "Requested By Contact (E-mail, phone, frequency)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	agencyApproverNameDef = &xscmsg.FieldDef{
		Tag:        "25a.",
		Annotation: "agency-approver-name",
		Label:      "Approved By Name",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	agencyApproverTitleDef = &xscmsg.FieldDef{
		Tag:        "25b.",
		Annotation: "agency-approver-title",
		Label:      "Approved By Title",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	agencyApproverContactDef = &xscmsg.FieldDef{
		Tag:        "25c.",
		Annotation: "agency-approver-contact",
		Label:      "Approved By Contact (E-mail, phone, frequency)",
		Comment:    "required",
		Validators: []xscmsg.Validator{xscform.ValidateRequired},
	}
	agencyApprovedDateDef = &xscmsg.FieldDef{
		Tag:        "26a.",
		Annotation: "agency-approved-date",
		Label:      "Approved By Date",
		Comment:    "required date",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateDate},
	}
	agencyApprovedTimeDef = &xscmsg.FieldDef{
		Tag:        "26b.",
		Annotation: "agency-approved-time",
		Label:      "Approved By Time",
		Comment:    "required time",
		Validators: []xscmsg.Validator{xscform.ValidateRequired, xscform.ValidateTime},
	}
)

func setResourcesRole(f *xscmsg.Field, msg *xscmsg.Message, strict bool) string {
	var btag = f.Def.Tag
	var etag = strings.Replace(btag, "b", "e", 1)
	var ftag = strings.Replace(btag, "b", "f", 1)
	var eval = msg.Field(etag).Value
	var fval = msg.Field(ftag).Value
	var bval = eval
	if fval != "" {
		bval += " / " + fval
	}
	if !strict {
		f.Value = bval
	} else if f.Value != bval {
		return fmt.Sprintf("The value of field %q is not consistent with the values of fields %q and %q.", btag, etag, ftag)
	}
	return ""
}

func validateType(f *xscmsg.Field, msg *xscmsg.Message, strict bool) string {
	var tag = f.Def.Tag
	var etag = tag[:len(tag)-2] + "e."
	var eval = msg.Field(etag).Value
	if eval == "" {
		if f.Value != "" {
			return fmt.Sprintf("Field %q must not have a value unless field %q has a value.", tag, etag)
		}
		return ""
	}
	var val = f.Value
	if val == "" || val == "Type IV" || val == "Type V" {
		return ""
	}
	if !strict {
		if strings.EqualFold(val, "Type IV") {
			f.Value = "Type IV"
			return ""
		}
		if strings.EqualFold(val, "Type V") {
			f.Value = "Type V"
			return ""
		}
	}
	if len(val) != 2 || (val[1] != '1' && val[1] != '2' && val[1] != '3') {
		return fmt.Sprintf("%q is not a valid value for field %q.", val, tag)
	}
	tcode := val[:1]
	if !strict {
		tcode = strings.ToUpper(tcode)
	}
	if tcode != "F" && tcode != "N" && tcode != "P" && tcode != "S" {
		return fmt.Sprintf("%q is not a valid value for field %q.", val, tag)
	}
	if tcode != eval[:1] {
		return fmt.Sprintf("The value %q for field %q is not compatible with the value %q for field %q.", val, tag, eval, etag)
	}
	f.Value = tcode + val[1:2]
	return ""
}
