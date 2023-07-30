//go:build packetpdf

package racesmar

import (
	_ "embed" // .
)

//go:embed XSC_RACES_MA_Req_Fillable_v20220129_p12.pdf
var pdfBaseEmbed []byte

func init() {
	pdfBase = pdfBaseEmbed
}
