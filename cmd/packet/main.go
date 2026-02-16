package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cio"
	"github.com/rothskeller/packet/cmd/packet/cmd"
	"github.com/rothskeller/packet/cmd/packet/osdep"
)

func main() {
	reopenLogFile()
	logInvocation()
	defaultCommands()
	if err := cmd.Execute(); err != nil {
		cio.Open().Error("%s", err)
	}
}

func logInvocation() {
	var attrs []slog.Attr

	attrs = append(attrs, slog.String("cmd", strings.Join(os.Args, " ")))
	if exe, _ := os.Executable(); exe != "" {
		attrs = append(attrs, slog.String("exe", exe))
	}
	ver := "(devel)"
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		ver = info.Main.Version
	}
	cmd.RootCmd.Version = ver
	attrs = append(attrs, slog.String("ver", ver))
	if osdep.IsAdmin() {
		attrs = append(attrs, slog.Bool("admin", true))
	}
	slog.LogAttrs(context.Background(), slog.LevelDebug, "START", attrs...)
}

func defaultCommands() {
	if len(os.Args) != 1 {
		return // They specified a command.
	}
	// Assign a default command based on the executable name.
	// If the user invoked pifo.exe without arguments, treat that as if they
	// invoked with "gui".
	if exe := strings.ToLower(filepath.Base(os.Args[0])); exe == "pifo" || exe == "pifo.exe" {
		os.Args = append(os.Args, "gui")
	} else if strings.Contains(exe, "install") || strings.Contains(exe, "setup") {
		os.Args = append(os.Args, "install")
	}
}
