package cmdutil

import (
	"strings"

	"github.com/rothskeller/packet/errors"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/field"
)

// MatchFieldName finds the message field that (best) matches the supplied field
// name.  If loose is true, it can be a partial, heuristic match.
func MatchFieldName(msg message.Message, in string, loose bool) (field.Field, error) {
	// First priority is a match on PIFO tag.
	for f := range msg.Fields() {
		if f.Tag() == in {
			return f, nil
		}
	}
	// Remaining comparisons are case-sensitive if the input contains any
	// uppercase letters.
	var caseSensitive bool
	var compare func(string, string) bool
	if strings.IndexFunc(in, func(r rune) bool { return r >= 'A' && r <= 'Z' }) >= 0 {
		compare, caseSensitive = func(a, b string) bool { return a == b }, true
	} else {
		compare, caseSensitive = strings.EqualFold, false
	}
	// Second priority is a case-insensitive match on full field name.
	for f := range msg.Fields() {
		if compare(f.Label(), in) {
			return f, nil
		}
	}
	// To avoid the need for quoting, we will also accept a case-insensitive
	// match on the full field name with spaces removed.
	if !strings.Contains(in, " ") {
		for f := range msg.Fields() {
			if compare(strings.ReplaceAll(f.Label(), " ", ""), in) {
				return f, nil
			}
		}
	}
	// Unless the loose flag is set, those are the only options.
	if !loose {
		return nil, errors.NewF("There is no field %q.", in)
	}
	// Now we look for fields whose name contains the same characters as the
	// input, in the same order, but also contains additional characters.
	// We return the field that has the smallest number of unmatched
	// capital letters in its name, and among those, the one with the
	// smallest number of unmatched characters.
	var match field.Field
	var missedUC, missedCH int
	for f := range msg.Fields() {
		if mUC, mCH, ok := matchFieldName(f.Label(), in, caseSensitive); ok {
			if match == nil || mUC < missedUC || (mUC == missedUC && mCH < missedCH) {
				match, missedUC, missedCH = f, mUC, mCH
			}
		}
	}
	if match == nil {
		return nil, errors.NewF("There is no field %q.", in)
	}
	return match, nil
}

// matchFieldName is a recursive function that determines whether the input is a
// valid shortening of the field name, and returns the heuristic scoring if so.
// It's a pretty expensive algorithm, but at human time scales it's negligible.
func matchFieldName(fname, in string, caseSensitive bool) (mUC, mCH int, ok bool) {
	if in == "" && fname == "" {
		// Nothing left of either string.  Perfect match.
		return 0, 0, true
	}
	if in == "" {
		// Nothing left of in, but we still have some fname.  Compute
		// the score.  Use a recursive call to get the score after
		// removing the first character of fname, then add the score
		// for that character.
		mUC, mCH, _ = matchFieldName(fname[1:], in, caseSensitive)
		if fname[0] >= 'A' && fname[0] <= 'Z' {
			mUC++
		}
		mCH++
		return mUC, mCH, true
	}
	if fname == "" {
		// Nothing left of fname, but we still have some in.  Not a
		// match at all.
		return 0, 0, false
	}
	var mUC1, mCH1, mUC2, mCH2 int
	var ok1, ok2 bool
	// If the lead characters of fname and in match, calculate the score
	// based on matching those two.
	if fname[0] == in[0] || (!caseSensitive && downcase(fname[0]) == downcase(in[0])) {
		mUC1, mCH1, ok1 = matchFieldName(fname[1:], in[1:], caseSensitive)
	}
	// Whether the lead characters match or not, also calculate the score
	// assuming they don't.
	mUC2, mCH2, ok2 = matchFieldName(fname[1:], in, caseSensitive)
	if fname[0] >= 'A' && fname[0] <= 'Z' {
		mUC2++
	}
	mCH2++
	// Return the better of the two scores.
	if !ok1 && !ok2 {
		return 0, 0, false
	}
	if !ok1 {
		return mUC2, mCH2, ok2
	}
	if !ok2 || mUC1 < mUC2 || (mUC1 == mUC2 && mCH1 < mCH2) {
		return mUC1, mCH1, ok1
	}
	return mUC2, mCH2, ok2
}

func downcase(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b + 'A' - 'a'
	}
	return b
}
