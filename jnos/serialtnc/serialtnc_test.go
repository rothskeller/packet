package serialtnc

import (
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"go.bug.st/serial"

	"github.com/rothskeller/packet/v4/jnos"
	"github.com/rothskeller/packet/v4/jnos/tnc"
)

// errPortClosed is returned by the fake port after it is closed.
var errPortClosed = errors.New("fake port closed")

// fakePort is an in-memory serial port with a simulated command-mode TNC on
// the other end.  It mimics go.bug.st/serial: a Read that times out returns
// (0, nil), a negative timeout other than NoTimeout is rejected, and Read and
// Write return errors after Close.
type fakePort struct {
	mu       sync.Mutex
	in       []byte
	wake     chan struct{}
	closedCh chan struct{}
	closed   bool
	timeout  time.Duration
	timeouts []time.Duration // every timeout successfully set
	line     []byte
	lines    []string // complete lines written (Ctrl-C recorded as "^C")
	noEcho   bool     // if set, the TNC doesn't echo
	// slowEcho, if set, makes the echo arrive in 8-byte chunks with this
	// delay between them.
	slowEcho time.Duration
	// respond returns the TNC output (after the echo) in response to a
	// complete line (or "^C").
	respond func(line string) string
}

func newFakePort(respond func(string) string) *fakePort {
	if respond == nil {
		respond = func(string) string { return "" }
	}
	return &fakePort{
		wake:     make(chan struct{}, 1),
		closedCh: make(chan struct{}),
		timeout:  serial.NoTimeout,
		respond:  respond,
	}
}

func (f *fakePort) feed(s string) {
	if s == "" {
		return
	}
	f.mu.Lock()
	if !f.closed {
		f.in = append(f.in, s...)
	}
	f.mu.Unlock()
	select {
	case f.wake <- struct{}{}:
	default:
	}
}

func (f *fakePort) Read(p []byte) (int, error) {
	f.mu.Lock()
	to := f.timeout
	f.mu.Unlock()
	var timer <-chan time.Time
	if to != serial.NoTimeout {
		tm := time.NewTimer(to)
		defer tm.Stop()
		timer = tm.C
	}
	for {
		f.mu.Lock()
		if f.closed {
			f.mu.Unlock()
			return 0, errPortClosed
		}
		if len(f.in) != 0 {
			n := copy(p, f.in)
			f.in = f.in[n:]
			f.mu.Unlock()
			return n, nil
		}
		f.mu.Unlock()
		select {
		case <-f.wake:
		case <-f.closedCh:
		case <-timer:
			return 0, nil
		}
	}
}

func (f *fakePort) Write(p []byte) (int, error) {
	var echo, out strings.Builder

	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return 0, errPortClosed
	}
	var lines []string
	for _, b := range p {
		switch b {
		case 3:
			lines = append(lines, "^C")
		case '\r':
			echo.WriteString("\r\n")
			lines = append(lines, string(f.line))
			f.line = f.line[:0]
		default:
			echo.WriteByte(b)
			f.line = append(f.line, b)
		}
	}
	f.lines = append(f.lines, lines...)
	noEcho, slowEcho, respond := f.noEcho, f.slowEcho, f.respond
	f.mu.Unlock()
	for _, line := range lines {
		out.WriteString(respond(line))
	}
	switch {
	case noEcho:
		f.feed(out.String())
	case slowEcho != 0:
		go func(e, o string) {
			for len(e) != 0 {
				n := min(8, len(e))
				time.Sleep(slowEcho)
				f.feed(e[:n])
				e = e[n:]
			}
			f.feed(o)
		}(echo.String(), out.String())
	default:
		f.feed(echo.String() + out.String())
	}
	return len(p), nil
}

func (f *fakePort) SetReadTimeout(t time.Duration) error {
	if t < 0 && t != serial.NoTimeout {
		return errors.New("invalid timeout value")
	}
	f.mu.Lock()
	f.timeout = t
	f.timeouts = append(f.timeouts, t)
	f.mu.Unlock()
	return nil
}

func (f *fakePort) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.closed {
		f.closed = true
		close(f.closedCh)
	}
	return nil
}

func (f *fakePort) isClosed() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.closed
}

func (f *fakePort) writtenLines() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.lines...)
}

// tncSim returns a respond function simulating a command-mode TNC whose reply
// to the connect command is connectReply.
func tncSim(connectReply string) func(string) string {
	return func(line string) string {
		switch {
		case line == "^C":
			return "\r\ncmd:"
		case strings.HasPrefix(line, "C "):
			return connectReply
		case line == "D":
			return "\r\n*** DISCONNECTED\r\ncmd:"
		case line == "B": // sent to the BBS in converse mode
			return ""
		default:
			return "cmd:"
		}
	}
}

func testTNC() *tnc.TNC {
	return &tnc.TNC{
		Name:                   "Test TNC",
		CommandPrompt:          "cmd:",
		ConnectedMessage:       "*** CONNECTED",
		DisconnectedMessage:    "*** DISCONNECTED",
		PreConnectCommands:     []string{"MONITOR OFF"},
		PostDisconnectCommands: []string{"CR ON"},
		MyCallCommand:          "MY",
		ConnectCommand:         "C",
		EchoTimeout:            50,
		ReplyTimeout:           100,
	}
}

// setRFTimeout shortens rfTimeout for the duration of the test.
func setRFTimeout(t *testing.T, d time.Duration) {
	old := rfTimeout
	rfTimeout = d
	t.Cleanup(func() { rfTimeout = old })
}

// within runs fn, failing the test if it doesn't finish within limit.  On
// failure it calls unblock (which should close the fake port) so that the
// goroutine running fn can exit.
func within(t *testing.T, limit time.Duration, unblock func(), fn func()) time.Duration {
	t.Helper()
	done := make(chan struct{})
	start := time.Now()
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
		return time.Since(start)
	case <-time.After(limit):
		unblock()
		<-done
		t.Fatalf("did not finish within %s", limit)
		return 0
	}
}

func contains(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}

// trickle feeds a beacon to the port every few milliseconds until the
// returned stop function is called.
func trickle(fp *fakePort) (stop func()) {
	quit := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		tick := time.NewTicker(5 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-quit:
				return
			case <-tick.C:
				fp.feed("KJ6XYZ>BEACON: monitor traffic\r\n")
			}
		}
	}()
	return func() { close(quit); <-done }
}

func TestNormalSession(t *testing.T) {
	setRFTimeout(t, 5*time.Second)
	fp := newFakePort(tncSim("\r\n*** CONNECTED to BBS-1\r\n[JNOS-2.0]\r\n(#1) >"))
	var (
		tr                  *Transport
		err, rerr, serr, ce error
		s1, s2              string
	)
	within(t, 10*time.Second, func() { fp.Close() }, func() {
		if tr, err = start(testTNC(), fp, "BBS-1", "MBOX", nil); err != nil {
			return
		}
		s1, rerr = tr.ReadUntil(">")
		serr = tr.Send("B")
		fp.feed("\r\n*** DISCONNECTED\r\ncmd:") // BBS hangs up
		s2, _ = tr.ReadUntil("\n")
		_, err2 := tr.ReadUntil("\n")
		if !errors.Is(err2, jnos.ErrDisconnected) {
			rerr = err2
		}
		ce = tr.Close()
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if rerr != nil || serr != nil || ce != nil {
		t.Fatalf("ReadUntil %v, Send %v, Close %v", rerr, serr, ce)
	}
	if !strings.HasSuffix(s1, "(#1) >") {
		t.Errorf("ReadUntil returned %q", s1)
	}
	if s2 != "\n" {
		t.Errorf("second ReadUntil returned %q", s2)
	}
	lines := fp.writtenLines()
	for _, want := range []string{"MONITOR OFF", "MY MBOX", "C BBS-1", "B", "CR ON"} {
		if !contains(lines, want) {
			t.Errorf("TNC never received %q; got %q", want, lines)
		}
	}
	if contains(lines, "D") || contains(lines, "^C") {
		t.Errorf("unexpected D or ^C after orderly disconnect: %q", lines)
	}
	if !fp.isClosed() {
		t.Error("port not closed")
	}
}

func TestTrickleTimesOut(t *testing.T) {
	fp := newFakePort(nil)
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	stop := trickle(fp)
	defer stop()
	var err error
	elapsed := within(t, 5*time.Second, func() { fp.Close() }, func() {
		_, err = tr.readUntil("NEVER", 200*time.Millisecond)
	})
	if err != jnos.ErrTimeout {
		t.Errorf("readUntil error %v, want ErrTimeout", err)
	}
	if elapsed > time.Second {
		t.Errorf("readUntil took %s with a 200ms allowance", elapsed)
	}
	// Same through the exported ReadUntil.
	setRFTimeout(t, 200*time.Millisecond)
	elapsed = within(t, 5*time.Second, func() { fp.Close() }, func() {
		_, err = tr.ReadUntil("NEVER")
	})
	if err != jnos.ErrTimeout {
		t.Errorf("ReadUntil error %v, want ErrTimeout", err)
	}
	if elapsed > time.Second {
		t.Errorf("ReadUntil took %s with a 200ms allowance", elapsed)
	}
}

func TestEchoTrickleTimesOut(t *testing.T) {
	fp := newFakePort(nil)
	fp.noEcho = true
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	stop := trickle(fp)
	defer stop()
	var err error
	elapsed := within(t, 5*time.Second, func() { fp.Close() }, func() {
		err = tr.Send("hello")
	})
	if err != ErrBadEcho {
		t.Errorf("Send error %v, want ErrBadEcho", err)
	}
	if elapsed > time.Second {
		t.Errorf("Send took %s with a ~60ms echo allowance", elapsed)
	}
}

func TestSlowEchoSucceeds(t *testing.T) {
	// A full 256-byte block echoed at about 9600 baud (8 bytes every 8ms)
	// must not be treated as a bad echo with the default echo timeout.
	def := testTNC()
	def.EchoTimeout = 0
	fp := newFakePort(nil)
	fp.slowEcho = 8 * time.Millisecond
	tr := newTransport(def, fp, nil)
	tr.connected = true
	var err error
	within(t, 5*time.Second, func() { fp.Close() }, func() {
		err = tr.Send(strings.Repeat("x", 300))
	})
	if err != nil {
		t.Errorf("Send: %v", err)
	}
}

func TestNegativeTimeouts(t *testing.T) {
	def := testTNC()
	def.EchoTimeout = -5
	def.ReplyTimeout = -1
	if err := def.Validate(); err == nil {
		t.Error("Validate accepted negative timeouts")
	}
	// Open must reject the definition before touching the port.
	if _, err := Open(def, "", "BBS-1", nil); err == nil || !strings.Contains(err.Error(), "invalid TNC definition") {
		t.Errorf("Open error %v, want invalid TNC definition", err)
	}
	// Even if a negative timeout gets past validation, it is replaced by
	// the default, and a silent TNC times out rather than hanging.
	tr := newTransport(def, newFakePort(nil), nil)
	if tr.echoTimeout != echoTimeout || tr.replyTimeout != tncTimeout {
		t.Errorf("timeouts %s/%s, want defaults", tr.echoTimeout, tr.replyTimeout)
	}
	fp := newFakePort(nil)
	fp.noEcho = true
	var err error
	within(t, 10*time.Second, func() { fp.Close() }, func() {
		_, err = start(def, fp, "BBS-1", "MBOX", nil)
	})
	if err == nil || !strings.Contains(err.Error(), "initial prompt") {
		t.Errorf("start error %v, want initial prompt timeout", err)
	}
	fp.mu.Lock()
	defer fp.mu.Unlock()
	if len(fp.timeouts) == 0 {
		t.Error("no read timeout was ever set")
	}
	for _, to := range fp.timeouts {
		if to <= 0 {
			t.Errorf("read timeout set to %s", to)
		}
	}
}

func TestConnectRefusedFailsFast(t *testing.T) {
	setRFTimeout(t, 30*time.Second)
	fp := newFakePort(tncSim("\r\n*** retry count exceeded\r\n*** DISCONNECTED\r\ncmd:"))
	var err error
	elapsed := within(t, 10*time.Second, func() { fp.Close() }, func() {
		_, err = start(testTNC(), fp, "BBS-1", "MBOX", nil)
	})
	if err == nil || !strings.Contains(err.Error(), "retry count exceeded") {
		t.Errorf("start error %v, want connection failure", err)
	}
	if elapsed > 3*time.Second {
		t.Errorf("failed connect took %s", elapsed)
	}
	lines := fp.writtenLines()
	if contains(lines, "D") {
		t.Errorf("sent D although the TNC reported the connect failed: %q", lines)
	}
	if !contains(lines, "CR ON") {
		t.Errorf("post-disconnect commands not sent: %q", lines)
	}
	if !fp.isClosed() {
		t.Error("port not closed")
	}
}

func TestCloseAfterFailedConnectSendsD(t *testing.T) {
	setRFTimeout(t, 300*time.Millisecond)
	fp := newFakePort(tncSim("")) // connect never answered
	var err error
	within(t, 10*time.Second, func() { fp.Close() }, func() {
		_, err = start(testTNC(), fp, "BBS-1", "MBOX", nil)
	})
	if err == nil || !strings.Contains(err.Error(), jnos.ErrTimeout.Error()) {
		t.Errorf("start error %v, want timeout", err)
	}
	lines := fp.writtenLines()
	var i int
	for i = 0; i < len(lines) && lines[i] != "C BBS-1"; i++ {
	}
	after := lines[min(i, len(lines)):]
	if !contains(after, "^C") || !contains(after, "D") {
		t.Errorf("Close after failed connect didn't send ^C and D: %q", lines)
	}
	if !contains(after, "CR ON") {
		t.Errorf("post-disconnect commands not sent: %q", lines)
	}
	if !fp.isClosed() {
		t.Error("port not closed")
	}
}

func TestAbortDuringReadUntil(t *testing.T) {
	setRFTimeout(t, 30*time.Second)
	fp := newFakePort(tncSim("\r\n*** CONNECTED to BBS-1\r\n"))
	var (
		tr  *Transport
		err error
	)
	within(t, 10*time.Second, func() { fp.Close() }, func() {
		tr, err = start(testTNC(), fp, "BBS-1", "MBOX", nil)
	})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	stop := trickle(fp) // keep data arriving so the race detector sees overlap
	defer stop()
	nlines := len(fp.writtenLines())
	var rerr error
	elapsed := within(t, 5*time.Second, func() { fp.Close() }, func() {
		go func() {
			time.Sleep(50 * time.Millisecond)
			tr.Abort()
			tr.Abort() // must be harmless
		}()
		_, rerr = tr.ReadUntil("NEVER")
	})
	if rerr == nil || rerr == jnos.ErrTimeout {
		t.Errorf("ReadUntil after Abort returned %v", rerr)
	}
	if elapsed > 2*time.Second {
		t.Errorf("ReadUntil took %s to notice Abort", elapsed)
	}
	if !fp.isClosed() {
		t.Error("Abort didn't close the port")
	}
	if _, err = tr.ReadUntil("x"); err == nil {
		t.Error("ReadUntil after Abort succeeded")
	}
	if err = tr.Send("x"); err == nil {
		t.Error("Send after Abort succeeded")
	}
	elapsed = within(t, 5*time.Second, func() { fp.Close() }, func() { tr.Close() })
	if elapsed > 500*time.Millisecond {
		t.Errorf("Close after Abort took %s", elapsed)
	}
	if lines := fp.writtenLines(); len(lines) != nlines {
		t.Errorf("TNC chatter after Abort: %q", lines[nlines:])
	}
}

func TestAbortDuringSend(t *testing.T) {
	def := testTNC()
	def.EchoTimeout = 30000
	fp := newFakePort(nil)
	fp.noEcho = true
	tr := newTransport(def, fp, nil)
	tr.connected = true
	var err error
	elapsed := within(t, 5*time.Second, func() { fp.Close() }, func() {
		go func() {
			time.Sleep(50 * time.Millisecond)
			tr.Abort()
		}()
		err = tr.Send("hello")
	})
	if err == nil {
		t.Error("Send after Abort succeeded")
	}
	if elapsed > 2*time.Second {
		t.Errorf("Send took %s to notice Abort", elapsed)
	}
	within(t, 2*time.Second, func() { fp.Close() }, func() { tr.Close() })
}

func TestAbortBeforeAnything(t *testing.T) {
	fp := newFakePort(nil)
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	tr.Abort()
	within(t, 2*time.Second, func() { fp.Close() }, func() {
		if _, err := tr.ReadUntil("x"); err == nil {
			t.Error("ReadUntil succeeded")
		}
		tr.Close()
		tr.Abort()
	})
}

func TestPendingCap(t *testing.T) {
	fp := newFakePort(nil)
	tr := newTransport(testTNC(), fp, nil)
	tr.connected = true
	fp.feed(strings.Repeat("x", maxPending+4096))
	var err error
	within(t, 10*time.Second, func() { fp.Close() }, func() {
		_, err = tr.readUntil("NEVER", 5*time.Second)
	})
	if err != errOverflow {
		t.Errorf("readUntil error %v, want errOverflow", err)
	}
}

func TestCheckPendingEarliest(t *testing.T) {
	tr := newTransport(testTNC(), newFakePort(nil), nil)
	tr.pending = []byte("abc*** DISCONNECTED\r\ncmd:*** CONNECTED")
	data, which := tr.checkPending([][]byte{[]byte("*** CONNECTED"), []byte("*** DISCONNECTED")})
	if which != 1 || data != "abc*** DISCONNECTED" {
		t.Errorf("checkPending = %q, %d", data, which)
	}
	if data, which = tr.checkPending([][]byte{nil}); which != -1 || data != "" {
		t.Errorf("empty target matched: %q, %d", data, which)
	}
}
