// Package telnet provides a transport layer for communicating with a JNOS BBS
// over a network.  The Open function, if successful, returns a Transport that
// can be passed to jnos.Connect.
package telnet

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
)

// trace is a flag indicating that data sent and retrieved should be echoed to
// stdout.  It is normally false, but can be set to true for debugging.
const trace = false

// bbsTimeout is the amount of time to wait for response from the BBS.
const bbsTimeout = 5 * time.Second

// dialTimeout is the amount of time to wait for the connection to the BBS to be
// established.
const dialTimeout = 30 * time.Second

// writeTimeout is the amount of time to wait for each Send to be written to
// the network.
const writeTimeout = 60 * time.Second

// Constant byte slices.
var (
	crlf = []byte{'\r', '\n'}
	lf   = []byte{'\n'}
)

// md5ChallengeRE is a regular expression matching the JNOS password prompt with
// an embedded MD5 challenge.
var md5ChallengeRE = regexp.MustCompile(`Password \[([0-9a-f]{1,8})\] : $`)

// Connect connects to the JNOS BBS at bbsAddress (host:port).  It logs into the
// specified BBS mailbox with the specified password.  If log is non-nil, all
// traffic on the connection is written to the log.
func Connect(bbsAddress, mailbox, password string, log io.Writer) (c *jnos.Conn, err error) {
	return ConnectContext(context.Background(), bbsAddress, mailbox, password, log)
}

// ConnectContext is Connect, except that canceling ctx before it returns
// abandons the connection attempt promptly.  (Once it returns, ctx no longer
// affects the connection.)
func ConnectContext(ctx context.Context, bbsAddress, mailbox, password string, log io.Writer) (c *jnos.Conn, err error) {
	var t *Transport

	if t, err = open(ctx, bbsAddress, mailbox, password, log); err != nil {
		return nil, err
	}
	stop := context.AfterFunc(ctx, t.Abort)
	c, err = jnos.Connect(t, "")
	stop()
	if err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return c, nil
}

// Open connects to the JNOS BBS at bbsAddress (host:port).  It logs into the
// specified BBS mailbox with the specified password.  If log is non-nil, all
// traffic on the connection is written to the log.
func Open(bbsAddress, mailbox, password string, log io.Writer) (t *Transport, err error) {
	return open(context.Background(), bbsAddress, mailbox, password, log)
}

// open is Open, except that canceling ctx before it returns abandons the
// connection attempt promptly.
func open(ctx context.Context, bbsAddress, mailbox, password string, log io.Writer) (t *Transport, err error) {
	var pwprompt string

	t = new(Transport)
	t.log = log
	if err = t.connectTelnet(ctx, bbsAddress); err != nil {
		return nil, err
	}
	defer context.AfterFunc(ctx, t.Abort)()
	// Wait for a login: prompt.
	if _, err = t.ReadUntil("login: "); err != nil {
		goto ERROR
	}
	// Send the mailbox.
	if err = t.Send(mailbox); err != nil {
		goto ERROR
	}
	// Wait for the colon at the end of the password prompt.
	if pwprompt, err = t.ReadUntil(": "); err != nil {
		goto ERROR
	}
	// Get the MD5 challenge from the password prompt and send the hashed
	// password response.
	if match := md5ChallengeRE.FindStringSubmatch(pwprompt); match != nil {
		var (
			challenge uint64
			buf       []byte
			sum       [16]byte
		)
		challenge, _ = strconv.ParseUint(match[1], 16, 32)
		buf = make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(challenge))
		buf = append(buf, password...)
		sum = md5.Sum(buf)
		password = hex.EncodeToString(sum[:])
	} else {
		slog.Error("JNOS password prompt did not include MD5 challenge")
		err = errors.New("JNOS password prompt did not include MD5 challenge")
		goto ERROR
	}
	if err = t.Send(password); err != nil {
		goto ERROR
	}
	return t, nil
ERROR:
	t.Close()
	return nil, fmt.Errorf("BBS connect: %s", err)
}

// Transport is the telnet (TCP) transport to the JNOS BBS.
type Transport struct {
	telnet    net.Conn
	readch    <-chan []byte
	pending   []byte
	connected bool
	log       io.Writer

	// done is closed (by Abort or Close) to tell the reader goroutine and
	// any blocked ReadUntil to stop.  readerDone is closed when the reader
	// goroutine has exited.
	done       chan struct{}
	readerDone chan struct{}
	stopOnce   sync.Once
	closeOnce  sync.Once
	closeErr   error

	// mu protects err and readErr.  err is the sticky terminal error,
	// recorded the first time ReadUntil sees the connection end; readErr
	// is the error that ended the reader goroutine, set before readch is
	// closed.
	mu      sync.Mutex
	err     error
	readErr error

	// logmu serializes all writes to log, which need not be safe for
	// concurrent use.  logClosed is set by Close so that nothing is written
	// to log after Close returns.
	logmu     sync.Mutex
	logClosed bool
}

// Compile-time checks that Transport implements the jnos interfaces.
var (
	_ jnos.Transport = (*Transport)(nil)
	_ jnos.Aborter   = (*Transport)(nil)
)

// ReadUntil reads data from the BBS until the specified string is seen, or a
// timeout occurs.  It returns the data that was read (even if it returns an
// error.)
func (t *Transport) ReadUntil(until string) (data string, err error) {
	var (
		untilb []byte
		timer  *time.Timer
	)
	if err = t.terminalErr(false); err != nil {
		return "", err
	}
	untilb = bytes.ReplaceAll([]byte(until), lf, crlf)
	// Since Go 1.23, Stop and Reset guarantee that no stale value is
	// received from timer.C afterward, so no draining is needed.
	timer = time.NewTimer(bbsTimeout)
	defer timer.Stop()
	for err == nil {
		if data = t.checkPending(untilb); data != "" {
			return data, nil
		}
		timer.Reset(bbsTimeout)
		select {
		case read, ok := <-t.readch:
			if !ok {
				err = t.setTerminalErr(t.getReadErr())
			} else {
				t.pending = append(t.pending, read...)
			}
		case <-t.done:
			err = t.setTerminalErr(jnos.ErrDisconnected)
		case <-timer.C:
			slog.Error("timeout")
			err = jnos.ErrTimeout
		}
	}
	traceReceived(t.pending)
	data = string(bytes.ReplaceAll(t.pending, crlf, lf))
	t.pending = t.pending[:0]
	return data, err
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
func (t *Transport) Send(data string) (err error) {
	var (
		tosend []byte
		count  int
	)
	if err = t.terminalErr(true); err != nil {
		return err
	}
	// tosend = bytes.ReplaceAll([]byte(data), lf, cr)
	tosend = []byte(data)
	if len(tosend) == 0 || tosend[len(tosend)-1] != '\n' {
		tosend = append(tosend, '\n')
	}
	if trace {
		for _, line := range traceLines(tosend) {
			fmt.Printf("Tx %-36s >\n", line)
		}
	}
	t.writeLog(tosend)
	if err = t.telnet.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		if t.stopped() {
			return jnos.ErrDisconnected
		}
		slog.Error("telnet.SetWriteDeadline", "err", err)
		return err
	}
	for len(tosend) != 0 {
		if count, err = t.telnet.Write(tosend); err != nil {
			if t.stopped() {
				return jnos.ErrDisconnected
			}
			slog.Error("telnet.Write", "err", err)
			return err
		}
		tosend = tosend[count:]
	}
	return nil
}

// Abort ends the connection immediately, without an orderly disconnection
// from the BBS.  Any ReadUntil or Send in progress on another goroutine
// returns promptly with jnos.ErrDisconnected, as does any later call.  Abort
// is safe to call at any time, more than once, and concurrently with the other
// methods.  Close must still be called afterward.
func (t *Transport) Abort() {
	t.stop()
}

// Close closes the connection to the BBS.  It waits for the background reader
// to exit, and nothing is written to the log after it returns.  It is safe to
// call more than once, and after Abort.
func (t *Transport) Close() error {
	t.closeOnce.Do(func() {
		var err = t.stop()
		if t.readerDone != nil {
			<-t.readerDone
		}
		t.logmu.Lock()
		t.logClosed = true
		t.logmu.Unlock()
		if err != nil {
			slog.Error("telnet.Close", "err", err)
			t.closeErr = err
		}
	})
	return t.closeErr
}

// stop closes the done channel and the socket, the first time it is called.
// It returns the error from closing the socket on that first call, and nil on
// any later call.
func (t *Transport) stop() (err error) {
	t.stopOnce.Do(func() {
		if t.done != nil {
			close(t.done)
		}
		if t.telnet != nil {
			err = t.telnet.Close()
		}
	})
	return err
}

// stopped returns whether Abort or Close has been called.
func (t *Transport) stopped() bool {
	if t.done == nil {
		return false
	}
	select {
	case <-t.done:
		return true
	default:
		return false
	}
}

// terminalErr returns the error that ended the connection, if it is known to
// have ended, or nil otherwise.  If includeReader is true, an error seen by
// the reader goroutine counts even if ReadUntil has not yet consumed it.
func (t *Transport) terminalErr(includeReader bool) error {
	if t.stopped() {
		return t.setTerminalErr(jnos.ErrDisconnected)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err != nil {
		return t.err
	}
	if includeReader && t.readErr != nil {
		return mapReadErr(t.readErr)
	}
	return nil
}

// setTerminalErr records err (mapped as needed) as the sticky terminal error,
// unless one is already recorded, and returns the recorded error.
func (t *Transport) setTerminalErr(err error) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.err == nil {
		if t.stopped() {
			err = jnos.ErrDisconnected
		}
		t.err = mapReadErr(err)
	}
	return t.err
}

// getReadErr returns the error that ended the reader goroutine.
func (t *Transport) getReadErr() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.readErr
}

// mapReadErr maps the error that ended the reader to the one that should be
// returned to callers.
func mapReadErr(err error) error {
	if err == nil || err == io.EOF {
		return jnos.ErrDisconnected
	}
	return err
}

// writeLog writes b to the log, if there is one, serializing all such writes
// and suppressing them after Close.
func (t *Transport) writeLog(b []byte) {
	if t.log == nil {
		return
	}
	t.logmu.Lock()
	defer t.logmu.Unlock()
	if !t.logClosed {
		t.log.Write(b)
	}
}

// UseVerboseReads returns whether it's appropriate to use verbose reads
// when communicating over this transport.
func (t *Transport) UseVerboseReads() bool { return true }

// connectTelnet opens a connection to the specified network address, and starts
// a background reader thread for it.
func (t *Transport) connectTelnet(ctx context.Context, address string) (err error) {
	var conn net.Conn

	dialer := net.Dialer{Timeout: dialTimeout}
	if conn, err = dialer.DialContext(ctx, "tcp", address); err != nil {
		slog.Error("net.Dial", "addr", address, "err", err)
		return err
	}
	t.start(conn)
	return nil
}

// start sets up the Transport to use the specified connection, and starts the
// background reader thread for it.
func (t *Transport) start(conn net.Conn) {
	var readch = make(chan []byte, 16)

	t.telnet = conn
	t.readch = readch
	t.done = make(chan struct{})
	t.readerDone = make(chan struct{})
	go t.reader(conn, readch)
}

// reader is a background thread constantly reading the connection.  Reader
// sends incoming data on the readch channel until an error occurs or the
// Transport is stopped, at which point it records the error, closes the readch
// channel, and exits.  Note that when the connection is closed, an error is
// expected.
func (t *Transport) reader(conn net.Conn, readch chan<- []byte) {
	var (
		buf   [1024]byte
		count int
		err   error
	)
	defer close(t.readerDone)
	defer close(readch)
	for {
		count, err = conn.Read(buf[:])
		if count != 0 {
			var out = make([]byte, count)
			copy(out, buf[:count])
			t.writeLog(out)
			select {
			case readch <- out:
			case <-t.done:
				t.setReadErr(jnos.ErrDisconnected)
				return
			}
		}
		if err != nil {
			if err != io.EOF && !t.stopped() {
				slog.Error("telnet.Read", "err", err)
			}
			t.setReadErr(err)
			return
		}
	}
}

// setReadErr records the error that ended the reader goroutine.
func (t *Transport) setReadErr(err error) {
	t.mu.Lock()
	t.readErr = err
	t.mu.Unlock()
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
