package logstorage

import (
	"fmt"
	"log/slog"
)

type serverLogger struct{}

func (l *serverLogger) Errorf(format string, v ...interface{}) {
	slog.Error(fmt.Sprintf(format, v...), "channel", "badger")
}

func (l *serverLogger) Warningf(format string, v ...interface{}) {
	slog.Warn(fmt.Sprintf(format, v...), "channel", "badger")
}

func (l *serverLogger) Infof(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...), "channel", "badger")
}

func (l *serverLogger) Debugf(format string, v ...interface{}) {
	slog.Debug(fmt.Sprintf(format, v...), "channel", "badger")
}
