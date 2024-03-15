package util

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// DebugfBytes logs a slice of bytes and its length.
func DebugfBytes(logger *slog.Logger, msg string, b []byte) {
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, DebugfBytes]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, msg, pcs[0])
	r.AddAttrs(slog.Attr{
		Key:   "length",
		Value: slog.IntValue(len(b)),
	})
	r.AddAttrs(slog.Attr{
		Key:   "data",
		Value: slog.StringValue(fmt.Sprintf("%v", b)),
	})
	_ = logger.Handler().Handle(context.Background(), r)
}

func NewLogger() *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		ReplaceAttr: replace,
		Level:       slog.LevelDebug,
	}))
}
