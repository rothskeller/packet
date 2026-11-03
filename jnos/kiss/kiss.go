// Package kiss provides a transport layer for communicating with a JNOS BBS
// over RF, by way of a serial connection to a TNC in KISS mode (such as the
// one built into the Kenwood TH-D75).  A KISS TNC only sends and receives raw
// AX.25 frames, so this package implements the AX.25 version 2.0 connected
// mode link itself.  The Connect function, if successful, returns a jnos.Conn
// for interaction with the BBS.
package kiss

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"

	"github.com/rothskeller/packet/v4/jnos"
)

// trace is a flag indicating that AX.25 frames sent and received should be
// logged.  It is normally false, but can be set to true for debugging.
const trace = false

// readTimeout is how long ReadUntil waits without hearing from the BBS.
const readTimeout = time.Minute

// closeTimeout is the longest Close waits for the link to shut down cleanly.
const closeTimeout = 2 * time.Minute

// maxDiscRetries is the maximum number of times to retry a disconnect request.
const maxDiscRetries = 3

// bitsPerSecond is the assumed RF data rate, used to estimate how long packets
// take to transmit.
const bitsPerSecond = 1200

// Link states.
const (
	stateConnecting = iota
	stateConnected
	stateDisconnecting
	stateDisconnected
)

// Connect connects to the JNOS BBS at bbsAddress, by way of a KISS TNC of the
// supplied type attached to serialPort, and returns an open jnos.Conn for
// interaction with it.  (bbsAddress should consist of a call sign, a dash, and
// a small integer SSID.)  It logs into the specified BBS mailbox, which is also
// used as the AX.25 call sign of the connection.  If callsign is set and is
// different from mailbox, it self-identifies periodically using that call sign
// for FCC compliance.  (callsign should be the licensed FCC call sign of the
// calling user.)  If log is set, all traffic is logged to it.
func Connect(dev *Device, serialPort, bbsAddress, mailbox, callsign string, log io.Writer) (c *jnos.Conn, err error) {
	var (
		t             *Transport
		local, remote address
		port          serial.Port
	)
	if callsign == mailbox {
		callsign = "" // no need to ident
	}
	if remote, err = parseAddress(bbsAddress); err != nil {
		return nil, err
	}
	if local, err = parseAddress(mailbox); err != nil {
		return nil, fmt.Errorf("%q cannot be used as an AX.25 call sign", mailbox)
	}
	if port, err = serial.Open(serialPort, &serial.Mode{}); err != nil {
		slog.Error("serial.Open", "port", serialPort, "err", err)
		return nil, fmt.Errorf("serial.Open: %s", err)
	}
	if t, err = open(port, dev, local, remote, log); err != nil {
		return nil, err
	}
	if c, err = jnos.Connect(t, callsign); err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return c, nil
}

// Transport is the KISS TNC transport to the JNOS BBS.
type Transport struct {
	port          io.ReadWriteCloser
	dev           *Device
	local, remote address
	log           io.Writer
	frames        chan []byte   // raw frames from the reader goroutine
	readErr       error         // why the reader goroutine stopped
	sendch        chan []byte   // data from Send to the link goroutine
	closech       chan struct{} // closed when Close is called
	closeOnce     sync.Once
	done          chan struct{} // closed when the link goroutine exits
	notify        chan struct{} // signaled when state or pending changes

	mu           sync.Mutex // guards the fields below
	state        int        // written only by the link goroutine
	err          error      // why the link went down
	pending      []byte     // received data not yet returned by ReadUntil
	lastActivity time.Time  // when we last sent a request or heard the BBS
	lastCR       bool       // last received byte was a CR

	// The remaining fields are used only by the link goroutine.
	vr          byte      // V(R): next expected N(S) from the BBS
	va          byte      // V(A): oldest unacknowledged N(S)
	outstanding [][]byte  // unacknowledged data; [i] has N(S) = va+i
	nsent       int       // number of outstanding frames sent this pass
	queue       []byte    // data waiting to be packetized
	t1          time.Time // acknowledgment timer deadline, or zero
	t2          time.Time // delayed acknowledgment deadline, or zero
	retries     int       // number of T1 expirations without progress
	recovery    bool      // we've polled the BBS and await its response
	remoteBusy  bool      // the BBS sent RNR
	rejSent     bool      // we've sent REJ and await the missing frame
	closing     bool      // Close has been called
}

// open starts an AX.25 connection from local to remote through a KISS TNC
// reachable over port.  It returns once the connection is established.
func open(port io.ReadWriteCloser, dev *Device, local, remote address, log io.Writer) (t *Transport, err error) {
	t = &Transport{
		port:         port,
		dev:          dev,
		local:        local,
		remote:       remote,
		log:          log,
		frames:       make(chan []byte, 16),
		sendch:       make(chan []byte),
		closech:      make(chan struct{}),
		done:         make(chan struct{}),
		notify:       make(chan struct{}, 1),
		lastActivity: time.Now(),
	}
	go t.reader()
	go t.link()
	t.mu.Lock()
	for t.state == stateConnecting {
		t.mu.Unlock()
		<-t.notify
		t.mu.Lock()
	}
	err = t.err
	t.mu.Unlock()
	if err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return t, nil
}

// ReadUntil reads data from the BBS until the specified string is seen, or a
// timeout occurs.  It returns the data that was read (even if it returns an
// error).
func (t *Transport) ReadUntil(until string) (data string, err error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for {
		if idx := bytes.Index(t.pending, []byte(until)); idx >= 0 {
			idx += len(until)
			data = string(t.pending[:idx])
			t.pending = t.pending[idx:]
			return data, nil
		}
		if t.state == stateDisconnected {
			err = cmp.Or(t.err, jnos.ErrDisconnected)
			break
		}
		wait := readTimeout - time.Since(t.lastActivity)
		if wait <= 0 {
			slog.Error("AX.25 read", "err", "timeout")
			err = jnos.ErrTimeout
			break
		}
		t.mu.Unlock()
		timer := time.NewTimer(wait)
		select {
		case <-t.notify:
		case <-timer.C:
		}
		timer.Stop()
		t.mu.Lock()
	}
	data = string(t.pending)
	t.pending = t.pending[:0]
	return data, err
}

// Send sends a string to the BBS.  It returns once the data has been queued
// for transmission.
func (t *Transport) Send(s string) (err error) {
	t.mu.Lock()
	if t.state != stateConnected {
		err = cmp.Or(t.err, jnos.ErrDisconnected)
	} else {
		t.lastActivity = time.Now()
	}
	t.mu.Unlock()
	if err != nil {
		return err
	}
	data := strings.ReplaceAll(s, "\n", "\r")
	if !strings.HasSuffix(data, "\r") {
		data += "\r"
	}
	if t.log != nil {
		io.WriteString(t.log, strings.ReplaceAll(data, "\r", "\n"))
	}
	select {
	case t.sendch <- []byte(data):
		return nil
	case <-t.done:
		t.mu.Lock()
		defer t.mu.Unlock()
		return cmp.Or(t.err, jnos.ErrDisconnected)
	}
}

// Close disconnects from the BBS (after any queued data has been sent) and
// closes the serial port.
func (t *Transport) Close() (err error) {
	t.closeOnce.Do(func() {
		close(t.closech)
		select {
		case <-t.done:
		case <-time.After(closeTimeout):
			slog.Error("AX.25 disconnect", "err", "timeout")
		}
		err = t.port.Close()
	})
	return err
}

// UseVerboseReads returns whether it's appropriate to use verbose reads
// when communicating over this transport.
func (t *Transport) UseVerboseReads() bool { return false }

// reader reads KISS frames from the TNC and passes them to the link goroutine.
func (t *Transport) reader() {
	var (
		dec kissDecoder
		buf = make([]byte, 1024)
	)
	defer close(t.frames)
	for {
		count, err := t.port.Read(buf)
		dec.feed(buf[:count], func(raw []byte) {
			select {
			case t.frames <- raw:
			case <-t.done:
			}
		})
		if err != nil {
			t.readErr = err
			return
		}
	}
}

// link runs the AX.25 data link.  It owns all of the link state, so that
// acknowledgments and retransmissions happen on time regardless of what the
// caller is doing.
func (t *Transport) link() {
	defer close(t.done)
	closech := t.closech
	t.sendU(ctlSABM|pfBit, true)
	t.t1 = time.Now().Add(t.dev.Frack)
	for {
		t.transmit()
		if t.closing && t.state == stateConnected && len(t.queue) == 0 && len(t.outstanding) == 0 {
			t.setState(stateDisconnecting, nil)
			t.retries = 0
			t.sendU(ctlDISC|pfBit, true)
			t.t1 = time.Now().Add(t.dev.Frack)
		}
		if t.state == stateDisconnected {
			return
		}
		var timer *time.Timer
		var timerC <-chan time.Time
		if next := t.nextDeadline(); !next.IsZero() {
			timer = time.NewTimer(time.Until(next))
			timerC = timer.C
		}
		select {
		case raw, ok := <-t.frames:
			if !ok {
				slog.Error("read from TNC", "err", t.readErr)
				t.setState(stateDisconnected, fmt.Errorf("read from TNC: %s", t.readErr))
			} else {
				t.receive(raw)
			}
		case data := <-t.sendch:
			t.queue = append(t.queue, data...)
		case <-closech:
			closech, t.closing = nil, true
			if t.state == stateConnecting {
				t.setState(stateDisconnected, jnos.ErrDisconnected)
			}
		case <-timerC:
			t.timeout()
		}
		if timer != nil {
			timer.Stop()
		}
	}
}

// nextDeadline returns the earliest running timer deadline, or zero if no
// timers are running.
func (t *Transport) nextDeadline() time.Time {
	switch {
	case t.t1.IsZero():
		return t.t2
	case t.t2.IsZero() || t.t1.Before(t.t2):
		return t.t1
	default:
		return t.t2
	}
}

// setState changes the link state, recording err as the reason if the link is
// going down.
func (t *Transport) setState(state int, err error) {
	t.mu.Lock()
	t.state = state
	if state == stateConnected {
		t.lastActivity = time.Now()
	}
	if err != nil && t.err == nil {
		t.err = err
	}
	t.mu.Unlock()
	t.signal()
}

// signal notifies any waiting reader that something has changed.
func (t *Transport) signal() {
	select {
	case t.notify <- struct{}{}:
	default:
	}
}

// transmit sends I frames: first any outstanding frames that need to be
// (re)sent, and then new frames from the queue, as long as the window allows.
func (t *Transport) transmit() {
	if t.state != stateConnected || t.remoteBusy || t.recovery {
		return
	}
	for {
		if t.nsent == len(t.outstanding) {
			if len(t.outstanding) >= t.dev.MaxFrame || len(t.queue) == 0 {
				break
			}
			count := min(len(t.queue), t.dev.PacLen)
			t.outstanding = append(t.outstanding, t.queue[:count:count])
			t.queue = t.queue[count:]
		}
		ns := (t.va + byte(t.nsent)) & 7
		t.sendFrame(&frame{command: true, control: t.vr<<5 | ns<<1, info: t.outstanding[t.nsent]})
		t.nsent++
		t.t2 = time.Time{} // the I frame acknowledges what we've received
		if t.state != stateConnected {
			return
		}
	}
	if t.nsent != 0 && t.t1.IsZero() {
		t.startT1()
	}
}

// startT1 starts the acknowledgment timer, allowing for the time it takes to
// transmit the unacknowledged frames.
func (t *Transport) startT1() {
	var size int

	for _, info := range t.outstanding[:t.nsent] {
		size += len(info) + 20 // allow for addresses, flags, etc.
	}
	t.t1 = time.Now().Add(t.dev.Frack + time.Duration(size*8)*time.Second/bitsPerSecond)
}

// receive handles a frame received from the TNC.
func (t *Transport) receive(raw []byte) {
	f, err := decodeFrame(raw)
	if err != nil || f.digis || f.src != t.remote || f.dest != t.local {
		return // not for us
	}
	if trace {
		slog.Debug("AX.25 recv", "control", fmt.Sprintf("%02X", f.control), "command", f.command, "len", len(f.info))
	}
	t.mu.Lock()
	t.lastActivity = time.Now()
	t.mu.Unlock()
	switch t.state {
	case stateConnecting:
		if f.isI() || f.isS() {
			break
		}
		switch f.uType() {
		case ctlUA:
			t.vr, t.va, t.t1, t.retries = 0, 0, time.Time{}, 0
			t.setState(stateConnected, nil)
		case ctlDM:
			slog.Error("AX.25 connect", "err", "rejected")
			t.setState(stateDisconnected, errors.New("BBS rejected the connection (it may be busy)"))
		}
	case stateConnected:
		t.receiveConnected(f)
	case stateDisconnecting:
		if !f.isI() && !f.isS() && (f.uType() == ctlUA || f.uType() == ctlDM) {
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		}
	}
}

// receiveConnected handles a frame received while connected.
func (t *Transport) receiveConnected(f *frame) {
	switch {
	case f.isI():
		t.ackUpTo(f.nr())
		if f.ns() != t.vr {
			// Out of sequence: an earlier frame was lost.  Ask
			// the BBS to resend from the one we're missing.
			if !t.rejSent || f.pf() {
				t.sendS(ctlREJ, false, f.pf())
				t.rejSent = true
			}
			return
		}
		t.vr = (t.vr + 1) & 7
		t.rejSent = false
		t.deliver(f.info)
		if f.pf() {
			t.sendS(ctlRR, false, true)
		} else {
			// Wait a moment before acknowledging, in case more
			// frames are coming.
			t.t2 = time.Now().Add(t.dev.RespTime)
		}
	case f.isS():
		t.remoteBusy = f.sType() == ctlRNR
		if !t.ackUpTo(f.nr()) {
			return
		}
		if f.command && f.pf() {
			// The BBS is polling us.
			t.sendS(ctlRR, false, true)
		}
		if !f.command && f.pf() && t.recovery {
			// The BBS answered our poll; resend whatever it
			// hasn't acknowledged.
			t.recovery, t.retries, t.nsent, t.t1 = false, 0, 0, time.Time{}
		} else if f.sType() == ctlREJ {
			t.nsent, t.t1 = 0, time.Time{}
		}
	default:
		switch f.uType() {
		case ctlDISC:
			t.sendU(ctlUA|f.control&pfBit, false)
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		case ctlDM:
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		case ctlSABM:
			// The BBS reset the link.  Anything outstanding is lost.
			slog.Warn("AX.25 link reset by BBS")
			t.sendU(ctlUA|f.control&pfBit, false)
			t.vr, t.va, t.outstanding, t.nsent, t.t1 = 0, 0, nil, 0, time.Time{}
		case ctlFRMR:
			slog.Error("AX.25 frame reject", "info", fmt.Sprintf("% X", f.info))
			t.sendU(ctlDISC|pfBit, true)
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		}
	}
}

// ackUpTo processes an N(R) received from the BBS, which acknowledges all of
// our frames before it.  It returns false if N(R) is invalid.
func (t *Transport) ackUpTo(nr byte) bool {
	count := int((nr - t.va) & 7)
	if count > t.nsent {
		slog.Warn("AX.25 invalid N(R)", "nr", nr, "va", t.va, "sent", t.nsent)
		return false
	}
	if count == 0 {
		return true
	}
	t.outstanding = t.outstanding[count:]
	t.nsent -= count
	t.va = nr
	t.retries = 0
	t.t1 = time.Time{}
	if t.nsent != 0 || t.recovery {
		t.startT1()
	}
	return true
}

// timeout handles the expiration of the T1 or T2 timer.
func (t *Transport) timeout() {
	now := time.Now()
	if !t.t2.IsZero() && !now.Before(t.t2) {
		t.sendS(ctlRR, false, false)
	}
	if t.t1.IsZero() || now.Before(t.t1) {
		return
	}
	t.t1 = time.Time{}
	t.retries++
	switch t.state {
	case stateConnecting:
		if t.retries > t.dev.Retry {
			slog.Error("AX.25 connect", "err", "no response")
			t.setState(stateDisconnected, errors.New("no response from BBS"))
			return
		}
		t.sendU(ctlSABM|pfBit, true)
		t.t1 = now.Add(t.dev.Frack)
	case stateConnected:
		if t.retries > t.dev.Retry {
			slog.Error("AX.25 link", "err", "BBS stopped responding")
			t.setState(stateDisconnected, jnos.ErrDisconnected)
			return
		}
		// Poll the BBS to find out what it has received.
		t.recovery = true
		t.sendS(ctlRR, true, true)
		t.t1 = now.Add(t.dev.Frack)
	case stateDisconnecting:
		if t.retries > maxDiscRetries {
			t.setState(stateDisconnected, jnos.ErrDisconnected)
			return
		}
		t.sendU(ctlDISC|pfBit, true)
		t.t1 = now.Add(t.dev.Frack)
	}
}

// deliver adds data received from the BBS to the pending buffer, converting
// CR line endings to LF.
func (t *Transport) deliver(info []byte) {
	t.mu.Lock()
	start := len(t.pending)
	for _, c := range info {
		if c == '\n' && t.lastCR {
			t.lastCR = false
			continue
		}
		t.lastCR = c == '\r'
		if c == '\r' {
			c = '\n'
		}
		t.pending = append(t.pending, c)
	}
	if t.log != nil {
		t.log.Write(t.pending[start:])
	}
	t.mu.Unlock()
	t.signal()
}

// sendS sends a supervisory frame of the specified type, carrying our current
// N(R).  It acknowledges everything we've received.
func (t *Transport) sendS(stype byte, command, pf bool) {
	control := stype | t.vr<<5
	if pf {
		control |= pfBit
	}
	t.sendFrame(&frame{command: command, control: control})
	t.t2 = time.Time{}
}

// sendU sends an unnumbered frame.
func (t *Transport) sendU(control byte, command bool) {
	t.sendFrame(&frame{command: command, control: control})
}

// sendFrame sends a frame to the BBS through the TNC.
func (t *Transport) sendFrame(f *frame) {
	f.dest, f.src = t.remote, t.local
	if trace {
		slog.Debug("AX.25 send", "control", fmt.Sprintf("%02X", f.control), "command", f.command, "len", len(f.info))
	}
	data := kissEncode(f.encode())
	for len(data) != 0 {
		count, err := t.port.Write(data)
		if err != nil {
			slog.Error("write to TNC", "err", err)
			t.setState(stateDisconnected, fmt.Errorf("write to TNC: %s", err))
			return
		}
		data = data[count:]
	}
}
