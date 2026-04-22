package geo

// type Coord struct {
// 	Lat float64 `json:"Lat"`
// 	Lng float64 `jsn:"Lng"`
// }

// type OverlapResult struct {
// 	AreaA       float64
// 	AreaB       float64
// 	OverlapArea float64
// 	RatioA      float64
// 	RatioB      float64
// }

// func makeValid(poly *geos.Geometry) *geos.Geometry {
// 	g, err := poly.Buffer(0)
// 	if err != nil {
// 		return poly
// 	}
// 	return g
// }

// func ComputeOverlap(a, b []Coord) (*OverlapResult, error) {

// 	pA, err := toGEOSPolygon(a)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pB, err := toGEOSPolygon(b)
// 	if err != nil {
// 		return nil, err
// 	}

// 	validA := makeValid(pA)
// 	validB := makeValid(pB)

// 	areaA, err := validA.Area()
// 	if err != nil {
// 		return nil, err
// 	}
// 	areaB, err := validB.Area()
// 	if err != nil {
// 		return nil, err
// 	}

// 	inter, err := validA.Intersection(validB)
// 	if err != nil {
// 		return nil, err
// 	}

// 	overlapArea, err := inter.Area()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &OverlapResult{
// 		AreaA:       areaA,
// 		AreaB:       areaB,
// 		OverlapArea: overlapArea,
// 		RatioA:      overlapArea / areaA,
// 		RatioB:      overlapArea / areaB,
// 	}, nil
// }

// func projectUTM(c Coord) (float64, float64) {
// 	// UTM constants
// 	zone := 35
// 	k0 := 0.9996

// 	// Convert degrees to radians
// 	lat := c.Lat * math.Pi / 180
// 	lon := c.Lng * math.Pi / 180

// 	lon0 := float64(zone*6-183) * math.Pi / 180 // central meridian

// 	a := 6378137.0
// 	f := 1 / 298.257223563
// 	e2 := f * (2 - f)
// 	ePrime2 := e2 / (1 - e2)

// 	N := a / math.Sqrt(1-e2*math.Sin(lat)*math.Sin(lat))
// 	T := math.Tan(lat) * math.Tan(lat)
// 	C := ePrime2 * math.Cos(lat) * math.Cos(lat)
// 	A := math.Cos(lat) * (lon - lon0)

// 	M := a * ((1-e2/4-3*e2*e2/64-5*e2*e2*e2/256)*lat -
// 		(3*e2/8+3*e2*e2/32+45*e2*e2*e2/1024)*math.Sin(2*lat) +
// 		(15*e2*e2/256+45*e2*e2*e2/1024)*math.Sin(4*lat) -
// 		(35*e2*e2*e2/3072)*math.Sin(6*lat))

// 	x := k0 * N * (A +
// 		(1-T+C)*A*A*A/6 +
// 		(5-18*T+T*T+72*C-58*ePrime2)*A*A*A*A*A/120)

// 	y := k0 * (M + N*math.Tan(lat)*
// 		(A*A/2+
// 			(5-T+9*C+4*C*C)*A*A*A*A/24+
// 			(61-58*T+T*T+600*C-330*ePrime2)*A*A*A*A*A*A/720))

// 	return x, y
// }

// func toGEOSPolygon(coords []Coord) (*geos.Geometry, error) {
// 	if len(coords) < 3 {
// 		return nil, fmt.Errorf("polygon requires at least 3 coordinates")
// 	}

// 	pts := make([]geos.Coord, 0, len(coords)+1)

// 	for _, c := range coords {
// 		x, y := projectUTM(c)
// 		pts = append(pts, geos.NewCoord(x, y))
// 	}

// 	// close ring
// 	x0, y0 := projectUTM(coords[0])
// 	pts = append(pts, geos.NewCoord(x0, y0))

// 	return geos.NewPolygon(pts)
// }

// var farms = []string{
// 	"Hamdi Abdel Latif Mohamed Ahmed",
// }

// func MapGetKml(farms []string, out string) error {

// 	fmt.Println(len(farms))

// 	maps, err := frappe.Get[types.MapRecord](nil, nil, nil)
// 	if err != nil {
// 		return err
// 	}

// 	MapsMap := make(map[string]types.MapRecord)
// 	for i := range maps {
// 		m := maps[i]
// 		if m.Farm == "" {
// 			continue
// 		}

// 		MapsMap[m.Farm] = m
// 		// if MapsMap[m.Farm].Area_in_feddan < m.Area_in_feddan {
// 		// }
// 	}

// 	maps = []types.MapRecord{}
// 	for _, m := range MapsMap {
// 		if slices.Contains(farms, m.Farm) {
// 			maps = append(maps, m)
// 		}
// 	}

// 	bytes, err := types.RecordsToKML(maps)
// 	if err != nil {
// 		return err
// 	}

// 	f, err := os.Create(out)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	if _, err := f.Write(bytes); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func example() {

// 	// * we need to login to frappe first

// 	err = MapGetKml(farms, "remaining.kml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	os.Exit(0)

// 	m, err := frappe.Get1[types.MapRecord]("ju3h6j906t")
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = m.Parse()
// 	if err != nil {
// 		panic(err)
// 	}

// 	p1 := Polygon{
// 		Coords: m.Coordinates,
// 	}

// 	area, err := p1.GeodesicArea()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(area / 4200)

// 	os.Exit(0)

// 	// err = MapGetKml(farms)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	kmlFile, err := os.Open("master_combined_farms_cleaned.kml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	kml := &utils.KML{}
// 	err = xml.NewDecoder(kmlFile).Decode(kml)
// 	if err != nil {
// 		panic(err)
// 	}

// 	farmsMap := map[string][]Coord{}

// 	for _, placemark := range kml.Document.Placemarks {
// 		coords := []Coord{}
// 		polyStr := placemark.Polygon.OuterBoundary.LinearRing.Coordinates
// 		for line := range strings.SplitSeq(polyStr, "\n") {
// 			line = strings.TrimSpace(line)
// 			if line == "" {
// 				continue
// 			}
// 			fields := strings.Split(line, ",")
// 			if len(fields) != 3 {
// 				continue
// 			}

// 			lat, err := strconv.ParseFloat(fields[0], 64)
// 			if err != nil {
// 				continue
// 			}
// 			lng, err := strconv.ParseFloat(fields[1], 64)
// 			if err != nil {
// 				continue
// 			}
// 			coords = append(coords, Coord{Lat: lat, Lng: lng})
// 		}

// 		_, ok := farmsMap[placemark.Name]
// 		if ok {
// 			fmt.Fprintln(os.Stderr, "Already exists:", placemark.Name)
// 			continue
// 		}

// 		if len(coords) < 3 {
// 			fmt.Fprintln(os.Stderr, "skipping", placemark.Name)
// 			fmt.Println(polyStr)
// 			break
// 		}

// 		farmsMap[placemark.Name] = coords

// 	}

// 	for name1, coords1 := range farmsMap {
// 		for name2, coords2 := range farmsMap {
// 			if name1 == name2 {
// 				continue
// 			}

// 			if len(coords1) < 3 {
// 				fmt.Fprintln(os.Stderr, "skipping", name1)
// 				continue
// 			}

// 			if len(coords2) < 3 {
// 				fmt.Fprintln(os.Stderr, "skipping", name2)
// 				continue
// 			}

// 			res, err := ComputeOverlap(coords1, coords2)
// 			if err != nil {
// 				fmt.Fprintln(os.Stderr, name1, name2, err)
// 				os.Exit(0)
// 				continue
// 			}

// 			if res.OverlapArea > 0 {

// 				parts := strings.Split(name1, "-")
// 				code1 := strings.TrimSpace(parts[len(parts)-1])

// 				parts = strings.Split(name2, "-")
// 				code2 := strings.TrimSpace(parts[len(parts)-1])

// 				fmt.Println(strings.Join([]string{code1, code2, fmt.Sprintf("%f", res.OverlapArea)}, "\t"))
// 			}
// 		}
// 	}

// }
