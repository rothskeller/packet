package cmd

import (
	"fmt"
	"strings"

	"github.com/rothskeller/packet/v4/cmd/packet/cio"
	"github.com/rothskeller/packet/v4/errors"
	"github.com/rothskeller/packet/v4/incident"
	"github.com/rothskeller/packet/v4/message"
	"github.com/rothskeller/packet/v4/message/address"
	"github.com/rothskeller/packet/v4/message/msgifc"
	"github.com/spf13/pflag"
)

const (
	manualSendSlug = `Display the JNOS command to send a message`
	manualSendHelp = `
usage: packet manual send [-f] «message-id»
  -f, --force  Send the message even if it is invalid

The "packet manual send" (or "man s") command displays the JNOS command to send the message identified by «message-id», which must be an unsent outgoing message.  The message must pass validation checks unless --force (or -f) is specified.

If standard input and output are terminals, the command will ask whether to mark the message as having been sent.  Otherwise, it can be so marked using the "packet mark sent" command.
`
)

func cmdManualSend(args []string) (err error) {
	var (
		force bool
		ident string
		msg   message.Message
		le    *incident.LogEntry
		to    []string
		eb    string
		c     = cio.Open()
	)
	flags := pflag.NewFlagSet("m-send", pflag.ContinueOnError)
	flags.BoolVarP(&force, "force", "f", false, "Send the message even if invalid")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdManualHelp([]string{"send"})
	} else if err != nil {
		c.Error(err)
		return usage(manualSendHelp)
	}
	if flags.NArg() != 1 {
		return usage(manualSendHelp)
	}
	registerForms()
	if err = incRead(func(i *incident.Incident) error {
		if i.Config.TacCall != "" {
			ident = i.Config.OpCall
		}
		if msg, le, err = matchMessage(i, flags.Arg(0), MMMessageOnly); err != nil {
			return err
		}
		return nil
	}); err != nil && err != errNoChange {
		return err
	}
	switch le.Status {
	case incident.StatusReceived:
		return errors.NewF("Message %s is a received message.", le.LocalMsgID)
	case incident.StatusSent:
		return errors.NewF("Message %s has already been sent.", le.LocalMsgID)
	}
	if !force {
		if err = message.ValidateMessage(msg, msgifc.VPacket); err != nil {
			return err
		}
	}
	if addrs, err := address.ParseList(msg.To()); err != nil {
		return errors.NewF("The message To: address is not valid.  %s", err)
	} else if len(addrs) == 0 {
		return errors.New("The message has no To: address.")
	} else {
		for _, addr := range addrs {
			to = append(to, addr.Address)
		}
	}
	// Generate the actual send command.
	if msg.Bulletin() {
		fmt.Print("SB ")
	} else if len(to) > 1 {
		fmt.Print("SC ")
	} else {
		fmt.Print("SP ")
	}
	fmt.Println(to[0])
	if len(to) > 1 {
		fmt.Println(strings.Join(to[1:], ","))
	}
	fmt.Println(msg.Subject().EncodedSubject())
	eb = msg.Payload().Encode()
	fmt.Print(eb)
	if !strings.HasSuffix(eb, "\n") {
		fmt.Println()
	}
	fmt.Print("/EX\n")
	if ident != "" {
		fmt.Printf("# DE %s\n", ident)
	}
	if !c.InputIsTerm || !c.OutputIsTerm {
		return nil
	}
	if _, mark, err := c.EditField("Mark message sent?", 0, "Yes", 3, []string{"Yes", "No"},
		"Specify whether the message has been successfully sent.",
		"", false, false, nil); err != nil || !strings.HasPrefix(strings.ToUpper(mark), "Y") {
		return nil
	}
	return incWrite(false, func(i *incident.Incident) (err error) {
		le = i.GetLogEntryByIdent(le.Ident)
		return i.MarkMessageSent(msg.(*message.DraftMessage), le)
	})
}
