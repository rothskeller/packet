package msgedit

import (
	"github.com/gdamore/tcell/v2"
)

var (
	// StyleCheckboxCheck is the style for the mark between the square
	// brackets in a checkbox when the checkbox is checked.
	StyleCheckboxCheck tcell.Style
	// StyleCheckboxFrame is the style for the square brackets around the
	// checkbox value in the message editing panel.
	StyleCheckboxFrame tcell.Style
	// StyleCheckboxHint is the style for the text describing how to change
	// the state of a checkbox.
	StyleCheckboxHint tcell.Style
	// StyleFieldName is the style for a field name in the message viewing
	// or editing panels.
	StyleFieldName tcell.Style
	// StyleFieldNameInvalid is the style for a field name in the message
	// editing panel when the field contains an invalid value.
	StyleFieldNameInvalid tcell.Style
	// StyleFieldNameInvalidSelected is the style for a field name in the
	// message editing panel when the field is selected and contains an
	// invalid value.
	StyleFieldNameInvalidSelected tcell.Style
	// StyleFieldNameSelected is the style for a field name in the message
	// editing panel when that field is selected.
	StyleFieldNameSelected tcell.Style
	// StyleInput is the style for an input field when it does not have the
	// focus.
	StyleInput tcell.Style
	// StyleInputFocus is the style for an input field when it does have the
	// focus.
	StyleInputFocus tcell.Style
	// StyleInputHint is the style for the comment or hint on an input
	// field.
	StyleInputHint tcell.Style
	// StyleInputOption is the style for a possible value (that is not the
	// current value) for an input field.
	StyleInputOption tcell.Style
	// StyleInputOptionArrow is the style for the arrow shown between an
	// input field and its list of possible values.
	StyleInputOptionArrow tcell.Style
	// StyleInputOptionSelected is the style for the current value of an
	// input field when it appears in an option list.
	StyleInputOptionSelected tcell.Style
	// StyleInputSelected is the style for selected text in an input field.
	StyleInputSelected tcell.Style
	// StylePanelHeader is the style for the title bar of a panel.
	StylePanelHeader tcell.Style
	// colorSelected is the background color for list selections.
	colorSelected tcell.Color
)

// InitStyles initializes the style variables based on the capabilities of the
// screen.
func InitStyles(screen tcell.Screen) {
	// On terminals with limited palettes, tcell automatically maps these
	// RGB values to palette colors (and I've checked that they are the
	// desired palette colors).  I'm not coding them as the palette colors
	// directly, because full-color terminals tend to substitute different
	// RGB values for the basic palette colors (see Wikipedia ANSI Escape
	// Code for details), and I don't want that.  This is done in an init
	// function so that, if future needs dictate, I can select different
	// colors based on the terminal's available palette.
	colorSelected = tcell.NewHexColor(0x000080)
	StyleCheckboxCheck = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleCheckboxFrame = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleCheckboxHint = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleFieldName = tcell.StyleDefault.Foreground(tcell.NewHexColor(0x00FFFF))
	StyleFieldNameInvalid = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFF0000))
	StyleFieldNameInvalidSelected = MakeSelected(StyleFieldNameInvalid)
	StyleFieldNameSelected = MakeSelected(StyleFieldName)
	StyleInput = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleInputFocus = MakeSelected(StyleInput)
	StyleInputHint = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleInputOption = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleInputOptionArrow = tcell.StyleDefault.Foreground(tcell.NewHexColor(0xFFFFFF))
	StyleInputOptionSelected = tcell.StyleDefault.Foreground(tcell.NewHexColor(0)).Background(tcell.NewHexColor(0xC0C0C0))
	StyleInputSelected = tcell.StyleDefault.Foreground(tcell.NewHexColor(0)).Background(tcell.NewHexColor(0xFFFFFF))
	StylePanelHeader = tcell.StyleDefault.Foreground(tcell.NewHexColor(0x00FF00))
}

// MakeSelected changes a style to reflect its being part of a selected row in
// a list.
func MakeSelected(s tcell.Style) tcell.Style {
	return s.Background(colorSelected)
}
