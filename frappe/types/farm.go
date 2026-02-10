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
	DocName_              string          `json:"docname"`
	Type                  string          `json:"type"`
	IsPlotOfSector        int             `json:"is_plot_of_sector"`
	ArabicName            string          `json:"arabic_name"`
	UpscalingProject      string          `json:"upscaling_project"`
	Subupscalingproject   string          `json:"subupscalingproject"`
	Region                string          `json:"region"`
	ParentFarm            string          `json:"parent_farm"`
	Category              string          `json:"gategory"`
	LeadingEngineers      int             `json:"leading_engineers"`
	FarmName              string          `json:"farm_name"`
	Area                  float64         `json:"farm_area__feddan"`
	IsInternalFarm        int             `json:"is_internal_farm"`
	Company               string          `json:"company"`
	FarmGroup             string          `json:"farm_group"`
	FarmId                string          `json:"farm_id"`
	FarmOwner             string          `json:"farm_owner"`
	FarmOwnership         string          `json:"farm_ownership"`
	Phone                 string          `json:"phone"`
	FarmOwnershipDocument string          `json:"farm_ownership_document"`
	TotalFarmers          int             `json:"total_farmers"`
	Latitude              string          `json:"latitude"`
	Longitude             string          `json:"longitude"`
	FarmApplicationName   string          `json:"farm_application"`
	FarmApplication       FarmApplication `json:"-"`
	YearOfReclamation     int             `json:"year_of_reclamation"`
	SoilStatus            string          `json:"soil_status"`
	TotalTrees            int             `json:"total_trees"`
	Lft                   int             `json:"lft"`
	Rgt                   int             `json:"rgt"`
	OldParent             string          `json:"old_parent"`
	SectorAreaInFeddan    float64         `json:"sector_area_in_feddan"`
	PlotAreaInFeddan      float64         `json:"plot_area_in_feddan"`
	SubPlotAreaInFeddan   float64         `json:"sub_plot_area_in_feddan"`
	IsGroup               int             `json:"is_group"`
	EnglishName           string          `json:"en_name"`
	CreationDate          string          `json:"creation_date"`
	Table28               []Unknown       `json:"table_28"`
	Certificates          []Unknown       `json:"certificates"`
	Workers               []Unknown       `json:"workers"`
	Table24               []Unknown       `json:"table_24"`
	Table26               []Unknown       `json:"table_26"`
	Tree                  []Unknown       `json:"tree"`
	Farmers               []FarmFarmer    `json:"farmers"`
	UserTags              string          `json:"_user_tags"`
	FarmCertificates      string          `json:"farm_certificates"`
	FarmStatus            string          `json:"farm_status"`
}

func (f Farm) DocTypeName() string {
	return "Farm"
}

func (f Farm) Update() {
	frappe.UpdateDoc(f)
}
