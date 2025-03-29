package sse

import (
	"context"
	"fmt"

	"bufio"
	"bytes"
	"io"
	"strings"
)

type EventStreamReader struct {
	scanner *bufio.Scanner
}

func NewEventStreamReader(r io.Reader) *EventStreamReader {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 64*1024)
	scanner.Split(func(data []byte, atEOF bool) (int, []byte, error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i, nlen := findDoubleNewLine(data); i >= 0 {
			return i + nlen, data[:i], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	})

	return &EventStreamReader{
		scanner: scanner,
	}
}

func (r *EventStreamReader) Next() (Event, error) {
	var e Event

	if r.scanner.Scan() {
		payload := r.scanner.Bytes()

		lines := bytes.Split(payload, []byte("\n"))
		for _, line := range lines {
			line = bytes.TrimSpace(line)
			if len(line) == 0 {
				continue
			}

			parts := bytes.SplitN(line, []byte(":"), 2)
			if len(parts) != 2 {
				return e, fmt.Errorf("invalid line: %q", line)
			}

			key := strings.TrimSpace(string(parts[0]))
			value := strings.TrimSpace(string(parts[1]))

			switch key {
			case "id":
				e.ID = value

			case "event":
				e.Type = value

			case "data":
				e.Data = value
			}
		}

		if e.Type == "" {
			return e, fmt.Errorf("missing event type")
		}

		if e.Data == "" {
			return e, fmt.Errorf("missing event data")
		}

		return e, nil
	}

	if err := r.scanner.Err(); err != nil {
		if err == context.Canceled {
			return e, io.EOF
		}

		return e, err
	}

	return e, io.EOF
}
