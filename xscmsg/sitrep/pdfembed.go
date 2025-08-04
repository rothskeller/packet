//go:build packetpdf

package sitrep

import (
	_ "embed" // .
)

//go:embed Situation_Report_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
