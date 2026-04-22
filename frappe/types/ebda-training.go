package types

type TrainingFarmer struct {
	BaseInnerTable
	FarmerName string `json:"farmer_name"`
}

type EbdaTraining struct {
	Base
	Farm       string           `json:"farm"`
	Topic      string           `json:"ebda_trainingname"`
	Status     string           `json:"ebda_trainingstatus"`
	FarmID     string           `json:"farmid"`
	ArabicName string           `json:"arabic_name"`
	Trainer    string           `json:"ebda_trainer"`
	Type       string           `json:"ebda_trainingtype"`
	Region     string           `json:"region"`
	Date       string           `json:"trainingdate"`
	Intro      string           `json:"intro"`
	Farmers    []TrainingFarmer `json:"farmers"`
}

func (e EbdaTraining) DocTypeName() string {
	return "Ebda Training"
}
