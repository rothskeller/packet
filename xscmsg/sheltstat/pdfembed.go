//go:build packetpdf

package sheltstat

import (
	_ "embed" // .
)

//go:embed XSC_SheltStat_v20190619_p12.pdf
var pdfBaseEmbed []byte

func init() {
	Type23.PDFBase = pdfBaseEmbed
}
