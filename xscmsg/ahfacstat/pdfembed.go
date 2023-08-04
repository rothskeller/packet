//go:build packetpdf

package ahfacstat

import (
	_ "embed" // .
)

//go:embed Allied_Health_Facility_Status_DEOC-9_v20180200_with_XSC_RACES_Routing_Slip_Fillable_v20190527.pdf
var pdfBaseEmbed []byte

func init() {
	pdfBase = pdfBaseEmbed
}
