package receipt

import (
	"fmt"
	"iter"
	"regexp"
	"strings"
	"time"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/body"
	"github.com/rothskeller/packet/message/cachetrack"
	"github.com/rothskeller/packet/message/field"
	"github.com/rothskeller/packet/message/msgifc"
	"github.com/rothskeller/packet/message/payload"
	"github.com/rothskeller/packet/message/subject"
)

var (
	deliveryReceiptMatchRE = regexp.MustCompile(`^\n*!LMI![^!]+!DR!`)
	deliveryReceiptRE      = regexp.MustCompile(`^\n*!LMI!([^!]+)!DR!(.+)\nYour Message\nTo: (.+)\nSubject: (.*)\nwas delivered on.*\nRecipient's Local.*\n`)
	ErrInvalidLMI          = errors.New("The local message ID must not contain newlines or exclamation points.")
	ErrNoDeliveryTime      = errors.New("The delivery time must be set.")
	ErrNoLMI               = errors.New("The local message ID must not be empty.")
	ErrDRSyntax            = errors.New("The body cannot be decoded as a delivery receipt because it does not conform to delivery receipt syntax.")
)

// NewDeliveryReceipt creates a new delivery receipt message (in draft state)
// with the specified parameters.  It returns an error if the parameters are
// invalid.  rmFrom is the from address for the message being receipted.  rmTo
// is the destination address for the message being receipted.  rmSubject is the
// encoded subject line of the message being receipted.  lmi is the message ID
// that we have assigned to the message being receipted (the "local message
// ID").  deliveryTime is the time that the message being receipted was
// delivered to us.  extraText is free-form text to be added to the end of the
// receipt message body.  (This is rarely used since Outpost users will never
// see it.)
func NewDeliveryReceipt(rmFrom, rmTo, rmSubject, lmi string, deliveryTime time.Time, extraText string) (m *message.DraftMessage, err error) {
	var (
		subj *subject.PlainSubject
		body *DeliveryReceiptBody
	)
	subj, _ = subject.NewPlainSubject("", "", "DELIVERED: "+rmSubject)
	subj.SetNonStandard()
	if body, err = newDeliveryReceiptBody(rmTo, rmSubject, lmi, deliveryTime, extraText); err != nil {
		return nil, err
	}
	m = message.NewDraftMessage(DeliveryReceipt, subj, payload.NewOutpostPayload(body, false), false)
	m.SetTo(rmFrom)
	m.SetReadyToSend(true)
	return m, nil
}

//--- MTYPE -------------------------------------------------------------------

type deliveryReceipt struct{ *message.BaseMType }

var DeliveryReceipt deliveryReceipt

func init() {
	DeliveryReceipt.BaseMType = message.NewBaseMType("a delivery receipt")
	message.RegisterType(DeliveryReceipt)
}

func (mt deliveryReceipt) Recognize(m message.Message) {
	if _, ok := m.Body().(*DeliveryReceiptBody); ok {
		m.SetType(DeliveryReceipt)
		m.Subject().(*subject.PlainSubject).SetNonStandard()
	}
}

//--- BODY --------------------------------------------------------------------

// A DeliveryReceiptBody is a body of a delivery receipt message.
type DeliveryReceiptBody struct {
	cachetrack.NoTracker
	body           string
	messageTo      string
	messageSubject string
	localMessageID string
	deliveryTime   string
	extraText      string
}

var _ body.Body = (*DeliveryReceiptBody)(nil)

// newDeliveryReceiptBody creates a new body for a delivery receipt message
// with the supplied contents.
func newDeliveryReceiptBody(rmTo, rmSubject, lmi string, deliveryTime time.Time, extraText string) (b *DeliveryReceiptBody, err error) {
	if rmTo == "" {
		err = errors.Join(err, ErrNoRMTo)
	} else if strings.ContainsAny(rmTo, "\r\n") {
		err = errors.Join(err, ErrNewlineInRMTo)
	}
	if rmSubject == "" {
		err = errors.Join(err, ErrNoRMSubject)
	} else if strings.ContainsAny(rmSubject, "\r\n") {
		err = errors.Join(err, ErrNewlineInRMSubject)
	}
	if lmi == "" {
		err = errors.Join(err, ErrNoLMI)
	} else if strings.ContainsAny(rmSubject, "!\r\n") {
		err = errors.Join(err, ErrInvalidLMI)
	}
	if deliveryTime.IsZero() {
		err = errors.Join(err, ErrNoDeliveryTime)
	}
	if err != nil {
		return nil, err
	}
	ftime := deliveryTime.Format("01/02/2006 15:04")
	return &DeliveryReceiptBody{
		body: fmt.Sprintf(
			"!LMI!%s!DR!%s\nYour Message\nTo: %s\nSubject: %s\nwas delivered on %[2]s\nRecipient's Local Message ID: %[1]s\n%[5]s",
			lmi, ftime, rmTo, rmSubject, extraText),
		messageTo:      rmTo,
		messageSubject: rmSubject,
		localMessageID: lmi,
		deliveryTime:   ftime,
		extraText:      extraText,
	}, nil
}

func init() { body.RegisterDecoder(decodeDeliveryReceiptBody) }
func decodeDeliveryReceiptBody(body string) (_ body.Body, err error) {
	if !deliveryReceiptMatchRE.MatchString(body) {
		return nil, nil
	}
	if match := deliveryReceiptRE.FindStringSubmatch(body); match == nil {
		return nil, ErrDRSyntax
	} else {
		return &DeliveryReceiptBody{
			body:           body,
			messageTo:      match[3],
			messageSubject: match[4],
			localMessageID: match[1],
			deliveryTime:   match[2],
			extraText:      body[len(match[0]):],
		}, nil
	}
}

// EncodedBody returns the encoded delivery receipt message body.
func (b *DeliveryReceiptBody) EncodedBody() string { return b.body }

// MessageTo returns the address to which the message being receipted was sent.
func (b *DeliveryReceiptBody) MessageTo() string { return b.messageTo }

// MessageSubject returns the subject of the message being receipted.
func (b *DeliveryReceiptBody) MessageSubject() string { return b.messageSubject }

// LocalMessageID returns the message ID assigned to the message being receipted
// by the receiver (i.e., by the sender of the receipt).
func (b *DeliveryReceiptBody) LocalMessageID() string { return b.localMessageID }

// DeliveryTime returns the date and time when the message being receipted was
// delivered to its receiver.  The date and time format are indeterminate.
func (b *DeliveryReceiptBody) DeliveryTime() string { return b.deliveryTime }

// ExtraText returns any extra text in the receipt message following the
// receipt itself.
func (b *DeliveryReceiptBody) ExtraText() string { return b.extraText }

func (b *DeliveryReceiptBody) Clone() body.Body { panic("should not be called") }
func (b *DeliveryReceiptBody) IsForm() bool     { return false }

func (b *DeliveryReceiptBody) Fields() iter.Seq[field.Field] {
	return func(yield func(field.Field) bool) {
		yield(receiptBodyField)
	}
}

var receiptBodyField = field.NewField("", "Body").
	Common(field.CDefaultBody).
	ValueFunc(func(m msgifc.Message) string { return m.Body().EncodedBody() }).
	MakeField()
