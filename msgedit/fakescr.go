package msgedit

import "github.com/gdamore/tcell/v2"

// A fakeScreen is a tcell.Screen into which controls can draw, which can then
// be copied onto a real Screen at a specified origin and size.  This is used to
// implement drawing with clipping.
type fakeScreen struct {
	backing tcell.Screen
	width   int
	height  int
	cells   []byte
	styles  []tcell.Style
	cursorX int
	cursorY int
}

func (f *fakeScreen) setSize(w, h int) {
	f.width, f.height = w, h
	f.cursorX, f.cursorY = -1, -1
	if w*h == len(f.cells) {
		return
	}
	f.cells = make([]byte, w*h)
	f.styles = make([]tcell.Style, w*h)
}

func (f *fakeScreen) copy(to tcell.Screen, sourceX, sourceY, destX, destY, width, height int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			ch := rune(f.cells[(sourceY+y)*f.width+sourceX+x])
			st := f.styles[(sourceY+y)*f.width+sourceX+x]
			to.SetContent(destX+x, destY+y, ch, nil, st)
		}
	}
}

// Most of the Screen methods are not needed and not implemented.
func (f *fakeScreen) Init() error                        { panic("not implemented") }
func (f *fakeScreen) Fini()                              { panic("not implemented") }
func (f *fakeScreen) Clear()                             { panic("not implemented") }
func (f *fakeScreen) Fill(_ rune, _ tcell.Style)         { panic("not implemented") }
func (f *fakeScreen) SetStyle(style tcell.Style)         { panic("not implemented") }
func (f *fakeScreen) SetCursorStyle(_ tcell.CursorStyle) { panic("not implemented") }
func (f *fakeScreen) ChannelEvents(ch chan<- tcell.Event, quit <-chan struct{}) {
	panic("not implemented")
}
func (f *fakeScreen) PollEvent() tcell.Event                    { panic("not implemented") }
func (f *fakeScreen) HasPendingEvent() bool                     { panic("not implemented") }
func (f *fakeScreen) PostEvent(ev tcell.Event) error            { panic("not implemented") }
func (f *fakeScreen) PostEventWait(ev tcell.Event)              { panic("not implemented") }
func (f *fakeScreen) EnableMouse(_ ...tcell.MouseFlags)         { panic("not implemented") }
func (f *fakeScreen) DisableMouse()                             { panic("not implemented") }
func (f *fakeScreen) EnablePaste()                              { panic("not implemented") }
func (f *fakeScreen) DisablePaste()                             { panic("not implemented") }
func (f *fakeScreen) Show()                                     { panic("not implemented") }
func (f *fakeScreen) Sync()                                     { panic("not implemented") }
func (f *fakeScreen) RegisterRuneFallback(r rune, subst string) { panic("not implemented") }
func (f *fakeScreen) UnregisterRuneFallback(r rune)             { panic("not implemented") }
func (f *fakeScreen) Resize(_ int, _ int, _ int, _ int)         { panic("not implemented") }
func (f *fakeScreen) Suspend() error                            { panic("not implemented") }
func (f *fakeScreen) Resume() error                             { panic("not implemented") }
func (f *fakeScreen) Beep() error                               { panic("not implemented") }
func (f *fakeScreen) SetSize(_ int, _ int)                      { panic("not implemented") }

// Some methods (capability queries) are passed on to the backing screen.
func (f *fakeScreen) HasMouse() bool          { return f.backing.HasMouse() }
func (f *fakeScreen) Colors() int             { return f.backing.Colors() }
func (f *fakeScreen) CharacterSet() string    { return f.backing.CharacterSet() }
func (f *fakeScreen) HasKey(k tcell.Key) bool { return f.backing.HasKey(k) }

// The remainder of the methods are actually implemented.
func (f *fakeScreen) SetCell(x int, y int, style tcell.Style, ch ...rune) {
	f.SetContent(x, y, ch[0], ch[1:], style)
}
func (f *fakeScreen) GetContent(x int, y int) (primary rune, combining []rune, style tcell.Style, width int) {
	ch := rune(f.cells[y*f.width+x])
	st := f.styles[y*f.width+x]
	return ch, nil, st, 1
}
func (f *fakeScreen) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	f.cells[y*f.width+x] = byte(primary)
	f.styles[y*f.width+x] = style
}
func (f *fakeScreen) CanDisplay(r rune, checkFallbacks bool) bool { return r >= 32 && r < 127 }
func (f *fakeScreen) Size() (width int, height int)               { return f.width, f.height }
func (f *fakeScreen) ShowCursor(x int, y int)                     { f.cursorX, f.cursorY = x, y }
func (f *fakeScreen) HideCursor()                                 { f.cursorX, f.cursorY = -1, -1 }
