package server

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"
)

// linkMonitorLines is the number of recent lines kept for each incident's
// Link Monitor.
const linkMonitorLines = 5000

// A linkMonitor holds the recent lines of a live display of an incident's BBS
// exchanges, for the Link Monitor panel of the incident window.  It is kept after each connection
// ends, so that the window can show what happened.
type linkMonitor struct {
	lines  []string
	first  int           // sequence number of lines[0]
	notify chan struct{} // closed when lines are added
}

// linkMonitors maps incident directories to their Link Monitors.  Entries are
// created only by addMonitorLine.  linkMonitorsMutex guards linkMonitors, the
// linkMonitor structures, and linkMonitorCreated.  It may be taken while
// holding bbsConnectionsMutex, but never the reverse.
var linkMonitors = map[string]*linkMonitor{}
var linkMonitorsMutex sync.Mutex

// linkMonitorCreated, if not nil, is closed (and reset to nil) when a new
// entry is added to linkMonitors.  Requests for a monitor that doesn't exist
// yet wait on it.
var linkMonitorCreated chan struct{}

// addMonitorLine adds a line to the Link Monitor for the incident in dir.
func addMonitorLine(dir, line string) {
	linkMonitorsMutex.Lock()
	defer linkMonitorsMutex.Unlock()
	m := linkMonitors[dir]
	if m == nil {
		m = new(linkMonitor)
		linkMonitors[dir] = m
		if linkMonitorCreated != nil {
			close(linkMonitorCreated)
			linkMonitorCreated = nil
		}
	}
	m.lines = append(m.lines, line)
	if len(m.lines) > 2*linkMonitorLines {
		drop := len(m.lines) - linkMonitorLines
		m.lines = slices.Clone(m.lines[drop:])
		m.first += drop
	}
	if m.notify != nil {
		close(m.notify)
		m.notify = nil
	}
}

// Monitor implements incident.BBSExchangeMonitor.
func (c *bbsConnection) Monitor(line string) { addMonitorLine(c.dir, line) }

// serveGetLinkMonitorLines handles GET /link-monitor/lines requests, which
// have a dir= parameter identifying the incident and a seq= parameter giving
// the number of lines the caller already has.  This is a long polling request.
// It returns the lines added since then, as soon as there are any, as JSON:
// {"seq": N, "lines": [...], "skipped": N}.  "seq" is the value to pass on the
// next request, and "skipped" is the number of lines too old to return.
func (s *Server) serveGetLinkMonitorLines(w http.ResponseWriter, r *http.Request) {
	var (
		dir    = r.FormValue("dir")
		seq, _ = strconv.Atoi(r.FormValue("seq"))
		resp   struct {
			Seq     int      `json:"seq"`
			Lines   []string `json:"lines"`
			Skipped int      `json:"skipped"`
		}
	)
	for {
		var notify chan struct{}

		linkMonitorsMutex.Lock()
		if m := linkMonitors[dir]; m == nil {
			// There's no monitor for dir (yet), so it's empty.
			// Wait for some monitor to be created, and check again.
			// (We don't create one here, since dir can be any
			// string.)
			seq = 0 // the server restarted since the caller's last request
			if linkMonitorCreated == nil {
				linkMonitorCreated = make(chan struct{})
			}
			notify = linkMonitorCreated
		} else {
			end := m.first + len(m.lines)
			if seq > end {
				seq = 0 // the server restarted since the caller's last request
			}
			if seq < end {
				if seq < m.first {
					resp.Skipped, seq = m.first-seq, m.first
				}
				resp.Seq, resp.Lines = end, slices.Clone(m.lines[seq-m.first:])
				linkMonitorsMutex.Unlock()
				break
			}
			if m.notify == nil {
				m.notify = make(chan struct{})
			}
			notify = m.notify
		}
		linkMonitorsMutex.Unlock()
		select {
		case <-s.stop:
			return
		case <-r.Context().Done():
			return
		case <-notify:
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(resp)
}

// monitorSeparator returns the line that starts a new connection in the Link
// Monitor.
func monitorSeparator() string {
	return "------------------------ " + time.Now().Format("2006-01-02 15:04:05") + " ------------------------"
}
