// Package report generates reports on the messages in a practice session.
package report

import (
	"time"

	"steve.rothskeller.net/packet/wppsvr/store"
)

// Store is an interface covering those methods of store.Store that are used in
// generating reports.
type Store interface {
	GetSessionMessages(int) []*store.Message
	GetSessions(start, end time.Time) []*store.Session
	UpdateSession(*store.Session)
	NextMessageID(string) string
}

// A Report contains all of the information that goes into a report about a
// practice session.  (This can include information from multiple sessions when
// a weekly summary is part of the report.)
type Report struct {
	SessionName            string
	SessionDate            string
	Preliminary            bool
	Parameters             string
	Modified               bool
	TotalMessages          int
	UniqueAddresses        int
	UniqueAddressesCorrect int
	PercentCorrect         int
	uniqueCallSigns        map[string]struct{}
	UniqueCallSigns        int
	UniqueCallSignsWeek    int
	Sources                []*Source
	CountedMessages        []*Message
	InvalidMessages        []*Message
	Participants           []string
	GenerationInfo         string
}

// A Source contains the information about a single source of messages in a
// Report.
type Source struct {
	Name          string
	Count         int
	SimulatedDown bool
}

// A Message contains the information about a single message in a Report.
type Message struct {
	FromAddress string
	Subject     string
	Problems    []string
}
