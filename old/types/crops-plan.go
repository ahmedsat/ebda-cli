package types

// Crops Plan Total Table
type CropsPlanTotalTable struct {
	Name string  `json:"name"`
	Crop string  `json:"crop"`
	Area float64 `json:"area_in_feddan"`
}

type CropPlan struct {
	ID         string                `json:"name"`
	Season     string                `json:"season"`
	Farm       string                `json:"farm"`
	TotalTable []CropsPlanTotalTable `json:"merged_crops_plan"`
}

func (CropPlan) DocTypeName() string {
	return "Crops Plan"
}
