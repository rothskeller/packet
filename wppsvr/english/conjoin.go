// Package english contains minor utilities for generating English-language
// text.
package english

import (
	"fmt"
	"strings"
)

// Conjoin returns a set of strings joined together with the specified
// conjunction and (of course) the Oxford comma.
func Conjoin(ss []string, conj string) string {
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return ss[0]
	case 2:
		return fmt.Sprintf("%s %s %s", ss[0], conj, ss[1])
	default:
		return fmt.Sprintf("%s %s, %s", strings.Join(ss[:len(ss)-1], ", "), conj, ss[len(ss)-1])
	}
}
