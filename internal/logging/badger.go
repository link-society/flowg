package logging

import (
	"fmt"
	"log/slog"
)

type BadgerLogger struct {
	Channel string
}

func (l *BadgerLogger) Errorf(format string, v ...interface{}) {
	slog.Error(fmt.Sprintf(format, v...), "channel", l.Channel)
}

func (l *BadgerLogger) Warningf(format string, v ...interface{}) {
	slog.Warn(fmt.Sprintf(format, v...), "channel", l.Channel)
}

func (l *BadgerLogger) Infof(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...), "channel", l.Channel)
}

func (l *BadgerLogger) Debugf(format string, v ...interface{}) {
	slog.Debug(fmt.Sprintf(format, v...), "channel", l.Channel)
}
