//go:build packetpdf

package shelter

import (
	_ "embed" // .
)

//go:embed Shelter_v20250812.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
