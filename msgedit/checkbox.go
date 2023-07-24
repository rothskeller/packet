package msgedit

/*
type meCheckbox struct {
	*tview.Box
	meFieldManager
	f        *xscmsg.Field
	msg      *xscmsg.Message
	col1     int
	selected bool
	finished func(tcell.Key) tview.Primitive
}

func newMECheckbox(f *xscmsg.Field, msg *xscmsg.Message, manager meFieldManager) meditField {
	return &meCheckbox{Box: tview.NewBox(), f: f, msg: msg, meFieldManager: manager}
}

func (mf *meCheckbox) getField() *xscmsg.Field { return mf.f }

func (mf *meCheckbox) getLines(width int) int { return 1 }

func (mf *meCheckbox) setFieldNameWidth(w int) { mf.col1 = w }

func (mf *meCheckbox) setSelected(sel bool) { mf.selected = sel }

func (mf *meCheckbox) Draw(screen tcell.Screen) {
	mf.DrawForSubclass(screen, mf)
	x, y, width, _ := mf.GetInnerRect()
	label := mf.f.Def.Label
	if label == "" {
		label = mf.f.Def.Tag
	}
	labelStyle := StyleFieldName
	if mf.selected {
		labelStyle = StyleFieldNameSelected
	}
	DrawWidth(screen, label, x, y, mf.col1+2, labelStyle)
	x, width = x+mf.col1+2, width-mf.col1-2
	frameStyle := StyleCheckboxFrame
	if mf.HasFocus() {
		frameStyle = MakeSelected(frameStyle)
	}
	screen.SetContent(x, y, '[', nil, frameStyle)
	screen.SetContent(x+2, y, ']', nil, frameStyle)
	if mf.f.Value != "" {
		screen.SetContent(x+1, y, 'X', nil, StyleCheckboxCheck)
	}
	if mf.selected && width >= 23 {
		Draw(screen, "(space to toggle)", x+6, y, StyleCheckboxHint)
	}
	if mf.HasFocus() {
		screen.ShowCursor(x+1, y)
	}
}

func (mf *meCheckbox) InputHandler() func(*tcell.EventKey, func(tview.Primitive)) {
	return mf.WrapInputHandler(func(evt *tcell.EventKey, setFocus func(tview.Primitive)) {
		switch evt.Key() {
		case tcell.KeyF1:
			mf.showHelp()
		case tcell.KeyEnter, tcell.KeyTab, tcell.KeyDown:
			mf.FieldFinished(tcell.KeyTab)
		case tcell.KeyBacktab, tcell.KeyUp:
			mf.FieldFinished(tcell.KeyBacktab)
		case tcell.KeyEsc, tcell.KeyF10:
			mf.FieldFinished(evt.Key())
		case tcell.KeyRune:
			switch evt.Rune() {
			case ' ':
				mf.toggle()
			}
		}
	})
}

func (mf *meCheckbox) MouseHandler() func(tview.MouseAction, *tcell.EventMouse, func(tview.Primitive)) (bool, tview.Primitive) {
	return mf.WrapMouseHandler(func(act tview.MouseAction, evt *tcell.EventMouse, setFocus func(tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, y, _, _ := mf.GetInnerRect()
		mx, my := evt.Position()
		if my == y && mx == x+mf.col1+3 && act == tview.MouseLeftClick {
			mf.toggle()
			return true, nil
		}
		if act == tview.MouseLeftDown && mf.InRect(mx, my) {
			setFocus(mf)
			return true, nil
		}
		return
	})
}

func (mf *meCheckbox) toggle() {
	if mf.f.Value == "" {
		mf.f.Value = "checked"
	} else {
		mf.f.Value = ""
	}
}

func (mf *meCheckbox) showHelp() {
	label := mf.f.Def.Label
	if label == "" {
		label = mf.f.Def.Tag
	}
	helpText := fmt.Sprintf("You are editing the value of the %q field of %s %s.", label, mf.msg.Type.Article, mf.msg.Type.Name)
	helpText = strings.Join(Wrap(helpText, 56), "\n")
	helpText += `

To toggle the value of this field, press the space bar.
When finished editing this field, press:
    Enter or Tab    to move to the next field
    Shift-Tab       to move to the previous field
When you are done editing the message, press:
    F10             to queue the message to be sent
    Esc             to leave the message as a draft`
	ShowHelp(mf.Box, "Field Entry Help", helpText, true, mf)
}
*/
