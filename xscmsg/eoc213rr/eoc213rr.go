package eoc213rr

import (
	"github.com/rothskeller/packet/pktmsg"
	"github.com/rothskeller/packet/typedmsg"
	"github.com/rothskeller/packet/xscmsg"
)

// Type is the type definition for a resource request form.
var Type = typedmsg.MessageType{
	Tag:       "EOC213RR",
	Name:      "EOC-213RR resource request form",
	Article:   "an",
	Create:    create,
	Recognize: recognize,
}

// EOR213RR is a resource request form.
type EOR213RR struct {
	*xscmsg.StdForm
	IncidentName        string
	DateInitiated       string
	TimeInitiated       string
	TrackingNumber      string
	RequestedBy         string
	PreparedBy          string
	ApprovedBy          string
	QtyUnit             string
	ResourceDescription string
	ResourceArrival     string
	Priority            string
	EstdCost            string
	DeliverTo           string
	DeliverToLocation   string
	Substitutes         string
	EquipmentOperator   bool
	Lodging             bool
	Fuel                bool
	FuelType            string
	Power               bool
	Meals               bool
	Maintenance         bool
	Water               bool
	Other               bool
	Instructions        string
}

// NewEOR213RR creates a new resource request form.
func NewEOR213RR() *EOR213RR {
	return &EOR213RR{Message: &xscmsg.Message{Message: new(pktmsg.Message)}}
}

func create() typedmsg.Message { return NewEOR213RR() }

func recognize(base *pktmsg.Message) typedmsg.Message {
	return &EOR213RR{Message: &xscmsg.Message{Message: base}}
}

// Type returns the type of the message.
func (m *EOR213RR) Type() *typedmsg.MessageType { return &Type }

// View returns the set of viewable fields of the message.
func (m *EOR213RR) View() []xscmsg.LabelValue {
	var lvs = m.ViewHeaders()
	lvs = append(lvs, xscmsg.LV("Subject", m.SubjectHeader))
	lvs = append(lvs, xscmsg.LV("Body", m.Body))
	return lvs
}

// Edit returns the set of editable fields of the message.
func (m *EOR213RR) Edit() []xscmsg.Field {
	if m.fields == nil {
		m.fields = []xscmsg.Field{
			xscmsg.NewToField(&m.To),
			xscmsg.NewOriginMsgIDField(&m.OriginMsgID),
			xscmsg.NewHandlingField(&m.Handling),
			&subjectField{xscmsg.NewBaseField(&m.Subject)},
			&bodyField{xscmsg.NewBaseField(&m.Body)},
		}
	}
	return m.fields
}

type subjectField struct{ *xscmsg.BaseField }

func (f subjectField) Label() string    { return "Subject" }
func (f subjectField) Size() (w, h int) { return 80, 1 }
func (f subjectField) Problem() string {
	if f.Value() == "" {
		return "The message subject is required."
	}
	return ""
}
func (f subjectField) Help() string {
	return "This is the message subject.  It is required."
}

type bodyField struct{ *xscmsg.BaseField }

func (f bodyField) Label() string    { return "Body" }
func (f bodyField) Size() (w, h int) { return 80, 10 }
func (f bodyField) Problem() string {
	if f.Value() == "" {
		return "The message body is required."
	}
	return ""
}
func (f bodyField) Help() string {
	return "This is the message body.  It is required."
}
