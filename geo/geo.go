package geo

import (
	"fmt"
	"math"
)

const EarthRadius = 6378137.0 // meters (WGS84)

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Triangle struct {
	A Coord
	B Coord
	C Coord
}

type Polygon struct {
	Coords []Coord
}

func toRad(d float64) float64 {
	return d * math.Pi / 180
}

func (p Polygon) GeodesicArea() (float64, error) {
	if len(p.Coords) < 3 {
		return 0, fmt.Errorf("polygon must have at least 3 points")
	}

	// Ensure closed polygon first
	if p.Coords[0] != p.Coords[len(p.Coords)-1] {
		p.Coords = append(p.Coords, p.Coords[0])
	}

	n := len(p.Coords)

	var area float64

	for i := 0; i < n-1; i++ {
		lat1 := toRad(p.Coords[i].Lat)
		lng1 := toRad(p.Coords[i].Lng)
		lat2 := toRad(p.Coords[i+1].Lat)
		lng2 := toRad(p.Coords[i+1].Lng)

		area += (lng2 - lng1) * (2 + math.Sin(lat1) + math.Sin(lat2))
	}

	area = area * EarthRadius * EarthRadius / 2.0

	return math.Abs(area), nil
}
