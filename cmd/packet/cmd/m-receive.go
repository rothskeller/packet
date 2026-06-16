package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/spf13/pflag"
)

const (
	manualReceiveSlug = `Log and record a manually received message`
	manualReceiveHelp = `
usage: packet manual receive [ < «filename» ]

The "packet manual receive" (or "man r") command logs and records a manually received message.  The received message is read from standard input up to end of file or a "/EX" line.  There are two common usage models:
  - ⇥The operator saves the received message into a file, then redirects that file into the "packet manual receive" command.
  - ⇥The operator copies the received message from a terminal emulator, then runs "packet manual receive" and pastes it.  After pasting, the operator ends the input with a "/EX" line or the OS-dependent end-of-file character (^Z on Windows, ^D elsewhere).

The resulting message will be added to the incident as a received message and a log entry for it will be added to the log.  If automatic delivery receipts are enabled in the incident configuration, a draft delivery receipt will be created for it, which the operator can manually send when desired.  (If automatic receipts are not enabled, the operator can generate one manually with the "packet manual dr" command.)
`
)

func cmdManualReceive(args []string) (err error) {
	var (
		sb   strings.Builder
		scan *bufio.Scanner
		dr   *message.DraftMessage
		drle *incident.LogEntry
		msg  *message.JustReceivedMessage
		le   *incident.LogEntry
		c    = cio.Open()
	)
	flags := pflag.NewFlagSet("m-receive", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdManualHelp([]string{"receive"})
	} else if err != nil {
		c.Error(err)
		return usage(manualReceiveHelp)
	}
	if flags.NArg() != 0 {
		return usage(manualReceiveHelp)
	}
	registerForms()
	scan = bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		line := scan.Text()
		if line == "/EX" {
			break
		}
		fmt.Fprintln(&sb, line)
	}
	if err = scan.Err(); err != nil {
		return err
	}
	err = incWrite(false, func(i *incident.Incident) error {
		if err = requiredConfig(i, "OpCall", "OpName", "Rx Message ID", "ConnectBBS"); err != nil {
			return err
		}
		if msg, err = message.NewJustReceivedMessage(sb.String(), i.Config.ConnectBBS, ""); err != nil {
			return err
		}
		if dr, le, err = i.ReceiveMessage(msg); err != nil {
			return err
		}
		if dr != nil {
			if drle, err = i.AddDraftMessage(dr); err != nil {
				return errors.NewF("Unable to queue delivery receipt: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if c.OutputIsTerm {
		c.Confirm("Message received as %s.", le.LocalMsgID)
	} else {
		fmt.Println(le.LocalMsgID)
	}
	if drle != nil {
		if c.OutputIsTerm {
			c.Confirm(`Delivery receipt queued; send with "packet manual send #%d".`, drle.Ident)
		} else {
			fmt.Printf("#%d\n", drle.Ident)
		}
	}
	return nil
}
