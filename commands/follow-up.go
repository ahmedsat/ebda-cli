package commands

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/atotto/clipboard"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/utils"
)

type FollowUpCommand struct {
	copy    bool
	results []types.FarmFollowUp
}

// Description implements [main.subcommand].
func (f *FollowUpCommand) Description() string {
	return "Calculating rates and printing results"
}

// Name implements [main.subcommand].
func (f *FollowUpCommand) Name() string {
	return "follow-up"
}

// Result implements [main.subcommand].
func (f *FollowUpCommand) Result() any {

	fmt.Fprintf(os.Stderr, "Printing results[%d]...\n", len(f.results))
	var res strings.Builder
	for _, result := range f.results {
		if !result.Rated {
			continue
		}
		fmt.Fprintln(&res, strings.Join([]string{
			result.Name,
			result.FarmCode,
			fmt.Sprintf("%f", result.RatePercent),
			strings.Join(result.Issues, " - "),
		}, "\t"))
	}

	if f.copy {
		clipboard.WriteAll(res.String())
		return "copied to clipboard"
	}

	return res.String()
}

// Run implements [main.subcommand].
func (f *FollowUpCommand) Run(args []string) error {

	fs := flag.NewFlagSet("follow-up", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	fs.Parse(args)

	f.copy = *copy

	results, err := frappe.Get[types.FarmFollowUp](nil, frappe.List{"name"})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Calculating rates...")
	counter := 1
	s := utils.NewSyncRunner(10, 0)
	for i := range results {
		s.Run(func() {
			err := results[i].Rate()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%)", counter, len(results), float64(counter)/float64(len(results))*100)
			counter++
		})
	}
	s.Wait()
	fmt.Fprintln(os.Stderr)

	fmt.Fprintln(os.Stderr, "Sorting results...")
	slices.SortFunc(results, func(f1, f2 types.FarmFollowUp) int {
		return int(f2.RatePercent*100) - int(f1.RatePercent*100)
	})

	f.results = results
	return nil
}

// Usage implements [main.subcommand].
func (f *FollowUpCommand) Usage() string {
	panic("unimplemented")
}

func FollowUp(args []string) error {

	results, err := frappe.Get[types.FarmFollowUp](nil, frappe.List{"name"})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Calculating rates...")
	counter := 1
	s := utils.NewSyncRunner(10, 0)
	for i := range results {
		s.Run(func() {
			err := results[i].Rate()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%)", counter, len(results), float64(counter)/float64(len(results))*100)
			counter++
		})
	}
	s.Wait()
	fmt.Fprintln(os.Stderr)

	fmt.Fprintln(os.Stderr, "Sorting results...")
	slices.SortFunc(results, func(f1, f2 types.FarmFollowUp) int {
		return int(f2.RatePercent*100) - int(f1.RatePercent*100)
	})

	fmt.Fprintln(os.Stderr, "Printing results...")
	for _, result := range results {
		if !result.Rated {
			continue
		}
		fmt.Println(strings.Join([]string{
			result.Name,
			result.FarmCode,
			fmt.Sprintf("%f", result.RatePercent),
			strings.Join(result.Issues, " - "),
		}, "\t"))
	}

	return nil
}
