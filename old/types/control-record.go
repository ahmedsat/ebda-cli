package types

type ControlTable struct {
	Name string `json:"name"`
}

type ControlRecord struct {
	Name         string         `json:"name"`
	Season       string         `json:"season"`
	Farm         string         `json:"farm"`
	ControlTable []ControlTable `json:"control_table"`
}

func (ControlRecord) DocTypeName() string { return "Control Record" }
