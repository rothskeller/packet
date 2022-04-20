// Package report generates reports on the messages in a practice session.
package report

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"steve.rothskeller.net/packet/wppsvr/config"
	"steve.rothskeller.net/packet/wppsvr/english"
	"steve.rothskeller.net/packet/wppsvr/store"
)

// This function returns the current time; it can be overridden by tests.
var now = func() time.Time { return time.Now() }

// Store is an interface covering those methods of store.Store that are used in
// generating reports.
type Store interface {
	GetSessionMessages(int) []*store.Message
	GetSessionsEnding(start, end time.Time) []*store.Session
	UpdateSession(*store.Session)
	NextMessageID(string) string
}

// Generate generates the report for the specified session.
func Generate(st Store, session *store.Session) string {
	report, _ := generate(st, session)
	return report
}

// generate generates the report for the specified session.  It returns the
// report text and the list of addresses of participants who checked in (i.e.,
// to whom the report should be sent).
func generate(st Store, session *store.Session) (report string, participants []string) {
	var (
		sb       strings.Builder
		messages []*store.Message
	)
	messages = st.GetSessionMessages(session.ID)
	reportTitle(&sb, session)
	reportParams(&sb, session)
	reportStatistics(&sb, st, session, messages)
	reportMessages(&sb, session, messages)
	reportGenInfo(&sb)
	return sb.String(), allFromAddresses(messages)
}

// reportTitle generates the title of the report.
func reportTitle(sb *strings.Builder, session *store.Session) {
	var (
		line1  string
		line2  string
		maxlen int
		pad1   int
		pad2   int
	)
	line1 = "SCCo ARES/RACES Packet Practice Report"
	if session.Running {
		line1 = "***PRELIMINARY*** " + line1
	}
	maxlen = len(line1)
	line2 = fmt.Sprintf("for %s on %s", session.Name, session.End.Format("Monday, January 2, 2006"))
	if len(line2) > maxlen {
		maxlen = len(line2)
	}
	pad1 = (maxlen - len(line1)) / 2
	pad2 = (maxlen - len(line2)) / 2
	fmt.Fprintf(sb, "==== %*s%-*s ====\n==== %*s%-*s ====\n\n",
		pad1, "", maxlen-pad1, line1, pad2, "", maxlen-pad2, line2)
}

// reportParams writes a description of the parameters of the practice session
// to the report.
func reportParams(sb *strings.Builder, session *store.Session) {
	var (
		article      string
		messageTypes []string
		wr           = english.NewWrapper(sb)
	)
	for i, id := range session.MessageTypes {
		var mt = config.LookupMessageType(id)

		messageTypes = append(messageTypes, mt.TypeName())
		if i == 0 {
			article = mt.TypeArticle()
		}
	}
	fmt.Fprintf(wr, "This practice session expected %s %s sent to %s at %s, %s.",
		article, english.Conjoin(messageTypes, "or"), session.CallSign,
		english.Conjoin(session.ToBBSes, "or"),
		timerange(session.Start, session.End))
	switch len(session.DownBBSes) {
	case 0:
		break
	case 1:
		fmt.Fprintf(wr, "  %s was simulated \"down\" for this practice session.",
			session.DownBBSes[0])
	default:
		fmt.Fprintf(wr, "  %s were simulated \"down\" for this practice session.",
			english.Conjoin(session.DownBBSes, "and"))
	}
	if session.Modified {
		wr.WriteString("\n\nNOTE: The practice session expectations were changed after some check-in messages were received.  The earlier check-in messages may have been evaluated with different criteria.")
	}
	wr.WriteString("\n\n")
	wr.Close()
}

// reportGenInfo reports when then report was generated and by what software.
func reportGenInfo(sb *strings.Builder) {
	if bi, ok := debug.ReadBuildInfo(); ok {
		fmt.Fprintf(sb, "This report was generated on %s by wppsvr version %s.\n",
			now().Format("Monday, January 2, 2006 at 15:04"), bi.Main.Version)
	} else {
		fmt.Fprintf(sb, "This report was generated on %s by wppsvr.\n",
			now().Format("Monday, January 2, 2006 at 15:04"))
	}
}

// timerange returns a formatted date/time range.
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

// allFromAddresses returns a de-duplicated list of all from addresses of the
// supplied messages.
func allFromAddresses(messages []*store.Message) (list []string) {
	var addresses = make(map[string]bool)

	for _, m := range messages {
		if m.FromAddress != "" {
			addresses[m.FromAddress] = true
		}
	}
	list = make([]string, 0, len(addresses))
	for address := range addresses {
		list = append(list, address)
	}
	return list
}
