//go:build packetpdf

package ics213

import (
	_ "embed" // .
)

//go:embed ICS-213_SCCo_Message_Form_Fillable_v20220119_p1.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
