package cmd

import (
	"fmt"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/spf13/pflag"
)

const (
	manualDRSlug = `Generate a draft delivery receipt for a received message`
	manualDRHelp = `
usage: packet manual dr      «message-id»
       packet manual receipt «message-id»

The "packet manual dr" (or "receipt") command generates a draft delivery receipt message for the named message.  That message must be a received message and there must not already be a delivery receipt for it.  This command is useful when automatic delivery receipts are turned off for manual operations, but a delivery receipt is needed for a specific message.  Once the draft delivery receipt is generated, it can be sent the same as any other outgoing message (see the "packet manual send" command).
`
)

func cmdManualDR(args []string) (err error) {
	var (
		msg  message.Message
		le   *incident.LogEntry
		dr   *message.DraftMessage
		drle *incident.LogEntry
		c    = cio.Open()
	)
	flags := pflag.NewFlagSet("m-dr", pflag.ContinueOnError)
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdManualHelp([]string{"dr"})
	} else if err != nil {
		c.Error(err)
		return usage(manualDRHelp)
	}
	if len(args) != 1 {
		return usage(manualDRHelp)
	}
	registerForms()
	if err = incWrite(false, func(i *incident.Incident) error {
		if msg, le, err = matchMessage(i, args[0], MMMessageOnly); err != nil {
			return err
		}
		if dr, err = i.MakeDeliveryReceipt(msg, le); err != nil {
			return err
		}
		dr.SetReadyToSend(true)
		if drle, err = i.AddDraftMessage(dr); err != nil {
			return err
		}
		return nil
	}); err != nil && err != errNoChange {
		return err
	}
	if c.OutputIsTerm {
		c.Confirm(`Delivery receipt queued.  Use "packet manual send #%d" to send it.`, drle.Ident)
	} else {
		fmt.Printf("#%d\n", drle.Ident)
	}
	return nil
}
