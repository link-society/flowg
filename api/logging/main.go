package logging

import "log/slog"

// Logger returns the logger of the API operations, derived from the root
// logger injected into each operation.
//
// It tags entries with the "api" channel so operation logs can be told apart
// from the rest of FlowG's output.
func Logger(rootLogger *slog.Logger) *slog.Logger {
	return rootLogger.With(slog.String("channel", "api"))
}
