//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"github.com/rothskeller/packet/packetver"
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
	sh.Run("scp", "SCCoPIFO.forms", "sccares:www/www/form-bundles/4.0.9/SCCoPIFO.forms")
}

func Dist() {
	mg.Deps(IncidentHTML, WindowsResources)
	os.MkdirAll("dist", 0777)
	sh.RunWith(map[string]string{"GOOS": "darwin", "GOARCH": "amd64"},
		mg.GoCmd(), "build", "-tags", "sccopifo", "-o", fmt.Sprintf("./dist/packet-v%s_darwin_amd64", packetver.Version), "./cmd/packet")
	sh.RunWith(map[string]string{"GOOS": "linux", "GOARCH": "amd64"},
		mg.GoCmd(), "build", "-tags", "sccopifo", "-o", fmt.Sprintf("./dist/packet-v%s_linux_amd64", packetver.Version), "./cmd/packet")
	sh.RunWith(map[string]string{"GOOS": "windows", "GOARCH": "amd64"},
		mg.GoCmd(), "build", "-tags", "sccopifo", "-o", fmt.Sprintf("./dist/packet-v%s_windows_amd64_console.exe", packetver.Version), "./cmd/packet")
	sh.RunWith(map[string]string{"GOOS": "windows", "GOARCH": "amd64"},
		mg.GoCmd(), "build", "-ldflags", "-H=windowsgui", "-tags", "sccopifo", "-o", fmt.Sprintf("./dist/packet-v%s_windows_amd64_gui.exe", packetver.Version), "./cmd/packet")
}

func IncidentHTML() {
	if newer, err := target.Dir("cmd/packet/server/incident.html", "cmd/packet/server/incident-html"); err != nil || newer {
		println("Updating incident.html.")
		sh.Run("bash", "-c", "go run -tags sccopifo ./cmd/merge-html cmd/packet/server/incident-html/* >cmd/packet/server/incident.html")
	}
}

func WindowsResources() {
	if newer, err := target.Path("cmd/packet/rsrc_windows_amd64.syso", "cmd/packet/winres/winres.json"); err != nil || newer {
		println("Building Windows resources.")
		sh.Run("go-winres", "make", "--in", "cmd/packet/winres.json", "--arch", "amd64", "--out", "cmd/packet/rsrc")
	}
}

func UpdateForms() {
	if stat, err := os.Stat("/Users/stever/.local/share/packet/4.0.9/SCCoPIFO"); err != nil {
		return
	} else if newer, err := target.DirNewer(stat.ModTime(), "forms/SCCoPIFO"); err != nil || newer {
		println("Removing old forms.")
		os.RemoveAll("/Users/stever/.local/share/packet/4.0.9/SCCoPIFO")
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

func A315() error {
	mg.Deps(IncidentHTML, WindowsResources)
	if _, err := os.Stat("/Volumes/str-a315"); err != nil {
		println("Mounting STR-A315.")
		if err = sh.Run("osascript", "-e", `mount volume "smb://str-a315/c"`); err != nil {
			return err
		}
	}
	if err := sh.RunWith(map[string]string{"GOOS": "windows"}, mg.GoCmd(), "build", "-tags", "sccopifo", "-o", "/Volumes/str-a315/PackItForms/packet.exe", "./cmd/packet"); err != nil {
		return err
	}
	if err := sh.RunWith(map[string]string{"GOOS": "windows"}, mg.GoCmd(), "build", "-ldflags", "-H=windowsgui", "-tags", "sccopifo", "-o", "/Volumes/str-a315/PackItForms/pifo.exe", "./cmd/packet"); err != nil {
		return err
	}
	return nil
}
