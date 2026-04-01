package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/utils"
	"github.com/atotto/clipboard"
)

type Map struct {
	maps []types.MapRecord
}

// Description implements [main.subcommand].
func (m *Map) Description() string {
	return "Get map records"
}

// Name implements [main.subcommand].
func (m *Map) Name() string {
	return "map"
}

// Run implements [main.subcommand].
func (ma *Map) Run(args []string) (r any, err error) {

	fs := flag.NewFlagSet("map", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)

	fmt.Fprintln(os.Stderr, "getting data")
	ma.maps, err = frappe.Get[types.MapRecord](nil, nil, nil)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stderr, "done")

	sbOut := strings.Builder{}
	sbErr := strings.Builder{}
	stdOut := utils.SyncIoWriter{
		Mutex:  sync.Mutex{},
		Writer: &sbOut,
	}
	stdErr := utils.SyncIoWriter{
		Mutex:  sync.Mutex{},
		Writer: &sbErr,
	}

	sbOut.WriteString("Farm\tname\tarea\n")
	total := len(ma.maps)
	i := 1
	runner := utils.NewSyncRunner(100, 0)
	for _, m := range ma.maps {
		runner.Run(func() (e error) {

			defer func() {
				fmt.Fprintf(os.Stderr, "\r[%d/%d]", i, total)
				i++
			}()
			err := m.Parse()
			if err != nil {
				fmt.Fprintln(&stdErr, ma.Name(), ":", err)
				return
			}

			fmt.Fprintf(&stdOut, "%s\t%s\t%0.2f\n", m.Farm, m.Name, m.Area_in_feddan)

			return
		})
	}

	err = runner.Wait()
	if err != nil {
		return
	}

	if *copy {
		clipboard.WriteAll(sbOut.String())
		r = "copied to clipboard"
		return
	} else {
		r = sbOut.String()
	}

	fmt.Fprintln(os.Stderr, sbErr.String())

	return
}

// Usage implements [main.subcommand].
func (m *Map) Usage() string {
	panic("unimplemented")
}
