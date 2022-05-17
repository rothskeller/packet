// Package kpc3plus provides a transport layer for communicating with a JNOS BBS
// over RF, by way of a serial connection to a Kantronics KPC 3 Plus TNC.  The
// Open function, if successful, returns a Transport that can be passed to
// jnos.Connect.
package kpc3plus

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/jacobsa/go-serial/serial"

	"steve.rothskeller.net/packet/jnos"
)

// trace is a flag indicating that data sent and retrieved should be echoed to
// stdout.  It is normally false, but can be set to true for debugging.
const trace = false

// echoTimeout is the amount of time to wait for an echo of data sent.
const echoTimeout = 200 * time.Millisecond

// tncTimeout is the amount of time to wait for response from the TNC
// (not involving RF traffic).
const tncTimeout = 500 * time.Millisecond

// rfTimeout is the amount of time to wait for a response from the BBS
// (requiring RF round trip).
const rfTimeout = time.Minute

// commandPrompt is the TNC command prompt that we wait for between commands.
const commandPrompt = "cmd:"

// Constant byte slices.
var (
	cr                  = []byte{'\r'}
	crlf                = []byte{'\r', '\n'}
	lf                  = []byte{'\n'}
	disconnectedMessage = []byte("*** DISCONNECTED\r\n")
)

// ErrBadEcho is returned when the echo of the sent data does not arrive within
// a reasonable amount of time.
var ErrBadEcho = errors.New("sent data not echoed correctly")

// preConnectCommands are commands to be sent to the TNC prior to connecting to
// the BBS.
var preConnectCommands = []string{
	"INTFACE TERMINAL",
	"CD SOFTWARE",
	"NEWMODE ON",
	"8BITCONV ON",
	"BEACON EVERY 0",
	"SLOTTIME 10",
	"PERSIST 63",
	"PACLEN 128",
	"MAXFRAME 2",
	"FRACK 6",
	"RETRY 8",
	"CHECK 30",
	"TXDELAY 40",
	"XFLOW OFF",
	"SENDPAC $05",
	"CR OFF",
	"PACTIME AFTER 2",
	"CPACTIME ON",
	"STREAMEV OFF",
	"STREAMSW $00",
	"UNPROTO IDENT",
	"MONITOR ON",
}

// postDisconnectCommands are commands to be sent to the TNC after we disconnect
// from the BBS.
var postDisconnectCommands = []string{
	"SENDPAC $0D",
	"CR ON",
	"PACTIME AFTER 10",
	"CPACTIME OFF",
	"STREAMSW $7C",
	"UNPROTO CQ",
}

// Connect connects to the JNOS BBS at bbsAddress, by way of a Kantronics KPC 3
// Plus TNC attached to serialPort, and returns an open jnos.Conn for
// interaction with it.  (bbsAddress should consist of a call sign, a dash, and
// a small integer SSID.)  It logs into the specified BBS mailbox.  If callsign
// is set, it self-identifies periodically using the that call sign for FCC
// compliance.  (callsign should be the licensed FCC call sign of the calling
// user.)
func Connect(serialPort, bbsAddress, mailbox, callsign string) (c *jnos.Conn, err error) {
	var t *Transport

	if t, err = open(serialPort, bbsAddress, mailbox, callsign); err != nil {
		return nil, err
	}
	if c, err = jnos.Connect(t); err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	c.IdentEvery(10*time.Minute-30*time.Second, fmt.Sprintf("DE %s", callsign))
	return c, nil
}

// Open opens a transport to the JNOS BBS at bbsAddress, by way of a Kantronics
// KPC 3 Plus TNC attached to serialPort.  (bbsAddress should consist of a call
// sign, a dash, and a small integer SSID.)  It logs into the mailbox
// corresponding to the specified callsign, which must be the licensed FCC call
// sign of the calling user.  (For connecting to other mailboxes, see the
// Connect function.)
func Open(serialPort, bbsAddress, callsign string) (t *Transport, err error) {
	return open(serialPort, bbsAddress, callsign, "")
}

// open is the common code between Connect and Open.
func open(serialPort, bbsAddress, mailbox, callsign string) (t *Transport, err error) {
	t = new(Transport)
	if err = t.connectSerial(serialPort); err != nil {
		return nil, err
	}
	// Send a newline and check that we get a command prompt.  We may get
	// other stuff ahead of it, and may get more than one command prompt.
	if err = t.sendRaw([]byte{'\r'}); err != nil {
		goto ERROR
	}
	if _, err = t.readUntil(commandPrompt, tncTimeout); err != nil {
		goto ERROR
	}
	// Apply the pre-connect settings.
	for _, c := range preConnectCommands {
		if err = t.send(c); err != nil {
			goto ERROR
		}
		if _, err = t.readUntil(commandPrompt, tncTimeout); err != nil {
			goto ERROR
		}
	}
	// Set the mailbox we want to connect to as our "call sign".
	t.callsign = callsign
	if err = t.send(fmt.Sprintf("MY %s\n", mailbox)); err != nil {
		goto ERROR
	}
	if _, err = t.readUntil(commandPrompt, tncTimeout); err != nil {
		goto ERROR
	}
	// Connect to the BBS.
	t.wasConnected = true
	if err = t.send(fmt.Sprintf("CONNECT %s\n", bbsAddress)); err != nil {
		goto ERROR
	}
	if _, err = t.readUntil(commandPrompt, tncTimeout); err != nil {
		goto ERROR
	}
	t.connected = true
	return t, nil
ERROR:
	t.Close()
	return nil, fmt.Errorf("BBS connect: %s", err)
}

// Transport is the KPC 3 Plus transport to the JNOS BBS.
type Transport struct {
	serial       io.ReadWriteCloser
	readch       <-chan []byte
	errch        <-chan error
	pending      []byte
	connected    bool
	wasConnected bool
	callsign     string
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
	var (
		untilb []byte
		timer  *time.Timer
	)
	untilb = bytes.ReplaceAll([]byte(until), lf, crlf)
	timer = time.NewTimer(timeout)
	timer.Stop()
	for err == nil {
		if data, err = t.checkDisconnected(); err != nil {
			return data, err
		}
		if data = t.checkPending(untilb); data != "" {
			return data, nil
		}
		timer.Reset(timeout)
		select {
		case read, ok := <-t.readch:
			if !ok {
				err = <-t.errch
			} else {
				t.pending = append(t.pending, read...)
			}
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			err = jnos.ErrTimeout
		}
	}
	traceReceived(t.pending)
	data = string(bytes.ReplaceAll(t.pending, crlf, lf))
	t.pending = t.pending[:0]
	return data, err
}

// checkDisconnected checks whether the pending input buffer includes a message
// from the TNC indicating that its connection to the BBS was lost.  If so, it
// returns a disconnected error, along with everything in the pending buffer
// *except* that message.
func (t *Transport) checkDisconnected() (data string, err error) {
	if !t.connected {
		return "", nil
	}
	if idx := bytes.Index(t.pending, disconnectedMessage); idx >= 0 {
		traceReceived(t.pending)
		t.pending = append(t.pending[:idx], t.pending[idx+len(disconnectedMessage):]...)
		data = string(bytes.ReplaceAll(t.pending, crlf, lf))
		t.pending = t.pending[:0]
		t.connected = false
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
		traceReceived(t.pending[:idx])
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
		tosend []byte
		echo   []byte
		timer  *time.Timer
	)
	tosend = bytes.ReplaceAll([]byte(data), lf, cr)
	if len(tosend) == 0 || tosend[len(tosend)-1] != '\r' {
		tosend = append(tosend, '\r')
	}
	if err = t.sendRaw(tosend); err != nil {
		return err
	}
	echo = bytes.ReplaceAll(tosend, cr, crlf)
	timer = time.NewTimer(echoTimeout)
	timer.Stop()
	for err == nil {
		if _, err = t.checkDisconnected(); err != nil {
			break
		}
		if idx := bytes.Index(t.pending, echo); idx >= 0 {
			t.pending = append(t.pending[:idx], t.pending[idx+len(echo):]...)
			return nil
		}
		timer.Reset(echoTimeout)
		select {
		case read, ok := <-t.readch:
			if !ok {
				err = <-t.errch
			} else {
				t.pending = append(t.pending, read...)
			}
			if !timer.Stop() {
				<-timer.C
			}
		case <-timer.C:
			err = ErrBadEcho
		}
	}
	return err
}

// sendRaw sends the specified data with no translation and no wait for echo.
func (t *Transport) sendRaw(data []byte) (err error) {
	var count int

	if trace {
		for _, line := range traceLines(data) {
			fmt.Printf("Tx %-36s >\n", line)
		}
	}
	for len(data) != 0 {
		if count, err = t.serial.Write(data); err != nil {
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
			return fmt.Errorf("unable to get back to command mode for cleanup: %s", err)
		}
		if _, err = t.readUntil(commandPrompt, tncTimeout); err != nil && err != jnos.ErrDisconnected {
			return fmt.Errorf("unable to get back to command mode for cleanup: %s", err)
		}
		// Then tell the TNC to disconnect from the BBS, and wait until
		// we get the *** DISCONNECTED message.
		if err = t.send("D"); err != nil && err != jnos.ErrDisconnected {
			return fmt.Errorf("cleanup: can't disconnect: %s", err)
		} else if err == nil {
			if _, err = t.readUntil(string(disconnectedMessage), rfTimeout); err != jnos.ErrDisconnected {
				return fmt.Errorf("ERROR: cleanup: expected ErrDisconnected, got %s", err)
			}
		}
	}
	// Send our FCC identification if needed.
	if t.wasConnected && t.callsign != "" {
		var abort bool
		if abort, err = t.postIdentify(); abort {
			return err
		}
	}
	// Apply all of the post-connect settings.
	if t.callsign != "" {
		if err2 := t.send(fmt.Sprintf("MY %s\n", t.callsign)); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore settings: %s", err2)
		}
		if _, err2 := t.readUntil(commandPrompt, tncTimeout); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore settings: %s", err2)
		}
	}
	for _, c := range postDisconnectCommands {
		if err2 := t.send(c); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore settings: %s", err2)
		}
		if _, err2 := t.readUntil(commandPrompt, tncTimeout); err == nil && err2 != nil {
			err = fmt.Errorf("cleanup: restore settings: %s", err2)
		}
	}
	return err
}

// postIdentify sends the FCC callsign in CONVERS mode after disconnecting from
// the BBS.  If it hits an error, it returns the error and also a flag
// indicating whether subsequent post-disconnect steps should be aborted.
func (t *Transport) postIdentify() (abort bool, err error) {
	// Enter converse mode.
	if err = t.send("CONV"); err != nil {
		return false, fmt.Errorf("cleanup: send FCC ID: %s", err)
	}
	// Send our call sign identification.
	var ident = fmt.Sprintf("DE %s\n", t.callsign)
	if err = t.send(ident); err != nil {
		err = fmt.Errorf("cleanup: send FCC ID: %s", err)
	} else {
		// Wait until we get the monitor message saying that the ident has gone
		// out on the air.  (If we don't wait, it stays in the transmit queue
		// until the next time we connect, and winds up being sent as input to
		// the next BBS we connect to!)
		if _, err = t.readUntil(ident, rfTimeout); err != nil {
			err = fmt.Errorf("cleanup: send FCC ID: %s", err)
		}
	}
	// Return to command mode.
	if err2 := t.sendRaw([]byte{3}); err2 != nil {
		return true, fmt.Errorf("cleanup: send FCC ID: exit CONVERS mode: %s", err2)
	}
	if _, err2 := t.readUntil(commandPrompt, tncTimeout); err2 != nil {
		return true, fmt.Errorf("cleanup: send FCC ID: exit CONVERS mode: %s", err)
	}
	return false, err
}

// UseVerboseReads returns whether it's appropriate to use verbose reads
// when communicating over this transport.
func (t *Transport) UseVerboseReads() bool { return false }

// connectSerial opens a connection to the specified serial port, and starts a
// background reader thread for it.
func (t *Transport) connectSerial(serialPort string) (err error) {
	var oo = serial.OpenOptions{
		PortName:          serialPort,
		BaudRate:          9600,
		DataBits:          8,
		StopBits:          1,
		ParityMode:        serial.PARITY_NONE,
		RTSCTSFlowControl: true,
		MinimumReadSize:   1,
	}
	if t.serial, err = serial.Open(oo); err != nil {
		return err
	}
	var readch = make(chan []byte, 16)
	var errch = make(chan error, 1)
	go t.reader(readch, errch)
	t.readch = readch
	t.errch = errch
	return nil
}

// reader is a background thread constantly reading the serial port.  Reader
// sends incoming data on the readch channel until an error occurs, at which
// point it closes the readch channel, sends the error on the errch channel,
// and exits.  Note that when the connection is closed, an error is expected.
func (t *Transport) reader(readch chan<- []byte, errch chan<- error) {
	var (
		buf   [1024]byte
		count int
		err   error
	)
	for {
		count, err = t.serial.Read(buf[:])
		if count != 0 {
			var out = make([]byte, count)
			copy(out, buf[:count])
			readch <- out
		}
		if err != nil {
			close(readch)
			errch <- err
			return
		}
	}
}

// traceReceived emits trace messages containing received data, if tracing is
// enabled.
func traceReceived(data []byte) {
	if !trace || len(data) == 0 {
		return
	}
	for _, line := range traceLines(data) {
		fmt.Println("Rx                                      <", line)
	}
}

// traceLines formats data appropriately for tracing.
func traceLines(b []byte) (lines []string) {
	var nl = bytes.IndexByte(b, '\n')
	if nl < 0 {
		nl = bytes.IndexByte(b, '\r')
	}
	for nl >= 0 {
		lines = append(lines, traceLine(b[:nl+1])...)
		b = b[nl+1:]
		nl = bytes.IndexByte(b, '\n')
		if nl < 0 {
			nl = bytes.IndexByte(b, '\r')
		}
	}
	if len(b) != 0 {
		lines = append(lines, traceLine(b)...)
	}
	return lines
}
func traceLine(b []byte) (lines []string) {
	var q = strconv.Quote(string(b))
	q = q[1 : len(q)-1]
	for len(q) >= 36 {
		lines = append(lines, q[:36])
		q = q[36:]
	}
	if len(q) != 0 {
		lines = append(lines, q)
	}
	return lines
}
