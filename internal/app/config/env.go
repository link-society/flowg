package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// applyEnv overrides the configuration with the environment variables that are
// set (an empty variable is treated as unset), as declared by the `env` struct
// tags of the schema (see [applyEnvToStruct]).
func (cfg *RootConfig) applyEnv() {
	// The storage backend is selected first, so that the settings of the
	// selected backend are then overridden by their own environment variables.
	backend := os.Getenv("FLOWG_STORAGE_BACKEND")
	if backend != "" && backend != cfg.Storage.Backend.Backend {
		cfg.Storage.Backend.selectBackend(backend)
	}

	applyEnvToStruct(reflect.ValueOf(cfg).Elem())
}

// applyEnvToStruct overrides the fields of s from the environment variables
// named by their `env` tags:
//
//   - `string`, `bool` and `[]string` (comma-separated items, trimmed) fields
//     are set to the parsed value of the variable, when set;
//   - pointer-to-struct fields (blocks) with an `env` tag are enabled or
//     disabled by that boolean variable, when set: an enabled block is
//     allocated with its defaults if missing, a disabled block is removed;
//   - struct fields and non-nil blocks are overridden recursively.
//
// An unparsable boolean is read as false. It panics on a field type that is not
// supported, which is a programming error.
func applyEnvToStruct(s reflect.Value) {
	t := s.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		value := s.Field(i)
		key, hasKey := field.Tag.Lookup("env")

		switch {
		case field.Type.Kind() == reflect.Struct:
			applyEnvToStruct(value)

		case field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct:
			if hasKey {
				if raw := os.Getenv(key); raw != "" {
					if !parseEnvBool(raw) {
						value.SetZero()
					} else if value.IsNil() {
						value.Set(reflect.New(field.Type.Elem()))
						applyDefaultsToStruct(value.Elem())
					}
				}
			}

			if !value.IsNil() {
				applyEnvToStruct(value.Elem())
			}

		case !hasKey:
			continue

		default:
			raw := os.Getenv(key)
			if raw == "" {
				continue
			}

			if err := setFromString(value, raw); err != nil {
				if value.Kind() != reflect.Bool {
					panic(fmt.Sprintf("invalid env for %s.%s: %v", t.Name(), field.Name, err))
				}

				value.SetBool(false)
			}
		}
	}
}

// parseEnvBool parses a boolean environment variable, an unparsable value being
// read as false.
func parseEnvBool(raw string) bool {
	value, err := strconv.ParseBool(raw)
	return err == nil && value
}
