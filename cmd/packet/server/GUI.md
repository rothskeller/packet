# Packet GUI Design

The centerpiece of the packet GUI is a chronological table of messages. There
are two possible formats: the ICS-309 format

    TIME  FROMCALL  FROMID--  TOCALL--  TOID----  SUBJECT

or the packet shell 1.x format

    TIME FROMID-- -> LOCALID-             SUBJECT
    TIME             LOCALID- -> TOID---- SUBJECT

The packet shell format is more compact and makes the distinction between sent
and received messages more obvious. It also is consistent with the TUI, which
will remain using that format because of its 80-column width limit. The
ICS-309 format is more familiar for packet users and consistent with what will
be printed.

Both formats should be evaluated on a number of special cases:

- Messages without a local ID.
- Messages that we hear but aren't the sender or receiver.
- Manual ICS-309 entries not connected to a packet message.
- Voice messages.
- Use of the software to create an ICS-309 without doing packet work.
- Display of bulletins.
- Edit in place vs. edit in dialog.

Setting that question aside for the moment, we then consider what surrounds
the table, the operations that can be performed, and how they're invoked. The
operations sort themselves into two categories, incident-global:

- Select incident
- View/change incident configuration
- Create incident archives
- Create incident ICS-309
- Connect to the BBS for send/receive

and per-message operations:

- Import a new message from a file
- Import a new message from a manual receive
- Create a new draft message (specify type)
- Open a finalized message as PDF
- Open a non-finalized message in an editor
- Delete a message
- Mark a non-finalized message ready or not ready to send
- Print a message
- Show the command needed to manually send a message
- Edit the ICS-309 entry for a message
- Add an ICS-309 entry not associated with a packet message
- Mark as "delivered"
- Mark as "followup needed"
- Generate delivery receipt

There are also UI options:

- Help
- Change font size
- Change table layout (if we support more than one)
- Change light/dark mode (if we support both)

We'll plan to use the standard interaction model where clicking on a message
selects it without taking any action on it, and the up/down/home/end/pgup/pgdn
keys also change the selection. Double-clicking on a message or hitting Enter
invokes the default action on it.

Two main interaction models come to mind from here. The first is that we have
a top menu bar with actions available on the selected message. In that case,
the default action is probably "Open" (as PDF for a finalized message, in an
editor otherwise). The table remains static.

The other option is that the default option inserts some vertical space between
the selected message and the one after it, and displays a panel of information
and action buttions in that space. The shortcuts of the action buttons would
work on the selected message even if the panel weren't open. Esc would close
the panel while leaving the message selected. There would also be >v icons at
the left to open and close the panel by mouse. I like that because it avoids
needing to open a dialog box for basic stuff, but I'm not sure how well it
would be received. Also, there are a lot of per-message operations to try to
squeeze into such a box.

The box might look like:

    Date         Time    FromCall From Msg ID  To Call  To Msg ID
    [--/--/----] [--:--] [------] [----------] [------] [----------]
    [-(Subject)------------------------------------------------------------]
    [X] Ready to Send [X] Delivered [X] Followup Needed [X] Voice

    [Open v] [Delete] [Print] [Manual Send] [Make Receipt]   [Cancel] [Save]

with lots of variations on the buttons and checkboxes available based on the
message state.

If we do that, then the top bar only needs a few buttons:

    [New Incident]  [Setup]  [Export]  [Connect]  [Help]

Now, let's consider the incident settings form.

    Incident Directory:    non-editable string
    Incident Name: [---] Activation Number [---]
    Op Start: [--/--/----] [--:--]  Op End: [--/--/----] [--:--]
    Operator:     Call Sign: [------] Name: [---]
    Tac Station:  Call Sign: [------] Name: [---]
    Starting Message ID:  Receive [--------] Send [--------]

    Defaults for New Messages:
    [The usual header To/From layout, position and location only]
    Body:  [---]

    BBS to Use: [------] via (x) Serial+TNC (x) Telnet (x) Manual
    Serial+TNC:  TNC Type: [KPC-3+ v]
                 Serial Port: [--- v]
                 BBS Address: [--------]
    Telnet:      BBS Address: [------]
                 BBS Password: [------] for (user: non-editable string)
    Manual:      [X] Generate delivery receipts automatically

    Bulletin Retrieval:
        Area     Check Every         Check Now
        [------] [---] hr [---] min  [X]
