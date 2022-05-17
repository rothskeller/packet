// Package jnos provides a library for interacting with a JNOS BBS.  To use it,
// first create a Transport that will handle the physical communication.  This
// repository contains three Transport implementations: kpc3plus, telnet, and
// simulator.  Others can be implemented elsewhere as long as they honor the
// Transport interface in this package.
//
// Once you have an open Transport, pass it to this package's Connect method to
// create a Conn.  Then you can call the Conn methods corresponding to JNOS BBS
// commands.  Call Close on the Conn when finished.
package jnos

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//--- TRANSPORT INTERFACE ------------------------------------------------------

// Transport is the interface that transport layers passed to Connect must
// satisfy.  All of the transport layers in the subpackages of this library
// honor this interface.
type Transport interface {
	// ReadUntil reads data from the BBS until the specified string is seen,
	// or a timeout occurs.  It returns the data that was read (even if it
	// returns an error.)  It returns ErrDisconnected if the connection to
	// the BBS was lost, ErrTimeout if the read timed out before the
	// specified string was found, and any other error if something else
	// goes wrong.
	ReadUntil(string) (string, error)
	// ReadUntilT reads data from the BBS until the specified string is
	// seen, or the specified timeout occurs.  It returns the data that was
	// read (even if it returns an error.)  It returns ErrDisconnected if
	// the connection to the BBS was lost, ErrTimeout if the read timed out
	// before the specified string was found, and any other error if
	// something else goes wrong.
	ReadUntilT(string, time.Duration) (string, error)
	// Send sends a string to the BBS.  It returns ErrDisconnected if the
	// connection to the BBS was lost, and any other error if something else
	// goes wrong.
	Send(string) error
	// Close closes the connection to the BBS.
	Close() error
}

// ErrDisconnected is returned when the connection disconnects unexpectedly.
var ErrDisconnected = errors.New("disconnected")

// ErrTimeout is returned when a read times out.
var ErrTimeout = errors.New("read timeout expired")

//--- JNOS CONN ----------------------------------------------------------------

// Conn is live connection to a JNOS BBS.
type Conn struct {
	t           Transport
	partialRead string
	ident       string
	identEvery  time.Duration
	nextIdent   time.Time
}

// jnosPromptRE is a regular expression matching the JNOS prompt.
var jnosPromptRE = regexp.MustCompile(`^\(#\d+\) >$`)

// Connect connects to the JNOS BBS over the supplied open Transport.
func Connect(t Transport) (c *Conn, err error) {
	c = &Conn{t: t}

	// Read and discard the connection message.
	if err = c.skipLinesUntilPrompt(); err != nil {
		return nil, err
	}
	// Turn off paging.
	if err = c.t.Send("XM 0\n"); err != nil {
		return nil, err
	}
	if err = c.skipLinesUntilPrompt(); err != nil {
		return nil, err
	}
	return c, nil
}

// IdentEvery causes the JNOS handler to send a comment with an identification
// every so often.  This can be used to send an FCC call sign if we're connected
// over the air with a tactical call and we're connected for longer than 10
// minutes.
func (c *Conn) IdentEvery(interval time.Duration, ident string) {
	c.ident = ident
	c.identEvery = interval
	c.nextIdent = time.Now().Add(interval)
}

// Send sends a private message.
func (c *Conn) Send(subject, body string, to ...string) (err error) {
	var cmd string

	defer c.maybeIdent()
	switch len(to) {
	case 0:
		panic("jnos.Send with no destination address")
	case 1:
		cmd = "SP"
	default:
		cmd = "SC"
	}
	if err = c.t.Send(fmt.Sprintf("%s %s\n", cmd, to[0])); err != nil {
		return err
	}
	if len(to) > 1 {
		if _, err = c.t.ReadUntil("Cc: "); err != nil {
			return err
		}
		if err = c.t.Send(strings.Join(to[1:], " ") + "\n"); err != nil {
			return err
		}
	}
	if _, err = c.t.ReadUntil("Subject:\n"); err != nil {
		return err
	}
	if err = c.t.Send(subject + "\n"); err != nil {
		return err
	}
	if _, err = c.t.ReadUntil("Enter message.  End with /EX or ^Z in first column (^A aborts):\n"); err != nil {
		return err
	}
	if !strings.HasSuffix(body, "\n") {
		body += "\n/EX\n"
	} else {
		body += "/EX\n"
	}
	if err = c.t.Send(body); err != nil {
		return err
	}
	if _, err = c.t.ReadUntil("Msg queued\n"); err != nil {
		return err
	}
	if err = c.skipLinesUntilPrompt(); err != nil {
		return err
	}
	return nil
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
var listHeaderNone = "None to list."
var listHeader = "St.  #  TO            FROM     DATE   SIZE SUBJECT"
var listMsgLineRE = regexp.MustCompile(`^[> ]([D ])([HYN]) (  \d| \d\d|\d+) ([^ ].{12}) ([^ ].{7}) ((?:Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d\d) (   \d|  \d\d| \d\d\d|\d+) (.*)`)

// List returns the list of messages.  If start is positive, only messages with
// that number or higher are returned.  If start is negative, only unread
// messages are returned.  If start is zero, all messages are returned.
func (c *Conn) List(start int) (ml *MessageList, err error) {
	var (
		cmd   string
		line  string
		match []string
	)
	defer c.maybeIdent()
	ml = new(MessageList)
	if start < 0 {
		cmd = "LM\n"
	} else if start > 0 {
		cmd = fmt.Sprintf("L %d\n", start)
	} else {
		cmd = "LA\n"
	}
	if err = c.t.Send(cmd); err != nil {
		return nil, err
	}
	// Line 1: Mail area: %s
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if !strings.HasPrefix(line, listLine1Prefix) {
		return nil, fmt.Errorf(`expected "Mail area: %%s" received %q`, line)
	}
	ml.Area = line[len(listLine1Prefix):]
	// Line 2: %d messages  -  %d new
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if match = listLine2RE.FindStringSubmatch(line); match == nil {
		return nil, fmt.Errorf(`expected "%%d message(s)  -  %%d new" received %q`, line)
	}
	ml.Count, _ = strconv.Atoi(match[1])
	ml.CountNew, _ = strconv.Atoi(match[2])
	// Line 3: blank
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if line != "" {
		return nil, fmt.Errorf(`expected "" received %q`, line)
	}
	// Line 4: possibly "None to list.", otherwise list header.
	if line, err = c.readLine(); err != nil {
		return nil, err
	} else if line == "None to list." {
		return nil, nil
	} else if line != listHeader {
		return nil, fmt.Errorf(`expected %q received %q`, listHeader, line)
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
		return nil, fmt.Errorf(`expected JNOS prompt, received %q`, line)
	}
	return ml, nil
}

var msgLine1RE = regexp.MustCompile(`^Message #(?:\d+) (?:\[Deleted|Held\])?$`)

// Read reads a single message given its number.  It returns nil if there is no
// such message.
func (c *Conn) Read(msgnum int, verbose bool) (msg string, err error) {
	var (
		cmd       string
		sb        strings.Builder
		line      string
		sawHeader bool
	)
	defer c.maybeIdent()
	if verbose {
		cmd = "V"
	} else {
		cmd = "R"
	}
	if err = c.t.Send(fmt.Sprintf("%s %d\n", cmd, msgnum)); err != nil {
		return "", err
	}
	if line, err = c.readLine(); err != nil {
		return "", err
	} else if strings.HasPrefix(line, "Invalid Message") || strings.HasPrefix(line, "No messages") {
		return "", c.skipLinesUntilPrompt()
	} else if !msgLine1RE.MatchString(line) {
		c.skipLinesUntilPrompt()
		return "", fmt.Errorf(`expected "Message %%d" received %q`, line)
	}
	for {
		if line, err = c.readLine(); err != nil {
			return "", err
		}
		if jnosPromptRE.MatchString(line) {
			if sawHeader {
				return sb.String(), nil
			}
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

	defer c.maybeIdent()
	sb.WriteString("K")
	for _, msgnum := range msgnums {
		fmt.Fprintf(&sb, " %d", msgnum)
	}
	sb.WriteByte('\n')
	if err = c.t.Send(sb.String()); err != nil {
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
		if !killMsgRE.MatchString(line) {
			c.skipLinesUntilPrompt()
			return fmt.Errorf(`expected "Msg %%d Killed." received %q`, line)
		}
	}
}

// maybeIdent checks for whether we should send an ident string, and if so, does
// so.  Errors are ignored.
func (c *Conn) maybeIdent() {
	if c.nextIdent.IsZero() || time.Now().Before(c.nextIdent) {
		return
	}
	c.nextIdent = time.Now().Add(c.identEvery)
	if err := c.t.Send(fmt.Sprintf("# %s\n", c.ident)); err != nil {
		return
	}
	c.skipLinesUntilPrompt()
}

// Close disconnects from the BBS and closes the underlying Transport
// connection.
func (c *Conn) Close() (err error) {
	var line string

	// Send the BYE command to the BBS.
	if err = c.t.Send("B\r"); err != nil {
		c.t.Close()
		return err
	}
	// Read until the underlying connection has been disconnected.
	if line, err = c.readLine(); err == nil {
		c.t.Close()
		return fmt.Errorf(`expected disconnect, received %q`, line)
	}
	if err != ErrDisconnected {
		c.t.Close()
		return fmt.Errorf(`expected disconnect, got error: %s`, err)
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
	response, err = c.t.ReadUntil("\n")
	if len(response) != 0 && response[len(response)-1] == '\n' {
		return response[:len(response)-1], err
	}
	return response, err
}
