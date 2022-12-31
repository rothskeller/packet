package eoc213rr

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/xscmsg"
	"github.com/rothskeller/packet/xscmsg/internal/xscform"
)

// Tag identifies EOC-213RR forms.
const Tag = "EOC213RR"

func init() {
	xscmsg.RegisterCreate(Tag, create)
	xscmsg.RegisterType(recognize)
}

func create() xscmsg.Message {
	return xscform.CreateForm(formtype, makeFields())
}

func recognize(msg *pktmsg.Message, form *pktmsg.Form) xscmsg.Message {
	if form == nil || form.FormType != formtype.HTML || xscmsg.OlderVersion(form.FormVersion, "2.0") {
		return nil
	}
	return xscform.AdoptForm(formtype, makeFields(), msg, form)
}

var formtype = &xscmsg.MessageType{
	Tag:     Tag,
	Name:    "EOC-213RR resource request form",
	Article: "an",
	HTML:    "form-scco-eoc-213rr.html",
	Version: "2.3",
}

func makeFields() []xscmsg.Field {
	return []xscmsg.Field{
		xscform.FOriginMessageNumber(),
		xscform.FDestinationMessageNumber(),
		xscform.FMessageDate(),
		xscform.FMessageTime(),
		xscform.FHandling(),
		&toICSPositionField{*xscform.FToICSPosition().(*xscform.Field)},
		xscform.FFromICSPosition(),
		&toLocationField{*xscform.FToLocation().(*xscform.Field)},
		xscform.FFromLocation(),
		xscform.FToName(),
		xscform.FFromName(),
		xscform.FToContact(),
		xscform.FFromContact(),
		xscform.NewField(incidentNameID, true),
		&xscform.DateField{Field: *xscform.NewField(dateInitiatedID, true)},
		&xscform.TimeField{Field: *xscform.NewField(timeInitiatedID, true)},
		xscform.NewField(trackingNumberID, false),
		xscform.NewField(requestedByID, true),
		xscform.NewField(preparedByID, false),
		xscform.NewField(approvedByID, false),
		xscform.NewField(qtyUnitID, true),
		xscform.NewField(resourceDescriptionID, true),
		xscform.NewField(resourceArrivalID, true),
		&xscform.ChoicesField{Field: *xscform.NewField(priorityID, true), Choices: priorityChoices},
		xscform.NewField(estdCostID, false),
		xscform.NewField(deliverToID, true),
		xscform.NewField(deliverToLocationID, true),
		xscform.NewField(substitutesID, false),
		&xscform.BooleanField{Field: *xscform.NewField(equipmentOperatorID, false)},
		&xscform.BooleanField{Field: *xscform.NewField(lodgingID, false)},
		&xscform.BooleanField{Field: *xscform.NewField(fuelID, false)},
		xscform.NewField(fuelTypeID, false),
		&xscform.BooleanField{Field: *xscform.NewField(powerID, false)},
		&xscform.BooleanField{Field: *xscform.NewField(mealsID, false)},
		&xscform.BooleanField{Field: *xscform.NewField(maintenanceID, false)},
		&xscform.BooleanField{Field: *xscform.NewField(waterID, false)},
		&xscform.BooleanField{Field: *xscform.NewField(otherID, false)},
		xscform.NewField(instructionsID, false),
		xscform.FOpRelayRcvd(),
		xscform.FOpRelaySent(),
		xscform.FOpName(),
		xscform.FOpCall(),
		xscform.FOpDate(),
		xscform.FOpTime(),
	}
}

var (
	incidentNameID = &xscmsg.FieldID{
		Tag:        "21.",
		Annotation: "incident-name",
		Label:      "1. Incident Name",
		Comment:    "required",
		Canonical:  xscmsg.FSubject,
	}
	dateInitiatedID = &xscmsg.FieldID{
		Tag:        "22.",
		Annotation: "date",
		Label:      "2. Date Initiated",
		Comment:    "required date",
	}
	timeInitiatedID = &xscmsg.FieldID{
		Tag:        "23.",
		Annotation: "time",
		Label:      "3. Time Initiated",
		Comment:    "required time",
	}
	trackingNumberID = &xscmsg.FieldID{
		Tag:        "24.",
		Annotation: "tracking-number",
		Label:      "4. Tracking Number (OA EOC)",
	}
	requestedByID = &xscmsg.FieldID{
		Tag:        "25.",
		Annotation: "requested-by",
		Label:      "5. Requested by (name, agency, position, email, phone)",
		Comment:    "required",
	}
	preparedByID = &xscmsg.FieldID{
		Tag:        "26.",
		Annotation: "prepared-by",
		Label:      "6. Prepared by (name, position, email, phone)",
	}
	approvedByID = &xscmsg.FieldID{
		Tag:        "27.",
		Annotation: "approved-by",
		Label:      "7. Approved By (name, position, email, phone)",
	}
	qtyUnitID = &xscmsg.FieldID{
		Tag:        "28.",
		Annotation: "qty-unit",
		Label:      "8. Qty/Unit",
		Comment:    "required",
	}
	resourceDescriptionID = &xscmsg.FieldID{
		Tag:        "29.",
		Annotation: "resource-description",
		Label:      "9. Resource Description (kind/type if applicable)",
		Comment:    "required",
	}
	resourceArrivalID = &xscmsg.FieldID{
		Tag:        "30.",
		Annotation: "resource-arrival",
		Label:      "10. Arrival (date/time)",
		Comment:    "required",
	}
	priorityID = &xscmsg.FieldID{
		Tag:        "31.",
		Annotation: "priority",
		Label:      "11. Priority",
		Comment:    "required: Now, High, Medium, Low",
	}
	priorityChoices = []string{"Now", "High", "Medium", "Low"}
	estdCostID      = &xscmsg.FieldID{
		Tag:        "32.",
		Annotation: "estd-cost", // it's resource-priority in PackItForms, but that's clearly wrong
		Label:      "12. Est'd Cost",
	}
	deliverToID = &xscmsg.FieldID{
		Tag:        "33.",
		Annotation: "deliver-to",
		Label:      "13. Deliver to (name, agency, position, email, phone)",
		Comment:    "required",
	}
	deliverToLocationID = &xscmsg.FieldID{
		Tag:        "34.",
		Annotation: "deliver-to-location",
		Label:      "14. Location (address or lat/long, site type)",
		Comment:    "required",
	}
	substitutesID = &xscmsg.FieldID{
		Tag:        "35.",
		Annotation: "substitutes",
		Label:      "15. Substitute/Suggested Sources (name, phone, website)",
	}
	equipmentOperatorID = &xscmsg.FieldID{
		Tag:        "36a.",
		Annotation: "equipment-operator",
		Label:      "16. Supplemental Requirements: Equipment Operator",
		Comment:    "boolean",
	}
	lodgingID = &xscmsg.FieldID{
		Tag:        "36b.",
		Annotation: "lodging",
		Label:      "16. Supplemental Requirements: Lodging",
		Comment:    "boolean",
	}
	fuelID = &xscmsg.FieldID{
		Tag:        "36c.",
		Annotation: "fuel",
		Label:      "16. Supplemental Requirements: Fuel",
		Comment:    "boolean",
	}
	fuelTypeID = &xscmsg.FieldID{
		Tag:        "36d.",
		Annotation: "fuel-type",
		Label:      "16. Supplemental Requirements: Fuel Type",
	}
	powerID = &xscmsg.FieldID{
		Tag:        "36e.",
		Annotation: "power",
		Label:      "16. Supplemental Requirements: Power",
		Comment:    "boolean",
	}
	mealsID = &xscmsg.FieldID{
		Tag:        "36f.",
		Annotation: "meals",
		Label:      "16. Supplemental Requirements: Meals",
		Comment:    "boolean",
	}
	maintenanceID = &xscmsg.FieldID{
		Tag:        "36g.",
		Annotation: "maintenance",
		Label:      "16. Supplemental Requirements: Maintenance",
		Comment:    "boolean",
	}
	waterID = &xscmsg.FieldID{
		Tag:        "36h.",
		Annotation: "water",
		Label:      "16. Supplemental Requirements: Water",
		Comment:    "boolean",
	}
	otherID = &xscmsg.FieldID{
		Tag:        "36i.",
		Annotation: "other",
		Label:      "16. Supplemental Requirements: Other",
		Comment:    "boolean",
	}
	instructionsID = &xscmsg.FieldID{
		Tag:        "37.",
		Annotation: "instructions",
		Label:      "17. Special Instructions",
	}
)

type toICSPositionField struct{ xscform.Field }

func (f *toICSPositionField) Default() string { return "Planning Section" }

type toLocationField struct{ xscform.Field }

func (f *toLocationField) Default() string { return "County EOC" }
