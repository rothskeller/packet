// Package report generates reports on the messages in a practice session.
package report

import (
	"time"

	"github.com/rothskeller/packet/wppsvr/store"
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
	MessageTypes           string
	SentTo                 string
	SentBefore             string
	SentAfter              string
	NotSentFrom            string
	Modified               bool
	TotalMessages          int
	UniqueAddresses        int
	UniqueAddressesCorrect int
	PercentCorrect         int
	uniqueCallSigns        map[string]struct{}
	UniqueCallSigns        int
	UniqueCallSignsWeek    int
	Sources                []*Source
	Jurisdictions          []*Jurisdiction
	Messages               []*Message
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

// A Jurisdiction contains the information about a single jurisdiction in a
// Report.
type Jurisdiction struct {
	Name  string
	Count int
}

// A Message contains the information about a single message in a Report.
type Message struct {
	ID           string
	FromCallSign string
	Prefix       string
	Suffix       string
	Source       string
	Multiple     bool
	Jurisdiction string
	Class        string
	Problem      string
}
