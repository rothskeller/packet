package xscmsg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rothskeller/packet/pktmsg"
)

// PlainTextTag identifies a plain text message.
const PlainTextTag = "plain"

func makePlainTextMessage(msg *pktmsg.Message) *Message {
	subject := ParseSubject(msg.Header.Get("Subject"))
	if subject == nil {
		subject = &XSCSubject{Subject: msg.Header.Get("Subject")}
	}
	return &Message{
		Type:       &plainTextMessageType,
		RawMessage: msg,
		Fields: []*Field{
			{
				Def:   &plainTextMsgNoField,
				Value: subject.MessageNumber,
			},
			{
				Def:   &plainTextHandlingField,
				Value: subject.HandlingOrder.String(),
			},
			{
				Def:   &plainTextSubjectField,
				Value: subject.Subject,
			},
			{
				Def:   &plainTextBodyField,
				Value: msg.Body,
			},
		},
	}
}

// CreatePlainTextMessage creates a new, empty plain text message.
func CreatePlainTextMessage() *Message {
	return &Message{
		Type: &plainTextMessageType,
		Fields: []*Field{
			{Def: &plainTextMsgNoField},
			{Def: &plainTextHandlingField},
			{Def: &plainTextSubjectField},
			{Def: &plainTextBodyField},
		},
	}
}

var plainTextMessageType = MessageType{
	Tag:     PlainTextTag,
	Name:    "plain text message",
	Article: "a",
	SubjectFunc: func(msg *Message) string {
		if msg.KeyField(FOriginMsgNo).Value != "" || msg.KeyField(FHandling).Value != "" {
			h, _ := ParseHandlingOrder(msg.KeyField(FHandling).Value)
			return msg.KeyField(FOriginMsgNo).Value + "_" + h.Code() + "_" + msg.KeyField(FSubject).Value
		}
		return msg.KeyField(FSubject).Value
	},
}

var plainTextMsgNoField = FieldDef{
	Tag:        string(FOriginMsgNo),
	Key:        FOriginMsgNo,
	Label:      "Origin Message Number",
	Flags:      Required,
	Validators: []Validator{ValidateMessageNumber},
}
var plainTextHandlingField = FieldDef{
	Tag:        string(FHandling),
	Label:      "Handling",
	Key:        FHandling,
	Validators: []Validator{ValidateChoices},
	Choices:    []string{"IMMEDIATE", "PRIORITY", "ROUTINE"},
	Flags:      Required,
}
var plainTextSubjectField = FieldDef{
	Tag:   string(FSubject),
	Key:   FSubject,
	Label: "Subject",
	Flags: Required,
}
var plainTextBodyField = FieldDef{
	Tag:   string(FBody),
	Key:   FBody,
	Label: "Body",
	Flags: Required | Multiline,
}

var messageNumberRE = regexp.MustCompile(`(?i)((?:[A-Z]{3}|[0-9][A-Z]{2}|[A-Z][0-9][A-Z])-)(\d+)([A-Z]?)$`)

// ValidateMessageNumber doesn't actually validate the message number, because
// PackItForms doesn't validate it, and we don't want to be raising errors that
// PackItForms doesn't.  However, in non-strict mode, if the message number is
// well formed, we can at least canonicalize it.
func ValidateMessageNumber(f *Field, _ *Message, strict bool) string {
	if !strict {
		if match := messageNumberRE.FindStringSubmatch(f.Value); match != nil {
			num, _ := strconv.Atoi(match[2])
			f.Value = fmt.Sprintf("%s%03d%s", strings.ToUpper(match[1]), num, strings.ToUpper(match[3]))
		}
	}
	return ""
}

// ValidateChoices ensures the field has one of the allowed values.  In
// non-strict mode, the values are case-insensitive, and any unambiguous prefix
// of an allowed value is accepted.
func ValidateChoices(f *Field, _ *Message, strict bool) string {
	var prefixOf string

	if f.Value == "" {
		return ""
	}
	for _, allowed := range f.Def.Choices {
		if strict && f.Value == allowed {
			return ""
		}
		if !strict && strings.EqualFold(allowed, f.Value) {
			return ""
		}
		if !strict && len(f.Value) < len(allowed) && strings.EqualFold(allowed[:len(f.Value)], f.Value) {
			if prefixOf == "" {
				prefixOf = allowed
			} else {
				prefixOf = "∅"
			}
		}
	}
	if prefixOf != "" && prefixOf != "∅" {
		f.Value = prefixOf
		return ""
	}
	return fmt.Sprintf("%q is not a valid value for field %q.", f.Value, f.Def.Tag)
}
