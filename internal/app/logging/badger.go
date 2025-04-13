package logging

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

type BadgerLogger struct {
	Channel string
}

var _ badger.Logger = (*BadgerLogger)(nil)

func (l *BadgerLogger) Errorf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Error(line, "channel", l.Channel)
	}
}

func (l *BadgerLogger) Warningf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Warn(line, "channel", l.Channel)
	}
}

func (l *BadgerLogger) Infof(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Info(line, "channel", l.Channel)
	}
}

func (l *BadgerLogger) Debugf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Debug(line, "channel", l.Channel)
	}
}
