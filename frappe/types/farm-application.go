package types

type FarmApplicationFarmer struct {
	BaseInnerTable
	Farmer             string  `json:"farmer"`
	OwnedAreaInFeddan  float64 `json:"owned_area_in_feddan"`
	Gender             string  `json:"gender"`
	FarmerPhoto        string  `json:"farmer_photo"`
	FarmerNationalId   string  `json:"farmer_national_id"`
	FarmerPhoneNumber  int     `json:"farmer_phone_number"`
	FarmerNationalIdNo string  `json:"farmer_national_id_no"`
}

type FarmApplication struct {
	Base
	Farmers      []FarmApplicationFarmer `json:"farmers"`
	TotalFarmers int                     `json:"total_farmers"`
}

func (f FarmApplication) DocTypeName() string {
	return "Farm Application"
}
