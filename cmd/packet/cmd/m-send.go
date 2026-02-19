package cmd

const (
	manualSendSlug = `Display the JNOS command to send a message`
	manualSendHelp = `
usage: packet manual send «message-id»

The "packet manual send" (or "man s") command displays the JNOS command to send the message identified by «message-id», which must be an unsent outgoing message.

If standard input and output are terminals, the command will ask whether to mark the message as having been sent.  Otherwise, it can be so marked using the "packet mark sent" command.
`
)

func cmdManualSend(args []string) (err error) {
	panic("not implemented")
}
