//go:build packetpdf

package racesmar

import (
	_ "embed" // .
)

//go:embed XSC_RACES_MA_Req_v20240711_V1_Test76_p12.pdf
var pdfBaseEmbed []byte

func init() {
	Type33.PDFBase = pdfBaseEmbed
}
