//go:build packetpdf

package ics213

import (
	_ "embed" // .
)

//go:embed ICS-213_SCCo_Message_Form_v20220119.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
