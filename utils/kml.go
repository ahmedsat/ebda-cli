package utils

import (
	"encoding/xml"
	"strings"
)

type KML struct {
	XMLName  xml.Name `xml:"kml"`
	Xmlns    string   `xml:"xmlns,attr"`
	Document Document `xml:"Document"`
}

// type Folder struct {
// }

type Document struct {
	Name string `xml:"name,omitempty"`
	// Folders []Folder `xml:"Folder"`
	Placemarks []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name        string  `xml:"name,omitempty"`
	Description string  `xml:"description,omitempty"`
	Style       Style   `xml:"Style,omitempty"`
	Polygon     Polygon `xml:"Polygon,omitempty"`
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

func KmlColor(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "7d00ff00" // default semi-transparent green
	}

	rr := hex[0:2]
	gg := hex[2:4]
	bb := hex[4:6]

	return "7d" + bb + gg + rr
}
