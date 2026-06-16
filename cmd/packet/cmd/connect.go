package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/spf13/pflag"
)

const (
	connectSlug = `Connect to the BBS to send and/or receive messages`
	connectHelp = `
usage: packet connect [-i]
  -i, --immediate  ⇥Send and receive only immediate messages

The "connect" (or "c") command makes a connection to the BBS and sends and/or receives messages.  With the --immediate (or -i) flag, only immediate messages are sent and/or received.

When receiving messages without the --immediate flag, any scheduled bulletin checks are performed as well.

The "connect" command lists all messages sent and received.  Run "packet help list" for details of the output format.
`
)

var ErrInterrupted = errors.New("connection interrupted by Ctrl-C")

func cmdConnect(args []string) (err error) {
	var (
		immediate bool
		dir       string
		myID      string
		ctx       context.Context
		cancel    func()
		sigintch  chan os.Signal
		updates   *connectUpdates
		finished  chan struct{}
		c         = cio.Open()
	)
	flags := pflag.NewFlagSet("connect", pflag.ContinueOnError)
	flags.BoolVarP(&immediate, "immediate", "i", false, "Send and receive only immediate messages")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"connect"})
	} else if err != nil {
		c.Error(err)
		return usage(connectHelp)
	}
	if flags.NArg() != 0 {
		return usage(connectHelp)
	}
	if err = incWrite(true, func(i *incident.Incident) error {
		if err = requiredConfig(i, "OpCall", "OpName", "Rx Message ID", "Connect*"); err != nil {
			return err
		}
		dir = i.Dir
		myID = i.Config.ActiveCall()
		return nil
	}); err != nil {
		return err
	}
	// We want to stop the connection if we get a control-C, so we need to
	// run it in a context we can cancel, and then we want to trap
	// interrupts.
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	sigintch = make(chan os.Signal, 10)
	signal.Notify(sigintch, os.Interrupt)
	defer drainSigInt(sigintch)
	// When the connection finishes, we need to break out of the select
	// loop, so we'll use a sentinel channel for that.
	finished = make(chan struct{})
	updates = &connectUpdates{c: c, finished: finished}
	// Run the connection.
	go incident.BBSExchange(ctx, dir, immediate, updates)
LOOP:
	for {
		select {
		case <-sigintch:
			cancel()
		case <-finished:
			break LOOP
		}
	}
	if len(updates.log) != 0 {
		c.EmitLogList(updates.log, false, false, true, myID)
	} else if updates.err == "" {
		c.Confirm("No messages sent or received.")
	}
	if updates.err != "" {
		return errors.New(updates.err)
	}
	return nil
}

type connectUpdates struct {
	c        *cio.CIO
	finished chan struct{}
	log      []*incident.LogEntry
	err      string
}

func (cu *connectUpdates) Progress(s string) { cu.c.Status("%s", s) }

func (cu *connectUpdates) LogEntry(le *incident.LogEntry) { cu.log = append(cu.log, le) }

func (cu *connectUpdates) Error(e string) {
	cu.c.ErrorF("%s", e)
	if cu.err == "" {
		cu.err = e
	}
}

func (cu *connectUpdates) Finished() {
	cu.c.Status("")
	close(cu.finished)
}

// drainSigInt stops trapping interrupt signals, drains the signal channel
// buffer, and closes the channel.
func drainSigInt(sigintch chan os.Signal) {
	var done bool
	signal.Reset(os.Interrupt)
	for !done {
		select {
		case <-sigintch:
			// no op
		default:
			done = true
		}
	}
	close(sigintch)
}
