package commands

import (
	_ "embed"
	"flag"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/utils"
	"github.com/atotto/clipboard"
)

const MapOk = ""

var exceptionMaps = []string{"EG/1247", "EG/1248", "EG/1249"}

type output struct {
	// target
	code           string
	arabic_name    string
	Region         string
	countOfFarmers int
	area           float64
	eng            string
	date           string

	// follow up
	countOfFollowUps int
	visitRate        float64
	issues           string
	Map              string
	hasSoil          bool
	// rate             float64

	// pgs
	countOfPGS int
	status     string
	auditorEng string
}

type Farm struct {
	farms           []types.Farm
	applications    []types.FarmApplication
	maps            []types.MapRecord
	solis           []types.SoilAnalysis
	PGSs            []kobo.PGSNew
	followUps       []types.FarmFollowUp
	io              *utils.SyncIoWriter
	applicationsMap utils.LockableMap[string, int]    // map[string]int
	farmsMap        utils.LockableMap[string, int]    // map[string]int
	codesMap        utils.LockableMap[string, int]    // map[string]int
	outputMap       utils.LockableMap[string, output] // map[string]output

	followUpFrom time.Time
	followUpTo   time.Time
}

// Description implements [main.subcommand].
func (f *Farm) Description() string {
	return "information about farms"
}

// Name implements [main.subcommand].
func (f *Farm) Name() string {
	return "farms"
}

// Run implements [main.subcommand].
func (f *Farm) Run(args []string) (err error) {

	fs := flag.NewFlagSet("farms", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	followUpFrom := fs.String("follow-up-from", "1-1-2022", "Date part of ISO")
	followUpTo := fs.String("to", time.Now().Format(utils.TimeLayout), "Date part of ISO")
	fs.Parse(args)

	from, err := time.Parse(utils.TimeLayout, *followUpFrom)
	if err != nil {
		return
	}
	to, err := time.Parse(utils.TimeLayout, *followUpTo)
	if err != nil {
		return
	}

	f.followUpFrom = from
	f.followUpTo = to

	f.farmsMap = utils.NewLockableMap[string, int]()        //make(map[string]int)
	f.applicationsMap = utils.NewLockableMap[string, int]() //make(map[string]int)
	f.codesMap = utils.NewLockableMap[string, int]()        //make(map[string]int)
	f.outputMap = utils.NewLockableMap[string, output]()    //make(map[string]output)

	f.io = &utils.SyncIoWriter{Writer: os.Stderr}
	runner := utils.NewSyncRunner(1, 0)
	fmt.Fprintln(f.io, "getting data")
	runner.Run(f.getApplications)
	runner.Run(f.getFarms)
	runner.Run(f.getMaps)
	runner.Run(f.getSolis)
	runner.Run(f.getFollowUps)
	runner.Run(f.getPGSs)

	err = runner.Wait()
	if err != nil {
		return err
	}

	err = runner.Wait()
	if err != nil {
		return err
	}

	sb := strings.Builder{}

	fmt.Fprintln(&sb, strings.Join([]string{
		"code",
		"arabic_name",
		"region",
		"countOfFarmers",
		"area",
		"eng",
		"date",
		"countOfFollowUps",
		// "visitRate",
		"issues",
		"rate",
		"countOfPGS",
		"status",
		"auditorEng",
	}, "\t"))
	for k, v := range f.outputMap.Map {

		rate := 2 + v.visitRate/float64(v.countOfFollowUps)

		if v.Map != MapOk {
			if v.Map == "" {
				v.Map = "بدون خريطة"
			}
			rate--
			v.issues = strings.Join([]string{v.issues, v.Map}, "\n")
		}

		if !v.hasSoil {
			rate--
			v.issues = strings.Join([]string{v.issues, "بدون عينة تربة"}, "\n")
		}

		fmt.Fprintln(&sb, strings.Join([]string{
			k,
			v.arabic_name,
			v.Region,
			fmt.Sprintf("%d", v.countOfFarmers),
			fmt.Sprintf("%.2f", v.area),
			v.eng,
			v.date,
			fmt.Sprintf("%d", v.countOfFollowUps),
			// fmt.Sprintf("%.2f", v.visitRate/float64(v.countOfFollowUps)),
			fmt.Sprintf("\"%s\"", strings.TrimSpace(v.issues)),
			fmt.Sprintf("%.2f", rate/3),
			fmt.Sprintf("%d", v.countOfPGS),
			fmt.Sprintf("\"%s\"", strings.TrimSpace(v.status)),
			fmt.Sprintf("\"%s\"", strings.TrimSpace(v.auditorEng)),
		}, "\t"))
	}

	if *copy {
		err = clipboard.WriteAll(sb.String())
		if err != nil {
			return err
		}
	} else {
		fmt.Print(sb.String())
	}

	return nil
}

func (f *Farm) getFollowUps() (err error) {
	fmt.Fprintln(f.io, "getting followUps data ")
	f.followUps, err = frappe.Get[types.FarmFollowUp](frappe.Filters{
		frappe.NewFilter("visit_date", frappe.Gte, f.followUpFrom.Format(time.DateOnly)),
		frappe.NewFilter("visit_date", frappe.Lte, f.followUpTo.AddDate(0, 0, 1).Format(time.DateOnly)), // offset by 1 day to include the last day
	}, frappe.List{"name"}, nil)
	if err != nil {
		return err
	}

	fmt.Fprintln(f.io, "followUps rates calculation...")
	runner := utils.NewSyncRunner(10, 100)
	atomicCount := atomic.Int64{}
	for i := range f.followUps {
		runner.Run(func() (err error) {
			defer func() {
				atomicCount.Add(1)
				fmt.Fprintf(f.io, "\rrating followUps %d/%d (%.2f%%)",
					atomicCount.Load(), len(f.followUps), float64(atomicCount.Load())/float64(len(f.followUps))*100)
			}()
			err = f.followUps[i].Rate()
			if err != nil {
				return
			}
			f.outputMap.Lock()
			output := f.outputMap.Map[f.followUps[i].FarmCode]

			output.countOfFollowUps++
			output.visitRate += f.followUps[i].RatePercent

			output.issues = strings.Join(append([]string{output.issues}, f.followUps[i].Issues...), "\n")

			f.outputMap.Map[f.followUps[i].FarmCode] = output

			f.outputMap.Unlock()
			return
		})
	}
	err = runner.Wait()
	if err != nil {
		return
	}
	fmt.Fprintln(f.io)
	return nil
}

func (f *Farm) getPGSs() (err error) {
	fmt.Fprintln(f.io, "getting PGSs data")
	start := 0
	res, err := kobo.GetAssetsExt[(kobo.PGSNew)](nil, 0, start)
	if err != nil {
		return
	}

	count := res.Count

	f.PGSs = append(f.PGSs, res.Results...)
	fmt.Fprintf(f.io, "\rPGS progress [%d : %d]", len(f.PGSs), count)

	for res.Next != "" {
		start += len(res.Results)
		res, err = kobo.GetAssetsExt[(kobo.PGSNew)](nil, 0, start)
		if err != nil {
			return
		}
		f.PGSs = append(f.PGSs, res.Results...)
		fmt.Fprintf(f.io, "\rPGS progress [%d : %d]", len(f.PGSs), count)
	}

	for i := range f.PGSs {
		f.outputMap.Lock()
		output := f.outputMap.Map[f.PGSs[i].AtHouse.FarmId]
		output.countOfPGS++
		output.status = fmt.Sprintf("%s\n%s", output.status, f.PGSs[i].ValidationStatus.Label)
		output.auditorEng = fmt.Sprintf("%s\n%s", output.auditorEng, f.PGSs[i].EngineerData.EngineerName)
		f.outputMap.Map[f.PGSs[i].AtHouse.FarmId] = output
		f.outputMap.Unlock()
	}

	fmt.Fprintln(f.io)
	return
}

func (f *Farm) getSolis() (err error) {
	fmt.Fprintln(f.io, "getting solis data")
	f.solis, err = frappe.Get[types.SoilAnalysis](nil, frappe.List{"farm", "location"}, nil)
	if err != nil {
		return err
	}

	for i := range f.solis {
		if f.solis[i].Location == "" {
			continue
		}
		f.farmsMap.Lock()
		farm := f.farms[f.farmsMap.Map[f.solis[i].Farm]]
		f.outputMap.Lock()
		output := f.outputMap.Map[farm.FarmId]

		output.hasSoil = true
		f.outputMap.Map[farm.FarmId] = output

		f.outputMap.Unlock()
		f.farmsMap.Unlock()
	}

	return nil
}

func (f *Farm) getMaps() (err error) {
	fmt.Fprintln(f.io, "getting maps data")
	f.maps, err = frappe.Get[types.MapRecord](nil, nil, nil)
	if err != nil {
		return
	}

	slices.SortFunc(f.maps, func(a, b types.MapRecord) int {
		// creationA, err := time.Parse(frappe.TimeLayout, a.Base.Creation)
		// utils.Assert(err != nil, err.Error())

		// creationB, err := time.Parse(frappe.TimeLayout, b.Base.Creation)
		// utils.Assert(err != nil, err.Error())

		// return int(creationA.Unix() - creationB.Unix())
		return int(a.Area_in_feddan - b.Area_in_feddan)
	})

	fmt.Fprintln(f.io, "parsing maps data")
	runner := utils.NewSyncRunner(100, 100)
	for i := range f.maps {
		runner.Run(func() (err error) {
			err = f.maps[i].Parse()
			if err != nil {
				return
			}

			f.farmsMap.Lock()
			defer f.farmsMap.Unlock()
			f.outputMap.Lock()
			defer f.outputMap.Unlock()

			farm := f.farms[f.farmsMap.Map[f.maps[i].Farm]]
			output := f.outputMap.Map[farm.FarmId]

			if math.Abs(f.maps[i].Area_in_feddan-farm.Area) > 0.25 {
				output.Map = fmt.Sprintf("مساحة غير متطابقة %0.2f %0.2f", f.maps[i].Area_in_feddan, farm.Area)
			} else {
				output.Map = MapOk
			}

			if slices.Contains(exceptionMaps, farm.FarmId) {
				output.Map = MapOk
			}

			f.outputMap.Map[farm.FarmId] = output
			return
		})
	}
	err = runner.Wait()
	if err != nil {
		return
	}
	return
}

func (f *Farm) getApplications() (err error) {
	fmt.Fprintln(f.io, "getting applications data")
	f.applications, err = frappe.Get[types.FarmApplication](nil, frappe.List{"name", "engineer_name", "user_name", "farm_name"}, nil)
	if err != nil {
		return err
	}

	for i := range f.applications {
		f.applicationsMap.Lock()
		f.applicationsMap.Map[f.applications[i].Name] = i
		f.applicationsMap.Unlock()
	}

	return nil
}

func (f *Farm) getFarms() (err error) {
	fmt.Fprintln(f.io, "getting farms data")
	f.farms, err = frappe.Get[types.Farm](frappe.Filters{
		frappe.NewFilter("type", frappe.Eq, "Farm"),
		frappe.NewFilter("farm_status", frappe.Neq, "Cancelled"),
	},
		frappe.List{
			"name",
			"farm_id",
			"arabic_name",
			"farm_application",
			"region",
			"total_farmers",
			"farm_area__feddan",
			"creation_date",
		}, nil)
	if err != nil {
		return err
	}

	for i := range f.farms {

		f.farmsMap.Lock()
		f.farmsMap.Map[f.farms[i].Name] = i

		f.codesMap.Lock()
		f.codesMap.Map[f.farms[i].FarmId] = i

		f.applicationsMap.Lock()
		app := f.applicationsMap.Map[f.farms[i].FarmApplication]

		f.outputMap.Lock()
		output := f.outputMap.Map[f.farms[i].FarmId]
		output.code = f.farms[i].FarmId
		output.arabic_name = f.farms[i].ArabicName
		output.Region = f.farms[i].Region
		output.countOfFarmers = f.farms[i].TotalFarmers
		output.area = f.farms[i].Area
		output.eng = f.applications[app].EngineerName
		output.date = f.farms[i].CreationDate
		f.outputMap.Map[f.farms[i].FarmId] = output

		f.farmsMap.Unlock()
		f.codesMap.Unlock()
		f.applicationsMap.Unlock()
		f.outputMap.Unlock()
	}
	return nil
}

// Usage implements [main.subcommand].
func (f *Farm) Usage() string {
	panic("unimplemented")
}
