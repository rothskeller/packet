package tnc

import (
	"encoding/json/v2"
	"errors"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// A TNC describes the details of interacting with a particular type of TNC.
type TNC struct {
	// Name is the name of the TNC type.
	Name string
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
		tncs[t.Name] = &t
	}
}

// Validate returns an error if the TNC definition is not usable: if any of
// the strings needed to interact with the TNC is empty (or only white
// space), or if a timeout is
// negative.  (A zero timeout means "use the default".)  An empty prompt or
// message string would match immediately or never, depending on how it is
// used, so it is never valid.
func (t *TNC) Validate() error {
	if t == nil {
		return errors.New("no TNC definition")
	}
	for _, f := range []struct{ name, value string }{
		{"Name", t.Name},
		{"CommandPrompt", t.CommandPrompt},
		{"ConnectedMessage", t.ConnectedMessage},
		{"DisconnectedMessage", t.DisconnectedMessage},
		{"MyCallCommand", t.MyCallCommand},
		{"ConnectCommand", t.ConnectCommand},
	} {
		if strings.TrimSpace(f.value) == "" {
			return errors.New(f.name + " is missing or blank")
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
