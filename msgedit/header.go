package msgedit

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// HeaderLine is a string of line-drawing characters, as long as we expect to
// need.
const HeaderLine = "════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════"

// A HeaderBar is the bar appearing at the top of a panel, consisting of a
// title on the left and a (possibly empty) set of buttons on the right, all
// overlaid on a lined background.
type HeaderBar struct {
	*tview.Box
	title    string
	keys     []string
	evts     []*tcell.EventKey
	labels   []string
	offsets  []int
	onButton func(string, func(tview.Primitive))
}

// NewHeaderBar creates a new header bar with the specified title.
func NewHeaderBar(title string) *HeaderBar {
	return &HeaderBar{Box: tview.NewBox(), title: title}
}

// SetOnButton sets the function to be called when a button is pressed.  The
// function is passed the button name (as passed to SetButtons) and the setFocus
// function.
func (b *HeaderBar) SetOnButton(fn func(string, func(tview.Primitive))) *HeaderBar {
	b.onButton = fn
	return b
}

// SetButtons sets the set of buttons for the header bar.  The arguments come in
// pairs, where the first is a key and the second is a button label. Keys are
// identified by a single character (for printable keys) or by name (for non-
// printing keys like ESC and F10).
func (b *HeaderBar) SetButtons(pairs ...string) *HeaderBar {
	if len(pairs)%2 != 0 {
		panic("SetButtons must have an even number of arguments")
	}
	b.keys, b.evts, b.labels = b.keys[:0], b.evts[:0], b.labels[:0]
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		label := pairs[i+1]
		b.keys = append(b.keys, key)
		switch {
		case len(key) == 1:
			b.evts = append(b.evts, tcell.NewEventKey(tcell.KeyRune, rune(key[0]), 0))
			if strings.HasPrefix(label, key) {
				b.labels = append(b.labels, "("+key+")"+label[1:])
			} else {
				b.labels = append(b.labels, key+"="+label)
			}
		case key == "ESC":
			b.evts = append(b.evts, tcell.NewEventKey(tcell.KeyESC, 0, 0))
			b.labels = append(b.labels, key+"="+label)
		case key == "F10":
			b.evts = append(b.evts, tcell.NewEventKey(tcell.KeyF10, 0, 0))
			b.labels = append(b.labels, key+"="+label)
		default:
			panic("unsupported key in arguments to SetButtons")
		}
	}
	b.offsets = make([]int, len(b.keys))
	return b
}

// Draw draws the header bar.
func (b *HeaderBar) Draw(screen tcell.Screen) {
	b.DrawForSubclass(screen, b)
	x, y, width, _ := b.GetInnerRect()
	Draw(screen, HeaderLine[:width*3], x, y, StylePanelHeader) // *3 because each character is 3 bytes long
	Draw(screen, " "+b.title+" ", x+4, y, StylePanelHeader)
	x = x + width - 4
	for i := len(b.labels) - 1; i >= 0; i-- {
		label := b.labels[i]
		x -= len(label) + 1
		b.offsets[i] = x
		Draw(screen, " "+label+" ", x-1, y, StylePanelHeader)
	}
}

// CheckKey checks to see if an input key event is one of the keys associated
// with the buttons.  It returns true if the key was handled.
func (b *HeaderBar) CheckKey(evt *tcell.EventKey, setFocus func(tview.Primitive)) bool {
	for i, e := range b.evts {
		if e.Key() == tcell.KeyRune && evt.Key() == tcell.KeyRune {
			r := evt.Rune()
			if r >= 'a' && r <= 'z' {
				r -= 'a' - 'A'
			}
			if r == e.Rune() {
				if b.onButton != nil {
					b.onButton(b.keys[i], setFocus)
				}
				return true
			}
		} else if e.Key() != tcell.KeyRune && e.Key() == evt.Key() {
			if b.onButton != nil {
				b.onButton(b.keys[i], setFocus)
			}
			return true
		}
	}
	return false
}

// CheckMouse checks to see if an input mouse event was a click on one of the
// keys.  It returns true if the click was handled.
func (b *HeaderBar) CheckMouse(act tview.MouseAction, evt *tcell.EventMouse, setFocus func(tview.Primitive)) bool {
	x, y := evt.Position()
	if bx, by, bw, _ := b.GetRect(); x < bx || x > bx+bw || y != by || act != tview.MouseLeftClick {
		return false
	}
	for i, off := range b.offsets {
		label := b.labels[i]
		if x >= off && x < off+len(label) {
			if b.onButton != nil {
				b.onButton(b.keys[i], setFocus)
			}
			return true
		}
	}
	return false
}
