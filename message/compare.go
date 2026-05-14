package message

import (
	"slices"

	"github.com/rothskeller/packet/message/field"
)

// Compare compares the actual message to the expected message.  It returns a
// numerator and denominator of a score, where 0% is no match and 100% is a
// perfect match.  It also returns a list of the detailed comparisons of each
// field.
func Compare(exp, act Message) (score, outOf int, fields []*field.ComparedField) {
	// First, we need to be sure that the two messages are of comparable
	// types.
	if !exp.Type().CanCompareAgainst(act.Type()) {
		return 0, 1, []*field.ComparedField{{
			Label: "Message Type", Score: 0, OutOf: 1,
			Expected: exp.Type().Name(), ExpectedMask: "*",
			Actual: act.Type().Name(), ActualMask: "*",
		}}
	}
	// Get the set of fields in the actual message.
	actfs := slices.Collect(act.Fields())
	// Walk through each field in the expected message.
	for expf := range exp.Fields() {
		var actf field.Field

		// If the expected field is a common one, look for the same
		// common field in the actual.
		if expf.Common() != "" {
			if idx := slices.IndexFunc(actfs, func(f field.Field) bool {
				return f != nil && f.Common() == expf.Common()
			}); idx >= 0 {
				actf = actfs[idx]
				actfs[idx] = nil
			}
		}
		// If we haven't found a match and the expected field has a tag,
		// look for the same tag in the actual.
		if actf == nil && expf.Tag() != "" {
			if idx := slices.IndexFunc(actfs, func(f field.Field) bool {
				return f != nil && f.Tag() == expf.Tag()
			}); idx >= 0 {
				actf = actfs[idx]
				actfs[idx] = nil
			}
		}
		// If we still haven't found an actual, there's nothing to
		// compare this expected against, so skip it.
		if actf == nil {
			continue
		}
		// Compare the fields.
		eval := expf.Value(exp)
		aval := actf.Value(act)
		if eval == "" && aval == "" {
			continue
		}
		cmp := expf.Compare(expf.Label(), eval, aval)
		if cmp == nil {
			continue
		}
		fields = append(fields, cmp)
		score += cmp.Score
		outOf += cmp.OutOf
	}
	return score, outOf, fields
}
