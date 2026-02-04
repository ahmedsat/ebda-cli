package types

type Farm struct {
	Name       string `json:"name"`
	ArabicName string `json:"arabic_name"`
	Region     string `json:"region"`
	Code       string `json:"farm_id"`
}

func (f Farm) DocTypeName() string {
	return "Farm"
}
