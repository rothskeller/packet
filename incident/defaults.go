package incident

import (
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"

	"github.com/rothskeller/packet/v4/cmd/packet/osdep"
)

// maxIncidentDirs is the maximum number of incident directories to keep in
// IncDefaults.IncidentDirs.
const maxIncidentDirs = 8

// The IncDefaults structure contains default values for some incident config
// settings:  those that tend to be the same for all incidents.  Most fields
// are updated whenever an incident is configured, and contain the most
// recently used value of that field.  TelnetPasswords is cumulative, retaining
// the most recently used password for any BBS and user.  IncidentDirs is
// cumulative with a limit.
type IncDefaults struct {
	// OpCall is the FCC call sign of the operator.
	OpCall string `json:"opCall,omitempty"`
	// OpName is the name of the operator.
	OpName string `json:"opName,omitempty"`
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
	// TCPAddresses is a map from ConnectBBS to ConnectAddress for Telnet
	// connections.
	TCPAddresses map[string]string `json:"tcpAddrs,omitempty"`
	// TelnetPasswords is a map from TelnetUser to TelnetPassword.
	TelnetPasswords map[string]string `json:"telnetPwds,omitempty"`
	// IncidentDirs is a list of recent incident directories, used to make
	// it easy for the user to return to them.
	IncidentDirs []string `json:"incidentDirs,omitempty"`
	// ViewFlags is a bitmask of flags describing how the incident log is
	// displayed.
	ViewFlags ViewFlag `json:"viewFlags,omitempty"`
}

// GetIncDefaults retrieves the current incident defaults.  It returns a valid
// structure even if there are errors or no saved defaults.
func GetIncDefaults() (idef *IncDefaults) {
	var (
		fname string
		fh    *os.File
		err   error
	)
	idef = new(IncDefaults)
	if fname = osdep.DefaultsFile; fname == "" {
		return idef
	}
	if fh, err = os.Open(fname); err != nil {
		if !os.IsNotExist(err) {
			slog.Error("os.Open", "f", fname, "err", err)
		}
		return idef
	}
	defer fh.Close()
	if err = osdep.ReadLock(fh); err != nil {
		slog.Error("osdep.ReadLock", "f", fname, "err", err)
		return idef
	}
	defer osdep.Unlock(fh)
	if err = json.UnmarshalRead(fh, idef, json.RejectUnknownMembers(true)); err != nil {
		slog.Error("json.UnmarshalRead", "f", fname, "err", err)
	}
	return idef
}

func configFromDefaults() (c *Config) {
	defs := GetIncDefaults()
	return &Config{
		OpCall:         defs.OpCall,
		OpName:         defs.OpName,
		ConnectType:    defs.ConnectType,
		ConnectBBS:     defs.ConnectBBS,
		ConnectAddress: defs.ConnectAddress,
		SerialPort:     defs.SerialPort,
		TNCType:        defs.TNCType,
		TelnetPassword: defs.TelnetPasswords[defs.OpCall],
		ViewFlags:      defs.ViewFlags,
	}
}

// UpdateIncDefaults updates the incident defaults based on the configuration
// of the receiver incident.  Errors are logged and ignored.
func (inc *Incident) UpdateIncDefaults() {
	var (
		fname string
		fh    *os.File
		idef  IncDefaults
		err   error
	)
	// Open and read the defaults file.
	if fname = osdep.DefaultsFile; fname == "" {
		return
	}
	if err = os.MkdirAll(filepath.Dir(fname), 0777); err != nil {
		slog.Error("os.MkdirAll", "d", filepath.Dir(fname), "err", err)
		return
	}
	if fh, err = os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0666); err != nil {
		slog.Error("os.OpenFile", "f", fname, "err", err)
		return
	}
	defer fh.Close()
	if err = osdep.WriteLock(fh); err != nil {
		slog.Error("osdep.WriteLock", "f", fname, "err", err)
		return
	}
	defer osdep.Unlock(fh)
	if err = json.UnmarshalRead(fh, &idef, json.RejectUnknownMembers(true)); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		slog.Error("json.UnmarshalRead", "f", fname, "err", err)
		return
	}
	// Update the incident defaults structure.
	if inc.Config.OpCall != "" {
		idef.OpCall = inc.Config.OpCall
	}
	if inc.Config.OpName != "" {
		idef.OpName = inc.Config.OpName
	}
	if inc.Config.ConnectType != "" {
		idef.ConnectType = inc.Config.ConnectType
	}
	if inc.Config.ConnectBBS != "" {
		idef.ConnectBBS = inc.Config.ConnectBBS
	}
	if inc.Config.ConnectAddress != "" {
		idef.ConnectAddress = inc.Config.ConnectAddress
	}
	if inc.Config.ConnectType == ConnectSerialTNC {
		if inc.Config.SerialPort != "" {
			idef.SerialPort = inc.Config.SerialPort
		}
		if inc.Config.TNCType != "" {
			idef.TNCType = inc.Config.TNCType
		}
	}
	if inc.Config.ConnectType == ConnectTelnet && inc.Config.ConnectBBS != "" && inc.Config.ConnectAddress != "" {
		if idef.TCPAddresses == nil {
			idef.TCPAddresses = make(map[string]string)
		}
		idef.TCPAddresses[inc.Config.ConnectBBS] = inc.Config.ConnectAddress
	}
	if inc.Config.ConnectType == ConnectTelnet && inc.Config.TelnetUser != "" && inc.Config.TelnetPassword != "" {
		if idef.TelnetPasswords == nil {
			idef.TelnetPasswords = make(map[string]string)
		}
		idef.TelnetPasswords[inc.Config.TelnetUser] = inc.Config.TelnetPassword
	}
	idef.IncidentDirs = slices.DeleteFunc(idef.IncidentDirs, func(dir string) bool { return dir == inc.Dir })
	idef.IncidentDirs = append(idef.IncidentDirs, inc.Dir)
	if len(idef.IncidentDirs) > maxIncidentDirs {
		idef.IncidentDirs = idef.IncidentDirs[len(idef.IncidentDirs)-maxIncidentDirs:]
	}
	idef.ViewFlags = inc.Config.ViewFlags
	// Rewind and truncate the file and write the structure.
	if _, err = fh.Seek(0, 0); err != nil {
		slog.Error("fh.Seek", "f", fname, "err", err)
		return
	}
	if err = fh.Truncate(0); err != nil {
		slog.Error("fh.Truncate", "f", fname, "err", err)
		return
	}
	if err = json.MarshalWrite(fh, &idef); err != nil {
		slog.Error("json.MarshalWrite", "f", fname, "err", err)
	}
}
