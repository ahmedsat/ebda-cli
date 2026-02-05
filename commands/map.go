package commands

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/atotto/clipboard"
)

type Map struct {
	copy bool
	maps []types.MapRecord
}

// Name implements [main.subcommand].
func (m *Map) Name() string {
	panic("unimplemented")
}

// Result implements [main.subcommand].
func (m *Map) Result() any {

	sb := strings.Builder{}
	sb.WriteString("Farm\tJsonCode\n")
	for _, m := range m.maps {
		fmt.Fprintf(&sb, "%s\t%s\n", m.Farm, m.Jsoncode)
	}

	if m.copy {
		clipboard.WriteAll(sb.String())
		return "copied to clipboard"
	}

	return sb.String()
}

// Run implements [main.subcommand].
func (m *Map) Run(args []string) (err error) {

	fs := flag.NewFlagSet("map", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)
	m.copy = *copy

	m.maps, err = frappe.Get[types.MapRecord](nil, frappe.List{"farm", "jsoncode"})
	if err != nil {
		return err
	}

	return
}

// Usage implements [main.subcommand].
func (m *Map) Usage() string {
	panic("unimplemented")
}
