//go:build packetpdf

package racesmar

import (
	_ "embed" // .
)

//go:embed XSC_RACES_MA_Req_v20220129.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
