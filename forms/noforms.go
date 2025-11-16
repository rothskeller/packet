//go:build !allforms && !losaltos && !milpitas && !report911 && !sccopifo && !sccopvt

package forms

import "embed"

var EmbeddedForms embed.FS
