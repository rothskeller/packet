package xscform

import (
	"github.com/rothskeller/packet/xscmsg"
)

// OriginMessageNumberDef is the definition of the XSC-standard Origin Message
// Number field.
var OriginMessageNumberDef = &xscmsg.FieldDef{
	Tag:        "MsgNo",
	Label:      "Origin Message Number",
	Comment:    "required message-number",
	Key:        xscmsg.FOriginMsgNo,
	Validators: []xscmsg.Validator{ValidateRequired, ValidateMessageNumber},
}

// DestinationMessageNumberDef is the definition of the XSC-standard Destination
// Message Number field.
var DestinationMessageNumberDef = &xscmsg.FieldDef{
	Tag:        "DestMsgNo",
	Label:      "Destination Message Number",
	Comment:    "message-number",
	Key:        xscmsg.FDestinationMsgNo,
	ReadOnly:   true,
	Validators: []xscmsg.Validator{ValidateMessageNumber},
}

// MessageDateDef is the definition of the XSC-standard Message Date field.
var MessageDateDef = &xscmsg.FieldDef{
	Tag:         "1a.",
	Annotation:  "date",
	Label:       "Date",
	Comment:     "required date",
	Validators:  []xscmsg.Validator{ValidateRequired, ValidateDate},
	DefaultFunc: DefaultDate,
}

// MessageTimeDef is the definition of the XSC-standard Message Time field.
var MessageTimeDef = &xscmsg.FieldDef{
	Tag:         "1b.",
	Annotation:  "time",
	Label:       "Time",
	Comment:     "required time",
	Validators:  []xscmsg.Validator{ValidateRequired, ValidateTime},
	DefaultFunc: DefaultTime,
}

// HandlingDef is the definition of the XSC-standard Handling field.
var HandlingDef = &xscmsg.FieldDef{
	Tag:        "5.",
	Annotation: "handling",
	Label:      "Handling",
	Comment:    "required: IMMEDIATE, PRIORITY, ROUTINE",
	Key:        xscmsg.FHandling,
	Validators: []xscmsg.Validator{ValidateRequired, ValidateChoices},
	Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
}

// ToICSPositionDef is the definition of the XSC-standard To ICS Position field.
var ToICSPositionDef = &xscmsg.FieldDef{
	Tag:        "7a.",
	Annotation: "to-ics-position",
	Label:      "To ICS Position",
	Comment:    "required",
	Key:        xscmsg.FToICSPosition,
	Validators: []xscmsg.Validator{ValidateRequired},
}

// FromICSPositionDef is the definition of the XSC-standard From ICS Position
// field.
var FromICSPositionDef = &xscmsg.FieldDef{
	Tag:        "8a.",
	Annotation: "from-ics-position",
	Label:      "From ICS Position",
	Comment:    "required",
	Validators: []xscmsg.Validator{ValidateRequired},
}

// ToLocationDef is the definition of the XSC-standard To Location field.
var ToLocationDef = &xscmsg.FieldDef{
	Tag:        "7b.",
	Annotation: "to-location",
	Label:      "To Location",
	Comment:    "required",
	Key:        xscmsg.FToLocation,
	Validators: []xscmsg.Validator{ValidateRequired},
}

// FromLocationDef is the definition of the XSC-standard From Location field.
var FromLocationDef = &xscmsg.FieldDef{
	Tag:        "8b.",
	Annotation: "from-location",
	Label:      "From Location",
	Comment:    "required",
	Validators: []xscmsg.Validator{ValidateRequired},
}

// ToNameDef is the definition of the XSC-standard To Name field.
var ToNameDef = &xscmsg.FieldDef{
	Tag:        "7c.",
	Annotation: "to-name",
	Label:      "To Name",
}

// FromNameDef is the definition of the XSC-standard From Name field.
var FromNameDef = &xscmsg.FieldDef{
	Tag:        "8c.",
	Annotation: "from-name",
	Label:      "From Name",
}

// ToContactDef is the definition of the XSC-standard To Contact field.
var ToContactDef = &xscmsg.FieldDef{
	Tag:        "7d.",
	Annotation: "to-contact",
	Label:      "To Contact Info",
}

// FromContactDef is the definition of the XSC-standard From Contact field.
var FromContactDef = &xscmsg.FieldDef{
	Tag:        "8d.",
	Annotation: "from-contact",
	Label:      "From Contact Info",
}

// OpRelayRcvdDef is the definition of the XSC-standard Operator Relay Received
// field.
var OpRelayRcvdDef = &xscmsg.FieldDef{
	Tag:   "OpRelayRcvd",
	Label: "Relay Rcvd",
}

// OpRelaySentDef is the definition of the XSC-standard Operator Relay Sent
// field.
var OpRelaySentDef = &xscmsg.FieldDef{
	Tag:   "OpRelaySent",
	Label: "Relay Sent",
}

// OpNameDef is the definition of the XSC-standard Operator Name field.
var OpNameDef = &xscmsg.FieldDef{
	Tag:   "OpName",
	Label: "Operator Name",
	Key:   xscmsg.FOpName,
}

// OpCallDef is the definition of the XSC-standard Operator Call field.
var OpCallDef = &xscmsg.FieldDef{
	Tag:        "OpCall",
	Label:      "Operator Call Sign",
	Comment:    "required call-sign",
	Key:        xscmsg.FOpCall,
	Validators: []xscmsg.Validator{ValidateRequired, ValidateCallSign},
}

// OpDateDef is the definition of the XSC-standard Operator Date field.
var OpDateDef = &xscmsg.FieldDef{
	Tag:         "OpDate",
	Label:       "Operator Date",
	Comment:     "date",
	Key:         xscmsg.FOpDate,
	Validators:  []xscmsg.Validator{ValidateDate},
	DefaultFunc: DefaultDate,
}

// OpTimeDef is the definition of the XSC-standard Operator Time field.
var OpTimeDef = &xscmsg.FieldDef{
	Tag:         "OpTime",
	Label:       "Operator Time",
	Comment:     "time",
	Key:         xscmsg.FOpTime,
	Validators:  []xscmsg.Validator{ValidateTime},
	DefaultFunc: DefaultTime,
}
