package types

type SoilAnalysis struct {
	Base

	Farm                 string `json:"farm"`
	FarmID               string `json:"custom_farm_id"`
	Location             string `json:"location"`
	CollectionDatetime   string `json:"collection_datetime"`
	NamingSeries         string `json:"naming_series"`
	SoilAnalysisCriteria []any  `json:"soil_analysis_criteria"`
}

func (s SoilAnalysis) DocTypeName() string {
	return "Soil Analysis"
}
