package basemsg

import (
	"strings"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
	"golang.org/x/exp/slices"
)

// EncodeSubject encodes the message subject line.
func (bm *BaseMessage) EncodeSubject() string {
	var msgid, handling, formtag, subject string

	if bm.FOriginMsgID != nil {
		msgid = *bm.FOriginMsgID
	}
	if bm.FHandling != nil {
		handling = *bm.FHandling
	}
	if bm.Form != nil {
		formtag = bm.Form.Tag
	}
	if bm.FSubject != nil {
		subject = *bm.FSubject
	}
	return common.EncodeSubject(msgid, handling, formtag, subject)
}

// EncodeBody encodes the message body, suitable for transmission or
// storage.
func (bm *BaseMessage) EncodeBody() string {
	var (
		sb     strings.Builder
		enc    *common.PIFOEncoder
		values []string
	)
	if bm.Form.HTML == "" {
		panic("BaseMessage.EncodeBody can only encode PackItForms; other message types must override")
	}
	enc = common.NewPIFOEncoder(&sb, bm.Form.HTML, bm.Form.Version)
	values = make([]string, len(bm.Form.FieldOrder))
	for _, f := range bm.Fields {
		if f.PIFOTag == "" || *f.Value == "" {
			continue
		}
		if idx := slices.Index(bm.Form.FieldOrder, f.PIFOTag); idx >= 0 {
			values[idx] = *f.Value
		} else {
			enc.Write(f.PIFOTag, *f.Value)
		}
	}
	for i, tag := range bm.Form.FieldOrder {
		if values[i] != "" {
			enc.Write(tag, values[i])
		}
	}
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}

// Decode decodes the message.  It returns whether the message was recognized
// and decoded.
func Decode(body string, versions []*FormVersion, create func(*FormVersion) message.Message) (msg message.Message) {
	var (
		form *common.PIFOForm
		bm   *BaseMessage
	)
	// Decode the form and check for an HTML/version combo we recognize.
	if form = common.DecodePIFO(body); form == nil {
		return nil // not a form or not encoded properly
	}
	if idx := slices.IndexFunc(versions, func(v *FormVersion) bool {
		return form.HTMLIdent == v.HTML && form.FormVersion == v.Version
	}); idx < 0 {
		return nil // not an HTML/version combo we recognize
	} else {
		// Record which form version we actually saw.
		msg = create(versions[idx])
	}
	bm = msg.(interface{ Base() *BaseMessage }).Base()
	bm.PIFOVersion = form.PIFOVersion
	for _, f := range bm.Fields {
		if f.PIFOTag == "" {
			continue // field is not part of PIFO encoding
		}
		*f.Value = form.TaggedValues[f.PIFOTag]
	}
	// TODO really should make sure there aren't any fields unaccounted for.
	return msg
}
