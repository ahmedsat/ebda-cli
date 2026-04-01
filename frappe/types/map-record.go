package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/geo"
	"github.com/ahmedsat/ebda-cli/utils"
)

type MapRecord struct {
	Base
	Name             string      `json:"name"`
	ShapeID          string      `json:"shape_id"`
	Type             string      `json:"type"`
	NameOfShape      string      `json:"name_of_shape"`
	Farm             string      `json:"farm"`
	Farm_application string      `json:"farm_application"`
	Season           string      `json:"season"`
	Posting_date     string      `json:"posting_date"`
	Area_in_feddan   float64     `json:"area_in_feddan"`
	Area_in_hectar   float64     `json:"area_in_hectar"`
	Color            string      `json:"color"`
	Jsoncode         string      `json:"jsoncode"`
	Coordinates      []geo.Coord `json:"-"`
	Parsed           bool        `json:"-"`
}

func (m MapRecord) DocTypeName() string { return "Map Records" }

func (m *MapRecord) Parse() error {
	if m.Parsed {
		return nil
	}

	// ! the shitty implementation of this field makes it some time quoted despite being a float
	// ! is have to do this shit to make it work
	m.Jsoncode = strings.ReplaceAll(m.Jsoncode, "\":\"", "\":")
	m.Jsoncode = strings.ReplaceAll(m.Jsoncode, "\"}", "}")
	m.Jsoncode = strings.ReplaceAll(m.Jsoncode, "\",\"", ",\"")

	m.Jsoncode = strings.TrimSpace(m.Jsoncode)
	if !strings.HasPrefix(m.Jsoncode, "[") {
		m.Jsoncode = "[" + m.Jsoncode + "]"
	}
	if err := json.Unmarshal([]byte(m.Jsoncode), &m.Coordinates); err != nil {
		return err
	}

	// if m.Farm != "" {
	// 	m.Farm = strings.TrimSpace(m.Farm)

	// 	f, err := frappe.Get1[Farm](m.Farm)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	m.Farm = f.Name

	// 	m.Name = fmt.Sprintf("%s - %s - %s", f.ArabicName, f.Region, f.FarmId)
	// }

	area, err := geo.Polygon{
		Coords: m.Coordinates,
	}.GeodesicArea()
	if err != nil {
		return err
	}
	m.Area_in_feddan = area / 4200

	m.Parsed = true

	return nil
}

func RecordsToKML(records []MapRecord) ([]byte, error) {
	var placemarks []utils.Placemark

	for i := range records {
		r := &records[i]

		if err := r.Parse(); err != nil {
			return nil, fmt.Errorf("shape %s: %w", r.ShapeID, err)
		}

		if len(r.Coordinates) < 3 {
			continue
		}
		coords := buildKMLCoordinates(r.Coordinates)

		farm, err := frappe.Get1[Farm](r.Farm)
		if err != nil {
			return nil, err
		}

		pm := utils.Placemark{
			Name: fmt.Sprintf("%s - %s - %s", farm.ArabicName, farm.Region, farm.FarmId),
			Description: fmt.Sprintf(
				"Farm: %s\nSeason: %s\nArea: %.2f Ha",
				r.Farm, r.Season, farm.Area*4200/10000,
			),
			Style: utils.Style{
				PolyStyle: utils.PolyStyle{
					Color: utils.KmlColor(r.Color),
				},
			},
			Polygon: utils.Polygon{
				OuterBoundary: utils.OuterBoundary{
					LinearRing: utils.LinearRing{
						Coordinates: coords,
					},
				},
			},
		}

		placemarks = append(placemarks, pm)

		fmt.Fprintf(os.Stderr, "\r%f%%", float64(i+1)/float64(len(records))*100)
	}

	k := utils.KML{
		Xmlns: "http://www.opengis.net/kml/2.2",
		Document: utils.Document{
			Name:       "Map Records",
			Placemarks: placemarks,
		},
	}

	// progress
	return xml.MarshalIndent(k, "", "  ")
}

func buildKMLCoordinates(coords []geo.Coord) string {
	var b strings.Builder

	for _, c := range coords {
		// KML = lng,lat
		fmt.Fprintf(&b, "%f,%f ", c.Lng, c.Lat)
	}

	// Close the ring
	first := coords[0]
	fmt.Fprintf(&b, "%f,%f", first.Lng, first.Lat)

	return strings.TrimSpace(b.String())
}
