//go:build packetpdf

package eoc213rr

import (
	_ "embed" // .
)

//go:embed XSC_EOC-213RR_v20170803_with_XSC_RACES_Routing_Slip_v20190527.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
