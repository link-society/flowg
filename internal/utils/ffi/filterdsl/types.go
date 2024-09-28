package filterdsl

import "link-society.com/flowg/internal/models"

type Filter interface {
	Evaluate(record *models.LogRecord) bool
}

type FilterAnd struct {
	Filters []Filter
}

type FilterOr struct {
	Filters []Filter
}

type FilterNot struct {
	Filter Filter
}

type FilterMatchField struct {
	Field string
	Value string
}

type FilterMatchFieldList struct {
	Field  string
	Values []string
}

func (f *FilterAnd) Evaluate(record *models.LogRecord) bool {
	for _, filter := range f.Filters {
		if !filter.Evaluate(record) {
			return false
		}
	}

	return true
}

func (f *FilterOr) Evaluate(record *models.LogRecord) bool {
	for _, filter := range f.Filters {
		if filter.Evaluate(record) {
			return true
		}
	}

	return false
}

func (f *FilterNot) Evaluate(record *models.LogRecord) bool {
	return !f.Filter.Evaluate(record)
}

func (f *FilterMatchField) Evaluate(record *models.LogRecord) bool {
	value, exists := record.Fields[f.Field]
	if !exists {
		return false
	}

	return value == f.Value
}

func (f *FilterMatchFieldList) Evaluate(record *models.LogRecord) bool {
	value, exists := record.Fields[f.Field]
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
