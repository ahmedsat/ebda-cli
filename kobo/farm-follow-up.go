package kobo

type FarmFollowupForm struct {
	Common
	FarmCode                    string `json:"farm_code"`                      // select_one_from_file
	Farm                        string `json:"farm"`                           // text
	Region                      string `json:"region"`                         // text
	AreaFeddan                  string `json:"area_feddan"`                    // decimal
	FarmOwner                   string `json:"farm_owner"`                     // text
	Phone                       string `json:"phone"`                          // text
	Latitude                    string `json:"latitude"`                       // text
	Longitude                   string `json:"longitude"`                      // text
	UpscalingProject            string `json:"upscaling_project"`              // text
	Subupscalingproject         string `json:"subupscalingproject"`            // text
	FarmGroup                   string `json:"farm_group"`                     // text
	Gps                         string `json:"gps"`                            // geoshape
	VisitDate                   string `json:"visit_date"`                     // date
	FollowerType                string `json:"follower_type"`                  // select_one
	FollowerName                string `json:"follower_name"`                  // select_one_from_file
	PictureOfFollower           string `json:"picture_of_follower"`            // image
	FarmersCount                string `json:"farmers_count"`                  // integer
	LastEbdaVisit               string `json:"last_ebda_visit"`                // select_one
	HasHealthInsuranceCard      string `json:"has_health_insurance_card"`      // select_one
	StorageExist                string `json:"storage_exist"`                  // select_one
	WarehousesNotes             string `json:"warehouses_notes"`               // text
	RecordsFarmBook             string `json:"records_farm_book"`              // select_one
	RecordImage                 string `json:"record_image"`                   // image
	IntercroppingOrGreenManure  string `json:"intercropping_or_green_manure"`  // select_one
	IntercroppingPercent        string `json:"intercropping_percent"`          // decimal
	PlantedTreesOrHedge         string `json:"planted_trees_or_hedge"`         // select_one
	TreesCount                  string `json:"trees_count"`                    // integer
	HasAnimals                  string `json:"has_animals"`                    // select_one
	AnimalsTypeCount            string `json:"animals_type_count"`             // text
	FertilizationSource         string `json:"fertilization_source"`           // select_one
	CompostSource               string `json:"compost_source"`                 // select_one
	CompostProduction           string `json:"compost_production"`             // decimal
	CompostQtys                 string `json:"compost_qtys"`                   // decimal
	UsesBioProducts             string `json:"uses_bio_products"`              // select_one
	CompQty                     string `json:"comp_qty"`                       // decimal
	Qurts                       string `json:"qurts"`                          // decimal
	Horns                       string `json:"horns"`                          // decimal
	QuantityUsed                string `json:"quantity_used"`                  // text
	PestsOrDiseasesLastSeason   string `json:"pests_or_diseases_last_season"`  // select_one
	PestControlMethod           string `json:"pest_control_method"`            // select_one
	BioControlDetails           string `json:"bio_control_details"`            // text
	ChemicalControlDetails      string `json:"chemical_control_details"`       // text
	UsesNaturalEnemies          string `json:"uses_natural_enemies"`           // select_one
	NaturalEnemiesDetails       string `json:"natural_enemies_details"`        // text
	WeedDisposal                string `json:"weed_disposal"`                  // select_one
	WeedDisposalOther           string `json:"weed_disposal_other"`            // text
	IrrigationMethod            string `json:"irrigation_method"`              // select_one
	IrrigationMethodOther       string `json:"irrigation_method_other"`        // text
	WaterSource                 string `json:"water_source"`                   // select_one
	WaterSourceOther            string `json:"water_source_other"`             // text
	WaterShortageOrSalinity     string `json:"water_shortage_or_salinity"`     // select_one
	EnergyType                  string `json:"energy_type"`                    // select_one
	EnergyDetails               string `json:"energy_details"`                 // text
	ClimateChallengesThisSeason string `json:"climate_challenges_this_season"` // select_one
	YieldVsLastSeason           string `json:"yield_vs_last_season"`           // select_one
	CurrentChallenges           string `json:"current_challenges"`             // text
	SupportNeeded               string `json:"support_needed"`                 // select_one
	SupportNeededOther          string `json:"support_needed_other"`           // text
	FollowerAssessment          string `json:"follower_assessment"`            // text
	FollowerRecommendations     string `json:"follower_recommendations"`       // text
}

func (s FarmFollowupForm) GetFormID() string { return "aJFgvX6ZBgEARrGLqAZeko" }
