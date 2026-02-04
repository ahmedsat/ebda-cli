package commands

import (
	"fmt"
	"strings"
)

type ListFlagString []string

func (l ListFlagString) String() string {
	if len(l) == 0 {
		return "[]"
	}

	sb := strings.Builder{}
	sb.WriteString("[")
	for i, item := range l {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q", item)) // quoted
	}
	sb.WriteString("]")
	return sb.String()
}

func (l *ListFlagString) Set(s string) error {
	*l = strings.Split(s, ",")
	return nil
}
