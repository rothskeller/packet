// get-pifo-html extracts the SCCoPIFO HTML files for each version of each form
// from the pack-it-forms repository.
//
// usage: get-pifo-html [pack-it-forms-dir] [output-dir]
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// formToPackage maps HTML filename to our package name for the form.  Note that
// in some cases multiple filenames map to the same package name.  Filenames
// that map to an empty string are not processed.  Filenames that are not in the
// map cause an error.
var formToPackage = map[string]string{
	"form-allied-health-facility-status.html":      "ahfacstat",
	"form-checkin-out.html":                        "",
	"form-ics213.html":                             "ics213",
	"form-los-altos-da-v2.html":                    "",
	"form-los-altos-da.html":                       "",
	"form-los-altos-emergency.html":                "",
	"form-los-altos-emergency.receiver.html":       "",
	"form-los-altos-urgent-report.html":            "",
	"form-milpitas-da.html":                        "",
	"form-milpitas-incident-report.html":           "",
	"form-oa-muni-status.html":                     "jurisstat",
	"form-oa-mutual-aid-request-v2.html":           "racesmar",
	"form-oa-mutual-aid-request-v2.read-only.html": "",
	"form-oa-mutual-aid-request.html":              "racesmar",
	"form-oa-shelter-status.html":                  "sheltstat",
	"form-scco-eoc-213rr-v2.html":                  "",
	"form-scco-eoc-213rr.html":                     "eoc213rr",
}

func main() {
	var pifodir = "../pack-it-forms"
	var outdir = "xscmsg"
	var err error

	log.SetFlags(0)
	if len(os.Args) > 1 {
		pifodir = os.Args[1]
	}
	if len(os.Args) > 2 {
		outdir = os.Args[2]
	}
	if outdir, err = filepath.Abs(outdir); err != nil {
		log.Fatal(err)
	}
	if err = os.Chdir(pifodir); err != nil {
		log.Fatal(err)
	}
	tagsExist := getExistingTags(pifodir)
	tagsRead := getTagsRead(outdir)
	for _, tag := range tagsExist {
		if !tagsRead[tag] {
			readHTMLFilesFromTag(pifodir, outdir, tag)
		}
	}
	saveTagsRead(outdir, tagsExist)
	restoreTopOfBranch()
}

func getExistingTags(dir string) []string {
	if err := os.Chdir(dir); err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("git", "tag", "--merged", "SCCo.2", "--contains", "vSCCo.22")
	list, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(strings.TrimSpace(string(list)), "\n")
}

func getTagsRead(dir string) map[string]bool {
	var seen = make(map[string]bool)
	list, err := os.ReadFile(filepath.Join(dir, "tags-read"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		log.Fatal(err)
	}
	for _, tag := range strings.Split(string(list), "\n") {
		seen[tag] = true
	}
	return seen
}

func readHTMLFilesFromTag(pifodir, outdir, tag string) {
	// Check out the tag.
	if err := exec.Command("git", "checkout", tag).Run(); err != nil {
		log.Fatal(err)
	}
	log.Printf("Tag %s:", tag)
	// Look for forms to process.
	forms, _ := filepath.Glob("form-*.html")
	for _, form := range forms {
		pkg, ok := formToPackage[form]
		if !ok {
			log.Fatalf("Tag %q contains an unknown form %s.", tag, form)
		}
		if pkg != "" {
			processForm(outdir, tag, form, pkg)
		}
	}
}

func processForm(outdir, tag, form, pkg string) {
	html, version := readHTML(form)
	if version == "" {
		log.Fatalf("no version found in %s %s", tag, form)
	}
	fname := fmt.Sprintf("%s.v%s.html", form[:len(form)-5], version)
	err := os.WriteFile(filepath.Join(outdir, pkg, fname), []byte(html), 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("  %s", form)
}

func saveTagsRead(outdir string, tags []string) {
	err := os.WriteFile(
		filepath.Join(outdir, "tags-read"),
		[]byte(strings.Join(tags, "\n")+"\n"),
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func restoreTopOfBranch() {
	if err := exec.Command("git", "checkout", "SCCo.2").Run(); err != nil {
		log.Fatal(err)
	}
}
