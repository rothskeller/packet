package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/rothskeller/packet/incident"
)

type bbsConnection struct {
	mutex   sync.Mutex
	dir     string
	ProgMsg string `json:"progress"`
	ErrMsg  string `json:"error"`
	Seq     int    `json:"seq"`
	cancel  func()
	notify  chan struct{}
}

var bbsConnections = map[string]*bbsConnection{}
var bbsConnectionsMutex sync.Mutex

// servePostConnectBBS handles POST /connect-bbs requests, which have a dir=
// parameter identifying the incident to connect with, and possibly an immOnly=
// parameter specifying that only immediate messages should be exchanged.  It
// always returns with 204.
func (s *Server) servePostConnectBBS(w http.ResponseWriter, r *http.Request) {
	var dir = r.FormValue("dir")
	bbsConnectionsMutex.Lock()
	if bbsConnections[dir] != nil { // connection already running
		bbsConnectionsMutex.Unlock()
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	conn := &bbsConnection{dir: dir, cancel: cancel}
	bbsConnections[dir] = conn
	bbsConnectionsMutex.Unlock()
	go incident.BBSExchange(ctx, r.FormValue("dir"), r.FormValue("immOnly") != "", conn)
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
	var update []byte
	var notify chan struct{}

RESTART:
	bbsConnectionsMutex.Lock()
	if c, ok := bbsConnections[dir]; !ok {
		bbsConnectionsMutex.Unlock()
		http.Error(w, "invalid connection directory", http.StatusBadRequest)
		return
	} else if c == nil {
		bbsConnectionsMutex.Unlock()
		w.WriteHeader(http.StatusGone)
		return
	} else if c.Seq > seq {
		update, _ = json.Marshal(c)
	} else {
		if c.notify == nil {
			c.notify = make(chan struct{})
		}
		notify = c.notify
	}
	bbsConnectionsMutex.Unlock()
	if notify != nil {
		select {
		case <-s.stop:
			return
		case <-r.Context().Done():
			return
		case <-notify:
			notify = nil
			goto RESTART
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(update)
}

func (c *bbsConnection) Progress(msg string) {
	c.mutex.Lock()
	c.ProgMsg = msg
	c.Seq++
	if c.notify != nil {
		close(c.notify)
		c.notify = nil
	}
	c.mutex.Unlock()
}

func (c *bbsConnection) LogEntry(*incident.LogEntry) {}

func (c *bbsConnection) Error(msg string) {
	c.mutex.Lock()
	if c.ErrMsg == "" {
		c.ErrMsg = msg
		c.Seq++
		if c.notify != nil {
			close(c.notify)
			c.notify = nil
		}
	}
	c.mutex.Unlock()
}

func (c *bbsConnection) Finished() {
	bbsConnectionsMutex.Lock()
	if c.notify != nil {
		close(c.notify)
	}
	bbsConnections[c.dir] = nil
	bbsConnectionsMutex.Unlock()
}
