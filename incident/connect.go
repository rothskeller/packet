package incident

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/rothskeller/packet/jnos"
	"github.com/rothskeller/packet/jnos/kpc3plus"
	"github.com/rothskeller/packet/jnos/telnet"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/receipt"
	"k8s.io/apimachinery/pkg/util/sets"
)

var connectionsRunning = sets.New[string]()
var connectionsRunningMutex sync.Mutex

// BBSExchange connects to a BBS and exchanges messages with it for the incident
// in the named directory (which should not be locked when the call is issued).
// The connection will be aborted if the supplied context is canceled.  If
// immOnly is true, only IMMEDIATE messages are transferred (and no bulletin
// checks are performed).  Progress updates and results are sent to the supplied
// watcher.  BBSExchange blocks, and should be run in a goroutine.
func BBSExchange(ctx context.Context, dir string, immOnly bool, updates BBSExchangeWatcher) {
	var (
		logfname string
		logf     *os.File
		e        exchange
		done     bool
		err      error
	)
	e.dir, e.immOnly, e.updates = dir, immOnly, updates
	// Read the incident configuration.
	if err := Read(dir, func(i *Incident) error {
		e.config = i.Config
		if !immOnly {
			e.areas = sets.New(i.BulletinAreasToCheck()...)
		}
		return nil
	}); err != nil {
		updates.Error(err.Error())
		updates.Finished()
		return
	}
	// Make sure we don't already have a conflicting connection running.
	connTo := e.config.ActiveCall() + "@" + e.config.ConnectBBS
	connectionsRunningMutex.Lock()
	if connectionsRunning.Has(connTo) {
		connectionsRunningMutex.Unlock()
		slog.Error("already connected", "to", connTo)
		updates.Error(fmt.Sprintf("there is already a connection to %s in progress", connTo))
		updates.Finished()
		return
	}
	if e.config.ConnectType == ConnectSerialTNC {
		if connectionsRunning.Has(e.config.SerialPort) {
			connectionsRunningMutex.Unlock()
			slog.Error("serial port in use", "port", e.config.SerialPort)
			updates.Error(fmt.Sprintf("can't connect using %s because it is in use", e.config.SerialPort))
			updates.Finished()
			return
		}
		connectionsRunning.Insert(e.config.SerialPort)
	}
	connectionsRunning.Insert(connTo)
	connectionsRunningMutex.Unlock()
	defer func() {
		connectionsRunningMutex.Lock()
		connectionsRunning.Delete(connTo)
		if e.config.ConnectType == ConnectSerialTNC {
			connectionsRunning.Delete(e.config.SerialPort)
		}
		connectionsRunningMutex.Unlock()
	}()
	// Open the log file.
	logfname = filepath.Join(dir, e.config.ConnectBBS+".log")
	if logf, err = os.OpenFile(logfname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err != nil {
		slog.Error("os.OpenFile", "f", logfname, "err", err)
		updates.Error(err.Error())
		updates.Finished()
		return
	}
	defer logf.Close()
	fmt.Fprintf(logf, "\n[[%s]]\n", time.Now().Format(time.RFC3339))
	// Establish the JNOS connection.
	switch e.config.ConnectType {
	case ConnectSerialTNC:
		e.progress("Connecting to %s@%s...", e.config.ActiveCall(), e.config.ConnectBBS)
		e.conn, err = kpc3plus.Connect(e.config.SerialPort, e.config.ConnectAddress, e.config.ActiveCall(), e.config.OpCall, logf)
	case ConnectTelnet:
		e.progress("Connecting to %s@%s...", e.config.TelnetUser, e.config.ConnectBBS)
		e.conn, err = telnet.Connect(e.config.ConnectAddress, e.config.TelnetUser, e.config.TelnetPassword, logf)
	default:
		slog.Error("invalid connect method", "method", e.config.ConnectType)
		err = fmt.Errorf("can't connect using method %s", e.config.ConnectType)
	}
	if err != nil {
		fmt.Fprintf(logf, "[[ERROR: %s]]\n", err)
		updates.Error(err.Error())
		updates.Finished()
		return
	}
	// Run the main connection loop.
	for !done {
		if done, err = e.step(ctx); err != nil {
			break
		}
	}
	if err != nil {
		updates.Error(err.Error())
	}
	// Close the connection.
	e.progress("Closing connection...")
	if err2 := e.conn.Close(); err2 != nil && err == nil {
		updates.Error(err2.Error())
	}
	updates.Finished()
}

// BBSExchangeWatcher is the interface required for the "updates" parameter to
// BBSExchange; it is an object to which progress of the BBS exchange is
// reported.
type BBSExchangeWatcher interface {
	// Progress is called with messages indicating the progress of the BBS
	// exchange.
	Progress(string)
	// Error is called if an error occurs during the BBS exchange.
	Error(string)
	// Finished is called when the BBS exchange has finished, whether or not
	// it was successful.  (It was successful if Error was never called.)
	Finished()
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
	select {
	case <-ctx.Done():
		slog.Info("connection aborted by user")
		if e.tokill != 0 {
			e.killMessage()
		}
		return false, errors.New("connection aborted by user")
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
	return Write(e.dir, func(i *Incident) error {
		return i.MarkMessageSent(dm, i.GetLogEntryByIdent(logident))
	})
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

	if list, err = e.conn.List(""); err != nil {
		return false, err
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

		if dr, err = i.ReceiveMessage(msg); err != nil {
			return err
		}
		if dr != nil {
			if _, err := i.AddDraftMessage(dr, false); err != nil {
				return fmt.Errorf("queueing delivery receipt: %s", err)
			}
		}
		return nil
	})
	if e.area == "" {
		e.tokill = index
	}
	return nil
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
