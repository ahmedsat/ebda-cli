package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/utils"
	"gorm.io/gorm"
)

type SubmissionState struct {
	gorm.Model
	FarmersDone  bool
	SoilDone     bool
	BoundaryDone bool
}

func init() {
	config.MigrationsList = append(config.MigrationsList, &SubmissionState{})
}

var userRegionMap = map[string]string{}

type Missing struct {
}

// Description implements [main.subcommand].
func (m *Missing) Description() string {
	return "Migrating data from kobo to frappe"
}

// Name implements [main.subcommand].
func (m *Missing) Name() string {
	return "missing"
}

// Result implements [main.subcommand].
func (m *Missing) Result() any {
	return nil
}

// Run implements [main.subcommand].
func (m *Missing) Run(args []string) error {
	flagSet := flag.NewFlagSet("missing", flag.ExitOnError)
	fix := flagSet.Bool("fix", false, "Fix not approved submissions")
	flagSet.Parse(args)

	fmt.Fprintln(os.Stderr, "getting data from kobo...")
	data, err := kobo.GetAssets[kobo.Collect]()
	if err != nil {
		return err
	}

	// a, err := kobo.GetAssetByID[kobo.Collect](646267313)
	// if err != nil {
	// 	return err
	// }
	// data := []kobo.Collect{a}

	if *fix {
		fmt.Fprintln(os.Stderr, "fixing not approved submissions...")
		for _, d := range data {
			if d.CollectValidationSate.Label == "Not Approved" {
				fmt.Println(d.ID)
				fmt.Println(d.Code)
				res, err := kobo.GetUpdateURL[kobo.Collect](d.ID)
				if err != nil {
					return err
				}

				fmt.Printf("%+v", res)
				// todo: ask for confirmation if confirmed update validation states
				fmt.Scanln()
			}
		}
		return nil
	}

	for i, d := range data {
		fmt.Fprintf(os.Stderr, "\rProgress {%d} [%d:%d] (%.2f%%)", d.ID, i+1, len(data), float64(i+1)/float64(len(data))*100)
		if d.CollectValidationSate.Label == "Approved" || d.CollectValidationSate.Label == "Not Approved" {
			continue
		}
		var submissionState SubmissionState
		err = config.DB.First(&submissionState, d.ID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == gorm.ErrRecordNotFound {
			submissionState.ID = uint(d.ID)
			err = config.DB.Create(&submissionState).Error
			if err != nil {
				return err
			}
		}

		if d.Farm == "" {
			_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusNotApproved)
			if err != nil {
				return err
			}
			continue
		}

		if !submissionState.FarmersDone {
			err = HandleFarmers(&d, &submissionState)
			if err != nil {
				return err
			}
		}

		if !submissionState.SoilDone {
			err = HandleSoil(&d, &submissionState)
			if err != nil {
				return err
			}
		}

		if !submissionState.BoundaryDone {
			err = HandleBoundary(&d, &submissionState)
			if err != nil {
				return err
			}
		}

		err = config.DB.Save(&submissionState).Error
		if err != nil {
			return err
		}

		if submissionState.FarmersDone && submissionState.SoilDone && submissionState.BoundaryDone {
			_, err = kobo.UpdateValidationState[kobo.Collect](d.ID, kobo.ValidationStatusApproved)
			if err != nil {
				return err
			}
		}

	}
	fmt.Fprintln(os.Stderr, "")

	return nil
}

func HandleBoundary(collect *kobo.Collect, submissionState *SubmissionState) error {

	if submissionState.BoundaryDone {
		return nil
	}
	area := collect.AreaNew
	if area == "" {
		area = collect.AreaOld
	}
	if area == "" {
		submissionState.BoundaryDone = true
		return nil
	}

	pointsStr := strings.Split(area, ";")

	type point struct {
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	}

	var points []point
	for _, p := range pointsStr {
		p = strings.TrimSpace(p)
		parts := strings.Split(p, " ")
		if len(parts) != 4 {
			return fmt.Errorf("invalid point: %s", p)
		}
		var p = point{
			Lat: parts[0],
			Lng: parts[1],
		}
		points = append(points, p)
	}

	pointsJson, err := json.Marshal(points)
	if err != nil {
		return err
	}

	pointsJson = pointsJson[1 : len(pointsJson)-1]

	maps, err := frappe.Get[types.MapRecord](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, collect.Farm)}, nil)
	if err != nil {
		return err
	}

	for _, m := range maps {
		if m.Color != "#181818" {
			continue
		}

		err := frappe.Delete[types.MapRecord](m.Name)
		if err != nil {
			return err
		}
	}

	m := types.MapRecord{
		Farm:     collect.Farm,
		Jsoncode: string(pointsJson),
	}
	_, err = frappe.Create(m)
	if err != nil {
		return err
	}

	submissionState.BoundaryDone = true
	return nil
}

func HandleSoil(collect *kobo.Collect, submissionState *SubmissionState) error {
	s := types.SoilAnalysis{
		Farm:               collect.Farm,
		FarmID:             collect.Code,
		CollectionDatetime: collect.Today,
		NamingSeries:       "Soil-kobo-.YY.-.MM.-",
	}
	features := []utils.Feature{}

	for _, s := range collect.Points {
		lat, lng, _, _, err := s.GeoInfo()
		if err != nil {
			return err
		}
		features = append(features, utils.NewPointFeature(s.PointStr, lat, lng))
	}

	geojson := utils.NewGeoJSON(features...)
	b, err := json.Marshal(geojson)
	if err != nil {
		return err
	}
	s.Location = string(b)

	s, err = frappe.Create(s)
	if err != nil {
		return err
	}

	submissionState.SoilDone = true

	return nil
}

func HandleFarmers(collect *kobo.Collect, submissionState *SubmissionState) error {
	return nil
}

// Usage implements [main.subcommand].
func (m *Missing) Usage() string {
	panic("unimplemented")
}
