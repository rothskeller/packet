// Package serialtnc provides a transport layer for communicating with a JNOS
// BBS over RF, by way of a serial connection to a TNC.  The Open function, if
// successful, returns a Transport that can be passed to jnos.Connect.
package serialtnc

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.bug.st/serial"

	"github.com/rothskeller/packet/v4/jnos"
	"github.com/rothskeller/packet/v4/jnos/tnc"
)

// echoTimeout is the default amount of time to wait for an echo of data sent,
// beyond the time needed to transfer the data over the serial line (see
// echoPerByte).
const echoTimeout = 200 * time.Millisecond

// echoPerByte is the additional time allowed for the echo of each byte sent.
// It is roughly twice the time needed to transfer a byte each way at 9600 baud.
const echoPerByte = 2 * time.Millisecond

// maxPending is the maximum amount of received data we will buffer while
// looking for an expected string.  Far more than this would be a sign that
// something is badly wrong.
const maxPending = 1 << 20

// sendBlockSize is the maximum size of data to send before waiting for echoes.
const sendBlockSize = 256

// tncTimeout is the default amount of time to wait for response from the TNC
// (not involving RF traffic).
const tncTimeout = 500 * time.Millisecond

// rfTimeout is the amount of time to wait for a response from the BBS
// (requiring RF round trip).  It is a variable only so that tests can shorten
// it.
var rfTimeout = time.Minute

// staleDisconnectWait is how long to wait for the TNC to finish disconnecting
// a stale connection (see setupTNC).  It is a variable only so that tests can
// shorten it.
var staleDisconnectWait = 10 * time.Second

// Constant byte slices.
var (
	cr   = []byte{'\r'}
	crlf = []byte{'\r', '\n'}
	lf   = []byte{'\n'}
)

// ErrBadEcho is returned when the echo of the sent data does not arrive within
// a reasonable amount of time.
var ErrBadEcho = errors.New("sent data not echoed correctly")

// errAborted is returned by all Transport methods after Abort is called.
var errAborted = errors.New("connection aborted")

// errOverflow is returned when too much data arrives without the expected
// string being seen.
var errOverflow = errors.New("too much unexpected data received from TNC")

// Connect connects to the JNOS BBS at bbsAddress, by way of a TNC of the
// supplied type attached to serialPort, and returns an open jnos.Conn for
// interaction with it.  (bbsAddress should consist of a call sign, a dash, and
// a small integer SSID.)  It logs into the specified BBS mailbox.  If callsign
// is set and is different from mailbox, it self-identifies periodically using
// that call sign for FCC compliance.  (callsign should be the licensed FCC call
// sign of the calling user.)  If log is set, all traffic except echo-backs is
// logged to it.
func Connect(tnc *tnc.TNC, serialPort, bbsAddress, mailbox, callsign string, log io.Writer) (c *jnos.Conn, err error) {
	return ConnectContext(context.Background(), tnc, serialPort, bbsAddress, mailbox, callsign, log)
}

// ConnectContext is Connect, except that canceling ctx before it returns
// abandons the connection attempt promptly (cleaning up the TNC as Close does
// after an Abort).  Once it returns, ctx no longer affects the connection.
func ConnectContext(ctx context.Context, tnc *tnc.TNC, serialPort, bbsAddress, mailbox, callsign string, log io.Writer) (c *jnos.Conn, err error) {
	var t *Transport

	if callsign == mailbox {
		callsign = "" // no need to ident
	}
	if t, err = open(tnc, serialPort, bbsAddress, mailbox, log, withContext(ctx)); err != nil {
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

// withContext is an option for start: canceling ctx aborts the Transport,
// until stopContext is called.
func withContext(ctx context.Context) func(*Transport) {
	return func(t *Transport) { t.stopCtx = context.AfterFunc(ctx, t.Abort) }
}

// stopContext stops the context passed to withContext from aborting the
// Transport.
func (t *Transport) stopContext() {
	if t.stopCtx != nil {
		t.stopCtx()
		t.stopCtx = nil
	}
}

// Open opens a transport to the JNOS BBS at bbsAddress, by way of a TNC of the
// supplied type attached to serialPort.  (bbsAddress should consist of a call
// sign, a dash, and a small integer SSID.)  It logs into the mailbox
// corresponding to the specified callsign, which must be the licensed FCC call
// sign of the calling user.  (For connecting to other mailboxes, see the
// Connect function.)  If log is set, all traffic except echo-backs is logged to
// it.
func Open(tnc *tnc.TNC, serialPort, bbsAddress string, log io.Writer) (t *Transport, err error) {
	return open(tnc, serialPort, bbsAddress, "", log)
}

// open is the common code between Connect and Open.
func open(tnc *tnc.TNC, serialPort, bbsAddress, mailbox string, log io.Writer, opts ...func(*Transport)) (t *Transport, err error) {
	var sp serial.Port

	if err = tnc.Validate(); err != nil {
		return nil, fmt.Errorf("TNC setup: invalid TNC definition: %s", err)
	}
	if sp, err = serial.Open(serialPort, &serial.Mode{}); err != nil {
		slog.Error("serial.Open", "port", serialPort, "err", err)
		return nil, fmt.Errorf("serial.Open: %s", err)
	}
	// After an Abort (which closes the port), Close reopens it to put the
	// TNC back in order.
	reopen := func(t *Transport) {
		t.reopen = func() (port, error) { return serial.Open(serialPort, &serial.Mode{}) }
	}
	return start(tnc, sp, bbsAddress, mailbox, log, append([]func(*Transport){reopen}, opts...)...)
}

// withReopen is an option for start: after an Abort, Close uses reopen to
// reopen the port and clean up the TNC.
func withReopen(reopen func() (port, error)) func(*Transport) {
	return func(t *Transport) { t.reopen = reopen }
}

// start sets up the TNC attached to the (open) port p, and connects to the
// BBS.  The port is closed if it fails.  opts are applied to the Transport
// first.
func start(def *tnc.TNC, p port, bbsAddress, mailbox string, log io.Writer, opts ...func(*Transport)) (t *Transport, err error) {
	t = newTransport(def, p, log)
	for _, opt := range opts {
		opt(t)
	}
	if err = t.setupTNC(mailbox); err != nil {
		t.stopContext()
		t.Close()
		return nil, fmt.Errorf("TNC setup: %s", err)
	}
	if err = t.connect(bbsAddress); err != nil {
		t.stopContext()
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return t, nil
}

// newTransport returns a new Transport for the TNC described by def, attached
// to the port p.
func newTransport(def *tnc.TNC, p port, log io.Writer) *Transport {
	t := &Transport{tnc: def, serial: p, log: log, readbuf: make([]byte, 1024)}
	// Zero (or invalid negative) timeouts mean "use the default".
	if def.EchoTimeout > 0 {
		t.echoTimeout = time.Duration(def.EchoTimeout) * time.Millisecond
	} else {
		t.echoTimeout = echoTimeout
	}
	if def.ReplyTimeout > 0 {
		t.replyTimeout = time.Duration(def.ReplyTimeout) * time.Millisecond
	} else {
		t.replyTimeout = tncTimeout
	}
	return t
}

// setupTNC gets a command prompt from the TNC, applies the pre-connect
// settings, and sets the TNC call sign to the mailbox we want to connect to.
func (t *Transport) setupTNC(mailbox string) (err error) {
	// Send a newline and check that we get a command prompt.  We may get
	// other stuff ahead of it, and may get more than one command prompt.
	// Empirically, the TNC doesn't always respond to the first CR with a
	// prompt; sometimes it takes more than one.  We'll give it three tries
	// before giving up.
	//
	// If the TNC doesn't answer the first newline with a prompt, it may
	// still be in converse mode from a session that didn't end cleanly
	// (the program was killed, say), so the later tries start with a
	// Ctrl-C to get back to command mode.  If that's what it took, the
	// TNC may also still be connected, so we tell it to disconnect.
	var tries int
	for tries = 0; tries < 3; tries++ {
		var send = []byte{'\r'}
		if tries != 0 {
			send = []byte{3, '\r'}
		}
		if err = t.sendRaw(send); err != nil {
			return fmt.Errorf("send initial newline: %s", err)
		}
		if _, err = t.readUntil(t.tnc.CommandPrompt, t.replyTimeout); err == nil {
			break
		} else if err == errAborted {
			return err
		}
	}
	if tries >= 3 {
		return fmt.Errorf("read initial prompt (after 3 tries): %s", err)
	}
	if tries != 0 {
		slog.Warn("TNC was not in command mode; disconnecting any stale connection")
		if err = t.send("D"); err != nil {
			return err
		}
		// If the TNC was connected, it reports the disconnection once
		// it's finished over the air, which must happen before we try
		// to connect: otherwise we'd take that report as the failure of
		// our connection.  If it wasn't connected, it doesn't report
		// one, so this wait is bounded.
		if data, _, err := t.readUntilAny(staleDisconnectWait, t.tnc.DisconnectedMessage); err == errAborted {
			return err
		} else if err == nil {
			slog.Info("stale connection disconnected", "tnc", data)
		}
		t.readUntil(t.tnc.CommandPrompt, t.replyTimeout) // eat a prompt if there is one
	}
	// Apply the pre-connect settings.
	for _, c := range t.tnc.PreConnectCommands {
		if err = t.send(c); err != nil {
			return err
		}
		if _, err = t.readUntil(t.tnc.CommandPrompt, t.replyTimeout); err != nil {
			return err
		}
	}
	// Set the mailbox we want to connect to as our "call sign".
	if err = t.send(fmt.Sprintf("%s %s\n", t.tnc.MyCallCommand, mailbox)); err != nil {
		return err
	}
	if _, err = t.readUntil(t.tnc.CommandPrompt, t.replyTimeout); err != nil {
		return err
	}
	return nil
}

// connect tells the TNC to connect to the BBS, and waits for the connection to
// be established.
func (t *Transport) connect(bbsAddress string) (err error) {
	var (
		data  string
		which int
	)
	// From here on, until we know the connect attempt has ended, Close
	// needs to tell the TNC to disconnect.
	t.connecting = true
	if err = t.send(fmt.Sprintf("%s %s\n", t.tnc.ConnectCommand, bbsAddress)); err != nil {
		return err
	}
	// Wait for the connection, but if the TNC reports that the connection
	// attempt failed (e.g. "*** retry count exceeded" followed by "***
	// DISCONNECTED"), fail right away.
	if data, which, err = t.readUntilAny(rfTimeout, t.tnc.ConnectedMessage, t.tnc.DisconnectedMessage); err != nil {
		return err
	}
	t.connecting = false
	if which != 0 {
		slog.Error("connect failed", "bbs", bbsAddress, "tnc", data)
		return fmt.Errorf("connection to %s failed: %s", bbsAddress, summarize(data))
	}
	t.connected = true
	return nil
}

// summarize returns a one-line summary of the TNC output in data, for use in
// an error message.
func summarize(data string) string {
	var lines []string
	for line := range strings.Lines(data) {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) > 3 {
		lines = lines[len(lines)-3:]
	}
	s := strings.Join(lines, "; ")
	if len(s) > 200 {
		s = s[len(s)-200:]
	}
	return s
}

// port is the subset of serial.Port used by the Transport.  (It is an
// interface so that tests can supply a fake port.)
type port interface {
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	SetReadTimeout(t time.Duration) error
	Close() error
}

// Transport is the KPC-3 Plus transport to the JNOS BBS.
type Transport struct {
	tnc          *tnc.TNC
	serial       port
	readbuf      []byte
	pending      []byte
	connected    bool
	connecting   bool // connect command sent, outcome not yet known
	log          io.Writer
	echoTimeout  time.Duration
	replyTimeout time.Duration
	sendTimeout  time.Duration
	aborted      atomic.Bool
	closeOnce    sync.Once
	reopen       func() (port, error) // reopens the port for cleanup after Abort, or nil
	stopCtx      func() bool          // stops a context from aborting us (see withContext), or nil
}

// ReadUntil reads data from the BBS until the specified string is seen, or a
// timeout occurs.  It returns the data that was read (even if it returns an
// error).
func (t *Transport) ReadUntil(until string) (s string, err error) {
	if t.aborted.Load() {
		return "", errAborted
	}
	if !t.connected {
		return "", jnos.ErrDisconnected
	}
	// For a timeout, use the rough amount of time needed to send the data
	// in the preceding Send call, if any, plus our usual overhead timeout.
	timeout := t.sendTimeout + rfTimeout
	t.sendTimeout = 0
	return t.readUntil(until, timeout)
}

// readUntil reads data from the TNC until the specified string is seen, or the
// specified timeout occurs.  It returns the data that was read (even if it
// returns an error).
func (t *Transport) readUntil(until string, timeout time.Duration) (data string, err error) {
	data, _, err = t.readUntilAny(timeout, until)
	return data, err
}

// readUntilAny reads data from the TNC until any of the specified strings is
// seen, or the specified timeout (measured from the start of the call)
// expires.  It returns the data that was read (even if it returns an error),
// and on success, the index of the string that was seen.
func (t *Transport) readUntilAny(timeout time.Duration, untils ...string) (data string, which int, err error) {
	var targets = make([][]byte, len(untils))
	for i, until := range untils {
		targets[i] = bytes.ReplaceAll([]byte(until), lf, crlf)
	}
	deadline := time.Now().Add(timeout)
	for {
		if t.aborted.Load() {
			err = errAborted
			break
		}
		if data, err = t.checkDisconnected(); err != nil {
			return data, -1, err
		}
		if data, which = t.checkPending(targets); which >= 0 {
			return data, which, nil
		}
		if err = t.readMore(deadline, true); err != nil {
			break
		}
	}
	if err == jnos.ErrTimeout {
		slog.Error("serial.Read", "err", "timeout")
	}
	data = string(bytes.ReplaceAll(t.pending, crlf, lf))
	t.pending = t.pending[:0]
	return data, -1, err
}

// readMore reads whatever data is available from the TNC (waiting for some,
// but not past the deadline) and adds it to the pending buffer, logging it if
// log is true.  It returns jnos.ErrTimeout if the deadline has passed.
func (t *Transport) readMore(deadline time.Time, log bool) (err error) {
	var (
		count     int
		remaining time.Duration
	)
	if remaining = time.Until(deadline); remaining <= 0 {
		return jnos.ErrTimeout
	}
	remaining = max(remaining, time.Millisecond)
	if err = t.serial.SetReadTimeout(remaining); err != nil {
		if t.aborted.Load() {
			return errAborted
		}
		slog.Error("serial.SetReadTimeout", "timeout", remaining, "err", err)
		return fmt.Errorf("set read timeout: %s", err)
	}
	count, err = t.serial.Read(t.readbuf)
	if t.aborted.Load() {
		return errAborted
	}
	if count > 0 {
		if log && t.log != nil {
			t.log.Write(t.readbuf[:count])
		}
		t.pending = append(t.pending, t.readbuf[:count]...)
		if len(t.pending) > maxPending {
			slog.Error("serial.Read", "err", errOverflow)
			return errOverflow
		}
	}
	if err != nil {
		slog.Error("serial.Read", "err", err)
		return err
	}
	// A zero-length read with no error is a timeout, but our caller will
	// notice that on the next call, if the deadline really has passed.
	return nil
}

// checkDisconnected checks whether the pending input buffer includes a message
// from the TNC indicating that its connection to the BBS was lost.  If so, it
// returns a disconnected error, along with everything in the pending buffer
// prior to that message.
func (t *Transport) checkDisconnected() (data string, err error) {
	if !t.connected || t.tnc.DisconnectedMessage == "" {
		return "", nil
	}
	if idx := bytes.Index(t.pending, []byte(t.tnc.DisconnectedMessage)); idx >= 0 {
		data = string(bytes.ReplaceAll(t.pending[:idx], crlf, lf))
		t.pending = t.pending[idx+len(t.tnc.DisconnectedMessage):]
		t.connected = false
		slog.Error("disconnected")
		return data, jnos.ErrDisconnected
	}
	return "", nil
}

// checkPending checks whether any of the specified strings is in the pending
// buffer.  If so, it extracts everything up through and including the first
// one found from the pending buffer and returns it, along with the index of
// the string found.  Otherwise, it returns an empty string and -1.  Empty
// strings never match.
func (t *Transport) checkPending(untils [][]byte) (data string, which int) {
	var start, end int

	which = -1
	for i, until := range untils {
		if len(until) == 0 {
			continue
		}
		if idx := bytes.Index(t.pending, until); idx >= 0 && (which < 0 || idx < start) {
			which, start, end = i, idx, idx+len(until)
		}
	}
	if which < 0 {
		return "", -1
	}
	data = string(bytes.ReplaceAll(t.pending[:end], crlf, lf))
	t.pending = t.pending[end:]
	return data, which
}

// Send sends a string to the BBS.
func (t *Transport) Send(s string) (err error) {
	if t.aborted.Load() {
		return errAborted
	}
	if !t.connected {
		return jnos.ErrDisconnected
	}
	return t.send(s)
}

// send sends a string to the TNC.  The method returns when the string has been
// echoed by the TNC, or after a timeout or other error.
func (t *Transport) send(data string) (err error) {
	var (
		tosend []byte
	)
	tosend = bytes.ReplaceAll([]byte(data), lf, cr)
	if len(tosend) == 0 || tosend[len(tosend)-1] != '\r' {
		tosend = append(tosend, '\r')
	}
	for len(tosend) != 0 {
		if len(tosend) > sendBlockSize {
			err = t.sendBlock(tosend[:sendBlockSize])
			tosend = tosend[sendBlockSize:]
		} else {
			err = t.sendBlock(tosend)
			tosend = nil
		}
		if err != nil {
			return err
		}
	}
	// The next ReadUntil should have an extra long timeout based on the
	// amount of time it takes to send this data.  We compute this very
	// roughly as 2½ seconds per packet.
	t.sendTimeout = 2500 * time.Millisecond * time.Duration(len(data)/128)
	return nil
}

// sendBlock sends a block of data to the TNC.  The method returns when the
// block has been echoed by the TNC, or after a timeout or other error.
func (t *Transport) sendBlock(tosend []byte) (err error) {
	var (
		echo     []byte
		plogoff  int
		deadline time.Time
	)
	plogoff = len(t.pending) // pending bytes already logged
	if err = t.sendRaw(tosend); err != nil {
		return err
	}
	echo = bytes.ReplaceAll(tosend, cr, crlf)
	// The echo must arrive within the echo timeout plus the time needed
	// to transfer the data, regardless of what else arrives meanwhile.
	deadline = time.Now().Add(t.echoTimeout + time.Duration(len(tosend))*echoPerByte)
	for {
		if t.aborted.Load() {
			return errAborted
		}
		if _, err = t.checkDisconnected(); err != nil {
			return err
		}
		if idx := bytes.Index(t.pending, echo); idx >= 0 {
			t.pending = append(t.pending[:idx], t.pending[idx+len(echo):]...)
			if t.log != nil && len(t.pending) > plogoff {
				t.log.Write(t.pending[plogoff:])
			}
			return nil
		}
		if err = t.readMore(deadline, false); err == jnos.ErrTimeout {
			slog.Error("bad echo", "exp", hex.EncodeToString(echo), "act", hex.EncodeToString(t.pending))
			return ErrBadEcho
		} else if err != nil {
			return err
		}
	}
}

// sendRaw sends the specified data with no translation and no wait for echo.
func (t *Transport) sendRaw(data []byte) (err error) {
	var count int

	if t.log != nil {
		t.log.Write(bytes.ReplaceAll(data, cr, crlf))
	}
	for len(data) != 0 {
		if t.aborted.Load() {
			return errAborted
		}
		if count, err = t.serial.Write(data); err != nil {
			if t.aborted.Load() {
				return errAborted
			}
			slog.Error("serial.Write", "err", err)
			return err
		}
		data = data[count:]
	}
	return nil
}

var _ jnos.Aborter = (*Transport)(nil)

// Abort ends the connection immediately.  It closes the serial port, so that
// any Transport call in progress on another goroutine returns promptly with an
// error; any later call (other than Close) also returns an error.  It is safe
// to call at any time, more than once, and concurrently with the other
// methods.  Close must still be called afterward, once the call in progress
// (if any) has returned: it reopens the port and makes a quick, best-effort
// cleanup of the TNC (see cleanupAfterAbort).
func (t *Transport) Abort() {
	t.aborted.Store(true)
	t.closePort()
}

// closePort closes the serial port, if it hasn't already been closed.
func (t *Transport) closePort() {
	t.closeOnce.Do(func() { t.serial.Close() })
}

// Close closes the connection to the BBS and resets the configuration of the
// TNC.
func (t *Transport) Close() (err error) {
	defer t.closePort()
	if t.aborted.Load() {
		return t.cleanupAfterAbort()
	}
	// Whatever goes wrong below, we carry on with the rest of the cleanup,
	// so that the TNC is left in as good a state as we can manage; the
	// first error is the one returned.
	if t.connected {
		// We're still connected (or at least, we haven't received a
		// *** DISCONNECTED message).  Send a Ctrl-C to get back into
		// command mode.  Give it a second try if there's no prompt.
		var gotPrompt bool
		for range 2 {
			if err2 := t.sendRaw([]byte{3}); err2 != nil {
				break
			}
			if _, err2 := t.readUntil(t.tnc.CommandPrompt, t.replyTimeout); err2 == nil || err2 == jnos.ErrDisconnected {
				gotPrompt = true
				break
			}
		}
		if !gotPrompt {
			slog.Error("unable to get back to TNC command mode for cleanup")
			err = errors.New("cleanup: unable to get back to TNC command mode")
		}
		// Then tell the TNC to disconnect from the BBS, and wait until
		// we get the *** DISCONNECTED message.
		if !t.connected {
			// We got it while getting back to command mode.
		} else if err2 := t.send("D"); err2 == jnos.ErrDisconnected {
			// We got it while sending D.
		} else if err2 != nil {
			slog.Error("can't disconnect", "err", err2)
			if err == nil {
				err = fmt.Errorf("cleanup: can't disconnect: %s", err2)
			}
		} else if _, err2 = t.readUntil(string(t.tnc.DisconnectedMessage), rfTimeout); err2 != jnos.ErrDisconnected {
			slog.Error("unexpected response", "exp", "disconnect", "act", err2)
			if err == nil {
				err = fmt.Errorf("cleanup: expected ErrDisconnected, got %s", err2)
			}
		}
		t.connected = false
	} else if t.connecting {
		// We sent a connect command, but never learned whether it
		// succeeded.  The TNC may still be retrying the connection, or
		// may have connected after we gave up on it and be sitting in
		// converse mode.  Get back to command mode (in case it's in
		// converse mode) and cancel the connection.  This is best
		// effort; we don't wait for the disconnection to complete.
		t.connecting = false
		if err2 := t.sendRaw([]byte{3}); err2 == nil {
			t.readUntil(t.tnc.CommandPrompt, t.replyTimeout) // may or may not get a prompt
			if err2 = t.send("D"); err2 != nil {
				slog.Error("cleanup: can't cancel connection", "err", err2)
			}
		}
	}
	// Apply all of the post-connect settings.
	t.readUntil(t.tnc.CommandPrompt, t.replyTimeout) // eat a prompt if there is one
	// Note: we used to MYCALL back to the FCC call sign at this point, but
	// some TNCs (Kenwood and Alinco embeddeds, primarily) had a problem
	// with doing it that soon (while they were still handling the
	// disconnect).  They'd send the disconnect ack with the changed call
	// sign!  Which would then leave JNOS thinking we're still connected,
	// and it would reject the next connection attempt (if right away) with
	// a "busy" message.  To avoid all that, we're leaving the call sign
	// unchanged.
	for _, c := range t.tnc.PostDisconnectCommands {
		if err2 := t.send(c); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore TNC settings: %s", err2)
		}
		if _, err2 := t.readUntil(t.tnc.CommandPrompt, t.replyTimeout); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore TNC settings: %s", err2)
		}
	}
	return err
}

// cleanupAfterAbort is Close after an Abort.  Abort closed the port, perhaps in
// the middle of a dialog with the TNC, which may be left in converse mode,
// connected to the BBS, with our settings in place.  So (if it can) this
// reopens the port, and with a fresh Transport, gets the TNC back to command
// mode, tells it to disconnect (without waiting for the disconnection to
// finish, since the user wants things to stop), and restores its settings.
// It's best effort: the result is logged, not returned.
func (t *Transport) cleanupAfterAbort() error {
	t.closePort()
	if t.reopen == nil {
		return nil
	}
	p, err := t.reopen()
	if err != nil {
		slog.Warn("TNC cleanup after abort: can't reopen the serial port", "err", err)
		return nil
	}
	clean := newTransport(t.tnc, p, t.log)
	// If we might be connected, Close sends Ctrl-C and D (without waiting
	// for the disconnection).  Either way, it restores the settings.
	clean.connecting = t.connected || t.connecting
	if err = clean.Close(); err != nil {
		slog.Warn("TNC cleanup after abort", "err", err)
	} else {
		slog.Info("TNC cleanup after abort: done")
	}
	return nil
}

// UseVerboseReads returns whether it's appropriate to use verbose reads
// when communicating over this transport.
func (t *Transport) UseVerboseReads() bool { return false }
