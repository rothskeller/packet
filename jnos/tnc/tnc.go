package tnc

import (
	"encoding/json/v2"
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
	} else {
		tncs[t.Name] = &t
	}
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
