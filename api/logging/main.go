package logging

import "log/slog"

// Logger returns the logger shared by every API operation.
//
// It tags entries with the "api" channel so operation logs can be told apart
// from the rest of FlowG's output. It resolves [slog.Default] lazily so it
// reflects the logging configuration in force when an operation runs.
func Logger() *slog.Logger {
	return slog.Default().With(slog.String("channel", "api"))
}
