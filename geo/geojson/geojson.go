package geojson

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Geometry interface {
	Type() string
}

type Point struct {
	Coordinates [2]float64
}

func (Point) Type() string { return "Point" }

type LineString struct {
	Coordinates [][2]float64
}

func (LineString) Type() string { return "LineString" }

type Polygon struct {
	Coordinates [][][2]float64 // rings
}

func (Polygon) Type() string { return "Polygon" }

type MultiPoint struct {
	Coordinates [][2]float64
}

func (MultiPoint) Type() string { return "MultiPoint" }

type MultiLineString struct {
	Coordinates [][][2]float64
}

func (MultiLineString) Type() string { return "MultiLineString" }

type MultiPolygon struct {
	Coordinates [][][][2]float64
}

func (MultiPolygon) Type() string { return "MultiPolygon" }

type geometryRaw struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

type GeometryWrapper struct {
	Geometry Geometry
}

func (g *GeometryWrapper) UnmarshalJSON(data []byte) error {
	var raw geometryRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch raw.Type {
	case "Point":
		var coords [2]float64
		if err := json.Unmarshal(raw.Coordinates, &coords); err != nil {
			return err
		}
		g.Geometry = Point{coords}

	case "LineString":
		var coords [][2]float64
		if err := json.Unmarshal(raw.Coordinates, &coords); err != nil {
			return err
		}
		g.Geometry = LineString{coords}

	case "Polygon":
		var coords [][][2]float64
		if err := json.Unmarshal(raw.Coordinates, &coords); err != nil {
			return err
		}
		g.Geometry = Polygon{coords}

	case "MultiPoint":
		var coords [][2]float64
		if err := json.Unmarshal(raw.Coordinates, &coords); err != nil {
			return err
		}
		g.Geometry = MultiPoint{coords}

	case "MultiLineString":
		var coords [][][2]float64
		if err := json.Unmarshal(raw.Coordinates, &coords); err != nil {
			return err
		}
		g.Geometry = MultiLineString{coords}

	case "MultiPolygon":
		var coords [][][][2]float64
		if err := json.Unmarshal(raw.Coordinates, &coords); err != nil {
			return err
		}
		g.Geometry = MultiPolygon{coords}

	default:
		return fmt.Errorf("unsupported geometry type: %s", raw.Type)
	}

	return nil
}

type Feature struct {
	Type       string          `json:"type"`
	Properties map[string]any  `json:"properties"`
	Geometry   GeometryWrapper `json:"geometry"`
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func DecodeFeatureCollection(data []byte) (*FeatureCollection, error) {
	var fc FeatureCollection
	if err := json.Unmarshal(data, &fc); err != nil {
		return nil, err
	}

	if fc.Type != "FeatureCollection" {
		return nil, errors.New("not a FeatureCollection")
	}

	return &fc, nil
}
