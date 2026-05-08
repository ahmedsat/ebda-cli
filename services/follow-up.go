package services

import (
	"sync/atomic"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/utils"
)

func LoadFollowUps(from, to time.Time, onProgress func(int)) ([]types.FarmFollowUp, error) {

	defer func() {
		if onProgress != nil {
			onProgress(100)
		}
	}()

	results, err := frappe.Get[types.FarmFollowUp](frappe.Filters{
		frappe.NewFilter("visit_date", frappe.Gte, from.Format(time.DateOnly)),
		frappe.NewFilter("visit_date", frappe.Lte, to.AddDate(0, 0, 1).Format(time.DateOnly)), // offset by 1 day to include the last day
	}, frappe.List{"name"}, nil)
	if err != nil {
		return nil, err
	}

	c := atomic.Int64{}
	s := utils.NewSyncRunner(10, 100)
	for i := range results {
		s.Run(func() (err error) {
			err = results[i].Rate()
			if err != nil {
				return
			}
			if onProgress != nil {
				onProgress(int(float64(c.Add(1)) / float64(len(results)) * 100))
			}
			return
		})
	}
	err = s.Wait()
	if err != nil {
		return nil, err
	}
	return results, nil
}
