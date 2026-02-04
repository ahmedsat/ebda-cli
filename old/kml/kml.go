package kml

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedsat/erp-reports-cli/types"
)

type KML struct {
	XMLName  xml.Name `xml:"kml"`
	Xmlns    string   `xml:"xmlns,attr"`
	Document Document `xml:"Document"`
}

type Document struct {
	Name       string      `xml:"name,omitempty"`
	Placemarks []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name        string   `xml:"name,omitempty"`
	Description string   `xml:"description,omitempty"`
	Style       *Style   `xml:"Style,omitempty"`
	Polygon     *Polygon `xml:"Polygon,omitempty"`
}

type Style struct {
	PolyStyle PolyStyle `xml:"PolyStyle"`
}

type PolyStyle struct {
	Color string `xml:"color,omitempty"`
}

type Polygon struct {
	OuterBoundary OuterBoundary `xml:"outerBoundaryIs"`
}

type OuterBoundary struct {
	LinearRing LinearRing `xml:"LinearRing"`
}

type LinearRing struct {
	Coordinates string `xml:"coordinates"`
}

func RecordsToKML(records []types.MapRecord) ([]byte, error) {
	var placemarks []Placemark

	for i := range records {
		r := &records[i]

		if err := r.Parse(); err != nil {
			return nil, fmt.Errorf("shape %s: %w", r.ShapeID, err)
		}

		if len(r.Coordinates) < 3 {
			continue
		}

		coords := buildKMLCoordinates(r.Coordinates)

		pm := Placemark{
			Name: r.Name,
			Description: fmt.Sprintf(
				"Farm: %s\nSeason: %s\nArea (feddan): %.2f",
				r.Farm, r.Season, r.Area_in_feddan,
			),
			Style: &Style{
				PolyStyle: PolyStyle{
					Color: kmlColor(r.Color),
				},
			},
			Polygon: &Polygon{
				OuterBoundary: OuterBoundary{
					LinearRing: LinearRing{
						Coordinates: coords,
					},
				},
			},
		}

		placemarks = append(placemarks, pm)

		fmt.Fprintf(os.Stderr, "\r%f%%", float64(i+1)/float64(len(records))*100)
	}

	k := KML{
		Xmlns: "http://www.opengis.net/kml/2.2",
		Document: Document{
			Name:       "Map Records",
			Placemarks: placemarks,
		},
	}

	// progress
	return xml.MarshalIndent(k, "", "  ")
}

func buildKMLCoordinates(coords []types.Coord) string {
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

func kmlColor(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "7d00ff00" // default semi-transparent green
	}

	rr := hex[0:2]
	gg := hex[2:4]
	bb := hex[4:6]

	return "7d" + bb + gg + rr
}
