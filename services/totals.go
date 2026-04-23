package services

import (
	"sort"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
)

type TotalsRow struct {
	Region  string
	Farms   int
	Farmers int
	Area    float64
}

type TotalsReport struct {
	Rows         []TotalsRow
	TotalFarms   int
	TotalFarmers int
	TotalArea    float64
}

func LoadTotalsReport(from, to time.Time) (TotalsReport, error) {
	report := TotalsReport{}

	farms, err := frappe.Get[types.Farm](
		frappe.Filters{
			frappe.NewFilter("type", frappe.Eq, "farm"),
			frappe.NewFilter("farm_status", frappe.Neq, "Cancelled"),
			frappe.NewFilter("creation_date", frappe.Gte, from.Format("2006-01-02")),
			frappe.NewFilter("creation_date", frappe.Lte, to.AddDate(0, 0, 1).Format("2006-01-02")),
		},
		[]string{"name", "region", "total_farmers", "farm_area__feddan"}, nil,
	)
	if err != nil {
		return report, err
	}

	byRegion := map[string]*TotalsRow{}
	for _, farm := range farms {
		row := byRegion[farm.Region]
		if row == nil {
			row = &TotalsRow{Region: farm.Region}
			byRegion[farm.Region] = row
		}

		row.Farms++
		row.Farmers += farm.TotalFarmers
		row.Area += farm.Area

		report.TotalFarms++
		report.TotalFarmers += farm.TotalFarmers
		report.TotalArea += farm.Area
	}

	report.Rows = make([]TotalsRow, 0, len(byRegion))
	for _, row := range byRegion {
		report.Rows = append(report.Rows, *row)
	}

	sort.Slice(report.Rows, func(i, j int) bool {
		return report.Rows[i].Region < report.Rows[j].Region
	})

	return report, nil
}
