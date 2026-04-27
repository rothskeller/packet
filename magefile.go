//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

// Default target.
var Default = GUI

func Forms() {
	os.Remove("SCCoPIFO.zip")
	os.Chdir("forms/SCCoPIFO")
	sh.Run("zip", "-r", "../../SCCoPIFO.zip", ".", "-x", "*.docx", "*.md")
	os.Chdir("../..")
	sh.Run(mg.GoCmd(), "run", "./cmd/sign-forms", "SCCoPIFO", "SCCoPIFO.zip")
	os.Remove("SCCoPIFO.zip")
	sh.Run("scp", "SCCoPIFO.forms", "sccares:www/www/form-bundles/4.0/SCCoPIFO.forms")
}

func IncidentHTML() {
	if newer, err := target.Dir("cmd/packet/server/incident.html", "cmd/packet/server/incident-html"); err != nil || newer {
		println("Updating incident.html.")
		sh.Run("bash", "-c", "go run -tags sccopifo ./cmd/merge-html cmd/packet/server/incident-html/* >cmd/packet/server/incident.html")
	}
}

func UpdateForms() {
	if stat, err := os.Stat("/Users/stever/.local/share/packet/4.0/SCCoPIFO"); err != nil {
		return
	} else if newer, err := target.DirNewer(stat.ModTime(), "forms/SCCoPIFO"); err != nil || newer {
		println("Removing old forms.")
		os.RemoveAll("/Users/stever/.local/share/packet/4.0/SCCoPIFO")
	}
}

func StopOldServer() error {
	println("Stopping old server.")
	return sh.Run(mg.GoCmd(), "run", "-tags", "sccopifo", "./cmd/packet", "server", "stop")
}

func GUI() error {
	mg.Deps(IncidentHTML, UpdateForms, StopOldServer)
	return sh.Run(mg.GoCmd(), "run", "-tags", "sccopifo", "./cmd/packet", "gui", "2025-11-MPMP")
}

func Build() error {
	mg.Deps(IncidentHTML, UpdateForms, StopOldServer)
	return sh.Run(mg.GoCmd(), "build", "-tags", "sccopifo", "./cmd/packet")
}
