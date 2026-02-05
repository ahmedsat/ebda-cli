package utils

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	Name string `json:"name"`
}

type Feature struct {
	Type       string `json:"type"`
	Properties `json:"properties"`
	Geometry   `json:"geometry"`
}

func NewPointFeature(name string, lat, lon float64) Feature {
	return Feature{
		Type: "Feature",
		Properties: Properties{
			Name: name,
		},
		Geometry: Geometry{
			Type:        "Point",
			Coordinates: []float64{lon, lat},
		},
	}
}

type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

func NewGeoJSON(features ...Feature) GeoJSON {
	return GeoJSON{
		Type:     "FeatureCollection",
		Features: features,
	}
}
