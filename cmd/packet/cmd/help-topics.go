package cmd

const configSlug = `incident configuration settings`
const configHelp = `
Configuration settings for the incident can be viewed with the "show config" command and changed with the "edit config" or "set config" commands.  These commands deal with the following configuration settings:

Incident Name
Activation Number
Operation Start
Operation End
    These are text placed at the top of generated ICS-309 communication logs.
Operator Call Sign
Operator Name
    These are the FCC call sign and name of the operator running the station. They are used for station identification during the BBS connection, as well as being filled into various forms.
Tactical Call Sign
Tactical Station Name
    These are the assigned call sign and name of the tactical station being operated.  They are filled into various forms.
Tx Message ID
Rx Message ID
    These are the message numbers to be assigned to the next outgoing and incoming message, respectively.  They are often set the same for interleaved message numbers.
Connection Method
    This is the method to use for connection to the BBS:  Manual, Serial+TNC, or Telnet.
BBS Call Sign
    This is the call sign of the BBS in use.  It is applicable for Manual and Telnet connection types.
Automatic Receipts
    This indicates whether outgoing delivery receipt messages should be automatically created for manually received messages.  It is applicable for the Manual connection type.
TNC Type
    This indicates the type of TNC in use for the Serial+TNC connection type.
Serial Port
    This identifies the serial port used for the Serial+TNC connection type.
BBS Address
    This identifies the address of the BBS.  For Serial+TNC connection types, it is an AX.25 address (e.g., W1XSC-1).  For Telnet connection types, it is a hostname:port string (e.g., w1xsc-gw.scc-ares-races.org:8000).  For Manual connection types, it is not applicable.
Bulletin Check 1-5
    These lines identify up to five bulletin areas to check for new bulletins, and the frequency at which bulletins should be checked in each.
Default To Address
    This is the address list to be added to the "To" field of any new message.
Default To ICS Position
Default To Location
Default From ICS Position
Default From Location
    These are default values for the addressing fields of new messages.
Default Body Text
    This is text to be added to the body text field of any new message, e.g., "**** This is drill traffic ****".
`

const filesSlug = `directory layout and file formats`
const filesHelp = `
Each incident has its own directory.  It contains the following files:

  LOC-111P.txt    ⇥message with local ID "LOC-111P", in RFC-5322 format
  LOC-111P.pdf    ⇥message with local ID "LOC-111P", in PDF format
  REM-222P.txt    ⇥symbolic link: remote ID "REM-222P" to local ID "LOC-111P"
  REM-222P.pdf    ⇥symbolic link: remote ID "REM-222P" to local ID "LOC-111P"
  ics309.pdf      ⇥ICS-309 communications log, in PDF format
  packet.json     ⇥incident configuration settings and log, in JSON format
  «bbscall».log   ⇥text file with log of all BBS communications
  .receipt33.txt  ⇥receipt message, in RFC-5322 format
  .unsent44.txt   ⇥unsent outgoing message, in RFC-5322 format
  .unsent44.pdf   ⇥unsent outgoing message, in PDF format
  .packet.lock    ⇥lock file for synchronization of incident changes

For messages that we received, LOC-111P.txt and LOC-111P.pdf contain the received message, and REM-222P.txt and REM-222P.pdf are named with the Origin Message ID of the received message.  (On Windows, the REM-222P links will not exist unless the user has been granted symbolic link creation privileges.)

For messages that we sent, LOC-111P.txt and LOC-111P.pdf contain the sent message, and REM-222P.txt and REM-222P.pdf are named with the destination stations' message IDs for the message we sent (which we pull from their delivery receipts).

For outgoing messages that we haven't sent yet, .unsent44.txt and .unsent44.pdf contain the message.  The number in the filename is the log entry number for the message.  The message has an "X-Packet-Queued: true" header if it is queued to be sent.

Receipt messages (incoming and outgoing) are stored in .receipt33.txt, where 33 is the log entry number for the message.
`

const scriptSlug = `how to use "packet" from scripts`
const scriptHelp = `
The "packet" commands provide script-friendly behavior when standard input and output are not a terminal.  In particular:
  - ⇥When standard output is not a terminal:
    - ⇥Confirmation and transient status messages are suppressed.
    - ⇥All colorization of the output is suppressed.
    - ⇥Commands that normally produce tables will produce CSV output.
    - ⇥Error messages are written to standard error instead of standard output.
  - ⇥When either standard input or standard output is not a terminal:
    - ⇥The "new" command prints to standard output the local message ID of the new message or log entry number of the new log entry, so that it can be used in subsequent commands.
    - ⇥The "new" command does not start an editor, and the "edit" command is not available.
    - ⇥The "set" command reads the new value of a field from standard input without any prompting or editing.
    - ⇥The "set" and "show" commands require the full name (or PIFO tag) of the field being changed or displayed, not an abbreviation.  This prevents scripts from being broken by future additions of fields or by future changes to the abbreviation heuristics.

For a script to send a message, it will usually follow this sequence:
  1. ⇥Run "packet set config" commands to set necessary configuration settings.
  2. ⇥Run "packet new" command to start a new message.
  3. ⇥Run "packet set" commands to set fields of that message.
  4. ⇥Run "packet mark ready" command to mark the message ready for sending.
  5. ⇥Run "packet connect" command to connect to the BBS and send it.
`
