// extract-pifo-fields generates «ident».def.go for each «ident».html in the
// source tree.  It also generates «pkgname».def.go in each package containing
// such files.
package main

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	log.SetFlags(0)
	var packages = make(map[string][]string)
	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".html") {
			log.Print(path)
			varname := extract(path)
			packages[filepath.Dir(path)] = append(packages[filepath.Dir(path)], varname)
		}
		return nil
	})
	for pkg, files := range packages {
		emitPackage(pkg, files)
	}
}
