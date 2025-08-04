//go:build packetpdf

package wssurvey

import (
	_ "embed" // .
)

//go:embed Windshield_Survey_v20250730.pdf
var pdfBaseEmbed []byte

func init() {
	Type.PDFBase = pdfBaseEmbed
}
