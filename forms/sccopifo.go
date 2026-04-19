//go:build sccopifo

package forms

import "embed"

// Unfortunately we have to list all of the embeds separately, because the
// SCCoPIFO repository also contains Word documents, and we don't want those.

//go:embed SCCoPIFO/*/*.form
//go:embed SCCoPIFO/*/*.html
//go:embed SCCoPIFO/*/*.pdf
//go:embed SCCoPIFO/definitions.html
//go:embed SCCoPIFO/pack-it-forms.*
//go:embed SCCoPIFO/README.txt
//go:embed SCCoPIFO/SCCoPIFO.*
//go:embed SCCoPIFO/update.json
var EmbeddedForms embed.FS
