package commands

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
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

	sb := strings.Builder{}
	sb.WriteString("Region\tFarms\tFarmers\tArea\n")

	region := map[string]struct {
		farms   int
		farmers int
		area    float64
	}{}

	var (
		totalFarms   int
		totalFarmers int
		totalArea    float64
	)

	// get all farms
	farms, err := frappe.Get[types.Farm](
		frappe.Filters{
			frappe.NewFilter("type", frappe.Eq, "farm"),
			frappe.NewFilter("farm_status", frappe.Neq, "Cancelled"),
			frappe.NewFilter("creation_date", frappe.Gte, form.Format("2006-01-02")),
			frappe.NewFilter("creation_date", frappe.Lte, to.AddDate(0, 0, 1).Format("2006-01-02")),
		},
		[]string{"name", "region", "total_farmers", "farm_area__feddan"}, nil)
	if err != nil {
		return
	}

	for _, farm := range farms {
		d := region[farm.Region]
		d.farms++
		totalFarms++
		d.farmers += farm.TotalFarmers
		totalFarmers += farm.TotalFarmers
		d.area += farm.Area
		totalArea += farm.Area
		region[farm.Region] = d
	}

	for region, date := range region {
		fmt.Fprintf(&sb, "%s\t%d\t%d\t%.2f\n", region, date.farms, date.farmers, date.area)
	}

	fmt.Fprintf(&sb, "Total\t%d\t%d\t%.2f\n", totalFarms, totalFarmers, totalArea)

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
