package logging

import (
	"io"
	"log/slog"
	"os"
)

var VERBOSE_LOGGING = false

func Setup(verbose bool) {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	VERBOSE_LOGGING = verbose
	opts := &slog.HandlerOptions{Level: level}
	slog.SetDefault(slog.New(newHandler(os.Stdout, opts)))
}

func Discard() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	slog.SetDefault(slog.New(newHandler(io.Discard, opts)))
}
