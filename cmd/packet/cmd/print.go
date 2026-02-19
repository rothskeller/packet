package cmd

const (
	printSlug = `Sends a message to the system printer`
	printHelp = `
usage: packet print «message-id»

The "packet print" command sends the identified message to the system default printer.
`
)

func cmdPrint(args []string) (err error) {
	panic("not implemented")
}
