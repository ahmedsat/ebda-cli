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
	copy    bool
	records []types.SoilAnalysis
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
func (s *Soil) Result() any {
	sb := strings.Builder{}

	fmt.Fprintf(&sb, "Farm\tLocation\n")
	for _, r := range s.records {
		fmt.Fprintf(&sb, "%s\t%s\n", r.Farm, r.Location)
	}

	if s.copy {
		clipboard.WriteAll(sb.String())
		return "copied to clipboard"
	}

	return sb.String()
}

// Run implements [main.subcommand].
func (s *Soil) Run(args []string) (err error) {
	fs := flag.NewFlagSet("soil", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)
	s.copy = *copy

	s.records, err = frappe.Get[types.SoilAnalysis](nil, frappe.List{"farm", "location"})
	if err != nil {
		return err
	}

	return
}

// Usage implements [main.subcommand].
func (s *Soil) Usage() string {
	panic("unimplemented")
}
