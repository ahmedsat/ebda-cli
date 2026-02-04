package types

type SeedsTable struct {
	Name string `json:"name"`
}

type SowingRecord struct {
	Name   string `json:"name"`
	Season string `json:"season"`
	Farm   string `json:"farm"`

	SeedsTable []SeedsTable `json:"sowing_table"`
}

func (s SowingRecord) DocTypeName() string {
	return "Sowing Record"
}
