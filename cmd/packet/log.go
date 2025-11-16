package main

import (
	"context"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rothskeller/packet/cmd/packet/osdep"
)

var logFH *os.File

// reopenLogFile opens a log file named with today's date.  It is called on
// command invocation and again just after midnight each day.
func reopenLogFile() {
	var (
		logpath string
		newFH   *os.File
		logger  *slog.Logger
		nextLog time.Time
		err     error
	)
	nextLog = time.Now()
	logpath = filepath.Join(osdep.LogsDir, nextLog.Format("2006-01-02")+".log")
	if err = os.MkdirAll(osdep.LogsDir, 0777); err != nil {
		newFH = os.Stderr
	} else if newFH, err = os.OpenFile(logpath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		newFH = os.Stderr
	} else if logFH != nil {
		slog.Info("switching to new log file", "file", logpath)
	}
	logger = slog.New(handler{slog.NewTextHandler(newFH, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 1 && groups[0] == slog.SourceKey {
				if a.Key == "function" {
					a.Key = "fn"
				} else {
					return slog.Attr{}
				}
			}
			if groups == nil && a.Key == slog.TimeKey {
				a.Key = "t"
				if a.Value.Kind() == slog.KindTime {
					str := a.Value.Time().Format("15:04:05")
					a.Value = slog.StringValue(str)
				}
			}
			if groups == nil && a.Key == slog.LevelKey {
				a.Key = "l"
			}
			if groups == nil && a.Key == slog.MessageKey {
				a.Key = "m"
			}
			return a
		},
	})})
	slog.SetDefault(logger)
	if newFH == os.Stderr {
		slog.Error("can't open log file (logging to stderr)", "file", logpath, "err", err)
	}
	if logFH != nil {
		logFH.Close()
	}
	logFH = newFH
	nextLog = time.Now().AddDate(0, 0, 1)
	nextLog = time.Date(nextLog.Year(), nextLog.Month(), nextLog.Day(), 0, 0, 0, 0, time.Local)
	go func() {
		time.Sleep(time.Until(nextLog))
		reopenLogFile()
	}()
}

// handler is our own variant of the standard TextHandler.
type handler struct {
	*slog.TextHandler
}

// Handle defers to TextHandler, but first, it adds the reporting function name
// to all records at levels other than Info.
func (h handler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level != slog.LevelInfo && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.Function != "" {
			r.AddAttrs(slog.String("fn", path.Base(f.Function)))
		}
	}
	return h.TextHandler.Handle(ctx, r)
}
