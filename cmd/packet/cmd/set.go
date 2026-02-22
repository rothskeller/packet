package cmd

import (
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/pseudomsg"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/spf13/pflag"
)

const (
	setSlug = `Set a value in a message, log entry, or incident configuration`
	setHelp = `
usage: packet set ⇥[flags] { «message-id» | «log-entry» | configuration } «field» [«value»]
  -f, --force    ⇥Sets the value even if it is invalid for the field
  -m, --message  ⇥Sets the value in the message, not the log entry

The "packet set" command sets the value of the named field in a message, a log entry, or the incident configuration.  If the first argument identifies a message, it must be for an unsent outgoing message, and the named field of that message is changed.  If it identifies a log entry, that log entry is changed.  If the first argument is "configuration" (or any abbreviation), the incident configuration is changed.

The second argument identifies the field of the item to be changed.

If a third argument is given, it is taken as the new value of the field; otherwise, the new value is read from the console or standard input.  It must be a valid value for the field unless the --force (or -f) flag is given.
`
)

func cmdSet(args []string) (err error) {
	var (
		force     bool
		asMessage bool
		msg       message.Message
		c         = cio.Open()
	)
	flags := pflag.NewFlagSet("set", pflag.ContinueOnError)
	flags.BoolVarP(&force, "force", "f", false, "Accept invalid value for the field")
	flags.BoolVarP(&asMessage, "message", "m", false, "Change message rather than log entry")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"set"})
	} else if err != nil {
		c.Error(err)
		return usage(setHelp)
	}
	if n := flags.NArg(); n < 2 {
		return usage(setHelp)
	}
	registerForms()
	if err = incWrite(strings.HasPrefix("configuration", flags.Arg(0)), func(i *incident.Incident) error {
		var (
			mmf      matchMessageFlag
			entry    *incident.LogEntry
			fld      field.Field
			problems map[field.Field]string
			value    string
			ispseudo bool
			newprob  error
		)
		// Get the message and field to be changed and validate them.
		if asMessage {
			mmf |= MMMessageOnly
		}
		if msg, entry, err = matchMessage(i, flags.Arg(0), mmf); err != nil {
			return err
		}
		switch msg.(type) {
		case *pseudomsg.ConfigMessage, *pseudomsg.LogEntryMessage:
			ispseudo = true
		}
		if fld, err = matchField(msg, flags.Arg(1), c.OutputIsTerm); err != nil {
			return err
		} else if _, ok := msg.(*message.JustReceivedMessage); ok {
			return errors.NewF("Received messages cannot be edited.")
		} else if _, ok := msg.(*message.ReceivedMessage); ok {
			return errors.NewF("Received messages cannot be edited.")
		} else if _, ok := msg.(*message.SentMessage); ok {
			return errors.NewF("Sent messages cannot be edited.")
		} else if !fld.Editable(msg, true) {
			return errors.NewF("Field %q is not editable.", fld.Label())
		}
		// Find out what problems already exist in the message.
		problems = make(map[field.Field]string)
		for f := range msg.Fields() {
			if err := f.Validate(msg, f, 0); err != nil {
				problems[f] = err.Error()
			}
		}
		// Get the value to set.
		if flags.NArg() > 2 {
			// Take the value from the command line and sanitize it.
			value = strings.Join(flags.Args()[2:], " ")
			value = strings.Map(func(r rune) rune { // ensure pure ASCII
				if (r >= ' ' && r <= '~') || r == '\n' {
					return r
				}
				if r == '\t' {
					return ' '
				}
				return -1
			}, value)
		} else {
			// Read the value from the console or stdin.
			var (
				choices []string
				width   int
			)
			for _, c := range fld.Choices(msg) {
				choices = append(choices, c.Human)
			}
			width, _ = fld.EditSize()
			_, value, err = c.EditField(
				fld.Label(), 0, fld.ToHuman(msg, fld.Value(msg)), width, choices, fld.EditHelp(), fld.EditHint(), fld.Multiline(), fld.Obscured(),
				func(s string) string { return fld.ToHuman(msg, fld.FromHuman(msg, s)) })
			if err != nil {
				return err
			}
		}
		// Either way, apply the value.
		fld.SetValue(msg, fld.FromHuman(msg, value))
		// Check for any new problems.
		for f := range msg.Fields() {
			if err := f.Validate(msg, f, 0); err != nil && err.Error() != problems[f] {
				newprob = errors.Join(newprob, err)
			}
		}
		if newprob != nil {
			if ispseudo || !force {
				return newprob
			}
			c.Error(newprob)
			c.Confirm("NOTE: applying the changes anyway since --force was used")
		}
		// Save the change
		switch msg := msg.(type) {
		case *message.DraftMessage:
			if err = i.UpdateDraftMessage(entry.Ident, msg); err != nil {
				return err
			}
		case *pseudomsg.LogEntryMessage:
			i.UpdateLogEntry(entry)
		case *pseudomsg.ConfigMessage:
			i.UpdateConfig(msg.Config)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
