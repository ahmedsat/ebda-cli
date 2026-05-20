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

 */

func main() {
	// area()
	// OverlapKml()
	// FakeFollowUp("EG/10033", "test eng")
	// Kml()
	// deletable()
	// deleteMap()
	// deleteMapByFarmCode()
}
