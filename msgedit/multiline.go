package msgedit

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rothskeller/packet/message"
)

// meMultiLine is a multi-line input control for a field.  It supports online
// help for the field.
type meMultiLine struct {
	*tview.Box
	// Parameters passed to the control:
	meFieldManager
	f *message.EditField
	// Data given to the control before drawing:
	col1     int
	selected bool
	// Others:
	area  *tview.TextArea
	lines int
}

// newMEMultiLine creates a new multi-line input control for the specified field of
// the specified message.  The control calls the supplied finished function
// whenever the user hits a key to leave the field; the key passed to that
// function will be one of Tab, Backtab, Esc, or F10.
func newMEMultiLine(f *message.EditField, manager meFieldManager) meditField {
	mf := &meMultiLine{Box: tview.NewBox(), f: f, meFieldManager: manager}
	if mf.f.Value != "" {
		mf.createArea()
	}
	return mf
}

func (mf *meMultiLine) getField() *message.EditField { return mf.f }

func (mf *meMultiLine) getLines(width int) int {
	if mf.area == nil {
		return 1
	}
	value := mf.area.GetText()
	mf.lines = strings.Count(value, "\n") + 2
	return mf.lines
}

func (mf *meMultiLine) setFieldNameWidth(w int) { mf.col1 = w }

func (mf *meMultiLine) setSelected(sel bool) { mf.selected = sel }

func (mf *meMultiLine) Draw(screen tcell.Screen) {
	mf.DrawForSubclass(screen, mf)
	x, y, width, height := mf.GetInnerRect()
	// Draw the label.
	labelStyle := StyleFieldName
	switch {
	case mf.selected && mf.f.Problem != "":
		labelStyle = StyleFieldNameInvalidSelected
	case mf.selected:
		labelStyle = StyleFieldNameSelected
	case mf.f.Problem != "":
		labelStyle = StyleFieldNameInvalid
	}
	DrawWidth(screen, mf.f.Label, x, y, mf.col1+2, labelStyle)
	// If we don't have a text area, the field is empty and has been since
	// this control was created.  Display help text and stop.
	if mf.area == nil {
		if mf.selected {
			Draw(screen, "Press ENTER to edit", x+mf.col1+2, y, StyleInputHint)
		}
		return
	}
	x, width = x+4, width-4
	y, height = y+1, height-1
	// When someone hits Enter at the end of the text area, it scrolls up
	// before we have a chance to grow the area to accommodate the new line.
	// Then, once we have the new line, it's not at the bottom because of
	// the scroll.  So right now, check to be sure that the area is not
	// overscrolled.
	rs, cs := mf.area.GetOffset()
	if maxScroll := (mf.lines - 1) - height; rs > maxScroll {
		mf.area.SetOffset(maxScroll, cs)
		rs = maxScroll
	}
	mf.area.SetRect(x, y, width, height)
	mf.area.Draw(screen)
}

func (mf *meMultiLine) createArea() {
	mf.area = tview.NewTextArea().
		SetSelectedStyle(StyleInputSelected).
		SetText(mf.f.Value, false).
		SetTextStyle(StyleInput).
		SetWrap(false)
	mf.area.SetFocusFunc(func() {
		mf.area.SetTextStyle(StyleInputFocus).
			SetSelectedStyle(StyleInputSelected)
	}).SetBlurFunc(func() {
		mf.area.SetTextStyle(StyleInput).
			SetSelectedStyle(StyleInput)
	})
	mf.area.SetFinishedFunc(mf.onFinished)
}

// Pass the focus on to the area when we receive it.
func (mf *meMultiLine) Focus(delegate func(tview.Primitive)) {
	if mf.area != nil {
		delegate(mf.area)
	} else {
		mf.Box.Focus(delegate)
	}
}

func (mf *meMultiLine) HasFocus() bool {
	if mf.area != nil {
		return mf.area.HasFocus()
	}
	return mf.Box.HasFocus()
}

func (mf *meMultiLine) selectStart() {
	if mf.area != nil {
		mf.area.Select(0, 0)
	}
}

func (mf *meMultiLine) selectEnd() {
	if mf.area != nil {
		mf.area.Select(99999, 99999)
	}
}

func (mf *meMultiLine) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return mf.WrapInputHandler(func(evt *tcell.EventKey, setFocus func(tview.Primitive)) {
		if mf.area == nil {
			switch evt.Key() {
			case tcell.KeyF1:
				mf.f.Value = strings.TrimRight(mf.area.GetText(), "\n")
				mf.applyEdits()
				mf.showHelp()
			case tcell.KeyTab, tcell.KeyDown:
				mf.FieldFinished(tcell.KeyTab)
			case tcell.KeyBacktab, tcell.KeyUp:
				mf.FieldFinished(tcell.KeyBacktab)
			case tcell.KeyEsc:
				mf.FieldFinished(tcell.KeyEsc)
			case tcell.KeyF10:
				mf.FieldFinished(tcell.KeyF10)
			case tcell.KeyEnter:
				mf.createArea()
				setFocus(mf.area)
			}
			return
		}
		switch evt.Key() {
		case tcell.KeyUp:
			if evt.Modifiers() == 0 {
				if fr, _, _, _ := mf.area.GetCursor(); fr == 0 {
					mf.onFinished(tcell.KeyBacktab)
					return
				}
			}
		case tcell.KeyDown:
			if evt.Modifiers() == 0 {
				// This is expensive, but...
				t := mf.area.GetText()
				lr := strings.Count(t, "\n")
				if _, _, tr, _ := mf.area.GetCursor(); tr == lr {
					mf.onFinished(tcell.KeyTab)
					return
				}
			}
		case tcell.KeyF10:
			mf.onFinished(tcell.KeyF10)
			return
		}
		mf.area.InputHandler()(evt, setFocus)
	})
}

func (mf *meMultiLine) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return mf.WrapMouseHandler(func(act tview.MouseAction, evt *tcell.EventMouse, setFocus func(tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if mf.area != nil {
			if consumed, capture = mf.area.MouseHandler()(act, evt, setFocus); consumed {
				return
			}
		}
		if mf.InRect(evt.Position()) && act == tview.MouseLeftDown {
			if mf.area == nil {
				mf.createArea()
			}
			setFocus(mf.area)
			return true, nil
		}
		return
	})
}

func (mf *meMultiLine) onFinished(key tcell.Key) {
	mf.f.Value = strings.TrimRight(mf.area.GetText(), "\n")
	switch key {
	case tcell.KeyTab, tcell.KeyDown:
		mf.FieldFinished(tcell.KeyTab)
	case tcell.KeyBacktab, tcell.KeyUp:
		mf.FieldFinished(tcell.KeyBacktab)
	case tcell.KeyEsc:
		mf.FieldFinished(tcell.KeyEsc)
	case tcell.KeyF10:
		mf.FieldFinished(tcell.KeyF10)
	}
}

func (mf *meMultiLine) showHelp() {
	helpText := mf.f.Help
	helpText = strings.Join(Wrap(helpText, 56), "\n")
	helpText += `

While editing this field, you can press:
    Ctrl-A or Home   to move to the start of the line
    Ctrl-B or PgUp   to move up one page
    Ctrl-E or End    to move to the end of the line
    Ctrl-D or Del    to erase the next character
    Ctrl-F or PgDn   to move down one page
    Ctrl-K           to erase to the end of the line
    Ctrl-L           to select all text
    Ctrl-Q           to copy selected text to clipboard
    Ctrl-U           to erase the entire line
    Ctrl-V           to paste clipboard into text
    Ctrl-W           to erase previous word
    Ctrl-X           to cut selected text to clipboard
    Ctrl-Y           redo the last undone change
    Ctrl-Z           undo the last change
    Alt-B or Ctrl-←  to move left one word
    Alt-F or Ctrl-→  to move right one word
    Shift-Arrow      to select text
When finished editing this field, press:
    Tab             to move to the next field
    Shift-Tab       to move to the previous field
When you are done editing the message, press:
    F10             to queue the message to be sent
    Esc             to leave the message as a draft`
	ShowHelp(mf.Box, "Field Entry Help", helpText, true, mf)
}
