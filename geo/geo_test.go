package geo_test

import (
	"math"
	"testing"

	"github.com/ahmedsat/ebda-cli/geo"
)

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

const floatTol = 1e-6 // relative tolerance for floating-point comparisons

func approxEqual(a, b, tol float64) bool {
	if b == 0 {
		return math.Abs(a) < tol
	}
	return math.Abs(a-b)/math.Abs(b) < tol
}

// A small square (~1 km side) near Cairo for realistic spherical tests.
// Coordinates are closed (first == last).
var squareCairo = []geo.Coord{
	{Lat: 30.0000, Lng: 31.0000},
	{Lat: 30.0090, Lng: 31.0000},
	{Lat: 30.0090, Lng: 31.0113},
	{Lat: 30.0000, Lng: 31.0113},
	{Lat: 30.0000, Lng: 31.0000}, // closed
}

// A triangle inside squareCairo – used for overlap / union tests.
var triangleInside = []geo.Coord{
	{Lat: 30.0010, Lng: 31.0010},
	{Lat: 30.0080, Lng: 31.0010},
	{Lat: 30.0045, Lng: 31.0100},
	{Lat: 30.0010, Lng: 31.0010},
}

// A polygon completely outside squareCairo.
var farAway = []geo.Coord{
	{Lat: 31.0000, Lng: 32.0000},
	{Lat: 31.0090, Lng: 32.0000},
	{Lat: 31.0090, Lng: 32.0113},
	{Lat: 31.0000, Lng: 32.0113},
	{Lat: 31.0000, Lng: 32.0000},
}

// ─────────────────────────────────────────────
// Constants
// ─────────────────────────────────────────────

func TestEarthRadius(t *testing.T) {
	const want = 6_378_137.0
	if geo.EarthRadius != want {
		t.Errorf("EarthRadius = %v, want %v", geo.EarthRadius, want)
	}
}

// ─────────────────────────────────────────────
// Distance
// ─────────────────────────────────────────────

func TestDistance_SamePoint(t *testing.T) {
	c := geo.Coord{Lat: 30.0, Lng: 31.0}
	if d := geo.Distance(c, c); d != 0 {
		t.Errorf("Distance of same point = %v, want 0", d)
	}
}

func TestDistance_KnownValue(t *testing.T) {
	// Cairo (30.0444° N, 31.2357° E) → Alexandria (31.2001° N, 29.9187° E)
	// Haversine ≈ 183 km  (accept ±5 km)
	cairo := geo.Coord{Lat: 30.0444, Lng: 31.2357}
	alex := geo.Coord{Lat: 31.2001, Lng: 29.9187}
	d := geo.Distance(cairo, alex)
	if d < 178_000 || d > 188_000 {
		t.Errorf("Distance Cairo→Alexandria = %.0f m, want ≈183 000 m", d)
	}
}

func TestDistance_Symmetry(t *testing.T) {
	a := geo.Coord{Lat: 30.0, Lng: 31.0}
	b := geo.Coord{Lat: 31.0, Lng: 32.0}
	if geo.Distance(a, b) != geo.Distance(b, a) {
		t.Error("Distance is not symmetric")
	}
}

func TestDistance_NonNegative(t *testing.T) {
	a := geo.Coord{Lat: -33.8688, Lng: 151.2093} // Sydney
	b := geo.Coord{Lat: 51.5074, Lng: -0.1278}   // London
	if d := geo.Distance(a, b); d < 0 {
		t.Errorf("Distance is negative: %v", d)
	}
}

// ─────────────────────────────────────────────
// NewPolygon
// ─────────────────────────────────────────────

func TestNewPolygon_Valid(t *testing.T) {
	_, err := geo.NewPolygon(squareCairo)
	if err != nil {
		t.Errorf("NewPolygon with valid coords returned error: %v", err)
	}
}

func TestNewPolygon_TooFewPoints(t *testing.T) {
	coords := []geo.Coord{
		{Lat: 30.0, Lng: 31.0},
		{Lat: 30.1, Lng: 31.0},
	}
	_, err := geo.NewPolygon(coords)
	if err == nil {
		t.Error("NewPolygon with 2 points should return an error")
	}
}

func TestNewPolygon_EmptySlice(t *testing.T) {
	_, err := geo.NewPolygon(nil)
	if err == nil {
		t.Error("NewPolygon with nil slice should return an error")
	}
}

// ─────────────────────────────────────────────
// Polygon.Valid
// ─────────────────────────────────────────────

func TestPolygon_Valid_OK(t *testing.T) {
	p, err := geo.NewPolygon(squareCairo)
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := p.Valid(); err != nil {
		t.Errorf("Valid() on a good polygon returned error: %v", err)
	}
}

func TestPolygon_Valid_Unclosed(t *testing.T) {
	// Ring that is not closed (first != last)
	coords := []geo.Coord{
		{Lat: 30.0000, Lng: 31.0000},
		{Lat: 30.0090, Lng: 31.0000},
		{Lat: 30.0090, Lng: 31.0113},
		{Lat: 30.0000, Lng: 31.0113},
		// deliberately omit the closing point
	}
	p := geo.Polygon{Ring: coords}
	// Valid() may or may not error on an unclosed ring depending on
	// implementation; we just ensure it doesn't panic.
	_ = p.Valid()
}

// ─────────────────────────────────────────────
// Polygon.SphericalArea
// ─────────────────────────────────────────────

func TestPolygon_SphericalArea_Positive(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	area, err := p.SphericalArea()
	if err != nil {
		t.Fatalf("SphericalArea error: %v", err)
	}
	if area <= 0 {
		t.Errorf("SphericalArea = %v, want > 0", area)
	}
}

func TestPolygon_SphericalArea_ReasonableOrder(t *testing.T) {
	// squareCairo is roughly 1 km × 1 km ≈ 1e6 m²; accept 0.5e6 – 2e6
	p, _ := geo.NewPolygon(squareCairo)
	area, err := p.SphericalArea()
	if err != nil {
		t.Fatalf("SphericalArea error: %v", err)
	}
	if area < 0.5e6 || area > 2e6 {
		t.Errorf("SphericalArea = %.0f m², expected ~1 000 000 m²", area)
	}
}

// ─────────────────────────────────────────────
// Polygon.AreaFeddan
// ─────────────────────────────────────────────

func TestPolygon_AreaFeddan_Positive(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	feddan, err := p.AreaFeddan()
	if err != nil {
		t.Fatalf("AreaFeddan error: %v", err)
	}
	if feddan <= 0 {
		t.Errorf("AreaFeddan = %v, want > 0", feddan)
	}
}

func TestPolygon_AreaFeddan_ConsistentWithSpherical(t *testing.T) {
	// 1 feddan = 4200.833 m²  (Egyptian feddan)
	const feddanInM2 = 4200.833
	p, _ := geo.NewPolygon(squareCairo)
	area, _ := p.SphericalArea()
	feddan, err := p.AreaFeddan()
	if err != nil {
		t.Fatalf("AreaFeddan error: %v", err)
	}
	expected := area / feddanInM2
	if !approxEqual(feddan, expected, 0.01) {
		t.Errorf("AreaFeddan = %v, expected %.4f (from SphericalArea/4200.833)", feddan, expected)
	}
}

// ─────────────────────────────────────────────
// Polygon.Centroid
// ─────────────────────────────────────────────

func TestPolygon_Centroid_InsidePolygon(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	c, err := p.Centroid()
	if err != nil {
		t.Fatalf("Centroid error: %v", err)
	}
	// Centroid of squareCairo should be roughly (30.0045, 31.00565)
	if !approxEqual(c.Lat, 30.0045, 0.01) {
		t.Errorf("Centroid.Lat = %v, want ~30.0045", c.Lat)
	}
	if !approxEqual(c.Lng, 31.00565, 0.01) {
		t.Errorf("Centroid.Lng = %v, want ~31.00565", c.Lng)
	}
}

func TestPolygon_Centroid_ValidCoords(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	c, err := p.Centroid()
	if err != nil {
		t.Fatalf("Centroid error: %v", err)
	}
	if c.Lat < -90 || c.Lat > 90 {
		t.Errorf("Centroid.Lat out of range: %v", c.Lat)
	}
	if c.Lng < -180 || c.Lng > 180 {
		t.Errorf("Centroid.Lng out of range: %v", c.Lng)
	}
}

// ─────────────────────────────────────────────
// Polygon.Overlaps
// ─────────────────────────────────────────────

func TestPolygon_Overlaps_WithItself(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	overlaps, err := p.Overlaps(p)
	if err != nil {
		t.Fatalf("Overlaps error: %v", err)
	}
	if !overlaps {
		t.Error("A polygon should overlap with itself")
	}
}

func TestPolygon_Overlaps_ContainedPolygon(t *testing.T) {
	outer, _ := geo.NewPolygon(squareCairo)
	inner, _ := geo.NewPolygon(triangleInside)
	overlaps, err := outer.Overlaps(inner)
	if err != nil {
		t.Fatalf("Overlaps error: %v", err)
	}
	if !overlaps {
		t.Error("Outer polygon should overlap with inner polygon")
	}
}

func TestPolygon_Overlaps_DisjointPolygons(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(farAway)
	overlaps, err := p.Overlaps(q)
	if err != nil {
		t.Fatalf("Overlaps error: %v", err)
	}
	if overlaps {
		t.Error("Disjoint polygons should not overlap")
	}
}

func TestPolygon_Overlaps_Symmetry(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(triangleInside)
	ab, err1 := p.Overlaps(q)
	ba, err2 := q.Overlaps(p)
	if err1 != nil || err2 != nil {
		t.Fatalf("Overlaps errors: %v, %v", err1, err2)
	}
	if ab != ba {
		t.Errorf("Overlaps not symmetric: p.Overlaps(q)=%v q.Overlaps(p)=%v", ab, ba)
	}
}

// ─────────────────────────────────────────────
// Polygon.OverlapArea
// ─────────────────────────────────────────────

func TestPolygon_OverlapArea_WithItself(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	selfArea, _ := p.SphericalArea()
	overlap, err := p.OverlapArea(p)
	if err != nil {
		t.Fatalf("OverlapArea error: %v", err)
	}
	if !approxEqual(overlap, selfArea, 0.01) {
		t.Errorf("OverlapArea(self) = %.0f, want ≈ SphericalArea = %.0f", overlap, selfArea)
	}
}

func TestPolygon_OverlapArea_Disjoint(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(farAway)
	overlap, err := p.OverlapArea(q)
	if err != nil {
		t.Fatalf("OverlapArea error: %v", err)
	}
	if overlap != 0 {
		t.Errorf("OverlapArea of disjoint polygons = %v, want 0", overlap)
	}
}

func TestPolygon_OverlapArea_ContainedIsSmaller(t *testing.T) {
	outer, _ := geo.NewPolygon(squareCairo)
	inner, _ := geo.NewPolygon(triangleInside)
	outerArea, _ := outer.SphericalArea()
	overlap, err := outer.OverlapArea(inner)
	if err != nil {
		t.Fatalf("OverlapArea error: %v", err)
	}
	if overlap > outerArea {
		t.Errorf("OverlapArea (%v) > outer area (%v)", overlap, outerArea)
	}
}

func TestPolygon_OverlapArea_NonNegative(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(triangleInside)
	overlap, err := p.OverlapArea(q)
	if err != nil {
		t.Fatalf("OverlapArea error: %v", err)
	}
	if overlap < 0 {
		t.Errorf("OverlapArea is negative: %v", overlap)
	}
}

// ─────────────────────────────────────────────
// Polygon.Union
// ─────────────────────────────────────────────

func TestPolygon_Union_ReturnsPolygons(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(triangleInside)
	result, err := p.Union(q)
	if err != nil {
		t.Fatalf("Union error: %v", err)
	}
	if len(result) == 0 {
		t.Error("Union returned empty slice")
	}
}

func TestPolygon_Union_Disjoint_TwoPolygons(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(farAway)
	result, err := p.Union(q)
	if err != nil {
		t.Fatalf("Union error: %v", err)
	}
	// Disjoint union should produce 2 polygons
	if len(result) != 2 {
		t.Errorf("Union of disjoint polygons = %d polygon(s), want 2", len(result))
	}
}

func TestPolygon_Union_AreaAtLeastMax(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(triangleInside)
	pArea, _ := p.SphericalArea()
	qArea, _ := q.SphericalArea()
	maxArea := math.Max(pArea, qArea)

	result, err := p.Union(q)
	if err != nil {
		t.Fatalf("Union error: %v", err)
	}

	var totalArea float64
	for _, poly := range result {
		a, err := poly.SphericalArea()
		if err != nil {
			t.Fatalf("SphericalArea on union result error: %v", err)
		}
		totalArea += a
	}
	if totalArea < maxArea*0.99 {
		t.Errorf("Union total area (%.0f) < max input area (%.0f)", totalArea, maxArea)
	}
}

func TestPolygon_Union_AllResultsValid(t *testing.T) {
	p, _ := geo.NewPolygon(squareCairo)
	q, _ := geo.NewPolygon(triangleInside)
	result, err := p.Union(q)
	if err != nil {
		t.Fatalf("Union error: %v", err)
	}
	for i, poly := range result {
		if err := poly.Valid(); err != nil {
			t.Errorf("Union result[%d].Valid() error: %v", i, err)
		}
	}
}
