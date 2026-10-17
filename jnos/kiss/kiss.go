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
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/serial"

	"github.com/rothskeller/packet/v4/jnos"
)

// trace is a flag indicating that AX.25 frames sent and received should be
// logged.  It is normally false, but can be set to true for debugging.
const trace = false

// readTimeout is how long ReadUntil waits without the link making progress:
// that is, without receiving data from the BBS or having our data
// acknowledged.  (Frames from the BBS that do neither don't count, so that a
// link that's failing but still exchanging supervisory frames times out.)  It
// applies only while the link is idle: while the link is still retrying our
// unacknowledged data, the link's own retry limit decides when the BBS is
// gone, so a slow BBS isn't cut off mid-recovery.  It is a variable so that
// tests can shorten it.
var readTimeout = time.Minute

// closeTimeout is the longest Close waits for the link to shut down cleanly.
const closeTimeout = 2 * time.Minute

// minWriteTimeout is the shortest time a write to the TNC may take before we
// decide the TNC is stuck and declare the link failed.  (The actual limit is
// the larger of this and the device's Frack.)  Serial port writes have no
// timeout of their own, and a stalled USB or Bluetooth device can block one
// forever.  It is a variable so that tests can shorten it.
var minWriteTimeout = 10 * time.Second

// portCloseTimeout is the longest Close waits for the serial port to close.
// It is a variable so that tests can shorten it.
var portCloseTimeout = 5 * time.Second

// connectMargin is added to the longest time the connection handshake should
// take, to give the overall limit on how long open waits for it.
const connectMargin = 5 * time.Second

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

// finalWriteTimeout is the longest we wait to hand the TNC the DISC we send when
// abandoning a link (see finalDISC).  It is a variable only so that tests can
// change it.
var finalWriteTimeout = 500 * time.Millisecond

// errAborted is returned by all Transport methods after Abort is called.
var errAborted = fmt.Errorf("connection aborted: %w", jnos.ErrDisconnected)

// errWriteTimeout is the link failure when a write to the TNC takes too long.
var errWriteTimeout = errors.New("write to TNC timed out")

// Connect connects to the JNOS BBS at bbsAddress, by way of a KISS TNC of the
// supplied type attached to serialPort, and returns an open jnos.Conn for
// interaction with it.  (bbsAddress should consist of a call sign, a dash, and
// a small integer SSID.)  It logs into the specified BBS mailbox, which is also
// used as the AX.25 call sign of the connection.  If callsign is set and is
// different from mailbox, it self-identifies periodically using that call sign
// for FCC compliance.  (callsign should be the licensed FCC call sign of the
// calling user.)  If log is set, all traffic is logged to it.  If monitor is
// set, it is called with each line of a packet monitor display of the AX.25
// frames exchanged with the BBS (e.g., "06:57:02T KO6NVY>W3XSC-1 <SABM C P>").
func Connect(dev *Device, serialPort, bbsAddress, mailbox, callsign string, log io.Writer, monitor func(line string)) (c *jnos.Conn, err error) {
	return ConnectContext(context.Background(), dev, serialPort, bbsAddress, mailbox, callsign, log, monitor)
}

// ConnectContext is Connect, except that canceling ctx before it returns
// abandons the connection attempt promptly (with the same cleanup as Close
// after an Abort, including a reset of a Kenwood TNC).  Once it returns, ctx
// no longer affects the connection.
func ConnectContext(ctx context.Context, dev *Device, serialPort, bbsAddress, mailbox, callsign string, log io.Writer, monitor func(line string)) (c *jnos.Conn, err error) {
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
	// Raise DTR and RTS: the TH-D75 (like many USB TNCs) ignores data from
	// the computer while DTR is low, and macOS leaves it low by default.
	// Some ports (e.g., Bluetooth ones) don't have these lines; that's fine.
	if err := port.SetDTR(true); err != nil {
		slog.Debug("serial.SetDTR", "port", serialPort, "err", err)
	}
	if err := port.SetRTS(true); err != nil {
		slog.Debug("serial.SetRTS", "port", serialPort, "err", err)
	}
	var opts []func(*Transport)
	if dev.KenwoodReset {
		// If the last session on this port never finished (the
		// program was killed in the middle of it, say), the TNC may
		// still hold frames from it; reset it before starting.
		marker := sessionMarker(serialPort)
		if marker != "" {
			if _, err := os.Stat(marker); err == nil {
				if err = resetBeforeStart(port, monitor); err != nil {
					port.Close()
					return nil, err
				}
			}
			markSession(marker)
		}
		opts = append(opts, func(t *Transport) {
			t.marker = marker
			t.reopen = func() (resetPort, error) {
				port, err := serial.Open(serialPort, &serial.Mode{})
				if err == nil {
					port.SetDTR(true)
					port.SetRTS(true)
				}
				return port, err
			}
		})
	}
	opts = append(opts, withContext(ctx))
	if t, err = open(port, dev, local, remote, log, monitor, opts...); err != nil {
		return nil, err
	}
	c, err = jnos.Connect(t, callsign)
	t.stopContext()
	if err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return c, nil
}

// withContext is an option for open: canceling ctx aborts the Transport, until
// stopContext is called.
func withContext(ctx context.Context) func(*Transport) {
	return func(t *Transport) { t.ctx = ctx }
}

// stopContext stops the context passed to withContext from aborting the
// Transport.
func (t *Transport) stopContext() {
	if t.stopCtx != nil {
		t.stopCtx()
		t.stopCtx = nil
	}
}

// Transport is the KISS TNC transport to the JNOS BBS.
type Transport struct {
	port             io.ReadWriteCloser
	dev              *Device // validated copy of the caller's Device
	local, remote    address
	log              io.Writer
	monitor          func(line string)
	readTimeout      time.Duration // see readTimeout
	writeTimeout     time.Duration // see minWriteTimeout
	portCloseTimeout time.Duration // see portCloseTimeout
	frames           chan []byte   // raw frames from the reader goroutine
	readErr          error         // why the reader goroutine stopped
	writech          chan []byte   // data from the link to the writer goroutine
	writeErr         chan error    // results of writes, from the writer goroutine
	wake             chan struct{} // signaled when Send queues data
	closech          chan struct{} // closed when Close is called
	closeOnce        sync.Once
	abortch          chan struct{} // closed when Abort is called
	abortOnce        sync.Once
	done             chan struct{} // closed when the link goroutine exits
	notify           chan struct{} // signaled when state or pending changes
	portOnce         sync.Once
	portClosed       chan struct{}             // closed when port.Close returns
	portErr          error                     // result of port.Close
	reopen           func() (resetPort, error) // reopens the port to reset a Kenwood TNC, or nil
	cleanEnd         atomic.Bool               // the session ended with a clean disconnection
	marker           string                    // file marking a session in progress (see sessionMarker), or ""
	ctx              context.Context           // aborts the connection attempt when canceled (see withContext), or nil
	stopCtx          func() bool               // stops ctx from aborting us
	lastSent         atomic.Int64              // when we last wrote a frame to the TNC (UnixNano)
	wakeups          atomic.Int64              // link goroutine loop count (for testing)

	outMu  sync.Mutex // guards log and monitor output, and the field below
	outOff bool       // Close has returned, so there must be no more output

	mu           sync.Mutex // guards the fields below
	state        int        // written only by the link goroutine
	err          error      // why the link went down
	aborted      bool       // Abort has been called
	pending      []byte     // received data not yet returned by ReadUntil
	sendq        []byte     // data from Send not yet taken by the link
	lastProgress time.Time  // when we last sent a request, received data, or had data acknowledged
	retrying     bool       // the link is still retrying our data (written by the link goroutine)
	lastCR       bool       // last received byte was a CR

	// The remaining fields are used only by the link goroutine.
	vr            byte      // V(R): next expected N(S) from the BBS
	va            byte      // V(A): oldest unacknowledged N(S)
	outstanding   [][]byte  // unacknowledged data; [i] has N(S) = va+i
	nsent         int       // number of outstanding frames sent this pass
	queue         []byte    // data waiting to be packetized
	t1            time.Time // acknowledgment timer deadline, or zero
	ackAt         time.Time // when to send a plain acknowledgment (quiet timer)
	holdUntil     time.Time // don't transmit before this (BBS turnaround)
	t1Origin      time.Time // T1's deadline when last started (before pushes)
	t1Pushed      time.Time // the deadline T1 was last pushed to
	respond       byte      // response owed to the BBS (ctlRR or ctlREJ), or 0
	respondF      bool      // the owed response must have the F bit set
	retries       int       // number of T1 expirations without progress
	recovery      bool      // we've polled the BBS and await its response
	remoteBusy    bool      // the BBS sent RNR
	rejSent       bool      // we've sent REJ and await the missing frame
	closing       bool      // Close has been called
	lastFrame     *frame    // the frame from the BBS being handled
	pendingResend bool      // resend unacknowledged data once the hold ends (see resendIfPolledWithout)
	lastDataSent  time.Time // when we last sent (or resent) an I frame
	linkEnded     bool      // the BBS has ended the link, or our DISCs have run out
	portFailed    bool      // reading from or writing to the TNC failed

	// Round trip timing (see t1Base).
	rtts      []time.Duration // the BBS's last rttWindow answer times
	notedT1   time.Duration   // the T1 interval last reported to the monitor
	pollAt    time.Time       // when we sent our outstanding poll, or zero
	pollAir   time.Duration   // the airtime of that poll
	pollMulti bool            // we've sent more than one poll since the last answer
	timing    bool            // timing the acknowledgment of the I frame below
	timedNS   byte            // N(S) of the I frame being timed
	timedAt   time.Time       // when it was sent
	timedAir  time.Duration   // its airtime
}

// open starts an AX.25 connection from local to remote through a KISS TNC
// reachable over port.  It returns once the connection is established.
func open(port io.ReadWriteCloser, dev *Device, local, remote address, log io.Writer, monitor func(string), opts ...func(*Transport)) (t *Transport, err error) {
	dev = dev.validated()
	t = &Transport{
		port:             port,
		dev:              dev,
		local:            local,
		remote:           remote,
		log:              log,
		monitor:          monitor,
		readTimeout:      readTimeout,
		writeTimeout:     max(dev.Frack, minWriteTimeout),
		portCloseTimeout: portCloseTimeout,
		frames:           make(chan []byte, 16),
		writech:          make(chan []byte),
		writeErr:         make(chan error, 1),
		wake:             make(chan struct{}, 1),
		closech:          make(chan struct{}),
		abortch:          make(chan struct{}),
		done:             make(chan struct{}),
		notify:           make(chan struct{}, 1),
		portClosed:       make(chan struct{}),
		lastProgress:     time.Now(),
	}
	for _, opt := range opts {
		opt(t)
	}
	go t.reader()
	go t.writer()
	go t.link()
	if t.ctx != nil {
		t.stopCtx = context.AfterFunc(t.ctx, t.Abort)
	}
	// The link gives up after Retry+1 unanswered SABMs, but don't count on
	// it: wait no longer than that (plus a margin) in any case.
	deadline := time.NewTimer(time.Duration(dev.Retry+2)*dev.InitialFrack + t.writeTimeout + connectMargin)
	defer deadline.Stop()
	t.mu.Lock()
	for t.state == stateConnecting && err == nil {
		t.mu.Unlock()
		select {
		case <-t.notify:
		case <-deadline.C:
			slog.Error("AX.25 connect", "err", "timeout")
			err = errors.New("timed out")
		}
		t.mu.Lock()
	}
	if t.err != nil {
		err = t.err
	}
	t.mu.Unlock()
	if err != nil {
		t.stopContext()
		t.Abort()
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
		if t.aborted {
			err = errAborted
			break
		}
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
		wait := t.readTimeout - time.Since(t.lastProgress)
		if wait <= 0 && t.retrying {
			// The link will either get our data through or give
			// up (and disconnect); wait for that.
			wait = t.readTimeout
		} else if wait <= 0 {
			slog.Error("AX.25 read", "err", "timeout")
			err = jnos.ErrTimeout
			break
		}
		t.mu.Unlock()
		timer := time.NewTimer(wait)
		select {
		case <-t.notify:
		case <-t.abortch:
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
// for transmission; it never waits for the link goroutine.
func (t *Transport) Send(s string) (err error) {
	data := strings.ReplaceAll(s, "\n", "\r")
	if !strings.HasSuffix(data, "\r") {
		data += "\r"
	}
	t.mu.Lock()
	switch {
	case t.aborted:
		err = errAborted
	case t.state != stateConnected:
		err = cmp.Or(t.err, jnos.ErrDisconnected)
	default:
		t.lastProgress = time.Now()
		t.writeLog([]byte(strings.ReplaceAll(data, "\r", "\n")))
		t.sendq = append(t.sendq, data...)
	}
	t.mu.Unlock()
	if err != nil {
		return err
	}
	select {
	case t.wake <- struct{}{}:
	default: // the link goroutine already has a wakeup pending
	}
	return nil
}

// Close disconnects from the BBS (after any queued data has been sent) and
// closes the serial port.  If the link doesn't shut down cleanly within
// closeTimeout, or the port doesn't close within portCloseTimeout, it gives up
// on them.  Either way, nothing is written to the log or the monitor after it
// returns.
func (t *Transport) Close() (err error) {
	t.closeOnce.Do(func() {
		close(t.closech)
		timer := time.NewTimer(closeTimeout)
		select {
		case <-t.done:
		case <-timer.C:
			slog.Error("AX.25 disconnect", "err", "timeout")
			t.Abort() // stop the link goroutine
		}
		timer.Stop()
		err = t.closePort()
		if t.reopen == nil || t.cleanEnd.Load() || t.resetTNC() {
			// The TNC is in good order; the next session needn't
			// reset it.
			clearSession(t.marker)
		}
		t.outMu.Lock()
		t.outOff = true
		t.outMu.Unlock()
	})
	return err
}

// resetTNC resets a Kenwood TNC after a session that didn't end cleanly (see
// kenwoodReset), so that it doesn't transmit frames it still holds after we've
// given up, and starts the next session clean.  It is best-effort, and reports
// the result to the monitor; it returns whether it succeeded.  The port must
// already be closed.
func (t *Transport) resetTNC() bool {
	if last := t.lastSent.Load(); last != 0 {
		if wait := resetSettle - time.Since(time.Unix(0, last)); wait > 0 {
			time.Sleep(wait) // let our last frame get on the air
		}
	}
	port, err := t.reopen()
	if err != nil {
		slog.Warn("TNC reset", "err", err)
		t.note(fmt.Sprintf("TNC reset failed: can't reopen the serial port: %s", err))
		return false
	}
	defer port.Close()
	cmd, err := kenwoodReset(port)
	if err != nil {
		slog.Warn("TNC reset", "err", err)
		t.note(fmt.Sprintf("TNC reset failed: %s.  Reset it on the radio: turn the TNC off and back to KISS.", err))
		return false
	}
	slog.Info("TNC reset", "reentered", cmd)
	t.note("TNC reset: left KISS mode, discarding anything it still held, and re-entered it with " + cmd)
	return true
}

// note reports a status line to the monitor, if there is one.
func (t *Transport) note(msg string) {
	t.output(func() { monitorNote(t.monitor, msg) })
}

// monitorNote reports a status line to monitor, if it isn't nil.
func monitorNote(monitor func(string), msg string) {
	if monitor != nil {
		monitor("*** " + time.Now().Format("15:04:05") + " " + msg)
	}
}

// Abort ends the connection immediately, without an orderly disconnection from
// the BBS.  Any Transport call in progress returns promptly with an error, as
// does any later call (except Close, which must still be called).
func (t *Transport) Abort() {
	t.abortOnce.Do(func() {
		t.mu.Lock()
		t.aborted = true
		if t.err == nil {
			t.err = errAborted
		}
		t.mu.Unlock()
		close(t.abortch)
		t.signal()
		// The link goroutine sends a final DISC as it exits (see
		// finalDISC); close the port once it has, or once it's had
		// long enough.
		go func() {
			timer := time.NewTimer(2 * finalWriteTimeout)
			defer timer.Stop()
			select {
			case <-t.done:
			case <-timer.C:
			}
			t.startPortClose()
		}()
	})
}

// finalDISC is called as the link goroutine exits.  If we're abandoning a link
// that the BBS may still think is up (because the user aborted, the BBS stopped
// answering, or we gave up connecting after the BBS may have heard us), it
// sends a DISC, without waiting for an answer, so that the BBS can end its
// side of the session rather than wait for its own timers.  (If the BBS can
// hear us but we can't hear it, that's the only way it learns we're gone.)  It
// waits no longer than finalWriteTimeout to hand the DISC to the TNC.
func (t *Transport) finalDISC() {
	if t.linkEnded || t.portFailed {
		return
	}
	f := &frame{dest: t.remote, src: t.local, command: true, control: ctlDISC | pfBit}
	t.monitorFrame('T', f)
	if t.finalWrite(kissEncode(f.encode())) {
		t.lastSent.Store(time.Now().UnixNano())
	}
}

// finalWrite writes data to the TNC as the link goroutine exits, waiting no
// longer than finalWriteTimeout.  It returns whether the write succeeded.
func (t *Transport) finalWrite(data []byte) bool {
	timer := time.NewTimer(finalWriteTimeout)
	defer timer.Stop()
	select {
	case t.writech <- data:
	case <-timer.C:
		return false
	}
	select {
	case err := <-t.writeErr:
		return err == nil
	case <-timer.C:
		return false
	}
}

// startPortClose starts closing the serial port, if that hasn't already been
// started.  It doesn't wait: a stalled device can block port.Close.
func (t *Transport) startPortClose() {
	t.portOnce.Do(func() {
		go func() {
			t.portErr = t.port.Close()
			close(t.portClosed)
		}()
	})
}

// closePort closes the serial port, waiting no longer than portCloseTimeout.
func (t *Transport) closePort() error {
	t.startPortClose()
	timer := time.NewTimer(t.portCloseTimeout)
	defer timer.Stop()
	select {
	case <-t.portClosed:
		return t.portErr
	case <-timer.C:
		slog.Error("close TNC port", "err", "timeout")
		return errors.New("close of TNC port timed out")
	}
}

// isAborted returns whether Abort has been called.
func (t *Transport) isAborted() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.aborted
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

// writer writes data from the link goroutine to the TNC, reporting the result
// of each write on writeErr.  The writes happen here, rather than in the link
// goroutine, so that the link goroutine can give up on one that's stuck.
func (t *Transport) writer() {
	for {
		select {
		case data := <-t.writech:
			var err error
			for len(data) != 0 && err == nil {
				var count int
				if count, err = t.port.Write(data); count == 0 && err == nil {
					err = io.ErrShortWrite
				}
				data = data[count:]
			}
			select {
			case t.writeErr <- err:
			case <-t.done:
				return
			}
		case <-t.done:
			return
		}
	}
}

// link runs the AX.25 data link.  It owns all of the link state, so that
// acknowledgments and retransmissions happen on time regardless of what the
// caller is doing.
func (t *Transport) link() {
	defer close(t.done)
	defer t.finalDISC()
	closech := t.closech
	t.configureTNC()
	t.sendU(ctlSABM|pfBit, true)
	t.t1 = time.Now().Add(t.t1Base())
	t.notedT1 = t.t1Base()
	for {
		t.wakeups.Add(1)
		select {
		case <-t.abortch:
			t.setState(stateDisconnected, errAborted)
			return
		default:
		}
		t.takeSent()
		t.sendResponse()
		t.resendPolled()
		t.transmit()
		// Unacknowledged data must always have T1 running (unless it's
		// about to be sent), or the link would stall silently, waiting
		// for an event that never comes.  This guards against any path
		// that clears T1 without sending.
		if t.state == stateConnected && len(t.outstanding) != 0 && t.t1.IsZero() && !t.haveIFrames() && !t.pendingResend {
			slog.Warn("AX.25 link: T1 not running with data outstanding; restarting it")
			t.t1 = time.Now().Add(t.frack())
		}
		t.busyTimer()
		if t.closing && t.state == stateConnected && len(t.queue) == 0 && len(t.outstanding) == 0 && !t.held() {
			t.setState(stateDisconnecting, nil)
			t.retries = 0
			t.sendU(ctlDISC|pfBit, true)
			t.t1 = time.Now().Add(t.t1Base())
		}
		if t.state == stateDisconnected {
			return
		}
		t.setRetrying(t.state == stateConnected && (len(t.outstanding) != 0 || len(t.queue) != 0 || t.recovery))
		var timer *time.Timer
		var timerC <-chan time.Time
		if next := t.nextDeadline(); !next.IsZero() {
			timer = time.NewTimer(time.Until(next))
			timerC = timer.C
		}
		select {
		case raw, ok := <-t.frames:
			if !ok {
				t.portFailed = true
				if !t.isAborted() {
					slog.Error("read from TNC", "err", t.readErr)
				}
				t.setState(stateDisconnected, fmt.Errorf("read from TNC: %s", t.readErr))
			} else {
				t.receive(raw)
			}
		case <-t.wake:
			// takeSent will pick up the data.
		case <-closech:
			closech, t.closing = nil, true
			if t.state == stateConnecting {
				t.setState(stateDisconnected, jnos.ErrDisconnected)
			}
		case <-t.abortch:
			// handled at the top of the loop
		case <-timerC:
			t.timeout()
		}
		if timer != nil {
			timer.Stop()
		}
	}
}

// takeSent moves data queued by Send into the link goroutine's queue.
func (t *Transport) takeSent() {
	t.mu.Lock()
	if len(t.sendq) != 0 {
		t.queue = append(t.queue, t.sendq...)
		t.sendq = t.sendq[:0]
	}
	t.mu.Unlock()
}

// nextDeadline returns the earliest time the link goroutine has something to
// do: T1, the quiet timer for a pending acknowledgment, or the end of the
// turnaround hold when something is waiting for it.  It returns zero if there
// is none.  The quiet timer and the hold count only while connected, since
// nothing they hold back is sent otherwise.
func (t *Transport) nextDeadline() (next time.Time) {
	next = t.t1
	earlier := func(d time.Time) {
		if !d.IsZero() && (next.IsZero() || d.Before(next)) {
			next = d
		}
	}
	if t.state != stateConnected {
		return next
	}
	if t.respond != 0 {
		earlier(t.ackAt)
	}
	if t.held() && (t.respond != 0 || t.haveIFrames() || t.pendingResend || (t.closing && len(t.queue) == 0 && len(t.outstanding) == 0)) {
		earlier(t.holdUntil)
	}
	return next
}

// held returns whether we must not transmit yet, because the BBS may not have
// finished switching from transmitting to receiving.
func (t *Transport) held() bool {
	return time.Now().Before(t.holdUntil)
}

// owesFinal returns whether we owe the BBS a response with the F bit (to its
// poll), which must go before any I frame: an I frame can't carry the F bit.
// (sendResponse sends it first in each pass of the link loop, but if the
// turnaround hold ends between its check and a later one in the same pass, the
// I frame would otherwise go first.)
func (t *Transport) owesFinal() bool {
	return t.respond != 0 && t.respondF
}

// haveIFrames returns whether transmit has I frames ready to send.
func (t *Transport) haveIFrames() bool {
	if t.state != stateConnected || t.remoteBusy || t.recovery {
		return false
	}
	return t.nsent < len(t.outstanding) || (len(t.queue) != 0 && len(t.outstanding) < t.dev.MaxFrame)
}

// busyTimer keeps T1 running while the BBS is busy and we have data for it,
// so that T1's expiration polls the BBS to find out whether it's still busy.
// (Otherwise, if the RR saying it's ready again were lost, we'd wait forever.)
func (t *Transport) busyTimer() {
	if t.state == stateConnected && t.remoteBusy && t.t1.IsZero() && (len(t.outstanding) != 0 || len(t.queue) != 0) {
		t.t1 = time.Now().Add(t.frack())
	}
}

// sendResponse sends the response owed to the BBS, if any.  Nothing goes
// until the turnaround hold has passed.  Then a REJ, or a response to a poll,
// goes right away (the TNC still waits for a clear channel).  A plain
// acknowledgment also waits for the quiet timer, in case the BBS has more
// frames to send, unless we have I frames to send, which carry it.
func (t *Transport) sendResponse() {
	if t.respond == 0 || t.state != stateConnected || t.held() {
		return
	}
	if t.respond == ctlRR && !t.respondF && (t.haveIFrames() || time.Now().Before(t.ackAt)) {
		return
	}
	t.sendS(t.respond, false, t.respondF)
}

// configureTNC sends the device's channel access settings to the TNC.
func (t *Transport) configureTNC() {
	var sent []string

	for _, setting := range []struct {
		name  string
		cmd   byte
		value int
		unit  time.Duration
	}{
		{"TXDELAY", kissTXDelay, int(t.dev.TXDelay / (10 * time.Millisecond)), 10 * time.Millisecond},
		{"PERSIST", kissPersist, t.dev.Persist, 0},
		{"SLOTTIME", kissSlotTime, int(t.dev.SlotTime / (10 * time.Millisecond)), 10 * time.Millisecond},
	} {
		if setting.value > 0 {
			value := min(setting.value, 255)
			t.writeTNC(kissFrame(setting.cmd, []byte{byte(value)}))
			if setting.unit != 0 {
				sent = append(sent, fmt.Sprintf("%s %d (%v)", setting.name, value, time.Duration(value)*setting.unit))
			} else {
				sent = append(sent, fmt.Sprintf("%s %d", setting.name, value))
			}
		}
	}
	if len(sent) != 0 && t.state != stateDisconnected {
		t.monitorLine('T', "KISS "+strings.Join(sent, ", "))
	}
}

// output calls fn, which writes to the log or calls the monitor, unless Close
// has returned.  It serializes all such output.
func (t *Transport) output(fn func()) {
	t.outMu.Lock()
	defer t.outMu.Unlock()
	if !t.outOff {
		fn()
	}
}

// writeLog writes data to the log, if there is one.
func (t *Transport) writeLog(data []byte) {
	if t.log != nil {
		t.output(func() { t.log.Write(data) })
	}
}

// monitorLine reports a line to the packet monitor, if there is one, prefixed
// with the time and the direction ('T' for transmitted, 'R' for received).
func (t *Transport) monitorLine(direction byte, line string) {
	if t.monitor != nil {
		line = time.Now().Format("15:04:05") + string(direction) + " " + line
		t.output(func() { t.monitor(line) })
	}
}

// monitorFrame reports a frame to the packet monitor, if there is one,
// followed by the text it carries.
func (t *Transport) monitorFrame(direction byte, f *frame) {
	if t.monitor == nil {
		return
	}
	lines := []string{time.Now().Format("15:04:05") + string(direction) + " " + f.String()}
	if f.hasInfo() && len(f.info) != 0 {
		text := strings.TrimRight(strings.ReplaceAll(strings.ReplaceAll(string(f.info), "\r\n", "\n"), "\r", "\n"), "\n")
		lines = append(lines, strings.Split(text, "\n")...)
	}
	t.output(func() {
		for _, line := range lines {
			t.monitor(line)
		}
	})
}

// setState changes the link state, recording err as the reason if the link is
// going down.
func (t *Transport) setState(state int, err error) {
	t.mu.Lock()
	t.state = state
	if state == stateConnected {
		t.lastProgress = time.Now()
	}
	if err != nil && t.err == nil {
		t.err = err
	}
	t.mu.Unlock()
	if state != stateConnected {
		// Responses are sent only while connected; don't leave one
		// owed (and its quiet timer running) after that.
		t.respond, t.respondF = 0, false
	}
	t.signal()
}

// signal notifies any waiting reader that something has changed.
func (t *Transport) signal() {
	select {
	case t.notify <- struct{}{}:
	default:
	}
}

// progress records that the link has made progress, for ReadUntil's timeout.
// setRetrying records whether the link is still trying to get our data to the
// BBS, for ReadUntil.  When that stops, it wakes ReadUntil to recompute its
// timeout.
func (t *Transport) setRetrying(retrying bool) {
	t.mu.Lock()
	changed := t.retrying != retrying
	t.retrying = retrying
	t.mu.Unlock()
	if changed && !retrying {
		t.signal()
	}
}

func (t *Transport) progress() {
	t.mu.Lock()
	t.lastProgress = time.Now()
	t.mu.Unlock()
}

// resetLink resets the link state variables to their initial values, as when
// the link is (re)established.
func (t *Transport) resetLink() {
	t.vr, t.va, t.outstanding, t.nsent, t.t1 = 0, 0, nil, 0, time.Time{}
	t.respond, t.respondF, t.retries = 0, false, 0
	t.recovery, t.remoteBusy, t.rejSent = false, false, false
	t.pollAt, t.pollMulti, t.timing, t.pendingResend = time.Time{}, false, false, false
}

// transmit sends I frames: first any outstanding frames that need to be
// (re)sent, and then new frames from the queue, as long as the window allows.
func (t *Transport) transmit() {
	if !t.haveIFrames() || t.held() || t.owesFinal() {
		return
	}
	for {
		var fresh bool
		if t.nsent == len(t.outstanding) {
			if len(t.outstanding) >= t.dev.MaxFrame || len(t.queue) == 0 {
				break
			}
			count := min(len(t.queue), t.dev.PacLen)
			t.outstanding = append(t.outstanding, t.queue[:count:count])
			t.queue = t.queue[count:]
			fresh = true
		}
		ns := (t.va + byte(t.nsent)) & 7
		info := t.outstanding[t.nsent]
		t.sendFrame(&frame{command: true, control: t.vr<<5 | ns<<1, info: info})
		t.lastDataSent = time.Now()
		// Time the acknowledgment of a frame sent only once (an
		// acknowledgment of a resent one is ambiguous).
		if fresh && !t.timing {
			t.timing, t.timedNS, t.timedAt, t.timedAir = true, ns, time.Now(), airtime(len(info))
		} else if !fresh && t.timing && t.timedNS == ns {
			t.timing = false
		}
		t.nsent++
		if t.respond == ctlRR && !t.respondF {
			t.respond = 0 // the I frame acknowledges what we've received
		}
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
	t.t1 = time.Now().Add(t.frack() + time.Duration(size*8)*time.Second/bitsPerSecond)
}

// maxFrack caps the backed-off T1 interval.
const maxFrack = 30 * time.Second

// airtime returns the time it takes to transmit a frame carrying n bytes of
// information, allowing for addresses, flags, etc.
func airtime(n int) time.Duration {
	return time.Duration((n+20)*8) * time.Second / bitsPerSecond
}

// frack returns the current T1 interval while connected: the base interval
// (see t1Base), doubled for each retry without progress (up to four times the
// base, but no more than maxFrack).  A slow BBS needs the backoff: retrying at
// a fixed short interval only crowds the channel and collides with the BBS's
// own retransmissions.
func (t *Transport) frack() time.Duration {
	base := t.t1Base()
	limit := max(min(4*base, maxFrack), base)
	return min(base<<min(t.retries, 2), limit)
}

// rttWindow is the number of recent BBS answer times that T1 is based on, and
// rttMinSamples the number needed before they're used.
const (
	rttWindow     = 5
	rttMinSamples = 3
)

// t1Base returns the base T1 interval: how long to wait for the BBS to answer
// before deciding that our frame (or its answer) was lost.  Until the BBS has
// answered rttMinSamples times, that's the device's InitialFrack, a patient
// guess.  After that, it's based on the last rttWindow answer times: the larger
// of twice their median (the BBS's typical answer time, which a single fast or
// slow answer doesn't move) and one and a half times the slowest of them (so
// that a BBS that sometimes takes much longer than usual, as W3XSC does when
// its channel is busy, doesn't get our retransmission while its answer is
// still on the way: a duplicate, which it rejects, sending the link into
// recovery for nothing).  A slow answer raises the interval at once; it comes
// down only once rttWindow faster answers have pushed the slow one out.  The
// interval is never shorter than Frack, and never longer than maxFrack.
func (t *Transport) t1Base() time.Duration {
	if len(t.rtts) < rttMinSamples {
		return max(t.dev.InitialFrack, t.dev.Frack)
	}
	median, slowest := rttStats(t.rtts)
	return min(max(t.dev.Frack, 2*median, slowest*3/2), maxFrack)
}

// rttStats returns the median and maximum of the answer times in rtts, which
// must not be empty.
func rttStats(rtts []time.Duration) (median, slowest time.Duration) {
	sorted := slices.Sorted(slices.Values(rtts))
	n := len(sorted)
	if n%2 == 1 {
		median = sorted[n/2]
	} else {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return median, sorted[n-1]
}

// sampleRTT records a measurement of the time the BBS took to answer us.
func (t *Transport) sampleRTT(rtt time.Duration) {
	rtt = max(rtt, 0)
	t.rtts = append(t.rtts, rtt)
	if len(t.rtts) > rttWindow {
		t.rtts = slices.Delete(t.rtts, 0, len(t.rtts)-rttWindow)
	}
	if trace {
		slog.Debug("AX.25 RTT", "sample", rtt, "window", t.rtts)
	}
	if base := t.t1Base(); (base - t.notedT1).Abs() >= time.Second {
		t.notedT1 = base
		median, slowest := rttStats(t.rtts)
		t.note(fmt.Sprintf("retry timer now %s (the BBS's last %d answers took %s typically, %s at most)",
			base.Round(100*time.Millisecond), len(t.rtts),
			median.Round(100*time.Millisecond), slowest.Round(100*time.Millisecond)))
	}
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
	t.monitorFrame('R', f)
	t.lastFrame = f
	// Any frame from the BBS restarts the quiet timer for a pending
	// acknowledgment, since the BBS may be sending a burst of frames, and
	// the turnaround hold before we transmit anything.
	t.ackAt = time.Now().Add(t.dev.RespTime)
	t.holdUntil = time.Now().Add(t.dev.Turnaround)
	// A BBS that's transmitting is alive, and retransmitting over it
	// would only collide with it: any frame from the BBS pushes T1 out,
	// but by no more than one interval past when T1 was started, so that a
	// BBS that keeps transmitting without making progress can't keep us
	// from ever retrying (and giving up).  (Handling the frame below may
	// stop or restart T1 anyway.)
	if t.state == stateConnected && !t.t1.IsZero() {
		if !t.t1.Equal(t.t1Pushed) {
			t.t1Origin = t.t1 // T1 was restarted since our last push
		}
		pushed := time.Now().Add(t.frack())
		if limit := t.t1Origin.Add(t.frack()); pushed.After(limit) {
			pushed = limit
		}
		if pushed.After(t.t1) {
			t.t1, t.t1Pushed = pushed, pushed
		}
	}
	switch t.state {
	case stateConnecting:
		if f.isI() || f.isS() {
			// The BBS is sending us data or acknowledgments on a new
			// link (N(S) and N(R) both 0): it accepted our SABM, and
			// we missed its UA.  Rather than wait for T1 and send
			// another SABM (which would reset the link the BBS has
			// already started using), take the link as up and
			// handle the frame.
			if f.nr() == 0 && (f.isS() || f.ns() == 0) {
				t.resetLink()
				t.setState(stateConnected, nil)
				t.note("connected (the BBS is sending to us, so it got our SABM; its UA was lost)")
				t.receiveConnected(f)
			}
			break
		}
		switch f.uType() {
		case ctlUA:
			// If we sent only one SABM, the UA's delay is the first
			// measure of how long the BBS takes to answer.
			sampled := !t.pollAt.IsZero() && !t.pollMulti
			rtt := time.Since(t.pollAt) - t.pollAir
			t.resetLink()
			t.setState(stateConnected, nil)
			if sampled {
				t.sampleRTT(rtt)
			}
		case ctlDM:
			t.linkEnded = true
			slog.Error("AX.25 connect", "err", "rejected")
			t.setState(stateDisconnected, errors.New("BBS rejected the connection (it may be busy)"))
		}
	case stateConnected:
		t.receiveConnected(f)
	case stateDisconnecting:
		if !f.isI() && !f.isS() && (f.uType() == ctlUA || f.uType() == ctlDM) {
			// Our orderly disconnection (sent only once all our data
			// was acknowledged) completed.
			t.cleanEnd.Store(true)
			t.linkEnded = true
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		}
	}
}

// receiveConnected handles a frame received while connected.
func (t *Transport) receiveConnected(f *frame) {
	switch {
	case f.isI():
		t.ackUpTo(f.nr())
		t.resendIfPolledWithout()
		if behind, ahead := (t.vr-f.ns())&7, (f.ns()-t.vr)&7; f.ns() != t.vr && behind <= ahead {
			// A frame we already have: the BBS resent it because
			// it didn't hear our acknowledgment.  Acknowledge it
			// again (not REJ: nothing is missing).
			if t.respond == 0 {
				t.respond = ctlRR
			}
			t.respondF = t.respondF || f.pf()
			return
		}
		if f.ns() != t.vr {
			// Out of sequence: an earlier frame was lost.  Ask
			// the BBS to resend from the one we're missing.
			if !t.rejSent || f.pf() {
				t.respond, t.respondF = ctlREJ, t.respondF || f.pf()
				t.rejSent = true
			}
			return
		}
		t.vr = (t.vr + 1) & 7
		t.rejSent = false
		t.deliver(f.info)
		// Acknowledge right away if the BBS polled for it; otherwise
		// when the channel goes quiet, or with our next I frame.
		t.respond, t.respondF = ctlRR, t.respondF || f.pf()
	case f.isS():
		defer t.resendIfPolledWithout()
		if !f.command && f.pf() && !t.pollAt.IsZero() {
			// The answer to our poll.  If we sent only one, it's a
			// measure of how long the BBS takes to answer.
			if !t.pollMulti {
				t.sampleRTT(time.Since(t.pollAt) - t.pollAir)
			}
			t.pollAt, t.pollMulti = time.Time{}, false
		}
		if !t.ackUpTo(f.nr()) {
			return
		}
		wasBusy := t.remoteBusy
		t.remoteBusy = f.sType() == ctlRNR
		if f.command && f.pf() {
			// The BBS is polling us.
			if t.respond == 0 {
				t.respond = ctlRR
			}
			t.respondF = true
		}
		if !f.command && f.pf() && t.recovery {
			// The BBS answered our poll; resend whatever it
			// hasn't acknowledged.  (This doesn't reset the retry
			// count: only an acknowledgment of new data does that,
			// in ackUpTo.  Otherwise, on a path where our short
			// polls get through but our long I frames don't, we'd
			// poll and resend forever.)
			t.recovery, t.nsent, t.t1 = false, 0, time.Time{}
		} else if f.sType() == ctlREJ {
			// Resend from N(R).  During timer recovery, that waits
			// for the answer to our poll, and T1 keeps running in
			// case the answer is lost.
			t.nsent = 0
			if !t.recovery {
				t.t1 = time.Time{}
			}
		} else if wasBusy && !t.remoteBusy && !t.recovery && t.nsent == 0 {
			// The BBS is no longer busy.  T1 was running only to
			// poll it; transmit restarts it for what it sends.
			t.t1 = time.Time{}
		}
	default:
		switch f.uType() {
		case ctlDISC:
			t.sendU(ctlUA|f.control&pfBit, false)
			// The BBS disconnected: cleanly, unless we still had data
			// it hadn't acknowledged (which the TNC may still hold).
			t.cleanEnd.Store(len(t.outstanding) == 0 && len(t.queue) == 0)
			t.linkEnded = true
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		case ctlDM:
			t.linkEnded = true
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		case ctlSABM:
			// The BBS reset the link.  Anything outstanding is lost.
			var lost int
			for _, info := range t.outstanding {
				lost += len(info)
			}
			slog.Warn("AX.25 link reset by BBS", "lost", lost)
			if lost != 0 {
				t.monitorLine('R', fmt.Sprintf("*** link reset by BBS; %d bytes of unacknowledged data lost", lost))
			}
			t.sendU(ctlUA|f.control&pfBit, false)
			t.resetLink()
		case ctlFRMR:
			slog.Error("AX.25 frame reject", "info", fmt.Sprintf("% X", f.info))
			t.sendU(ctlDISC|pfBit, true)
			t.setState(stateDisconnected, jnos.ErrDisconnected)
		}
	}
}

// polledWithoutGap is the shortest time since we last sent data before a poll
// from the BBS that doesn't acknowledge it counts as evidence that it was lost
// (see resendIfPolledWithout), when the BBS's answer time hasn't been measured.
// It is a variable so that tests can change it.
var polledWithoutGap = 3 * time.Second

// resendIfPolledWithout is called after the N(R) of a frame from the BBS has
// been processed.  If the frame is a poll (a command with P), and it doesn't
// acknowledge data we sent long enough ago for the BBS to have answered it
// (its typical answer time, plus the data's airtime), then the BBS is telling
// us it never got that data: it polls only when its own timer has run out
// waiting to hear from us.  So we resend the data now, right after our answer
// to the poll and in the same turn at the channel, as QtTermTCP does, rather
// than waiting for our own T1 to expire.  (W3XSC polls this way, resending its
// last frame, whenever it misses our acknowledgment; waiting for T1 cost
// 6-13 seconds each time.)
func (t *Transport) resendIfPolledWithout() {
	if t.lastFrame == nil || !t.lastFrame.command || !t.lastFrame.pf() {
		return
	}
	if len(t.outstanding) == 0 || t.nsent == 0 || t.remoteBusy || t.state != stateConnected {
		return
	}
	gap := polledWithoutGap
	if len(t.rtts) >= rttMinSamples {
		gap, _ = rttStats(t.rtts)
	}
	var size int
	for _, info := range t.outstanding[:t.nsent] {
		size += len(info)
	}
	if time.Since(t.lastDataSent) < gap+airtime(size) {
		return // the BBS may not have had time to answer it yet
	}
	// The resend happens in resendPolled, once our answer to the poll has
	// gone and the turnaround hold has passed.
	t.pendingResend = true
	t.note(fmt.Sprintf("the BBS's poll shows it missed our data (N(R)=%d); resending it after answering", t.lastFrame.nr()))
}

// resendPolled carries out a resend requested by resendIfPolledWithout.  In
// timer recovery (we've polled the BBS and await its answer), the resend is a
// poll itself, as on a T1 expiry; otherwise transmit resends the
// unacknowledged frames.  Either way, T1 restarts from the resend.
func (t *Transport) resendPolled() {
	if !t.pendingResend {
		return
	}
	if t.state != stateConnected || len(t.outstanding) == 0 || t.remoteBusy {
		t.pendingResend = false
		return
	}
	if t.held() || t.owesFinal() {
		return // nextDeadline wakes us when the hold ends
	}
	t.pendingResend = false
	if !t.recovery {
		t.nsent, t.t1 = 0, time.Time{} // transmit resends and restarts T1
		return
	}
	info := t.outstanding[0]
	if t.timing && t.timedNS == t.va {
		t.timing = false // its acknowledgment is now ambiguous
	}
	t.sendFrame(&frame{command: true, control: t.vr<<5 | pfBit | t.va<<1, info: info})
	t.lastDataSent = time.Now()
	t.nsent = max(t.nsent, 1)
	if t.respond == ctlRR && !t.respondF {
		t.respond = 0 // the I frame acknowledges what we've received
	}
	t.t1 = time.Now().Add(t.frack() + airtime(len(info)))
}

// ackUpTo processes an N(R) received from the BBS, which acknowledges all of
// our frames before it.  It returns false if N(R) is invalid.  Every
// outstanding frame has been sent at least once (even if a REJ or a poll
// response has us resending from an earlier one), so N(R) can acknowledge any
// of them.
func (t *Transport) ackUpTo(nr byte) bool {
	count := int((nr - t.va) & 7)
	if count > len(t.outstanding) {
		slog.Warn("AX.25 invalid N(R)", "nr", nr, "va", t.va, "outstanding", len(t.outstanding))
		return false
	}
	if count == 0 {
		return true
	}
	if t.timing && int((t.timedNS-t.va)&7) < count {
		t.timing = false
		t.sampleRTT(time.Since(t.timedAt) - t.timedAir)
	}
	t.outstanding = t.outstanding[count:]
	t.nsent = max(t.nsent-count, 0)
	t.va = nr
	t.retries = 0
	t.progress()
	t.t1 = time.Time{}
	if len(t.outstanding) == 0 && t.recovery {
		// Everything we sent has been acknowledged, so there's
		// nothing left to recover: no need to poll the BBS (an extra
		// transmission on a busy channel) to finish timer recovery.
		t.recovery, t.pollAt, t.pollMulti = false, time.Time{}, false
	}
	if t.nsent != 0 || t.recovery {
		t.startT1()
	}
	return true
}

// timeout handles the expiration of the T1 timer.  (The quiet timer for a
// pending acknowledgment is handled by sendResponse.)
func (t *Transport) timeout() {
	now := time.Now()
	if t.t1.IsZero() || now.Before(t.t1) {
		return
	}
	if t.state == stateConnected && t.held() {
		// Don't poll while the BBS may still be turning around; wait
		// (without using up a retry).
		t.t1 = t.holdUntil
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
		t.t1 = now.Add(t.t1Base())
	case stateConnected:
		if t.retries > t.dev.Retry {
			slog.Error("AX.25 link", "err", "BBS stopped responding")
			t.setState(stateDisconnected, jnos.ErrDisconnected)
			return
		}
		// Poll the BBS to find out what it has received.  If we have
		// data it hasn't acknowledged (and it isn't busy), the poll is
		// the oldest of that data, resent with the P bit, as QtTermTCP
		// does: on a lossy path, that puts the data itself in the
		// quiet slot where our frames are most likely to be heard, so
		// a single frame getting through makes progress.  Otherwise
		// it's an RR.
		t.recovery = true
		if len(t.outstanding) != 0 && !t.remoteBusy {
			info := t.outstanding[0]
			if t.timing && t.timedNS == t.va {
				t.timing = false // its acknowledgment is now ambiguous
			}
			t.sendFrame(&frame{command: true, control: t.vr<<5 | pfBit | t.va<<1, info: info})
			t.lastDataSent = time.Now()
			t.nsent = max(t.nsent, 1)
			if t.respond == ctlRR && !t.respondF {
				t.respond = 0 // the I frame acknowledges what we've received
			}
			t.t1 = now.Add(t.frack() + airtime(len(info)))
		} else {
			t.sendS(ctlRR, true, true)
			t.t1 = now.Add(t.frack())
		}
	case stateDisconnecting:
		if t.retries > maxDiscRetries {
			t.linkEnded = true // we've said goodbye enough
			t.setState(stateDisconnected, jnos.ErrDisconnected)
			return
		}
		t.sendU(ctlDISC|pfBit, true)
		t.t1 = now.Add(t.t1Base())
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
	t.lastProgress = time.Now()
	if len(t.pending) > start {
		t.writeLog(t.pending[start:])
	}
	t.mu.Unlock()
	t.signal()
}

// sendS sends a supervisory frame of the specified type, carrying our current
// N(R).  It acknowledges everything we've received, so it satisfies any owed
// response, except that only a response frame can carry an owed F bit.
func (t *Transport) sendS(stype byte, command, pf bool) {
	control := stype | t.vr<<5
	if pf {
		control |= pfBit
	}
	t.sendFrame(&frame{command: command, control: control})
	if !command || !t.respondF {
		t.respond, t.respondF = 0, false
	}
}

// sendU sends an unnumbered frame.
func (t *Transport) sendU(control byte, command bool) {
	t.sendFrame(&frame{command: command, control: control})
}

// sendFrame sends a frame to the BBS through the TNC (unless the link has
// already failed).
func (t *Transport) sendFrame(f *frame) {
	if t.state == stateDisconnected {
		return
	}
	f.dest, f.src = t.remote, t.local
	if trace {
		slog.Debug("AX.25 send", "control", fmt.Sprintf("%02X", f.control), "command", f.command, "len", len(f.info))
	}
	t.monitorFrame('T', f)
	t.writeTNC(kissEncode(f.encode()))
	if f.command && f.pf() && (t.state == stateConnected || t.state == stateConnecting) {
		// A poll: time the BBS's answer (see receiveConnected).
		t.pollMulti = !t.pollAt.IsZero()
		t.pollAt, t.pollAir = time.Now(), airtime(len(f.info))
	}
}

// writeTNC writes KISS-encoded data to the TNC, by way of the writer
// goroutine.  If the write fails, or takes longer than writeTimeout, the link
// fails.
func (t *Transport) writeTNC(data []byte) {
	if t.state == stateDisconnected {
		return // the link already failed
	}
	if trace {
		slog.Debug("KISS write", "bytes", fmt.Sprintf("% X", data))
	}
	timer := time.NewTimer(t.writeTimeout)
	defer timer.Stop()
	select {
	case t.writech <- data:
		select {
		case err := <-t.writeErr:
			if err != nil {
				t.portFailed = true
				if !t.isAborted() {
					slog.Error("write to TNC", "err", err)
				}
				t.setState(stateDisconnected, fmt.Errorf("write to TNC: %s", err))
			} else {
				t.lastSent.Store(time.Now().UnixNano())
			}
			return
		case <-timer.C:
		case <-t.abortch:
			t.setState(stateDisconnected, errAborted)
			return
		}
	case <-timer.C:
	case <-t.abortch:
		t.setState(stateDisconnected, errAborted)
		return
	}
	slog.Error("write to TNC", "err", "timeout")
	t.portFailed = true
	t.setState(stateDisconnected, errWriteTimeout)
}
