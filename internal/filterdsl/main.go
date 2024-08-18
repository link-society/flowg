package filterdsl

import (
	"encoding/json"

	"link-society.com/flowg/internal/storage"
)

func astToFilter(ast map[string]interface{}) storage.Filter {
	if val, exists := ast["$and"]; exists {
		v := val.([]interface{})
		filters := make([]storage.Filter, len(v))
		for i, filter := range v {
			filters[i] = astToFilter(filter.(map[string]interface{}))
		}
		return &storage.AndFilter{Filters: filters}
	}

	if val, exists := ast["$or"]; exists {
		v := val.([]interface{})
		filters := make([]storage.Filter, len(v))
		for i, filter := range v {
			filters[i] = astToFilter(filter.(map[string]interface{}))
		}
		return &storage.OrFilter{Filters: filters}
	}

	if val, exists := ast["$not"]; exists {
		return &storage.NotFilter{Filter: astToFilter(val.(map[string]interface{}))}
	}

	if val, exists := ast["$eq"]; exists {
		v := val.(map[string]interface{})
		field := v["field"].(string)
		value := v["value"].(string)
		return &storage.FieldExact{Field: field, Value: value}
	}

	if val, exists := ast["$in"]; exists {
		v := val.(map[string]interface{})
		field := v["field"].(string)
		iValues := v["values"].([]interface{})
		values := make([]string, len(iValues))
		for i, value := range iValues {
			values[i] = value.(string)
		}
		return &storage.FieldIn{Field: field, Values: values}
	}

	panic("unreachable")
}

func Compile(input string) (storage.Filter, error) {
	output, err := compile(input)
	if err != nil {
		return nil, err
	}

	var ast map[string]interface{}
	if err := json.Unmarshal([]byte(output), &ast); err != nil {
		return nil, &UnmarshalError{Reason: err}
	}

	return astToFilter(ast), nil
}
