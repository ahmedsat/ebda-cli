package utils

import (
	"fmt"
	"strings"
)

// filters=[["field1", "=", "value1"], ["field2", ">", "value2"]]

type FilterComparator string

const (
	Eq  FilterComparator = "="
	Gt  FilterComparator = ">"
	Lt  FilterComparator = "<"
	Gte FilterComparator = ">="
	Lte FilterComparator = "<="
	Neq FilterComparator = "!="
)

type Filter struct {
	Field    string
	Operator FilterComparator
	Value    string
}

func NewFilter(field string, operator FilterComparator, value string) Filter {
	return Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

func (f Filter) String() string {
	return fmt.Sprintf("[\"%s\",\"%s\",\"%s\"]", f.Field, f.Operator, f.Value)
}

type Filters []Filter

func (f Filters) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for _, filter := range f {
		sb.WriteString(filter.String())
		sb.WriteString(",")
	}
	// remove last comma
	temp := sb.String()[:len(sb.String())-1]
	sb.Reset()
	sb.WriteString(temp)

	sb.WriteString("]")
	return sb.String()
}

type List []string

func (l List) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for _, item := range l {
		sb.WriteString("\"" + item + "\"")
		sb.WriteString(",")
	}
	temp := sb.String()[:len(sb.String())-1]
	sb.Reset()
	sb.WriteString(temp)
	sb.WriteString("]")
	return sb.String()
}
