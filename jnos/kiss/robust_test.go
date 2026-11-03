package kiss

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
)

// These tests cover the link's behavior on bad paths and with misbehaving
// devices.

// shortTimeouts shortens the package's timeouts for the duration of a test.
// (They are read only by open, so connections opened before the test, whose
// goroutines may still be running, aren't affected.)
func shortTimeouts(t *testing.T, read, write, portClose time.Duration) {
	oldRead, oldWrite, oldClose := readTimeout, minWriteTimeout, portCloseTimeout
	readTimeout, minWriteTimeout, portCloseTimeout = read, write, portClose
	t.Cleanup(func() { readTimeout, minWriteTimeout, portCloseTimeout = oldRead, oldWrite, oldClose })
}

// within runs fn, failing the test if it doesn't return within limit.
func within(t *testing.T, limit time.Duration, what string, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		fn()
	}()
	select {
	case <-done:
	case <-time.After(limit):
		t.Fatalf("%s did not return within %v", what, limit)
	}
}

// connectTest opens a Transport to a fakeBBS, which accepts the connection.
// The Transport is aborted and closed when the test ends.
func connectTest(t *testing.T, dev *Device, port func(net.Conn) io.ReadWriteCloser, monitor func(string)) (*Transport, *fakeBBS) {
	t.Helper()
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go b.accept()
	var rwc io.ReadWriteCloser = client
	if port != nil {
		rwc = port(client)
	}
	tr, err := open(rwc, dev, testClient, testBBS, nil, monitor)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		tr.Abort()
		tr.Close()
		server.Close()
	})
	return tr, b
}

// next returns the next frame from the client, failing the test if there is
// none within limit.
func next(t *testing.T, b *fakeBBS, limit time.Duration, what string) *frame {
	t.Helper()
	select {
	case f, ok := <-b.frames:
		if !ok {
			t.Fatalf("%s: connection closed", what)
		}
		return f
	case <-time.After(limit):
		t.Fatalf("%s: no frame from client within %v", what, limit)
	}
	return nil
}

func isPoll(f *frame) bool            { return (f.isS() || f.isI()) && f.command && f.pf() }
func isU(f *frame, u byte) bool       { return !f.isI() && !f.isS() && f.uType() == u }
func isIFrame(f *frame, ns byte) bool { return f.isI() && f.ns() == ns }

// stallPort is a port whose writes (once stall is set) and close (once
// blockClose is set) block until release is closed, like a stalled USB or
// Bluetooth device.
type stallPort struct {
	net.Conn
	stall, blockClose atomic.Bool
	release           chan struct{}
}

func newStallPort(t *testing.T) func(net.Conn) io.ReadWriteCloser {
	return func(c net.Conn) io.ReadWriteCloser {
		p := &stallPort{Conn: c, release: make(chan struct{})}
		t.Cleanup(func() { close(p.release) })
		lastStallPort = p
		return p
	}
}

// lastStallPort is the stallPort most recently created by newStallPort.  (Tests
// don't run in parallel.)
var lastStallPort *stallPort

func (p *stallPort) Write(b []byte) (int, error) {
	if p.stall.Load() {
		<-p.release
		return 0, io.ErrClosedPipe
	}
	return p.Conn.Write(b)
}

func (p *stallPort) Close() error {
	if p.blockClose.Load() {
		<-p.release
	}
	return p.Conn.Close()
}

// TestPollLivelock checks that a path on which our polls get through, but our
// I frames don't, ends in failure, rather than polling and resending forever.
// The BBS answers each poll with "RR F R0" and keeps the channel busy, but
// never receives any data.
func TestPollLivelock(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Retry, dev.PacLen, dev.MaxFrame = 100*time.Millisecond, 3, 8, 1
	tr, b := connectTest(t, &dev, nil, nil)
	stop := make(chan struct{})
	defer close(stop)
	var polls atomic.Int32
	go func() {
		tick := time.NewTicker(50 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case f, ok := <-b.frames:
				if !ok {
					return
				}
				if isPoll(f) {
					polls.Add(1)
					b.sendS(ctlRR, true) // V(R) stays 0
				}
			case <-tick.C:
				b.sendS(ctlRR, false)
			case <-stop:
				return
			}
		}
	}()
	if err := tr.Send("hello there"); err != nil {
		t.Fatal(err)
	}
	var err error
	within(t, 10*time.Second, "ReadUntil", func() { _, err = tr.ReadUntil("\n") })
	if err == nil {
		t.Fatal("ReadUntil succeeded")
	}
	if got := polls.Load(); got != int32(dev.Retry) {
		t.Errorf("client polled %d times, expected %d", got, dev.Retry)
	}
}

// TestProgressTimeout checks that ReadUntil times out when the link makes no
// progress, even though the BBS keeps sending (supervisory) frames.
func TestProgressTimeout(t *testing.T) {
	shortTimeouts(t, 500*time.Millisecond, minWriteTimeout, portCloseTimeout)
	tr, b := connectTest(t, testDevice, nil, nil)
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		tick := time.NewTicker(100 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case f, ok := <-b.frames:
				if !ok {
					return
				}
				b.handle(f)
			case <-tick.C:
				b.send(&frame{command: true, control: ctlRR | pfBit | b.vr<<5})
			case <-stop:
				return
			}
		}
	}()
	if err := tr.Send("command"); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	var err error
	within(t, 5*time.Second, "ReadUntil", func() { _, err = tr.ReadUntil("never") })
	if err != jnos.ErrTimeout {
		t.Errorf("got %v, expected timeout", err)
	}
	if elapsed := time.Since(start); elapsed < 400*time.Millisecond {
		t.Errorf("timed out after only %v", elapsed)
	}
}

// TestREJDuringRecovery checks that a REJ received while we're waiting for
// the answer to a poll doesn't stop the link: T1 keeps running, and an N(R)
// acknowledging frames sent before the REJ is still accepted.
func TestREJDuringRecovery(t *testing.T) {
	dev := *testDevice
	dev.PacLen = 8
	tr, b := connectTest(t, &dev, nil, nil)
	tr.Send("0123456789ABCDE") // two I frames, with the CR
	for ns := range byte(2) {
		if f := next(t, b, 2*time.Second, "I frame"); !isIFrame(f, ns) {
			t.Fatalf("got control %02X, expected I frame S%d", f.control, ns)
		}
	}
	// The BBS stays silent, so the client polls.
	if f := next(t, b, 3*time.Second, "first poll"); !isPoll(f) {
		t.Fatalf("got control %02X, expected poll", f.control)
	}
	// A REJ R0 (without F) arrives, and the answer to the poll is lost.
	b.send(&frame{control: ctlREJ})
	if f := next(t, b, 3*time.Second, "second poll"); !isPoll(f) {
		t.Fatalf("got control %02X, expected poll", f.control)
	}
	// This time the BBS answers: it has both frames after all.
	b.vr = 2
	b.sendS(ctlRR, true)
	tr.Send("next")
	for {
		f := next(t, b, 3*time.Second, "next I frame")
		if f.isI() {
			if f.ns() != 2 || string(f.info) != "next\r" {
				t.Fatalf("got I frame S%d %q, expected S2 \"next\\r\"", f.ns(), f.info)
			}
			break
		}
	}
}

// TestRNRThenRR checks that the client stops sending while the BBS is busy,
// and resumes when it's ready again.
func TestRNRThenRR(t *testing.T) {
	tr, b := connectTest(t, testDevice, nil, nil)
	tr.Send("hello")
	if f := next(t, b, 2*time.Second, "I frame"); !isIFrame(f, 0) {
		t.Fatalf("got control %02X, expected I frame S0", f.control)
	}
	// The BBS acknowledges it but is busy.  It polls, so that we know when
	// the client has heard.
	b.vr = 1
	b.send(&frame{command: true, control: ctlRNR | pfBit | b.vr<<5})
	if f := next(t, b, 2*time.Second, "poll answer"); !f.isS() || f.command || !f.pf() {
		t.Fatalf("got control %02X, expected response with F", f.control)
	}
	tr.Send("world")
	timeout := time.After(150 * time.Millisecond)
quiet:
	for {
		select {
		case f := <-b.frames:
			if f.isI() {
				t.Fatal("client sent I frame while BBS busy")
			}
		case <-timeout:
			break quiet
		}
	}
	b.sendS(ctlRR, false)
	for {
		f := next(t, b, 2*time.Second, "I frame after RR")
		if isPoll(f) {
			b.sendS(ctlRR, true)
		} else if f.isI() {
			if f.ns() != 1 || string(f.info) != "world\r" {
				t.Fatalf("got I frame S%d %q, expected S1 \"world\\r\"", f.ns(), f.info)
			}
			break
		}
	}
}

// TestRNRLostRR checks that, while the BBS is busy and we have data for it,
// the client polls it, so that a lost RR doesn't stop the link forever; and
// that it keeps polling after an "RNR F" answer.
func TestRNRLostRR(t *testing.T) {
	tr, b := connectTest(t, testDevice, nil, nil)
	tr.Send("hello")
	if f := next(t, b, 2*time.Second, "I frame"); !isIFrame(f, 0) {
		t.Fatalf("got control %02X, expected I frame S0", f.control)
	}
	b.vr = 1
	b.send(&frame{command: true, control: ctlRNR | pfBit | b.vr<<5})
	if f := next(t, b, 2*time.Second, "poll answer"); !f.isS() || f.command || !f.pf() {
		t.Fatalf("got control %02X, expected response with F", f.control)
	}
	tr.Send("world")
	// The BBS's RR is lost; the client must poll.
	if f := next(t, b, 2*time.Second, "first poll"); !isPoll(f) {
		t.Fatalf("got control %02X, expected poll", f.control)
	}
	b.sendS(ctlRNR, true) // still busy
	if f := next(t, b, 2*time.Second, "second poll"); !isPoll(f) {
		t.Fatalf("got control %02X, expected poll", f.control)
	}
	b.sendS(ctlRR, true) // ready now
	f := next(t, b, 2*time.Second, "I frame after RR")
	if !f.isI() || f.ns() != 1 || string(f.info) != "world\r" {
		t.Fatalf("got control %02X %q, expected I frame S1 \"world\\r\"", f.control, f.info)
	}
}

// TestSABMReset checks that when the BBS resets the link (while the client is
// in timer recovery, with data outstanding), the client resets all of its
// link state and the link stays usable.
func TestSABMReset(t *testing.T) {
	var (
		mu    sync.Mutex
		lines []string
	)
	monitor := func(line string) {
		mu.Lock()
		lines = append(lines, line)
		mu.Unlock()
	}
	tr, b := connectTest(t, testDevice, nil, monitor)
	tr.Send("one")
	if f := next(t, b, 2*time.Second, "I frame"); !isIFrame(f, 0) {
		t.Fatalf("got control %02X, expected I frame S0", f.control)
	}
	if f := next(t, b, 3*time.Second, "poll"); !isPoll(f) {
		t.Fatalf("got control %02X, expected poll", f.control)
	}
	b.send(&frame{command: true, control: ctlSABM | pfBit})
	if f := next(t, b, 2*time.Second, "UA"); !isU(f, ctlUA) || f.command || !f.pf() {
		t.Fatalf("got control %02X, expected UA F", f.control)
	}
	tr.Send("two")
	var f *frame
	for f = next(t, b, 2*time.Second, "I frame after reset"); !f.isI(); f = next(t, b, 2*time.Second, "I frame after reset") {
		if isPoll(f) {
			t.Fatal("client still polling after reset")
		}
	}
	if f.ns() != 0 || f.nr() != 0 || string(f.info) != "two\r" {
		t.Fatalf("got I frame R%d S%d %q, expected R0 S0 \"two\\r\"", f.nr(), f.ns(), f.info)
	}
	b.handle(f)
	b.waitFor("two")
	wrote := make(chan struct{})
	go func() {
		b.write("ok\n")
		close(wrote)
	}()
	within(t, 3*time.Second, "ReadUntil", func() {
		if s, err := tr.ReadUntil("ok\n"); err != nil || s != "ok\n" {
			t.Errorf("got %q, %v", s, err)
		}
	})
	within(t, 3*time.Second, "BBS write", func() { <-wrote })
	mu.Lock()
	defer mu.Unlock()
	if !strings.Contains(strings.Join(lines, "\n"), "link reset by BBS; 4 bytes of unacknowledged data lost") {
		t.Errorf("monitor didn't report the lost data; got:\n%s", strings.Join(lines, "\n"))
	}
}

// TestNoSpinWhileDisconnecting checks that the link goroutine doesn't spin
// while waiting for the BBS to answer our DISC with an acknowledgment owed.
func TestNoSpinWhileDisconnecting(t *testing.T) {
	dev := *testDevice
	dev.RespTime = 300 * time.Millisecond
	tr, b := connectTest(t, &dev, nil, nil)
	b.send(&frame{command: true, control: 0, info: []byte("hello\r")})
	if s, err := tr.ReadUntil("\n"); err != nil || s != "hello\n" {
		t.Fatalf("got %q, %v", s, err)
	}
	// The client owes an acknowledgment, but disconnects first.  The BBS
	// never answers the DISC.
	var discs atomic.Int32
	go func() {
		for f := range b.frames {
			if isU(f, ctlDISC) {
				discs.Add(1)
			}
		}
	}()
	before := tr.wakeups.Load()
	within(t, 5*time.Second, "Close", func() { tr.Close() })
	if loops := tr.wakeups.Load() - before; loops > 100 {
		t.Errorf("link loop ran %d times while disconnecting", loops)
	}
	if got := discs.Load(); got != maxDiscRetries+1 {
		t.Errorf("client sent %d DISCs, expected %d", got, maxDiscRetries+1)
	}
}

// TestStalledWrite checks that when writes to the TNC block forever, Send,
// ReadUntil and Close all return in bounded time (even though the port's
// Close blocks too), and that nothing is logged or monitored after Close.
func TestStalledWrite(t *testing.T) {
	shortTimeouts(t, readTimeout, 200*time.Millisecond, 200*time.Millisecond)
	var closed, lateOutput atomic.Bool
	monitor := func(string) {
		if closed.Load() {
			lateOutput.Store(true)
		}
	}
	tr, b := connectTest(t, testDevice, newStallPort(t), monitor)
	port := lastStallPort
	go func() {
		for range b.frames {
		}
	}()
	port.stall.Store(true)
	port.blockClose.Store(true)
	within(t, time.Second, "Send", func() {
		for range 5 {
			tr.Send("hello")
		}
	})
	var err error
	within(t, 3*time.Second, "ReadUntil", func() { _, err = tr.ReadUntil("\n") })
	if err != errWriteTimeout {
		t.Errorf("ReadUntil returned %v, expected %v", err, errWriteTimeout)
	}
	if err = tr.Send("more"); err == nil {
		t.Error("Send succeeded after link failure")
	}
	within(t, 2*time.Second, "Close", func() { err = tr.Close() })
	closed.Store(true)
	if err == nil {
		t.Error("Close succeeded although the port didn't close")
	}
	time.Sleep(100 * time.Millisecond)
	if lateOutput.Load() {
		t.Error("monitor called after Close returned")
	}
}

// TestStalledWriteAtConnect checks that open gives up when writes to the TNC
// block forever from the start.
func TestStalledWriteAtConnect(t *testing.T) {
	shortTimeouts(t, readTimeout, 200*time.Millisecond, 200*time.Millisecond)
	client, server := net.Pipe()
	defer server.Close()
	b := newFakeBBS(t, server)
	go func() {
		for range b.frames {
		}
	}()
	port := newStallPort(t)(client).(*stallPort)
	port.stall.Store(true)
	var err error
	within(t, 3*time.Second, "open", func() { _, err = open(port, testDevice, testClient, testBBS, nil, nil) })
	if err == nil || !strings.Contains(err.Error(), "timed out") {
		t.Errorf("got %v, expected write timeout", err)
	}
}

// TestAbort checks that Abort makes a blocked ReadUntil return promptly, makes
// later calls fail, sends a single DISC (so that the BBS can end its side of the
// session) without waiting for an answer, and that Close after it returns
// promptly.
func TestAbort(t *testing.T) {
	tr, b := connectTest(t, testDevice, nil, nil)
	var discs atomic.Int32
	drained := make(chan struct{})
	go func() {
		for f := range b.frames {
			if isU(f, ctlDISC) {
				discs.Add(1)
			}
		}
		close(drained)
	}()
	result := make(chan error, 1)
	go func() {
		_, err := tr.ReadUntil("never")
		result <- err
	}()
	time.Sleep(100 * time.Millisecond)
	tr.Abort()
	tr.Abort() // harmless
	select {
	case err := <-result:
		if !errors.Is(err, jnos.ErrDisconnected) {
			t.Errorf("ReadUntil returned %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("ReadUntil didn't return after Abort")
	}
	if err := tr.Send("more"); err == nil {
		t.Error("Send succeeded after Abort")
	}
	if _, err := tr.ReadUntil("x"); err == nil {
		t.Error("ReadUntil succeeded after Abort")
	}
	within(t, 500*time.Millisecond, "Close", func() { tr.Close() })
	within(t, time.Second, "BBS drain", func() { <-drained })
	if got := discs.Load(); got != 1 {
		t.Errorf("client sent %d DISCs after Abort, want 1", got)
	}
}

// TestAbortStalled checks that Abort and Close return promptly even while the
// link goroutine is stuck in a write and the port's Close blocks.
func TestAbortStalled(t *testing.T) {
	shortTimeouts(t, readTimeout, minWriteTimeout, 200*time.Millisecond)
	tr, b := connectTest(t, testDevice, newStallPort(t), nil)
	port := lastStallPort
	go func() {
		for range b.frames {
		}
	}()
	port.stall.Store(true)
	port.blockClose.Store(true)
	tr.Send("hello") // the link goroutine is now stuck writing this
	result := make(chan error, 1)
	go func() {
		_, err := tr.ReadUntil("never")
		result <- err
	}()
	time.Sleep(100 * time.Millisecond)
	within(t, 100*time.Millisecond, "Abort", tr.Abort)
	select {
	case err := <-result:
		if err != errAborted {
			t.Errorf("ReadUntil returned %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("ReadUntil didn't return after Abort")
	}
	within(t, time.Second, "Close", func() { tr.Close() })
}

// TestKISSDecoderOverflow checks that the KISS decoder discards an over-long
// frame, and doesn't let garbage without FENDs grow its buffer without bound.
func TestKISSDecoderOverflow(t *testing.T) {
	var (
		dec kissDecoder
		got [][]byte
	)
	emit := func(f []byte) { got = append(got, f) }
	good := []byte("a good frame")
	stream := append([]byte{fend, kissData}, bytes.Repeat([]byte{'A'}, 3*maxKISSFrame)...)
	stream = append(stream, kissEncode(good)...)
	for len(stream) != 0 {
		count := min(len(stream), 1000)
		dec.feed(stream[:count], emit)
		stream = stream[count:]
	}
	if len(got) != 1 || !bytes.Equal(got[0], good) {
		t.Errorf("got %d frames, expected just the good one", len(got))
	}
	junk := bytes.Repeat([]byte{0x55}, 1000)
	for range 1000 {
		dec.feed(junk, emit)
	}
	if cap(dec.buf) > 2*maxKISSFrame {
		t.Errorf("decoder buffer grew to %d bytes", cap(dec.buf))
	}
	got = nil
	dec.feed(kissEncode(good), emit)
	if len(got) != 1 || !bytes.Equal(got[0], good) {
		t.Errorf("after junk: got %d frames, expected just the good one", len(got))
	}
}

// TestDeviceValidated checks that device parameters are forced into range.
func TestDeviceValidated(t *testing.T) {
	for _, tt := range []struct{ in, want Device }{
		{Device{MaxFrame: 0, PacLen: 0, Frack: 0, Retry: 0, RespTime: -1}, Device{MaxFrame: 1, PacLen: 1, Frack: minFrack, InitialFrack: minFrack, Retry: 1}},
		{Device{MaxFrame: 8, PacLen: 1000, Frack: time.Second, Retry: 3}, Device{MaxFrame: 7, PacLen: 256, Frack: time.Second, InitialFrack: time.Second, Retry: 3}},
		{Device{MaxFrame: 1, PacLen: 1, Frack: time.Second, InitialFrack: time.Minute, Retry: 1}, Device{MaxFrame: 1, PacLen: 1, Frack: time.Second, InitialFrack: maxFrack, Retry: 1}},
		{*testDevice, func() Device { d := *testDevice; d.InitialFrack = d.Frack; return d }()},
	} {
		in := tt.in
		if got := in.validated(); *got != tt.want {
			t.Errorf("%+v: got %+v, want %+v", tt.in, *got, tt.want)
		}
		if in != tt.in {
			t.Errorf("validated modified its receiver")
		}
	}
}
