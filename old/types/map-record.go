package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ahmedsat/erp-reports-cli/erp"
)

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type MapRecord struct {
	Name             string  `json:"-"`
	ShapeID          string  `json:"shape_id"`
	Type             string  `json:"type"`
	NameOfShape      string  `json:"name_of_shape"`
	Farm             string  `json:"farm"`
	Farm_application string  `json:"farm_application"`
	Season           string  `json:"season"`
	Posting_date     string  `json:"posting_date"`
	Area_in_feddan   float64 `json:"area_in_feddan"`
	Area_in_hectar   float64 `json:"area_in_hectar"`
	Color            string  `json:"color"`
	Jsoncode         string  `json:"jsoncode"`
	Coordinates      []Coord `json:"-"`
	Parsed           bool    `json:"-"`
}

func (m MapRecord) DocTypeName() string { return "Map Records" }
func (m *MapRecord) Parse() error {
	if m.Parsed {
		return nil
	}

	m.Jsoncode = strings.TrimSpace(m.Jsoncode)
	if !strings.HasPrefix(m.Jsoncode, "[") {
		m.Jsoncode = "[" + m.Jsoncode + "]"
	}
	if err := json.Unmarshal([]byte(m.Jsoncode), &m.Coordinates); err != nil {
		return err
	}

	if m.Farm != "" {
		m.Farm = strings.TrimSpace(m.Farm)

		f, err := erp.Get1[Farm](m.Farm)
		if err != nil {
			return err
		}
		m.Farm = f.Name

		m.Name = fmt.Sprintf("%s - %s - %s", f.ArabicName, f.Region, f.Code)
	}

	m.Parsed = true

	return nil
}
