// Package report generates reports on the messages in a practice session.
package report

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/analyze"
	"github.com/rothskeller/packet/wppsvr/config"
	"github.com/rothskeller/packet/wppsvr/english"
	"github.com/rothskeller/packet/wppsvr/store"
)

// This function returns the current time; it can be overridden by tests.
var now = func() time.Time { return time.Now() }

// Generate generates the report for the specified session.
func Generate(st Store, session *store.Session) *Report {
	var (
		r        Report
		messages []*store.Message
		count    int
	)
	messages = st.GetSessionMessages(session.ID)
	count = len(messages)
	messages = removeDroppedMessages(messages)
	r.DroppedCount = count - len(messages)
	generateTitle(&r, session)
	generateParams(&r, session)
	generateStatistics(&r, session, messages)
	generateWeekSummary(&r, st, session)
	generateMessages(&r, session, messages)
	generateGenInfo(&r)
	generateParticipants(&r, messages)
	return &r
}

// removeDroppedMessages removes from the messages list any messages that
// should be excluded from the report (e.g., delivery receipts).
func removeDroppedMessages(messages []*store.Message) []*store.Message {
	j := 0
	for _, m := range messages {
		if m.Actions&config.ActionDropMsg == 0 {
			messages[j] = m
			j++
		}
	}
	return messages[:j]
}

// generateTitle generates the title of the report.
func generateTitle(r *Report, session *store.Session) {
	r.SessionName = session.Name
	r.SessionDate = session.End.Format("Monday, January 2, 2006")
	if session.Running {
		r.Preliminary = true
	}
}

// generateParams adds a description of the parameters of the practice session
// to the report.
func generateParams(r *Report, session *store.Session) {
	for _, id := range session.MessageTypes {
		r.MessageTypes = append(r.MessageTypes, config.LookupMessageType(id).TypeName())
	}
	r.SentTo = fmt.Sprintf("%s at %s", session.CallSign, english.Conjoin(session.ToBBSes, "or"))
	r.SentAfter = session.Start.Format("Mon 2006-01-02 15:04")
	r.SentBefore = session.End.Format("Mon 2006-01-02 15:04")
	r.NotSentFrom = english.Conjoin(session.DownBBSes, "or")
	r.Modified = session.Modified
}

// generateWeekSummary looks up all sessions that end in the same week as the
// specified session.  If the specified session is the last one of those, it
// generates a count of unique call signs across all of the sessions.
func generateWeekSummary(r *Report, st Store, session *store.Session) {
	var (
		ostart   time.Time
		sessions []*store.Session
		unique   = make(map[string]struct{})
	)
	// If the specified session isn't part of the official week, do nothing.
	if session.ExcludeFromWeek {
		return
	}
	// Get all of the sessions of the week.
	ostart = session.End.AddDate(0, 0, -int(session.End.Weekday()))
	ostart = time.Date(ostart.Year(), ostart.Month(), ostart.Day(), 0, 0, 0, 0, time.Local)
	sessions = st.GetSessions(ostart, ostart.AddDate(0, 0, 7))
	// Remove the ones that aren't officially part of the week.
	j := 0
	for _, s := range sessions {
		if !s.ExcludeFromWeek {
			sessions[j] = s
			j++
		}
	}
	sessions = sessions[:j]
	// If our specified session is not the last one on the list, or is the
	// only one on the list, do nothing.
	if len(sessions) < 2 || session.ID != sessions[len(sessions)-1].ID {
		return
	}
	// OK, we do want to generate the week's list of unique call signs.
	// Start with our own, and remove it from the list of sessions.
	// Put our own call signs into the map.
	for cs := range r.uniqueCallSigns {
		unique[cs] = struct{}{}
	}
	sessions = sessions[:len(sessions)-1]
	// Now add the unique call signs from the other sessions in the week.
	for _, osession := range sessions {
		oreport := Generate(st, osession)
		for cs := range oreport.uniqueCallSigns {
			unique[cs] = struct{}{}
		}
	}
	r.UniqueCallSignsWeek = len(unique)
}

// generateStatistics scans the messages accumulated in the session and computes
// the statistics that we will display.
func generateStatistics(r *Report, session *store.Session, messages []*store.Message) {
	var sources = make(map[string]int)
	var jurisdictions = make(map[string]int)
	var mtypes = make(map[string]int)

	r.uniqueCallSigns = make(map[string]struct{})
	messages, r.InvalidCount, r.ReplacedCount = removeInvalidAndReplaced(messages)
	for _, m := range messages {
		if m.FromBBS != "" {
			sources[m.FromBBS]++
		} else if strings.HasSuffix(strings.ToLower(m.FromAddress), "@winlink.org") {
			sources["Winlink"]++
		} else {
			sources["Email"]++
		}
		if m.Jurisdiction != "" {
			jurisdictions[m.Jurisdiction]++
		} else {
			jurisdictions["~~~"]++ // chosen to sort after anything real
		}
		mtypes[m.MessageType]++
		if m.FromCallSign != "" {
			r.uniqueCallSigns[m.FromCallSign] = struct{}{}
		}
		if m.Actions&config.ActionError != 0 {
			r.ErrorCount++
		} else if len(m.Problems) > 1 || (len(m.Problems) == 1 && m.Problems[0] != analyze.ProblemMultipleMessagesFromAddress) {
			r.WarningCount++
		} else {
			r.OKCount++
		}
	}
	r.UniqueCallSigns = len(r.uniqueCallSigns)
	r.Sources = make([]*Source, 0, len(sources))
	for source, count := range sources {
		r.Sources = append(r.Sources, &Source{
			Name:          source,
			Count:         count,
			SimulatedDown: wasSimulatedDown(session, source),
		})
	}
	sort.Slice(r.Sources, func(i, j int) bool { return r.Sources[i].Name < r.Sources[j].Name })
	r.Jurisdictions = make([]*Count, 0, len(jurisdictions))
	for jurisdiction, count := range jurisdictions {
		r.Jurisdictions = append(r.Jurisdictions, &Count{
			Name:  jurisdiction,
			Count: count,
		})
	}
	if len(r.Jurisdictions) == 1 && r.Jurisdictions[0].Name == "~~~" {
		r.Jurisdictions = nil
	}
	sort.Slice(r.Jurisdictions, func(i, j int) bool { return r.Jurisdictions[i].Name < r.Jurisdictions[j].Name })
	if len(r.Jurisdictions) != 0 && r.Jurisdictions[len(r.Jurisdictions)-1].Name == "~~~" {
		r.Jurisdictions[len(r.Jurisdictions)-1].Name = "???"
	}
	r.MTypeCounts = make([]*Count, 0, len(mtypes))
	for mtype, count := range mtypes {
		r.MTypeCounts = append(r.MTypeCounts, &Count{
			Name:  mtype,
			Count: count,
		})
	}
	sort.Slice(r.MTypeCounts, func(i, j int) bool { return r.MTypeCounts[i].Name < r.MTypeCounts[j].Name })
}

// removeInvalidAndReplaced removes invalid and replaced messages from the list
// of messages.
func removeInvalidAndReplaced(messages []*store.Message) (out []*store.Message, invalid, replaced int) {
	var (
		msgidx    int
		outidx    int
		addresses = make(map[string]bool)
	)
	out = make([]*store.Message, len(messages))
	outidx = len(messages)
	for msgidx = len(messages) - 1; msgidx >= 0; msgidx-- {
		var m = messages[msgidx]
		if addresses[m.FromAddress] {
			replaced++
			continue
		}
		if m.Actions&config.ActionDontCount != 0 {
			invalid++
			continue
		}
		addresses[m.FromAddress] = true
		outidx--
		out[outidx] = m
	}
	return out[outidx:], invalid, replaced
}

// wasSimulatedDown returns whether the specified BBS was simulated down for the
// practice session.
func wasSimulatedDown(session *store.Session, bbs string) bool {
	for _, down := range session.DownBBSes {
		if down == bbs {
			return true
		}
	}
	return false
}

// generateMessages generates the lists of valid and invalid check-in messages
// that appear in the report.
func generateMessages(r *Report, session *store.Session, messages []*store.Message) {
	var actions = config.Get().ProblemActionFlags

	messages = removeReplaced(messages)
	for _, m := range messages {
		var rm Message

		rm.ID = m.LocalID
		rm.FromCallSign = m.FromCallSign
		if len(m.FromCallSign) > 2 {
			if m.FromCallSign[1] >= '0' && m.FromCallSign[1] <= '9' {
				rm.Prefix, rm.Suffix = m.FromCallSign[:2], m.FromCallSign[2:]
			} else {
				rm.Prefix, rm.Suffix = m.FromCallSign[:3], m.FromCallSign[3:]
			}
		} else if idx := strings.IndexByte(m.FromAddress, '@'); idx > 2 {
			rm.Prefix, rm.Suffix = m.FromAddress[:3], m.FromAddress[3:idx]
		} else if idx >= 0 {
			rm.Prefix = m.FromAddress[:idx]
		} else {
			rm.Prefix, rm.Suffix = "???", "???"
		}
		if m.FromBBS != "" {
			rm.Source = m.FromBBS
		} else if strings.HasSuffix(strings.ToLower(m.FromAddress), "@winlink.org") {
			rm.Source = "Winlink"
		} else {
			rm.Source = "Email"
		}
		rm.Jurisdiction = m.Jurisdiction
		if m.Actions&config.ActionDontCount != 0 {
			rm.Class = "invalid"
		} else if m.Actions&config.ActionError != 0 {
			rm.Class = "error"
		} else {
			rm.Class = "ok"
		}
		for _, p := range m.Problems {
			if p == analyze.ProblemMultipleMessagesFromAddress {
				rm.Multiple = true
				continue
			}
			if act, ok := actions[p]; ok {
				if act&config.ActionReport == 0 {
					continue
				}
				p = analyze.ProblemLabel[p]
				if rm.Class == "ok" {
					rm.Class = "warning"
				}
			} else {
				// Problem was imported from old packet NCO
				// scripts.
			}
			if rm.Problem == "" {
				rm.Problem = p
			} else {
				rm.Problem = "multiple issues"
			}
		}
		r.Messages = append(r.Messages, &rm)
	}
	sort.Slice(r.Messages, func(i, j int) bool { return compareMessages(r.Messages[i], r.Messages[j]) })
}
func compareMessages(a, b *Message) bool {
	if a.FromCallSign != "" && b.FromCallSign == "" {
		return true
	}
	if a.FromCallSign == "" && b.FromCallSign != "" {
		return false
	}
	if a.FromCallSign != "" {
		if a.Suffix != b.Suffix {
			return a.Suffix < b.Suffix
		}
		return a.Prefix < b.Prefix
	}
	if a.Prefix != b.Prefix {
		return a.Prefix < b.Prefix
	}
	return a.Suffix < b.Suffix
}

// removeReplaced removes all but the last message from each address.  If more
// than one message is found from a given address, a MultipleMessagesFromAddress
// problem code is added to the one that is kept.
func removeReplaced(messages []*store.Message) (out []*store.Message) {
	var (
		msgidx    int
		outidx    int
		addresses = make(map[string]*store.Message)
	)
	out = make([]*store.Message, len(messages))
	outidx = len(messages)
	for msgidx = len(messages) - 1; msgidx >= 0; msgidx-- {
		m := messages[msgidx]
		if m.FromAddress == "" {
			outidx--
			out[outidx] = m
			continue
		}
		if keeper := addresses[m.FromAddress]; keeper == nil {
			outidx--
			out[outidx] = m
			addresses[m.FromAddress] = m
		} else if len(keeper.Problems) == 0 ||
			keeper.Problems[len(keeper.Problems)-1] != analyze.ProblemMultipleMessagesFromAddress {
			keeper.Problems = append(keeper.Problems, analyze.ProblemMultipleMessagesFromAddress)
		}
	}
	return out[outidx:]
}

// generateGenInfo records when then report was generated and by what software.
func generateGenInfo(r *Report) {
	if bi, ok := debug.ReadBuildInfo(); ok && bi.Main.Version != "" {
		r.GenerationInfo = fmt.Sprintf("This report was generated on %s by wppsvr version %s.",
			now().Format("Monday, January 2, 2006 at 15:04"), bi.Main.Version)
	} else {
		r.GenerationInfo = fmt.Sprintf("This report was generated on %s by wppsvr.",
			now().Format("Monday, January 2, 2006 at 15:04"))
	}
}

// generateParticipants returns a de-duplicated list of all from addresses of the
// supplied messages.
func generateParticipants(r *Report, messages []*store.Message) {
	var addresses = make(map[string]bool)

	for _, m := range messages {
		if m.FromAddress != "" {
			addresses[m.FromAddress] = true
		}
	}
	r.Participants = make([]string, 0, len(addresses))
	for address := range addresses {
		r.Participants = append(r.Participants, address)
	}
}
