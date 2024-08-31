package logstorage

type Filter interface {
	Evaluate(entry *LogEntry) bool
}

type AndFilter struct {
	Filters []Filter
}

type OrFilter struct {
	Filters []Filter
}

type NotFilter struct {
	Filter Filter
}

type FieldExact struct {
	Field string
	Value string
}

type FieldIn struct {
	Field  string
	Values []string
}

func (f *AndFilter) Evaluate(entry *LogEntry) bool {
	for _, filter := range f.Filters {
		if !filter.Evaluate(entry) {
			return false
		}
	}

	return true
}

func (f *OrFilter) Evaluate(entry *LogEntry) bool {
	for _, filter := range f.Filters {
		if filter.Evaluate(entry) {
			return true
		}
	}

	return false
}

func (f *NotFilter) Evaluate(entry *LogEntry) bool {
	return !f.Filter.Evaluate(entry)
}

func (f *FieldExact) Evaluate(entry *LogEntry) bool {
	value, exists := entry.Fields[f.Field]
	if !exists {
		return false
	}

	return value == f.Value
}

func (f *FieldIn) Evaluate(entry *LogEntry) bool {
	value, exists := entry.Fields[f.Field]
	if !exists {
		return false
	}

	for _, v := range f.Values {
		if v == value {
			return true
		}
	}

	return false
}
