// Package serialtnc provides a transport layer for communicating with a JNOS
// BBS over RF, by way of a serial connection to a TNC.  The Open function, if
// successful, returns a Transport that can be passed to jnos.Connect.
package serialtnc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"go.bug.st/serial"

	"github.com/rothskeller/packet/jnos"
	"github.com/rothskeller/packet/jnos/tnc"
)

// echoTimeout is the amount of time to wait for an echo of data sent.
const echoTimeout = 200 * time.Millisecond

// tncTimeout is the amount of time to wait for response from the TNC
// (not involving RF traffic).
const tncTimeout = 500 * time.Millisecond

// rfTimeout is the amount of time to wait for a response from the BBS
// (requiring RF round trip).
const rfTimeout = time.Minute

// Constant byte slices.
var (
	cr   = []byte{'\r'}
	crlf = []byte{'\r', '\n'}
	lf   = []byte{'\n'}
)

// ErrBadEcho is returned when the echo of the sent data does not arrive within
// a reasonable amount of time.
var ErrBadEcho = errors.New("sent data not echoed correctly")

// Connect connects to the JNOS BBS at bbsAddress, by way of a TNC of the
// supplied type attached to serialPort, and returns an open jnos.Conn for
// interaction with it.  (bbsAddress should consist of a call sign, a dash, and
// a small integer SSID.)  It logs into the specified BBS mailbox.  If callsign
// is set and is different from mailbox, it self-identifies periodically using
// that call sign for FCC compliance.  (callsign should be the licensed FCC call
// sign of the calling user.)  If log is set, all traffic except echo-backs is
// logged to it.
func Connect(tnc *tnc.TNC, serialPort, bbsAddress, mailbox, callsign string, log io.Writer) (c *jnos.Conn, err error) {
	var t *Transport

	if callsign == mailbox {
		callsign = "" // no need to ident
	}
	if t, err = open(tnc, serialPort, bbsAddress, mailbox, log); err != nil {
		return nil, err
	}
	if c, err = jnos.Connect(t, callsign); err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return c, nil
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
func open(tnc *tnc.TNC, serialPort, bbsAddress, mailbox string, log io.Writer) (t *Transport, err error) {
	t = &Transport{tnc: tnc, log: log}
	if t.serial, err = serial.Open(serialPort, &serial.Mode{}); err != nil {
		slog.Error("serial.Open", "port", serialPort, "err", err)
		return nil, fmt.Errorf("serial.Open: %s", err)
	}
	// Send a newline and check that we get a command prompt.  We may get
	// other stuff ahead of it, and may get more than one command prompt.
	// Empirically, the TNC doesn't always respond to the first CR with a
	// prompt; sometimes it takes more than one.  We'll give it three tries
	// before giving up.
	var tries int
	for tries = 0; tries < 3; tries++ {
		if err = t.sendRaw([]byte{'\r'}); err != nil {
			err = fmt.Errorf("send initial newline: %s", err)
			goto TNCERROR
		}
		if _, err = t.readUntil(tnc.CommandPrompt, tncTimeout); err == nil {
			break
		}
	}
	if tries >= 3 {
		err = fmt.Errorf("read initial prompt (after 3 tries): %s", err)
		goto TNCERROR
	}
	// Apply the pre-connect settings.
	for _, c := range tnc.PreConnectCommands {
		if err = t.send(c); err != nil {
			goto TNCERROR
		}
		if _, err = t.readUntil(tnc.CommandPrompt, tncTimeout); err != nil {
			goto TNCERROR
		}
	}
	// Set the mailbox we want to connect to as our "call sign".
	if err = t.send(fmt.Sprintf("%s %s\n", tnc.MyCallCommand, mailbox)); err != nil {
		goto TNCERROR
	}
	if _, err = t.readUntil(tnc.CommandPrompt, tncTimeout); err != nil {
		goto TNCERROR
	}
	// Connect to the BBS.
	t.wasConnected = true
	if err = t.send(fmt.Sprintf("%s %s\n", tnc.ConnectCommand, bbsAddress)); err != nil {
		goto BBSERROR
	}
	if _, err = t.readUntil(tnc.ConnectedMessage, rfTimeout); err != nil {
		goto BBSERROR
	}
	t.connected = true
	return t, nil
TNCERROR:
	t.Close()
	return nil, fmt.Errorf("TNC setup: %s", err)
BBSERROR:
	t.Close()
	return nil, fmt.Errorf("BBS connect: %s", err)
}

// Transport is the KPC-3 Plus transport to the JNOS BBS.
type Transport struct {
	tnc          *tnc.TNC
	serial       serial.Port
	readbuf      []byte
	pending      []byte
	connected    bool
	wasConnected bool
	log          io.Writer
}

// ReadUntil reads data from the BBS until the specified string is seen, or a
// timeout occurs.  It returns the data that was read (even if it returns an
// error).
func (t *Transport) ReadUntil(until string) (s string, err error) {
	if !t.connected {
		return "", jnos.ErrDisconnected
	}
	return t.readUntil(until, rfTimeout)
}

// ReadUntilT reads data from the BBS until the specified string is seen, or the
// specified timeout occurs.  It returns the data that was read (even if it
// returns an error).
func (t *Transport) ReadUntilT(until string, timeout time.Duration) (s string, err error) {
	if !t.connected {
		return "", jnos.ErrDisconnected
	}
	return t.readUntil(until, timeout)
}

// readUntil reads data from the TNC until the specified string is seen, or the
// specified timeout occurs.  It returns the data that was read (even if it
// returns an error).
func (t *Transport) readUntil(until string, timeout time.Duration) (data string, err error) {
	var untilb = bytes.ReplaceAll([]byte(until), lf, crlf)
	if t.readbuf == nil {
		t.readbuf = make([]byte, 1024)
	}
	for err == nil {
		var count int

		if data, err = t.checkDisconnected(); err != nil {
			return data, err
		}
		if data = t.checkPending(untilb); data != "" {
			return data, nil
		}
		t.serial.SetReadTimeout(timeout)
		count, err = t.serial.Read(t.readbuf)
		if count != 0 {
			if t.log != nil {
				t.log.Write(t.readbuf[:count])
			}
			t.pending = append(t.pending, t.readbuf[:count]...)
		} else if err == nil {
			slog.Error("serial.Read", "err", "timeout")
			err = jnos.ErrTimeout
		}
	}
	data = string(bytes.ReplaceAll(t.pending, crlf, lf))
	t.pending = t.pending[:0]
	return data, err
}

// checkDisconnected checks whether the pending input buffer includes a message
// from the TNC indicating that its connection to the BBS was lost.  If so, it
// returns a disconnected error, along with everything in the pending buffer
// prior to that message.
func (t *Transport) checkDisconnected() (data string, err error) {
	if !t.connected {
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

// checkPending checks whether the specified string is in the pending buffer.
// If so, it extracts everything up through and including that string from the
// pending buffer and returns it.  Otherwise, it returns an empty string.
func (t *Transport) checkPending(until []byte) (data string) {
	// Look for the desired string in our pending buffer.  If it's
	// there, return everything up through and including it.
	if idx := bytes.Index(t.pending, until); idx >= 0 {
		idx += len(until)
		data = string(bytes.ReplaceAll(t.pending[:idx], crlf, lf))
		t.pending = t.pending[idx:]
		return data
	}
	return ""
}

// Send sends a string to the BBS.
func (t *Transport) Send(s string) (err error) {
	if !t.connected {
		return jnos.ErrDisconnected
	}
	return t.send(s)
}

// send sends a string to the TNC.  The method returns when the string has been
// echoed by the TNC, or after a timeout or other error.
func (t *Transport) send(data string) (err error) {
	var (
		tosend  []byte
		echo    []byte
		plogoff int
	)
	plogoff = len(t.pending) // pending bytes already logged
	tosend = bytes.ReplaceAll([]byte(data), lf, cr)
	if len(tosend) == 0 || tosend[len(tosend)-1] != '\r' {
		tosend = append(tosend, '\r')
	}
	if err = t.sendRaw(tosend); err != nil {
		return err
	}
	echo = bytes.ReplaceAll(tosend, cr, crlf)
	for err == nil {
		var count int

		if _, err = t.checkDisconnected(); err != nil {
			break
		}
		if idx := bytes.Index(t.pending, echo); idx >= 0 {
			t.pending = append(t.pending[:idx], t.pending[idx+len(echo):]...)
			if t.log != nil && len(t.pending) > plogoff {
				t.log.Write(t.pending[plogoff:])
			}
			return nil
		}
		t.serial.SetReadTimeout(echoTimeout)
		count, err = t.serial.Read(t.readbuf)
		if count != 0 {
			t.pending = append(t.pending, t.readbuf[:count]...)
		} else if err == nil {
			slog.Error("bad echo")
			err = ErrBadEcho
		}
	}
	return err
}

// sendRaw sends the specified data with no translation and no wait for echo.
func (t *Transport) sendRaw(data []byte) (err error) {
	var count int

	if t.log != nil {
		t.log.Write(bytes.ReplaceAll(data, cr, crlf))
	}
	for len(data) != 0 {
		if count, err = t.serial.Write(data); err != nil {
			slog.Error("serial.Write", "err", err)
			return err
		}
		data = data[count:]
	}
	return nil
}

// Close closes the connection to the BBS and resets the configuration of the
// TNC.
func (t *Transport) Close() (err error) {
	defer t.serial.Close()
	if t.connected {
		// We're still connected (or at least, we haven't received a
		// *** DISCONNECTED message).  Send a Ctrl-C to get back into
		// command mode.
		if err = t.sendRaw([]byte{3}); err != nil {
			slog.Error("unable to get back to TNC command mode for cleanup")
			return fmt.Errorf("unable to get back to TNC command mode for cleanup: %s", err)
		}
		if _, err = t.readUntil(t.tnc.CommandPrompt, tncTimeout); err != nil && err != jnos.ErrDisconnected {
			slog.Error("unable to get back to TNC command mode for cleanup")
			return fmt.Errorf("unable to get back to TNC command mode for cleanup: %s", err)
		}
		// Then tell the TNC to disconnect from the BBS, and wait until
		// we get the *** DISCONNECTED message.
		if err = t.send("D"); err != nil && err != jnos.ErrDisconnected {
			slog.Error("can't disconnect")
			return fmt.Errorf("cleanup: can't disconnect: %s", err)
		} else if err == nil {
			if _, err = t.readUntil(string(t.tnc.DisconnectedMessage), rfTimeout); err != jnos.ErrDisconnected {
				slog.Error("unexpected response", "exp", "disconnect", "act", err)
				return fmt.Errorf("ERROR: cleanup: expected ErrDisconnected, got %s", err)
			}
		}
	}
	// Apply all of the post-connect settings.
	t.readUntil(t.tnc.CommandPrompt, tncTimeout) // eat a prompt if there is one
	// Note: we used to MYCALL back to the FCC call sign at this point, but
	// some TNCs (Kenwood and Alinco embeddeds, primarily) had a problem
	// with doing it that soon (while they were still handling the
	// disconnect).  They'd send the disconnect ack with the changed call
	// sign!  Which would then leave JNOS thinking we're still connected,
	// and it would reject the next connection attempt (if right away) with
	// a "busy" message.  To avoid all that, we're leaving the call sign
	// unchanged.``
	for _, c := range t.tnc.PostDisconnectCommands {
		if err2 := t.send(c); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore TNC settings: %s", err2)
		}
		if _, err2 := t.readUntil(t.tnc.CommandPrompt, tncTimeout); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore TNC settings: %s", err2)
		}
	}
	return err
}

// UseVerboseReads returns whether it's appropriate to use verbose reads
// when communicating over this transport.
func (t *Transport) UseVerboseReads() bool { return false }
