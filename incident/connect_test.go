package incident

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
	"github.com/rothskeller/packet/v4/message"
)

// testTimeout bounds every wait in these tests, so that they fail rather than
// hang if the code under test is broken.
const testTimeout = 10 * time.Second

func TestLineWriterConcurrent(t *testing.T) {
	var (
		mu    sync.Mutex
		lines []string
		wg    sync.WaitGroup
	)
	lw := &lineWriter{fn: func(s string) {
		mu.Lock()
		lines = append(lines, s)
		mu.Unlock()
	}}
	const writers, perWriter = 8, 200
	for w := range writers {
		wg.Go(func() {
			for n := range perWriter {
				// Write each line in two pieces, so that a lack of
				// locking would tear lines.
				line := fmt.Sprintf("writer %d line %d", w, n)
				lw.Write([]byte(line))
				lw.Write([]byte("\r\n"))
			}
		})
	}
	wg.Wait()
	lw.Write([]byte("partial"))
	lw.Flush()
	lw.Flush() // no-op
	if len(lines) != writers*perWriter+1 {
		t.Fatalf("got %d lines, want %d", len(lines), writers*perWriter+1)
	}
	if lines[len(lines)-1] != "partial" {
		t.Errorf("last line %q, want %q", lines[len(lines)-1], "partial")
	}
	// Writes are atomic, so every line is either a whole line or (when
	// two writers' halves interleave, which the transports never do) a
	// concatenation; but with locking, nothing is ever lost or garbled.
	var total int
	for _, l := range lines[:len(lines)-1] {
		total += strings.Count(l, "writer ")
	}
	if total != writers*perWriter {
		t.Errorf("found %d line texts, want %d", total, writers*perWriter)
	}
}

// fakeTransport is a jnos.Transport that simulates a JNOS BBS.  onSend is
// called with each string sent to it, and returns the BBS's response and
// whether a read for anything beyond that response should block (until
// Abort) rather than report a disconnect.
type fakeTransport struct {
	mu      sync.Mutex
	onSend  func(string) (resp string, block bool)
	textlog io.Writer
	out     string
	block   bool
	sent    []string
	closed  bool
	aborted chan struct{}
	abort   sync.Once
	blocked chan struct{} // closed when a read first blocks
	blockOn sync.Once
	panicRd bool // panic (once) on a read that would block
}

func newFakeTransport(onSend func(string) (string, bool)) *fakeTransport {
	return &fakeTransport{
		onSend:  onSend,
		out:     "Welcome\n(#1) >\n",
		aborted: make(chan struct{}),
		blocked: make(chan struct{}),
	}
}

func (ft *fakeTransport) ReadUntil(s string) (string, error) {
	ft.mu.Lock()
	if idx := strings.Index(ft.out, s); idx >= 0 {
		data := ft.out[:idx+len(s)]
		ft.out = ft.out[idx+len(s):]
		tl := ft.textlog
		ft.mu.Unlock()
		if tl != nil {
			io.WriteString(tl, data)
		}
		return data, nil
	}
	if !ft.block {
		ft.mu.Unlock()
		return "", jnos.ErrDisconnected
	}
	if ft.panicRd {
		ft.panicRd, ft.block = false, false
		ft.mu.Unlock()
		panic("fake transport panic")
	}
	ft.mu.Unlock()
	ft.blockOn.Do(func() { close(ft.blocked) })
	select {
	case <-ft.aborted:
		return "", jnos.ErrDisconnected
	case <-time.After(testTimeout):
		return "", jnos.ErrTimeout
	}
}

func (ft *fakeTransport) Send(s string) error {
	select {
	case <-ft.aborted:
		return jnos.ErrDisconnected
	default:
	}
	ft.mu.Lock()
	defer ft.mu.Unlock()
	ft.sent = append(ft.sent, s)
	resp, block := ft.onSend(s)
	ft.out += resp
	ft.block = block
	return nil
}

func (ft *fakeTransport) Close() error {
	ft.mu.Lock()
	ft.closed = true
	ft.mu.Unlock()
	return nil
}

func (ft *fakeTransport) UseVerboseReads() bool { return false }

func (ft *fakeTransport) Abort() { ft.abort.Do(func() { close(ft.aborted) }) }

func (ft *fakeTransport) sentCommands() []string {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	return slices.Clone(ft.sent)
}

func (ft *fakeTransport) isClosed() bool {
	ft.mu.Lock()
	defer ft.mu.Unlock()
	return ft.closed
}

// recorder is a BBSExchangeWatcher (and BBSExchangeMonitor) that records what
// it's told.
type recorder struct {
	mu         sync.Mutex
	ft         *fakeTransport
	progress   []string
	errors     []string
	monitor    []string
	logentries int
	finishes   int
	// state observed at the first Finished call
	runningAtFinish  int
	closedAtFinish   bool
	finished         chan struct{}
	onProgress       func(string)
	reportedAfterFin bool
}

func newRecorder(ft *fakeTransport) *recorder {
	return &recorder{ft: ft, finished: make(chan struct{})}
}

func (r *recorder) Progress(s string) {
	r.mu.Lock()
	r.progress = append(r.progress, s)
	r.reportedAfterFin = r.reportedAfterFin || r.finishes != 0
	fn := r.onProgress
	r.mu.Unlock()
	if fn != nil {
		fn(s)
	}
}

func (r *recorder) LogEntry(le *LogEntry) {
	if le == nil {
		panic("LogEntry(nil)")
	}
	r.mu.Lock()
	r.logentries++
	r.mu.Unlock()
}

func (r *recorder) Error(s string) {
	r.mu.Lock()
	r.errors = append(r.errors, s)
	r.reportedAfterFin = r.reportedAfterFin || r.finishes != 0
	r.mu.Unlock()
}

func (r *recorder) Monitor(line string) {
	r.mu.Lock()
	r.monitor = append(r.monitor, line)
	r.mu.Unlock()
}

func (r *recorder) Finished() {
	connectionsRunningMutex.Lock()
	running := connectionsRunning.Len()
	connectionsRunningMutex.Unlock()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.finishes++
	if r.finishes == 1 {
		r.runningAtFinish = running
		r.closedAtFinish = r.ft == nil || r.ft.isClosed()
		close(r.finished)
	}
}

// newTestIncident creates an incident in a temporary directory.
func newTestIncident(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := Create(dir, func(i *Incident) error {
		i.Config.OpCall = "KC6RSC"
		i.Config.OpName = "Test Operator"
		i.Config.TxMessageID = "XSC-100P"
		i.Config.RxMessageID = "XSC-100M"
		i.Config.ConnectType = ConnectTelnet
		i.Config.ConnectBBS = "W1XSC"
		i.Config.TelnetUser = "kc6rsc"
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return dir
}

// useFakeTransport makes BBSExchange connect through ft for the rest of the
// test.
func useFakeTransport(t *testing.T, ft *fakeTransport) {
	saved := connectBBS
	connectBBS = func(e *exchange, _ context.Context, textlog io.Writer, monitor func(string)) (*jnos.Conn, error) {
		ft.mu.Lock()
		ft.textlog = textlog
		ft.mu.Unlock()
		return jnos.Connect(ft, "")
	}
	t.Cleanup(func() { connectBBS = saved })
}

// runExchange runs BBSExchange in a goroutine and waits (with a time limit)
// for it to return.  If before is not nil, it is called (in the test
// goroutine) after starting the exchange.
func runExchange(t *testing.T, ctx context.Context, dir string, rec *recorder, before func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		BBSExchange(ctx, dir, true, rec)
		close(done)
	}()
	if before != nil {
		before()
	}
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("BBSExchange did not return")
	}
	rec.mu.Lock()
	defer rec.mu.Unlock()
	if rec.finishes != 1 {
		t.Errorf("Finished called %d times, want 1", rec.finishes)
	}
	if rec.reportedAfterFin {
		t.Error("Progress or Error called after Finished")
	}
	if rec.runningAtFinish != 0 {
		t.Errorf("%d connectionsRunning entries still reserved at Finished", rec.runningAtFinish)
	}
	if !rec.closedAtFinish {
		t.Error("transport not closed before Finished")
	}
}

// bbs is a simple onSend function for a BBS with an empty mailbox.
func bbs(s string) (string, bool) {
	switch {
	case s == "XM 0\n":
		return "(#1) >\n", false
	case s == "LA\n":
		return "No messages\n(#1) >\n", false
	}
	return "", false // includes "B\r": the BBS disconnects
}

func TestBBSExchangeNormal(t *testing.T) {
	dir := newTestIncident(t)
	ft := newFakeTransport(bbs)
	useFakeTransport(t, ft)
	rec := newRecorder(ft)
	runExchange(t, context.Background(), dir, rec, nil)
	if len(rec.errors) != 0 {
		t.Errorf("unexpected errors %q", rec.errors)
	}
	if sent := ft.sentCommands(); !slices.Contains(sent, "B\r") {
		t.Errorf("no orderly disconnect; sent %q", sent)
	}
	if !slices.Contains(rec.monitor, "No messages") {
		t.Errorf("monitor lines %q missing BBS text", rec.monitor)
	}
}

func TestBBSExchangeAbortInterruptsBlockedCall(t *testing.T) {
	dir := newTestIncident(t)
	ft := newFakeTransport(func(s string) (string, bool) {
		if s == "LA\n" {
			return "", true // the BBS never answers
		}
		return bbs(s)
	})
	useFakeTransport(t, ft)
	rec := newRecorder(ft)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := time.Now()
	runExchange(t, ctx, dir, rec, func() {
		select {
		case <-ft.blocked:
		case <-time.After(testTimeout):
			t.Error("exchange never blocked")
		}
		cancel()
	})
	if elapsed := time.Since(start); elapsed > testTimeout/2 {
		t.Errorf("abort took %s", elapsed)
	}
	if !slices.Equal(rec.errors, []string{errAborted.Error()}) {
		t.Errorf("errors %q, want %q", rec.errors, errAborted.Error())
	}
	for _, s := range ft.sentCommands() {
		if s == "B\r" || strings.HasPrefix(s, "K") {
			t.Errorf("sent %q to the BBS after abort", s)
		}
	}
}

func TestBBSExchangePanicCleanup(t *testing.T) {
	dir := newTestIncident(t)
	ft := newFakeTransport(func(s string) (string, bool) {
		if s == "LA\n" {
			return "", true
		}
		return bbs(s)
	})
	ft.panicRd = true
	useFakeTransport(t, ft)
	rec := newRecorder(ft)
	runExchange(t, context.Background(), dir, rec, nil)
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "software error") {
		t.Errorf("errors %q, want one software error", rec.errors)
	}
	// A second exchange must be possible: nothing is left reserved.
	ft2 := newFakeTransport(bbs)
	useFakeTransport(t, ft2)
	rec2 := newRecorder(ft2)
	runExchange(t, context.Background(), dir, rec2, nil)
	if len(rec2.errors) != 0 {
		t.Errorf("second exchange errors %q", rec2.errors)
	}
}

func TestBBSExchangeSentMessageDeleted(t *testing.T) {
	dir := newTestIncident(t)
	var ident int
	if err := Write(dir, func(i *Incident) error {
		dm := message.PlainMessage.NewDraft().(*message.DraftMessage)
		dm.SetTo("W1XSC")
		dm.Subject().SetSubjectSummary("test message")
		dm.Subject().SetSubjectHandling("IMMEDIATE") // runExchange uses immOnly
		le, err := i.AddDraftMessage(dm)
		if err != nil {
			return err
		}
		le.Status = StatusQueued
		ident = le.Ident
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	var deleteErr error
	ft := newFakeTransport(func(s string) (string, bool) {
		switch {
		case strings.HasPrefix(s, "SP "):
			return "Subject:\n", false
		case strings.HasSuffix(s, "/EX\n"):
			// The user deletes the message while it's being sent.
			// (onSend runs on the exchange goroutine with no
			// incident lock held.)
			deleteErr = Write(dir, func(i *Incident) error {
				return i.DeleteMessage(i.GetLogEntryByIdent(ident))
			})
			return "Msg queued\n(#1) >\n", false
		case !strings.HasSuffix(s, "\r") && !strings.HasPrefix(s, "XM") && !strings.HasPrefix(s, "LA"):
			if strings.Contains(s, "test message") { // the subject line
				return "Enter message.  End with /EX or ^Z in first column (^A aborts):\n", false
			}
		}
		return bbs(s)
	})
	useFakeTransport(t, ft)
	rec := newRecorder(ft)
	runExchange(t, context.Background(), dir, rec, nil)
	if deleteErr != nil {
		t.Fatalf("deleting message: %s", deleteErr)
	}
	if len(rec.errors) != 0 {
		t.Errorf("unexpected errors %q", rec.errors)
	}
	if rec.logentries != 0 {
		t.Errorf("LogEntry called %d times for a deleted message", rec.logentries)
	}
	if err := Read(dir, func(i *Incident) error {
		if le := i.GetLogEntryByIdent(ident); le == nil || le.Status != StatusDeleted {
			t.Errorf("log entry %+v, want deleted", le)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
