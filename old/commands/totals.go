package commands

import (
	"fmt"

	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

type TotalsOptions struct {
	GlobalOptions
}

func Totals(opt TotalsOptions) (err error) {

	err = opt.Validate()
	if err != nil {
		return
	}

	filters := utils.Filters{
		utils.NewFilter("type", utils.Eq, "farm"),
		utils.NewFilter("creation_date", utils.Gte, opt.From.Format("2006-01-02")),
		// ? to is excluded, so we have to add 1 day
		utils.NewFilter("creation_date", utils.Lte, opt.To.AddDate(0, 0, 1).Format("2006-01-02")),
	}

	if !opt.IncludeCanceled {
		filters = append(filters, utils.NewFilter("farm_status", utils.Neq, "Cancelled"))
	}

	farms, err := erp.Get[Farm](filters, []string{"name", "region", "total_farmers", "farm_area__feddan"})
	if err != nil {
		return
	}

	regions := make(map[string]struct {
		count        int
		area         float64
		totalFarmers int
	})
	for _, farm := range farms {
		region := regions[farm.Region]
		region.count++
		region.area += farm.Area
		region.totalFarmers += farm.TotalFarmers
		regions[farm.Region] = region

		total := regions["total"]
		total.count++
		total.area += farm.Area
		total.totalFarmers += farm.TotalFarmers
		regions["total"] = total
	}

	t := utils.TableBase{}

	t.SetHeader("region", "total_farmers", "count of farms", "area")
	for k, region := range regions {
		if k == "total" {
			continue
		}

		t.AppendRow(
			k,
			fmt.Sprintf("%d", region.totalFarmers),
			fmt.Sprintf("%d", region.count),
			fmt.Sprintf("%.2f", region.area),
		)
	}

	t.AppendRow("total", fmt.Sprintf("%d", regions["total"].totalFarmers), fmt.Sprintf("%d", regions["total"].count), fmt.Sprintf("%.2f", regions["total"].area))

	data := opt.TablePrinter(&t)
	opt.Print(data)

	return nil
}
