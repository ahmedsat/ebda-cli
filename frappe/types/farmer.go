package types

type Farmer struct {
	Base
	FarmerName string `json:"farmer_name"`
	Phone      string `json:"phone"`
	Gender     string `json:"gender"`
}

func (f Farmer) DocTypeName() string {
	return "Farmer"
}
