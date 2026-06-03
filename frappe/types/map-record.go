package types

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ahmedsat/ebda-cli/frappe"
	"github.com/ahmedsat/ebda-cli/geo"
	"github.com/ahmedsat/ebda-cli/geo/kml"
)

type MapRecord struct {
	Base
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
	// ! i have to do this shit to make it work
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

	if m.Coordinates[0] != m.Coordinates[len(m.Coordinates)-1] {
		m.Coordinates = append(m.Coordinates, m.Coordinates[0])
	}

	area, err := geo.Polygon{
		Ring: m.Coordinates,
	}.SphericalArea()
	if err != nil {
		return fmt.Errorf("%s: %s", m.Name, err)
	}
	m.Area_in_feddan = area / 4200

	m.Parsed = true

	return nil
}

func MapRecordsToKML(records []MapRecord) ([]byte, error) {

	slices.SortFunc(records, func(m1, m2 MapRecord) int { return strings.Compare(m1.Farm, m2.Farm) })

	var placemarks []kml.Placemark

	for i := range records {
		r := &records[i]

		if err := r.Parse(); err != nil {
			return nil, fmt.Errorf("shape %s: %w", r.ShapeID, err)
		}

		if len(r.Coordinates) < 3 {
			continue
		}
		coords := buildKMLCoordinates(r.Coordinates)

		farm, err := frappe.GetCached1[Farm](r.Farm)
		if err != nil {
			return nil, err
		}

		pm := kml.Placemark{
			Name: fmt.Sprintf("%s - %s - %s - %s", farm.ArabicName, farm.Region, farm.FarmId, r.Name),
			Description: fmt.Sprintf(
				"Creator: %s\nFarm: %s\nSeason: %s\nArea: %.2f Fed",
				r.Owner, r.Farm, r.Season, farm.Area,
			),
			Style: kml.Style{
				PolyStyle: kml.PolyStyle{
					Color: kml.KmlColor(r.Color),
				},
			},
			Polygon: kml.Polygon{
				OuterBoundary: kml.OuterBoundary{
					LinearRing: kml.LinearRing{
						Coordinates: coords,
					},
				},
			},
		}

		placemarks = append(placemarks, pm)

		fmt.Fprintf(os.Stderr, "\r%f%%", float64(i+1)/float64(len(records))*100)
	}

	k := kml.KML{
		Xmlns: "http://www.opengis.net/kml/2.2",
		Document: kml.Document{
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

func GetMapColored(farm, color string) (result []MapRecord, err error) {
	result, err = frappe.Get[MapRecord](frappe.Filters{frappe.NewFilter("farm", frappe.Eq, farm)}, nil, nil)
	if err != nil {
		return
	}

	for i := range result {
		err = result[i].Parse()
		if err != nil {
			return
		}
		result[i].Color = color
	}

	return
}
