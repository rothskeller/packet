package racesmar

import (
	"time"

	"github.com/rothskeller/packet/message/common"
)

// New creates a new RACES mutual aid request form with default values.
func New() *RACESMAR {
	return &RACESMAR{
		StdFields: common.StdFields{
			FormVersion: "2.3",
			MessageDate: time.Now().Format("01/02/2006"),
			Handling:    "ROUTINE",
			ToLocation:  "County EOC",
		},
	}
}
