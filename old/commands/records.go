package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/types"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

func printBoolSlice(b []bool) {
	for _, v := range b {
		fmt.Printf("%t\t", v)
	}
	fmt.Println()
}

func Records(args []string) (err error) {

	fs := flag.NewFlagSet("records", flag.ExitOnError)

	season := fs.String("season", "", "season")
	output := fs.String("output", "", "output")
	fs.Parse(args)

	if *season == "" {
		fs.Usage()
		return fmt.Errorf("season is required")
	}

	var outputFile *os.File = os.Stdout

	if *output != "" {
		outputFile, err = os.Create(*output)
		if err != nil {
			return err
		}
	}
	args = fs.Args()

	sb := strings.Builder{}
	sb.WriteString("Code\tFertilization Record\tSowing Record\tIrrigation Record\tFarm Operations\tControl Record\tHarvest Record\n")

	for i, code := range args {
		if code == "" {
			continue
		}

		fmt.Fprintf(&sb, "%s\t", code)
		var farms []types.Farm
		farms, err = erp.Get[types.Farm](utils.Filters{utils.NewFilter("farm_id", utils.Eq, code)}, []string{"name"})
		if err != nil {
			return
		}
		if len(farms) == 0 {
			return fmt.Errorf("farm not found")
		}
		if len(farms) > 1 {
			return fmt.Errorf("more than one farm found")
		}

		r, err := CheckCropsPlan(farms[0].Name, *season)
		if err != nil {
			return err
		}
		if !r {
			sb.WriteString("false\tfalse\tfalse\tfalse\tfalse\tfalse\n")
			continue
		}

		// Fertilization Record
		r, err = CheckFertilizationRecord(farms[0].Name, *season)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sb, "%t\t", r)

		// Sowing Record
		r, err = CheckSowingRecord(farms[0].Name, *season)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sb, "%t\t", r)

		// Irrigation Record
		r, err = CheckIrrigationRecord(farms[0].Name, *season)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sb, "%t\t", r)

		// Farm Operations Record
		r, err = CheckFarmOperationsRecord(farms[0].Name, *season)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sb, "%t\t", r)

		// Control Record
		r, err = CheckControlRecord(farms[0].Name, *season)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sb, "%t\t", r)

		// Harvest Record
		r, err = CheckHarvestRecord(farms[0].Name, *season)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sb, "%t\n", r)

		fmt.Fprintf(os.Stderr, "\rProgress: [%d/%d] %0.2f%%", i+1, len(args), float64(i)/float64(len(args))*100.0)
	}

	_, err = outputFile.WriteString(sb.String())
	if err != nil {
		return err
	}

	return
}

func CheckCropsPlan(farmName, season string) (bool, error) {
	cropsPlans, err := erp.Get[types.CropPlan](
		utils.Filters{utils.NewFilter("farm", utils.Eq, farmName), utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"})
	if err != nil {
		return false, err
	}

	if len(cropsPlans) == 0 {
		return false, nil
	}

	if len(cropsPlans) > 1 {
		return false, fmt.Errorf("more than one crops plan found")
	}

	cp, err := erp.Get1[types.CropPlan](cropsPlans[0].ID)
	if err != nil {
		return false, err
	}

	if len(cp.TotalTable) == 0 {
		return false, nil
	}

	return true, nil
}

func CheckFertilizationRecord(farmName, season string) (bool, error) {
	fertilizationRecords, err := erp.Get[types.FertilizationRecord](
		utils.Filters{
			utils.NewFilter("farm", utils.Eq, farmName),
			utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"},
	)
	if err != nil {
		return false, err
	}

	if len(fertilizationRecords) == 0 {
		return false, nil
	}

	if len(fertilizationRecords) > 1 {
		return false, fmt.Errorf("more than one fertilization record found")
	}

	fr, err := erp.Get1[types.FertilizationRecord](fertilizationRecords[0].Name)
	if err != nil {
		return false, err
	}

	if len(fr.MergedFertilizations) == 0 {
		return false, nil
	}

	return true, nil
}

func CheckSowingRecord(s, season string) (bool, error) {
	sowingRecords, err := erp.Get[types.SowingRecord](
		utils.Filters{
			utils.NewFilter("farm", utils.Eq, s),
			utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"},
	)
	if err != nil {
		return false, err
	}

	if len(sowingRecords) == 0 {
		return false, nil
	}

	if len(sowingRecords) > 1 {
		return false, fmt.Errorf("more than one sowing record found")
	}

	sr, err := erp.Get1[types.SowingRecord](sowingRecords[0].Name)
	if err != nil {
		return false, err
	}

	if len(sr.SeedsTable) == 0 {
		return false, nil
	}

	return true, nil

}

func CheckIrrigationRecord(s, season string) (bool, error) {
	irrigationRecords, err := erp.Get[types.IrrigationRecord](
		utils.Filters{
			utils.NewFilter("farm", utils.Eq, s),
			utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"},
	)
	if err != nil {
		return false, err
	}

	if len(irrigationRecords) == 0 {
		return false, nil
	}

	if len(irrigationRecords) > 1 {
		return false, fmt.Errorf("more than one irrigation record found")
	}

	ir, err := erp.Get1[types.IrrigationRecord](irrigationRecords[0].Name)
	if err != nil {
		return false, err
	}

	if len(ir.IrrigationTable) == 0 {
		return false, nil
	}

	return true, nil

}

func CheckFarmOperationsRecord(s, season string) (bool, error) {
	farmOperationsRecords, err := erp.Get[types.FarmOperationsRecord](
		utils.Filters{
			utils.NewFilter("farm", utils.Eq, s),
			utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"},
	)
	if err != nil {
		return false, err
	}

	if len(farmOperationsRecords) == 0 {
		return false, nil
	}

	if len(farmOperationsRecords) > 1 {
		return false, fmt.Errorf("more than one farm operations record found")
	}

	fo, err := erp.Get1[types.FarmOperationsRecord](farmOperationsRecords[0].Name)
	if err != nil {
		return false, err
	}

	if len(fo.FarmOperationsTable) == 0 {
		return false, nil
	}

	return true, nil
}

func CheckControlRecord(s, season string) (bool, error) {
	controlRecords, err := erp.Get[types.ControlRecord](
		utils.Filters{
			utils.NewFilter("farm", utils.Eq, s),
			utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"},
	)
	if err != nil {
		return false, err
	}

	if len(controlRecords) == 0 {
		return false, nil
	}

	if len(controlRecords) > 1 {
		return false, fmt.Errorf("more than one control record found")
	}

	cr, err := erp.Get1[types.ControlRecord](controlRecords[0].Name)
	if err != nil {
		return false, err
	}

	if len(cr.ControlTable) == 0 {
		return false, nil
	}

	return true, nil
}

func CheckHarvestRecord(s, season string) (bool, error) {
	harvestRecords, err := erp.Get[types.HarvestRecord](
		utils.Filters{
			utils.NewFilter("farm", utils.Eq, s),
			utils.NewFilter("season", utils.Eq, season)},
		[]string{"name", "season", "farm"},
	)
	if err != nil {
		return false, err
	}

	if len(harvestRecords) == 0 {
		return false, nil
	}

	if len(harvestRecords) > 1 {
		return false, fmt.Errorf("more than one harvest record found")
	}

	hr, err := erp.Get1[types.HarvestRecord](harvestRecords[0].Name)
	if err != nil {
		return false, err
	}

	if len(hr.HarvestTable) == 0 {
		return false, nil
	}

	return true, nil
}
