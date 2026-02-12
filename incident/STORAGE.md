# Incident Storage Details

Each packet incident is stored in its own directory, with the contents
described here.

## Lock File

The `.packet.lock` file controls access to everything else. Packet software
does not change any file in the directory except while holding a write lock on
this file. Packet software does not read any file in the directory except
while holding a read lock on this file. (Of course, non-packet software like
browsers and PDF readers won't honor this, and will read files without getting
a lock.)

The `.packet.lock` file contains a sequence number that is incremented any time
the incident files are changed in any semantic way. (Generating a PDF for an
existing message doesn't count, for example.) It is written with the
equivalent of "%08d" so that it can be overwritten without needing a truncate
operation.

Software that needs to watch for incident changes will watch for writes to this
file. (Note: I've confirmed empirically that merely opening or locking a file
for write does not trigger an fsnotify WRITE event; only an actual write does.
It's perhaps worth noting that on Windows, the fsnotify WRITE event does not
occur until the file is closed after the write.)

## State File

The incident state is stored in `packet.json`, which as noted above is only
read or written when a lock is held on the lock file. It is a JSON file
containing the incident configuration, the bulletin check timestamps, and a log
of all messages. For coding simplicity, it is always read and written as a
complete unit. (I may look into incremental writes someday if performance
becomes an issue.)

## Message Files

Each packet message (draft, sent, or received) is stored in a text file in the
incident directory. The filename is derivable from the log entry, as follows:

- For an unsent message, the filename is ".unsent${IDENT}.txt", with the ident
  of the corresponding log entry formatted as %04d.
- For a receipt, the filename is ".receipt${IDENT}.txt, with the ident of the
  corresponding log entry formatted as %04d.
- For all other messages, the filename is "${LMI}.txt", with the local message
  ID of the message.

The message files contain the packet message in RFC-5322 (email) format.
Received messages can be distinguished by having a Received: header. Outgoing
messages will have a Date: header if they have been sent and not if they
haven't.

Several non-standard headers are added to message files. X-Packet-Receipt is
added to a sent message with details from a received receipt for it.
X-Packet-Bulletin is added to a sent message to indicate that it is a bulletin.
X-Packet-Queued is added to an unsent outgoing message to indicate that it is
ready to send

## PDF Files

Sent and received messages, other than receipts, are automatically rendered in
PDF format for viewing and printing. The PDF version of a message has the same
basename as the message file (i.e., its local message ID) with a `.pdf`
extension.

## Message Links

If a message has discernable remote message IDs, symbolic links are created from
"${RMI}.txt" to "${LMI}.txt" and "${RMI}.pdf" to "${LMI}.pdf".  If there is a
conflict between an RMI and any LMI, the LMI takes precedence and the symbolic
links that would conflict are not created.  Errors creating the symbolic links
are ignored (such as on Windows systems without base user permissions to create
symbolic links).

## ICS-309 Log

On demand, the packet software can generate an ICS-309 communications log in
both CSV and PDF formats. This is stored in `ICS309.csv` and `ICS309.pdf`
respectively. These files are removed automatically when the incident changes,
to ensure they are never stale. Note the lack of hyphen in the filenames, so
that they can't conflict with a message number.

## Archive Files

On demand, the packet software can generate an archive in both text and PDF
formats. This is stored in `all.txt` and `all.pdf`, respectively. These files
are removed automatically when the incident changes, to ensure they are never
stale.
