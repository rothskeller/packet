# Packet Manager (pktmgr)

Package `pktmgr` manages collections of related packet messages.  These
collections are known as "incidents" because each collection contains the set of
messages from a single incident, represented on a single ICS-309 form.  The
package provides:

* storage of messages on disk
* assignment of local message IDs
* matching receipt messages up with the messages they describe
* generation of ICS-309 forms
* managing the state of outgoing messages (draft, queued, sent, receipted)
* generation of delivery receipts for received (non-bulletin) messages
* population of operator fields on forms

## Major Objects and Methods

The primary object type is an `Incident`, which is created by calling
`NewIncident`.  When a new incident is created, all messages are read from
disk into memory.

Incidents contain a collection of `MAndR` objects, which stands for
message-and-receipts.  Each one contains a single human-generated message (field
`M`), an optional delivery receipt (field `DR`), and an optional read receipt
(field `RR`).  Each `MAndR` has a local message ID in its `LMI` field.  To
retrieve a message by local ID, call the incident's `GetLMI` method.  To
iterate through all MAndRs in the incident, in LMI order, call the incident's
`GetIndex` method with index numbers ranging from zero to `Count`.

Other `Incident` methods include:

* AddDraft to add a new outgoing message in draft state
* Receive to add a newly received message (and generate receipt)
* Remove to remove an outgoing message in draft or queued state
* HasBulletin to detect whether the incident already has a particular bulletin
* Refresh to reload all messages from disk

The three messages in an MandR (`M`, `DR`, and `RR`) are of type `Message`,
which extends `xscmsg.Message` and provides all the same fields and methods.  In
addition, it provides metadata fields (`BBS`, `Bulletin`, `Received`, `From`,
and `To` for received messages; `Sent`, `From`, and `To` for outgoing messages);
methods to detect the message state (`IsReceived`, `IsBulletin`, `IsSent`,
`IsReady`); and methods to change the state of an outgoing message (`MarkDraft`,
`MarkReady`, `MarkSent`).  Finally, it provides `SetOperatorFields` to update
the operator fields of the form if the message has them.  After making any
changes to messages, call the MAndR's `Save` method to save them to disk.

## Incident Storage

Each incident is stored in a separate directory.  `pktmgr` expects the current
working directory to be set to the storage directory for the incident prior to
calling any `pktmgr` function.

The human-generated messages are stored in files named `«LMI».txt`, where
`«LMI»` is the local message ID assigned to them.  If the message is a
recognized form type, the rendered form is stored in `«LMI».pdf` as well.  The
message's corresponding delivery and read receipts are stored in `«LMI».DR.txt`
and `«LMI».RR.txt`, respectively.

When a remote message ID (`«RMI»`)for a human-generated message is known — i.e.,
when it can be parsed from the subject line of a received message, or when a
delivery receipt is received for a sent message — a symbolic link is created
from `«RMI».txt` to `«LMI».txt`, and when appropriate from `«RMI».pdf` to
`«LMI».pdf`.  This allows people looking in the incident directory to find a
message using either ID.

Each message is stored in RFC 5322 email format.  Outgoing messages have `To`
and `Subject` headers, and may also have `From` and `Date` headers.  The
presence of both `From` and `Date` indicates that the message has been sent.
The presence of `From` without `Date` indicates that the message is queued for
sending, but has not yet been sent.  The absence of either `From` or `Date`
indicates that the message is a draft, not ready to be sent.

Received messages have the `From`, `To`, `Date`, and `Subject` headers they were
received with, plus a `Received` header added by `pktmgr`.  (Other headers in
the received message, if any, are discarded.)  The `Received` header marks the
message as incoming rather than outgoing, and provides information about what
BBS and bulletin area it was retrieved from and when it was received.

In addition to the message files, each incident directory also contains an
`ICS-309.csv` file and an `ICS-309.pdf` file.  `pktmgr` keeps these files up to
date with a current ICS-309 form listing all messages in the incident.
