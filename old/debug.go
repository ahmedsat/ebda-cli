//go:build !release

package main

func init() {
	subcommands = append(subcommands, "training", "map", "pgs", "records")
}
