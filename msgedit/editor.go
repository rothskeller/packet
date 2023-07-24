package msgedit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
	"github.com/rothskeller/packet/message/plaintext"
)

// messageEditor is the message editing panel.
type messageEditor struct {
	*tview.Box
	*manager
	env      *envelope.Envelope
	msg      message.IEdit
	header   *HeaderBar
	fakescr  *fakeScreen
	tofield  *message.EditField
	lmifield *message.EditField
	fields   []meditField
	offsets  []int
	scroll   int
	selected int
	loose    bool // true when selected item does not have to be in scroll
}
type meditField interface {
	tview.Primitive
	getField() *message.EditField
	getLines(int) int
	setFieldNameWidth(int)
	setSelected(bool)
}

// newEditor returns a new message editing panel.
func newEditor(lmi string, env *envelope.Envelope, msg message.IEdit, manager *manager) *messageEditor {
	var title string

	me := &messageEditor{
		Box:     tview.NewBox(),
		env:     env,
		msg:     msg,
		manager: manager,
		fakescr: new(fakeScreen),
	}
	// Format the title.
	if _, ok := msg.(*plaintext.PlainText); ok {
		if title = lmi; title == "" {
			title = "New Message"
		}
	} else if lmi == "" {
		title = "New " + msg.Type().Tag
	} else {
		title = lmi + " (" + msg.Type().Tag + ")"
	}
	// Set up the header bar.
	me.header = NewHeaderBar(title).SetOnButton(me.onButton)
	// Set up an editor field for the To: line.
	me.applyToField()
	me.fields = append(me.fields, newMEInput(me.tofield, me))
	// Set up all of the editor fields.
	for _, f := range msg.EditFields() {
		if f.Multiline {
			me.fields = append(me.fields, newMEMultiLine(f, me))
		} else {
			me.fields = append(me.fields, newMEInput(f, me))
		}
		if f.LocalMessageID {
			me.lmifield = f
		}
	}
	if me.lmifield.Problem != "" {
		me.header.SetButtons()
	} else if me.tofield.Problem != "" {
		me.header.SetButtons("ESC", "Draft")
	} else {
		me.header.SetButtons("ESC", "Draft", "F10", "Send")
	}
	// Calculate the maximum field name length.
	var col1 int
	for _, f := range me.fields {
		label := f.getField().Label
		if l := len(label); l > col1 {
			col1 = l
		}
	}
	for _, f := range me.fields {
		f.setFieldNameWidth(col1)
	}
	me.offsets = make([]int, len(me.fields)+1)
	return me
}

// Draw draws the message editing panel.
func (me *messageEditor) Draw(screen tcell.Screen) {
	me.DrawForSubclass(screen, me)
	x, y, width, height := me.GetInnerRect()
	me.header.SetRect(x, y, width, 1)
	me.header.Draw(screen)
	me.drawMessage(screen, x, y+1, width, height-1)
}

func (me *messageEditor) drawMessage(screen tcell.Screen, x, y, width, height int) {
	// Find out how many lines each field wants to have, and set the
	// necessary size for the fake screen.
	for i, f := range me.fields {
		f.setSelected(i == me.selected)
		me.offsets[i+1] = me.offsets[i] + f.getLines(width)
	}
	me.fakescr.backing = screen
	me.fakescr.setSize(width, me.offsets[len(me.fields)])
	// Draw the fields into the fake screen.
	for i, f := range me.fields {
		f.SetRect(0, me.offsets[i], width, me.offsets[i+1]-me.offsets[i])
		f.Draw(me.fakescr)
	}
	// If we're doing strict scrolling — i.e., last change was something
	// other than a mouse wheel scroll — then make sure the currently
	// selected item is scrolled into view.
	if !me.loose {
		var currentLine = me.offsets[me.selected]
		if mf, ok := me.fields[me.selected].(*meMultiLine); ok && mf.area != nil {
			_, _, tr, _ := mf.area.GetCursor()
			currentLine += tr + 1
		}
		if me.scroll > currentLine {
			me.scroll = currentLine
		}
		if me.scroll+height <= currentLine {
			me.scroll = currentLine - height + 1
		}
	}
	// Make sure scroll is sane.
	if me.scroll > me.offsets[len(me.fields)]-height {
		me.scroll = me.offsets[len(me.fields)] - height
	}
	if me.scroll < 0 {
		me.scroll = 0
	}
	// Draw the fake screen onto the real screen at the appropriate place.
	dh := height
	if dh > me.offsets[len(me.fields)] {
		dh = me.offsets[len(me.fields)]
	}
	me.fakescr.copy(screen, 0, me.scroll, x, y, width, dh)
	// Set the correct cursor location.
	if me.fakescr.cursorX < 0 || me.fakescr.cursorY < me.scroll || me.fakescr.cursorY >= me.scroll+height {
		screen.HideCursor()
	} else {
		screen.ShowCursor(me.fakescr.cursorX+x, me.fakescr.cursorY+y-me.scroll)
	}
}

// Focus receives focus on the message editing panel.
func (me *messageEditor) Focus(delegate func(tview.Primitive)) {
	delegate(me.fields[me.selected])
}

// HasFocus returns whether the message editing panel has focus.
func (me *messageEditor) HasFocus() bool {
	return me.fields[me.selected].HasFocus()
}

// InputHandler handles keystrokes in the message editing panel.
func (me *messageEditor) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	me.loose = false
	return me.WrapInputHandler(func(evt *tcell.EventKey, setFocus func(tview.Primitive)) {
		me.fields[me.selected].InputHandler()(evt, setFocus)
	})
}

// MouseHandler handles mouse actions in the message editing panel.
func (me *messageEditor) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return me.WrapMouseHandler(func(act tview.MouseAction, evt *tcell.EventMouse, setFocus func(tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if me.header.CheckMouse(act, evt, setFocus) {
			return true, nil
		}
		mx, my := evt.Position()
		if !me.InRect(mx, my) {
			return
		}
		x, y, _, _ := me.GetInnerRect()
		if my > y {
			switch act {
			case tview.MouseScrollUp:
				me.scroll--
				me.loose = true
				return true, nil
			case tview.MouseScrollDown:
				me.scroll++
				me.loose = true
				return true, nil
			}
		}
		if !me.HasFocus() {
			return
		}
		adjme := tcell.NewEventMouse(mx-x, my-y-1+me.scroll, evt.Buttons(), evt.Modifiers())
		for i, f := range me.fields {
			if consumed, capture = f.MouseHandler()(act, adjme, setFocus); consumed {
				me.loose = false
				if f.HasFocus() {
					me.selected = i
				}
				if capture != nil {
					capture = me
				}
				return
			}
		}
		return
	})
}

func (me *messageEditor) onButton(button string, setFocus func(tview.Primitive)) {
	me.applyToField()
	me.msg.ApplyEdits()
	switch button {
	case "ESC":
		me.saveAsDraft()
		me.loose = false
	case "F10":
		me.saveAsReady()
		me.loose = false
	}
}

// applyEdits applies edits to the fields of the message.
func (me *messageEditor) applyEdits() {
	me.applyToField()
	me.msg.ApplyEdits()
}

// FieldFinished is called by a field when it receives a key that should
// transfer focus outside the field.  The key passed to it is always Tab,
// Backtab, Escape, or F10; other equivalent keys are translated to those.
func (me *messageEditor) FieldFinished(key tcell.Key) {
	me.applyEdits()
	if me.lmifield.Problem != "" {
		me.header.SetButtons()
	} else if me.tofield.Problem != "" {
		me.header.SetButtons("ESC", "Draft")
	} else {
		me.header.SetButtons("ESC", "Draft", "F10", "Send")
	}
	switch key {
	case tcell.KeyTab:
		if me.selected < len(me.fields)-1 {
			me.selected++
		} else {
			me.selected = 0
		}
		f := me.fields[me.selected]
		if f, ok := f.(*meMultiLine); ok {
			f.selectStart()
		}
		me.SetFocus(f)
	case tcell.KeyBacktab:
		if me.selected > 0 {
			me.selected--
		} else {
			me.selected = len(me.fields) - 1
		}
		f := me.fields[me.selected]
		if f, ok := f.(*meMultiLine); ok {
			f.selectEnd()
		}
		me.SetFocus(f)
	case tcell.KeyEsc:
		me.saveAsDraft()
	case tcell.KeyF10:
		me.saveAsReady()
	}
}

func (me *messageEditor) saveAsDraft() {
	if me.lmifield.Problem != "" {
		return
	}
	me.env.ReadyToSend = false
	me.RemovePanel(me)
}

func (me *messageEditor) saveAsReady() {
	if me.lmifield.Problem != "" || me.tofield.Problem != "" {
		return
	}
	valid := true
	for _, f := range me.fields {
		if f.getField().Problem != "" {
			valid = false
			break
		}
	}
	if valid {
		me.env.ReadyToSend = true
		me.RemovePanel(me)
		return
	}
	modal := tview.NewModal().
		SetText("This message has invalid fields.  Are you sure you want to send it this way?").
		AddButtons([]string{"Save as Draft", "Queue for Sending"}).
		SetDoneFunc(func(button int, _ string) {
			me.CloseDialog()
			me.env.ReadyToSend = button == 1
			me.RemovePanel(me)
		})
	modal.SetTitle("Invalid Message")
	me.OpenDialog(modal)
}
