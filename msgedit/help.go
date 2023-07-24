package msgedit

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DialogManager is the interface that must be satisfied by the dialog manager
// passed to AddHelp or ShowHelp.
type DialogManager interface {
	ScreenSize() (int, int)
	OpenDialog(tview.Primitive)
	CloseDialog()
}

// AddHelp adds help text capabilities to a box.  Whenever F1 is pressed while
// the box has focus, a help dialog will be displayed with the specified title
// and text.  The text should be pre-formatted to 56-column width, and may
// include tview color codes.  The help dialog is dismissed when the user
// presses Enter, ESC, or Tab.
func AddHelp(p *tview.Box, title, help string, manager DialogManager) {
	capture := p.GetInputCapture()
	p.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		if evt.Key() == tcell.KeyF1 {
			ShowHelp(p, title, help, false, manager)
			return nil
		}
		if capture != nil {
			evt = capture(evt)
		}
		return evt
	})
}

// ShowHelp displays a help dialog with the specified title and text.  The text
// should be pre-formatted to a 56-column width, and may include tview color
// codes.  The help dialog is dismissed when the user presses Enter, ESC, or
// Tab.  If fullScreen is true, the help dialog will be given the full screen.
// Otherwise, the help dialog will be placed above or below the specified box,
// whichever gives it sufficient space.
func ShowHelp(p *tview.Box, title, help string, fullScreen bool, manager DialogManager) {
	var x, y, width, height int
	// The desired height of the dialog is the number of lines of help text,
	// plus 4 for border and padding.
	height = strings.Count(help, "\n") + 4
	if !strings.HasSuffix(help, "\n") {
		height++
	}
	// Calculate the number of rows available above and below the box.
	screenWidth, screenHeight := manager.ScreenSize()
	if fullScreen {
		y = (screenHeight - height) / 2
		if y < 2 {
			y = 2
		}
		if height > screenHeight-y-2 {
			height = screenHeight - y - 2
		}
	} else {
		_, py, _, pheight := p.GetRect()
		rowsAbove := py
		rowsBelow := screenHeight - py - pheight
		// Determine vertical location of dialog.
		if height <= rowsBelow {
			y = py + pheight
		} else if rowsBelow >= rowsAbove {
			y = py + pheight
			height = rowsBelow
		} else {
			y = 0
			height = rowsAbove
		}
	}
	// Determine horizontal location of dialog.
	if screenWidth < 60 {
		x, width = 0, screenWidth
	} else {
		x, width = (screenWidth-60)/2, 60
	}
	// Create the dialog.
	dialog := tview.NewTextView().
		SetText(help).
		SetDoneFunc(func(_ tcell.Key) { manager.CloseDialog() }).
		SetDynamicColors(true).
		SetSize(height-4, width-4).
		SetTextColor(tview.Styles.PrimaryTextColor).
		SetWrap(false)
	dialog.SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle(" " + title + " ").
		SetMouseCapture(func(act tview.MouseAction, evt *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if dialog.InRect(evt.Position()) && act == tview.MouseLeftClick {
				manager.CloseDialog()
			}
			return act, evt
		})
	// Place it as desired.
	modal := tview.NewGrid().
		SetColumns(x, width, 0).
		SetRows(y, height, 0).
		AddItem(dialog, 1, 1, 1, 1, 0, 0, true)
	// Show the dialog.
	manager.OpenDialog(modal)
}
