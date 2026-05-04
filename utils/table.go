package utils

type Table struct {
	Headers []string
	Rows    [][]string
}

func NewTable(headers ...string) Table {
	return Table{
		Headers: headers,
		Rows:    [][]string{},
	}
}

func NewTableWithRows(headers []string, rows [][]string) Table {
	return Table{
		Headers: headers,
		Rows:    rows,
	}
}

func (t Table) AppendRow(row ...string) {
	t.Rows = append(t.Rows, row)
}
