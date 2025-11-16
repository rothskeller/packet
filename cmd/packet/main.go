package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"

	"github.com/rothskeller/packet/cmd/packet/cmd"
	"github.com/rothskeller/packet/cmd/packet/osdep"
)

func main() {
	reopenLogFile()
	logInvocation()
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
