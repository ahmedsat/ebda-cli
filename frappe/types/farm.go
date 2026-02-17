package types

import "github.com/ahmedsat/ebda-cli/frappe"

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

type Farm struct {
	Base
	Farmers         []FarmFarmer `json:"farmers"`
	FarmId          string       `json:"farm_id"`
	ArabicName      string       `json:"arabic_name"`
	Region          string       `json:"region"`
	TotalFarmers    int          `json:"total_farmers"`
	FarmOwner       string       `json:"farm_owner"`
	Area            float64      `json:"farm_area__feddan"`
	CreationDate    string       `json:"creation_date"`
	FarmApplication string       `json:"farm_application"`
}

func (f Farm) DocTypeName() string {
	return "Farm"
}

func (f *Farm) Update() (err error) {
	f, err = frappe.UpdateDoc(f)
	return
}
