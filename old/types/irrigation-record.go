package types

type IrrigationTable struct {
	Name string `json:"name"`
}

type IrrigationRecord struct {
	Name            string            `json:"name"`
	Season          string            `json:"season"`
	Farm            string            `json:"farm"`
	IrrigationTable []IrrigationTable `json:"irrigation_table"`
}

func (IrrigationRecord) DocTypeName() string { return "Irrigation Record" }
