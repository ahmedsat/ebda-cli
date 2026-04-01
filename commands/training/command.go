package training

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/utils"
)

type Training struct{}

type TrainingEntry struct {
	Regions []string       `json:"regions"`
	Farms   []string       `json:"farms"`
	Codes   []string       `json:"codes"`
	Modules map[string]int `json:"modules"`
}

var TMap = map[string]TrainingEntry{}
var mu = sync.Mutex{}

func Print(data map[string]TrainingEntry) (sb strings.Builder) {

	sb.WriteString("Farmer\tRegions\tFarms\tCodes\tModule\tCount\n")

	for k, v := range data {
		for m, c := range v.Modules {
			fmt.Fprintf(
				&sb,
				"%s\t%s\t%s\t%s\t%s\t%d\n",
				k,
				strings.Join(v.Regions, ","),
				strings.Join(v.Farms, ","),
				strings.Join(v.Codes, ","),
				m,
				c,
			)
		}
	}

	return
}

func SaveToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(TMap)
}

func LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&TMap)

	return err
}

// Description implements [main.subcommand].
func (t *Training) Description() (desc string) {
	return "Training commands"
}

// Name implements [main.subcommand].
func (t *Training) Name() (name string) { return "training" }

// Run implements [main.subcommand].
func (t *Training) Run(args []string) (result any, err error) {

	if len(args) == 0 {
		return nil, errors.New("missing subcommand")
	}

	switch args[0] {
	case "get-data":
		if len(args) > 1 {
			return "", getData(args[1])
		}
		return "", getData("training.json")
	case "filter":

		if len(args) < 2 {
			return "", errors.New("training file not set")
		}

		err = LoadFromFile(args[1])
		if err != nil {
			return
		}

		expr := `region("Minya") AND (module("Marking") And Not module("Conversion))`

		opts, err := Compile(expr)
		if err != nil {
			panic(err)
		}

		filtered, err := filter(opts)
		if err != nil {
			return "", err
		}

		sb := Print(filtered)
		return sb.String(), nil

	default:
		return "", errors.New("unavailable commands: " + args[0])
	}

}

func filter(opts []Opt) (result map[string]TrainingEntry, err error) {

	result = make(map[string]TrainingEntry)

	for k, v := range TMap {
		ok, err := Evaluate(v, opts)
		if err != nil {
			return nil, err
		}
		if ok {
			result[k] = v
		}
	}

	return
}

// Usage implements [main.subcommand].
func (t *Training) Usage() (usage string) {
	panic("unimplemented")
}

func getData(file string) error {

	if !strings.HasSuffix(file, ".json") {
		file = file + ".json"
	}

	farms, err := frappe.Get[types.Farm](
		frappe.Filters{frappe.NewFilter("type", frappe.Eq, "Farm"), frappe.NewFilter("farm_status", frappe.Neq, "Cancelled")},
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	runner := utils.NewSyncRunner(10, 100)

	c := atomic.Int64{}
	for _, farm := range farms {
		runner.Run(func() (err error) {
			fmt.Fprintf(os.Stderr, "\r%d/%d", c.Add(1), len(farms))

			farm, err := frappe.Get1[types.Farm](farm.Name)
			if err != nil {
				return
			}

			for _, farmer := range farm.Farmers {
				mu.Lock()
				entry := TMap[farmer.FarmerName]

				if !slices.Contains(entry.Regions, farm.Region) {
					entry.Regions = append(entry.Regions, farm.Region)
				}

				if !slices.Contains(entry.Farms, farm.Name) {
					entry.Farms = append(entry.Farms, farm.Name)
				}

				if !slices.Contains(entry.Codes, farm.FarmId) {
					entry.Codes = append(entry.Codes, farm.FarmId)
				}

				TMap[farmer.FarmerName] = entry
				mu.Unlock()
			}
			return
		})
	}
	err = runner.Wait()
	if err != nil {
		return err
	}

	ts, err := frappe.Get[types.EbdaTraining](nil, nil, nil)
	if err != nil {
		return err
	}

	c.Store(0)
	for _, t := range ts {
		runner.Run(func() (err error) {
			fmt.Fprintf(os.Stderr, "\r%d/%d", c.Add(1), len(ts))

			t, err := frappe.Get1[types.EbdaTraining](t.Name)
			if err != nil {
				return err
			}

			for _, farmer := range t.Farmers {
				mu.Lock()
				entry := TMap[farmer.FarmerName]
				if entry.Modules == nil {
					entry.Modules = map[string]int{}
				}
				entry.Modules[t.Topic]++
				TMap[farmer.FarmerName] = entry
				mu.Unlock()
			}
			return
		})
	}
	err = runner.Wait()
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr)

	err = SaveToFile(file)
	if err != nil {
		return err
	}

	return nil
}
