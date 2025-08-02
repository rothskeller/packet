//go:build veocipdf

package notrep

import (
	_ "embed" // .
)

//go:embed Notable_Report_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
