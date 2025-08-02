//go:build veocipdf

package dmgasmt

import (
	_ "embed" // .
)

//go:embed Damage_Assessment_v20250731.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
