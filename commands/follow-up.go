//go:build !release

package commands

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
	"github.com/ahmedsat/ebda-cli/frappe/types"
	"github.com/ahmedsat/ebda-cli/kobo"
	"github.com/ahmedsat/ebda-cli/services"
	"github.com/atotto/clipboard"
	"gorm.io/gorm"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/utils"
)

type FollowUpCommand struct {
}

// Description implements [main.subcommand].
func (f *FollowUpCommand) Description() string {
	return "Calculating rates and printing results"
}

// Name implements [main.subcommand].
func (f *FollowUpCommand) Name() string {
	return "follow-up"
}

// Run implements [main.subcommand].
func (f *FollowUpCommand) Run(args []string) error {

	fs := flag.NewFlagSet("follow-up", flag.ExitOnError)
	copy := fs.Bool("copy", false, "Copy to clipboard")
	migrate := fs.Bool("migrate", false, "Migrate data")
	followUpFrom := fs.String("from", "1-1-2022", "Date part of ISO")
	followUpTo := fs.String("to", time.Now().Format(utils.TimeLayout), "Date part of ISO")
	fs.Parse(args)

	if *migrate {
		return migrateFollowUp()
	}

	from, err := time.Parse(utils.TimeLayout, *followUpFrom)
	if err != nil {
		return err
	}
	to, err := time.Parse(utils.TimeLayout, *followUpTo)
	if err != nil {
		return err
	}

	results, err := services.LoadFollowUps(from, to, func(i int) {
		fmt.Fprintf(os.Stderr, "\r%d%%", i)
	})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Sorting results...")
	slices.SortFunc(results, func(f1, f2 types.FarmFollowUp) int {
		return int(f2.RatePercent*100) - int(f1.RatePercent*100)
	})

	fmt.Fprintf(os.Stderr, "Printing results[%d]...\n", len(results))
	var res strings.Builder
	fmt.Fprintln(&res, strings.Join([]string{
		"ID",
		"Farm Code",
		"FollowerName",
		"Rate",
		"Issues",
		"Visit Date",
		"Email",
	}, "\t"))
	for _, result := range results {
		if !result.Rated {
			continue
		}
		fmt.Fprintln(&res, strings.Join([]string{
			result.Name,
			result.FarmCode,
			result.FollowerName,
			fmt.Sprintf("%f", result.RatePercent),
			strings.Join(result.Issues, " - "),
			result.VisitDate,
			result.Owner,
		}, "\t"))
	}

	if *copy {
		clipboard.WriteAll(res.String())
		fmt.Println("copied to clipboard")
		return nil
	}

	fmt.Print(res.String())
	return nil

}

// Usage implements [main.subcommand].
func (f *FollowUpCommand) Usage() string {
	panic("unimplemented")
}

type SubmissionFollowUpState struct {
	gorm.Model
	FormIDOnFrappe string
	Done           bool
}

func init() {
	config.MigrationsList = append(config.MigrationsList, &SubmissionFollowUpState{})
}

func migrateFollowUp() error {
	query := kobo.Query{
		kobo.ValidationStatusKey: nil,
		"_id":                    732413723,
	}
	forms, err := kobo.GetAssets[kobo.FarmFollowupForm](query)
	if err != nil {
		return err
	}

	for _, form := range forms {
		var state SubmissionFollowUpState
		err = config.DB.First(&state, form.ID).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if err == gorm.ErrRecordNotFound {
			state.ID = uint(form.ID)
			err = config.DB.Create(&state).Error
			if err != nil {
				return err
			}
		}

		if state.Done {
			return nil
		}

		var followUp types.FarmFollowUp
		if state.FormIDOnFrappe != "" {
			followUp, err = frappe.GetCached1[types.FarmFollowUp](state.FormIDOnFrappe)
			if err != nil {
				return err
			}
		}

		if followUp.Name == "" {
			newFollowUp := types.FarmFollowUp{
				Farm:     form.Farm,
				FarmName: form.Farm,
				FarmCode: form.FarmCode,
				Region:   form.Region,
				// AreaFeddan:       form.AreaFeddan,
				Latitude:         form.Latitude,
				Longitude:        form.Longitude,
				UpscalingProject: form.UpscalingProject,
				// SubUpscalingProject:         form.SubUpscalingProject,
				FarmGroup:         form.FarmGroup,
				FarmOwner:         form.FarmOwner,
				Phone:             form.Phone,
				PictureOfFollower: form.PictureOfFollower,
				WarehousesNotes:   form.WarehousesNotes,
				RecordImage:       form.RecordImage,
				AnimalsTypeCount:  form.AnimalsTypeCount,
				// UsedQuantities:              form.UsedQuantities,
				BioControlDetails:      form.BioControlDetails,
				ChemicalControlDetails: form.ChemicalControlDetails,
				NaturalEnemiesDetails:  form.NaturalEnemiesDetails,
				WeedDisposalOther:      form.WeedDisposalOther,
				IrrigationMethodOther:  form.IrrigationMethodOther,
				WaterSourceOther:       form.WaterSourceOther,
				EnergyDetails:          form.EnergyDetails,
				// SoilAnalysisDetails:         form.SoilAnalysisDetails,
				CurrentChallenges:  form.CurrentChallenges,
				SupportNeededOther: form.SupportNeededOther,
				FollowerAssessment: form.FollowerAssessment,
				// Recommendations:             form.Recommendations,
				GPS:          form.Gps,
				VisitDate:    form.VisitDate,
				FollowerType: form.FollowerType,
				FollowerName: form.FollowerName,
				// FarmersCount:                form.FarmersCount,
				// LastEBDAVisit:               form.LastEBDAVisit,
				// HasHealthInsuranceCard:      form.HasHealthInsuranceCard,
				StorageExist: form.StorageExist,
				// RecordsFarmBook:             form.RecordsFarmBook,
				IntercroppingOrGreenManure: form.IntercroppingOrGreenManure,
				// IntercroppingPercent:        form.IntercroppingPercent,
				PlantedTreesOrHedge: form.PlantedTreesOrHedge,
				// TreesCount:                  form.TreesCount,
				HasAnimals:          form.HasAnimals,
				FertilizationSource: form.FertilizationSource,
				CompostSource:       form.CompostSource,
				// CompostProduction:           form.CompostProduction,
				// CompostQtys:                 form.CompostQtys,
				UsesBioProducts: form.UsesBioProducts,
				// CompQty:                     form.CompQty,
				// Qurts:                       form.Qurts,
				// Horns:                       form.Horns,
				PestsOrDiseasesLastSeason: form.PestsOrDiseasesLastSeason,
				PestControlMethod:         form.PestControlMethod,
				// UsesNaturalEnemies:          form.UsesNaturalEnemies,
				WeedDisposal:     form.WeedDisposal,
				IrrigationMethod: form.IrrigationMethod,
				WaterSource:      form.WaterSource,
				// WaterShortageOrSalinity:     form.WaterShortageOrSalinity,
				EnergyType: form.EnergyType,
				// ClimateChallengesThisSeason: form.ClimateChallengesThisSeason,
				YieldVsLastSeason: form.YieldVsLastSeason,
				// NeedCertSupport:             form.NeedCertSupport,
				SupportNeeded:       form.SupportNeeded,
				Doctype:             "",
				ServicesUsed:        []any{},
				FarmersNames:        []types.FarmerFollowUp{},
				BiosProductsDetails: []types.BioProductFollowUp{},
				CurrentCrops:        []types.CropFollowUp{},
				RatePercent:         0,
				Issues:              []string{},
				Rated:               false,
			}

			t := form.Start
			if t == "" {
				t = form.SubmissionTime
			}

			newFollowUp.Creation = t[:10]

			checks := []utils.Check{{Name: "لايوجد موقع", Ok: newFollowUp.GPS != "", Weight: 3},
				{Name: "اسم المتابع غير موجود", Ok: newFollowUp.FollowerName != "", Weight: 5},
				{Name: "صورة المتابع مع المزارعين غير موجودة", Ok: newFollowUp.PictureOfFollower != "", Weight: 5},
				{Name: "لا يوجد محاصيل", Ok: len(newFollowUp.CurrentCrops) != 0, Weight: 3},
				{Name: "معدل انتاج الكمبوست غير موجود", Ok: newFollowUp.CompostProduction > 0, Weight: 3},
				{Name: "كمية الكمبوست غير موجودة", Ok: newFollowUp.CompostQtys > 0, Weight: 3},
				{Name: "لم يتم ذكر التحديات الحالية", Ok: newFollowUp.CurrentChallenges != "", Weight: 1},
				{Name: "لم يتم ذكر تقييم المراجع", Ok: newFollowUp.FollowerAssessment != "", Weight: 1},
				{Name: "لم يتم ذكر التوصيات", Ok: newFollowUp.Recommendations != "", Weight: 4},
				{Name: "لم يتم ذكر هل يوجد مخزن ام لا", Ok: newFollowUp.StorageExist != "", Weight: 3}}
			if newFollowUp.RecordsFarmBook != 0 {
				checks = append(checks, utils.Check{Name: "صورة دفتر المزرعة غير موجودة", Ok: newFollowUp.RecordImage != "1", Weight: 3})
			}
			if newFollowUp.StorageExist == "نعم" {
				checks = append(checks, utils.Check{Name: "لم يتم ذكر محتويات المخزن", Ok: newFollowUp.WarehousesNotes != "", Weight: 3})
			}
			if newFollowUp.IntercroppingOrGreenManure == "نعم" {
				checks = append(checks, utils.Check{Name: "لم يتم ذكر نسبة التحميل", Ok: newFollowUp.IntercroppingPercent > 0, Weight: 3})
			}
			if newFollowUp.PlantedTreesOrHedge == "نعم" {
				checks = append(checks, utils.Check{Name: "لم يتم ذكر عدد الاشجار الجديدة", Ok: newFollowUp.TreesCount > 0, Weight: 3})
			}
			if newFollowUp.HasAnimals == "نعم" {
				checks = append(checks, utils.Check{Name: "لم يتم ذكر انواع الحيوانات", Ok: newFollowUp.AnimalsTypeCount != "", Weight: 3})
			}
			if newFollowUp.UsesBioProducts == "نعم" {
				checks = append(checks, utils.Check{Name: "لم يتم ذكر المنتجات الحيوية المستخدمة", Ok: len(newFollowUp.BiosProductsDetails) > 0, Weight: 3})
			}

			if newFollowUp.VisitDate == "" {
				checks = append(checks, utils.Check{Name: "لا يوجد تاريخ زيارة", Ok: true, Weight: 3})
			} else {
				creation, err := time.Parse(time.DateOnly, strings.Split(newFollowUp.Creation, " ")[0])
				if err != nil {
					return err
				}
				visitDate, err := time.Parse(time.DateOnly, strings.Split(newFollowUp.VisitDate, " ")[0])
				if err != nil {
					return err
				}
				const permissibleTimeDiff = 24 * time.Hour
				if visitDate.Sub(creation) > permissibleTimeDiff {
					checks = append(checks, utils.Check{Name: "تلاعب بتاريخ الزيارة", Ok: false, Weight: 0})
				}
				if creation.Sub(visitDate) > permissibleTimeDiff {
					checks = append(checks, utils.Check{Name: "تلاعب بتاريخ الزيارة", Ok: false, Weight: 0})
				}
			}

			var (
				totalWeight  float64
				filledWeight float64
			)
			for _, c := range checks {
				totalWeight += c.Weight
				if !c.Ok {
					filledWeight += c.Weight
					newFollowUp.Issues = append(newFollowUp.Issues, c.Name)
				}
			}
			newFollowUp.RatePercent = 1 - (filledWeight / totalWeight)

			fmt.Println(newFollowUp.RatePercent)
			fmt.Println(newFollowUp.Issues)

			// followUp, err = frappe.Create()
			// if err != nil {
			// 	return err
			// }
			// state.FormIDOnFrappe = followUp.Name
			// err = config.DB.Save(&state).Error
			// if err != nil {
			// 	return err
			// }
		}

	}

	return nil
}
