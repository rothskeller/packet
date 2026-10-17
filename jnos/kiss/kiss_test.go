package kiss

import (
	"bytes"
	"cmp"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
)

var testDevice = &Device{
	Name:     "test",
	PacLen:   32,
	MaxFrame: 2,
	Frack:    300 * time.Millisecond,
	Retry:    5,
	RespTime: 50 * time.Millisecond,
	TXDelay:  400 * time.Millisecond,
	Persist:  63,
	SlotTime: 100 * time.Millisecond,
}

var (
	testBBS    = address{call: "W6XSC", ssid: 1}
	testClient = address{call: "KC6RSC"}
)

func TestAddress(t *testing.T) {
	for _, s := range []string{"W6XSC-1", "kc6rsc", "W6XSC-15"} {
		if a, err := parseAddress(s); err != nil {
			t.Errorf("%s: %v", s, err)
		} else if back, _, _ := decodeAddress(appendAddress(nil, a, false, true)); back != a {
			t.Errorf("%s: round trip gave %v", s, back)
		}
	}
	for _, s := range []string{"", "W6XSC-16", "TOOLONGCALL", "W6-XSC", "W6XSC-"} {
		if _, err := parseAddress(s); err == nil {
			t.Errorf("%q: expected error", s)
		}
	}
}

func TestKISSFraming(t *testing.T) {
	data := []byte{1, fend, 2, fesc, 3, tfend, tfesc}
	var got [][]byte
	var dec kissDecoder
	stream := append([]byte("junk"), kissEncode(data)...)
	stream = append(stream, kissEncode([]byte{9, 9})...)
	for _, c := range stream { // one byte at a time, to test buffering
		dec.feed([]byte{c}, func(f []byte) { got = append(got, f) })
	}
	if len(got) != 2 || !bytes.Equal(got[0], data) || !bytes.Equal(got[1], []byte{9, 9}) {
		t.Errorf("got %v", got)
	}
}

func TestFrameEncoding(t *testing.T) {
	f := &frame{dest: testBBS, src: testClient, command: true, control: 3<<5 | pfBit | 5<<1, info: []byte("hello")}
	g, err := decodeFrame(f.encode())
	if err != nil {
		t.Fatal(err)
	}
	if g.dest != testBBS || g.src != testClient || !g.command || !g.isI() || g.ns() != 5 || g.nr() != 3 || !g.pf() || string(g.info) != "hello" {
		t.Errorf("got %+v", g)
	}
	f = &frame{dest: testBBS, src: testClient, command: false, control: ctlRR | 6<<5}
	if g, err = decodeFrame(f.encode()); err != nil {
		t.Fatal(err)
	}
	if g.command || !g.isS() || g.sType() != ctlRR || g.nr() != 6 || g.pf() {
		t.Errorf("got %+v", g)
	}
}

// fakeBBS simulates a JNOS BBS on the far side of a KISS TNC.  It implements
// just enough of AX.25 to exercise the Transport, using a simple go-back-N
// sender with a window of 4 frames.
type fakeBBS struct {
	t      *testing.T
	conn   net.Conn
	frames chan *frame
	vs, va byte // send state
	vr     byte // receive state
	unack  [][]byte
	line   []byte
	// dropFromClient lists the (1-based) I frames from the client that
	// are "lost on the air" the first time they're sent.
	dropFromClient map[int]bool
	clientI        int
	// dropToClient lists the (1-based) I frames to the client that are
	// lost the first time they're sent.
	dropToClient map[int]bool
	toClientI    int
	mu           sync.Mutex
	received     []string // complete lines received from the client
	// Counts of recovery actions taken by the client.
	clientREJs, clientPolls int
	// resendAfter is how long to wait for an acknowledgment before
	// resending (default 1s).
	resendAfter time.Duration
	// events lists, in order, the KISS commands and SABMs the client sent.
	events []string
}

func newFakeBBS(t *testing.T, conn net.Conn) *fakeBBS {
	b := &fakeBBS{t: t, conn: conn, frames: make(chan *frame, 64)}
	go func() {
		var dec kissDecoder
		dec.command = func(cmd byte, data []byte) {
			b.mu.Lock()
			b.events = append(b.events, fmt.Sprintf("cmd %02X % X", cmd, data))
			b.mu.Unlock()
		}
		buf := make([]byte, 512)
		for {
			count, err := conn.Read(buf)
			dec.feed(buf[:count], func(raw []byte) {
				if f, err := decodeFrame(raw); err == nil && f.dest == testBBS && f.src == testClient {
					if !f.isI() && !f.isS() && f.uType() == ctlSABM {
						b.mu.Lock()
						b.events = append(b.events, "SABM")
						b.mu.Unlock()
					}
					b.frames <- f
				}
			})
			if err != nil {
				close(b.frames)
				return
			}
		}
	}()
	return b
}

func (b *fakeBBS) send(f *frame) {
	f.dest, f.src = testClient, testBBS
	b.conn.Write(kissEncode(f.encode()))
}

func (b *fakeBBS) sendS(stype byte, pf bool) {
	control := stype | b.vr<<5
	if pf {
		control |= pfBit
	}
	b.send(&frame{control: control})
}

// sendI sends the frame at index i of the unacknowledged list, unless the
// test calls for it to be lost.
func (b *fakeBBS) sendI(i int, pf bool) {
	ns := (b.va + byte(i)) & 7
	control := b.vr<<5 | ns<<1
	if pf {
		control |= pfBit
	}
	b.toClientI++
	if b.dropToClient[b.toClientI] {
		return
	}
	b.send(&frame{command: true, control: control, info: b.unack[i]})
}

// write sends text to the client and waits until the client acknowledges all
// of it, handling anything else the client sends meanwhile.
func (b *fakeBBS) write(s string) {
	data := []byte(strings.ReplaceAll(s, "\n", "\r"))
	for len(data) != 0 {
		count := min(len(data), 20)
		b.unack = append(b.unack, data[:count])
		data = data[count:]
	}
	for i := range b.unack {
		if i < 4 {
			b.sendI(i, false)
		}
	}
	b.vs = (b.va + byte(min(len(b.unack), 4))) & 7
	for len(b.unack) != 0 {
		select {
		case f, ok := <-b.frames:
			if !ok {
				return
			}
			b.handle(f)
		case <-time.After(cmp.Or(b.resendAfter, time.Second)):
			// Lost acknowledgment or data: resend with poll.
			for i := range min(len(b.unack), 4) {
				b.sendI(i, i == min(len(b.unack), 4)-1)
			}
		}
	}
}

// handle processes a frame from the client.
func (b *fakeBBS) handle(f *frame) {
	switch {
	case f.isI():
		b.clientI++
		if f.pf() {
			b.clientPolls++ // timer recovery resends data with P
		}
		if b.dropFromClient[b.clientI] {
			delete(b.dropFromClient, b.clientI)
			return
		}
		b.ack(f.nr())
		if f.ns() != b.vr {
			b.sendS(ctlREJ, f.pf())
			return
		}
		b.vr = (b.vr + 1) & 7
		b.sendS(ctlRR, f.pf())
		for _, c := range f.info {
			if c == '\r' {
				b.mu.Lock()
				b.received = append(b.received, string(b.line))
				b.mu.Unlock()
				b.line = b.line[:0]
			} else {
				b.line = append(b.line, c)
			}
		}
	case f.isS():
		b.ack(f.nr())
		if f.sType() == ctlREJ {
			b.clientREJs++
		}
		if f.command && f.pf() {
			b.clientPolls++
			b.sendS(ctlRR, true)
		}
		if f.sType() == ctlREJ || (!f.command && f.pf()) {
			for i := range min(len(b.unack), 4) {
				b.sendI(i, false)
			}
		}
	case f.uType() == ctlSABM:
		b.send(&frame{control: ctlUA | pfBit})
	case f.uType() == ctlDISC:
		b.send(&frame{control: ctlUA | pfBit})
	}
}

func (b *fakeBBS) ack(nr byte) {
	count := int((nr - b.va) & 7)
	if count > len(b.unack) {
		b.t.Errorf("client acknowledged unsent frame: N(R)=%d V(A)=%d", nr, b.va)
		return
	}
	b.unack = b.unack[count:]
	b.va = nr
	// Keep the window full.
	for i := int((b.vs - b.va) & 7); i < min(len(b.unack), 4); i++ {
		b.sendI(i, false)
		b.vs = (b.vs + 1) & 7
	}
}

// waitFor handles client frames until the client has sent the given line.
func (b *fakeBBS) waitFor(line string) {
	deadline := time.After(10 * time.Second)
	for {
		b.mu.Lock()
		for i, l := range b.received {
			if l == line {
				b.received = b.received[i+1:]
				b.mu.Unlock()
				return
			}
		}
		b.mu.Unlock()
		select {
		case f, ok := <-b.frames:
			if !ok {
				b.t.Fatalf("connection closed waiting for %q", line)
			}
			b.handle(f)
		case <-deadline:
			b.t.Fatalf("timeout waiting for %q", line)
		}
	}
}

// accept waits for the client's SABM and accepts the connection.
func (b *fakeBBS) accept() {
	for f := range b.frames {
		if !f.isI() && !f.isS() && f.uType() == ctlSABM {
			b.send(&frame{control: ctlUA | pfBit})
			return
		}
	}
}

// disconnect disconnects from the client.
func (b *fakeBBS) disconnect() {
	b.send(&frame{command: true, control: ctlDISC | pfBit})
	for f := range b.frames {
		if !f.isI() && !f.isS() && f.uType() == ctlUA {
			return
		}
	}
}

// runSession runs a JNOS-like session against a client using jnos.Conn: login
// banner, paging off, sending one private message, and logging off.
func runSession(t *testing.T, b *fakeBBS, body string) {
	var lines []string
	b.accept()
	b.write("[JNOS-2.0p-B1FHIM$]\nYou have 0 messages.\n(#1) >\n")
	b.waitFor("XM 0")
	b.write("(#1) >\n")
	b.waitFor("SP pkt@w6xsc.ampr.org")
	b.write("Subject:\n")
	b.waitFor("Test subject")
	b.write("Enter message.  End with /EX or ^Z in first column (^A aborts):\n")
	for _, line := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
		b.waitFor(line)
		lines = append(lines, line)
	}
	b.waitFor("/EX")
	if strings.Join(lines, "\n")+"\n" != body {
		t.Errorf("BBS received wrong message body")
	}
	b.write("Msg queued\n(#1) >\n")
	b.waitFor("B")
	b.disconnect()
}

func testSession(t *testing.T, b *fakeBBS, client net.Conn) {
	body := strings.Repeat("The quick brown fox jumps over the lazy dog.\n", 10)
	finished := make(chan struct{})
	go func() {
		runSession(t, b, body)
		close(finished)
	}()
	defer func() { <-finished }()
	tr, err := open(client, testDevice, testClient, testBBS, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	c, err := jnos.Connect(tr, "")
	if err != nil {
		t.Fatal(err)
	}
	if err = c.Send("Test subject", body, "pkt@w6xsc.ampr.org"); err != nil {
		t.Fatal(err)
	}
	if err = c.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestSession(t *testing.T) {
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	testSession(t, b, client)
	// The TNC's channel access settings go first: TXDELAY 40 (400ms),
	// persistence 63, slot time 10 (100ms).
	b.mu.Lock()
	defer b.mu.Unlock()
	want := []string{"cmd 01 28", "cmd 02 3F", "cmd 03 0A", "SABM"}
	if len(b.events) < len(want) || strings.Join(b.events[:len(want)], ", ") != strings.Join(want, ", ") {
		t.Errorf("client started with %q, want %q", b.events, want)
	}
}

func TestAckTriggers(t *testing.T) {
	dev := *testDevice
	dev.RespTime = 200 * time.Millisecond
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go b.accept()
	tr, err := open(client, &dev, testClient, testBBS, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	next := func(what string) (f *frame, waited time.Duration) {
		start := time.Now()
		select {
		case f = <-b.frames:
			return f, time.Since(start)
		case <-time.After(3 * time.Second):
			t.Fatalf("%s: no frame from client", what)
			return nil, 0
		}
	}
	expectRR := func(what string, f *frame, nr byte, final bool) {
		if !f.isS() || f.sType() != ctlRR || f.command || f.nr() != nr || f.pf() != final {
			t.Errorf("%s: got control %02X command=%v, want RR response N(R)=%d F=%v", what, f.control, f.command, nr, final)
		}
	}
	iframe := func(ns byte, poll bool) *frame {
		f := &frame{command: true, control: ns << 1, info: []byte("text\r")}
		if poll {
			f.control |= pfBit
		}
		return f
	}

	// A plain I frame is acknowledged once the channel has been quiet.
	b.send(iframe(0, false))
	f, waited := next("plain ack")
	expectRR("plain ack", f, 1, false)
	if waited < 150*time.Millisecond {
		t.Errorf("plain ack: sent after %v, before the channel had been quiet for %v", waited, dev.RespTime)
	}

	// A burst gets one acknowledgment, after the last frame.
	b.send(iframe(1, false))
	time.Sleep(100 * time.Millisecond)
	b.send(iframe(2, false))
	f, waited = next("burst ack")
	expectRR("burst ack", f, 3, false)
	if waited < 150*time.Millisecond {
		t.Errorf("burst ack: sent %v after the last frame, before the channel had been quiet for %v", waited, dev.RespTime)
	}

	// A poll is answered right away.
	b.send(&frame{command: true, control: ctlRR | pfBit})
	f, waited = next("poll response")
	expectRR("poll response", f, 3, true)
	if waited > 100*time.Millisecond {
		t.Errorf("poll response: took %v", waited)
	}

	// So is an I frame with the poll bit.
	b.send(iframe(3, true))
	f, waited = next("polled I frame")
	expectRR("polled I frame", f, 4, true)
	if waited > 100*time.Millisecond {
		t.Errorf("polled I frame: took %v", waited)
	}

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

func TestSessionWithLostFrames(t *testing.T) {
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	// Lose some of the client's frames, including the last frame of the
	// message body (which can only be recovered by polling), and some of
	// the BBS's frames.
	b.dropFromClient = map[int]bool{3: true, 6: true, 17: true}
	b.dropToClient = map[int]bool{2: true, 5: true}
	testSession(t, b, client)
	if b.clientPolls == 0 || b.clientREJs == 0 {
		t.Errorf("client sent %d polls and %d REJs; expected some of each", b.clientPolls, b.clientREJs)
	}
	t.Logf("client sent %d polls and %d REJs", b.clientPolls, b.clientREJs)
}

func TestConnectRejected(t *testing.T) {
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go func() {
		for f := range b.frames {
			if f.uType() == ctlSABM {
				b.send(&frame{control: ctlDM | pfBit})
			}
		}
	}()
	if _, err := open(client, testDevice, testClient, testBBS, nil, nil); err == nil || !strings.Contains(err.Error(), "rejected") {
		t.Errorf("got %v, expected rejection", err)
	}
	server.Close()
}

func TestNoResponse(t *testing.T) {
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	var sabms int
	counted := make(chan struct{})
	go func() {
		for f := range b.frames {
			if f.uType() == ctlSABM {
				sabms++
			}
		}
		close(counted)
	}()
	if _, err := open(client, testDevice, testClient, testBBS, nil, nil); err == nil || !strings.Contains(err.Error(), "no response") {
		t.Errorf("got %v, expected no response", err)
	}
	server.Close()
	<-counted
	if sabms != testDevice.Retry+1 {
		t.Errorf("sent %d SABMs, expected %d", sabms, testDevice.Retry+1)
	}
}

func TestIgnoresOtherStations(t *testing.T) {
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go func() {
		b.accept()
		// Traffic between other stations, and to us from someone
		// else, must be ignored.
		b.conn.Write(kissEncode((&frame{dest: address{call: "N0CALL"}, src: testBBS, command: true, control: 0, info: []byte("nope\r")}).encode()))
		b.conn.Write(kissEncode((&frame{dest: testClient, src: address{call: "N0CALL"}, command: true, control: 0, info: []byte("nope\r")}).encode()))
		b.write("hello\n")
		b.waitFor("B")
		b.disconnect()
	}()
	tr, err := open(client, testDevice, testClient, testBBS, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if s, err := tr.ReadUntil("\n"); err != nil || s != "hello\n" {
		t.Errorf("got %q, %v", s, err)
	}
	tr.Send("B")
	if _, err := tr.ReadUntil("\n"); err != jnos.ErrDisconnected {
		t.Errorf("got %v, expected disconnect", err)
	}
	tr.Close()
}

func TestFrameString(t *testing.T) {
	ko6nvy, w3xsc := address{call: "KO6NVY"}, address{call: "W3XSC", ssid: 1}
	// The expected strings are from a QtTermTCP monitor capture.
	for _, tt := range []struct {
		command bool
		control byte
		want    string
	}{
		{true, ctlSABM | pfBit, "KO6NVY>W3XSC-1 <SABM C P>"},
		{false, ctlUA | pfBit, "KO6NVY>W3XSC-1 <UA R F>"},
		{true, 1<<5 | pfBit, "KO6NVY>W3XSC-1 <I C P R1 S0>"},
		{true, 2<<5 | pfBit | 1<<1, "KO6NVY>W3XSC-1 <I C P R2 S1>"},
		{false, ctlRR | 1<<5, "KO6NVY>W3XSC-1 <RR R R1>"},
		{false, ctlRR | pfBit | 2<<5, "KO6NVY>W3XSC-1 <RR R F R2>"},
		{true, ctlRR | pfBit | 2<<5, "KO6NVY>W3XSC-1 <RR C P R2>"},
		{false, ctlREJ | pfBit | 2<<5, "KO6NVY>W3XSC-1 <REJ R F R2>"},
		{true, ctlDISC | pfBit, "KO6NVY>W3XSC-1 <DISC C P>"},
		{false, ctlDM | pfBit, "KO6NVY>W3XSC-1 <DM R F>"},
	} {
		f := &frame{src: ko6nvy, dest: w3xsc, command: tt.command, control: tt.control}
		if got := f.String(); got != tt.want {
			t.Errorf("control %02X command=%v: got %q, want %q", tt.control, tt.command, got, tt.want)
		}
	}
}

func TestMonitor(t *testing.T) {
	var (
		mu    sync.Mutex
		lines []string
	)
	monitor := func(line string) {
		mu.Lock()
		lines = append(lines, line)
		mu.Unlock()
	}
	client, server := net.Pipe()
	b := newFakeBBS(t, server)
	go func() {
		b.accept()
		b.write("Welcome\n(#1) >\n")
		b.waitFor("B")
		b.disconnect()
	}()
	tr, err := open(client, testDevice, testClient, testBBS, nil, monitor)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = tr.ReadUntil("(#1) >\n"); err != nil {
		t.Fatal(err)
	}
	tr.Send("B")
	tr.ReadUntil("\n")
	tr.Close()
	mu.Lock()
	defer mu.Unlock()
	// Strip the timestamps ("15:04:05") from the frame lines.
	var got []string
	for _, line := range lines {
		if len(line) > 8 && line[2] == ':' && line[5] == ':' {
			line = line[8:]
		}
		got = append(got, line)
	}
	for _, want := range []string{
		"T KISS TXDELAY 40 (400ms), PERSIST 63, SLOTTIME 10 (100ms)",
		"T KC6RSC>W6XSC-1 <SABM C P>",
		"R W6XSC-1>KC6RSC <UA R F>",
		"R W6XSC-1>KC6RSC <I C R0 S0>",
		"Welcome",
		"(#1) >",
		"T KC6RSC>W6XSC-1 <I C R1 S0>",
		"B",
		"R W6XSC-1>KC6RSC <DISC C P>",
		"T KC6RSC>W6XSC-1 <UA R F>",
	} {
		if !slices.Contains(got, want) {
			t.Errorf("monitor is missing %q; got:\n%s", want, strings.Join(got, "\n"))
		}
	}
}
