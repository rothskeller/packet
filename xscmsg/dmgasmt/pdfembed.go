//go:build packetpdf

package dmgasmt

import (
	_ "embed" // .
)

//go:embed Damage_Assessment_v20250812.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
