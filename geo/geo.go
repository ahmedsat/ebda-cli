package geo

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

const EarthRadius = 6378137.0

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Polygon struct {
	Ring []Coord
}

func toRad(d float64) float64 {
	return d * math.Pi / 180
}

func equal(a, b Coord) bool {
	const eps = 1e-9
	return math.Abs(a.Lat-b.Lat) < eps && math.Abs(a.Lng-b.Lng) < eps
}

func Distance(a, b Coord) float64 {
	lat1 := toRad(a.Lat)
	lng1 := toRad(a.Lng)
	lat2 := toRad(b.Lat)
	lng2 := toRad(b.Lng)

	dlat := lat2 - lat1
	dlng := lng2 - lng1

	sinDLat := math.Sin(dlat / 2)
	sinDLng := math.Sin(dlng / 2)

	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLng*sinDLng
	return 2 * EarthRadius * math.Asin(math.Sqrt(h))
}

func NewPolygon(coords []Coord) (Polygon, error) {
	if len(coords) < 3 {
		return Polygon{}, errors.New("polygon must have at least 3 points")
	}

	ring := append([]Coord(nil), coords...)

	if !equal(ring[0], ring[len(ring)-1]) {
		ring = append(ring, ring[0])
	}

	return Polygon{Ring: ring}, nil
}

func (p Polygon) Valid() error {
	if len(p.Ring) < 4 {
		return fmt.Errorf("invalid polygon: not enough points")
	}
	if !equal(p.Ring[0], p.Ring[len(p.Ring)-1]) {
		return fmt.Errorf("polygon is not closed")
	}
	return nil
}

func (p Polygon) SphericalArea() (float64, error) {
	if err := p.Valid(); err != nil {
		return 0, err
	}

	coords := append([]Coord(nil), p.Ring...)

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

func (p Polygon) Overlaps(other Polygon) (bool, error) {
	if err := p.Valid(); err != nil {
		return false, err
	}
	if err := other.Valid(); err != nil {
		return false, err
	}

	for i := 0; i < len(p.Ring)-1; i++ {
		a1 := p.Ring[i]
		a2 := p.Ring[i+1]

		for j := 0; j < len(other.Ring)-1; j++ {
			b1 := other.Ring[j]
			b2 := other.Ring[j+1]

			if segmentsIntersect(a1, a2, b1, b2) {
				return true, nil
			}
		}
	}

	if pointInPolygon(p.Ring[0], other.Ring) {
		return true, nil
	}
	if pointInPolygon(other.Ring[0], p.Ring) {
		return true, nil
	}

	return false, nil
}

func orientation(a, b, c Coord) float64 {
	return (b.Lng-a.Lng)*(c.Lat-a.Lat) - (b.Lat-a.Lat)*(c.Lng-a.Lng)
}

func onSegment(a, b, c Coord) bool {
	return math.Min(a.Lng, b.Lng) <= c.Lng && c.Lng <= math.Max(a.Lng, b.Lng) &&
		math.Min(a.Lat, b.Lat) <= c.Lat && c.Lat <= math.Max(a.Lat, b.Lat)
}

func segmentsIntersect(p1, p2, q1, q2 Coord) bool {
	o1 := orientation(p1, p2, q1)
	o2 := orientation(p1, p2, q2)
	o3 := orientation(q1, q2, p1)
	o4 := orientation(q1, q2, p2)

	if o1*o2 < 0 && o3*o4 < 0 {
		return true
	}

	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}
	if o2 == 0 && onSegment(p1, p2, q2) {
		return true
	}
	if o3 == 0 && onSegment(q1, q2, p1) {
		return true
	}
	if o4 == 0 && onSegment(q1, q2, p2) {
		return true
	}

	return false
}

func pointInPolygon(pt Coord, ring []Coord) bool {
	inside := false
	n := len(ring)

	for i, j := 0, n-1; i < n; j, i = i, i+1 {
		xi, yi := ring[i].Lng, ring[i].Lat
		xj, yj := ring[j].Lng, ring[j].Lat

		intersect := ((yi > pt.Lat) != (yj > pt.Lat)) &&
			(pt.Lng < (xj-xi)*(pt.Lat-yi)/(yj-yi)+xi)

		if intersect {
			inside = !inside
		}
	}

	return inside
}

func (p Polygon) Centroid() (Coord, error) {
	if err := p.Valid(); err != nil {
		return Coord{}, err
	}

	var area float64
	var cx, cy float64

	n := len(p.Ring)

	for i := 0; i < n-1; i++ {
		x0 := p.Ring[i].Lng
		y0 := p.Ring[i].Lat
		x1 := p.Ring[i+1].Lng
		y1 := p.Ring[i+1].Lat

		cross := x0*y1 - x1*y0
		area += cross
		cx += (x0 + x1) * cross
		cy += (y0 + y1) * cross
	}

	area *= 0.5

	if math.Abs(area) < 1e-12 {
		var sumLat, sumLng float64
		for i := 0; i < n-1; i++ {
			sumLat += p.Ring[i].Lat
			sumLng += p.Ring[i].Lng
		}
		count := float64(n - 1)
		return Coord{
			Lat: sumLat / count,
			Lng: sumLng / count,
		}, nil
	}

	cx /= (6 * area)
	cy /= (6 * area)

	return Coord{
		Lat: cy,
		Lng: cx,
	}, nil
}

func polygonClip(subject, clip []Coord) []Coord {
	output := append([]Coord(nil), subject...)

	for i := 0; i < len(clip)-1; i++ {
		input := append([]Coord(nil), output...)
		output = nil

		A := clip[i]
		B := clip[i+1]

		if len(input) == 0 {
			break
		}

		prev := input[len(input)-1]

		for _, curr := range input {
			if inside(curr, A, B) {
				if !inside(prev, A, B) {
					output = append(output, intersection(prev, curr, A, B))
				}
				output = append(output, curr)
			} else if inside(prev, A, B) {
				output = append(output, intersection(prev, curr, A, B))
			}
			prev = curr
		}
	}

	return output
}

func inside(p, a, b Coord) bool {
	return (b.Lng-a.Lng)*(p.Lat-a.Lat)-(b.Lat-a.Lat)*(p.Lng-a.Lng) >= 0
}

func intersection(p1, p2, q1, q2 Coord) Coord {
	A1 := p2.Lat - p1.Lat
	B1 := p1.Lng - p2.Lng
	C1 := A1*p1.Lng + B1*p1.Lat

	A2 := q2.Lat - q1.Lat
	B2 := q1.Lng - q2.Lng
	C2 := A2*q1.Lng + B2*q1.Lat

	det := A1*B2 - A2*B1

	if math.Abs(det) < 1e-12 {
		return p2
	}

	x := (B2*C1 - B1*C2) / det
	y := (A1*C2 - A2*C1) / det

	return Coord{Lat: y, Lng: x}
}

func (p Polygon) Union(other Polygon) ([]Polygon, error) {
	if err := p.Valid(); err != nil {
		return nil, err
	}
	if err := other.Valid(); err != nil {
		return nil, err
	}

	overlap, _ := p.Overlaps(other)
	if !overlap {
		return []Polygon{p, other}, nil
	}

	points := collectUnionPoints(p.Ring, other.Ring)

	if len(points) < 3 {
		return nil, fmt.Errorf("union failed: insufficient points")
	}

	hull := convexHull(points)

	poly, err := NewPolygon(hull)
	if err != nil {
		return nil, err
	}

	return []Polygon{poly}, nil
}

func collectUnionPoints(a, b []Coord) []Coord {
	var pts []Coord

	for _, p := range a[:len(a)-1] {
		if !pointInPolygon(p, b) {
			pts = append(pts, p)
		}
	}

	for _, p := range b[:len(b)-1] {
		if !pointInPolygon(p, a) {
			pts = append(pts, p)
		}
	}

	for i := 0; i < len(a)-1; i++ {
		for j := 0; j < len(b)-1; j++ {
			if segmentsIntersect(a[i], a[i+1], b[j], b[j+1]) {
				pts = append(pts, intersection(a[i], a[i+1], b[j], b[j+1]))
			}
		}
	}

	return pts
}

func convexHull(points []Coord) []Coord {
	if len(points) < 3 {
		return points
	}

	sort.Slice(points, func(i, j int) bool {
		if points[i].Lng == points[j].Lng {
			return points[i].Lat < points[j].Lat
		}
		return points[i].Lng < points[j].Lng
	})

	var lower []Coord
	for _, p := range points {
		for len(lower) >= 2 && orientation(lower[len(lower)-2], lower[len(lower)-1], p) <= 0 {
			lower = lower[:len(lower)-1]
		}
		lower = append(lower, p)
	}

	var upper []Coord
	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		for len(upper) >= 2 && orientation(upper[len(upper)-2], upper[len(upper)-1], p) <= 0 {
			upper = upper[:len(upper)-1]
		}
		upper = append(upper, p)
	}

	hull := append(lower[:len(lower)-1], upper[:len(upper)-1]...)
	return append(hull, hull[0])
}

type triangle [3]Coord

func (p Polygon) OverlapArea(other Polygon) (float64, error) {
	if err := p.Valid(); err != nil {
		return 0, err
	}
	if err := other.Valid(); err != nil {
		return 0, err
	}

	ring1 := ensureCCW(append([]Coord(nil), p.Ring[:len(p.Ring)-1]...))
	ring2 := ensureCCW(append([]Coord(nil), other.Ring[:len(other.Ring)-1]...))

	t1 := triangulate(ring1)
	t2 := triangulate(ring2)

	var total float64

	for _, a := range t1 {
		for _, b := range t2 {
			poly := triangleIntersection(a, b)
			if len(poly) >= 3 {
				total += polygonAreaMeters(poly)
			}
		}
	}

	return total, nil
}

func triangulate(points []Coord) []triangle {
	var tris []triangle

	pts := append([]Coord(nil), points...)

	if len(pts) < 3 {
		return nil
	}

	for len(pts) >= 3 {
		earFound := false

		n := len(pts)

		for i := 0; i < n; i++ {
			prev := pts[(i+n-1)%n]
			curr := pts[i]
			next := pts[(i+1)%n]

			if !isConvexCorner(prev, curr, next) {
				continue
			}

			ear := triangle{prev, curr, next}

			if containsAnyPoint(ear, pts, i) {
				continue
			}

			tris = append(tris, ear)

			pts = append(pts[:i], pts[i+1:]...)
			earFound = true
			break
		}

		if !earFound {
			break
		}
	}

	return tris
}

//	func isConvexCorner(a, b, c Coord) bool {
//		return orientation(a, b, c) < 0
//	}
func isConvexCorner(a, b, c Coord) bool {
	return orientation(a, b, c) > 0
}

func containsAnyPoint(t triangle, pts []Coord, skip int) bool {
	for i, p := range pts {
		if i == skip || i == (skip+1)%len(pts) || i == (skip-1+len(pts))%len(pts) {
			continue
		}
		if pointInTriangle(p, t) {
			return true
		}
	}
	return false
}

func triangleIntersection(a, b triangle) []Coord {
	poly := []Coord{a[0], a[1], a[2], a[0]}
	clip := []Coord{b[0], b[1], b[2], b[0]}

	return polygonClip(poly, clip)
}

func pointInTriangle(p Coord, t triangle) bool {
	a, b, c := t[0], t[1], t[2]

	o1 := orientation(a, b, p)
	o2 := orientation(b, c, p)
	o3 := orientation(c, a, p)

	hasNeg := (o1 < 0) || (o2 < 0) || (o3 < 0)
	hasPos := (o1 > 0) || (o2 > 0) || (o3 > 0)

	return !(hasNeg && hasPos)
}

func polygonAreaMeters(coords []Coord) float64 {
	if len(coords) < 3 {
		return 0
	}

	ref := coords[0]

	var area float64

	for i := range coords {
		j := (i + 1) % len(coords)

		x1, y1 := project(ref, coords[i])
		x2, y2 := project(ref, coords[j])

		area += x1*y2 - x2*y1
	}

	return math.Abs(area) / 2
}

func project(ref, p Coord) (float64, float64) {
	latScale := 111320.0
	lngScale := 111320.0 * math.Cos(toRad(ref.Lat))

	x := (p.Lng - ref.Lng) * lngScale
	y := (p.Lat - ref.Lat) * latScale

	return x, y
}

func ensureCCW(points []Coord) []Coord {
	if polygonSignedArea(points) < 0 {

		for i, j := 0, len(points)-1; i < j; i, j = i+1, j-1 {
			points[i], points[j] = points[j], points[i]
		}
	}
	return points
}

func polygonSignedArea(pts []Coord) float64 {
	var a float64
	for i := range pts {
		j := (i + 1) % len(pts)
		a += pts[i].Lng*pts[j].Lat - pts[j].Lng*pts[i].Lat
	}
	return a / 2
}
