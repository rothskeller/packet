//go:build packetpdf

package roadcl

import (
	_ "embed" // .
)

//go:embed Road_Closure_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
