package commands

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/geo"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/services"
	"github.com/ahmedsat/ebda-cli/utils"
	"github.com/gen2brain/beeep"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const (
	sheetID       = "1Dxoh8DxTy4lidAFt4YtgfZrcFP7CTmBVBD7VEDWo7rY"
	configRange   = "config"
	farmsRange    = "auto-farms"
	followUpRange = "auto-follow-up"
	pgsRange      = "auto-pgs"
	fixNamesRange = "fix-names"
	mapsRange     = "auto-maps"
	soilRange     = "auto-soil"
	errorsRange   = "errors"
)

type Update struct {
	From time.Time
	To   time.Time

	NewFarmsFrom time.Time
	NewFarmsTo   time.Time

	FollowUpFrom time.Time
	FollowUpTo   time.Time

	MapOverlapTolerance float64
	MapAreaTolerance    float64

	*utils.SyncRunner
	*services.Sheet

	errMu sync.Mutex
	errs  [][]any

	progress *mpb.Progress
	wg       sync.WaitGroup
}

func (u *Update) logErr(where string, err error) {
	if err == nil {
		return
	}
	log.Printf("[ERR] %s: %v", where, err)

	u.errMu.Lock()
	u.errs = append(u.errs, []any{
		time.Now().Format(time.RFC3339),
		where,
		err.Error(),
	})
	u.errMu.Unlock()
}

func (u *Update) guard(name string, fn func() error) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				u.logErr(name, fmt.Errorf("panic: %v\n%s", r, string(debug.Stack())))
				err = nil
			}
		}()
		if e := fn(); e != nil {
			u.logErr(name, e)
		}
		return nil // swallow non-critical errors
	}
}

func (u *Update) flushErrors() error {
	ctx := context.Background()

	if len(u.errs) == 0 {
		return u.ClearAndUpdateRange(ctx, errorsRange, [][]any{{"Time", "Stage", "Error"}})
	}

	values := make([][]any, 0, len(u.errs)+1)
	values = append(values, []any{"Time", "Stage", "Error"})
	values = append(values, u.errs...)

	return u.ClearAndUpdateRange(ctx, errorsRange, values)
}

func (u *Update) Configure() error {
	u.SyncRunner = utils.NewSyncRunner(10, 100)

	sheet, err := services.NewSheet(context.Background(), sheetID)
	if err != nil {
		return err
	}
	u.Sheet = sheet

	rng, err := sheet.ReadRange(context.Background(), configRange)
	if err != nil {
		return err
	}

	for _, row := range rng {
		if len(row) < 1 {
			continue
		}
		switch row[0] {
		case "key":
			continue
		case "from":
			u.From, err = time.Parse(utils.TimeLayout, row[1].(string))
		case "to":
			u.To, err = time.Parse(utils.TimeLayout, row[1].(string))
		case "new-farms-from":
			u.NewFarmsFrom, err = time.Parse(utils.TimeLayout, row[1].(string))
		case "new-farms-to":
			u.NewFarmsTo, err = time.Parse(utils.TimeLayout, row[1].(string))
		case "follow-up-from":
			u.FollowUpFrom, err = time.Parse(utils.TimeLayout, row[1].(string))
		case "follow-up-to":
			u.FollowUpTo, err = time.Parse(utils.TimeLayout, row[1].(string))
		case "map-overlap-tolerance":
			u.MapOverlapTolerance, err = strconv.ParseFloat(row[1].(string), 64)
			if err == nil {
				u.MapOverlapTolerance *= 4200
			}
		case "map-area-tolerance":
			u.MapAreaTolerance, err = strconv.ParseFloat(row[1].(string), 64)
			if err == nil {
				u.MapAreaTolerance *= 4200
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown config key: %s\n", row[0])
			continue
		}
		if err != nil {
			return err
		}
	}

	u.progress = mpb.New()

	return nil
}

func (u *Update) Description() string { panic("unimplemented") }
func (u *Update) Name() string        { return "update" }

func (u *Update) runWithBar(name string, fn func() error) {
	u.wg.Add(1)

	bar := u.progress.AddSpinner(1,
		mpb.PrependDecorators(decor.Name(name+": ")),
		mpb.AppendDecorators(
			decor.OnAbort(decor.OnComplete(decor.Name("running"), "done"), "failed"),
		),
	)

	go func() {
		defer u.wg.Done()

		err := func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic: %v\n%s", r, string(debug.Stack()))
				}
			}()
			return fn()
		}()

		if err != nil {
			u.logErr(name, err)
			bar.Abort(false) // shows failed state
			return
		}

		bar.Increment() // shows done
	}()
}

func (u *Update) runJob(name string, fn func() error) {
	u.wg.Go(func() {
		_ = u.guard(name, fn)()
	})
}

func (u *Update) Run(args []string) error {
	defer beeep.Alert("ebda-cli", "ebda-cli is done", nil)

	fs := flag.NewFlagSet("update", flag.ExitOnError)
	skipNewFarms := fs.Bool("skip-new-farms", false, "Skip updating new farms totals")
	skipFollowUp := fs.Bool("skip-follow-up", false, "Skip updating follow-up audits")
	skipPGS := fs.Bool("skip-pgs", false, "Skip updating PGS audits")
	skipSoils := fs.Bool("skip-soils", false, "Skip updating soil analysis")
	skipMaps := fs.Bool("skip-maps", false, "Skip updating map validation")
	fs.Parse(args)

	if err := u.Configure(); err != nil {
		return err
	}

	if !*skipNewFarms {
		u.runWithBar("NewFarms", u.NewFarms)
	}
	if !*skipSoils {
		u.runWithBar("Soils", u.Soils)
	}
	if !*skipFollowUp {
		u.runJob("FollowUp", u.FollowUp)
	}
	if !*skipPGS {
		u.runJob("PGS", u.PGS)
	}
	if !*skipMaps {
		u.runJob("Maps", u.Maps)
	}

	u.wg.Wait()
	u.progress.Shutdown()

	return u.flushErrors()
}

func (u *Update) Usage() string { panic("unimplemented") }

func (u *Update) NewFarms() error {

	ctx := context.Background()

	names, err := u.ReadRange(ctx, fixNamesRange)
	if err != nil {
		return err
	}

	namesMap := make(map[string]string)
	for _, name := range names {
		if len(name) < 2 {
			continue
		}
		namesMap[name[0].(string)] = name[1].(string)
	}

	values := [][]any{
		{
			utils.Trans("Farm Name"),
			utils.Trans("Region"),
			utils.Trans("Farm Code"),
			utils.Trans("Farmers Count"),
			utils.Trans("Area Feddan"),
			utils.Trans("Engineer Name")},
	}

	farms, err := frappe.Get[types.Farm](frappe.Filters{
		frappe.NewFilter("type", frappe.Eq, "farm"),
		frappe.NewFilter("farm_status", frappe.Neq, "Cancelled"),
		frappe.NewFilter("creation_date", frappe.Gte, u.NewFarmsFrom.Format("2006-01-02")),
		frappe.NewFilter("creation_date", frappe.Lte, u.NewFarmsTo.AddDate(0, 0, 1).Format("2006-01-02")),
	},
		[]string{"name", "arabic_name", "region", "farm_id", "total_farmers", "farm_area__feddan", "farm_application"}, nil)
	if err != nil {
		return err
	}

	missingSet := make(map[string]struct{})
	for _, farm := range farms {
		app, err := frappe.GetCached1[types.FarmApplication](farm.FarmApplicationID)
		if err != nil {
			return err
		}

		engName, ok := namesMap[app.EngineerName]
		if !ok {
			missingSet[app.EngineerName] = struct{}{}
		}

		values = append(values, []any{
			farm.ArabicName,
			utils.Trans(farm.Region),
			farm.FarmId,
			farm.TotalFarmers,
			farm.Area,
			engName,
		})
	}

	err = u.ClearAndUpdateRange(context.Background(), farmsRange, values)
	if err != nil {
		return err
	}

	// missing names
	if len(missingSet) > 0 {
		toAppend := make([][]any, 0, len(missingSet))
		for name := range missingSet {
			toAppend = append(toAppend, []any{name})
		}
		if err := u.Append(ctx, fixNamesRange, toAppend); err != nil {
			return err
		}
	}

	return nil
}

func (u *Update) FollowUp() (err error) {

	ctx := context.Background()

	names, err := u.ReadRange(ctx, fixNamesRange)
	if err != nil {
		return err
	}

	namesMap := make(map[string]string)
	for _, name := range names {
		if len(name) < 2 {
			continue
		}
		namesMap[name[0].(string)] = name[1].(string)
	}

	values := [][]any{
		{"Visit ID", "Farm Code", "Visit Date", "Creation", "Rate", "Issues", "Follower Name"},
	}

	bar := u.progress.AddBar(100,
		mpb.PrependDecorators(
			decor.Name("FollowUp: "),
			decor.Percentage(),
		),
	)

	progressCh := make(chan int, 1000)
	done := make(chan struct{})

	go func() {
		for n := range progressCh {
			bar.SetCurrent(int64(n))
		}
		defer close(done)
	}()
	defer func() {
		close(progressCh)
		<-done
		if err != nil {
			bar.Abort(true)
			return
		}
		bar.SetTotal(100, true)
	}()

	followUps, err := services.LoadFollowUps(u.FollowUpFrom, u.FollowUpTo, func(n int) {
		progressCh <- n
	})
	if err != nil {
		return err
	}

	missingSet := make(map[string]struct{})

	for _, followUp := range followUps {

		engName, ok := namesMap[followUp.FollowerName]
		if !ok {
			missingSet[followUp.FollowerName] = struct{}{}
		}

		values = append(values, []any{
			followUp.Name,
			followUp.FarmCode,
			followUp.VisitDate,
			followUp.Creation,
			followUp.RatePercent,
			strings.Join(followUp.Issues, "\n"),
			engName,
		})
	}
	err = u.ClearAndUpdateRange(context.Background(), followUpRange, values)
	if err != nil {
		return err
	}

	// missing names
	if len(missingSet) > 0 {
		toAppend := make([][]any, 0, len(missingSet))
		for name := range missingSet {
			toAppend = append(toAppend, []any{name})
		}
		if err := u.Append(ctx, fixNamesRange, toAppend); err != nil {
			return err
		}
	}

	return
}

func (u *Update) PGS() (err error) {
	ctx := context.Background()

	names, err := u.ReadRange(ctx, fixNamesRange)
	if err != nil {
		return err
	}

	namesMap := make(map[string]string)
	for _, name := range names {
		if len(name) < 2 {
			continue
		}
		namesMap[name[0].(string)] = name[1].(string)
	}

	if err := u.Clear(ctx, pgsRange); err != nil {
		return err
	}

	header := []any{"Submission ID", "Date", "Farms Code", "Validation Status", "Engineer"}

	total, dataCh, errCh := kobo.StreamAssets[kobo.PGSNew](nil)

	bar := u.progress.AddBar(int64(total),
		mpb.PrependDecorators(
			decor.Name("PGS: "),
			decor.Percentage(),
		),
	)

	progressCh := make(chan struct{}, 1000)
	done := make(chan struct{})

	// UI updater
	go func() {
		for range progressCh {
			bar.Increment()
		}
		close(done)
	}()
	defer func() {
		close(progressCh)
		<-done
		if err != nil {
			bar.Abort(true)
			return
		}
		bar.SetTotal(int64(total), true)
	}()

	missingSet := make(map[string]struct{})

	// batching
	batch := make([][]any, 0, 100)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		err := u.Append(ctx, pgsRange, batch)
		batch = batch[:0]
		return err
	}

	// write header once
	if err := u.Append(ctx, pgsRange, [][]any{header}); err != nil {
		return err
	}

	for dataCh != nil || errCh != nil {
		select {
		case audit, ok := <-dataCh:
			if !ok {
				dataCh = nil
				continue
			}

			row := []any{
				audit.ID,
				audit.AtHouse.VisitDate,
				audit.AtHouse.FarmId,
				audit.ValidationStatus.Label,
			}

			if name, ok := namesMap[audit.EngineerData.EngineerName]; ok {
				row = append(row, name)
			} else {
				row = append(row, "")
				missingSet[audit.EngineerData.EngineerName] = struct{}{}
			}

			batch = append(batch, row)

			progressCh <- struct{}{}

			if len(batch) >= 100 {
				if err := flush(); err != nil {
					return err
				}
			}

		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			if err != nil {
				return err
			}
		}
	}

	// final flush
	if err := flush(); err != nil {
		return err
	}

	// missing names
	if len(missingSet) > 0 {
		toAppend := make([][]any, 0, len(missingSet))
		for name := range missingSet {
			toAppend = append(toAppend, []any{name})
		}
		if err := u.Append(ctx, fixNamesRange, toAppend); err != nil {
			return err
		}
	}

	return nil
}

type MapData struct {
	Name     string
	FarmID   string
	Area     float64
	MapsIds  []string
	Polygons []geo.Polygon
}

func (md MapData) Overlapped(other MapData) (float64, error) {

	for _, p1 := range md.Polygons {
		for _, p2 := range other.Polygons {
			area, err := p1.OverlapArea(p2)
			if err != nil {
				return 0, err
			}
			if area > 0 {
				return area, nil
			}
		}
	}

	return 0, nil
}

func (u *Update) Maps() (err error) {

	mapsList, err := frappe.Get[types.MapRecord](nil, nil, nil)
	if err != nil {
		u.logErr("Maps", err)
		return err
	}

	bar := u.progress.AddBar(int64(len(mapsList)),
		mpb.PrependDecorators(
			decor.Name("Maps (build): "),
			decor.Percentage(),
		),
	)

	progressCh := make(chan struct{}, 1000)
	done := make(chan struct{})

	go func() {
		for range progressCh {
			bar.Increment()
		}
		close(done)
	}()
	buildBarClosed := false
	finishBuildBar := func() {
		if buildBarClosed {
			return
		}
		buildBarClosed = true
		close(progressCh)
		<-done
		if err != nil {
			bar.Abort(true)
			return
		}
		bar.SetTotal(int64(len(mapsList)), true)
	}
	defer finishBuildBar()

	MapsMap := utils.NewLockableMap[string, MapData]()
	runner := utils.NewSyncRunner(10, 100)
	c := atomic.Int64{}
	for _, m := range mapsList {
		runner.Run(func() error {
			defer func() { progressCh <- struct{}{} }()
			if m.Farm == "" {
				return nil
			}

			err := m.Parse()
			if err != nil {
				return err
			}

			MapsMap.RLock()
			data, ok := MapsMap.Map[m.Farm]
			MapsMap.RUnlock()
			if !ok {
				farm, err := frappe.GetCached1[types.Farm](m.Farm)
				if farm.FarmStatus == "Cancelled" {
					return nil
				}
				if err != nil {
					return err
				}
				data = MapData{
					Name:     farm.ArabicName,
					FarmID:   farm.FarmId,
					MapsIds:  []string{m.Name},
					Area:     farm.Area * 4200,
					Polygons: []geo.Polygon{{Ring: m.Coordinates}},
				}
				MapsMap.Lock()
				MapsMap.Map[m.Farm] = data
				MapsMap.Unlock()
			} else {
				data.MapsIds = append(data.MapsIds, m.Name)
				overlapped := false
				for j, poly := range data.Polygons {
					polys, err := poly.Union(geo.Polygon{Ring: m.Coordinates})
					if err != nil {
						return err
					}
					if len(polys) == 1 {
						data.Polygons[j] = polys[0]
						overlapped = true
						break
					}
				}
				if !overlapped {
					data.Polygons = append(data.Polygons, geo.Polygon{Ring: m.Coordinates})
				}
			}

			MapsMap.Lock()
			MapsMap.Map[m.Farm] = data
			MapsMap.Unlock()
			return nil
		})
	}

	if err := runner.Wait(); err != nil {
		return err
	}
	finishBuildBar()

	bar2 := u.progress.AddBar(int64(len(MapsMap.Map)),
		mpb.PrependDecorators(
			decor.Name("Maps (validate): "),
			decor.Percentage(),
		),
	)
	progressCh2 := make(chan struct{}, 1000)
	done2 := make(chan struct{})

	go func() {
		for range progressCh2 {
			bar2.Increment()
		}
		close(done2)
	}()
	validateBarClosed := false
	finishValidateBar := func() {
		if validateBarClosed {
			return
		}
		validateBarClosed = true
		close(progressCh2)
		<-done2
		if err != nil {
			bar2.Abort(true)
			return
		}
		bar2.SetTotal(int64(len(MapsMap.Map)), true)
	}
	defer finishValidateBar()

	values := [][]any{{"Code", "Issues"}}
	mut := sync.Mutex{}

	c.Store(0)
	for _, data := range MapsMap.Map {
		runner.Run(func() error {
			defer func() { progressCh2 <- struct{}{} }()
			row := []any{data.FarmID}
			if len(data.MapsIds) == 0 {
				row = append(row, "No maps")
				mut.Lock()
				values = append(values, row)
				mut.Unlock()
				return nil
			}

			issues := []string{}
			totalArea := 0.0
			for _, poly := range data.Polygons {
				area, err := poly.SphericalArea()
				if err != nil {
					u.logErr("Maps", err)
					continue
				}
				totalArea += area
			}

			if math.Abs(data.Area-totalArea) > u.MapAreaTolerance {
				issues = append(issues, fmt.Sprintf(utils.Trans("Area is %0.2f fed, should be %0.2f fed"), totalArea/4200, data.Area/4200))
			}

			for _, v := range MapsMap.Map {
				if v.FarmID == data.FarmID {
					continue
				}
				overlapped, err := data.Overlapped(v)
				if err != nil {
					u.logErr("Maps", err)
					continue
				}
				if overlapped > u.MapOverlapTolerance {
					issues = append(issues, fmt.Sprintf(utils.Trans("Overlapped with %s (%s) by %0.2f fed"), v.Name, v.FarmID, overlapped/4200))
				}
			}

			row = append(row, strings.Join(issues, "\n"))
			mut.Lock()
			values = append(values, row)
			mut.Unlock()
			return nil
		})
	}

	if err := runner.Wait(); err != nil {
		u.logErr("Maps", err)
		return err
	}
	finishValidateBar()

	err = u.ClearAndUpdateRange(context.Background(), mapsRange, values)
	if err != nil {
		u.logErr("Maps", err)
	}
	return err
}

func (u *Update) Soils() error {

	soils, err := frappe.Get[types.SoilAnalysis](nil, nil, nil)
	if err != nil {
		return err
	}

	values := [][]any{{"Name", "Farm", "Farm ID", "Location"}}

	for _, soil := range soils {
		values = append(values, []any{soil.Name, soil.Farm, soil.FarmID, soil.Location})
	}

	return u.ClearAndUpdateRange(context.Background(), soilRange, values)
}
