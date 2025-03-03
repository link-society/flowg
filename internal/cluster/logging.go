package cluster

import (
	"context"
	"log"
	"log/slog"
	"strings"
)

type memberlistLoggerWriter struct {
	parent *slog.Logger
}

func (w *memberlistLoggerWriter) Write(p []byte) (n int, err error) {
	msg := string(p)
	level := slog.LevelInfo

	switch {
	case strings.HasPrefix(msg, "[DEBUG]"):
		level = slog.LevelDebug
		msg = strings.TrimPrefix(msg, "[DEBUG]")

	case strings.HasPrefix(msg, "[INFO]"):
		level = slog.LevelInfo
		msg = strings.TrimPrefix(msg, "[INFO]")

	case strings.HasPrefix(msg, "[WARN]"):
		level = slog.LevelWarn
		msg = strings.TrimPrefix(msg, "[WARN]")

	case strings.HasPrefix(msg, "[ERR]"):
		level = slog.LevelError
		msg = strings.TrimPrefix(msg, "[ERR]")
	}

	w.parent.Log(context.Background(), level, msg)

	return len(p), nil
}

func newMemberlistLogger(parent *slog.Logger) *log.Logger {
	return log.New(&memberlistLoggerWriter{parent: parent}, "", 0)
}
