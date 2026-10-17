package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// testTimeout bounds every wait in these tests, so that they fail rather than
// hang if the code under test is broken.
const testTimeout = 10 * time.Second

func newTestServer() *Server { return &Server{stop: make(chan struct{})} }

// poll issues a GET /connect-progress request (with a time limit) and returns
// the status code and decoded body.  It may be called on any goroutine.
func poll(t *testing.T, s *Server, dir string, seq int) (code int, body map[string]any) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	q := url.Values{"dir": {dir}, "seq": {strconv.Itoa(seq)}}
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/connect-progress?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	s.serveGetConnectProgress(w, r)
	if ctx.Err() != nil {
		t.Errorf("poll(%q, %d) timed out", dir, seq) // (may be on another goroutine; no Fatal)
		return -1, nil
	}
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Errorf("poll body %q: %s", w.Body.String(), err)
		}
	}
	return w.Code, body
}

// addTestConnection registers a connection for dir, as servePostConnectBBS
// does, without starting an exchange.
func addTestConnection(t *testing.T, dir string) *bbsConnection {
	c := &bbsConnection{dir: dir, cancel: func() {}}
	bbsConnectionsMutex.Lock()
	bbsConnections[dir] = c
	bbsConnectionsMutex.Unlock()
	t.Cleanup(func() {
		bbsConnectionsMutex.Lock()
		delete(bbsConnections, dir)
		bbsConnectionsMutex.Unlock()
		linkMonitorsMutex.Lock()
		delete(linkMonitors, dir)
		linkMonitorsMutex.Unlock()
	})
	return c
}

func TestConnectionProgressSequencing(t *testing.T) {
	s := newTestServer()
	dir := "/test/sequencing"
	c := addTestConnection(t, dir)

	// A poller waiting for an update is woken by Progress.
	got := make(chan map[string]any, 1)
	go func() {
		_, body := poll(t, s, dir, 0)
		got <- body
	}()
	time.Sleep(20 * time.Millisecond) // let the poller start waiting
	c.Progress("Connecting...")
	select {
	case body := <-got:
		if body["progress"] != "Connecting..." || body["seq"] != float64(1) {
			t.Errorf("first update %v", body)
		}
	case <-time.After(testTimeout):
		t.Fatal("poller not woken by Progress")
	}
	// An error, then Finished: the error is delivered, then 410.
	c.Error("it broke")
	c.Finished()
	if code, body := poll(t, s, dir, 1); code != http.StatusOK || body["error"] != "it broke" {
		t.Errorf("after error: %d %v", code, body)
	}
	if code, _ := poll(t, s, dir, 2); code != http.StatusGone {
		t.Errorf("after delivered error: %d, want 410", code)
	}
	// Late calls must not panic (e.g., double close of notify).
	c.Finished()
	c.Error("another")
	c.Progress("late")
}

func TestConnectionFinishedWakesPoller(t *testing.T) {
	s := newTestServer()
	dir := "/test/finished-wakes"
	c := addTestConnection(t, dir)
	c.Progress("Working...")
	code := make(chan int, 1)
	go func() {
		c, _ := poll(t, s, dir, 1)
		code <- c
	}()
	time.Sleep(20 * time.Millisecond)
	c.Finished()
	select {
	case got := <-code:
		if got != http.StatusGone {
			t.Errorf("got %d, want 410", got)
		}
	case <-time.After(testTimeout):
		t.Fatal("poller not woken by Finished")
	}
	// A poller with a stale, larger seq for a finished connection that
	// still has an undelivered error must not wait forever.
	c2 := addTestConnection(t, dir)
	c2.Error("e")
	c2.Finished()
	if got, _ := poll(t, s, dir, 99); got != http.StatusGone {
		t.Errorf("stale seq: got %d, want 410", got)
	}
}

// waitFinished waits (with a time limit) for the connection registered for dir
// to finish, and returns it.
func waitFinished(t *testing.T, dir string) *bbsConnection {
	t.Helper()
	deadline := time.Now().Add(testTimeout)
	for time.Now().Before(deadline) {
		bbsConnectionsMutex.Lock()
		c := bbsConnections[dir]
		finished := c != nil && c.finished
		bbsConnectionsMutex.Unlock()
		if finished {
			return c
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("connection for %s did not finish", dir)
	return nil
}

func TestConnectAfterUnfetchedError(t *testing.T) {
	s := newTestServer()
	// A directory that isn't an incident: the exchange fails immediately,
	// without any network activity.
	dir := t.TempDir()
	t.Cleanup(func() {
		bbsConnectionsMutex.Lock()
		delete(bbsConnections, dir)
		bbsConnectionsMutex.Unlock()
		linkMonitorsMutex.Lock()
		delete(linkMonitors, dir)
		linkMonitorsMutex.Unlock()
	})
	post := func() {
		r := httptest.NewRequest(http.MethodPost, "/connect-bbs?"+url.Values{"dir": {dir}}.Encode(), nil)
		w := httptest.NewRecorder()
		s.servePostConnectBBS(w, r)
		if w.Code != http.StatusNoContent {
			t.Fatalf("POST /connect-bbs: %d", w.Code)
		}
	}
	post()
	first := waitFinished(t, dir)
	bbsConnectionsMutex.Lock()
	unsent := first.closeAfterUpdate
	bbsConnectionsMutex.Unlock()
	if !unsent {
		t.Fatal("first connection should be holding an unfetched error")
	}
	// Nobody fetched the error.  A new connect must start a new
	// connection, not be ignored as "already running".
	post()
	second := waitFinished(t, dir)
	if second == first {
		t.Fatal("second POST /connect-bbs did not start a new connection")
	}
	if code, body := poll(t, s, dir, 0); code != http.StatusOK || body["error"] == "" {
		t.Errorf("second connection's error not delivered: %d %v", code, body)
	}
	// The superseded connection's late calls must not disturb the new
	// state.
	first.Finished()
	bbsConnectionsMutex.Lock()
	cur := bbsConnections[dir]
	bbsConnectionsMutex.Unlock()
	if cur != nil {
		t.Errorf("after delivering error, entry is %v, want nil", cur)
	}
}

func TestLinkMonitorUnknownDir(t *testing.T) {
	s := newTestServer()
	dir := "/test/no-such-monitor"
	t.Cleanup(func() {
		linkMonitorsMutex.Lock()
		delete(linkMonitors, dir)
		linkMonitorsMutex.Unlock()
	})
	get := func(ctx context.Context) *httptest.ResponseRecorder {
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/link-monitor/lines?"+url.Values{"dir": {dir}, "seq": {"0"}}.Encode(), nil)
		w := httptest.NewRecorder()
		s.serveGetLinkMonitorLines(w, r)
		return w
	}
	// A request for an unknown dir waits, and doesn't create an entry.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	get(ctx)
	cancel()
	linkMonitorsMutex.Lock()
	_, exists := linkMonitors[dir]
	linkMonitorsMutex.Unlock()
	if exists {
		t.Fatal("GET /link-monitor/lines created a monitor for an unknown dir")
	}
	// A waiting request is woken when the monitor is created.
	ctx, cancel = context.WithTimeout(context.Background(), testTimeout)
	defer cancel()
	result := make(chan *httptest.ResponseRecorder, 1)
	go func() { result <- get(ctx) }()
	time.Sleep(20 * time.Millisecond)
	addMonitorLine("/test/some-other-dir", "unrelated") // wakes, rechecks, waits again
	addMonitorLine(dir, "hello")
	t.Cleanup(func() {
		linkMonitorsMutex.Lock()
		delete(linkMonitors, "/test/some-other-dir")
		linkMonitorsMutex.Unlock()
	})
	select {
	case w := <-result:
		var resp struct {
			Seq   int      `json:"seq"`
			Lines []string `json:"lines"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("response %q: %s", w.Body.String(), err)
		}
		if resp.Seq != 1 || len(resp.Lines) != 1 || resp.Lines[0] != "hello" {
			t.Errorf("response %+v", resp)
		}
	case <-time.After(testTimeout):
		t.Fatal("waiting request not woken when the monitor was created")
	}
}

func TestServerStopTwice(t *testing.T) {
	s := newTestServer()
	r := httptest.NewRequest(http.MethodPost, "/stop", nil)
	s.handleStop(httptest.NewRecorder(), r)
	s.handleStop(httptest.NewRecorder(), r)
	s.requestStop()
	select {
	case <-s.stop:
	default:
		t.Error("stop channel not closed")
	}
}

func TestGUIRequestClearsOutpostOnly(t *testing.T) {
	s := newTestServer()
	s.outpostOnly.Store(true)
	called := false
	h := s.gui(func(http.ResponseWriter, *http.Request) { called = true })
	h(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	if !called || s.outpostOnly.Load() {
		t.Errorf("called=%v outpostOnly=%v", called, s.outpostOnly.Load())
	}
}

// TestAbortBBSConnections checks that stopping the server aborts running
// exchanges and waits for them to finish (so they can clean up their TNCs),
// but no longer than the limit.
func TestAbortBBSConnections(t *testing.T) {
	finishing := addTestConnection(t, "/test/abort-finishing")
	finishing.done = make(chan struct{})
	finishing.cancel = func() { go func() { time.Sleep(100 * time.Millisecond); finishing.Finished() }() }
	stuck := addTestConnection(t, "/test/abort-stuck")
	stuck.done = make(chan struct{})
	var canceled atomic.Bool
	stuck.cancel = func() { canceled.Store(true) }
	finished := addTestConnection(t, "/test/abort-finished")
	finished.Finished()
	finished.cancel = func() { t.Error("canceled a finished connection") }

	start := time.Now()
	abortBBSConnections(500 * time.Millisecond)
	elapsed := time.Since(start)
	if !canceled.Load() {
		t.Error("running connection not canceled")
	}
	select {
	case <-finishing.done:
	default:
		t.Error("didn't wait for the finishing connection")
	}
	if elapsed < 400*time.Millisecond || elapsed > 2*time.Second {
		t.Errorf("waited %v for a stuck connection, want about the 500ms limit", elapsed)
	}
}
