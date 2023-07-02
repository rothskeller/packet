# XSC Message Design

The sub-packages of package xscmsg define message types that are registered with
the typedmsg package, and can therefore be used to create new outgoing messages
and/or recognize existing messages of those types.  Package xscmsg itself
provides a base implementation and common definitions used by many of those
message types.

All of the message types defined in this package tree honor the typedmsg.Message
interface.  Many of them honor additional interfaces as well:

## xscmsg.IMessage

    type IMessage interface {
        typedmsg.Message
  	    OpCall() string
  	    SetDestinationMsgID(string)
  	    SetOperator(received bool, call, name string)
  	    SetOriginMsgID(string)
  	    SetTactical(call, name string)
  	    ToICSPosition() string
  	    ToLocation() string
  	    Validate() []string
        View() []LabelValue
    }

The xscmsg.IMessage interface is honored by message types representing messages
written by humans.  Non-human messages, like delivery receipts, do not honor it.
It provides:

* Accessors and mutators for common fields.  These are no-ops if the specific
  message type doesn't have that field.
* A Validate method, which verifies whether the message meets the requirements
  imposed by XSC-standard packet software (Outpost and PackItForms).
* A View method, which renders the message as an ordered list of
  (label, value) pairs suitable for viewing in a message viewer.

## xscmsg.Editable

    type Editable interface {
        Message
        Edit() []Field
    }

The xscmsg.Editable interface is honored for human message types of which new
messages can be composed and sent.  In practice, almost all xscmsg.Message types
honor xscmsg.Editable as well, but the UnknownForm type does not, and some types
representing obsolete versions of forms do not.  Editable provides an ordered
list of fields that can be edited.  Each field honors the Field interface:

    type Field interface {
        Label() string
        Value() string
        SetValue(string)
        Size() (width, height int)
        Problem() string
        Choices() []string
        Help() string
        Hint() string
    }

## Type Hierarchy

Message types are built through an embedding hierarchy, as follows.  Types shown
in parentheses are not registered message types.

    (pktmsg.Message)
        delivrcpt.DeliveryReceipt
        readrcpt.ReadReceipt
        (xscmsg.Message)
            plaintext.PlainText
            checkin.CheckInMessage
            checkout.CheckOutMessage
            (xscmsg.BaseForm)
                unkform.UnknownForm
                ics213.ICS213Form
                (xscmsg.StdForm)
                    eoc213rr.EOC213RRForm
                    havbed.HAvBedForm             (in xscmsgpvt)
                    medresreq.MedResReqForm       (in xscmsgpvt)
                    racesmar.RacesMARForm
                    (xscmsg.UpdateCompleteForm)
                        ahfacstat.AHFacStatForm
                        jurisstat.JurisStatForm
                        medfacstat.MedFacStatForm (in xscmsgpvt)
                        sheltstat.SheltStatForm

Each of the concrete message types is in its own package which, as a matter of
convention, exports the following symbols (using EOC-213RR as an example):

    EOC213RRForm - the message type
    EOC213RRType - the typedmsg.MessageType structure describing the type
    NewEOC213RR  - a method to create a new message of the type

Fields are built through a wrapping chain, as follows.  All of them have
xscmsg.BaseField as the innermost type:

    xscmsg.BaseField

They may then be wrapped in one of the types that handle standard types:

    xscmsg.BooleanField
    xscmsg.CardinalNumberField
    xscmsg.ChoicesField
    xscmsg.DateField
    xscmsg.FCCCallField
    xscmsg.MessageNumberField
    xscmsg.PhoneNumberField
    xscmsg.RealNumberField
    xscmsg.TacticalCallField
    xscmsg.TimeField
    xscmsg.UnknownField

Then they may be wrapped in one of the types that handles required presence:

    xscmsg.RequiredField
    xscmsg.RequiredIfCompleteField
