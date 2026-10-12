package config

import (
	"os"

	"strconv"
	"strings"
)

// getEnvListString reads a comma-separated environment variable into a slice,
// trimming whitespace around each item, or returns defaultValue when unset.
func getEnvListString(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	items := strings.Split(value, ",")
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
	}

	return items
}

// getEnvString returns the value of the environment variable key, or
// defaultValue when it is unset or empty.
func getEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvBool parses the environment variable key as a boolean, falling back to
// defaultValue when it is unset and to false when it cannot be parsed.
func getEnvBool(key string, defaultValue bool) bool {
	stringVal := os.Getenv(key)
	if stringVal == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(stringVal)
	if err == nil {
		return value
	}
	return false
}
