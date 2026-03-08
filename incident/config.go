package incident

import (
	"context"
	"encoding/json/jsontext"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"strings"
	"time"
)

// A Config represents the configuration of an incident, i.e., all of the data
// and behavior that are not message-specific.
type Config struct {
	// IncidentName is the name of the incident on the ICS-309 log.
	IncidentName string `json:"incName,omitempty"`
	// ActivationNum is the activation number of the incident on the
	// ICS-309 log.
	ActivationNum string `json:"actNum,omitempty"`
	// OpStart is the start time of the operational period on the ICS-309
	// log.
	OpStart time.Time `json:"opStart,omitzero,format:'2006-01-02T15:04'"`
	// OpEnd is the end time of the operational period on the ICS-309 log.
	OpEnd time.Time `json:"opEnd,omitzero,format:'2006-01-02T15:04'"`
	// OpCall is the FCC call sign of the operator.
	OpCall string `json:"opCall,omitempty"`
	// OpName is the name of the operator.
	OpName string `json:"opName,omitempty"`
	// TacCall is the tactical call sign of the local station, if any.
	TacCall string `json:"tacCall,omitempty"`
	// TacName is the name of the local station, if any.
	TacName string `json:"tacName,omitempty"`
	// TxMessageID is the start of the message number sequence for outgoing
	// messages.
	TxMessageID string `json:"txMsgID,omitempty"`
	// RxMessageID is the start of the message number sequence for incoming
	// messages.
	RxMessageID string `json:"rxMsgID,omitempty"`
	// DefaultTo is the default To: destination for outgoing messages.
	DefaultTo string `json:"defTo,omitempty"`
	// DefaultToPos is the default value for the To ICS Position field for
	// outgoing forms messages.
	DefaultToPos string `json:"defToPos,omitempty"`
	// DefaultToLoc is the default value for the To Location field for
	// outgoing forms messages.
	DefaultToLoc string `json:"defToLoc,omitempty"`
	// DefaultFromPos is the default value for the From ICS Position field
	// for outgoing forms messages.
	DefaultFromPos string `json:"defFromPos,omitempty"`
	// DefaultFromLoc is the default value for the From Location field for
	// outgoing forms messages.
	DefaultFromLoc string `json:"defFromLoc,omitempty"`
	// DefaultBody is default contents for the primary body field of
	// outgoing messages (often used for "**** This is drill traffic ****").
	DefaultBody string `json:"defBody,omitempty"`
	// ConnectType is the method to use to connect to the BBS.  Options are
	// ConnectNone, ConnectSerialTNC, or ConnectTelnet.
	ConnectType string `json:"connType,omitempty"`
	// ConnectBBS is the name of the BBS that we are connecting to.
	ConnectBBS string `json:"connBBS,omitempty"`
	// ConnectAddress is the address of the BBS that we're connecting to.
	// If ConnectType is ConnectSerialTNC, this should be an AX.25 address
	// (e.g., "W5XSC-1").  If ConnectType is ConnectTelnet, this should be
	// a hostname:port or ipaddr:port string.
	ConnectAddress string `json:"connAddr,omitempty"`
	// SerialPort is the pathname of the serial port to use to communicate
	// with the TNC (e.g., "COM3" on Windows or "/dev/tty.something" on
	// Mac or Linux).  It is relevant only when ConnectType is
	// ConnectSerialTNC.
	SerialPort string `json:"serialPort,omitempty"`
	// TNCType is the type of TNC in use.  Currently the only supported
	// value is TNCKPC3Plus.  Relevant only when ConnectType is
	// ConnectSerialTNC.
	TNCType string `json:"tncType,omitempty"`
	// TelnetUser is the username to use to log in to the BBS.  Relevant
	// only when ConnectType is ConnectTelnet.
	TelnetUser string `json:"telnetUser,omitempty"`
	// TelnetPassword is the username to use to log in to the BBS.
	// Relevant only when ConnectType is ConnectTelnet.
	TelnetPassword string `json:"telnetPwd,omitempty"`
	// BulletinChecks is the set of bulletin areas that should be checked,
	// and the check frequency for each.
	BulletinChecks map[string]CheckFrequency `json:"bulletinChecks,omitempty"`
	// NoSendReceipts is a flag indicating that delivery receipts should
	// not be automatically generated and sent for received messages.
	NoSendReceipts bool `json:"noSendReceipts,omitempty"`
	// ViewFlags is a bitmask of flags describing how the incident log is
	// displayed.
	ViewFlags ViewFlag `json:"viewFlags,omitempty"`
	// AlLowVoice is a flag indicating that log entries can be marked as
	// voice messages to be put on a separate ICS-309.
	AllowVoice bool `json:"allowVoice,omitempty"`
}

// A ViewFlag is a flag (or bitmask of flags) describing how the incident log
// is displayed.
type ViewFlag uint8

// Values for ViewFlag
const (
	// ViewFull is a flag indicating that the incident log view should be in
	// full (vs. compact) format.
	ViewFull ViewFlag = 1 << iota
	// ViewReceipts is a flag indicating that the incident log view should
	// include receipt messages.
	ViewReceipts
	// ViewLarge is a flag indicating that the incident log view should use
	// a large font.
	ViewLarge
)

func (f ViewFlag) MarshalJSONTo(enc *jsontext.Encoder) (err error) {
	return enc.WriteToken(jsontext.String(f.String()))
}
func (f ViewFlag) String() string {
	var sb strings.Builder

	if f&ViewFull != 0 {
		sb.WriteByte('F')
	}
	if f&ViewLarge != 0 {
		sb.WriteByte('L')
	}
	if f&ViewReceipts != 0 {
		sb.WriteByte('R')
	}
	return sb.String()
}

func (f *ViewFlag) UnmarshalJSONFrom(dec *jsontext.Decoder) (err error) {
	if tok, err := dec.ReadToken(); err != nil {
		return err
	} else if tok.Kind() != '"' {
		return errors.New(`key "viewFlags" must map to a string`)
	} else if *f, err = ParseViewFlags(tok.String()); err != nil {
		return fmt.Errorf(`%s in key "viewFlags"`, err)
	} else {
		return nil
	}
}
func ParseViewFlags(s string) (f ViewFlag, err error) {
	for _, r := range s {
		switch r {
		case 'F':
			f |= ViewFull
		case 'L':
			f |= ViewLarge
		case 'R':
			f |= ViewReceipts
		default:
			return 0, fmt.Errorf(`unknown flag "%s"`, string(r))
		}
	}
	return f, nil
}

type CheckFrequency struct {
	time.Duration `json:"d,format:sec"`
}

// Methods for connecting to a BBS (i.e., values for Config.ConnectType).
const (
	// ConnectNone means there is no connection to the BBS.
	ConnectNone = ""
	// ConnectSerialTNC is connection to a BBS over the air using a serial
	// connection to a TNC.
	ConnectSerialTNC = "serial-tnc"
	// ConnectTelnet is connection to a BBS over the Internet using the
	// Telnet protocol.
	ConnectTelnet = "telnet"
)

// Types of TNCs for use with ConnectSerialTNC (i.e., values of Config.TNCType).
const (
	TNCKPC3Plus = "KPC3+" // Kantronics KPC-3 Plus
)

// ActiveCall returns the active call sign in the configuration.
func (c *Config) ActiveCall() string {
	if c.TacCall != "" {
		return c.TacCall
	}
	return c.OpCall
}

// ActiveName returns the active station name in the configuration.
func (c *Config) ActiveName() string {
	if c.TacCall != "" {
		return c.TacName
	}
	return c.OpName
}

// FromAddress returns the From address for messages sent using the config.
func (c *Config) FromAddress() string {
	return fmt.Sprintf("%s@%s.scc-ares-races.org", strings.ToLower(c.ActiveCall()), strings.ToLower(c.ConnectBBS))
}

// Clone creates a clone of the configuration.
func (c *Config) Clone() (n *Config) {
	n = new(Config)
	*n = *c
	n.BulletinChecks = maps.Clone(c.BulletinChecks)
	return n
}

// UpdateConfig applies a new configuration to the incident, logging the
// changes made.
func (inc *Incident) UpdateConfig(c *Config) {
	var attrs []slog.Attr

	if c.IncidentName != inc.Config.IncidentName {
		attrs = append(attrs, slog.String("IncidentName", c.IncidentName))
	}
	if c.ActivationNum != inc.Config.ActivationNum {
		attrs = append(attrs, slog.String("ActivationNum", c.ActivationNum))
	}
	if !c.OpStart.Equal(inc.Config.OpStart) {
		attrs = append(attrs, slog.String("OpStart", c.OpStart.Format("2006-01-02T15:04")))
	}
	if !c.OpEnd.Equal(inc.Config.OpEnd) {
		attrs = append(attrs, slog.String("OpEnd", c.OpEnd.Format("2006-01-02T15:04")))
	}
	if c.OpCall != inc.Config.OpCall {
		attrs = append(attrs, slog.String("OpCall", c.OpCall))
	}
	if c.OpName != inc.Config.OpName {
		attrs = append(attrs, slog.String("OpName", c.OpName))
	}
	if c.TacCall != inc.Config.TacCall {
		attrs = append(attrs, slog.String("TacCall", c.TacCall))
	}
	if c.TacName != inc.Config.TacName {
		attrs = append(attrs, slog.String("TacName", c.TacName))
	}
	if c.TxMessageID != inc.Config.TxMessageID {
		attrs = append(attrs, slog.String("TxMessageID", c.TxMessageID))
	}
	if c.RxMessageID != inc.Config.RxMessageID {
		attrs = append(attrs, slog.String("RxMessageID", c.RxMessageID))
	}
	if c.DefaultTo != inc.Config.DefaultTo {
		attrs = append(attrs, slog.String("DefaultTo", c.DefaultTo))
	}
	if c.DefaultToPos != inc.Config.DefaultToPos {
		attrs = append(attrs, slog.String("DefaultToPos", c.DefaultToPos))
	}
	if c.DefaultToLoc != inc.Config.DefaultToLoc {
		attrs = append(attrs, slog.String("DefaultToLoc", c.DefaultToLoc))
	}
	if c.DefaultFromPos != inc.Config.DefaultFromPos {
		attrs = append(attrs, slog.String("DefaultFromPos", c.DefaultFromPos))
	}
	if c.DefaultFromLoc != inc.Config.DefaultFromLoc {
		attrs = append(attrs, slog.String("DefaultFromLoc", c.DefaultFromLoc))
	}
	if c.DefaultBody != inc.Config.DefaultBody {
		attrs = append(attrs, slog.String("DefaultBody", c.DefaultBody))
	}
	if c.ConnectType != inc.Config.ConnectType {
		attrs = append(attrs, slog.String("ConnectType", c.ConnectType))
	}
	if c.ConnectBBS != inc.Config.ConnectBBS {
		attrs = append(attrs, slog.String("ConnectBBS", c.ConnectBBS))
	}
	if c.ConnectAddress != inc.Config.ConnectAddress {
		attrs = append(attrs, slog.String("ConnectAddress", c.ConnectAddress))
	}
	if c.SerialPort != inc.Config.SerialPort {
		attrs = append(attrs, slog.String("SerialPort", c.SerialPort))
	}
	if c.TNCType != inc.Config.TNCType {
		attrs = append(attrs, slog.String("TNCType", c.TNCType))
	}
	if c.TelnetUser != inc.Config.TelnetUser {
		attrs = append(attrs, slog.String("TelnetUser", c.TelnetUser))
	}
	if c.TelnetPassword != inc.Config.TelnetPassword {
		attrs = append(attrs, slog.String("TelnetPassword", "(redacted)"))
	}
	for area, freq := range c.BulletinChecks {
		if freq != inc.Config.BulletinChecks[area] {
			attrs = append(attrs, slog.Int("BulletinChecks."+area, int(freq.Duration/time.Minute)))
		}
	}
	for area := range inc.Config.BulletinChecks {
		if _, ok := c.BulletinChecks[area]; !ok {
			attrs = append(attrs, slog.String("BulletinChecks."+area, "(disabled)"))
		}
	}
	if c.NoSendReceipts != inc.Config.NoSendReceipts {
		attrs = append(attrs, slog.Bool("NoSendReceipts", c.NoSendReceipts))
	}
	if c.ViewFlags != inc.Config.ViewFlags {
		attrs = append(attrs, slog.String("ViewFlags", c.ViewFlags.String()))
	}
	if len(attrs) != 0 {
		slog.LogAttrs(context.Background(), slog.LevelInfo, "incident configuration changed", attrs...)
	}
	inc.Config = c
}
