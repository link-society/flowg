package filterdsl

import "encoding/json"

func astToFilter(ast map[string]interface{}) Filter {
	if val, exists := ast["$and"]; exists {
		v := val.([]interface{})
		filters := make([]Filter, len(v))
		for i, filter := range v {
			filters[i] = astToFilter(filter.(map[string]interface{}))
		}
		return &FilterAnd{Filters: filters}
	}

	if val, exists := ast["$or"]; exists {
		v := val.([]interface{})
		filters := make([]Filter, len(v))
		for i, filter := range v {
			filters[i] = astToFilter(filter.(map[string]interface{}))
		}
		return &FilterOr{Filters: filters}
	}

	if val, exists := ast["$not"]; exists {
		return &FilterNot{Filter: astToFilter(val.(map[string]interface{}))}
	}

	if val, exists := ast["$eq"]; exists {
		v := val.(map[string]interface{})
		field := v["field"].(string)
		value := v["value"].(string)
		return &FilterMatchField{Field: field, Value: value}
	}

	if val, exists := ast["$in"]; exists {
		v := val.(map[string]interface{})
		field := v["field"].(string)
		iValues := v["values"].([]interface{})
		values := make([]string, len(iValues))
		for i, value := range iValues {
			values[i] = value.(string)
		}
		return &FilterMatchFieldList{Field: field, Values: values}
	}

	panic("unreachable")
}

func Compile(input string) (Filter, error) {
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
