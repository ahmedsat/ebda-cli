package kobo

type ValidationStatus struct {
	Label string `json:"label"`
}

type PGSNew struct {
	FormID           string `json:"at_house/farm_id"`
	EngName          string `json:"at_house/__011"`
	VisitDate        string `json:"at_house/visit_date"`
	ValidationStatus `json:"_validation_status"`
}

func (pgs PGSNew) GetFormID() string { return "aX4NJWgge6tooXjfSYXhrq" }
