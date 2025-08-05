//go:build packetpdf

package notrep

import (
	_ "embed" // .
)

//go:embed Notable_Report_v20250804.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
