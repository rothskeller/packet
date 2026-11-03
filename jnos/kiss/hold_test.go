package kiss

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
)

// holdHarness connects a Transport, using dev, to a fake BBS, and returns
// helpers for exchanging single frames with it.
func holdHarness(t *testing.T, dev *Device) (b *fakeBBS, tr *Transport, next func(what string) (*frame, time.Duration), cleanup func()) {
	t.Helper()
	client, server := net.Pipe()
	b = newFakeBBS(t, server)
	go b.accept()
	tr, err := open(client, dev, testClient, testBBS, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	next = func(what string) (*frame, time.Duration) {
		start := time.Now()
		select {
		case f := <-b.frames:
			return f, time.Since(start)
		case <-time.After(5 * time.Second):
			t.Fatalf("%s: no frame from client", what)
			return nil, 0
		}
	}
	cleanup = func() {
		go func() {
			for f := range b.frames {
				if !f.isI() && !f.isS() && f.uType() == ctlDISC {
					b.send(&frame{control: ctlUA | pfBit})
				}
			}
		}()
		tr.Close()
		server.Close()
	}
	return b, tr, next, cleanup
}

func TestTurnaroundHold(t *testing.T) {
	const hold = 300 * time.Millisecond
	dev := *testDevice
	dev.RespTime = 50 * time.Millisecond
	dev.Turnaround = hold
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()
	// Wait out the hold that the UA started.
	time.Sleep(hold + 50*time.Millisecond)

	// A response to a poll waits for the turnaround hold, though it would
	// otherwise go right away.
	b.send(&frame{command: true, control: ctlRR | pfBit})
	f, waited := next("poll response")
	if !f.isS() || f.sType() != ctlRR || f.command || !f.pf() {
		t.Errorf("poll response: got control %02X command=%v, want RR response with F", f.control, f.command)
	}
	if waited < hold*9/10 {
		t.Errorf("poll response sent after %v, before the %v turnaround hold", waited, hold)
	}
	if waited > hold+250*time.Millisecond {
		t.Errorf("poll response took %v, much longer than the %v hold", waited, hold)
	}

	// So does data sent right after hearing the BBS.  (Here the I frame
	// also carries the acknowledgment of the BBS's frame.)
	time.Sleep(hold + 50*time.Millisecond)
	b.send(&frame{command: true, control: 0, info: []byte("prompt\r")})
	start := time.Now()
	tr.Send("XM 0")
	f, _ = next("data after BBS frame")
	waited = time.Since(start)
	if !f.isI() || f.nr() != 1 {
		t.Errorf("data: got control %02X, want an I frame acknowledging N(R)=1", f.control)
	}
	if waited < hold*9/10 {
		t.Errorf("data sent after %v, before the %v turnaround hold", waited, hold)
	}
}

func TestDuplicateFrameAcknowledged(t *testing.T) {
	dev := *testDevice
	dev.RespTime = 50 * time.Millisecond
	b, _, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	// The BBS sends its first frame, and we acknowledge it.
	b.send(&frame{command: true, control: 0, info: []byte("banner\r")})
	f, _ := next("first acknowledgment")
	if !f.isS() || f.sType() != ctlRR || f.nr() != 1 {
		t.Fatalf("first acknowledgment: got control %02X, want RR N(R)=1", f.control)
	}
	// It didn't hear that, so it resends the frame with a poll.  We
	// already have it: the answer is RR (with F), not REJ.
	b.send(&frame{command: true, control: pfBit, info: []byte("banner\r")})
	f, _ = next("duplicate response")
	if !f.isS() || f.sType() != ctlRR || f.command || !f.pf() || f.nr() != 1 {
		t.Errorf("duplicate frame: got control %02X command=%v, want RR response with F and N(R)=1", f.control, f.command)
	}
	// And the same without a poll.
	b.send(&frame{command: true, control: 0, info: []byte("banner\r")})
	f, _ = next("duplicate acknowledgment")
	if !f.isS() || f.sType() != ctlRR || f.nr() != 1 {
		t.Errorf("duplicate frame without poll: got control %02X, want RR N(R)=1", f.control)
	}

	// A frame from ahead of what we expect is still a lost frame: REJ.
	b.send(&frame{command: true, control: 2 << 1, info: []byte("later\r")})
	f, _ = next("out-of-sequence response")
	if !f.isS() || f.sType() != ctlREJ || f.nr() != 1 {
		t.Errorf("frame N(S)=2 when expecting 1: got control %02X, want REJ N(R)=1", f.control)
	}
}

// TestReadWaitsForLinkRetries checks that while the link is still retrying our
// data, ReadUntil doesn't time out on its own: the link's retry limit decides
// when the BBS is gone.  (On a slow BBS, the read timeout used to end the
// session while the link was still recovering.)
func TestReadWaitsForLinkRetries(t *testing.T) {
	const read = 300 * time.Millisecond
	shortTimeouts(t, read, minWriteTimeout, portCloseTimeout)
	dev := *testDevice
	dev.Frack, dev.Retry = 200*time.Millisecond, 4
	tr, b := connectTest(t, &dev, nil, nil)
	stop := make(chan struct{})
	defer close(stop)
	go func() {
		// The BBS hears our polls and answers them, but never gets
		// our data.
		for {
			select {
			case f, ok := <-b.frames:
				if !ok {
					return
				}
				if f.command && f.pf() {
					b.sendS(ctlRR, true)
				}
			case <-stop:
				return
			}
		}
	}()
	if err := tr.Send("SP PKTMON"); err != nil {
		t.Fatal(err)
	}
	start := time.Now()
	var err error
	within(t, 10*time.Second, "ReadUntil", func() { _, err = tr.ReadUntil("Subject:\n") })
	elapsed := time.Since(start)
	if errors.Is(err, jnos.ErrTimeout) {
		t.Errorf("ReadUntil timed out after %v while the link was still retrying", elapsed)
	} else if !errors.Is(err, jnos.ErrDisconnected) {
		t.Errorf("got %v, expected the link to give up (disconnected)", err)
	}
	if elapsed < 2*read {
		t.Errorf("ReadUntil returned after %v, before the link had used its retries", elapsed)
	}
}

// TestTimeoutResendsDataWithPoll checks that when T1 expires with data
// unacknowledged, the poll is that data resent with the P bit (as QtTermTCP
// does), not a bare RR, and that the BBS's answer to it makes progress.
func TestTimeoutResendsDataWithPoll(t *testing.T) {
	dev := *testDevice
	dev.Frack = 300 * time.Millisecond
	dev.RespTime = 50 * time.Millisecond
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	tr.Send("XM 0")
	f, _ := next("first transmission")
	if !f.isI() || f.pf() || f.ns() != 0 || string(f.info) != "XM 0\r" {
		t.Fatalf("first transmission: got control %02X info %q, want I N(S)=0 without P carrying XM 0", f.control, f.info)
	}
	// That one is lost.  When T1 expires, the poll must be the data again.
	f, waited := next("timeout poll")
	if !f.isI() || !f.command || !f.pf() || f.ns() != 0 || string(f.info) != "XM 0\r" {
		t.Fatalf("timeout poll: got control %02X command=%v info %q, want the I frame N(S)=0 resent with P", f.control, f.command, f.info)
	}
	if waited < dev.Frack*9/10 {
		t.Errorf("timeout poll sent after %v, before T1 (%v) expired", waited, dev.Frack)
	}
	// The BBS hears it and acknowledges it in its answer to the poll.
	b.vr = 1
	b.sendS(ctlRR, true)
	// The link is healthy again: new data goes out normally, numbered
	// after the acknowledged frame, without a poll.
	tr.Send("LA")
	f, _ = next("next data")
	if !f.isI() || f.pf() || f.ns() != 1 || string(f.info) != "LA\r" {
		t.Errorf("next data: got control %02X info %q, want I N(S)=1 without P carrying LA", f.control, f.info)
	}
}

// TestT1Backoff checks that each retry without progress waits longer than the
// last, up to the cap, so that a slow BBS isn't swamped with retransmissions.
func TestT1Backoff(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Retry, dev.Turnaround = 200*time.Millisecond, 4, 0
	_, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	tr.Send("XM 0")
	next("first transmission")
	var gaps []time.Duration
	for i := range 3 {
		f, waited := next("retry")
		if !f.isI() || !f.command || !f.pf() {
			t.Fatalf("retry %d: got control %02X, want the data resent with P", i+1, f.control)
		}
		gaps = append(gaps, waited)
	}
	// The first retry comes after Frack (plus airtime), the second after
	// twice that, the third after four times (the cap here).
	if gaps[1] < gaps[0]*3/2 {
		t.Errorf("second retry gap %v isn't backed off from the first (%v)", gaps[1], gaps[0])
	}
	if gaps[2] < gaps[1]*3/2 {
		t.Errorf("third retry gap %v isn't backed off from the second (%v)", gaps[2], gaps[1])
	}
	if gaps[2] > 4*dev.Frack+400*time.Millisecond {
		t.Errorf("third retry gap %v exceeds the 4×Frack cap", gaps[2])
	}
}

// TestT1PushedByBBSFrames checks that while the BBS keeps transmitting (here,
// resending a frame we already have, as W3XSC does when it misses our
// acknowledgment), we don't retransmit over it: T1 counts quiet time.
func TestT1PushedByBBSFrames(t *testing.T) {
	dev := *testDevice
	dev.Frack, dev.Turnaround, dev.RespTime = 300*time.Millisecond, 0, 20*time.Millisecond
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	// The BBS sends a frame; we acknowledge it.
	b.send(&frame{command: true, control: 0, info: []byte("prompt\r")})
	next("acknowledgment")
	// We send data, which the BBS doesn't get.
	tr.Send("SP PKTMON")
	if f, _ := next("data"); !f.isI() {
		t.Fatalf("data: got control %02X, want an I frame", f.control)
	}
	// For 700ms, the BBS keeps resending its prompt every 150ms.  We
	// answer those, but must not retransmit, though T1 (300ms plus the
	// data's airtime, about 500ms) would otherwise have expired.
	stop := time.After(700 * time.Millisecond)
	tick := time.NewTicker(150 * time.Millisecond)
	defer tick.Stop()
	polls := 0
BUSY:
	for {
		select {
		case <-tick.C:
			b.send(&frame{command: true, control: 0, info: []byte("prompt\r")})
		case f := <-b.frames:
			if f.command && f.pf() {
				polls++
			}
		case <-stop:
			break BUSY
		}
	}
	if polls != 0 {
		t.Errorf("retransmitted %d times while the BBS was transmitting", polls)
	}
	// Once the BBS goes quiet, we retransmit after T1.
	for {
		f, waited := next("retransmission after quiet")
		if f.command && f.pf() {
			if !f.isI() || string(f.info) != "SP PKTMON\r" {
				t.Errorf("retransmission: got control %02X info %q, want the data with P", f.control, f.info)
			}
			if waited > 3*dev.Frack {
				t.Errorf("retransmission came %v after the BBS went quiet", waited)
			}
			break
		}
	}
}

// TestResendWhenPolledWithoutAck checks that when the BBS polls us (here, by
// resending its last frame with P, as W3XSC does) without acknowledging data
// we sent long enough ago, we answer the poll and resend the data right away,
// rather than waiting for T1.
func TestResendWhenPolledWithoutAck(t *testing.T) {
	old := polledWithoutGap
	polledWithoutGap = 200 * time.Millisecond
	t.Cleanup(func() { polledWithoutGap = old })
	dev := *testDevice
	dev.Frack, dev.InitialFrack, dev.Turnaround, dev.RespTime = 5*time.Second, 5*time.Second, 0, 20*time.Millisecond
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	// The BBS sends a prompt, which our data acknowledges.
	b.send(&frame{command: true, control: 0, info: []byte("prompt\r")})
	tr.ReadUntil("prompt\n")
	tr.Send("XM 0")
	if f, _ := next("data"); !f.isI() || f.nr() != 1 {
		t.Fatalf("data: got control %02X, want I acknowledging N(R)=1", f.control)
	}
	// The BBS heard neither: it resends its prompt, with P, N(R)=0.
	time.Sleep(400 * time.Millisecond)
	b.send(&frame{command: true, control: pfBit, info: []byte("prompt\r")})
	f, _ := next("poll answer")
	if !f.isS() || f.command || !f.pf() || f.nr() != 1 {
		t.Fatalf("poll answer: got control %02X command=%v, want RR response with F, N(R)=1", f.control, f.command)
	}
	f, waited := next("resend")
	if !f.isI() || f.ns() != 0 || string(f.info) != "XM 0\r" {
		t.Fatalf("resend: got control %02X info %q, want the data again", f.control, f.info)
	}
	if waited > time.Second {
		t.Errorf("resent the data %v after the poll; want right away, not after T1", waited)
	}
	b.vr = 1
	b.sendS(ctlRR, false)
}

// TestNoResendWhenPolledTooSoon checks that a poll from the BBS that comes
// before it could have answered our data doesn't make us resend it.
func TestNoResendWhenPolledTooSoon(t *testing.T) {
	old := polledWithoutGap
	polledWithoutGap = 2 * time.Second
	t.Cleanup(func() { polledWithoutGap = old })
	dev := *testDevice
	dev.Frack, dev.InitialFrack, dev.Turnaround, dev.RespTime = 5*time.Second, 5*time.Second, 0, 20*time.Millisecond
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	tr.Send("XM 0")
	next("data")
	b.send(&frame{command: true, control: ctlRR | pfBit})
	if f, _ := next("poll answer"); !f.isS() || !f.pf() {
		t.Fatalf("poll answer: got control %02X", f.control)
	}
	select {
	case f := <-b.frames:
		if f.isI() {
			t.Errorf("resent the data (control %02X) though the BBS couldn't have answered it yet", f.control)
		}
	case <-time.After(700 * time.Millisecond):
	}
	b.vr = 0
	b.send(&frame{control: ctlRR | 1<<5})
}

// TestResendWhenPolledDuringRecovery checks the case of the 09:43 session:
// our data and our T1 poll were both lost, and then the BBS polled us (by
// resending its banner with P).  We must answer and resend our data (as a
// poll), not go silent: T1 must keep running.
func TestResendWhenPolledDuringRecovery(t *testing.T) {
	old := polledWithoutGap
	polledWithoutGap = 100 * time.Millisecond
	t.Cleanup(func() { polledWithoutGap = old })
	dev := *testDevice
	// T1 is long compared to the other delays, so that the BBS's poll
	// comes well before T1 would expire again, even on a busy machine.
	dev.Frack, dev.InitialFrack, dev.Turnaround, dev.RespTime, dev.Retry = time.Second, time.Second, 200*time.Millisecond, 20*time.Millisecond, 10
	b, tr, next, cleanup := holdHarness(t, &dev)
	defer cleanup()

	b.send(&frame{command: true, control: 0, info: []byte("banner\r")})
	tr.ReadUntil("banner\n")
	tr.Send("XM 0")
	next("data")
	if f, _ := next("T1 poll"); !f.isI() || !f.pf() {
		t.Fatalf("T1 poll: got control %02X, want the data resent with P", f.control)
	}
	// Now in timer recovery.  The BBS, having heard nothing, polls us a
	// while later.
	time.Sleep(400 * time.Millisecond)
	b.send(&frame{command: true, control: pfBit, info: []byte("banner\r")})
	f, waited := next("poll answer")
	if !f.isS() || f.command || !f.pf() || f.nr() != 1 {
		t.Fatalf("poll answer: got control %02X, want RR response with F, N(R)=1", f.control)
	}
	if waited < dev.Turnaround*9/10 {
		t.Errorf("answered after %v, before the turnaround hold", waited)
	}
	f, waited = next("resend")
	if !f.isI() || !f.pf() || f.ns() != 0 {
		t.Fatalf("resend: got control %02X, want the data resent with P", f.control)
	}
	if waited > 500*time.Millisecond {
		t.Errorf("resent %v after the poll answer; want right after it", waited)
	}
	// And if that's lost too, T1 is still running: another resend comes.
	if f, _ = next("next T1 poll"); !f.isI() || !f.pf() {
		t.Fatalf("after the resend: got control %02X, want another poll from T1", f.control)
	}
	b.vr = 1
	b.sendS(ctlRR, true)
}
