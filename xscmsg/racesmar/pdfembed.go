//go:build packetpdf

package racesmar

import (
	_ "embed" // .
)

//go:embed XSC_RACES_MA_Req_v20220129.pdf
var pdfBaseEmbed24 []byte

//go:embed XSC_RACES_MA_Req_v20240711_V1_Test76_p12.pdf
var pdfBaseEmbed33 []byte

func init() {
	Type24.PDFBase = pdfBaseEmbed24
	Type33.PDFBase = pdfBaseEmbed33
}
