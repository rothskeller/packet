// Package jnos provides a library for interacting with a JNOS BBS.  To use it,
// first create a Transport that will handle the physical communication.  This
// repository contains two Transport implementations: jnos/kpc3plus and
// jnos/telnet.  Others can be implemented elsewhere as long as they honor the
// Transport interface in this package.
//
// Once you have an open Transport, pass it to this package's Connect method to
// create a Conn.  Then you can call the Conn methods corresponding to JNOS BBS
// commands.  Call Close on the Conn when finished.
package jnos

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//--- TRANSPORT INTERFACE ------------------------------------------------------

// Transport is the interface that transport layers passed to Connect must
// satisfy.  All of the transport layers in the subpackages of this library
// honor this interface.
type Transport interface {
	// ReadUntil reads data from the BBS until the specified string is seen,
	// or a timeout occurs.  It returns the data that was read (even if it
	// returns an error).  It returns ErrDisconnected if the connection to
	// the BBS was lost, ErrTimeout if the read timed out before the
	// specified string was found, and any other error if something else
	// goes wrong.
	ReadUntil(string) (string, error)
	// Send sends a string to the BBS.  It returns ErrDisconnected if the
	// connection to the BBS was lost, and any other error if something else
	// goes wrong.
	Send(string) error
	// Close closes the connection to the BBS.
	Close() error
	// UseVerboseReads returns whether it's appropriate to use verbose reads
	// when communicating over this transport.
	UseVerboseReads() bool
}

// Aborter is an optional interface that a Transport can implement.  Abort ends
// the connection immediately, without an orderly disconnection from the BBS,
// and makes any Transport call in progress on another goroutine return
// promptly with an error, as will any later call.  Abort must be safe to call
// at any time, more than once, and concurrently with the other methods; after
// it, Close must still be called, and must return promptly.
type Aborter interface {
	Abort()
}

// ErrDisconnected is returned when the connection disconnects unexpectedly.
var ErrDisconnected = errors.New("disconnected")

// ErrTimeout is returned when a read times out.
var ErrTimeout = errors.New("read timeout expired")

// errAborted is the terminal error of a Conn that has been aborted.  It wraps
// ErrDisconnected, since the link is gone.
var errAborted = fmt.Errorf("connection aborted: %w", ErrDisconnected)

// errClosed is the terminal error of a Conn that has been closed.
var errClosed = errors.New("connection closed")

//--- JNOS CONN ----------------------------------------------------------------

// Conn is live connection to a JNOS BBS.
type Conn struct {
	t          Transport
	ident      string
	identEvery time.Duration
	nextIdent  time.Time
	aborted    atomic.Bool
	// err is the sticky terminal error.  Once it is set (because the
	// transport reported ErrDisconnected, or the Conn was closed), every
	// later method returns it without touching the transport.  It is only
	// accessed by the goroutine calling Conn methods; Abort, which may run
	// on another goroutine, uses aborted instead.
	err error
}

// Abort ends the connection immediately, without an orderly disconnection from
// the BBS.  It may be called from another goroutine while a Conn method is in
// progress; that method returns promptly with an error, as will any later
// method call.  Close must still be called afterward.
func (c *Conn) Abort() {
	if c.aborted.Swap(true) {
		return
	}
	if a, ok := c.t.(Aborter); ok {
		a.Abort()
	}
}

// terminalErr returns the sticky terminal error of the Conn, if any.
func (c *Conn) terminalErr() error {
	if c.aborted.Load() {
		return errAborted
	}
	return c.err
}

// record notes the error returned by a transport call, making it sticky if it
// means the link is gone.  It returns the error the caller should report.
func (c *Conn) record(err error) error {
	if err == nil {
		return nil
	}
	if c.aborted.Load() {
		return errAborted
	}
	// ErrTimeout is deliberately not sticky: the link may still be up (a
	// slow RF path), and the caller may retry or close in an orderly way.
	if errors.Is(err, ErrDisconnected) && c.err == nil {
		c.err = err
	}
	return err
}

// send sends a string to the BBS, unless the Conn has a terminal error.
func (c *Conn) send(s string) error {
	if err := c.terminalErr(); err != nil {
		return err
	}
	return c.record(c.t.Send(s))
}

// readUntil reads from the BBS until the specified string, unless the Conn has
// a terminal error.
func (c *Conn) readUntil(s string) (string, error) {
	if err := c.terminalErr(); err != nil {
		return "", err
	}
	data, err := c.t.ReadUntil(s)
	return data, c.record(err)
}

// drain reads and discards the rest of a response through the next prompt,
// after a response format error, so that the next command starts in sync.  It
// returns fmtErr; any error from the drain itself is not reported (but a
// disconnect is still recorded as sticky by readUntil).
func (c *Conn) drain(fmtErr error) error {
	_ = c.skipLinesUntilPrompt()
	return fmtErr
}

// identIfOK is deferred by the BBS command methods: it sends an ident if one
// is due, but only if the command succeeded.
func (c *Conn) identIfOK(err *error) {
	if *err == nil {
		c.maybeIdent()
	}
}

// jnosPromptRE is a regular expression matching the JNOS prompt.
var jnosPromptRE = regexp.MustCompile(`^(?:Area: \S+ )?\(#\d+\) >$`)

// Connect connects to the JNOS BBS over the supplied open Transport.  If an
// ident is supplied, it is sent as a JSON comment when connected, periodically
// during the connection, and before closing the connection.
func Connect(t Transport, ident string) (c *Conn, err error) {
	c = &Conn{t: t, ident: ident}

	// Read and discard the connection message.
	if err = c.skipLinesUntilPrompt(); err != nil {
		return nil, err
	}
	// Send the ident if any.
	if ident != "" {
		c.identEvery = 8 * time.Minute
		c.maybeIdent()
	}
	// Turn off paging.
	if err = c.send("XM 0\n"); err != nil {
		return nil, err
	}
	if err = c.skipLinesUntilPrompt(); err != nil {
		return nil, err
	}
	return c, nil
}

// Send sends a private message.  Note that the "to" addresses must be bare
// addresses without any display name or angle brackets.
func (c *Conn) Send(subject, body string, to ...string) (err error) {
	var cmd string

	switch len(to) {
	case 0:
		panic("jnos.Send with no destination address")
	case 1:
		cmd = "SP"
	default:
		cmd = "SC"
	}
	defer c.identIfOK(&err)
	if err = c.send(fmt.Sprintf("%s %s\n", cmd, to[0])); err != nil {
		return err
	}
	if len(to) > 1 {
		if _, err = c.readUntil("Cc: "); err != nil {
			return err
		}
		if err = c.send(strings.Join(to[1:], " ") + "\n"); err != nil {
			return err
		}
	}
	if _, err = c.readUntil("Subject:\n"); err != nil {
		return err
	}
	if err = c.send(subject + "\n"); err != nil {
		return err
	}
	if _, err = c.readUntil("Enter message.  End with /EX or ^Z in first column (^A aborts):\n"); err != nil {
		return err
	}
	if !strings.HasSuffix(body, "\n") {
		body += "\n/EX\n"
	} else {
		body += "/EX\n"
	}
	if err = c.send(body); err != nil {
		return err
	}
	if _, err = c.readUntil("Msg queued\n"); err != nil {
		return err
	}
	if err = c.skipLinesUntilPrompt(); err != nil {
		return err
	}
	return nil
}

// SendBulletin sends a bulletin message.  Note that the "to" address must be
// a bare address without any display name or angle brackets.
func (c *Conn) SendBulletin(subject, body string, to string) (err error) {
	defer c.identIfOK(&err)
	if err = c.send(fmt.Sprintf("SB %s\n", to)); err != nil {
		return err
	}
	if _, err = c.readUntil("Subject:\n"); err != nil {
		return err
	}
	if err = c.send(subject + "\n"); err != nil {
		return err
	}
	if _, err = c.readUntil("Enter message.  End with /EX or ^Z in first column (^A aborts):\n"); err != nil {
		return err
	}
	if !strings.HasSuffix(body, "\n") {
		body += "\n/EX\n"
	} else {
		body += "/EX\n"
	}
	if err = c.send(body); err != nil {
		return err
	}
	if _, err = c.readUntil("Msg queued\n"); err != nil {
		return err
	}
	if err = c.skipLinesUntilPrompt(); err != nil {
		return err
	}
	return nil
}

// SetArea switches to the specified message area.
func (c *Conn) SetArea(area string) (err error) {
	var line string

	defer c.identIfOK(&err)
	if err = c.send(fmt.Sprintf("A %s\n", area)); err != nil {
		return err
	}
	for {
		if line, err = c.readLine(); err != nil {
			return err
		}
		if jnosPromptRE.MatchString(line) {
			return nil
		}
		if strings.Contains(line, "No such message area") {
			slog.Error("no such message area", "area", area)
			return c.drain(errors.New("no such message area: " + area))
		}
	}
}

// MessageList contains a message list, retrieved with the List call.
type MessageList struct {
	Area     string
	Count    int
	CountNew int
	Messages []*MessageInfo
}

// MessageInfo contains the details of a single message in a MessageList.
type MessageInfo struct {
	// Number is the message number.  It is stable only during the BBS
	// connection.
	Number int
	// Deleted indicates that the message has been Killed and will be
	// deleted when the connection is closed.
	Deleted bool
	// Held indicates whether the message has been held.
	Held bool
	// Read indicates whether the message has been read.
	Read bool
	// ToPrefix is a truncated copy of the To: line of the message.
	ToPrefix string
	// FromPrefix is a truncated copy of the From: line of the message.
	FromPrefix string
	// Date is a partial date of the message.  It can be matched with the
	// "Jan 02" format string for time.Parse.
	Date string
	// Size is the size of the message in bytes.
	Size int
	// SubjectPrefix is a truncated copy of the message subject.
	SubjectPrefix string
}

var listLine1Prefix = "Mail area: "
var listLine2RE = regexp.MustCompile(`^(\d+) message(?:s)?  -  (\d+) new$`)
var listHeaderNone1 = "None to list."
var listHeaderNone2 = "No messages"
var listHeader = "St.  #  TO            FROM     DATE   SIZE SUBJECT"
var listMsgLineRE = regexp.MustCompile(`^[> ]([D ])([HYN]) (  \d| \d\d|\d+) ([^ ].{12}) ([^ ].{7}) ((?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (?: \d|\d\d)) (   \d|  \d\d| \d\d\d|\d+) (.*)`)

// List returns the list of messages.  If "to" is set, only messages sent to
// addresses containing that string are returned.  List returns nil if there are
// no messages.
func (c *Conn) List(to string) (ml *MessageList, err error) {
	var (
		cmd   string
		line  string
		match []string
	)
	defer c.identIfOK(&err)
	ml = new(MessageList)
	if to != "" {
		cmd = fmt.Sprintf("L> %s\n", to)
	} else {
		cmd = "LA\n"
	}
	if err = c.send(cmd); err != nil {
		return nil, err
	}
	// Line 1: Mail area: %s
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if line == listHeaderNone1 || line == listHeaderNone2 {
		return nil, c.skipLinesUntilPrompt()
	} else if !strings.HasPrefix(line, listLine1Prefix) {
		slog.Error("unexpected response", "exp", "Mail area: %s", "act", line)
		return nil, c.drain(fmt.Errorf(`expected "Mail area: %%s" received %q`, line))
	}
	ml.Area = line[len(listLine1Prefix):]
	// Line 2: %d messages  -  %d new
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if match = listLine2RE.FindStringSubmatch(line); match == nil {
		slog.Error("unexpected response", "exp", "%d message(s)  -  %d new", "act", line)
		return nil, c.drain(fmt.Errorf(`expected "%%d message(s)  -  %%d new" received %q`, line))
	}
	ml.Count, _ = strconv.Atoi(match[1])
	ml.CountNew, _ = strconv.Atoi(match[2])
	// Line 3: blank
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if line != "" {
		slog.Error("unexpected response", "exp", "", "act", line)
		return nil, c.drain(fmt.Errorf(`expected "" received %q`, line))
	}
	// Line 4: possibly "None to list.", otherwise list header.
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if line == listHeaderNone1 || line == listHeaderNone2 {
		return nil, c.skipLinesUntilPrompt()
	} else if line != listHeader {
		slog.Error("unexpected response", "exp", listHeader, "act", line)
		return nil, c.drain(fmt.Errorf(`expected %q received %q`, listHeader, line))
	}
	// Subsequent lines: messages.
	if line, err = c.readLine(); err != nil {
		return nil, err
	}
	match = listMsgLineRE.FindStringSubmatch(line)
	for match != nil {
		var mi MessageInfo
		if match[1] == "D" {
			mi.Deleted = true
		}
		switch match[2] {
		case "H":
			mi.Held = true
		case "Y":
			mi.Read = true
		}
		mi.Number, _ = strconv.Atoi(strings.TrimSpace(match[3]))
		mi.ToPrefix = strings.TrimSpace(match[4])
		mi.FromPrefix = strings.TrimSpace(match[5])
		mi.Date = match[6]
		mi.Size, _ = strconv.Atoi(strings.TrimSpace(match[7]))
		mi.SubjectPrefix = strings.TrimSpace(match[8])
		ml.Messages = append(ml.Messages, &mi)
		if line, err = c.readLine(); err != nil {
			return nil, err
		}
		match = listMsgLineRE.FindStringSubmatch(line)
	}
	// After the message lines, we may see a new mail line.
	if strings.HasPrefix(line, "You have new mail.") {
		if line, err = c.readLine(); err != nil {
			return nil, err
		}
	}
	// We should now be looking at the JNOS prompt.
	if !jnosPromptRE.MatchString(line) {
		slog.Error("unexpected response", "exp", "JNOS prompt", "act", line)
		return nil, c.drain(fmt.Errorf(`expected JNOS prompt, received %q`, line))
	}
	return ml, nil
}

var msgLine1RE = regexp.MustCompile(`^Message #\d+ (?:\[(?:Deleted|Held)\])?$`)
var removeNewMailRE = regexp.MustCompile(`\nYou have new mail.*\n$`)

// Read reads a single message given its number.  It returns nil if there is no
// such message.
func (c *Conn) Read(msgnum int) (msg string, err error) {
	var (
		cmd       string
		sb        strings.Builder
		line      string
		sawHeader bool
	)
	if err = c.terminalErr(); err != nil {
		return "", err
	}
	defer c.identIfOK(&err)
	if c.t.UseVerboseReads() {
		cmd = "V"
	} else {
		cmd = "R"
	}
	if err = c.send(fmt.Sprintf("%s %d\n", cmd, msgnum)); err != nil {
		return "", err
	}
	if line, err = c.readLine(); err != nil {
		return "", err
	} else if strings.HasPrefix(line, "Invalid Message") || strings.HasPrefix(line, "No messages") {
		return "", c.skipLinesUntilPrompt()
	} else if !msgLine1RE.MatchString(line) {
		slog.Error("unexpected response", "exp", "Message %d", "act", line)
		return "", c.drain(fmt.Errorf(`expected "Message %%d" received %q`, line))
	}
	for {
		if line, err = c.readLine(); err != nil {
			return "", err
		}
		if jnosPromptRE.MatchString(line) {
			if sawHeader {
				msg = sb.String()
				msg = removeNewMailRE.ReplaceAllLiteralString(msg, "\n")
				return msg, nil
			}
			slog.Error("unexpected response", "exp", "message body", "act", "JNOS prompt")
			return "", fmt.Errorf(`expected message body, received JNOS prompt`)
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
		if line == "" {
			sawHeader = true
		}
	}
}

var killMsgRE = regexp.MustCompile(`^Msg \d+ Killed.$`)

// Kill kills one or more messages.
func (c *Conn) Kill(msgnums ...int) (err error) {
	var sb strings.Builder

	defer c.identIfOK(&err)
	sb.WriteString("K")
	for _, msgnum := range msgnums {
		fmt.Fprintf(&sb, " %d", msgnum)
	}
	sb.WriteByte('\n')
	if err = c.send(sb.String()); err != nil {
		return err
	}
	for {
		var line string

		if line, err = c.readLine(); err != nil {
			return err
		}
		if jnosPromptRE.MatchString(line) {
			return nil
		}
		if strings.HasPrefix(line, "You have new mail.") {
			continue
		}
		if !killMsgRE.MatchString(line) {
			slog.Error("unexpected response", "exp", "Msg %d Killed.", "act", line)
			return c.drain(fmt.Errorf(`expected "Msg %%d Killed." received %q`, line))
		}
	}
}

// maybeIdent checks for whether we should send an ident string, and if so, does
// so.  Errors are ignored.  No ident is sent if the Conn has a terminal error.
func (c *Conn) maybeIdent() {
	if c.ident == "" || time.Now().Before(c.nextIdent) || c.terminalErr() != nil {
		return
	}
	c.nextIdent = time.Now().Add(c.identEvery)
	if err := c.send(fmt.Sprintf("# DE %s\n", c.ident)); err != nil {
		return
	}
	_ = c.skipLinesUntilPrompt()
}

// Close disconnects from the BBS and closes the underlying Transport
// connection.  If the connection has already been lost or aborted, it just
// closes the Transport.
func (c *Conn) Close() (err error) {
	var line string

	if c.err == errClosed {
		return nil // already closed
	}
	defer func() { c.err = errClosed }()
	if c.terminalErr() != nil {
		return c.t.Close()
	}
	// Send an ident if needed.
	if c.ident != "" {
		c.nextIdent = time.Time{}
		c.maybeIdent()
		if c.terminalErr() != nil {
			return c.t.Close()
		}
	}
	// Send the BYE command to the BBS.
	if err = c.send("B\r"); err != nil {
		c.t.Close()
		return err
	}
	// Read until the underlying connection has been disconnected.  Any
	// read error counts as a disconnect: transports report the BBS
	// dropping the link in various ways (ErrDisconnected, "no response
	// from BBS", ECONNRESET, ...).
	if line, err = c.readLine(); err == nil {
		c.t.Close()
		slog.Error("unexpected response", "exp", "disconnect", "act", line)
		return fmt.Errorf(`expected disconnect, received %q`, line)
	}
	if !errors.Is(err, ErrDisconnected) {
		slog.Info("treating read error after BYE as disconnect", "err", err)
	}
	return c.t.Close()
}

func (c *Conn) skipLinesUntilPrompt() (err error) {
	for {
		var line string

		if line, err = c.readLine(); err != nil {
			return err
		}
		if jnosPromptRE.MatchString(line) {
			return nil
		}
	}
}

func (c *Conn) readLine() (response string, err error) {
	response, err = c.readUntil("\n")
	if len(response) != 0 && response[len(response)-1] == '\n' {
		return response[:len(response)-1], err
	}
	return response, err
}
