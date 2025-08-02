//go:build veocipdf

package shelter

import (
	_ "embed" // .
)

//go:embed Shelter_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
