package cmd

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
	panic("not implemented")
}
