//go:build packetpdf

package cpodsite

import (
	_ "embed" // .
)

//go:embed CPOD_Site_Information_v20250812.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
