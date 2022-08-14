package webserver

import (
	_ "embed" // -
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rothskeller/packet/wppsvr/report"
	"github.com/rothskeller/packet/wppsvr/store"
)

//go:embed "front.html"
var frontHTML []byte

// serveFrontPage handles GET / requests.  These could be either page requests,
// asking for HTML, or fetch requests, asking for JSON.
func (ws *webserver) serveFrontPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Cache-Control", "nostore")
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(frontHTML)
		w.Write([]byte(`<script>var wppdata=`))
		ws.emitFrontJSON(w)
		w.Write([]byte(`;</script>`))
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		ws.emitFrontJSON(w)
	}
}

// emitFrontJSON emits the JSON data for the front page.
func (ws *webserver) emitFrontJSON(w io.Writer) {
	var week1, week2, nextUpdate time.Time
	var reload time.Duration

	switch now := time.Now(); now.Weekday() {
	case time.Sunday, time.Monday, time.Tuesday:
		// On Sunday through Tuesday, show last week and this week.
		week2 = now.AddDate(0, 0, -int(now.Weekday()))
	default:
		// On Wednesday through Saturday, show this week and next week.
		week2 = now.AddDate(0, 0, 7-int(now.Weekday()))
	}
	week1 = week2.AddDate(0, 0, -7)
	io.WriteString(w, "{")
	ws.emitFrontJSONWeek(w, "week1", week1)
	nextUpdate = ws.emitFrontJSONWeek(w, "week2", week2)
	reload = time.Until(nextUpdate) + time.Minute
	if reload > 5*time.Minute {
		reload = 5 * time.Minute
	}
	fmt.Fprintf(w, `"reload":%d}`, reload/time.Second)
}

// emitFrontJSONWeek emits the data for a single week.
func (ws *webserver) emitFrontJSONWeek(w io.Writer, name string, date time.Time) (nextUpdate time.Time) {
	var specs, svecs *store.Session

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	for _, session := range ws.st.GetSessions(date, date.AddDate(0, 0, 7)) {
		switch session.CallSign {
		case "PKTMON":
			specs = session
		case "PKTTUE":
			svecs = session
		}
	}
	fmt.Fprintf(w, `"%s":"%s",`, name, date.Format("January 2"))
	if specs != nil {
		r := report.Generate(ws.st, specs)
		fmt.Fprintf(w, `"%sspecs":{"date":"%s","count":%d`, name, specs.End.Format("January 2"), r.UniqueCallSigns)
		if count := r.OKCount + r.WarningCount + r.ErrorCount; count != 0 {
			pct := (r.OKCount + r.WarningCount) * 100 / count
			fmt.Fprintf(w, `,"correct":%d`, pct)
		}
		if specs.Running {
			io.WriteString(w, `,"preliminary":true`)
		}
		io.WriteString(w, `},`)
	}
	if svecs != nil {
		r := report.Generate(ws.st, svecs)
		fmt.Fprintf(w, `"%ssvecs":{"date":"%s","count":%d`, name, svecs.End.Format("January 2"), r.UniqueCallSigns)
		if count := r.OKCount + r.WarningCount + r.ErrorCount; count != 0 {
			pct := (r.OKCount + r.WarningCount) * 100 / count
			fmt.Fprintf(w, `,"correct":%d`, pct)
		}
		if svecs.Running {
			io.WriteString(w, `,"preliminary":true`)
		}
		io.WriteString(w, `},`)
		if r.UniqueCallSignsWeek != 0 {
			fmt.Fprintf(w, `"%scombined":%d,`, name, r.UniqueCallSignsWeek)
		}
	}
	if specs != nil && specs.Running {
		return specs.End
	} else if svecs != nil && svecs.Running {
		return svecs.End
	} else {
		return time.Now().Add(5 * time.Minute)
	}
}
