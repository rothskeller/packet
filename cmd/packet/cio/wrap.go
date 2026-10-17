package cio

import (
	"strings"
	"unicode/utf8"
)

// IndentMarker marks the indentation level for wrapped lines in a block of
// running text.  If a line contains this character (before the first wrap
// point), all wrapped continuations of that line will be indented to the
// location of this character.  This character is not included in the output.
const IndentMarker = '⇥'

// WrapText wraps a string to fit within the screen width.  It gives special
// control over the indentation of wrapped lines.  By default, wrapped lines are
// indented the same amount as the initial unwrapped line.  However, if the
// pre-wrapped line contains a "⇥" character that appears before the first wrap
// point, all wrapped lines are indented to the location of that character
// (which is not included in the output).
func (cio *CIO) WrapText(s string) (wrapped string) {
	width := max(cio.Width, 1)
	for line := range strings.SplitSeq(strings.TrimRight(s, "\n"), "\n") {
		var indent int

		line = strings.TrimRight(line, " ")
		if indent = strings.IndexRune(line, IndentMarker); indent >= 0 {
			// Remove the indent marker from the line.
			line = line[:indent] + line[indent+3:] // 3 bytes in UTF-8
			if indent >= width {
				indent = -1
			}
		}
		if indent < 0 {
			indent = strings.IndexFunc(line, func(r rune) bool { return r != ' ' })
			// The only way that comes out -1 is for an empty line,
			// which will never wrap, so its indent doesn't matter.
		}
		// Continuation lines need room for some text after the indent, or
		// the wrapping would never end.
		if indent < 0 || indent >= width {
			indent = 0
		}
		pad := strings.Repeat(" ", indent)
		for len(line) > width {
			// Break at the last space that fits, unless the indented
			// remainder would be no shorter (as when the space is
			// within the indentation): that would never finish.
			idx := strings.LastIndexByte(line[:width], ' ')
			if rest := strings.TrimLeft(line[max(idx, 0)+1:], " "); idx > 0 && indent+len(rest) < len(line) {
				wrapped += line[:idx] + "\n"
				line = pad + rest
			} else {
				// No usable space: break mid-word at the screen width,
				// but not in the middle of a multibyte character.
				cut := width
				for cut > indent+1 && !utf8.RuneStart(line[cut]) {
					cut--
				}
				wrapped += line[:cut] + "\n"
				line = pad + line[cut:]
			}
		}
		wrapped += line + "\n"
	}
	return wrapped
}

// wrap wraps a string to fit within the specified width.  It wraps at
// whitespace where possible, and mid-word if necessary.  The returned offsets
// array gives the offset into the original string of the start of each line.
func wrap(s string, width int) (lines []string, offsets []int) {
	width = max(width, 1) // a width less than 1 would never finish
	offsets = append(offsets, 0)
	offset := 0

	for _, line := range strings.Split(s, "\n") {
		for len(line) > width {
			if idx := strings.LastIndexByte(line[:width], ' '); idx > 0 {
				lines = append(lines, line[:idx+1])
				line = line[idx+1:]
				offset += idx + 1
				offsets = append(offsets, offset)
			} else {
				lines = append(lines, line[:width])
				line = line[width:]
				offset += width
				offsets = append(offsets, offset)
			}
		}
		lines = append(lines, line)
		offset += len(line) + 1 // account for newline
		offsets = append(offsets, offset)
	}
	// Note that this leaves the final offset at len(s)+1, since we
	// accounted for a newline on the final line that isn't actually in the
	// string.  However, this is good, because it means cursorToXY(len(s))
	// will return the end of the last line rather than the start of the
	// nonexistent one.
	return lines, offsets
}
