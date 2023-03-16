package pktmgr

import (
	"regexp"
	"strconv"
	"sync"
)

// Config is the configuration for an incident.
type Config struct {
	// IncidentName is the name of the incident, recorded in the ICS-309.
	IncidentName string
	// ActivationNum is the activation number for the incident, recorded in
	// the ICS-309.
	ActivationNum string
	// OpStartDate is the date of the start of the operational period,
	// recorded in the ICS-309.
	OpStartDate string
	// OpStartTime is the time of the start of the operational period,
	// recorded in the ICS-309.
	OpStartTime string
	// OpEndDate is the date of the end of the operational period, recorded
	// in the ICS-309.
	OpEndDate string
	// OpEndTime is the time of the end of the operational period, recorded
	// in the ICS-309.
	OpEndTime string
	// BBS is the name of the operator's home BBS.
	BBS string
	// OpCall is the operator's FCC call sign.
	OpCall string
	// OpName is the operator's name.
	OpName string
	// TacCall is the tactical call sign for the station.  It is optional.
	TacCall string `json:",omitempty"`
	// TacName is the name of the tactical station.  It should be provided
	// if and only if TacCall is provided.
	TacName string `json:",omitempty"`
	// StartMsgID is the local message ID for the first message sent or
	// received.  It must have the form XXX-###S, where S is either "P" or
	// "M".
	StartMsgID string
	// DefBody is an optional string added to the body of any new message.
	DefBody string `json:",omitempty"`
	// BackgroundPDF is an indicator that PDF generation should happen in a
	// background goroutine.  If BackgroundPDF is set, the goroutine will
	// acquire the lock before generating the PDF and release it when
	// finished.  If BackgroundPDF is nil, PDF generation will happen in the
	// foreground.
	BackgroundPDF *sync.Mutex
	// msgIDPrefix is the prefix of local message IDs, derived from
	// StartMsgID.
	msgIDPrefix string
	// startMsgNum is the sequence number in StartMsgID.
	startMsgNum int
	// msgIDSuffix is the suffix of local message IDs, derived from
	// StartMsgID.
	msgIDSuffix string
	// callsign is the local call sign:  TacCall if specified, and OpCall
	// otherwise.
	callsign string
	// name is the local station name:  TacName if TacCall is specified, and
	// OpName otherwise.
	name string
}

var startMsgIDRE = regexp.MustCompile(`^((?:[A-Z][A-Z0-9][A-Z0-9]|[0-9][A-Z][A-Z])-)(\d+)([PM])$`)

func (c *Config) fillin() {
	if c.TacCall != "" {
		c.callsign, c.name = c.TacCall, c.TacName
	} else {
		c.callsign, c.name = c.OpCall, c.OpName
	}
	match := startMsgIDRE.FindStringSubmatch(c.StartMsgID)
	c.msgIDPrefix = match[1]
	c.startMsgNum, _ = strconv.Atoi(match[2])
	c.msgIDSuffix = match[3]
}
