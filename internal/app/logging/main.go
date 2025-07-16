package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

var VERBOSE_LOGGING = false

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

	VERBOSE_LOGGING = verbose
	opts := &slog.HandlerOptions{Level: level}
	slog.SetDefault(slog.New(newHandler(os.Stdout, opts)))
}

func Discard() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	slog.SetDefault(slog.New(newHandler(io.Discard, opts)))
}
