package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/rothskeller/packet/v4/incident"
)

// A bbsConnection is the server's record of a BBS exchange, and the watcher
// for it.  All of its fields (other than dir and cancel, which never change)
// are guarded by bbsConnectionsMutex.  Lock order: bbsConnectionsMutex may be
// held while taking linkMonitorsMutex, never the reverse.
type bbsConnection struct {
	dir              string
	ProgMsg          string `json:"progress"`
	ErrMsg           string `json:"error"`
	Seq              int    `json:"seq"`
	unsentErr        bool   // ErrMsg has not been sent to a poller yet
	finished         bool   // the exchange has finished
	closeAfterUpdate bool   // finished; remove from bbsConnections once the error is sent
	cancel           func()
	notify           chan struct{} // closed when Seq changes or exchange finishes
	done             chan struct{} // closed when the exchange finishes (may be nil in tests)
}

// bbsConnections maps incident directories to their connections.  A key with
// a nil value is an incident whose last connection has ended, for which
// GET /connect-progress returns 410 Gone.
var bbsConnections = map[string]*bbsConnection{}
var bbsConnectionsMutex sync.Mutex

// servePostConnectBBS handles POST /connect-bbs requests, which have a dir=
// parameter identifying the incident to connect with, and possibly an immOnly=
// parameter specifying that only immediate messages should be exchanged.  It
// always returns with 204.
func (s *Server) servePostConnectBBS(w http.ResponseWriter, r *http.Request) {
	var dir = r.FormValue("dir")
	bbsConnectionsMutex.Lock()
	if c := bbsConnections[dir]; c != nil && !c.finished { // connection already running
		bbsConnectionsMutex.Unlock()
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Any finished connection still recorded here is one whose error was
	// never fetched by a poller.  The user is starting a new connection,
	// which supersedes it, so its error is dropped.  (It was also written
	// to the Link Monitor, and the log file, when it happened.)
	ctx, cancel := context.WithCancel(context.Background())
	conn := &bbsConnection{dir: dir, cancel: cancel, done: make(chan struct{})}
	addMonitorLine(dir, monitorSeparator())
	bbsConnections[dir] = conn
	bbsConnectionsMutex.Unlock()
	go incident.BBSExchange(ctx, dir, r.FormValue("immOnly") != "", conn)
	w.WriteHeader(http.StatusNoContent)
}

// servePostConnectAbort handles POST /connect-abort requests, which has a dir=
// parameter identifying the incident whose connection should be aborted.  It
// always returns 204.
func (s *Server) servePostConnectAbort(w http.ResponseWriter, r *http.Request) {
	var dir = r.FormValue("dir")
	bbsConnectionsMutex.Lock()
	if bbsConnections[dir] != nil {
		bbsConnections[dir].cancel()
	}
	bbsConnectionsMutex.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

// serveGetConnectProgress handles GET /connect-progress requests, which have a
// dir= parameter identifying the incident owning the connection and a seq=
// parameter indicating the sequence number of the last update received.  This
// is a long polling request.  It returns with a BBS connection update when one
// becomes available, or a 410 Gone status if the connection has ended.
func (s *Server) serveGetConnectProgress(w http.ResponseWriter, r *http.Request) {
	var dir = r.FormValue("dir")
	var seq, _ = strconv.Atoi(r.FormValue("seq"))

	for {
		var update []byte
		var notify chan struct{}

		bbsConnectionsMutex.Lock()
		c, ok := bbsConnections[dir]
		switch {
		case !ok:
			bbsConnectionsMutex.Unlock()
			http.Error(w, "invalid connection directory", http.StatusBadRequest)
			return
		case c == nil, c.finished && c.Seq <= seq:
			// The connection has ended, and there's nothing more
			// for this caller to see.
			bbsConnectionsMutex.Unlock()
			w.WriteHeader(http.StatusGone)
			return
		case c.Seq > seq:
			update, _ = json.Marshal(c)
			c.unsentErr = false
			if c.closeAfterUpdate {
				bbsConnections[dir] = nil
			}
		default:
			if c.notify == nil {
				c.notify = make(chan struct{})
			}
			notify = c.notify
		}
		bbsConnectionsMutex.Unlock()
		if notify == nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(update)
			return
		}
		select {
		case <-s.stop:
			return
		case <-r.Context().Done():
			return
		case <-notify:
		}
	}
}

// wake wakes any pollers waiting on c.  bbsConnectionsMutex must be held.
func (c *bbsConnection) wake() {
	if c.notify != nil {
		close(c.notify)
		c.notify = nil
	}
}

func (c *bbsConnection) Progress(msg string) {
	addMonitorLine(c.dir, "*** "+time.Now().Format("15:04:05")+" "+msg)
	bbsConnectionsMutex.Lock()
	defer bbsConnectionsMutex.Unlock()
	c.ProgMsg = msg
	c.Seq++
	c.wake()
}

func (c *bbsConnection) LogEntry(*incident.LogEntry) {}

func (c *bbsConnection) Error(msg string) {
	addMonitorLine(c.dir, "*** "+time.Now().Format("15:04:05")+" ERROR: "+msg)
	bbsConnectionsMutex.Lock()
	defer bbsConnectionsMutex.Unlock()
	if c.ErrMsg == "" {
		c.ErrMsg = msg
		c.unsentErr = true
		c.Seq++
		c.wake()
		if c.finished && bbsConnections[c.dir] == c {
			// An error after Finished (which BBSExchange doesn't
			// do, but tolerate it): make sure it's delivered.
			c.closeAfterUpdate = true
		}
	}
}

// Finished implements incident.BBSExchangeWatcher.  Calls after the first are
// ignored.
func (c *bbsConnection) Finished() {
	bbsConnectionsMutex.Lock()
	defer bbsConnectionsMutex.Unlock()
	if c.finished {
		return
	}
	c.finished = true
	c.wake()
	if c.done != nil {
		close(c.done)
	}
	if bbsConnections[c.dir] != c {
		return // superseded by a newer connection
	}
	if c.unsentErr {
		c.closeAfterUpdate = true
	} else {
		bbsConnections[c.dir] = nil
	}
}

// abortBBSConnections aborts all running BBS exchanges, and waits (no longer
// than limit) for them to finish.  It's called when the server is stopping, so
// that each exchange gets to put its TNC back in order (see the transports'
// Close methods) before the program exits, rather than leaving the TNC
// connected, in converse mode, or holding frames.
func abortBBSConnections(limit time.Duration) {
	var waits []chan struct{}
	bbsConnectionsMutex.Lock()
	for _, c := range bbsConnections {
		if c != nil && !c.finished {
			c.cancel()
			if c.done != nil {
				waits = append(waits, c.done)
			}
		}
	}
	bbsConnectionsMutex.Unlock()
	if len(waits) == 0 {
		return
	}
	slog.Info("stopping server: aborting BBS connections", "count", len(waits))
	timer := time.NewTimer(limit)
	defer timer.Stop()
	for _, done := range waits {
		select {
		case <-done:
		case <-timer.C:
			slog.Warn("stopping server: BBS connections did not finish", "limit", limit)
			return
		}
	}
}
