//go:build packetpdf

package cpodupd

import (
	_ "embed" // .
)

//go:embed CPOD_Commodities_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
