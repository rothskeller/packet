package incident

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/v4/jnos"
	"github.com/rothskeller/packet/v4/jnos/kiss"
	"github.com/rothskeller/packet/v4/jnos/serialtnc"
	"github.com/rothskeller/packet/v4/jnos/telnet"
	"github.com/rothskeller/packet/v4/jnos/tnc"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/address"
	"github.com/rothskeller/packet/v4/message/field"
	"github.com/rothskeller/packet/v4/message/receipt"
	"k8s.io/apimachinery/pkg/util/sets"
)

var connectionsRunning = sets.New[string]()
var connectionsRunningMutex sync.Mutex

// errAborted is the error reported when the user aborts the connection.
var errAborted = errors.New("connection aborted by user")

// BBSExchange connects to a BBS and exchanges messages with it for the incident
// in the named directory (which should not be locked when the call is issued).
// The connection will be aborted if the supplied context is canceled.  If
// immOnly is true, only IMMEDIATE messages are transferred (and no bulletin
// checks are performed).  Progress updates and results are sent to the supplied
// watcher; its Finished method is called exactly once, last.  BBSExchange
// blocks, and should be run in a goroutine.
func BBSExchange(ctx context.Context, dir string, immOnly bool, updates BBSExchangeWatcher) {
	var e = &exchange{dir: dir, immOnly: immOnly, updates: updates}

	// All cleanup, including the call to Finished, is done by this single
	// deferred call, on every path (including a panic).
	defer e.cleanup(ctx)
	if err := e.run(ctx); err != nil {
		e.fail(err.Error())
	}
}

// run performs the BBS exchange.  It returns an error if the exchange fails.
// Resources that it acquires are recorded in e, and released by e.cleanup.
func (e *exchange) run(ctx context.Context) (err error) {
	// Read the incident configuration.
	if err = Read(e.dir, func(i *Incident) error {
		e.config = i.Config
		if !e.immOnly {
			e.areas = sets.New(i.BulletinAreasToCheck()...)
		}
		return nil
	}); err != nil {
		return err
	}
	// Make sure we don't already have a conflicting connection running.
	connTo := e.config.ActiveCall() + "@" + e.config.ConnectBBS
	connectionsRunningMutex.Lock()
	if connectionsRunning.Has(connTo) {
		connectionsRunningMutex.Unlock()
		slog.Error("already connected", "to", connTo)
		return fmt.Errorf("there is already a connection to %s in progress", connTo)
	}
	if e.config.UsesSerialPort() {
		if connectionsRunning.Has(e.config.SerialPort) {
			connectionsRunningMutex.Unlock()
			slog.Error("serial port in use", "port", e.config.SerialPort)
			return fmt.Errorf("can't connect using %s because it is in use", e.config.SerialPort)
		}
		connectionsRunning.Insert(e.config.SerialPort)
		e.reserved = append(e.reserved, e.config.SerialPort)
	}
	connectionsRunning.Insert(connTo)
	e.reserved = append(e.reserved, connTo)
	connectionsRunningMutex.Unlock()
	// Open the log file.
	logfname := filepath.Join(e.dir, e.config.ConnectBBS+".log")
	if e.logf, err = os.OpenFile(logfname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
		slog.Error("os.OpenFile", "f", logfname, "err", err)
		e.logf = nil
		return err
	}
	fmt.Fprintf(e.logf, "\n[[%s]]\n", time.Now().Format(time.RFC3339))
	// If the watcher wants a live monitor of the exchange, KISS connections
	// report the AX.25 frames to it; other connections report the text
	// exchanged with the BBS.
	var (
		monitor func(string)
		textlog io.Writer = e.logf
	)
	if m, ok := e.updates.(BBSExchangeMonitor); ok {
		monitor = m.Monitor
		if e.config.ConnectType != ConnectKISS {
			e.monitorlw = &lineWriter{fn: monitor}
			textlog = io.MultiWriter(e.logf, e.monitorlw)
		}
	}
	// Establish the JNOS connection.
	conn, err := connectBBS(e, ctx, textlog, monitor)
	if err != nil {
		fmt.Fprintf(e.logf, "[[ERROR: %s]]\n", err)
		if ctx.Err() != nil {
			return errAborted
		}
		return err
	}
	// From here on, canceling the context aborts the connection, making
	// any Conn call in progress return promptly with an error.
	e.conn = conn
	e.stopAbort = context.AfterFunc(ctx, conn.Abort)
	// Run the main connection loop.
	for done := false; !done; {
		if done, err = e.step(ctx); err != nil {
			break
		}
	}
	if err != nil && ctx.Err() != nil {
		// Whatever error we got is a consequence of the abort.
		err = errAborted
	}
	return err
}

// connectBBS establishes the JNOS connection for an exchange.  It is a variable
// so that tests can substitute a fake connection.
var connectBBS = (*exchange).connect

// connect establishes the JNOS connection, using the transport specified in
// the configuration.  textlog receives the text exchanged with the BBS, and
// monitor (if not nil) receives the Link Monitor display of KISS frames.
// Canceling ctx abandons the connection attempt.
func (e *exchange) connect(ctx context.Context, textlog io.Writer, monitor func(string)) (conn *jnos.Conn, err error) {
	switch e.config.ConnectType {
	case ConnectSerialTNC:
		tncconf := tnc.Get(e.config.TNCType)
		if tncconf == nil {
			slog.Error("no such TNC type", "type", e.config.TNCType)
			return nil, fmt.Errorf("The configuration specifies a TNC type of %q which is not defined.", e.config.TNCType)
		}
		e.progress("Connecting to %s@%s...", e.config.ActiveCall(), e.config.ConnectBBS)
		conn, err = serialtnc.ConnectContext(ctx, tncconf, e.config.SerialPort, e.config.ConnectAddress, e.config.ActiveCall(), e.config.OpCall, textlog)
	case ConnectKISS:
		dev := kiss.GetDevice(e.config.TNCType)
		if dev == nil {
			slog.Error("no such KISS device type", "type", e.config.TNCType)
			return nil, fmt.Errorf("The configuration specifies a KISS TNC type of %q which is not defined.", e.config.TNCType)
		}
		e.progress("Connecting to %s@%s...", e.config.ActiveCall(), e.config.ConnectBBS)
		conn, err = kiss.ConnectContext(ctx, dev, e.config.SerialPort, e.config.ConnectAddress, e.config.ActiveCall(), e.config.OpCall, e.logf, monitor)
	case ConnectTelnet:
		e.progress("Connecting to %s@%s...", e.config.TelnetUser, e.config.ConnectBBS)
		conn, err = telnet.ConnectContext(ctx, e.config.ConnectAddress, e.config.TelnetUser, e.config.TelnetPassword, textlog)
	default:
		slog.Error("invalid connect method", "method", e.config.ConnectType)
		err = fmt.Errorf("can't connect using method %s", e.config.ConnectType)
	}
	return conn, err
}

// cleanup is deferred by BBSExchange.  It recovers from any panic, closes the
// connection if it's still open (aborting it rather than disconnecting in an
// orderly fashion if the context was canceled), releases the reservations of
// the BBS and serial port, flushes the Link Monitor and closes the log file,
// and finally calls the watcher's Finished method.
func (e *exchange) cleanup(ctx context.Context) {
	if p := recover(); p != nil {
		slog.Error("panic in BBSExchange", "err", p, "stack", string(debug.Stack()))
		e.fail("A software error occurred.  See log file for details.")
	}
	if e.conn != nil {
		if ctx.Err() != nil {
			e.conn.Abort() // makes Close prompt, with no BBS dialog
		} else {
			e.safely("reporting progress", func() { e.progress("Closing connection...") })
		}
		e.safely("closing connection", func() {
			if err := e.conn.Close(); err != nil && !e.failed {
				if ctx.Err() != nil {
					err = errAborted
				}
				e.fail(err.Error())
			}
		})
		if e.stopAbort != nil {
			e.stopAbort()
		}
		e.conn = nil
	}
	if len(e.reserved) != 0 {
		connectionsRunningMutex.Lock()
		connectionsRunning.Delete(e.reserved...)
		connectionsRunningMutex.Unlock()
		e.reserved = nil
	}
	if e.monitorlw != nil {
		e.safely("flushing monitor", e.monitorlw.Flush)
	}
	if e.logf != nil {
		e.logf.Close()
		e.logf = nil
	}
	e.updates.Finished()
}

// safely calls fn, recovering from (and logging) any panic in it, so that the
// rest of the cleanup still happens.
func (e *exchange) safely(what string, fn func()) {
	defer func() {
		if p := recover(); p != nil {
			slog.Error("panic in BBSExchange "+what, "err", p, "stack", string(debug.Stack()))
			if !e.failed {
				e.fail("A software error occurred.  See log file for details.")
			}
		}
	}()
	fn()
}

// fail reports an error to the watcher.
func (e *exchange) fail(msg string) {
	e.failed = true
	e.updates.Error(msg)
}

// BBSExchangeWatcher is the interface required for the "updates" parameter to
// BBSExchange; it is an object to which progress of the BBS exchange is
// reported.
type BBSExchangeWatcher interface {
	// Progress is called with messages indicating the progress of the BBS
	// exchange.
	Progress(string)
	// LogEntry is called with the log entry for a received or transmitted
	// message.
	LogEntry(*LogEntry)
	// Error is called if an error occurs during the BBS exchange.
	Error(string)
	// Finished is called when the BBS exchange has finished, whether or not
	// it was successful.  (It was successful if Error was never called.)
	Finished()
}

// BBSExchangeMonitor is an optional extension of BBSExchangeWatcher.  If the
// watcher passed to BBSExchange implements it, Monitor is called with each line
// of a live display of the exchange: for KISS connections, a packet monitor
// display of the AX.25 frames exchanged with the BBS (e.g., "06:57:02T
// KO6NVY>W3XSC-1 <SABM C P>", followed by the text of each I frame); for other
// connections, the text exchanged with the BBS.
type BBSExchangeMonitor interface {
	Monitor(line string)
}

// lineWriter is an io.Writer that calls fn with each complete line written to
// it, without its line ending.  It is safe for concurrent use (transports may
// write to it from more than one goroutine); fn is called with the mutex held,
// so lines are reported in order and fn is never called concurrently.
type lineWriter struct {
	mu  sync.Mutex
	fn  func(string)
	buf []byte
}

func (lw *lineWriter) Write(p []byte) (int, error) {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	lw.buf = append(lw.buf, p...)
	for {
		idx := bytes.IndexByte(lw.buf, '\n')
		if idx < 0 {
			return len(p), nil
		}
		lw.fn(strings.TrimSuffix(string(lw.buf[:idx]), "\r"))
		lw.buf = lw.buf[idx+1:]
	}
}

// Flush reports any partial line remaining.
func (lw *lineWriter) Flush() {
	lw.mu.Lock()
	defer lw.mu.Unlock()
	if len(lw.buf) != 0 {
		lw.fn(strings.TrimSuffix(string(lw.buf), "\r"))
		lw.buf = nil
	}
}

type exchange struct {
	dir      string
	immOnly  bool
	updates  BBSExchangeWatcher
	config   *Config
	conn     *jnos.Conn
	areas    sets.Set[string]
	area     string // XND@XSC
	listto   string // XND
	mailbox  string // ALLXSC
	indexes  []int
	msgindex int
	tokill   int

	// Resources held during the exchange, released by cleanup.
	logf      *os.File
	monitorlw *lineWriter
	reserved  []string    // entries in connectionsRunning
	stopAbort func() bool // unregisters the abort-on-cancel hook
	failed    bool        // an error has been reported to updates
}

// bbsExchangeStep takes one step in the BBS exchange process.  If there's a
// message waiting to be sent, it will do that.  Otherwise, if there's a message
// to retrieve in the current area, it will do that.  Otherwise, if there's
// another BBS area to check, it will switch to that are look for messages
// there.  The function returns whether all BBS exchange steps have completed,
// and may return an error if something goes wrong.
func (e *exchange) step(ctx context.Context) (done bool, err error) {
	// First, is the context canceled, meaning the user aborted the
	// connection?
	// (The connection has already been aborted, so we don't try to talk
	// to the BBS, e.g. to kill a message we just read.)
	select {
	case <-ctx.Done():
		slog.Info("connection aborted by user")
		return false, errAborted
	default: // continue
	}
	// Is there a message waiting to be sent?
	if tosend, logident, err := e.messageToSend(); err != nil {
		return false, err
	} else if tosend != nil {
		return false, e.send(tosend, logident)
	}
	// Is there a message we just read that needs to be killed?
	if e.tokill != 0 {
		return false, e.killMessage()
	}
	// If we don't have a list of message indexes to retrieve, and we need
	// one, get it.
	if e.indexes == nil {
		if e.immOnly {
			return e.getImmediateIndexes()
		} else if e.area != "" {
			if err = e.getBulletinIndexes(); err != nil {
				return false, err
			} else {
				return false, nil
			}
		}
	}
	// If we have a non-empty list of indexes, read the next one.
	if len(e.indexes) != 0 {
		index := e.indexes[0]
		e.indexes = e.indexes[1:]
		return false, e.readIndex(index)
	}
	// If we have a nil list of indexes, it means we're in the main mailbox
	// without immOnly and we're just reading in order.
	if e.indexes == nil && e.msgindex < 9999 {
		e.msgindex++
		if err = e.readIndex(e.msgindex); err == errNoSuchIndex {
			e.msgindex = 9999
			return len(e.areas) == 0, nil
		}
		return false, err
	}
	// If we've gotten here, we've finished reading the current mailbox.
	// Switch to the next one.
	return e.switchArea()
}

func (e *exchange) messageToSend() (tosend *message.DraftMessage, logident int, err error) {
	err = Read(e.dir, func(i *Incident) (err error) {
		for _, le := range i.Log {
			if le.Status != StatusQueued {
				continue
			}
			if e.immOnly && le.Flags&FImmediate == 0 {
				continue
			}
			if msg, err := i.GetMessageFromLogEntry(le); msg != nil {
				tosend = msg.(*message.DraftMessage)
				logident = le.Ident
				break
			} else if err != nil {
				return err
			}
		}
		return nil
	})
	return tosend, logident, err
}

func (e *exchange) send(dm *message.DraftMessage, logident int) (err error) {
	var le *LogEntry

	// First, update the OpDate and OpTime fields if any.
	now := time.Now()
	for f := range dm.Fields() {
		switch f.Common() {
		case field.COperatorDate:
			f.SetValue(dm, now.Format("01/02/2006"))
		case field.COperatorTime:
			f.SetValue(dm, now.Format("15:04"))
		}
	}
	var to []string
	if addrs, err := address.ParseList(dm.To()); err != nil || len(addrs) == 0 {
		slog.Error("invalid To", "to", dm.To(), "ident", logident)
		return errors.New("invalid or empty To: address list")
	} else {
		for _, addr := range addrs {
			to = append(to, addr.Address)
		}
	}
	if dm.MType == receipt.DeliveryReceipt {
		e.progress("Sending delivery receipt for %s...", dm.Body().(*receipt.DeliveryReceiptBody).LocalMessageID())
	} else {
		e.progress("Sending message %s...", dm.Subject().SubjectMessageID())
	}
	if dm.Bulletin() {
		err = e.conn.SendBulletin(dm.Subject().EncodedSubject(), dm.Payload().Encode(), to[0])
	} else {
		err = e.conn.Send(dm.Subject().EncodedSubject(), dm.Payload().Encode(), to...)
	}
	if err != nil {
		return err
	}

	if err = Write(e.dir, func(i *Incident) error {
		if le = i.GetLogEntryByIdent(logident); le == nil || le.Status == StatusDeleted {
			// The user deleted the message while we were sending
			// it (minutes, over RF).  It has been sent, but there's
			// nothing left to mark as sent (a deleted log entry has
			// lost its subject, addresses, and flags), and we don't
			// resurrect what the user deleted.  It won't be sent
			// again, since it's no longer in the queue.
			slog.Warn("sent message was deleted during send", "ident", logident, "subject", dm.Subject().EncodedSubject())
			le = nil
			return nil
		}
		return i.MarkMessageSent(dm, le)
	}); err != nil {
		return err
	}
	if le != nil {
		e.updates.LogEntry(le)
	} else {
		e.progress("Message %s was sent, but it had been deleted, so it is not logged.", dm.Subject().SubjectMessageID())
	}
	return nil
}

func (e *exchange) killMessage() (err error) {
	e.progress("Deleting message %d from BBS...", e.tokill)
	tokill := e.tokill
	e.tokill = 0
	return e.conn.Kill(tokill)
}

var immediateRE = regexp.MustCompile(`^[^ _]{5,10}_I_`)

func (e *exchange) getImmediateIndexes() (done bool, err error) {
	var list *jnos.MessageList

	if list, err = e.conn.List(""); err != nil || list == nil {
		return true, err
	}
	for _, item := range list.Messages {
		if immediateRE.MatchString(item.SubjectPrefix) {
			e.indexes = append(e.indexes, item.Number)
		}
	}
	return e.indexes == nil, nil
}

func (e *exchange) getBulletinIndexes() (err error) {
	var list *jnos.MessageList

	if e.listto != "" {
		e.progress("Listing bulletins to %s...", e.area)
	} else {
		e.progress("Listing bulletins in %s...", e.area)
	}
	if list, err = e.conn.List(e.listto); err != nil {
		return err
	}
	e.indexes = []int{}
	if list == nil {
		return nil
	}
	return Read(e.dir, func(i *Incident) error {
	LIST:
		for _, lm := range list.Messages {
			for _, le := range i.Log {
				if le.Status != StatusReceived || le.Flags&FBulletin == 0 {
					continue
				}
				subject := le.Subject
				if len(subject) > 35 {
					subject = subject[:35]
				}
				subject = strings.TrimSpace(subject)
				if subject == lm.SubjectPrefix {
					continue LIST
				}
			}
			e.indexes = append(e.indexes, lm.Number)
		}
		return nil
	})
}

var errNoSuchIndex = errors.New("no message exists with the requested index")

func (e *exchange) readIndex(index int) (err error) {
	var (
		raw string
		msg *message.JustReceivedMessage
		le  *LogEntry
	)

	if e.mailbox != "" {
		e.progress("Reading message %d in %s...", index, e.mailbox)
	} else {
		e.progress("Reading message %d...", index)
	}
	if raw, err = e.conn.Read(index); err != nil {
		return err
	} else if raw == "" {
		return errNoSuchIndex
	}
	if msg, err = message.NewJustReceivedMessage(raw, e.config.ConnectBBS, e.area); err != nil {
		return err
	}
	err = Write(e.dir, func(i *Incident) error {
		var dr *message.DraftMessage

		if dr, le, err = i.ReceiveMessage(msg); err != nil {
			return err
		}
		if dr != nil {
			if _, err := i.AddDraftMessage(dr); err != nil {
				return fmt.Errorf("queueing delivery receipt: %s", err)
			}
		}
		return nil
	})
	if err == nil && e.area == "" {
		e.tokill = index
	}
	if err == nil {
		e.updates.LogEntry(le)
	}
	return err
}

func (e *exchange) switchArea() (done bool, err error) {
	// If we're currently in an area, we need to mark it as having been
	// checked.
	if e.area != "" {
		if err = Write(e.dir, func(i *Incident) error {
			i.BulletinAreaChecked(e.area, time.Now())
			return nil
		}); err != nil {
			return false, err
		}
	}
	// If there are no more areas to check, we're done.
	if e.area, _ = e.areas.PopAny(); e.area == "" {
		return true, nil
	}
	// Switch areas.
	e.indexes = nil
	var found bool
	if e.listto, e.mailbox, found = strings.Cut(strings.ToUpper(e.area), "@"); !found {
		e.listto, e.mailbox = "", e.listto
	}
	if !strings.HasPrefix(e.mailbox, "ALL") && (!strings.HasPrefix(e.mailbox, "XSC") || e.mailbox == "XSC") {
		e.mailbox = "ALL" + e.mailbox
	}
	e.progress("Switching to %s...", e.mailbox)
	return false, e.conn.SetArea(e.mailbox)
}

func (e *exchange) progress(f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	slog.Info(msg)
	e.updates.Progress(msg)
}
