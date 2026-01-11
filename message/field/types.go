package field

import (
	"regexp"
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message/address"
	"github.com/rothskeller/packet/message/messageid"
	"github.com/rothskeller/packet/message/msgifc"
)

//------------------------------------------------------------------------------

// NewAddressList creates a factory for a new address list field with the
// specified label.
func NewAddressList(label string) (ff *FieldFactory) {
	ff = NewField("", label)
	ff.tf = addressList{&ff.f}
	ff.f.editHeight = 1
	return ff
}

type addressList struct{ *field }

func (f addressList) Validate(m msgifc.Message, fi msgifc.Field, flags msgifc.ValidateFlags) (err error) {
	if err = f.validatePresence(m, fi); err != nil {
		return err
	}
	if flags&msgifc.VPIFOOnly != 0 {
		if _, err := address.ParseList(f.Value(m)); err != nil {
			return errors.NewF("The %q field does not contain a valid address list.", f.label)
		}
	}
	return f.validateCustom(m, fi, flags)
}

func (f addressList) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	if addrs, err := address.ParseList(s); err == nil {
		strs := make([]string, len(addrs))
		for i, a := range addrs {
			strs[i] = a.String()
		}
		s = strings.Join(strs, ", ")
	}
	return s
}

var fccCallSignRE = regexp.MustCompile(`^(?:A[A-L][0-9][A-Z]{1,3}|[KNW][0-9][A-Z]{2,3}|[KNW][A-Z][0-9][A-Z]{1,3})$`)

// NewFCCCallSign creates a factory for a new FCC call sign field with the
// specified tag and label.
func NewFCCCallSign(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 6
	ff.f.editHeight = 1
	ff.f.validateFunc = ValidateFCCCallSign
	ff.tf = fccCallSign{&ff.f}
	return ff
}

type fccCallSign struct{ *field }

func ValidateFCCCallSign(m msgifc.Message, fi msgifc.Field, flags msgifc.ValidateFlags) (err error) {
	if value := fi.Value(m); value != "" && flags&msgifc.VPIFOOnly != 0 && !fccCallSignRE.MatchString(value) {
		return errors.NewF("The %q field does not contain a valid FCC call sign.", fi.Label())
	}
	return nil
}

func (f fccCallSign) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	return strings.ToUpper(strings.TrimSpace(s))
}

// NewMessageID creates a factory for a new message ID field with the specified
// tag and label.  If packet is true, the message ID is validated to have a
// letter suffix.
func NewMessageID(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 9
	ff.f.editHeight = 1
	ff.f.editHint = "XXX-###P"
	ff.f.validateFunc = ValidateMessageID
	ff.tf = messageID{&ff.f}
	return ff
}

type messageID struct {
	*field
}

func ValidateMessageID(m msgifc.Message, fi msgifc.Field, flags msgifc.ValidateFlags) (err error) {
	if value := fi.Value(m); value != "" && flags&msgifc.VPIFOOnly != 0 {
		if _, _, _, err := messageid.Decode(value, false, flags&msgifc.VPacket != 0); err != nil {
			if flags&msgifc.VPacket != 0 {
				return errors.NewF("The %q field does not contain a valid packet message ID (XXX-###X).", fi.Label())
			} else {
				return errors.NewF("The %q field does not contain a valid message ID (XXX-### or XXX-###X).", fi.Label())
			}
		}
	}
	return nil
}

func (f messageID) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	s = strings.TrimSpace(s)
	if prefix, seq, suffix, err := messageid.Decode(s, true, false); err == nil {
		s, _ = messageid.Encode(prefix, seq, suffix)
	}
	return s
}

var tacticalCallSignRE = regexp.MustCompile(`^[A-Z][A-Z0-9]{4,5}$`)

// NewTacticalCallSign creates a factory for a new tactical call sign field with
// the specified tag and label.
func NewTacticalCallSign(tag, label string) (ff *FieldFactory) {
	ff = NewField(tag, label)
	ff.f.editWidth = 6
	ff.f.editHeight = 1
	ff.f.validateFunc = ValidateTacticalCallSign
	ff.tf = tacticalCallSign{&ff.f}
	return ff
}

type tacticalCallSign struct{ *field }

func ValidateTacticalCallSign(m msgifc.Message, fi msgifc.Field, flags msgifc.ValidateFlags) (err error) {
	if value := fi.Value(m); value != "" && flags&msgifc.VPIFOOnly != 0 && !tacticalCallSignRE.MatchString(value) {
		return errors.NewF("The %q field does not contain a valid tactical call sign.", fi.Label())
	}
	return nil
}

func (f tacticalCallSign) FromHuman(s string) string {
	if f.field.fromHumanFunc != nil {
		return f.field.fromHumanFunc(s)
	}
	return strings.ToUpper(strings.TrimSpace(s))
}
