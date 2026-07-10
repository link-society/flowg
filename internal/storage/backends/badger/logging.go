package badger

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

type badgerLogger struct {
	Channel string
}

var _ badger.Logger = (*badgerLogger)(nil)

func (l *badgerLogger) Errorf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Error(line, "channel", l.Channel)
	}
}

func (l *badgerLogger) Warningf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Warn(line, "channel", l.Channel)
	}
}

func (l *badgerLogger) Infof(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func (l *badgerLogger) Debugf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		slog.Debug(line, "channel", l.Channel)
	}
}
