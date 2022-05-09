// Package report generates reports on the messages in a practice session.
package report

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/analyze"
	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/english"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// This function returns the current time; it can be overridden by tests.
var now = func() time.Time { return time.Now() }

// Generate generates the report for the specified session.
func Generate(st Store, session *store.Session) *Report {
	var (
		r        Report
		messages []*store.Message
	)
	messages = st.GetSessionMessages(session.ID)
	messages = removeDroppedMessages(messages)
	generateTitle(&r, session)
	generateParams(&r, session)
	generateStatistics(&r, session, messages)
	generateMessages(&r, session, messages)
	generateGenInfo(&r)
	generateParticipants(&r, messages)
	if session.GenerateWeekSummary {
		generateWeekSummary(&r, st, session)
	}
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
	var (
		sb           strings.Builder
		article      string
		messageTypes []string
	)
	for i, id := range session.MessageTypes {
		var mt = config.LookupMessageType(id)

		messageTypes = append(messageTypes, mt.TypeName())
		if i == 0 {
			article = mt.TypeArticle()
		}
	}
	fmt.Fprintf(&sb, "This practice session expected %s %s sent to %s at %s, %s.",
		article, english.Conjoin(messageTypes, "or"), session.CallSign,
		english.Conjoin(session.ToBBSes, "or"),
		timerange(session.Start, session.End))
	switch len(session.DownBBSes) {
	case 0:
		break
	case 1:
		fmt.Fprintf(&sb, "  %s was simulated \"down\" for this practice session.",
			session.DownBBSes[0])
	default:
		fmt.Fprintf(&sb, "  %s were simulated \"down\" for this practice session.",
			english.Conjoin(session.DownBBSes, "and"))
	}
	r.Parameters = sb.String()
	r.Modified = session.Modified
}
func timerange(start, end time.Time) string {
	if start.Year() == end.Year() && start.Month() == end.Month() && start.Day() == end.Day() {
		return fmt.Sprintf("between %s and %s",
			start.Format("15:04"),
			end.Format("15:04 on Monday, January 2, 2006"))
	}
	if start.Year() == end.Year() {
		return fmt.Sprintf("between %s and %s",
			start.Format("15:04 on Monday, January 2"),
			end.Format("15:04 on Monday, January 2, 2006"))
	}
	return fmt.Sprintf("between %s and %s",
		start.Format("15:04 on Monday, January 2, 2006"),
		end.Format("15:04 on Monday, January 2, 2006"))
}

// generateStatistics scans the messages accumulated in the session and computes
// the statistics that we will display.
func generateStatistics(r *Report, session *store.Session, messages []*store.Message) {
	var (
		otherSourcesCount int
		fromBBS           = make(map[string]int)
	)
	r.uniqueCallSigns = make(map[string]struct{})
	r.TotalMessages = len(messages)
	messages = removeInvalidAndReplaced(messages)
	r.UniqueAddresses = len(messages)
	for _, m := range messages {
		if m.FromBBS == "" {
			otherSourcesCount++
		} else {
			fromBBS[m.FromBBS]++
		}
		if m.FromCallSign != "" {
			r.uniqueCallSigns[m.FromCallSign] = struct{}{}
		}
		if m.Actions&config.ActionError == 0 {
			r.UniqueAddressesCorrect++
		}
	}
	r.UniqueCallSigns = len(r.uniqueCallSigns)
	if r.UniqueAddresses != 0 {
		r.PercentCorrect = r.UniqueAddressesCorrect * 100 / r.UniqueAddresses
	}
	r.Sources = make([]*Source, 0, len(fromBBS))
	for source, count := range fromBBS {
		r.Sources = append(r.Sources, &Source{
			Name:          source,
			Count:         count,
			SimulatedDown: wasSimulatedDown(session, source),
		})
	}
	sort.Slice(r.Sources, func(i, j int) bool { return r.Sources[i].Name < r.Sources[j].Name })
	if otherSourcesCount != 0 {
		r.Sources = append(r.Sources, &Source{
			Name:  "Other sources",
			Count: otherSourcesCount,
		})
	}
}

// removeInvalidAndReplaced removes invalid and replaced messages from the list
// of messages.
func removeInvalidAndReplaced(messages []*store.Message) (out []*store.Message) {
	var (
		msgidx    int
		outidx    int
		addresses = make(map[string]bool)
	)
	out = make([]*store.Message, len(messages))
	outidx = len(messages)
	for msgidx = len(messages) - 1; msgidx >= 0; msgidx-- {
		var m = messages[msgidx]
		if m.Actions&config.ActionDontCount != 0 {
			continue
		}
		if addresses[m.FromAddress] {
			continue
		}
		addresses[m.FromAddress] = true
		outidx--
		out[outidx] = m
	}
	return out[outidx:]
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

		rm.FromAddress = m.FromAddress
		if m.FromAddress == "" && m.Subject == "" {
			rm.Subject = fmt.Sprintf("[unparseable message with hash %s]", m.Hash)
		} else {
			rm.Subject = m.Subject
		}
		for _, p := range m.Problems {
			if act, ok := actions[p]; ok {
				if act&config.ActionReport != 0 {
					rm.Problems = append(rm.Problems, analyze.ProblemLabel[p])
				}
			} else {
				// Problem was imported from old packet NCO
				// scripts.
				rm.Problems = append(rm.Problems, p)
			}
		}
		if m.Actions&config.ActionDontCount == 0 {
			r.CountedMessages = append(r.CountedMessages, &rm)
		} else {
			r.InvalidMessages = append(r.InvalidMessages, &rm)
		}
	}
	sort.Slice(r.CountedMessages, func(i, j int) bool { return compareMessages(r.CountedMessages[i], r.CountedMessages[j]) })
	sort.Slice(r.InvalidMessages, func(i, j int) bool { return compareMessages(r.InvalidMessages[i], r.InvalidMessages[j]) })
}
func compareMessages(a, b *Message) bool {
	as := strings.ToLower(a.FromAddress)
	bs := strings.ToLower(b.FromAddress)
	return as < bs
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

// generateWeekSummary looks up all other sessions that ended before the
// specified session, but during the same week as it.  It generates a count of
// unique call signs across all of the sessions.
func generateWeekSummary(r *Report, st Store, session *store.Session) {
	var (
		ostart time.Time
		unique = make(map[string]struct{})
	)
	// Put our own call signs into the map.
	for cs := range r.uniqueCallSigns {
		unique[cs] = struct{}{}
	}
	// Calculate the start of the date range for interesting sessions, by
	// rewinding to Sunday of the same week as the argument session, and
	// then rewinding to midnight on that date.
	ostart = session.End.AddDate(0, 0, -int(session.End.Weekday()))
	ostart = time.Date(ostart.Year(), ostart.Month(), ostart.Day(), 0, 0, 0, 0, time.Local)
	// Get all of the sessions in that range.
	for _, osession := range st.GetSessionsEnding(ostart, session.End) {
		if osession.ExcludeFromWeekSummary { // e.g., PKTEST session
			continue
		}
		oreport := Generate(st, osession)
		for cs := range oreport.uniqueCallSigns {
			unique[cs] = struct{}{}
		}
	}
	r.UniqueCallSignsWeek = len(unique)
}
