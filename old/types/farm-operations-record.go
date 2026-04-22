package types

type FarmOperationsTable struct {
	Name string `json:"name"`
}

type FarmOperationsRecord struct {
	Name                string                `json:"name"`
	Season              string                `json:"season"`
	Farm                string                `json:"farm"`
	FarmOperationsTable []FarmOperationsTable `json:"operations_table"`
}

func (FarmOperationsRecord) DocTypeName() string { return "Farm Operations Record" }
