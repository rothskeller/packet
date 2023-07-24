package msgedit

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/rothskeller/packet/envelope"
	"github.com/rothskeller/packet/message"
)

// RunEditor runs the message editor on the provided message.
func RunEditor(title string, env *envelope.Envelope, msg message.IEdit) {
	var (
		mgr   manager
		medit *messageEditor
	)
	medit = newEditor(title, env, msg, &mgr)
	mgr.pages = tview.NewPages().AddPage("main", medit, true, true)
	mgr.app = tview.NewApplication().SetRoot(mgr.pages, true).EnableMouse(true)
	mgr.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		InitStyles(screen)
		mgr.app.SetBeforeDrawFunc(nil)
		return false
	})
	mgr.app.Run()
}

type manager struct {
	app               *tview.Application
	pages             *tview.Pages
	dialogFocusReturn tview.Primitive
}

func (m *manager) ScreenSize() (width, height int) {
	_, _, width, height = m.pages.GetRect()
	return
}
func (m *manager) OpenDialog(p tview.Primitive) {
	m.dialogFocusReturn = m.app.GetFocus()
	m.pages.AddPage("dialog", p, true, true)
}
func (m *manager) CloseDialog() {
	m.pages.RemovePage("dialog")
	m.app.SetFocus(m.dialogFocusReturn)
}
func (m *manager) RemovePanel(tview.Primitive) {
	m.app.Stop()
}
func (m *manager) SetFocus(p tview.Primitive) {
	m.app.SetFocus(p)
}
