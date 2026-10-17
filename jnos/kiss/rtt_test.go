package kiss

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// TestT1Base checks the arithmetic of the adaptive T1 interval.
func TestT1Base(t *testing.T) {
	const frack, initial = 8 * time.Second, 15 * time.Second
	newT := func() *Transport { return &Transport{dev: &Device{Frack: frack, InitialFrack: initial}} }
	samples := func(t *Transport, secs ...float64) {
		for _, s := range secs {
			t.sampleRTT(time.Duration(s * float64(time.Second)))
		}
	}
	sec := func(s float64) time.Duration { return time.Duration(s * float64(time.Second)) }

	// Until there are rttMinSamples answers: the patient initial value,
	// however fast they were.  (One fast answer used to drop the wait from
	// 15 to 8 seconds.)
	tr := newT()
	if got := tr.t1Base(); got != initial {
		t.Errorf("unmeasured: %v, want %v", got, initial)
	}
	samples(tr, 2, 2)
	if got := tr.t1Base(); got != initial {
		t.Errorf("two answers: %v, want %v", got, initial)
	}
	// A BBS that answers quickly and steadily: never shorter than Frack.
	samples(tr, 2, 2, 2, 2)
	if got := tr.t1Base(); got != frack {
		t.Errorf("fast BBS: %v, want the %v floor", got, frack)
	}
	// W3XSC's answers (from a real session): one slow one among the last
	// five sets the interval to 1.5 times it.
	tr = newT()
	samples(tr, 2.2, 2.5, 2.4, 2.6, 6.5)
	if got := tr.t1Base(); got != sec(9.75) {
		t.Errorf("slow answer in the window: %v, want 9.75s", got)
	}
	// It stays until five faster answers have pushed it out ...
	samples(tr, 2, 2, 2, 2)
	if got := tr.t1Base(); got != sec(9.75) {
		t.Errorf("slow answer still in the window: %v, want 9.75s", got)
	}
	samples(tr, 2)
	if got := tr.t1Base(); got != frack {
		t.Errorf("slow answer out of the window: %v, want %v", got, frack)
	}
	// ... but a single slow answer raises it at once.
	samples(tr, 7)
	if got := tr.t1Base(); got != sec(10.5) {
		t.Errorf("new slow answer: %v, want 10.5s", got)
	}
	// A BBS that's slow every time: twice its median.
	tr = newT()
	samples(tr, 5, 6, 5.5, 6, 5)
	if got := tr.t1Base(); got != 11*time.Second {
		t.Errorf("steadily slow BBS: %v, want 11s", got)
	}
	// A very slow answer: capped.
	samples(tr, 60)
	if got := tr.t1Base(); got != maxFrack {
		t.Errorf("very slow BBS: %v, want the %v cap", got, maxFrack)
	}
	// Backoff applies on top of the measured base, up to the cap.
	tr = newT()
	samples(tr, 2.2, 2.5, 6.5, 2.4, 2.6)
	base := tr.t1Base()
	tr.retries = 1
	if got := tr.frack(); got != min(2*base, maxFrack) {
		t.Errorf("first retry: %v, want %v", got, min(2*base, maxFrack))
	}
	// The median of an even number of answers is the mean of the middle two.
	if median, slowest := rttStats([]time.Duration{4, 1, 3, 2}); median != 2 || slowest != 4 {
		t.Errorf("rttStats = %v, %v; want 2, 4", median, slowest)
	}
}

// TestT1AdaptsToSlowBBS checks that once the BBS has been measured answering
// slowly, we wait longer than Frack before retransmitting, and that the
// monitor reports the change.
func TestT1AdaptsToSlowBBS(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Turnaround, dev.RespTime, dev.Retry = 300*time.Millisecond, 0, 20*time.Millisecond, 8
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()
	frameAir := airtime(len("m0\r"))

	// The BBS acknowledges each frame after a variable delay, always
	// before our T1 (300ms plus airtime) would expire.
	for i, delay := range []time.Duration{250, 420, 250, 420, 250, 420} {
		tr.Send(fmt.Sprintf("m%d", i))
		f, _ := next(fmt.Sprintf("frame %d", i))
		if !f.isI() || f.pf() {
			t.Fatalf("frame %d: got control %02X, want a first transmission", i, f.control)
		}
		time.Sleep(delay * time.Millisecond)
		b.vr = byte(i+1) & 7
		b.sendS(ctlRR, false)
	}
	// Now a frame the BBS doesn't answer.  The retransmission must wait
	// longer than the unmeasured T1 would have.
	tr.Send("m6")
	next("frame 6")
	f, waited := next("retransmission")
	if !f.isI() || !f.pf() {
		t.Fatalf("retransmission: got control %02X, want the data resent with P", f.control)
	}
	if unmeasured := dev.Frack + frameAir; waited < unmeasured+150*time.Millisecond {
		t.Errorf("retransmitted after %v; want well beyond the unmeasured %v, since the BBS takes up to %v to answer",
			waited, unmeasured, 420*time.Millisecond)
	}
	if waited > maxFrack {
		t.Errorf("retransmitted after %v", waited)
	}
	// Answer it, so that Close needn't wait out the retries.
	b.vr = 7
	b.sendS(ctlRR, true)
}

// TestInitialFrack checks that until the BBS's answer time has been measured,
// we wait InitialFrack, not Frack, before resending: here, a SABM.
func TestInitialFrack(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.InitialFrack, dev.Retry = 200*time.Millisecond, 700*time.Millisecond, 3
	if v := dev.validated(); v.InitialFrack != 700*time.Millisecond {
		t.Fatalf("validated InitialFrack = %v", v.InitialFrack)
	}
	if v := (&Device{Frack: time.Second}).validated(); v.InitialFrack != time.Second {
		t.Errorf("zero InitialFrack validated to %v, want Frack", v.InitialFrack)
	}
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	gap := make(chan time.Duration, 1)
	go func() {
		var first time.Time
		for f := range b.frames {
			if !isU(f, ctlSABM) {
				continue
			}
			if first.IsZero() {
				first = time.Now() // lost
				continue
			}
			gap <- time.Since(first)
			b.send(&frame{control: ctlUA | pfBit})
			break
		}
		for range b.frames {
		}
	}()
	tr, err := open(client, &dev, testClient, testBBS, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { tr.Abort(); tr.Close(); server.Close() }()
	if g := <-gap; g < 600*time.Millisecond {
		t.Errorf("resent the SABM after %v, before InitialFrack (%v)", g, dev.InitialFrack)
	}
}

// TestNoPollWhenAllAcknowledged checks that once everything we sent has been
// acknowledged, we don't poll the BBS to finish timer recovery (as happened
// at 09:04:54: an RR poll, an extra transmission, with nothing outstanding).
func TestNoPollWhenAllAcknowledged(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Turnaround, dev.RespTime = 300*time.Millisecond, 0, 20*time.Millisecond
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()
	tr.Send("XM 0")
	next("first transmission")
	if f, _ := next("retransmission"); !f.isI() || !f.pf() {
		t.Fatalf("retransmission: got control %02X, want the data with P", f.control)
	}
	// The BBS answers with data (a command, without F) that acknowledges
	// ours, rather than with an F response to the poll.
	b.vr = 1
	b.send(&frame{command: true, control: 1 << 5, info: []byte("prompt\r")})
	deadline := time.After(1500 * time.Millisecond)
	for {
		select {
		case f := <-b.frames:
			if f.command && f.pf() {
				t.Fatalf("polled (control %02X) with nothing outstanding", f.control)
			}
		case <-deadline:
			return
		}
	}
}
