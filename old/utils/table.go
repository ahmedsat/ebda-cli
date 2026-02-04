package utils

import (
	"strings"
)

type Table interface {
	GetHeader() []string
	GetRows() [][]string
	SetHeader(header ...string)
	AppendRow(row ...string)
	AppendRows(rows ...[]string)
}

type TablePrinter func(Table) string

type TableBase struct {
	Header []string
	Rows   [][]string
}

func (t *TableBase) GetHeader() []string {
	return t.Header
}

func (t *TableBase) GetRows() [][]string {
	return t.Rows
}

func (t *TableBase) SetHeader(header ...string) {
	t.Header = header
}

func (t *TableBase) AppendRow(row ...string) {
	t.Rows = append(t.Rows, row)
}

func (t *TableBase) AppendRows(rows ...[]string) {
	t.Rows = append(t.Rows, rows...)
}

func TablePrinterCsv(t Table) string {

	sb := strings.Builder{}

	UpdateF(
		t.GetHeader(),
		TruePredicate[string](),
		func(s string) string { return "\"" + s + "\"" },
	)

	sb.WriteString(strings.Join(t.GetHeader(), ",") + "\n")
	for _, row := range t.GetRows() {
		UpdateF(
			row,
			TruePredicate[string](),
			func(s string) string { return "\"" + s + "\"" },
		)
		sb.WriteString(strings.Join(row, ",") + "\n")
	}

	// remove last newline
	str := sb.String()
	sb.Reset()
	sb.WriteString(str[:len(str)-1])

	return sb.String()
}
func TablePrinterTsv(t Table) string {
	sb := strings.Builder{}
	UpdateF(
		t.GetHeader(),
		TruePredicate[string](),
		func(s string) string { return "\"" + s + "\"" },
	)
	sb.WriteString(strings.Join(t.GetHeader(), "\t") + "\n")
	for _, row := range t.GetRows() {
		UpdateF(
			row,
			TruePredicate[string](),
			func(s string) string { return "\"" + s + "\"" },
		)
		sb.WriteString(strings.Join(row, "\t") + "\n")
	}
	return sb.String()

}
