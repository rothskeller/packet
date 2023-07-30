//go:build packetpdf

package incident

import (
	_ "embed" // .
)

//go:embed ICS-309.pdf
var ics309pdf []byte
