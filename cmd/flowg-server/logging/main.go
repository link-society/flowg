package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"

	applog "link-society.com/flowg/internal/app/logging"
)

// Setup builds the root slog logger, writing to stdout using the
// correlation-id aware handler from the internal/app/logging package, and makes
// it the process-wide default. The returned logger is the one injected into the
// application.
//
// When verbose is true, the log level is forced to debug. Otherwise the level is
// derived from levelName ("debug", "info", "warn" or "error", defaulting to
// "info").
func Setup(verbose bool, levelName string) *slog.Logger {
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

	opts := &slog.HandlerOptions{Level: level}
	logger := slog.New(applog.NewHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return logger
}

// Discard configures the process-wide default slog logger to drop every record.
// It is primarily used to silence logging output during tests.
func Discard() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	slog.SetDefault(slog.New(applog.NewHandler(io.Discard, opts)))
}
