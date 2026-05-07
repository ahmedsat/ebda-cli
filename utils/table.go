package utils

import (
	"fmt"
	"strings"
	"sync"
)

type Table struct {
	Headers []any
	Rows    [][]any
	sync.Mutex
}

func NewTable(headers ...any) Table {
	return Table{
		Headers: headers,
		Rows:    [][]any{},
	}
}

func NewTableWithRows(headers []any, rows [][]any) Table {
	return Table{
		Headers: headers,
		Rows:    rows,
	}
}

func (t *Table) AppendRow(row ...any) {
	t.Lock()
	defer t.Unlock()
	t.Rows = append(t.Rows, row)
}

func (t *Table) TSV() string {
	sb := strings.Builder{}
	for _, header := range t.Headers {
		fmt.Fprintf(&sb, "%v\t", header)
	}
	sb.WriteString("\n")
	for _, row := range t.Rows {
		for _, col := range row {
			fmt.Fprintf(&sb, "%v\t", col)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *Table) CSV() string {
	sb := strings.Builder{}
	for _, header := range t.Headers {
		fmt.Fprintf(&sb, "%v,", header)
	}
	sb.WriteString("\n")
	for _, row := range t.Rows {
		for _, col := range row {
			fmt.Fprintf(&sb, "%v,", col)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *Table) WriteToGoogleSheet(spreadSheetID string, rng string) (err error) {
	panic("not implemented")
}
