//go:build veocipdf

package resreq

import (
	_ "embed" // .
)

//go:embed Resource_Request_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
