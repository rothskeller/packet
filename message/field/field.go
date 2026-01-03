package field

import (
	"github.com/rothskeller/packet/message/msgifc"
)

// Field is the interface satisfied by all message fields.  Note that this is
// more than just PackItForms fields; there are fields relating to metadata as
// well, and non-forms messages have fields also.
type Field = msgifc.Field

type ChoicePair = msgifc.ChoicePair

// A ComparedField structure represents the comparison of a single field.
type ComparedField struct {
	// Label is the field label.
	Label string
	// Score is the comparison score for this field.  0 <= Score <= OutOf.
	Score int
	// OutOf is the maximum possible score for this field, i.e., the score
	// for this field if its contents match exactly.
	OutOf int
	// Expected is the value of this field in the expected message (i.e.,
	// the receiver of the Compare method), formatted for human viewing.
	Expected string
	// ExpectedMask is a string describing which characters of Expected are
	// different from those in Actual.  Space characters in the mask
	// correspond to characters in Expected that are properly matched by
	// Actual.  "~" characters in the mask correspond to characters in
	// Expected that have minor differences in Actual.  All other characters
	// in the mask correspond to significant differences.  If ExpectedMask
	// is shorter than Expected, the last character of ExpectedMask is
	// implicitly repeated.
	ExpectedMask string
	// Actual is the value of this field in the actual message (i.e., the
	// argument of the Compare method), formatted for human viewing.
	Actual string
	// ActualMask is a string describing which characters of Actual are
	// different from those in Expected.  Space characters in the mask
	// correspond to characters in Actual that properly match Expected.  "~"
	// characters in the mask correspond to characters in Actual that have
	// minor differences with Expected.  All other characters in the mask
	// correspond to significant differences.  If ActualMask is shorter than
	// Actual, the last character of ActualMask is implicitly repeated.
	ActualMask string
}
