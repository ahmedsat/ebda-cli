package types

import (
	"errors"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
)

type FarmFarmer struct {
	BaseInnerTable
	FarmerName            string  `json:"farmer_name"`
	TotalArea             float64 `json:"total_area"`
	AreaUnit              string  `json:"area_unit"`
	Visa                  string  `json:"visa"`
	FarmerNationalIdImage string  `json:"farmer_national_id"`
	NationalIdNumber      string  `json:"national_id"`
	Phone                 string  `json:"phone"`
}

type FarmWorker struct {
	BaseInnerTable
	Count  int    `json:"worker"`
	Gender string `json:"gender"`
	Age    string `json:"age"`
}

type Farm struct {
	Base
	Name__              string       `json:"docname"`
	Farmers             []FarmFarmer `json:"farmers"`
	FarmId              string       `json:"farm_id"`
	ArabicName          string       `json:"arabic_name"`
	Region              string       `json:"region"`
	TotalFarmers        int          `json:"total_farmers"`
	FarmOwner           string       `json:"farm_owner"`
	Area                float64      `json:"farm_area__feddan"`
	CreationDate        string       `json:"creation_date"`
	FarmApplicationID   string       `json:"farm_application"`
	Type                string       `json:"type"`
	IsPlotOfSector      byte         `json:"is_plot_of_sector"`
	UpscalingProject    string       `json:"upscaling_project"`
	Subupscalingproject string       `json:"subupscalingproject"`
	ParentFarm          string       `json:"parent_farm"`
	Category            string       `json:"gategory"`
	LeadingEngineers    byte         `json:"leading_engineers"`
	FarmName            string       `json:"farm_name"`
	IsInternalFarm      byte         `json:"is_internal_farm"`
	Latitude            string       `json:"latitude"`
	Longitude           string       `json:"longitude"`
	FarmStatus          string       `json:"farm_status"`
	Workers             []FarmWorker `json:"workers"`
}

func (f Farm) DocTypeName() string {
	return "Farm"
}

func (f *Farm) Update() (err error) {
	f, err = frappe.UpdateDoc(f)
	return
}

func GetFarmByCode(code string) (f Farm, err error) {

	if code == "" {
		return
	}

	if !strings.HasPrefix(code, "EG/") {
		code = "EG/" + code
	}

	farms, err := frappe.Get[Farm](frappe.Filters{frappe.NewFilter("farm_id", frappe.Eq, code)}, nil, nil)
	if err != nil {
		return
	}
	if len(farms) > 1 {
		err = errors.New("more than one farm found")
	}

	if len(farms) < 1 {
		err = errors.New("no farm found")
	}

	f = farms[0]

	return
}
