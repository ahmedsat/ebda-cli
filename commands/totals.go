package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/ahmedsat/ebda-cli/services"
	"github.com/ahmedsat/ebda-cli/utils"
	"github.com/atotto/clipboard"
)

type Totals struct{}

const cliTimeFormat = "2-1-2006"

// Run implements [main.subcommand].
func (t *Totals) Run(args []string) (err error) {
	fs := flag.NewFlagSet("totals", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	upload := fs.Bool("upload", false, "Upload to google sheet")
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

	tb := utils.NewTable("Region", "Farms", "Farmers", "Area")
	for _, row := range report.Rows {
		tb.AppendRow(row.Region, row.Farms, row.Farmers, row.Area)
	}
	tb.AppendRow("Total", report.TotalFarms, report.TotalFarmers, report.TotalArea)

	if *upload {
		err = tb.WriteToGoogleSheet(
			"11tXfIz9o_cgD-czMTQRRLF9JEkmvNTF5QSmdY6lVQFs",
			"sheet1!A1",
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	if *copy {
		clipboard.WriteAll(tb.TSV())
		fmt.Println("copied to clipboard")
	}

	if *upload || *copy {
		return
	}

	fmt.Print(tb.TSV())

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
