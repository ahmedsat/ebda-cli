package utils_test

import (
	"testing"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestNewPointFeature(t *testing.T) {
	f := utils.NewPointFeature("test", 30.0, 31.0)
	if f.Type != "Feature" {
		t.Fatalf("Type = %q, want Feature", f.Type)
	}
	if f.Geometry.Type != "Point" {
		t.Fatalf("Geometry.Type = %q, want Point", f.Geometry.Type)
	}
	if f.Properties.Name != "test" {
		t.Fatalf("Properties.Name = %q, want test", f.Properties.Name)
	}
	if got := f.Geometry.Coordinates; len(got) != 2 || got[0] != 31.0 || got[1] != 30.0 {
		t.Fatalf("coordinates = %v, want [31 30]", got)
	}
}

func TestNewGeoJSON(t *testing.T) {
	first := utils.NewPointFeature("first", 0, 0)
	second := utils.NewPointFeature("second", 1, 1)
	g := utils.NewGeoJSON(first, second)
	if g.Type != "FeatureCollection" {
		t.Fatalf("Type = %q, want FeatureCollection", g.Type)
	}
	if len(g.Features) != 2 {
		t.Fatalf("Features len = %d, want 2", len(g.Features))
	}
	if g.Features[0].Properties.Name != "first" || g.Features[1].Properties.Name != "second" {
		t.Fatalf("feature order not preserved: %v", g.Features)
	}
}
