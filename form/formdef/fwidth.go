package formdef

import "github.com/rothskeller/pdf/v2"

// CharWidth returns the approximate number of characters that can be emitted
// in a field based on its PDF rendering.  (It's approximate because PDF
// rendering usually uses a proportional font.)  The function returns zero if
// a width cannot be determined.
func (fd *FieldDef) CharWidth(n int) int {
	var text TextRenderer

	// If the nth PDF renderer for the field is not a text renderer, skip.
	if n < 0 || n >= len(fd.PDF) {
		return 0
	}
	if tr, ok := fd.PDF[n].Renderer.(TextRenderer); ok {
		text = tr
	} else {
		return 0
	}
	// What is the width of a "0" in the font used in that renderer?
	zero, _, _ := pdf.MeasureText("0", text.Font, text.FontSize)
	if zero == 0 { // no font metrics
		return 0
	}
	// How many of those fit in the width of the PDF render area?
	return int((text.Rectangle.URX - text.Rectangle.LLX) / zero)
}
