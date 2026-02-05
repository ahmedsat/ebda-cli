package types

type Base struct {
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Creation   string `json:"creation"`
	Modified   string `json:"modified"`
	ModifiedBy string `json:"modified_by"`
	DocStatus  int    `json:"docstatus"`
	Idx        int    `json:"idx"`
	Doctype    string `json:"doctype"`
}
