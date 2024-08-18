package logging

import "log/slog"

func Setup() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	slog.SetDefault(slog.New(newHandler(opts)))
}
