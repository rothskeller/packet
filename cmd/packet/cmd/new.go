package cmd

const (
	newSlug = `Create a new outgoing message or log entry`
	newHelp = `
usage: packet new ⇥[-flags] [«new-message-type»] [«new-message-id»]
  -c, --copy «message-id»    ⇥Send a copy of an existing message
  -l, --log                  ⇥Create a log entry not associated with a message
  -r, --reply «message-id»   ⇥Send a reply to an existing received message
  -x, --resend «message-id»  ⇥Correct and resend an already-sent message

The "packet new" (or "n") command creates a new outgoing message or log entry.  If standard input and output are terminals, the new message or log entry will be opened for editing (see "packet help edit" for details).  Otherwise, the local message ID for the new message (or log entry number for the new log entry) will be printed to standard output, and subsequent "set" commands can be used to populate it.

When the --copy (or -c) flag is given, the new message will be an exact copy of the named source message except for being given a new local message ID.  A «new-message-type» argument is not allowed.

When the --log (or -l) flag is given, a new log entry is created that is not associated with any message.  No arguments are allowed.

When the --reply (or -r) flag is given, the new message will have the same handling order, subject line, and body as the named source message, which must be a received message.  Its "To" address will be set to the "From" address of the source message.  The To and From ICS Position, Location, Name, and Contact fields, if any, will be swapped.  The new message will have the same message type as the source message unless a «new-message-type» is given on the command line.  If the message type has a "Reference" field, it will be filled with the source message's origin message ID.

When the --resend (or -x) flag is given, the new message will be identical to the source message except with the suffix of the message number changed to "R" (or the next available letter if "R" has already been used).  The source message must be an already-sent message.

When no flags are given, an empty message of «new-message-type» is created.  «new-message-type» must be the tag name or create key of one of the supported message types; if omitted, a plain text message is created.  Use "packet forms list" to get a list of supported forms.

If a «new-message-id» is provided on the command line, the new message is created with that local message ID.  The sequence number in it will be incremented as needed to make it unique.  The «new-message-id» may be just an integer, in which case the message number and prefix in the incident configuration are used (see "packet help config").  If no «new-message-id» is given, one will be automatically assigned based on the incident configuration.
`
)

func cmdNew(args []string) (err error) {
	panic("not implemented")
}
