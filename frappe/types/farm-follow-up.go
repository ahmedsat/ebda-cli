package types

import (
	"strings"
	"time"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/utils"
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
	Base

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

	RatePercent float64  `json:"-"`
	Issues      []string `json:"-"`
	Rated       bool     `json:"-"`
}

func (f FarmFollowUp) DocTypeName() string {
	return "Farm FollowUp"
}

// const sep = "\t"

func (f *FarmFollowUp) Rate() error {
	if f.Rated {
		return nil
	}

	follow, err := frappe.Get1[FarmFollowUp](f.Name)
	if err != nil {
		return err
	}

	*f = follow

	checks := []utils.Check{
		{Name: "لايوجد موقع", Ok: f.GPS != "", Weight: 3},
		{Name: "اسم المتابع غير موجود", Ok: f.FollowerName != "", Weight: 5},
		{Name: "صورة المتابع مع المزارعين غير موجودة", Ok: f.PictureOfFollower != "", Weight: 5},
		// {Name: "عدد المزارعين غير مطابق لاسمائهم", Ok: f.FarmersCount == len(f.FarmersNames), Weight: 0},
		{Name: "لا يوجد محاصيل", Ok: len(f.CurrentCrops) != 0, Weight: 3},
		{Name: "معدل انتاج الكمبوست غير موجود", Ok: f.CompostProduction > 0, Weight: 3},
		{Name: "كمية الكمبوست غير موجودة", Ok: f.CompostQtys > 0, Weight: 3},
		{Name: "لم يتم ذكر التحديات الحالية", Ok: f.CurrentChallenges != "", Weight: 1},
		{Name: "لم يتم ذكر تقييم المراجع", Ok: f.FollowerAssessment != "", Weight: 1},
		{Name: "لم يتم ذكر التوصيات", Ok: f.Recommendations != "", Weight: 4},
		{Name: "لم يتم ذكر هل يوجد مخزن ام لا", Ok: f.StorageExist != "", Weight: 3},
	}

	if f.RecordsFarmBook != 0 {
		checks = append(checks, utils.Check{Name: "صورة دفتر المزرعة غير موجودة", Ok: f.RecordImage != "1", Weight: 3})
	}

	if f.StorageExist == "نعم" {
		checks = append(checks, utils.Check{Name: "لم يتم ذكر محتويات المخزن", Ok: f.WarehousesNotes != "", Weight: 3})
	}

	if f.IntercroppingOrGreenManure == "نعم" {
		checks = append(checks, utils.Check{Name: "لم يتم ذكر نسبة التحميل", Ok: f.IntercroppingPercent > 0, Weight: 3})
	}

	if f.PlantedTreesOrHedge == "نعم" {
		checks = append(checks, utils.Check{Name: "لم يتم ذكر عدد الاشجار الجديدة", Ok: f.TreesCount > 0, Weight: 3})
	}

	if f.HasAnimals == "نعم" {
		checks = append(checks, utils.Check{Name: "لم يتم ذكر انواع الحيوانات", Ok: f.AnimalsTypeCount != "", Weight: 3})
	}

	// BiosProductsDetails
	if f.UsesBioProducts == "نعم" {
		checks = append(checks, utils.Check{Name: "لم يتم ذكر المنتجات الحيوية المستخدمة", Ok: len(f.BiosProductsDetails) > 0, Weight: 3})
	}

	if f.VisitDate == "" {
		checks = append(checks, utils.Check{Name: "لا يوجد تاريخ زيارة", Ok: true, Weight: 3})
	} else {
		creation, err := time.Parse(frappe.TimeLayout, strings.Split(f.Creation, " ")[0])
		if err != nil {
			return err
		}

		visitDate, err := time.Parse(frappe.TimeLayout, strings.Split(f.VisitDate, " ")[0])
		if err != nil {
			return err
		}

		const permissibleTimeDiff = 24 * time.Hour

		// time diff
		if visitDate.Sub(creation) > permissibleTimeDiff {
			checks = append(checks, utils.Check{Name: "تلاعب بتاريخ الزيارة", Ok: false, Weight: 0})
		}

		// time diff
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
			f.Issues = append(f.Issues, c.Name)
		}
	}

	f.RatePercent = 1 - (filledWeight / totalWeight)

	f.Rated = true
	return nil
}
