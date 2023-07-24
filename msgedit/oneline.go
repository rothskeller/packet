package msgedit

import (
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rothskeller/packet/message"
)

type meFieldManager interface {
	DialogManager
	applyEdits()
	FieldFinished(tcell.Key)
}

// meInput is a single-line input control for a field.  It supports online
// help for the field, autocomplete of allowed/suggested choices, and selection
// from a narrow list of allowed/suggested choices.
type meInput struct {
	*tview.Box
	// Parameters passed to the control:
	meFieldManager
	f *message.EditField
	// Data given to the control before drawing:
	col1     int
	selected bool
	// Others:
	showingChoices bool // whether discrete choice list is showing
	hastyped       bool // whether the user has typed anything in the input
	selstart       int  // character on which the selection starts
	selend         int  // character before which the selection ends
	cursor         int  // location of cursor (always == selstart or selend)
	scroll         int  // amount of horizontal scroll
}

// newMEInput creates a new single-line input control for the specified field of
// the specified message.  The control calls the supplied finished function
// whenever the user hits a key to leave the field; the key passed to that
// function will be one of Tab, Backtab, Esc, or F10.
func newMEInput(f *message.EditField, manager meFieldManager) meditField {
	mf := &meInput{Box: tview.NewBox(), f: f, meFieldManager: manager}
	mf.selend = len(mf.f.Value)
	mf.cursor = mf.selend
	return mf
}

func (mf *meInput) getField() *message.EditField { return mf.f }

func (mf *meInput) getLines(width int) int { return 1 }

func (mf *meInput) setFieldNameWidth(w int) { mf.col1 = w }

func (mf *meInput) setSelected(sel bool) { mf.selected = sel }

func (mf *meInput) Draw(screen tcell.Screen) {
	mf.DrawForSubclass(screen, mf)
	x, y, width, _ := mf.GetInnerRect()
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
	x, width = x+mf.col1+2, width-mf.col1-2
	// Draw the comment if any.
	if mf.f.Value == "" && mf.f.Hint != "" && mf.selected {
		Draw(screen, "("+mf.f.Hint+")", x+width-len(mf.f.Hint)-2, y, StyleInputHint)
		width -= len(mf.f.Hint) + 4
	}
	// Special case for small restricted lists.  If we have focus, and the
	// input is either empty or one of the valid choices (selected from the
	// restricted list, not typed), and there is room to display the
	// restricted list, do so.
	if mf.HasFocus() && !mf.hastyped && len(mf.f.Choices) != 0 && (mf.f.Value == "" || inList(mf.f.Choices, mf.f.Value)) {
		choicesLen := 5
		for _, c := range mf.f.Choices {
			choicesLen += len(c) + 2
		}
		if width >= choicesLen {
			mf.drawChoices(screen, x, y, width)
			return
		}
	}
	mf.showingChoices = false
	// Draw input area.
	inputStyle := StyleInput
	if mf.HasFocus() {
		inputStyle = StyleInputFocus
	}
	if _, inputBg, _ := inputStyle.Decompose(); inputBg != tcell.ColorDefault {
		for index := 0; index < width; index++ {
			screen.SetContent(x+index, y, ' ', nil, inputStyle)
		}
	}
	// Make sure the selection is sane.
	if mf.selstart < 0 {
		mf.selstart = 0
	}
	if mf.selstart > len(mf.f.Value) {
		mf.selstart = len(mf.f.Value)
	}
	if mf.selend < mf.selstart {
		mf.selend = mf.selstart
	}
	if mf.selend > len(mf.f.Value) {
		mf.selend = len(mf.f.Value)
	}
	// Make sure the cursor is sane.
	if mf.cursor < mf.selstart {
		mf.cursor = mf.selstart
	}
	if mf.cursor > mf.selend {
		mf.cursor = mf.selend
	}
	// Make sure the scroll includes the cursor.
	if mf.scroll+width <= mf.cursor {
		mf.scroll = mf.cursor - width + 1
	}
	if mf.scroll > mf.cursor {
		mf.scroll = mf.cursor
	}
	// Make sure the scroll is sane.
	if mf.scroll > len(mf.f.Value)-width+1 {
		mf.scroll = len(mf.f.Value) - width + 1
	}
	if mf.scroll < 0 {
		mf.scroll = 0
	}
	// Draw entered text.
	DrawWidth(screen, mf.f.Value[mf.scroll:], x, y, width, inputStyle)
	// Draw selection and cursor if focused.
	if mf.HasFocus() {
		for i := mf.selstart; i < mf.selend; i++ {
			x := x + i - mf.scroll
			c, _, _, _ := screen.GetContent(x, y)
			screen.SetContent(x, y, c, nil, StyleInputSelected)
		}
		screen.ShowCursor(x+mf.cursor-mf.scroll, y)
	}
}

// drawChoices draws the input control in the restricted choices mode.  This
// mode appears when the control has focus and there is room to display all of
// the choices for the field.  However, it is not displayed if the field has a
// value that isn't one of the choices, or if the value was typed during the
// current focus session rather than selected from the list.
func (mf *meInput) drawChoices(screen tcell.Screen, x, y, width int) {
	mf.cursor, mf.selstart, mf.selend = 0, 0, 0
	// Draw input area.
	for index := 0; index < 4; index++ {
		screen.SetContent(x+index, y, ' ', nil, StyleInputFocus)
	}
	if mf.f.Value == "" {
		screen.ShowCursor(x, y)
		screen.SetContent(x+5, y, '→', nil, StyleInputOptionArrow)
	} else {
		screen.HideCursor()
		screen.SetContent(x+5, y, '←', nil, StyleInputOptionArrow)
	}
	x += 7
	// Draw the choices.
	for _, c := range mf.f.Choices {
		style := StyleInputOption
		if mf.f.Value == c {
			style = StyleInputOptionSelected
		}
		Draw(screen, c, x, y, style)
		x += len(c) + 2
	}
	mf.showingChoices = true
	mf.cursor, mf.selstart, mf.selend = 0, 0, 0
}

// Select the entire field contents when receiving focus.
func (mf *meInput) Focus(delegate func(tview.Primitive)) {
	mf.Box.Focus(delegate)
	mf.selstart = 0
	mf.selend = len(mf.f.Value)
	mf.cursor = mf.selend
}

// Validate and auto-correct the field value when losing focus.  Also reset the
// per-focus-session flag indicating whether the user typed the current value.
func (mf *meInput) Blur() {
	mf.Box.Blur()
	mf.applyEdits()
	mf.hastyped = false
}

func (mf *meInput) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return mf.WrapInputHandler(func(evt *tcell.EventKey, setFocus func(tview.Primitive)) {
		switch evt.Key() {
		case tcell.KeyF1:
			mf.applyEdits()
			mf.showHelp()
		case tcell.KeyEnter, tcell.KeyTab, tcell.KeyDown:
			mf.FieldFinished(tcell.KeyTab)
		case tcell.KeyBacktab, tcell.KeyUp:
			mf.FieldFinished(tcell.KeyBacktab)
		case tcell.KeyEsc:
			mf.FieldFinished(tcell.KeyEsc)
		case tcell.KeyF10:
			mf.FieldFinished(tcell.KeyF10)
		case tcell.KeyHome:
			mf.selstart = 0
			mf.cursor = 0
			if evt.Modifiers()&tcell.ModShift == 0 {
				mf.selend = 0
			}
		case tcell.KeyEnd:
			mf.selend = len(mf.f.Value)
			mf.cursor = mf.selend
			if evt.Modifiers()&tcell.ModShift == 0 {
				mf.selstart = mf.selend
			}
		case tcell.KeyLeft:
			if evt.Modifiers()&tcell.ModShift != 0 { // Shift-Left extends selection
				if mf.cursor == mf.selstart && mf.selstart > 0 {
					mf.selstart--
					mf.cursor--
				} else if mf.cursor == mf.selend && mf.selend > mf.selstart {
					mf.selend--
				}
			} else if mf.showingChoices { // In choices mode, choose previous choice
				cs := mf.f.Choices
				for i, c := range cs {
					if mf.f.Value == c {
						if i == 0 {
							mf.f.Value = ""
						} else {
							mf.f.Value = cs[i-1]
						}
						break
					}
				}
			} else { // Otherwise move cursor left
				if mf.cursor > 0 {
					mf.cursor--
				}
				mf.selstart = mf.cursor
				mf.selend = mf.cursor
			}
		case tcell.KeyRight:
			if evt.Modifiers()&tcell.ModShift != 0 { // Shift-Right extends selection
				if mf.cursor == mf.selend && mf.selend < len(mf.f.Value) {
					mf.selend++
					mf.cursor++
				} else if mf.cursor == mf.selstart && mf.selstart < mf.selend {
					mf.selstart++
				}
			} else if mf.showingChoices { // In choices mode, choose next choice
				cs := mf.f.Choices
				if mf.f.Value == "" {
					mf.f.Value = cs[0]
				} else {
					for i, c := range cs {
						if mf.f.Value == c && i < len(cs)-1 {
							mf.f.Value = cs[i+1]
							break
						}
					}
				}
			} else { // Otherwise move cursor right
				if mf.cursor < len(mf.f.Value) {
					mf.cursor++
				}
				mf.selstart = mf.cursor
				mf.selend = mf.cursor
			}
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			if mf.showingChoices { // In choices mode, clear choice
				mf.f.Value = ""
			} else if mf.selstart != mf.selend { // Have selection, remove it
				mf.f.Value = mf.f.Value[:mf.selstart] + mf.f.Value[mf.selend:]
				mf.cursor = mf.selstart
				mf.selend = mf.selstart
				mf.hastyped = mf.f.Value != ""
			} else if mf.cursor != 0 { // Remove character before selection
				mf.f.Value = mf.f.Value[:mf.cursor-1] + mf.f.Value[mf.cursor:]
				mf.cursor--
				mf.selstart = mf.cursor
				mf.selend = mf.cursor
				mf.hastyped = mf.f.Value != ""
			}
		case tcell.KeyDelete, tcell.KeyCtrlD:
			if mf.showingChoices { // In choices mode, clear choice
				mf.f.Value = ""
			} else if mf.selstart != mf.selend { // Have selection, remove it
				mf.f.Value = mf.f.Value[:mf.selstart] + mf.f.Value[mf.selend:]
				mf.cursor = mf.selstart
				mf.selend = mf.selstart
				mf.hastyped = mf.f.Value != ""
			} else if mf.cursor < len(mf.f.Value) { // Remove character after cursor
				mf.f.Value = mf.f.Value[:mf.cursor] + mf.f.Value[mf.cursor+1:]
				mf.hastyped = mf.f.Value != ""
			}
		case tcell.KeyCtrlA:
			mf.selstart, mf.cursor, mf.selend = 0, 0, 0
		case tcell.KeyCtrlE:
			mf.selend = len(mf.f.Value)
			mf.cursor, mf.selstart = mf.selend, mf.selend
		case tcell.KeyCtrlK:
			mf.f.Value = mf.f.Value[:mf.cursor]
			mf.hastyped = mf.f.Value != ""
		case tcell.KeyCtrlU:
			mf.f.Value = ""
			mf.hastyped = false
		case tcell.KeyRune:
			if evt.Modifiers()&^tcell.ModShift != 0 {
				break
			}
			mf.hastyped = true
			if mf.selstart != mf.selend {
				mf.f.Value = mf.f.Value[:mf.selstart] + mf.f.Value[mf.selend:]
				mf.cursor = mf.selstart
			}
			mf.f.Value = mf.f.Value[:mf.cursor] + string(evt.Rune()) + mf.f.Value[mf.cursor:]
			mf.cursor++
			mf.selstart = mf.cursor
			mf.selend = mf.cursor
			if cs := mf.f.Choices; len(cs) != 0 && mf.cursor == len(mf.f.Value) {
				lower := strings.ToLower(mf.f.Value)
				for _, c := range cs {
					if strings.HasPrefix(strings.ToLower(c), lower) {
						mf.f.Value = c
						mf.selend = len(mf.f.Value)
						break
					}
				}
			}
		}
	})
}

func (mf *meInput) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return mf.WrapMouseHandler(func(act tview.MouseAction, evt *tcell.EventMouse, setFocus func(tview.Primitive)) (consumed bool, capture tview.Primitive) {
		mx, my := evt.Position()
		if !mf.InRect(mx, my) {
			return
		}
		x, _, _, _ := mf.GetInnerRect()
		if act == tview.MouseLeftDown && !mf.HasFocus() {
			setFocus(mf)
			return true, nil
		}
		if mf.showingChoices && act == tview.MouseLeftDown {
			mx -= x + mf.col1 + 2
			if mx >= 0 && mx < 4 {
				mf.f.Value = ""
				return true, nil
			}
			mx -= 7
			for _, c := range mf.f.Choices {
				if mx >= 0 && mx < len(c) {
					mf.f.Value = c
					return true, nil
				}
				mx -= len(c) + 2
			}
			return
		}
		if act == tview.MouseLeftClick {
			loc := mx - (x + mf.col1 + 2) + mf.scroll
			if loc >= 0 {
				if loc > len(mf.f.Value) {
					loc = len(mf.f.Value)
				}
				mf.selstart, mf.selend, mf.cursor = loc, loc, loc
				consumed = true
			}
		}
		return
	})
}

func (mf *meInput) showHelp() {
	helpText := strings.Join(Wrap(mf.f.Help, 56), "\n")
	helpText += `

While editing this field, you can press:
    Ctrl-A or Home  to move to the start of the line
    Ctrl-E or End   to move to the end of the line
    Ctrl-D or Del   to erase the next character
    Ctrl-K          to erase to the end of the line
    Ctrl-U          to erase the entire line
When finished editing this field, press:
    Enter, Tab, ↓   to move to the next field
    Shift-Tab, ↑    to move to the previous field
When you are done editing the message, press:
    F10             to queue the message to be sent
    Esc             to leave the message as a draft`
	ShowHelp(mf.Box, "Field Entry Help", helpText, true, mf)
}

func choiceList(choices []string) string {
	// We want to format the choices in a table, with as few rows as
	// possible, sorted phone book style.  Start by making sure the list is
	// sorted.
	sort.Strings(choices)
	var rows = []string{"    "}
OUTER:
	for { // repeat trying to format into successively higher row counts
		row := 0
		max := 0
		for _, choice := range choices {
			if row == len(rows) {
				for row := range rows {
					rows[row] = rows[row] + spaces[:max+3-len(rows[row])]
				}
				row = 0
			}
			rows[row] += choice
			if len(rows[row]) > max {
				max = len(rows[row])
			}
			if max > 56 {
				for row := range rows {
					rows[row] = "    "
				}
				rows = append(rows, "    ")
				continue OUTER
			}
			row++
		}
		return strings.Join(rows, "\n")
	}
}

func inList[T comparable](list []T, item T) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
