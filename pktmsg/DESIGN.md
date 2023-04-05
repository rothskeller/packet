# Packet Message Design

The main export of the pktmsg package is the Message type, which is an interface
satisfied by all packet messages:

    type Message interface {
        BBSRxDate() BBSRxDateField
        Body() BodyField
        FormHTML() FormHTMLField
        FormVersion() FormVersionField
        FromAddr() FromAddrField
        NotPlainText() NotPlainTextField
        OutpostFlags() OutpostFlagsField
        PIFOVersion() PIFOVersionField
        ReturnAddr() ReturnAddrField
        RxArea() RxAreaField
        RxBBS() RxBBSField
        RxDate() RxDateField
        SentDate() SentDateField
        Subject() SubjectField
        ToAddrs() ToAddrsField

        TaggedField(string) Field
        TaggedFields(func(string, Field))

        Save() string
        Transmit() (to []string, subject, body string)
    }

This interface contains a set of accessor methods for well-known fields, a pair
of methods for fetching and iterating PackItForms tagged fields, and a pair of
methods for rendering the message into encoded forms.

## Message Fields

The 15 fields listed in the Message interface are common fields, shared by many
messages and often accessed programmatically.  They are made easy to access by
providing each one with its own exported interface and its own accessor method
in the Message interface.  Not every message has every one of these fields; the
accessor methods will return nil when a message does not have the requested
field.

All of the well-known field interfaces extend the Field interface, and most of
them extend the SettableField interface:

    type Field interface {
        Value() string
    }
    type SettableField interface {
        Field
        SetValue(string)
    }

Some well-known fields have additional methods for value conversion.  For
example, date fields have methods for converting to and from time.Time values.

In addition to the well-known fields, PackItForms messages also have tagged
fields, identified by their PackItForms field tags.  The TaggedField method can
be used to access a tagged field given its tag, and the TaggedFields method
iterates through all of the tagged fields of the message.  The implementations
of all tagged fields extend the SettableField interface.

## Implementation Types

The baseTx structure provides the base implementation for all Messages.  It
includes the Body, FromAddr, SentDate, Subject, and ToAddrs fields, and base
implementations of all Message methods.

The outpostMessage structure extends baseTx, adding an understanding of Outpost
flags in the body of a message.  These are !B64!, !URG!, !RDR!, and !RRR!.  The
!B64! flag is interpreted and used to modify the FBody field.  The others are
stored in the OutpostFlags field.

The baseRx structure extends outpostMessage, adding the fields that are
meaningful only for received messages:  RxArea, RxBBS, and RxDate.  RxArea is
present only for received bulletins.

The baseRetrieved structure extends baseRx, adding the fields that are
meaningful only for immediately received messages — that is, messages that we
just now retrieved from a BBS, as opposed to previously retrieved messages that
we just read from local storage.  These fields are not preserved in local
storage.  These include ReturnAddr and (when appropriate) BBSRxDate and/or
NotPlainText.

Finally, the pifoMessage structure is a layer on top of outpostMessage, baseRx,
or baseRetrieved that understands PackItForms form encodings (although it knows
nothing about specific forms).  pifoMessage does not expose the underlying Body
field.  Instead, it exposes the form metadata fields, FormHTML, FormVersion, and
PIFOVersion.  It also exposes each of the fields of the form itself, as tagged
fields identified by their PackItForms field tag.

## Message Creation

The pktmsg package provides four functions for creating messages:

    NewMessage() creates a new outpostMessage, with the fields appropriate for
        an outgoing message.
    NewForm() creates a new pifoMessage, with the fields appropriate for an
        outgoing message.  The FPIFOVersion field is given a default value.
    ParseMessage() parses a message read from local storage and returns the
        appropriate message type (or an error).
    ReceiveMessage() parses a message instantly received from JNOS and returns
        the appropriate message type (or an error).
