package basemsg

import (
	"fmt"
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
			if f.EditHelp == "" {
				continue // not editable
			}
			var choices = f.Choices.ListHuman()
			var editWidth = f.EditWidth
			for _, c := range choices {
				if len(c) > editWidth {
					editWidth = len(c)
				}
			}
			var ef = message.EditField{
				Label:          f.Label,
				Width:          editWidth,
				Multiline:      f.Multiline,
				LocalMessageID: bm.FOriginMsgID == f.Value,
				Help:           f.EditHelp,
				Hint:           f.EditHint,
				Choices:        choices,
			}
			bm.editFields = append(bm.editFields, &ef)
			bm.editFieldMap[f] = &ef
		}
	}
	for _, f := range bm.Fields {
		if ef := bm.editFieldMap[f]; ef != nil {
			ef.Value = f.EditValue(f)
			if ef.Problem = validatePresence(f); ef.Problem == "" {
				ef.Problem = f.EditValid(f)
			}
		}
	}
	return bm.editFields
}

// ApplyEdits applies the revised Values in the EditFields to the message.
func (bm *BaseMessage) ApplyEdits() {
	for _, f := range bm.Fields {
		if ef := bm.editFieldMap[f]; ef != nil {
			f.EditApply(f, ef.Value)
			ef.Value = f.EditValue(f)
		}
	}
	for _, f := range bm.Fields {
		if ef := bm.editFieldMap[f]; ef != nil {
			if ef.Problem = validatePresence(f); ef.Problem == "" {
				ef.Problem = f.EditValid(f)
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

// ApplyMessageNumber applies an edited value to a message number field.
func ApplyMessageNumber(f *Field, v string) {
	v = strings.ToUpper(v)
	if match := messageNumberLooseRE.FindStringSubmatch(v); match != nil {
		num, _ := strconv.Atoi(match[2])
		v = fmt.Sprintf("%s-%03d%s", match[1], num, match[3])
	}
	*f.Value = v
}

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
