package unkform

import (
	"strings"

	"github.com/rothskeller/packet/message/common"
)

// EncodeSubject encodes the message subject.
func (f *UnknownForm) EncodeSubject() string {
	return common.EncodeSubject(f.OriginMsgID, f.Handling, f.FormTag, f.Subject)
}

// EncodeBody encodes the message body.
func (f *UnknownForm) EncodeBody() string {
	var (
		sb  strings.Builder
		enc *common.PIFOEncoder
	)
	enc = common.NewPIFOEncoder(&sb, f.FormHTML, f.FormVersion)
	for tag, value := range f.TaggedValues {
		enc.Write(tag, value)
	}
	if err := enc.Close(); err != nil {
		panic(err)
	}
	return sb.String()
}
