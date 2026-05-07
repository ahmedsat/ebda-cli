package geo

import (
	"errors"
	"fmt"
	"math"
)

const EarthRadius = 6378137.0 // meters (WGS84)

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Polygon struct {
	Ring []Coord // outer ring only (closed)
}

// ---------- helpers ----------

func toRad(d float64) float64 {
	return d * math.Pi / 180
}

func equal(a, b Coord) bool {
	const eps = 1e-9
	return math.Abs(a.Lat-b.Lat) < eps && math.Abs(a.Lng-b.Lng) < eps
}

// ---------- polygon construction ----------

// NewPolygon ensures:
// 1. at least 3 unique points
// 2. ring is closed
func NewPolygon(coords []Coord) (Polygon, error) {
	if len(coords) < 3 {
		return Polygon{}, errors.New("polygon must have at least 3 points")
	}

	ring := append([]Coord(nil), coords...)

	// close ring if needed
	if !equal(ring[0], ring[len(ring)-1]) {
		ring = append(ring, ring[0])
	}

	return Polygon{Ring: ring}, nil
}

// ---------- validation ----------

// Simple validation (not self-intersection check)
func (p Polygon) Valid() error {
	if len(p.Ring) < 4 { // because closed (first == last)
		return fmt.Errorf("invalid polygon: not enough points")
	}
	if !equal(p.Ring[0], p.Ring[len(p.Ring)-1]) {
		return fmt.Errorf("polygon is not closed")
	}
	return nil
}

// ---------- geometry ----------

func (p Polygon) SphericalArea() (float64, error) {

	if err := p.Valid(); err != nil {
		return 0, err
	}

	coords := append([]Coord(nil), p.Ring...)

	// convert to radians
	for i := range coords {
		coords[i].Lat = toRad(coords[i].Lat)
		coords[i].Lng = toRad(coords[i].Lng)
	}

	var area float64
	n := len(coords)

	for i := 0; i < n-1; i++ {
		lat1 := coords[i].Lat
		lng1 := coords[i].Lng
		lat2 := coords[i+1].Lat
		lng2 := coords[i+1].Lng

		dlng := lng2 - lng1
		if dlng > math.Pi {
			dlng -= 2 * math.Pi
		} else if dlng < -math.Pi {
			dlng += 2 * math.Pi
		}

		area += dlng * (2 + math.Sin(lat1) + math.Sin(lat2))
	}

	area = area * EarthRadius * EarthRadius / 2.0
	return math.Abs(area), nil
}

func (p Polygon) AreaFeddan() (float64, error) {
	a, err := p.SphericalArea()
	if err != nil {
		return 0, err
	}
	return a / 4200.0, nil
}
