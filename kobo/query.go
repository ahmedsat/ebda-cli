package kobo

import "encoding/json"

type Query map[string]any

func (q Query) String() string {
	b, _ := json.Marshal(q)
	return string(b)
}
