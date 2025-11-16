//go:build sccopifo && !sccopvt

package forms

import "embed"

//go:embed SCCoPIFO
var EmbeddedForms embed.FS
