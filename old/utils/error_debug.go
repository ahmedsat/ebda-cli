//go:build !release

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func HandelErr(err error) {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		if ok {
			fmt.Fprintf(os.Stderr, "%s:%d\n", relativePath(file), line)
		} else {
			fmt.Fprintln(os.Stderr, "unknown")
		}
		os.Exit(1)
	}
}

func relativePath(path string) string {
	cwd, err := os.Getwd()
	HandelErr(err)
	rel, err := filepath.Rel(cwd, path)
	HandelErr(err)
	return rel
}
