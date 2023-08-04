package message

import (
	"testing"
)

var compareTextTests = []struct {
	name    string
	exp     string
	act     string
	score   int
	outOf   int
	expMask string
	actMask string
}{
	{
		"empty match",
		"",
		"",
		1, 1,
		"",
		"",
	},
	{
		"nonempty match",
		"Hello World",
		"Hello World",
		4, 4,
		"           ",
		"           ",
	},
	{
		"simple mismatch",
		"Hello",
		"World",
		0, 2,
		"*****",
		"*****",
	},
	{
		"missing word",
		"Hello World",
		"World",
		2, 4,
		"******     ",
		"     ",
	},
	{
		"added word",
		"World",
		"Hello World",
		0, 2,
		"     ",
		"******     ",
	},
	{
		"allow all caps",
		"Hello World",
		"HELLO WORLD",
		4, 4,
		"           ",
		"           ",
	},
	{
		"allow all lowercase",
		"Hello World",
		"hello world",
		4, 4,
		"           ",
		"           ",
	},
	{
		"disallow other case",
		"Hello World",
		"hELLo World",
		3, 4,
		"~~~~~      ",
		"~~~~~      ",
	},
	{
		"require exact case",
		"¡Hello World",
		"hello World",
		3, 4,
		"~~~~~      ",
		"~~~~~      ",
	},
	{
		"allow newlines and multiple spaces",
		"Hello World",
		"Hello  \n  World",
		4, 4,
		"           ",
		"           ",
	},
	{
		"match with newline",
		"Hello\nWorld",
		"Hello  \n  World",
		4, 4,
		"           ",
		"           ",
	},
	{
		"require newline",
		"Hello\nWorld",
		"Hello World",
		3, 4,
		"     ~     ",
		"     ~     ",
	},
	{
		"require empty line",
		"Hello\n\nWorld",
		"Hello World",
		3, 4,
		"     ~~     ",
		"     ~     ",
	},
	{
		"unexpected empty line",
		"Hello World",
		"Hello\n\nWorld",
		3, 4,
		"     ~     ",
		"     ~~     ",
	},
	{
		"match with punctuation",
		"Hello, World",
		"Hello, World",
		6, 6,
		"            ",
		"            ",
	},
	{
		"space before punctuation",
		"Hello, World",
		"Hello , World",
		5, 6,
		"            ",
		"     ~       ",
	},
	{
		"no space before punctuation",
		"Hello , World",
		"Hello, World",
		5, 6,
		"     ~       ",
		"            ",
	},
}

func TestCompareText(t *testing.T) {
	for _, tt := range compareTextTests {
		t.Run(tt.name, func(t *testing.T) {
			cf := CompareText(tt.name, tt.exp, tt.act)
			if cf.Score != tt.score || cf.OutOf != tt.outOf {
				t.Errorf("score was %d/%d, should be %d/%d", cf.Score, cf.OutOf, tt.score, tt.outOf)
			}
			if cf.ExpectedMask != tt.expMask || cf.ActualMask != tt.actMask {
				t.Errorf("Expected: %q\nWantMask: %q\nHaveMask: %q\nActual:   %q\nWantMask: %q\nHaveMask: %q\n",
					tt.exp, tt.expMask, cf.ExpectedMask, tt.act, tt.actMask, cf.ActualMask)
			}
		})
	}
}
