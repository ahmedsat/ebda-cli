package main

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
)

func init() {
	err := config.Configure()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

/*

 */

func main() {

	follows, err := frappe.Get[types.FarmFollowUp](frappe.Filters{
		frappe.NewFilter("name", frappe.Like, "%FollowUp-%"),
	}, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	codes := []string{}

	for i, f := range follows {
		fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%) ", i, len(follows), float64(i)/float64(len(follows))*100)
		if slices.Contains(codes, f.FarmCode) {
			f.Farm = "test"
			r, err := frappe.UpdateDoc(f)
			if err != nil {
				fmt.Printf("f:%s =>> %s", f.Name, err)
			}
			if r.Farm != "test" {
				fmt.Printf("f:%s", f.Name)
			}
			continue
		}
		date, err := time.Parse(time.DateTime, f.Creation)
		if err != nil {
			fmt.Printf("f:%s =>> %s", f.Name, err)
			continue
		}
		f.VisitDate = date.Format(time.DateOnly)
		_, err = frappe.UpdateDoc(f)
		if err != nil {
			fmt.Printf("f:%s =>> %s", f.Name, err)
			continue
		}
		codes = append(codes, f.FarmCode)
	}

	// area()
	// OverlapKml()
	// FakeFollowUp("EG/10033", "test eng")
	// Kml()
	// deletable()
	// deleteMap()
	// deleteMapByFarmCode()
}
