package commands

import (
	"fmt"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
)

type Info struct {
	sb strings.Builder
}

// Description implements [main.subcommand].
func (i *Info) Description() string {
	return "General purpose information"
}

// Name implements [main.subcommand].
func (i *Info) Name() string {
	return "info"
}

// Result implements [main.subcommand].
func (i *Info) Result() any {
	return i.sb.String()
}

// Run implements [main.subcommand].
func (i *Info) Run(args []string) (err error) {

	if len(args) < 1 {
		return
	}

	if strings.HasPrefix(args[0], "EG/") {
		code := args[0]
		args = args[1:]
		var farms []types.Farm
		farms, err = frappe.Get[types.Farm](frappe.Filters{frappe.NewFilter("farm_id", frappe.Eq, code)}, nil)
		if err != nil {
			return err
		}
		switch len(farms) {
		case 0:
			i.sb.WriteString("No farm found\n")
		case 1:
			farm := farms[0]
			fmt.Fprintf(&i.sb, "Farm: %s\n", farm.ArabicName)
			fmt.Fprintf(&i.sb, "Owner: %s\n", farm.FarmOwner)
			fmt.Fprintf(&i.sb, "Region: %s\n", farm.Region)
			fmt.Fprintf(&i.sb, "Total farmers: %d\n", farm.TotalFarmers)
			fmt.Fprintf(&i.sb, "Area: %.2f\n", farm.Area)
			fmt.Fprintf(&i.sb, "Creation date: %s\n", farm.CreationDate)
		default:
			i.sb.WriteString("Multiple farms found\n")
		}
		return
	}

	if len(args) != 0 {
		i.Run(args)
	}
	return nil
}

// Usage implements [main.subcommand].
func (i *Info) Usage() string {
	panic("unimplemented")
}
