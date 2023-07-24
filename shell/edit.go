package shell

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/msgedit"
)

// cmdEdit implements the edit command.
func cmdEdit(args []string) bool {
	var (
		lmi  string
		env  *envelope.Envelope
		msg  message.Message
		emsg message.IEdit
		err  error
	)
	// Check arguments
	if len(args) != 1 {
		io.WriteString(os.Stderr, "usage: packet edit <message-id>\n")
		return false
	}
	switch lmis := expandMessageID(args[0], false); len(lmis) {
	case 0:
		fmt.Fprintf(os.Stderr, "ERROR: no such message %q\n", args[0])
		return false
	case 1:
		lmi = lmis[0]
	default:
		fmt.Fprintf(os.Stderr, "ERROR: %q is ambiguous (%s)\n", args[0], strings.Join(lmis, ", "))
		return false
	}
	// Read and parse file to be edited.
	if env, msg, err = readMessage(lmi); err != nil {
		return false
	}
	if m, ok := msg.(message.IEdit); ok {
		emsg = m
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: editing is not supported for %ss\n", msg.Type().Name)
		return false
	}
	if env.IsFinal() {
		if env.IsReceived() {
			io.WriteString(os.Stderr, "ERROR: cannot edit a received message\n")
		} else {
			io.WriteString(os.Stderr, "ERROR: cannot edit a message that has been sent\n")
		}
		return false
	}
	// Start the editor.
	return editMessage(lmi, env, emsg)
}

// helpEdit prints the help message for the edit command.
func helpEdit() {
	io.WriteString(os.Stdout, `The "edit" (or "e") command edits an unsent message.
    usage: packet edit <message-id>
<message-id> is the local message ID of the message to edit.  It can be just
    the numeric part if that is unambiguous.
Editing is the default action for unsent messages, so the "edit" keyword can
be omitted.
    The message editor is a full-screen, scrollable two-column list of fields,
with field names on the left and field values on the right.  To move the
cursor around in the list, use the mouse, or the arrow, Enter, Tab, Shift-Tab,
PgUp, PgDn, Home, and End keys.  Press F1 for help on the currently selected
field.  When finished editing, press either Esc or F10.  Esc leaves the
message in draft state, not queued to be sent.  F10 leaves the message queued
to be sent.
`)
}

// editMessage starts a message editor for a message.
func editMessage(lmi string, env *envelope.Envelope, msg message.IEdit) bool {
	var err error

	// Run the editor.
	msgedit.RunEditor(lmi, env, msg)
	// Check for a change to the LMI.
	newlmi := msg.GetOriginID()
	if newlmi != lmi {
		if unique := getNextMessageID(newlmi); unique != newlmi {
			fmt.Fprintf(os.Stderr, "WARNING: message ID %s is already in use; using %s instead\n", newlmi, unique)
			newlmi = unique
			msg.SetOriginID(newlmi)
		}
		if lmi != "" {
			os.Remove(lmi + ".txt")
			os.Remove(lmi + ".pdf")
		}
		lmi = newlmi
	}
	// Save the resulting message.
	if err = saveMessage(lmi, env, msg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return false
	}
	return true
}
