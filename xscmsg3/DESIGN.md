# XSC Message Design

Package xscmsg adds Santa Clara County protocols and tactics on top of the
generic facilities in package pktmsg.

## Message Interface

The xscmsg.Message interface extends the pktmsg.Message interface by adding the
concept of named message types, iterators for viewable and editable fields, and
accessors for all of the well-known fields in XSC messages:

    type Message interface {
        pktmsg.Message
        Type() MessageType

        ViewableFields(func(ViewableField))
        EditableFields(func(EditableField))

        Body() BodyField
        DestinationMsgID() DestinationMsgIDField
        FormTag() FormTagField
        Handling() HandlingField
        MessageDate() MessageDateField
        MessageTime() MessageTimeField
        OpCall() OpCallField
        OpDate() OpDateField
        OpName() OpNameField
        OpTime() OpTimeField        
        OriginMsgID() OriginMsgIDField
        RawSubject() pktmsg.SubjectField
        Reference() ReferenceField
        Retrieved() RetrievedField
        Severity() SeverityField
        Subject() SubjectField
        TacCall() TacCallField
        TacName() TacNameField
        ToICSPosition() ToICSPositionField
        ToLocation() ToLocationField
    }

## Message Fields

Fields are more complex in the xscmsg package than they are in the pktmsg
package.  All fields implement the basic pktmsg.Field interface, and most of
them implement the pktmsg.SettableField interface.  But there are a number of
additional interfaces that extend those and that are implemented by some xscmsg
Fields:

    INTERFACE         METHODS
    ViewableField     Label() string
    SizedField        Size() (width, height int)
    DefaultedField    Default() string
    EditableField     Help(pktmsg.Message) string
                      Validate(m pktmsg.Message, pifo bool) string
    HintedField       Hint() string
    ChoicesField      Choices(pktmsg.Message) []string
    CalculatedField   Calculate(pktmsg.Message)

ViewableFields are those that are shown to an end user in a message viewer.
CalculatedFields are those whose values are calculated rather than provided by
an end user.  EditableFields (except for those that also implement
CalculatedField) are those that are presented to an end user for editing in a
message editor.  SizedFields, DefaultedFields, HintedFields, and ChoicesFields
all add additional capabilities to ViewableFields or EditableFields.

## Implementation Types

The xscMessage structure extends pktmsg.Message to handle the SCCo packet
message subject line standard.  It adds the OriginMsgID and Handling fields, and
when appropriate for the message, also adds Severity and/or FormTag fields.  It
remaps the pktmsg Subject field to RawSubject, and changes the Subject field to
contains only the actual message subject.  It does not honor xscmsg.Message; it
still only honors pktmsg.Message.

The baseMessage structure extends xscMessage to implement the full
xscmsg.Message interface.  To do so, it provides a named message type, and it
extends several of the standard pktmsg fields with xscmsg equivalents that
implement some of the additional field interfaces described above.

The plainText structure extends baseMessage with a specific message type
appropriate for plain text messages.

The unknownForm structure extends baseMessage with a specific message type
appropriate for PackItForms messages with an unrecognized form type.

## Message Type Registration

Package xscmsg provides a means for other packages to register handlers for
specific named message types.  These can be used to recognize existing messages
of those types as well as create new outgoing messages of those types.
Subpackages of xscmsg register handlers for all standard XSC form types.

## Message Creation

The xscmsg package provides four functions for creating messages:

    NewMessage() creates a new, outgoing plain text message.
    NewForm() creates a new, outgoing PackItForms message.
    ParseMessage() parses a message read from local storage and returns the
        appropriate message type (or an error).
    ReceiveMessage() parses a message instantly received from JNOS and returns
        the appropriate message type (or an error).

Each of these extends the like-named function in the pktmsg package, instead
returning a message that satisfies xscmsg.Message as well as pktmsg.Message.
The ParseMessage and ReceiveMessage functions will return a named message type
registered by another package if that package recognizes the message.
