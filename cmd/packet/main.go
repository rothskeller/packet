package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cmd"
	"github.com/rothskeller/packet/cmd/packet/osdep"
)

func main() {
	reopenLogFile()
	logInvocation()
	redirectPIFOtoGUI()
	cmd.Execute()
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

func redirectPIFOtoGUI() {
	// If the user invoked pifo.exe without arguments, treat that as if they
	// invoked with "gui".
	if exe := strings.ToLower(filepath.Base(os.Args[0])); exe == "pifo" || exe == "pifo.exe" {
		if len(os.Args) == 1 {
			os.Args = append(os.Args, "gui")
		}
	}
}
