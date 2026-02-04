//go:build !release

package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/ahmedsat/erp-reports-cli/erp"
	"github.com/ahmedsat/erp-reports-cli/kml"
	"github.com/ahmedsat/erp-reports-cli/types"
	"github.com/paulsmith/gogeos/geos"
)

var mapCommands = []string{"center", "delete", "overlap", "overlap-list", "area-list"}

func Map(args []string) (err error) {
	switch args[0] {
	case "center":
		err = MapGetCenter(args[1])
	case "delete":

		for i := 2; i < len(args); i++ {
			res, err := erp.DeleteDoc("/api/resource/Map Records/" + args[i])
			if err != nil {
				return err
			}
			if res {
				fmt.Printf("deleted %s\n", args[i])
			} else {
				fmt.Printf("failed to delete %s\n", args[i])
			}
		}
	case "overlap":
		if len(args) < 3 {
			return fmt.Errorf("not enough arguments")
		}
		err = MapGetOverlap(args[1], args[2])
	case "overlap-list":
		err = MapGetOverlapAll()
	case "area-list":
		err = MapGetAreaAll()
	case "kml":
		err = MapGetKml(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "available commands: %s\n", mapCommands)
		err = fmt.Errorf("unknown map command: %s", args[0])
	}
	return
}

func MapGetKml(args []string) error {
	maps, err := erp.Get[types.MapRecord](nil, nil)
	if err != nil {
		return err
	}

	MapsMap := make(map[string]types.MapRecord)
	for i := range maps {
		m := maps[i]
		if m.Farm == "" {
			continue
		}

		if MapsMap[m.Farm].Area_in_feddan < m.Area_in_feddan {
			MapsMap[m.Farm] = m
		}
	}

	maps = []types.MapRecord{}
	for _, m := range MapsMap {
		maps = append(maps, m)
	}

	bytes, err := kml.RecordsToKML(maps)
	if err != nil {
		return err
	}

	for _, arg := range args {
		f, err := os.Create(arg)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.Write(bytes); err != nil {
			return err
		}
	}

	return nil
}

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

var cordsCache = make(map[string][]Coord)

func GetCoords(id string) ([]Coord, error) {

	if c, ok := cordsCache[id]; ok {
		return c, nil
	} else {
		fmt.Fprintf(os.Stderr, "caching %d,%s\r", len(cordsCache), id)
	}

	res, err := erp.GetDoc("Map Records", id)
	if err != nil {
		return nil, err
	}
	type Resp struct {
		Data struct {
			JSONCode string `json:"jsoncode"`
		} `json:"data"`
	}

	var r Resp
	if err := json.Unmarshal(res, &r); err != nil {
		return nil, err
	}

	inner := strings.TrimSpace(r.Data.JSONCode)
	if !strings.HasPrefix(inner, "[") {
		inner = "[" + inner + "]"
	}

	var coords []Coord
	if err := json.Unmarshal([]byte(inner), &coords); err != nil {
		return nil, err
	}

	cordsCache[id] = coords

	return coords, nil
}

func MapGetCenter(id string) error {
	coords, err := GetCoords(id)
	if err != nil {
		return err
	}

	if len(coords) == 0 {
		return fmt.Errorf("no data")
	}

	var x, y, z float64

	for _, c := range coords {
		lat := c.Lat * math.Pi / 180
		lng := c.Lng * math.Pi / 180

		x += math.Cos(lat) * math.Cos(lng)
		y += math.Cos(lat) * math.Sin(lng)
		z += math.Sin(lat)
	}

	n := float64(len(coords))
	x /= n
	y /= n
	z /= n

	lng := math.Atan2(y, x)
	hyp := math.Sqrt(x*x + y*y)
	lat := math.Atan2(z, hyp)

	fmt.Printf("center: %f %f\n", lat*180/math.Pi, lng*180/math.Pi)
	return nil
}

func toGEOSPolygon(coords []Coord) (*geos.Geometry, error) {
	if len(coords) < 3 {
		return nil, fmt.Errorf("polygon requires at least 3 coordinates")
	}

	pts := make([]geos.Coord, 0, len(coords)+1)

	for _, c := range coords {
		x, y := projectUTM(c)
		pts = append(pts, geos.NewCoord(x, y))
	}

	// close ring
	x0, y0 := projectUTM(coords[0])
	pts = append(pts, geos.NewCoord(x0, y0))

	return geos.NewPolygon(pts)
}

func makeValid(poly *geos.Geometry) *geos.Geometry {
	g, err := poly.Buffer(0)
	if err != nil {
		return poly
	}
	return g
}

type OverlapResult struct {
	AreaA       float64
	AreaB       float64
	OverlapArea float64
	RatioA      float64
	RatioB      float64
}

func (r *OverlapResult) String() string {
	return fmt.Sprintf(
		"Area A: %0.2f f\nArea B: %0.2f f\nOverlap: %0.2f f\nRatio A: %0.2f%%\nRatio B: %0.2f%%",
		AreaToFeddan(r.AreaA), AreaToFeddan(r.AreaB), AreaToFeddan(r.OverlapArea),
		r.RatioA*100, r.RatioB*100,
	)
}

func ComputeOverlap(a, b []Coord) (*OverlapResult, error) {

	pA, err := toGEOSPolygon(a)
	if err != nil {
		return nil, err
	}
	pB, err := toGEOSPolygon(b)
	if err != nil {
		return nil, err
	}

	validA := makeValid(pA)
	validB := makeValid(pB)

	areaA, err := validA.Area()
	if err != nil {
		return nil, err
	}
	areaB, err := validB.Area()
	if err != nil {
		return nil, err
	}

	inter, err := validA.Intersection(validB)
	if err != nil {
		return nil, err
	}

	overlapArea, err := inter.Area()
	if err != nil {
		return nil, err
	}

	return &OverlapResult{
		AreaA:       areaA,
		AreaB:       areaB,
		OverlapArea: overlapArea,
		RatioA:      overlapArea / areaA,
		RatioB:      overlapArea / areaB,
	}, nil
}

func AreaToFeddan(areaM2 float64) float64 {
	return areaM2 / 4200.0
}

func MapGetOverlap(a, b string) error {
	coordsA, err := GetCoords(a)
	if err != nil {
		return err
	}
	coordsB, err := GetCoords(b)
	if err != nil {
		return err
	}

	res, err := ComputeOverlap(coordsA, coordsB)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}

func MapGetOverlapAll() error {
	// read all map records form stdin

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		coordsA, err := GetCoords(fields[0])
		if err != nil {
			return err
		}
		coordsB, err := GetCoords(fields[1])
		if err != nil {
			return err
		}

		res, err := ComputeOverlap(coordsA, coordsB)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\t%s\t-\n", fields[0], fields[1])
			return err
		}
		if res.OverlapArea > 0 {
			fmt.Printf("%s\t%s\t%f\n", fields[0], fields[1], AreaToFeddan(res.OverlapArea))
		}

		// fmt.Fprintf(os.Stderr, "%s\r", line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return nil
}

func MapGetAreaAll() error {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		coords, err := GetCoords(fields[0])
		if err != nil {
			return err
		}

		g, err := toGEOSPolygon(coords)
		if err != nil {
			return err
		}

		valid := makeValid(g)
		area, err := valid.Area()
		if err != nil {
			return err
		}

		fmt.Printf("%s\t%f\n", fields[0], AreaToFeddan(area))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return nil
}

func projectUTM(c Coord) (float64, float64) {
	// UTM constants
	zone := 35
	k0 := 0.9996

	// Convert degrees to radians
	lat := c.Lat * math.Pi / 180
	lon := c.Lng * math.Pi / 180

	lon0 := float64(zone*6-183) * math.Pi / 180 // central meridian

	a := 6378137.0
	f := 1 / 298.257223563
	e2 := f * (2 - f)
	ePrime2 := e2 / (1 - e2)

	N := a / math.Sqrt(1-e2*math.Sin(lat)*math.Sin(lat))
	T := math.Tan(lat) * math.Tan(lat)
	C := ePrime2 * math.Cos(lat) * math.Cos(lat)
	A := math.Cos(lat) * (lon - lon0)

	M := a * ((1-e2/4-3*e2*e2/64-5*e2*e2*e2/256)*lat -
		(3*e2/8+3*e2*e2/32+45*e2*e2*e2/1024)*math.Sin(2*lat) +
		(15*e2*e2/256+45*e2*e2*e2/1024)*math.Sin(4*lat) -
		(35*e2*e2*e2/3072)*math.Sin(6*lat))

	x := k0 * N * (A +
		(1-T+C)*A*A*A/6 +
		(5-18*T+T*T+72*C-58*ePrime2)*A*A*A*A*A/120)

	y := k0 * (M + N*math.Tan(lat)*
		(A*A/2+
			(5-T+9*C+4*C*C)*A*A*A*A/24+
			(61-58*T+T*T+600*C-330*ePrime2)*A*A*A*A*A*A/720))

	return x, y
}
