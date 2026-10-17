package cio

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

// wrapWithin calls WrapText with the specified width, failing the test if it
// doesn't finish promptly.
func wrapWithin(t *testing.T, width int, s string) string {
	t.Helper()
	done := make(chan string, 1)
	go func() { done <- (&CIO{Width: width}).WrapText(s) }()
	select {
	case wrapped := <-done:
		return wrapped
	case <-time.After(2 * time.Second):
		t.Fatalf("WrapText(width %d, %q) did not finish", width, s)
		return ""
	}
}

func TestWrapText(t *testing.T) {
	got := wrapWithin(t, 20, "ERROR: ⇥the quick brown fox jumps over the lazy dog")
	want := "ERROR: the quick\n       brown fox\n       jumps over\n       the lazy dog\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
	got = wrapWithin(t, 20, "  indented text that wraps around")
	want = "  indented text\n  that wraps around\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestWrapTextLongWord(t *testing.T) {
	// A warning containing a word too long for the line after the
	// indentation, such as a long path: this used to loop forever, since
	// the only space in reach was in the indentation.
	path := "/private/tmp/" + strings.Repeat("long-directory-name/", 6) + "forms"
	got := wrapWithin(t, 80, "WARNING: ⇥The program could not use "+path+".  It will use embedded forms.")
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	for i, line := range lines {
		if len(line) > 80 {
			t.Errorf("line %d is %d characters long: %q", i, len(line), line)
		}
		if i != 0 && !strings.HasPrefix(line, strings.Repeat(" ", 9)) {
			t.Errorf("line %d is not indented to the marker: %q", i, line)
		}
	}
	var joined strings.Builder
	for _, line := range lines {
		joined.WriteString(strings.TrimLeft(line, " "))
	}
	if !strings.Contains(joined.String(), path) {
		t.Errorf("the path was not kept intact:\n%s", got)
	}
}

func TestWrapTextEdgeCases(t *testing.T) {
	for _, tt := range []struct {
		width int
		text  string
	}{
		{10, strings.Repeat(" ", 12) + "abc def ghi"},     // indentation wider than the screen
		{10, "ERROR: ⇥" + strings.Repeat("x", 30)},        // no spaces at all after the marker
		{10, "ERROR:   ⇥word " + strings.Repeat("y", 30)}, // spaces only within the indentation
		{20, "abcdefghijkl⇥" + strings.Repeat("z", 40)},   // marker past most of the width
		{1, "hello world"},
		{0, "hello world"},
		{-5, "hello world"},
		{10, strings.Repeat("é", 30)},                      // multibyte characters
		{400, "WARNING: ⇥" + strings.Repeat("word ", 150)}, // wider than the old 300-space pad
		{8, ""},
		{8, "\n\n"},
		{2, "    ⇥"},      // only spaces before a marker past the width: used to panic
		{3, "      ⇥   "}, // likewise
		{5, "        ⇥x"}, // used to hang
	} {
		got := wrapWithin(t, tt.width, tt.text)
		if !utf8.ValidString(got) {
			t.Errorf("width %d, %q: wrapping split a character: %q", tt.width, tt.text, got)
		}
		if squash(got) != squash(strings.ReplaceAll(tt.text, "⇥", "")) {
			t.Errorf("width %d, %q: text was lost: %q", tt.width, tt.text, got)
		}
	}
}

// squash removes all spaces and newlines, leaving the text that wrapping must
// preserve.
func squash(s string) string {
	return strings.NewReplacer(" ", "", "\n", "").Replace(s)
}
