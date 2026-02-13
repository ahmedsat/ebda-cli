package kobo

import (
	"fmt"
	"strings"
)

type Query map[string]fmt.Stringer

func (q Query) String() string {
	sb := strings.Builder{}
	sb.WriteRune('{')
	for k, v := range q {
		fmt.Fprintf(&sb, "%q", k)
		sb.WriteRune(':')
		if v == nil {
			sb.WriteString("null")
		} else {
			fmt.Fprintf(&sb, "%q", v)
		}
		sb.WriteRune(',')
	}
	// remove last comma
	temp := sb.String()[:len(sb.String())-1]
	sb.Reset()
	sb.WriteString(temp)

	sb.WriteRune('}')
	return sb.String()
}

func StringValue(s string) fmt.Stringer {
	sb := strings.Builder{}
	sb.WriteString(s)
	return &sb
}
