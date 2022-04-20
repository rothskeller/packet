// Package telnet provides a transport layer for communicating with a JNOS BBS
// over a network.  The Open function, if successful, returns a Transport that
// can be passed to jnos.Connect.
package telnet

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"time"

	"steve.rothskeller.net/packet/jnos"
)

// trace is a flag indicating that data sent and retrieved should be echoed to
// stdout.  It is normally false, but can be set to true for debugging.
const trace = false

// bbsTimeout is the amount of time to wait for response from the BBS.
const bbsTimeout = 5 * time.Second

// Constant byte slices.
var (
	crlf = []byte{'\r', '\n'}
	lf   = []byte{'\n'}
)

// md5ChallengeRE is a regular expression matching the JNOS password prompt with
// an embedded MD5 challenge.
var md5ChallengeRE = regexp.MustCompile(`Password \[([0-9a-f]{1,8})\] : $`)

// Connect connects to the JNOS BBS at bbsAddress (host:port).  It logs into the
// specified BBS mailbox with the specified password.
func Connect(bbsAddress, mailbox, password string) (c *jnos.Conn, err error) {
	var t *Transport

	if t, err = Open(bbsAddress, mailbox, password); err != nil {
		return nil, err
	}
	if c, err = jnos.Connect(t); err != nil {
		t.Close()
		return nil, fmt.Errorf("BBS connect: %s", err)
	}
	return c, nil
}

// Open connects to the JNOS BBS at bbsAddress (host:port).  It logs into the
// specified BBS mailbox with the specified password.
func Open(bbsAddress, mailbox, password string) (t *Transport, err error) {
	var pwprompt string

	t = new(Transport)
	if err = t.connectTelnet(bbsAddress); err != nil {
		return nil, err
	}
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

// Transport is the KPC 3 Plus transport to the JNOS BBS.
type Transport struct {
	telnet    io.ReadWriteCloser
	readch    <-chan []byte
	errch     <-chan error
	pending   []byte
	connected bool
}

// ReadUntil readch data from the BBS until the specified string is seen, or a
// timeout occurs.  It returns the data that was read (even if it returns an
// error.)
func (t *Transport) ReadUntil(until string) (data string, err error) {
	var (
		untilb []byte
		timer  *time.Timer
	)
	untilb = bytes.ReplaceAll([]byte(until), lf, crlf)
	timer = time.NewTimer(bbsTimeout)
	timer.Stop()
	for err == nil {
		if data = t.checkPending(untilb); data != "" {
			return data, nil
		}
		timer.Reset(bbsTimeout)
		select {
		case read, ok := <-t.readch:
			if !ok {
				if err = <-t.errch; err == io.EOF {
					err = jnos.ErrDisconnected
				}
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
	for len(tosend) != 0 {
		if count, err = t.telnet.Write(tosend); err != nil {
			return err
		}
		tosend = tosend[count:]
	}
	return nil
}

// Close closes the connection to the BBS.
func (t *Transport) Close() error {
	return t.telnet.Close()
}

// connectTelnet opens a connection to the specified network address, and starts
// a background reader thread for it.
func (t *Transport) connectTelnet(address string) (err error) {
	if t.telnet, err = net.Dial("tcp", address); err != nil {
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
		count, err = t.telnet.Read(buf[:])
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
