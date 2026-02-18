package cmdutil

import (
	"fmt"
	"os/exec"
	"runtime"
)

// ShowPDF opens the specified PDF file in the system default PDF viewer.
func ShowPDF(filename string) (err error) {
	var showcmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		showcmd = exec.Command("cmd.exe", "/C", filename)
	case "darwin":
		showcmd = exec.Command("open", filename)
	default:
		showcmd = exec.Command("xdg-open", filename)
	}
	if err := showcmd.Start(); err != nil {
		return fmt.Errorf("starting PDF viewer: %s", err)
	}
	go func() { showcmd.Wait() }()
	return nil
}
