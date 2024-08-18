package logging

import "log/slog"

func Setup(verbose bool) {
	var level slog.Level
	if verbose {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	slog.SetDefault(slog.New(newHandler(opts)))
}
