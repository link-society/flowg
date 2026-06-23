package utils

import (
	"fmt"
	"os"
	"slices"

	"github.com/go-logfmt/logfmt"

	"link-society.com/flowg/internal/models"
)

// Printer renders log records to standard output as logfmt lines.
type Printer struct {
	encoder *logfmt.Encoder
}

// NewPrinter builds a Printer that writes to standard output.
func NewPrinter() *Printer {
	return &Printer{
		encoder: logfmt.NewEncoder(os.Stdout),
	}
}

// Print encodes one log record, emitting its fields in a stable (sorted) order
// so successive records line up.
func (p *Printer) Print(log models.LogRecord) error {
	if err := p.encoder.EncodeKeyval("@timestamp", log.Timestamp); err != nil {
		return fmt.Errorf("could not encode timestamp: %w", err)
	}

	var fieldNames []string
	for fieldName := range log.Fields {
		fieldNames = append(fieldNames, fieldName)
	}
	slices.Sort(fieldNames)

	for _, fieldName := range fieldNames {
		fieldValue := log.Fields[fieldName]
		if err := p.encoder.EncodeKeyval(fieldName, fieldValue); err != nil {
			return fmt.Errorf("could not encode field '%s': %w", fieldName, err)
		}
	}

	if err := p.encoder.EndRecord(); err != nil {
		return fmt.Errorf("could not flush record: %w", err)
	}

	return nil
}
