package cmd

const (
	markSlug = `Change flags on a message or log entry`
	markHelp = `
usage: packet mark ⇥[-f] { «message-id» | «log-entry» } [not] «flag»...
       packet mark ⇥«message-id» sent
  -f, --force  ⇥Apply the change even if the message is invalid

The "packet mark" command sets flags on a message or log entry.  The flags are:
  - ⇥"dr" (or "d") means that a delivery receipt is expected for the message and has not yet been received from the recipient identified in this log entry.  The message must be a sent message.
  - ⇥"followup" (or "f") means that the log entry needs followup.
  - ⇥"ready" (or "r") means the message is ready to send at the next BBS connection.  The message must be an unsent outgoing message.
If the keyword "not" is used, the command clears the flags instead of setting them.

"packet mark ... ready" will not mark a message with validation errors ready unless the --force (or -f) flag is given.

"packet mark ... sent" is available only when the incident connection type is Manual, and is used to indicate that the operator has manually sent the identified message (usually after a "packet manual send" command to get the JNOS command with which to send it).  The message must be an unsent outgoing message, and will be marked as having been sent at the current time.  This is not reversible.
`
)

func cmdMark(args []string) (err error) {
	panic("not implemented")
}
