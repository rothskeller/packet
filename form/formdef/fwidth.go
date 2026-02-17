package formdef

import "github.com/rothskeller/pdf/v2"

// CharWidth returns the approximate number of characters that can be emitted
// in a field based on its PDF rendering.  (It's approximate because PDF
// rendering usually uses a proportional font.)  The function returns zero if
// a width cannot be determined.
func (fd *FieldDef) CharWidth() int {
	var text TextRenderer

	// Find the first PDF text rendering for the field.
	for _, r := range fd.PDF {
		if tr, ok := r.Renderer.(TextRenderer); ok {
			text = tr
			break
		}
	}
	if text.Page == 0 { // no text renderer found
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
