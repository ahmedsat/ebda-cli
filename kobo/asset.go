package kobo

type AssetsList struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
}

func (AssetsList) GetFormID() string { panic("assets list is not a typical asset") }
