package cluster

import (
	"context"
	"fmt"

	"io"

	"log"
	"log/slog"

	"github.com/hashicorp/go-hclog"
)

type raftLogger struct {
	logger *slog.Logger
	name   string
	args   []any
}

func newRaftLogger(name string, backend *slog.Logger, args ...any) *raftLogger {
	return &raftLogger{
		logger: backend.With(args...),
		name:   name,
		args:   args,
	}
}

func (l *raftLogger) Log(level hclog.Level, msg string, args ...any) {
	args = append([]any{slog.String("domain", l.name)}, args...)

	switch level {
	case hclog.Trace:
		fallthrough
	case hclog.Debug:
		l.logger.Debug(msg, args...)

	case hclog.NoLevel:
		fallthrough
	case hclog.Info:
		l.logger.Info(msg, args...)

	case hclog.Warn:
		l.logger.Warn(msg, args...)

	case hclog.Error:
		l.logger.Error(msg, args...)
	}
}

func (l *raftLogger) Trace(msg string, args ...any) {
	l.Log(hclog.Trace, msg, args...)
}

func (l *raftLogger) Debug(msg string, args ...any) {
	l.Log(hclog.Debug, msg, args...)
}

func (l *raftLogger) Info(msg string, args ...any) {
	l.Log(hclog.Info, msg, args...)
}

func (l *raftLogger) Warn(msg string, args ...any) {
	l.Log(hclog.Warn, msg, args...)
}

func (l *raftLogger) Error(msg string, args ...any) {
	l.Log(hclog.Error, msg, args...)
}

func (l *raftLogger) IsTrace() bool {
	return l.logger.Enabled(context.Background(), slog.LevelDebug)
}

func (l *raftLogger) IsDebug() bool {
	return l.logger.Enabled(context.Background(), slog.LevelDebug)
}

func (l *raftLogger) IsInfo() bool {
	return l.logger.Enabled(context.Background(), slog.LevelInfo)
}

func (l *raftLogger) IsWarn() bool {
	return l.logger.Enabled(context.Background(), slog.LevelWarn)
}

func (l *raftLogger) IsError() bool {
	return l.logger.Enabled(context.Background(), slog.LevelError)
}

func (l *raftLogger) ImpliedArgs() []any {
	return l.args
}

func (l *raftLogger) With(args ...any) hclog.Logger {
	return &raftLogger{
		logger: l.logger.With(args...),
		name:   l.name,
		args:   append(l.args, args...),
	}
}

func (l *raftLogger) Name() string {
	return l.name
}

func (l *raftLogger) Named(name string) hclog.Logger {
	return &raftLogger{
		logger: l.logger,
		name:   fmt.Sprintf("%s.%s", l.name, name),
		args:   l.args,
	}
}

func (l *raftLogger) ResetNamed(name string) hclog.Logger {
	return &raftLogger{
		logger: l.logger,
		name:   name,
		args:   nil,
	}
}

func (l *raftLogger) SetLevel(level hclog.Level) {
	// noop
}

func (l *raftLogger) GetLevel() hclog.Level {
	return hclog.NoLevel
}

func (l *raftLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	var lvl slog.Level
	switch opts.ForceLevel {
	case hclog.Trace:
		fallthrough
	case hclog.Debug:
		lvl = slog.LevelDebug

	case hclog.NoLevel:
		fallthrough
	case hclog.Info:
		lvl = slog.LevelInfo

	case hclog.Warn:
		lvl = slog.LevelWarn

	case hclog.Error:
		lvl = slog.LevelError

	default:
		lvl = slog.LevelInfo
	}

	return slog.NewLogLogger(l.logger.Handler(), lvl)
}

func (l *raftLogger) StandardWriter(*hclog.StandardLoggerOptions) io.Writer {
	return l
}

func (l *raftLogger) Write(p []byte) (int, error) {
	l.logger.Info(string(p))
	return len(p), nil
}
