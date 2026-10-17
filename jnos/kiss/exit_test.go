package kiss

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestFinalDISCAfterLinkFailure checks that when we give up on a BBS that
// stopped answering, we send it a final DISC, so that if it can still hear us,
// it ends its side of the session.
func TestFinalDISCAfterLinkFailure(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Retry, dev.Turnaround = 100*time.Millisecond, 2, 0
	tr, b := connectTest(t, &dev, nil, nil)
	var discs atomic.Int32
	drained := make(chan struct{})
	go func() {
		for f := range b.frames { // the BBS hears everything, answers nothing
			if isU(f, ctlDISC) {
				discs.Add(1)
			}
		}
		close(drained)
	}()
	tr.Send("SP PKTMON")
	within(t, 5*time.Second, "ReadUntil", func() { tr.ReadUntil("never") })
	within(t, time.Second, "Close", func() { tr.Close() })
	if got := discs.Load(); got != 1 {
		t.Errorf("sent %d DISCs after the link failed, want 1", got)
	}
}

// TestNoFinalDISCAfterBBSDisconnect checks that when the BBS ends the session,
// we don't send it a DISC afterward.
func TestNoFinalDISCAfterBBSDisconnect(t *testing.T) {
	tr, b := connectTest(t, testDevice, nil, nil)
	var discs atomic.Int32
	go func() {
		b.write("bye\n")
		b.send(&frame{command: true, control: ctlDISC | pfBit})
		for f := range b.frames {
			if isU(f, ctlDISC) {
				discs.Add(1)
			}
		}
	}()
	within(t, 5*time.Second, "ReadUntil", func() { tr.ReadUntil("never") })
	within(t, time.Second, "Close", func() { tr.Close() })
	time.Sleep(50 * time.Millisecond)
	if got := discs.Load(); got != 0 {
		t.Errorf("sent %d DISCs after the BBS disconnected", got)
	}
}

// TestFinalDISCAfterCanceledConnect checks that closing while still connecting
// sends a DISC, in case the BBS heard our SABM and we missed its UA.
func TestFinalDISCAfterCanceledConnect(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Retry = 100*time.Millisecond, 1
	client, server := net.Pipe()
	defer server.Close()
	b := newFakeBBS(t, server)
	var discs atomic.Int32
	go func() {
		for f := range b.frames { // the BBS never answers
			if isU(f, ctlDISC) {
				discs.Add(1)
			}
		}
	}()
	if _, err := open(client, &dev, testClient, testBBS, nil, nil); err == nil {
		t.Fatal("connect succeeded")
	}
	time.Sleep(50 * time.Millisecond)
	if got := discs.Load(); got != 1 {
		t.Errorf("sent %d DISCs after giving up on the connection, want 1", got)
	}
}

// TestSessionMarker checks the session marker's name and life cycle: it's
// cleared after a clean session, and kept when the TNC couldn't be reset, so
// that the next session resets it first.
func TestSessionMarker(t *testing.T) {
	if m := sessionMarker("/dev/cu.usbmodem 1234"); m != "" && !strings.HasSuffix(m, filepath.Join("packet", "kiss-session-cu.usbmodem_1234")) {
		t.Errorf("sessionMarker = %q", m)
	}
	fastReset(t)
	run := func(radio *fakeRadio, clean bool) (kept bool) {
		marker := filepath.Join(t.TempDir(), "marker")
		markSession(marker)
		var reopened int
		client, server := net.Pipe()
		defer server.Close()
		b := newFakeBBS(t, server)
		go func() {
			b.accept()
			if clean {
				b.disconnect()
			} else {
				// The BBS drops us once we send something.
				for f := range b.frames {
					if f.isI() {
						b.send(&frame{control: ctlDM | pfBit})
						break
					}
				}
			}
			for range b.frames {
			}
		}()
		withMarker := func(t *Transport) { t.marker = marker }
		tr, err := open(client, testDevice, testClient, testBBS, nil, nil, withRadio(radio, &reopened), withMarker)
		if err != nil {
			t.Fatal(err)
		}
		if !clean {
			tr.Send("hello")
		}
		tr.ReadUntil("never")
		tr.Close()
		_, err = os.Stat(marker)
		return err == nil
	}
	if run(newFakeRadio("0"), true) {
		t.Error("marker kept after a clean session")
	}
	if run(newFakeRadio("0"), false) {
		t.Error("marker kept after an unclean session whose reset worked")
	}
	silent := newFakeRadio("0")
	silent.silent = true
	if !run(silent, false) {
		t.Error("marker cleared though the TNC reset failed")
	}
}

func TestResetBeforeStart(t *testing.T) {
	fastReset(t)
	var lines []string
	monitor := func(line string) { lines = append(lines, line) }
	radio := newFakeRadio("1")
	if err := resetBeforeStart(radio, monitor); err != nil {
		t.Fatal(err)
	}
	if !radio.kiss || !strings.Contains(strings.Join(lines, "\n"), "re-entered it with TN 2,1") {
		t.Errorf("reset before start: radio in KISS %v, monitor:\n%s", radio.kiss, strings.Join(lines, "\n"))
	}
	silent := newFakeRadio("0")
	silent.silent = true
	if err := resetBeforeStart(silent, monitor); err == nil || !strings.Contains(err.Error(), "back to KISS mode") {
		t.Errorf("failed reset before start: got %v", err)
	}
}

// TestAbortWhileConnecting checks that canceling the context while the
// connection is still being attempted ends it promptly, with a final DISC and
// (for a Kenwood TNC) a reset, rather than waiting out the SABM retries.
func TestAbortWhileConnecting(t *testing.T) {
	fastReset(t)
	dev := *testDevice
	dev.Frack, dev.Retry = 2*time.Second, 5 // the retries would take 12s
	radio := newFakeRadio("0")
	var reopened int
	client, server := net.Pipe()
	defer server.Close()
	b := newFakeBBS(t, server)
	var discs atomic.Int32
	go func() {
		for f := range b.frames { // the BBS never answers
			if isU(f, ctlDISC) {
				discs.Add(1)
			}
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(200*time.Millisecond, cancel)
	start := time.Now()
	_, err := open(client, &dev, testClient, testBBS, nil, nil, withRadio(radio, &reopened), withContext(ctx))
	if err == nil {
		t.Fatal("connect succeeded")
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("connect took %v to notice the abort", elapsed)
	}
	if reopened != 1 {
		t.Errorf("TNC reset %d times after the aborted connect, want 1", reopened)
	}
	time.Sleep(50 * time.Millisecond)
	if got := discs.Load(); got != 1 {
		t.Errorf("sent %d DISCs after the aborted connect, want 1", got)
	}
}

// TestConnectedWithoutUA checks that when the BBS's UA is lost but its first
// I frame arrives, we take the link as up right away (without another SABM,
// which would reset it), and deliver and acknowledge that frame.
func TestConnectedWithoutUA(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.InitialFrack, dev.Turnaround, dev.RespTime = 2*time.Second, 2*time.Second, 0, 20*time.Millisecond
	client, server := net.Pipe()
	defer server.Close()
	b := newFakeBBS(t, server)
	sabms := make(chan int, 1)
	acked := make(chan *frame, 1)
	go func() {
		var count int
		for f := range b.frames {
			switch {
			case isU(f, ctlSABM):
				count++
				if count == 1 {
					// The UA is lost; the banner isn't.
					b.send(&frame{command: true, control: 0, info: []byte("Welcome\r(#1) >\r")})
				}
			case f.isS() && f.nr() == 1:
				sabms <- count
				acked <- f
			}
		}
	}()
	start := time.Now()
	tr, err := open(client, &dev, testClient, testBBS, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { tr.Abort(); tr.Close() }()
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("connect took %v, waiting for T1 though the BBS was sending", elapsed)
	}
	var data string
	within(t, 2*time.Second, "ReadUntil", func() { data, err = tr.ReadUntil("(#1) >\n") })
	if err != nil || data != "Welcome\n(#1) >\n" {
		t.Errorf("ReadUntil = %q, %v", data, err)
	}
	select {
	case n := <-sabms:
		if n != 1 {
			t.Errorf("sent %d SABMs, want 1", n)
		}
	case <-time.After(2 * time.Second):
		t.Error("the BBS's frame was never acknowledged")
	}
}

// TestStaleFrameWhileConnecting checks that a frame from the middle of some
// other link (nonzero sequence numbers) doesn't count as the BBS accepting our
// connection.
func TestStaleFrameWhileConnecting(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.InitialFrack, dev.Retry = 200*time.Millisecond, 200*time.Millisecond, 1
	client, server := net.Pipe()
	defer server.Close()
	b := newFakeBBS(t, server)
	go func() {
		for f := range b.frames {
			if isU(f, ctlSABM) {
				b.send(&frame{command: true, control: 3<<5 | 2<<1, info: []byte("stale\r")})
			}
		}
	}()
	if _, err := open(client, &dev, testClient, testBBS, nil, nil); err == nil {
		t.Error("connected on a stale frame")
	}
}
