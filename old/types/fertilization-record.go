package types

type MergedFertilization struct {
	Name       string  `json:"name"`       // "name": "topmcf45vf",
	Fertilizer string  `json:"fertilizer"` // "fertilizer": "Compost",
	Quantity   float64 `json:"qty"`        // "qty": 472.0,
	UOM        string  `json:"uom"`        // "uom": "Ton",
}

type FertilizationRecord struct {
	Name                 string                `json:"name"`
	Season               string                `json:"season"`
	Farm                 string                `json:"farm"`
	Receipt              string                `json:"receipt"`
	MergedFertilizations []MergedFertilization `json:"merged_fertilization_table"`
}

func (FertilizationRecord) DocTypeName() string { return "Fertilization Record" }
