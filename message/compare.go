package message

import (
	"regexp"
	"strconv"
	"strings"
)

// A CompareField structure represents a single field in the comparison of two
// messages as performed by the Compare method.
type CompareField struct {
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

// Compare compares two messages.  It returns a score indicating how closely
// they match, and the detailed comparisons of each field in the message.  The
// comparison is not symmetric:  the receiver of the call is the "expected"
// message and the argument is the "actual" message.
func (bm *BaseMessage) Compare(actual Message) (score int, outOf int, cfields []*CompareField) {
	var act = actual.Base()

	if act.Type != bm.Type {
		return 0, 1, []*CompareField{{
			Label: "Message Type",
			Score: 0, OutOf: 1,
			Expected:     bm.Type.Name,
			ExpectedMask: "*",
			Actual:       act.Type.Name,
			ActualMask:   "*",
		}}
	}
	for i, expf := range bm.Fields {
		actf := act.Fields[i]
		if expf.Value != nil && actf.Value != nil && (*expf.Value != "" || *actf.Value != "") {
			comp := expf.Compare(expf.Label, *expf.Value, *actf.Value)
			if comp == nil {
				continue // omit from comparison
			}
			if comp.OutOf == 0 {
				comp.OutOf = 1
			}
			score, outOf = score+comp.Score, outOf+comp.OutOf
			cfields = append(cfields, comp)
		}
	}
	return score, outOf, cfields
}

// CompareNone is a "comparison" function that causes the field to be omitted
// from the comparison.
func CompareNone(string, string, string) *CompareField {
	return nil
}

// CompareCardinal compares two values for a field that is supposed to contain
// a cardinal number.
func CompareCardinal(label, exp, act string) (c *CompareField) {
	if eval, err := strconv.Atoi(exp); err == nil {
		if aval, err := strconv.Atoi(act); err == nil {
			if eval == aval {
				return &CompareField{
					Label: label, Expected: exp, Actual: act, Score: 2, OutOf: 2,
					ExpectedMask: " ", ActualMask: " ",
				}
			}
		}
	}
	return CompareExact(label, exp, act)
}

// CompareCheckbox compares two values for a checkbox field.
func CompareCheckbox(label, exp, act string) (c *CompareField) {
	c = &CompareField{
		Label: label, Expected: exp, Actual: act, OutOf: 2,
	}
	if exp == act {
		c.ExpectedMask, c.ActualMask, c.Score = " ", " ", 2
		if exp == "false" {
			c.Expected, c.Actual = "", ""
		}
	} else {
		c.ExpectedMask, c.ActualMask, c.Score = "*", "*", 0
		if exp == "" || exp == "false" {
			c.Expected = "not checked"
		}
		if act == "" || act == "false" {
			c.Actual = "not checked"
		}
	}
	return c
}

// CompareExact compares two values for a field, which are expected to be
// exactly identical.
func CompareExact(label, exp, act string) (c *CompareField) {
	c = &CompareField{
		Label: label, Expected: exp, Actual: act, OutOf: 2,
	}
	if exp == act {
		c.ExpectedMask, c.ActualMask, c.Score = " ", " ", 2
	} else {
		c.ExpectedMask, c.ActualMask, c.Score = "*", "*", 0
		if c.Expected == "" {
			c.Expected = "(not set)"
		}
		if c.Actual == "" {
			c.Actual = "(not set)"
		}
	}
	return c
}

// CompareExactMap compares two values for a field, which are expected to be
// exactly identical.  It maps the values for human display.
func CompareExactMap(label, exp, act string, mapping map[string]string) (c *CompareField) {
	var ok bool

	c = &CompareField{
		Label: label, OutOf: 2,
	}
	if c.Expected, ok = mapping[exp]; !ok {
		c.Expected = exp
	}
	if c.Actual, ok = mapping[act]; !ok {
		c.Actual = act
	}
	if exp == act {
		c.ExpectedMask, c.ActualMask, c.Score = " ", " ", 2
	} else {
		c.ExpectedMask, c.ActualMask, c.Score = "*", "*", 0
		if c.Expected == "" {
			c.Expected = "(not set)"
		}
		if c.Actual == "" {
			c.Actual = "(not set)"
		}
	}
	return c
}

var dateRE = regexp.MustCompile(`^(\d?\d)([-/.])(\d?\d)([-/.])(20)?(\d\d)$`)

// CompareDate compares two values for a date field.
func CompareDate(label, exp, act string) (c *CompareField) {
	c = &CompareField{
		Label: label, Expected: exp, Actual: act, Score: 2, OutOf: 2,
	}
	if exp == act {
		c.ExpectedMask, c.ActualMask, c.Score = " ", " ", 2
		return c
	}
	expM := dateRE.FindStringSubmatch(exp)
	actM := dateRE.FindStringSubmatch(act)
	if expM == nil || actM == nil {
		c.ExpectedMask, c.ActualMask, c.Score = "*", "*", 0
		if c.Expected == "" {
			c.Expected = "(not set)"
		}
		if c.Actual == "" {
			c.Actual = "(not set)"
		}
		return c
	}
	switch {
	case expM[1] == actM[1]:
		c.ExpectedMask += strings.Repeat(" ", len(expM[1]))
		c.ActualMask += strings.Repeat(" ", len(actM[1]))
	case expM[1] == "0"+actM[1]:
		c.ExpectedMask += "~ "
		c.ActualMask += " "
		c.Score = min(c.Score, 1)
	case "0"+expM[1] == actM[1]:
		c.ExpectedMask += " "
		c.ActualMask += "~ "
		c.Score = min(c.Score, 1)
	default:
		c.ExpectedMask += strings.Repeat("*", len(expM[1]))
		c.ActualMask += strings.Repeat("*", len(actM[1]))
		c.Score = 0
	}
	if expM[2] == actM[2] {
		c.ExpectedMask += " "
		c.ActualMask += " "
	} else {
		c.ExpectedMask += "~"
		c.ActualMask += "~"
		c.Score = min(c.Score, 1)
	}
	switch {
	case expM[3] == actM[3]:
		c.ExpectedMask += strings.Repeat(" ", len(expM[3]))
		c.ActualMask += strings.Repeat(" ", len(actM[3]))
	case expM[3] == "0"+actM[3]:
		c.ExpectedMask += "~ "
		c.ActualMask += " "
		c.Score = min(c.Score, 1)
	case "0"+expM[3] == actM[3]:
		c.ExpectedMask += " "
		c.ActualMask += "~ "
		c.Score = min(c.Score, 1)
	default:
		c.ExpectedMask += strings.Repeat("*", len(expM[3]))
		c.ActualMask += strings.Repeat("*", len(actM[3]))
		c.Score = 0
	}
	if expM[4] == actM[4] {
		c.ExpectedMask += " "
		c.ActualMask += " "
	} else {
		c.ExpectedMask += "~"
		c.ActualMask += "~"
		c.Score = min(c.Score, 1)
	}
	switch {
	case expM[5] == "" && actM[5] == "":
		break
	case expM[5] == actM[5]:
		c.ExpectedMask += "  "
		c.ActualMask += "  "
	case expM[5] == "":
		c.ActualMask += "~~"
		c.Score = min(c.Score, 1)
	default:
		c.ExpectedMask += "~~"
		c.Score = min(c.Score, 1)
	}
	if expM[6] == actM[6] {
		c.ExpectedMask += "  "
		c.ActualMask += "  "
	} else {
		c.ExpectedMask += "**"
		c.ActualMask += "**"
		c.Score = 0
	}
	return c
}

// CompareReal compares two values for a field that is supposed to contain a
// real number.
func CompareReal(label, exp, act string) (c *CompareField) {
	if eval, err := strconv.ParseFloat(exp, 64); err == nil {
		if aval, err := strconv.ParseFloat(act, 64); err == nil {
			if eval == aval {
				return &CompareField{
					Label: label, Expected: exp, Actual: act, Score: 2, OutOf: 2,
					ExpectedMask: " ", ActualMask: " ",
				}
			}
		}
	}
	return CompareExact(label, exp, act)
}

var timeRE = regexp.MustCompile(`^(\d?\d)(:?)(\d\d)$`)

// CompareTime compares two values for a time field.
func CompareTime(label, exp, act string) (c *CompareField) {
	c = &CompareField{
		Label: label, Expected: exp, Actual: act, Score: 2, OutOf: 2,
	}
	if exp == act {
		c.ExpectedMask, c.ActualMask, c.Score = " ", " ", 2
		return c
	}
	expM := timeRE.FindStringSubmatch(exp)
	actM := timeRE.FindStringSubmatch(act)
	if expM == nil || actM == nil {
		c.ExpectedMask, c.ActualMask, c.Score = "*", "*", 0
		if c.Expected == "" {
			c.Expected = "(not set)"
		}
		if c.Actual == "" {
			c.Actual = "(not set)"
		}
		return c
	}
	switch {
	case expM[1] == actM[1]:
		c.ExpectedMask += strings.Repeat(" ", len(expM[1]))
		c.ActualMask += strings.Repeat(" ", len(actM[1]))
	case expM[1] == "0"+actM[1]:
		c.ExpectedMask += "~ "
		c.ActualMask += " "
		c.Score = min(c.Score, 1)
	case "0"+expM[1] == actM[1]:
		c.ExpectedMask += " "
		c.ActualMask += "~ "
		c.Score = min(c.Score, 1)
	default:
		c.ExpectedMask += strings.Repeat("*", len(expM[1]))
		c.ActualMask += strings.Repeat("*", len(actM[1]))
		c.Score = 0
	}
	if expM[2] == actM[2] {
		c.ExpectedMask += strings.Repeat(" ", len(expM[2]))
		c.ActualMask += strings.Repeat(" ", len(actM[2]))
	} else {
		c.ExpectedMask += strings.Repeat("~", len(expM[2]))
		c.ActualMask += strings.Repeat("~", len(actM[2]))
		c.Score = min(c.Score, 1)
	}
	if expM[3] == actM[3] {
		c.ExpectedMask += "  "
		c.ActualMask += "  "
	} else {
		c.ExpectedMask += "**"
		c.ActualMask += "**"
		c.Score = 0
	}
	return c
}

// ComparePhoneNumber compares two values for a phone number field.
func ComparePhoneNumber(label, exp, act string) (c *CompareField) {
	c = &CompareField{
		Label: label, Expected: exp, Actual: act, Score: 2, OutOf: 2,
	}
	enums := strings.Map(digitsOnly, exp)
	anums := strings.Map(digitsOnly, act)
	if enums != anums {
		c.ExpectedMask, c.ActualMask, c.Score = "*", "*", 0
		if c.Expected == "" {
			c.Expected = "(not set)"
		}
		if c.Actual == "" {
			c.Actual = "(not set)"
		}
	} else if exp != act {
		// This is a bit simplistic — it would be nice to mark
		// specifically the punctuation that's different, rather than
		// the whole string — but it's probably not worth the effort.
		c.ExpectedMask, c.ActualMask, c.Score = "~", "~", 1
	}
	c.ExpectedMask, c.ActualMask, c.Score = " ", " ", 2
	return c
}
func digitsOnly(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}

// CompareText compares two values for a textual field.  Textual fields use a
// somewhat loose comparison, as follows:
//   - If a group starts with "¡", a case-sensitive comparison is used for that
//     group (without the "¡").  Otherwise, a group matches if the actual has
//     the same case as the expected, is in all caps, or is in all lowercase.
//   - Runs of spaces are treated as a single whitespace for comparison.
//   - Newlines in expected must be matched by newlines in actual; the presence
//     or absence of spaces on either side is ignored.
//   - Runs of spaces without a newline in expected can be matched by a single
//     (but not multiple) newline in actual, optionally surrounded by spaces.
func CompareText(label, exp, act string) (c *CompareField) {
	var et, at []token

	c = &CompareField{
		Label: label, Expected: exp, Actual: act,
	}
	// First, we split both "exp" and "act" into tokens.
	et = textCompareSplit(exp, true)
	at = textCompareSplit(act, false)
	// Streamline common case of comparing empties.
	if len(et) == 0 && len(at) == 0 {
		c.Score, c.OutOf = 1, 1
		return c
	}
	if len(et) == 0 && len(at) != 0 {
		et = append(et, token{tok: "(not set)"})
	}
	if len(et) != 0 && len(at) == 0 {
		at = append(at, token{tok: "(not set)"})
	}
	// Max score is 2 points for each group, or 1 point for empty expected.
	if c.OutOf = len(et) * 2; c.OutOf != 0 {
		c.Score = c.OutOf
	} else {
		c.OutOf = 1
	}
	// Calculate the LCS matrix of the two strings.
	matrix := make([][]int, len(et)+1)
	for i := 0; i <= len(et); i++ {
		matrix[i] = make([]int, len(at)+1)
		for j := 0; j <= len(at); j++ {
			matrix[i][j] = 0
		}
	}
	for i := 1; i <= len(et); i++ {
		for j := 1; j <= len(at); j++ {
			if strings.EqualFold(et[i-1].tok, at[j-1].tok) {
				matrix[i][j] = matrix[i-1][j-1] + 1
			} else {
				matrix[i][j] = max(matrix[i][j-1], matrix[i-1][j])
			}
		}
	}
	// Run the diff between the strings.
	et, at = diffTokens(et, at, matrix, len(et), len(at))
	// Permit soft newlines.
	et, at = permitSoftNewlines(et, at)
	// Calculate the score.
	if c.Score = c.Score - overallPenalty(et, at); c.Score < 0 {
		c.Score = 0
	}
	// Calculate the human-visible diff strings.
	c.ExpectedMask, c.ActualMask = makeMaskStrings(et, at)
	return c
}

type token struct {
	exactCase bool
	tok       string
	sep       string // "", " ", "\n", or "\n\n".
}

// textCompareSplit breaks the supplied string into tokens.  A token is a run of
// non-space characters.  As an exception, if a run of non-whitespace characters
// ends in a punctuation character, the punctuation character is considered a
// separate token.
//
// Each token has an associated separator character string, describing the
// whitespace that follows the token.  (The whitespace preceding the first token
// is ignored.)  It can act the following values:
//   - "" indicates no whitespace after the token.  This happens for the final
//     token in the string: trailing whitespace is ignored.  It also happens for
//     tokens that ended with a punctuation character as noted above.
//   - " " indicates that the token was followed by one or more spaces with no
//     newlines.
//   - "\n" indicates that the token was followed by a single newline, possibly
//     with spaces on one or both sides.
//   - "\n\n" indicates that the token was followed by a string of whitespace
//     containing multiple newlines (it doesn't matter how many).
//
// Each token has an exactCase flag, set if the token string begins with the "¡"
// exact case marker.
func textCompareSplit(s string, exactCaseAllowed bool) (tokens []token) {
	s = strings.TrimSpace(s) // ignore leading and trailing whitespace
	for s != "" {
		var (
			tok       string
			exactCase bool
		)
		if idx := strings.IndexAny(s, " \n"); idx >= 0 {
			tok, s = s[:idx], s[idx:]
		} else {
			tok, s = s, ""
		}
		if exactCase = exactCaseAllowed && strings.HasPrefix(tok, "¡"); exactCase {
			tok = tok[2:]
		}
		if len(tok) > 1 && strings.IndexByte(",:;?!", tok[len(tok)-1]) >= 0 {
			tokens = append(tokens, token{exactCase, tok[:len(tok)-1], ""})
			tok, exactCase = tok[len(tok)-1:], false
		}
		if idx := strings.IndexFunc(s, func(r rune) bool {
			return r != ' ' && r != '\n'
		}); idx >= 0 {
			if nl1 := strings.IndexByte(s[:idx], '\n'); nl1 >= 0 {
				if strings.IndexByte(s[nl1+1:idx], '\n') >= 0 {
					tokens = append(tokens, token{exactCase, tok, "\n\n"})
				} else {
					tokens = append(tokens, token{exactCase, tok, "\n"})
				}
			} else {
				tokens = append(tokens, token{exactCase, tok, " "})
			}
			s = s[idx:]
		} else {
			tokens = append(tokens, token{exactCase, tok, ""})
		}
	}
	return tokens
}

// diffTokens is a recursive function that performs a diff between two sets of
// tokens using the supplied LCS matrix.  It returns two new parallel lists of
// tokens: one with the tokens from the first set, in order, with "ø" tokens
// interspersed; the other with the tokens from the second set, in order, with
// "ø" tokens interspersed.  The two are aligned such that the tokens they
// act in common are at the same indexes.
func diffTokens(et, at []token, matrix [][]int, wl, hl int) (oet, oat []token) {
	if wl > 0 && hl > 0 && strings.EqualFold(et[wl-1].tok, at[hl-1].tok) {
		oet, oat = diffTokens(et, at, matrix, wl-1, hl-1)
		oet = append(oet, et[wl-1])
		oat = append(oat, at[hl-1])
	} else if hl > 0 && (wl == 0 || matrix[wl][hl-1] > matrix[wl-1][hl]) {
		oet, oat = diffTokens(et, at, matrix, wl, hl-1)
		oet = append(oet, token{false, "", ""})
		oat = append(oat, at[hl-1])
	} else if wl > 0 && (hl == 0 || matrix[wl][hl-1] <= matrix[wl-1][hl]) {
		oet, oat = diffTokens(et, at, matrix, wl-1, hl)
		oet = append(oet, et[wl-1])
		oat = append(oat, token{false, "", ""})
	} else {
		oet = make([]token, wl)
		copy(oet, et)
		oat = make([]token, hl)
		copy(oat, at)
	}
	return oet, oat
}

// permitSoftNewlines revises the result of diffTokens to allow for cases where
// et has a space and at has a newline.  That's acceptable.
func permitSoftNewlines(et, at []token) (_, _ []token) {
	for i := range et {
		if et[i].sep != " " || at[i].sep != "\n" {
			continue
		}
		if !strings.EqualFold(et[i].tok, at[i].tok) {
			continue
		}
		at[i].sep = " "
	}
	return et, at
}

// overallPenalty calculates the penalty for all of the diffTokens.  It is
// changePenalties + max(addPenalties, removePenalties).
func overallPenalty(et, at []token) (overall int) {
	var adds, removes int

	for i := range et {
		var pen = penalty(et[i], at[i])
		if et[i].tok == "" {
			adds += pen
		} else if at[i].tok == "" {
			removes += pen
		} else {
			overall += pen
		}
	}
	return overall + max(adds, removes)
}

// penalty returns the amount subtracted from the possible score for the string,
// based on the difference between the supplied "act" and "exp" tokens.  It
// returns zero if the tokens match, 1 for spacing or case sensitivity problems,
// and 2 for a more significant mismatch.
func penalty(exp, act token) int {
	switch {
	case !strings.EqualFold(exp.tok, act.tok):
		return 2
	case exp.sep != act.sep:
		return 1
	case exp.exactCase && exp.tok != act.tok:
		return 1
	case exp.tok != act.tok && strings.ToLower(exp.tok) != act.tok && strings.ToUpper(exp.tok) != act.tok:
		return 1
	default:
		return 0
	}
}

// makeMaskStrings creates the mask strings that indicate the locations of
// differences and their severities.
func makeMaskStrings(et, at []token) (em, am string) {
	for i := range et {
		if et[i].tok != "" && at[i].tok != "" { // same token
			if penalty(token{et[i].exactCase, et[i].tok, ""}, token{false, at[i].tok, ""}) != 0 {
				em += strings.Repeat("~", len(et[i].tok))
				am += strings.Repeat("~", len(at[i].tok))
			} else {
				em += strings.Repeat(" ", len(et[i].tok))
				am += strings.Repeat(" ", len(at[i].tok))
			}
			if et[i].sep != at[i].sep {
				em += strings.Repeat("~", len(et[i].sep))
				am += strings.Repeat("~", len(at[i].sep))
			} else {
				em += strings.Repeat(" ", len(et[i].sep))
				am += strings.Repeat(" ", len(at[i].sep))
			}
		} else {
			em += strings.Repeat("*", len(et[i].tok)+len(et[i].sep))
			am += strings.Repeat("*", len(at[i].tok)+len(at[i].sep))
		}
	}
	return em, am
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
