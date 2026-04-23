package commands

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/services"
	"github.com/atotto/clipboard"
)

type Totals struct{}

const cliTimeFormat = "2-1-2006"

// Run implements [main.subcommand].
func (t *Totals) Run(args []string) (err error) {
	fs := flag.NewFlagSet("totals", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	formStr := fs.String("from", "1-1-2022", "Date part of ISO")
	toSt := fs.String("to", time.Now().Format(cliTimeFormat), "Date part of ISO")
	fs.Parse(args)

	// parse date
	form, err := time.Parse(cliTimeFormat, *formStr)
	if err != nil {
		return
	}
	to, err := time.Parse(cliTimeFormat, *toSt)
	if err != nil {
		return
	}

	report, err := services.LoadTotalsReport(form, to)
	if err != nil {
		return
	}

	sb := strings.Builder{}
	sb.WriteString("Region\tFarms\tFarmers\tArea\n")
	for _, row := range report.Rows {
		fmt.Fprintf(&sb, "%s\t%d\t%d\t%.2f\n", row.Region, row.Farms, row.Farmers, row.Area)
	}
	fmt.Fprintf(&sb, "Total\t%d\t%d\t%.2f\n", report.TotalFarms, report.TotalFarmers, report.TotalArea)

	if *copy {
		clipboard.WriteAll(sb.String())
		fmt.Println("copied to clipboard")
		return
	}
	fmt.Print(sb.String())
	return
}

// Description implements [main.subcommand].
func (t *Totals) Description() (desc string) {
	return "Get totals region -> total farmers -> total area"
}

// Name implements [main.subcommand].
func (t *Totals) Name() (name string) {
	return "totals"
}

// Usage implements [main.subcommand].
func (t *Totals) Usage() (usage string) {
	panic("unimplemented")
}
