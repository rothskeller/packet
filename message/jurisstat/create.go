package jurisstat

import (
	"time"

	"github.com/rothskeller/packet/message/common"
)

// New creates a new OA jurisdiction status form with default values.
func New() *JurisStat {
	return &JurisStat{
		StdFields: common.StdFields{
			FormVersion: "2.2",
			MessageDate: time.Now().Format("01/02/2006"),
			Handling:    "IMMEDIATE",
			ToLocation:  "County EOC",
		},
	}
}
