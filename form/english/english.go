package english

import "strings"

// Article returns "a" or "an", whichever is the proper article for the supplied
// string, which is assumed to be a noun phrase.  It uses heuristics and may not
// always be correct.
func Article(s string) string {
	if s == "" {
		return ""
	}
	switch strings.ToLower(s[:1]) {
	case "a", "e", "i", "o", "u":
		return "an"
	default:
		return "a"
	}
	// This algorithm is ridiculously inaccurate.  There are many words that
	// start with a vowel that sounds like a consonant ("one", "uniform")
	// and should take "a".  And there are many words that start with a
	// consonant that sounds like a vowel ("hour") and should take "an".
	// And acronyms may go by how the name of the first letter is
	// pronounced, in which case ones starting with "F", "H", "M", "N", "R",
	// "S", and "X" should take "an".  However, none of these matter for the
	// names of any currently known forms, so we'll ignore them until they
	// cause a problem.
}

// ToLower converts an English phrase to lower case, but it leaves acronyms
// (defined as words that start with an uppercase letter followed by anything
// other than a lowercase letter) in their original case.  It does not preserve
// exact spacing in the phrase.
func ToLower(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) < 2 || word[0] < 'A' || word[0] > 'Z' || (word[1] >= 'a' && word[1] <= 'z') {
			words[i] = strings.ToLower(word)
		}
	}
	return strings.Join(words, " ")
}
