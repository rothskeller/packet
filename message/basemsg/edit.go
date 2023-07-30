package basemsg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/common"
)

// TODO defining these functions on BaseMessage means that only messages for
// which they're appropriate can leverage BaseMessage.

// EditFields returns the set of editable fields of the message.  Callers may
// change the Value in each field, but must otherwise treat the set as
// read-only.  Changing the Value in a field does not affect the underlying
// message until ApplyEdits is called.
func (bm *BaseMessage) EditFields() []*message.EditField {
	if bm.editFields == nil {
		bm.editFieldMap = make(map[*Field]*message.EditField)
		for _, f := range bm.Fields {
			if f.EditWidth == 0 {
				continue // not editable
			}
			var ef = message.EditField{
				Label:          f.Label,
				Width:          f.EditWidth,
				Multiline:      f.Multiline,
				LocalMessageID: bm.FOriginMsgID == f.Value,
				Help:           f.EditHelp,
				Hint:           f.EditHint,
			}
			if f.Choices != nil {
				ef.Choices = f.Choices.ListHuman()
			}
			bm.editFields = append(bm.editFields, &ef)
			bm.editFieldMap[f] = &ef
		}
	}
	for _, f := range bm.Fields {
		if ef := bm.editFieldMap[f]; ef != nil {
			if f.EditValue != nil {
				ef.Value = f.EditValue(f)
			} else if f.Choices != nil {
				ef.Value = f.Choices.ToHuman(*f.Value)
			} else {
				ef.Value = *f.Value
			}
			if p := validatePresence(f); p != "" {
				ef.Problem = p
			} else if f.EditValid != nil {
				ef.Problem = f.EditValid(f)
			} else if f.PIFOValid != nil {
				ef.Problem = f.PIFOValid(f)
			} else {
				ef.Problem = ""
			}
		}
	}
	return bm.editFields
}

// ApplyEdits applies the revised Values in the EditFields to the message.
func (bm *BaseMessage) ApplyEdits() {
	for _, f := range bm.Fields {
		if ef := bm.editFieldMap[f]; ef != nil {
			value := strings.TrimSpace(ef.Value)
			if f.EditApply != nil {
				f.EditApply(f, value)
			} else if f.Choices != nil {
				*f.Value = f.Choices.ToPIFO(value)
			} else {
				*f.Value = value
			}
			if f.EditValue != nil {
				ef.Value = f.EditValue(f)
			} else if f.Choices != nil {
				ef.Value = f.Choices.ToHuman(*f.Value)
			} else {
				ef.Value = *f.Value
			}
		}
	}
	for _, f := range bm.Fields {
		if ef := bm.editFieldMap[f]; ef != nil {
			if p := validatePresence(f); p != "" {
				ef.Problem = p
			} else if f.EditValid != nil {
				ef.Problem = f.EditValid(f)
			} else if f.PIFOValid != nil {
				ef.Problem = f.PIFOValid(f)
			} else {
				ef.Problem = ""
			}
		}
	}
}

// ApplyCardinal applies an edited value to a cardinal number field.
func ApplyCardinal(f *Field, v string) {
	if n, err := strconv.Atoi(v); err == nil {
		v = strconv.Itoa(n)
	}
	*f.Value = v
}

var dateLooseRE = regexp.MustCompile(`^(0?[1-9]|1[0-2])[-./](0?[1-9]|[12][0-9]|3[01])[-./](?:20)?([0-9][0-9])$`)

// ApplyDate applies an edited value to a date field.
func ApplyDate(f *Field, v string) {
	if match := dateLooseRE.FindStringSubmatch(v); match != nil {
		// Add leading zeroes and set delimiter to slash.
		v = fmt.Sprintf("%02s/%02s/20%s", match[1], match[2], match[3])
		// Correct values that are out of range, e.g. 06/31 => 07/01.
		if t, err := time.ParseInLocation("01/02/2006", v, time.Local); err == nil {
			v = t.Format("01/02/2006")
		}
	}
	*f.Value = v
}

// ApplyDateTime applies an edited value to a date/time field.
func ApplyDateTime(date, time *string, v string) {
	words := strings.Fields(v)
	if len(words) > 0 {
		var f = Field{Value: date}
		ApplyDate(&f, words[0])
	}
	if len(words) > 1 {
		var f = Field{Value: time}
		ApplyTime(&f, strings.Join(words[1:], " "))
	}
}

// ValueDateTime returns the formatted value of a date/time field.
func ValueDateTime(date, time string) string {
	return common.SmartJoin(date, time, " ")
}

var messageNumberLooseRE = regexp.MustCompile(`^([A-Z0-9]{3})-(\d+)([A-Z]?)$`)

// ApplyMessageNumber applies an edited value to a message number field.
func ApplyMessageNumber(f *Field, v string) {
	v = strings.ToUpper(v)
	if match := messageNumberLooseRE.FindStringSubmatch(v); match != nil {
		num, _ := strconv.Atoi(match[2])
		v = fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	*f.Value = v
}

var timeLooseRE = regexp.MustCompile(`^([1-9]:|[01][0-9]:?|2[0-4]:?)([0-5][0-9])$`)

// ApplyTime applies an edited value to a time field.
func ApplyTime(f *Field, v string) {
	if match := timeLooseRE.FindStringSubmatch(v); match != nil {
		// Add colon if needed.
		if !strings.HasSuffix(match[1], ":") {
			match[1] += ":"
		}
		// Add leading zero to hour if needed.
		v = fmt.Sprintf("%03s%s", match[1], match[2])
	}
	*f.Value = v
}
