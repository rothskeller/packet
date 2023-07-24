package msgedit

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

const spaces = "                                                                                                                                                                                                                                                "

// Wrap word-wraps the input string to the specified width, and returns a list
// of wrapped lines.
func Wrap(s string, w int) (l []string) {
	if s == "" {
		l = append(l, "")
	}
	for s != "" {
		if idx := strings.IndexByte(s, '\n'); idx >= 0 && idx < w {
			l = append(l, s[:idx])
			s = s[idx+1:]
			continue
		}
		if len(s) <= w {
			l = append(l, s)
			break
		}
		if idx := strings.LastIndexByte(s[:w], ' '); idx >= 0 {
			l = append(l, s[:idx])
			s = s[idx+1:]
		} else {
			l = append(l, s[:w])
			s = s[w:]
		}
	}
	return l
}

// Draw draws the specified string at the specified location with the specified
// style.
func Draw(screen tcell.Screen, text string, x, y int, style tcell.Style) {
	for _, r := range text {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

// DrawWidth draws the specified string at the specified location with the
// specified style.  The string is truncated or right-padded as needed to meet
// the specified width.
func DrawWidth(screen tcell.Screen, text string, x, y, width int, style tcell.Style) {
	if len(text) >= width {
		Draw(screen, text[:width], x, y, style)
	} else {
		Draw(screen, text, x, y, style)
		Draw(screen, spaces[:width-len(text)], x+len(text), y, style)
	}
}
