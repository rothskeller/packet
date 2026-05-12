//go:build allforms

package forms

import "embed"

//go:embed SCCoPIFO/*/*.form
//go:embed SCCoPIFO/*/*.html
//go:embed SCCoPIFO/*/*.pdf
//go:embed SCCoPIFO/definitions.html
//go:embed SCCoPIFO/README.txt
//go:embed SCCoPIFO/SCCoPIFO.*
//go:embed SCCoPIFO/update.json
var EmbeddedForms embed.FS
