package cmd

import (
	"slices"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/pseudomsg"
	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/incident"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/msgifc"
	"github.com/spf13/pflag"
)

const (
	editSlug = `Edit an unsent message, log entry, or incident configuration`
	editHelp = `
usage: packet edit ⇥[flags] { «message-id» | «log-entry» | configuration } [«field»]
  -e, --errors   ⇥edit only the fields that have errors
  -m, --message  ⇥edit the message, not the log entry

The "edit" (or "e") command edits an unsent message, a log entry, or the incident configuration.  It presents each field in turn and allows that field's value to be changed.  Note that the "edit" command cannot be used in scripted mode (see "packet help script").

To edit the incident configuration, give the word "configuration" or any abbreviation of it.

The "edit" command normally starts with the first field of the message, log entry, or configuration (or the first that has an error, if --errors or -e is used).  If a «field» is specified, editing begins with that field instead.

Usage of the editor depends on the capabilities of the standard output device (e.g., the terminal).  If it is fully capable, the following keys can be used:
    Ctrl-A, Home    ⇥move cursor to beginning of line (*)
    Ctrl-B, ←       ⇥move cursor to the left (*)(+)
    Ctrl-C          ⇥abort the edit and do not save any changes
    Ctrl-D, Delete  ⇥delete the character under the cursor
    Ctrl-E, End     ⇥move cursor to end of line (*)
    Ctrl-F, →       ⇥move cursor to the right (*)(+)
    Ctrl-H, Backsp  ⇥delete the character before the cursor
    Ctrl-I, Tab     ⇥save this field and move to the next field
    Shift-Tab       ⇥save this field and move to the previous field
    Ctrl-K          ⇥delete the remainder of the current line
    Ctrl-L          ⇥redraw the editor (in case of screen corruption)
    Ctrl-M, Enter   ⇥multi-line fields:  enter a newline
                    ⇥single-line fields: save field and move to next field
    Ctrl-N, ↓       ⇥move cursor down one line (*)
    Ctrl-P, ↑       ⇥move cursor up one line (*)
    Ctrl-U          ⇥delete the entire contents of the field
    Ctrl-V + Enter  ⇥enter a newline in a normally single-line field
    ESC             ⇥save this field and exit the editor
    F1              ⇥display online help for the field
    (*) with Shift, extend the selection in the direction of movement
    (+) arrows with Ctrl move by words instead of characters
Some fields have a discrete set of possible or recommended values.  For those fields, the editor will show the set of values and allow you to select from among them using the arrow keys.  Or, if you prefer to type, the editor will autocomplete your entry from that set.

If you enter an invalid value for a field, an appropriate error will be shown and you will be asked to enter that field again.  When editing messages, if you hit Enter on the value to confirm it, the value will be kept despite the error.  (When editing log entries or incident configurations, errors must be corrected.)

When finished editing a message, if the message is fully valid and not already marked ready to send, the editor will ask whether to mark it.  If the message is already marked ready to send but is not valid, it will be unmarked.  (To forcibly mark a message ready to send, use the "packet mark ready" command.)
`
)

func cmdEdit(args []string) (err error) {
	var (
		errorsOnly bool
		asMessage  bool
		msg        message.Message
		c          = cio.Open()
	)
	if !c.InputIsTerm || !c.OutputIsTerm {
		return errors.New("Editing is supported only when stdin/stdout is a terminal.")
	}
	flags := pflag.NewFlagSet("edit", pflag.ContinueOnError)
	flags.BoolVarP(&errorsOnly, "errors", "e", false, "Edit only fields with errors")
	flags.BoolVarP(&asMessage, "message", "m", false, "Change message rather than log entry")
	flags.Usage = func() {} // we do our own
	if err = flags.Parse(args); err == pflag.ErrHelp {
		return cmdHelp([]string{"edit"})
	} else if err != nil {
		c.Error(err)
		return usage(editHelp)
	}
	if n := flags.NArg(); n < 1 || n > 2 {
		return usage(editHelp)
	}
	registerForms()
	if err = incWrite(strings.HasPrefix("configuration", flags.Arg(0)), func(i *incident.Incident) error {
		var (
			mmf   matchMessageFlag
			entry *incident.LogEntry
			fld   field.Field
		)
		// Get the message to be changed, and starting field if any, and
		// validate them.
		if asMessage {
			mmf |= MMMessageOnly
		}
		if msg, entry, err = matchMessage(i, flags.Arg(0), mmf); err != nil {
			return err
		}
		switch msg.(type) {
		case *pseudomsg.ConfigMessage, *pseudomsg.LogEntryMessage, *message.DraftMessage:
			// OK
		case *message.JustReceivedMessage, *message.ReceivedMessage:
			return errors.NewF("Received messages cannot be edited.")
		case *message.SentMessage:
			return errors.NewF("Sent messages cannot be edited.")
		}
		if flags.NArg() == 2 {
			if fld, err = matchField(msg, flags.Arg(1), c.OutputIsTerm); err != nil {
				return err
			}
			if !fld.Editable(msg, true) {
				return errors.NewF("Field %q is not editable.", fld.Label())
			}
		}
		return doEdit(c, i, entry, msg, fld, errorsOnly)
	}); err != nil {
		return err
	}
	return nil
}

// doEdit is the common code between edit and new.
func doEdit(c *cio.CIO, i *incident.Incident, entry *incident.LogEntry, msg message.Message, start field.Field, errorsOnly bool) (err error) {
	var (
		fields     []field.Field
		fld        field.Field
		wasQueued  bool
		startSeen  = start == nil
		labelWidth = 14 // "Ready to Send?"
	)
	// Build the list of fields to be edited.
	for f := range msg.Fields() {
		if f.EditHelp() != "" {
			fields = append(fields, f)
			labelWidth = max(labelWidth, len(f.Label()))
			if f == start {
				startSeen = true
			}
			if fld == nil && startSeen && f.Editable(msg, f == start) && (!errorsOnly || f.Validate(msg, f, msgifc.VPacket) != nil) {
				fld = f
			}
		}
	}
	if fld == nil {
		// Only happens if they said errorsOnly and we didn't find any.
		return errors.New("There are no fields with errors.")
	}
	// Is the message ready to send before editing?
	if msg, ok := msg.(*message.DraftMessage); ok {
		wasQueued = msg.ReadyToSend()
		fields = append(fields, newReadyToSendField())
	}
	if start == nil && !errorsOnly {
		c.StartEdit() // Editor instructions
	}
LOOP: // Run the editor loop.
	for {
		var (
			value      string
			valueWidth int
			choices    []string
			result     cio.EditResult
			newvalue   string
			first      = true
		)
		value = fld.ToHuman(fld.Value(msg))
		valueWidth, _ = fld.EditSize()
		for _, c := range fld.Choices(msg) {
			choices = append(choices, c.Human)
		}
		for {
			if result, newvalue, err = c.EditField(fld.Label(), labelWidth, value, valueWidth, choices, fld.EditHelp(),
				fld.EditHint(), fld.Multiline(), fld.Obscured(),
				func(s string) string { return fld.ToHuman(fld.FromHuman(s)) }); err != nil {
				return err
			}
			fld.SetValue(msg, fld.FromHuman(newvalue))
			if err = fld.Validate(msg, fld, msgifc.VPacket); err == nil {
				break
			}
			if !first && newvalue == value && !errors.IsType[pseudomsg.NoBypassValidationError](err) {
				err = nil
				break
			}
			first = false
			c.Error(err)
		}
		switch result {
		case cio.ResultDone:
			break LOOP
		case cio.ResultNext:
			idx := slices.Index(fields, fld) + 1
			for idx < len(fields) {
				fld = fields[idx]
				if fld.Editable(msg, false) && (!errorsOnly || fld.Validate(msg, fld, msgifc.VPacket) != nil) {
					break
				}
				idx++
			}
			if idx >= len(fields) {
				break LOOP
			}
		case cio.ResultPrevious:
			idx := slices.Index(fields, fld) - 1
			for idx >= 0 {
				fld = fields[idx]
				if fld.Editable(msg, false) && (!errorsOnly || fld.Validate(msg, fld, msgifc.VPacket) != nil) {
					break
				}
				idx--
			}
			if idx < 0 {
				break LOOP
			}
		default:
			panic("unknown result code")
		}
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
	// Notify the user if we took the message out of the queue.
	if wasQueued && !msg.(*message.DraftMessage).ReadyToSend() {
		c.Confirm("NOTE: This message has invalid fields and is no longer marked ready to send.")
	}
	return nil
}

func newReadyToSendField() (f field.Field) {
	f = field.NewField("", "Ready to Send?").
		AllowedValues("Yes", "No").
		EditHelp(`This indicates whether the message should be sent during the next BBS connection.`).
		EditableWhen(func(msg msgifc.Message, _ bool) bool {
			if message.ValidateMessage(msg, msgifc.VPacket) == nil {
				return true
			} else {
				msg.(*message.DraftMessage).SetReadyToSend(false)
				return false
			}
		}).
		ValueFunc(func(_ msgifc.Message) string { return "Yes" }).
		SetValueFunc(func(msg msgifc.Message, s string) {
			msg.(*message.DraftMessage).SetReadyToSend(f.FromHuman(s) == "Yes")
		}).
		MakeField()
	return f
}
