package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"

	applog "link-society.com/flowg/internal/app/logging"
)

// Setup configures the process-wide default slog logger to write to stdout
// using the correlation-id aware handler from the internal/app/logging package.
//
// When verbose is true, the log level is forced to debug and the access-log
// middleware is told to dump response bodies for failed requests. Otherwise the
// level is derived from levelName ("debug", "info", "warn" or "error",
// defaulting to "info").
func Setup(verbose bool, levelName string) {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		switch strings.ToLower(levelName) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	}

	applog.VERBOSE_LOGGING = verbose
	opts := &slog.HandlerOptions{Level: level}
	slog.SetDefault(slog.New(applog.NewHandler(os.Stdout, opts)))
}

// Discard configures the process-wide default slog logger to drop every record.
// It is primarily used to silence logging output during tests.
func Discard() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	slog.SetDefault(slog.New(applog.NewHandler(io.Discard, opts)))
}
