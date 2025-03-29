package timex

import (
	"fmt"
	"time"
)

var unitMap = map[string]time.Duration{
	"ns": 1,
	"us": 1e3,
	"µs": 1e3,
	"ms": 1e6,
	"s":  1e9,
	"m":  60 * 1e9,
	"h":  60 * 60 * 1e9,
	"d":  24 * 60 * 60 * 1e9,
	"w":  7 * 24 * 60 * 60 * 1e9,
	"y":  365 * 24 * 60 * 60 * 1e9,
}

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}

	// parse number + unit
	var qty int64
	var unit string

	if _, err := fmt.Sscanf(s, "%d%s", &qty, &unit); err != nil {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	if unit == "" {
		return 0, fmt.Errorf("missing unit in duration: %s", s)
	}

	factor, ok := unitMap[unit]
	if !ok {
		return 0, fmt.Errorf("invalid unit in duration: %s", unit)
	}

	d := time.Duration(qty) * factor

	return d, nil
}
