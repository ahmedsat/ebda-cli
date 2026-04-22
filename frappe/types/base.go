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

func (b Base) DocName() string {
	return b.Name
}

type Unknown map[string]any

type BaseInnerTable struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	Creation    string `json:"creation"`
	Modified    string `json:"modified"`
	ModifiedBy  string `json:"modified_by"`
	DocStatus   int    `json:"docstatus"`
	Idx         int    `json:"idx"`
	Parent      string `json:"parent"`
	ParentField string `json:"parentfield"`
	ParentType  string `json:"parenttype"`
	Doctype     string `json:"doctype"`
}
