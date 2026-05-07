package services

import (
	"fmt"
	"os"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/utils"
)

func LoadFollowUps(from time.Time, to time.Time) ([]types.FarmFollowUp, error) {
	results, err := frappe.Get[types.FarmFollowUp](frappe.Filters{
		frappe.NewFilter("visit_date", frappe.Gte, from.Format(time.DateOnly)),
		frappe.NewFilter("visit_date", frappe.Lte, to.AddDate(0, 0, 1).Format(time.DateOnly)), // offset by 1 day to include the last day
	}, frappe.List{"name"}, nil)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stderr, "Calculating rates...")
	counter := 1
	s := utils.NewSyncRunner(10, 0)
	for i := range results {
		s.Run(func() (err error) {
			err = results[i].Rate()
			if err != nil {
				return
			}
			fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%)", counter, len(results), float64(counter)/float64(len(results))*100)
			counter++
			return
		})
	}
	err = s.Wait()
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return nil, err
	}
	return results, nil
}
