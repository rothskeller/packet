package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/store"
)

type statistics struct {
	totalMessages          int
	uniqueAddresses        int
	uniqueAddressesCorrect int
	percentCorrect         int
	uniqueCallSigns        int
	uniqueCallSignsWeek    int
	sourceNames            []string
	sourceCounts           []int
	otherSourcesCount      int
}

// reportStatistics reports the statistical summary of practice messages.
func reportStatistics(sb *strings.Builder, st Store, session *store.Session, messages []*store.Message) {
	var (
		stats      *statistics
		countWidth int
	)
	stats = tabulate(st, session, messages)
	countWidth = numberWidth(stats.totalMessages)
	if cw := numberWidth(stats.uniqueCallSignsWeek); cw > countWidth {
		countWidth = cw
	}
	fmt.Fprintf(sb, "Total messages:     %*d\n", countWidth, stats.totalMessages)
	fmt.Fprintf(sb, "Unique addresses:   %*d\n", countWidth, stats.uniqueAddresses)
	if stats.uniqueAddresses != 0 {
		fmt.Fprintf(sb, "Correct messages:   %*d  (%d%%)\n", countWidth, stats.uniqueAddressesCorrect, stats.percentCorrect)
	}
	fmt.Fprintf(sb, "Unique call signs:  %*d  [report this count to the net]\n", countWidth, stats.uniqueCallSigns)
	if session.GenerateWeekSummary {
		fmt.Fprintf(sb, "  for the week:     %*d\n", countWidth, stats.uniqueCallSignsWeek)
	}
	sb.WriteString("Messages from:\n")
	for i, source := range stats.sourceNames {
		if wasSimulatedDown(session, source) {
			fmt.Fprintf(sb, "  %s:%*s%*d  (simulated down)\n",
				source, 17-len(source), "", countWidth, stats.sourceCounts[i])
		} else {
			fmt.Fprintf(sb, "  %s:%*s%*d\n", source, 17-len(source), "", countWidth, stats.sourceCounts[i])
		}
	}
	if stats.otherSourcesCount != 0 {
		fmt.Fprintf(sb, "  Other sources:    %*d\n", countWidth, stats.otherSourcesCount)
	}
	sb.WriteByte('\n')
}

// tabulate scans the messages accumulated in the session and computes the
// statistics that we will display.
func tabulate(st Store, session *store.Session, messages []*store.Message) (stats *statistics) {
	var (
		callsigns = make(map[string]bool)
		fromBBS   = make(map[string]int)
	)
	stats = new(statistics)
	stats.totalMessages = len(messages)
	messages = removeInvalidAndReplaced(messages)
	stats.uniqueAddresses = len(messages)
	for _, m := range messages {
		if m.FromBBS == "" {
			stats.otherSourcesCount++
		} else {
			fromBBS[m.FromBBS]++
		}
		if m.FromCallSign != "" {
			callsigns[m.FromCallSign] = true
		}
		if m.Correct {
			stats.uniqueAddressesCorrect++
		}
	}
	if stats.uniqueAddresses != 0 {
		stats.percentCorrect = stats.uniqueAddressesCorrect * 100 / stats.uniqueAddresses
	}
	stats.uniqueCallSigns = len(callsigns)
	stats.sourceNames = make([]string, 0, len(fromBBS))
	for source := range fromBBS {
		stats.sourceNames = append(stats.sourceNames, source)
	}
	sort.Strings(stats.sourceNames)
	stats.sourceCounts = make([]int, len(fromBBS))
	for i, source := range stats.sourceNames {
		stats.sourceCounts[i] = fromBBS[source]
	}
	if session.GenerateWeekSummary {
		stats.uniqueCallSignsWeek = uniqueCallSignsWeek(st, session, callsigns)
	}
	return stats
}

// uniqueCallSignsWeek looks up all other sessions that ended before the
// specified session, but during the same week as it; adds their unique call
// signs to the ones already in the supplied map; and returns the total number
// of unique call signs.
func uniqueCallSignsWeek(st Store, session *store.Session, callsigns map[string]bool) int {
	var (
		messages []*store.Message
		ostart   time.Time
	)
	// Calculate the start of the date range for interesting sessions, by
	// rewinding to Sunday of the same week as the argument session, and
	// the rewinding to midnight on that date.
	ostart = session.End.AddDate(0, 0, -int(session.End.Weekday()))
	ostart = time.Date(ostart.Year(), ostart.Month(), ostart.Day(), 0, 0, 0, 0, time.Local)
	// Get all of the sessions in that range.
	for _, osession := range st.GetSessionsEnding(ostart, session.End) {
		if osession.ExcludeFromWeekSummary { // e.g., PKTEST session
			continue
		}
		messages = st.GetSessionMessages(osession.ID)
		messages = removeInvalidAndReplaced(messages)
		for _, m := range messages {
			if m.FromCallSign != "" {
				callsigns[m.FromCallSign] = true
			}
		}
	}
	return len(callsigns)
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
		if !m.Valid {
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

// numberWidth returns the width of an integer in characters.
func numberWidth(i int) int {
	if i > 99 {
		return 3 // none of the numbers we deal with are over 999
	}
	if i > 9 {
		return 2
	}
	return 1 // none of them are negative, either
}
