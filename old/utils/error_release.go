//go:build release

package utils

import (
	"fmt"
	"os"
)

func HandelErr(err error) {
	if err != nil {
		// Production behavior: crash silently or log normally
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
