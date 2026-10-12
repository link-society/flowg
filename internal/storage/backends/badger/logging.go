package badger

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dgraph-io/badger/v4"
)

// badgerLogger forwards BadgerDB's log output to a slog logger, one record per
// line. Info messages are logged at debug level.
type badgerLogger struct {
	logger *slog.Logger
}

var _ badger.Logger = (*badgerLogger)(nil)

func (l *badgerLogger) Errorf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		l.logger.Error(line)
	}
}

func (l *badgerLogger) Warningf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		l.logger.Warn(line)
	}
}

func (l *badgerLogger) Infof(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func (l *badgerLogger) Debugf(format string, v ...interface{}) {
	content := fmt.Sprintf(format, v...)
	lines := strings.FieldsFunc(content, func(c rune) bool { return c == '\n' || c == '\r' })
	for _, line := range lines {
		l.logger.Debug(line)
	}
}
