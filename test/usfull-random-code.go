package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/utils"
)

func deleteMapByFarmCode() {
	for _, m := range os.Args[1:] {
		f, err := types.GetFarmByCode(m)
		if err != nil {
			fmt.Println(m, err)
			continue
		}
		maps, err := frappe.Get[types.MapRecord](frappe.Filters{
			frappe.NewFilter("farm", frappe.Eq, f.Name),
		}, nil, nil)
		if err != nil {
			fmt.Println(m, err)
			continue
		}

		for _, m := range maps {
			err := frappe.Delete[types.MapRecord](m.Name)
			if err != nil {
				fmt.Println(m.Name, err)
			}
		}
		fmt.Println(m)
	}
}

func deleteOldMap() {
	for _, m := range os.Args[1:] {
		f, err := types.GetFarmByCode(m)
		if err != nil {
			fmt.Println(m, err)
			continue
		}
		maps, err := frappe.Get[types.MapRecord](frappe.Filters{
			frappe.NewFilter("farm", frappe.Eq, f.Name),
		}, nil, nil)
		if err != nil {
			fmt.Println(m, err)
			continue
		}

		for _, m := range maps {
			if m.Color == "#181818" {
				continue
			}
			err := frappe.Delete[types.MapRecord](m.Name)
			if err != nil {
				fmt.Println(m.Name, err)
			}
		}
		fmt.Println(m)
	}
}

func deleteMap() {
	for _, m := range os.Args[1:] {
		err := frappe.Delete[types.MapRecord](m)
		if err != nil {
			fmt.Println(m, err)
		}
	}
}

func deletable() {
	maps, err := frappe.Get[types.MapRecord](frappe.Filters{
		frappe.NewFilter("color", frappe.Neq, "#181818"),
	}, nil, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	runner := utils.NewSyncRunner(10, 100)
	c := atomic.Int64{}
	for _, m := range maps {
		runner.Run(func() error {
			fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%)", c.Add(1), len(maps), float64(c.Load())/float64(len(maps))*100)
			if m.Farm == "" {
				return nil
			}
			farm, err := frappe.GetCached1[types.Farm](m.Farm)
			if err != nil {
				return err
			}
			if farm.Region == "Galvina" {
				fmt.Println(strings.Join([]string{
					farm.ArabicName,
					farm.Region,
					m.Name,
					m.Color,
				}, "\t"))
			}
			return nil
		})
	}

	err = runner.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%)\n", c.Load(), len(maps), float64(c.Load())/float64(len(maps))*100)
}

func FakeFollowUp(code, eng string) {

	farm, err := types.GetFarmByCode(code)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	compostRate := rand.Intn(7) + 8
	f := types.FarmFollowUp{
		Farm:               farm.Name,
		FollowerType:       "مهندس زراعى",
		GPS:                fmt.Sprintf("%s,%s", farm.Latitude, farm.Longitude),
		FollowerName:       eng,
		CurrentCrops:       []types.CropFollowUp{{Crops: "Corn"}},
		CompostProduction:  float64(compostRate),
		CompostQtys:        farm.Area * float64(compostRate),
		CurrentChallenges:  "لا يوجد",
		FollowerAssessment: "المزرعة بحالة جيدة",
		Recommendations:    "لا يوجد",
		StorageExist:       "لا",
		VisitDate:          time.Now().Format(time.DateOnly),
	}

	f.Name = fmt.Sprintf("FollowUp-%s-%s", f.Farm, time.Now().Format(time.DateTime))

	_, err = frappe.Create(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func OverlapKml() {
	var maps []types.MapRecord
	for _, code := range os.Args[1:] {
		f, err := types.GetFarmByCode(code)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		l, err := types.GetMapColored(f.Name, RandomHexColor(f.FarmId))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		maps = append(maps, l...)
	}
	bytes, err := types.MapRecordsToKML(maps)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = os.WriteFile("out.kml", bytes, 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func RandomHexColor(seed string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(seed))
	r := rand.New(rand.NewSource(int64(h.Sum64())))
	return fmt.Sprintf("#%06x", r.Intn(0x1000000))
}

func Kml() {
	maps, err := frappe.Get[types.MapRecord](nil, nil, nil)
	if err != nil {
		panic(err)
	}

	// regions := map[string][]types.MapRecord{}
	regions := utils.NewLockableMap[string, []types.MapRecord]()
	runner := utils.NewSyncRunner(50, 100)

	i := atomic.Int64{}
	for _, m := range maps {
		runner.Run(func() error {
			defer func() {
				i := i.Add(1)
				fmt.Printf("\r%d/%d (%0.2f%%)", i, len(maps), float64(i)/float64(len(maps))*100)
			}()
			if m.Farm == "" {
				return nil
			}
			f, err := frappe.GetCached1[types.Farm](m.Farm)
			if err != nil {
				return err
			}

			m.Color = RandomHexColor(f.FarmId)
			regions.Lock()
			regions.Map[f.Region] = append(regions.Map[f.Region], m)
			regions.Unlock()
			return nil
		})
	}

	runner.Wait()
	fmt.Println()

	runner = utils.NewSyncRunner(len(regions.Map), 0)
	for k, v := range regions.Map {
		runner.Run(func() error {
			bytes, err := types.MapRecordsToKML(v)
			if err != nil {
				return err
			}
			err = os.WriteFile(fmt.Sprintf("kml/%s.kml", k), bytes, 0666)
			if err != nil {
				return err
			}
			return nil
		})
	}
	runner.Wait()
}

func area() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "no map id")
	}

	for _, id := range os.Args[1:] {
		record, err := frappe.GetCached1[types.MapRecord](id)
		if err != nil {
			fmt.Printf("Fetch error: %s => %s", id, err)
			continue
		}

		err = record.Parse()
		if err != nil {
			fmt.Printf("Parse error: %s => %s", id, err)
			continue
		}

		fmt.Printf("%s => %0.2f\n", id, record.Area_in_feddan)
	}

}
