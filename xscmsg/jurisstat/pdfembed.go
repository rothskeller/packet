//go:build packetpdf

package jurisstat

import (
	_ "embed" // .
)

//go:embed XSC_JurisStat_Fillable_v20190528_p123.pdf
var pdfBaseEmbed []byte

func init() {
	Type22.PDFBase = pdfBaseEmbed
}
