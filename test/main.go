package main

import (
	"fmt"
	"os"

	"github.com/ahmedsat/ebda-cli/config"
)

func init() {
	err := config.Configure()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

/*
EG/1300 EG/1351 EG/3467
*/

type point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func main() {

	// asset, err := kobo.GetAssetByID[kobo.PGSNew](738885331)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }

	// fmt.Println(time.Unix(int64(asset.ValidationStatus.Timestamp), 0))
	// fmt.Println(asset.ValidationStatus.ByWhom)

	// area()
	// OverlapKml()
	// FakeFollowUp("EG/10033", "test eng")
	// Kml()
	// deletable()
	// deleteMap()
	// deleteMapByFarmCode()
	deleteOldMap()
}
