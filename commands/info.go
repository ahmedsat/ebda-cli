package commands

import (
	"fmt"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
)

type Info struct {
}

// Description implements [main.subcommand].
func (i *Info) Description() string {
	return "General purpose information"
}

// Name implements [main.subcommand].
func (i *Info) Name() string {
	return "info"
}

// Run implements [main.subcommand].
func (i *Info) Run(args []string) (err error) {

	if len(args) < 1 {
		return fmt.Errorf("no code provided")
	}

	if strings.HasPrefix(args[0], "EG/") {
		code := args[0]
		args = args[1:]
		var farms []types.Farm
		farms, err = frappe.Get[types.Farm](frappe.Filters{frappe.NewFilter("farm_id", frappe.Eq, code)}, nil, nil)
		if err != nil {
			return err
		}
		switch len(farms) {
		case 0:
			fmt.Println("No farm found")
		case 1:
			farm := farms[0]
			fmt.Printf("Farm: %s\n", farm.ArabicName)
			fmt.Printf("Owner: %s\n", farm.FarmOwner)
			fmt.Printf("Region: %s\n", farm.Region)
			fmt.Printf("Total farmers: %d\n", farm.TotalFarmers)
			fmt.Printf("Area: %.2f\n", farm.Area)
			fmt.Printf("Creation date: %s\n", farm.CreationDate)
		default:
			fmt.Println("Multiple farms found")
		}
		return
	}

	if len(args) != 0 {
		i.Run(args)
	}

	return
}

// Usage implements [main.subcommand].
func (i *Info) Usage() string {
	panic("unimplemented")
}
