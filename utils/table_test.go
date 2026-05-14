package utils_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestTableAppendAndRender(t *testing.T) {
	tbl := utils.NewTable("Name", "City")
	tbl.AppendRow("Alice", "Cairo")

	if len(tbl.Headers) != 2 || len(tbl.Rows) != 1 {
		t.Fatalf("table shape = headers %d rows %d, want headers 2 rows 1", len(tbl.Headers), len(tbl.Rows))
	}
	if csv := tbl.CSV(); !strings.Contains(csv, "Name,") || !strings.Contains(csv, "Alice,") {
		t.Fatalf("CSV missing expected content: %q", csv)
	}
	if tsv := tbl.TSV(); !strings.Contains(tsv, "Name\t") || !strings.Contains(tsv, "Alice\t") {
		t.Fatalf("TSV missing expected content: %q", tsv)
	}
}

func TestTableAppendRowConcurrent(t *testing.T) {
	tbl := utils.NewTable("val")
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			tbl.AppendRow(n)
		}(i)
	}
	wg.Wait()
	if len(tbl.Rows) != 100 {
		t.Fatalf("Rows len = %d, want 100", len(tbl.Rows))
	}
}
