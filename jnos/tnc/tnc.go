package tnc

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/rothskeller/packet/v4/jnos/kiss"
)

// A TNC describes the details of interacting with a particular type of TNC.
type TNC struct {
	// Name is the name of the TNC type.
	Name string
	// Robust contains the TNC parameters if the TNC is one with a built-in
	// AX.25 implementation that we issue commands to.  It should be set if
	// and only if KISS is not set.
	Robust *RobustTNC `json:"embed"`
	// KISS contains the TNC parameters if the TNC is a KISS-mode only TNC.
	// It should be set if and only if Robust is not set.
	KISS *kiss.Device
}

// A RobustTNC describes the details of interacting with a particular type of
// full-featured TNC (i.e., one with a built-in AX.25 implementation).
type RobustTNC struct {
	// CommandPrompt is the TNC command prompt that we wait for between
	// commands.
	CommandPrompt string
	// ConnectedMessage is the message from the TNC that indicates it is
	// connected to a BBS.
	ConnectedMessage string
	// DisconnectedMessage is the message from the TNC that indicates it has
	// lost its connection to the BBS.
	DisconnectedMessage string
	// PreConnectCommands is the list of commands sent to the TNC prior to
	// establishing a connection.  Note this does not include MyCallCommand.
	PreConnectCommands []string
	// PostDisconnectCommands is the list of commands sent to the TNC after
	// disconnecting from a BBS.
	PostDisconnectCommands []string
	// MyCallCommand is the command to set the call sign used for
	// connections.
	MyCallCommand string
	// ConnectCommand is the command to establish a connection.
	ConnectCommand string
	// DisconnectCommand is the command to disconnect a connection.
	DisconnectCommand string
	// EchoTimeout is the number of milliseconds to wait for the TNC to echo
	// data sent to it.
	EchoTimeout int
	// ReplyTimeout is the number of milliseconds to wait for the TNC to
	// reply to commands sent to it (not involving RF).
	ReplyTimeout int
}

var fetchOnce sync.Once

func fetch() {
	slog.Debug("Fetching TNC definitions from file system.")
	if tncdir := TNCDir(); tncdir == "" {
		slog.Warn("TNCDir", "err", "no TNC directory")
	} else if ents, err := os.ReadDir(tncdir); err != nil {
		slog.Warn("os.ReadDir", "d", tncdir, "err", err)
	} else {
		for _, ent := range ents {
			if strings.HasSuffix(ent.Name(), ".json") {
				fetchOne(filepath.Join(tncdir, ent.Name()))
			}
		}
	}
}

func fetchOne(path string) {
	var t TNC
	if by, err := os.ReadFile(path); err != nil {
		slog.Warn("os.ReadFile", "f", path, "err", err)
	} else if err = json.Unmarshal(by, &t, json.MatchCaseInsensitiveNames(true), json.RejectUnknownMembers(true)); err != nil {
		slog.Warn("json.Unmarshal", "f", path, "err", err)
	} else if err = t.Validate(); err != nil {
		slog.Warn("invalid TNC definition", "f", path, "err", err)
	} else {
		if t.Robust != nil && t.Robust.DisconnectCommand == "" {
			// DisconnectCommand was added after initial release, so
			// there may be configs where it isn't set.
			t.Robust.DisconnectCommand = "D"
		}
		tncs[t.Name] = &t
	}
}

func (t *TNC) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("Name is missing or blank")
	}
	if (t.Robust == nil) == (t.KISS == nil) {
		return errors.New("exactly one of Robust or KISS must be set")
	}
	if t.Robust != nil {
		return t.Robust.Validate()
	} else {
		return t.KISS.Validate()
	}
}

func (t *RobustTNC) Validate() error {
	for _, f := range []struct{ name, value string }{
		{"CommandPrompt", t.CommandPrompt},
		{"ConnectedMessage", t.ConnectedMessage},
		{"DisconnectedMessage", t.DisconnectedMessage},
		{"MyCallCommand", t.MyCallCommand},
		{"ConnectCommand", t.ConnectCommand},
		{"DisconnectCommand", t.DisconnectCommand},
	} {
		if strings.TrimSpace(f.value) == "" {
			return fmt.Errorf("%s is missing or blank", f.name)
		}
	}
	if t.EchoTimeout < 0 {
		return errors.New("EchoTimeout is negative")
	}
	if t.ReplyTimeout < 0 {
		return errors.New("ReplyTimeout is negative")
	}
	return nil
}

// AllTNCs returns an alphabetical list of the names of all supported TNCs.
func AllTNCs() []string {
	fetchOnce.Do(fetch)
	var all = slices.Collect(maps.Keys(tncs))
	slices.Sort(all)
	return all
}

// Get returns the TNC definition for the TNC with the specified name, or nil
// if there is none.
func Get(name string) *TNC {
	fetchOnce.Do(fetch)
	return tncs[name]
}
