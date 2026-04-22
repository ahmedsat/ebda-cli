package types

type HarvestTable struct {
	Name string `json:"name"`
}

type HarvestRecord struct {
	Name         string         `json:"name"`
	Season       string         `json:"season"`
	Farm         string         `json:"farm"`
	HarvestTable []HarvestTable `json:"harvest_table"`
}

func (hr HarvestRecord) DocTypeName() string {
	return "Harvest Record"
}
