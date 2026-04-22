package commands

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/utils"
)

type BioProductFollowUp struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Creation   string `json:"creation"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modified_by"`
	DocStatus  int    `json:"docstatus"`
	Idx        int    `json:"idx"`

	Fertilizers string `json:"fertilizers"`

	Parent      string `json:"parent"`
	ParentField string `json:"parentfield"`
	ParentType  string `json:"parenttype"`
	Doctype     string `json:"doctype"`
}

type CropFollowUp struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Creation   string `json:"creation"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modified_by"`
	DocStatus  int    `json:"docstatus"`
	Idx        int    `json:"idx"`

	Crops string `json:"crops"`

	Parent      string `json:"parent"`
	ParentField string `json:"parentfield"`
	ParentType  string `json:"parenttype"`
	Doctype     string `json:"doctype"`
}

type FarmerFollowUp struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Creation   string `json:"creation"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modified_by"`
	DocStatus  int    `json:"docstatus"`
	Idx        int    `json:"idx"`

	Farmer string `json:"farmer"`
	Status string `json:"status"`

	Parent      string `json:"parent"`
	ParentField string `json:"parentfield"`
	ParentType  string `json:"parenttype"`
	Doctype     string `json:"doctype"`
}

type FarmFollowUp struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Creation   string `json:"creation"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modified_by"`
	DocStatus  int    `json:"docstatus"`
	Idx        int    `json:"idx"`

	Farm                string  `json:"farm"`
	FarmName            string  `json:"farm_name"`
	FarmCode            string  `json:"farm_code"`
	Region              string  `json:"region"`
	AreaFeddan          float64 `json:"area_feddan"`
	Latitude            string  `json:"latitude"`
	Longitude           string  `json:"longitude"`
	UpscalingProject    string  `json:"upscaling_project"`
	SubUpscalingProject string  `json:"subupscalingproject"`
	FarmGroup           string  `json:"farm_group"`

	FarmOwner              string `json:"farm_owner"`
	Phone                  string `json:"phone"`
	PictureOfFollower      string `json:"picture_of_follower"`
	WarehousesNotes        string `json:"warehouses_notes"`
	RecordImage            string `json:"record_image"`
	AnimalsTypeCount       string `json:"animals_type_count"`
	UsedQuantities         string `json:"الكميات_المستخدمة"`
	BioControlDetails      string `json:"bio_control_details"`
	ChemicalControlDetails string `json:"chemical_control_details"`
	NaturalEnemiesDetails  string `json:"natural_enemies_details"`
	WeedDisposalOther      string `json:"weed_disposal_other"`
	IrrigationMethodOther  string `json:"irrigation_method_other"`
	WaterSourceOther       string `json:"water_source_other"`
	EnergyDetails          string `json:"energy_details"`
	SoilAnalysisDetails    string `json:"soil_analysis_details"`
	CurrentChallenges      string `json:"current_challenges"`
	SupportNeededOther     string `json:"support_needed_other"`
	FollowerAssessment     string `json:"follower_assessment"`
	Recommendations        string `json:"follower_recommendations"`

	GPS       string `json:"gps"`
	VisitDate string `json:"visit_date"`

	FollowerType  string `json:"follower_type"`
	FollowerName  string `json:"follower_name"`
	FarmersCount  int    `json:"farmers_count"`
	LastEBDAVisit string `json:"last_ebda_visit"`

	HasHealthInsuranceCard int    `json:"has_health_insurance_card"`
	StorageExist           string `json:"storage_exist"`
	RecordsFarmBook        int    `json:"records_farm_book"`

	IntercroppingOrGreenManure string  `json:"intercropping_or_green_manure"`
	IntercroppingPercent       float64 `json:"intercropping_percent"`

	PlantedTreesOrHedge string `json:"planted_trees_or_hedge"`
	TreesCount          int    `json:"trees_count"`
	HasAnimals          string `json:"has_animals"`

	FertilizationSource string  `json:"fertilization_source"`
	CompostSource       string  `json:"compost_source"`
	CompostProduction   float64 `json:"compost_production"`
	CompostQtys         float64 `json:"compost_qtys"`

	UsesBioProducts string  `json:"uses_bio_products"`
	CompQty         float64 `json:"comp_qty"`
	Qurts           float64 `json:"qurts"`
	Horns           float64 `json:"horns"`

	PestsOrDiseasesLastSeason string `json:"pests_or_diseases_last_season"`
	PestControlMethod         string `json:"pest_control_method"`
	UsesNaturalEnemies        int    `json:"uses_natural_enemies"`

	WeedDisposal            string `json:"weed_disposal"`
	IrrigationMethod        string `json:"irrigation_method"`
	WaterSource             string `json:"water_source"`
	WaterShortageOrSalinity int    `json:"water_shortage_or_salinity"`

	EnergyType                  string `json:"energy_type"`
	ClimateChallengesThisSeason int    `json:"climate_challenges_this_season"`
	YieldVsLastSeason           string `json:"yield_vs_last_season"`

	NeedCertSupport int    `json:"need_cert_support"`
	SupportNeeded   string `json:"support_needed"`

	Doctype      string `json:"doctype"`
	ServicesUsed []any  `json:"services_used"`

	FarmersNames        []FarmerFollowUp     `json:"farmers_names"`
	BiosProductsDetails []BioProductFollowUp `json:"bios_products_details"`
	CurrentCrops        []CropFollowUp       `json:"curent_crops"`

	rate  float64 `json:"-"`
	rated bool    `json:"-"`
	issue string  `json:"-"`
}

func (f FarmFollowUp) DocTypeName() string {
	return "Farm FollowUp"
}

const sep = "\t"

func (f *FarmFollowUp) Rate() error {
	if f.rated {
		return nil
	}

	follow, err := erp.Get1[FarmFollowUp](f.Name)
	if err != nil {
		return err
	}

	*f = follow

	type check struct {
		name   string
		ok     bool
		weight float64
	}

	checks := []check{
		{"لايوجد موقع", f.GPS != "", 3},
		{"لا يوجد تاريخ زيارة", f.VisitDate != "", 3},
		{"اسم المتابع غير موجود", f.FollowerName != "", 5},
		{"صورة المتابع مع المزارعين غير موجودة", f.PictureOfFollower != "", 5},
		{"عدد المزارعين غير مطابق لاسمائهم", f.FarmersCount == len(f.FarmersNames), 3},
		{"لا يوجد محاصيل", len(f.CurrentCrops) != 0, 3},
		{"معدل انتاج الكمبوست غير موجود", f.CompostProduction > 0, 3},
		{"كمية الكمبوست غير موجودة", f.CompostQtys > 0, 3},
		{"لم يتم ذكر التحديات الحالية", f.CurrentChallenges != "", 1},
		{"لم يتم ذكر تقييم المراجع", f.FollowerAssessment != "", 1},
		{"لم يتم ذكر التوصيات", f.Recommendations != "", 4},
		{"لم يتم ذكر هل يوجد مخزن ام لا", f.StorageExist != "", 3},
	}

	if f.RecordsFarmBook != 0 {
		checks = append(checks, check{"صورة دفتر المزرعة غير موجودة", f.RecordImage != "1", 3})
	}

	if f.StorageExist == "نعم" {
		checks = append(checks, check{"لم يتم ذكر محتويات المخزن", f.WarehousesNotes != "", 3})
	}

	if f.IntercroppingOrGreenManure == "نعم" {
		checks = append(checks, check{"لم يتم ذكر نسبة التحميل", f.IntercroppingPercent > 0, 3})
	}

	if f.PlantedTreesOrHedge == "نعم" {
		checks = append(checks, check{"لم يتم ذكر عدد الاشجار الجديدة", f.TreesCount > 0, 3})
	}

	if f.HasAnimals == "نعم" {
		checks = append(checks, check{"لم يتم ذكر انواع الحيوانات", f.AnimalsTypeCount != "", 3})
	}

	// BiosProductsDetails
	if f.UsesBioProducts == "نعم" {
		checks = append(checks, check{"لم يتم ذكر المنتجات الحيوية المستخدمة", len(f.BiosProductsDetails) > 0, 3})
	}

	var (
		totalWeight   float64
		filledWeight  float64
		missingFields []string
	)

	for _, c := range checks {
		totalWeight += c.weight
		if !c.ok {
			filledWeight += c.weight
			missingFields = append(missingFields, c.name)
		}
	}

	// Failure rate: 0.0 (perfect) → 1.0 (total failure)
	f.rate = 1 - (filledWeight / totalWeight)

	switch {
	case f.rate == 1:
		f.issue = "كامل"
	case f.rate > 0.75:
		f.issue = "نقص قليل" + sep + strings.Join(missingFields, ", ")
	case f.rate > 0.5:
		f.issue = "ناقص" + sep + strings.Join(missingFields, ", ")
	default:
		f.issue = "غير مكتمل" + sep + strings.Join(missingFields, ", ")
	}

	f.rated = true
	return nil
}

func FollowUp(args []string) error {

	results, err := erp.Get[FarmFollowUp](nil, utils.List{"name"})
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "Calculating rates...")
	counter := 1
	s := utils.NewSyncRunner(10, 0)
	for i := range results {
		s.Run(func() {
			err := results[i].Rate()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Fprintf(os.Stderr, "\r%d/%d (%.2f%%)", counter, len(results), float64(counter)/float64(len(results))*100)
			counter++
		})
	}
	s.Wait()
	fmt.Fprintln(os.Stderr)

	fmt.Fprintln(os.Stderr, "Sorting results...")
	slices.SortFunc(results, func(f1, f2 FarmFollowUp) int {
		return int(f2.rate*100) - int(f1.rate*100)
	})

	fmt.Fprintln(os.Stderr, "Printing results...")
	for _, result := range results {
		if !result.rated {
			continue
		}
		fmt.Println(strings.Join([]string{result.Name, result.FarmCode, fmt.Sprintf("%f", result.rate), result.issue}, sep))
	}

	return nil
}
