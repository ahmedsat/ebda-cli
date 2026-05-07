package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/services"
	"github.com/ahmedsat/ebda-cli/utils"
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

	*utils.SyncRunner
	*services.Sheet

	errMu sync.Mutex
	errs  [][]any
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
		default:
			fmt.Fprintf(os.Stderr, "unknown config key: %s\n", row[0])
			continue
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Update) Description() string { panic("unimplemented") }
func (u *Update) Name() string        { return "update" }

func (u *Update) Run(args []string) error {
	if err := u.Configure(); err != nil {
		return err
	}

	// wrap all jobs
	// u.SyncRunner.Run(u.guard("NewFarms", u.NewFarms))
	// u.SyncRunner.Run(u.guard("FollowUp", u.FollowUp))
	// u.SyncRunner.Run(u.guard("PGS", u.PGS))
	u.SyncRunner.Run(u.guard("Maps", u.Maps))
	// u.SyncRunner.Run(u.guard("Soils", u.Soils))

	_ = u.Wait() // never fail whole run

	// single sheet write for all errors
	return u.flushErrors()
}

func (u *Update) Usage() string { panic("unimplemented") }

func (u *Update) NewFarms() error {
	values := [][]any{
		{utils.Trans("Region"), utils.Trans("Farms Count"), utils.Trans("Farmers Count"), utils.Trans("Area Feddan")},
	}

	report, err := services.LoadTotalsReport(u.NewFarmsFrom, u.NewFarmsTo)
	if err != nil {
		return err
	}

	for _, row := range report.Rows {
		values = append(values, []any{utils.Trans(row.Region), row.Farms, row.Farmers, row.Area})
	}

	values = append(values, []any{utils.Trans("Total"), report.TotalFarms, report.TotalFarmers, report.TotalArea})

	return u.ClearAndUpdateRange(context.Background(), farmsRange, values)
}

func (u *Update) FollowUp() error {

	values := [][]any{
		{"Visit ID", "Farm Code", "Visit Date", "Creation", "Rate", "Issues"},
	}

	followUps, err := services.LoadFollowUps(u.FollowUpFrom, u.FollowUpTo)
	if err != nil {
		return err
	}

	for _, followUp := range followUps {
		values = append(values, []any{
			followUp.Name,
			followUp.FarmCode,
			followUp.VisitDate,
			followUp.Creation,
			followUp.RatePercent,
			strings.Join(followUp.Issues, "\n"),
		})
	}

	return u.ClearAndUpdateRange(context.Background(), followUpRange, values)
}

func (u *Update) PGS() error {
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

	values := [][]any{
		{"Submission ID", "Date", "Farms Code", "Validation Status", "Engineer"},
	}

	dataCh, errCh := kobo.StreamAssets[kobo.PGSNew](nil)

	missingSet := make(map[string]struct{})
	processed := 0

	for {
		select {
		case audit, ok := <-dataCh:
			if !ok {
				if len(values) > 1 {
					if err := u.Append(ctx, pgsRange, values[1:]); err != nil {
						return err
					}
				}

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

			values = append(values, row)
			processed++

			if processed%50 == 0 {
				log.Printf("PGS processed: %d", processed)
			}

		case err := <-errCh:
			if err != nil {
				return err
			}
		}
	}
}

func (u *Update) Maps() error {
	mapsList, err := frappe.Get[types.MapRecord](nil, nil, nil)
	if err != nil {
		u.logErr("Maps", err)
		return err
	}

	MapsMap := make(map[string][]types.MapRecord)
	for i := range mapsList {
		err := mapsList[i].Parse()
		if err != nil {
			u.logErr("Maps", err)
		}
		m := mapsList[i]
		if m.Farm == "" {
			continue
		}

		MapsMap[m.Farm] = append(MapsMap[m.Farm], m)
	}

	mapsList = []types.MapRecord{}
	mapsListMut := sync.Mutex{}
	runner := utils.NewSyncRunner(50, 100)

	c := atomic.Int64{}
	for farm, maps := range MapsMap {
		runner.Run(func() error {
			defer fmt.Printf("\r[%d:%d]", c.Add(1), len(MapsMap))
			farm, err := frappe.Get1[types.Farm](farm)
			if err != nil {
				u.logErr("Maps", err)
				return nil
			}

			if len(maps) > 1 {
				u.logErr("Maps", fmt.Errorf("Farm (%s) has more than one Map (%d)", farm.FarmId, len(maps)))
			}

			m := types.MapRecord{}
			for _, m_ := range maps {
				// m.Area_in_feddan += m_.Area_in_feddan
				if m_.Area_in_feddan > m.Area_in_feddan {
					m = m_
				}
			}
			m.Farm = farm.FarmId
			mapsListMut.Lock()
			mapsList = append(mapsList, m)
			mapsListMut.Unlock()
			return nil
		})
	}
	err = runner.Wait()
	if err != nil {
		u.logErr("Maps", err)
	}
	fmt.Println()

	values := [][]any{{"Name", "Farm", "Area"}}
	for _, mapRecord := range mapsList {
		row := []any{mapRecord.Name, mapRecord.Farm, mapRecord.Area_in_feddan}
		values = append(values, row)
	}

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
