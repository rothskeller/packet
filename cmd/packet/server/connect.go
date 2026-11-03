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

type bbsConnection struct {
	dir              string
	ProgMsg          string `json:"progress"`
	ErrMsg           string `json:"error"`
	Seq              int    `json:"seq"`
	unsentErr        bool
	finished         bool
	closeAfterUpdate bool
	cancel           func()
	notify           chan struct{}
	done             chan struct{}
}

// bbsConnections maps incident directories to their connections.  A key with a
// nil value is an incident whose last connection has ended, for which GET
// /connect-progress returns 410 Gone.
var bbsConnections = map[string]*bbsConnection{}
var bbsConnectionsMutex sync.Mutex

// servePostConnectBBS handles POST /connect-bbs requests, which have a dir=
// parameter identifying the incident to connect with, and possibly an immOnly=
// parameter specifying that only immediate messages should be exchanged.  It
// always returns with 204.
func (s *Server) servePostConnectBBS(w http.ResponseWriter, r *http.Request) {
	var dir = r.FormValue("dir")
	bbsConnectionsMutex.Lock()
	if c := bbsConnections[dir]; c != nil && !c.finished {
		bbsConnectionsMutex.Unlock()
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Any finished connection still recorded here is one whose error was
	// never fetched by a poller.  The user is starting a new connection,
	// which supersedes it, so its error is dropped.
	ctx, cancel := context.WithCancel(context.Background())
	conn := &bbsConnection{dir: dir, cancel: cancel, done: make(chan struct{})}
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
		if c, ok := bbsConnections[dir]; !ok {
			bbsConnectionsMutex.Unlock()
			http.Error(w, "invalid connection directory", http.StatusBadRequest)
			return
		} else if c == nil || (c.finished && c.Seq <= seq) {
			bbsConnectionsMutex.Unlock()
			w.WriteHeader(http.StatusGone)
			return
		} else if c.Seq > seq {
			update, _ = json.Marshal(c)
			c.unsentErr = false
			if c.closeAfterUpdate {
				bbsConnections[dir] = nil
			}
		} else {
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
	bbsConnectionsMutex.Lock()
	defer bbsConnectionsMutex.Unlock()
	c.ProgMsg = msg
	c.Seq++
	c.wake()
}

func (c *bbsConnection) LogEntry(*incident.LogEntry) {}

func (c *bbsConnection) Error(msg string) {
	bbsConnectionsMutex.Lock()
	defer bbsConnectionsMutex.Unlock()
	if c.ErrMsg == "" {
		c.ErrMsg = msg
		c.unsentErr = true
		c.Seq++
		c.wake()
		if c.finished && bbsConnections[c.dir] == c {
			c.closeAfterUpdate = true
		}
	}
}

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
// that each exchange gets to put its TNC back in order before the program
// exits.
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
