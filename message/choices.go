package message

import (
	"strings"

	"slices"
)

// A ChoiceMapper provides a set of choices for a Field and maps between
// internal (PIFO) and human representations of them.
type ChoiceMapper interface {
	// IsHuman returns whether the supplied string is a human representation
	// of one of the choices in the set.
	IsHuman(string) bool
	// IsPIFO returns whether the supplied string is a PIFO representation
	// of one of the choices in the set.
	IsPIFO(string) bool
	// ToHuman translates the supplied choice from PIFO representation to
	// human representation.  If the supplied string is not a PIFO
	// representation of one of the choices in the set, ToHuman returns its
	// argument unchanged.
	ToHuman(string) string
	// ToPIFO translates the supplied choice from human representation to
	// PIFO representation.  If the supplied string is not a human
	// representation of one of the choices in the set, ToPIFO returns its
	// argument unchanged.
	ToPIFO(string) string
	// ListHuman returns a list of the human representations of the choices
	// in the set.
	ListHuman() []string
}

////////////////////////////////////////////////////////////////////////////////

// NoChoices is a zero implementation of ChoiceMapper for fields that don't have
// choices.
type NoChoices struct{}

func (NoChoices) IsHuman(s string) bool   { return false }
func (NoChoices) IsPIFO(s string) bool    { return false }
func (NoChoices) ToHuman(s string) string { return s }
func (NoChoices) ToPIFO(s string) string  { return s }
func (NoChoices) ListHuman() []string     { return nil }

////////////////////////////////////////////////////////////////////////////////

// Choices is a wrapper around []string to make it a valid ChoiceMapper.  Each
// string in the slice is a choice whose PIFO and human representations are the
// same.
type Choices []string

func (c Choices) IsHuman(s string) bool {
	return slices.ContainsFunc(c, func(cs string) bool { return strings.EqualFold(s, cs) })
}
func (c Choices) IsPIFO(s string) bool    { return slices.Contains(c, s) }
func (c Choices) ToHuman(s string) string { return s }
func (c Choices) ToPIFO(s string) string {
	var match string
	if s == "" {
		return s
	}
	for _, v := range c {
		if len(v) >= len(s) && strings.EqualFold(v[:len(s)], s) {
			if match != "" {
				return s
			}
			match = v
		}
	}
	if match != "" {
		return match
	}
	return s
}
func (c Choices) ListHuman() []string { return c }

////////////////////////////////////////////////////////////////////////////////

// ChoicePairs is a wrapper around []string to make it a valid ChoiceMapper.
// Each pair of strings in the slice is a PIFO and a human representation of a
// choice.
type ChoicePairs []string

func (c ChoicePairs) IsHuman(s string) bool {
	for i := 1; i < len(c); i += 2 {
		if strings.EqualFold(c[i], s) {
			return true
		}
	}
	return false
}
func (c ChoicePairs) IsPIFO(s string) bool {
	for i := 0; i < len(c); i += 2 {
		if c[i] == s {
			return true
		}
	}
	return false
}
func (c ChoicePairs) ToHuman(s string) string {
	for i := 0; i < len(c)-1; i += 2 {
		if c[i] == s {
			return c[i+1]
		}
	}
	return s
}
func (c ChoicePairs) ToPIFO(s string) string {
	var match string
	for i := 1; i < len(c); i += 2 {
		if len(c[i]) >= len(s) && strings.EqualFold(c[i][:len(s)], s) {
			if match != "" {
				return s
			}
			match = c[i-1]
		}
	}
	if match != "" {
		return match
	}
	return s
}
func (c ChoicePairs) ListHuman() (human []string) {
	human = make([]string, len(c)/2)
	for i := 1; i < len(c); i += 2 {
		human[i/2] = c[i]
	}
	return human
}
