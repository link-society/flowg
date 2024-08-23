package storage

import "fmt"

func intersectKeysMap(a, b map[string]struct{}) map[string]struct{} {
	result := map[string]struct{}{}

	for key := range a {
		if _, exists := b[key]; exists {
			result[key] = struct{}{}
		}
	}

	return result
}

func unionKeysMap(a, b map[string]struct{}) map[string]struct{} {
	result := map[string]struct{}{}

	for key := range a {
		result[key] = struct{}{}
	}

	for key := range b {
		result[key] = struct{}{}
	}

	return result
}

func differenceKeysMap(a, b map[string]struct{}) map[string]struct{} {
	result := map[string]struct{}{}

	for key := range a {
		if _, exists := b[key]; !exists {
			result[key] = struct{}{}
		}
	}

	return result
}

func sliceToMap(slice []string) map[string]struct{} {
	result := map[string]struct{}{}

	for _, key := range slice {
		result[key] = struct{}{}
	}

	return result
}

func mapToSlice(m map[string]struct{}) []string {
	result := []string{}

	for key := range m {
		result = append(result, key)
	}

	return result
}

func flatten(
	prefix string,
	input interface{},
	output map[string]string,
) {
	switch v := input.(type) {
	case map[string]interface{}:
		for key, item := range v {
			fullKey := key
			if prefix != "" {
				fullKey = fmt.Sprintf("%s.%s", prefix, key)
			}
			flatten(fullKey, item, output)
		}

	case []interface{}:
		for i, item := range v {
			fullKey := fmt.Sprintf("%s[%d]", prefix, i)
			flatten(fullKey, item, output)
		}

	default:
		output[prefix] = fmt.Sprintf("%v", v)
	}
}
