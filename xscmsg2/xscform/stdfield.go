package xscform

import (
	"github.com/rothskeller/packet/xscmsg"
)

// OriginMessageNumberDef is the definition of the XSC-standard Origin Message
// Number field.
var OriginMessageNumberDef = &xscmsg.FieldDef{
	Tag:        "MsgNo",
	Label:      "Origin Message Number",
	Key:        xscmsg.FOriginMsgNo,
	Validators: []xscmsg.Validator{ValidateMessageNumber},
	Flags:      xscmsg.Required,
}

// DestinationMessageNumberDef is the definition of the XSC-standard Destination
// Message Number field.
var DestinationMessageNumberDef = &xscmsg.FieldDef{
	Tag:        "DestMsgNo",
	Label:      "Destination Message Number",
	Key:        xscmsg.FDestinationMsgNo,
	Flags:      xscmsg.Readonly,
	Validators: []xscmsg.Validator{ValidateMessageNumber},
}

// MessageDateDef is the definition of the XSC-standard Message Date field.
var MessageDateDef = &xscmsg.FieldDef{
	Tag:         "1a.",
	Label:       "Date",
	Comment:     "MM/DD/YYYY",
	Validators:  []xscmsg.Validator{ValidateDate},
	Flags:       xscmsg.Required,
	DefaultFunc: DefaultDate,
}

// MessageTimeDef is the definition of the XSC-standard Message Time field.
var MessageTimeDef = &xscmsg.FieldDef{
	Tag:         "1b.",
	Label:       "Time",
	Comment:     "HH:MM",
	Validators:  []xscmsg.Validator{ValidateTime},
	Flags:       xscmsg.Required,
	DefaultFunc: DefaultTime,
}

// HandlingDef is the definition of the XSC-standard Handling field.
var HandlingDef = &xscmsg.FieldDef{
	Tag:        "5.",
	Label:      "Handling",
	Key:        xscmsg.FHandling,
	Validators: []xscmsg.Validator{ValidateChoices},
	Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
	Flags:      xscmsg.Required,
}

// ToICSPositionDef is the definition of the XSC-standard To ICS Position field.
var ToICSPositionDef = &xscmsg.FieldDef{
	Tag:   "7a.",
	Label: "To ICS Position",
	Key:   xscmsg.FToICSPosition,
	Flags: xscmsg.Required,
}

// FromICSPositionDef is the definition of the XSC-standard From ICS Position
// field.
var FromICSPositionDef = &xscmsg.FieldDef{
	Tag:   "8a.",
	Label: "From ICS Position",
	Flags: xscmsg.Required,
}

// ToLocationDef is the definition of the XSC-standard To Location field.
var ToLocationDef = &xscmsg.FieldDef{
	Tag:   "7b.",
	Label: "To Location",
	Key:   xscmsg.FToLocation,
	Flags: xscmsg.Required,
}

// FromLocationDef is the definition of the XSC-standard From Location field.
var FromLocationDef = &xscmsg.FieldDef{
	Tag:   "8b.",
	Label: "From Location",
	Flags: xscmsg.Required,
}

// ToNameDef is the definition of the XSC-standard To Name field.
var ToNameDef = &xscmsg.FieldDef{
	Tag:   "7c.",
	Label: "To Name",
}

// FromNameDef is the definition of the XSC-standard From Name field.
var FromNameDef = &xscmsg.FieldDef{
	Tag:   "8c.",
	Label: "From Name",
}

// ToContactDef is the definition of the XSC-standard To Contact field.
var ToContactDef = &xscmsg.FieldDef{
	Tag:   "7d.",
	Label: "To Contact Info",
}

// FromContactDef is the definition of the XSC-standard From Contact field.
var FromContactDef = &xscmsg.FieldDef{
	Tag:   "8d.",
	Label: "From Contact Info",
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
	Key:        xscmsg.FOpCall,
	Validators: []xscmsg.Validator{ValidateCallSign},
	Flags:      xscmsg.Required,
}

// OpDateDef is the definition of the XSC-standard Operator Date field.
var OpDateDef = &xscmsg.FieldDef{
	Tag:         "OpDate",
	Label:       "Operator Date",
	Comment:     "MM/DD/YYYY",
	Key:         xscmsg.FOpDate,
	Validators:  []xscmsg.Validator{ValidateDate},
	DefaultFunc: DefaultDate,
}

// OpTimeDef is the definition of the XSC-standard Operator Time field.
var OpTimeDef = &xscmsg.FieldDef{
	Tag:         "OpTime",
	Label:       "Operator Time",
	Comment:     "HH:MM",
	Key:         xscmsg.FOpTime,
	Validators:  []xscmsg.Validator{ValidateTime},
	DefaultFunc: DefaultTime,
}
