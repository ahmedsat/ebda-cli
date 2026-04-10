package web

import (
	"fmt"

	"github.com/ahmedsat/ebda-cli/kobo"
)

var (
	count      = 0
	globeError error
)

func SyncPGS() {
	items, errs := kobo.StreamAssets[kobo.PGSNew](nil)

	for item := range items {
		fmt.Println(item.AtHouseFarmId, "=>", item.ValidationStatus.Label)
		count++
	}

	if err := <-errs; err != nil {
		globeError = err
	}
}
