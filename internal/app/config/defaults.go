package config

import (
	"reflect"

	"fmt"
	"strconv"
	"strings"
)

// newWithDefaults allocates a T and fills it with [applyDefaults].
func newWithDefaults[T any]() *T {
	v := new(T)
	applyDefaults(v)
	return v
}

// applyDefaults fills the struct pointed to by v from the `default` tags of its
// fields:
//
//   - `string`, `bool` and `[]string` (comma-separated items, trimmed) fields
//     are set to the parsed tag value;
//   - pointer-to-struct fields tagged `default:"{}"` are allocated and filled
//     recursively, untagged ones are left nil (optional blocks);
//   - struct fields are filled recursively.
//
// Fields without a `default` tag keep their zero value. It panics on an invalid
// tag, which is a programming error.
func applyDefaults(v any) {
	applyDefaultsToStruct(reflect.ValueOf(v).Elem())
}

func applyDefaultsToStruct(s reflect.Value) {
	t := s.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		value := s.Field(i)
		tag, hasTag := field.Tag.Lookup("default")

		switch {
		case field.Type.Kind() == reflect.Struct:
			applyDefaultsToStruct(value)

		case !hasTag:
			continue

		case field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct:
			if tag != "{}" {
				panic(fmt.Sprintf(
					"invalid default for %s.%s: blocks only accept `default:\"{}\"`",
					t.Name(), field.Name,
				))
			}

			value.Set(reflect.New(field.Type.Elem()))
			applyDefaultsToStruct(value.Elem())

		default:
			if err := setFromString(value, tag); err != nil {
				panic(fmt.Sprintf("invalid default for %s.%s: %v", t.Name(), field.Name, err))
			}
		}
	}
}

// setFromString parses raw into v, which must be a `string`, a `bool` or a
// `[]string` (comma-separated items, trimmed; empty when raw is empty).
func setFromString(v reflect.Value, raw string) error {
	switch {
	case v.Kind() == reflect.String:
		v.SetString(raw)

	case v.Kind() == reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		v.SetBool(b)

	case v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String:
		items := []string{}
		if raw != "" {
			for item := range strings.SplitSeq(raw, ",") {
				items = append(items, strings.TrimSpace(item))
			}
		}
		v.Set(reflect.ValueOf(items))

	default:
		return fmt.Errorf("unsupported type %s", v.Type())
	}

	return nil
}
