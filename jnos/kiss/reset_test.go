package kiss

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

// fakeRadio simulates a Kenwood radio's serial port, for resetting its TNC.
type fakeRadio struct {
	mu       sync.Mutex
	kiss     bool   // the TNC is in KISS mode
	second   string // the second TN setting
	silent   bool   // never answer anything
	pending  []byte // data for Read
	got      []string
	readWait time.Duration
	closed   bool
	partial  []byte // bytes of an unfinished command (out of KISS mode)
}

func newFakeRadio(second string) *fakeRadio {
	return &fakeRadio{kiss: true, second: second}
}

func (r *fakeRadio) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if bytes.Equal(p, []byte{fend, 0xFF, fend}) {
		r.got = append(r.got, "KISS exit")
		if !r.kiss {
			// Out of KISS mode, these bytes start a (garbage)
			// command, as the real radio treats them.
			r.partial = append(r.partial, p...)
		}
		r.kiss = false
		return len(p), nil
	}
	cmd := strings.TrimSuffix(string(p), "\r")
	r.got = append(r.got, cmd)
	if r.kiss || r.silent {
		return len(p), nil // in KISS mode, CAT commands are ignored
	}
	garbage := len(r.partial) != 0
	r.partial = nil
	switch {
	case garbage:
		r.pending = append(r.pending, "?\r"...)
	case cmd == "TN":
		r.pending = append(r.pending, "TN 0,"+r.second+"\r"...)
	case strings.HasPrefix(cmd, "TN 2,"):
		r.kiss, r.second = true, strings.TrimPrefix(cmd, "TN 2,")
		r.pending = append(r.pending, cmd+"\r"...)
	default:
		r.pending = append(r.pending, "?\r"...)
	}
	return len(p), nil
}

func (r *fakeRadio) Read(p []byte) (int, error) {
	deadline := time.Now().Add(r.readTimeout())
	for {
		r.mu.Lock()
		if len(r.pending) != 0 {
			n := copy(p, r.pending)
			r.pending = r.pending[n:]
			r.mu.Unlock()
			return n, nil
		}
		r.mu.Unlock()
		if time.Now().After(deadline) {
			return 0, nil // timeout, as go.bug.st/serial reports it
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (r *fakeRadio) readTimeout() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.readWait
}

func (r *fakeRadio) SetReadTimeout(d time.Duration) error {
	r.mu.Lock()
	r.readWait = d
	r.mu.Unlock()
	return nil
}

func (r *fakeRadio) Close() error {
	r.mu.Lock()
	r.closed = true
	r.mu.Unlock()
	return nil
}

func (r *fakeRadio) commands() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.got...)
}

// fastReset shortens the reset timings for a test.
func fastReset(t *testing.T) {
	oldSettle, oldExit, oldReply := resetSettle, resetExitWait, resetReplyWait
	resetSettle, resetExitWait, resetReplyWait = 50*time.Millisecond, 20*time.Millisecond, 300*time.Millisecond
	t.Cleanup(func() { resetSettle, resetExitWait, resetReplyWait = oldSettle, oldExit, oldReply })
}

func TestKenwoodReset(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("1")
	cmd, err := kenwoodReset(radio)
	if err != nil {
		t.Fatal(err)
	}
	// The second TN setting the radio reported is kept.
	if cmd != "TN 2,1" {
		t.Errorf("re-entered KISS with %q, want %q", cmd, "TN 2,1")
	}
	if want := []string{"KISS exit", "", "TN", "TN 2,1"}; strings.Join(radio.commands(), "|") != strings.Join(want, "|") {
		t.Errorf("sent %q, want %q", radio.commands(), want)
	}
	if !radio.kiss {
		t.Error("the radio was left out of KISS mode")
	}
}

func TestKenwoodResetNoAnswer(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("0")
	radio.silent = true
	start := time.Now()
	if _, err := kenwoodReset(radio); !errors.Is(err, errNoAnswer) {
		t.Errorf("got %v, want no answer", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("took %v", elapsed)
	}
	// Without an answer to TN, it must not guess at re-entering KISS.
	for _, cmd := range radio.commands() {
		if strings.HasPrefix(cmd, "TN 2") {
			t.Errorf("sent %q though the radio didn't answer", cmd)
		}
	}
}

// withRadio returns an open option that resets TNCs through radio, counting
// how often the port was reopened.
func withRadio(radio *fakeRadio, reopened *int) func(*Transport) {
	var mu sync.Mutex
	return func(t *Transport) {
		t.reopen = func() (resetPort, error) {
			mu.Lock()
			*reopened++
			mu.Unlock()
			return radio, nil
		}
	}
}

func TestResetAfterFailedConnect(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("0")
	var reopened int
	var mu sync.Mutex
	var lines []string
	monitor := func(line string) { mu.Lock(); lines = append(lines, line); mu.Unlock() }
	dev := *testDevice
	dev.Frack, dev.Retry = 100*time.Millisecond, 1
	client, server := net.Pipe()
	defer server.Close()
	b := newFakeBBS(t, server)
	go func() {
		for range b.frames { // the BBS never answers
		}
	}()
	if _, err := open(client, &dev, testClient, testBBS, nil, monitor, withRadio(radio, &reopened)); err == nil {
		t.Fatal("connect succeeded")
	}
	if reopened != 1 || strings.Join(radio.commands(), "|") != "KISS exit||TN|TN 2,0" {
		t.Errorf("reopened %d times and sent %q; want a reset", reopened, radio.commands())
	}
	mu.Lock()
	defer mu.Unlock()
	if !strings.Contains(strings.Join(lines, "\n"), "TNC reset: left KISS mode") {
		t.Errorf("the monitor doesn't report the reset:\n%s", strings.Join(lines, "\n"))
	}
}

func TestNoResetAfterCleanEnd(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("0")
	var reopened int
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go func() {
		b.accept()
		b.write("Welcome\n(#1) >\n")
		b.waitFor("B")
		b.disconnect()
	}()
	tr, err := open(client, testDevice, testClient, testBBS, nil, nil, withRadio(radio, &reopened))
	if err != nil {
		t.Fatal(err)
	}
	tr.ReadUntil("(#1) >\n")
	tr.Send("B")
	tr.ReadUntil("\n")
	tr.Close()
	if reopened != 0 {
		t.Errorf("reset the TNC (sent %q) after a clean disconnection", radio.commands())
	}
}

func TestResetAfterDisconnectWithData(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("0")
	var reopened int
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go func() {
		b.accept()
		// The BBS never hears our data, and gives up on us.
		for f := range b.frames {
			if f.isI() {
				b.send(&frame{command: true, control: ctlDISC | pfBit})
				break
			}
		}
		for range b.frames {
		}
	}()
	tr, err := open(client, testDevice, testClient, testBBS, nil, nil, withRadio(radio, &reopened))
	if err != nil {
		t.Fatal(err)
	}
	tr.Send("SP PKTMON")
	tr.ReadUntil("\n")
	tr.Close()
	if reopened != 1 {
		t.Errorf("didn't reset the TNC after the BBS disconnected with our data unacknowledged")
	}
}

// TestKenwoodResetOutOfKISS checks the reset of a radio whose TNC isn't in
// KISS mode: it takes the KISS exit bytes as the start of a command, and
// would answer the TN query with "?" (as happened at 09:02:46).
func TestKenwoodResetOutOfKISS(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("0")
	radio.kiss = false
	cmd, err := kenwoodReset(radio)
	if err != nil {
		t.Fatal(err)
	}
	if cmd != "TN 2,0" || !radio.kiss {
		t.Errorf("re-entered KISS with %q (in KISS: %v)", cmd, radio.kiss)
	}
}

// TestNoResetAfterLogoffFullDuplex checks the usual end of a session with a
// full-duplex TNC: we send "B", and JNOS disconnects without acknowledging it.
// The TNC holds nothing, so there's nothing to reset.
func TestNoResetAfterLogoffFullDuplex(t *testing.T) {
	fastReset(t)
	radio := newFakeRadio("0")
	var reopened int
	dev := *testDevice
	dev.FullDuplex = true
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go func() {
		b.accept()
		for f := range b.frames {
			if f.isI() { // our "B": disconnect without acknowledging it
				b.send(&frame{command: true, control: ctlDISC | pfBit})
				break
			}
		}
		for range b.frames {
		}
	}()
	tr, err := open(client, &dev, testClient, testBBS, nil, nil, withRadio(radio, &reopened))
	if err != nil {
		t.Fatal(err)
	}
	tr.Send("B")
	tr.ReadUntil("\n")
	tr.Close()
	if reopened != 0 {
		t.Errorf("reset the TNC (sent %q) after the BBS's ordinary disconnection", radio.commands())
	}
}
