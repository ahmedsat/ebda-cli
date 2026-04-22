package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/atotto/clipboard"
)

type Soil struct {
}

// Description implements [main.subcommand].
func (s *Soil) Description() string {
	return "Get soil records"
}

// Name implements [main.subcommand].
func (s *Soil) Name() string {
	return "soil"
}

// Result implements [main.subcommand].

// Run implements [main.subcommand].
func (s *Soil) Run(args []string) (err error) {

	fs := flag.NewFlagSet("soil", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "Farm\tLocation\n")

	records, err := frappe.Get[types.SoilAnalysis](nil, frappe.List{"farm", "location"}, nil)
	if err != nil {
		return err
	}
	for _, r := range records {
		fmt.Fprintf(&sb, "%s\t%s\n", r.Farm, r.Location)
	}

	if *copy {
		clipboard.WriteAll(sb.String())
		fmt.Println("copied to clipboard")
		return
	}

	fmt.Print(sb.String())

	return
}

// Usage implements [main.subcommand].
func (s *Soil) Usage() string {
	panic("unimplemented")
}
