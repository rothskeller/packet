//go:build veocipdf

package cpodsite

import (
	_ "embed" // .
)

//go:embed CPOD_Site_Information_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
