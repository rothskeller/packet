package cmd

const (
	manualDRSlug = `Generate a draft delivery receipt for a received message`
	manualDRHelp = `
usage: packet manual dr      «message-id»
       packet manual receipt «message-id»

The "packet manual dr" (or "receipt") command generates a draft delivery receipt message for the named message.  That message must be a received message and there must not already be a delivery receipt for it.  This command is useful when automatic delivery receipts are turned off for manual operations, but a delivery receipt is needed for a specific message.  Once the draft delivery receipt is generated, it can be sent the same as any other outgoing message (see the "packet manual send" command).
`
)

func cmdManualDR(args []string) (err error) {
	panic("not implemented")
}
